package album

import (
	"context"
	"music/internal/modules/catalog/common"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Album, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Album, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, uid)
}

func (s *Service) List(ctx context.Context, limit, offset int, opts ...ListOptions) ([]Album, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	return s.repo.List(ctx, limit, offset, o)
}

func (s *Service) ListByArtist(ctx context.Context, artistID string) ([]Album, error) {
	uid, err := common.ParseUUID(artistID)
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, 100, 0, ListOptions{ArtistID: &uid})
}

func (s *Service) ListTracks(ctx context.Context, id string, limit, offset int) ([]AlbumTrack, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListTracks(ctx, uid, limit, offset)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Album, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, uid, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, uid)
}
func (s *Service) ListArtists(ctx context.Context, id string) ([]AlbumArtist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListArtists(ctx, uid)
}

func (s *Service) ReplaceArtists(ctx context.Context, id string, artists []AlbumArtistRequest) ([]AlbumArtist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ReplaceArtists(ctx, uid, artists)
}
