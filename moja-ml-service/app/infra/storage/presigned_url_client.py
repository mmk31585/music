"""Storage client for ``STORAGE_MODE=presigned_url``.

The Go backend passes a presigned HTTP(S) download URL.  The client
streams the download to a temp file to avoid loading the whole file
into memory at once (important for RAM-constrained VPS workers).
"""

from __future__ import annotations

import hashlib
import logging
import os
import uuid

import httpx

from app.infra.storage.base import AudioStorageClient

logger = logging.getLogger(__name__)

# Maximum time in seconds to wait for a single download to complete.
# Audio files can be large and the VPS may have a slow internet link.
DOWNLOAD_TIMEOUT_S = 120


class PresignedUrlStorageClient(AudioStorageClient):
    """Downloads a presigned URL to a local temp file for processing.

    The file is saved to *dest_dir* with a generated filename derived
    from the URL hash to avoid collisions.  The caller is responsible
    for deleting the temp file after processing (the task's ``finally``
    block handles this — see :mod:`app.workers.tasks.lyrics_tasks`).
    """

    def __init__(self) -> None:
        self._client = httpx.Client(timeout=DOWNLOAD_TIMEOUT_S, follow_redirects=True)

    def fetch_to_local_path(self, source: str, dest_dir: str) -> str:
        """Stream-download *source* (a presigned URL) to *dest_dir*.

        Args:
            source: HTTP(S) presigned URL.
            dest_dir: Local directory to write the temp file into.

        Returns:
            Absolute path to the downloaded temp file.

        Raises:
            httpx.HTTPStatusError: If the server returned a non-2xx
                status (e.g. expired presigned URL).
            httpx.TimeoutException: If the download exceeds the
                configured timeout (120 s).
            IOError: If writing to *dest_dir* fails.
        """
        os.makedirs(dest_dir, exist_ok=True)

        # Derive a collision-resistant temp filename from the URL hash.
        url_hash = hashlib.sha256(source.encode()).hexdigest()[:16]
        local_path = os.path.join(dest_dir, f"audio_{url_hash}_{uuid.uuid4().hex[:8]}")

        logger.info("Streaming download from presigned URL → %s", local_path)

        with self._client.stream("GET", source) as response:
            response.raise_for_status()  # raise on 4xx/5xx

            # Log a warning for non-audio Content-Type (informational only).
            content_type = response.headers.get("content-type", "")
            if content_type and not content_type.startswith("audio/"):
                logger.warning(
                    "Unexpected Content-Type for audio download: %s "
                    "(URL hash: %s)",
                    content_type,
                    url_hash,
                )

            with open(local_path, "wb") as f:
                for chunk in response.iter_bytes(chunk_size=8192):
                    f.write(chunk)

        logger.info("Download complete: %s (%d bytes)", local_path, os.path.getsize(local_path))
        return local_path
