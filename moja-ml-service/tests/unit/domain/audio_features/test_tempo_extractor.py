"""Tests for the tempo extractor module."""

from __future__ import annotations

import numpy as np
import pytest

from app.domain.audio_features.tempo_extractor import (
    TempoExtractionError,
    extract_tempo,
)


def test_returns_bpm_for_plausible_value(mocker) -> None:
    """A valid BPM within 40-220 should be returned."""
    mock_librosa = mocker.patch("app.domain.audio_features.tempo_extractor.librosa")
    mock_librosa.beat.beat_track.return_value = (120.0, np.array([1, 5, 10, 15, 20]))

    y = np.random.rand(48000 * 10)  # 10s at 48kHz
    result = extract_tempo(y, sr=48000)

    assert result == 120.0


def test_returns_none_for_implausible_bpm(mocker) -> None:
    """Out-of-range BPM values should return None."""
    mock_librosa = mocker.patch("app.domain.audio_features.tempo_extractor.librosa")
    mock_librosa.beat.beat_track.return_value = (300.0, np.array([1, 5, 10]))

    y = np.random.rand(48000 * 10)
    result = extract_tempo(y, sr=48000)

    assert result is None


def test_returns_none_for_zero_bpm(mocker) -> None:
    """Zero BPM (no beat detected) should return None."""
    mock_librosa = mocker.patch("app.domain.audio_features.tempo_extractor.librosa")
    mock_librosa.beat.beat_track.return_value = (0.0, np.array([]))

    y = np.random.rand(48000 * 10)
    result = extract_tempo(y, sr=48000)

    assert result is None


def test_returns_none_for_very_short_audio() -> None:
    """Less than 1 second of audio should return None without calling librosa."""
    y = np.random.rand(100)  # Way less than 1s at any reasonable sample rate
    result = extract_tempo(y, sr=48000)

    assert result is None


def test_returns_none_for_too_few_beat_frames(mocker) -> None:
    """Fewer than 2 beat frames indicates unreliable estimate."""
    mock_librosa = mocker.patch("app.domain.audio_features.tempo_extractor.librosa")
    mock_librosa.beat.beat_track.return_value = (100.0, np.array([1]))

    y = np.random.rand(48000 * 10)
    result = extract_tempo(y, sr=48000)

    assert result is None


def test_wraps_librosa_errors_in_tempo_extraction_error(mocker) -> None:
    """Librosa failures should be wrapped in TempoExtractionError."""
    mock_librosa = mocker.patch("app.domain.audio_features.tempo_extractor.librosa")
    mock_librosa.beat.beat_track.side_effect = RuntimeError("librosa crashed")

    y = np.random.rand(48000 * 10)

    with pytest.raises(TempoExtractionError) as exc_info:
        extract_tempo(y, sr=48000)

    assert "librosa crashed" in str(exc_info.value)


def test_returns_none_for_negative_bpm(mocker) -> None:
    """Negative BPM should return None."""
    mock_librosa = mocker.patch("app.domain.audio_features.tempo_extractor.librosa")
    mock_librosa.beat.beat_track.return_value = (-1.0, np.array([1, 5, 10]))

    y = np.random.rand(48000 * 10)
    result = extract_tempo(y, sr=48000)

    assert result is None
