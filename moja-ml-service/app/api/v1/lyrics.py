"""
API endpoints for lyrics job management.

- ``POST /api/v1/lyrics/jobs`` — enqueue a new lyrics-generation job
  (idempotent: same ``track_id`` returns existing job).
- ``GET /api/v1/lyrics/jobs/{job_id}`` — query job status and results.
"""

from __future__ import annotations

import uuid
from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, field_validator, model_validator
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.infra.db.models import LyricsJob
from app.infra.db.session import get_db
from app.workers.celery_app import celery_app

router = APIRouter(prefix="/lyrics", tags=["lyrics"])


# ── Request / response schemas ─────────────────────────────────────────


class CreateLyricsJobRequest(BaseModel):
    """Payload for enqueuing a new lyrics-generation job."""

    track_id: str
    audio_file_path: str | None = None
    audio_download_url: str | None = None
    track_title: str | None = None
    track_artist: str | None = None

    @model_validator(mode="after")
    def _require_one_audio_source(self) -> CreateLyricsJobRequest:
        if not self.audio_file_path and not self.audio_download_url:
            raise ValueError(
                "Either audio_file_path or audio_download_url must be provided"
            )
        return self


class LyricsJobResponse(BaseModel):
    """Summary response returned when creating a job."""

    job_id: str
    track_id: str
    status: str
    created_at: datetime | None = None

    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, LyricsJob):
            data.job_id = str(data.id)  # type: ignore[attr-defined]
        return data


class LyricsJobDetailResponse(BaseModel):
    """Full job detail including results if completed."""

    job_id: str
    track_id: str
    status: str
    audio_file_path: str | None = None
    audio_download_url: str | None = None
    track_title: str | None = None
    track_artist: str | None = None
    result_lrc_content: str | None = None
    result_plain_text: str | None = None
    result_confidence: float | None = None
    result_detected_language: str | None = None
    result_segment_count: int | None = None
    error_message: str | None = None
    retry_count: int = 0
    callback_delivered: bool = False
    created_at: datetime | None = None
    updated_at: datetime | None = None

    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, LyricsJob):
            data.job_id = str(data.id)  # type: ignore[attr-defined]
        return data

    @field_validator("retry_count", "callback_delivered", mode="before")
    @classmethod
    def _none_to_default(cls, v: object) -> object:
        return 0 if v is None else v


# ── Endpoints ──────────────────────────────────────────────────────────


@router.post("/jobs", response_model=LyricsJobResponse, status_code=status.HTTP_202_ACCEPTED)
async def create_lyrics_job(
    req: CreateLyricsJobRequest,
    db: AsyncSession = Depends(get_db),
) -> LyricsJobResponse:
    """Enqueue a new lyrics-generation job.

    **Idempotent:** if a job for this ``track_id`` already exists and is
    not in a terminal ``failed`` state, the existing job is returned
    instead of creating a duplicate. The idempotency key is derived as
    ``lyrics:{track_id}`` — if the Go backend enqueues the same track
    twice (e.g. an admin double-clicks), the second request returns the
    first job's status with ``202 Accepted``.
    """
    idempotency_key = f"lyrics:{req.track_id}"

    # Check for existing non-failed job
    result = await db.execute(
        select(LyricsJob).where(
            LyricsJob.idempotency_key == idempotency_key,
            LyricsJob.status.in_(["queued", "downloading", "transcribing", "formatting", "callback_pending"]),
        )
    )
    existing_job = result.scalar_one_or_none()

    if existing_job is not None:
        return LyricsJobResponse.model_validate(existing_job)

    # Check for a failed job with same key — reset and re-queue it
    result = await db.execute(
        select(LyricsJob).where(
            LyricsJob.idempotency_key == idempotency_key,
            LyricsJob.status == "failed",
        )
    )
    failed_job = result.scalar_one_or_none()

    if failed_job is not None:
        # Reset the failed job for retry
        failed_job.status = "queued"
        failed_job.error_message = None
        failed_job.retry_count = 0
        failed_job.audio_file_path = req.audio_file_path
        failed_job.audio_download_url = req.audio_download_url
        failed_job.track_title = req.track_title
        failed_job.track_artist = req.track_artist
        await db.commit()
        await db.refresh(failed_job)

        # Enqueue Celery task
        celery_app.send_task(
            "app.workers.tasks.lyrics_tasks.generate_lyrics_task",
            args=[str(failed_job.id)],
        )

        return LyricsJobResponse.model_validate(failed_job)

    # Create new job
    job = LyricsJob(
        id=uuid.uuid4(),
        track_id=req.track_id,
        idempotency_key=idempotency_key,
        status="queued",
        audio_file_path=req.audio_file_path,
        audio_download_url=req.audio_download_url,
        track_title=req.track_title,
        track_artist=req.track_artist,
    )
    db.add(job)

    # Commit BEFORE enqueuing — never enqueue before the row exists.
    # If Celery picks up the task before the row is visible (due to
    # transaction isolation), the task will fail to find the job.
    await db.commit()
    await db.refresh(job)

    # Enqueue Celery task
    celery_app.send_task(
        "app.workers.tasks.lyrics_tasks.generate_lyrics_task",
        args=[str(job.id)],
    )

    return LyricsJobResponse.model_validate(job)


@router.get("/jobs/{job_id}", response_model=LyricsJobDetailResponse)
async def get_lyrics_job(
    job_id: str,
    db: AsyncSession = Depends(get_db),
) -> LyricsJobDetailResponse:
    """Return full job state including result fields if completed."""
    try:
        uid = uuid.UUID(job_id)
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Invalid job ID format",
        )

    result = await db.execute(select(LyricsJob).where(LyricsJob.id == uid))
    job = result.scalar_one_or_none()

    if job is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Job not found",
        )

    return LyricsJobDetailResponse.model_validate(job)


@router.get("/jobs/by-track/{track_id}", response_model=LyricsJobDetailResponse)
async def get_lyrics_job_by_track(
    track_id: str,
    db: AsyncSession = Depends(get_db),
) -> LyricsJobDetailResponse:
    """Return the latest lyrics job for a track, or 404 if none exists.

    Useful for the admin UI to poll AI generation progress.
    """
    result = await db.execute(
        select(LyricsJob)
        .where(LyricsJob.track_id == track_id)
        .order_by(LyricsJob.created_at.desc())
        .limit(1)
    )
    job = result.scalar_one_or_none()

    if job is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="No job found for this track",
        )

    return LyricsJobDetailResponse.model_validate(job)
