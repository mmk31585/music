package track

import (
	"context"
	"database/sql"
	"log/slog"
	"music/internal/modules/catalog/common"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func normalizeTrackArtists(reqArtists []TrackArtistRequest, fallbackArtistID uuid.UUID) []TrackArtistRequest {
	if len(reqArtists) == 0 && fallbackArtistID != uuid.Nil {
		return []TrackArtistRequest{
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

func primaryArtistID(artists []TrackArtistRequest, fallback uuid.UUID) uuid.UUID {
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

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Track, error) {
	slog.InfoContext(ctx, "creating track",
		"title", req.Title,
		"album_id", req.AlbumID,
		"artist_id", req.ArtistID,
		"artist_ids_count", len(req.ArtistIDs),
		"artists_count", len(req.Artists),
		"credits_count", len(req.Credits),
	)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin transaction", "error", err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err == nil {
			slog.DebugContext(ctx, "transaction rolled back")
		}
	}()

	fallbackArtistID := req.ArtistID
	if fallbackArtistID == uuid.Nil && len(req.ArtistIDs) > 0 {
		fallbackArtistID = req.ArtistIDs[0]
	}

	artistsInput := req.Artists
	if len(artistsInput) == 0 && len(req.Credits) > 0 {
		artistsInput = req.Credits
	}

	artists := normalizeTrackArtists(artistsInput, fallbackArtistID)

	mainArtistID := primaryArtistID(artists, fallbackArtistID)
	if mainArtistID == uuid.Nil {
		slog.WarnContext(ctx, "invalid main artist id")
		return nil, common.ErrInvalidInput
	}

	slug := common.Slugify(req.Title)

	explicit := false
	if req.Explicit != nil {
		explicit = *req.Explicit
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	slog.DebugContext(ctx, "prepared track data",
		"slug", slug,
		"explicit", explicit,
		"is_public", isPublic,
	)

	var item Track

	err = tx.GetContext(ctx, &item, `
		INSERT INTO tracks (
			artist_id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			audio_media_id,
			cover_media_id,
			is_public
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING
			id,
			artist_id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			audio_media_id,
			cover_media_id,
			play_count,
			is_public,
			created_at,
			updated_at
	`,
		mainArtistID,
		req.AlbumID,
		req.Title,
		slug,
		req.DurationSeconds,
		req.TrackNumber,
		explicit,
		req.AudioURL,
		req.CoverURL,
		req.AudioMediaID,
		req.CoverMediaID,
		isPublic,
	)

	if err != nil {
		slog.ErrorContext(ctx, "failed to insert track",
			"title", req.Title,
			"artist_id", mainArtistID,
			"error", err,
		)
		return nil, common.MapPGError(err)
	}

	slog.InfoContext(ctx, "track inserted",
		"track_id", item.ID,
	)

	if err := r.replaceArtistsTx(ctx, tx, item.ID, artists); err != nil {
		slog.ErrorContext(ctx, "failed to replace track artists",
			"track_id", item.ID,
			"error", err,
		)
		return nil, err
	}

	slog.DebugContext(ctx, "track artists updated",
		"track_id", item.ID,
		"artists_count", len(artists),
	)

	if err := r.replaceGenresTx(ctx, tx, item.ID, req.GenreIDs); err != nil {
		slog.ErrorContext(ctx, "failed to replace track genres",
			"track_id", item.ID,
			"error", err,
		)
		return nil, err
	}

	slog.DebugContext(ctx, "track genres updated",
		"track_id", item.ID,
		"genres_count", len(req.GenreIDs),
	)

	if err := tx.Commit(); err != nil {
		slog.ErrorContext(ctx, "failed to commit transaction",
			"track_id", item.ID,
			"error", err,
		)
		return nil, err
	}

	slog.InfoContext(ctx, "track created successfully",
		"track_id", item.ID,
		"title", item.Title,
	)

	result, err := r.GetByID(ctx, item.ID)
	if err != nil {
		slog.ErrorContext(ctx, "track created but failed to reload",
			"track_id", item.ID,
			"error", err,
		)
		return nil, err
	}

	slog.DebugContext(ctx, "track loaded successfully",
		"track_id", item.ID,
	)

	return result, nil
}
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Track, error) {
	var item Track
	err := r.db.GetContext(ctx, &item, `
		SELECT
			id,
			artist_id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			audio_media_id,
			cover_media_id,
			play_count,
			is_public,
			created_at,
			updated_at
		FROM tracks
		WHERE id = $1
	`, id)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := r.hydrate(ctx, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) Random(ctx context.Context, limit int) ([]Track, error) {
	var items []Track
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			artist_id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			audio_media_id,
			cover_media_id,
			play_count,
			is_public,
			created_at,
			updated_at
		FROM tracks
		WHERE is_public = TRUE
		ORDER BY RANDOM()
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}

	for i := range items {
		if err := r.hydrate(ctx, &items[i]); err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int, publicOnly bool) ([]Track, error) {
	var items []Track

	query := `
		SELECT
			id,
			artist_id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			audio_media_id,
			cover_media_id,
			play_count,
			is_public,
			created_at,
			updated_at
		FROM tracks
	`

	args := []any{}
	if publicOnly {
		query += ` WHERE is_public = TRUE`
	}

	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, err
	}

	for i := range items {
		if err := r.hydrate(ctx, &items[i]); err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Track, error) {
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

	var item Track
	err = tx.GetContext(ctx, &item, `
		UPDATE tracks
		SET
			album_id = COALESCE($2, album_id),
			title = COALESCE($3, title),
			slug = $4,
			duration_seconds = COALESCE($5, duration_seconds),
			track_number = COALESCE($6, track_number),
			explicit = COALESCE($7, explicit),
			audio_url = COALESCE($8, audio_url),
			cover_url = COALESCE($9, cover_url),
			audio_media_id = COALESCE($10, audio_media_id),
			cover_media_id = COALESCE($11, cover_media_id),
			is_public = COALESCE($12, is_public),
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			artist_id,
			album_id,
			title,
			slug,
			duration_seconds,
			track_number,
			explicit,
			audio_url,
			cover_url,
			audio_media_id,
			cover_media_id,
			play_count,
			is_public,
			created_at,
			updated_at
	`,
		id,
		req.AlbumID,
		req.Title,
		slug,
		req.DurationSeconds,
		req.TrackNumber,
		req.Explicit,
		req.AudioURL,
		req.CoverURL,
		req.AudioMediaID,
		req.CoverMediaID,
		req.IsPublic,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.MapPGError(err)
	}

	if req.Artists != nil {
		artists := normalizeTrackArtists(req.Artists, uuid.Nil)
		if len(artists) == 0 {
			return nil, common.ErrInvalidInput
		}

		if err := r.replaceArtistsTx(ctx, tx, id, artists); err != nil {
			return nil, err
		}

		newPrimary := primaryArtistID(artists, uuid.Nil)
		if newPrimary != uuid.Nil {
			if _, err := tx.ExecContext(ctx, `
				UPDATE tracks
				SET artist_id = $2, updated_at = NOW()
				WHERE id = $1
			`, id, newPrimary); err != nil {
				return nil, common.MapPGError(err)
			}
		}
	}

	if req.GenreIDs != nil {
		if err := r.replaceGenresTx(ctx, tx, id, req.GenreIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, item.ID)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tracks WHERE id = $1`, id)
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

func (r *Repository) hydrate(ctx context.Context, item *Track) error {
	artists, err := r.listArtistsByTrack(ctx, item.ID)
	if err != nil {
		return err
	}
	item.Artists = artists

	genres, err := r.listGenresByTrack(ctx, item.ID)
	if err != nil {
		return err
	}
	item.Genres = genres

	return nil
}

func (r *Repository) replaceArtistsTx(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID, artists []TrackArtistRequest) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM track_artists WHERE track_id = $1`, trackID); err != nil {
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
			INSERT INTO track_artists (
				track_id,
				artist_id,
				role,
				position
			)
			VALUES ($1, $2, $3, $4)
		`, trackID, a.ArtistID, role, position); err != nil {
			return common.MapPGError(err)
		}
	}

	return nil
}

func (r *Repository) listArtistsByTrack(ctx context.Context, trackID uuid.UUID) ([]TrackArtist, error) {
	var items []TrackArtist
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			a.id AS artist_id,
			a.name,
			a.slug,
			ta.role,
			ta.position
		FROM track_artists ta
		INNER JOIN artists a ON a.id = ta.artist_id
		WHERE ta.track_id = $1
		ORDER BY ta.position ASC, ta.role ASC, a.name ASC
	`, trackID)

	return items, err
}

func (r *Repository) replaceGenresTx(ctx context.Context, tx *sqlx.Tx, trackID uuid.UUID, genreIDs []uuid.UUID) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM track_genres WHERE track_id = $1`, trackID); err != nil {
		return common.MapPGError(err)
	}

	for _, gid := range genreIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO track_genres (track_id, genre_id)
			VALUES ($1, $2)
		`, trackID, gid); err != nil {
			return common.MapPGError(err)
		}
	}

	return nil
}

func (r *Repository) listGenresByTrack(ctx context.Context, trackID uuid.UUID) ([]Genre, error) {
	var items []Genre
	err := r.db.SelectContext(ctx, &items, `
		SELECT g.id, g.name, g.slug, g.created_at
		FROM genres g
		INNER JOIN track_genres tg ON tg.genre_id = g.id
		WHERE tg.track_id = $1
		ORDER BY g.name ASC
	`, trackID)

	return items, err
}
func normalizeTrackCredits(credits []TrackCreditRequest) []TrackCreditRequest {
	for i := range credits {
		if credits[i].Position <= 0 {
			credits[i].Position = i + 1
		}
	}

	return credits
}

func (r *Repository) ListCredits(ctx context.Context, trackID uuid.UUID) ([]TrackCredit, error) {
	var items []TrackCredit

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			tc.id,
			tc.track_id,
			tc.artist_id,
			a.name AS artist_name,
			a.slug AS artist_slug,
			tc.credit_type,
			tc.position,
			tc.created_at
		FROM track_credits tc
		INNER JOIN artists a ON a.id = tc.artist_id
		WHERE tc.track_id = $1
		ORDER BY tc.position ASC, tc.credit_type ASC, a.name ASC
	`, trackID)

	return items, err
}

func (r *Repository) ReplaceCredits(ctx context.Context, trackID uuid.UUID, credits []TrackCreditRequest) ([]TrackCredit, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM tracks
			WHERE id = $1
		)
	`, trackID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM track_credits
		WHERE track_id = $1
	`, trackID); err != nil {
		return nil, common.MapPGError(err)
	}

	credits = normalizeTrackCredits(credits)

	for _, c := range credits {
		if c.ArtistID == uuid.Nil || c.CreditType == "" {
			return nil, common.ErrInvalidInput
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO track_credits (
				track_id,
				artist_id,
				credit_type,
				position
			)
			VALUES ($1, $2, $3, $4)
		`,
			trackID,
			c.ArtistID,
			c.CreditType,
			c.Position,
		); err != nil {
			return nil, common.MapPGError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.ListCredits(ctx, trackID)
}
func (r *Repository) ListArtists(ctx context.Context, trackID uuid.UUID) ([]TrackArtist, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM tracks
			WHERE id = $1
		)
	`, trackID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	return r.listArtistsByTrack(ctx, trackID)
}

func (r *Repository) ReplaceArtists(ctx context.Context, trackID uuid.UUID, artists []TrackArtistRequest) ([]TrackArtist, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM tracks
			WHERE id = $1
		)
	`, trackID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	artists = normalizeTrackArtists(artists, uuid.Nil)
	if len(artists) == 0 {
		return nil, common.ErrInvalidInput
	}

	for _, a := range artists {
		if a.ArtistID == uuid.Nil {
			return nil, common.ErrInvalidInput
		}
	}

	if err := r.replaceArtistsTx(ctx, tx, trackID, artists); err != nil {
		return nil, err
	}

	newPrimary := primaryArtistID(artists, uuid.Nil)
	if newPrimary != uuid.Nil {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tracks
			SET artist_id = $2, updated_at = NOW()
			WHERE id = $1
		`, trackID, newPrimary); err != nil {
			return nil, common.MapPGError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.listArtistsByTrack(ctx, trackID)
}
