"""
Integration test for the real openl3 / librosa audio embedding extractor.

This test loads an **actual** openl3 model and processes a real audio file.
It is **slow** (CPU inference can take 30s–5min depending on model size and
file length) and requires a sample audio file to exist at
``tests/fixtures/sample_audio.mp3``.

The human operator must provide this file. A suitable sample can be any
short audio clip (10–30 seconds) in MP3 or WAV format.

Skipped by default. Run explicitly::

    uv run pytest -m integration -v
"""

from __future__ import annotations

import os

import pytest

from app.domain.audio_features.embedding_extractor import (
    AudioEmbeddingError,
    AudioEmbeddingExtractor,
)

SAMPLE_FILE = os.path.join(os.path.dirname(__file__), "..", "fixtures", "sample_audio.mp3")


@pytest.mark.integration
@pytest.mark.slow
def test_real_embedding_extraction_on_sample_file() -> None:
    """Load a real openl3 model and extract embedding from a sample file."""
    if not os.path.isfile(SAMPLE_FILE):
        pytest.skip(
            f"Sample audio file not found at {SAMPLE_FILE}. "
            "Provide a short audio clip (10-30s) in MP3 or WAV format to enable this test."
        )

    extractor = AudioEmbeddingExtractor(
        model_variant="mel256",
        embedding_size=512,
        device="cpu",
    )

    result = extractor.extract(SAMPLE_FILE)

    # Basic sanity checks
    assert len(result.embedding) > 0, "Embedding is empty"
    assert result.embedding_dim == len(result.embedding)
    assert result.embedding_dim > 0
    assert result.model_name == "openl3-mel256-music-512"
    assert result.duration_analyzed_seconds > 0

    # Embedding values should be finite floats
    for val in result.embedding:
        assert isinstance(val, float), f"Non-float value: {val}"
        assert val == val, "NaN in embedding"  # noqa: PLR0124


@pytest.mark.integration
@pytest.mark.slow
def test_real_embedding_nonexistent_file() -> None:
    """Extracting from a non-existent file should raise AudioEmbeddingError."""
    extractor = AudioEmbeddingExtractor(
        model_variant="mel256",
        embedding_size=512,
        device="cpu",
    )

    with pytest.raises(AudioEmbeddingError):
        extractor.extract("/nonexistent/path/to/file.mp3")
