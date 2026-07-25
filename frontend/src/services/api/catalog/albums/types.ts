import { z } from 'zod'
import { IdSchema } from '../common'

const AlbumArtistSchema = z.object({
  artist_id: IdSchema.optional().nullable().default(null),
  name: z.string(),
  slug: z.string().optional(),
  role: z.string().optional(),
  position: z.number().optional(),
})

export const AlbumSchema = z.object({
  id: IdSchema,
  title: z.string(),
  slug: z.string().optional(),
  cover_url: z.string().optional().nullable().default(null),
  cover_media_id: z.string().optional().nullable().default(null),
  artist_id: IdSchema.optional().nullable().default(null),
  artist_name: z.string().optional().nullable().default(null),
  release_date: z.string().optional().nullable().default(null),
  album_type: z.string().optional().nullable().default(null),
  artists: z.array(AlbumArtistSchema).optional().nullable().default([]),
  track_count: z.number().default(0),
  created_at: z.string().optional().nullable().default(null),
  updated_at: z.string().optional().nullable().default(null),
})

export type Album = {
  id: string | number
  title: string
  slug?: string
  cover_url: string | null
  cover_media_id: string | null
  artist_id: string | number | null
  artist_name: string | null
  release_date: string | null
  album_type: string | null
  genre?: string | null
  track_count: number
  created_at: string | null
  updated_at: string | null
}

export interface AlbumArtistRequest {
  artist_id: string | number
  role?: string
  position?: number
}

export interface AlbumCreatePayload {
  title: string
  artist_id?: string | number
  artists?: AlbumArtistRequest[]
  cover_url?: string | null
  cover_media_id?: string | null
  release_date?: string | null
  album_type?: string | null
}

export interface AlbumUpdatePayload {
  title?: string | null
  artist_id?: string | number | null
  artists?: AlbumArtistRequest[]
  cover_url?: string | null
  cover_media_id?: string | null
  release_date?: string | null
  album_type?: string | null
}
