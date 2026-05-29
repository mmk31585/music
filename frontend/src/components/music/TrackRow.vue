<template>
  <div
    class="group grid grid-cols-[48px_1fr_auto] items-center gap-4 rounded-xl px-3 py-2 transition hover:bg-white/10"
    :class="isCurrent ? 'bg-white/10' : ''"
  >
    <button
      type="button"
      class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-[#1db954] hover:text-black"
      :disabled="loadingThisTrack"
      @click="handlePlay"
    >
      <i v-if="loadingThisTrack" class="pi pi-spin pi-spinner text-sm" />
      <i v-else :class="buttonIcon" class="text-sm" />
    </button>

    <div class="min-w-0">
      <div class="truncate text-sm font-semibold text-white">
        {{ title }}
      </div>

      <div class="truncate text-xs text-slate-400">
        {{ artistName }}
      </div>
    </div>

    <div class="flex items-center gap-3 text-xs text-slate-400">
      <span v-if="durationLabel">{{ durationLabel }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'

const props = defineProps<{
  track: any
}>()

const player = usePlayer()
const playerApi = usePlayerApi()

const title = computed(() => props.track.title || 'Untitled')
const artistName = computed(() => {
  return (
    props.track.artistName ||
    props.track.artist_name ||
    props.track.artist?.name ||
    props.track.artist ||
    'Unknown artist'
  )
})

const durationSeconds = computed(() => {
  return props.track.durationSeconds ?? props.track.duration_seconds ?? props.track.duration ?? null
})

const durationLabel = computed(() => {
  const total = Number(durationSeconds.value)
  if (!Number.isFinite(total) || total <= 0) return ''

  const minutes = Math.floor(total / 60)
  const seconds = Math.floor(total % 60)

  return `${minutes}:${String(seconds).padStart(2, '0')}`
})

const playbackTrack = computed<PlaybackTrack>(() => {
  const id = String(props.track.id)

  return {
    id,
    title: title.value,
    artistName: artistName.value,
    albumTitle:
      props.track.albumTitle || props.track.album_title || props.track.album?.title || null,
    coverUrl: props.track.coverUrl || props.track.cover_url || props.track.cover || null,
    durationSeconds: durationSeconds.value,
    streamUrl: playerApi.getTrackStreamUrl(id),
  }
})

const isCurrent = computed(() => {
  return player.currentTrack.value?.id === playbackTrack.value.id
})

const loadingThisTrack = computed(() => {
  return isCurrent.value && (player.isLoadingTrack.value || player.isBuffering.value)
})

const buttonIcon = computed(() => {
  if (isCurrent.value && player.isPlaying.value) return 'pi pi-pause'
  return 'pi pi-play'
})

async function handlePlay() {
  await player.toggleTrack(playbackTrack.value)
}
</script>
