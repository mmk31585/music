package player

import "context"

func (s *Service) DispatchPlayStarted(ctx context.Context, trackID string) {
	go s.TrackPlayed(context.Background(), trackID)
}
