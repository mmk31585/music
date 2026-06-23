"""
HMAC-signed HTTP callback client for delivering results to the Go backend.

Security model
--------------
All callbacks from this service to the Go backend are signed with
**HMAC-SHA256** using a shared secret known only to these two services.

The signing scheme is::

    signature = HMAC-SHA256(secret, "{timestamp}.{raw_body}")

The Go backend **must** verify the signature against the **raw request
body bytes** (not a re-parsed/re-serialised version of the JSON, which
could differ in key ordering or whitespace and break the signature
check).

See ``docs/webhook-contract.md`` for the full specification,
including the Go-side verification procedure.
"""

from __future__ import annotations

import hashlib
import hmac
import json
import logging
import time
from typing import Any

import httpx

logger = logging.getLogger(__name__)


class CallbackDeliveryError(Exception):
    """Base exception for callback delivery failures.

    Transient subclasses trigger Celery ``self.retry()``;
    permanent subclasses fail the job immediately.
    """


class CallbackTransientError(CallbackDeliveryError):
    """Transient failure — network issue, 5xx from Go, etc.

    The Celery task catches this and calls ``self.retry()``.
    """


class CallbackRejectedError(CallbackDeliveryError):
    """Permanent rejection from Go (4xx response).

    The Celery task does **NOT** retry — the job is marked ``failed``.
    """


class GoCallbackClient:
    """Deliver lyrics-processing results to the Go backend via HMAC-signed POST.

    Args:
        callback_url: Base URL of the Go callback endpoint
            (e.g. ``http://go-backend:8080/api/v1``).  The client
            appends ``/internal/v1/lyrics/callback``.
        hmac_secret: Shared HMAC secret (must match Go's copy).
    """

    # Maximum age of a callback timestamp, in seconds.  Go should also
    # enforce this for replay protection.
    TIMESTAMP_TOLERANCE_S: int = 300

    def __init__(self, callback_url: str, hmac_secret: str) -> None:
        self._callback_url = callback_url.rstrip("/") + "/internal/v1/lyrics/callback"
        self._hmac_secret = hmac_secret
        self._client = httpx.Client(timeout=30.0)

    # ── Public API ────────────────────────────────────────────────────

    def deliver_lyrics_result(
        self,
        track_id: str,
        lrc_content: str,
        plain_text: str,
        confidence: float,
        detected_language: str,
        whisper_model_version: str,
    ) -> bool:
        """POST lyrics results to Go's webhook endpoint.

        Args:
            track_id: The track identifier (matching Go's catalog).
            lrc_content: LRC-format lyrics with timestamps.
            plain_text: Plain-text transcription.
            confidence: Overall confidence score 0.0–1.0.
            detected_language: ISO 639-1 language code.
            whisper_model_version: Model size string (e.g. ``"medium-int8"``).

        Returns:
            ``True`` if Go responded with 2xx, ``False`` otherwise.

        Raises:
            CallbackRejectedError: On 4xx responses (permanent
                rejection — do not retry).
            CallbackTransientError: On 5xx, timeouts, or connection
                failures.  The Celery task catches this and calls
                ``self.retry()``.
        """
        timestamp = str(int(time.time()))
        payload = self._build_payload(
            track_id=track_id,
            lrc_content=lrc_content,
            plain_text=plain_text,
            confidence=confidence,
            detected_language=detected_language,
            whisper_model_version=whisper_model_version,
        )
        body_bytes = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).encode()

        signature = self._sign_payload(timestamp, body_bytes)

        headers = {
            "Content-Type": "application/json",
            "X-Moja-Signature": signature,
            "X-Moja-Timestamp": timestamp,
        }

        logger.info(
            "Sending callback to %s (track=%s, confidence=%.3f)",
            self._callback_url,
            track_id,
            confidence,
        )

        try:
            response = self._client.post(
                self._callback_url,
                content=body_bytes,
                headers=headers,
            )
        except httpx.TimeoutException:
            logger.warning("Callback to %s timed out", self._callback_url)
            raise CallbackTransientError("Callback request timed out")
        except httpx.RequestError as exc:
            logger.warning("Callback to %s failed: %s", self._callback_url, exc)
            raise CallbackTransientError(f"Callback request failed: {exc}") from exc

        if response.status_code in (200, 204):
            logger.info("Callback accepted by Go (%d)", response.status_code)
            return True

        if 400 <= response.status_code < 500:
            # 4xx = permanent rejection; do NOT retry.
            logger.error(
                "Callback permanently rejected by Go (%d): %s",
                response.status_code,
                response.text[:500],
            )
            raise CallbackRejectedError(
                f"Callback rejected (HTTP {response.status_code}): {response.text[:200]}"
            )

        # 5xx or other = transient; caller should retry.
        logger.warning(
            "Callback got transient error from Go (%d): %s",
            response.status_code,
            response.text[:500],
        )
        raise CallbackTransientError(
            f"Callback transient error (HTTP {response.status_code})"
        )

    # ── Internal helpers ──────────────────────────────────────────────

    @staticmethod
    def _build_payload(
        track_id: str,
        lrc_content: str,
        plain_text: str,
        confidence: float,
        detected_language: str,
        whisper_model_version: str,
    ) -> dict[str, Any]:
        """Construct the JSON-serialisable payload dict.

        Uses ``ensure_ascii=False`` and ``separators=(",", ":")`` during
        serialisation so that Persian text is preserved and the output is
        compact (no extra whitespace that could shift the HMAC signature).
        """
        return {
            "track_id": track_id,
            "lrc_content": lrc_content,
            "plain_text": plain_text,
            "confidence": confidence,
            "detected_language": detected_language,
            "whisper_model_version": whisper_model_version,
        }

    def _sign_payload(self, timestamp: str, body_bytes: bytes) -> str:
        """HMAC-SHA256 of ``"{timestamp}.{raw_body}"``, hex-encoded.

        .. caution::

            Go **must** verify this signature against the **raw body
            bytes**, not a re-serialised version of the parsed JSON.
            Key ordering and whitespace in the JSON serialisation must
            match exactly.
        """
        message = f"{timestamp}.{body_bytes.decode()}".encode()
        return hmac.new(
            self._hmac_secret.encode(),
            message,
            hashlib.sha256,
        ).hexdigest()
