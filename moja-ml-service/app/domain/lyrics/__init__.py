"""Lyrics domain — transcription, normalisation, LRC formatting, scoring."""

from app.domain.lyrics.confidence import compute_overall_confidence
from app.domain.lyrics.lrc_formatter import format_as_lrc
from app.domain.lyrics.persian_normalizer import normalize_persian_text
from app.domain.lyrics.service import (
    LyricsGenerationError,
    LyricsGenerationResult,
    LyricsGenerationService,
)
from app.domain.lyrics.transcriber import (
    Transcriber,
    TranscriptionError,
    TranscriptionResult,
    TranscriptSegment,
    get_transcriber,
)

__all__ = [
    "Transcriber",
    "TranscriptSegment",
    "TranscriptionError",
    "TranscriptionResult",
    "get_transcriber",
    "LyricsGenerationError",
    "LyricsGenerationResult",
    "LyricsGenerationService",
    "normalize_persian_text",
    "format_as_lrc",
    "compute_overall_confidence",
]
