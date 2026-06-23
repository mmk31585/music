"""Tests for the embedding Celery task helpers.

Follows the pattern from test_lyrics_tasks.py — tests the helper
functions with mocked DB session.
"""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from app.workers.tasks.embedding_tasks import (
    _bump_retry,
    _fail_job,
    _set_status,
    warm_up_embedding_model,
)


@pytest.fixture
def mock_job() -> MagicMock:
    job = MagicMock()
    job.status = "queued"
    job.retry_count = 0
    job.error_message = None
    return job


@pytest.fixture
def mock_session() -> MagicMock:
    session = MagicMock()
    session.commit.return_value = None
    return session


class TestSetStatus:
    def test_sets_status_on_job(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "extracting")
        assert mock_job.status == "extracting"

    def test_calls_session_add(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "completed")
        mock_session.add.assert_called_once_with(mock_job)

    def test_does_not_call_commit(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "downloading")
        mock_session.commit.assert_not_called()


class TestFailJob:
    def test_sets_failed_status(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "Corrupt audio")
        assert mock_job.status == "failed"

    def test_stores_error_message(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "openl3 error")
        assert mock_job.error_message == "openl3 error"

    def test_calls_session_add(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "error")
        mock_session.add.assert_called_once_with(mock_job)

    def test_calls_session_commit(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "error")
        mock_session.commit.assert_called_once()


class TestBumpRetry:
    def test_increments_retry_count(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 1

    def test_increments_from_existing(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        mock_job.retry_count = 2
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 3

    def test_stores_error_message(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "connection refused")
        assert mock_job.error_message == "connection refused"


class TestWarmUpEmbeddingModel:
    def test_calls_get_embedding_extractor(self, mocker) -> None:
        mock_get = mocker.patch(
            "app.workers.tasks.embedding_tasks.get_embedding_extractor"
        )
        mock_settings = mocker.patch(
            "app.workers.tasks.embedding_tasks.settings"
        )

        warm_up_embedding_model()

        mock_get.assert_called_once_with(mock_settings)

    def test_logs_on_success(self, mocker) -> None:
        mocker.patch(
            "app.workers.tasks.embedding_tasks.get_embedding_extractor"
        )
        mock_logger = mocker.patch(
            "app.workers.tasks.embedding_tasks.logger"
        )

        warm_up_embedding_model()

        mock_logger.info.assert_any_call("OpenL3 model warm-up complete")

    def test_logs_on_failure(self, mocker) -> None:
        mocker.patch(
            "app.workers.tasks.embedding_tasks.get_embedding_extractor",
            side_effect=RuntimeError("OOM"),
        )
        mock_logger = mocker.patch(
            "app.workers.tasks.embedding_tasks.logger"
        )

        warm_up_embedding_model()

        mock_logger.exception.assert_called_once()
        assert "failed to load" in mock_logger.exception.call_args[0][0].lower()
