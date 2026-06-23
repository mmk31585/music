"""Tests for SimilarityService.

Uses mocked async DB sessions — no real postgres or pgvector needed.
"""

from __future__ import annotations

from unittest.mock import AsyncMock, MagicMock

import pytest

from app.domain.audio_features.similarity import (
    GENRE_BOOST,
    SeedTrackNotFoundError,
    SimilarityQueryResult,
    SimilarityService,
)
from app.infra.db.models import TrackEmbedding


@pytest.fixture
def mock_session() -> MagicMock:
    session = MagicMock()
    session.execute = AsyncMock()
    session.execute.return_value = MagicMock()
    session.execute.return_value.scalar_one_or_none = MagicMock(return_value=None)
    session.execute.return_value.all = MagicMock(return_value=[])
    return session


class TestFindSimilar:
    """Tests for SimilarityService.find_similar()."""

    @pytest.mark.asyncio
    async def test_raises_for_unknown_seed(self, mock_session: MagicMock) -> None:
        """Unknown seed track_id should raise SeedTrackNotFoundError."""
        mock_session.execute.return_value.scalar_one_or_none.return_value = None

        service = SimilarityService(mock_session)
        with pytest.raises(SeedTrackNotFoundError):
            await service.find_similar(track_id="nonexistent")

    @pytest.mark.asyncio
    async def test_returns_similar_tracks(self, mock_session: MagicMock) -> None:
        """Happy path: seed found, NN query returns candidates."""
        # First query: seed exists
        seed_mock = MagicMock(spec=TrackEmbedding)
        seed_mock.embedding = [0.1, 0.2, 0.3]
        seed_mock.genre_ids = ["rock"]
        mock_session.execute.return_value.scalar_one_or_none.return_value = seed_mock

        # Second query: NN returns two candidates
        mock_session.execute.return_value.all.return_value = [
            ("cand-1", 0.85, ["rock"]),
            ("cand-2", 0.72, ["jazz"]),
        ]

        service = SimilarityService(mock_session)
        result = await service.find_similar(track_id="seed-1", limit=10)

        assert isinstance(result, SimilarityQueryResult)
        assert len(result.tracks) == 2
        assert result.tracks[0].track_id == "cand-1"
        assert result.tracks[1].track_id == "cand-2"
        assert result.seed_track_id == "seed-1"
        assert result.total_candidates == 2

    @pytest.mark.asyncio
    async def test_genre_boost_affects_ranking(self, mock_session: MagicMock) -> None:
        """Genre-matching candidates should get a score boost."""
        seed_mock = MagicMock(spec=TrackEmbedding)
        seed_mock.embedding = [0.1, 0.2, 0.3]
        seed_mock.genre_ids = ["rock", "pop"]
        mock_session.execute.return_value.scalar_one_or_none.return_value = seed_mock

        mock_session.execute.return_value.all.return_value = [
            ("cand-1", 0.80, ["rock"]),
            ("cand-2", 0.79, ["jazz"]),
        ]

        service = SimilarityService(mock_session)
        result = await service.find_similar(track_id="seed-1", limit=10)

        assert result.tracks[0].track_id == "cand-1"
        assert result.tracks[0].similarity_score == pytest.approx(0.80 + GENRE_BOOST, rel=1e-3)
        assert "genre" in result.tracks[0].matched_on
        assert "genre" not in result.tracks[1].matched_on

    @pytest.mark.asyncio
    async def test_genre_filter_hard_filters(self, mock_session: MagicMock) -> None:
        """Genre filter should hard-filter, not boost."""
        seed_mock = MagicMock(spec=TrackEmbedding)
        seed_mock.embedding = [0.1, 0.2, 0.3]
        seed_mock.genre_ids = ["rock"]
        mock_session.execute.return_value.scalar_one_or_none.return_value = seed_mock

        mock_session.execute.return_value.all.return_value = [
            ("cand-1", 0.80, ["rock"]),
        ]

        service = SimilarityService(mock_session)
        result = await service.find_similar(
            track_id="seed-1", limit=10, genre_filter=["rock"]
        )

        assert len(result.tracks) == 1
        assert result.tracks[0].track_id == "cand-1"


class TestFindSimilarToMultiple:
    """Tests for SimilarityService.find_similar_to_multiple()."""

    @pytest.mark.asyncio
    async def test_centroid_search_from_multiple_seeds(
        self, mock_session: MagicMock
    ) -> None:
        seed_1 = MagicMock(spec=TrackEmbedding)
        seed_1.embedding = [0.1, 0.2, 0.3]
        seed_1.genre_ids = ["rock"]

        seed_2 = MagicMock(spec=TrackEmbedding)
        seed_2.embedding = [0.4, 0.5, 0.6]
        seed_2.genre_ids = ["pop"]

        # Seed lookup returns each seed in sequence
        mock_session.execute.return_value.scalar_one_or_none.side_effect = [
            seed_1,
            seed_2,
        ]

        # NN query
        mock_session.execute.return_value.all.return_value = [
            ("cand-1", 0.90, ["rock"]),
            ("cand-2", 0.75, ["jazz"]),
        ]

        service = SimilarityService(mock_session)
        result = await service.find_similar_to_multiple(
            track_ids=["seed-1", "seed-2"], limit=10
        )

        assert len(result.tracks) == 2
        assert result.tracks[0].track_id == "cand-1"
