<template>
  <div
    class="group grid grid-cols-[48px_1fr_auto] items-center gap-4 rounded-xl px-3 py-2.5 transition-all duration-200 hover:bg-white/[0.08]"
    :class="isCurrent ? 'bg-white/[0.10] shadow-[inset_3px_0_0_#1db954]' : ''"
  >
    <button
      type="button"
      class="relative flex h-11 w-11 items-center justify-center overflow-hidden rounded-xl bg-white/10 text-white transition-all duration-200 hover:scale-105 hover:bg-[#1db954] hover:text-black disabled:cursor-wait disabled:opacity-70"
      :disabled="loadingThisTrack"
      @click="handlePlay"
    >
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="title"
        class="absolute inset-0 h-full w-full object-cover opacity-60 transition group-hover:opacity-35"
        loading="lazy"
      />

      <span class="relative z-10 flex items-center justify-center">
        <i v-if="loadingThisTrack" class="pi pi-spin pi-spinner text-sm" />

        <span
          v-else-if="isCurrent && player.isPlaying.value"
          class="flex h-4 items-end gap-[2px]"
          aria-label="Playing"
        >
          <span class="eq-bar h-2" />
          <span class="eq-bar h-4 animation-delay-150" />
          <span class="eq-bar h-3 animation-delay-300" />
        </span>

        <i v-else :class="buttonIcon" class="text-sm" />
      </span>
    </button>

    <div class="min-w-0">
      <div
        class="truncate text-sm font-semibold transition"
        :class="isCurrent ? 'text-[#1db954]' : 'text-white'"
      >
        {{ title }}
      </div>

      <div class="mt-0.5 truncate text-xs text-slate-400">
        {{ artistName }}
      </div>
    </div>

    <div class="flex items-center gap-4 text-xs text-slate-400">
      <span v-if="durationLabel" class="tabular-nums">
        {{ durationLabel }}
      </span>

      <button
        type="button"
        class="hidden rounded-full p-2 text-slate-400 transition hover:bg-white/10 hover:text-white group-hover:block"
        title="More"
      >
        <i class="pi pi-ellipsis-h" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'

const props = defineProps<{
  track: any
  index?: number
  queue?: any[]
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

const coverUrl = computed(() => {
  return (
    props.track.coverUrl ||
    props.track.cover_url ||
    props.track.cover ||
    props.track.album?.coverUrl ||
    props.track.album?.cover_url ||
    null
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
    albumTitle: props.track.albumTitle || props.track.album_title || props.track.album?.title || null,
    coverUrl: coverUrl.value,
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

<style scoped>
.eq-bar {
  width: 3px;
  border-radius: 999px;
  background: currentColor;
  animation: equalizer 850ms ease-in-out infinite alternate;
}

.animation-delay-150 {
  animation-delay: 150ms;
}

.animation-delay-300 {
  animation-delay: 300ms;
}

@keyframes equalizer {
  from {
    transform: scaleY(0.45);
    opacity: 0.6;
  }

  to {
    transform: scaleY(1);
    opacity: 1;
  }
}
</style>
