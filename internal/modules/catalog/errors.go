package catalog

import "errors"

var (
	ErrArtistNotFound = errors.New("artist not found")
	ErrAlbumNotFound  = errors.New("album not found")
	ErrTrackNotFound  = errors.New("track not found")
	ErrGenreNotFound  = errors.New("genre not found")

	ErrDuplicateArtistSlug = errors.New("artist slug already exists")
	ErrDuplicateAlbumSlug  = errors.New("album slug already exists")
	ErrDuplicateTrackSlug  = errors.New("track slug already exists")
	ErrDuplicateGenreSlug  = errors.New("genre slug already exists")
	ErrDuplicateGenreName  = errors.New("genre name already exists")

	ErrInvalidGenreIDs = errors.New("one or more genre ids are invalid")
)
