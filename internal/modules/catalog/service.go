package catalog

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// --------------------
// Slugs
// --------------------

func (s *Service) uniqueArtistSlug(ctx context.Context, name string, excludeID *uuid.UUID) (string, error) {
	base := makeSlug(name)
	slug := base

	for i := 2; i <= 1000; i++ {
		exists, err := s.repo.ArtistSlugExists(ctx, slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = base + "-" + strconv.Itoa(i)
	}

	return "", ErrDuplicateArtistSlug
}

func (s *Service) uniqueAlbumSlug(ctx context.Context, artistID uuid.UUID, title string, excludeID *uuid.UUID) (string, error) {
	base := makeSlug(title)
	slug := base

	for i := 2; i <= 1000; i++ {
		exists, err := s.repo.AlbumSlugExists(ctx, artistID, slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = base + "-" + strconv.Itoa(i)
	}

	return "", ErrDuplicateAlbumSlug
}

func (s *Service) uniqueTrackSlug(ctx context.Context, artistID uuid.UUID, title string, excludeID *uuid.UUID) (string, error) {
	base := makeSlug(title)
	slug := base

	for i := 2; i <= 1000; i++ {
		exists, err := s.repo.TrackSlugExists(ctx, artistID, slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = base + "-" + strconv.Itoa(i)
	}

	return "", ErrDuplicateTrackSlug
}

func (s *Service) uniqueGenreSlug(ctx context.Context, name string, excludeID *uuid.UUID) (string, error) {
	base := makeSlug(name)
	slug := base

	for i := 2; i <= 1000; i++ {
		exists, err := s.repo.GenreSlugExists(ctx, slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = base + "-" + strconv.Itoa(i)
	}

	return "", ErrDuplicateGenreSlug
}

// --------------------
// Artists
// --------------------

func (s *Service) CreateArtist(ctx context.Context, req CreateArtistRequest) (*Artist, error) {
	slug, err := s.uniqueArtistSlug(ctx, req.Name, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	isVerified := false
	if req.IsVerified != nil {
		isVerified = *req.IsVerified
	}

	monthlyListeners := int64(0)
	if req.MonthlyListeners != nil {
		monthlyListeners = *req.MonthlyListeners
	}
	fmt.Print(req)
	artist := &Artist{
		ID:               uuid.New(),
		Name:             req.Name,
		Slug:             slug,
		Bio:              req.Bio,
		ImageURL:         req.ImageURL,
		IsVerified:       isVerified,
		MonthlyListeners: monthlyListeners,
		CreatedAt:        now,
		UpdatedAt:        &now,
	}

	if err := s.repo.CreateArtist(ctx, artist); err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *Service) GetArtist(ctx context.Context, id uuid.UUID) (*Artist, error) {
	return s.repo.GetArtistByID(ctx, id)
}

func (s *Service) ListArtists(ctx context.Context, filter ArtistListFilter) ([]Artist, int, error) {
	artists, err := s.repo.ListArtists(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountArtists(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return artists, total, nil
}

func (s *Service) UpdateArtist(ctx context.Context, id uuid.UUID, req UpdateArtistRequest) (*Artist, error) {
	artist, err := s.repo.GetArtistByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		slug, err := s.uniqueArtistSlug(ctx, *req.Name, &id)
		if err != nil {
			return nil, err
		}
		artist.Name = *req.Name
		artist.Slug = slug
	}

	if req.Bio != nil {
		artist.Bio = req.Bio
	}

	if req.ImageURL != nil {
		artist.ImageURL = req.ImageURL
	}

	if req.IsVerified != nil {
		artist.IsVerified = *req.IsVerified
	}

	if req.MonthlyListeners != nil {
		artist.MonthlyListeners = *req.MonthlyListeners
	}

	now := time.Now().UTC()
	artist.UpdatedAt = &now

	if err := s.repo.UpdateArtist(ctx, artist); err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *Service) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteArtist(ctx, id)
}

// --------------------
// Albums
// --------------------

func (s *Service) CreateAlbum(ctx context.Context, req CreateAlbumRequest) (*Album, error) {
	if _, err := s.repo.GetArtistByID(ctx, req.ArtistID); err != nil {
		return nil, err
	}

	slug, err := s.uniqueAlbumSlug(ctx, req.ArtistID, req.Title, nil)
	if err != nil {
		return nil, err
	}

	var releaseDate *time.Time
	if req.ReleaseDate != nil && *req.ReleaseDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err != nil {
			return nil, err
		}
		releaseDate = &parsed
	}
	now := time.Now().UTC()
	albumType := "album"
	if req.AlbumType != nil && *req.AlbumType != "" {
		albumType = *req.AlbumType
	}

	album := &Album{
		ID:          uuid.New(),
		ArtistID:    req.ArtistID,
		Title:       req.Title,
		Slug:        slug,
		CoverURL:    req.CoverURL,
		ReleaseDate: releaseDate,
		AlbumType:   albumType,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	if err := s.repo.CreateAlbum(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}

func (s *Service) GetAlbum(ctx context.Context, id uuid.UUID) (*Album, error) {
	return s.repo.GetAlbumByID(ctx, id)
}

func (s *Service) ListAlbums(ctx context.Context, filter AlbumListFilter) ([]Album, int, error) {
	albums, err := s.repo.ListAlbums(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountAlbums(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return albums, total, nil
}

func (s *Service) UpdateAlbum(ctx context.Context, id uuid.UUID, req UpdateAlbumRequest) (*Album, error) {
	album, err := s.repo.GetAlbumByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		slug, err := s.uniqueAlbumSlug(ctx, album.ArtistID, *req.Title, &id)
		if err != nil {
			return nil, err
		}
		album.Title = *req.Title
		album.Slug = slug
	}

	if req.CoverURL != nil {
		album.CoverURL = req.CoverURL
	}

	if req.ReleaseDate != nil {
		if *req.ReleaseDate == "" {
			album.ReleaseDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.ReleaseDate)
			if err != nil {
				return nil, err
			}
			album.ReleaseDate = &parsed
		}
	}

	if req.AlbumType != nil {
		album.AlbumType = *req.AlbumType
	}

	now := time.Now().UTC()
	album.UpdatedAt = &now

	if err := s.repo.UpdateAlbum(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}

func (s *Service) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAlbum(ctx, id)
}

// --------------------
// Tracks
// --------------------

func (s *Service) CreateTrack(ctx context.Context, req CreateTrackRequest) (*Track, error) {
	if _, err := s.repo.GetArtistByID(ctx, req.ArtistID); err != nil {
		return nil, err
	}

	if req.AlbumID != nil {
		if _, err := s.repo.GetAlbumByID(ctx, *req.AlbumID); err != nil {
			return nil, err
		}
	}

	if err := s.validateGenreIDs(ctx, req.GenreIDs); err != nil {
		return nil, err
	}

	slug, err := s.uniqueTrackSlug(ctx, req.ArtistID, req.Title, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	explicit := false
	if req.Explicit != nil {
		explicit = *req.Explicit
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	track := &Track{
		ID:              uuid.New(),
		ArtistID:        req.ArtistID,
		AlbumID:         req.AlbumID,
		Title:           req.Title,
		Slug:            slug,
		DurationSeconds: req.DurationSeconds,
		TrackNumber:     req.TrackNumber,
		Explicit:        explicit,
		AudioURL:        req.AudioURL,
		CoverURL:        req.CoverURL,
		PlayCount:       0,
		IsPublic:        isPublic,
		CreatedAt:       now,
		UpdatedAt:       &now,
	}
	log.Print(track)
	if err := s.repo.CreateTrack(ctx, track, req.GenreIDs); err != nil {
		return nil, err
	}

	genresByTrack, err := s.repo.GetGenresByTrackIDs(ctx, []uuid.UUID{track.ID})
	if err == nil {
		track.Genres = genresByTrack[track.ID]
	}

	return track, nil
}

func (s *Service) GetTrack(ctx context.Context, id uuid.UUID, includePrivate bool) (*Track, error) {
	return s.repo.GetTrackByID(ctx, id, includePrivate)
}

func (s *Service) ListTracks(ctx context.Context, filter TrackListFilter) ([]Track, int, error) {
	tracks, err := s.repo.ListTracks(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountTracks(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	ids := make([]uuid.UUID, 0, len(tracks))
	for _, track := range tracks {
		ids = append(ids, track.ID)
	}

	genresByTrack, err := s.repo.GetGenresByTrackIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	for i := range tracks {
		tracks[i].Genres = genresByTrack[tracks[i].ID]
	}

	return tracks, total, nil
}

func (s *Service) UpdateTrack(ctx context.Context, id uuid.UUID, req UpdateTrackRequest) (*Track, error) {
	track, err := s.repo.GetTrackByID(ctx, id, true)
	if err != nil {
		return nil, err
	}

	if req.AlbumID != nil {
		if _, err := s.repo.GetAlbumByID(ctx, *req.AlbumID); err != nil {
			return nil, err
		}
		track.AlbumID = req.AlbumID
	}

	if req.Title != nil {
		slug, err := s.uniqueTrackSlug(ctx, track.ArtistID, *req.Title, &id)
		if err != nil {
			return nil, err
		}
		track.Title = *req.Title
		track.Slug = slug
	}

	if req.DurationSeconds != nil {
		track.DurationSeconds = *req.DurationSeconds
	}

	if req.TrackNumber != nil {
		track.TrackNumber = req.TrackNumber
	}

	if req.Explicit != nil {
		track.Explicit = *req.Explicit
	}

	if req.AudioURL != nil {
		track.AudioURL = req.AudioURL
	}

	if req.CoverURL != nil {
		track.CoverURL = req.CoverURL
	}

	if req.IsPublic != nil {
		track.IsPublic = *req.IsPublic
	}

	if err := s.validateGenreIDs(ctx, req.GenreIDs); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	track.UpdatedAt = &now

	if err := s.repo.UpdateTrack(ctx, track, req.GenreIDs); err != nil {
		return nil, err
	}

	genresByTrack, err := s.repo.GetGenresByTrackIDs(ctx, []uuid.UUID{track.ID})
	if err == nil {
		track.Genres = genresByTrack[track.ID]
	}

	return track, nil
}

func (s *Service) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTrack(ctx, id)
}

// --------------------
// Genres
// --------------------

func (s *Service) CreateGenre(ctx context.Context, req CreateGenreRequest) (*Genre, error) {
	slug, err := s.uniqueGenreSlug(ctx, req.Name, nil)
	if err != nil {
		return nil, err
	}

	genre := &Genre{
		ID:        uuid.New(),
		Name:      req.Name,
		Slug:      slug,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateGenre(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *Service) GetGenre(ctx context.Context, id uuid.UUID) (*Genre, error) {
	return s.repo.GetGenreByID(ctx, id)
}

func (s *Service) ListGenres(ctx context.Context) ([]Genre, error) {
	return s.repo.ListGenres(ctx)
}

func (s *Service) UpdateGenre(ctx context.Context, id uuid.UUID, req UpdateGenreRequest) (*Genre, error) {
	genre, err := s.repo.GetGenreByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		slug, err := s.uniqueGenreSlug(ctx, *req.Name, &id)
		if err != nil {
			return nil, err
		}

		genre.Name = *req.Name
		genre.Slug = slug
	}

	if err := s.repo.UpdateGenre(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

func (s *Service) DeleteGenre(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteGenre(ctx, id)
}

func (s *Service) validateGenreIDs(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	unique := make(map[uuid.UUID]struct{}, len(ids))
	for _, id := range ids {
		unique[id] = struct{}{}
	}

	deduped := make([]uuid.UUID, 0, len(unique))
	for id := range unique {
		deduped = append(deduped, id)
	}

	count, err := s.repo.CountGenresByIDs(ctx, deduped)
	if err != nil {
		return err
	}

	if count != len(deduped) {
		return ErrInvalidGenreIDs
	}

	return nil
}

// --------------------
// Search
// --------------------

func (s *Service) Search(ctx context.Context, q string) (*SearchResponse, error) {
	if q == "" {
		return &SearchResponse{
			Artists: []Artist{},
			Albums:  []Album{},
			Tracks:  []Track{},
		}, nil
	}

	return s.repo.Search(ctx, q, 10)
}
