import { z } from 'zod'

export const RecommendationTrackSchema = z.object({
  id: z.string().catch(''),
  title: z.string().catch(''),
  artist_id: z.string().nullable().optional().default(null),
  artist_name: z.string().nullable().optional().default(null),
  album_id: z.string().nullable().optional().default(null),
  album_title: z.string().nullable().optional().default(null),
  genre: z.string().nullable().optional().default(null),
  cover_url: z.string().nullable().optional().default(null),
  audio_url: z.string().nullable().optional().default(null),
  duration_seconds: z.number().nullable().optional().default(0),
  score: z.number().nullable().optional().default(0),
})

export type RecommendationTrack = z.infer<typeof RecommendationTrackSchema>

export const RecommendationResponseSchema = z.object({
  type: z.string().optional().nullable().default(null),
  items: z.array(RecommendationTrackSchema).catch([]),
  limit: z.number().default(0),
})

export type RecommendationResponse = z.infer<typeof RecommendationResponseSchema>

export const HomeFeedSectionSchema = z.object({
  id: z.string().catch(''),
  title: z.string().catch(''),
  subtitle: z.string().optional().nullable().default(null),
  type: z.string().optional().nullable().default(null),
  items: z.array(RecommendationTrackSchema).catch([]),
  seed_track: RecommendationTrackSchema.nullable().optional().default(null),
})

export type HomeFeedSection = z.infer<typeof HomeFeedSectionSchema>

export const HomeFeedResponseSchema = z.object({
  sections: z.array(HomeFeedSectionSchema).catch([]),
})

export type HomeFeedResponse = z.infer<typeof HomeFeedResponseSchema>

export const DiscoverWeeklyPlaylistMetaSchema = z.object({
  generated_at: z.string().optional().nullable().default(null),
  week_of: z.string().optional().nullable().default(null),
  track_count: z.number().default(0),
})

export type DiscoverWeeklyPlaylistMeta = z.infer<typeof DiscoverWeeklyPlaylistMetaSchema>

export const DiscoverWeeklyResponseSchema = z.object({
  playlist: DiscoverWeeklyPlaylistMetaSchema.optional().nullable().default({ generated_at: null, week_of: null, track_count: 0 }),
  tracks: z.array(RecommendationTrackSchema).catch([]),
})

export type DiscoverWeeklyResponse = z.infer<typeof DiscoverWeeklyResponseSchema>

export const TrackStatSchema = z.object({
  track_id: z.string().catch(''),
  title: z.string().catch(''),
  artist_name: z.string().catch(''),
  play_count: z.number().default(0),
})

export const ArtistStatSchema = z.object({
  artist_id: z.string().catch(''),
  artist_name: z.string().catch(''),
  play_count: z.number().default(0),
})

export const GenreStatSchema = z.object({
  genre_name: z.string().catch(''),
  play_count: z.number().default(0),
})

export const ListeningStatsSchema = z.object({
  period_label: z.string().catch(''),
  total_minutes_listened: z.number().default(0),
  total_tracks_played: z.number().default(0),
  top_tracks: z.array(TrackStatSchema).catch([]),
  top_artists: z.array(ArtistStatSchema).catch([]),
  top_genres: z.array(GenreStatSchema).catch([]),
  unique_artists_count: z.number().default(0),
  longest_listening_streak_days: z.number().default(0),
  discovery_score: z.number().default(0),
})

export type ListeningStats = z.infer<typeof ListeningStatsSchema>
