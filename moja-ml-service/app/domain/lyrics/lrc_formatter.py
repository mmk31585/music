"""
LRC (LyRiCs) format converter.

Converts a list of ``TranscriptSegment`` objects into a standard ``.lrc``
file as a string. The format is::

    [ti:Track Title]
    [ar:Artist Name]
    [by:Moja AI Lyrics]
    [mm:ss.xx]اولین خط متن
    [mm:ss.xx]خط دوم متن

This is a **pure function** — no I/O, no side effects. It only depends on
``TranscriptSegment`` and ``normalize_persian_text``.
"""

from __future__ import annotations

from app.domain.lyrics.persian_normalizer import normalize_persian_text
from app.domain.lyrics.transcriber import TranscriptSegment


def format_as_lrc(
    segments: list[TranscriptSegment],
    track_title: str | None = None,
    track_artist: str | None = None,
) -> str:
    """Convert transcript segments into standard ``.lrc`` format.

    Args:
        segments: Transcribed segments (will be sorted by ``start`` time).
        track_title: Optional track title for the ``[ti]`` header tag.
        track_artist: Optional artist name for the ``[ar]`` header tag.

    Returns:
        A complete LRC file as a string (lines joined by ``\\n``).
    """
    lines: list[str] = []

    # ── Optional header tags ────────────────────────────────────────
    if track_title is not None:
        lines.append(f"[ti:{track_title}]")
    if track_artist is not None:
        lines.append(f"[ar:{track_artist}]")
    lines.append("[by:Moja AI Lyrics]")

    # ── Sort segments by start time (defensive) ─────────────────────
    sorted_segments = sorted(segments, key=lambda s: s.start)

    # ── Format each segment ─────────────────────────────────────────
    for seg in sorted_segments:
        raw_text = normalize_persian_text(seg.text)

        # Skip lines that become empty after normalisation (e.g. filler artifacts)
        if not raw_text:
            continue

        timestamp = _format_lrc_timestamp(seg.start)
        lines.append(f"{timestamp}{raw_text}")

    return "\n".join(lines)


# ── Internal helpers ───────────────────────────────────────────────────


def _format_lrc_timestamp(seconds: float) -> str:
    """Format a floating-point time in seconds to ``[mm:ss.xx]``.

    Guarantees:
        - Two-digit minutes (padded with leading zero if needed).
        - Two-digit seconds (padded).
        - Two-digit centiseconds (truncated, not rounded).
    """
    # Guard against negative timestamps
    if seconds < 0:
        seconds = 0.0

    # Use an integer-based approach to avoid floating-point drift.
    # Multiply by 100 to work in centiseconds, then decompose.
    total_cs = int(seconds * 100)
    centiseconds = total_cs % 100
    total_secs = total_cs // 100
    minutes = total_secs // 60
    secs = total_secs % 60

    return f"[{minutes:02d}:{secs:02d}.{centiseconds:02d}]"
