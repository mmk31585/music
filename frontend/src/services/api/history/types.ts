import { z } from 'zod'

export const HistoryItemSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  track_id: z.string(),
  played_at: z.string(),
  duration: z.number(),
  completed: z.boolean(),
  track_title: z.string().nullable().optional(),
  track_duration: z.number().nullable().optional(),
  track_cover_url: z.string().nullable().optional(),
  artist_name: z.string().nullable().optional(),
})

export type HistoryItem = z.infer<typeof HistoryItemSchema>

export const HistoryResponseSchema = z.object({
  items: z.array(HistoryItemSchema).catch([]),
  pagination: z.object({
    limit: z.number(),
    offset: z.number(),
    count: z.number(),
    has_more: z.boolean(),
  }),
})

export type HistoryResponse = z.infer<typeof HistoryResponseSchema>
