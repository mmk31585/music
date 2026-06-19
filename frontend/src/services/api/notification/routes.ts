import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { NotificationApiRoutes } from './enums'

export const NotificationType = {
  Welcome: 'welcome',
  SubscriptionPurchased: 'subscription_purchased',
  PlaylistCreated: 'playlist_created',
  ArtistPublishedTrack: 'artist_published_track',
  PlaylistShared: 'playlist_shared',
  SubscriptionRenewed: 'subscription_renewed',
  UserFollowed: 'user_followed',
} as const

export type NotificationType = (typeof NotificationType)[keyof typeof NotificationType]

export interface NotificationResponse {
  id: string
  type: string
  title: string
  body: string
  entityType?: string | null
  entityId?: string | null
  payload?: Record<string, any>
  isRead: boolean
  readAt?: string | null
  createdAt: string
}

export interface ListNotificationsResponse {
  notifications: NotificationResponse[]
  unreadCount: number
}

export interface MarkReadResponse {
  message: string
}

export const useNotificationApi = () => {
  const getNotifications = async (
    limit = 20,
    offset = 0,
    config?: UseRequestConfig<ListNotificationsResponse>,
  ) => {
    return useRequest<ListNotificationsResponse>(
      NotificationApiRoutes.LIST,
      { method: 'GET', params: { limit, offset } },
      { silent: true, ...config },
    )
  }

  const markRead = async (id: string, config?: UseRequestConfig<MarkReadResponse>) => {
    return useRequest<MarkReadResponse>(
      NotificationApiRoutes.MARK_READ.replace(':id', id),
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const markAllRead = async (config?: UseRequestConfig<MarkReadResponse>) => {
    return useRequest<MarkReadResponse>(
      NotificationApiRoutes.MARK_ALL_READ,
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  return { getNotifications, markRead, markAllRead }
}

// Legacy standalone functions — prefer useNotificationApi
export async function getNotificationsLegacy(limit = 20, offset = 0): Promise<ListNotificationsResponse> {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
  const res = await fetch(`/api/v1/notifications?${params}`, { credentials: 'include' })
  if (!res.ok) throw new Error('Failed to fetch notifications')
  return res.json()
}

export async function markReadLegacy(id: string): Promise<MarkReadResponse> {
  const res = await fetch(`/api/v1/notifications/${id}/read`, {
    method: 'POST',
    credentials: 'include',
  })
  if (!res.ok) throw new Error('Failed to mark notification as read')
  return res.json()
}

export async function markAllReadLegacy(): Promise<MarkReadResponse> {
  const res = await fetch('/api/v1/notifications/read-all', {
    method: 'POST',
    credentials: 'include',
  })
  if (!res.ok) throw new Error('Failed to mark all as read')
  return res.json()
}
