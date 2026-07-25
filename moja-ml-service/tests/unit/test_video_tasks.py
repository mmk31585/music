"""Tests for the video processing Celery task helpers.

Follows the pattern from test_lyrics_tasks.py.
"""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from app.workers.tasks.video_tasks import (
    _bump_retry,
    _fail_job,
    _set_status,
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
    def test_sets_status(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "processing")
        assert mock_job.status == "processing"

    def test_calls_session_add(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _set_status(mock_session, mock_job, "completed")
        mock_session.add.assert_called_once_with(mock_job)


class TestFailJob:
    def test_sets_failed_status(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "FFmpeg error")
        assert mock_job.status == "failed"

    def test_stores_error_message(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "corrupt video")
        assert mock_job.error_message == "corrupt video"

    def test_calls_commit(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _fail_job(mock_session, mock_job, "error")
        mock_session.commit.assert_called_once()


class TestBumpRetry:
    def test_increments_count(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 1

    def test_increments_from_existing(self, mock_session: MagicMock, mock_job: MagicMock) -> None:
        mock_job.retry_count = 1
        _bump_retry(mock_session, mock_job, "timeout")
        assert mock_job.retry_count == 2
