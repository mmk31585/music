"""
Integration test for the real Whisper transcriber.

This test loads an **actual** faster-whisper model and transcribes a real
audio file. It is **slow** (CPU inference can take 30s–5min depending on
model size and file length) and requires a sample audio file to exist at
``tests/fixtures/sample_persian_audio.mp3``.

The human operator must provide this file. A suitable sample can be any
short Persian audio clip (10–30 seconds) in MP3 format.

Skipped by default. Run explicitly::

    uv run pytest -m integration -v
"""

import os

import pytest

from app.domain.lyrics.transcriber import Transcriber, TranscriptionError

SAMPLE_FILE = os.path.join(os.path.dirname(__file__), "..", "fixtures", "sample_persian_audio.mp3")


@pytest.mark.integration
@pytest.mark.slow
def test_real_whisper_transcription_on_sample_file() -> None:
    """Load a real Whisper model (tiny, for speed) and transcribe a sample file."""
    if not os.path.isfile(SAMPLE_FILE):
        pytest.skip(
            f"Sample audio file not found at {SAMPLE_FILE}. "
            "Provide a short Persian MP3 clip (10-30s) to enable this test."
        )

    transcriber = Transcriber(
        model_size="tiny",
        compute_type="int8",
        device="cpu",
        language_hint="fa",
    )

    result = transcriber.transcribe(SAMPLE_FILE)

    # Basic sanity checks
    assert len(result.segments) > 0, "Transcription returned zero segments"
    assert result.detected_language is not None
    assert result.duration_seconds > 0

    # Each non-empty segment should have plausible text and timing
    for seg in result.segments:
        if seg.text.strip():
            assert seg.end > seg.start, f"Segment end <= start: {seg}"
            assert 0.0 <= seg.confidence <= 1.0, f"Confidence out of range: {seg.confidence}"


@pytest.mark.integration
@pytest.mark.slow
def test_real_whisper_transcription_nonexistent_file() -> None:
    """Transcribing a non-existent file should raise TranscriptionError."""
    transcriber = Transcriber(
        model_size="tiny",
        compute_type="int8",
        device="cpu",
    )

    with pytest.raises(TranscriptionError):
        transcriber.transcribe("/nonexistent/path/to/file.mp3")
