"""Tests for the metadata signature builder."""

from __future__ import annotations

import pytest

from app.domain.audio_features.metadata_features import (
    MetadataValidationError,
    build_metadata_signature,
)


def test_builds_valid_signature() -> None:
    """Valid inputs produce a correctly structured dict."""
    result = build_metadata_signature(
        genre_ids=["rock", "indie"],
        release_year=2023,
        bpm=120.0,
    )
    assert result == {
        "genre_ids": ["rock", "indie"],
        "release_year": 2023,
        "bpm": 120.0,
    }


def test_allows_none_fields() -> None:
    """Optional None fields should be passed through as None."""
    result = build_metadata_signature(genre_ids=["pop"])
    assert result["genre_ids"] == ["pop"]
    assert result["release_year"] is None
    assert result["bpm"] is None


def test_filters_invalid_genres() -> None:
    """Empty or non-string genre entries should be filtered."""
    result = build_metadata_signature(
        genre_ids=["jazz", "", "   ", None, 123],
    )
    assert result["genre_ids"] == ["jazz"]


def test_rejects_out_of_range_year() -> None:
    """Years outside 1900-2030 should be rejected."""
    with pytest.raises(MetadataValidationError, match="release_year out of range"):
        build_metadata_signature(genre_ids=[], release_year=1899)

    with pytest.raises(MetadataValidationError, match="release_year out of range"):
        build_metadata_signature(genre_ids=[], release_year=2031)


def test_rejects_non_int_year() -> None:
    """Non-integer year should be rejected."""
    with pytest.raises(MetadataValidationError, match="release_year must be an int"):
        build_metadata_signature(genre_ids=[], release_year="2023")


def test_rejects_out_of_range_bpm() -> None:
    """BPM outside 40-220 should be rejected."""
    with pytest.raises(MetadataValidationError, match="bpm out of range"):
        build_metadata_signature(genre_ids=[], bpm=10.0)

    with pytest.raises(MetadataValidationError, match="bpm out of range"):
        build_metadata_signature(genre_ids=[], bpm=300.0)


def test_rejects_non_numeric_bpm() -> None:
    """Non-numeric BPM should be rejected."""
    with pytest.raises(MetadataValidationError, match="bpm must be numeric"):
        build_metadata_signature(genre_ids=[], bpm="fast")


def test_rejects_non_list_genre_ids() -> None:
    """genre_ids must be a list."""
    with pytest.raises(MetadataValidationError, match="genre_ids must be a list"):
        build_metadata_signature(genre_ids="jazz")  # type: ignore[arg-type]
