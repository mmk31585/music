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
    artists: z.array(AlbumArtistSchema).optional().nullable(),
  })
  .transform((album) => ({
    id: album.id,
    title: album.title,
    cover_url: album.coverUrl ?? null,
    artist_id: album.artistId ?? null,
    artist_name: album.artists?.find((a) => a.role === 'primary')?.name ?? album.artists?.[0]?.name ?? null,
    release_date: album.releaseDate ?? null,
    track_count: 0,
  }))

export type Album = z.infer<typeof AlbumSchema>

export interface AlbumCreatePayload {
  title: string
  cover_url?: string | null
  artist_id?: string | number | null
}

export interface AlbumUpdatePayload {
  title?: string
  cover_url?: string | null
  artist_id?: string | number | null
}
