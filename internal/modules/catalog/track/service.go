package track

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

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Track, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Track, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, uid)
}

func (s *Service) List(ctx context.Context, limit, offset int, publicOnly bool) ([]Track, error) {
	return s.repo.List(ctx, limit, offset, publicOnly)
}

func (s *Service) Random(ctx context.Context, limit int) ([]Track, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.Random(ctx, limit)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Track, error) {
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
func (s *Service) ListCredits(ctx context.Context, id string) ([]TrackCredit, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListCredits(ctx, uid)
}

func (s *Service) ReplaceCredits(ctx context.Context, id string, credits []TrackCreditRequest) ([]TrackCredit, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ReplaceCredits(ctx, uid, credits)
}
func (s *Service) ListArtists(ctx context.Context, id string) ([]TrackArtist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListArtists(ctx, uid)
}

func (s *Service) ReplaceArtists(ctx context.Context, id string, artists []TrackArtistRequest) ([]TrackArtist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ReplaceArtists(ctx, uid, artists)
}
