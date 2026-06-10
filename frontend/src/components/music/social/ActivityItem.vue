<template>
  <div
    class="group glass hover:glass-hover flex items-start gap-3 rounded-2xl px-4 py-3 transition-all"
  >
    <div
      class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-white/60 transition"
      :class="iconBgClass"
    >
      <i class="text-sm" :class="typeIcon" />
    </div>
    <div class="min-w-0 flex-1">
      <p class="text-sm text-white/70">
        <span class="font-semibold text-white">{{ item.user_display_name }}</span>
        {{ actionLabel }}
        <span v-if="item.target_name" class="font-medium text-white/90">{{
          item.target_name
        }}</span>
      </p>
      <p class="mt-0.5 text-xs text-white/30">{{ timeAgo }}</p>
    </div>
    <button
      v-if="showPlayButton"
      type="button"
      class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-white/10 text-white/50 opacity-0 transition group-hover:opacity-100 hover:bg-[#1db954] hover:text-black"
      @click="handlePlay"
    >
      <i class="pi pi-play-fill text-xs" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'

const props = defineProps<{
  item: {
    id: string
    type: string
    target_type: string
    target_name?: string | null
    target_id?: string | number | null
    user_display_name: string
    created_at: string
  }
}>()

const player = usePlayer()
const playerApi = usePlayerApi()

const showPlayButton = computed(
  () =>
    props.item.type === 'upload' && props.item.target_type === 'track' && !!props.item.target_id,
)

const typeIcon = computed(() => {
  switch (props.item.type) {
    case 'follow':
      return 'pi pi-user-plus'
    case 'like':
      return 'pi pi-heart'
    case 'upload':
      return 'pi pi-upload'
    case 'playlist_create':
      return 'pi pi-list'
    default:
      return 'pi pi-star'
  }
})

const iconBgClass = computed(() => {
  switch (props.item.type) {
    case 'follow':
      return 'bg-blue-500/10 text-blue-400'
    case 'like':
      return 'bg-pink-500/10 text-pink-400'
    case 'upload':
      return 'bg-[#1db954]/10 text-[#1db954]'
    case 'playlist_create':
      return 'bg-purple-500/10 text-purple-400'
    default:
      return 'bg-white/5 text-white/60'
  }
})

const actionLabel = computed(() => {
  switch (props.item.type) {
    case 'follow':
      return 'followed'
    case 'like':
      return 'liked'
    case 'upload':
      return 'uploaded'
    case 'playlist_create':
      return 'created playlist'
    default:
      return props.item.type
  }
})

const timeAgo = computed(() => {
  const now = Date.now()
  const created = new Date(props.item.created_at).getTime()
  const diff = now - created
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d ago`
  return new Date(props.item.created_at).toLocaleDateString()
})

function handlePlay() {
  if (!props.item.target_id) return
  const track = {
    id: String(props.item.target_id),
    title: props.item.target_name || 'Unknown',
    artistName: '',
    albumTitle: null,
    coverUrl: null,
    durationSeconds: null,
    streamUrl: playerApi.getTrackStreamUrl(String(props.item.target_id)),
  }
  player.setQueueAndPlay([track], 0)
}
</script>
