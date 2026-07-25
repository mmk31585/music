import { z } from 'zod'

export const PlaylistTrackItemSchema = z.object({
  playlist_track_id: z.string(),
  track_id: z.string(),
  position: z.number(),
  title: z.string(),
  artist_name: z.string().nullable().optional(),
  album_title: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  audio_url: z.string().nullable().optional(),
  duration_seconds: z.number().nullable().optional(),
})

export type PlaylistTrackItem = z.infer<typeof PlaylistTrackItemSchema>

export const PlaylistListItemSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  name: z.string(),
  description: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  is_public: z.boolean(),
  is_collaborative: z.boolean().optional().default(false),
  track_count: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type PlaylistListItem = z.infer<typeof PlaylistListItemSchema>

export const CollaboratorResponseSchema = z.object({
  user_id: z.string(),
  added_at: z.string(),
  is_creator: z.boolean(),
})

export const PlaylistDetailSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  name: z.string(),
  description: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  is_public: z.boolean(),
  is_collaborative: z.boolean().optional().default(false),
  created_at: z.string(),
  updated_at: z.string(),
  tracks: z.array(PlaylistTrackItemSchema).optional(),
  collaborators: z.array(CollaboratorResponseSchema).optional(),
})

export type CollaboratorResponse = z.infer<typeof CollaboratorResponseSchema>

export type PlaylistDetail = z.infer<typeof PlaylistDetailSchema>

export interface CreatePlaylistPayload {
  name: string
  description?: string | null
  cover_url?: string | null
  is_public?: boolean
}

export interface AddTrackPayload {
  track_id: string
}
