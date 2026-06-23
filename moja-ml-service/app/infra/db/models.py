"""
SQLAlchemy ORM models for moja-ml-service's OWN database.

This database is separate from the Go backend's catalog database. It tracks
job lifecycle for ML processing tasks.
"""

from __future__ import annotations

import uuid
from datetime import datetime

from pgvector.sqlalchemy import Vector
from sqlalchemy import ARRAY, String, Text, func
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    """Declarative base for all models in this service's database."""


class LyricsJob(Base):
    """Tracks the lifecycle of a single lyrics-generation job."""

    __tablename__ = "lyrics_jobs"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    track_id: Mapped[str] = mapped_column(index=True)
    # idempotency_key is deterministic: f"lyrics:{track_id}"
    idempotency_key: Mapped[str] = mapped_column(unique=True, index=True)

    # Status values: queued | downloading | transcribing | formatting
    #                | callback_pending | completed | failed
    status: Mapped[str] = mapped_column(default="queued")

    # Audio source — exactly one should be non-null per the storage mode
    audio_file_path: Mapped[str | None] = mapped_column(nullable=True)
    audio_download_url: Mapped[str | None] = mapped_column(nullable=True)

    # Optional metadata for LRC headers
    track_title: Mapped[str | None] = mapped_column(nullable=True)
    track_artist: Mapped[str | None] = mapped_column(nullable=True)

    # Results (populated on successful completion)
    result_lrc_content: Mapped[str | None] = mapped_column(Text, nullable=True)
    result_plain_text: Mapped[str | None] = mapped_column(Text, nullable=True)
    result_confidence: Mapped[float | None] = mapped_column(nullable=True)
    result_detected_language: Mapped[str | None] = mapped_column(nullable=True)
    result_segment_count: Mapped[int | None] = mapped_column(nullable=True)

    # Error handling
    error_message: Mapped[str | None] = mapped_column(nullable=True)
    retry_count: Mapped[int] = mapped_column(default=0)
    callback_delivered: Mapped[bool] = mapped_column(default=False)
    whisper_model_version: Mapped[str | None] = mapped_column(nullable=True)

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        default=func.now(),
        onupdate=func.now(),
    )


class CoverJob(Base):
    """Tracks the lifecycle of a cover image optimization job."""

    __tablename__ = "cover_jobs"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    album_id: Mapped[str] = mapped_column(index=True)
    idempotency_key: Mapped[str] = mapped_column(unique=True, index=True)

    # Status: queued | downloading | optimizing | storing | callback_pending | completed | failed
    status: Mapped[str] = mapped_column(default="queued")

    # Source
    cover_url: Mapped[str | None] = mapped_column(nullable=True)
    cover_path: Mapped[str | None] = mapped_column(nullable=True)
    track_id: Mapped[str | None] = mapped_column(nullable=True)

    # Results (populated on successful completion)
    result_cover_url: Mapped[str | None] = mapped_column(nullable=True)
    result_thumb_url: Mapped[str | None] = mapped_column(nullable=True)
    result_med_url: Mapped[str | None] = mapped_column(nullable=True)
    original_width: Mapped[int | None] = mapped_column(nullable=True)
    original_height: Mapped[int | None] = mapped_column(nullable=True)
    webp_bytes_saved: Mapped[int | None] = mapped_column(nullable=True)

    # Error handling
    error_message: Mapped[str | None] = mapped_column(nullable=True)
    retry_count: Mapped[int] = mapped_column(default=0)
    callback_delivered: Mapped[bool] = mapped_column(default=False)

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        default=func.now(),
        onupdate=func.now(),
    )


class TrackEmbedding(Base):
    """Persists a single track's audio embedding vector and associated metadata.

    The ``embedding`` column uses pgvector's ``Vector(512)`` type for
    efficient cosine-distance nearest-neighbor queries.  An HNSW index
    on this column enables fast approximate nearest neighbor (ANN) search
    at the scale of thousands to low millions of tracks.

    Cosine distance is the standard comparison metric for audio embeddings
    from models like OpenL3 (Phase 1): direction matters more than
    magnitude, which makes cosine similarity the natural choice over
    Euclidean (L2) distance for this embedding family.
    """

    __tablename__ = "track_embeddings_audio"

    track_id: Mapped[str] = mapped_column(primary_key=True)  # Go's track UUID
    embedding: Mapped[list[float]] = mapped_column(Vector(512))
    embedding_model: Mapped[str] = mapped_column()  # e.g. "openl3-mel256-music-512"
    tempo_bpm: Mapped[float | None] = mapped_column(nullable=True)
    genre_ids: Mapped[list[str]] = mapped_column(ARRAY(String), default=list)
    release_year: Mapped[int | None] = mapped_column(nullable=True)

    extracted_at: Mapped[datetime] = mapped_column(default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        default=func.now(),
        onupdate=func.now(),
    )


class VideoProcessingJob(Base):
    """Tracks the lifecycle of a single video-processing job.

    Status values:
        queued | downloading | processing | uploading | completed | failed

    The ``idempotency_key`` is deterministic: ``f"video:{video_id}"``.
    """

    __tablename__ = "video_processing_jobs"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    video_id: Mapped[str] = mapped_column(index=True)  # Go's video UUID
    idempotency_key: Mapped[str] = mapped_column(unique=True, index=True)
    status: Mapped[str] = mapped_column(default="queued")

    # Source file references
    raw_video_path: Mapped[str | None] = mapped_column(nullable=True)
    raw_video_download_url: Mapped[str | None] = mapped_column(nullable=True)
    audio_file_path: Mapped[str | None] = mapped_column(nullable=True)
    audio_download_url: Mapped[str | None] = mapped_column(nullable=True)
    track_start_ms: Mapped[int] = mapped_column(default=0)

    # Results (populated on successful processing)
    result_video_path: Mapped[str | None] = mapped_column(nullable=True)
    result_thumbnail_path: Mapped[str | None] = mapped_column(nullable=True)
    result_duration_ms: Mapped[int | None] = mapped_column(nullable=True)
    result_aspect_ratio: Mapped[str | None] = mapped_column(nullable=True)

    # Error handling
    error_message: Mapped[str | None] = mapped_column(nullable=True)
    retry_count: Mapped[int] = mapped_column(default=0)
    callback_delivered: Mapped[bool] = mapped_column(default=False)

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        default=func.now(),
        onupdate=func.now(),
    )


class EmbeddingJob(Base):
    """Tracks the lifecycle of a single embedding-extraction job.

    Mirrors the structure of ``LyricsJob`` — same idempotency-key pattern,
    same status transitions, same retry / error-handling discipline.
    """

    __tablename__ = "embedding_jobs"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    track_id: Mapped[str] = mapped_column(index=True)
    # idempotency_key is deterministic: f"embedding:{track_id}"
    idempotency_key: Mapped[str] = mapped_column(unique=True, index=True)

    # Status values: queued | downloading | extracting | completed | failed
    status: Mapped[str] = mapped_column(default="queued")

    # Audio source — exactly one should be non-null per the storage mode
    audio_file_path: Mapped[str | None] = mapped_column(nullable=True)
    audio_download_url: Mapped[str | None] = mapped_column(nullable=True)

    # Optional metadata passed alongside the audio source
    genre_ids: Mapped[list[str]] = mapped_column(ARRAY(String), default=list)
    release_year: Mapped[int | None] = mapped_column(nullable=True)

    # Error handling
    error_message: Mapped[str | None] = mapped_column(nullable=True)
    retry_count: Mapped[int] = mapped_column(default=0)

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        default=func.now(),
        onupdate=func.now(),
    )
