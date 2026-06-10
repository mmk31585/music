package gamification

import "time"

type XPTransaction struct {
	ID           int64     `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Amount       int       `json:"amount" db:"amount"`
	BalanceAfter int       `json:"balance_after" db:"balance_after"`
	Source       string    `json:"source" db:"source"`
	Metadata     *string   `json:"metadata" db:"metadata"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type LevelDefinition struct {
	Level        int    `json:"level" db:"level"`
	XPRequired   int64  `json:"xp_required" db:"xp_required"`
	Title        string `json:"title" db:"title"`
	TitlePersian string `json:"title_persian" db:"title_persian"`
	Privileges   string `json:"privileges" db:"privileges"`
}

type Badge struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IconURL     string `json:"icon_url" db:"icon_url"`
	Category    string `json:"category" db:"category"`
	Rarity      string `json:"rarity" db:"rarity"`
	Criteria    string `json:"criteria" db:"criteria"`
	XPReward    int    `json:"xp_reward" db:"xp_reward"`
	IsHidden    bool   `json:"is_hidden" db:"is_hidden"`
}

type UserBadge struct {
	UserID      string    `json:"user_id" db:"user_id"`
	BadgeID     string    `json:"badge_id" db:"badge_id"`
	EarnedAt    time.Time `json:"earned_at" db:"earned_at"`
	IsDisplayed bool      `json:"is_displayed" db:"is_displayed"`
	Badge       *Badge    `json:"badge,omitempty"`
}

type DailyChallenge struct {
	ID            string    `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	ChallengeType string    `json:"challenge_type" db:"challenge_type"`
	TargetCount   int       `json:"target_count" db:"target_count"`
	XPReward      int       `json:"xp_reward" db:"xp_reward"`
	BadgeRewardID *string   `json:"badge_reward_id" db:"badge_reward_id"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	ValidFrom     time.Time `json:"valid_from" db:"valid_from"`
	ValidUntil    time.Time `json:"valid_until" db:"valid_until"`
}

type UserChallenge struct {
	UserID      string          `json:"user_id" db:"user_id"`
	ChallengeID string          `json:"challenge_id" db:"challenge_id"`
	Progress    int             `json:"progress" db:"progress"`
	IsCompleted bool            `json:"is_completed" db:"is_completed"`
	CompletedAt *time.Time      `json:"completed_at" db:"completed_at"`
	Challenge   *DailyChallenge `json:"challenge,omitempty"`
}

type LeaderboardEntry struct {
	Rank      int    `json:"rank"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Score     int64  `json:"score"`
}

type XPRequest struct {
	Amount int    `json:"amount" binding:"required,min=1,max=1000"`
	Source string `json:"source" binding:"required,max=50"`
}

type UserProfileResponse struct {
	UserID       string          `json:"user_id"`
	Username     string          `json:"username"`
	AvatarURL    string          `json:"avatar_url"`
	Level        int             `json:"level"`
	CurrentXP    int64           `json:"current_xp"`
	NextLevelXP  int64           `json:"next_level_xp"`
	TotalXP      int64           `json:"total_xp"`
	Title        string          `json:"title"`
	TitlePersian string          `json:"title_persian"`
	Rank         int             `json:"rank"`
	Badges       []UserBadge     `json:"badges"`
	Challenges   []UserChallenge `json:"challenges"`
}

type LeaderboardResponse struct {
	Type    string             `json:"type"`
	Entries []LeaderboardEntry `json:"entries"`
}

const XPForStream int = 1
const XPForLike int = 5
const XPForShare int = 10
const XPForContribution int = 50
const XPForTranslation int = 30
const XPForLogin int = 15
const XPForDailyChallenge int = 50

func CalculateLevel(xp int64) (level int, title, titlePersian string, nextXP int64) {
	levels := []struct {
		Level        int
		XPRequired   int64
		Title        string
		TitlePersian string
	}{
		{1, 0, "Novice Listener", "شنونده تازه‌کار"},
		{2, 100, "Curious Ear", "گوش کنجکاو"},
		{3, 520, "Melody Seeker", "جوینده ملودی"},
		{4, 1600, "Rhythm Catcher", "ریتم‌گیر"},
		{5, 3900, "Harmony Lover", "عاشق هارمونی"},
		{6, 8000, "Tone Explorer", "کاوشگر نغمه"},
		{7, 14000, "Music Wanderer", "سرگردان موسیقی"},
		{8, 22000, "Note Collector", "جمع‌آورنده نت"},
		{9, 33000, "Groove Master", "استاد گروو"},
		{10, 48000, "Beat Connoisseur", "خبره ضرب‌آهنگ"},
		{11, 66000, "Vocal Virtuoso", "استاد آواز"},
		{12, 88000, "Lyric Sage", "فرزانه شعر"},
		{13, 115000, "Dastgah Navigator", "راهبر دستگاه"},
		{14, 147000, "Genre Weaver", "بافنده سبک"},
		{15, 185000, "Orchestra Commander", "فرمانده ارکستر"},
		{16, 230000, "Melody Architect", "معمار ملودی"},
		{17, 283000, "Music Scholar", "دانشور موسیقی"},
		{18, 345000, "Sound Alchemist", "کیمیاگر صدا"},
		{19, 417000, "Legendary Listener", "شنونده افسانه‌ای"},
		{20, 500000, "Persian Gem", "گوهر پارسی"},
	}

	for i := len(levels) - 1; i >= 0; i-- {
		if xp >= levels[i].XPRequired {
			l := levels[i]
			if i+1 < len(levels) {
				return l.Level, l.Title, l.TitlePersian, levels[i+1].XPRequired
			}
			return l.Level, l.Title, l.TitlePersian, l.XPRequired
		}
	}
	return 1, "Novice Listener", "شنونده تازه‌کار", levels[1].XPRequired
}
