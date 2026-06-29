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
  duration_seconds: z.number().default(0),
  audio_url: z.string().optional().nullable().default(null),
  cover_url: z.string().optional().nullable().default(null),
  artist_id: IdSchema.optional().nullable().default(null),
  album_id: IdSchema.optional().nullable().default(null),
  genre_id: IdSchema.optional().nullable().default(null),
  artist_name: z.string().optional().nullable().default(null),
  album_title: z.string().optional().nullable().default(null),
  genres: z.array(GenreSchema).optional().default([]),
  play_count: z.number().default(0),
  track_number: z.number().optional().nullable().default(null),
  explicit: z.boolean().default(false),
})

export type Track = {
  id: string | number
  title: string
  duration_seconds: number
  audio_url: string | null
  cover_url: string | null
  artist_id: string | number | null
  album_id: string | number | null
  genre_id: string | number | null
  artist_name: string | null
  album_title: string | null
  genres: any[]
  play_count: number
  track_number: number | null
  explicit: boolean
}
export type TrackArtistRequest = z.infer<typeof TrackArtistRequestSchema>

export interface TrackCreatePayload {
  title: string
  artist_id: string | number
  artists?: TrackArtistRequest[]
  album_id?: string | number | null
  duration_seconds?: number | null
  audio_url?: string | null
  cover_url?: string | null
  genre_ids?: Array<string | number>
  track_number?: number | null
  explicit?: boolean | null
  is_public?: boolean
}

export interface TrackUpdatePayload {
  title?: string
  artist_id?: string | number | null
  artists?: TrackArtistRequest[]
  album_id?: string | number | null
  duration_seconds?: number | null
  audio_url?: string | null
  cover_url?: string | null
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

export function toTrackMutationPayload(payload: Partial<Track>): TrackCreatePayload {
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
    album_id: payload.album_id ?? null,
    duration_seconds: payload.duration_seconds ?? 0,
    audio_url: payload.audio_url ?? null,
    cover_url: payload.cover_url ?? null,
    genre_ids: payload.genre_id ? [payload.genre_id] : [],
    track_number: payload.track_number ?? null,
    explicit: payload.explicit ?? false,
  }
}
