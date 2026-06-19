package social

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"music/internal/modules/playlist"
)

type Service struct {
	repo             Repository
	partyBroadcaster *PartyBroadcaster
	roomBroadcaster  *RoomBroadcaster
	playlistSvc      *playlist.Service
	clubSvc          *ClubService
	discussionSvc    *DiscussionService
}

func NewService(repo Repository, partyBroadcaster *PartyBroadcaster, roomBroadcaster *RoomBroadcaster) *Service {
	return &Service{
		repo:             repo,
		partyBroadcaster: partyBroadcaster,
		roomBroadcaster:  roomBroadcaster,
		discussionSvc:    NewDiscussionService(repo),
	}
}

func NewServiceWithPlaylist(repo Repository, partyBroadcaster *PartyBroadcaster, roomBroadcaster *RoomBroadcaster, playlistSvc *playlist.Service) *Service {
	svc := &Service{
		repo:             repo,
		partyBroadcaster: partyBroadcaster,
		roomBroadcaster:  roomBroadcaster,
		playlistSvc:      playlistSvc,
		discussionSvc:    NewDiscussionService(repo),
	}
	svc.clubSvc = NewClubService(repo, playlistSvc, svc)
	return svc
}

func (s *Service) Follow(ctx context.Context, followerID, followedID string) error {
	fUID, err := uuid.Parse(followerID)
	if err != nil {
		return fmt.Errorf("invalid follower id: %w", err)
	}
	fUID2, err := uuid.Parse(followedID)
	if err != nil {
		return fmt.Errorf("invalid followed id: %w", err)
	}
	if fUID == fUID2 {
		return fmt.Errorf("cannot follow yourself")
	}
	return s.repo.Follow(ctx, fUID, fUID2)
}

func (s *Service) Unfollow(ctx context.Context, followerID, followedID string) error {
	fUID, _ := uuid.Parse(followerID)
	fUID2, _ := uuid.Parse(followedID)
	return s.repo.Unfollow(ctx, fUID, fUID2)
}

func (s *Service) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]UserFollow, int, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetFollowers(ctx, uid, limit, offset)
}

func (s *Service) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]UserFollow, int, error) {
	uid, _ := uuid.Parse(userID)
	return s.repo.GetFollowing(ctx, uid, limit, offset)
}

func (s *Service) IsFollowing(ctx context.Context, followerID, followedID string) (bool, error) {
	fUID, _ := uuid.Parse(followerID)
	fUID2, _ := uuid.Parse(followedID)
	return s.repo.IsFollowing(ctx, fUID, fUID2)
}

func (s *Service) GetFeed(ctx context.Context, userID string, limit, offset int, types string) ([]ActivityFeedItem, error) {
	uid, _ := uuid.Parse(userID)
	var typeFilter []string
	if types != "" {
		typeFilter = strings.Split(types, ",")
	}
	return s.repo.GetFeed(ctx, uid, limit, offset, typeFilter)
}

func (s *Service) RecordActivity(ctx context.Context, userID, activityType, targetID, targetType string, metadata map[string]interface{}) error {
	uid, _ := uuid.Parse(userID)
	meta, _ := json.Marshal(metadata)
	activity := Activity{
		ID:         uuid.New(),
		UserID:     uid,
		Type:       activityType,
		TargetID:   targetID,
		TargetType: targetType,
		Metadata:   meta,
		CreatedAt:  time.Now().UTC(),
	}
	return s.repo.InsertActivity(ctx, activity)
}

// --- Listening Parties ---

func (s *Service) CreateParty(ctx context.Context, req CreatePartyRequest, hostID string) (*ListeningParty, error) {
	hostUUID, _ := uuid.Parse(hostID)
	now := time.Now().UTC()
	party := &ListeningParty{
		HostID:    hostUUID,
		Title:     req.Title,
		IsPublic:  req.IsPublic,
		Status:    "active",
		StartedAt: now,
		CreatedAt: now,
	}
	if req.Description != "" {
		party.Description = &req.Description
	}
	if req.TrackID != nil && *req.TrackID != "" {
		tid, _ := uuid.Parse(*req.TrackID)
		party.CurrentTrackID = &tid
	}
	if err := s.repo.CreateParty(ctx, party); err != nil {
		return nil, err
	}
	// Auto-create a room with the same ID as the party so that
	// room queue operations (suggest, vote, queue state) work.
	room := &LiveRoom{
		ID:          party.ID,
		HostID:      hostUUID,
		Title:       req.Title,
		IsPublic:    req.IsPublic,
		Status:      "live",
		CreatedAt:   now,
		Description: party.Description,
	}
	if err := s.repo.CreateRoom(ctx, room); err != nil {
		// Room creation failure should not block party creation
		zap.L().Error("failed to create room for party", zap.String("party_id", party.ID.String()), zap.Error(err))
	}
	// Host auto-joins the room
	_ = s.repo.JoinRoom(ctx, party.ID, hostUUID, "host")
	// Host auto-joins the party
	_ = s.repo.JoinParty(ctx, party.ID, hostUUID)
	if s.partyBroadcaster != nil {
		s.partyBroadcaster.UserJoined(party.ID, hostUUID)
	}
	return party, nil
}

func (s *Service) GetParty(ctx context.Context, id string) (*ListeningParty, error) {
	uid, _ := uuid.Parse(id)
	return s.repo.GetParty(ctx, uid)
}

func (s *Service) ListActiveParties(ctx context.Context, limit, offset int) ([]ListeningParty, error) {
	return s.repo.ListActiveParties(ctx, limit, offset)
}

func (s *Service) UpdatePartyStatus(ctx context.Context, id, status string, trackID *string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid party id: %w", err)
	}
	var tid *uuid.UUID
	if trackID != nil && *trackID != "" {
		t, err := uuid.Parse(*trackID)
		if err == nil {
			tid = &t
		}
		// Non-UUID track IDs (numeric IDs, etc.) are silently ignored —
		// the status update still proceeds without changing current_track_id.
	}
	if err := s.repo.UpdatePartyStatus(ctx, uid, status, tid); err != nil {
		return err
	}
	if s.partyBroadcaster != nil && status == "ended" {
		s.partyBroadcaster.PartyEnded(uid)
	}
	return nil
}

func (s *Service) UpdatePartyPosition(ctx context.Context, id, trackID string, positionMs int64) error {
	uid, _ := uuid.Parse(id)
	var tid *uuid.UUID
	if trackID != "" {
		t, _ := uuid.Parse(trackID)
		tid = &t
	}
	return s.repo.UpdatePartyPosition(ctx, uid, tid, positionMs)
}

func (s *Service) JoinParty(ctx context.Context, partyID, userID string) error {
	puid, _ := uuid.Parse(partyID)
	uid, _ := uuid.Parse(userID)
	if err := s.repo.JoinParty(ctx, puid, uid); err != nil {
		return err
	}
	if s.partyBroadcaster != nil {
		s.partyBroadcaster.UserJoined(puid, uid)
	}
	return nil
}

func (s *Service) LeaveParty(ctx context.Context, partyID, userID string) error {
	puid, _ := uuid.Parse(partyID)
	uid, _ := uuid.Parse(userID)
	if err := s.repo.LeaveParty(ctx, puid, uid); err != nil {
		return err
	}
	if s.partyBroadcaster != nil {
		s.partyBroadcaster.UserLeft(puid, uid)
	}
	return nil
}

func (s *Service) GetPartyParticipants(ctx context.Context, partyID string) ([]ListeningPartyParticipant, error) {
	puid, _ := uuid.Parse(partyID)
	return s.repo.GetPartyParticipants(ctx, puid)
}

// --- Live Rooms ---

func (s *Service) CreateRoom(ctx context.Context, req CreateRoomRequest, hostID string) (*LiveRoom, error) {
	hostUUID, _ := uuid.Parse(hostID)
	room := &LiveRoom{
		HostID:    hostUUID,
		Title:     req.Title,
		IsPublic:  req.IsPublic,
		Status:    "live",
		CreatedAt: time.Now().UTC(),
	}
	if req.Description != "" {
		room.Description = &req.Description
	}
	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return nil, err
	}
	_ = s.repo.JoinRoom(ctx, room.ID, hostUUID, "host")
	if s.roomBroadcaster != nil {
		s.roomBroadcaster.ParticipantJoined(room.ID, hostUUID, "host")
	}
	return room, nil
}

func (s *Service) GetRoom(ctx context.Context, id string) (*LiveRoom, error) {
	uid, _ := uuid.Parse(id)
	return s.repo.GetRoom(ctx, uid)
}

func (s *Service) ListActiveRooms(ctx context.Context, limit, offset int) ([]LiveRoom, error) {
	return s.repo.ListActiveRooms(ctx, limit, offset)
}

func (s *Service) UpdateRoomStatus(ctx context.Context, id, status string) error {
	uid, _ := uuid.Parse(id)
	return s.repo.UpdateRoomStatus(ctx, uid, status)
}

func (s *Service) JoinRoom(ctx context.Context, roomID, userID string) error {
	rid, _ := uuid.Parse(roomID)
	uid, _ := uuid.Parse(userID)
	if err := s.repo.JoinRoom(ctx, rid, uid, "listener"); err != nil {
		return err
	}
	if s.roomBroadcaster != nil {
		s.roomBroadcaster.ParticipantJoined(rid, uid, "listener")
	}
	return nil
}

func (s *Service) LeaveRoom(ctx context.Context, roomID, userID string) error {
	rid, _ := uuid.Parse(roomID)
	uid, _ := uuid.Parse(userID)
	if err := s.repo.LeaveRoom(ctx, rid, uid); err != nil {
		return err
	}
	if s.roomBroadcaster != nil {
		s.roomBroadcaster.ParticipantLeft(rid, uid)
	}
	return nil
}

func (s *Service) GetRoomParticipants(ctx context.Context, roomID string) ([]LiveRoomParticipant, error) {
	rid, _ := uuid.Parse(roomID)
	return s.repo.GetRoomParticipants(ctx, rid)
}

func (s *Service) AddToRoomQueue(ctx context.Context, roomID, trackID, addedBy string) error {
	rid, _ := uuid.Parse(roomID)
	tid, _ := uuid.Parse(trackID)
	uid, _ := uuid.Parse(addedBy)
	// Get current max position
	queue, _ := s.repo.GetRoomQueue(ctx, rid)
	pos := len(queue)
	item := &LiveRoomQueueItem{
		RoomID:    rid,
		TrackID:   tid,
		AddedBy:   uid,
		Position:  pos,
		Status:    "queued",
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.AddToRoomQueue(ctx, item)
}

func (s *Service) GetRoomQueue(ctx context.Context, roomID string) ([]LiveRoomQueueItem, error) {
	rid, _ := uuid.Parse(roomID)
	return s.repo.GetRoomQueue(ctx, rid)
}

// --- Music Clubs ---

func (s *Service) GetClub(ctx context.Context, id string) (*MusicClub, error) {
	uid, _ := uuid.Parse(id)
	return s.repo.GetClub(ctx, uid)
}

func (s *Service) ListClubs(ctx context.Context, limit, offset int) ([]MusicClub, error) {
	return s.repo.ListClubs(ctx, limit, offset)
}

func (s *Service) IsClubMember(ctx context.Context, clubID, userID string) (bool, error) {
	cid, _ := uuid.Parse(clubID)
	uid, _ := uuid.Parse(userID)
	return s.repo.IsClubMember(ctx, cid, uid)
}

func (s *Service) IsClubAdmin(ctx context.Context, clubID, userID string) (bool, error) {
	cid, _ := uuid.Parse(clubID)
	uid, _ := uuid.Parse(userID)
	return s.repo.IsClubAdmin(ctx, cid, uid)
}

func (s *Service) GetClubMembers(ctx context.Context, clubID string) ([]MusicClubMember, error) {
	cid, _ := uuid.Parse(clubID)
	return s.repo.GetClubMembers(ctx, cid)
}

func (s *Service) CreateClubPost(ctx context.Context, clubID, userID, content string) (*MusicClubPost, error) {
	cid, _ := uuid.Parse(clubID)
	uid, _ := uuid.Parse(userID)
	now := time.Now().UTC()
	post := &MusicClubPost{
		ClubID:    cid,
		UserID:    uid,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.CreateClubPost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) GetClubPosts(ctx context.Context, clubID string, limit, offset int) ([]MusicClubPost, error) {
	cid, _ := uuid.Parse(clubID)
	return s.repo.GetClubPosts(ctx, cid, limit, offset)
}

// --- Discussions ---

func (s *Service) CreateDiscussion(ctx context.Context, req CreateDiscussionRequest, userID string) (*Discussion, error) {
	uid, _ := uuid.Parse(userID)
	tid, _ := uuid.Parse(req.TargetID)
	now := time.Now().UTC()
	d := &Discussion{
		UserID:     uid,
		TargetType: req.TargetType,
		TargetID:   tid,
		Content:    req.Content,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if req.ParentID != "" {
		pid, _ := uuid.Parse(req.ParentID)
		d.ParentID = &pid
	}
	if err := s.repo.CreateDiscussion(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDiscussions(ctx context.Context, targetType, targetID string, limit, offset int) ([]Discussion, error) {
	tid, _ := uuid.Parse(targetID)
	return s.repo.GetDiscussions(ctx, targetType, tid, limit, offset)
}

func (s *Service) GetDiscussionReplies(ctx context.Context, parentID string) ([]Discussion, error) {
	pid, _ := uuid.Parse(parentID)
	return s.repo.GetDiscussionReplies(ctx, pid)
}

// --- Track Ratings ---

func (s *Service) CreateRating(ctx context.Context, req CreateRatingRequest, userID string) (*TrackRating, error) {
	uid, _ := uuid.Parse(userID)
	tid, _ := uuid.Parse(req.TrackID)
	now := time.Now().UTC()
	rating := &TrackRating{
		ID:        uuid.New(),
		UserID:    uid,
		TrackID:   tid,
		Rating:    req.Rating,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if req.Review != "" {
		rating.Review = &req.Review
	}
	if err := s.repo.CreateRating(ctx, rating); err != nil {
		return nil, err
	}
	return rating, nil
}

func (s *Service) GetTrackRatings(ctx context.Context, trackID string) ([]TrackRating, error) {
	tid, _ := uuid.Parse(trackID)
	return s.repo.GetTrackRatings(ctx, tid)
}

func (s *Service) GetTrackRatingAverage(ctx context.Context, trackID string) (float64, int, error) {
	tid, _ := uuid.Parse(trackID)
	return s.repo.GetTrackRatingAverage(ctx, tid)
}

// --- Club Enhancements (Phase 5) ---

func (s *Service) ListClubsWithGenre(ctx context.Context, genre string, limit, offset int) ([]MusicClub, error) {
	if genre == "" {
		return s.repo.ListClubs(ctx, limit, offset)
	}
	return s.repo.ListClubsByGenre(ctx, genre, limit, offset)
}

func (s *Service) GetClubDetail(ctx context.Context, clubID, userID string) (*ClubDetailResponse, error) {
	if s.clubSvc == nil {
		return nil, fmt.Errorf("club service not initialized")
	}
	return s.clubSvc.GetClubDetail(ctx, clubID, userID)
}

func (s *Service) LaunchListeningParty(ctx context.Context, clubID, initiatorID string, req LaunchPartyRequest) (*ListeningParty, error) {
	if s.clubSvc == nil {
		return nil, fmt.Errorf("club service not initialized")
	}
	return s.clubSvc.LaunchListeningParty(ctx, clubID, initiatorID, req)
}

func (s *Service) CreateClub(ctx context.Context, req CreateClubRequest, ownerID string) (*Club, error) {
	if s.clubSvc == nil {
		return nil, fmt.Errorf("club service not initialized")
	}
	return s.clubSvc.CreateClub(ctx, ownerID, req)
}

func (s *Service) JoinClub(ctx context.Context, clubID, userID string) error {
	if s.clubSvc == nil {
		return fmt.Errorf("club service not initialized")
	}
	return s.clubSvc.JoinClub(ctx, clubID, userID)
}

func (s *Service) LeaveClub(ctx context.Context, clubID, userID string) error {
	if s.clubSvc == nil {
		return fmt.Errorf("club service not initialized")
	}
	return s.clubSvc.LeaveClub(ctx, clubID, userID)
}

// --- Club Discussions (Phase 6) ---

func (s *Service) CreateClubDiscussion(ctx context.Context, clubID, authorID, title, body string) (*ClubDiscussion, error) {
	return s.discussionSvc.CreateDiscussion(ctx, clubID, authorID, title, body)
}

func (s *Service) ListClubDiscussions(ctx context.Context, clubID string, limit, offset int) ([]ClubDiscussion, error) {
	return s.discussionSvc.ListDiscussions(ctx, clubID, limit, offset)
}

func (s *Service) GetClubDiscussion(ctx context.Context, discussionID string) (*ClubDiscussion, error) {
	return s.discussionSvc.GetDiscussion(ctx, discussionID)
}

func (s *Service) CreateClubDiscussionReply(ctx context.Context, discussionID, authorID, body string) (*ClubDiscussionReply, error) {
	return s.discussionSvc.CreateReply(ctx, discussionID, authorID, body)
}

func (s *Service) GetClubDiscussionReplies(ctx context.Context, discussionID string) ([]ClubDiscussionReply, error) {
	return s.discussionSvc.GetReplies(ctx, discussionID)
}

func (s *Service) DeleteClubDiscussion(ctx context.Context, discussionID, userID string) error {
	return s.discussionSvc.DeleteDiscussion(ctx, discussionID, userID)
}

func (s *Service) DeleteClubDiscussionReply(ctx context.Context, replyID, userID string) error {
	return s.discussionSvc.DeleteReply(ctx, replyID, userID)
}
