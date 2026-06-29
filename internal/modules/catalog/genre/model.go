package genre

import (
	"time"

	"github.com/google/uuid"
)

type Genre struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateRequest struct {
	Name string `db:"name" json:"name" validate:"required,min=2,max=100"`
}

type UpdateRequest struct {
	Name *string `db:"name" json:"name" validate:"omitempty,min=2,max=100"`
}
