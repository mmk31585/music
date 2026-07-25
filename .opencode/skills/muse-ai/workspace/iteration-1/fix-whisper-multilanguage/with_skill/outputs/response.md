# Fix: Whisper Multi-Language Lyrics Transcription

## Root Cause

The garbled Turkish output is caused by a single architectural bug in `moja-ml-service/app/domain/lyrics/service.py:114-119`:

```python
plain_text = " ".join(
    normalize_persian_text(seg.text)
    for seg in transcription.segments
)
```

The `normalize_persian_text()` function is applied **unconditionally** to every segment's text, regardless of the detected language. This means Turkish songs — which use Latin script with characters like `ç`, `ğ`, `ı`, `ö`, `ş`, `ü` — get passed through `hazm`'s Persian normalizer, which:

1. **Mangles Turkish Latin characters** — `hazm`'s `unicodes_replacement` and `persian_style` transform apply Arabic/Persian character substitutions to Turkish text, corrupting it.
2. **Mis-applies ZWNJ insertion** — The `correct_spacing=True` and `seperate_mi=True` flags insert zero-width non-joiners around patterns that match Persian morphology but don't exist in Turkish.
3. **Strips content via `_WHISPER_BRACKET_PATTERN`** — The regex `\[[^\]]*\]` is designed to remove Whisper artifacts like `[موسیقی]`, but it also strips legitimate bracketed content or sung cues in any language. Turkish Whisper output may contain different artifact patterns (e.g. `[Müzik]`, `[Alkış]`).

Currently, the `transcriber.py` correctly auto-detects language per track (`language=None` at line 120 triggers Whisper's language detection), and the `TranscriptionResult.detected_language` field is populated. But `service.py` ignores it when choosing normalization.

Additionally, `config.py:30` sets `whisper_language_hint = "fa"`, but this is passed only to the `Transcriber.__init__` and never actually used (`transcriber.py:120` hardcodes `language=None`). This field is dead config that misleads future engineers.

---

## Implementation Plan

### 1. Make `normalize_persian_text()` language-aware

**File:** `moja-ml-service/app/domain/lyrics/persian_normalizer.py`

Add an optional `language` parameter that gates the Persian normalization to only apply when the detected language is Persian:

```python
# At module level, define the languages that get Persian normalization:
PERSIAN_LANGUAGE_CODES = {"fa", "pes", "prs"}

def normalize_persian_text(text: str, language: str | None = None) -> str:
    if not text:
        return ""

    # Non-Persian text: only strip Whisper artifacts, skip hazm
    if language and language.lower() not in PERSIAN_LANGUAGE_CODES:
        text = _WHISPER_BRACKET_PATTERN.sub("", text)
        text = _collapse_whitespace(text)
        return text.strip()

    # Persian text: full pipeline (existing logic)
    text = _fix_arabic_chars(text)
    text = _HAZM_NORMALIZER.normalize(text)
    text = _WHISPER_BRACKET_PATTERN.sub("", text)
    text = _collapse_whitespace(text)
    return text.strip()
```

Design rationale:
- For Turkish, English, Arabic, Kurdish, etc., we do **not** apply `hazm` or `_fix_arabic_chars` — those transforms are Persian-specific and destructive to Turkish Latin text.
- We still strip Whisper bracket artifacts universally (they occur in all languages).
- We still collapse whitespace universally.
- The function remains **idempotent** for Persian text; for non-Persian text, the only non-idempotent step is bracket removal, which is idempotent by definition (after one pass no brackets remain).

### 2. Update `service.py` to pass detected language

**File:** `moja-ml-service/app/domain/lyrics/service.py`

Change the plain_text construction block (lines 114-119) to:

```python
plain_text = " ".join(
    normalize_persian_text(seg.text, language=transcription.detected_language)
    for seg in transcription.segments
)
```

Also, do the same in `lrc_formatter.py:52` — the `format_as_lrc()` function also calls `normalize_persian_text(seg.text)` for each segment. Since `format_as_lrc()` doesn't have access to the detected language, the cleanest approach is to thread the language through. Optionally, we can avoid normalization inside `lrc_formatter.py` altogether and do it in `service.py` pre-formatting. The current architecture has normalization happening twice (once in `lrc_formatter.py:52`, once in `service.py:116`), which is a hidden performance cost.

A cleaner refactor: **normalize once** in `service.py`, and pass pre-normalized segments to `format_as_lrc()`.

**Update `lrc_formatter.py`:**

Remove the `normalize_persian_text` import and call — accept pre-normalized text:

```python
def format_as_lrc(
    segments: list[TranscriptSegment],
    track_title: str | None = None,
    track_artist: str | None = None,
) -> str:
    ...
    for seg in sorted_segments:
        raw_text = seg.text  # already normalized by caller
        if not raw_text:
            continue
        timestamp = _format_lrc_timestamp(seg.start)
        lines.append(f"{timestamp}{raw_text}")
    return "\n".join(lines)
```

**Update `service.py`** to normalize segments before passing to `format_as_lrc()`:

```python
# Step 1.5: Normalize segment text based on detected language
lang = transcription.detected_language
for seg in transcription.segments:
    seg.text = normalize_persian_text(seg.text, language=lang)

# Step 2: Format as LRC (now receives pre-normalized text)
lrc_content = format_as_lrc(
    segments=transcription.segments,
    track_title=track_title,
    track_artist=track_artist,
)

# Step 3: Build plain text (normalized concatenation, idempotent on already-clean text)
plain_text = " ".join(seg.text for seg in transcription.segments)
```

This makes normalization happen exactly **once** per segment and keeps the LRC and plain-text outputs consistent.

### 3. Clean up dead config

**File:** `moja-ml-service/app/config.py`

Remove or rename `whisper_language_hint` — it's never consumed by the transcriber pipeline:

```python
# Remove this line (line 30):
whisper_language_hint: str = "fa"
```

**File:** `moja-ml-service/app/domain/lyrics/transcriber.py`

Remove the `language_hint` parameter from `__init__` and `get_transcriber()`, since the transcribe method always passes `language=None` for auto-detection. The field `self._language_hint` at line 92 and its fallback at line 132 (`self._language_hint or "unknown"`) is dead code.

Alternatively, if you want to keep the option for forcing a language (useful for testing or for tracks with known language metadata), make the language overridable at the `transcribe()` method level rather than at construction:

```python
def transcribe(self, audio_path: str, language: str | None = None) -> TranscriptionResult:
    ...
    segments_gen, info = self._model.transcribe(
        audio_path,
        language=language,  # None = auto-detect; str = force
        word_timestamps=True,
        vad_filter=True,
        vad_parameters=dict(min_silence_duration_ms=500),
    )
```

---

## Tests to Add / Update

### `tests/unit/domain/lyrics/test_persian_normalizer.py`

```python
# Language-aware behavior

def test_non_persian_text_passes_through() -> None:
    """Turkish text should NOT be modified by Persian normalization."""
    turkish = "Gönlümün eşiğinde bir başka sen varsın"
    result = normalize_persian_text(turkish, language="tr")
    assert result == turkish

def test_english_text_passes_through() -> None:
    """English text should NOT be modified."""
    english = "Hello world, this is a song."
    result = normalize_persian_text(english, language="en")
    assert result == english

def test_non_persian_removes_brackets() -> None:
    """Whisper artifacts should still be removed for non-Persian."""
    text = "[Music] Hello world"
    result = normalize_persian_text(text, language="en")
    assert result == "Hello world"

def test_persian_no_language_fallback_works() -> None:
    """When language is None, assume Persian (backward compat)."""
    result = normalize_persian_text("علي", language=None)
    assert result == "علی"

def test_arabic_text_not_normalized_as_persian() -> None:
    """Arabic (ar) should skip hazm normalizer."""
    arabic = "يقول كلام"
    result = normalize_persian_text(arabic, language="ar")
    # hazm not applied; bracket removal + whitespace only
    assert "ي" in result  # Arabic yeh kept intact

def test_turkish_special_chars_preserved() -> None:
    """Turkish chars ç, ğ, ı, ö, ş, ü should survive."""
    turkish = "Çiğdem şekerli bir öyküydü"
    result = normalize_persian_text(turkish, language="tr")
    assert "ç" in result
    assert "ğ" in result
    assert "ı" in result
    assert "ö" in result
    assert "ş" in result
    assert "ü" in result
```

### `tests/unit/domain/lyrics/test_service.py`

```python
def test_service_respects_detected_language(mocker) -> None:
    """Non-Persian segments should NOT be Persian-normalized."""
    mock_transcriber = mocker.Mock(spec=Transcriber)
    mock_transcriber.transcribe.return_value = TranscriptionResult(
        segments=[
            TranscriptSegment(start=1.0, end=2.0, text="Gönlümün eşiğinde", confidence=0.95),
            TranscriptSegment(start=3.0, end=4.0, text="bir başka sen varsın", confidence=0.93),
        ],
        detected_language="tr",
        language_probability=0.97,
        duration_seconds=180.0,
    )

    service = LyricsGenerationService(transcriber=mock_transcriber)
    result = service.generate(audio_path="/fake/path.mp3")

    assert "Gönlümün eşiğinde" in result.lrc_content
    assert "bir başka sen varsın" in result.plain_text
    assert "ğ" in result.lrc_content  # Turkish chars preserved
    assert "ş" in result.plain_text
    assert result.detected_language == "tr"
```

---

## Files to Modify (Summary)

| File | Change |
|------|--------|
| `moja-ml-service/app/domain/lyrics/persian_normalizer.py` | Add `language` parameter; gate hazm + char fix behind Persian-only; strip brackets for all |
| `moja-ml-service/app/domain/lyrics/service.py` | Pass `detected_language` to normalization; move normalization before LRC formatting; deduplicate normalization call |
| `moja-ml-service/app/domain/lyrics/lrc_formatter.py` | Remove `normalize_persian_text` call (accept pre-normalized text instead) |
| `moja-ml-service/app/domain/lyrics/transcriber.py` | Remove dead `language_hint` parameter; make language overridable at `transcribe()` level |
| `moja-ml-service/app/config.py` | Remove `whisper_language_hint` config field |
| `moja-ml-service/tests/unit/domain/lyrics/test_persian_normalizer.py` | Add multilingual test cases |
| `moja-ml-service/tests/unit/domain/lyrics/test_service.py` | Add test verifying language-aware normalization |

---

## Pipeline Flow After Fix

```
audio file → Transcriber.transcribe()
             ↓
       detected_language = "tr" (or "en", "fa", etc.)
             ↓
       normalize_persian_text(seg.text, language="tr")
         → strips [brackets], collapses whitespace → NO hazm
             ↓
       format_as_lrc(pre_normalized_segments)
             ↓
       compute_overall_confidence()
             ↓
       LyricsGenerationResult (Turkish text intact)
```

For Persian songs (`language="fa"`), the existing pipeline runs unchanged — hazm normalization, Arabic→Persian character substitution, ZWNJ insertion, etc. The change is purely additive: non-Persian text bypasses the Persian-specific transforms while still getting basic cleanup.
