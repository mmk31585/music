"""
Celery tasks for audio embedding extraction.

This module follows the exact structure of ``lyrics_tasks.py``:

- A ``worker_process_init`` signal that pre-loads the OpenL3 model into
  the process-wide singleton at worker startup (not at first-task time).
- The ``extract_embedding_task`` task that orchestrates: resolve audio →
  extract embedding → upsert into ``track_embeddings`` → complete.

Design
------
- The task takes only a ``job_id`` (not the full payload). The DB row is
  the single source of truth. This keeps Celery messages small and makes
  the idempotency guarantee reliable.
- Dependencies (DB session, extractor, storage client) are obtained from
  factories, keeping the task body testable via mocking.
- Extraction failures are NOT auto-retried (a corrupt audio file will
  fail the same way every retry). Only download failures (transient
  network issues) use ``self.retry()``.
- Storage client is reused via ``get_storage_client()`` — the same
  abstraction used by ``lyrics_tasks.py`` (not a duplicate).
"""

from __future__ import annotations

import logging
import os

from celery import Task
from celery.signals import worker_process_init
from sqlalchemy import text

from app.config import settings
from app.domain.audio_features.embedding_extractor import (
    AudioEmbeddingError,
    get_embedding_extractor,
)
from app.domain.audio_features.service import (
    AudioFeatureExtractionError,
    AudioFeatureExtractionService,
)
from app.infra.db.models import EmbeddingJob
from app.infra.db.sync_session import get_sync_db
from app.infra.storage.factory import get_storage_client

logger = logging.getLogger(__name__)


# ── Warm-up signal ──────────────────────────────────────────────────────


@worker_process_init.connect
def warm_up_embedding_model(**kwargs) -> None:  # noqa: ARG001
    """Load the OpenL3 model at worker startup.

    This ensures:
    1. The first real job doesn't pay the multi-second model-load penalty.
    2. Startup confirms the model can load on this machine, rather than
       failing silently at first-task-time.
    """
    logger.info("Warming up OpenL3 model via worker_process_init signal...")
    try:
        get_embedding_extractor(settings)
        logger.info("OpenL3 model warm-up complete")
    except Exception:
        logger.exception(
            "OpenL3 model failed to load at worker startup — "
            "tasks will fail with AudioEmbeddingError"
        )


# ── Celery app import ──────────────────────────────────────────────────


from app.workers.celery_app import celery_app as _celery_app  # noqa: E402

# ── Task ────────────────────────────────────────────────────────────────


@_celery_app.task(
    bind=True,
    max_retries=2,
    default_retry_delay=60,
    name="app.workers.tasks.embedding_tasks.extract_embedding_task",
    acks_late=True,
    reject_on_worker_lost=True,
)
def extract_embedding_task(self: Task, job_id: str) -> None:
    """Execute the full embedding-extraction pipeline for a given job.

    Flow::

        load job → resolve audio → extract embedding
        → upsert into track_embeddings → mark completed

    Args:
        job_id: The UUID of the ``EmbeddingJob`` row to process.
    """
    # Track temp file path so the finally block can clean up.  Will be
    # set only in ``presigned_url`` mode (shared_volume files are owned
    # by Go and must NEVER be deleted by this service).
    local_temp_path: str | None = None

    with get_sync_db() as session:
        # ── Step 1: Load job ────────────────────────────────────────
        job = session.get(EmbeddingJob, job_id)
        if job is None:
            logger.error("Embedding job %s not found in database — giving up", job_id)
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

            # ── Step 3: Extract embedding ───────────────────────────
            _set_status(session, job, "extracting")
            session.commit()

            extractor = get_embedding_extractor(settings)
            service = AudioFeatureExtractionService(extractor)
            result = service.extract_all(audio_path=audio_path)

            # ── Step 4: Upsert embedding into track_embeddings ──────
            _set_status(session, job, "completed")

            # Use raw SQL for the INSERT ... ON CONFLICT DO UPDATE pattern
            # because SQLAlchemy's ORM upsert is not straightforward with
            # Vector columns and ARRAY fields.
            session.execute(
                text("""
                    INSERT INTO track_embeddings_audio
                        (track_id, embedding, embedding_model, tempo_bpm,
                         genre_ids, release_year, extracted_at, updated_at)
                    VALUES
                        (:track_id, :embedding, :embedding_model, :tempo_bpm,
                         :genre_ids, :release_year, now(), now())
                    ON CONFLICT (track_id) DO UPDATE SET
                        embedding = EXCLUDED.embedding,
                        embedding_model = EXCLUDED.embedding_model,
                        tempo_bpm = EXCLUDED.tempo_bpm,
                        genre_ids = EXCLUDED.genre_ids,
                        release_year = EXCLUDED.release_year,
                        updated_at = now()
                """),
                {
                    "track_id": job.track_id,
                    "embedding": result.embedding,
                    "embedding_model": result.embedding_model,
                    "tempo_bpm": result.tempo_bpm,
                    "genre_ids": job.genre_ids if job.genre_ids else [],
                    "release_year": job.release_year,
                },
            )

            # Mark the job as completed
            _set_status(session, job, "completed")
            session.commit()

            logger.info(
                "Embedding job %s completed: track=%s model=%s dim=%d tempo=%s",
                job_id,
                job.track_id,
                result.embedding_model,
                result.embedding_dim,
                f"{result.tempo_bpm:.1f}" if result.tempo_bpm is not None else "None",
            )

        # ── Error handling ──────────────────────────────────────────

        except (
            FileNotFoundError,
            ValueError,
        ) as exc:
            # Permanent storage errors — bad path, path traversal, etc.
            logger.warning(
                "Embedding job %s failed (permanent storage error): %s", job_id, exc
            )
            _fail_job(session, job, str(exc))

        except (AudioEmbeddingError, AudioFeatureExtractionError) as exc:
            # Extraction failures are NOT retried — a malformed audio
            # file will fail identically every time; retrying would waste
            # precious CPU on a constrained box.
            logger.warning(
                "Embedding job %s failed during extraction (not retried): %s",
                job_id,
                exc,
            )
            _fail_job(session, job, str(exc))

        except self.MaxRetriesExceededError:
            logger.error("Embedding job %s exhausted retries", job_id)
            _fail_job(session, job, "Exhausted retries")

        except Exception as exc:
            # Transient errors (network timeouts, etc.) — retry.
            logger.warning(
                "Embedding job %s failed (retry %d/%d): %s",
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


def _set_status(session, job: EmbeddingJob, new_status: str) -> None:
    """Transition the job to a new status."""
    job.status = new_status
    session.add(job)


def _fail_job(session, job: EmbeddingJob, error_message: str) -> None:
    """Transition the job to ``failed`` without retrying."""
    job.status = "failed"
    job.error_message = error_message
    session.add(job)
    session.commit()


def _bump_retry(session, job: EmbeddingJob, error_message: str) -> None:
    """Increment retry count and store the latest error."""
    job.retry_count = (job.retry_count or 0) + 1
    job.error_message = error_message
    session.add(job)
