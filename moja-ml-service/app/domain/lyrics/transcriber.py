"""
Transcriber wrapper around faster-whisper.

Domain layer — no HTTP or Celery knowledge. This module provides:

- ``Transcriber`` class: wraps ``faster_whisper.WhisperModel`` with a clean
  interface and typed result dataclasses.
- ``get_transcriber()``: process-wide singleton accessor so the model is
  loaded exactly once per process.
- ``TranscriptionError``: custom exception that wraps any raw exception from
  faster-whisper / ctranslate2.
"""

from __future__ import annotations

import logging
from dataclasses import dataclass, field
from typing import TYPE_CHECKING

from faster_whisper import WhisperModel

if TYPE_CHECKING:
    from app.config import Settings

logger = logging.getLogger(__name__)


# ── Custom exceptions ──────────────────────────────────────────────────


class TranscriptionError(Exception):
    """Raised when transcription fails for any reason.

    Wraps the original exception as the cause so the chain is preserved::

        raise TranscriptionError("whisper crashed") from original_exc
    """


# ── Dataclasses ────────────────────────────────────────────────────────


@dataclass
class TranscriptSegment:
    """A single transcribed segment with precise timing."""

    start: float  # seconds
    end: float  # seconds
    text: str
    confidence: float  # avg_logprob normalized to 0-1


@dataclass
class TranscriptionResult:
    """Complete transcription output for one audio file."""

    segments: list[TranscriptSegment] = field(default_factory=list)
    detected_language: str = ""
    language_probability: float = 0.0
    duration_seconds: float = 0.0


# ── Transcriber ────────────────────────────────────────────────────────


class Transcriber:
    """Wraps a faster-whisper model with a typed interface.

    The model is loaded *once* at construction time. Clients should use the
    module-level ``get_transcriber()`` singleton to avoid redundant loads.
    """

    def __init__(
        self,
        model_size: str,
        compute_type: str,
        device: str,
        language_hint: str | None = None,
    ) -> None:
        logger.info(
            "Loading Whisper model: size=%s device=%s compute=%s",
            model_size,
            device,
            compute_type,
        )
        try:
            self._model = WhisperModel(model_size, device=device, compute_type=compute_type)
        except Exception as exc:
            raise TranscriptionError(
                f"Failed to load Whisper model '{model_size}' on {device}: {exc}"
            ) from exc
        self._language_hint = language_hint
        logger.info("Whisper model loaded successfully")

    def transcribe(self, audio_path: str) -> TranscriptionResult:
        """Transcribe an audio file and return typed results.

        Word-level timestamps are enabled (``word_timestamps=True``) because
        LRC formatting needs precise line timing.

        Args:
            audio_path: Path to an audio file readable by ffmpeg.

        Returns:
            A ``TranscriptionResult`` with per-segment timing and text.

        Raises:
            TranscriptionError: If transcription fails for any reason.
        """
        try:
            # faster-whisper's ``language`` param:
            #   - If ``None``: force auto-detect.
            #   - If a string: force that language model (Whisper will NOT
            #     auto-detect; it runs the specified language's decoder).
            #
            # We pass ``None`` so Whisper auto-detects per track — this is
            # critical for a multi-language catalog where English, Turkish,
            # Arabic etc. songs coexist with Persian. Forcing ``fa`` makes
            # non-Persian tracks produce garbled output.
            language = None

            segments_gen, info = self._model.transcribe(
                audio_path,
                language=language,
                word_timestamps=True,
                vad_filter=True,
                vad_parameters=dict(min_silence_duration_ms=500),
            )

            segments_list = list(segments_gen)

            detected_language = getattr(info, "language", self._language_hint or "unknown")
            language_probability = getattr(info, "language_probability", 0.0)
            duration = getattr(info, "duration", 0.0)

            segments = [
                TranscriptSegment(
                    start=seg.start,
                    end=seg.end,
                    text=seg.text.strip(),
                    confidence=_normalize_logprob(seg.avg_logprob),
                )
                for seg in segments_list
            ]

            return TranscriptionResult(
                segments=segments,
                detected_language=detected_language,
                language_probability=language_probability,
                duration_seconds=duration,
            )

        except TranscriptionError:
            raise
        except Exception as exc:
            raise TranscriptionError(
                f"Whisper transcription failed for '{audio_path}': {exc}"
            ) from exc


# ── Helpers ────────────────────────────────────────────────────────────


def _normalize_logprob(avg_logprob: float) -> float:
    """Convert Whisper's ``avg_logprob`` (roughly -20 to 0) to a 0-1 score.

    The mapping is heuristic:
      - -16.0 or lower → 0.0
      -   0.0          → 1.0
      - linear clamp in between.

    This is not a probability — it is a normalised confidence score suitable
    for the ``TranscriptSegment.confidence`` field.
    """
    # Typical range for decent transcriptions: -5 to 0
    # Clamp at -16 (very bad) to 0 (perfect)
    clamped = max(-16.0, min(0.0, avg_logprob))
    return (clamped + 16.0) / 16.0


# ── Singleton accessor ─────────────────────────────────────────────────


_transcriber_instance: Transcriber | None = None


def get_transcriber(settings: Settings) -> Transcriber:
    """Return the process-wide ``Transcriber`` singleton.

    The model is loaded on the **first** call and cached thereafter. Both
    the FastAPI readiness probe (Phase 1) and the Celery worker (Phase 3)
    should call this function to get a warm model without reloading it.
    """
    global _transcriber_instance  # noqa: PLW0603

    if _transcriber_instance is None:
        _transcriber_instance = Transcriber(
            model_size=settings.whisper_model_size,
            compute_type=settings.whisper_compute_type,
            device=settings.whisper_device,
            language_hint=settings.whisper_language_hint,
        )
    return _transcriber_instance
