package player

import "context"

func (s *Service) DispatchPlayStarted(
	userID string,
	track *PlaybackTrack,
	duration int,
	completed bool,
	source string,
) {
	go s.TrackPlayed(
		context.Background(),
		userID,
		track,
		duration,
		completed,
		source,
	)
}
