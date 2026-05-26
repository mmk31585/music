<template>
  <div
    class="group flex items-center gap-4 rounded-2xl px-4 py-3 transition hover:bg-white/5"
    :class="isCurrent ? 'bg-white/10' : ''"
    role="button"
    tabindex="0"
    @click="playTrack"
    @keydown.enter.prevent="playTrack"
    @keydown.space.prevent="playTrack"
  >
    <div
      class="flex h-10 w-10 items-center justify-center rounded-lg bg-white/10 text-sm text-slate-300"
    >
      <i v-if="isCurrent && player.isPlaying" class="pi pi-volume-up text-emerald-300" />
      <span v-else>{{ index + 1 }}</span>
    </div>

    <div class="min-w-0 flex-1">
      <p class="truncate font-medium text-white">{{ track.title }}</p>
      <p class="truncate text-sm text-slate-400">
        {{ track.artist_name || 'Unknown artist' }}
        <span v-if="track.album_title">• {{ track.album_title }}</span>
      </p>
    </div>

    <div class="hidden text-sm text-slate-400 md:block">
      {{ formatDuration(track.duration_seconds) }}
    </div>

    <button
      class="rounded-full bg-[#1db954] p-3 text-black transition disabled:cursor-not-allowed disabled:bg-white/10 disabled:text-slate-500 md:opacity-0 md:group-hover:opacity-100"
      :disabled="!track.audio_url"
      :aria-label="track.audio_url ? `Play ${track.title}` : `${track.title} has no audio file`"
      @click.stop="playTrack"
    >
      <i :class="isCurrent && player.isPlaying ? 'pi pi-pause' : 'pi pi-play-fill'" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Track } from '@/services/api/catalog'
import { usePlayerStore } from '@/stores'

const props = defineProps<{
  track: Track
  index: number
  queue?: Track[]
}>()

const player = usePlayerStore()
const isCurrent = computed(() => player.currentTrack?.id === props.track.id)

function formatDuration(value?: number | null) {
  if (!value) return '--:--'
  const mins = Math.floor(value / 60)
  const secs = value % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}

function playTrack() {
  if (!props.track.audio_url) return
  void player.playTrack(props.track, props.queue ?? [props.track])
}
</script>
