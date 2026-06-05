import { z } from 'zod'
import { IdSchema } from '../common'

export const GenreSchema = z.object({
  id: IdSchema,
  name: z.string(),
  slug: z.string().optional().nullable(),
  track_count: z.number().optional().nullable(),
}).transform((genre) => ({
  id: genre.id,
  name: genre.name,
  slug: genre.slug ?? null,
  track_count: genre.track_count ?? 0,
}))

export type Genre = z.infer<typeof GenreSchema>

export interface GenreCreatePayload {
  name: string
}

export interface GenreUpdatePayload {
  name?: string
}
