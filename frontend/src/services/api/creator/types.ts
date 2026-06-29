import { z } from 'zod'

export const CreatorStatsSchema = z.object({
  user_id: z.string().catch(''),
  total_plays: z.number().default(0),
  unique_listeners: z.number().default(0),
  total_followers: z.number().default(0),
  total_tracks: z.number().default(0),
  total_albums: z.number().default(0),
  total_playlists: z.number().default(0),
  estimated_revenue: z.number().default(0),
  last_calculated: z.string().optional().nullable().default(null),
})

export const CreatorDailyStatSchema = z.object({
  id: z.string().catch(''),
  user_id: z.string().catch(''),
  date: z.string().catch(''),
  plays: z.number().default(0),
  listeners: z.number().default(0),
  likes: z.number().default(0),
  follows: z.number().default(0),
  shares: z.number().default(0),
  revenue_cents: z.number().default(0),
})

export const TrackStatsSchema = z.object({
  track_id: z.string().catch(''),
  title: z.string().catch(''),
  total_plays: z.number().default(0),
  total_likes: z.number().default(0),
  duration: z.number().default(0),
  created_at: z.string().optional().nullable().default(null),
})

export const OverviewResponseSchema = z.object({
  success: z.boolean(),
  data: CreatorStatsSchema.nullable().optional(),
  error: z.string().nullable().optional(),
})

export const DailyStatsResponseSchema = z.object({
  success: z.boolean(),
  data: z.array(CreatorDailyStatSchema).optional(),
  error: z.string().nullable().optional(),
})

export const TrackStatsResponseSchema = z.object({
  success: z.boolean(),
  data: z.array(TrackStatsSchema).optional(),
  error: z.string().nullable().optional(),
})

// --- Zod schemas for remaining creator data types ---

export const EarningsBreakdownSchema = z.object({
  total_revenue: z.number(),
  stream_revenue: z.number(),
  tip_revenue: z.number(),
  subscription_revenue: z.number(),
  pending_payout: z.number(),
  last_payout: z.number(),
  last_payout_date: z.string().optional(),
})

export const PayoutSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  amount: z.number(),
  method: z.string(),
  status: z.string(),
  created_at: z.string(),
  paid_at: z.string().optional(),
})

export const PayoutMethodSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  type: z.string(),
  details: z.string(),
  is_active: z.boolean(),
})

export const TopListenerSchema = z.object({
  user_id: z.string(),
  username: z.string(),
  avatar_url: z.string().optional(),
  play_count: z.number(),
})

export const GeographicStatSchema = z.object({
  country: z.string(),
  city: z.string(),
  listeners: z.number(),
  plays: z.number(),
})

export const AudienceOverviewSchema = z.object({
  total_listeners: z.number(),
  new_listeners_7d: z.number(),
  repeat_rate: z.number(),
})

export const AudienceDataSchema = z.object({
  overview: AudienceOverviewSchema,
  top_listeners: z.array(TopListenerSchema),
  geographics: z.array(GeographicStatSchema),
})

export const AlbumStatsSchema = z.object({
  id: z.string(),
  title: z.string(),
  release_year: z.number(),
  track_count: z.number(),
  total_plays: z.number(),
  cover_url: z.string().optional(),
})

export const CreatorPlaylistItemSchema = z.object({
  id: z.string(),
  name: z.string(),
  track_count: z.number(),
  is_public: z.boolean(),
  cover_url: z.string().optional(),
})

export const CreatorContentDataSchema = z.object({
  tracks: z.array(TrackStatsSchema),
  albums: z.array(AlbumStatsSchema),
  playlists: z.array(CreatorPlaylistItemSchema),
})

// --- TypeScript interfaces (inferred from Zod schemas where possible) ---

export type EarningsBreakdown = z.infer<typeof EarningsBreakdownSchema>
export type Payout = z.infer<typeof PayoutSchema>
export type PayoutMethod = z.infer<typeof PayoutMethodSchema>
export type TopListener = z.infer<typeof TopListenerSchema>
export type GeographicStat = z.infer<typeof GeographicStatSchema>
export type AudienceOverview = z.infer<typeof AudienceOverviewSchema>
export type AudienceData = z.infer<typeof AudienceDataSchema>
export type AlbumStats = z.infer<typeof AlbumStatsSchema>
export type CreatorPlaylistItem = z.infer<typeof CreatorPlaylistItemSchema>
export type CreatorContentData = z.infer<typeof CreatorContentDataSchema>

export interface TrackUpdateRequest {
  title?: string
  persian_title?: string
  genre_ids?: string[]
  explicit?: boolean
  track_number?: number
  lyrics?: string
}

export interface AlbumUpdateRequest {
  title?: string
  persian_title?: string
  description?: string
  album_type?: string
  release_date?: string
  genre_ids?: string[]
}

export type CreatorStats = z.infer<typeof CreatorStatsSchema>
export type CreatorDailyStat = z.infer<typeof CreatorDailyStatSchema>
export type TrackStats = z.infer<typeof TrackStatsSchema>
export type OverviewResponse = z.infer<typeof OverviewResponseSchema>
export type DailyStatsResponse = z.infer<typeof DailyStatsResponseSchema>
export type TrackStatsResponse = z.infer<typeof TrackStatsResponseSchema>
