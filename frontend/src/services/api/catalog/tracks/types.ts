import { z } from 'zod'
import { IdSchema } from '../common'
import { GenreSchema } from '../genres'

export const TrackArtistRequestSchema = z.object({
  artistId: IdSchema,
  role: z.string(),
  position: z.number().int().nonnegative().optional(),
})

const TrackArtistSchema = z.object({
  artistId: IdSchema,
  name: z.string(),
  slug: z.string().optional(),
  role: z.string().optional(),
  position: z.number().optional(),
})

export const TrackSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    durationSeconds: z.number().optional().nullable(),
    audioUrl: z.string().optional().nullable(),
    coverUrl: z.string().optional().nullable(),
    artistId: IdSchema.optional().nullable(),
    albumId: IdSchema.optional().nullable(),
    genres: z.array(GenreSchema).optional().nullable(),
    playCount: z.number().optional().nullable(),
    trackNumber: z.number().optional().nullable(),
    explicit: z.boolean().optional().nullable(),
    artists: z.array(TrackArtistSchema).optional().nullable(),
  })
  .transform((track) => ({
    id: track.id,
    title: track.title,
    duration_seconds: track.durationSeconds ?? 0,
    audio_url: track.audioUrl ?? null,
    cover_url: track.coverUrl ?? null,
    artist_id: track.artistId ?? null,
    album_id: track.albumId ?? null,
    genre_id: track.genres?.[0]?.id ?? null,
    artist_name: track.artists?.find((a) => a.role === 'primary')?.name ?? track.artists?.[0]?.name ?? null,
    album_title: null,
    genres: track.genres ?? [],
    play_count: track.playCount ?? 0,
    track_number: track.trackNumber ?? null,
    explicit: track.explicit ?? false,
  }))

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
  artistId: string | number
  artists?: TrackArtistRequest[]
  albumId?: string | number | null
  durationSeconds?: number | null
  audioUrl?: string | null
  coverUrl?: string | null
  genreIds?: Array<string | number>
  trackNumber?: number | null
  explicit?: boolean | null
  isPublic?: boolean
}

export interface TrackUpdatePayload {
  title?: string
  artistId?: string | number | null
  artists?: TrackArtistRequest[]
  albumId?: string | number | null
  durationSeconds?: number | null
  audioUrl?: string | null
  coverUrl?: string | null
  genreIds?: Array<string | number>
  trackNumber?: number | null
  explicit?: boolean | null
  isPublic?: boolean | null
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
    artistId: payload.artist_id,
    artists: payload.artist_id
      ? [
        {
          artistId: payload.artist_id,
          role: 'primary',
          position: 0,
        },
      ]
      : [],
    albumId: payload.album_id ?? null,
    durationSeconds: payload.duration_seconds ?? 0,
    audioUrl: payload.audio_url ?? null,
    coverUrl: payload.cover_url ?? null,
    genreIds: payload.genre_id ? [payload.genre_id] : [],
    trackNumber: payload.track_number ?? null,
    explicit: payload.explicit ?? false,
  }
}
