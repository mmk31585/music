package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	opensearch "github.com/opensearch-project/opensearch-go"
)

type Service struct {
	client *opensearch.Client
	repo   *Repository
}

func NewService(client *opensearch.Client, repo *Repository) *Service {
	return &Service{client: client, repo: repo}
}

func (s *Service) Search(ctx context.Context, query string, limit int) (*SearchResponse, error) {
	query = strings.TrimSpace(query)

	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	resp := &SearchResponse{
		Query:     query,
		Tracks:    []TrackResult{},
		Albums:    []AlbumResult{},
		Artists:   []ArtistResult{},
		Playlists: []PlaylistResult{},
	}

	if query == "" {
		return resp, nil
	}

	tracks, err := s.searchTracks(ctx, query, limit)
	if err != nil && s.repo != nil {
		tracks, _ = s.repo.SearchTracks(ctx, query, limit)
	}
	resp.Tracks = tracks
	if resp.Tracks == nil {
		resp.Tracks = []TrackResult{}
	}

	albums, err := s.searchAlbums(ctx, query, limit)
	if err != nil && s.repo != nil {
		albums, _ = s.repo.SearchAlbums(ctx, query, limit)
	}
	resp.Albums = albums
	if resp.Albums == nil {
		resp.Albums = []AlbumResult{}
	}

	artists, err := s.searchArtists(ctx, query, limit)
	if err != nil && s.repo != nil {
		artists, _ = s.repo.SearchArtists(ctx, query, limit)
	}
	resp.Artists = artists
	if resp.Artists == nil {
		resp.Artists = []ArtistResult{}
	}

	playlists, err := s.searchPlaylists(ctx, query, limit)
	if err != nil && s.repo != nil {
		playlists, _ = s.repo.SearchPlaylists(ctx, query, limit)
	}
	resp.Playlists = playlists
	if resp.Playlists == nil {
		resp.Playlists = []PlaylistResult{}
	}

	return resp, nil
}

type osTrackDoc struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	ArtistID  *string `json:"artist_id"`
	Artist    string  `json:"artist_name"`
	AlbumID   *string `json:"album_id"`
	Album     string  `json:"album_title"`
	CoverURL  *string `json:"cover_url"`
	AudioURL  *string `json:"audio_url"`
	Duration  int     `json:"duration"`
	Explicit  bool    `json:"is_explicit"`
	Year      int     `json:"year"`
	CreatedAt string  `json:"created_at"`
}

type osAlbumDoc struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist_name"`
	ArtistID string `json:"artist_id"`
	Year     int    `json:"release_year"`
	Type     string `json:"album_type"`
	CoverURL string `json:"cover_url"`
}

type osArtistDoc struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Bio              string `json:"bio"`
	Verified         bool   `json:"is_verified"`
	MonthlyListeners int    `json:"monthly_listeners"`
	ImageURL         string `json:"image_url"`
}

type osPlaylistDoc struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	UserID      string `json:"user_id"`
	IsPublic    bool   `json:"is_public"`
}

func (s *Service) searchTracks(ctx context.Context, q string, limit int) ([]TrackResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("opensearch not available")
	}
	body := map[string]any{
		"size": limit,
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     q,
				"fields":    []string{"title^3", "artist_name^2", "album_title"},
				"fuzziness": "AUTO",
			},
		},
	}

	respBody, err := s.doSearch(ctx, "tracks", body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Hits struct {
			Hits []struct {
				Source osTrackDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Hits.Hits) == 0 {
		return []TrackResult{}, nil
	}

	if s.repo != nil {
		ids := make([]string, 0, len(parsed.Hits.Hits))
		for _, hit := range parsed.Hits.Hits {
			ids = append(ids, hit.Source.ID)
		}
		return s.repo.SearchTracksByIDs(ctx, ids)
	}

	results := make([]TrackResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		src := hit.Source
		r := TrackResult{
			ID:              src.ID,
			Title:           src.Title,
			DurationSeconds: src.Duration,
			Explicit:        src.Explicit,
			CreatedAt:       src.CreatedAt,
			Artists:         []TrackArtistResult{},
			Genres:          []GenreResult{},
		}
		if src.ArtistID != nil && *src.ArtistID != "" {
			r.ArtistID = src.ArtistID
		}
		if src.AudioURL != nil {
			r.AudioURL = src.AudioURL
		}
		results = append(results, r)
	}
	return results, nil
}

func (s *Service) searchAlbums(ctx context.Context, q string, limit int) ([]AlbumResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("opensearch not available")
	}
	body := map[string]any{
		"size": limit,
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     q,
				"fields":    []string{"title^3", "artist_name^2"},
				"fuzziness": "AUTO",
			},
		},
	}

	respBody, err := s.doSearch(ctx, "albums", body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Hits struct {
			Hits []struct {
				Source osAlbumDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Hits.Hits) == 0 {
		return []AlbumResult{}, nil
	}

	if s.repo != nil {
		ids := make([]string, 0, len(parsed.Hits.Hits))
		for _, hit := range parsed.Hits.Hits {
			ids = append(ids, hit.Source.ID)
		}
		return s.repo.SearchAlbumsByIDs(ctx, ids)
	}

	results := make([]AlbumResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		src := hit.Source
		results = append(results, AlbumResult{
			ID:    src.ID,
			Title: src.Title,
			Artists: []TrackArtistResult{{
				ArtistID: src.ArtistID,
				Name:     src.Artist,
				Role:     "primary",
				Position: 1,
			}},
		})
	}
	return results, nil
}

func (s *Service) searchArtists(ctx context.Context, q string, limit int) ([]ArtistResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("opensearch not available")
	}
	body := map[string]any{
		"size": limit,
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     q,
				"fields":    []string{"name^3"},
				"fuzziness": "AUTO",
			},
		},
	}

	respBody, err := s.doSearch(ctx, "artists", body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Hits struct {
			Hits []struct {
				Source osArtistDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Hits.Hits) == 0 {
		return []ArtistResult{}, nil
	}

	if s.repo != nil {
		ids := make([]string, 0, len(parsed.Hits.Hits))
		for _, hit := range parsed.Hits.Hits {
			ids = append(ids, hit.Source.ID)
		}
		return s.repo.SearchArtistsByIDs(ctx, ids)
	}

	results := make([]ArtistResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		src := hit.Source
		results = append(results, ArtistResult{
			ID:               src.ID,
			Name:             src.Name,
			Bio:              strPtr(src.Bio),
			CoverURL:         strPtr(src.ImageURL),
			IsVerified:       src.Verified,
			MonthlyListeners: src.MonthlyListeners,
		})
	}
	return results, nil
}

func (s *Service) searchPlaylists(ctx context.Context, q string, limit int) ([]PlaylistResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("opensearch not available")
	}
	body := map[string]any{
		"size": limit,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query":     q,
							"fields":    []string{"name^3", "description"},
							"fuzziness": "AUTO",
						},
					},
				},
				"filter": []any{
					map[string]any{
						"term": map[string]any{
							"is_public": true,
						},
					},
				},
			},
		},
	}

	respBody, err := s.doSearch(ctx, "playlists", body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Hits struct {
			Hits []struct {
				Source osPlaylistDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	results := make([]PlaylistResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		src := hit.Source
		results = append(results, PlaylistResult{
			ID:          src.ID,
			Name:        src.Name,
			Description: strPtr(src.Description),
			CoverURL:    strPtr(src.CoverURL),
			UserID:      strPtr(src.UserID),
			IsPublic:    boolPtr(src.IsPublic),
		})
	}

	return results, nil
}

func (s *Service) doSearch(ctx context.Context, index string, body map[string]any) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(index),
		s.client.Search.WithBody(bytes.NewReader(raw)),
		s.client.Search.WithTrackTotalHits(false),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("opensearch search error on index %s: %s", index, string(respBody))
	}

	return respBody, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
