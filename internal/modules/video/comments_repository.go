package video

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// CommentRepository handles DB queries for video comments.
type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ListByVideo returns all comments for a video, newest first, limited.
func (r *CommentRepository) ListByVideo(ctx context.Context, videoID uuid.UUID, limit, offset int) ([]VideoComment, int, error) {
	// Count
	var total int
	err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM video_comments WHERE video_id = $1`, videoID)
	if err != nil {
		return nil, 0, err
	}

	var comments []VideoComment
	err = r.db.SelectContext(ctx, &comments, `
		SELECT id, video_id, user_id, content, created_at, updated_at
		FROM video_comments
		WHERE video_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, videoID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if comments == nil {
		comments = []VideoComment{}
	}
	return comments, total, nil
}

// Create inserts a new comment.
func (r *CommentRepository) Create(ctx context.Context, videoID, userID uuid.UUID, content string) (*VideoComment, error) {
	var c VideoComment
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO video_comments (video_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, video_id, user_id, content, created_at, updated_at
	`, videoID, userID, content).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetUserDisplayName fetches the display name or username for a user ID.
func (r *CommentRepository) GetUserDisplayName(ctx context.Context, userID uuid.UUID) (string, error) {
	var name sql.NullString
	err := r.db.GetContext(ctx, &name, `
		SELECT COALESCE(display_name, username) FROM users WHERE id = $1
	`, userID)
	if err != nil {
		return "", err
	}
	if name.Valid {
		return name.String, nil
	}
	return "Unknown", nil
}

// AdminListByVideo returns all comments for a video with author name, ordered newest first.
func (r *CommentRepository) AdminListByVideo(ctx context.Context, videoID uuid.UUID, limit, offset int) ([]AdminCommentItem, int, error) {
	var total int
	err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM video_comments WHERE video_id = $1`, videoID)
	if err != nil {
		return nil, 0, err
	}

	var items []AdminCommentItem
	err = r.db.SelectContext(ctx, &items, `
		SELECT
			c.id, c.video_id, c.user_id, c.content,
			c.created_at, c.updated_at,
			COALESCE(u.display_name, u.username, 'Unknown') AS author_name
		FROM video_comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.video_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, videoID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = []AdminCommentItem{}
	}
	return items, total, nil
}

// UpdateComment updates a comment's content.
func (r *CommentRepository) UpdateComment(ctx context.Context, id uuid.UUID, content string) (*VideoComment, error) {
	var c VideoComment
	err := r.db.QueryRowxContext(ctx, `
		UPDATE video_comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, video_id, user_id, content, created_at, updated_at
	`, content, id).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// DeleteComment removes a comment by id.
func (r *CommentRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM video_comments WHERE id = $1`, id)
	return err
}
