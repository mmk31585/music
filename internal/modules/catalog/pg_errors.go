package catalog

import (
	"errors"

	"github.com/jackc/pgconn"
)

const (
	pgUniqueViolation = "23505"
)

func mapPostgresError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code == pgUniqueViolation {
		switch pgErr.ConstraintName {
		case "artists_slug_key":
			return ErrDuplicateArtistSlug
		case "albums_artist_id_slug_key":
			return ErrDuplicateAlbumSlug
		case "tracks_artist_id_slug_key":
			return ErrDuplicateTrackSlug
		case "genres_slug_key":
			return ErrDuplicateGenreSlug
		case "genres_name_key":
			return ErrDuplicateGenreName
		default:
			return err
		}
	}

	return err
}
