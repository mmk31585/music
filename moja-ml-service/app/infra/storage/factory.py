"""Factory: select the right storage client based on ``STORAGE_MODE``."""

from __future__ import annotations

from app.config import Settings
from app.infra.storage.base import AudioStorageClient
from app.infra.storage.presigned_url_client import PresignedUrlStorageClient
from app.infra.storage.shared_volume_client import SharedVolumeStorageClient


def get_storage_client(settings: Settings) -> AudioStorageClient:
    """Return the appropriate storage client for the configured mode.

    Args:
        settings: Application settings (must have ``storage_mode`` set).

    Returns:
        A concrete :class:`AudioStorageClient` implementation.

    Raises:
        ValueError: If ``storage_mode`` is not recognised.
    """
    mode = settings.storage_mode
    if mode == "shared_volume":
        return SharedVolumeStorageClient(audio_shared_path=settings.audio_shared_path)
    elif mode == "presigned_url":
        return PresignedUrlStorageClient()
    raise ValueError(f"Unknown storage_mode: {mode!r}")
