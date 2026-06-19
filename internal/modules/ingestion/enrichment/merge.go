package enrichment

import "strings"

func MergeResults(mb *MusicBrainzResult, lfm *LastFMResult, spot *SpotifyResult, lrc *LRCLibResult) EnrichmentResult {
	result := EnrichmentResult{
		MusicBrainz: mb,
		LastFM:      lfm,
		Spotify:     spot,
		Attempted:   true,
	}

	result.Suggestions = mergeTitle(mb, spot)
	result.Suggestions = append(result.Suggestions, mergeArtist(mb)...)
	result.Suggestions = append(result.Suggestions, mergeAlbum(mb)...)
	result.Suggestions = append(result.Suggestions, mergeYear(mb)...)
	result.Suggestions = append(result.Suggestions, mergeDuration(mb)...)
	result.Suggestions = append(result.Suggestions, mergeGenres(mb, lfm)...)
	result.Suggestions = append(result.Suggestions, mergePlayCount(lfm)...)
	result.Suggestions = append(result.Suggestions, mergePopularity(spot)...)
	result.Suggestions = append(result.Suggestions, mergeArtistBio(lfm)...)
	result.Suggestions = append(result.Suggestions, mergeArtistImage(lfm)...)
	result.Suggestions = append(result.Suggestions, mergeAlbumCover(lfm)...)
	result.Suggestions = append(result.Suggestions, mergeSimilarArtists(lfm)...)
	result.Suggestions = append(result.Suggestions, mergeLyrics(lrc)...)

	return result
}

func mergeTitle(mb *MusicBrainzResult, spot *SpotifyResult) []EnrichedSuggestion {
	if mb == nil && spot == nil {
		return nil
	}
	if mb != nil && mb.Title != "" {
		return []EnrichedSuggestion{{
			Field:      "title",
			Value:      mb.Title,
			Source:     SourceMusicBrainz,
			Confidence: ConfidenceExact,
		}}
	}
	return nil
}

func mergeArtist(mb *MusicBrainzResult) []EnrichedSuggestion {
	if mb == nil || mb.ArtistName == "" {
		return nil
	}
	suggs := []EnrichedSuggestion{{
		Field:      "artist",
		Value:      mb.ArtistName,
		Source:     SourceMusicBrainz,
		Confidence: ConfidenceExact,
	}}
	if mb.ArtistMBID != "" {
		suggs = append(suggs, EnrichedSuggestion{
			Field:      "artist_mbid",
			Value:      mb.ArtistMBID,
			Source:     SourceMusicBrainz,
			Confidence: ConfidenceExact,
		})
	}
	return suggs
}

func mergeAlbum(mb *MusicBrainzResult) []EnrichedSuggestion {
	if mb == nil || mb.AlbumName == "" {
		return nil
	}
	suggs := []EnrichedSuggestion{{
		Field:      "album",
		Value:      mb.AlbumName,
		Source:     SourceMusicBrainz,
		Confidence: ConfidenceExact,
	}}
	if mb.AlbumMBID != "" {
		suggs = append(suggs, EnrichedSuggestion{
			Field:      "album_mbid",
			Value:      mb.AlbumMBID,
			Source:     SourceMusicBrainz,
			Confidence: ConfidenceExact,
		})
	}
	return suggs
}

func mergeYear(mb *MusicBrainzResult) []EnrichedSuggestion {
	if mb == nil || mb.ReleaseYear == 0 {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "year",
		Value:      mb.ReleaseYear,
		Source:     SourceMusicBrainz,
		Confidence: ConfidenceExact,
	}}
}

func mergeDuration(mb *MusicBrainzResult) []EnrichedSuggestion {
	if mb == nil || mb.Duration == 0 {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "duration_seconds",
		Value:      mb.Duration,
		Source:     SourceMusicBrainz,
		Confidence: ConfidenceExact,
	}}
}

func mergeGenres(mb *MusicBrainzResult, lfm *LastFMResult) []EnrichedSuggestion {
	var suggs []EnrichedSuggestion
	seen := map[string]bool{}

	if mb != nil {
		for _, g := range mb.Genres {
			g = strings.TrimSpace(g)
			if g != "" && !seen[g] {
				seen[g] = true
				suggs = append(suggs, EnrichedSuggestion{
					Field:      "genre",
					Value:      g,
					Source:     SourceMusicBrainz,
					Confidence: ConfidenceFuzzy,
				})
			}
		}
	}

	if lfm != nil {
		for _, g := range lfm.Tags {
			g = strings.TrimSpace(g)
			if g != "" && !seen[g] {
				seen[g] = true
				suggs = append(suggs, EnrichedSuggestion{
					Field:      "genre",
					Value:      g,
					Source:     SourceLastFM,
					Confidence: ConfidenceFuzzy,
				})
			}
		}
	}

	return suggs
}

func mergePlayCount(lfm *LastFMResult) []EnrichedSuggestion {
	if lfm == nil || lfm.PlayCount == 0 {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "play_count",
		Value:      lfm.PlayCount,
		Source:     SourceLastFM,
		Confidence: ConfidenceExact,
	}}
}

func mergePopularity(spot *SpotifyResult) []EnrichedSuggestion {
	if spot == nil {
		return nil
	}
	suggs := []EnrichedSuggestion{{
		Field:      "spotify_id",
		Value:      spot.SpotifyID,
		Source:     SourceSpotify,
		Confidence: ConfidenceExact,
	}}
	if spot.Popularity > 0 {
		suggs = append(suggs, EnrichedSuggestion{
			Field:      "popularity",
			Value:      spot.Popularity,
			Source:     SourceSpotify,
			Confidence: ConfidenceExact,
		})
	}
	if spot.PreviewURL != "" {
		suggs = append(suggs, EnrichedSuggestion{
			Field:      "preview_url",
			Value:      spot.PreviewURL,
			Source:     SourceSpotify,
			Confidence: ConfidenceExact,
		})
	}
	if spot.AlbumCoverURL != "" {
		suggs = append(suggs, EnrichedSuggestion{
			Field:      "album_cover_url",
			Value:      spot.AlbumCoverURL,
			Source:     SourceSpotify,
			Confidence: ConfidenceExact,
		})
	}
	if spot.ArtistImageURL != "" {
		suggs = append(suggs, EnrichedSuggestion{
			Field:      "artist_image_url",
			Value:      spot.ArtistImageURL,
			Source:     SourceSpotify,
			Confidence: ConfidenceExact,
		})
	}
	return suggs
}

func mergeArtistBio(lfm *LastFMResult) []EnrichedSuggestion {
	if lfm == nil || lfm.ArtistBio == "" {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "artist_bio",
		Value:      lfm.ArtistBio,
		Source:     SourceLastFM,
		Confidence: ConfidenceFuzzy,
	}}
}

func mergeLyrics(lrc *LRCLibResult) []EnrichedSuggestion {
	if lrc == nil {
		return nil
	}
	content := strings.ReplaceAll(lrc.SyncedLyrics, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	if content == "" {
		content = strings.ReplaceAll(lrc.PlainLyrics, "\r\n", "\n")
		content = strings.ReplaceAll(content, "\r", "\n")
	}
	if content == "" {
		return nil
	}
	// Detect LRC format from content itself, since LRCLIB doesn't always
	// include the "synced" boolean in its response.
	lyricsType := "plain"
	if looksLikeLRC(content) {
		lyricsType = "lrc"
	} else if lrc.Synced {
		lyricsType = "lrc"
	}
	return []EnrichedSuggestion{{
		Field:      "lyrics",
		Value:      content,
		Source:     "lrclib",
		Confidence: ConfidenceExact,
	}, {
		Field:      "lyrics_type",
		Value:      lyricsType,
		Source:     "lrclib",
		Confidence: ConfidenceExact,
	}}
}

func mergeSimilarArtists(lfm *LastFMResult) []EnrichedSuggestion {
	if lfm == nil || len(lfm.SimilarArtists) == 0 {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "similar_artists",
		Value:      lfm.SimilarArtists,
		Source:     SourceLastFM,
		Confidence: ConfidenceFuzzy,
	}}
}

func mergeArtistImage(lfm *LastFMResult) []EnrichedSuggestion {
	if lfm == nil || lfm.ArtistImageURL == "" {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "artist_image_url",
		Value:      lfm.ArtistImageURL,
		Source:     SourceLastFM,
		Confidence: ConfidenceFuzzy,
	}}
}

// looksLikeLRC returns true if the content starts with a timestamp bracket,
// indicating LRC (synced) lyrics format rather than plain text.
func looksLikeLRC(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return strings.HasPrefix(line, "[")
	}
	return false
}

func mergeAlbumCover(lfm *LastFMResult) []EnrichedSuggestion {
	if lfm == nil || lfm.AlbumCoverURL == "" {
		return nil
	}
	return []EnrichedSuggestion{{
		Field:      "album_cover_url_lastfm",
		Value:      lfm.AlbumCoverURL,
		Source:     SourceLastFM,
		Confidence: ConfidenceFuzzy,
	}}
}
