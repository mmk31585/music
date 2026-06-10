package social

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Follow(ctx context.Context, followerID, followedID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error)
	GetFollowerCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error)
	InsertActivity(ctx context.Context, activity Activity) error
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int, types []string) ([]ActivityFeedItem, error)
	CreateParty(ctx context.Context, p *ListeningParty) error
	GetParty(ctx context.Context, id uuid.UUID) (*ListeningParty, error)
	ListActiveParties(ctx context.Context, limit, offset int) ([]ListeningParty, error)
	UpdatePartyStatus(ctx context.Context, id uuid.UUID, status string, trackID *uuid.UUID) error
	UpdatePartyPosition(ctx context.Context, id uuid.UUID, trackID *uuid.UUID, positionMs int64) error
	JoinParty(ctx context.Context, partyID, userID uuid.UUID) error
	LeaveParty(ctx context.Context, partyID, userID uuid.UUID) error
	GetPartyParticipants(ctx context.Context, partyID uuid.UUID) ([]ListeningPartyParticipant, error)
	CreateRoom(ctx context.Context, room *LiveRoom) error
	GetRoom(ctx context.Context, id uuid.UUID) (*LiveRoom, error)
	ListActiveRooms(ctx context.Context, limit, offset int) ([]LiveRoom, error)
	UpdateRoomStatus(ctx context.Context, id uuid.UUID, status string) error
	JoinRoom(ctx context.Context, roomID, userID uuid.UUID, role string) error
	LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error
	GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]LiveRoomParticipant, error)
	AddToRoomQueue(ctx context.Context, item *LiveRoomQueueItem) error
	GetRoomQueue(ctx context.Context, roomID uuid.UUID) ([]LiveRoomQueueItem, error)
	CreateClub(ctx context.Context, club *MusicClub) error
	GetClub(ctx context.Context, id uuid.UUID) (*MusicClub, error)
	ListClubs(ctx context.Context, limit, offset int) ([]MusicClub, error)
	JoinClub(ctx context.Context, clubID, userID uuid.UUID) error
	LeaveClub(ctx context.Context, clubID, userID uuid.UUID) error
	IsClubMember(ctx context.Context, clubID, userID uuid.UUID) (bool, error)
	IsClubAdmin(ctx context.Context, clubID, userID uuid.UUID) (bool, error)
	GetClubMembers(ctx context.Context, clubID uuid.UUID) ([]MusicClubMember, error)
	CreateClubPost(ctx context.Context, post *MusicClubPost) error
	GetClubPosts(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]MusicClubPost, error)
	CreateDiscussion(ctx context.Context, d *Discussion) error
	GetDiscussions(ctx context.Context, targetType string, targetID uuid.UUID, limit, offset int) ([]Discussion, error)
	GetDiscussionReplies(ctx context.Context, parentID uuid.UUID) ([]Discussion, error)
	CreateRating(ctx context.Context, rating *TrackRating) error
	GetTrackRatings(ctx context.Context, trackID uuid.UUID) ([]TrackRating, error)
	GetTrackRatingAverage(ctx context.Context, trackID uuid.UUID) (float64, int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// --- Follow ---

func (r *repository) Follow(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_follows (follower_id, followed_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, followerID, followedID)
	return err
}

func (r *repository) Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_follows
		WHERE follower_id = $1 AND followed_id = $2
	`, followerID, followedID)
	return err
}

func (r *repository) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1 FROM user_follows
			WHERE follower_id = $1 AND followed_id = $2
		)
	`, followerID, followedID)
	return exists, err
}

func (r *repository) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM user_follows WHERE followed_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	var items []UserFollow
	if err := r.db.SelectContext(ctx, &items, `
		SELECT follower_id, followed_id, created_at
		FROM user_follows
		WHERE followed_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	var items []UserFollow
	if err := r.db.SelectContext(ctx, &items, `
		SELECT follower_id, followed_id, created_at
		FROM user_follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) GetFollowerCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM user_follows WHERE followed_id = $1`, userID)
	return count, err
}

func (r *repository) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`, userID)
	return count, err
}

func (r *repository) InsertActivity(ctx context.Context, activity Activity) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO activities (id, user_id, type, target_id, target_type, metadata)
		VALUES (:id, :user_id, :type, :target_id, :target_type, :metadata)
	`, activity)
	return err
}

func (r *repository) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int, types []string) ([]ActivityFeedItem, error) {
	query := `
		SELECT
			a.id, a.user_id, a.type, a.target_id, a.target_type, a.metadata, a.created_at,
			u.display_name AS user_display_name,
			u.avatar_url AS user_avatar_url
		FROM activities a
		JOIN users u ON u.id = a.user_id
		WHERE a.user_id IN (
			SELECT followed_id FROM user_follows WHERE follower_id = $1
			UNION ALL
			SELECT $1
		)
	`

	args := []interface{}{userID, limit, offset}

	if len(types) > 0 {
		query += ` AND a.type = ANY($4) `
		args = append(args, types)
	}

	query += ` ORDER BY a.created_at DESC LIMIT $2 OFFSET $3 `

	var items []ActivityFeedItem
	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}

// --- Listening Parties ---

func (r *repository) CreateParty(ctx context.Context, p *ListeningParty) error {
	p.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO listening_parties (id, host_id, title, description, cover_url, is_public, status, current_track_id, current_position_ms, started_at)
		VALUES (:id, :host_id, :title, :description, :cover_url, :is_public, :status, :current_track_id, :current_position_ms, :started_at)
	`, p)
	return err
}

func (r *repository) GetParty(ctx context.Context, id uuid.UUID) (*ListeningParty, error) {
	var p ListeningParty
	err := r.db.GetContext(ctx, &p, `
		SELECT lp.*, COUNT(lpp.id) AS participant_count
		FROM listening_parties lp
		LEFT JOIN listening_party_participants lpp ON lpp.party_id = lp.id AND lpp.is_active = true
		WHERE lp.id = $1
		GROUP BY lp.id
	`, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) ListActiveParties(ctx context.Context, limit, offset int) ([]ListeningParty, error) {
	var items []ListeningParty
	err := r.db.SelectContext(ctx, &items, `
		SELECT lp.*, COUNT(lpp.id) AS participant_count
		FROM listening_parties lp
		LEFT JOIN listening_party_participants lpp ON lpp.party_id = lp.id AND lpp.is_active = true
		WHERE lp.status IN ('active', 'paused') AND lp.is_public = true
		GROUP BY lp.id
		ORDER BY lp.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *repository) UpdatePartyStatus(ctx context.Context, id uuid.UUID, status string, trackID *uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_parties
		SET status = $2,
		    ended_at = CASE WHEN $2 = 'ended' THEN NOW() ELSE ended_at END,
		    current_track_id = CASE WHEN $3 IS NOT NULL THEN $3 ELSE current_track_id END
		WHERE id = $1
	`, id, status, trackID)
	return err
}

func (r *repository) UpdatePartyPosition(ctx context.Context, id uuid.UUID, trackID *uuid.UUID, positionMs int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_parties SET current_track_id = $2, current_position_ms = $3
		WHERE id = $1
	`, id, trackID, positionMs)
	return err
}

func (r *repository) JoinParty(ctx context.Context, partyID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO listening_party_participants (party_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (party_id, user_id) DO UPDATE SET is_active = true, left_at = NULL
	`, partyID, userID)
	return err
}

func (r *repository) LeaveParty(ctx context.Context, partyID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_party_participants
		SET is_active = false, left_at = NOW()
		WHERE party_id = $1 AND user_id = $2
	`, partyID, userID)
	return err
}

func (r *repository) GetPartyParticipants(ctx context.Context, partyID uuid.UUID) ([]ListeningPartyParticipant, error) {
	var items []ListeningPartyParticipant
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM listening_party_participants
		WHERE party_id = $1 AND is_active = true
		ORDER BY joined_at ASC
	`, partyID)
	return items, err
}

// --- Live Rooms ---

func (r *repository) CreateRoom(ctx context.Context, room *LiveRoom) error {
	room.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO live_rooms (id, host_id, title, description, cover_url, is_public, status, created_at)
		VALUES (:id, :host_id, :title, :description, :cover_url, :is_public, :status, :created_at)
	`, room)
	return err
}

func (r *repository) GetRoom(ctx context.Context, id uuid.UUID) (*LiveRoom, error) {
	var room LiveRoom
	err := r.db.GetContext(ctx, &room, `SELECT * FROM live_rooms WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *repository) ListActiveRooms(ctx context.Context, limit, offset int) ([]LiveRoom, error) {
	var items []LiveRoom
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM live_rooms
		WHERE status = 'live' AND is_public = true
		ORDER BY listener_count DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *repository) UpdateRoomStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE live_rooms SET status = $2, ended_at = CASE WHEN $2 = 'ended' THEN NOW() ELSE ended_at END
		WHERE id = $1
	`, id, status)
	return err
}

func (r *repository) JoinRoom(ctx context.Context, roomID, userID uuid.UUID, role string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO live_room_participants (room_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, user_id) DO UPDATE SET is_active = true, left_at = NULL, role = $3
	`, roomID, userID, role)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE live_rooms SET listener_count = (SELECT COUNT(*) FROM live_room_participants WHERE room_id = $1 AND is_active = true) WHERE id = $1`, roomID)
	return err
}

func (r *repository) LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE live_room_participants
		SET is_active = false, left_at = NOW()
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE live_rooms SET listener_count = (SELECT COUNT(*) FROM live_room_participants WHERE room_id = $1 AND is_active = true) WHERE id = $1`, roomID)
	return err
}

func (r *repository) GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]LiveRoomParticipant, error) {
	var items []LiveRoomParticipant
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM live_room_participants
		WHERE room_id = $1 AND is_active = true
		ORDER BY role ASC, joined_at ASC
	`, roomID)
	return items, err
}

func (r *repository) AddToRoomQueue(ctx context.Context, item *LiveRoomQueueItem) error {
	item.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO live_room_queue (id, room_id, track_id, added_by, position, status, created_at)
		VALUES (:id, :room_id, :track_id, :added_by, :position, :status, :created_at)
	`, item)
	return err
}

func (r *repository) GetRoomQueue(ctx context.Context, roomID uuid.UUID) ([]LiveRoomQueueItem, error) {
	var items []LiveRoomQueueItem
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM live_room_queue
		WHERE room_id = $1 AND status = 'queued'
		ORDER BY position ASC
	`, roomID)
	return items, err
}

// --- Music Clubs ---

func (r *repository) CreateClub(ctx context.Context, club *MusicClub) error {
	club.ID = uuid.New()
	club.MemberCount = 1
	club.CreatedAt = time.Now().UTC()
	club.UpdatedAt = club.CreatedAt
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.NamedExecContext(ctx, `
		INSERT INTO music_clubs (id, name, description, cover_url, created_by, is_public, max_members, member_count, created_at, updated_at)
		VALUES (:id, :name, :description, :cover_url, :created_by, :is_public, :max_members, :member_count, :created_at, :updated_at)
	`, club); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO music_club_members (club_id, user_id, role)
		VALUES ($1, $2, 'admin')
	`, club.ID, club.CreatedBy); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetClub(ctx context.Context, id uuid.UUID) (*MusicClub, error) {
	var club MusicClub
	err := r.db.GetContext(ctx, &club, `SELECT * FROM music_clubs WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *repository) ListClubs(ctx context.Context, limit, offset int) ([]MusicClub, error) {
	var items []MusicClub
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_clubs
		WHERE is_public = true
		ORDER BY member_count DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *repository) JoinClub(ctx context.Context, clubID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO music_club_members (club_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, clubID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE music_clubs SET member_count = (SELECT COUNT(*) FROM music_club_members WHERE club_id = $1), updated_at = NOW()
		WHERE id = $1
	`, clubID)
	return err
}

func (r *repository) LeaveClub(ctx context.Context, clubID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM music_club_members WHERE club_id = $1 AND user_id = $2
	`, clubID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE music_clubs SET member_count = (SELECT COUNT(*) FROM music_club_members WHERE club_id = $1), updated_at = NOW()
		WHERE id = $1
	`, clubID)
	return err
}

func (r *repository) IsClubMember(ctx context.Context, clubID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (SELECT 1 FROM music_club_members WHERE club_id = $1 AND user_id = $2)
	`, clubID, userID)
	return exists, err
}

func (r *repository) IsClubAdmin(ctx context.Context, clubID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (SELECT 1 FROM music_club_members WHERE club_id = $1 AND user_id = $2 AND role IN ('admin', 'moderator'))
	`, clubID, userID)
	return exists, err
}

func (r *repository) GetClubMembers(ctx context.Context, clubID uuid.UUID) ([]MusicClubMember, error) {
	var items []MusicClubMember
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_club_members WHERE club_id = $1 ORDER BY role ASC, joined_at ASC
	`, clubID)
	return items, err
}

func (r *repository) CreateClubPost(ctx context.Context, post *MusicClubPost) error {
	post.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO music_club_posts (id, club_id, user_id, content, created_at, updated_at)
		VALUES (:id, :club_id, :user_id, :content, :created_at, :updated_at)
	`, post)
	return err
}

func (r *repository) GetClubPosts(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]MusicClubPost, error) {
	var items []MusicClubPost
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_club_posts WHERE club_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, clubID, limit, offset)
	return items, err
}

// --- Discussions ---

func (r *repository) CreateDiscussion(ctx context.Context, d *Discussion) error {
	d.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO discussions (id, user_id, target_type, target_id, content, parent_id, created_at, updated_at)
		VALUES (:id, :user_id, :target_type, :target_id, :content, :parent_id, :created_at, :updated_at)
	`, d)
	return err
}

func (r *repository) GetDiscussions(ctx context.Context, targetType string, targetID uuid.UUID, limit, offset int) ([]Discussion, error) {
	var items []Discussion
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM discussions
		WHERE target_type = $1 AND target_id = $2 AND parent_id IS NULL
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, targetType, targetID, limit, offset)
	return items, err
}

func (r *repository) GetDiscussionReplies(ctx context.Context, parentID uuid.UUID) ([]Discussion, error) {
	var items []Discussion
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM discussions WHERE parent_id = $1 ORDER BY created_at ASC
	`, parentID)
	return items, err
}

// --- Track Ratings ---

func (r *repository) CreateRating(ctx context.Context, rating *TrackRating) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO track_ratings (id, user_id, track_id, rating, review, created_at, updated_at)
		VALUES (:id, :user_id, :track_id, :rating, :review, :created_at, :updated_at)
		ON CONFLICT (user_id, track_id) DO UPDATE SET rating = EXCLUDED.rating, review = EXCLUDED.review, updated_at = NOW()
	`, rating)
	return err
}

func (r *repository) GetTrackRatings(ctx context.Context, trackID uuid.UUID) ([]TrackRating, error) {
	var items []TrackRating
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM track_ratings WHERE track_id = $1 ORDER BY created_at DESC`, trackID)
	return items, err
}

func (r *repository) GetTrackRatingAverage(ctx context.Context, trackID uuid.UUID) (float64, int, error) {
	var avg float64
	var count int
	err := r.db.GetContext(ctx, &avg, `SELECT COALESCE(AVG(rating), 0) FROM track_ratings WHERE track_id = $1`, trackID)
	if err != nil {
		return 0, 0, err
	}
	err = r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM track_ratings WHERE track_id = $1`, trackID)
	return avg, count, err
}
