package social

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"music/internal/modules/playlist"
)

var (
	ErrOwnerCannotLeave = errors.New("مالک باید قبل از خروج، مالکیت رو منتقل کنه")
	ErrNotClubMember    = errors.New("شما عضو این کلاب نیستید")
	ErrClubNotFound     = errors.New("کلاب یافت نشد")
	ErrAlreadyMember    = errors.New("شما قبلاً عضو این کلاب هستید")
	ErrLaunchNotMember  = errors.New("فقط اعضای کلاب می‌تونن مهمونی شروع کنن")
)

// PlaylistCollaborator is the subset of playlist.Service that ClubService needs.
type PlaylistCollaborator interface {
	CreatePlaylist(ctx context.Context, req playlist.CreatePlaylistRequest, userID uuid.UUID) (playlist.Playlist, error)
	SetCollaborative(ctx context.Context, playlistID string, collab bool) error
	AddCollaborator(ctx context.Context, playlistID, userID string) error
	RemoveCollaborator(ctx context.Context, playlistID, userID string) error
	ListPlaylistTracks(ctx context.Context, playlistID string) ([]playlist.PlaylistTrackItem, error)
}

// ClubService wraps club-specific operations that need external deps.
// It delegates to the main Service for common operations.
type ClubService struct {
	repo        Repository
	playlistSvc PlaylistCollaborator
	partySvc    *Service // for CreateParty
}

func NewClubService(repo Repository, playlistSvc PlaylistCollaborator, partySvc *Service) *ClubService {
	return &ClubService{repo: repo, playlistSvc: playlistSvc, partySvc: partySvc}
}

func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, s)
	s = strings.Trim(s, "-_")
	if s == "" {
		s = uuid.New().String()[:8]
	}
	return s
}

func (cs *ClubService) CreateClub(ctx context.Context, ownerID string, req CreateClubRequest) (*Club, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner id: %w", err)
	}

	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}

	now := time.Now().UTC()

	// Create a collaborative playlist as the club's backing playlist
	playlistName := req.Name + " — کلاب"
	desc := fmt.Sprintf("پلی‌لیست مشترک کلاب %s", req.Name)
	createReq := playlist.CreatePlaylistRequest{
		Name:        playlistName,
		Description: &desc,
		IsPublic:    true,
	}
	p, err := cs.playlistSvc.CreatePlaylist(ctx, createReq, ownerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to create backing playlist: %w", err)
	}

	// Mark playlist as collaborative
	if err := cs.playlistSvc.SetCollaborative(ctx, p.ID.String(), true); err != nil {
		return nil, fmt.Errorf("failed to set playlist as collaborative: %w", err)
	}

	// Add owner as collaborator on the playlist
	if err := cs.playlistSvc.AddCollaborator(ctx, p.ID.String(), ownerID); err != nil {
		return nil, fmt.Errorf("failed to add owner as collaborator: %w", err)
	}

	// Create the club
	club := &MusicClub{
		Name:        req.Name,
		Slug:        slug,
		Description: nil,
		CoverURL:    nil,
		Genre:       nil,
		PlaylistID:  &p.ID,
		CreatedBy:   ownerUUID,
		IsPublic:    true,
		MaxMembers:  1000,
		MemberCount: 1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if req.Description != "" {
		club.Description = &req.Description
	}
	if req.Genre != "" {
		club.Genre = &req.Genre
	}
	if !req.IsPublic {
		club.IsPublic = false
	}
	if req.MaxMembers > 0 {
		club.MaxMembers = req.MaxMembers
	}

	if err := cs.repo.CreateClub(ctx, club); err != nil {
		return nil, fmt.Errorf("failed to create club: %w", err)
	}

	return toClub(club), nil
}

func (cs *ClubService) JoinClub(ctx context.Context, clubID, userID string) error {
	cid, err := uuid.Parse(clubID)
	if err != nil {
		return fmt.Errorf("invalid club id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	// Check membership
	isMember, err := cs.repo.IsClubMember(ctx, cid, uid)
	if err != nil {
		return err
	}
	if isMember {
		return nil // idempotent
	}

	// Get club to find playlist_id
	club, err := cs.repo.GetClub(ctx, cid)
	if err != nil {
		return ErrClubNotFound
	}

	if err := cs.repo.JoinClub(ctx, cid, uid); err != nil {
		return err
	}

	// Add as playlist collaborator
	if club.PlaylistID != nil {
		_ = cs.playlistSvc.AddCollaborator(ctx, club.PlaylistID.String(), userID)
	}

	return nil
}

func (cs *ClubService) LeaveClub(ctx context.Context, clubID, userID string) error {
	cid, err := uuid.Parse(clubID)
	if err != nil {
		return fmt.Errorf("invalid club id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	// Owner cannot leave
	isAdmin, err := cs.repo.IsClubAdmin(ctx, cid, uid)
	if err != nil {
		return err
	}
	if isAdmin {
		return ErrOwnerCannotLeave
	}

	// Get club to find playlist_id
	club, err := cs.repo.GetClub(ctx, cid)
	if err != nil {
		return ErrClubNotFound
	}

	if err := cs.repo.LeaveClub(ctx, cid, uid); err != nil {
		return err
	}

	// Remove as playlist collaborator
	if club.PlaylistID != nil {
		_ = cs.playlistSvc.RemoveCollaborator(ctx, club.PlaylistID.String(), userID)
	}

	return nil
}

func (cs *ClubService) LaunchListeningParty(ctx context.Context, clubID, initiatorID string, req LaunchPartyRequest) (*ListeningParty, error) {
	cid, err := uuid.Parse(clubID)
	if err != nil {
		return nil, fmt.Errorf("invalid club id: %w", err)
	}
	uid, err := uuid.Parse(initiatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	// Verify membership
	isMember, err := cs.repo.IsClubMember(ctx, cid, uid)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrLaunchNotMember
	}

	club, err := cs.repo.GetClub(ctx, cid)
	if err != nil {
		return nil, ErrClubNotFound
	}

	// Create party using existing service
	party, err := cs.partySvc.CreateParty(ctx, CreatePartyRequest{
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	}, initiatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to create party: %w", err)
	}

	// Seed the party queue with club playlist tracks
	if club.PlaylistID != nil {
		tracks, err := cs.playlistSvc.ListPlaylistTracks(ctx, club.PlaylistID.String())
		if err == nil {
			for i, t := range tracks {
				if i == 0 {
					// Set first track as current
					tid := t.TrackID
					_ = cs.partySvc.UpdatePartyStatus(ctx, party.ID.String(), "active", uid.String(), &[]string{tid.String()}[0])
				}
				_ = cs.repo.AddToRoomQueue(ctx, &LiveRoomQueueItem{
					RoomID:    party.ID,
					TrackID:   t.TrackID,
					AddedBy:   uid,
					Position:  i + 1,
					Status:    "queued",
					CreatedAt: time.Now().UTC(),
				})
			}
		}
	}

	return party, nil
}

func (cs *ClubService) GetClubDetail(ctx context.Context, clubID, requestingUserID string) (*ClubDetailResponse, error) {
	cid, err := uuid.Parse(clubID)
	if err != nil {
		return nil, fmt.Errorf("invalid club id: %w", err)
	}

	club, err := cs.repo.GetClub(ctx, cid)
	if err != nil {
		return nil, ErrClubNotFound
	}

	members, err := cs.repo.GetClubMembers(ctx, cid)
	if err != nil {
		members = []MusicClubMember{}
	}

	posts, err := cs.repo.GetClubPosts(ctx, cid, 50, 0)
	if err != nil {
		posts = []MusicClubPost{}
	}

	isMember := false
	memberRole := ""
	if requestingUserID != "" {
		ruid, err := uuid.Parse(requestingUserID)
		if err == nil {
			isMember, _ = cs.repo.IsClubMember(ctx, cid, ruid)
			if isMember {
				for _, m := range members {
					if m.UserID == ruid {
						memberRole = m.Role
						break
					}
				}
			}
		}
	}

	trackCount := 0
	if club.PlaylistID != nil {
		tracks, err := cs.playlistSvc.ListPlaylistTracks(ctx, club.PlaylistID.String())
		if err == nil {
			trackCount = len(tracks)
		}
	}

	return &ClubDetailResponse{
		Club:       *club,
		IsMember:   isMember,
		MemberRole: memberRole,
		Members:    members,
		Posts:      posts,
		PostCount:  len(posts),
		TrackCount: trackCount,
	}, nil
}

// Club mirrors MusicClub for the new Club response.
// This is the API-facing model matching the task spec.
type Club struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	CoverURL    string    `json:"cover_url,omitempty"`
	Genre       string    `json:"genre,omitempty"`
	PlaylistID  string    `json:"playlist_id,omitempty"`
	CreatedBy   string    `json:"created_by"`
	MemberCount int       `json:"member_count"`
	IsPublic    bool      `json:"is_public"`
	MaxMembers  int       `json:"max_members"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toClub(mc *MusicClub) *Club {
	c := &Club{
		ID:          mc.ID.String(),
		Name:        mc.Name,
		Slug:        mc.Slug,
		CreatedBy:   mc.CreatedBy.String(),
		MemberCount: mc.MemberCount,
		IsPublic:    mc.IsPublic,
		MaxMembers:  mc.MaxMembers,
		CreatedAt:   mc.CreatedAt,
		UpdatedAt:   mc.UpdatedAt,
	}
	if mc.Description != nil {
		c.Description = *mc.Description
	}
	if mc.CoverURL != nil {
		c.CoverURL = *mc.CoverURL
	}
	if mc.Genre != nil {
		c.Genre = *mc.Genre
	}
	if mc.PlaylistID != nil {
		c.PlaylistID = mc.PlaylistID.String()
	}
	return c
}
