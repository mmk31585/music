"""
Confidence scoring for whole-transcription quality assessment.

The computed score will be surfaced to the Moja admin (via the Smart Upload
Review UI) so operators know whether to trust auto-generated lyrics or
review them manually.

This module is a **pure function** — no I/O, no side effects.
"""

from __future__ import annotations

from app.domain.lyrics.transcriber import TranscriptSegment


def compute_overall_confidence(segments: list[TranscriptSegment]) -> float:
    """Aggregate per-segment confidence into a single 0.0–1.0 score.

    The aggregation is a **duration-weighted average**: a long, confident
    segment contributes more to the final score than a short, confident one.

    Args:
        segments: Transcribed segments with per-segment confidence scores.

    Returns:
        A float in [0.0, 1.0] representing overall transcription quality.
        Returns ``0.0`` for an empty segment list.
    """
    if not segments:
        return 0.0

    total_weight = 0.0
    weighted_sum = 0.0

    for seg in segments:
        duration = seg.end - seg.start
        if duration < 0:
            # Defensive: guard against negative-duration segments
            duration = 0.0
        weighted_sum += seg.confidence * duration
        total_weight += duration

    if total_weight <= 0.0:
        # All segments had zero duration — fall back to simple average
        return sum(seg.confidence for seg in segments) / len(segments)

    return weighted_sum / total_weight
