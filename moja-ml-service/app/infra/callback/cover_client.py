"""
HMAC-signed callback client for delivering cover optimization results to Go.

Extends the pattern from ``go_client.py`` but for the covers callback
endpoint (``/internal/v1/covers/callback``).

See ``docs/webhook-contract.md`` for the HMAC signing specification.
"""

from __future__ import annotations

import hashlib
import hmac
import json
import logging
import time

import httpx

from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
)

logger = logging.getLogger(__name__)


class GoCallbackCoverClient:
    """Deliver cover-optimization results to the Go backend via HMAC-signed POST.

    Args:
        callback_url: Base URL of the Go callback endpoint
            (e.g. ``http://go-backend:8080/api/v1``).
        hmac_secret: Shared HMAC secret (must match Go's copy).
    """

    TIMESTAMP_TOLERANCE_S: int = 300

    def __init__(self, callback_url: str, hmac_secret: str) -> None:
        self._callback_url = callback_url.rstrip("/") + "/internal/v1/covers/callback"
        self._hmac_secret = hmac_secret
        self._client = httpx.Client(timeout=30.0)

    def deliver_cover_result(
        self,
        album_id: str,
        cover_url: str,
        thumb_url: str,
        med_url: str,
        track_id: str | None = None,
    ) -> bool:
        """POST cover optimization results to Go's webhook endpoint.

        Args:
            album_id: The album identifier.
            cover_url: URL/path to the optimized WebP original.
            thumb_url: URL/path to the 100x100 WebP thumbnail.
            med_url: URL/path to the 300x300 WebP medium.
            track_id: Optional track ID to also update track cover.

        Returns:
            ``True`` if Go responded with 2xx.

        Raises:
            CallbackRejectedError: On 4xx (permanent).
            CallbackTransientError: On 5xx/timeout (retryable).
        """
        timestamp = str(int(time.time()))
        payload = {
            "album_id": album_id,
            "cover_url": cover_url,
            "thumb_url": thumb_url,
            "med_url": med_url,
            "track_id": track_id,
        }
        body_bytes = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).encode()

        signature = self._sign_payload(timestamp, body_bytes)

        headers = {
            "Content-Type": "application/json",
            "X-Moja-Signature": signature,
            "X-Moja-Timestamp": timestamp,
        }

        logger.info(
            "Sending cover callback to %s (album=%s, cover=%s)",
            self._callback_url,
            album_id,
            cover_url,
        )

        try:
            response = self._client.post(
                self._callback_url,
                content=body_bytes,
                headers=headers,
            )
        except httpx.TimeoutException:
            logger.warning("Cover callback to %s timed out", self._callback_url)
            raise CallbackTransientError("Cover callback request timed out")
        except httpx.RequestError as exc:
            logger.warning("Cover callback to %s failed: %s", self._callback_url, exc)
            raise CallbackTransientError(f"Cover callback request failed: {exc}") from exc

        if response.status_code in (200, 204):
            logger.info("Cover callback accepted by Go (%d)", response.status_code)
            return True

        if 400 <= response.status_code < 500:
            logger.error(
                "Cover callback permanently rejected by Go (%d): %s",
                response.status_code,
                response.text[:500],
            )
            raise CallbackRejectedError(
                f"Cover callback rejected (HTTP {response.status_code}): {response.text[:200]}"
            )

        logger.warning(
            "Cover callback got transient error from Go (%d): %s",
            response.status_code,
            response.text[:500],
        )
        raise CallbackTransientError(
            f"Cover callback transient error (HTTP {response.status_code})"
        )

    def _sign_payload(self, timestamp: str, body_bytes: bytes) -> str:
        """HMAC-SHA256 of ``\"{timestamp}.{raw_body}\"``, hex-encoded."""
        message = f"{timestamp}.{body_bytes.decode()}".encode()
        return hmac.new(
            self._hmac_secret.encode(),
            message,
            hashlib.sha256,
        ).hexdigest()
