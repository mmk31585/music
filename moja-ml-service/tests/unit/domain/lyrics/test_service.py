"""Tests for the LyricsGenerationService orchestrator.

These tests use a **mocked** Transcriber — no real Whisper model is loaded.
That is an integration-test concern (``tests/integration/``).
"""

import pytest

from app.domain.lyrics.service import LyricsGenerationError, LyricsGenerationService
from app.domain.lyrics.transcriber import Transcriber, TranscriptionResult, TranscriptSegment


def test_service_orchestrates_full_pipeline(mocker) -> None:
    """The service should call transcribe, format, and score correctly."""
    mock_transcriber = mocker.Mock(spec=Transcriber)
    mock_transcriber.transcribe.return_value = TranscriptionResult(
        segments=[
            TranscriptSegment(start=1.0, end=2.0, text="سلام", confidence=0.9),
            TranscriptSegment(start=3.0, end=4.0, text="دنیا", confidence=0.8),
        ],
        detected_language="fa",
        language_probability=0.95,
        duration_seconds=180.0,
    )

    service = LyricsGenerationService(transcriber=mock_transcriber)
    result = service.generate(audio_path="/fake/path.mp3")

    assert "سلام" in result.lrc_content
    assert "دنیا" in result.lrc_content
    assert result.detected_language == "fa"
    assert result.confidence > 0
    assert result.segment_count == 2
    assert "سلام دنیا" in result.plain_text

    mock_transcriber.transcribe.assert_called_once_with("/fake/path.mp3")


def test_service_includes_headers_when_provided(mocker) -> None:
    """Track title and artist should appear in LRC output."""
    mock_transcriber = mocker.Mock(spec=Transcriber)
    mock_transcriber.transcribe.return_value = TranscriptionResult(
        segments=[
            TranscriptSegment(start=1.0, end=2.0, text="سلام", confidence=0.9),
        ],
        detected_language="fa",
        language_probability=0.95,
        duration_seconds=60.0,
    )

    service = LyricsGenerationService(transcriber=mock_transcriber)
    result = service.generate(
        audio_path="/fake/path.mp3",
        track_title="آهنگ تست",
        track_artist="خواننده تست",
    )

    assert "[ti:آهنگ تست]" in result.lrc_content
    assert "[ar:خواننده تست]" in result.lrc_content


def test_service_wraps_transcriber_errors(mocker) -> None:
    """Raw Transcriber errors should be wrapped in LyricsGenerationError."""
    mock_transcriber = mocker.Mock(spec=Transcriber)
    mock_transcriber.transcribe.side_effect = Exception("whisper crashed")

    service = LyricsGenerationService(transcriber=mock_transcriber)
    with pytest.raises(LyricsGenerationError) as exc_info:
        service.generate(audio_path="/fake/path.mp3")

    assert "whisper crashed" in str(exc_info.value)


def test_service_empty_transcription(mocker) -> None:
    """Empty transcription (zero segments) should not crash."""
    mock_transcriber = mocker.Mock(spec=Transcriber)
    mock_transcriber.transcribe.return_value = TranscriptionResult(
        segments=[],
        detected_language="fa",
        language_probability=0.0,
        duration_seconds=0.0,
    )

    service = LyricsGenerationService(transcriber=mock_transcriber)
    result = service.generate(audio_path="/fake/path.mp3")

    assert result.lrc_content == ""
    assert result.plain_text == ""
    assert result.confidence == 0.0
    assert result.segment_count == 0


def test_service_normalizes_persian_text(mocker) -> None:
    """Arabic characters in segments should be normalized in output."""
    mock_transcriber = mocker.Mock(spec=Transcriber)
    mock_transcriber.transcribe.return_value = TranscriptionResult(
        segments=[
            TranscriptSegment(start=1.0, end=2.0, text="علي", confidence=0.9),
        ],
        detected_language="fa",
        language_probability=0.95,
        duration_seconds=60.0,
    )

    service = LyricsGenerationService(transcriber=mock_transcriber)
    result = service.generate(audio_path="/fake/path.mp3")

    # "علي" (Arabic yeh) should become "علی" (Persian yeh)
    assert "علی" in result.lrc_content
    assert "علي" not in result.lrc_content
    assert "علی" in result.plain_text
