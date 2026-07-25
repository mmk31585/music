package stats

import "context"

// Service provides catalog stats operations.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll returns aggregate counts across catalog entities.
func (s *Service) GetAll(ctx context.Context) (StatsResponse, error) {
	return s.repo.CountAll(ctx)
}
