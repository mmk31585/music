<template>
  <div
    class="flex items-start gap-4 rounded-xl px-4 py-3 transition hover:bg-white/[0.04]"
    :class="{ 'opacity-50': notification.isRead }"
  >
    <div
      class="mt-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-full"
      :class="iconBgClass"
    >
      <i aria-hidden="true" :class="iconClass" class="text-sm" />
    </div>

    <div class="min-w-0 flex-1 cursor-pointer" role="button" tabindex="0" @click="handleClick" @keydown.enter="handleClick" @keydown.space.prevent="handleClick">
      <p class="text-sm leading-relaxed text-white">
        <span class="font-semibold">{{ notification.title }}</span>
        <span class="ml-1 text-slate-300">{{ notification.body }}</span>
      </p>
      <p class="mt-0.5 text-xs text-slate-500">
        {{ timeAgo }}
      </p>
    </div>

    <button
      v-if="!notification.isRead"
      type="button"
      class="shrink-0 self-center rounded-full bg-white/10 px-3 py-1 text-xs font-bold text-white transition hover:bg-white/15"
      @click.stop="handleMarkRead"
    >
      Mark read
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { NotificationResponse } from '@/services/api/notification'
import { useNotificationApi } from '@/services/api/notification'
import { useToast } from 'primevue/usetoast'

const notifApi = useNotificationApi()

const props = defineProps<{
  notification: NotificationResponse
}>()

const emit = defineEmits<{
  read: [id: string]
  action: [notification: NotificationResponse]
}>()

const toast = useToast()

const iconMap: Record<string, { icon: string; bg: string }> = {
  welcome: { icon: 'pi pi-star', bg: 'bg-amber-500/20 text-amber-400' },
  subscription_purchased: { icon: 'pi pi-crown', bg: 'bg-[#1db954]/20 text-[#1db954]' },
  playlist_created: { icon: 'pi pi-list', bg: 'bg-blue-500/20 text-blue-400' },
  artist_published_track: { icon: 'pi pi-discord', bg: 'bg-purple-500/20 text-purple-400' },
  playlist_shared: { icon: 'pi pi-share-alt', bg: 'bg-sky-500/20 text-sky-400' },
  subscription_renewed: { icon: 'pi pi-sync', bg: 'bg-emerald-500/20 text-emerald-400' },
  user_followed: { icon: 'pi pi-user-plus', bg: 'bg-pink-500/20 text-pink-400' },
}

const iconClass = computed(() => iconMap[props.notification.type]?.icon || 'pi pi-bell')
const iconBgClass = computed(
  () => iconMap[props.notification.type]?.bg || 'bg-white/10 text-slate-400',
)

const timeAgo = computed(() => {
  const ts = props.notification.createdAt
  if (!ts) return ''
  const diff = Date.now() - new Date(ts).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'Just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d ago`
  return new Date(ts).toLocaleDateString()
})

function handleClick() {
  emit('action', props.notification)
}

async function handleMarkRead() {
  try {
    await notifApi.markRead(props.notification.id)
    emit('read', props.notification.id)
  } catch (err) {
    console.error('Failed to mark as read:', err)
    toast.add({ severity: 'error', summary: 'Failed to mark as read', life: 2000 })
  }
}
</script>
