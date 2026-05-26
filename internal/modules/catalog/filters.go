package catalog

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"music/internal/common/pagination"
)

type ArtistListFilter struct {
	Pagination pagination.Params
	Query      string
	Verified   *bool
}

type AlbumListFilter struct {
	Pagination pagination.Params
	Query      string
	ArtistID   *uuid.UUID
	AlbumType  string
}

type TrackListFilter struct {
	Pagination     pagination.Params
	Query          string
	ArtistID       *uuid.UUID
	AlbumID        *uuid.UUID
	GenreID        *uuid.UUID
	IsPublic       *bool
	IncludePrivate bool
}

// ParseArtistListFilter extracts artist list filters from Gin context.
func ParseArtistListFilter(c *gin.Context) ArtistListFilter {
	var verified *bool
	if v := c.Query("verified"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			verified = &parsed
		}
	}

	return ArtistListFilter{
		Pagination: pagination.FromRequest(c.Request), // reuse existing helper with underlying request
		Query:      strings.TrimSpace(c.Query("q")),
		Verified:   verified,
	}
}

// ParseAlbumListFilter extracts album list filters from Gin context.
func ParseAlbumListFilter(c *gin.Context) AlbumListFilter {
	var artistID *uuid.UUID
	if v := c.Query("artistId"); v != "" {
		parsed, err := uuid.Parse(v)
		if err == nil {
			artistID = &parsed
		}
	}

	return AlbumListFilter{
		Pagination: pagination.FromRequest(c.Request),
		Query:      strings.TrimSpace(c.Query("q")),
		ArtistID:   artistID,
		AlbumType:  strings.TrimSpace(c.Query("albumType")),
	}
}

// ParseTrackListFilter extracts track list filters from Gin context.
func ParseTrackListFilter(c *gin.Context, includePrivate bool) TrackListFilter {
	var artistID *uuid.UUID
	if v := c.Query("artistId"); v != "" {
		parsed, err := uuid.Parse(v)
		if err == nil {
			artistID = &parsed
		}
	}

	var albumID *uuid.UUID
	if v := c.Query("albumId"); v != "" {
		parsed, err := uuid.Parse(v)
		if err == nil {
			albumID = &parsed
		}
	}

	var genreID *uuid.UUID
	if v := c.Query("genreId"); v != "" {
		parsed, err := uuid.Parse(v)
		if err == nil {
			genreID = &parsed
		}
	}

	var isPublic *bool
	if includePrivate {
		if v := c.Query("isPublic"); v != "" {
			parsed, err := strconv.ParseBool(v)
			if err == nil {
				isPublic = &parsed
			}
		}
	}

	return TrackListFilter{
		Pagination:     pagination.FromRequest(c.Request),
		Query:          strings.TrimSpace(c.Query("q")),
		ArtistID:       artistID,
		AlbumID:        albumID,
		GenreID:        genreID,
		IsPublic:       isPublic,
		IncludePrivate: includePrivate,
	}
}
