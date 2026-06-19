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
	if err != nil {
		tracks, _ = s.repo.SearchTracks(ctx, query, limit)
	}
	resp.Tracks = tracks
	if resp.Tracks == nil {
		resp.Tracks = []TrackResult{}
	}

	albums, err := s.searchAlbums(ctx, query, limit)
	if err != nil {
		albums, _ = s.repo.SearchAlbums(ctx, query, limit)
	}
	resp.Albums = albums
	if resp.Albums == nil {
		resp.Albums = []AlbumResult{}
	}

	artists, err := s.searchArtists(ctx, query, limit)
	if err != nil {
		artists, _ = s.repo.SearchArtists(ctx, query, limit)
	}
	resp.Artists = artists
	if resp.Artists == nil {
		resp.Artists = []ArtistResult{}
	}

	playlists, err := s.searchPlaylists(ctx, query, limit)
	if err != nil {
		playlists, _ = s.repo.SearchPlaylists(ctx, query, limit)
	}
	resp.Playlists = playlists
	if resp.Playlists == nil {
		resp.Playlists = []PlaylistResult{}
	}

	return resp, nil
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
				Source TrackResult `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	results := make([]TrackResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		results = append(results, hit.Source)
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
				Source AlbumResult `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	results := make([]AlbumResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		results = append(results, hit.Source)
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
				Source ArtistResult `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	results := make([]ArtistResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		results = append(results, hit.Source)
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
				Source PlaylistResult `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	results := make([]PlaylistResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		results = append(results, hit.Source)
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
