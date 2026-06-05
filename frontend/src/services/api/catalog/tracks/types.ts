import { z } from 'zod'
import { IdSchema } from '../common'
import { GenreSchema } from '../genres'

export const TrackArtistRequestSchema = z.object({
  artist_id: IdSchema,
  role: z.string(),
  position: z.number().int().nonnegative().optional(),
})

export const TrackSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    duration_seconds: z.number().optional().nullable(),
    audio_url: z.string().optional().nullable(),
    cover_url: z.string().optional().nullable(),
    artist_id: IdSchema.optional().nullable(),
    album_id: IdSchema.optional().nullable(),
    genre_id: IdSchema.optional().nullable(),
    artist_name: z.string().optional().nullable(),
    album_title: z.string().optional().nullable(),
    genres: z.array(GenreSchema).optional().nullable(),
    play_count: z.number().optional().nullable(),
    track_number: z.number().optional().nullable(),
    explicit: z.boolean().optional().nullable(),
  })
  .transform((track) => ({
    id: track.id,
    title: track.title,
    duration_seconds: track.duration_seconds ?? 0,
    audio_url: track.audio_url ?? null,
    cover_url: track.cover_url ?? null,
    artist_id: track.artist_id ?? null,
    album_id: track.album_id ?? null,
    genre_id: track.genre_id ?? track.genres?.[0]?.id ?? null,
    artist_name: track.artist_name ?? null,
    album_title: track.album_title ?? null,
    genres: track.genres ?? [],
    play_count: track.play_count ?? 0,
    track_number: track.track_number ?? null,
    explicit: track.explicit ?? false,
  }))

export type Track = z.infer<typeof TrackSchema>
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
