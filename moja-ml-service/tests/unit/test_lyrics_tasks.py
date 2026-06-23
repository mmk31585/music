"""Tests for the Celery task helpers and worker signal.

The task itself (``generate_lyrics_task``) is tested indirectly via
the API endpoint test that verifies Celery is invoked. Full end-to-end
task tests would require a running Celery worker and are better suited
for integration tests (``pytest.mark.integration``).

What's tested here:
  - ``_set_status`` / ``_fail_job`` / ``_bump_retry`` helper functions
    with mocked DB session.
  - ``warm_up_whisper_model`` signal — verifies it calls
    ``get_transcriber`` and logs the result.
"""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from app.workers.tasks.lyrics_tasks import (
    _bump_retry,
    _fail_job,
    _set_status,
    warm_up_whisper_model,
)


@pytest.fixture
def mock_job() -> MagicMock:
    """A mock LyricsJob with the fields our helpers touch."""
    job = MagicMock()
    job.status = "queued"
    job.retry_count = 0
    job.error_message = None
    return job


@pytest.fixture
def mock_session() -> MagicMock:
    """A mock DB session that records ``add`` calls."""
    session = MagicMock()
    session.commit.return_value = None
    return session


# ── _set_status ─────────────────────────────────────────────────────────


class TestSetStatus:
    def test_sets_status_on_job(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "transcribing")
        assert mock_job.status == "transcribing"

    def test_calls_session_add(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "completed")
        mock_session.add.assert_called_once_with(mock_job)

    def test_does_not_call_commit(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "formatting")
        mock_session.commit.assert_not_called()


# ── _fail_job ───────────────────────────────────────────────────────────


class TestFailJob:
    def test_sets_failed_status(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "Something went wrong")
        assert mock_job.status == "failed"

    def test_stores_error_message(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "Network error")
        assert mock_job.error_message == "Network error"

    def test_calls_session_add(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "error")
        mock_session.add.assert_called_once_with(mock_job)

    def test_calls_session_commit(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "error")
        mock_session.commit.assert_called_once()


# ── _bump_retry ─────────────────────────────────────────────────────────


class TestBumpRetry:
    def test_increments_retry_count(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 1

    def test_increments_from_zero(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        mock_job.retry_count = 0
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 1

    def test_increments_existing_count(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        mock_job.retry_count = 2
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 3

    def test_stores_error_message(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "connection refused")
        assert mock_job.error_message == "connection refused"

    def test_calls_session_add(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "err")
        mock_session.add.assert_called_once_with(mock_job)

    def test_does_not_call_commit(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "err")
        mock_session.commit.assert_not_called()


# ── warm_up_whisper_model ────────────────────────────────────────────────


class TestWarmUpWhisperModel:
    def test_calls_get_transcriber(self, mocker) -> None:
        """The signal handler should load the Whisper model through get_transcriber."""
        mock_get = mocker.patch("app.workers.tasks.lyrics_tasks.get_transcriber")
        mock_settings = mocker.patch("app.workers.tasks.lyrics_tasks.settings")

        warm_up_whisper_model()

        mock_get.assert_called_once_with(mock_settings)

    def test_logs_on_success(self, mocker) -> None:
        """Successful warm-up should log completion."""
        mocker.patch("app.workers.tasks.lyrics_tasks.get_transcriber")
        mock_logger = mocker.patch("app.workers.tasks.lyrics_tasks.logger")

        warm_up_whisper_model()

        mock_logger.info.assert_any_call("Whisper model warm-up complete")

    def test_logs_on_failure(self, mocker) -> None:
        """Failed warm-up should log an error, not crash."""
        mocker.patch(
            "app.workers.tasks.lyrics_tasks.get_transcriber",
            side_effect=RuntimeError("OOM"),
        )
        mock_logger = mocker.patch("app.workers.tasks.lyrics_tasks.logger")

        # Should NOT raise
        warm_up_whisper_model()

        mock_logger.exception.assert_called_once()
        assert "failed to load" in mock_logger.exception.call_args[0][0].lower()
