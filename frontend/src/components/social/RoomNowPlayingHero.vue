<template>
  <div
    class="room-now-playing-hero relative overflow-hidden rounded-2xl p-6"
    :class="{ 'min-h-[160px]': !nowPlaying }"
  >
    <!-- Aurora background -->
    <div
      class="pointer-events-none absolute inset-0 opacity-30"
      :style="auroraStyle"
    />
    <div class="pointer-events-none absolute inset-0 bg-gradient-to-b from-black/60 to-black/20" />

    <div class="relative z-10">
      <h2 class="mb-4 text-xs font-bold uppercase tracking-wider text-white/30">
        Now Playing
      </h2>

      <div v-if="isLoading" class="flex items-center gap-4">
        <div class="h-24 w-24 shrink-0 animate-pulse rounded-2xl bg-white/10" />
        <div class="min-w-0 flex-1 space-y-2">
          <div class="h-5 w-3/4 animate-pulse rounded bg-white/10" />
          <div class="h-4 w-1/2 animate-pulse rounded bg-white/10" />
          <div class="h-3 w-1/3 animate-pulse rounded bg-white/10" />
        </div>
      </div>

      <div v-else-if="nowPlaying" class="flex flex-col gap-4 sm:flex-row sm:items-center sm:gap-6">
        <div class="relative mx-auto sm:mx-0">
          <div class="disc" :class="{ spinning: isPlaying }">
            <div class="disc-inner">
              <img
                v-if="nowPlaying.track.cover_url"
                :src="nowPlaying.track.cover_url"
                :alt="nowPlaying.track.title"
              />
              <div v-else class="flex h-full w-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-headphones text-2xl text-white/40" />
              </div>
            </div>
            <div class="disc-hole" />
          </div>
        </div>

        <div class="min-w-0 flex-1 text-center sm:text-left">
          <p class="truncate text-lg font-bold text-white">{{ nowPlaying.track.title }}</p>
          <p class="truncate text-sm text-white/50">{{ nowPlaying.track.artist_name }}</p>
          <p v-if="nowPlaying.track.album_title" class="mt-0.5 truncate text-xs text-white/30">
            {{ nowPlaying.track.album_title }}
          </p>

          <div class="mt-3 flex items-center justify-center gap-3 sm:justify-start">
            <span class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs"
              :class="sourceBadgeClass"
            >
              <i aria-hidden="true" v-if="nowPlaying.source === 'autofill'" class="pi pi-refresh text-xs" />
              <i aria-hidden="true" v-else class="pi pi-thumbs-up text-xs" />
              {{ sourceLabel }}
            </span>
            <button
              class="flex h-10 w-10 items-center justify-center rounded-full transition hover:scale-105 active:scale-95"
              :class="isPlaying
                ? 'bg-white/10 text-white hover:bg-white/20'
                : 'bg-[#1db954]/20 text-[#1db954] hover:bg-[#1db954]/30 play-pulse'"
              @click="$emit('toggle-play')"
              aria-label="Toggle play"
            >
              <i aria-hidden="true" v-if="isPlaying" class="pi pi-pause text-lg" />
              <i aria-hidden="true" v-else class="pi pi-play ml-0.5 text-lg" />
            </button>
            <span class="text-xs text-white/30">
              <span v-if="isPlaying" class="text-green-400">● Playing</span>
              <span v-else class="text-white/40">⏸ Paused</span>
            </span>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center gap-3 py-8 text-sm text-white/30">
        <i aria-hidden="true" class="pi pi-headphones text-3xl" />
        <span>Waiting for a track to play...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RoomNowPlaying } from '@/services/api/social/room-queue'

const props = defineProps<{
  nowPlaying: RoomNowPlaying | null
  isLoading: boolean
  isPlaying: boolean
}>()

defineEmits<{
  'toggle-play': []
}>()

const auroraStyle = computed(() => {
  if (!props.nowPlaying?.track.cover_url) return {}
  return {
    background: `radial-gradient(ellipse at 50% 0%, rgba(29,185,84,0.15) 0%, transparent 70%)`,
  }
})

const sourceLabel = computed(() => {
  if (!props.nowPlaying) return ''
  return props.nowPlaying.source === 'autofill'
    ? 'پخش خودکار از پرشنیده‌ها'
    : 'پیشنهاد شده'
})

const sourceBadgeClass = computed(() => {
  if (!props.nowPlaying) return ''
  return props.nowPlaying.source === 'autofill'
    ? 'bg-white/10 text-white/50'
    : 'bg-green-500/10 text-green-400'
})
</script>

<style scoped>
.room-now-playing-hero {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.07);
  transition: border-color 0.3s ease;
}

.room-now-playing-hero:has(.spinning) {
  border-color: rgba(29, 185, 84, 0.2);
}

.disc {
  width: 96px;
  height: 96px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.disc-inner {
  width: 100%;
  height: 100%;
  border-radius: 999px;
  overflow: hidden;
  background: #1a1a1a;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.disc-inner img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.disc.spinning .disc-inner {
  animation: spinDisc 6s linear infinite;
}

.disc-hole {
  position: absolute;
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.5);
  border: 2px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(4px);
}

@keyframes spinDisc {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Pulse animation for the play button when track is loaded but not playing */
@keyframes playPulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(29, 185, 84, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(29, 185, 84, 0);
  }
}

.play-pulse {
  animation: playPulse 2s ease-in-out infinite;
}
</style>
