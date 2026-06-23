import { z } from 'zod'
import { IdSchema } from '../common'

const AlbumArtistSchema = z.object({
  artistId: IdSchema,
  name: z.string(),
  slug: z.string().optional(),
  role: z.string().optional(),
  position: z.number().optional(),
})

export const AlbumSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    coverUrl: z.string().optional().nullable(),
    artistId: IdSchema.optional().nullable(),
    releaseDate: z.string().optional().nullable(),
    genre: z.string().optional().nullable(),
    artists: z.array(AlbumArtistSchema).optional().nullable(),
  })
  .transform((album) => ({
    id: album.id,
    title: album.title,
    cover_url: album.coverUrl ?? null,
    artist_id: album.artistId ?? null,
    artist_name: album.artists?.find((a) => a.role === 'primary')?.name ?? album.artists?.[0]?.name ?? null,
    release_date: album.releaseDate ?? null,
    genre: album.genre ?? null,
    track_count: 0,
  }))

export type Album = {
  id: string | number
  title: string
  cover_url: string | null
  artist_id: string | number | null
  artist_name: string | null
  release_date: string | null
  genre: string | null
  track_count: number
}

export interface AlbumCreatePayload {
  title: string
  coverUrl?: string | null
  artistId?: string | number | null
}

export interface AlbumUpdatePayload {
  title?: string
  coverUrl?: string | null
  artistId?: string | number | null
}
