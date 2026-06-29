import { z } from 'zod'

export const HistoryItemSchema = z.object({
  id: z.string().catch(''),
  user_id: z.string().catch(''),
  track_id: z.string().catch(''),
  played_at: z.string().catch(''),
  duration: z.number().default(0),
  completed: z.boolean().default(false),
  track_title: z.string().nullable().optional().default(null),
  track_duration: z.number().nullable().optional().default(0),
  track_cover_url: z.string().nullable().optional().default(null),
  artist_name: z.string().nullable().optional().default(null),
})

export type HistoryItem = z.infer<typeof HistoryItemSchema>

export const HistoryResponseSchema = z.object({
  items: z.array(HistoryItemSchema).catch([]),
  pagination: z.object({
    limit: z.number().default(0),
    offset: z.number().default(0),
    count: z.number().default(0),
    has_more: z.boolean().default(false),
  }).optional().nullable().default({ limit: 0, offset: 0, count: 0, has_more: false }),
})

export type HistoryResponse = z.infer<typeof HistoryResponseSchema>
