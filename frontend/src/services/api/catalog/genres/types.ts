import { z } from 'zod'
import { IdSchema } from '../common'

export const GenreSchema = z.object({
  id: IdSchema,
  name: z.string(),
  slug: z.string().optional().nullable().default(null),
  track_count: z.number().default(0),
})

export type Genre = z.infer<typeof GenreSchema>

export interface GenreCreatePayload {
  name: string
}

export interface GenreUpdatePayload {
  name?: string
}
