"""Tests for the LRC formatter."""

from app.domain.lyrics.lrc_formatter import format_as_lrc
from app.domain.lyrics.transcriber import TranscriptSegment


def test_basic_lrc_format() -> None:
    """A single segment should produce one timestamped line."""
    segments = [
        TranscriptSegment(start=12.34, end=15.0, text="سلام دنیا", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    assert "[00:12.34]سلام دنیا" in result


def test_multiple_segments() -> None:
    """Multiple segments should each produce a timestamped line."""
    segments = [
        TranscriptSegment(start=1.0, end=2.0, text="خط اول", confidence=0.9),
        TranscriptSegment(start=3.0, end=4.0, text="خط دوم", confidence=0.8),
    ]
    result = format_as_lrc(segments)
    assert "[00:01.00]خط اول" in result
    assert "[00:03.00]خط دوم" in result


def test_includes_header_tags_when_provided() -> None:
    """Header tags should appear when track title and artist are given."""
    result = format_as_lrc([], track_title="آهنگ تست", track_artist="خواننده تست")
    assert "[ti:آهنگ تست]" in result
    assert "[ar:خواننده تست]" in result
    assert "[by:Moja AI Lyrics]" in result


def test_omits_ti_tag_when_title_not_provided() -> None:
    """[ti] tag should be omitted when track_title is None."""
    result = format_as_lrc([], track_artist="خواننده")
    assert "[ti:" not in result
    assert "[ar:خواننده]" in result


def test_omits_ar_tag_when_artist_not_provided() -> None:
    """[ar] tag should be omitted when track_artist is None."""
    result = format_as_lrc([], track_title="آهنگ")
    assert "[ar:" not in result
    assert "[ti:آهنگ]" in result


def test_omits_all_header_tags_when_none_provided() -> None:
    """Header tags should all be omitted when neither title nor artist."""
    result = format_as_lrc([])
    assert "[ti:" not in result
    assert "[ar:" not in result
    # [by] should always be included
    assert "[by:Moja AI Lyrics]" in result


def test_skips_empty_lines_after_normalization() -> None:
    """Segments that become empty after normalization should be skipped."""
    segments = [
        TranscriptSegment(start=1.0, end=2.0, text="[موسیقی]", confidence=0.5),
        TranscriptSegment(start=3.0, end=4.0, text="سلام", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    lines = [line for line in result.split("\n") if line.startswith("[")]
    # Only one timestamp line (the "سلام" one), plus the [by] header
    timestamp_lines = [line for line in lines if line.startswith("[0")]
    assert len(timestamp_lines) == 1, f"Expected 1 timestamp line, got {timestamp_lines}"


def test_sorts_segments_by_start_time() -> None:
    """Out-of-order segments should be sorted by start time."""
    segments = [
        TranscriptSegment(start=5.0, end=6.0, text="دوم", confidence=0.9),
        TranscriptSegment(start=1.0, end=2.0, text="اول", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    first_line_idx = result.index("اول")
    second_line_idx = result.index("دوم")
    assert first_line_idx < second_line_idx, "Segments not sorted by start time"


def test_timestamp_centiseconds() -> None:
    """Timestamps should have two-digit centiseconds."""
    segments = [
        TranscriptSegment(start=12.345, end=15.0, text="تست", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    # 12.345 → 12 seconds + 34 centiseconds (truncated, not rounded)
    assert "[00:12.34]" in result


def test_timestamp_negative_start() -> None:
    """Negative start times should be clamped to zero."""
    segments = [
        TranscriptSegment(start=-5.0, end=2.0, text="تست", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    assert "[00:00.00]" in result


def test_timestamp_long_duration() -> None:
    """Timestamps over a minute should show correct minutes."""
    segments = [
        TranscriptSegment(start=125.0, end=130.0, text="تست", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    # 125 seconds = 2 minutes 5 seconds
    assert "[02:05.00]" in result


def test_normalizes_persian_text_in_output() -> None:
    """Arabic characters in segment text should be normalized in output."""
    segments = [
        TranscriptSegment(start=1.0, end=2.0, text="علي", confidence=0.9),
    ]
    result = format_as_lrc(segments)
    # "علي" (with Arabic yeh) should become "علی" (with Persian yeh)
    assert "علی" in result
    assert "علي" not in result


def test_skips_all_artifact_segments() -> None:
    """Only artifact segments should produce no timestamp lines."""
    segments = [
        TranscriptSegment(start=1.0, end=2.0, text="[موسیقی]", confidence=0.5),
        TranscriptSegment(start=3.0, end=4.0, text="[آهنگ]", confidence=0.4),
    ]
    result = format_as_lrc(segments)
    timestamp_lines = [line for line in result.split("\n") if line.startswith("[0")]
    assert len(timestamp_lines) == 0
