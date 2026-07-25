"""Tests for the AudioFeatureExtractionService orchestrator.

These tests use **mocked** dependencies — no real openl3 / librosa model
is loaded. That is an integration-test concern.
"""

from __future__ import annotations

import numpy as np
import pytest

from app.domain.audio_features.embedding_extractor import (
    AudioEmbeddingExtractor,
    AudioEmbeddingResult,
)
from app.domain.audio_features.service import (
    AudioFeatureExtractionError,
    AudioFeatureExtractionService,
)


def test_service_orchestrates_extraction_pipeline(mocker) -> None:
    """The service should load audio once, extract embedding + tempo."""
    mock_librosa = mocker.patch("app.domain.audio_features.service.librosa")

    sr = 48000
    mock_librosa.load.return_value = (np.random.rand(int(sr * 60)).astype(np.float32), sr)

    mock_extractor = mocker.Mock(spec=AudioEmbeddingExtractor)
    mock_extractor.extract.return_value = AudioEmbeddingResult(
        embedding=[0.1, 0.2, 0.3],
        embedding_dim=3,
        model_name="openl3-mel256-music-512",
        duration_analyzed_seconds=30.0,
    )

    mock_tempo = mocker.patch("app.domain.audio_features.service.extract_tempo")
    mock_tempo.return_value = 120.0

    service = AudioFeatureExtractionService(extractor=mock_extractor)
    result = service.extract_all(audio_path="/fake/path.mp3")

    assert result.embedding == [0.1, 0.2, 0.3]
    assert result.embedding_dim == 3
    assert result.embedding_model == "openl3-mel256-music-512"
    assert result.tempo_bpm == 120.0
    assert result.extraction_duration_seconds > 0

    # Verify audio loaded once
    mock_librosa.load.assert_called_once_with("/fake/path.mp3", sr=48000, mono=True)

    # Verify both extractors called
    mock_extractor.extract.assert_called_once_with("/fake/path.mp3")
    mock_tempo.assert_called_once()


def test_service_handles_missing_tempo_gracefully(mocker) -> None:
    """Tempo extraction failure should not crash the pipeline."""
    mock_librosa = mocker.patch("app.domain.audio_features.service.librosa")
    sr = 48000
    mock_librosa.load.return_value = (np.random.rand(int(sr * 60)).astype(np.float32), sr)

    mock_extractor = mocker.Mock(spec=AudioEmbeddingExtractor)
    mock_extractor.extract.return_value = AudioEmbeddingResult(
        embedding=[0.1, 0.2],
        embedding_dim=2,
        model_name="openl3-mel256-music-512",
        duration_analyzed_seconds=30.0,
    )

    mock_tempo = mocker.patch("app.domain.audio_features.service.extract_tempo")
    mock_tempo.return_value = None  # Ambient track, no tempo

    service = AudioFeatureExtractionService(extractor=mock_extractor)
    result = service.extract_all(audio_path="/fake/path.mp3")

    assert result.embedding == [0.1, 0.2]
    assert result.tempo_bpm is None


def test_service_wraps_embedding_errors(mocker) -> None:
    """Embedding extractor errors should be wrapped."""
    mock_librosa = mocker.patch("app.domain.audio_features.service.librosa")
    sr = 48000
    mock_librosa.load.return_value = (np.random.rand(int(sr * 60)).astype(np.float32), sr)

    mock_extractor = mocker.Mock(spec=AudioEmbeddingExtractor)
    mock_extractor.extract.side_effect = Exception("openl3 crashed")

    service = AudioFeatureExtractionService(extractor=mock_extractor)

    with pytest.raises(AudioFeatureExtractionError) as exc_info:
        service.extract_all(audio_path="/fake/path.mp3")

    assert "openl3 crashed" in str(exc_info.value)


def test_service_wraps_librosa_errors(mocker) -> None:
    """librosa load failure should be wrapped."""
    mock_librosa = mocker.patch("app.domain.audio_features.service.librosa")
    mock_librosa.load.side_effect = Exception("ffmpeg not found")

    mock_extractor = mocker.Mock(spec=AudioEmbeddingExtractor)
    service = AudioFeatureExtractionService(extractor=mock_extractor)

    with pytest.raises(AudioFeatureExtractionError) as exc_info:
        service.extract_all(audio_path="/bad/file.mp3")

    assert "ffmpeg not found" in str(exc_info.value)


def test_service_handles_empty_audio(mocker) -> None:
    """Zero-sample audio should raise."""
    mock_librosa = mocker.patch("app.domain.audio_features.service.librosa")
    mock_librosa.load.return_value = (np.array([], dtype=np.float32), 48000)

    mock_extractor = mocker.Mock(spec=AudioEmbeddingExtractor)
    service = AudioFeatureExtractionService(extractor=mock_extractor)

    with pytest.raises(AudioFeatureExtractionError) as exc_info:
        service.extract_all(audio_path="/fake/empty.mp3")

    assert "zero samples" in str(exc_info.value)
