package creator

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStats(ctx context.Context, userID uuid.UUID) (*CreatorStats, error) {
	var stats CreatorStats
	err := r.db.GetContext(ctx, &stats, `
		SELECT
			user_id, total_plays, unique_listeners, total_followers,
			total_tracks, total_albums, total_playlists, estimated_revenue,
			last_calculated
		FROM creator_stats
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *Repository) GetDailyStats(ctx context.Context, userID uuid.UUID, from, to string, limit int) ([]CreatorDailyStat, error) {
	args := []interface{}{userID}
	query := `
		SELECT id, user_id, date, plays, listeners, likes, follows, shares, revenue_cents
		FROM creator_daily_stats
		WHERE user_id = $1
	`

	if from != "" {
		args = append(args, from)
		query += fmt.Sprintf(` AND date >= $%d`, len(args))
	}
	if to != "" {
		args = append(args, to)
		query += fmt.Sprintf(` AND date <= $%d`, len(args))
	}

	args = append(args, limit)
	query += fmt.Sprintf(` ORDER BY date DESC LIMIT $%d`, len(args))

	var items []CreatorDailyStat
	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}

func (r *Repository) GetTrackStats(ctx context.Context, userID uuid.UUID) ([]TrackStats, error) {
	var items []TrackStats
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			t.id AS track_id,
			t.title,
			COALESCE(tp.total_plays, 0) AS total_plays,
			COALESCE(r.likes, 0) AS total_likes,
			t.duration_seconds AS duration,
			t.created_at
		FROM tracks t
		LEFT JOIN (
			SELECT target_id, COUNT(*) AS likes
			FROM reactions
			WHERE target_type = 'track' AND type = 'like'
			GROUP BY target_id
		) r ON r.target_id = t.id::text
		LEFT JOIN (
			SELECT track_id, COUNT(*) AS total_plays
			FROM track_plays
			GROUP BY track_id
		) tp ON tp.track_id = t.id
		WHERE t.artist_id = $1 OR t.id IN (
			SELECT track_id FROM track_artists WHERE artist_id = $1
		)
		ORDER BY total_plays DESC
	`, userID)
	return items, err
}

func (r *Repository) RefreshStats(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO creator_stats (user_id, total_plays, unique_listeners, total_followers,
			total_tracks, total_albums, total_playlists, estimated_revenue, last_calculated)
		SELECT
			$1 AS user_id,
			COALESCE((SELECT COUNT(*) FROM track_plays WHERE user_id = $1), 0) AS total_plays,
			COALESCE((SELECT COUNT(DISTINCT user_id) FROM track_plays WHERE track_id IN (
				SELECT id FROM tracks WHERE artist_id = $1
			)), 0) AS unique_listeners,
			COALESCE((SELECT COUNT(*) FROM user_follows WHERE followed_id = $1), 0) AS total_followers,
			COALESCE((SELECT COUNT(*) FROM tracks WHERE artist_id = $1), 0) AS total_tracks,
			COALESCE((SELECT COUNT(*) FROM albums WHERE artist_id = $1), 0) AS total_albums,
			COALESCE((SELECT COUNT(*) FROM playlists WHERE user_id = $1), 0) AS total_playlists,
			0 AS estimated_revenue,
			NOW() AS last_calculated
		ON CONFLICT (user_id)
		DO UPDATE SET
			total_plays = EXCLUDED.total_plays,
			unique_listeners = EXCLUDED.unique_listeners,
			total_followers = EXCLUDED.total_followers,
			total_tracks = EXCLUDED.total_tracks,
			total_albums = EXCLUDED.total_albums,
			total_playlists = EXCLUDED.total_playlists,
			last_calculated = NOW()
	`, userID)
	return err
}

func (r *Repository) IsCreator(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1 FROM tracks WHERE artist_id = $1
			UNION
			SELECT 1 FROM artists WHERE id IN (
				SELECT artist_id FROM track_artists WHERE artist_id = $1
			)
		)
	`, userID)
	return exists, err
}

// --- Earnings ---

func (r *Repository) GetEarningsBreakdown(ctx context.Context, userID uuid.UUID) (*EarningsBreakdown, error) {
	var e EarningsBreakdown
	err := r.db.GetContext(ctx, &e, `
		SELECT
			COALESCE((SELECT SUM(amount) FROM tips WHERE artist_id = $1 AND status = 'completed'), 0) AS tip_revenue,
			COALESCE((SELECT SUM(amount) FROM subscription_payments WHERE artist_id = $1 AND status = 'paid'), 0) AS subscription_revenue,
			COALESCE((SELECT SUM(estimated_revenue) FROM creator_stats WHERE user_id = $1), 0) AS stream_revenue,
			COALESCE((SELECT SUM(amount) FROM payouts WHERE user_id = $1 AND status = 'pending'), 0) AS pending_payout,
			COALESCE((SELECT amount FROM payouts WHERE user_id = $1 AND status = 'paid' ORDER BY paid_at DESC LIMIT 1), 0) AS last_payout,
			(SELECT paid_at FROM payouts WHERE user_id = $1 AND status = 'paid' ORDER BY paid_at DESC LIMIT 1) AS last_payout_date
	`, userID)
	if err != nil {
		return nil, err
	}
	e.TotalRevenue = e.StreamRevenue + e.TipRevenue + e.SubscriptionRev
	return &e, nil
}

func (r *Repository) GetPayoutHistory(ctx context.Context, userID uuid.UUID, limit int) ([]Payout, error) {
	var items []Payout
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, user_id, amount, method, status, created_at, paid_at
		FROM payouts WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2
	`, userID, limit)
	return items, err
}

func (r *Repository) GetPayoutMethods(ctx context.Context, userID uuid.UUID) ([]PayoutMethod, error) {
	var items []PayoutMethod
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, user_id, type, details, is_active
		FROM payout_methods WHERE user_id = $1 AND is_active = true
	`, userID)
	return items, err
}

// --- Audience ---

func (r *Repository) GetTopListeners(ctx context.Context, userID uuid.UUID, limit int) ([]TopListener, error) {
	var items []TopListener
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			tp.user_id, u.display_name AS username, u.avatar_url,
			COUNT(*) AS play_count
		FROM track_plays tp
		JOIN tracks t ON t.id = tp.track_id
		LEFT JOIN users u ON u.id = tp.user_id
		WHERE t.artist_id = $1
		GROUP BY tp.user_id, u.display_name, u.avatar_url
		ORDER BY play_count DESC
		LIMIT $2
	`, userID, limit)
	return items, err
}

func (r *Repository) GetGeographicStats(ctx context.Context, userID uuid.UUID) ([]GeographicStat, error) {
	var items []GeographicStat
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			COALESCE(tp.country, 'Unknown') AS country,
			COALESCE(tp.city, 'Unknown') AS city,
			COUNT(DISTINCT tp.user_id) AS listeners,
			COUNT(*) AS plays
		FROM track_plays tp
		JOIN tracks t ON t.id = tp.track_id
		WHERE t.artist_id = $1
		GROUP BY tp.country, tp.city
		ORDER BY plays DESC
		LIMIT 20
	`, userID)
	return items, err
}

func (r *Repository) GetAudienceOverview(ctx context.Context, userID uuid.UUID) (*AudienceOverview, error) {
	var total int
	_ = r.db.GetContext(ctx, &total, `
		SELECT COUNT(DISTINCT user_id) FROM track_plays tp
		JOIN tracks t ON t.id = tp.track_id WHERE t.artist_id = $1
	`, userID)

	var new7d int
	_ = r.db.GetContext(ctx, &new7d, `
		SELECT COUNT(DISTINCT tp.user_id) FROM track_plays tp
		JOIN tracks t ON t.id = tp.track_id
		WHERE t.artist_id = $1 AND tp.created_at >= NOW() - INTERVAL '7 days'
		AND tp.user_id NOT IN (
			SELECT DISTINCT tp2.user_id FROM track_plays tp2
			JOIN tracks t2 ON t2.id = tp2.track_id
			WHERE t2.artist_id = $1 AND tp2.created_at < NOW() - INTERVAL '7 days'
		)
	`, userID)

	var repeatRate float64
	_ = r.db.GetContext(ctx, &repeatRate, `
		SELECT COALESCE(
			(SELECT COUNT(*) * 1.0 / NULLIF(COUNT(DISTINCT user_id), 0)
			FROM (
				SELECT user_id, COUNT(*) AS plays
				FROM track_plays tp
				JOIN tracks t ON t.id = tp.track_id
				WHERE t.artist_id = $1
				GROUP BY user_id
			) sub WHERE plays > 1), 0)
	`, userID)

	return &AudienceOverview{
		TotalListeners: total,
		NewListeners:   new7d,
		RepeatRate:     repeatRate,
	}, nil
}

// --- Content Management ---

func (r *Repository) GetCreatorContent(ctx context.Context, userID uuid.UUID) (*CreatorContent, error) {
	var tracks []TrackStats
	err := r.db.SelectContext(ctx, &tracks, `
		SELECT t.id AS track_id, t.title, COALESCE(tp.total_plays, 0) AS total_plays,
			COALESCE(r.likes, 0) AS total_likes, t.duration_seconds AS duration, t.created_at
		FROM tracks t
		LEFT JOIN (SELECT target_id, COUNT(*) AS likes FROM reactions WHERE target_type = 'track' AND type = 'like' GROUP BY target_id) r ON r.target_id = t.id::text
		LEFT JOIN (SELECT track_id, COUNT(*) AS total_plays FROM track_plays GROUP BY track_id) tp ON tp.track_id = t.id
		WHERE t.artist_id = $1
		ORDER BY t.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}

	var albums []AlbumStats
	err = r.db.SelectContext(ctx, &albums, `
		SELECT a.id, a.title, EXTRACT(YEAR FROM a.release_date)::int AS release_year,
			(SELECT COUNT(*) FROM tracks WHERE album_id = a.id) AS track_count,
			COALESCE(SUM(tp.total_plays), 0) AS total_plays, a.cover_url
		FROM albums a
		LEFT JOIN tracks t ON t.album_id = a.id
		LEFT JOIN (SELECT track_id, COUNT(*) AS total_plays FROM track_plays GROUP BY track_id) tp ON tp.track_id = t.id
		WHERE a.artist_id = $1
		GROUP BY a.id, a.title, a.release_date, a.cover_url
		ORDER BY a.release_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}

	var playlists []CreatorPlaylist
	err = r.db.SelectContext(ctx, &playlists, `
		SELECT p.id, p.name,
			(SELECT COUNT(*) FROM playlist_tracks WHERE playlist_id = p.id) AS track_count,
			p.is_public, p.cover_url
		FROM playlists p
		WHERE p.user_id = $1
		ORDER BY p.name ASC
	`, userID)
	if err != nil {
		return nil, err
	}

	return &CreatorContent{
		Tracks:    tracks,
		Albums:    albums,
		Playlists: playlists,
	}, nil
}

func (r *Repository) UpdateTrack(ctx context.Context, trackID uuid.UUID, req TrackUpdateRequest) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE tracks SET
			title = COALESCE(NULLIF(:title, ''), title),
			persian_title = COALESCE(NULLIF(:persian_title, ''), persian_title),
			explicit = COALESCE(:explicit, explicit),
			track_number = :track_number,
			lyrics = :lyrics
		WHERE id = :id
	`, map[string]interface{}{
		"id":            trackID,
		"title":         req.Title,
		"persian_title": req.PersianTitle,
		"explicit":      req.Explicit,
		"track_number":  req.TrackNumber,
		"lyrics":        req.Lyrics,
	})
	return err
}

func (r *Repository) UpdateAlbum(ctx context.Context, albumID uuid.UUID, req AlbumUpdateRequest) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE albums SET
			title = COALESCE(NULLIF(:title, ''), title),
			persian_title = COALESCE(NULLIF(:persian_title, ''), persian_title),
			description = :description,
			album_type = COALESCE(NULLIF(:album_type, ''), album_type),
			release_date = COALESCE(NULLIF(:release_date, '')::date, release_date)
		WHERE id = :id
	`, map[string]interface{}{
		"id":            albumID,
		"title":         req.Title,
		"persian_title": req.PersianTitle,
		"description":   req.Description,
		"album_type":    req.AlbumType,
		"release_date":  req.ReleaseDate,
	})
	return err
}

func (r *Repository) DeleteTrack(ctx context.Context, trackID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tracks WHERE id = $1 AND artist_id = $2`, trackID, userID)
	return err
}

func (r *Repository) DeleteAlbum(ctx context.Context, albumID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM albums WHERE id = $1 AND artist_id = $2`, albumID, userID)
	return err
}
