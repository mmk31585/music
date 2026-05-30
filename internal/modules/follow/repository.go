package follow

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrArtistNotFound = errors.New("artist not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	exists, err := r.userExists(ctx, followeeID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO user_follows (follower_id, followee_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, followee_id) DO NOTHING
	`, followerID, followeeID)

	return err
}

func (r *Repository) UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_follows
		WHERE follower_id = $1 AND followee_id = $2
	`, followerID, followeeID)

	return err
}

func (r *Repository) FollowArtist(ctx context.Context, userID, artistID uuid.UUID) error {
	exists, err := r.artistExists(ctx, artistID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrArtistNotFound
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO artist_follows (user_id, artist_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, artist_id) DO NOTHING
	`, userID, artistID)

	return err
}

func (r *Repository) UnfollowArtist(ctx context.Context, userID, artistID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM artist_follows
		WHERE user_id = $1 AND artist_id = $2
	`, userID, artistID)

	return err
}

func (r *Repository) IsFollowingUser(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM user_follows
			WHERE follower_id = $1 AND followee_id = $2
		)
	`, followerID, followeeID)
	return exists, err
}

func (r *Repository) IsFollowingArtist(ctx context.Context, userID, artistID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM artist_follows
			WHERE user_id = $1 AND artist_id = $2
		)
	`, userID, artistID)
	return exists, err
}

func (r *Repository) userExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`, userID)
	return exists, err
}

func (r *Repository) artistExists(ctx context.Context, artistID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM artists
			WHERE id = $1
		)
	`, artistID)
	return exists, err
}
