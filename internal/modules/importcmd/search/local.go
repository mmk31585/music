package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type LocalProvider struct {
	db *sqlx.DB
}

func NewLocalProvider(db *sqlx.DB) *LocalProvider {
	return &LocalProvider{db: db}
}

func (p *LocalProvider) Name() string { return "local" }

func (p *LocalProvider) Search(ctx context.Context, q SearchQuery) ([]Result, error) {
	query := q.Raw
	if query == "" {
		query = strings.TrimSpace(q.Title + " " + q.Artist)
	}

	sqlQuery := `
		SELECT
			t.id,
			t.title,
			t.duration_seconds,
			t.cover_url,
			COALESCE(a.name, '') AS artist_name,
			COALESCE(alb.title, '') AS album_title
		FROM tracks t
		LEFT JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums alb ON alb.id = t.album_id
		WHERE t.title ILIKE '%' || $1 || '%'
		   OR a.name ILIKE '%' || $1 || '%'
		LIMIT 10
	`

	rows, err := p.db.QueryxContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, fmt.Errorf("local search: %w", err)
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var id int64
		var title, artistName, albumTitle, coverURL string
		var duration *int

		if err := rows.Scan(&id, &title, &duration, &coverURL, &artistName, &albumTitle); err != nil {
			continue
		}

		dur := 0
		if duration != nil {
			dur = *duration
		}

		results = append(results, Result{
			Title:     title,
			Artist:    artistName,
			Album:     albumTitle,
			Duration:  dur,
			Thumbnail: coverURL,
			Source:    "local",
		})
	}

	return results, nil
}
