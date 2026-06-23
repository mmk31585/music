"""
HMAC-signed callback client for delivering video processing results to Go.

Extends the pattern from ``go_client.py`` but for the video callback
endpoint (``/internal/v1/video/callback``).

See ``docs/webhook-contract.md`` for the HMAC signing specification.
"""

from __future__ import annotations

import hashlib
import hmac
import json
import logging
import time
from typing import Any

import httpx

from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
)

logger = logging.getLogger(__name__)


class GoCallbackVideoClient:
    """Deliver video-processing results to the Go backend via HMAC-signed POST.

    Args:
        callback_url: Base URL of the Go callback endpoint
            (e.g. ``http://go-backend:8080/api/v1``).
        hmac_secret: Shared HMAC secret (must match Go's copy).
    """

    TIMESTAMP_TOLERANCE_S: int = 300

    def __init__(self, callback_url: str, hmac_secret: str) -> None:
        self._callback_url = (
            callback_url.rstrip("/") + "/internal/v1/video/callback"
        )
        self._hmac_secret = hmac_secret
        self._client = httpx.Client(timeout=30.0)

    def deliver_video_result(
        self,
        video_id: str,
        final_video_path: str,
        thumbnail_path: str,
        duration_ms: int,
        aspect_ratio: str,
        processing_status: str = "completed",
        error_message: str | None = None,
    ) -> bool:
        """POST video processing results to Go's webhook endpoint.

        Args:
            video_id: The video identifier (matching Go's video UUID).
            final_video_path: Path/URL to the processed video file.
            thumbnail_path: Path/URL to the extracted thumbnail.
            duration_ms: Duration of the processed video in milliseconds.
            aspect_ratio: Detected aspect ratio (e.g. ``"16:9"``).
            processing_status: ``"completed"`` or ``"failed"``.
            error_message: Optional error description if failed.

        Returns:
            ``True`` if Go responded with 2xx.

        Raises:
            CallbackRejectedError: On 4xx (permanent — do not retry).
            CallbackTransientError: On 5xx/timeout (retryable).
        """
        timestamp = str(int(time.time()))
        payload = self._build_payload(
            video_id=video_id,
            final_video_path=final_video_path,
            thumbnail_path=thumbnail_path,
            duration_ms=duration_ms,
            aspect_ratio=aspect_ratio,
            processing_status=processing_status,
            error_message=error_message,
        )
        body_bytes = json.dumps(
            payload, ensure_ascii=False, separators=(",", ":")
        ).encode()

        signature = self._sign_payload(timestamp, body_bytes)

        headers = {
            "Content-Type": "application/json",
            "X-Moja-Signature": signature,
            "X-Moja-Timestamp": timestamp,
        }

        logger.info(
            "Sending video callback to %s (video_id=%s, status=%s)",
            self._callback_url,
            video_id,
            processing_status,
        )

        try:
            response = self._client.post(
                self._callback_url,
                content=body_bytes,
                headers=headers,
            )
        except httpx.TimeoutException:
            logger.warning(
                "Video callback to %s timed out", self._callback_url
            )
            raise CallbackTransientError(
                "Video callback request timed out"
            )
        except httpx.RequestError as exc:
            logger.warning(
                "Video callback to %s failed: %s", self._callback_url, exc
            )
            raise CallbackTransientError(
                f"Video callback request failed: {exc}"
            ) from exc

        if response.status_code in (200, 204):
            logger.info(
                "Video callback accepted by Go (%d)", response.status_code
            )
            return True

        if 400 <= response.status_code < 500:
            logger.error(
                "Video callback permanently rejected by Go (%d): %s",
                response.status_code,
                response.text[:500],
            )
            raise CallbackRejectedError(
                f"Video callback rejected (HTTP {response.status_code}): "
                f"{response.text[:200]}"
            )

        logger.warning(
            "Video callback got transient error from Go (%d): %s",
            response.status_code,
            response.text[:500],
        )
        raise CallbackTransientError(
            f"Video callback transient error (HTTP {response.status_code})"
        )

    # ── Internal helpers ──────────────────────────────────────────────

    @staticmethod
    def _build_payload(
        video_id: str,
        final_video_path: str,
        thumbnail_path: str,
        duration_ms: int,
        aspect_ratio: str,
        processing_status: str = "completed",
        error_message: str | None = None,
    ) -> dict[str, Any]:
        """Construct the JSON-serialisable payload dict."""
        return {
            "video_id": video_id,
            "final_video_path": final_video_path,
            "thumbnail_path": thumbnail_path,
            "duration_ms": duration_ms,
            "aspect_ratio": aspect_ratio,
            "processing_status": processing_status,
            "error_message": error_message,
        }

    def _sign_payload(self, timestamp: str, body_bytes: bytes) -> str:
        """HMAC-SHA256 of ``\"{timestamp}.{raw_body}\"``, hex-encoded."""
        message = f"{timestamp}.{body_bytes.decode()}".encode()
        return hmac.new(
            self._hmac_secret.encode(),
            message,
            hashlib.sha256,
        ).hexdigest()
