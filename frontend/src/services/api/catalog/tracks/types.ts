import { z } from 'zod'
import { IdSchema } from '../common'
import { GenreSchema } from '../genres'

export const TrackArtistRequestSchema = z.object({
  artist_id: IdSchema,
  role: z.string(),
  position: z.number().int().nonnegative().optional(),
})

const TrackArtistSchema = z.object({
  artist_id: IdSchema.optional().nullable().default(null),
  name: z.string(),
  slug: z.string().optional(),
  role: z.string().optional(),
  position: z.number().optional(),
})

export const TrackSchema = z.object({
  id: IdSchema,
  title: z.string(),
  slug: z.string().optional(),
  duration_seconds: z.number().default(0),
  audio_url: z.string().optional().nullable().default(null),
  cover_url: z.string().optional().nullable().default(null),
  audio_media_id: z.string().optional().nullable().default(null),
  cover_media_id: z.string().optional().nullable().default(null),
  artist_id: IdSchema.optional().nullable().default(null),
  artist_name: z.string().optional().nullable().default(null),
  album_id: IdSchema.optional().nullable().default(null),
  artists: z.array(TrackArtistSchema).optional().default([]),
  genres: z.array(GenreSchema).optional().default([]),
  play_count: z.number().default(0),
  track_number: z.number().optional().nullable().default(null),
  explicit: z.boolean().default(false),
  is_public: z.boolean().optional().default(true),
  created_at: z.string().optional().nullable().default(null),
  updated_at: z.string().optional().nullable().default(null),
})

export type Track = {
  id: string | number
  title: string
  slug?: string
  duration_seconds: number
  audio_url: string | null
  cover_url: string | null
  audio_media_id: string | null
  cover_media_id: string | null
  artist_id: string | number | null
  artist_name: string | null
  album_id: string | number | null
  album_title?: string | null
  artists?: Array<{ artist_id: string | number | null; name: string; slug?: string; role?: string; position?: number }>
  genres: any[]
  play_count: number
  track_number: number | null
  explicit: boolean
  is_public: boolean
  created_at: string | null
  updated_at: string | null
}
export type TrackArtistRequest = z.infer<typeof TrackArtistRequestSchema>

export interface TrackCreatePayload {
  title: string
  artist_id: string | number
  artists?: TrackArtistRequest[]
  artist_ids?: Array<string | number>
  credits?: TrackArtistRequest[]
  album_id?: string | number | null
  duration_seconds?: number | null
  audio_url?: string | null
  cover_url?: string | null
  audio_media_id?: string | null
  cover_media_id?: string | null
  genre_ids?: Array<string | number>
  track_number?: number | null
  explicit?: boolean | null
  is_public?: boolean
}

export interface TrackUpdatePayload {
  title?: string | null
  artist_id?: string | number | null
  artists?: TrackArtistRequest[]
  album_id?: string | number | null
  duration_seconds?: number | null
  audio_url?: string | null
  cover_url?: string | null
  audio_media_id?: string | null
  cover_media_id?: string | null
  genre_ids?: Array<string | number>
  track_number?: number | null
  explicit?: boolean | null
  is_public?: boolean | null
  /** Names of nullable fields to set to NULL (backend uses these to distinguish 'clear' from 'absent') */
  clear_fields?: string[]
}

export interface TrackUploadPayload extends TrackCreatePayload {
  audioFile: File
}

export function toTrackMutationPayload(payload: Partial<Track> & { genre_ids?: Array<string | number>; artist_ids?: Array<string | number> }): TrackCreatePayload {
  if (!payload.title) {
    throw new Error('title is required')
  }

  if (!payload.artist_id) {
    throw new Error('artist_id is required')
  }

  return {
    title: payload.title,
    artist_id: payload.artist_id,
    artists: payload.artist_id
      ? [
        {
          artist_id: payload.artist_id,
          role: 'primary',
          position: 0,
        },
      ]
      : [],
    artist_ids: payload.artist_ids ?? [],
    album_id: payload.album_id ?? null,
    duration_seconds: payload.duration_seconds ?? 0,
    audio_url: payload.audio_url ?? null,
    cover_url: payload.cover_url ?? null,
    audio_media_id: payload.audio_media_id ?? null,
    cover_media_id: payload.cover_media_id ?? null,
    genre_ids: payload.genre_ids ?? [],
    track_number: payload.track_number ?? null,
    explicit: payload.explicit ?? false,
  }
}
