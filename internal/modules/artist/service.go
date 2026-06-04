package artist

import (
	"context"
	"music/internal/modules/catalog/transport"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Artist, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Artist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, uid)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Artist, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Artist, error) {
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

func (s *Service) Search(ctx context.Context, q string, limit, offset int) ([]Artist, error) {
	return s.repo.Search(ctx, q, limit, offset)
}
