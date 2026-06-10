package creator

import (
	"time"

	"github.com/google/uuid"
)

type CreatorStats struct {
	UserID           uuid.UUID `db:"user_id" json:"user_id"`
	TotalPlays       int64     `db:"total_plays" json:"total_plays"`
	UniqueListeners  int64     `db:"unique_listeners" json:"unique_listeners"`
	TotalFollowers   int       `db:"total_followers" json:"total_followers"`
	TotalTracks      int       `db:"total_tracks" json:"total_tracks"`
	TotalAlbums      int       `db:"total_albums" json:"total_albums"`
	TotalPlaylists   int       `db:"total_playlists" json:"total_playlists"`
	EstimatedRevenue int64     `db:"estimated_revenue" json:"estimated_revenue"`
	LastCalculated   time.Time `db:"last_calculated" json:"last_calculated"`
}

type CreatorDailyStat struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserID       uuid.UUID `db:"user_id" json:"user_id"`
	Date         time.Time `db:"date" json:"date"`
	Plays        int       `db:"plays" json:"plays"`
	Listeners    int       `db:"listeners" json:"listeners"`
	Likes        int       `db:"likes" json:"likes"`
	Follows      int       `db:"follows" json:"follows"`
	Shares       int       `db:"shares" json:"shares"`
	RevenueCents int       `db:"revenue_cents" json:"revenue_cents"`
}

type TrackStats struct {
	TrackID    uuid.UUID `db:"track_id" json:"track_id"`
	Title      string    `db:"title" json:"title"`
	TotalPlays int64     `db:"total_plays" json:"total_plays"`
	TotalLikes int64     `db:"total_likes" json:"total_likes"`
	Duration   int       `db:"duration" json:"duration"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Overview struct {
	Stats      CreatorStats       `json:"stats"`
	DailyStats []CreatorDailyStat `json:"daily_stats,omitempty"`
}

// Earnings
type EarningsBreakdown struct {
	TotalRevenue    int64      `json:"total_revenue"`
	StreamRevenue   int64      `json:"stream_revenue"`
	TipRevenue      int64      `json:"tip_revenue"`
	SubscriptionRev int64      `json:"subscription_revenue"`
	PendingPayout   int64      `json:"pending_payout"`
	LastPayout      int64      `json:"last_payout"`
	LastPayoutDate  *time.Time `json:"last_payout_date,omitempty"`
}

type Payout struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	Amount    int64      `db:"amount" json:"amount"`
	Method    string     `db:"method" json:"method"`
	Status    string     `db:"status" json:"status"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	PaidAt    *time.Time `db:"paid_at" json:"paid_at,omitempty"`
}

type PayoutMethod struct {
	ID       uuid.UUID `db:"id" json:"id"`
	UserID   uuid.UUID `db:"user_id" json:"user_id"`
	Type     string    `db:"type" json:"type"`
	Details  string    `db:"details" json:"details"`
	IsActive bool      `db:"is_active" json:"is_active"`
}

// Audience
type TopListener struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	PlayCount int64     `db:"play_count" json:"play_count"`
}

type GeographicStat struct {
	Country   string `db:"country" json:"country"`
	City      string `db:"city" json:"city"`
	Listeners int    `db:"listeners" json:"listeners"`
	Plays     int    `db:"plays" json:"plays"`
}

type AudienceOverview struct {
	TopListeners   []TopListener    `json:"top_listeners"`
	Geographics    []GeographicStat `json:"geographics"`
	TotalListeners int              `json:"total_listeners"`
	NewListeners   int              `json:"new_listeners_7d"`
	RepeatRate     float64          `json:"repeat_rate"`
}

// Content Management
type TrackUpdateRequest struct {
	Title        string   `json:"title"`
	PersianTitle string   `json:"persian_title"`
	GenreIDs     []string `json:"genre_ids"`
	Explicit     *bool    `json:"explicit"`
	TrackNumber  *int     `json:"track_number"`
	Lyrics       string   `json:"lyrics"`
}

type AlbumUpdateRequest struct {
	Title        string   `json:"title"`
	PersianTitle string   `json:"persian_title"`
	Description  string   `json:"description"`
	AlbumType    string   `json:"album_type"`
	ReleaseDate  string   `json:"release_date"`
	GenreIDs     []string `json:"genre_ids"`
}

type CreatorContent struct {
	Tracks    []TrackStats      `json:"tracks"`
	Albums    []AlbumStats      `json:"albums"`
	Playlists []CreatorPlaylist `json:"playlists"`
}

type AlbumStats struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	ReleaseYear int       `db:"release_year" json:"release_year"`
	TrackCount  int       `db:"track_count" json:"track_count"`
	TotalPlays  int64     `db:"total_plays" json:"total_plays"`
	CoverURL    *string   `db:"cover_url" json:"cover_url,omitempty"`
}

type CreatorPlaylist struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	TrackCount int       `db:"track_count" json:"track_count"`
	IsPublic   bool      `db:"is_public" json:"is_public"`
	CoverURL   *string   `db:"cover_url" json:"cover_url,omitempty"`
}
