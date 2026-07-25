import { z } from 'zod'

export const QueueItemSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  track_id: z.string(),
  position: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type QueueItem = z.infer<typeof QueueItemSchema>

export const QueueResponseSchema = z.object({
  items: z.array(QueueItemSchema),
  count: z.number(),
})

export type QueueResponse = z.infer<typeof QueueResponseSchema>

export const AddTrackResponseSchema = z.object({
  item: QueueItemSchema,
})

export type AddTrackResponse = z.infer<typeof AddTrackResponseSchema>

export const AddTrackPayloadSchema = z.object({
  track_id: z.string(),
  mode: z.enum(['next', 'later']).optional().default('later'),
})

export type AddTrackPayload = z.infer<typeof AddTrackPayloadSchema>

export const ReorderQueueItemSchema = z.object({
  id: z.string(),
  position: z.number().min(1),
})

export const ReorderQueuePayloadSchema = z.object({
  items: z.array(ReorderQueueItemSchema).min(1),
})

export type ReorderQueuePayload = z.infer<typeof ReorderQueuePayloadSchema>
