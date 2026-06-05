import { z } from 'zod'
import { IdSchema } from '../common'

export const AlbumSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    cover_url: z.string().optional().nullable(),
    artist_id: IdSchema.optional().nullable(),
    artist_name: z.string().optional().nullable(),
    release_date: z.string().optional().nullable(),
    track_count: z.number().optional().nullable(),
  })
  .transform((album) => ({
    id: album.id,
    title: album.title,
    cover_url: album.cover_url ?? null,
    artist_id: album.artist_id ?? null,
    artist_name: album.artist_name ?? null,
    release_date: album.release_date ?? null,
    track_count: album.track_count ?? 0,
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
