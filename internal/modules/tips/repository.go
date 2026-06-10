package tips

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, tip Tip) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO tips (id, sender_id, artist_id, track_id, amount_cents, currency, message, status)
		VALUES (:id, :sender_id, :artist_id, :track_id, :amount_cents, :currency, :message, :status)
	`, tip)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Tip, error) {
	var tip Tip
	err := r.db.GetContext(ctx, &tip, `
		SELECT id, sender_id, artist_id, track_id, amount_cents, currency, message, status, provider, provider_pay_id, paid_at, created_at
		FROM tips WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &tip, err
}

func (r *Repository) MarkPaid(ctx context.Context, id, provider, providerPayID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tips SET status = 'completed', provider = $2, provider_pay_id = $3, paid_at = NOW()
		WHERE id = $1
	`, id, provider, providerPayID)
	return err
}

func (r *Repository) ListBySender(ctx context.Context, senderID string, limit, offset int) ([]Tip, error) {
	var items []Tip
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, sender_id, artist_id, track_id, amount_cents, currency, message, status, provider, provider_pay_id, paid_at, created_at
		FROM tips WHERE sender_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, senderID, limit, offset)
	return items, err
}

func (r *Repository) ListByArtist(ctx context.Context, artistID string, limit, offset int) ([]Tip, error) {
	var items []Tip
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, sender_id, artist_id, track_id, amount_cents, currency, message, status, provider, provider_pay_id, paid_at, created_at
		FROM tips WHERE artist_id = $1 AND status = 'completed'
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, artistID, limit, offset)
	return items, err
}

func (r *Repository) TotalForArtist(ctx context.Context, artistID string) (int64, error) {
	var total int64
	err := r.db.GetContext(ctx, &total, `
		SELECT COALESCE(SUM(amount_cents), 0) FROM tips WHERE artist_id = $1 AND status = 'completed'
	`, artistID)
	return total, err
}
