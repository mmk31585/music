"""
API endpoints for similarity search and embedding job management.

- ``GET  /api/v1/similarity/tracks/{track_id}`` — find similar tracks.
- ``POST /api/v1/similarity/from-set`` — find similar tracks (multi-seed).
- ``GET  /api/v1/similarity/tempo/bulk`` — bulk tempo_bpm lookup.
- ``POST /api/v1/embeddings/jobs`` — enqueue embedding-extraction job.
- ``GET  /api/v1/embeddings/jobs/{job_id}`` — query job status.
"""

from __future__ import annotations

import uuid
from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, Query, status
from pydantic import BaseModel, field_validator, model_validator
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.domain.audio_features.similarity import (
    SeedTrackNotFoundError,
    SimilarityError,
    SimilarityService,
)
from app.infra.db.models import EmbeddingJob
from app.infra.db.session import get_db
from app.workers.celery_app import celery_app

# Two routers: one for embedding jobs (prefix /embeddings),
# one for similarity (prefix /similarity). Kept in the same
# file because they share schemas and are conceptually related.
emb_router = APIRouter(prefix="/embeddings", tags=["embeddings"])
sim_router = APIRouter(prefix="/similarity", tags=["similarity"])


# ── Embedding job schemas ──────────────────────────────────────────────


class CreateEmbeddingJobRequest(BaseModel):
    track_id: str
    audio_file_path: str | None = None
    audio_download_url: str | None = None
    genre_ids: list[str] = []
    release_year: int | None = None

    @model_validator(mode="after")
    def _require_one_audio_source(self) -> CreateEmbeddingJobRequest:
        if not self.audio_file_path and not self.audio_download_url:
            raise ValueError(
                "Either audio_file_path or audio_download_url must be provided"
            )
        return self


class EmbeddingJobResponse(BaseModel):
    job_id: str
    track_id: str
    status: str
    created_at: datetime | None = None
    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, EmbeddingJob):
            data.job_id = str(data.id)
        return data


class EmbeddingJobDetailResponse(BaseModel):
    job_id: str
    track_id: str
    status: str
    audio_file_path: str | None = None
    audio_download_url: str | None = None
    genre_ids: list[str] = []
    release_year: int | None = None
    error_message: str | None = None
    retry_count: int = 0
    created_at: datetime | None = None
    updated_at: datetime | None = None
    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, EmbeddingJob):
            data.job_id = str(data.id)
        return data

    @field_validator("retry_count", mode="before")
    @classmethod
    def _none_to_default(cls, v: object) -> object:
        return 0 if v is None else v


# ── Similarity schemas ─────────────────────────────────────────────────


class SimilarTrackResponse(BaseModel):
    track_id: str
    similarity_score: float
    matched_on: list[str] = []


class SimilarTracksResponse(BaseModel):
    tracks: list[SimilarTrackResponse]
    seed_track_id: str = ""
    total_candidates: int = 0


class SimilarFromSetRequest(BaseModel):
    track_ids: list[str]
    limit: int = 20
    exclude: list[str] = []

    @field_validator("track_ids")
    @classmethod
    def _limit_track_ids(cls, v: list[str]) -> list[str]:
        if len(v) > 50:
            raise ValueError("Maximum 50 track IDs allowed per request")
        return v


# ── Embedding job endpoints ────────────────────────────────────────────


@emb_router.post(
    "/jobs",
    response_model=EmbeddingJobResponse,
    status_code=status.HTTP_202_ACCEPTED,
)
async def create_embedding_job(
    req: CreateEmbeddingJobRequest,
    db: AsyncSession = Depends(get_db),
) -> EmbeddingJobResponse:
    idempotency_key = f"embedding:{req.track_id}"

    result = await db.execute(
        select(EmbeddingJob).where(
            EmbeddingJob.idempotency_key == idempotency_key,
            EmbeddingJob.status.in_(["queued", "downloading", "extracting"]),
        )
    )
    existing_job = result.scalar_one_or_none()

    if existing_job is not None:
        return EmbeddingJobResponse.model_validate(existing_job)

    result = await db.execute(
        select(EmbeddingJob).where(
            EmbeddingJob.idempotency_key == idempotency_key,
            EmbeddingJob.status == "failed",
        )
    )
    failed_job = result.scalar_one_or_none()

    if failed_job is not None:
        failed_job.status = "queued"
        failed_job.error_message = None
        failed_job.retry_count = 0
        failed_job.audio_file_path = req.audio_file_path
        failed_job.audio_download_url = req.audio_download_url
        failed_job.genre_ids = req.genre_ids
        failed_job.release_year = req.release_year
        await db.commit()
        await db.refresh(failed_job)

        celery_app.send_task(
            "app.workers.tasks.embedding_tasks.extract_embedding_task",
            args=[str(failed_job.id)],
        )

        return EmbeddingJobResponse.model_validate(failed_job)

    job = EmbeddingJob(
        id=uuid.uuid4(),
        track_id=req.track_id,
        idempotency_key=idempotency_key,
        status="queued",
        audio_file_path=req.audio_file_path,
        audio_download_url=req.audio_download_url,
        genre_ids=req.genre_ids,
        release_year=req.release_year,
    )
    db.add(job)

    await db.commit()
    await db.refresh(job)

    celery_app.send_task(
        "app.workers.tasks.embedding_tasks.extract_embedding_task",
        args=[str(job.id)],
    )

    return EmbeddingJobResponse.model_validate(job)


@emb_router.get("/jobs/{job_id}", response_model=EmbeddingJobDetailResponse)
async def get_embedding_job(
    job_id: str,
    db: AsyncSession = Depends(get_db),
) -> EmbeddingJobDetailResponse:
    try:
        uid = uuid.UUID(job_id)
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Invalid job ID format",
        )

    result = await db.execute(select(EmbeddingJob).where(EmbeddingJob.id == uid))
    job = result.scalar_one_or_none()

    if job is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Job not found",
        )

    return EmbeddingJobDetailResponse.model_validate(job)


@emb_router.get(
    "/jobs/by-track/{track_id}",
    response_model=EmbeddingJobDetailResponse,
)
async def get_embedding_job_by_track(
    track_id: str,
    db: AsyncSession = Depends(get_db),
) -> EmbeddingJobDetailResponse:
    """Return the latest embedding job for a track, or 404 if none exists."""
    result = await db.execute(
        select(EmbeddingJob)
        .where(EmbeddingJob.track_id == track_id)
        .order_by(EmbeddingJob.created_at.desc())
        .limit(1)
    )
    job = result.scalar_one_or_none()

    if job is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="No embedding job found for this track",
        )

    return EmbeddingJobDetailResponse.model_validate(job)


# ── Similarity endpoints ───────────────────────────────────────────────


@sim_router.get("/tracks/{track_id}", response_model=SimilarTracksResponse)
async def get_similar_tracks(
    track_id: str,
    limit: int = 20,
    genre: str | None = None,
    db: AsyncSession = Depends(get_db),
) -> SimilarTracksResponse:
    if limit < 1:
        limit = 20
    if limit > 100:
        limit = 100

    genre_filter = [genre] if genre else None

    service = SimilarityService(db)
    try:
        result = await service.find_similar(
            track_id=track_id,
            limit=limit,
            genre_filter=genre_filter,
        )
    except SeedTrackNotFoundError as exc:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(exc),
        )
    except SimilarityError as exc:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(exc),
        )

    return SimilarTracksResponse(
        tracks=[
            SimilarTrackResponse(
                track_id=t.track_id,
                similarity_score=t.similarity_score,
                matched_on=t.matched_on,
            )
            for t in result.tracks
        ],
        seed_track_id=result.seed_track_id,
        total_candidates=result.total_candidates,
    )


@sim_router.post("/from-set", response_model=SimilarTracksResponse)
async def get_similar_to_track_set(
    req: SimilarFromSetRequest,
    db: AsyncSession = Depends(get_db),
) -> SimilarTracksResponse:
    limit = req.limit
    if limit < 1:
        limit = 20
    if limit > 100:
        limit = 100

    service = SimilarityService(db)
    try:
        result = await service.find_similar_to_multiple(
            track_ids=req.track_ids,
            limit=limit,
            exclude_track_ids=req.exclude,
        )
    except SeedTrackNotFoundError as exc:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(exc),
        )
    except SimilarityError as exc:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(exc),
        )

    return SimilarTracksResponse(
        tracks=[
            SimilarTrackResponse(
                track_id=t.track_id,
                similarity_score=t.similarity_score,
                matched_on=t.matched_on,
            )
            for t in result.tracks
        ],
        seed_track_id=result.seed_track_id,
        total_candidates=result.total_candidates,
    )


# ── Bulk tempo endpoint ─────────────────────────────────────────────────


class BulkTempoResponse(BaseModel):
    tempos: dict[str, float | None]


@sim_router.get("/tempo/bulk", response_model=BulkTempoResponse)
async def get_bulk_tempos(
    track_ids: str = Query(..., description="Comma-separated track UUIDs"),
    db: AsyncSession = Depends(get_db),
) -> BulkTempoResponse:
    """Return tempo_bpm for a batch of tracks.

    Used by the Go backend's taste profile recomputation.  Track IDs
    without embeddings are returned as ``None``.
    """
    ids = [tid.strip() for tid in track_ids.split(",") if tid.strip()]
    if not ids:
        return BulkTempoResponse(tempos={})

    service = SimilarityService(db)
    tempos = await service.get_bulk_tempos(ids)
    return BulkTempoResponse(tempos=tempos)
