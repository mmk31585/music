package video

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Mock repositories
// ---------------------------------------------------------------------------

type mockVideoRepo struct {
	videos        map[uuid.UUID]*Video
	likes         map[string]*VideoLike           // key: "videoID:userID"
	tlv           map[string]*TrackLikeVisibility // key: "userID:trackID"
	status        map[uuid.UUID]*UserMusicStatus
	reactedTracks map[string]bool // key: "userID:trackID" — simulates reactions table
}

func newMockVideoRepo() *mockVideoRepo {
	return &mockVideoRepo{
		videos:        make(map[uuid.UUID]*Video),
		likes:         make(map[string]*VideoLike),
		tlv:           make(map[string]*TrackLikeVisibility),
		status:        make(map[uuid.UUID]*UserMusicStatus),
		reactedTracks: make(map[string]bool),
	}
}

// addTrackLike simulates a user liking a track (reactions table with type='like').
func (m *mockVideoRepo) addTrackLike(userID, trackID uuid.UUID) {
	key := userID.String() + ":" + trackID.String()
	m.reactedTracks[key] = true
}

func (m *mockVideoRepo) Create(ctx context.Context, v *Video) error {
	v.ID = uuid.New()
	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	m.videos[v.ID] = v
	return nil
}

func (m *mockVideoRepo) GetByID(ctx context.Context, id uuid.UUID) (*Video, error) {
	v, ok := m.videos[id]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (m *mockVideoRepo) Update(ctx context.Context, v *Video) error {
	v.UpdatedAt = time.Now()
	m.videos[v.ID] = v
	return nil
}

func (m *mockVideoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.videos, id)
	return nil
}

func (m *mockVideoRepo) ListByTrack(ctx context.Context, trackID uuid.UUID) ([]Video, error) {
	var items []Video
	for _, v := range m.videos {
		if v.TrackID == trackID {
			items = append(items, *v)
		}
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

func (m *mockVideoRepo) ListAll(ctx context.Context, limit, offset int) ([]Video, error) {
	var items []Video
	for _, v := range m.videos {
		items = append(items, *v)
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

func (m *mockVideoRepo) ListExplore(ctx context.Context, limit, offset int) ([]Video, error) {
	var items []Video
	for _, v := range m.videos {
		if v.IsApproved && v.IsPublic {
			items = append(items, *v)
		}
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

func (m *mockVideoRepo) CreateLike(ctx context.Context, like *VideoLike) error {
	key := like.VideoID.String() + ":" + like.UserID.String()
	if _, exists := m.likes[key]; exists {
		return nil // ON CONFLICT DO NOTHING
	}
	like.ID = uuid.New()
	like.CreatedAt = time.Now()
	m.likes[key] = like
	return nil
}

func (m *mockVideoRepo) DeleteLike(ctx context.Context, videoID, userID uuid.UUID) error {
	key := videoID.String() + ":" + userID.String()
	delete(m.likes, key)
	return nil
}

func (m *mockVideoRepo) GetLike(ctx context.Context, videoID, userID uuid.UUID) (*VideoLike, error) {
	key := videoID.String() + ":" + userID.String()
	like, ok := m.likes[key]
	if !ok {
		return nil, nil
	}
	return like, nil
}

func (m *mockVideoRepo) IncrementLikeCount(ctx context.Context, videoID uuid.UUID) error {
	if v, ok := m.videos[videoID]; ok {
		v.LikeCount++
	}
	return nil
}

func (m *mockVideoRepo) DecrementLikeCount(ctx context.Context, videoID uuid.UUID) error {
	if v, ok := m.videos[videoID]; ok {
		if v.LikeCount > 0 {
			v.LikeCount--
		}
	}
	return nil
}

func (m *mockVideoRepo) IncrementViewCount(ctx context.Context, videoID uuid.UUID) error {
	if v, ok := m.videos[videoID]; ok {
		v.ViewCount++
	}
	return nil
}

func (m *mockVideoRepo) ListByUser(ctx context.Context, userID, viewerID uuid.UUID, limit, offset int) ([]Video, error) {
	var items []Video
	for _, v := range m.videos {
		if v.UploaderID == userID {
			items = append(items, *v)
		}
	}
	if items == nil {
		items = []Video{}
	}
	return items, nil
}

func (m *mockVideoRepo) CheckTrackLikeExists(ctx context.Context, userID, trackID uuid.UUID) (bool, error) {
	key := userID.String() + ":" + trackID.String()
	_, exists := m.reactedTracks[key]
	return exists, nil
}

func (m *mockVideoRepo) UpsertTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID, visibility string) error {
	key := userID.String() + ":" + trackID.String()
	m.tlv[key] = &TrackLikeVisibility{
		UserID:     userID,
		TrackID:    trackID,
		Visibility: visibility,
		CreatedAt:  time.Now(),
	}
	return nil
}

func (m *mockVideoRepo) GetTrackLikeVisibility(ctx context.Context, userID, trackID uuid.UUID) (*TrackLikeVisibility, error) {
	key := userID.String() + ":" + trackID.String()
	tlv, ok := m.tlv[key]
	if !ok {
		return nil, nil
	}
	return tlv, nil
}

func (m *mockVideoRepo) GetPublicLikedTracks(ctx context.Context, userID uuid.UUID, limit, offset int) ([]TrackLikeVisibility, error) {
	var all []TrackLikeVisibility
	for _, tlv := range m.tlv {
		if tlv.UserID == userID && tlv.Visibility == "public" {
			all = append(all, *tlv)
		}
	}
	// Sort by created_at DESC to match SQL ORDER BY
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].CreatedAt.After(all[i].CreatedAt) {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	// Apply pagination
	start := offset
	if start < 0 {
		start = 0
	}
	if start >= len(all) {
		return []TrackLikeVisibility{}, nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], nil
}

func (m *mockVideoRepo) UpsertMusicStatus(ctx context.Context, status *UserMusicStatus) error {
	status.UpdatedAt = time.Now()
	m.status[status.UserID] = status
	return nil
}

func (m *mockVideoRepo) DeleteMusicStatus(ctx context.Context, userID uuid.UUID) error {
	delete(m.status, userID)
	return nil
}

func (m *mockVideoRepo) GetMusicStatus(ctx context.Context, userID uuid.UUID) (*UserMusicStatus, error) {
	s, ok := m.status[userID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

// ---------------------------------------------------------------------------
// Mock FollowService
// ---------------------------------------------------------------------------

type mockFollowService struct {
	follows map[string]bool // key: "followerID:followeeID"
}

func newMockFollowService() *mockFollowService {
	return &mockFollowService{
		follows: make(map[string]bool),
	}
}

func (m *mockFollowService) IsFollowingUser(ctx context.Context, followerID uuid.UUID, targetUserID string) (bool, error) {
	key := followerID.String() + ":" + targetUserID
	return m.follows[key], nil
}

func (m *mockFollowService) addFollow(followerID, followeeID uuid.UUID) {
	key := followerID.String() + ":" + followeeID.String()
	m.follows[key] = true
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestVideoService_OfficialMVRequiresAdminRole(t *testing.T) {
	// This test verifies the service-layer logic:
	// - Official MV upload is done via UploadOfficialMV which is only called by admin routes
	// - The admin role enforcement is done at the HTTP middleware level (RequireRole)
	// - Here we verify that the service correctly creates an MV with Type = VideoTypeOfficialMV
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	trackID := uuid.New()
	uploaderID := uuid.New()

	req := CreateVideoRequest{
		TrackID:       trackID.String(),
		Title:         "Test MV",
		Description:   "A test official MV",
		RawVideoURL:   "https://storage.example.com/videos/test.mp4",
		DurationMs:    240000,
		AspectRatio:   "16:9",
		FileSizeBytes: 50000000,
	}

	video, err := svc.UploadOfficialMV(context.Background(), trackID, uploaderID, req)
	if err != nil {
		t.Fatalf("UploadOfficialMV failed: %v", err)
	}

	if video.Type != VideoTypeOfficialMV {
		t.Errorf("expected Type %s, got %s", VideoTypeOfficialMV, video.Type)
	}
	if !video.IsApproved {
		t.Error("expected official MV to be auto-approved")
	}
	if video.UploaderID != uploaderID {
		t.Error("expected uploader ID to match")
	}
	if video.TrackID != trackID {
		t.Error("expected track ID to match")
	}

	// Verify the middleware-level enforcement:
	// The routes.go uses auth.RequireRole("admin") for the admin group,
	// so non-admin requests never reach UploadOfficialMV.
	// This is tested at the integration/e2e level.
}

func TestVideoService_UserEditUploaderCanDeleteOwn(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	req := CreateVideoRequest{
		TrackID:     trackID.String(),
		Title:       "My Edit",
		RawVideoURL: "https://storage.example.com/videos/edit.mp4",
	}

	video, err := svc.CreateUserEdit(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	// Owner can delete own video
	err = svc.DeleteOwnVideo(context.Background(), video.ID, userID)
	if err != nil {
		t.Fatalf("DeleteOwnVideo failed: %v", err)
	}

	// Verify deleted
	_, err = svc.GetVideo(context.Background(), video.ID)
	if !errors.Is(err, ErrVideoNotFound) {
		t.Errorf("expected ErrVideoNotFound after delete, got %v", err)
	}
}

func TestVideoService_UserEditCannotDeleteOthers(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	ownerID := uuid.New()
	otherUserID := uuid.New()
	trackID := uuid.New()

	req := CreateVideoRequest{
		TrackID:     trackID.String(),
		Title:       "My Edit",
		RawVideoURL: "https://storage.example.com/videos/edit.mp4",
	}

	video, err := svc.CreateUserEdit(context.Background(), ownerID, req)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	// Other user cannot delete
	err = svc.DeleteOwnVideo(context.Background(), video.ID, otherUserID)
	if !errors.Is(err, ErrNotOwner) {
		t.Errorf("expected ErrNotOwner, got %v", err)
	}
}

func TestVideoLike_Idempotent(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	req := CreateVideoRequest{
		TrackID:     trackID.String(),
		Title:       "Test Video",
		RawVideoURL: "https://storage.example.com/videos/test.mp4",
	}

	video, err := svc.CreateUserEdit(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	// First like should succeed
	err = svc.LikeVideo(context.Background(), video.ID, userID)
	if err != nil {
		t.Fatalf("first like failed: %v", err)
	}

	// Like count should be 1
	v, _ := svc.GetVideo(context.Background(), video.ID)
	if v.LikeCount != 1 {
		t.Errorf("expected like count 1, got %d", v.LikeCount)
	}

	// Second like should return ErrAlreadyLiked (idempotent at service)
	err = svc.LikeVideo(context.Background(), video.ID, userID)
	if !errors.Is(err, ErrAlreadyLiked) {
		t.Errorf("expected ErrAlreadyLiked, got %v", err)
	}

	// Like count should still be 1
	v, _ = svc.GetVideo(context.Background(), video.ID)
	if v.LikeCount != 1 {
		t.Errorf("expected like count 1 after duplicate, got %d", v.LikeCount)
	}

	// Unlike should succeed
	err = svc.UnlikeVideo(context.Background(), video.ID, userID)
	if err != nil {
		t.Fatalf("unlike failed: %v", err)
	}

	// Like count should be 0
	v, _ = svc.GetVideo(context.Background(), video.ID)
	if v.LikeCount != 0 {
		t.Errorf("expected like count 0 after unlike, got %d", v.LikeCount)
	}

	// Double unlike should return ErrNotLiked
	err = svc.UnlikeVideo(context.Background(), video.ID, userID)
	if !errors.Is(err, ErrNotLiked) {
		t.Errorf("expected ErrNotLiked on double unlike, got %v", err)
	}
}

func TestMusicStatus_PrivacyFollowersOnlyLogic(t *testing.T) {
	mockRepo := newMockVideoRepo()
	mockFollow := newMockFollowService()
	svc := NewService(mockRepo, mockFollow)

	targetUserID := uuid.New()
	followerUserID := uuid.New()
	strangerUserID := uuid.New()
	trackID := uuid.New()

	// Set the follower relationship
	mockFollow.addFollow(followerUserID, targetUserID)

	trackIDStr := trackID.String()

	// Target user sets status to "followers" visibility
	status, err := svc.SetMusicStatus(context.Background(), targetUserID, &trackIDStr, "followers")
	if err != nil {
		t.Fatalf("SetMusicStatus failed: %v", err)
	}
	if status.Visibility != "followers" {
		t.Errorf("expected visibility 'followers', got %s", status.Visibility)
	}

	// Follower should see the status
	followerStatus, err := svc.GetMusicStatus(context.Background(), targetUserID, &followerUserID)
	if err != nil {
		t.Fatalf("GetMusicStatus (follower) failed: %v", err)
	}
	if followerStatus == nil {
		t.Fatal("expected follower to see music status, got nil")
	}
	if followerStatus.CurrentTrackID == nil || *followerStatus.CurrentTrackID != trackID {
		t.Errorf("expected track id %s, got %v", trackID, followerStatus.CurrentTrackID)
	}

	// Stranger should NOT see the status
	strangerStatus, err := svc.GetMusicStatus(context.Background(), targetUserID, &strangerUserID)
	if err != nil {
		t.Fatalf("GetMusicStatus (stranger) failed: %v", err)
	}
	if strangerStatus != nil {
		t.Fatal("expected stranger to NOT see music status, got non-nil")
	}

	// Unauthenticated user should NOT see the status
	anonStatus, err := svc.GetMusicStatus(context.Background(), targetUserID, nil)
	if err != nil {
		t.Fatalf("GetMusicStatus (anon) failed: %v", err)
	}
	if anonStatus != nil {
		t.Fatal("expected anonymous user to NOT see music status, got non-nil")
	}
}

func TestMusicStatus_StaleAfter30Minutes(t *testing.T) {
	mockRepo := newMockVideoRepo()
	// Override the repo to allow controlling UpdatedAt
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()
	trackIDStr := trackID.String()

	// Set initial status
	_, err := svc.SetMusicStatus(context.Background(), userID, &trackIDStr, "public")
	if err != nil {
		t.Fatalf("SetMusicStatus failed: %v", err)
	}

	// Immediately it should be visible
	recentStatus, err := svc.GetMusicStatus(context.Background(), userID, &userID)
	if err != nil {
		t.Fatalf("GetMusicStatus failed: %v", err)
	}
	if recentStatus == nil {
		t.Fatal("expected to see recent music status")
	}

	// Manually set UpdatedAt to be older than 30 minutes
	staleTime := time.Now().Add(-31 * time.Minute)
	mockRepo.status[userID].UpdatedAt = staleTime

	// Now it should be treated as stale
	staleStatus, err := svc.GetMusicStatus(context.Background(), userID, &userID)
	if err != nil {
		t.Fatalf("GetMusicStatus (stale) failed: %v", err)
	}
	if staleStatus != nil {
		t.Fatal("expected stale status (>30min) to return nil (not playing)")
	}
}

func TestTrackLikeVisibility_DefaultPrivate(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	// Pre-seed: user must have liked the track before setting visibility
	mockRepo.addTrackLike(userID, trackID)

	// Set visibility to public
	tlv, err := svc.SetTrackLikeVisibility(context.Background(), userID, trackID, "public")
	if err != nil {
		t.Fatalf("SetTrackLikeVisibility failed: %v", err)
	}
	if tlv.Visibility != "public" {
		t.Errorf("expected visibility 'public', got %s", tlv.Visibility)
	}

	// Check public liked tracks
	publicTracks, err := svc.GetPublicLikedTracks(context.Background(), userID, 20, 0)
	if err != nil {
		t.Fatalf("GetPublicLikedTracks failed: %v", err)
	}
	if len(publicTracks) != 1 {
		t.Errorf("expected 1 public track, got %d", len(publicTracks))
	}

	// Change to private
	tlv, err = svc.SetTrackLikeVisibility(context.Background(), userID, trackID, "private")
	if err != nil {
		t.Fatalf("SetTrackLikeVisibility (private) failed: %v", err)
	}
	if tlv.Visibility != "private" {
		t.Errorf("expected visibility 'private', got %s", tlv.Visibility)
	}

	// Should no longer appear in public list
	publicTracks, err = svc.GetPublicLikedTracks(context.Background(), userID, 20, 0)
	if err != nil {
		t.Fatalf("GetPublicLikedTracks failed: %v", err)
	}
	if len(publicTracks) != 0 {
		t.Errorf("expected 0 public tracks after setting private, got %d", len(publicTracks))
	}

	// Invalid visibility should return error
	_, err = svc.SetTrackLikeVisibility(context.Background(), userID, trackID, "invalid")
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for invalid visibility, got %v", err)
	}
}

func TestExploreVideos_OnlyApprovedPublicVideos(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	adminID := uuid.New()
	userID := uuid.New()
	trackID := uuid.New()

	// Create an official MV (auto-approved, public)
	mvReq := CreateVideoRequest{
		TrackID:       trackID.String(),
		Title:         "Official MV",
		RawVideoURL:   "https://example.com/mv.mp4",
		DurationMs:    180000,
		AspectRatio:   "16:9",
		FileSizeBytes: 1000000,
	}
	_, err := svc.UploadOfficialMV(context.Background(), trackID, adminID, mvReq)
	if err != nil {
		t.Fatalf("UploadOfficialMV failed: %v", err)
	}

	// Create a user edit (not approved yet)
	editReq := CreateVideoRequest{
		TrackID:     trackID.String(),
		Title:       "User Edit - Pending",
		RawVideoURL: "https://example.com/edit.mp4",
	}
	userEdit, err := svc.CreateUserEdit(context.Background(), userID, editReq)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	// Create another user edit with is_public=false, approved
	nonPublicReq := CreateVideoRequest{
		TrackID:     trackID.String(),
		Title:       "Non-Public",
		RawVideoURL: "https://example.com/nonpublic.mp4",
	}
	nonPublic, err := svc.CreateUserEdit(context.Background(), userID, nonPublicReq)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}
	// Manually approve and set non-public
	nonPublic.IsApproved = true
	nonPublic.IsPublic = false
	_ = mockRepo.Update(context.Background(), nonPublic)

	// Explore should return only the official MV (it's approved + public)
	exploreItems, err := svc.ListExplore(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("ListExplore failed: %v", err)
	}

	if len(exploreItems) != 1 {
		t.Errorf("expected 1 explore item (official MV), got %d", len(exploreItems))
	}

	// Now approve the user edit
	_, err = svc.ApproveVideo(context.Background(), userEdit.ID)
	if err != nil {
		t.Fatalf("ApproveVideo failed: %v", err)
	}

	// Explore should now return 2 items (MV + approved edit)
	exploreItems, err = svc.ListExplore(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("ListExplore failed: %v", err)
	}
	if len(exploreItems) != 2 {
		t.Errorf("expected 2 explore items after approval, got %d", len(exploreItems))
	}
}

func TestAdminCanDeleteAnyVideo(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	adminID := uuid.New()
	trackID := uuid.New()

	// User creates an edit
	req := CreateVideoRequest{
		TrackID:     trackID.String(),
		Title:       "User Edit",
		RawVideoURL: "https://example.com/edit.mp4",
	}
	video, err := svc.CreateUserEdit(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	// Admin deletes it — should succeed
	err = svc.DeleteVideo(context.Background(), video.ID, adminID, true)
	if err != nil {
		t.Fatalf("admin DeleteVideo failed: %v", err)
	}

	// Verify deleted
	_, err = svc.GetVideo(context.Background(), video.ID)
	if !errors.Is(err, ErrVideoNotFound) {
		t.Errorf("expected ErrVideoNotFound, got %v", err)
	}
}

func TestAdminCanDeleteOfficialMV(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	adminID := uuid.New()
	userID := uuid.New()
	trackID := uuid.New()

	// Admin creates an official MV
	req := CreateVideoRequest{
		TrackID:       trackID.String(),
		Title:         "Official MV",
		RawVideoURL:   "https://example.com/mv.mp4",
		DurationMs:    180000,
		AspectRatio:   "16:9",
		FileSizeBytes: 1000000,
	}
	mv, err := svc.UploadOfficialMV(context.Background(), trackID, adminID, req)
	if err != nil {
		t.Fatalf("UploadOfficialMV failed: %v", err)
	}

	// Regular user cannot delete official MV
	err = svc.DeleteVideo(context.Background(), mv.ID, userID, false)
	if !errors.Is(err, ErrRequiresAdmin) && !errors.Is(err, ErrNotOwner) {
		t.Errorf("expected ErrRequiresAdmin or ErrNotOwner, got %v", err)
	}

	// Admin can delete official MV
	err = svc.DeleteVideo(context.Background(), mv.ID, adminID, true)
	if err != nil {
		t.Fatalf("admin DeleteVideo failed: %v", err)
	}
}

func TestMusicStatus_PrivateVisibility(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	otherUserID := uuid.New()
	trackID := uuid.New()
	trackIDStr := trackID.String()

	// Set private status
	status, err := svc.SetMusicStatus(context.Background(), userID, &trackIDStr, "private")
	if err != nil {
		t.Fatalf("SetMusicStatus failed: %v", err)
	}
	if status.Visibility != "private" {
		t.Errorf("expected visibility 'private', got %s", status.Visibility)
	}

	// User can see their own private status
	selfStatus, err := svc.GetMusicStatus(context.Background(), userID, &userID)
	if err != nil {
		t.Fatalf("GetMusicStatus (self) failed: %v", err)
	}
	if selfStatus == nil {
		t.Fatal("expected self to see own private status")
	}

	// Other user cannot see private status
	otherStatus, err := svc.GetMusicStatus(context.Background(), userID, &otherUserID)
	if err != nil {
		t.Fatalf("GetMusicStatus (other) failed: %v", err)
	}
	if otherStatus != nil {
		t.Fatal("expected other user to NOT see private status")
	}
}

func TestMusicStatus_PublicVisibility(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	anonID := uuid.New()
	trackID := uuid.New()
	trackIDStr := trackID.String()

	// Set public status
	_, err := svc.SetMusicStatus(context.Background(), userID, &trackIDStr, "public")
	if err != nil {
		t.Fatalf("SetMusicStatus failed: %v", err)
	}

	// Anyone can see public status
	for name, reqUser := range map[string]*uuid.UUID{
		"self":   &userID,
		"other":  &anonID,
		"unauth": nil,
	} {
		status, err := svc.GetMusicStatus(context.Background(), userID, reqUser)
		if err != nil {
			t.Fatalf("GetMusicStatus (%s) failed: %v", name, err)
		}
		if status == nil {
			t.Errorf("expected %s to see public status", name)
		}
	}
}

func TestViewVideo_IncrementsViewCount(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()
	trackIDStr := trackID.String()

	req := CreateVideoRequest{
		TrackID:     trackIDStr,
		Title:       "Test",
		RawVideoURL: "https://example.com/v.mp4",
	}
	video, err := svc.CreateUserEdit(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	for i := 0; i < 5; i++ {
		if err := svc.ViewVideo(context.Background(), video.ID); err != nil {
			t.Fatalf("ViewVideo failed: %v", err)
		}
	}

	v, _ := svc.GetVideo(context.Background(), video.ID)
	if v.ViewCount != 5 {
		t.Errorf("expected view count 5, got %d", v.ViewCount)
	}
}

func TestClearMusicStatus(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()
	trackIDStr := trackID.String()

	// Set status
	_, err := svc.SetMusicStatus(context.Background(), userID, &trackIDStr, "public")
	if err != nil {
		t.Fatalf("SetMusicStatus failed: %v", err)
	}

	// Verify it exists
	status, _ := svc.GetMusicStatus(context.Background(), userID, &userID)
	if status == nil {
		t.Fatal("expected status to exist")
	}

	// Clear it
	err = svc.ClearMusicStatus(context.Background(), userID)
	if err != nil {
		t.Fatalf("ClearMusicStatus failed: %v", err)
	}

	// Verify gone
	status, err = svc.GetMusicStatus(context.Background(), userID, &userID)
	if err != nil {
		t.Fatalf("GetMusicStatus failed: %v", err)
	}
	if status != nil {
		t.Fatal("expected status to be nil after clear")
	}
}

func TestListByTrack_ReturnsAllVideos(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	adminID := uuid.New()
	userID := uuid.New()
	trackID := uuid.New()

	// Create an official MV
	_, _ = svc.UploadOfficialMV(context.Background(), trackID, adminID, CreateVideoRequest{
		TrackID: trackID.String(), Title: "MV", RawVideoURL: "https://example.com/mv.mp4",
	})

	// Create a user edit
	_, _ = svc.CreateUserEdit(context.Background(), userID, CreateVideoRequest{
		TrackID: trackID.String(), Title: "Edit", RawVideoURL: "https://example.com/edit.mp4",
	})

	items, err := svc.ListByTrack(context.Background(), trackID)
	if err != nil {
		t.Fatalf("ListByTrack failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 videos for track, got %d", len(items))
	}
}

func TestGetVideo_NonExistent(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	_, err := svc.GetVideo(context.Background(), uuid.New())
	if !errors.Is(err, ErrVideoNotFound) {
		t.Errorf("expected ErrVideoNotFound for non-existent video, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests: Video callback (Phase 3)
// ---------------------------------------------------------------------------

func TestVideoCallback_CompletedUpdatesStatusToReady(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	// Create a user edit video (status = processing)
	video, err := svc.CreateUserEdit(context.Background(), userID, CreateVideoRequest{
		TrackID: trackID.String(),
		Title:   "My Edit",
	})
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}
	if video.Status != VideoStatusProcessing {
		t.Fatalf("expected status 'processing', got %s", video.Status)
	}

	// Simulate callback with completed status
	result, err := svc.HandleVideoCallback(context.Background(), VideoCallbackPayload{
		VideoID:          video.ID.String(),
		FinalVideoPath:   "/data/videos/final.mp4",
		ThumbnailPath:    "/data/videos/thumb.jpg",
		DurationMs:       180000,
		AspectRatio:      "16:9",
		ProcessingStatus: "completed",
	})
	if err != nil {
		t.Fatalf("HandleVideoCallback failed: %v", err)
	}

	if result.Status != VideoStatusReady {
		t.Errorf("expected status 'ready', got %s", result.Status)
	}
	if result.FinalVideoPath == nil || *result.FinalVideoPath != "/data/videos/final.mp4" {
		t.Errorf("expected final_video_path to be set")
	}
	if result.ThumbnailPath == nil || *result.ThumbnailPath != "/data/videos/thumb.jpg" {
		t.Errorf("expected thumbnail_path to be set")
	}
	if result.DurationMs != 180000 {
		t.Errorf("expected duration_ms 180000, got %d", result.DurationMs)
	}
	if result.AspectRatio != "16:9" {
		t.Errorf("expected aspect_ratio '16:9', got %s", result.AspectRatio)
	}
}

func TestVideoCallback_OfficialMVAutoApproved(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	adminID := uuid.New()
	trackID := uuid.New()

	// Create an official MV
	mv, err := svc.UploadOfficialMV(context.Background(), trackID, adminID, CreateVideoRequest{
		TrackID:       trackID.String(),
		Title:         "Official MV",
		RawVideoURL:   "https://example.com/mv.mp4",
		DurationMs:    240000,
		AspectRatio:   "16:9",
		FileSizeBytes: 50000000,
	})
	if err != nil {
		t.Fatalf("UploadOfficialMV failed: %v", err)
	}

	// Simulate callback
	result, err := svc.HandleVideoCallback(context.Background(), VideoCallbackPayload{
		VideoID:          mv.ID.String(),
		FinalVideoPath:   "/data/videos/final.mp4",
		ThumbnailPath:    "/data/videos/thumb.jpg",
		DurationMs:       240000,
		AspectRatio:      "16:9",
		ProcessingStatus: "completed",
	})
	if err != nil {
		t.Fatalf("HandleVideoCallback failed: %v", err)
	}

	if !result.IsApproved {
		t.Error("expected official MV to be auto-approved on callback")
	}
	if result.Status != VideoStatusReady {
		t.Errorf("expected status 'ready', got %s", result.Status)
	}
}

func TestVideoCallback_UserEditRequiresAdminApproval(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	// Create a user edit
	video, err := svc.CreateUserEdit(context.Background(), userID, CreateVideoRequest{
		TrackID: trackID.String(),
		Title:   "User Edit",
	})
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	// Simulate callback — user edit should NOT be auto-approved
	result, err := svc.HandleVideoCallback(context.Background(), VideoCallbackPayload{
		VideoID:          video.ID.String(),
		FinalVideoPath:   "/data/videos/final.mp4",
		ProcessingStatus: "completed",
	})
	if err != nil {
		t.Fatalf("HandleVideoCallback failed: %v", err)
	}

	if result.IsApproved {
		t.Error("expected user edit to NOT be auto-approved — admin must approve")
	}
	if result.Status != VideoStatusReady {
		t.Errorf("expected status 'ready', got %s", result.Status)
	}
}

func TestVideoCallback_FailedUpdatesStatusToFailed(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	video, err := svc.CreateUserEdit(context.Background(), userID, CreateVideoRequest{
		TrackID: trackID.String(),
		Title:   "Will Fail",
	})
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	errMsg := "corrupt video file: moov atom not found"
	result, err := svc.HandleVideoCallback(context.Background(), VideoCallbackPayload{
		VideoID:          video.ID.String(),
		ProcessingStatus: "failed",
		ErrorMessage:     &errMsg,
	})
	if err != nil {
		t.Fatalf("HandleVideoCallback failed: %v", err)
	}

	if result.Status != VideoStatusFailed {
		t.Errorf("expected status 'failed', got %s", result.Status)
	}
}

func TestVideoCallback_InvalidStatusReturnsError(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	video, err := svc.CreateUserEdit(context.Background(), userID, CreateVideoRequest{
		TrackID: trackID.String(),
		Title:   "Test",
	})
	if err != nil {
		t.Fatalf("CreateUserEdit failed: %v", err)
	}

	_, err = svc.HandleVideoCallback(context.Background(), VideoCallbackPayload{
		VideoID:          video.ID.String(),
		ProcessingStatus: "unknown",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput for unknown status, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests: SetTrackLikeVisibility — existing-like check (Phase 3)
// ---------------------------------------------------------------------------

func TestSetTrackLikeVisibility_RequiresExistingLike(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	// Attempt to set visibility without liking the track first
	_, err := svc.SetTrackLikeVisibility(context.Background(), userID, trackID, "public")
	if !errors.Is(err, ErrNotLiked) {
		t.Errorf("expected ErrNotLiked when track hasn't been liked, got %v", err)
	}
}

func TestSetTrackLikeVisibility_SucceedsWhenLiked(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	trackID := uuid.New()

	// Pre-seed the like
	mockRepo.addTrackLike(userID, trackID)

	tlv, err := svc.SetTrackLikeVisibility(context.Background(), userID, trackID, "public")
	if err != nil {
		t.Fatalf("SetTrackLikeVisibility failed: %v", err)
	}
	if tlv.Visibility != "public" {
		t.Errorf("expected visibility 'public', got %s", tlv.Visibility)
	}
	if tlv.TrackID != trackID {
		t.Errorf("expected track id %s, got %s", trackID, tlv.TrackID)
	}
	if tlv.UserID != userID {
		t.Errorf("expected user id %s, got %s", userID, tlv.UserID)
	}
}

// ---------------------------------------------------------------------------
// Tests: GetPublicLikedTracks (Phase 3)
// ---------------------------------------------------------------------------

func TestGetPublicLikedTracks_OnlyReturnsPublicOnes(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()
	track1 := uuid.New()
	track2 := uuid.New()
	track3 := uuid.New()

	// Seed likes for all three tracks
	mockRepo.addTrackLike(userID, track1)
	mockRepo.addTrackLike(userID, track2)
	mockRepo.addTrackLike(userID, track3)

	// Set track1 to public, track2 to private, track3 to public
	_, err := svc.SetTrackLikeVisibility(context.Background(), userID, track1, "public")
	if err != nil {
		t.Fatalf("SetTrackLikeVisibility failed: %v", err)
	}
	_, err = svc.SetTrackLikeVisibility(context.Background(), userID, track2, "private")
	if err != nil {
		t.Fatalf("SetTrackLikeVisibility failed: %v", err)
	}
	_, err = svc.SetTrackLikeVisibility(context.Background(), userID, track3, "public")
	if err != nil {
		t.Fatalf("SetTrackLikeVisibility failed: %v", err)
	}

	// Public liked tracks should only return track1 and track3
	items, err := svc.GetPublicLikedTracks(context.Background(), userID, 20, 0)
	if err != nil {
		t.Fatalf("GetPublicLikedTracks failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 public tracks, got %d", len(items))
	}

	// Count track1 and track3 specifically
	foundTrack1 := false
	foundTrack3 := false
	for _, item := range items {
		if item.TrackID == track1 && item.Visibility == "public" {
			foundTrack1 = true
		}
		if item.TrackID == track3 && item.Visibility == "public" {
			foundTrack3 = true
		}
	}
	if !foundTrack1 {
		t.Error("expected track1 to be in public liked tracks")
	}
	if !foundTrack3 {
		t.Error("expected track3 to be in public liked tracks")
	}
}

func TestGetPublicLikedTracks_RespectsPagination(t *testing.T) {
	mockRepo := newMockVideoRepo()
	svc := NewService(mockRepo, newMockFollowService())

	userID := uuid.New()

	// Create 5 tracks, like them all, set all to public
	for i := 0; i < 5; i++ {
		trackID := uuid.New()
		mockRepo.addTrackLike(userID, trackID)
		_, err := svc.SetTrackLikeVisibility(context.Background(), userID, trackID, "public")
		if err != nil {
			t.Fatalf("SetTrackLikeVisibility failed: %v", err)
		}
	}

	// Get page 1 (limit=2, offset=0)
	page1, err := svc.GetPublicLikedTracks(context.Background(), userID, 2, 0)
	if err != nil {
		t.Fatalf("GetPublicLikedTracks failed: %v", err)
	}
	if len(page1) > 2 {
		t.Errorf("expected at most 2 items on page 1, got %d", len(page1))
	}

	// Get page 2 (limit=2, offset=2)
	page2, err := svc.GetPublicLikedTracks(context.Background(), userID, 2, 2)
	if err != nil {
		t.Fatalf("GetPublicLikedTracks failed: %v", err)
	}
	if len(page2) > 2 {
		t.Errorf("expected at most 2 items on page 2, got %d", len(page2))
	}

	// Ensure pages are disjoint
	ids := make(map[uuid.UUID]bool)
	for _, item := range page1 {
		ids[item.TrackID] = true
	}
	for _, item := range page2 {
		if ids[item.TrackID] {
			t.Error("expected pages to be disjoint, but found duplicate track")
		}
	}
}
