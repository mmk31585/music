"""
API endpoints for video job management.

- ``POST /api/v1/video/jobs`` — enqueue a new video-processing job
  (idempotent: same ``video_id`` returns existing job).
- ``GET /api/v1/video/jobs/{job_id}`` — query job status and results.
"""

from __future__ import annotations

import uuid
from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, field_validator, model_validator
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.infra.db.models import VideoProcessingJob
from app.infra.db.session import get_db
from app.workers.celery_app import celery_app

router = APIRouter(prefix="/video", tags=["video"])


# ── Request / response schemas ─────────────────────────────────────────


class CreateVideoJobRequest(BaseModel):
    """Payload for enqueuing a new video-processing job."""

    video_id: str
    raw_video_path: str | None = None
    raw_video_download_url: str | None = None
    audio_file_path: str | None = None
    audio_download_url: str | None = None
    track_start_ms: int = 0

    @model_validator(mode="after")
    def _require_one_video_source(self) -> CreateVideoJobRequest:
        if (
            not self.raw_video_path
            and not self.raw_video_download_url
        ):
            raise ValueError(
                "Either raw_video_path or raw_video_download_url "
                "must be provided"
            )
        return self

    @model_validator(mode="after")
    def _require_one_audio_source(self) -> CreateVideoJobRequest:
        if (
            not self.audio_file_path
            and not self.audio_download_url
        ):
            raise ValueError(
                "Either audio_file_path or audio_download_url "
                "must be provided"
            )
        return self


class VideoJobResponse(BaseModel):
    """Summary response returned when creating a job."""

    job_id: str
    video_id: str
    status: str
    created_at: datetime | None = None

    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, VideoProcessingJob):
            data.job_id = str(data.id)  # type: ignore[attr-defined]
        return data


class VideoJobDetailResponse(BaseModel):
    """Full job detail including results if completed."""

    job_id: str
    video_id: str
    status: str
    raw_video_path: str | None = None
    raw_video_download_url: str | None = None
    audio_file_path: str | None = None
    audio_download_url: str | None = None
    track_start_ms: int = 0
    result_video_path: str | None = None
    result_thumbnail_path: str | None = None
    result_duration_ms: int | None = None
    result_aspect_ratio: str | None = None
    error_message: str | None = None
    retry_count: int = 0
    callback_delivered: bool = False
    created_at: datetime | None = None
    updated_at: datetime | None = None

    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, VideoProcessingJob):
            data.job_id = str(data.id)  # type: ignore[attr-defined]
        return data

    @field_validator("retry_count", "callback_delivered", mode="before")
    @classmethod
    def _none_to_default(cls, v: object) -> object:
        return 0 if v is None else v


# ── Endpoints ──────────────────────────────────────────────────────────


@router.post(
    "/jobs",
    response_model=VideoJobResponse,
    status_code=status.HTTP_202_ACCEPTED,
)
async def create_video_job(
    req: CreateVideoJobRequest,
    db: AsyncSession = Depends(get_db),
) -> VideoJobResponse:
    """Enqueue a new video-processing job.

    **Idempotent:** if a job for this ``video_id`` already exists and is
    not in a terminal ``failed`` state, the existing job is returned
    instead of creating a duplicate. The idempotency key is derived as
    ``video:{video_id}``.
    """
    idempotency_key = f"video:{req.video_id}"

    # Check for existing non-failed job
    result = await db.execute(
        select(VideoProcessingJob).where(
            VideoProcessingJob.idempotency_key == idempotency_key,
            VideoProcessingJob.status.in_(
                [
                    "queued",
                    "downloading",
                    "processing",
                    "uploading",
                ]
            ),
        )
    )
    existing_job = result.scalar_one_or_none()

    if existing_job is not None:
        return VideoJobResponse.model_validate(existing_job)

    # Check for a failed job with same key — reset and re-queue it
    result = await db.execute(
        select(VideoProcessingJob).where(
            VideoProcessingJob.idempotency_key == idempotency_key,
            VideoProcessingJob.status == "failed",
        )
    )
    failed_job = result.scalar_one_or_none()

    if failed_job is not None:
        # Reset the failed job for retry
        failed_job.status = "queued"
        failed_job.error_message = None
        failed_job.retry_count = 0
        failed_job.raw_video_path = req.raw_video_path
        failed_job.raw_video_download_url = req.raw_video_download_url
        failed_job.audio_file_path = req.audio_file_path
        failed_job.audio_download_url = req.audio_download_url
        failed_job.track_start_ms = req.track_start_ms
        await db.commit()
        await db.refresh(failed_job)

        # Enqueue Celery task
        celery_app.send_task(
            "app.workers.tasks.video_tasks.process_video_task",
            args=[str(failed_job.id)],
        )

        return VideoJobResponse.model_validate(failed_job)

    # Create new job
    job = VideoProcessingJob(
        id=uuid.uuid4(),
        video_id=req.video_id,
        idempotency_key=idempotency_key,
        status="queued",
        raw_video_path=req.raw_video_path,
        raw_video_download_url=req.raw_video_download_url,
        audio_file_path=req.audio_file_path,
        audio_download_url=req.audio_download_url,
        track_start_ms=req.track_start_ms,
    )
    db.add(job)

    # Commit BEFORE enqueuing — never enqueue before the row exists
    await db.commit()
    await db.refresh(job)

    # Enqueue Celery task
    celery_app.send_task(
        "app.workers.tasks.video_tasks.process_video_task",
        args=[str(job.id)],
    )

    return VideoJobResponse.model_validate(job)


@router.get(
    "/jobs/{job_id}", response_model=VideoJobDetailResponse
)
async def get_video_job(
    job_id: str,
    db: AsyncSession = Depends(get_db),
) -> VideoJobDetailResponse:
    """Return full job state including result fields if completed."""
    try:
        uid = uuid.UUID(job_id)
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Invalid job ID format",
        )

    result = await db.execute(
        select(VideoProcessingJob).where(
            VideoProcessingJob.id == uid
        )
    )
    job = result.scalar_one_or_none()

    if job is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Job not found",
        )

    return VideoJobDetailResponse.model_validate(job)
