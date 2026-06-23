import { z } from 'zod'

export const RecommendationTrackSchema = z.object({
  id: z.string(),
  title: z.string(),
  artist_id: z.string().nullable().optional(),
  artist_name: z.string().nullable().optional(),
  album_id: z.string().nullable().optional(),
  album_title: z.string().nullable().optional(),
  genre: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  audio_url: z.string().nullable().optional(),
  duration_seconds: z.number().nullable().optional(),
  score: z.number().nullable().optional(),
})

export type RecommendationTrack = z.infer<typeof RecommendationTrackSchema>

export const RecommendationResponseSchema = z.object({
  type: z.string(),
  items: z.array(RecommendationTrackSchema).catch([]),
  limit: z.number(),
})

export type RecommendationResponse = z.infer<typeof RecommendationResponseSchema>

export const HomeFeedSectionSchema = z.object({
  id: z.string(),
  title: z.string(),
  subtitle: z.string().optional(),
  type: z.string(),
  items: z.array(RecommendationTrackSchema).catch([]),
  seed_track: RecommendationTrackSchema.nullable().optional(),
})

export type HomeFeedSection = z.infer<typeof HomeFeedSectionSchema>

export const HomeFeedResponseSchema = z.object({
  sections: z.array(HomeFeedSectionSchema).catch([]),
})

export type HomeFeedResponse = z.infer<typeof HomeFeedResponseSchema>

export const DiscoverWeeklyPlaylistMetaSchema = z.object({
  generated_at: z.string(),
  week_of: z.string(),
  track_count: z.number(),
})

export type DiscoverWeeklyPlaylistMeta = z.infer<typeof DiscoverWeeklyPlaylistMetaSchema>

export const DiscoverWeeklyResponseSchema = z.object({
  playlist: DiscoverWeeklyPlaylistMetaSchema,
  tracks: z.array(RecommendationTrackSchema).catch([]),
})

export type DiscoverWeeklyResponse = z.infer<typeof DiscoverWeeklyResponseSchema>

export const TrackStatSchema = z.object({
  track_id: z.string(),
  title: z.string(),
  artist_name: z.string(),
  play_count: z.number(),
})

export const ArtistStatSchema = z.object({
  artist_id: z.string(),
  artist_name: z.string(),
  play_count: z.number(),
})

export const GenreStatSchema = z.object({
  genre_name: z.string(),
  play_count: z.number(),
})

export const ListeningStatsSchema = z.object({
  period_label: z.string(),
  total_minutes_listened: z.number(),
  total_tracks_played: z.number(),
  top_tracks: z.array(TrackStatSchema).catch([]),
  top_artists: z.array(ArtistStatSchema).catch([]),
  top_genres: z.array(GenreStatSchema).catch([]),
  unique_artists_count: z.number(),
  longest_listening_streak_days: z.number(),
  discovery_score: z.number(),
})

export type ListeningStats = z.infer<typeof ListeningStatsSchema>
