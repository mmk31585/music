package video

import (
	"context"
	"database/sql"
	"errors"
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

// ---------------------------------------------------------------------------
// Video CRUD
// ---------------------------------------------------------------------------

func (r *Repository) Create(ctx context.Context, v *Video) error {
	v.ID = uuid.New()
	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO videos (
			id, track_id, uploader_id, type, status,
			title, description,
			raw_video_path, final_video_path, thumbnail_path,
			duration_ms, aspect_ratio, file_size_bytes,
			track_start_ms, track_end_ms,
			view_count, like_count, is_public, is_approved,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7,
			$8, $9, $10,
			$11, $12, $13,
			$14, $15,
			0, 0, $16, $17,
			$18, $19
		)
	`,
		v.ID, v.TrackID, v.UploaderID, string(v.Type), string(v.Status),
		v.Title, v.Description,
		v.RawVideoPath, v.FinalVideoPath, v.ThumbnailPath,
		v.DurationMs, v.AspectRatio, v.FileSizeBytes,
		v.TrackStartMs, v.TrackEndMs,
		v.IsPublic, v.IsApproved,
		v.CreatedAt, v.UpdatedAt,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Video, error) {
	var v Video
	err := r.db.GetContext(ctx, &v, `
		SELECT id, track_id, uploader_id, type, status,
			title, description,
			raw_video_path, final_video_path, thumbnail_path,
			duration_ms, aspect_ratio, file_size_bytes,
			track_start_ms, track_end_ms,
			view_count, like_count, is_public, is_approved,
			created_at, updated_at
		FROM videos
		WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) Update(ctx context.Context, v *Video) error {
	v.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE videos SET
			status = $2,
			title = $3,
			description = $4,
			raw_video_path = $5,
			final_video_path = $6,
			thumbnail_path = $7,
			duration_ms = $8,
			aspect_ratio = $9,
			file_size_bytes = $10,
			track_start_ms = $11,
			track_end_ms = $12,
			view_count = $13,
			like_count = $14,
			is_public = $15,
			is_approved = $16,
			updated_at = $17
		WHERE id = $1
	`,
		v.ID, string(v.Status), v.Title, v.Description,
		v.RawVideoPath, v.FinalVideoPath, v.ThumbnailPath,
		v.DurationMs, v.AspectRatio, v.FileSizeBytes,
		v.TrackStartMs, v.TrackEndMs,
		v.ViewCount, v.LikeCount, v.IsPublic, v.IsApproved,
		v.UpdatedAt,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM videos WHERE id = $1`, id)
	return err
}

func (r *Repository) ListByTrack(ctx context.Context, trackID uuid.UUID) ([]Video, error) {
	var items []Video
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, track_id, uploader_id, type, status,
			title, description,
			raw_video_path, final_video_path, thumbnail_path,
			duration_ms, aspect_ratio, file_size_bytes,
			track_start_ms, track_end_ms,
			view_count, like_count, is_public, is_approved,
			created_at, updated_at
		FROM videos
		WHERE track_id = $1
		ORDER BY created_at DESC
	`, trackID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

func (r *Repository) ListExplore(ctx context.Context, limit, offset int) ([]Video, error) {
	var items []Video
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, track_id, uploader_id, type, status,
			title, description,
			raw_video_path, final_video_path, thumbnail_path,
			duration_ms, aspect_ratio, file_size_bytes,
			track_start_ms, track_end_ms,
			view_count, like_count, is_public, is_approved,
			created_at, updated_at
		FROM videos
		WHERE is_approved = true AND is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

// ---------------------------------------------------------------------------
// Video Likes
// ---------------------------------------------------------------------------

func (r *Repository) CreateLike(ctx context.Context, like *VideoLike) error {
	like.ID = uuid.New()
	like.CreatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO video_likes (id, video_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (video_id, user_id) DO NOTHING
	`, like.ID, like.VideoID, like.UserID, like.CreatedAt)
	return err
}

func (r *Repository) DeleteLike(ctx context.Context, videoID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM video_likes
		WHERE video_id = $1 AND user_id = $2
	`, videoID, userID)
	return err
}

func (r *Repository) GetLike(ctx context.Context, videoID, userID uuid.UUID) (*VideoLike, error) {
	var like VideoLike
	err := r.db.GetContext(ctx, &like, `
		SELECT id, video_id, user_id, created_at
		FROM video_likes
		WHERE video_id = $1 AND user_id = $2
	`, videoID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *Repository) IncrementLikeCount(ctx context.Context, videoID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE videos SET like_count = like_count + 1, updated_at = NOW()
		WHERE id = $1
	`, videoID)
	return err
}

func (r *Repository) DecrementLikeCount(ctx context.Context, videoID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE videos SET like_count = GREATEST(like_count - 1, 0), updated_at = NOW()
		WHERE id = $1
	`, videoID)
	return err
}

func (r *Repository) IncrementViewCount(ctx context.Context, videoID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE videos SET view_count = view_count + 1, updated_at = NOW()
		WHERE id = $1
	`, videoID)
	return err
}

func (r *Repository) ListAll(ctx context.Context, limit, offset int) ([]Video, error) {
	var items []Video
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, track_id, uploader_id, type, status,
			title, description,
			raw_video_path, final_video_path, thumbnail_path,
			duration_ms, aspect_ratio, file_size_bytes,
			track_start_ms, track_end_ms,
			view_count, like_count, is_public, is_approved,
			created_at, updated_at
		FROM videos
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID, viewerID uuid.UUID, limit, offset int) ([]Video, error) {
	var items []Video
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, track_id, uploader_id, type, status,
			title, description,
			raw_video_path, final_video_path, thumbnail_path,
			duration_ms, aspect_ratio, file_size_bytes,
			track_start_ms, track_end_ms,
			view_count, like_count, is_public, is_approved,
			created_at, updated_at
		FROM videos
		WHERE uploader_id = $1
		  AND (is_public = true OR is_approved = true OR uploader_id = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, viewerID, limit, offset)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

// ---------------------------------------------------------------------------
//  Like visibility
// ---------------------------------------------------------------------------

// CheckTrackLikeExists verifies that the user has actually liked this track
// (via the polymorphic reactions table) before allowing visibility changes.
func (r *Repository) CheckTrackLikeExists(ctx context.Context, userID, trackID uuid.UUID) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
		FROM reactions
		WHERE user_id = $1 AND target_id = $2 AND target_type = 'track' AND type = 'like'
	`, userID, trackID.String())
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) UpsertTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID, visibility string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO track_like_visibility (user_id, track_id, visibility, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, track_id)
		DO UPDATE SET visibility = $3
	`, userID, trackID, visibility)
	return err
}

func (r *Repository) GetTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID) (*TrackLikeVisibility, error) {
	var tlv TrackLikeVisibility
	err := r.db.GetContext(ctx, &tlv, `
		SELECT user_id, track_id, visibility, created_at
		FROM track_like_visibility
		WHERE user_id = $1 AND track_id = $2
	`, userID, trackID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tlv, nil
}

func (r *Repository) GetPublicLikedTracks(ctx context.Context, userID uuid.UUID, limit, offset int) ([]TrackLikeVisibility, error) {
	var items []TrackLikeVisibility
	err := r.db.SelectContext(ctx, &items, `
		SELECT user_id, track_id, visibility, created_at
		FROM track_like_visibility
		WHERE user_id = $1 AND visibility = 'public'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []TrackLikeVisibility{}
	}
	return items, nil
}

// ---------------------------------------------------------------------------
// User Music Status
// ---------------------------------------------------------------------------

func (r *Repository) UpsertMusicStatus(ctx context.Context, status *UserMusicStatus) error {
	status.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_music_status (user_id, current_track_id, visibility, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET current_track_id = $2, visibility = $3, updated_at = $4
	`, status.UserID, status.CurrentTrackID, status.Visibility, status.UpdatedAt)
	return err
}

func (r *Repository) DeleteMusicStatus(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_music_status WHERE user_id = $1
	`, userID)
	return err
}

func (r *Repository) GetMusicStatus(ctx context.Context, userID uuid.UUID) (*UserMusicStatus, error) {
	var s UserMusicStatus
	err := r.db.GetContext(ctx, &s, `
		SELECT user_id, current_track_id, visibility, updated_at
		FROM user_music_status
		WHERE user_id = $1
	`, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetTrackCoverURL returns the cover_url for a track given its ID.
func (r *Repository) GetTrackCoverURL(ctx context.Context, trackID uuid.UUID) (*string, error) {
	var coverURL *string
	err := r.db.GetContext(ctx, &coverURL, `
		SELECT cover_url FROM tracks WHERE id = $1
	`, trackID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return coverURL, nil
}

// GetTrackCovers returns a map of trackID -> cover_url for the given slice of track IDs.
func (r *Repository) GetTrackCovers(ctx context.Context, trackIDs []uuid.UUID) (map[uuid.UUID]*string, error) {
	if len(trackIDs) == 0 {
		return map[uuid.UUID]*string{}, nil
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, cover_url FROM tracks WHERE id = ANY($1)
	`, trackIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*string, len(trackIDs))
	for rows.Next() {
		var id uuid.UUID
		var coverURL *string
		if err := rows.Scan(&id, &coverURL); err != nil {
			return nil, err
		}
		result[id] = coverURL
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
