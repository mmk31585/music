package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const progressPrefix = "import:job:"

type ProgressStore struct {
	rdb redis.UniversalClient
}

func NewProgressStore(rdb redis.UniversalClient) *ProgressStore {
	return &ProgressStore{rdb: rdb}
}

func (s *ProgressStore) key(jobID string) string {
	return progressPrefix + jobID
}

func (s *ProgressStore) Set(ctx context.Context, jobID string, status Status, progress int, stage string, draftID string, errMsg string) error {
	data, _ := json.Marshal(ProgressUpdate{
		JobID:    jobID,
		Status:   status,
		Progress: progress,
		Stage:    stage,
		DraftID:  draftID,
		Error:    errMsg,
	})
	return s.rdb.Set(ctx, s.key(jobID), data, 24*time.Hour).Err()
}

func (s *ProgressStore) Get(ctx context.Context, jobID string) (*ProgressUpdate, error) {
	data, err := s.rdb.Get(ctx, s.key(jobID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var pu ProgressUpdate
	if err := json.Unmarshal(data, &pu); err != nil {
		return nil, fmt.Errorf("parse progress: %w", err)
	}
	return &pu, nil
}

func (s *ProgressStore) Delete(ctx context.Context, jobID string) error {
	return s.rdb.Del(ctx, s.key(jobID)).Err()
}
