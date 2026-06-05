package album

import (
	"context"
	"database/sql"
	"music/internal/modules/catalog/common"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func parseDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, common.ErrInvalidInput
	}

	return &t, nil
}

func normalizeAlbumArtists(reqArtists []AlbumArtistRequest, fallbackArtistID uuid.UUID) []AlbumArtistRequest {
	if len(reqArtists) == 0 && fallbackArtistID != uuid.Nil {
		return []AlbumArtistRequest{
			{
				ArtistID: fallbackArtistID,
				Role:     "primary",
				Position: 1,
			},
		}
	}

	for i := range reqArtists {
		if reqArtists[i].Role == "" {
			reqArtists[i].Role = "primary"
		}
		if reqArtists[i].Position <= 0 {
			reqArtists[i].Position = i + 1
		}
	}

	return reqArtists
}

func primaryArtistID(artists []AlbumArtistRequest, fallback uuid.UUID) uuid.UUID {
	if fallback != uuid.Nil {
		return fallback
	}

	for _, a := range artists {
		if a.Role == "primary" {
			return a.ArtistID
		}
	}

	if len(artists) > 0 {
		return artists[0].ArtistID
	}

	return uuid.Nil
}

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Album, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	artists := normalizeAlbumArtists(req.Artists, req.ArtistID)
	mainArtistID := primaryArtistID(artists, req.ArtistID)
	if mainArtistID == uuid.Nil {
		return nil, common.ErrInvalidInput
	}

	slug := common.Slugify(req.Title)

	releaseDate, err := parseDate(req.ReleaseDate)
	if err != nil {
		return nil, err
	}

	albumType := "album"
	if req.AlbumType != nil && *req.AlbumType != "" {
		albumType = *req.AlbumType
	}

	var item Album
	err = tx.GetContext(ctx, &item, `
		INSERT INTO albums (
			artist_id,
			title,
			slug,
			cover_url,
			cover_media_id,
			release_date,
			album_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			artist_id,
			title,
			slug,
			cover_url,
			cover_media_id,
			release_date,
			album_type,
			created_at,
			updated_at
	`,
		mainArtistID,
		req.Title,
		slug,
		req.CoverURL,
		req.CoverMediaID,
		releaseDate,
		albumType,
	)
	if err != nil {
		return nil, common.MapPGError(err)
	}

	if err := r.replaceArtistsTx(ctx, tx, item.ID, artists); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, item.ID)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Album, error) {
	var item Album
	err := r.db.GetContext(ctx, &item, `
		SELECT
			id,
			artist_id,
			title,
			slug,
			cover_url,
			cover_media_id,
			release_date,
			album_type,
			created_at,
			updated_at
		FROM albums
		WHERE id = $1
	`, id)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	artists, err := r.listArtistsByAlbum(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Artists = artists

	return &item, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Album, error) {
	var items []Album
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			artist_id,
			title,
			slug,
			cover_url,
			cover_media_id,
			release_date,
			album_type,
			created_at,
			updated_at
		FROM albums
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range items {
		artists, err := r.listArtistsByAlbum(ctx, items[i].ID)
		if err != nil {
			return nil, err
		}
		items[i].Artists = artists
	}

	return items, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Album, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	title := current.Title
	if req.Title != nil {
		title = *req.Title
	}
	slug := common.Slugify(title)

	releaseDate, err := parseDate(req.ReleaseDate)
	if err != nil && req.ReleaseDate != nil {
		return nil, err
	}

	albumType := current.AlbumType
	if req.AlbumType != nil && *req.AlbumType != "" {
		albumType = *req.AlbumType
	}

	var item Album
	err = tx.GetContext(ctx, &item, `
		UPDATE albums
		SET
			title = COALESCE($2, title),
			slug = $3,
			cover_url = COALESCE($4, cover_url),
			cover_media_id = COALESCE($5, cover_media_id),
			release_date = COALESCE($6, release_date),
			album_type = $7,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			artist_id,
			title,
			slug,
			cover_url,
			cover_media_id,
			release_date,
			album_type,
			created_at,
			updated_at
	`,
		id,
		req.Title,
		slug,
		req.CoverURL,
		req.CoverMediaID,
		releaseDate,
		albumType,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.MapPGError(err)
	}

	if req.Artists != nil {
		artists := normalizeAlbumArtists(req.Artists, uuid.Nil)
		if len(artists) == 0 {
			return nil, common.ErrInvalidInput
		}

		if err := r.replaceArtistsTx(ctx, tx, id, artists); err != nil {
			return nil, err
		}

		newPrimary := primaryArtistID(artists, uuid.Nil)
		if newPrimary != uuid.Nil {
			if _, err := tx.ExecContext(ctx, `
				UPDATE albums
				SET artist_id = $2, updated_at = NOW()
				WHERE id = $1
			`, id, newPrimary); err != nil {
				return nil, common.MapPGError(err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, item.ID)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM albums WHERE id = $1`, id)
	if err != nil {
		return common.MapPGError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return common.ErrNotFound
	}

	return nil
}

func (r *Repository) replaceArtistsTx(ctx context.Context, tx *sqlx.Tx, albumID uuid.UUID, artists []AlbumArtistRequest) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM album_artists WHERE album_id = $1`, albumID); err != nil {
		return common.MapPGError(err)
	}

	for i, a := range artists {
		role := a.Role
		if role == "" {
			role = "primary"
		}

		position := a.Position
		if position <= 0 {
			position = i + 1
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO album_artists (
				album_id,
				artist_id,
				role,
				position
			)
			VALUES ($1, $2, $3, $4)
		`, albumID, a.ArtistID, role, position); err != nil {
			return common.MapPGError(err)
		}
	}

	return nil
}

func (r *Repository) listArtistsByAlbum(ctx context.Context, albumID uuid.UUID) ([]AlbumArtist, error) {
	var items []AlbumArtist
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			a.id AS artist_id,
			a.name,
			a.slug,
			aa.role,
			aa.position
		FROM album_artists aa
		INNER JOIN artists a ON a.id = aa.artist_id
		WHERE aa.album_id = $1
		ORDER BY aa.position ASC, aa.role ASC, a.name ASC
	`, albumID)

	return items, err
}

func (r *Repository) ListTracks(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]AlbumTrack, error) {
	var items []AlbumTrack
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			play_count,
			is_public,
			created_at
		FROM tracks
		WHERE album_id = $1
		  AND is_public = TRUE
		ORDER BY
			track_number ASC NULLS LAST,
			created_at ASC
		LIMIT $2 OFFSET $3
	`, albumID, limit, offset)

	return items, err
}
func (r *Repository) ListArtists(ctx context.Context, albumID uuid.UUID) ([]AlbumArtist, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM albums
			WHERE id = $1
		)
	`, albumID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	return r.listArtistsByAlbum(ctx, albumID)
}

func (r *Repository) ReplaceArtists(ctx context.Context, albumID uuid.UUID, artists []AlbumArtistRequest) ([]AlbumArtist, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM albums
			WHERE id = $1
		)
	`, albumID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	artists = normalizeAlbumArtists(artists, uuid.Nil)
	if len(artists) == 0 {
		return nil, common.ErrInvalidInput
	}

	for _, a := range artists {
		if a.ArtistID == uuid.Nil {
			return nil, common.ErrInvalidInput
		}
	}

	if err := r.replaceArtistsTx(ctx, tx, albumID, artists); err != nil {
		return nil, err
	}

	newPrimary := primaryArtistID(artists, uuid.Nil)
	if newPrimary != uuid.Nil {
		if _, err := tx.ExecContext(ctx, `
			UPDATE albums
			SET artist_id = $2, updated_at = NOW()
			WHERE id = $1
		`, albumID, newPrimary); err != nil {
			return nil, common.MapPGError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.listArtistsByAlbum(ctx, albumID)
}
