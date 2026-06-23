"""
Celery tasks for cover image optimization.

Flow
----
1. Receive ``job_id`` → load ``CoverJob`` from DB.
2. Resolve cover source image (URL download or shared volume path).
3. Optimize with Pillow → WebP variants (original, thumb, med).
4. Store optimized files to shared storage.
5. Call back to Go backend with new URLs.
6. Mark job ``completed``.
"""

from __future__ import annotations

import logging
import os
import uuid
from pathlib import Path

import httpx
from celery import Task
from sqlalchemy import create_engine
from sqlalchemy.orm import Session as SASession

from app.config import settings
from app.domain.covers.optimizer import CoverOptimizationError, CoverOptimizer
from app.infra.callback.cover_client import GoCallbackCoverClient
from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
)
from app.infra.db.models import CoverJob
from app.workers.celery_app import celery_app

logger = logging.getLogger(__name__)

# ── Sync DB session (Celery workers don't use async) ───────────────────

_engine: object | None = None


def _get_sync_session() -> SASession:
    global _engine  # noqa: PLW0603
    if _engine is None:
        _engine = create_engine(settings.database_url.replace("+asyncpg", ""))
    return SASession(bind=_engine)


# ── Callback client ────────────────────────────────────────────────────

_callback_client: GoCallbackCoverClient | None = None


def _get_callback_client() -> GoCallbackCoverClient:
    global _callback_client  # noqa: PLW0603
    if _callback_client is None:
        _callback_client = GoCallbackCoverClient(
            callback_url=settings.go_backend_callback_url,
            hmac_secret=settings.webhook_hmac_secret,
        )
    return _callback_client


# ── Optimizer singleton ────────────────────────────────────────────────

_optimizer: CoverOptimizer | None = None


def _get_optimizer() -> CoverOptimizer:
    global _optimizer  # noqa: PLW0603
    if _optimizer is None:
        _optimizer = CoverOptimizer()
    return _optimizer


# ── Helpers ────────────────────────────────────────────────────────────


def _resolve_cover_bytes(cover_url: str | None, cover_path: str | None) -> bytes:
    """Download or read the cover image, return raw bytes.

    Priority: URL download > shared volume path.
    """
    if cover_url:
        logger.info("Downloading cover from URL: %s", cover_url[:80])
        try:
            resp = httpx.get(cover_url, timeout=60.0, follow_redirects=True)
            resp.raise_for_status()
            return resp.content
        except Exception as exc:
            raise CoverOptimizationError(
                f"Failed to download cover from {cover_url}: {exc}"
            ) from exc

    if cover_path:
        full_path = cover_path
        if not os.path.isabs(cover_path):
            full_path = os.path.join(settings.audio_shared_path, cover_path)

        logger.info("Reading cover from path: %s", full_path)
        try:
            with open(full_path, "rb") as f:
                return f.read()
        except Exception as exc:
            raise CoverOptimizationError(
                f"Failed to read cover from {full_path}: {exc}"
            ) from exc

    raise CoverOptimizationError("No cover_url or cover_path provided")


def _resolve_storage_path(album_id: str) -> str:
    """Determine where optimized covers should be written.

    For ``shared_volume`` mode, this is the catalog-covers directory
    relative to the Go backend's storage base.

    Returns the base path (without extension) for the album's cover.
    """
    # In shared_volume mode, the Go backend's uploads directory should
    # be accessible. The catalog-covers structure is:
    #   catalog-covers/{album_id}/cover.webp
    #   catalog-covers/{album_id}/cover_thumb.webp
    #   catalog-covers/{album_id}/cover_med.webp
    if settings.storage_mode == "shared_volume":
        # Assume we can write directly to the shared volume.
        # The Go backend's storage base is typically at:
        #   {audio_shared_path}/../catalog-covers/
        # But for simplicity, write relative to the storage base.
        # This path must match what Go's storage.GetURL produces.
        base_dir = settings.audio_shared_path
        # Go up one level to reach the uploads base dir
        parent = str(Path(base_dir).parent)
        cover_dir = os.path.join(parent, "catalog-covers", album_id)
        os.makedirs(cover_dir, exist_ok=True)
        return os.path.join(cover_dir, "cover")

    # For presigned_url mode, return the key prefix and let the
    # callback handler deal with S3 upload.
    return f"catalog-covers/{album_id}/cover"


def _write_webp_files(
    base_path: str,
    result: tuple[bytes, int, int] | None,
    suffix: str,
) -> str | None:
    """Write a WebP variant to disk, return its file path.

    In production with S3, this would upload to the bucket instead.
    """
    if result is None:
        return None

    data, w, h = result
    filepath = f"{base_path}{suffix}.webp"

    # Only write for shared_volume mode
    if settings.storage_mode == "shared_volume":
        with open(filepath, "wb") as f:
            f.write(data)
        logger.info("Wrote %s (%dx%d, %d bytes)", filepath, w, h, len(data))

    return filepath


# ── Task ───────────────────────────────────────────────────────────────


@celery_app.task(
    bind=True,
    max_retries=2,
    default_retry_delay=60,
    acks_late=True,
)
def optimize_cover_task(self: Task, job_id: str) -> None:
    """Optimize a cover image: download, process, store, callback.

    This task is the cover equivalent of ``generate_lyrics_task``.
    """
    session = _get_sync_session()
    try:
        # 1. Load job
        uid = uuid.UUID(job_id)
        job = session.get(CoverJob, uid)
        if job is None:
            logger.error("CoverJob %s not found", job_id)
            return

        # 2. Set status to downloading
        job.status = "downloading"
        session.commit()

        # 3. Resolve cover bytes
        cover_data = _resolve_cover_bytes(job.cover_url, job.cover_path)
        original_size = len(cover_data)

        # 4. Optimize
        job.status = "optimizing"
        session.commit()

        optimizer = _get_optimizer()
        opt_result = optimizer.optimize(cover_data)

        # 5. Store optimized files
        job.status = "storing"
        session.commit()

        base_path = _resolve_storage_path(job.album_id)
        original_path = _write_webp_files(base_path, opt_result.original, "")
        thumb_path = _write_webp_files(base_path, opt_result.thumb, "_thumb")
        med_path = _write_webp_files(base_path, opt_result.med, "_med")

        # Track savings
        webp_size = len(opt_result.original[0]) if opt_result.original else 0
        bytes_saved = original_size - webp_size if webp_size > 0 else 0

        # Save results to job
        job.result_cover_url = original_path
        job.result_thumb_url = thumb_path
        job.result_med_url = med_path
        if opt_result.original:
            job.original_width = opt_result.original[1]
            job.original_height = opt_result.original[2]
        job.webp_bytes_saved = bytes_saved
        session.commit()

        # 6. Callback to Go
        job.status = "callback_pending"
        session.commit()

        client = _get_callback_client()
        client.deliver_cover_result(
            album_id=job.album_id,
            cover_url=original_path or "",
            thumb_url=thumb_path or "",
            med_url=med_path or "",
            track_id=job.track_id,
        )

        # 7. Mark completed
        job.status = "completed"
        job.callback_delivered = True
        session.commit()

        logger.info(
            "Cover optimization complete for album=%s (original=%d → webp=%d, saved=%d bytes)",
            job.album_id,
            original_size,
            webp_size,
            bytes_saved,
        )

    except (CallbackRejectedError, CoverOptimizationError) as exc:
        logger.error("Permanent failure for job %s: %s", job_id, exc)
        if job is not None:
            job.status = "failed"
            job.error_message = str(exc)
            session.commit()

    except CallbackTransientError as exc:
        logger.warning("Transient failure for job %s, retrying: %s", job_id, exc)
        if job is not None:
            job.retry_count = (job.retry_count or 0) + 1
            session.commit()
        raise self.retry(exc=exc)

    except Exception as exc:
        logger.error("Unexpected error for job %s: %s", job_id, exc, exc_info=True)
        if job is not None:
            job.status = "failed"
            job.error_message = f"Unexpected error: {exc}"
            session.commit()
        raise

    finally:
        session.close()


@celery_app.task(acks_late=True)
def refresh_stale_covers(max_age_days: int = 7) -> int:
    """Periodic task: find albums whose covers haven't been optimized
    recently and enqueue optimization jobs for them.

    This is called by Celery Beat on a schedule. In production, this
    should query the Go backend's catalog to find albums with covers
    older than ``max_age_days``. For now, it logs a heartbeat.

    Returns the number of covers refreshed (0 = heartbeat only).
    """
    logger.info(
        "Cover refresh heartbeat: max_age_days=%d (queries Go catalog in production)",
        max_age_days,
    )
    # TODO: In production, call Go backend API to get stale albums:
    #   GET /api/v1/admin/covers/stale?days=7
    # Then enqueue optimize_cover_task for each.
    return 0
