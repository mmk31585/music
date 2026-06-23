"""Tests for ``PresignedUrlStorageClient``.

Uses ``respx`` to mock HTTP responses so no real network calls are made.
"""

from __future__ import annotations

import os
from pathlib import Path

import pytest
import respx

from app.infra.storage.presigned_url_client import PresignedUrlStorageClient


class TestPresignedUrlStorageClient:
    """Tests for PresignedUrlStorageClient."""

    def test_downloads_file_to_dest_dir(self, tmp_path: Path) -> None:
        """A successful download writes the file to dest_dir."""
        client = PresignedUrlStorageClient()

        download_url = "https://storage.example.com/audio/test.mp3?X-Amz-Signature=abc"

        with respx.mock:
            route = respx.get(download_url).respond(
                200,
                content_type="audio/mpeg",
                content=b"fake audio bytes",
            )

            result = client.fetch_to_local_path(download_url, dest_dir=str(tmp_path))

        assert route.called
        assert os.path.isfile(result)
        assert os.path.getsize(result) == len(b"fake audio bytes")
        with open(result, "rb") as f:
            assert f.read() == b"fake audio bytes"

    def test_download_streams_large_content(self, tmp_path: Path) -> None:
        """Download handles larger files without loading into memory."""
        client = PresignedUrlStorageClient()

        download_url = "https://storage.example.com/big-audio.wav"
        large_content = b"x" * 100_000

        with respx.mock:
            route = respx.get(download_url).respond(
                200,
                content_type="audio/wav",
                content=large_content,
            )

            result = client.fetch_to_local_path(download_url, dest_dir=str(tmp_path))

        assert route.called
        assert os.path.getsize(result) == 100_000

    def test_raises_on_http_error(self, tmp_path: Path) -> None:
        """A non-2xx response raises httpx.HTTPStatusError."""
        client = PresignedUrlStorageClient()

        download_url = "https://storage.example.com/expired.mp3"

        with respx.mock:
            respx.get(download_url).respond(403)

            with pytest.raises(Exception) as exc_info:
                client.fetch_to_local_path(download_url, dest_dir=str(tmp_path))

        assert "403" in str(exc_info.value) or "Forbidden" in str(exc_info.value)

    def test_warns_on_non_audio_content_type(
        self, tmp_path: Path, caplog: pytest.LogCaptureFixture
    ) -> None:
        """A non-audio content type logs a warning but doesn't fail."""
        import logging

        caplog.set_level(logging.WARNING)

        client = PresignedUrlStorageClient()

        download_url = "https://storage.example.com/audio.mp3"
        content = b"some data"

        with respx.mock:
            respx.get(download_url).respond(
                200,
                content_type="application/octet-stream",
                content=content,
            )

            result = client.fetch_to_local_path(download_url, dest_dir=str(tmp_path))

        assert os.path.isfile(result)
        assert any("unexpected" in msg.lower() for msg in caplog.messages)
        assert "application/octet-stream" in " ".join(caplog.messages)

    def test_empty_audio_content_type_no_warning(
        self, tmp_path: Path, caplog: pytest.LogCaptureFixture
    ) -> None:
        """Missing content-type header does not trigger the audio warning."""
        import logging

        caplog.set_level(logging.WARNING)

        client = PresignedUrlStorageClient()

        download_url = "https://storage.example.com/audio.mp3"

        with respx.mock:
            respx.get(download_url).respond(
                200,
                content=b"data",
            )

            client.fetch_to_local_path(download_url, dest_dir=str(tmp_path))

        audio_warnings = [
            m for m in caplog.messages if "content" in m.lower()
        ]
        assert len(audio_warnings) == 0

    def test_creates_dest_dir_if_not_exists(self, tmp_path: Path) -> None:
        """dest_dir is created automatically if it doesn't exist."""
        nonexistent_dir = str(tmp_path / "new" / "nested" / "dir")

        client = PresignedUrlStorageClient()

        download_url = "https://example.com/audio.mp3"
        with respx.mock:
            respx.get(download_url).respond(200, content=b"data")

            result = client.fetch_to_local_path(download_url, dest_dir=nonexistent_dir)

        assert os.path.isfile(result)
