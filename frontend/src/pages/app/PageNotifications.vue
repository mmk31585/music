<template>
  <div class="mx-auto w-full max-w-3xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">Updates</p>
        <h1 class="mt-1 text-3xl font-black text-white">Notifications</h1>
      </div>
      <button
        v-if="unreadCount > 0"
        type="button"
        class="rounded-full bg-white/10 px-4 py-2 text-xs font-bold text-white transition hover:bg-white/15"
        @click="handleMarkAllRead"
      >
        Mark all read
      </button>
    </div>

    <section v-if="today.length" class="mb-8">
      <p class="mb-3 text-xs font-semibold tracking-wider text-slate-500 uppercase">Today</p>
      <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
        <NotificationItem
          v-for="n in today"
          :key="n.id"
          :notification="n"
          @read="handleRead"
          @action="handleAction"
        />
      </div>
    </section>

    <section v-if="thisWeek.length" class="mb-8">
      <p class="mb-3 text-xs font-semibold tracking-wider text-slate-500 uppercase">This Week</p>
      <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
        <NotificationItem
          v-for="n in thisWeek"
          :key="n.id"
          :notification="n"
          @read="handleRead"
          @action="handleAction"
        />
      </div>
    </section>

    <section v-if="earlier.length" class="mb-8">
      <p class="mb-3 text-xs font-semibold tracking-wider text-slate-500 uppercase">Earlier</p>
      <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
        <NotificationItem
          v-for="n in earlier"
          :key="n.id"
          :notification="n"
          @read="handleRead"
          @action="handleAction"
        />
      </div>
    </section>

    <div
      v-if="!notifications.length && !loading && !error"
      class="flex flex-col items-center gap-4 rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-20 text-center"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-full bg-white/10">
        <i aria-hidden="true" class="pi pi-bell text-2xl text-slate-400" />
      </div>
      <h3 class="text-xl font-bold text-white">All caught up!</h3>
      <p class="max-w-sm text-sm text-slate-400">
        You have no notifications yet. Follow artists and interact with music to get updates.
      </p>
    </div>

    <div
      v-if="error"
      class="flex flex-col items-center gap-4 rounded-2xl border border-red-500/20 bg-red-500/5 px-6 py-16 text-center"
    >
      <i aria-hidden="true" class="pi pi-exclamation-triangle text-3xl text-red-400" />
      <h3 class="text-xl font-bold text-white">Failed to load</h3>
      <p class="text-sm text-slate-400">{{ error }}</p>
      <button
        type="button"
        class="rounded-full bg-white/10 px-5 py-2 text-sm font-bold text-white transition hover:bg-white/15"
        @click="loadNotifications"
      >
        Retry
      </button>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 5" :key="i" class="flex items-center gap-4">
        <div class="h-9 w-9 animate-pulse rounded-full bg-white/[0.06]" />
        <div class="flex-1 space-y-2">
          <div class="h-4 w-3/4 animate-pulse rounded bg-white/[0.06]" />
          <div class="h-3 w-1/4 animate-pulse rounded bg-white/[0.06]" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import NotificationItem from '@/components/music/social/NotificationItem.vue'
import { useNotificationApi } from '@/services/api/notification'
import type { NotificationResponse } from '@/services/api/notification'
import { wsClient } from '@/services/socket/client'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const toast = useToast()
const notifApi = useNotificationApi()

const loading = ref(true)
const error = ref('')
const notifications = ref<NotificationResponse[]>([])

const unreadCount = computed(() => notifications.value.filter((n) => !n.isRead).length)

const grouped = computed(() => {
  const now = Date.now()
  const today: NotificationResponse[] = []
  const thisWeek: NotificationResponse[] = []
  const earlier: NotificationResponse[] = []
  for (const n of notifications.value) {
    const diff = now - new Date(n.createdAt).getTime()
    if (diff < 86400000) today.push(n)
    else if (diff < 604800000) thisWeek.push(n)
    else earlier.push(n)
  }
  return { today, thisWeek, earlier }
})

const today = computed(() => grouped.value.today)
const thisWeek = computed(() => grouped.value.thisWeek)
const earlier = computed(() => grouped.value.earlier)

let unsubNotification: (() => void) | null = null

onMounted(() => {
  loadNotifications()
  unsubNotification = wsClient.on('notification', (msg) => {
    const n = msg.payload as NotificationResponse
    if (n && n.id) {
      notifications.value = [n, ...notifications.value]
    }
  })
  wsClient.connect()
})

onUnmounted(() => {
  unsubNotification?.()
})

async function loadNotifications() {
  loading.value = true
  error.value = ''
  try {
    const res = await notifApi.getNotifications()
    notifications.value = res.notifications
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}

async function handleMarkAllRead() {
  try {
    await notifApi.markAllRead()
    notifications.value = notifications.value.map((n) => ({ ...n, isRead: true }))
    toast.add({ severity: 'success', summary: 'All marked as read', life: 2000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to mark all as read', life: 2000 })
  }
}

function handleRead(id: string) {
  const n = notifications.value.find((x) => x.id === id)
  if (n) n.isRead = true
}

function handleAction(notification: NotificationResponse) {
  const et = notification.entityType
  const eid = notification.entityId
  if (et && eid) {
    if (et === 'track') router.push(`/track/${eid}`)
    else if (et === 'playlist') router.push(`/playlist/${eid}`)
    else if (et === 'subscription') router.push('/subscription')
    else if (et === 'user') router.push(`/user/${eid}`)
    else if (et === 'artist') router.push(`/artist/${eid}`)
  }
}
</script>
