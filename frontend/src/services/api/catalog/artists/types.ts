import { z } from 'zod'
import { IdSchema } from '../common'

export const ArtistSchema = z
  .object({
    id: IdSchema,
    name: z.string(),
    bio: z.string().optional().nullable(),
    image_url: z.string().optional().nullable(),
    is_verified: z.boolean().optional().nullable(),
    monthly_listeners: z.number().optional().nullable(),
  })
  .transform((artist) => ({
    id: artist.id,
    name: artist.name,
    bio: artist.bio ?? null,
    image_url: artist.image_url ?? null,
    is_verified: artist.is_verified ?? false,
    monthly_listeners: artist.monthly_listeners ?? 0,
  }))

export type Artist = z.infer<typeof ArtistSchema>

export interface ArtistCreatePayload {
  name: string
  bio?: string | null
  image_url?: string | null
}

export interface ArtistUpdatePayload {
  name?: string
  bio?: string | null
  image_url?: string | null
}
