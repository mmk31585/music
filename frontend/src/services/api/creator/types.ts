import { z } from 'zod'

export const CreatorStatsSchema = z.object({
  user_id: z.string(),
  total_plays: z.number(),
  unique_listeners: z.number(),
  total_followers: z.number(),
  total_tracks: z.number(),
  total_albums: z.number(),
  total_playlists: z.number(),
  estimated_revenue: z.number(),
  last_calculated: z.string(),
})

export const CreatorDailyStatSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  date: z.string(),
  plays: z.number(),
  listeners: z.number(),
  likes: z.number(),
  follows: z.number(),
  shares: z.number(),
  revenue_cents: z.number(),
})

export const TrackStatsSchema = z.object({
  track_id: z.string(),
  title: z.string(),
  total_plays: z.number(),
  total_likes: z.number(),
  duration: z.number(),
  created_at: z.string(),
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

// --- New types (no Zod) ---

export interface EarningsBreakdown {
  total_revenue: number
  stream_revenue: number
  tip_revenue: number
  subscription_revenue: number
  pending_payout: number
  last_payout: number
  last_payout_date?: string
}

export interface Payout {
  id: string
  user_id: string
  amount: number
  method: string
  status: string
  created_at: string
  paid_at?: string
}

export interface PayoutMethod {
  id: string
  user_id: string
  type: string
  details: string
  is_active: boolean
}

export interface TopListener {
  user_id: string
  username: string
  avatar_url?: string
  play_count: number
}

export interface GeographicStat {
  country: string
  city: string
  listeners: number
  plays: number
}

export interface AudienceOverview {
  total_listeners: number
  new_listeners_7d: number
  repeat_rate: number
}

export interface AudienceData {
  overview: AudienceOverview
  top_listeners: TopListener[]
  geographics: GeographicStat[]
}

export interface AlbumStats {
  id: string
  title: string
  release_year: number
  track_count: number
  total_plays: number
  cover_url?: string
}

export interface CreatorPlaylistItem {
  id: string
  name: string
  track_count: number
  is_public: boolean
  cover_url?: string
}

export interface CreatorContentData {
  tracks: TrackStats[]
  albums: AlbumStats[]
  playlists: CreatorPlaylistItem[]
}

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
