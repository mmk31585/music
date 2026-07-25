import { ref, onMounted, onUnmounted } from 'vue'
import { useUserAuthStore } from '@/stores'
import { client } from '@/composables/useRequest'
import { wsClient } from '@/services/socket/client'
import axios from 'axios'
import type { NotificationResponse } from '@/services/api/notification/routes'

export function useNotificationCount() {
  const store = useUserAuthStore()
  const unreadCount = ref(0)
  let unreadInterval: ReturnType<typeof setInterval> | null = null
  let unsubNotif: (() => void) | null = null

  async function fetchUnreadCount() {
    if (!store.isAuthenticated) return
    try {
      const res = await client.get('/notifications', { params: { limit: 1 } })
      unreadCount.value = res.data?.data?.unreadCount ?? 0
    } catch (err) {
      if (axios.isAxiosError(err) && err.response?.status === 401) return
      console.error('Failed to fetch unread count:', err)
    }
  }

  function startPolling() {
    fetchUnreadCount()
    unreadInterval = setInterval(fetchUnreadCount, 30000)
  }

  function stopPolling() {
    if (unreadInterval) {
      clearInterval(unreadInterval)
      unreadInterval = null
    }
  }

  function onVisibilityChange() {
    if (document.hidden) {
      stopPolling()
    } else if (store.isAuthenticated) {
      startPolling()
    }
  }

  onMounted(() => {
    if (store.isAuthenticated) {
      startPolling()
      unsubNotif = wsClient.on('notification', (msg) => {
        const n = msg.payload as NotificationResponse
        if (n && n.id && !n.isRead) {
          unreadCount.value++
        }
      })
      wsClient.connect()
    }
    document.addEventListener('visibilitychange', onVisibilityChange)
  })

  onUnmounted(() => {
    stopPolling()
    document.removeEventListener('visibilitychange', onVisibilityChange)
    unsubNotif?.()
  })

  return { unreadCount }
}
