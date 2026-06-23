"""
Celery tasks for video processing (audio replacement).

Follows the exact structure of ``lyrics_tasks.py`` and ``embedding_tasks.py``:

- The ``process_video_task`` task orchestrates: resolve video + audio →
  audio replacement → extract thumbnail → detect metadata → callback → complete.
- Video processing failures are typically NOT retried (a corrupt video file
  will fail identically every retry). Only transient network issues during
  download / callback use ``self.retry()``.

Storage mode decisions
----------------------
- ``shared_volume`` mode: raw_video_path and audio_file_path are local paths
  on a shared mount. The processed output goes to a subdirectory within the
  shared volume so Go can serve it directly. Raw files are NEVER deleted.
- ``presigned_url`` mode: both video and audio are downloaded via URL.
  Temp files are cleaned up in the ``finally`` block.

Storage of processed output
---------------------------
After ``AudioReplacer.process()``, the result video and thumbnail are written
to ``<audio_shared_path>/.video-output/<video_id>/`` (shared_volume mode).
This allows Go to pick them up from the same shared filesystem.

For ``presigned_url`` mode, the output paths are local temp files. The
callback delivers the local path, and Go is responsible for uploading to
final storage (S3) and moving to the correct bucket prefix. (Phase 3
implements the Go receiver.)
"""

from __future__ import annotations

import logging
import os

from celery import Task

from app.config import settings
from app.domain.video_processing.audio_replacer import (
    VideoProcessingError,
    get_audio_replacer,
)
from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
)
from app.infra.callback.video_client import GoCallbackVideoClient
from app.infra.db.models import VideoProcessingJob
from app.infra.db.sync_session import get_sync_db
from app.infra.storage.factory import get_storage_client

logger = logging.getLogger(__name__)


# ── Celery app import ──────────────────────────────────────────────────


from app.workers.celery_app import celery_app as _celery_app  # noqa: E402

# ── Task ────────────────────────────────────────────────────────────────


@_celery_app.task(
    bind=True,
    max_retries=1,  # video processing failures are usually not retriable
    default_retry_delay=120,
    name="app.workers.tasks.video_tasks.process_video_task",
    acks_late=True,
    reject_on_worker_lost=True,
)
def process_video_task(self: Task, job_id: str) -> None:
    """Execute the full video-processing pipeline for a given job.

    Flow::

        load job → download video + audio → audio replacement
        → thumbnail extraction → callback → mark completed

    Args:
        job_id: The UUID of the ``VideoProcessingJob`` row to process.
    """
    # Track temp file paths so the finally block can clean up. Will be
    # set only in ``presigned_url`` mode (shared_volume files are owned
    # by Go and must NEVER be deleted by this service).
    local_video_temp: str | None = None
    local_audio_temp: str | None = None
    # Output paths from processing — only cleaned up in presigned_url mode
    # after the callback is delivered (since shared_volume outputs are Go's
    # responsibility).
    local_output_paths: list[str] = []

    with get_sync_db() as session:
        # ── Step 1: Load job ────────────────────────────────────────
        job = session.get(VideoProcessingJob, job_id)
        if job is None:
            logger.error(
                "VideoProcessingJob %s not found in database — giving up",
                job_id,
            )
            return  # Don't retry — the job doesn't exist

        try:
            # ── Step 2: Resolve / download video and audio ──────────
            _set_status(session, job, "downloading")
            session.commit()

            storage_client = get_storage_client(settings)
            temp_dir = os.path.join(
                settings.audio_shared_path, ".ml-temp"
            )

            # Resolve video source
            video_source = (
                job.raw_video_path or job.raw_video_download_url
            )
            if not video_source:
                raise ValueError(
                    "Job has neither raw_video_path nor "
                    "raw_video_download_url"
                )

            local_video_temp = storage_client.fetch_to_local_path(
                video_source, dest_dir=temp_dir
            )
            video_path = local_video_temp

            # Resolve audio source
            audio_source = (
                job.audio_file_path or job.audio_download_url
            )
            if not audio_source:
                raise ValueError(
                    "Job has neither audio_file_path nor "
                    "audio_download_url"
                )

            local_audio_temp = storage_client.fetch_to_local_path(
                audio_source, dest_dir=temp_dir
            )
            audio_path = local_audio_temp

            # ── Step 3: Audio replacement ───────────────────────────
            _set_status(session, job, "processing")
            session.commit()

            # Determine output directory based on storage mode
            if settings.storage_mode == "shared_volume":
                # Write processed files to a known location within the
                # shared volume so Go can serve them.
                output_dir = os.path.join(
                    settings.audio_shared_path,
                    ".video-output",
                    job.video_id,
                )
            else:
                # presigned_url mode: write to temp dir
                output_dir = os.path.join(
                    settings.audio_shared_path, ".video-output"
                )

            replacer = get_audio_replacer()
            result = replacer.process(
                video_path=video_path,
                audio_path=audio_path,
                track_start_ms=job.track_start_ms,
                output_dir=output_dir,
            )

            # Track output paths for cleanup (presigned_url mode only)
            local_output_paths = [
                result.output_video_path,
                result.thumbnail_path,
            ]

            # Store results in the job
            job.result_video_path = result.output_video_path
            job.result_thumbnail_path = result.thumbnail_path
            job.result_duration_ms = result.duration_ms
            job.result_aspect_ratio = result.aspect_ratio

            # ── Step 4: Callback to Go backend ──────────────────────
            _set_status(session, job, "uploading")
            session.commit()

            callback_client = _get_callback_client()
            delivered = callback_client.deliver_video_result(
                video_id=job.video_id,
                final_video_path=result.output_video_path,
                thumbnail_path=result.thumbnail_path,
                duration_ms=result.duration_ms,
                aspect_ratio=result.aspect_ratio,
                processing_status="completed",
            )

            if not delivered:
                raise CallbackTransientError(
                    "Video callback returned non-success status"
                )

            # ── Step 5: Mark completed ──────────────────────────────
            _set_status(session, job, "completed")
            job.callback_delivered = True
            session.commit()

            logger.info(
                "VideoProcessingJob %s completed: video=%s, "
                "duration=%dms, aspect=%s",
                job_id,
                job.video_id,
                result.duration_ms,
                result.aspect_ratio,
            )

        # ── Error handling ──────────────────────────────────────────

        except (
            FileNotFoundError,
            ValueError,
        ) as exc:
            # Permanent storage errors — bad path, path traversal, etc.
            logger.warning(
                "VideoProcessingJob %s failed "
                "(permanent storage error): %s",
                job_id,
                exc,
            )
            _fail_job(session, job, str(exc))
            _send_failed_callback(session, job, str(exc))

        except VideoProcessingError as exc:
            # Processing failures are NOT retried — a corrupt video or
            # audio file will fail identically every time.
            logger.warning(
                "VideoProcessingJob %s failed during processing "
                "(not retried): %s",
                job_id,
                exc,
            )
            _fail_job(session, job, str(exc))
            _send_failed_callback(session, job, str(exc))

        except CallbackRejectedError as exc:
            # 4xx from Go — permanent rejection (bad signature, unknown
            # video_id, etc.). Do NOT retry.
            logger.warning(
                "VideoProcessingJob %s callback permanently "
                "rejected: %s",
                job_id,
                exc,
            )
            _fail_job(session, job, str(exc))

        except self.MaxRetriesExceededError:
            logger.error(
                "VideoProcessingJob %s exhausted retries", job_id
            )
            _fail_job(
                session, job, "Exhausted retries"
            )
            _send_failed_callback(session, job, "Exhausted retries")

        except Exception as exc:
            # Transient errors (network timeouts, 5xx from Go, etc.) —
            # retry.
            logger.warning(
                "VideoProcessingJob %s failed "
                "(retry %d/%d): %s",
                job_id,
                self.request.retries,
                self.max_retries,
                exc,
            )
            _bump_retry(session, job, str(exc))
            session.commit()
            raise self.retry(exc=exc) from exc

        else:
            if job.status != "completed":
                _set_status(session, job, "completed")
            session.commit()

        finally:
            # ── Cleanup: temp file removal ───────────────────────────
            # Never delete files in shared_volume mode — those are owned
            # by the Go backend.  Only clean up presigned_url downloads.
            if settings.storage_mode == "presigned_url":
                to_clean = [
                    local_video_temp,
                    local_audio_temp,
                ] + local_output_paths

                for path in to_clean:
                    if path is not None:
                        try:
                            os.remove(path)
                            logger.debug(
                                "Cleaned up temp file %s", path
                            )
                        except OSError:
                            logger.warning(
                                "Failed to remove temp %s", path
                            )


# ── Helpers ────────────────────────────────────────────────────────────


def _set_status(
    session, job: VideoProcessingJob, new_status: str
) -> None:
    """Transition the job to a new status."""
    job.status = new_status
    session.add(job)


def _fail_job(
    session, job: VideoProcessingJob, error_message: str
) -> None:
    """Transition the job to ``failed`` without retrying."""
    job.status = "failed"
    job.error_message = error_message
    session.add(job)
    session.commit()


def _bump_retry(
    session, job: VideoProcessingJob, error_message: str
) -> None:
    """Increment retry count and store the latest error."""
    job.retry_count = (job.retry_count or 0) + 1
    job.error_message = error_message
    session.add(job)


# ── Callback client singleton ──────────────────────────────────────────


_callback_client: GoCallbackVideoClient | None = None


def _get_callback_client() -> GoCallbackVideoClient:
    """Return the singleton callback client (lazy init)."""
    global _callback_client  # noqa: PLW0603
    if _callback_client is None:
        _callback_client = GoCallbackVideoClient(
            callback_url=settings.go_backend_callback_url,
            hmac_secret=settings.webhook_hmac_secret,
        )
    return _callback_client


def _send_failed_callback(
    session, job: VideoProcessingJob, error_message: str
) -> None:
    """Send a failed-status callback to Go if the job has results."""
    try:
        client = _get_callback_client()
        client.deliver_video_result(
            video_id=job.video_id,
            final_video_path=job.result_video_path or "",
            thumbnail_path=job.result_thumbnail_path or "",
            duration_ms=job.result_duration_ms or 0,
            aspect_ratio=job.result_aspect_ratio or "",
            processing_status="failed",
            error_message=error_message,
        )
    except Exception:
        logger.exception(
            "Failed to send failure callback for job %s", job.id
        )
