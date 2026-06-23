"""
Persian text normalisation for Whisper transcription output.

This is a **standalone, idempotent** function. It does NOT depend on any
other module in this package.

Transformation pipeline:
  1. Arabic → Persian character substitution.
  2. ZWNJ (Zero-Width Non-Joiner / نیم‌فاصله) normalisation via ``hazm``.
  3. Whisper filler-artifact removal (bracketed non-speech cues).
  4. Whitespace cleanup.

Design notes
------------
- hazm's default normaliser is tuned for formal news-style Persian text,
  which aggressively normalises punctuation and diacritics. Song lyrics
  frequently use repeated punctuation for stylistic emphasis (e.g.
  "آه………………" or "!!!"), so we disable or relax those steps.
- The Arabic-teh-marbuta (ة → ه) mapping is genuinely ambiguous: in some
  dialects it should stay as ة (e.g. family-name suffix), but in song
  lyrics the vast majority of occurrences are word-final and should become
  ه. We apply a simple heuristic: convert ة to ه only at the end of a word.
- Numerals are kept as-is (Latin digits). Persian digits are not introduced.
"""

from __future__ import annotations

import re

from hazm import Normalizer

# ── Character translation tables ───────────────────────────────────────

# Arabic yeh (ﻱ / ي) → Persian yeh (ی)
_ARABIC_YEH = "\u064a"  # Arabic yeh
_PERSIAN_YEH = "\u06cc"  # Persian yeh

# Arabic kaf (ﻙ / ك) → Persian kaf (ک)
_ARABIC_KAF = "\u0643"  # Arabic kaf
_PERSIAN_KAF = "\u06a9"  # Persian kaf

# Arabic teh marbuta (ة) → Persian heh (ه)  — word-final only
_ARABIC_TEH_MARBUTA = "\u0629"
_PERSIAN_HEH = "\u0647"

# ── Whisper filler-artifact patterns ───────────────────────────────────
# Whisper sometimes transcribes non-speech audio as bracketed text:
#   [موسیقی], [آهنگ], [خنده], [تشویق], etc.
_WHISPER_BRACKET_PATTERN = re.compile(r"\[[^\]]*\]")

# ── hazm normalizer (configured once at module load) ───────────────────
# hazm's ``Normalizer`` is tuned for formal news text by default.
# We adjust flags for song lyrics:
#
#   - ``persian_numbers=False``      → keep Latin digits as-is (no ۲۰۲۴)
#   - ``decrease_repeated_chars=False`` → preserve poetic repetition
#     (e.g. "آه………………!" should NOT become "آه!")
#   - ``remove_specials_chars=False`` → keep punctuation used in lyrics
#   - ``remove_diacritics=True``      → strip Arabic diacritics (common
#     in Whisper output but not part of written Persian)
#   - ``correct_spacing=True``        → insert ZWNJ in compound words
#   - ``persian_style=True``          → Arabic → Persian character mapping
#   - ``seperate_mi=True``            → separate می with ZWNJ
#   - ``unicodes_replacement=True``   → normalise various Unicode variants
_HAZM_NORMALIZER = Normalizer(
    persian_numbers=False,
    decrease_repeated_chars=False,
    remove_specials_chars=False,
    remove_diacritics=True,
    correct_spacing=True,
    persian_style=True,
    seperate_mi=True,
    unicodes_replacement=True,
)


def normalize_persian_text(text: str) -> str:
    """Clean raw Whisper Persian transcription output.

    The function is **idempotent**::

        normalize_persian_text(normalize_persian_text(x)) == normalize_persian_text(x)
    """
    if not text:
        return ""

    # 1. Arabic → Persian character normalisation
    text = _fix_arabic_chars(text)

    # 2. ZWNJ normalisation via hazm
    text = _HAZM_NORMALIZER.normalize(text)

    # 3. Remove Whisper bracketed filler artifacts
    text = _WHISPER_BRACKET_PATTERN.sub("", text)

    # 4. Whitespace cleanup
    text = _collapse_whitespace(text)
    text = text.strip()

    return text


# ── Internal helpers ───────────────────────────────────────────────────


def _fix_arabic_chars(text: str) -> str:
    """Substitute Arabic characters with their Persian equivalents."""
    # Replace Arabic yeh with Persian yeh
    text = text.replace(_ARABIC_YEH, _PERSIAN_YEH)
    # Replace Arabic kaf with Persian kaf
    text = text.replace(_ARABIC_KAF, _PERSIAN_KAF)
    # Replace word-final Arabic teh marbuta with Persian heh
    # Heuristic: a word-final ة before a word boundary (space, end-of-string,
    # or punctuation) should become ه.
    # We use a lookbehind for a Persian/Arabic letter and lookahead for a
    # word boundary.
    text = re.sub(
        rf"(?<=[\u0600-\u06FF]){_ARABIC_TEH_MARBUTA}(?=\s|$|[\.\,\!\?\;\:])",
        _PERSIAN_HEH,
        text,
    )
    # Also handle the case where ة is at the very end of the string after a letter
    text = re.sub(
        rf"(?<=[\u0600-\u06FF]){_ARABIC_TEH_MARBUTA}$",
        _PERSIAN_HEH,
        text,
    )
    return text


def _collapse_whitespace(text: str) -> str:
    """Collapse multiple consecutive whitespace chars into a single space."""
    return re.sub(r"\s+", " ", text)
