import { z } from 'zod'

export const NotificationResponseSchema = z.object({
  id: z.string(),
  type: z.string(),
  title: z.string(),
  body: z.string(),
  entityType: z.string().nullable().optional(),
  entityId: z.string().nullable().optional(),
  payload: z.record(z.string(), z.any()).optional(),
  isRead: z.boolean(),
  readAt: z.string().nullable().optional(),
  createdAt: z.string(),
})

export const ListNotificationsResponseSchema = z.object({
  notifications: z.array(NotificationResponseSchema),
  unreadCount: z.number(),
})

export const MarkReadResponseSchema = z.object({
  message: z.string(),
})

export type NotificationResponse = z.infer<typeof NotificationResponseSchema>
export type ListNotificationsResponse = z.infer<typeof ListNotificationsResponseSchema>
export type MarkReadResponse = z.infer<typeof MarkReadResponseSchema>
