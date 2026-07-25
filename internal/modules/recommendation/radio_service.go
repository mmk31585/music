package recommendation

import (
	"context"
	"fmt"
)

const (
	defaultBatchSize    = 10
	genreBoost          = 0.05
	candidateMultiplier = 3
)

type RadioService struct {
	repo     RadioRepository
	mlClient *SimilarityMLClient
}

func NewRadioService(repo RadioRepository, mlClient *SimilarityMLClient) *RadioService {
	return &RadioService{repo: repo, mlClient: mlClient}
}

func (s *RadioService) StartRadio(ctx context.Context, userID, seedTrackID string, seedType string) (*RadioSession, []RadioTrackItem, error) {
	// End any existing active session for this user
	existing, _ := s.repo.GetActiveSession(ctx, userID)
	if existing != nil {
		s.repo.UpdateLastActive(ctx, existing.ID)
	}

	session := &RadioSession{
		UserID:         userID,
		SeedTrackID:    seedTrackID,
		SeedType:       seedType,
		PlayedTrackIDs: []string{},
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("create session: %w", err)
	}

	if err := s.repo.AppendPlayedTrack(ctx, session.ID, seedTrackID); err != nil {
		return nil, nil, fmt.Errorf("append seed: %w", err)
	}

	tracks, err := s.getBatch(ctx, session, defaultBatchSize)
	if err != nil {
		return nil, nil, err
	}

	return session, tracks, nil
}

func (s *RadioService) GetNextBatch(ctx context.Context, sessionID string, count int) ([]RadioTrackItem, error) {
	if count <= 0 {
		count = defaultBatchSize
	}

	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	tracks, err := s.getBatch(ctx, session, count)
	if err != nil {
		return nil, err
	}

	s.repo.UpdateLastActive(ctx, sessionID)
	return tracks, nil
}

func (s *RadioService) EndRadio(ctx context.Context, sessionID string) error {
	return s.repo.UpdateLastActive(ctx, sessionID)
}

// getBatch implements the fallback chain:
// 1. Get similar tracks via ML service (cosine similarity)
// 2. If not enough results, relax genre boost
// 3. If still not enough, fall back to popular tracks in the seed's genre
func (s *RadioService) getBatch(ctx context.Context, session *RadioSession, count int) ([]RadioTrackItem, error) {
	// Step 1: Try ML similarity
	var resultIDs []string

	if s.mlClient != nil {
		similar, err := s.mlClient.GetSimilarTracks(ctx, session.SeedTrackID, count*candidateMultiplier, session.PlayedTrackIDs)
		if err == nil && len(similar) > 0 {
			for _, t := range similar {
				resultIDs = append(resultIDs, t.TrackID)
				if len(resultIDs) >= count {
					break
				}
			}
		}
	}

	// Step 2: If still not enough, widen request to ML with exclude only (no genre boost concern)
	if len(resultIDs) < count && s.mlClient != nil {
		more, err := s.mlClient.GetSimilarTracks(ctx, session.SeedTrackID, count*candidateMultiplier*2, session.PlayedTrackIDs)
		if err == nil && len(more) > 0 {
			seen := make(map[string]struct{}, len(resultIDs))
			for _, id := range resultIDs {
				seen[id] = struct{}{}
			}
			for _, t := range more {
				if _, ok := seen[t.TrackID]; !ok {
					resultIDs = append(resultIDs, t.TrackID)
					seen[t.TrackID] = struct{}{}
					if len(resultIDs) >= count {
						break
					}
				}
			}
		}
	}

	// Step 3: Fallback to popular tracks
	if len(resultIDs) < count {
		popular, err := s.repo.GetPopularTracks(ctx, "", count)
		if err == nil && len(popular) > 0 {
			seen := make(map[string]struct{}, len(resultIDs))
			for _, id := range resultIDs {
				seen[id] = struct{}{}
			}
			for _, t := range popular {
				if _, ok := seen[t.ID]; !ok {
					resultIDs = append(resultIDs, t.ID)
					seen[t.ID] = struct{}{}
					if len(resultIDs) >= count {
						break
					}
				}
			}
		}
	}

	if len(resultIDs) == 0 {
		return []RadioTrackItem{}, nil
	}

	tracks, err := s.repo.GetTracksByIDs(ctx, resultIDs)
	if err != nil {
		return nil, fmt.Errorf("get track details: %w", err)
	}

	return tracks, nil
}
