"""Tests for similarity and embedding job API endpoints."""

from __future__ import annotations

import uuid
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from fastapi.testclient import TestClient
from sqlalchemy.ext.asyncio import AsyncSession

from app.domain.audio_features.similarity import (
    SimilarityQueryResult,
    SimilarTrack,
)
from app.infra.db.models import EmbeddingJob
from app.infra.db.session import get_db
from app.main import app


@pytest.fixture
def mock_session() -> MagicMock:
    session = MagicMock(spec=AsyncSession)
    session.execute = AsyncMock()
    session.execute.return_value = MagicMock()
    session.execute.return_value.scalar_one_or_none = MagicMock(return_value=None)
    session.execute.return_value.scalars = MagicMock()
    session.execute.return_value.all = MagicMock(return_value=[])
    session.add = MagicMock()
    session.commit = AsyncMock()
    session.refresh = AsyncMock()
    session.get = AsyncMock()
    return session


@pytest.fixture
def client(mock_session: MagicMock) -> TestClient:
    async def override_get_db() -> AsyncSession:
        return mock_session  # type: ignore[return-value]

    app.dependency_overrides[get_db] = override_get_db
    transport = TestClient(app)
    yield transport
    app.dependency_overrides.pop(get_db, None)


class TestCreateEmbeddingJob:
    """Tests for POST /api/v1/embeddings/jobs."""

    def test_created_job_has_job_id(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        with patch("app.api.v1.similarity.celery_app.send_task") as mock_send:
            resp = client.post(
                "/api/v1/embeddings/jobs",
                json={"track_id": "test-123", "audio_file_path": "/some/path.mp3"},
            )

        assert resp.status_code == 202
        data = resp.json()
        assert "job_id" in data
        assert data["track_id"] == "test-123"
        assert data["status"] == "queued"
        mock_send.assert_called_once()
        assert (
            mock_send.call_args[0][0]
            == "app.workers.tasks.embedding_tasks.extract_embedding_task"
        )

    def test_idempotent_returns_existing(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        job_id = uuid.uuid4()
        existing = EmbeddingJob(
            id=job_id,
            track_id="test-123",
            idempotency_key="embedding:test-123",
            status="queued",
            audio_file_path="/some/path.mp3",
            retry_count=0,
            genre_ids=[],
        )
        mock_session.execute.return_value.scalar_one_or_none.return_value = existing

        with patch("app.api.v1.similarity.celery_app.send_task") as mock_send:
            resp = client.post(
                "/api/v1/embeddings/jobs",
                json={"track_id": "test-123", "audio_file_path": "/some/path.mp3"},
            )

        assert resp.status_code == 202
        data = resp.json()
        assert data["job_id"] == str(job_id)
        mock_send.assert_not_called()

    def test_allow_new_after_failure(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        with patch("app.api.v1.similarity.celery_app.send_task") as mock_send:
            resp = client.post(
                "/api/v1/embeddings/jobs",
                json={"track_id": "test-123", "audio_file_path": "/some/path.mp3"},
            )

        assert resp.status_code == 202
        mock_send.assert_called_once()

    def test_missing_audio_source_returns_422(self, client: TestClient) -> None:
        resp = client.post(
            "/api/v1/embeddings/jobs", json={"track_id": "test-123"}
        )
        assert resp.status_code == 422

    def test_accepts_genre_ids_and_release_year(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        with patch("app.api.v1.similarity.celery_app.send_task"):
            resp = client.post(
                "/api/v1/embeddings/jobs",
                json={
                    "track_id": "test-456",
                    "audio_download_url": "https://example.com/audio.mp3",
                    "genre_ids": ["rock", "pop"],
                    "release_year": 2023,
                },
            )

        assert resp.status_code == 202
        data = resp.json()
        assert data["track_id"] == "test-456"


class TestGetEmbeddingJob:
    """Tests for GET /api/v1/embeddings/jobs/{job_id}."""

    def test_get_existing_job(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        job_id = uuid.uuid4()
        existing = EmbeddingJob(
            id=job_id,
            track_id="test-123",
            idempotency_key="embedding:test-123",
            status="queued",
            audio_file_path="/some/path.mp3",
            retry_count=0,
            genre_ids=[],
        )
        mock_session.execute.return_value.scalar_one_or_none.return_value = existing

        resp = client.get(f"/api/v1/embeddings/jobs/{job_id}")
        assert resp.status_code == 200
        data = resp.json()
        assert data["job_id"] == str(job_id)
        assert data["track_id"] == "test-123"
        assert data["status"] == "queued"

    def test_get_nonexistent_job_returns_404(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        mock_session.execute.return_value.scalar_one_or_none.return_value = None
        resp = client.get(f"/api/v1/embeddings/jobs/{uuid.uuid4()}")
        assert resp.status_code == 404

    def test_get_with_invalid_uuid_returns_404(self, client: TestClient) -> None:
        resp = client.get("/api/v1/embeddings/jobs/not-a-uuid")
        assert resp.status_code == 404


class TestGetSimilarTracks:
    """Tests for GET /api/v1/similarity/tracks/{track_id}."""

    def test_returns_similar_tracks(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        from app.domain.audio_features.similarity import SimilarityService

        mock_result = SimilarityQueryResult(
            tracks=[
                SimilarTrack(track_id="cand-1", similarity_score=0.85, matched_on=["audio"]),
                SimilarTrack(
                    track_id="cand-2",
                    similarity_score=0.72,
                    matched_on=["audio", "genre"],
                ),
            ],
            seed_track_id="seed-1",
            total_candidates=2,
        )

        with patch.object(SimilarityService, "find_similar", new=AsyncMock()) as mock_find:
            mock_find.return_value = mock_result
            resp = client.get("/api/v1/similarity/tracks/seed-1?limit=10")

        assert resp.status_code == 200
        data = resp.json()
        assert len(data["tracks"]) == 2
        assert data["tracks"][0]["track_id"] == "cand-1"
        assert data["tracks"][0]["similarity_score"] == 0.85
        assert data["seed_track_id"] == "seed-1"

    def test_with_genre_filter(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        from app.domain.audio_features.similarity import SimilarityService

        mock_result = SimilarityQueryResult(
            tracks=[SimilarTrack(track_id="cand-1", similarity_score=0.85, matched_on=["audio"])],
            seed_track_id="seed-1",
            total_candidates=1,
        )

        with patch.object(SimilarityService, "find_similar", new=AsyncMock()) as mock_find:
            mock_find.return_value = mock_result
            resp = client.get("/api/v1/similarity/tracks/seed-1?genre=rock")

        assert resp.status_code == 200
        assert mock_find.call_args[1].get("genre_filter") == ["rock"]


class TestSimilarFromSet:
    """Tests for POST /api/v1/similarity/from-set."""

    def test_returns_similar_tracks(
        self, client: TestClient, mock_session: MagicMock
    ) -> None:
        from app.domain.audio_features.similarity import SimilarityService

        mock_result = SimilarityQueryResult(
            tracks=[SimilarTrack(track_id="cand-1", similarity_score=0.90, matched_on=["audio"])],
            seed_track_id="seed-1,seed-2",
            total_candidates=1,
        )

        with patch.object(
            SimilarityService, "find_similar_to_multiple", new=AsyncMock()
        ) as mock_find:
            mock_find.return_value = mock_result
            resp = client.post(
                "/api/v1/similarity/from-set",
                json={"track_ids": ["seed-1", "seed-2"], "limit": 10},
            )

        assert resp.status_code == 200
        data = resp.json()
        assert len(data["tracks"]) == 1
        assert data["tracks"][0]["track_id"] == "cand-1"
