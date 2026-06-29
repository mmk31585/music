"""
Celery tasks for audio embedding extraction (refactored).

Uses BaseTaskOrchestrator for consistent error handling, status transitions,
and cleanup. Keeps model warm-up signal and embedding-specific logic.
"""

from __future__ import annotations

import logging
import os
from typing import TYPE_CHECKING

from celery import Task
from celery.signals import worker_process_init
from sqlalchemy import text
from sqlalchemy.orm import Session as SASession

from app.config import settings
from app.domain.audio_features.embedding_extractor import (
    get_embedding_extractor,
)
from app.domain.audio_features.service import (
    AudioFeatureExtractionService,
)
from app.infra.db.models import EmbeddingJob
from app.workers.tasks.base import BaseTaskOrchestrator, StorageResolutionMixin

if TYPE_CHECKING:
    from app.domain.audio_features.service import AudioFeatureExtractionResult

logger = logging.getLogger(__name__)

# ── Celery app import ──────────────────────────────────────────────────

from app.workers.celery_app import celery_app as _celery_app  # noqa: E402

# ── Warm-up signal ──────────────────────────────────────────────────────


@worker_process_init.connect
def warm_up_embedding_model(**kwargs) -> None:  # noqa: ARG001
    """Load the OpenL3 model at worker startup."""
    logger.info("Warming up OpenL3 model via worker_process_init signal...")
    try:
        get_embedding_extractor(settings)
        logger.info("OpenL3 model warm-up complete")
    except Exception:
        logger.exception(
            "OpenL3 model failed to load at worker startup — "
            "tasks will fail with AudioEmbeddingError"
        )

# ── Orchestrator ────────────────────────────────────────────────────────


class EmbeddingTaskOrchestrator(BaseTaskOrchestrator[EmbeddingJob], StorageResolutionMixin):
    """Orchestrator for embedding extraction jobs."""

    def _get_job_model(self) -> type[EmbeddingJob]:
        return EmbeddingJob

    def _process_job(self, session: SASession, job: EmbeddingJob) -> None:
        """Execute embedding extraction pipeline."""
        # ── Step 1: Resolve audio file ──────────────────────────────
        self._set_status(session, job, "downloading")
        session.commit()

        temp_dir = os.path.join(settings.audio_shared_path, "temp-downloads")
        os.makedirs(temp_dir, exist_ok=True)

        local_path = self.resolve_file(
            download_url=job.audio_download_url,
            file_path=job.audio_file_path,
            temp_dir=temp_dir,
            settings=settings,
        )

        # ── Step 2: Extract features ─────────────────────────────────
        self._set_status(session, job, "processing")
        session.commit()

        extractor = get_embedding_extractor(settings)
        service = AudioFeatureExtractionService(extractor=extractor)
        result: AudioFeatureExtractionResult = service.extract_all(local_path)

        logger.info(
            "Extracted embedding for job %s: dim=%d tempo=%s",
            self.job_id,
            result.embedding_dim,
            f"{result.tempo_bpm:.1f}" if result.tempo_bpm else "None",
        )

        # ── Step 3: Upsert into track_embeddings_audio ───────────────
        self._set_status(session, job, "upserting")
        session.commit()

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

        self._set_status(session, job, "completed")
        session.commit()

        logger.info(
            "Embedding job %s completed: track=%s model=%s",
            self.job_id,
            job.track_id,
            result.embedding_model,
        )

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
    """Extract audio embedding for a track."""
    orchestrator = EmbeddingTaskOrchestrator(celery_task=self, job_id=job_id)
    orchestrator.execute()
