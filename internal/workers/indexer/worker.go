package indexer

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	osclient "music/internal/platform/opensearch"
)

type trackDoc struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	PersianTitle string   `json:"persian_title"`
	Artist       string   `json:"artist_name"`
	Album        string   `json:"album_title"`
	Lyrics       string   `json:"lyrics"`
	Genres       []string `json:"genre_names"`
	Duration     int      `json:"duration"`
	Year         int      `json:"year"`
	Popularity   float64  `json:"popularity"`
	Explicit     bool     `json:"is_explicit"`
	CreatedAt    string   `json:"created_at"`
}

type artistDoc struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	PersianName      string   `json:"persian_name"`
	Bio              string   `json:"bio"`
	PersianBio       string   `json:"persian_bio"`
	MonthlyListeners int      `json:"monthly_listeners"`
	Verified         bool     `json:"is_verified"`
	Genres           []string `json:"genre_names"`
}

type albumDoc struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	PersianTitle string   `json:"persian_title"`
	Artist       string   `json:"artist_name"`
	Year         int      `json:"release_year"`
	Type         string   `json:"album_type"`
	Genres       []string `json:"genre_names"`
}

type Worker struct {
	db       *sqlx.DB
	osClient *osclient.Client
	logger   *zap.Logger
}

func NewWorker(db *sqlx.DB, osClient *osclient.Client, logger *zap.Logger) *Worker {
	return &Worker{db: db, osClient: osClient, logger: logger}
}

func (w *Worker) Name() string { return "indexer" }

func (w *Worker) Run(ctx context.Context) error {
	if w.osClient == nil || !w.osClient.IsHealthy() {
		w.logger.Warn("opensearch not available, indexer worker idling")
		return nil
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.indexTracks(ctx); err != nil {
				w.logger.Error("track indexing failed", zap.Error(err))
			}
			if err := w.indexArtists(ctx); err != nil {
				w.logger.Error("artist indexing failed", zap.Error(err))
			}
			if err := w.indexAlbums(ctx); err != nil {
				w.logger.Error("album indexing failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) indexTracks(ctx context.Context) error {
	rows, err := w.db.QueryContext(ctx, `
		SELECT
			t.id, t.title, COALESCE(t.persian_title, '') AS persian_title,
			COALESCE(a.name, '') AS artist_name, COALESCE(al.title, '') AS album_title,
			COALESCE(t.lyrics, '') AS lyrics,
			t.duration_seconds, t.year,
			COALESCE(tp.play_count, 0) AS popularity,
			t.explicit, t.created_at::text
		FROM tracks t
		LEFT JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN (SELECT track_id, COUNT(*) AS play_count FROM listening_history GROUP BY track_id) tp ON tp.track_id = t.id
		WHERE t.id NOT IN (SELECT entity_id FROM search_sync WHERE entity_type = 'track' AND synced = true AND synced_at > NOW() - INTERVAL '1 hour')
		ORDER BY t.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var doc trackDoc
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.PersianTitle, &doc.Artist, &doc.Album,
			&doc.Lyrics, &doc.Duration, &doc.Year, &doc.Popularity, &doc.Explicit, &doc.CreatedAt); err != nil {
			continue
		}
		if doc.PersianTitle == "" {
			doc.PersianTitle = normalizePersian(doc.Title)
		}

		_ = w.db.SelectContext(ctx, &doc.Genres, `
			SELECT g.name FROM genres g
			JOIN track_genres tg ON tg.genre_id = g.id
			WHERE tg.track_id = $1
		`, doc.ID)

		if err := w.osClient.IndexDocument(ctx, "tracks", doc.ID, doc); err != nil {
			w.logger.Warn("failed to index track", zap.String("track_id", doc.ID), zap.Error(err))
			continue
		}

		w.db.ExecContext(ctx, `
			INSERT INTO search_sync (entity_type, entity_id, action, synced, synced_at)
			VALUES ('track', $1, 'index', true, NOW())
			ON CONFLICT (entity_type, entity_id) DO UPDATE SET synced = true, synced_at = NOW()
		`, doc.ID)
	}

	return rows.Err()
}

func (w *Worker) indexArtists(ctx context.Context) error {
	rows, err := w.db.QueryContext(ctx, `
		SELECT a.id, a.name, COALESCE(a.persian_name, '') AS persian_name,
			COALESCE(a.bio, ''), COALESCE(a.persian_bio, '') AS persian_bio,
			a.monthly_listeners, a.is_verified
		FROM artists a
		WHERE a.id NOT IN (SELECT entity_id FROM search_sync WHERE entity_type = 'artist' AND synced = true AND synced_at > NOW() - INTERVAL '1 hour')
		LIMIT 100
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var doc artistDoc
		if err := rows.Scan(&doc.ID, &doc.Name, &doc.PersianName, &doc.Bio, &doc.PersianBio,
			&doc.MonthlyListeners, &doc.Verified); err != nil {
			continue
		}
		if doc.PersianName == "" {
			doc.PersianName = normalizePersian(doc.Name)
		}

		if err := w.osClient.IndexDocument(ctx, "artists", doc.ID, doc); err != nil {
			w.logger.Warn("failed to index artist", zap.String("artist_id", doc.ID), zap.Error(err))
			continue
		}

		w.db.ExecContext(ctx, `
			INSERT INTO search_sync (entity_type, entity_id, action, synced, synced_at)
			VALUES ('artist', $1, 'index', true, NOW())
			ON CONFLICT (entity_type, entity_id) DO UPDATE SET synced = true, synced_at = NOW()
		`, doc.ID)
	}
	return rows.Err()
}

func (w *Worker) indexAlbums(ctx context.Context) error {
	rows, err := w.db.QueryContext(ctx, `
		SELECT al.id, al.title, COALESCE(al.persian_title, '') AS persian_title,
			COALESCE(a.name, '') AS artist_name,
			EXTRACT(YEAR FROM al.release_date)::int AS release_year,
			COALESCE(al.album_type, 'album') AS album_type
		FROM albums al
		LEFT JOIN artists a ON a.id = al.artist_id
		WHERE al.id NOT IN (SELECT entity_id FROM search_sync WHERE entity_type = 'album' AND synced = true AND synced_at > NOW() - INTERVAL '1 hour')
		LIMIT 100
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var doc albumDoc
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.PersianTitle, &doc.Artist, &doc.Year, &doc.Type); err != nil {
			continue
		}
		if doc.PersianTitle == "" {
			doc.PersianTitle = normalizePersian(doc.Title)
		}

		if err := w.osClient.IndexDocument(ctx, "albums", doc.ID, doc); err != nil {
			w.logger.Warn("failed to index album", zap.String("album_id", doc.ID), zap.Error(err))
			continue
		}

		w.db.ExecContext(ctx, `
			INSERT INTO search_sync (entity_type, entity_id, action, synced, synced_at)
			VALUES ('album', $1, 'index', true, NOW())
			ON CONFLICT (entity_type, entity_id) DO UPDATE SET synced = true, synced_at = NOW()
		`, doc.ID)
	}
	return rows.Err()
}

func normalizePersian(text string) string {
	replacer := strings.NewReplacer(
		"ي", "ی",
		"ك", "ک",
		"ة", "ه",
	)
	return replacer.Replace(text)
}
