package artist

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Artist, error) {
	slug := common.Slugify(req.Name)

	var item Artist
	err := r.db.GetContext(ctx, &item, `
		INSERT INTO artists (
			id,
			name,
			slug,
			bio,
			image_url,
			avatar_media_id,
			banner_media_id,
			country,
			is_verified,
			monthly_listeners
		)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, COALESCE($8, FALSE), COALESCE($9, 0))
		RETURNING
			id,
			name,
			slug,
			bio,
			image_url,
			avatar_media_id,
			banner_media_id,
			country,
			is_verified,
			monthly_listeners,
			created_at,
			updated_at
	`,
		req.Name,
		slug,
		req.Bio,
		req.ImageURL,
		req.AvatarMediaID,
		req.BannerMediaID,
		req.Country,
		req.IsVerified,
		req.MonthlyListeners,
	)
	if err != nil {
		return nil, common.MapPGError(err)
	}

	return &item, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	var item Artist
	err := r.db.GetContext(ctx, &item, `
		SELECT
			id,
			name,
			slug,
			bio,
			image_url,
			avatar_media_id,
			banner_media_id,
			country,
			is_verified,
			monthly_listeners,
			created_at,
			updated_at
		FROM artists
		WHERE id = $1
	`, id)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Artist, error) {
	var items []Artist
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			name,
			slug,
			bio,
			image_url,
			avatar_media_id,
			banner_media_id,
			country,
			is_verified,
			monthly_listeners,
			created_at,
			updated_at
		FROM artists
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	return items, err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Artist, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := current.Name
	if req.Name != nil {
		name = *req.Name
	}
	slug := common.Slugify(name)

	var item Artist
	err = r.db.GetContext(ctx, &item, `
		UPDATE artists
		SET
			name = COALESCE($2, name),
			slug = $3,
			bio = COALESCE($4, bio),
			image_url = COALESCE($5, image_url),
			avatar_media_id = COALESCE($6, avatar_media_id),
			banner_media_id = COALESCE($7, banner_media_id),
			country = COALESCE($8, country),
			is_verified = COALESCE($9, is_verified),
			monthly_listeners = COALESCE($10, monthly_listeners),
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			name,
			slug,
			bio,
			image_url,
			avatar_media_id,
			banner_media_id,
			country,
			is_verified,
			monthly_listeners,
			created_at,
			updated_at
	`,
		id,
		req.Name,
		slug,
		req.Bio,
		req.ImageURL,
		req.AvatarMediaID,
		req.BannerMediaID,
		req.Country,
		req.IsVerified,
		req.MonthlyListeners,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.MapPGError(err)
	}

	return &item, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM artists WHERE id = $1`, id)
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

func (r *Repository) Search(ctx context.Context, q string, limit, offset int) ([]Artist, error) {
	var items []Artist
	pattern := fmt.Sprintf("%%%s%%", q)

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			id,
			name,
			slug,
			bio,
			image_url,
			avatar_media_id,
			banner_media_id,
			country,
			is_verified,
			monthly_listeners,
			created_at,
			updated_at
		FROM artists
		WHERE name ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, pattern, limit, offset)

	return items, err
}
func (r *Repository) ListTracks(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]ArtistTrack, error) {
	var items []ArtistTrack

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			t.id,
			t.title,
			t.slug,
			t.album_id,
			t.duration_seconds,
			t.track_number,
			t.explicit,
			t.audio_url,
			t.cover_url,
			t.play_count,
			ta.role AS artist_role,
			ta.position AS artist_position,
			t.created_at
		FROM track_artists ta
		INNER JOIN tracks t ON t.id = ta.track_id
		WHERE ta.artist_id = $1
		  AND t.is_public = TRUE
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`, artistID, limit, offset)

	return items, err
}

func (r *Repository) ListAppearsOn(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]ArtistTrack, error) {
	var items []ArtistTrack

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			t.id,
			t.title,
			t.slug,
			t.album_id,
			t.duration_seconds,
			t.track_number,
			t.explicit,
			t.audio_url,
			t.cover_url,
			t.play_count,
			ta.role AS artist_role,
			ta.position AS artist_position,
			t.created_at
		FROM track_artists ta
		INNER JOIN tracks t ON t.id = ta.track_id
		WHERE ta.artist_id = $1
		  AND ta.role <> 'primary'
		  AND t.is_public = TRUE
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`, artistID, limit, offset)

	return items, err
}

func (r *Repository) ListAlbumsByType(ctx context.Context, artistID uuid.UUID, albumTypes []string, limit, offset int) ([]ArtistAlbum, error) {
	var items []ArtistAlbum

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			al.id,
			al.title,
			al.slug,
			al.cover_url,
			al.release_date,
			al.album_type,
			aa.role AS artist_role,
			al.created_at
		FROM album_artists aa
		INNER JOIN albums al ON al.id = aa.album_id
		WHERE aa.artist_id = $1
		  AND al.album_type = ANY($2)
		ORDER BY al.release_date DESC NULLS LAST, al.created_at DESC
		LIMIT $3 OFFSET $4
	`, artistID, albumTypes, limit, offset)

	return items, err
}

func (r *Repository) GetOverview(ctx context.Context, artistID uuid.UUID) (*Overview, error) {
	artist, err := r.GetByID(ctx, artistID)
	if err != nil {
		return nil, err
	}

	topTracks, err := r.ListTopTracks(ctx, artistID, 10, 0)
	if err != nil {
		return nil, err
	}

	albums, err := r.ListAlbumsByType(ctx, artistID, []string{"album", "ep", "compilation", "live"}, 10, 0)
	if err != nil {
		return nil, err
	}

	singles, err := r.ListAlbumsByType(ctx, artistID, []string{"single"}, 10, 0)
	if err != nil {
		return nil, err
	}

	appearsOn, err := r.ListAppearsOn(ctx, artistID, 10, 0)
	if err != nil {
		return nil, err
	}

	return &Overview{
		Artist:    artist,
		TopTracks: topTracks,
		Albums:    albums,
		Singles:   singles,
		AppearsOn: appearsOn,
	}, nil
}
func (r *Repository) ListTopTracks(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]ArtistTrack, error) {
	var items []ArtistTrack

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			t.id,
			t.title,
			t.slug,
			t.album_id,
			t.duration_seconds,
			t.track_number,
			t.explicit,
			t.audio_url,
			t.cover_url,
			t.play_count,
			COALESCE(ta.role, 'primary') AS artist_role,
			COALESCE(ta.position, att.position) AS artist_position,
			t.created_at
		FROM artist_top_tracks att
		INNER JOIN tracks t ON t.id = att.track_id
		LEFT JOIN track_artists ta
			ON ta.track_id = t.id
		   AND ta.artist_id = att.artist_id
		WHERE att.artist_id = $1
		  AND t.is_public = TRUE
		ORDER BY att.position ASC, t.play_count DESC, t.created_at DESC
		LIMIT $2 OFFSET $3
	`, artistID, limit, offset)

	return items, err
}
func normalizeArtistTopTracks(tracks []ArtistTopTrackRequest) []ArtistTopTrackRequest {
	for i := range tracks {
		if tracks[i].Position <= 0 {
			tracks[i].Position = i + 1
		}
	}

	return tracks
}

func (r *Repository) ReplaceTopTracks(ctx context.Context, artistID uuid.UUID, tracks []ArtistTopTrackRequest) ([]ArtistTrack, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM artists
			WHERE id = $1
		)
	`, artistID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	tracks = normalizeArtistTopTracks(tracks)

	for _, item := range tracks {
		if item.TrackID == uuid.Nil {
			return nil, common.ErrInvalidInput
		}
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM artist_top_tracks
		WHERE artist_id = $1
	`, artistID); err != nil {
		return nil, common.MapPGError(err)
	}

	for _, item := range tracks {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO artist_top_tracks (
				artist_id,
				track_id,
				position
			)
			VALUES ($1, $2, $3)
		`,
			artistID,
			item.TrackID,
			item.Position,
		); err != nil {
			return nil, common.MapPGError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.ListTopTracks(ctx, artistID, 50, 0)
}

func (r *Repository) ListRelated(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]RelatedArtist, error) {
	var items []RelatedArtist

	err := r.db.SelectContext(ctx, &items, `
		SELECT
			a.id,
			a.name,
			a.slug,
			a.bio,
			a.image_url,
			a.avatar_media_id,
			a.banner_media_id,
			a.is_verified,
			a.monthly_listeners,
			ra.score AS related_score,
			ra.source AS related_source,
			ra.created_at AS related_at
		FROM related_artists ra
		INNER JOIN artists a ON a.id = ra.related_artist_id
		WHERE ra.artist_id = $1
		ORDER BY ra.score DESC, a.monthly_listeners DESC, a.name ASC
		LIMIT $2 OFFSET $3
	`, artistID, limit, offset)

	return items, err
}

func (r *Repository) ReplaceRelated(ctx context.Context, artistID uuid.UUID, related []RelatedArtistRequest) ([]RelatedArtist, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM artists
			WHERE id = $1
		)
	`, artistID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM related_artists
		WHERE artist_id = $1
	`, artistID); err != nil {
		return nil, common.MapPGError(err)
	}

	for _, item := range related {
		if item.RelatedArtistID == uuid.Nil || item.RelatedArtistID == artistID {
			return nil, common.ErrInvalidInput
		}

		source := item.Source
		if source == "" {
			source = "manual"
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO related_artists (
				artist_id,
				related_artist_id,
				score,
				source
			)
			VALUES ($1, $2, $3, $4)
		`,
			artistID,
			item.RelatedArtistID,
			item.Score,
			source,
		); err != nil {
			return nil, common.MapPGError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.ListRelated(ctx, artistID, 50, 0)
}
