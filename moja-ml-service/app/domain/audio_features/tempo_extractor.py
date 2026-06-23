"""
BPM / tempo extraction from audio files.

While the audio file is already loaded for embedding extraction, we
extract tempo using librosa's onset-based beat tracker. This avoids
re-processing the file a second time and makes tempo available for
metadata-based filtering in Phase 2, as well as for DJ-mode crossfade
logic later.

Domain layer — no HTTP or Celery knowledge.
"""

from __future__ import annotations

import logging

import librosa
import numpy as np

logger = logging.getLogger(__name__)


class TempoExtractionError(Exception):
    """Raised when tempo extraction fails for reasons other than
    "this track is too ambient to have a beat."""


def extract_tempo(y: np.ndarray, sr: int) -> float | None:
    """Estimate the tempo (BPM) of an audio signal using onset detection.

    Uses ``librosa.beat.beat_track()`` for a standard onset-based BPM
    estimate. This is cheap compared to embedding extraction — it runs
    on the already-loaded audio array, avoiding a second file read.

    Args:
        y: Audio time series (numpy array), as loaded by librosa.
        sr: Sample rate of ``y``, in Hz.

    Returns:
        Estimated BPM as a float, or ``None`` if:
        - The track is too ambient / beatless for a confident estimate.
        - The estimated BPM is outside the plausible range (40–220 BPM).
        - librosa could not determine a tempo (onset envelope too flat).

    Raises:
        TempoExtractionError: If librosa itself fails (not a "no beat"
            situation, but an actual error like a corrupt array).
    """
    from app.domain.audio_features.embedding_extractor import BPM_MAX, BPM_MIN

    if len(y) < sr:  # Less than 1 second of audio
        logger.debug("Audio too short for tempo extraction (< 1s)")
        return None

    try:
        tempo, beat_frames = librosa.beat.beat_track(y=y, sr=sr)

        # beat_track returns a 0-d or 1-d array in some librosa versions
        if hasattr(tempo, "item"):
            tempo = tempo.item()

        if not isinstance(tempo, (int, float)) or tempo is None:
            return None

        tempo = float(tempo)

        # No beat detected or implausible value
        if tempo <= 0 or tempo < BPM_MIN or tempo > BPM_MAX:
            logger.debug(
                "Tempo %.1f BPM outside plausible range [%.0f-%.0f] — returning None",
                tempo,
                BPM_MIN,
                BPM_MAX,
            )
            return None

        # Confidence check: if very few beat frames detected relative
        # to track length, the estimate is unreliable.
        if beat_frames is not None and hasattr(beat_frames, "__len__"):
            if len(beat_frames) < 2:
                logger.debug("Too few beat frames (%d) — returning None", len(beat_frames))
                return None

        n_frames = len(beat_frames) if beat_frames is not None else 0
        logger.debug("Extracted tempo: %.1f BPM (frames=%d)", tempo, n_frames)
        return round(tempo, 1)

    except Exception as exc:
        raise TempoExtractionError(f"Tempo extraction failed: {exc}") from exc
