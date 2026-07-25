"""
Base task orchestrator for all Celery tasks.

Provides common infrastructure:
- Database session management
- Job status transitions
- Error handling patterns
- Retry logic
- Cleanup hooks
"""

from __future__ import annotations

import logging
from typing import TYPE_CHECKING, Any, Callable, Generic, TypeVar

from celery import Task
from sqlalchemy.orm import Session as SASession

from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
)
from app.infra.db.sync_session import get_sync_db

if TYPE_CHECKING:
    from app.infra.db.models import CoverJob, EmbeddingJob, LyricsJob, VideoProcessingJob

logger = logging.getLogger(__name__)

# Type variable for job models
JobT = TypeVar("JobT", bound="CoverJob | EmbeddingJob | LyricsJob | VideoProcessingJob")


class BaseTaskOrchestrator(Generic[JobT]):
    """Base orchestrator for processing jobs with consistent patterns.
    
    Handles:
    - Loading job from DB
    - Status transitions
    - Error classification (permanent vs transient)
    - Retry logic
    - Cleanup
    
    Subclasses implement:
    - `_process_job()` — the actual work
    - `_get_job_model()` — return the SQLAlchemy model class
    """

    def __init__(self, celery_task: Task, job_id: str):
        self.task = celery_task
        self.job_id = job_id
        self.cleanup_hooks: list[Callable[[], None]] = []

    def execute(self) -> None:
        """Main execution flow with error handling and cleanup."""
        with get_sync_db() as session:
            job = self._load_job(session)
            if job is None:
                logger.error(
                    "%s job %s not found — giving up",
                    self._get_job_type(),
                    self.job_id,
                )
                return

            try:
                self._process_job(session, job)
                
                if job.status != "completed":
                    self._set_status(session, job, "completed")
                session.commit()

            except CallbackRejectedError as exc:
                # 4xx from Go — permanent rejection
                logger.warning(
                    "%s job %s callback rejected (permanent): %s",
                    self._get_job_type(),
                    self.job_id,
                    exc,
                )
                self._fail_job(session, job, str(exc))

            except (CallbackTransientError, OSError, ConnectionError) as exc:
                # Transient errors — retry
                logger.warning(
                    "%s job %s failed (retry %d/%d): %s",
                    self._get_job_type(),
                    self.job_id,
                    self.task.request.retries,
                    self.task.max_retries,
                    exc,
                )
                self._bump_retry(session, job, str(exc))
                session.commit()
                raise self.task.retry(exc=exc) from exc

            except Exception as exc:
                # Processing failures — permanent (don't retry corrupt data)
                logger.exception(
                    "%s job %s failed (not retried): %s",
                    self._get_job_type(),
                    self.job_id,
                    exc,
                )
                self._fail_job(session, job, str(exc))
                self._send_failed_callback(session, job, str(exc))

            finally:
                self._cleanup()

    def _load_job(self, session: SASession) -> JobT | None:
        """Load the job from the database."""
        model_class = self._get_job_model()
        return session.get(model_class, self.job_id)

    def _set_status(self, session: SASession, job: JobT, new_status: str) -> None:
        """Transition job to a new status."""
        job.status = new_status  # type: ignore
        session.add(job)

    def _fail_job(self, session: SASession, job: JobT, error_message: str) -> None:
        """Mark job as failed without retrying."""
        job.status = "failed"  # type: ignore
        job.error_message = error_message  # type: ignore
        session.add(job)
        session.commit()

    def _bump_retry(self, session: SASession, job: JobT, error_message: str) -> None:
        """Increment retry count and store error."""
        job.retry_count = (job.retry_count or 0) + 1  # type: ignore
        job.error_message = error_message  # type: ignore
        session.add(job)

    def register_cleanup(self, cleanup_fn: Callable[[], None]) -> None:
        """Register a cleanup function to run at the end."""
        self.cleanup_hooks.append(cleanup_fn)

    def _cleanup(self) -> None:
        """Execute all registered cleanup hooks."""
        for hook in self.cleanup_hooks:
            try:
                hook()
            except Exception:
                logger.exception("Cleanup hook failed")

    # ── Abstract methods (subclasses must implement) ────────────────────

    def _process_job(self, session: SASession, job: JobT) -> None:
        """Execute the actual job processing logic.
        
        Raises:
            Any exception to trigger retry or failure handling.
        """
        raise NotImplementedError

    def _get_job_model(self) -> type[JobT]:
        """Return the SQLAlchemy model class for this job type."""
        raise NotImplementedError

    def _get_job_type(self) -> str:
        """Return a human-readable job type name for logging."""
        return self._get_job_model().__name__

    def _send_failed_callback(self, session: SASession, job: JobT, error_message: str) -> None:
        """Optional: send failure callback to Go backend.
        
        Override in subclasses that need to notify Go of failures.
        """
        pass  # Default: no callback


class StorageResolutionMixin:
    """Mixin for tasks that need to resolve audio/video files from storage."""

    def resolve_file(
        self,
        download_url: str | None,
        file_path: str | None,
        temp_dir: str,
        settings: Any,
    ) -> str:
        """Resolve file location and return local path.
        
        Returns:
            Local filesystem path to the file.
            
        Side effects:
            - Downloads file if using presigned_url mode
            - Registers cleanup hook for temp files
        """
        from app.infra.storage.factory import get_storage_client
        
        storage_client = get_storage_client(settings)
        source = file_path or download_url
        if not source:
            raise ValueError("Job has neither file_path nor download_url")

        local_path = storage_client.resolve(source, dest_dir=temp_dir)
        
        # Register cleanup for downloaded temp files
        if settings.storage_mode == "presigned_url" and local_path:
            self.register_cleanup(lambda: self._cleanup_temp_file(local_path, settings))
        
        return local_path

    @staticmethod
    def _cleanup_temp_file(path: str, settings: Any) -> None:
        """Remove temp file if in presigned_url mode."""
        import os
        
        if settings.storage_mode == "presigned_url":
            try:
                os.remove(path)
                logger.debug("Cleaned up temp file %s", path)
            except OSError:
                logger.warning("Failed to remove temp file %s", path)
