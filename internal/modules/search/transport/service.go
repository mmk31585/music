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
}

func NewService(client *opensearch.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Search(ctx context.Context, query string, limit int) (*SearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchResponse{
			Query:     query,
			Tracks:    []TrackResult{},
			Albums:    []AlbumResult{},
			Artists:   []ArtistResult{},
			Playlists: []PlaylistResult{},
		}, nil
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	tracks, err := s.searchTracks(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	albums, err := s.searchAlbums(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	artists, err := s.searchArtists(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	playlists, err := s.searchPlaylists(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Query:     query,
		Tracks:    tracks,
		Albums:    albums,
		Artists:   artists,
		Playlists: playlists,
	}, nil
}

func (s *Service) searchTracks(ctx context.Context, q string, limit int) ([]TrackResult, error) {
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

	respBody, err := s.doSearch(ctx, "artist", body)
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
