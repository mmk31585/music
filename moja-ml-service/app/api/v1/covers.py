"""
API endpoints for cover image optimization jobs.

- ``POST /api/v1/covers/jobs`` — enqueue a cover optimization job
  (idempotent: same ``album_id`` returns existing job).
- ``GET  /api/v1/covers/jobs/{job_id}`` — query job status and results.
"""

from __future__ import annotations

import uuid
from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, field_validator, model_validator
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.infra.db.models import CoverJob
from app.infra.db.session import get_db
from app.workers.celery_app import celery_app

router = APIRouter(prefix="/covers", tags=["covers"])


# ── Request / response schemas ─────────────────────────────────────────


class CreateCoverJobRequest(BaseModel):
    """Payload for enqueuing a new cover optimization job."""

    album_id: str
    cover_url: str | None = None
    cover_path: str | None = None
    track_id: str | None = None  # optional — also update track cover


class CoverJobResponse(BaseModel):
    """Summary response returned when creating a job."""

    job_id: str
    album_id: str
    status: str
    created_at: datetime | None = None

    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, CoverJob):
            data.job_id = str(data.id)  # type: ignore[attr-defined]
        return data


class CoverJobDetailResponse(BaseModel):
    """Full job detail including optimized URLs if completed."""

    job_id: str
    album_id: str
    status: str
    cover_url: str | None = None
    cover_path: str | None = None
    track_id: str | None = None
    result_cover_url: str | None = None
    result_thumb_url: str | None = None
    result_med_url: str | None = None
    original_format: str | None = None
    original_width: int | None = None
    original_height: int | None = None
    webp_bytes_saved: int | None = None
    error_message: str | None = None
    retry_count: int = 0
    callback_delivered: bool = False
    created_at: datetime | None = None
    updated_at: datetime | None = None

    model_config = {"from_attributes": True}

    @model_validator(mode="before")
    @classmethod
    def _map_id_to_job_id(cls, data: object) -> object:
        if isinstance(data, CoverJob):
            data.job_id = str(data.id)  # type: ignore[attr-defined]
        return data

    @field_validator("retry_count", "callback_delivered", mode="before")
    @classmethod
    def _none_to_default(cls, v: object) -> object:
        return 0 if v is None else v


# ── Endpoints ──────────────────────────────────────────────────────────


@router.post("/jobs", response_model=CoverJobResponse, status_code=status.HTTP_202_ACCEPTED)
async def create_cover_job(
    req: CreateCoverJobRequest,
    db: AsyncSession = Depends(get_db),
) -> CoverJobResponse:
    """Enqueue a cover optimization job.

    **Idempotent:** if a job for this ``album_id`` already exists and is
    not in a terminal ``failed`` state, the existing job is returned
    instead of creating a duplicate.
    """
    idempotency_key = f"cover:{req.album_id}"

    # Check for existing non-failed job
    result = await db.execute(
        select(CoverJob).where(
            CoverJob.idempotency_key == idempotency_key,
            CoverJob.status != "failed",
        )
    )
    existing_job = result.scalar_one_or_none()

    if existing_job is not None:
        return CoverJobResponse.model_validate(existing_job)

    # Create new job
    job = CoverJob(
        id=uuid.uuid4(),
        album_id=req.album_id,
        idempotency_key=idempotency_key,
        status="queued",
        cover_url=req.cover_url,
        cover_path=req.cover_path,
        track_id=req.track_id,
    )
    db.add(job)

    await db.commit()
    await db.refresh(job)

    # Enqueue Celery task
    celery_app.send_task(
        "app.workers.tasks.cover_tasks.optimize_cover_task",
        args=[str(job.id)],
    )

    return CoverJobResponse.model_validate(job)


@router.get("/jobs/{job_id}", response_model=CoverJobDetailResponse)
async def get_cover_job(
    job_id: str,
    db: AsyncSession = Depends(get_db),
) -> CoverJobDetailResponse:
    """Return full job state including optimized URLs if completed."""
    try:
        uid = uuid.UUID(job_id)
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Invalid job ID format",
        )

    result = await db.execute(select(CoverJob).where(CoverJob.id == uid))
    job = result.scalar_one_or_none()

    if job is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Job not found",
        )

    return CoverJobDetailResponse.model_validate(job)
