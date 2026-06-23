"""
Orchestration layer for audio feature extraction.

``AudioFeatureExtractionService`` is the **single entry point** that
Phase 2's Celery task will call. It wires together:

  1. ``AudioEmbeddingExtractor.extract()``  → embedding vector
  2. ``extract_tempo()``                    → BPM (reusing the loaded audio)

The service does **not** know about HTTP, Celery, S3, or shared volumes.
It only knows "here is a file on disk". The outer layer (storage / worker)
is responsible for getting the audio to a local path before calling this
service.
"""

from __future__ import annotations

import logging
import time
from dataclasses import dataclass, field

import librosa

from app.domain.audio_features.embedding_extractor import (
    AudioEmbeddingExtractor,
    AudioEmbeddingResult,
)
from app.domain.audio_features.tempo_extractor import extract_tempo

logger = logging.getLogger(__name__)


# ── Custom exceptions ─────────────────────────────────────────────────


class AudioFeatureExtractionError(Exception):
    """Raised when the full audio feature extraction pipeline fails.

    Wraps any internal exception (openl3 crash, I/O error, etc.) so
    callers see a single exception type.
    """


# ── Result dataclass ──────────────────────────────────────────────────


@dataclass
class AudioFeatureExtractionResult:
    """The final output of the audio feature extraction pipeline."""

    embedding: list[float] = field(default_factory=list)
    embedding_dim: int = 0
    embedding_model: str = ""
    tempo_bpm: float | None = None
    extraction_duration_seconds: float = 0.0


# ── Service ───────────────────────────────────────────────────────────


class AudioFeatureExtractionService:
    """Orchestrate the full audio feature extraction pipeline.

    Usage::

        service = AudioFeatureExtractionService(extractor=extractor)
        result = service.extract_all(audio_path="/path/to/track.mp3")
    """

    def __init__(self, extractor: AudioEmbeddingExtractor) -> None:
        self._extractor = extractor

    def extract_all(self, audio_path: str) -> AudioFeatureExtractionResult:
        """Run the full pipeline: load audio → extract embedding → extract tempo.

        Loads the audio file ONCE and reuses the loaded signal for both
        embedding and tempo extraction to avoid redundant I/O.

        Args:
            audio_path: Absolute path to a local audio file.

        Returns:
            An ``AudioFeatureExtractionResult`` with the embedding vector,
            model info, and optional BPM.

        Raises:
            AudioFeatureExtractionError: If any step in the pipeline fails.
        """
        from app.domain.audio_features.embedding_extractor import OPENL3_SR

        start_time = time.time()

        try:
            # ── Step 1: Load audio once ─────────────────────────────
            logger.info("Loading audio for feature extraction: %s", audio_path)
            try:
                y, sr = librosa.load(audio_path, sr=OPENL3_SR, mono=True)
            except Exception as exc:
                raise AudioFeatureExtractionError(
                    f"Failed to load audio file '{audio_path}': {exc}"
                ) from exc

            if len(y) == 0:
                raise AudioFeatureExtractionError(
                    f"Audio file '{audio_path}' loaded zero samples"
                )

            # ── Step 2: Extract embedding (reuse loaded audio) ──────
            logger.info("Extracting embedding from %s", audio_path)
            try:
                emb_result: AudioEmbeddingResult = self._extractor.extract_from_array(y, sr)
            except Exception as exc:
                raise AudioFeatureExtractionError(
                    f"Embedding extraction failed for '{audio_path}': {exc}"
                ) from exc

            # ── Step 3: Extract tempo (reuse loaded audio) ──────────
            tempo_bpm: float | None = None
            try:
                tempo_bpm = extract_tempo(y, sr)
                logger.info(
                    "Tempo extraction: %.1f BPM for %s",
                    tempo_bpm if tempo_bpm is not None else 0.0,
                    audio_path,
                )
            except Exception as exc:
                logger.warning(
                    "Tempo extraction failed for %s (non-fatal): %s",
                    audio_path,
                    exc,
                )

            elapsed = time.time() - start_time
            logger.info(
                "Feature extraction complete for %s: dim=%d tempo=%s duration=%.1fs",
                audio_path,
                emb_result.embedding_dim,
                f"{tempo_bpm:.1f}" if tempo_bpm is not None else "None",
                elapsed,
            )

            return AudioFeatureExtractionResult(
                embedding=emb_result.embedding,
                embedding_dim=emb_result.embedding_dim,
                embedding_model=emb_result.model_name,
                tempo_bpm=tempo_bpm,
                extraction_duration_seconds=elapsed,
            )

        except AudioFeatureExtractionError:
            raise
        except Exception as exc:
            raise AudioFeatureExtractionError(
                f"Audio feature extraction failed for '{audio_path}': {exc}"
            ) from exc
