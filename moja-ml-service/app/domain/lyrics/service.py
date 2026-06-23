"""
Orchestration layer for the lyrics-generation pipeline.

``LyricsGenerationService`` is the **single entry point** that Phase 3's
Celery task will call. It wires together:

  1. ``Transcriber.transcribe()``  → raw segments
  2. ``format_as_lrc()``           → LRC text
  3. ``compute_overall_confidence()`` → quality score

The service does **not** know about HTTP, Celery, S3, or shared volumes.
It only knows "here is a file on disk". The outer layer (storage / worker)
is responsible for getting the audio to a local path before calling this
service.
"""

from __future__ import annotations

import logging
from dataclasses import dataclass

from app.domain.lyrics.confidence import compute_overall_confidence
from app.domain.lyrics.lrc_formatter import format_as_lrc
from app.domain.lyrics.persian_normalizer import normalize_persian_text
from app.domain.lyrics.transcriber import Transcriber, TranscriptionError

logger = logging.getLogger(__name__)


# ── Custom exceptions ──────────────────────────────────────────────────


class LyricsGenerationError(Exception):
    """Raised when the full lyrics-generation pipeline fails.

    Wraps any internal exception (Whisper crash, I/O error, etc.) so callers
    see a single exception type.
    """


# ── Result dataclass ───────────────────────────────────────────────────


@dataclass
class LyricsGenerationResult:
    """The final output of the lyrics generation pipeline."""

    lrc_content: str = ""
    plain_text: str = ""
    confidence: float = 0.0
    detected_language: str = ""
    segment_count: int = 0


# ── Service ────────────────────────────────────────────────────────────


class LyricsGenerationService:
    """Orchestrate the full lyrics-generation pipeline.

    Usage::

        service = LyricsGenerationService(transcriber=transcriber)
        result = service.generate(audio_path="/path/to/track.mp3")
    """

    def __init__(self, transcriber: Transcriber) -> None:
        self._transcriber = transcriber

    def generate(
        self,
        audio_path: str,
        track_title: str | None = None,
        track_artist: str | None = None,
    ) -> LyricsGenerationResult:
        """Run the full pipeline: transcribe → normalise → format → score.

        Args:
            audio_path: Absolute path to a local audio file.
            track_title: Optional track title for LRC ``[ti]`` header.
            track_artist: Optional artist name for LRC ``[ar]`` header.

        Returns:
            A ``LyricsGenerationResult`` with LRC content, plain text,
            and quality score.

        Raises:
            LyricsGenerationError: If any step in the pipeline fails.
        """
        try:
            # ── Step 1: Transcribe ──────────────────────────────────
            logger.info("Starting transcription for %s", audio_path)
            transcription = self._transcriber.transcribe(audio_path)
            logger.info(
                "Transcription complete: %d segments, language=%s (p=%.2f)",
                len(transcription.segments),
                transcription.detected_language,
                transcription.language_probability,
            )

            if not transcription.segments:
                logger.warning("Transcription returned zero segments for %s", audio_path)
                return LyricsGenerationResult(
                    detected_language=transcription.detected_language,
                )

            # ── Step 2: Format as LRC ──────────────────────────────
            lrc_content = format_as_lrc(
                segments=transcription.segments,
                track_title=track_title,
                track_artist=track_artist,
            )

            # ── Step 3: Build plain text (normalized concatenation) ─
            plain_text = " ".join(
                normalize_persian_text(seg.text)
                for seg in transcription.segments
            )
            plain_text = plain_text.strip()

            # ── Step 4: Compute confidence ──────────────────────────
            confidence = compute_overall_confidence(transcription.segments)

            return LyricsGenerationResult(
                lrc_content=lrc_content,
                plain_text=plain_text,
                confidence=confidence,
                detected_language=transcription.detected_language,
                segment_count=len(transcription.segments),
            )

        except LyricsGenerationError:
            raise
        except TranscriptionError:
            raise
        except Exception as exc:
            logger.exception("Lyrics generation failed for %s", audio_path)
            raise LyricsGenerationError(
                f"Lyrics generation failed for '{audio_path}': {exc}"
            ) from exc
