"""
Similarity search service for pgvector-based audio embedding queries.

Queries the ``track_embeddings`` table for nearest neighbors using cosine
distance (the standard metric for OpenL3-style audio embeddings —
direction matters more than magnitude for this embedding family).  A
genre-matching metadata boost is applied for re-ranking.

Domain layer — no HTTP or Celery knowledge.  Uses the async SQLAlchemy
session for compatibility with the FastAPI application (Celery tasks use
a synchronous session, so they should not call this service directly).
"""

from __future__ import annotations

import logging
from dataclasses import dataclass, field
from typing import Any

import numpy as np
from sqlalchemy import select, text
from sqlalchemy.ext.asyncio import AsyncSession

from app.infra.db.models import TrackEmbedding

logger = logging.getLogger(__name__)

# ── Tunable constants ──────────────────────────────────────────────────

# Small score boost applied when a candidate track shares a genre with the
# seed track.  Tuned to be noticeable but not dominant — enough to break
# ties between equally audio-similar tracks without overwhelming the
# primary audio-similarity signal.
GENRE_BOOST: float = 0.05

# When fetching candidate nearest neighbors from pgvector, request this
# many extra candidates (multiples of ``limit``) so that the metadata
# boost / filtering step has room to re-rank within a reasonable set.
# 3x is a conservative default — on a 10M-track catalog this means 60
# candidates for a limit=20 query, trivially fast with the HNSW index.
CANDIDATE_MULTIPLIER: int = 3


# ── Custom exceptions ─────────────────────────────────────────────────


class SimilarityError(Exception):
    """Raised when similarity search fails for any reason."""


class SeedTrackNotFoundError(SimilarityError):
    """Raised when the seed track's embedding is not in the index."""


# ── Dataclasses ───────────────────────────────────────────────────────


@dataclass
class SimilarTrack:
    """One track found by a similarity search query.

    Attributes:
        track_id: The Go backend's track UUID.
        similarity_score: Cosine-similarity-based relevance, 0.0–1.0
            (higher = more similar), **after** metadata boost.
        matched_on: Transparency tags explaining WHY this track matched,
            e.g. ``["audio"]`` or ``["audio", "genre"]``.  Useful for
            later UX ("Recommended because it sounds like X" vs
            "Recommended because it shares a genre").
    """

    track_id: str
    similarity_score: float = 0.0
    matched_on: list[str] = field(default_factory=list)


@dataclass
class SimilarityQueryResult:
    """The full result of a similarity query, including debug info."""

    tracks: list[SimilarTrack] = field(default_factory=list)
    seed_track_id: str = ""
    seed_genre_ids: list[str] = field(default_factory=list)
    total_candidates: int = 0


# ── Service ───────────────────────────────────────────────────────────


class SimilarityService:
    """Query pgvector for similar tracks and apply metadata re-ranking.

    Usage::

        service = SimilarityService(db_session)
        results = await service.find_similar(
            track_id="abc-123",
            limit=20,
            genre_filter=["rock", "pop"],
        )
    """

    def __init__(self, db_session: AsyncSession) -> None:
        self._session = db_session

    async def _get_seed_embedding(
        self, track_id: str
    ) -> tuple[list[float], list[str]]:
        """Load the seed track's embedding and genre IDs.

        Returns:
            A tuple of (embedding_vector, genre_ids).

        Raises:
            SeedTrackNotFoundError: If the track is not in the embedding index.
        """
        result = await self._session.execute(
            select(TrackEmbedding).where(TrackEmbedding.track_id == track_id)
        )
        row = result.scalar_one_or_none()
        if row is None:
            raise SeedTrackNotFoundError(
                f"Track {track_id} has no embedding in the index. "
                "Run an embedding-extraction job first."
            )
        return row.embedding, row.genre_ids or []

    async def _query_nearest_neighbors(
        self,
        seed_embedding: list[float],
        seed_track_id: str | None = None,
        exclude_track_ids: list[str] | None = None,
        genre_filter: list[str] | None = None,
        limit: int = 20,
    ) -> list[dict[str, Any]]:
        """Query pgvector HNSW index for nearest neighbors using cosine distance.

        Returns raw rows with track_id, similarity, genre_ids so the
        caller can apply metadata boost and re-rank.

        Args:
            seed_embedding: The 512-dim query vector.
            seed_track_id: The seed track ID to exclude from results.
            exclude_track_ids: Additional tracks to exclude from results.
            genre_filter: If provided, ONLY tracks matching ANY of these
                genres are returned (hard filter, not a boost).
            limit: Maximum number of results.

        Returns:
            List of dicts with keys: track_id, similarity, genre_ids.
        """
        candidate_limit = limit * CANDIDATE_MULTIPLIER

        # Build the WHERE clause dynamically.
        conditions: list[str] = []
        params: dict[str, Any] = {
            "query_vec": seed_embedding,
            "candidate_limit": candidate_limit,
        }

        if seed_track_id:
            conditions.append("track_id != :seed_track_id")
            params["seed_track_id"] = seed_track_id

        if exclude_track_ids:
            conditions.append("track_id != ALL(:exclude_ids)")
            params["exclude_ids"] = exclude_track_ids

        if genre_filter:
            conditions.append("genre_ids && :genre_filter")
            params["genre_filter"] = genre_filter

        where_clause = " AND ".join(conditions) if conditions else "TRUE"

        # Cosine distance: ``embedding <=> :query_vec`` returns 0.0 for
        # identical vectors, 2.0 for opposite.  ``1 - distance`` gives
        # a similarity score in [0, 1].
        sql = text(f"""
            SELECT
                track_id,
                1 - (embedding <=> :query_vec) AS similarity,
                genre_ids
            FROM track_embeddings_audio
            WHERE {where_clause}
            ORDER BY embedding <=> :query_vec
            LIMIT :candidate_limit
        """)

        result = await self._session.execute(sql, params)
        rows = result.all()
        return [
            {
                "track_id": row[0],
                "similarity": float(row[1]) if row[1] is not None else 0.0,
                "genre_ids": row[2] or [],
            }
            for row in rows
        ]

    async def find_similar(
        self,
        track_id: str,
        limit: int = 20,
        genre_filter: list[str] | None = None,
        exclude_track_ids: list[str] | None = None,
    ) -> SimilarityQueryResult:
        """Find tracks similar to a single seed track.

        Process:
        1. Look up the seed track's embedding and genre IDs.
        2. Query pgvector HNSW for nearest neighbors (fetching extra
           candidates for re-ranking room).
        3. Apply genre-boost: tracks sharing at least one genre with the
           seed get ``+GENRE_BOOST`` on their similarity score.
        4. Apply ``exclude_track_ids`` filter (e.g. exclude tracks already
           in the user's recent history).
        5. Apply ``genre_filter`` if provided (hard filter, not a boost).
        6. Sort by final score descending, truncate to ``limit``.

        Args:
            track_id: The seed track's UUID.
            limit: Max results to return (default 20).
            genre_filter: Optional hard genre filter — only tracks matching
                ANY of these genre IDs are returned.
            exclude_track_ids: Optional list of track IDs to exclude from
                results.

        Returns:
            A ``SimilarityQueryResult`` with the ranked tracks.
        """
        seed_embedding, seed_genres = await self._get_seed_embedding(track_id)

        candidates = await self._query_nearest_neighbors(
            seed_embedding=seed_embedding,
            seed_track_id=track_id,
            exclude_track_ids=exclude_track_ids,
            genre_filter=genre_filter,
            limit=limit,
        )

        # ── Apply genre boost and build result ────────────────────
        scored: list[SimilarTrack] = []
        for cand in candidates:
            matched_on = ["audio"]
            score = cand["similarity"]

            # Boost: if the candidate shares any genre with the seed
            if seed_genres and any(g in seed_genres for g in cand["genre_ids"]):
                score += GENRE_BOOST
                matched_on.append("genre")

            scored.append(
                SimilarTrack(
                    track_id=cand["track_id"],
                    similarity_score=round(min(score, 1.0), 6),
                    matched_on=matched_on,
                )
            )

        # Sort by score descending and truncate
        scored.sort(key=lambda t: t.similarity_score, reverse=True)
        scored = scored[:limit]

        return SimilarityQueryResult(
            tracks=scored,
            seed_track_id=track_id,
            seed_genre_ids=seed_genres,
            total_candidates=len(candidates),
        )

    async def find_similar_to_multiple(
        self,
        track_ids: list[str],
        limit: int = 20,
        exclude_track_ids: list[str] | None = None,
    ) -> SimilarityQueryResult:
        """Find tracks similar to MULTIPLE seed tracks (centroid approach).

        Averages the embeddings of all seed tracks into one centroid
        vector, then queries pgvector for the nearest neighbors of that
        centroid.  Genre boosting uses the union of all seed genre IDs.

        This is what taste profile / Home feed recommendations will use
        — a user's taste isn't one track, it's an aggregate of many.

        Args:
            track_ids: List of seed track UUIDs.
            limit: Max results to return (default 20).
            exclude_track_ids: Optional list of track IDs to exclude
                (including the seed tracks themselves).

        Returns:
            A ``SimilarityQueryResult`` with the ranked tracks.
        """
        if not track_ids:
            raise ValueError("At least one track_id is required")

        # Load all seed embeddings
        embeddings: list[list[float]] = []
        union_genres: set[str] = set()
        for tid in track_ids:
            emb, genres = await self._get_seed_embedding(tid)
            embeddings.append(emb)
            union_genres.update(genres)

        # Compute centroid (average all embeddings)
        centroid = np.mean(embeddings, axis=0).tolist()

        # Ensure seed tracks are excluded
        all_exclude = list(set(exclude_track_ids or []) | set(track_ids))

        candidates = await self._query_nearest_neighbors(
            seed_embedding=centroid,
            exclude_track_ids=all_exclude,
            limit=limit,
        )

        # ── Apply genre boost ─────────────────────────────────────
        scored: list[SimilarTrack] = []
        for cand in candidates:
            matched_on = ["audio"]
            score = cand["similarity"]

            if union_genres and any(g in union_genres for g in cand["genre_ids"]):
                score += GENRE_BOOST
                matched_on.append("genre")

            scored.append(
                SimilarTrack(
                    track_id=cand["track_id"],
                    similarity_score=round(min(score, 1.0), 6),
                    matched_on=matched_on,
                )
            )

        scored.sort(key=lambda t: t.similarity_score, reverse=True)
        scored = scored[:limit]

        return SimilarityQueryResult(
            tracks=scored,
            seed_track_id=",".join(track_ids),
            seed_genre_ids=list(union_genres),
            total_candidates=len(candidates),
        )

    async def get_bulk_tempos(
        self,
        track_ids: list[str],
    ) -> dict[str, float | None]:
        """Return tempo_bpm for a batch of track IDs.

        Used by the Go backend's taste profile recomputation to compute
        ``avg_tempo_preference``.  Missing embeddings return ``None``.

        Args:
            track_ids: List of track UUIDs to look up.

        Returns:
            Dict mapping track_id → tempo_bpm (or None if not found).
        """
        if not track_ids:
            return {}

        result = await self._session.execute(
            select(TrackEmbedding.track_id, TrackEmbedding.tempo_bpm).where(
                TrackEmbedding.track_id.in_(track_ids)
            )
        )
        rows = result.all()
        tempos: dict[str, float | None] = {tid: None for tid in track_ids}
        for row in rows:
            tempos[row[0]] = row[1]
        return tempos
