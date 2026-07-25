import { z } from 'zod'

const IdSchema = z.union([z.string(), z.number()])

// ── Track summary (embedded in VideoItem) ──────────────────────────────

export const TrackSummarySchema = z.object({
  id: IdSchema,
  title: z.string(),
  artist_name: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  duration_seconds: z.number().nullable().optional(),
})

export type TrackSummary = z.infer<typeof TrackSummarySchema>

// ── Video uploader ────────────────────────────────────────────────────

export const UploaderSchema = z.object({
  id: IdSchema,
  username: z.string(),
  avatar_url: z.string().nullable().optional(),
})

// ── Video item (from explore endpoint) ─────────────────────────────────

export const VideoItemSchema = z.object({
  id: IdSchema,
  type: z.enum(['official_mv', 'user_edit']),
  track_id: z.string().optional(),
  track: TrackSummarySchema.optional(),
  uploader: UploaderSchema.optional(),
  /** Backend may send uploader_id (string) instead of uploader object */
  uploader_id: z.union([z.string(), z.number()]).optional(),
  title: z.string(),
  description: z.string(),
  thumbnail_url: z.string().nullable().optional(),
  final_video_url: z.string().nullable().optional(),
  /** Backend may send thumbnail_path / final_video_path instead of _url */
  thumbnail_path: z.string().nullable().optional(),
  final_video_path: z.string().nullable().optional(),
  /** Track cover art URL — used as fallback poster when video has no thumbnail */
  track_cover_url: z.string().nullable().optional(),
  duration_ms: z.number(),
  aspect_ratio: z.string(),
  view_count: z.number(),
  like_count: z.number(),
  is_liked: z.boolean().optional().default(false),
  /** Backend may send is_public / is_approved booleans */
  is_public: z.boolean().optional(),
  is_approved: z.boolean().optional(),
  track_start_ms: z.number().optional(),
  track_end_ms: z.number().optional(),
  status: z.enum(['processing', 'ready', 'failed']).optional(),
  created_at: z.string(),
  updated_at: z.string().optional(),
})

export type VideoItem = z.infer<typeof VideoItemSchema>

// ── Explore response ───────────────────────────────────────────────────

export const ExploreResponseSchema = z.object({
  items: z.array(VideoItemSchema),
  limit: z.number(),
  offset: z.number(),
  count: z.number(),
})

export type ExploreResponse = z.infer<typeof ExploreResponseSchema>

// ── Music status response ──────────────────────────────────────────────

export const MusicStatusResponseSchema = z.object({
  playing: z.boolean(),
  current_track_id: z.string().nullable().optional(),
  updated_at: z.string().nullable().optional(),
})

export type MusicStatusResponse = z.infer<typeof MusicStatusResponseSchema>

// ── Upload edit payload ────────────────────────────────────────────────

export const CreateEditPayloadSchema = z.object({
  track_id: IdSchema,
  title: z.string().optional(),
  description: z.string().optional(),
  raw_video_url: z.string(),
  track_start_ms: z.number().optional(),
  track_end_ms: z.number().optional(),
  duration_ms: z.number().optional(),
  aspect_ratio: z.string().optional(),
  file_size_bytes: z.number().optional(),
})

export type CreateEditPayload = z.infer<typeof CreateEditPayloadSchema>

// ── Upload edit response ───────────────────────────────────────────────

export const UploadEditResponseSchema = z.object({
  id: IdSchema.optional(),
  video_id: IdSchema.optional(),
  job_id: z.string().optional(),
})

export type UploadEditResponse = z.infer<typeof UploadEditResponseSchema>

// ── Video job status response ──────────────────────────────────────────

export const VideoJobStatusSchema = z.object({
  id: z.string(),
  status: z.enum(['pending', 'processing', 'completed', 'failed', 'timeout']),
  video_id: z.string().optional(),
  error: z.string().optional(),
})

export type VideoJobStatus = z.infer<typeof VideoJobStatusSchema>

// ── Like visibility toggle ─────────────────────────────────────────────

export const LikeVisibilityPayloadSchema = z.object({
  visibility: z.enum(['public', 'private']),
})

export type LikeVisibilityPayload = z.infer<typeof LikeVisibilityPayloadSchema>

// ── Liked track item (from public endpoint) ────────────────────────────

export const LikedTrackItemSchema = z.object({
  id: IdSchema,
  track_id: IdSchema,
  track_title: z.string(),
  artist_name: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  duration_seconds: z.number().nullable().optional(),
  is_public: z.boolean().optional(),
  liked_at: z.string(),
})

export type LikedTrackItem = z.infer<typeof LikedTrackItemSchema>
