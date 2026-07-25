import { z } from 'zod'
import { IdSchema } from '../common'

export const ArtistSchema = z.object({
  id: IdSchema,
  name: z.string(),
  slug: z.string().optional(),
  bio: z.string().optional().nullable(),
  image_url: z.string().optional().nullable(),
  avatar_media_id: z.string().optional().nullable().default(null),
  banner_media_id: z.string().optional().nullable().default(null),
  country: z.string().optional().nullable().default(null),
  is_verified: z.boolean().default(false),
  monthly_listeners: z.number().default(0),
  created_at: z.string().optional().nullable().default(null),
  updated_at: z.string().optional().nullable().default(null),
})

export type Artist = z.infer<typeof ArtistSchema>

export interface ArtistCreatePayload {
  name: string
  bio?: string | null
  image_url?: string | null
  avatar_media_id?: string | null
  banner_media_id?: string | null
  country?: string | null
  is_verified?: boolean | null
  monthly_listeners?: number | null
}

export interface ArtistUpdatePayload {
  name?: string | null
  bio?: string | null
  image_url?: string | null
  avatar_media_id?: string | null
  banner_media_id?: string | null
  country?: string | null
  is_verified?: boolean | null
  monthly_listeners?: number | null
}
