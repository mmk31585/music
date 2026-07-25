"""Tests for the Persian text normalizer.

These tests use real Persian text fixtures — not mocks — to verify that
``normalize_persian_text`` correctly handles the transformations described
in the module docstring.
"""

from app.domain.lyrics.persian_normalizer import normalize_persian_text


def test_normalize_arabic_yeh_to_persian() -> None:
    """Arabic yeh (ي) should become Persian yeh (ی)."""
    # Arabic ي vs Persian ی
    result = normalize_persian_text("علي")
    assert result == "علی"


def test_normalize_arabic_kaf_to_persian() -> None:
    """Arabic kaf (ك) should become Persian kaf (ک)."""
    result = normalize_persian_text("كتاب")
    assert result == "کتاب"


def test_normalize_arabic_teh_marbuta_to_heh() -> None:
    """Word-final Arabic teh marbuta (ة) should become Persian heh (ه)."""
    result = normalize_persian_text("فاطمة")
    assert result == "فاطمه"


def test_zwnj_insertion_for_compound_verbs() -> None:
    """Compound verbs should get ZWNJ (نیم‌فاصله) inserted."""
    result = normalize_persian_text("میخواهم برم")
    # The ZWNJ character (U+200C) should be present between می and خواهم
    assert "\u200c" in result, f"Expected ZWNJ in normalized text: {result!r}"


def test_zwnj_in_mikonam() -> None:
    """می‌کنم should have ZWNJ after می."""
    result = normalize_persian_text("میکنم")
    assert "\u200c" in result, f"Expected ZWNJ in normalized text: {result!r}"


def test_removes_whisper_music_artifacts() -> None:
    """Whisper bracketed artifacts like [موسیقی] should be removed."""
    result = normalize_persian_text("[موسیقی] سلام دنیا")
    assert result == "سلام دنیا"


def test_removes_whisper_laugh_artifacts() -> None:
    """Whisper [خنده] artifact should be removed."""
    result = normalize_persian_text("سلام [خنده] دنیا")
    assert result == "سلام دنیا"


def test_removes_whisper_applause_artifacts() -> None:
    """Whisper [تشویق] artifact should be removed."""
    result = normalize_persian_text("[تشویق] آهنگ جدید")
    assert result == "آهنگ جدید"


def test_removes_multiple_bracketed_artifacts() -> None:
    """Multiple bracketed artifacts should all be removed."""
    result = normalize_persian_text("[موسیقی] [آهنگ] متن آهنگ")
    assert result == "متن آهنگ"


def test_idempotent() -> None:
    """Applying the normalizer twice should produce the same result."""
    text = "علي كتاب ميخواند"
    once = normalize_persian_text(text)
    twice = normalize_persian_text(once)
    assert once == twice, f"Idempotency failed: {once!r} != {twice!r}"


def test_idempotent_with_artifacts() -> None:
    """Artifact removal should be idempotent (already clean)."""
    text = normalize_persian_text("[موسیقی] سلام دنیا")
    twice = normalize_persian_text(text)
    assert text == twice


def test_collapses_whitespace() -> None:
    """Multiple spaces should be collapsed into one."""
    result = normalize_persian_text("سلام    دنیا")
    assert result == "سلام دنیا"


def test_collapses_tabs_and_newlines() -> None:
    """Tabs and newlines should be collapsed into single spaces."""
    result = normalize_persian_text("سلام\t\tدنیا\nخوبی")
    assert result == "سلام دنیا خوبی"


def test_empty_string() -> None:
    """Empty string input should return empty string."""
    assert normalize_persian_text("") == ""


def test_preserves_actual_numbers() -> None:
    """Latin digits (e.g. years) should be preserved as-is."""
    result = normalize_persian_text("سال 2024")
    assert "2024" in result


def test_preserves_punctuation() -> None:
    """Stylistic punctuation should be preserved for lyrics."""
    result = normalize_persian_text("آه………………!")
    # Ellipsis-like sequences and exclamation marks should survive
    assert "آه" in result
    assert "!" in result


def test_strips_leading_trailing_whitespace() -> None:
    """Leading and trailing whitespace should be stripped."""
    result = normalize_persian_text("  سلام دنیا  ")
    assert result == "سلام دنیا"


def test_mixed_arabic_persian() -> None:
    """Mixed Arabic/Persian characters should all become Persian."""
    result = normalize_persian_text("يقول كلام")
    assert "ي" not in result  # No Arabic yeh
    assert "ك" not in result  # No Arabic kaf
    assert "ی" in result  # Has Persian yeh
    assert "ک" in result  # Has Persian kaf
