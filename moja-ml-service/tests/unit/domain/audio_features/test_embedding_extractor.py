"""Tests for AudioEmbeddingExtractor.

These tests use a **mocked** openl3 / librosa — no real model is loaded.
That is an integration-test concern (``tests/integration/``).
"""

from __future__ import annotations

import numpy as np
import pytest

from app.domain.audio_features.embedding_extractor import (
    AudioEmbeddingError,
    AudioEmbeddingExtractor,
    AudioEmbeddingResult,
)


def test_extract_averages_multiple_windows(mocker) -> None:
    """Multi-window extraction should average per-window embeddings."""
    mock_librosa = mocker.patch("app.domain.audio_features.embedding_extractor.librosa")

    # Simulate a 120-second track at 48kHz
    sr = 48000
    duration = 120.0
    n_samples = int(sr * duration)
    mock_librosa.load.return_value = (np.random.rand(n_samples).astype(np.float32), sr)

    # Mock openl3 instance — set via _openl3 attribute rather than module patch
    mock_openl3 = mocker.MagicMock()
    mock_openl3.get_audio_embedding.side_effect = [
        (np.array([[0.1, 0.2, 0.3]]), np.array([0.0])),   # window 1
        (np.array([[0.4, 0.5, 0.6]]), np.array([0.0])),   # window 2
        (np.array([[0.7, 0.8, 0.9]]), np.array([0.0])),   # window 3
    ]

    extractor = AudioEmbeddingExtractor(model_variant="mel256", embedding_size=512)
    extractor._model_loaded = True
    extractor._openl3 = mock_openl3

    result = extractor.extract("/fake/path.mp3")

    assert isinstance(result, AudioEmbeddingResult)
    assert len(result.embedding) == 3
    expected = [(0.1 + 0.4 + 0.7) / 3, (0.2 + 0.5 + 0.8) / 3, (0.3 + 0.6 + 0.9) / 3]
    assert result.embedding == pytest.approx(expected, rel=1e-5)
    assert result.embedding_dim == 3
    assert result.model_name == "openl3-mel256-music-512"


def test_short_track_uses_single_window(mocker) -> None:
    """A track under 30s should not attempt 3 windows."""
    mock_librosa = mocker.patch("app.domain.audio_features.embedding_extractor.librosa")

    # Simulate a 10-second track
    sr = 48000
    n_samples = int(sr * 10)
    mock_librosa.load.return_value = (np.random.rand(n_samples).astype(np.float32), sr)

    mock_openl3 = mocker.MagicMock()
    mock_openl3.get_audio_embedding.return_value = (
        np.array([[0.5, 0.5, 0.5]]),
        np.array([0.0]),
    )

    extractor = AudioEmbeddingExtractor()
    extractor._model_loaded = True
    extractor._openl3 = mock_openl3

    result = extractor.extract("/fake/short.mp3")

    assert mock_openl3.get_audio_embedding.call_count == 1
    assert result.embedding_dim == 3


def test_wraps_extraction_errors(mocker) -> None:
    """librosa failure should be wrapped in AudioEmbeddingError."""
    mock_librosa = mocker.patch("app.domain.audio_features.embedding_extractor.librosa")
    mock_librosa.load.side_effect = Exception("file not found")

    extractor = AudioEmbeddingExtractor()
    extractor._model_loaded = True

    with pytest.raises(AudioEmbeddingError) as exc_info:
        extractor.extract("/nonexistent/file.mp3")

    assert "file not found" in str(exc_info.value)


def test_empty_audio_raises_error(mocker) -> None:
    """Zero-sample audio should raise AudioEmbeddingError."""
    mock_librosa = mocker.patch("app.domain.audio_features.embedding_extractor.librosa")
    mock_librosa.load.return_value = (np.array([], dtype=np.float32), 48000)

    extractor = AudioEmbeddingExtractor()
    extractor._model_loaded = True

    with pytest.raises(AudioEmbeddingError) as exc_info:
        extractor.extract("/fake/empty.mp3")

    assert "zero samples" in str(exc_info.value)


def test_get_embedding_extractor_singleton(mocker) -> None:
    """get_embedding_extractor should return the same instance on repeated calls."""
    mocker.patch("app.domain.audio_features.embedding_extractor.AudioEmbeddingExtractor")

    import app.domain.audio_features.embedding_extractor as mod
    from app.config import Settings

    mod._extractor_instance = None

    settings = Settings()  # type: ignore[call-arg]
    get_embedding_extractor = mod.get_embedding_extractor

    first = get_embedding_extractor(settings)
    second = get_embedding_extractor(settings)

    assert first is second
    assert mod._extractor_instance is first
