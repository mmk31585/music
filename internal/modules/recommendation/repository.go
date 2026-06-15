package recommendation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

var ErrTrackNotFound = errors.New("track not found")

type TrackMeta struct {
	ID       string
	ArtistID *string
	AlbumID  *string
	Genre    *string
}

type Repository interface {
	GetPopularTracks(ctx context.Context, limit int) ([]TrackItem, error)
	GetBestTracks(ctx context.Context, limit int) ([]TrackItem, error)
	GetRecentTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error)

	GetTrackMeta(ctx context.Context, trackID string) (*TrackMeta, error)
	GetSimilarTracksByMeta(ctx context.Context, trackID string, artistID, albumID, genre *string, limit int) ([]TrackItem, error)

	GetTracksByArtist(ctx context.Context, artistID string, limit int) ([]TrackItem, error)
	GetTracksByGenre(ctx context.Context, genre string, limit int) ([]TrackItem, error)

	GetFollowedArtistIDs(ctx context.Context, userID string, limit int) ([]string, error)
	GetTopArtistIDs(ctx context.Context, userID string, limit int) ([]string, error)
	GetTopGenres(ctx context.Context, userID string, limit int) ([]string, error)
	GetLikedTrackIDs(ctx context.Context, userID string, limit int) ([]string, error)

	GetTracksFromArtists(ctx context.Context, artistIDs []string, limit int) ([]TrackItem, error)
	GetTracksFromGenres(ctx context.Context, genres []string, limit int) ([]TrackItem, error)
	GetTracksByIDs(ctx context.Context, trackIDs []string) ([]TrackItem, error)
	GetPopularTrackIDs(ctx context.Context, limit int) ([]string, error)
	GetUserAffinities(ctx context.Context, userID string, targetType string, limit int) ([]UserAffinity, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func baseTrackSelect() string {
	return `
SELECT
	t.id,
	t.title,
	t.artist_id,
	a.name AS artist_name,
	t.album_id,
	al.title AS album_title,
	(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
	t.cover_url,
	t.audio_url,
	t.duration_seconds
FROM tracks t
LEFT JOIN artists a ON a.id = t.artist_id
LEFT JOIN albums al ON al.id = t.album_id
`
}

func (r *repository) GetPopularTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	query := `
SELECT
	t.id,
	t.title,
	t.artist_id,
	a.name AS artist_name,
	t.album_id,
	al.title AS album_title,
	(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
	t.cover_url,
	t.audio_url,
	t.duration_seconds
FROM tracks t
JOIN play_history ph ON ph.track_id = t.id
LEFT JOIN artists a ON a.id = t.artist_id
LEFT JOIN albums al ON al.id = t.album_id
WHERE ph.played_at >= NOW() - INTERVAL '30 days'
GROUP BY t.id, t.title, t.artist_id, a.name, t.album_id, al.title, t.cover_url, t.audio_url, t.duration_seconds
ORDER BY COUNT(ph.id) DESC, t.title ASC
LIMIT $1
`
	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, limit); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetBestTracks(ctx context.Context, limit int) ([]TrackItem, error) {
	query := `
SELECT
	t.id,
	t.title,
	t.artist_id,
	a.name AS artist_name,
	t.album_id,
	al.title AS album_title,
	(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
	t.cover_url,
	t.audio_url,
	t.duration_seconds,
	(COALESCE(lt.like_count, 0) * 3 + COALESCE(ph.play_count, 0))::float AS score
FROM tracks t
LEFT JOIN artists a ON a.id = t.artist_id
LEFT JOIN albums al ON al.id = t.album_id
LEFT JOIN (
	SELECT track_id, COUNT(*) AS like_count
	FROM liked_tracks
	GROUP BY track_id
) lt ON lt.track_id = t.id
LEFT JOIN (
	SELECT track_id, COUNT(*) AS play_count
	FROM play_history
	WHERE played_at >= NOW() - INTERVAL '30 days'
	GROUP BY track_id
) ph ON ph.track_id = t.id
ORDER BY score DESC NULLS LAST, t.title ASC
LIMIT $1
`
	var rows []struct {
		ID              string          `db:"id"`
		Title           string          `db:"title"`
		ArtistID        *string         `db:"artist_id"`
		ArtistName      *string         `db:"artist_name"`
		AlbumID         *string         `db:"album_id"`
		AlbumTitle      *string         `db:"album_title"`
		Genre           *string         `db:"genre"`
		CoverURL        *string         `db:"cover_url"`
		AudioURL        *string         `db:"audio_url"`
		DurationSeconds *int            `db:"duration_seconds"`
		Score           sql.NullFloat64 `db:"score"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, limit); err != nil {
		return nil, err
	}

	items := make([]TrackItem, 0, len(rows))
	for _, row := range rows {
		item := TrackItem{
			ID:              row.ID,
			Title:           row.Title,
			ArtistID:        row.ArtistID,
			ArtistName:      row.ArtistName,
			AlbumID:         row.AlbumID,
			AlbumTitle:      row.AlbumTitle,
			Genre:           row.Genre,
			CoverURL:        row.CoverURL,
			AudioURL:        row.AudioURL,
			DurationSeconds: row.DurationSeconds,
		}
		if row.Score.Valid {
			score := row.Score.Float64
			item.Score = &score
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *repository) GetRecentTracks(ctx context.Context, userID string, limit int) ([]TrackItem, error) {
	query := `
SELECT
	x.id,
	x.title,
	x.artist_id,
	x.artist_name,
	x.album_id,
	x.album_title,
	x.genre,
	x.cover_url,
	x.audio_url,
	x.duration_seconds
FROM (
	SELECT DISTINCT ON (t.id)
		t.id,
		t.title,
		t.artist_id,
		a.name AS artist_name,
		t.album_id,
		al.title AS album_title,
		(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
		t.cover_url,
		t.audio_url,
		t.duration_seconds,
		ph.played_at
	FROM play_history ph
	JOIN tracks t ON t.id = ph.track_id
	LEFT JOIN artists a ON a.id = t.artist_id
	LEFT JOIN albums al ON al.id = t.album_id
	WHERE ph.user_id = $1
	ORDER BY t.id, ph.played_at DESC
) x
ORDER BY x.played_at DESC
LIMIT $2
`
	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, userID, limit); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetTrackMeta(ctx context.Context, trackID string) (*TrackMeta, error) {
	query := `
SELECT t.id, t.artist_id, t.album_id, (SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre
FROM tracks t
WHERE t.id = $1
LIMIT 1
`
	var meta TrackMeta
	err := r.db.GetContext(ctx, &meta, query, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTrackNotFound
		}
		return nil, err
	}
	return &meta, nil
}

func (r *repository) GetSimilarTracksByMeta(ctx context.Context, trackID string, artistID, albumID, genre *string, limit int) ([]TrackItem, error) {
	scoreParts := make([]string, 0, 3)
	args := []any{trackID}
	argPos := 2

	if artistID != nil && *artistID != "" {
		scoreParts = append(scoreParts, fmt.Sprintf("CASE WHEN t.artist_id = $%d THEN 5 ELSE 0 END", argPos))
		args = append(args, *artistID)
		argPos++
	}
	if albumID != nil && *albumID != "" {
		scoreParts = append(scoreParts, fmt.Sprintf("CASE WHEN t.album_id = $%d THEN 3 ELSE 0 END", argPos))
		args = append(args, *albumID)
		argPos++
	}
	if genre != nil && *genre != "" {
		subq := fmt.Sprintf("SELECT 1 FROM track_genres tg_%[1]d JOIN genres g_%[1]d ON g_%[1]d.id = tg_%[1]d.genre_id WHERE tg_%[1]d.track_id = t.id AND g_%[1]d.name = $%[1]d", argPos)
		scoreParts = append(scoreParts, fmt.Sprintf("CASE WHEN EXISTS (%s) THEN 4 ELSE 0 END", subq))
		args = append(args, *genre)
		argPos++
	}

	if len(scoreParts) == 0 {
		return []TrackItem{}, nil
	}

	scoreExpr := strings.Join(scoreParts, " + ")
	query := fmt.Sprintf(`
SELECT
	t.id,
	t.title,
	t.artist_id,
	a.name AS artist_name,
	t.album_id,
	al.title AS album_title,
	(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
	t.cover_url,
	t.audio_url,
	t.duration_seconds,
	(%s)::float AS score
FROM tracks t
LEFT JOIN artists a ON a.id = t.artist_id
LEFT JOIN albums al ON al.id = t.album_id
WHERE t.id <> $1
ORDER BY score DESC, t.title ASC
LIMIT $%d
`, scoreExpr, argPos)

	args = append(args, limit)

	var rows []struct {
		ID              string          `db:"id"`
		Title           string          `db:"title"`
		ArtistID        *string         `db:"artist_id"`
		ArtistName      *string         `db:"artist_name"`
		AlbumID         *string         `db:"album_id"`
		AlbumTitle      *string         `db:"album_title"`
		Genre           *string         `db:"genre"`
		CoverURL        *string         `db:"cover_url"`
		AudioURL        *string         `db:"audio_url"`
		DurationSeconds *int            `db:"duration_seconds"`
		Score           sql.NullFloat64 `db:"score"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	items := make([]TrackItem, 0, len(rows))
	for _, row := range rows {
		item := TrackItem{
			ID:              row.ID,
			Title:           row.Title,
			ArtistID:        row.ArtistID,
			ArtistName:      row.ArtistName,
			AlbumID:         row.AlbumID,
			AlbumTitle:      row.AlbumTitle,
			Genre:           row.Genre,
			CoverURL:        row.CoverURL,
			AudioURL:        row.AudioURL,
			DurationSeconds: row.DurationSeconds,
		}
		if row.Score.Valid {
			score := row.Score.Float64
			item.Score = &score
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *repository) GetTracksByArtist(ctx context.Context, artistID string, limit int) ([]TrackItem, error) {
	query := baseTrackSelect() + `
WHERE t.artist_id = $1
ORDER BY t.title ASC
LIMIT $2
`
	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, artistID, limit); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetTracksByGenre(ctx context.Context, genre string, limit int) ([]TrackItem, error) {
	query := baseTrackSelect() + `
WHERE EXISTS (SELECT 1 FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id AND g.name = $1)
ORDER BY t.title ASC
LIMIT $2
`
	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, genre, limit); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetFollowedArtistIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	query := `
SELECT artist_id
FROM followed_artists
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
`
	var ids []string
	if err := r.db.SelectContext(ctx, &ids, query, userID, limit); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *repository) GetTopArtistIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	query := `
SELECT t.artist_id
FROM play_history ph
JOIN tracks t ON t.id = ph.track_id
WHERE ph.user_id = $1
  AND t.artist_id IS NOT NULL
GROUP BY t.artist_id
ORDER BY COUNT(*) DESC
LIMIT $2
`
	var ids []string
	if err := r.db.SelectContext(ctx, &ids, query, userID, limit); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *repository) GetTopGenres(ctx context.Context, userID string, limit int) ([]string, error) {
	query := `
SELECT g.name
FROM play_history ph
JOIN tracks t ON t.id = ph.track_id
JOIN track_genres tg ON tg.track_id = t.id
JOIN genres g ON g.id = tg.genre_id
WHERE ph.user_id = $1
  AND g.name IS NOT NULL
  AND g.name <> ''
GROUP BY g.name
ORDER BY COUNT(*) DESC
LIMIT $2
`
	var genres []string
	if err := r.db.SelectContext(ctx, &genres, query, userID, limit); err != nil {
		return nil, err
	}
	return genres, nil
}

func (r *repository) GetLikedTrackIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	query := `
SELECT track_id
FROM liked_tracks
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
`
	var ids []string
	if err := r.db.SelectContext(ctx, &ids, query, userID, limit); err != nil {
		return nil, err
	}
	return ids, nil
}

func buildInClause(startPos int, values []string) (string, []any) {
	parts := make([]string, 0, len(values))
	args := make([]any, 0, len(values))
	for i, v := range values {
		parts = append(parts, fmt.Sprintf("$%d", startPos+i))
		args = append(args, v)
	}
	return strings.Join(parts, ","), args
}

func (r *repository) GetTracksFromArtists(ctx context.Context, artistIDs []string, limit int) ([]TrackItem, error) {
	if len(artistIDs) == 0 {
		return []TrackItem{}, nil
	}

	inClause, args := buildInClause(1, artistIDs)
	query := baseTrackSelect() + fmt.Sprintf(`
WHERE t.artist_id IN (%s)
ORDER BY t.title ASC
LIMIT $%d
`, inClause, len(args)+1)

	args = append(args, limit)

	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetTracksFromGenres(ctx context.Context, genres []string, limit int) ([]TrackItem, error) {
	if len(genres) == 0 {
		return []TrackItem{}, nil
	}

	inClause, args := buildInClause(1, genres)
	query := fmt.Sprintf(`
SELECT
	t.id,
	t.title,
	t.artist_id,
	a.name AS artist_name,
	t.album_id,
	al.title AS album_title,
	(SELECT g.name FROM track_genres tg JOIN genres g ON g.id = tg.genre_id WHERE tg.track_id = t.id LIMIT 1) AS genre,
	t.cover_url,
	t.audio_url,
	t.duration_seconds
FROM tracks t
LEFT JOIN artists a ON a.id = t.artist_id
LEFT JOIN albums al ON al.id = t.album_id
WHERE EXISTS (SELECT 1 FROM track_genres tg_2 JOIN genres g_2 ON g_2.id = tg_2.genre_id WHERE tg_2.track_id = t.id AND g_2.name IN (%s))
ORDER BY t.title ASC
LIMIT $%d
`, inClause, len(args)+1)

	args = append(args, limit)

	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetTracksByIDs(ctx context.Context, trackIDs []string) ([]TrackItem, error) {
	if len(trackIDs) == 0 {
		return []TrackItem{}, nil
	}

	inClause, args := buildInClause(1, trackIDs)
	query := baseTrackSelect() + fmt.Sprintf(`
WHERE t.id IN (%s)
`, inClause)

	var items []TrackItem
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetPopularTrackIDs(ctx context.Context, limit int) ([]string, error) {
	query := `
SELECT t.id
FROM tracks t
JOIN play_history ph ON ph.track_id = t.id
WHERE ph.played_at >= NOW() - INTERVAL '30 days'
GROUP BY t.id
ORDER BY COUNT(ph.id) DESC, t.id ASC
LIMIT $1
`
	var ids []string
	if err := r.db.SelectContext(ctx, &ids, query, limit); err != nil {
		return nil, err
	}
	return ids, nil
}

type scoredTrack struct {
	Item  TrackItem
	Score float64
}

func sortScoredTracks(items map[string]*scoredTrack) []TrackItem {
	list := make([]scoredTrack, 0, len(items))
	for _, item := range items {
		list = append(list, *item)
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Score == list[j].Score {
			return list[i].Item.Title < list[j].Item.Title
		}
		return list[i].Score > list[j].Score
	})

	result := make([]TrackItem, 0, len(list))
	for _, item := range list {
		score := item.Score
		item.Item.Score = &score
		result = append(result, item.Item)
	}

	return result
}
