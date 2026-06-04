package artist

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	Name             string     `db:"name" json:"name"`
	Slug             string     `db:"slug" json:"slug"`
	Bio              *string    `db:"bio" json:"bio,omitempty"`
	ImageURL         *string    `db:"image_url" json:"imageUrl,omitempty"`
	IsVerified       bool       `db:"is_verified" json:"isVerified"`
	MonthlyListeners int64      `db:"monthly_listeners" json:"monthlyListeners"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type CreateRequest struct {
	Name             string  `db:"name" json:"name" validate:"required,min=1,max=200"`
	Bio              *string `db:"bio" json:"bio" validate:"omitempty,max=5000"`
	ImageURL         *string `db:"image_url" json:"image_url" validate:"omitempty,url"`
	IsVerified       *bool   `db:"is_verified" json:"is_verified"`
	MonthlyListeners *int64  `db:"monthly_listeners" json:"monthly_listeners" validate:"omitempty,min=0"`
}

type UpdateRequest struct {
	Name             *string `db:"name" json:"name" validate:"omitempty,min=1,max=200"`
	Bio              *string `db:"bio" json:"bio" validate:"omitempty,max=5000"`
	ImageURL         *string `db:"image_url" json:"image_url" validate:"omitempty,url"`
	IsVerified       *bool   `db:"is_verified" json:"is_verified"`
	MonthlyListeners *int64  `db:"monthly_listeners" json:"monthly_listeners" validate:"omitempty,min=0"`
}
