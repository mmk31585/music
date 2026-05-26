package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// --------------------
// Slug existence
// --------------------

func (r *Repository) ArtistSlugExists(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM artists WHERE slug = $1`
	args := []interface{}{slug}

	if excludeID != nil {
		query += ` AND id <> $2`
		args = append(args, *excludeID)
	}

	query += `)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	return exists, err
}

func (r *Repository) AlbumSlugExists(ctx context.Context, artistID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM albums WHERE artist_id = $1 AND slug = $2`
	args := []interface{}{artistID, slug}

	if excludeID != nil {
		query += ` AND id <> $3`
		args = append(args, *excludeID)
	}

	query += `)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	return exists, err
}

func (r *Repository) TrackSlugExists(ctx context.Context, artistID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM tracks WHERE artist_id = $1 AND slug = $2`
	args := []interface{}{artistID, slug}

	if excludeID != nil {
		query += ` AND id <> $3`
		args = append(args, *excludeID)
	}

	query += `)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	return exists, err
}

func (r *Repository) GenreSlugExists(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM genres WHERE slug = $1`
	args := []interface{}{slug}

	if excludeID != nil {
		query += ` AND id <> $2`
		args = append(args, *excludeID)
	}

	query += `)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	return exists, err
}

// --------------------
// Artists
// --------------------

func (r *Repository) CreateArtist(ctx context.Context, artist *Artist) error {
	query := `
		INSERT INTO artists (
			id, name, slug, bio, image_url, is_verified,
			monthly_listeners, created_at, updated_at
		) VALUES (
			:id, :name, :slug, :bio, :image_url, :is_verified,
			:monthly_listeners, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, artist)
	return mapPostgresError(err)
}

func (r *Repository) GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	var artist Artist

	err := r.db.GetContext(ctx, &artist, `
		SELECT id, name, slug, bio, image_url, is_verified,
		       monthly_listeners, created_at, updated_at
		FROM artists
		WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArtistNotFound
		}
		return nil, err
	}

	return &artist, nil
}

func (r *Repository) ListArtists(ctx context.Context, filter ArtistListFilter) ([]Artist, error) {
	where, args := buildArtistWhere(filter)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)
	limitPos := len(args) - 1
	offsetPos := len(args)

	query := fmt.Sprintf(`
		SELECT id, name, slug, bio, image_url, is_verified,
		       monthly_listeners, created_at, updated_at
		FROM artists
		%s
		ORDER BY monthly_listeners DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, limitPos, offsetPos)

	var artists []Artist
	err := r.db.SelectContext(ctx, &artists, query, args...)
	return artists, err
}

func (r *Repository) CountArtists(ctx context.Context, filter ArtistListFilter) (int, error) {
	where, args := buildArtistWhere(filter)

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM artists
		%s
	`, where)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func buildArtistWhere(filter ArtistListFilter) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}

	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR slug ILIKE $%d)", len(args), len(args)))
	}

	if filter.Verified != nil {
		args = append(args, *filter.Verified)
		conditions = append(conditions, fmt.Sprintf("is_verified = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (r *Repository) UpdateArtist(ctx context.Context, artist *Artist) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE artists
		SET name = $2,
		    slug = $3,
		    bio = $4,
		    image_url = $5,
		    is_verified = $6,
		    monthly_listeners = $7,
		    updated_at = $8
		WHERE id = $1
	`,
		artist.ID,
		artist.Name,
		artist.Slug,
		artist.Bio,
		artist.ImageURL,
		artist.IsVerified,
		artist.MonthlyListeners,
		artist.UpdatedAt,
	)

	if err != nil {
		return mapPostgresError(err)
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrArtistNotFound
	}

	return err
}

func (r *Repository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM artists WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrArtistNotFound
	}

	return err
}

// --------------------
// Albums
// --------------------

func (r *Repository) CreateAlbum(ctx context.Context, album *Album) error {
	query := `
		INSERT INTO albums (
			id, artist_id, title, slug, cover_url,
			release_date, album_type, created_at, updated_at
		) VALUES (
			:id, :artist_id, :title, :slug, :cover_url,
			:release_date, :album_type, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, album)
	return mapPostgresError(err)
}

func (r *Repository) GetAlbumByID(ctx context.Context, id uuid.UUID) (*Album, error) {
	var album Album

	err := r.db.GetContext(ctx, &album, `
		SELECT id, artist_id, title, slug, cover_url,
		       release_date, album_type, created_at, updated_at
		FROM albums
		WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAlbumNotFound
		}
		return nil, err
	}

	return &album, nil
}

func (r *Repository) ListAlbums(ctx context.Context, filter AlbumListFilter) ([]Album, error) {
	where, args := buildAlbumWhere(filter)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)
	limitPos := len(args) - 1
	offsetPos := len(args)

	query := fmt.Sprintf(`
		SELECT id, artist_id, title, slug, cover_url,
		       release_date, album_type, created_at, updated_at
		FROM albums
		%s
		ORDER BY release_date DESC NULLS LAST, created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, limitPos, offsetPos)

	var albums []Album
	err := r.db.SelectContext(ctx, &albums, query, args...)
	return albums, err
}

func (r *Repository) CountAlbums(ctx context.Context, filter AlbumListFilter) (int, error) {
	where, args := buildAlbumWhere(filter)

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM albums
		%s
	`, where)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func buildAlbumWhere(filter AlbumListFilter) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}

	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR slug ILIKE $%d)", len(args), len(args)))
	}

	if filter.ArtistID != nil {
		args = append(args, *filter.ArtistID)
		conditions = append(conditions, fmt.Sprintf("artist_id = $%d", len(args)))
	}

	if filter.AlbumType != "" {
		args = append(args, filter.AlbumType)
		conditions = append(conditions, fmt.Sprintf("album_type = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (r *Repository) UpdateAlbum(ctx context.Context, album *Album) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE albums
		SET title = $2,
		    slug = $3,
		    cover_url = $4,
		    release_date = $5,
		    album_type = $6,
		    updated_at = $7
		WHERE id = $1
	`,
		album.ID,
		album.Title,
		album.Slug,
		album.CoverURL,
		album.ReleaseDate,
		album.AlbumType,
		album.UpdatedAt,
	)

	if err != nil {
		return mapPostgresError(err)
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrAlbumNotFound
	}

	return err
}

func (r *Repository) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM albums WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrAlbumNotFound
	}

	return err
}

// --------------------
// Tracks
// --------------------

func (r *Repository) CreateTrack(ctx context.Context, track *Track, genreIDs []uuid.UUID) error {
	return r.withTx(ctx, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO tracks (
				id, artist_id, album_id, title, slug,
				duration_seconds, track_number, explicit,
				audio_url, cover_url, play_count, is_public,
				created_at, updated_at
			) VALUES (
				:id, :artist_id, :album_id, :title, :slug,
				:duration_seconds, :track_number, :explicit,
				:audio_url, :cover_url, :play_count, :is_public,
				:created_at, :updated_at
			)
		`

		if _, err := tx.NamedExecContext(ctx, query, track); err != nil {
			return mapPostgresError(err)
		}

		return r.replaceTrackGenresTx(ctx, tx, track.ID, genreIDs)
	})
}

func (r *Repository) GetTrackByID(ctx context.Context, id uuid.UUID, includePrivate bool) (*Track, error) {
	var track Track

	query := `
		SELECT id, artist_id, album_id, title, slug,
		       duration_seconds, track_number, explicit,
		       audio_url, cover_url, play_count, is_public,
		       created_at, updated_at
		FROM tracks
		WHERE id = $1
	`

	args := []interface{}{id}

	if !includePrivate {
		query += ` AND is_public = TRUE`
	}

	err := r.db.GetContext(ctx, &track, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTrackNotFound
		}
		return nil, err
	}

	genres, err := r.GetGenresByTrackIDs(ctx, []uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	track.Genres = genres[id]

	return &track, nil
}

func (r *Repository) ListTracks(ctx context.Context, filter TrackListFilter) ([]Track, error) {
	where, args := buildTrackWhere(filter)

	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)
	limitPos := len(args) - 1
	offsetPos := len(args)

	query := fmt.Sprintf(`
		SELECT DISTINCT t.id, t.artist_id, t.album_id, t.title, t.slug,
		       t.duration_seconds, t.track_number, t.explicit,
		       t.audio_url, t.cover_url, t.play_count, t.is_public,
		       t.created_at, t.updated_at
		FROM tracks t
		%s
		ORDER BY t.play_count DESC, t.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, limitPos, offsetPos)

	var tracks []Track
	err := r.db.SelectContext(ctx, &tracks, query, args...)
	return tracks, err
}

func (r *Repository) CountTracks(ctx context.Context, filter TrackListFilter) (int, error) {
	where, args := buildTrackWhere(filter)

	query := fmt.Sprintf(`
		SELECT COUNT(DISTINCT t.id)
		FROM tracks t
		%s
	`, where)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func buildTrackWhere(filter TrackListFilter) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	joins := ""

	if filter.GenreID != nil {
		joins += ` JOIN track_genres tg ON tg.track_id = t.id `
		args = append(args, *filter.GenreID)
		conditions = append(conditions, fmt.Sprintf("tg.genre_id = $%d", len(args)))
	}

	if !filter.IncludePrivate {
		conditions = append(conditions, "t.is_public = TRUE")
	} else if filter.IsPublic != nil {
		args = append(args, *filter.IsPublic)
		conditions = append(conditions, fmt.Sprintf("t.is_public = $%d", len(args)))
	}

	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		conditions = append(conditions, fmt.Sprintf("(t.title ILIKE $%d OR t.slug ILIKE $%d)", len(args), len(args)))
	}

	if filter.ArtistID != nil {
		args = append(args, *filter.ArtistID)
		conditions = append(conditions, fmt.Sprintf("t.artist_id = $%d", len(args)))
	}

	if filter.AlbumID != nil {
		args = append(args, *filter.AlbumID)
		conditions = append(conditions, fmt.Sprintf("t.album_id = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return joins, args
	}

	return joins + " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *Repository) UpdateTrack(ctx context.Context, track *Track, genreIDs []uuid.UUID) error {
	return r.withTx(ctx, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, `
			UPDATE tracks
			SET album_id = $2,
			    title = $3,
			    slug = $4,
			    duration_seconds = $5,
			    track_number = $6,
			    explicit = $7,
			    audio_url = $8,
			    cover_url = $9,
			    is_public = $10,
			    updated_at = $11
			WHERE id = $1
		`,
			track.ID,
			track.AlbumID,
			track.Title,
			track.Slug,
			track.DurationSeconds,
			track.TrackNumber,
			track.Explicit,
			track.AudioURL,
			track.CoverURL,
			track.IsPublic,
			track.UpdatedAt,
		)

		if err != nil {
			return mapPostgresError(err)
		}

		rows, err := res.RowsAffected()
		if err == nil && rows == 0 {
			return ErrTrackNotFound
		}
		if err != nil {
			return err
		}

		return r.replaceTrackGenresTx(ctx, tx, track.ID, genreIDs)
	})
}

func (r *Repository) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tracks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrTrackNotFound
	}

	return err
}

func (r *Repository) GetGenresByTrackIDs(ctx context.Context, trackIDs []uuid.UUID) (map[uuid.UUID][]Genre, error) {
	result := make(map[uuid.UUID][]Genre)

	if len(trackIDs) == 0 {
		return result, nil
	}

	query, args, err := sqlx.In(`
		SELECT tg.track_id,
		       g.id,
		       g.name,
		       g.slug,
		       g.created_at
		FROM track_genres tg
		JOIN genres g ON g.id = tg.genre_id
		WHERE tg.track_id IN (?)
		ORDER BY g.name ASC
	`, trackIDs)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	type row struct {
		TrackID   uuid.UUID `db:"track_id"`
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		Slug      string    `db:"slug"`
		CreatedAt any       `db:"created_at"`
	}

	rows := []struct {
		TrackID   uuid.UUID    `db:"track_id"`
		ID        uuid.UUID    `db:"id"`
		Name      string       `db:"name"`
		Slug      string       `db:"slug"`
		CreatedAt sql.NullTime `db:"created_at"`
	}{}

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	for _, item := range rows {
		if item.CreatedAt.Valid {
			result[item.TrackID] = append(result[item.TrackID], Genre{
				ID:        item.ID,
				Name:      item.Name,
				Slug:      item.Slug,
				CreatedAt: item.CreatedAt.Time,
			})
		}
	}

	return result, nil
}

func (r *Repository) replaceTrackGenresTx(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID, genreIDs []uuid.UUID) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM track_genres WHERE track_id = $1`, trackID); err != nil {
		return err
	}

	if len(genreIDs) == 0 {
		return nil
	}

	for _, genreID := range genreIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO track_genres (track_id, genre_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, trackID, genreID)
		if err != nil {
			return err
		}
	}

	return nil
}

// --------------------
// Genres
// --------------------

func (r *Repository) CreateGenre(ctx context.Context, genre *Genre) error {
	query := `
		INSERT INTO genres (
			id, name, slug, created_at
		) VALUES (
			:id, :name, :slug, :created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, genre)
	return mapPostgresError(err)
}

func (r *Repository) GetGenreByID(ctx context.Context, id uuid.UUID) (*Genre, error) {
	var genre Genre

	err := r.db.GetContext(ctx, &genre, `
		SELECT id, name, slug, created_at
		FROM genres
		WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGenreNotFound
		}
		return nil, err
	}

	return &genre, nil
}

func (r *Repository) ListGenres(ctx context.Context) ([]Genre, error) {
	var genres []Genre

	err := r.db.SelectContext(ctx, &genres, `
		SELECT id, name, slug, created_at
		FROM genres
		ORDER BY name ASC
	`)

	return genres, err
}

func (r *Repository) UpdateGenre(ctx context.Context, genre *Genre) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE genres
		SET name = $2,
		    slug = $3
		WHERE id = $1
	`,
		genre.ID,
		genre.Name,
		genre.Slug,
	)

	if err != nil {
		return mapPostgresError(err)
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrGenreNotFound
	}

	return err
}

func (r *Repository) DeleteGenre(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM genres WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrGenreNotFound
	}

	return err
}

func (r *Repository) CountGenresByIDs(ctx context.Context, ids []uuid.UUID) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	query, args, err := sqlx.In(`
		SELECT COUNT(*)
		FROM genres
		WHERE id IN (?)
	`, ids)
	if err != nil {
		return 0, err
	}

	query = r.db.Rebind(query)

	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

// --------------------
// Search
// --------------------

func (r *Repository) Search(ctx context.Context, query string, limit int) (*SearchResponse, error) {
	result := &SearchResponse{}

	search := "%" + query + "%"

	if err := r.db.SelectContext(ctx, &result.Artists, `
		SELECT id, name, slug, bio, image_url, is_verified,
		       monthly_listeners, created_at, updated_at
		FROM artists
		WHERE name ILIKE $1 OR slug ILIKE $1
		ORDER BY monthly_listeners DESC, created_at DESC
		LIMIT $2
	`, search, limit); err != nil {
		return nil, err
	}

	if err := r.db.SelectContext(ctx, &result.Albums, `
		SELECT id, artist_id, title, slug, cover_url,
		       release_date, album_type, created_at, updated_at
		FROM albums
		WHERE title ILIKE $1 OR slug ILIKE $1
		ORDER BY release_date DESC NULLS LAST, created_at DESC
		LIMIT $2
	`, search, limit); err != nil {
		return nil, err
	}

	if err := r.db.SelectContext(ctx, &result.Tracks, `
		SELECT id, artist_id, album_id, title, slug,
		       duration_seconds, track_number, explicit,
		       audio_url, cover_url, play_count, is_public,
		       created_at, updated_at
		FROM tracks
		WHERE is_public = TRUE
		  AND (title ILIKE $1 OR slug ILIKE $1)
		ORDER BY play_count DESC, created_at DESC
		LIMIT $2
	`, search, limit); err != nil {
		return nil, err
	}

	return result, nil
}

// --------------------
// Transaction helper
// --------------------

func (r *Repository) withTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
