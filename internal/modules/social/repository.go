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
	ListClubsByGenre(ctx context.Context, genre string, limit, offset int) ([]MusicClub, error)
	UpdateClubPlaylistID(ctx context.Context, clubID, playlistID uuid.UUID) error
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

	// Room queue voting
	SuggestTrack(ctx context.Context, c *QueueCandidate) error
	GetCandidates(ctx context.Context, roomID uuid.UUID) ([]QueueCandidate, error)
	GetCandidateByID(ctx context.Context, candidateID uuid.UUID) (*QueueCandidate, error)
	CastVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error
	RemoveVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error
	HasVoted(ctx context.Context, candidateID, userID uuid.UUID) (bool, error)
	GetNowPlaying(ctx context.Context, roomID uuid.UUID) (*RoomNowPlaying, error)
	SetNowPlaying(ctx context.Context, np *RoomNowPlaying) error
	RemoveCandidate(ctx context.Context, candidateID uuid.UUID) error
	GetCandidatesWithVoteState(ctx context.Context, roomID, userID uuid.UUID) ([]CandidateWithVoteState, error)
	LockRoomQueue(ctx context.Context, roomID uuid.UUID) (func(), error)
	PickRandomTrack(ctx context.Context) (*uuid.UUID, error)
	GetRoomMembers(ctx context.Context, roomID uuid.UUID) (int, error)

	// Stage & Raise-Hand
	GetHandRaise(ctx context.Context, roomID, userID uuid.UUID) (*HandRaise, error)
	RaiseHand(ctx context.Context, roomID, userID uuid.UUID) error
	LowerHand(ctx context.Context, roomID, userID uuid.UUID) error
	UpdateHandRaiseStatus(ctx context.Context, id uuid.UUID, status string) error
	ListPendingHandRaises(ctx context.Context, roomID uuid.UUID) ([]HandRaise, error)
	GetStageMember(ctx context.Context, roomID, userID uuid.UUID) (*StageMember, error)
	SetStageMember(ctx context.Context, m *StageMember) error
	RemoveStageMember(ctx context.Context, roomID, userID uuid.UUID) error
	ListStageSpeakers(ctx context.Context, roomID uuid.UUID) ([]StageMember, error)
	GetStageHost(ctx context.Context, roomID uuid.UUID) (*StageMember, error)
	UpdateStageMemberMuted(ctx context.Context, roomID, userID uuid.UUID, muted bool) error

	// Club Discussions (Phase 6)
	CreateClubDiscussion(ctx context.Context, d *ClubDiscussion) error
	ListClubDiscussions(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]ClubDiscussion, error)
	GetClubDiscussion(ctx context.Context, id uuid.UUID) (*ClubDiscussion, error)
	DeleteClubDiscussion(ctx context.Context, id uuid.UUID) error
	IncrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error
	DecrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error
	CreateClubDiscussionReply(ctx context.Context, r *ClubDiscussionReply) error
	GetClubDiscussionReplies(ctx context.Context, discussionID uuid.UUID) ([]ClubDiscussionReply, error)
	GetClubDiscussionReply(ctx context.Context, id uuid.UUID) (*ClubDiscussionReply, error)
	DeleteClubDiscussionReply(ctx context.Context, id uuid.UUID) error
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
		INSERT INTO user_follows (follower_id, followee_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, followerID, followedID)
	return err
}

func (r *repository) Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_follows
		WHERE follower_id = $1 AND followee_id = $2
	`, followerID, followedID)
	return err
}

func (r *repository) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1 FROM user_follows
			WHERE follower_id = $1 AND followee_id = $2
		)
	`, followerID, followedID)
	return exists, err
}

func (r *repository) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM user_follows WHERE followee_id = $1`, userID); err != nil {
		return nil, 0, err
	}

	var items []UserFollow
	if err := r.db.SelectContext(ctx, &items, `
		SELECT follower_id, followee_id, created_at
		FROM user_follows
		WHERE followee_id = $1
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
		SELECT follower_id, followee_id, created_at
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
		`SELECT COUNT(*) FROM user_follows WHERE followee_id = $1`, userID)
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
			SELECT followee_id FROM user_follows WHERE follower_id = $1
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
		INSERT INTO music_clubs (id, name, slug, description, cover_url, genre, playlist_id, created_by, is_public, max_members, member_count, created_at, updated_at)
		VALUES (:id, :name, :slug, :description, :cover_url, :genre, :playlist_id, :created_by, :is_public, :max_members, :member_count, :created_at, :updated_at)
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

func (r *repository) ListClubsByGenre(ctx context.Context, genre string, limit, offset int) ([]MusicClub, error) {
	var items []MusicClub
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM music_clubs
		WHERE is_public = true AND genre = $1
		ORDER BY member_count DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`, genre, limit, offset)
	return items, err
}

func (r *repository) UpdateClubPlaylistID(ctx context.Context, clubID, playlistID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE music_clubs SET playlist_id = $2, updated_at = NOW() WHERE id = $1
	`, clubID, playlistID)
	return err
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

// --- Room Queue Voting ---

func (r *repository) SuggestTrack(ctx context.Context, c *QueueCandidate) error {
	c.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO room_queue_candidates (id, room_id, track_id, suggested_by, vote_count, created_at)
		VALUES (:id, :room_id, :track_id, :suggested_by, :vote_count, :created_at)
		ON CONFLICT (room_id, track_id) DO NOTHING
	`, c)
	return err
}

func (r *repository) GetCandidates(ctx context.Context, roomID uuid.UUID) ([]QueueCandidate, error) {
	var items []QueueCandidate
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, room_id, track_id, suggested_by, vote_count, created_at
		FROM room_queue_candidates
		WHERE room_id = $1
		ORDER BY vote_count DESC, created_at ASC
	`, roomID)
	return items, err
}

func (r *repository) GetCandidateByID(ctx context.Context, candidateID uuid.UUID) (*QueueCandidate, error) {
	var c QueueCandidate
	err := r.db.GetContext(ctx, &c, `
		SELECT id, room_id, track_id, suggested_by, vote_count, created_at
		FROM room_queue_candidates
		WHERE id = $1
	`, candidateID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) CastVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS(SELECT 1 FROM room_queue_votes WHERE candidate_id = $1 AND user_id = $2)
	`, candidateID, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO room_queue_votes (candidate_id, user_id)
		VALUES ($1, $2)
	`, candidateID, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE room_queue_candidates
		SET vote_count = vote_count + 1
		WHERE id = $1
	`, candidateID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) RemoveVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM room_queue_votes
		WHERE candidate_id = $1 AND user_id = $2
	`, candidateID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE room_queue_candidates
		SET vote_count = GREATEST(vote_count - 1, 0)
		WHERE id = $1
	`, candidateID)

	return tx.Commit()
}

func (r *repository) HasVoted(ctx context.Context, candidateID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS(SELECT 1 FROM room_queue_votes WHERE candidate_id = $1 AND user_id = $2)
	`, candidateID, userID)
	return exists, err
}

func (r *repository) GetNowPlaying(ctx context.Context, roomID uuid.UUID) (*RoomNowPlaying, error) {
	var np RoomNowPlaying
	err := r.db.GetContext(ctx, &np, `
		SELECT room_id, track_id, started_at, suggested_by, source
		FROM room_now_playing
		WHERE room_id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &np, nil
}

func (r *repository) SetNowPlaying(ctx context.Context, np *RoomNowPlaying) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO room_now_playing (room_id, track_id, started_at, suggested_by, source)
		VALUES (:room_id, :track_id, :started_at, :suggested_by, :source)
		ON CONFLICT (room_id) DO UPDATE SET
			track_id = EXCLUDED.track_id,
			started_at = EXCLUDED.started_at,
			suggested_by = EXCLUDED.suggested_by,
			source = EXCLUDED.source
	`, np)
	return err
}

func (r *repository) RemoveCandidate(ctx context.Context, candidateID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM room_queue_candidates WHERE id = $1
	`, candidateID)
	return err
}

func (r *repository) GetCandidatesWithVoteState(ctx context.Context, roomID, userID uuid.UUID) ([]CandidateWithVoteState, error) {
	var items []CandidateWithVoteState
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			c.id, c.room_id, c.track_id, c.suggested_by, c.vote_count, c.created_at,
			CASE WHEN v.id IS NOT NULL THEN true ELSE false END AS has_voted
		FROM room_queue_candidates c
		LEFT JOIN room_queue_votes v ON v.candidate_id = c.id AND v.user_id = $2
		WHERE c.room_id = $1
		ORDER BY c.vote_count DESC, c.created_at ASC
	`, roomID, userID)
	return items, err
}

func (r *repository) LockRoomQueue(ctx context.Context, roomID uuid.UUID) (func(), error) {
	_, err := r.db.ExecContext(ctx, `
		SELECT pg_advisory_xact_lock(hashtext($1))
	`, "room_queue:"+roomID.String())
	if err != nil {
		return nil, err
	}
	return func() {}, nil
}

func (r *repository) PickRandomTrack(ctx context.Context) (*uuid.UUID, error) {
	var trackID uuid.UUID
	err := r.db.GetContext(ctx, &trackID, `
		SELECT id FROM tracks
		WHERE is_active = true
		ORDER BY RANDOM()
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}
	return &trackID, nil
}

func (r *repository) GetRoomMembers(ctx context.Context, roomID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM live_room_participants
		WHERE room_id = $1 AND is_active = true
	`, roomID)
	return count, err
}

// --- Stage & Raise-Hand ---

func (r *repository) GetHandRaise(ctx context.Context, roomID, userID uuid.UUID) (*HandRaise, error) {
	var hr HandRaise
	err := r.db.GetContext(ctx, &hr, `
		SELECT id, room_id, user_id, status, created_at
		FROM room_hand_raises
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err != nil {
		return nil, err
	}
	return &hr, nil
}

func (r *repository) RaiseHand(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO room_hand_raises (room_id, user_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (room_id, user_id) DO NOTHING
	`, roomID, userID)
	return err
}

func (r *repository) LowerHand(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM room_hand_raises
		WHERE room_id = $1 AND user_id = $2 AND status = 'pending'
	`, roomID, userID)
	return err
}

func (r *repository) UpdateHandRaiseStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE room_hand_raises SET status = $2 WHERE id = $1
	`, id, status)
	return err
}

func (r *repository) ListPendingHandRaises(ctx context.Context, roomID uuid.UUID) ([]HandRaise, error) {
	var items []HandRaise
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, room_id, user_id, status, created_at
		FROM room_hand_raises
		WHERE room_id = $1 AND status = 'pending'
		ORDER BY created_at ASC
	`, roomID)
	return items, err
}

func (r *repository) GetStageMember(ctx context.Context, roomID, userID uuid.UUID) (*StageMember, error) {
	var m StageMember
	err := r.db.GetContext(ctx, &m, `
		SELECT room_id, user_id, role, joined_at, muted
		FROM room_stage_members
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) SetStageMember(ctx context.Context, m *StageMember) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO room_stage_members (room_id, user_id, role, joined_at, muted)
		VALUES (:room_id, :user_id, :role, :joined_at, :muted)
		ON CONFLICT (room_id, user_id) DO UPDATE SET
			role = EXCLUDED.role,
			muted = EXCLUDED.muted,
			joined_at = EXCLUDED.joined_at
	`, m)
	return err
}

func (r *repository) RemoveStageMember(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM room_stage_members
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	return err
}

func (r *repository) ListStageSpeakers(ctx context.Context, roomID uuid.UUID) ([]StageMember, error) {
	var items []StageMember
	err := r.db.SelectContext(ctx, &items, `
		SELECT room_id, user_id, role, joined_at, muted
		FROM room_stage_members
		WHERE room_id = $1 AND role = 'speaker'
		ORDER BY joined_at ASC
	`, roomID)
	return items, err
}

func (r *repository) GetStageHost(ctx context.Context, roomID uuid.UUID) (*StageMember, error) {
	var m StageMember
	err := r.db.GetContext(ctx, &m, `
		SELECT room_id, user_id, role, joined_at, muted
		FROM room_stage_members
		WHERE room_id = $1 AND role = 'host'
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) UpdateStageMemberMuted(ctx context.Context, roomID, userID uuid.UUID, muted bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE room_stage_members SET muted = $3
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID, muted)
	return err
}

// --- Club Discussions (Phase 6) ---

func (r *repository) CreateClubDiscussion(ctx context.Context, d *ClubDiscussion) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO club_discussions (id, club_id, author_id, title, body, reply_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, d.ID, d.ClubID, d.AuthorID, d.Title, d.Body, d.ReplyCount, d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *repository) ListClubDiscussions(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]ClubDiscussion, error) {
	var items []ClubDiscussion
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, club_id, author_id, title, body, reply_count, created_at, updated_at
		FROM club_discussions
		WHERE club_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, clubID, limit, offset)
	return items, err
}

func (r *repository) GetClubDiscussion(ctx context.Context, id uuid.UUID) (*ClubDiscussion, error) {
	var d ClubDiscussion
	err := r.db.GetContext(ctx, &d, `
		SELECT id, club_id, author_id, title, body, reply_count, created_at, updated_at
		FROM club_discussions
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) DeleteClubDiscussion(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM club_discussions WHERE id = $1`, id)
	return err
}

func (r *repository) IncrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE club_discussions SET reply_count = reply_count + 1 WHERE id = $1
	`, id)
	return err
}

func (r *repository) DecrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE club_discussions SET reply_count = GREATEST(reply_count - 1, 0) WHERE id = $1
	`, id)
	return err
}

func (r *repository) CreateClubDiscussionReply(ctx context.Context, reply *ClubDiscussionReply) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO club_discussion_replies (id, discussion_id, author_id, body, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, reply.ID, reply.DiscussionID, reply.AuthorID, reply.Body, reply.CreatedAt)
	return err
}

func (r *repository) GetClubDiscussionReplies(ctx context.Context, discussionID uuid.UUID) ([]ClubDiscussionReply, error) {
	var items []ClubDiscussionReply
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, discussion_id, author_id, body, created_at
		FROM club_discussion_replies
		WHERE discussion_id = $1
		ORDER BY created_at ASC
	`, discussionID)
	return items, err
}

func (r *repository) GetClubDiscussionReply(ctx context.Context, id uuid.UUID) (*ClubDiscussionReply, error) {
	var reply ClubDiscussionReply
	err := r.db.GetContext(ctx, &reply, `
		SELECT id, discussion_id, author_id, body, created_at
		FROM club_discussion_replies
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (r *repository) DeleteClubDiscussionReply(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM club_discussion_replies WHERE id = $1`, id)
	return err
}
