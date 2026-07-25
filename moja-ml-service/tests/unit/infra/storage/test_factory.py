"""Tests for ``get_storage_client`` factory."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from app.config import Settings
from app.infra.storage.factory import get_storage_client
from app.infra.storage.presigned_url_client import PresignedUrlStorageClient
from app.infra.storage.shared_volume_client import SharedVolumeStorageClient


class TestGetStorageClient:
    """Tests for the storage client factory."""

    @staticmethod
    def _make_settings(storage_mode: str) -> Settings:
        """Build a minimal Settings with the given storage_mode."""
        mock = MagicMock(spec=Settings)
        mock.storage_mode = storage_mode
        mock.audio_shared_path = "/mnt/audio-uploads"
        return mock

    def test_shared_volume_mode(self) -> None:
        """shared_volume mode returns SharedVolumeStorageClient."""
        client = get_storage_client(self._make_settings("shared_volume"))
        assert isinstance(client, SharedVolumeStorageClient)

    def test_presigned_url_mode(self) -> None:
        """presigned_url mode returns PresignedUrlStorageClient."""
        client = get_storage_client(self._make_settings("presigned_url"))
        assert isinstance(client, PresignedUrlStorageClient)

    def test_unknown_mode_raises(self) -> None:
        """An unsupported storage_mode raises ValueError."""
        settings = self._make_settings("s3_direct")
        with pytest.raises(ValueError, match="Unknown storage_mode"):
            get_storage_client(settings)