"""
Metadata feature vector for tracks.

This is the SECOND half of "similarity from both sources" — a separate,
much cheaper, function that structures the track's known metadata as a
queryable signature. Unlike the audio embedding (dense vector stored in
the vector index), these stay as normal Postgres columns and are used
for boosting/filtering during Phase 2's similarity query.

Domain layer — no HTTP or Celery knowledge.
"""

from dataclasses import dataclass

# ── Custom exceptions ─────────────────────────────────────────────────


class MetadataValidationError(ValueError):
    """Raised when metadata values fail validation."""


# ── Dataclasses ───────────────────────────────────────────────────────


@dataclass
class MetadataFeatures:
    """Structured metadata about a track for similarity boosting/filtering.

    Fields:
        genre_ids: The track's genre tag(s), as string IDs. Used for
            genre-based boosting in similarity queries.
        release_year: The track's release year, if known. Enables
            year-range filtering.
        bpm: Beats per minute, if already extracted. Used for tempo-based
            filtering (e.g. "find similar tracks but only upbeat ones").
    """

    genre_ids: list[str]
    release_year: int | None = None
    bpm: float | None = None


# ── Builder ───────────────────────────────────────────────────────────


def build_metadata_signature(
    genre_ids: list[str],
    release_year: int | None = None,
    bpm: float | None = None,
) -> dict:
    """Normalize and validate track metadata into a structured signature.

    This does NOT produce a dense vector — it returns a plain dict that
    will be stored as normal Postgres columns. Phase 2's similarity query
    will combine:
      - vector distance on audio embeddings (primary signal)
      - metadata-based boost/filter (secondary signal)
    e.g. "find audio-similar tracks, but boost ones that also share a
    genre tag, and allow filtering to a specific release-year range."

    Args:
        genre_ids: List of genre ID strings. Empty lists are allowed
            (the similarity query simply won't apply genre boosting).
        release_year: Four-digit year, or None. Must be between 1900
            and 2030 if provided.
        bpm: Beats per minute, or None. Must be between ``BPM_MIN``
            (40) and ``BPM_MAX`` (220) if provided, matching the
            sanity range used by ``tempo_extractor.extract_tempo()``.

    Returns:
        A dict with keys ``genre_ids``, ``release_year``, ``bpm``,
        with validated/normalized values.

    Raises:
        MetadataValidationError: If any value is outside acceptable range.
    """
    from app.domain.audio_features.embedding_extractor import BPM_MAX, BPM_MIN

    if not isinstance(genre_ids, list):
        raise MetadataValidationError("genre_ids must be a list")

    valid_genres = [g for g in genre_ids if isinstance(g, str) and g.strip()]
    if len(valid_genres) != len(genre_ids):
        logger = __import__("logging").getLogger(__name__)
        logger.warning("Ignoring %d invalid genre entries", len(genre_ids) - len(valid_genres))

    if release_year is not None:
        if not isinstance(release_year, int):
            raise MetadataValidationError(f"release_year must be an int, got {type(release_year)}")
        if release_year < 1900 or release_year > 2030:
            raise MetadataValidationError(
                f"release_year out of range [1900-2030]: {release_year}"
            )

    if bpm is not None:
        if not isinstance(bpm, (int, float)):
            raise MetadataValidationError(f"bpm must be numeric, got {type(bpm)}")
        if bpm < BPM_MIN or bpm > BPM_MAX:
            raise MetadataValidationError(
                f"bpm out of range [{BPM_MIN}-{BPM_MAX}]: {bpm}"
            )

    return {
        "genre_ids": valid_genres,
        "release_year": release_year,
        "bpm": bpm,
    }
