package gamification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrBadgeNotFound     = errors.New("badge not found")
	ErrChallengeNotFound = errors.New("challenge not found")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserXP(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM xp_transactions WHERE user_id = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	return total, err
}

func (r *Repository) AddXP(ctx context.Context, userID string, amount int, source string, metadata *string) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	var currentTotal int64
	err = tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount), 0) FROM xp_transactions WHERE user_id = $1`, userID).Scan(&currentTotal)
	if err != nil {
		return 0, err
	}

	newBalance := currentTotal + int64(amount)

	_, err = tx.ExecContext(ctx,
		`INSERT INTO xp_transactions (user_id, amount, balance_after, source, metadata) VALUES ($1, $2, $3, $4, $5)`,
		userID, amount, newBalance, source, metadata,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert xp: %w", err)
	}

	return newBalance, tx.Commit()
}

func (r *Repository) GetUserBadges(ctx context.Context, userID string) ([]UserBadge, error) {
	query := `
		SELECT ub.user_id, ub.badge_id, ub.earned_at, ub.is_displayed,
			b.id, b.name, b.description, b.icon_url, b.category, b.rarity, b.criteria, b.xp_reward, b.is_hidden
		FROM user_badges ub
		JOIN badges b ON b.id = ub.badge_id
		WHERE ub.user_id = $1
		ORDER BY ub.earned_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []UserBadge
	for rows.Next() {
		var ub UserBadge
		var b Badge
		if err := rows.Scan(&ub.UserID, &ub.BadgeID, &ub.EarnedAt, &ub.IsDisplayed,
			&b.ID, &b.Name, &b.Description, &b.IconURL, &b.Category, &b.Rarity, &b.Criteria, &b.XPReward, &b.IsHidden,
		); err != nil {
			return nil, err
		}
		ub.Badge = &b
		items = append(items, ub)
	}
	return items, rows.Err()
}

func (r *Repository) AwardBadge(ctx context.Context, userID, badgeID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_badges (user_id, badge_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, badgeID,
	)
	return err
}

func (r *Repository) GetAllBadges(ctx context.Context) ([]Badge, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, COALESCE(description, ''), COALESCE(icon_url, ''), category, rarity, criteria, xp_reward, is_hidden FROM badges ORDER BY rarity, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Badge
	for rows.Next() {
		var b Badge
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IconURL, &b.Category, &b.Rarity, &b.Criteria, &b.XPReward, &b.IsHidden); err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, rows.Err()
}

func (r *Repository) GetActiveChallenges(ctx context.Context) ([]DailyChallenge, error) {
	query := `SELECT id, title, COALESCE(description, ''), challenge_type, target_count, xp_reward, badge_reward_id, is_active, valid_from, valid_until FROM daily_challenges WHERE is_active = true AND valid_from <= NOW() AND valid_until >= NOW() ORDER BY valid_until ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DailyChallenge
	for rows.Next() {
		var c DailyChallenge
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.ChallengeType, &c.TargetCount, &c.XPReward, &c.BadgeRewardID, &c.IsActive, &c.ValidFrom, &c.ValidUntil); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (r *Repository) GetUserChallenges(ctx context.Context, userID string) ([]UserChallenge, error) {
	query := `
		SELECT uc.user_id, uc.challenge_id, uc.progress, uc.is_completed, uc.completed_at,
			dc.id, dc.title, COALESCE(dc.description, ''), dc.challenge_type, dc.target_count, dc.xp_reward, dc.badge_reward_id, dc.is_active, dc.valid_from, dc.valid_until
		FROM user_challenges uc
		JOIN daily_challenges dc ON dc.id = uc.challenge_id
		WHERE uc.user_id = $1 AND dc.is_active = true AND dc.valid_from <= NOW() AND dc.valid_until >= NOW()
		ORDER BY dc.valid_until ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []UserChallenge
	for rows.Next() {
		var uc UserChallenge
		var dc DailyChallenge
		if err := rows.Scan(&uc.UserID, &uc.ChallengeID, &uc.Progress, &uc.IsCompleted, &uc.CompletedAt,
			&dc.ID, &dc.Title, &dc.Description, &dc.ChallengeType, &dc.TargetCount, &dc.XPReward, &dc.BadgeRewardID, &dc.IsActive, &dc.ValidFrom, &dc.ValidUntil,
		); err != nil {
			return nil, err
		}
		uc.Challenge = &dc
		items = append(items, uc)
	}
	return items, rows.Err()
}

func (r *Repository) UpsertChallengeProgress(ctx context.Context, userID, challengeID string, progress int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_challenges (user_id, challenge_id, progress) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, challenge_id) DO UPDATE SET progress = GREATEST(user_challenges.progress, $3)
	`, userID, challengeID, progress)
	return err
}

func (r *Repository) CompleteChallenge(ctx context.Context, userID, challengeID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_challenges SET is_completed = true, completed_at = NOW() WHERE user_id = $1 AND challenge_id = $2`,
		userID, challengeID,
	)
	return err
}

func (r *Repository) GetLeaderboard(ctx context.Context, lbType string, limit int) ([]LeaderboardEntry, error) {
	var query string
	switch lbType {
	case "xp_weekly":
		query = `
			SELECT u.id AS user_id, u.username, COALESCE(u.avatar_url, ''), COALESCE(SUM(xt.amount), 0) AS score
			FROM users u
			JOIN xp_transactions xt ON xt.user_id = u.id
			WHERE xt.created_at >= NOW() - INTERVAL '7 days'
			GROUP BY u.id, u.username, u.avatar_url
			ORDER BY score DESC LIMIT $1
		`
	case "xp_monthly":
		query = `
			SELECT u.id AS user_id, u.username, COALESCE(u.avatar_url, ''), COALESCE(SUM(xt.amount), 0) AS score
			FROM users u
			JOIN xp_transactions xt ON xt.user_id = u.id
			WHERE xt.created_at >= NOW() - INTERVAL '30 days'
			GROUP BY u.id, u.username, u.avatar_url
			ORDER BY score DESC LIMIT $1
		`
	case "streams":
		query = `
			SELECT u.id AS user_id, u.username, COALESCE(u.avatar_url, ''), COUNT(lh.id) AS score
			FROM users u
			JOIN listening_history lh ON lh.user_id = u.id
			WHERE lh.played_at >= NOW() - INTERVAL '7 days'
			GROUP BY u.id, u.username, u.avatar_url
			ORDER BY score DESC LIMIT $1
		`
	case "contributions":
		query = `
			SELECT u.id AS user_id, u.username, COALESCE(u.avatar_url, ''), COUNT(c.id) AS score
			FROM users u
			JOIN contributions c ON c.contributor_id = u.id
			WHERE c.status = 'approved' AND c.created_at >= NOW() - INTERVAL '30 days'
			GROUP BY u.id, u.username, u.avatar_url
			ORDER BY score DESC LIMIT $1
		`
	default:
		query = `
			SELECT u.id AS user_id, u.username, COALESCE(u.avatar_url, ''), COALESCE(SUM(xt.amount), 0) AS score
			FROM users u
			LEFT JOIN xp_transactions xt ON xt.user_id = u.id
			GROUP BY u.id, u.username, u.avatar_url
			ORDER BY score DESC LIMIT $1
		`
	}

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.AvatarURL, &e.Score); err != nil {
			return nil, err
		}
		e.Rank = rank
		rank++
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *Repository) GetUserRank(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*) + 1 FROM (
			SELECT user_id, SUM(amount) as total FROM xp_transactions GROUP BY user_id HAVING SUM(amount) > (
				SELECT COALESCE(SUM(amount), 0) FROM xp_transactions WHERE user_id = $1
			)
		) ranked
	`
	var rank int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&rank)
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (r *Repository) CheckAndAwardBadges(ctx context.Context, userID string) ([]Badge, error) {
	var newBadges []Badge

	badges, err := r.GetAllBadges(ctx)
	if err != nil {
		return nil, err
	}

	for _, badge := range badges {
		var earned bool
		err := r.db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM user_badges WHERE user_id = $1 AND badge_id = $2)`,
			userID, badge.ID,
		).Scan(&earned)
		if err != nil || earned {
			continue
		}

		// Check criteria using raw SQL for each badge type
		var met bool
		switch badge.Name {
		case "First Stream":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 1 FROM listening_history WHERE user_id = $1`, userID).Scan(&met)
		case "Night Owl":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 100 FROM listening_history WHERE user_id = $1 AND EXTRACT(HOUR FROM played_at)::int BETWEEN 0 AND 5`, userID).Scan(&met)
		case "Explorer":
			r.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT g.id) >= 20 FROM listening_history lh JOIN tracks t ON t.id = lh.track_id JOIN track_genres tg ON tg.track_id = t.id JOIN genres g ON g.id = tg.genre_id WHERE lh.user_id = $1`, userID).Scan(&met)
		case "Marathon":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 1000 FROM listening_history WHERE user_id = $1 AND played_at >= NOW() - INTERVAL '7 days'`, userID).Scan(&met)
		case "Wordsmith":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 10 FROM contributions WHERE contributor_id = $1 AND status = 'approved' AND contribution_type = 'lyrics'`, userID).Scan(&met)
		case "Polyglot":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 5 FROM contributions WHERE contributor_id = $1 AND status = 'approved' AND contribution_type = 'translation'`, userID).Scan(&met)
		case "Curator":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 10 FROM playlists WHERE owner_id = $1`, userID).Scan(&met)
		case "Networker":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 100 FROM follows WHERE followee_id = $1`, userID).Scan(&met)
		case "Community Hero":
			r.db.QueryRowContext(ctx, `SELECT COUNT(*) >= 50 FROM moderation_queue WHERE resolved_by = $1`, userID).Scan(&met)
		case "Chart Climber":
			rank, _ := r.GetUserRank(ctx, userID)
			met = rank <= 10
		}

		if met {
			if err := r.AwardBadge(ctx, userID, badge.ID); err == nil {
				newBadges = append(newBadges, badge)
				if badge.XPReward > 0 {
					meta := fmt.Sprintf(`{"badge_id": "%s", "badge_name": "%s"}`, badge.ID, badge.Name)
					metaPtr := &meta
					r.AddXP(ctx, userID, badge.XPReward, fmt.Sprintf("badge_%s", badge.Name), metaPtr)
				}
			}
		}
	}

	return newBadges, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (username, avatarURL string, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT username, COALESCE(avatar_url, '') FROM users WHERE id = $1`, userID).Scan(&username, &avatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", ErrUserNotFound
	}
	return
}
