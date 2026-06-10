package reactions

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, userID uuid.UUID, targetID, targetType, reactionType string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO reactions (user_id, target_id, target_type, type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, target_id, target_type)
		DO UPDATE SET type = $4, created_at = NOW()
	`, userID, targetID, targetType, reactionType)
	return err
}

func (r *Repository) Remove(ctx context.Context, userID uuid.UUID, targetID, targetType string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM reactions
		WHERE user_id = $1 AND target_id = $2 AND target_type = $3
	`, userID, targetID, targetType)
	return err
}

func (r *Repository) GetUserReaction(ctx context.Context, userID uuid.UUID, targetID, targetType string) (*Reaction, error) {
	var reaction Reaction
	err := r.db.GetContext(ctx, &reaction, `
		SELECT id, user_id, target_id, target_type, type, created_at
		FROM reactions
		WHERE user_id = $1 AND target_id = $2 AND target_type = $3
		LIMIT 1
	`, userID, targetID, targetType)
	if err != nil {
		return nil, err
	}
	return &reaction, nil
}

func (r *Repository) GetCounts(ctx context.Context, targetID, targetType string) ([]ReactionCount, error) {
	var counts []ReactionCount
	err := r.db.SelectContext(ctx, &counts, `
		SELECT target_id, target_type, type, COUNT(*) as count
		FROM reactions
		WHERE target_id = $1 AND target_type = $2
		GROUP BY target_id, target_type, type
	`, targetID, targetType)
	return counts, err
}

func (r *Repository) GetUserReactions(ctx context.Context, userID uuid.UUID, targetType string, limit, offset int) ([]Reaction, error) {
	var items []Reaction
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, user_id, target_id, target_type, type, created_at
		FROM reactions
		WHERE user_id = $1 AND target_type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, targetType, limit, offset)
	return items, err
}

func (r *Repository) BulkGetCounts(ctx context.Context, targetIDs []string, targetType string) (map[string]int, error) {
	query := `
		SELECT target_id, COUNT(*) as count
		FROM reactions
		WHERE target_id = ANY($1) AND target_type = $2
		GROUP BY target_id
	`
	rows, err := r.db.QueryContext(ctx, query, targetIDs, targetType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var targetID string
		var count int
		if err := rows.Scan(&targetID, &count); err != nil {
			return nil, err
		}
		result[targetID] = count
	}
	return result, rows.Err()
}
