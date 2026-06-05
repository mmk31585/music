import { z } from 'zod'

export const LyricsSchema = z.object({
  id: z.union([z.string(), z.number()]),
  track_id: z.union([z.string(), z.number()]).optional().nullable(),
  content: z.string().optional().nullable(),
  language: z.string().optional().nullable(),
  type: z.string().optional().nullable(),
  created_at: z.string().optional().nullable(),
  updated_at: z.string().optional().nullable(),
})

export type Lyrics = z.infer<typeof LyricsSchema>

export type CreateLyricsPayload = {
  track_id: string | number
  content: string
  language?: string | null
  type?: string | null
}

export type UpdateLyricsPayload = {
  track_id?: string | number | null
  content?: string | null
  language?: string | null
  type?: string | null
}
