"""Tests for the lyrics job API endpoints.

These tests use mocked DB sessions and Celery — no real database or
worker is needed. Overrides are wired via ``app.dependency_overrides``,
the standard FastAPI pattern for test injection.

We use the **sync** ``TestClient`` (``from fastapi.testclient``) so that
all the async machinery is handled internally by Starlette — no need for
``@pytest.mark.asyncio`` or async fixtures.
"""

from __future__ import annotations

import uuid
from datetime import datetime
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from fastapi.testclient import TestClient
from sqlalchemy.ext.asyncio import AsyncSession

from app.infra.db.models import LyricsJob
from app.infra.db.session import get_db
from app.main import app

# ── Fixtures ────────────────────────────────────────────────────────────


@pytest.fixture
def mock_session() -> MagicMock:
    """Return a mock async DB session."""
    session = MagicMock(spec=AsyncSession)
    # execute is AsyncMock (needed for ``await``), but its return_value
    # must be a plain MagicMock — otherwise ``.scalar_one_or_none()``
    # returns a coroutine instead of a value.
    session.execute = AsyncMock()
    session.execute.return_value = MagicMock()
    session.add = MagicMock()
    session.commit = AsyncMock()
    session.refresh = AsyncMock()
    return session


@pytest.fixture
def client(mock_session: MagicMock) -> TestClient:  # type: ignore[misc]
    """Create a sync TestClient with the mock session wired in.

    We override ``get_db`` so that it returns our mock session, bypassing
    the real database entirely.
    """

    # We need an async generator so FastAPI's dependency injection can
    # ``__aiter__`` and ``__anext__`` it. A simple ``async def`` that
    # yields works because TestClient drives the event loop internally.
    async def override_get_db() -> AsyncSession:
        return mock_session  # type: ignore[return-value]

    app.dependency_overrides[get_db] = override_get_db
    transport = TestClient(app)
    yield transport
    app.dependency_overrides.pop(get_db, None)


# ── POST /api/v1/lyrics/jobs ────────────────────────────────────────────


class TestCreateLyricsJob:
    """Tests for the job creation endpoint."""

    def test_created_job_has_job_id(
        self,
        client: TestClient,
        mock_session: MagicMock,
    ) -> None:
        """A newly created job should return a ``job_id`` field."""
        # Simulate no existing job
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        with patch("app.api.v1.lyrics.celery_app.send_task") as mock_send:
            resp = client.post(
                "/api/v1/lyrics/jobs",
                json={
                    "track_id": "test-123",
                    "audio_file_path": "/some/path.mp3",
                },
            )

        assert resp.status_code == 202
        data = resp.json()
        assert "job_id" in data
        assert isinstance(data["job_id"], str)
        assert data["track_id"] == "test-123"
        assert data["status"] == "queued"
        mock_send.assert_called_once()
        assert (
            mock_send.call_args[0][0]
            == "app.workers.tasks.lyrics_tasks.generate_lyrics_task"
        )

    def test_idempotent_returns_existing(
        self,
        client: TestClient,
        mock_session: MagicMock,
    ) -> None:
        """Creating the same track_id twice should return the existing job."""
        job_id = uuid.uuid4()
        now = __import__("datetime").datetime.now()
        existing = LyricsJob(
            id=job_id,
            track_id="test-123",
            idempotency_key="lyrics:test-123",
            status="queued",
            audio_file_path="/some/path.mp3",
            retry_count=0,
            callback_delivered=False,
            created_at=now,
            updated_at=now,
        )
        mock_session.execute.return_value.scalar_one_or_none.return_value = existing

        with patch("app.api.v1.lyrics.celery_app.send_task") as mock_send:
            resp = client.post(
                "/api/v1/lyrics/jobs",
                json={
                    "track_id": "test-123",
                    "audio_file_path": "/some/path.mp3",
                },
            )

        assert resp.status_code == 202
        data = resp.json()
        assert data["job_id"] == str(job_id)
        # No new Celery task should be enqueued
        mock_send.assert_not_called()

    def test_allow_new_after_failure(
        self,
        client: TestClient,
        mock_session: MagicMock,
    ) -> None:
        """A failed job for a track_id should allow creating a new one."""
        # Existing job is "failed" — query returns None (filtered out)
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        with patch("app.api.v1.lyrics.celery_app.send_task") as mock_send:
            resp = client.post(
                "/api/v1/lyrics/jobs",
                json={
                    "track_id": "test-123",
                    "audio_file_path": "/some/path.mp3",
                },
            )

        assert resp.status_code == 202
        # A new job should be created
        mock_send.assert_called_once()

    def test_missing_audio_source_returns_422(
        self,
        client: TestClient,
    ) -> None:
        """Missing both audio_file_path and audio_download_url => 422."""
        resp = client.post(
            "/api/v1/lyrics/jobs",
            json={"track_id": "test-123"},
        )
        assert resp.status_code == 422


# ── GET /api/v1/lyrics/jobs/{job_id} ────────────────────────────────────


class TestGetLyricsJob:
    """Tests for the job detail endpoint."""

    def test_get_existing_job(
        self,
        client: TestClient,
        mock_session: MagicMock,
    ) -> None:
        """An existing job should return full detail."""
        job_id = uuid.uuid4()
        now = datetime.now()
        existing = LyricsJob(
            id=job_id,
            track_id="test-123",
            idempotency_key="lyrics:test-123",
            status="queued",
            audio_file_path="/some/path.mp3",
            retry_count=0,
            callback_delivered=False,
            created_at=now,
            updated_at=now,
        )
        mock_session.execute.return_value.scalar_one_or_none.return_value = existing

        resp = client.get(f"/api/v1/lyrics/jobs/{job_id}")
        assert resp.status_code == 200
        data = resp.json()
        assert data["job_id"] == str(job_id)
        assert data["track_id"] == "test-123"
        assert data["status"] == "queued"

    def test_get_nonexistent_job_returns_404(
        self,
        client: TestClient,
        mock_session: MagicMock,
    ) -> None:
        """A job that doesn't exist should return 404."""
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        resp = client.get(f"/api/v1/lyrics/jobs/{uuid.uuid4()}")
        assert resp.status_code == 404

    def test_get_with_invalid_uuid_returns_404(
        self,
        client: TestClient,
    ) -> None:
        """An invalid UUID format should return 404 before hitting DB."""
        resp = client.get("/api/v1/lyrics/jobs/not-a-uuid")
        assert resp.status_code == 404
