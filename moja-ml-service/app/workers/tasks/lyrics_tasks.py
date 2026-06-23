"""
Celery tasks for lyrics generation.

This module wires together:

- A ``worker_process_init`` signal that pre-loads the Whisper model into
  the process-wide singleton at worker startup (not at first-task time).
- The ``generate_lyrics_task`` task that orchestrates: resolve audio →
  transcribe → format → callback → complete.

Design
------
- The task takes only a ``job_id`` (not the full payload). The DB row is
  the single source of truth for all job details. This keeps Celery
  messages small and makes the idempotency guarantee reliable.
- Dependencies (DB session, transcriber, storage client, callback client)
  are obtained from factories, keeping the task body testable via mocking.
- Transcription failures are NOT auto-retried (a broken audio file will
  break the same way every time). Only download/callback failures
  (transient network issues) use ``self.retry()``.
- Storage and callback clients are **synchronous** -- Celery tasks run
  synchronously and there's no benefit to async in a concurrency=1 worker.
"""

from __future__ import annotations

import logging
import os

from celery import Task
from celery.signals import worker_process_init

from app.config import settings
from app.domain.lyrics.service import LyricsGenerationService
from app.domain.lyrics.transcriber import TranscriptionError, get_transcriber
from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
    GoCallbackClient,
)
from app.infra.db.models import LyricsJob
from app.infra.db.sync_session import get_sync_db
from app.infra.storage.factory import get_storage_client

logger = logging.getLogger(__name__)


# ── Warm-up signal ──────────────────────────────────────────────────────


@worker_process_init.connect
def warm_up_whisper_model(**kwargs) -> None:  # noqa: ARG001
    """Load the Whisper model into the process-wide singleton at worker startup.

    This ensures:
    1. The first real job doesn't pay the multi-second model-load penalty.
    2. Startup confirms (loudly, via log) that the model can actually load
       on this machine, rather than failing silently at first-task-time.
    """
    logger.info("Warming up Whisper model via worker_process_init signal...")
    try:
        get_transcriber(settings)
        logger.info("Whisper model warm-up complete")
    except Exception:
        logger.exception(
            "Whisper model failed to load at worker startup — "
            "tasks will fail with TranscriptionError"
        )


# ── Celery app import ──────────────────────────────────────────────────


from app.workers.celery_app import celery_app as _celery_app  # noqa: E402

# ── Task ────────────────────────────────────────────────────────────────


@_celery_app.task(
    bind=True,
    max_retries=2,
    default_retry_delay=60,
    name="app.workers.tasks.lyrics_tasks.generate_lyrics_task",
    acks_late=True,
    reject_on_worker_lost=True,
)
def generate_lyrics_task(self: Task, job_id: str) -> None:
    """Execute the full lyrics-generation pipeline for a given job.

    Flow::

        load job → resolve audio → transcribe → format
        → store results → callback → mark completed

    Args:
        job_id: The UUID of the ``LyricsJob`` row to process.
    """
    # Track temp file path so the finally block can clean up.  Will be
    # set only in ``presigned_url`` mode (shared_volume files are owned
    # by Go and must NEVER be deleted by this service).
    local_temp_path: str | None = None

    with get_sync_db() as session:
        # ── Step 1: Load job ────────────────────────────────────────
        job = session.get(LyricsJob, job_id)
        if job is None:
            logger.error("Job %s not found in database — giving up", job_id)
            return  # Don't retry — the job doesn't exist

        try:
            # ── Step 2: Resolve / download audio ─────────────────────
            _set_status(session, job, "downloading")
            session.commit()

            storage_client = get_storage_client(settings)
            # Determine the audio source (shared_volume path or presigned URL).
            source = job.audio_file_path or job.audio_download_url
            if not source:
                raise ValueError(
                    "Job has neither audio_file_path nor audio_download_url"
                )

            # Temporary directory for downloaded files (used only in
            # presigned_url mode; shared_volume mode ignores dest_dir).
            temp_dir = os.path.join(settings.audio_shared_path, ".ml-temp")
            local_temp_path = storage_client.fetch_to_local_path(
                source, dest_dir=temp_dir
            )
            audio_path = local_temp_path

            # ── Step 3: Transcribe ───────────────────────────────────
            _set_status(session, job, "transcribing")
            session.commit()

            transcriber = get_transcriber(settings)
            service = LyricsGenerationService(transcriber)
            result = service.generate(
                audio_path=audio_path,
                track_title=job.track_title,
                track_artist=job.track_artist,
            )

            # ── Step 4: Store results ────────────────────────────────
            _set_status(session, job, "formatting")
            job.result_lrc_content = result.lrc_content
            job.result_plain_text = result.plain_text
            job.result_confidence = result.confidence
            job.result_detected_language = result.detected_language
            job.result_segment_count = result.segment_count
            job.whisper_model_version = settings.whisper_model_size

            # ── Step 5: Callback to Go backend ───────────────────────
            _set_status(session, job, "callback_pending")
            session.commit()

            callback_client = GoCallbackClient(
                settings.go_backend_callback_url,
                settings.webhook_hmac_secret,
            )
            delivered = callback_client.deliver_lyrics_result(
                track_id=job.track_id,
                lrc_content=result.lrc_content,
                plain_text=result.plain_text,
                confidence=result.confidence,
                detected_language=result.detected_language,
                whisper_model_version=settings.whisper_model_size,
            )

            if not delivered:
                # Shouldn't normally happen (callback client raises on
                # failure), but protect against unexpected contract changes.
                raise CallbackTransientError("Callback returned non-success status")

            # ── Step 6: Mark completed ───────────────────────────────
            _set_status(session, job, "completed")
            job.callback_delivered = True
            session.commit()

            logger.info(
                "Job %s completed: %d segments, language=%s, confidence=%.3f",
                job_id,
                result.segment_count,
                result.detected_language,
                result.confidence,
            )

        # ── Error handling ──────────────────────────────────────────

        except (
            FileNotFoundError,
            ValueError,
        ) as exc:
            # Permanent storage errors — bad path, path traversal, etc.
            logger.warning("Job %s failed (permanent storage error): %s", job_id, exc)
            _fail_job(session, job, str(exc))

        except TranscriptionError as exc:
            # Transcription failures are NOT retried — a malformed audio
            # file will fail identically every time; retrying would waste
            # precious CPU on a constrained box.
            logger.warning(
                "Job %s failed during transcription (not retried): %s", job_id, exc
            )
            _fail_job(session, job, str(exc))

        except CallbackRejectedError as exc:
            # 4xx from Go — permanent rejection (bad signature, unknown
            # track_id, etc.).  Do NOT retry.
            logger.warning("Job %s callback permanently rejected: %s", job_id, exc)
            _fail_job(session, job, str(exc))

        except self.MaxRetriesExceededError:
            logger.error("Job %s exhausted retries", job_id)
            _fail_job(session, job, "Exhausted retries")

        except Exception as exc:
            # Transient errors (network timeouts, 5xx from Go, etc.) — retry.
            logger.warning(
                "Job %s failed (retry %d/%d): %s",
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
            if local_temp_path is not None and settings.storage_mode == "presigned_url":
                try:
                    os.remove(local_temp_path)
                    logger.debug("Cleaned up temp file %s", local_temp_path)
                except OSError:
                    logger.warning("Failed to remove temp %s", local_temp_path)


# ── Helpers ────────────────────────────────────────────────────────────


def _set_status(session, job: LyricsJob, new_status: str) -> None:
    """Transition the job to a new status."""
    job.status = new_status
    session.add(job)


def _fail_job(session, job: LyricsJob, error_message: str) -> None:
    """Transition the job to ``failed`` without retrying."""
    job.status = "failed"
    job.error_message = error_message
    session.add(job)
    session.commit()


def _bump_retry(session, job: LyricsJob, error_message: str) -> None:
    """Increment retry count and store the latest error."""
    job.retry_count = (job.retry_count or 0) + 1
    job.error_message = error_message
    session.add(job)
