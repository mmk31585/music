# Fix: Whisper multi-language transcription garbling non-Persian output

## Root cause

The Persian normalizer (`persian_normalizer.py`) is applied **unconditionally** to all
transcription output, regardless of the detected language. This happens in two places:

- `service.py:115-118` — plain-text construction calls `normalize_persian_text(seg.text)` for every segment.
- `lrc_formatter.py:52` — LRC formatting calls `normalize_persian_text(seg.text)` for every segment.

When the audio is Turkish (or English, Arabic, etc.), the hazm normalizer:

1. **Inserts ZWNJ** (`\u200c`) into Latin words via `correct_spacing=True` and
   `seperate_mi=True` — e.g. "biliyorum" becomes "bili\u200cyorum".
2. **Strips diacritics** via `remove_diacritics=True` — changes Turkish `ş`, `ç`, `ğ`, `ü`, `ö`, `ı`
   to their ASCII base forms.
3. **Mangles Unicode** via `unicodes_replacement=True` — applies Persian-specific
   character substitutions to Latin ranges.

The transcriber **already detects language** correctly (`detected_language` in
`TranscriptionResult`). The pipeline is simply ignoring it.

---

## Fix strategy

Thread the detected language through the pipeline and **skip Persian-specific
normalization for non-Arabic-script languages**.

### Step 1 — Make `normalize_persian_text` language-aware

Add an optional `language` parameter. When the language is **not** one of the
Arabic-script languages (`fa`, `ar`, `ur`, `ps`, `sd`, `ku`, etc.), skip the
Persian-specific steps (hazm, Arabic→Persian char mapping) and only do universal
cleanup (bracket artifact removal + whitespace collapse).

```python
# Arabic-script BCP-47 codes that should get Persian normalization
_PERSO_ARABIC_SCRIPTS = {"fa", "ar", "ur", "ps", "sd", "ku", "ckb", "bal", "glk", "mzn"}


def normalize_persian_text(text: str, language: str = "fa") -> str:
    if not text:
        return ""

    # Universal cleanup: bracket artifacts + whitespace
    text = _WHISPER_BRACKET_PATTERN.sub("", text)
    text = _collapse_whitespace(text)
    text = text.strip()

    # Persian-specific: skip for non-Arabic-script languages
    if language and language.split("-")[0].lower() not in _PERSO_ARABIC_SCRIPTS:
        return text

    # 1. Arabic → Persian character normalisation
    text = _fix_arabic_chars(text)
    # 2. ZWNJ via hazm
    text = _HAZM_NORMALIZER.normalize(text)
    # 3. Whitespace cleanup (again after hazm may introduce extra spaces)
    text = _collapse_whitespace(text)
    text = text.strip()

    return text
```

> **Note on bracket removal ordering**: Moving it before the language guard means
> `[موسیقی]` artifacts are stripped even for Persian tracks in the early-return
> path. The haghe 2x application for Persian is harmless (idempotent), but you
> can also keep it only after hazm if you want — the difference is negligible.

### Step 2 — Thread language through `lrc_formatter.py`

The `format_as_lrc` function currently has no language parameter. Add one and
pass it to `normalize_persian_text`:

```python
def format_as_lrc(
    segments: list[TranscriptSegment],
    track_title: str | None = None,
    track_artist: str | None = None,
    language: str = "fa",
) -> str:
    ...
    for seg in sorted_segments:
        raw_text = normalize_persian_text(seg.text, language=language)
        ...
```

### Step 3 — Thread language through `service.py`

Pass the detected language to both the LRC formatter and plain-text builder:

```python
# Step 2: Format as LRC
lrc_content = format_as_lrc(
    segments=transcription.segments,
    track_title=track_title,
    track_artist=track_artist,
    language=transcription.detected_language,
)

# Step 3: Build plain text (normalized concatenation)
plain_text = " ".join(
    normalize_persian_text(seg.text, language=transcription.detected_language)
    for seg in transcription.segments
)
```

No changes needed in `transcriber.py` — it already returns `detected_language`.

---

## Turkish-specific Whisper tuning (optional but recommended)

faster-whisper's `large-v3` model handles Turkish well, but `tiny`/`base`/`small`
can struggle. Consider these additional improvements:

1. **Pass `language="tr"` for Turkish tracks** if you have metadata (language tag
   from the Go backend) before calling `transcribe()`. This forces Whisper to
   use the Turkish decoder rather than auto-detect, which eliminates the risk of
   the model switching mid-recording.

   ```python
   language = None if not known_language else known_language
   ```

   The transcriber already has `language_hint` support — wire it to the Celery
   task's job metadata if available.

2. **Enable `beam_size=5`** for better accuracy on non-English languages at the
   cost of ~30% slower inference:

   ```python
   self._model.transcribe(audio_path, beam_size=5, ...)
   ```

3. **Increase `vad_parameters` min_silence_duration** for Turkish (which has
   different phoneme timing than Persian). The current `500ms` is fine for both.

---

## Edge cases handled

| Case | Before | After |
|---|---|---|
| Turkish song auto-detected as `tr` | ZWNJ everywhere, diacritics stripped | Text passed through cleanly |
| Persian song with Arabic yeh/kaf | Normalized to Persian | Same (no change) |
| Arabic song (Arabic script, not Persian) | Over-normalized by hazm Persian rules | Gets Persian-style normalization — acceptable since Arabic speakers read the result |
| English song in `en` | hazm applied to English | Only bracket removal + whitespace |
| Song with mixed languages (e.g. Turkish verse + English chorus) | Language = whichever Whisper detected first; some segments garbled | Dominant language wins; best-effort |
| Bracketed artifacts in Turkish | Removed by hazm's bracket pattern (working) | Moved to universal path — still removed |

---

## Test coverage to add

1. `test_normalize_turkish_preserves_chars()` — Turkish lyric with `ş`, `ç`, `ğ`
   passed with `language="tr"` should return identical text (minus brackets).
2. `test_normalize_english_no_zwnj()` — English text with `language="en"` should
   contain no ZWNJ characters.
3. `test_normalize_persian_still_works()` — Existing Persian tests pass when
   `language="fa"` is passed explicitly.
4. `test_format_lrc_passes_language()` — `format_as_lrc` with `language="tr"`
   doesn't mangle Turkish in LRC output.
5. `test_service_threads_language()` — Mock a `TranscriptionResult` with
   `detected_language="tr"` and assert plain_text is not Persian-normalized.
