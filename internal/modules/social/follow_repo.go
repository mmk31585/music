package social

import (
	"context"

	"github.com/google/uuid"
)

// --- Follow ---

func (r *repository) Follow(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_follows (follower_id, followee_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, followerID, followedID)
	return err
}

func (r *repository) Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_follows
		WHERE follower_id = $1 AND followee_id = $2
	`, followerID, followedID)
	return err
}

func (r *repository) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1 FROM user_follows
			WHERE follower_id = $1 AND followee_id = $2
		)
	`, followerID, followedID)
	return exists, err
}

func (r *repository) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM user_follows WHERE followee_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	var items []UserFollow
	if err := r.db.SelectContext(ctx, &items, `
		SELECT follower_id, followee_id, created_at
		FROM user_follows
		WHERE followee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	var items []UserFollow
	if err := r.db.SelectContext(ctx, &items, `
		SELECT follower_id, followee_id, created_at
		FROM user_follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) GetFollowerCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM user_follows WHERE followee_id = $1`, userID)
	return count, err
}

func (r *repository) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`, userID)
	return count, err
}

func (r *repository) InsertActivity(ctx context.Context, activity Activity) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO activities (id, user_id, type, target_id, target_type, metadata)
		VALUES (:id, :user_id, :type, :target_id, :target_type, :metadata)
	`, activity)
	return err
}

func (r *repository) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int, types []string) ([]ActivityFeedItem, error) {
	query := `
		SELECT
			a.id, a.user_id, a.type, a.target_id, a.target_type, a.metadata, a.created_at,
			u.display_name AS user_display_name,
			u.avatar_url AS user_avatar_url,
			COALESCE(t.title, al.title, ar.name, p.title, '') AS target_name,
			COALESCE(t.cover_url, al.cover_url, ar.image_url, p.cover_url, '') AS target_image_url
		FROM activities a
		JOIN users u ON u.id = a.user_id
		LEFT JOIN tracks t ON a.target_type = 'track' AND a.target_id = t.id::text
		LEFT JOIN albums al ON a.target_type = 'album' AND a.target_id = al.id::text
		LEFT JOIN artists ar ON a.target_type = 'artist' AND a.target_id = ar.id::text
		LEFT JOIN listening_parties p ON a.target_type = 'party' AND a.target_id = p.id::text
		WHERE a.user_id IN (
			SELECT followee_id FROM user_follows WHERE follower_id = $1
			UNION ALL
			SELECT $1
		)
	`

	args := []interface{}{userID, limit, offset}

	if len(types) > 0 {
		query += ` AND a.type = ANY($4) `
		args = append(args, types)
	}

	query += ` ORDER BY a.created_at DESC LIMIT $2 OFFSET $3 `

	items := make([]ActivityFeedItem, 0)
	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}
