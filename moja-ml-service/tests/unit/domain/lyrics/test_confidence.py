"""Tests for confidence scoring."""

from app.domain.lyrics.confidence import compute_overall_confidence
from app.domain.lyrics.transcriber import TranscriptSegment


def test_weighted_average_by_duration() -> None:
    """A long low-confidence segment pulls the average down more than a short confident one."""
    segments = [
        TranscriptSegment(start=0, end=1, text="a", confidence=1.0),  # 1s at 1.0
        TranscriptSegment(start=1, end=31, text="b", confidence=0.5),  # 30s at 0.5
    ]
    result = compute_overall_confidence(segments)
    # weighted = (1*1.0 + 30*0.5) / 31 = (1 + 15) / 31 = 16/31 ≈ 0.516
    assert 0.45 < result < 0.55, f"Expected ~0.516, got {result}"


def test_equal_duration_equal_confidence() -> None:
    """Equal-duration segments with equal confidence should average to that confidence."""
    segments = [
        TranscriptSegment(start=0, end=10, text="a", confidence=0.8),
        TranscriptSegment(start=10, end=20, text="b", confidence=0.8),
    ]
    result = compute_overall_confidence(segments)
    assert result == 0.8


def test_empty_segments_returns_zero() -> None:
    """Empty segment list should return 0.0."""
    assert compute_overall_confidence([]) == 0.0


def test_single_segment() -> None:
    """A single segment should return its confidence."""
    segments = [
        TranscriptSegment(start=0, end=10, text="test", confidence=0.75),
    ]
    result = compute_overall_confidence(segments)
    assert result == 0.75


def test_perfect_confidence() -> None:
    """All segments with confidence 1.0 should return 1.0."""
    segments = [
        TranscriptSegment(start=0, end=5, text="a", confidence=1.0),
        TranscriptSegment(start=5, end=10, text="b", confidence=1.0),
    ]
    result = compute_overall_confidence(segments)
    assert result == 1.0


def test_zero_confidence() -> None:
    """All segments with confidence 0.0 should return 0.0."""
    segments = [
        TranscriptSegment(start=0, end=5, text="a", confidence=0.0),
        TranscriptSegment(start=5, end=10, text="b", confidence=0.0),
    ]
    result = compute_overall_confidence(segments)
    assert result == 0.0


def test_zero_duration_segments() -> None:
    """Zero-duration segments should not cause division issues."""
    segments = [
        TranscriptSegment(start=0, end=0, text="a", confidence=1.0),
        TranscriptSegment(start=0, end=0, text="b", confidence=0.5),
    ]
    result = compute_overall_confidence(segments)
    # Falls back to simple average: (1.0 + 0.5) / 2 = 0.75
    assert result == 0.75


def test_single_zero_duration_segment() -> None:
    """A single zero-duration segment should return its confidence."""
    segments = [
        TranscriptSegment(start=0, end=0, text="a", confidence=0.9),
    ]
    result = compute_overall_confidence(segments)
    assert result == 0.9
