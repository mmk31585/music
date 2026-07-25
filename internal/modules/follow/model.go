package follow

import (
	"time"

	"github.com/google/uuid"
)

type UserFollow struct {
	FollowerID uuid.UUID `db:"follower_id" json:"follower_id"`
	FolloweeID uuid.UUID `db:"followee_id" json:"followee_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type ArtistFollow struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	ArtistID  uuid.UUID `db:"artist_id" json:"artist_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
