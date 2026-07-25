package stats

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository provides catalog count queries.
type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CountAll returns the total number of tracks, albums, artists, and genres.
func (r *Repository) CountAll(ctx context.Context) (StatsResponse, error) {
	var s StatsResponse

	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tracks`).Scan(&s.TotalTracks)
	if err != nil {
		return s, err
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM albums`).Scan(&s.TotalAlbums)
	if err != nil {
		return s, err
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM artists`).Scan(&s.TotalArtists)
	if err != nil {
		return s, err
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM genres`).Scan(&s.TotalGenres)
	if err != nil {
		return s, err
	}

	return s, nil
}
