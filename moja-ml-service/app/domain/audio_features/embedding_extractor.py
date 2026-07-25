"""
Audio embedding extractor using openl3.

Domain layer — no HTTP or Celery knowledge. This module provides:

- ``AudioEmbeddingExtractor``: wraps openl3 with a clean interface and
  typed result dataclass. Extracts a 512-dim embedding from each track
  by averaging three 30-second windows (at 25%, 50%, 75% of duration)
  for CPU-cost-bounded processing regardless of track length.
- ``get_embedding_extractor()``: process-wide singleton accessor so the
  model is loaded exactly once per process.
- ``AudioEmbeddingError``: custom exception wrapping any raw exception
  from openl3 / tensorflow / librosa.

openl3 / tensorflow-cpu installation note
------------------------------------------
If ``pip install openl3>=0.4.2 tensorflow-cpu>=2.15`` fails in the
Docker environment (e.g. incompatible glibc, missing AVX instructions),
the fallback is ``essentia`` (`pip install essentia`). Essentia is a
C++-based audio analysis library that can extract MFCC-style or
temporal feature embeddings on CPU without TensorFlow.

The embedding dimensionality and downstream Phase 2 vector index
assumptions will change with essentia — if you switch, update:
  - ``AudioEmbeddingResult.embedding_dim`` field value
  - ``AudioEmbeddingResult.model_name`` to ``"essentia-mfcc"``
  - The Phase 2 vector index dimension expectation
"""

from __future__ import annotations

import logging
from dataclasses import dataclass, field
from typing import TYPE_CHECKING

import librosa
import numpy as np

if TYPE_CHECKING:
    from app.config import Settings

logger = logging.getLogger(__name__)

# ── OpenL3 constants ──────────────────────────────────────────────────

# openl3 expects mono audio at this sample rate
OPENL3_SR: int = 48000

# Content (non-silent) windows: we process up to this many per track
MAX_ANALYSIS_WINDOWS: int = 3
# Each window duration in seconds
WINDOW_DURATION_S: float = 30.0

# BPM sanity range for tempo extraction (used by tempo_extractor too)
BPM_MIN: float = 40.0
BPM_MAX: float = 220.0


# ── Custom exceptions ─────────────────────────────────────────────────


class AudioEmbeddingError(Exception):
    """Raised when audio embedding extraction fails for any reason.

    Wraps the original exception as the cause so the chain is preserved::

        raise AudioEmbeddingError("openl3 crashed") from original_exc
    """


# ── Dataclasses ───────────────────────────────────────────────────────


@dataclass
class AudioEmbeddingResult:
    """The fixed-dimension embedding vector for one audio track."""

    embedding: list[float] = field(default_factory=list)
    embedding_dim: int = 0
    model_name: str = ""
    duration_analyzed_seconds: float = 0.0


# ── Embedding extractor ───────────────────────────────────────────────


class AudioEmbeddingExtractor:
    """Wraps openl3 with a typed interface for music embedding extraction.

    The model is loaded *once* at construction time. Clients should use the
    module-level ``get_embedding_extractor()`` singleton to avoid redundant
    loads.
    """

    def __init__(
        self,
        model_variant: str = "mel256",
        embedding_size: int = 512,
        device: str = "cpu",
    ) -> None:
        self._model_variant = model_variant
        self._embedding_size = embedding_size
        self._device = device
        self._model_loaded = False

        # Input shape expected by openl3: (n_windows, n_samples)
        # We compute n_samples once here.
        self._input_duration = WINDOW_DURATION_S
        self._input_samples = int(OPENL3_SR * WINDOW_DURATION_S)

        logger.info(
            "AudioEmbeddingExtractor created: variant=%s embedding_size=%d device=%s",
            model_variant,
            embedding_size,
            device,
        )

    def _ensure_model_loaded(self) -> None:
        """Lazy-load the openl3 model on first use.

        Mirrors the lazy-load pattern in Transcriber (transcriber.py) but
        uses a method-level flag rather than a class-level singleton because
        the module-level ``get_embedding_extractor()`` handles the singleton.

        openl3 is imported lazily because it pulls in TensorFlow (heavy)
        and may not be installed in dev/test environments.
        """
        if self._model_loaded:
            return

        try:
            import openl3  # noqa: F811
        except ImportError as exc:
            raise AudioEmbeddingError(
                "openl3 is not installed. "
                "Run: uv add openl3 tensorflow-cpu librosa\n"
                "If installation fails, see the fallback note in this module's docstring."
            ) from exc

        try:
            openl3.load_audio_embedding_model(
                input_repr=self._model_variant,
                content_type="music",
                embedding_size=self._embedding_size,
            )
            self._openl3 = openl3
            self._model_loaded = True
            logger.info(
                "openl3 model loaded: variant=%s embedding_size=%d",
                self._model_variant,
                self._embedding_size,
            )
        except Exception as exc:
            raise AudioEmbeddingError(
                f"Failed to load openl3 model (variant={self._model_variant}, "
                f"embedding_size={self._embedding_size}): {exc}"
            ) from exc

    def extract(self, audio_path: str) -> AudioEmbeddingResult:
        return self._extract_from_file(audio_path)

    def extract_from_array(self, y: np.ndarray, sr: int) -> AudioEmbeddingResult:
        return self._extract_from_array(y, sr)

    def _extract_from_file(self, audio_path: str) -> AudioEmbeddingResult:
        try:
            try:
                y, sr = librosa.load(audio_path, sr=OPENL3_SR, mono=True)
            except Exception as exc:
                raise AudioEmbeddingError(
                    f"Failed to load audio file '{audio_path}': {exc}"
                ) from exc

            if len(y) == 0:
                raise AudioEmbeddingError(
                    f"Audio file '{audio_path}' loaded zero samples"
                )

            return self._process_audio(y, sr)

        except AudioEmbeddingError:
            raise
        except Exception as exc:
            raise AudioEmbeddingError(
                f"Audio embedding extraction failed for '{audio_path}': {exc}"
            ) from exc

    def _extract_from_array(self, y: np.ndarray, sr: int) -> AudioEmbeddingResult:
        return self._process_audio(y, sr)

    def _process_audio(
        self, y: np.ndarray, sr: int
    ) -> AudioEmbeddingResult:
        total_duration = len(y) / sr

        windows: list[np.ndarray] = []

        if total_duration <= WINDOW_DURATION_S:
            windows.append(y)
            logger.debug(
                "Track shorter than window (%.1fs < %.1fs): single window",
                total_duration,
                WINDOW_DURATION_S,
            )
        else:
            fractions = [0.25, 0.50, 0.75]
            for frac in fractions:
                start_sample = int(frac * len(y) - self._input_samples / 2)
                start_sample = max(0, start_sample)
                end_sample = start_sample + self._input_samples
                if end_sample > len(y):
                    end_sample = len(y)
                    start_sample = max(0, end_sample - self._input_samples)
                window = y[start_sample:end_sample]
                if len(window) < self._input_samples:
                    pad_len = self._input_samples - len(window)
                    window = np.pad(window, (0, pad_len), mode="constant")
                windows.append(window)

        self._ensure_model_loaded()

        all_embeddings: list[np.ndarray] = []
        total_analyzed = 0.0

        for w in windows:
            emb, _ = self._openl3.get_audio_embedding(
                w,
                sr=OPENL3_SR,
                input_repr=self._model_variant,
                embedding_size=self._embedding_size,
            )
            emb_mean = np.mean(emb, axis=0)
            all_embeddings.append(emb_mean)
            total_analyzed += len(w) / sr

        if not all_embeddings:
            raise AudioEmbeddingError("No embeddings were generated")

        final_embedding = np.mean(all_embeddings, axis=0).tolist()

        return AudioEmbeddingResult(
            embedding=final_embedding,
            embedding_dim=len(final_embedding),
            model_name=f"openl3-{self._model_variant}-music-{self._embedding_size}",
            duration_analyzed_seconds=total_analyzed,
        )


# ── Singleton accessor ─────────────────────────────────────────────────


_extractor_instance: AudioEmbeddingExtractor | None = None


def get_embedding_extractor(settings: Settings) -> AudioEmbeddingExtractor:
    """Return the process-wide ``AudioEmbeddingExtractor`` singleton.

    The model is loaded on the **first** call and cached thereafter. Both
    the FastAPI readiness probe and the Celery worker should call this
    function to get a warm model without reloading it.

    Same singleton pattern as ``get_transcriber()`` in transcriber.py.
    """
    global _extractor_instance  # noqa: PLW0603

    if _extractor_instance is None:
        _extractor_instance = AudioEmbeddingExtractor(
            model_variant=settings.audio_embedding_model_variant,
            embedding_size=settings.audio_embedding_size,
            device=settings.audio_embedding_device,
        )
    return _extractor_instance
