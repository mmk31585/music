"""Storage client for ``STORAGE_MODE=shared_volume``.

The Go backend's uploads directory is bind-mounted (read-only) at
``AUDIO_SHARED_PATH``.  Audio files are already locally accessible;
no copying is needed — just path resolution + validation.
"""

from __future__ import annotations

import logging
import os

from app.infra.storage.base import AudioStorageClient

logger = logging.getLogger(__name__)


class SharedVolumeStorageClient(AudioStorageClient):
    """Resolves relative paths within a shared filesystem mount.

    The ``source`` is a path relative to *audio_shared_path*
    (e.g. ``"track-audio/abc123.mp3"``).  The client constructs the
    absolute path, validates it lives within the shared root (anti-path-
    traversal), and confirms the file exists.
    """

    def __init__(self, audio_shared_path: str) -> None:
        self._root = os.path.normpath(audio_shared_path)

    def fetch_to_local_path(self, source: str, dest_dir: str = "") -> str:  # noqa: ARG002
        """Resolve *source* relative to the shared volume root.

        Args:
            source: Path relative to ``audio_shared_path``.
            dest_dir: Ignored — the file is read in-place from the
                      mounted volume.

        Returns:
            The absolute, validated local path to the audio file.

        Raises:
            ValueError: If *source* attempts a path traversal.
            FileNotFoundError: If the resolved path does not exist.
        """
        full_path = os.path.normpath(os.path.join(self._root, source))

        # ── Path traversal prevention ─────────────────────────────
        if not full_path.startswith(self._root):
            raise ValueError(
                f"Path traversal attempt detected: {source!r} "
                f"(resolved to {full_path!r}, "
                f"expected under {self._root!r})"
            )

        # ── File existence check ──────────────────────────────────
        if not os.path.isfile(full_path):
            raise FileNotFoundError(
                f"Audio file not found at {full_path!r} "
                f"(from source={source!r})"
            )

        logger.debug("Resolved shared-volume path: %s", full_path)
        return full_path
