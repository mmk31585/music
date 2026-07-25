package history

import (
	"testing"
)

func TestClassifySignal_ShortPlayIsSkip(t *testing.T) {
	result := ClassifySignal(15_000, 200_000, false)
	if result != SignalTypeSkipNegative {
		t.Errorf("expected skip_negative, got %s", result)
	}
}

func TestClassifySignal_SkipAtExactThreshold(t *testing.T) {
	// 29,999ms is still below the 30s threshold → skip
	result := ClassifySignal(29_999, 200_000, false)
	if result != SignalTypeSkipNegative {
		t.Errorf("expected skip_negative, got %s", result)
	}
}

func TestClassifySignal_HighCompletionIsPositive(t *testing.T) {
	result := ClassifySignal(160_000, 200_000, false)
	if result != SignalTypeCompletePositive {
		t.Errorf("expected complete_positive, got %s", result)
	}
}

func TestClassifySignal_ExactCompletionThreshold(t *testing.T) {
	// 80% exactly → complete_positive
	result := ClassifySignal(160_000, 200_000, false)
	if result != SignalTypeCompletePositive {
		t.Errorf("expected complete_positive, got %s", result)
	}
}

func TestClassifySignal_ReplayOverridesOtherSignals(t *testing.T) {
	// Even a short skip should be classified as replay_strong if it's a replay
	result := ClassifySignal(5_000, 200_000, true)
	if result != SignalTypeReplayStrong {
		t.Errorf("expected replay_strong, got %s", result)
	}
}

func TestClassifySignal_PartialPlayIsNeutral(t *testing.T) {
	// 60s played out of 200s → 30% completion → not skip (over 30s), not complete → neutral
	result := ClassifySignal(60_000, 200_000, false)
	if result != SignalTypePartialNeutral {
		t.Errorf("expected partial_neutral, got %s", result)
	}
}

func TestClassifySignal_ZeroDurationTrackNoCrash(t *testing.T) {
	result := ClassifySignal(10_000, 0, false)
	if result != SignalTypePartialNeutral {
		t.Errorf("expected partial_neutral for zero-duration track, got %s", result)
	}
}
