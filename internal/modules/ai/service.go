package ai

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo    *Repository
	ai      AIClient
	logger  *zap.Logger
	enabled bool
}

func NewService(repo *Repository, ai AIClient, logger *zap.Logger, enabled bool) *Service {
	return &Service{repo: repo, ai: ai, logger: logger, enabled: enabled}
}

func (s *Service) GenerateEmbedding(ctx context.Context, trackID string) (*TrackEmbedding, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	uid, err := uuid.Parse(trackID)
	if err != nil {
		return nil, fmt.Errorf("invalid track ID: %w", err)
	}

	meta, err := s.repo.GetTrackMetadata(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("get track metadata: %w", err)
	}

	input := fmt.Sprintf("%s %s %s", meta.Title, meta.Artist, meta.Genre)
	vec, err := s.ai.GenerateEmbedding(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("generate embedding: %w", err)
	}

	if err := s.repo.UpsertEmbedding(ctx, uid, vec, "v1"); err != nil {
		return nil, fmt.Errorf("save embedding: %w", err)
	}

	return &TrackEmbedding{
		TrackID:      uid,
		Embedding:    vec,
		ModelVersion: "v1",
		UpdatedAt:    time.Now(),
	}, nil
}

func (s *Service) AnalyzeMood(ctx context.Context, trackID string) (*TrackMood, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	uid, err := uuid.Parse(trackID)
	if err != nil {
		return nil, fmt.Errorf("invalid track ID: %w", err)
	}

	meta, err := s.repo.GetTrackMetadata(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("get track metadata: %w", err)
	}

	mood, err := s.ai.AnalyzeMood(ctx, meta.Title, meta.Artist, meta.Genre)
	if err != nil {
		return nil, fmt.Errorf("analyze mood: %w", err)
	}
	mood.TrackID = uid

	if err := s.repo.UpsertMood(ctx, *mood); err != nil {
		return nil, fmt.Errorf("save mood: %w", err)
	}

	return mood, nil
}

func (s *Service) GetMood(ctx context.Context, trackID string) (*TrackMood, error) {
	uid, err := uuid.Parse(trackID)
	if err != nil {
		return nil, fmt.Errorf("invalid track ID: %w", err)
	}

	mood, err := s.repo.GetMood(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("mood not found: %w", err)
	}
	return mood, nil
}

func (s *Service) GetSimilarByMood(ctx context.Context, trackID string, mood string, limit int) ([]TrackMeta, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	if mood != "" {
		return s.repo.GetTracksByMood(ctx, mood, limit)
	}

	uid, err := uuid.Parse(trackID)
	if err != nil {
		return nil, fmt.Errorf("invalid track ID: %w", err)
	}

	sourceMood, err := s.repo.GetMood(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("source track mood not found: %w", err)
	}

	energyRange := 0.2
	valenceRange := 0.2

	return s.repo.GetTracksByMoodRange(ctx,
		math.Max(0, sourceMood.Energy-energyRange),
		math.Min(1, sourceMood.Energy+energyRange),
		math.Max(0, sourceMood.Valence-valenceRange),
		math.Min(1, sourceMood.Valence+valenceRange),
		limit+1,
	)
}

func (s *Service) GetSimilarByEmbedding(ctx context.Context, trackID string, limit int, embeddingSpace string) ([]TrackMeta, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if embeddingSpace == "" {
		embeddingSpace = "audio"
	}

	uid, err := uuid.Parse(trackID)
	if err != nil {
		return nil, fmt.Errorf("invalid track ID: %w", err)
	}

	emb, err := s.repo.GetEmbedding(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("embedding not found: %w", err)
	}

	tracks, err := s.repo.GetSimilarByEmbedding(ctx, []float64(emb.Embedding), limit+1, embeddingSpace)
	if err != nil {
		return nil, err
	}

	filtered := make([]TrackMeta, 0, limit)
	for _, t := range tracks {
		if t.ID != trackID {
			filtered = append(filtered, t)
			if len(filtered) >= limit {
				break
			}
		}
	}
	return filtered, nil
}

func (s *Service) GeneratePlaylist(ctx context.Context, req GeneratePlaylistRequest, userID uuid.UUID) (*AIPlaylistResponse, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI features are disabled")
	}

	start := time.Now()

	prompt := req.Prompt
	if req.Mood != "" {
		prompt = fmt.Sprintf("%s (mood: %s)", prompt, req.Mood)
	}
	if req.Activity != "" {
		prompt = fmt.Sprintf("%s (activity: %s)", prompt, req.Activity)
	}
	if req.Genre != "" {
		prompt = fmt.Sprintf("%s (genre: %s)", prompt, req.Genre)
	}

	var candidates []TrackMeta

	if req.SeedTrackID != "" {
		uid, err := uuid.Parse(req.SeedTrackID)
		if err == nil {
			if emb, err := s.repo.GetEmbedding(ctx, uid); err == nil {
				candidates, err = s.repo.GetSimilarByEmbedding(ctx, []float64(emb.Embedding), 50, "audio")
				if err != nil {
					s.logger.Warn("embedding similarity failed, falling back to mood", zap.Error(err))
				}
			}
		}
	}

	if len(candidates) == 0 && req.Mood != "" {
		candidates = s.getCandidatesForMood(ctx, req.Mood, 50)
	}

	if len(candidates) == 0 {
		var err error
		candidates, err = s.repo.GetTracksWithMood(ctx, 50)
		if err != nil {
			return nil, fmt.Errorf("get candidate tracks: %w", err)
		}
	}

	// If mood is specified but candidates weren't pre-filtered (e.g. AI fetched all tracks),
	// apply heuristic in-memory filtering before sending to AI for selection
	if req.Mood != "" {
		filtered := filterTracksByMood(candidates, req.Mood)
		if len(filtered) < len(candidates) {
			s.logger.Info("mood playlist: applied in-memory heuristic filter",
				zap.String("mood", req.Mood),
				zap.Int("before", len(candidates)),
				zap.Int("after", len(filtered)))
			candidates = filtered
		}
	}

	selectedIDs, err := s.ai.GeneratePlaylist(ctx, prompt, candidates)
	if err != nil {
		s.logger.Warn("AI playlist generation failed, using random selection", zap.Error(err))
		selectedIDs = s.fallbackSelect(candidates, 20)
	}

	tracks, err := s.repo.GetTracksByIDs(ctx, selectedIDs)
	if err != nil {
		return nil, fmt.Errorf("get selected tracks: %w", err)
	}

	if len(tracks) > 20 {
		tracks = tracks[:20]
	}

	items := make([]TrackItem, len(tracks))
	for i, t := range tracks {
		items[i] = TrackItem{
			ID:       t.ID,
			Title:    t.Title,
			Artist:   t.Artist,
			Duration: t.Duration,
			CoverURL: t.CoverURL,
			Energy:   t.Energy,
			Valence:  t.Valence,
		}
	}

	latency := int(time.Since(start).Milliseconds())

	playlistID := uuid.New()
	log := GenerationLog{
		UserID:     userID,
		PlaylistID: &playlistID,
		Prompt:     redactPII(prompt),
		TrackCount: len(items),
		ModelUsed:  "ai-v1",
		LatencyMs:  latency,
	}
	_ = s.repo.LogGeneration(ctx, log)

	name := s.generatePlaylistName(req)

	return &AIPlaylistResponse{
		ID:          fmt.Sprintf("ai_%d", time.Now().Unix()),
		Name:        name,
		Description: fmt.Sprintf("Generated from: %s", req.Prompt),
		Tracks:      items,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// moodGenres maps mood names to relevant genre names (case-insensitive).
// Used as a proxy when no tracks have AI mood analysis data, so we can still
// filter by genre affinity instead of returning all tracks.
// var moodGenres = map[string][]string{
// 	"energetic":  {"Rock", "Metal", "Electronic", "Dance", "Funk", "Pop Rock"},
// 	"happy":      {"Pop", "Dance", "Funk", "Reggae", "Pop Rock", "Indie"},
// 	"chill":      {"Ambient", "Lo-Fi", "Jazz", "R&B", "Soul", "Folk"},
// 	"calm":       {"Classical", "Ambient", "New Age", "Instrumental"},
// 	"sad":        {"Blues", "Folk", "Soul", "Indie", "Traditional"},
// 	"focus":      {"Classical", "Ambient", "Lo-Fi", "Instrumental", "Jazz"},
// 	"romantic":   {"R&B", "Soul", "Jazz", "Pop", "Fusion"},
// 	"intense":    {"Metal", "Rock", "Electronic", "Fusion", "Punk"},
// 	"confident":  {"Hip Hop", "Rap", "Pop", "Rock", "Electronic", "Dance"},
// 	"sleep":      {"Ambient", "Classical", "New Age", "Instrumental", "Lo-Fi"},
// }

// moodToRange maps mood names to energy/valence ranges for numerical mood filtering.
// Used when tracks don't have explicit mood tags yet but have mood analysis data.
type moodRange struct {
	MinEnergy, MaxEnergy   float64
	MinValence, MaxValence float64
	MinTempo, MaxTempo     float64 // optional: further narrow by BPM
	MinDanceability        float64
}

var moodRanges = map[string]moodRange{
	"energetic": {0.65, 1.0, 0.4, 1.0, 120, 200, 0.3},
	"happy":     {0.4, 1.0, 0.6, 1.0, 100, 180, 0.4},
	"chill":     {0.0, 0.5, 0.3, 0.8, 60, 100, 0.0},
	"calm":      {0.0, 0.35, 0.3, 0.7, 50, 90, 0.0},
	"sad":       {0.0, 0.4, 0.0, 0.5, 50, 100, 0.0},
	"focus":     {0.3, 0.6, 0.2, 0.5, 80, 140, 0.0},
	"romantic":  {0.25, 0.55, 0.4, 0.8, 60, 120, 0.0},
	"intense":   {0.7, 1.0, 0.0, 0.5, 130, 200, 0.0},
	"confident": {0.5, 1.0, 0.5, 1.0, 90, 160, 0.4},
	"sleep":     {0.0, 0.25, 0.0, 0.4, 30, 80, 0.0},
}

// moodGenres maps mood names to relevant genre names (matching DB genre names).
// Used as a proxy when no tracks have AI mood analysis data (energy/valence = 0).
var moodGenres = map[string][]string{
	"energetic": {"Rock", "Metal", "Electronic", "Dance", "Funk", "Pop Rock"},
	"happy":     {"Pop", "Dance", "Funk", "Reggae", "Pop Rock", "Indie"},
	"chill":     {"Ambient", "Lo-Fi", "Jazz", "R&B", "Soul", "Folk"},
	"calm":      {"Classical", "Ambient", "New Age", "Instrumental"},
	"sad":       {"Blues", "Folk", "Soul", "Indie", "Traditional"},
	"focus":     {"Classical", "Ambient", "Lo-Fi", "Instrumental", "Jazz"},
	"romantic":  {"R&B", "Soul", "Jazz", "Pop", "Fusion"},
	"intense":   {"Metal", "Rock", "Electronic", "Fusion", "Punk"},
	"confident": {"Hip Hop", "Rap", "Pop", "Rock", "Electronic", "Dance"},
	"sleep":     {"Ambient", "Classical", "New Age", "Instrumental", "Lo-Fi"},
}

func (s *Service) CleanupOldLogs(ctx context.Context) error {
	_, err := s.repo.db.ExecContext(ctx, `DELETE FROM ai_generation_log WHERE created_at < NOW() - INTERVAL '90 days'`)
	return err
}

// getCandidatesForMood attempts to find tracks matching a mood using progressively
// broader strategies: exact mood tag → energy/valence range → genre proxy → all tracks.
func (s *Service) getCandidatesForMood(ctx context.Context, mood string, limit int) []TrackMeta {
	// Strategy 1: Exact mood tag match (requires AI mood analysis to have run)
	candidates, err := s.repo.GetTracksByMood(ctx, mood, limit)
	if err != nil {
		s.logger.Warn("mood tag query failed", zap.String("mood", mood), zap.Error(err))
	}
	if len(candidates) > 0 {
		s.logger.Info("mood playlist: found candidates via mood tags",
			zap.String("mood", mood), zap.Int("count", len(candidates)))
		return candidates
	}

	// Strategy 2: Energy/valence range based on mood mapping (requires mood analysis data)
	if r, ok := moodRanges[mood]; ok {
		candidates, err = s.repo.GetTracksByMoodRange(ctx, r.MinEnergy, r.MaxEnergy, r.MinValence, r.MaxValence, limit)
		if err != nil {
			s.logger.Warn("mood range query failed", zap.String("mood", mood), zap.Error(err))
		}
		if len(candidates) > 0 {
			s.logger.Info("mood playlist: found candidates via energy/valence range",
				zap.String("mood", mood), zap.Int("count", len(candidates)))
			return candidates
		}
	}

	// Strategy 3: Genre-based proxy — when no mood analysis data exists, narrow by genre affinity
	if genres, ok := moodGenres[mood]; ok && len(genres) > 0 {
		candidates, err = s.repo.GetTracksByGenres(ctx, genres, limit)
		if err != nil {
			s.logger.Warn("genre proxy query failed", zap.String("mood", mood), zap.Error(err))
		}
		if len(candidates) > 0 {
			s.logger.Info("mood playlist: found candidates via genre proxy",
				zap.String("mood", mood), zap.Strings("genres", genres), zap.Int("count", len(candidates)))
			return candidates
		}
	}

	// Strategy 4: No mood data at all — fetch broad set with mood data for in-memory filtering
	s.logger.Warn("mood playlist: no mood data or genre matches found for mood, fetching all tracks with mood info",
		zap.String("mood", mood))
	candidates, err = s.repo.GetTracksWithMood(ctx, limit*2)
	if err != nil {
		s.logger.Warn("fallback track query failed", zap.Error(err))
	}
	return candidates
}

// filterTracksByMood applies heuristic mood filtering on in-memory tracks.
// Uses energy/valence ranges first, then falls back to genre-based proxy
// when no tracks have mood analysis data.
func filterTracksByMood(tracks []TrackMeta, mood string) []TrackMeta {
	if len(tracks) == 0 {
		return tracks
	}

	r, ok := moodRanges[mood]
	if !ok {
		return tracks
	}

	// Check if any tracks have mood data
	hasMoodData := false
	for _, t := range tracks {
		if t.Energy != 0 || t.Valence != 0 {
			hasMoodData = true
			break
		}
	}

	if hasMoodData {
		// Filter by energy/valence ranges
		filtered := make([]TrackMeta, 0, len(tracks))
		for _, t := range tracks {
			if t.Energy == 0 && t.Valence == 0 {
				continue
			}
			if t.Energy >= r.MinEnergy && t.Energy <= r.MaxEnergy &&
				t.Valence >= r.MinValence && t.Valence <= r.MaxValence {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) > 0 {
			return filtered
		}
	}

	// No mood data or all filtered out — fall back to genre-based proxy
	if genres, ok := moodGenres[mood]; ok && len(genres) > 0 {
		genreSet := make(map[string]bool, len(genres))
		for _, g := range genres {
			genreSet[g] = true
		}
		filtered := make([]TrackMeta, 0, len(tracks))
		for _, t := range tracks {
			if t.Genre != "" && genreSet[t.Genre] {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) > 0 {
			return filtered
		}
	}

	// Nothing worked — return original
	return tracks
}

func redactPII(input string) string {
	if input == "" {
		return input
	}
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	input = re.ReplaceAllString(input, "[EMAIL REDACTED]")
	re = regexp.MustCompile(`(\+98|0)?9\d{9}`)
	input = re.ReplaceAllString(input, "[PHONE REDACTED]")
	re = regexp.MustCompile(`\b\d{10}\b`)
	input = re.ReplaceAllString(input, "[ID REDACTED]")
	return input
}

func (s *Service) fallbackSelect(tracks []TrackMeta, limit int) []string {
	if len(tracks) == 0 {
		return nil
	}
	indices := rand.Perm(len(tracks))
	count := limit
	if count > len(tracks) {
		count = len(tracks)
	}
	ids := make([]string, count)
	for i := 0; i < count; i++ {
		ids[i] = tracks[indices[i]].ID
	}
	return ids
}

func (s *Service) generatePlaylistName(req GeneratePlaylistRequest) string {
	parts := []string{}
	if req.Mood != "" {
		parts = append(parts, req.Mood)
	}
	if req.Activity != "" {
		parts = append(parts, req.Activity)
	}
	if req.Genre != "" {
		parts = append(parts, req.Genre)
	}
	if len(parts) > 0 {
		title := parts[0]
		if len(title) > 0 {
			runes := []rune(title)
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
			title = string(runes)
		}
		return fmt.Sprintf("%s Vibes", title)
	}
	if req.Prompt != "" {
		if len(req.Prompt) > 40 {
			return req.Prompt[:40] + "..."
		}
		return req.Prompt
	}
	return "AI-Generated Playlist"
}

func (s *Service) ProcessBatchEmbeddings(ctx context.Context, limit int) (int, error) {
	tracks, err := s.repo.GetTracksWithoutEmbeddings(ctx, limit)
	if err != nil {
		return 0, fmt.Errorf("get tracks without embeddings: %w", err)
	}

	processed := 0
	for _, t := range tracks {
		select {
		case <-ctx.Done():
			return processed, ctx.Err()
		default:
		}

		rateLimit := time.NewTicker(200 * time.Millisecond)
		<-rateLimit.C
		rateLimit.Stop()

		_, err := s.GenerateEmbedding(ctx, t.ID)
		if err != nil {
			s.logger.Warn("embedding generation failed", zap.String("track_id", t.ID), zap.Error(err))
			continue
		}
		processed++
	}

	return processed, nil
}

func (s *Service) ProcessBatchMoods(ctx context.Context, limit int) (int, error) {
	tracks, err := s.repo.GetTracksWithoutMoods(ctx, limit)
	if err != nil {
		return 0, fmt.Errorf("get tracks without moods: %w", err)
	}

	processed := 0
	for _, t := range tracks {
		select {
		case <-ctx.Done():
			return processed, ctx.Err()
		default:
		}

		rateLimit := time.NewTicker(500 * time.Millisecond)
		<-rateLimit.C
		rateLimit.Stop()

		_, err := s.AnalyzeMood(ctx, t.ID)
		if err != nil {
			s.logger.Warn("mood analysis failed", zap.String("track_id", t.ID), zap.Error(err))
			continue
		}
		processed++
	}

	return processed, nil
}
