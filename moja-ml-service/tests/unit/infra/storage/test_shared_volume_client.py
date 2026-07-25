"""Tests for ``SharedVolumeStorageClient``.

Key behaviours to verify:
1. Path resolution and file existence.
2. Path traversal rejection.
3. Non-existent file raises ``FileNotFoundError``.
"""

from __future__ import annotations

import os

import pytest

from app.infra.storage.shared_volume_client import SharedVolumeStorageClient


class TestSharedVolumeStorageClient:
    """Tests for SharedVolumeStorageClient."""

    def test_resolves_valid_path(self, tmp_path) -> None:
        """Given a real file under the root, fetch_to_local_path returns it."""
        root = tmp_path / "shared"
        root.mkdir()
        audio_file = root / "track-audio" / "song.mp3"
        audio_file.parent.mkdir(parents=True)
        audio_file.write_text("fake audio content")

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        result = client.fetch_to_local_path("track-audio/song.mp3")

        assert result == os.path.normpath(str(audio_file))

    def test_accepts_absolute_source_within_root(self, tmp_path) -> None:
        """A source that resolves inside the root is accepted."""
        root = tmp_path / "shared"
        root.mkdir()
        audio_file = root / "song.mp3"
        audio_file.write_text("content")

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        result = client.fetch_to_local_path("song.mp3")
        assert result == os.path.normpath(str(audio_file))

    def test_rejects_path_traversal(self, tmp_path) -> None:
        """A source with '..' that escapes the root raises ValueError."""
        root = tmp_path / "shared"
        root.mkdir()

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        with pytest.raises(ValueError, match="Path traversal"):
            client.fetch_to_local_path("../../../etc/passwd")

    def test_rejects_traversal_via_symlink_pattern(self, tmp_path) -> None:
        """Even complex traversal attempts are caught."""
        root = tmp_path / "shared"
        root.mkdir()

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        with pytest.raises(ValueError, match="Path traversal"):
            client.fetch_to_local_path("foo/../../../../../etc/shadow")

    def test_raises_file_not_found(self, tmp_path) -> None:
        """A non-existent file raises FileNotFoundError."""
        root = tmp_path / "shared"
        root.mkdir()

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        with pytest.raises(FileNotFoundError, match="not found"):
            client.fetch_to_local_path("nonexistent.mp3")

    def test_directory_is_not_a_file(self, tmp_path) -> None:
        """A directory path raises FileNotFoundError (os.path.isfile check)."""
        root = tmp_path / "shared"
        root.mkdir()
        (root / "subdir").mkdir()

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        with pytest.raises(FileNotFoundError):
            client.fetch_to_local_path("subdir")

    def test_dest_dir_parameter_ignored(self, tmp_path) -> None:
        """dest_dir is accepted but ignored in shared_volume mode."""
        root = tmp_path / "shared"
        root.mkdir()
        audio_file = root / "song.mp3"
        audio_file.write_text("content")

        client = SharedVolumeStorageClient(audio_shared_path=str(root))

        result = client.fetch_to_local_path("song.mp3", dest_dir="/some/random/path")
        assert result == os.path.normpath(str(audio_file))
