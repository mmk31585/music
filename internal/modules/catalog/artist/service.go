package artist

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
func (s *Service) GetOverview(ctx context.Context, id string) (*Overview, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetOverview(ctx, uid)
}

func (s *Service) ListTracks(ctx context.Context, id string, limit, offset int) ([]ArtistTrack, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListTracks(ctx, uid, limit, offset)
}

func (s *Service) ListAlbums(ctx context.Context, id string, limit, offset int) ([]ArtistAlbum, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListAlbumsByType(ctx, uid, []string{"album", "ep", "compilation", "live"}, limit, offset)
}

func (s *Service) ListSingles(ctx context.Context, id string, limit, offset int) ([]ArtistAlbum, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListAlbumsByType(ctx, uid, []string{"single"}, limit, offset)
}

func (s *Service) ListAppearsOn(ctx context.Context, id string, limit, offset int) ([]ArtistTrack, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListAppearsOn(ctx, uid, limit, offset)
}
func (s *Service) ListTopTracks(ctx context.Context, id string, limit, offset int) ([]ArtistTrack, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListTopTracks(ctx, uid, limit, offset)
}

func (s *Service) ListRelated(ctx context.Context, id string, limit, offset int) ([]RelatedArtist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ListRelated(ctx, uid, limit, offset)
}

func (s *Service) ReplaceRelated(ctx context.Context, id string, related []RelatedArtistRequest) ([]RelatedArtist, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ReplaceRelated(ctx, uid, related)
}
func (s *Service) ReplaceTopTracks(ctx context.Context, id string, tracks []ArtistTopTrackRequest) ([]ArtistTrack, error) {
	uid, err := common.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.ReplaceTopTracks(ctx, uid, tracks)
}
