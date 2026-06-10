package tips

import "time"

type Tip struct {
	ID            string     `db:"id" json:"id"`
	SenderID      string     `db:"sender_id" json:"sender_id"`
	ArtistID      string     `db:"artist_id" json:"artist_id"`
	TrackID       *string    `db:"track_id" json:"track_id,omitempty"`
	AmountCents   int64      `db:"amount_cents" json:"amount_cents"`
	Currency      string     `db:"currency" json:"currency"`
	Message       *string    `db:"message" json:"message,omitempty"`
	Status        string     `db:"status" json:"status"`
	Provider      *string    `db:"provider" json:"provider,omitempty"`
	ProviderPayID *string    `db:"provider_pay_id" json:"provider_pay_id,omitempty"`
	PaidAt        *time.Time `db:"paid_at" json:"paid_at,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}
