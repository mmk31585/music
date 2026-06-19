package playlist

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrForbiddenPlaylistAccess = errors.New("forbidden playlist access")
	ErrInvalidTrackPosition    = errors.New("invalid track position")
	ErrInvalidPlaylistName     = errors.New("invalid playlist name")
)

type RepositoryInterface interface {
	CreatePlaylist(ctx context.Context, req CreatePlaylistRequest, userID uuid.UUID) (Playlist, error)
	UpdatePlaylist(ctx context.Context, playlistID, userID uuid.UUID, req UpdatePlaylistRequest) (Playlist, error)
	DeletePlaylist(ctx context.Context, playlistID, userID uuid.UUID) error
	GetPlaylistByID(ctx context.Context, playlistID uuid.UUID) (Playlist, error)
	ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]PlaylistTrackItem, error)
	ListPublicPlaylists(ctx context.Context) ([]PlaylistListItemResponse, error)
	ListUserPlaylists(ctx context.Context, userID uuid.UUID) ([]PlaylistListItemResponse, error)
	AddTrack(ctx context.Context, playlistID, trackID uuid.UUID) error
	RemoveTrack(ctx context.Context, playlistID, trackID uuid.UUID) error
	ReorderTrack(ctx context.Context, playlistID, trackID uuid.UUID, newPosition int) error
	IsCollaborator(ctx context.Context, playlistID, userID string) (bool, error)
	IsCollaborativePlaylist(ctx context.Context, playlistID string) (bool, error)
	SetCollaborative(ctx context.Context, playlistID string, collab bool) error
	AddCollaborator(ctx context.Context, playlistID, userID string) error
	RemoveCollaborator(ctx context.Context, playlistID, userID string) error
	ListCollaborators(ctx context.Context, playlistID string) ([]CollaboratorResponse, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePlaylist(ctx context.Context, req CreatePlaylistRequest, userID uuid.UUID) (Playlist, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return Playlist{}, ErrInvalidPlaylistName
	}

	return s.repo.CreatePlaylist(ctx, req, userID)
}

func (s *Service) UpdatePlaylist(ctx context.Context, playlistID, userID uuid.UUID, req UpdatePlaylistRequest) (Playlist, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return Playlist{}, ErrInvalidPlaylistName
	}

	return s.repo.UpdatePlaylist(ctx, playlistID, userID, req)
}

func (s *Service) DeletePlaylist(ctx context.Context, playlistID, userID uuid.UUID) error {
	p, err := s.repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbiddenPlaylistAccess
	}
	return s.repo.DeletePlaylist(ctx, playlistID, userID)
}

func (s *Service) GetPlaylist(ctx context.Context, playlistID uuid.UUID, requesterID *uuid.UUID) (Playlist, []PlaylistTrackItem, error) {
	p, err := s.repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		return Playlist{}, nil, err
	}

	if !p.IsPublic {
		if requesterID == nil || *requesterID != p.UserID {
			return Playlist{}, nil, ErrForbiddenPlaylistAccess
		}
	}

	tracks, err := s.repo.ListPlaylistTracks(ctx, playlistID)
	if err != nil {
		return Playlist{}, nil, err
	}

	return p, tracks, nil
}

func (s *Service) ListPublicPlaylists(ctx context.Context) ([]PlaylistListItemResponse, error) {
	return s.repo.ListPublicPlaylists(ctx)
}

func (s *Service) ListMyPlaylists(ctx context.Context, userID uuid.UUID) ([]PlaylistListItemResponse, error) {
	return s.repo.ListUserPlaylists(ctx, userID)
}

func (s *Service) AddTrack(ctx context.Context, playlistID, userID, trackID uuid.UUID) error {
	p, err := s.repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbiddenPlaylistAccess
	}

	return s.repo.AddTrack(ctx, playlistID, trackID)
}

func (s *Service) RemoveTrack(ctx context.Context, playlistID, userID, trackID uuid.UUID) error {
	p, err := s.repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbiddenPlaylistAccess
	}

	return s.repo.RemoveTrack(ctx, playlistID, trackID)
}

func (s *Service) ReorderTrack(ctx context.Context, playlistID, userID, trackID uuid.UUID, newPosition int) error {
	p, err := s.repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbiddenPlaylistAccess
	}

	return s.repo.ReorderTrack(ctx, playlistID, trackID, newPosition)
}

// Collaborator delegation methods
func (s *Service) IsCollaborator(ctx context.Context, playlistID, userID string) (bool, error) {
	return s.repo.IsCollaborator(ctx, playlistID, userID)
}

func (s *Service) IsCollaborativePlaylist(ctx context.Context, playlistID string) (bool, error) {
	return s.repo.IsCollaborativePlaylist(ctx, playlistID)
}

func (s *Service) SetCollaborative(ctx context.Context, playlistID string, collab bool) error {
	return s.repo.SetCollaborative(ctx, playlistID, collab)
}

func (s *Service) AddCollaborator(ctx context.Context, playlistID, userID string) error {
	return s.repo.AddCollaborator(ctx, playlistID, userID)
}

func (s *Service) RemoveCollaborator(ctx context.Context, playlistID, userID string) error {
	return s.repo.RemoveCollaborator(ctx, playlistID, userID)
}

func (s *Service) ListCollaborators(ctx context.Context, playlistID string) ([]CollaboratorResponse, error) {
	return s.repo.ListCollaborators(ctx, playlistID)
}

func (s *Service) ListPlaylistTracks(ctx context.Context, playlistID string) ([]PlaylistTrackItem, error) {
	pid, err := uuid.Parse(playlistID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListPlaylistTracks(ctx, pid)
}
