package history

// ClassifySignal classifies a playback event into a signal type based on
// played duration, track duration, and replay-in-session detection.
//
// Thresholds (locked decisions — do not change without spec change):
//   - <30s played          → skip_negative
//   - ≥80% of total played  → complete_positive
//   - replay in session    → replay_strong (overrides all others)
//   - everything else      → partial_neutral
func ClassifySignal(playedDurationMs, trackDurationMs int64, isReplayInSession bool) string {
	if isReplayInSession {
		return SignalTypeReplayStrong
	}
	if trackDurationMs > 0 && playedDurationMs < SkipThresholdMs {
		return SignalTypeSkipNegative
	}
	if trackDurationMs > 0 {
		completion := float64(playedDurationMs) / float64(trackDurationMs)
		if completion >= CompletionThreshold {
			return SignalTypeCompletePositive
		}
	}
	return SignalTypePartialNeutral
}
