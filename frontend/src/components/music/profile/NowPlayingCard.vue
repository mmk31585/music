<template>
  <div
    class="relative overflow-hidden rounded-2xl bg-surface-overlay p-4 ring-1 ring-border-default backdrop-blur-sm transition-all duration-300"
    :class="musicStatus?.playing ? 'ring-spotify/20' : ''"
  >
    <!-- Ambient glow when playing -->
    <div
      v-if="musicStatus?.playing"
      class="pointer-events-none absolute -inset-4 bg-accent/5 blur-3xl motion-safe:animate-pulse"
      aria-hidden="true"
      style="animation-duration: 3s;"
    />

    <div class="relative z-10 flex items-center gap-4">
      <!-- Vinyl cover art -->
      <div class="relative shrink-0">
        <div
          class="h-16 w-16 overflow-hidden rounded-full shadow-lg ring-2"
          :class="musicStatus?.playing ? 'ring-accent/40' : 'ring-border-default'"
        >
          <img
            v-if="currentTrackInfo?.coverUrl"
            :src="currentTrackInfo.coverUrl"
            alt=""
            class="h-full w-full object-cover motion-safe:animate-spin-slow"
            :class="{ 'animate-paused': !musicStatus?.playing }"
            loading="lazy"
          />
          <div v-else class="flex h-full w-full items-center justify-center bg-surface-overlay">
            <Music aria-hidden="true" class="text-lg text-muted"  />
          </div>
        </div>
        <!-- Center dot -->
        <div class="pointer-events-none absolute inset-0 flex items-center justify-center" aria-hidden="true">
          <div class="h-2.5 w-2.5 rounded-full bg-primary shadow" />
        </div>
      </div>

      <!-- Track info -->
      <div class="min-w-0 flex-1">
        <div v-if="musicStatus?.playing" class="flex items-center gap-2 mb-0.5">
          <span class="flex gap-0.5" aria-hidden="true">
            <span
              v-for="i in 4" :key="i"
              class="h-3 w-0.5 rounded-full bg-spotify motion-safe:animate-equalizer"
              :style="{ animationDelay: `${i * 80}ms` }"
            />
          </span>
          <span class="text-[10px] font-semibold uppercase tracking-wider text-accent">Now Playing</span>
        </div>
        <div v-else class="flex items-center gap-1.5 mb-0.5">
          <PauseCircle aria-hidden="true" class="text-[10px] text-muted"  />
          <span class="text-[10px] text-muted">Last played</span>
          <span v-if="lastPlayedLabel" class="text-[10px] text-muted">&middot; {{ lastPlayedLabel }}</span>
        </div>

        <p class="truncate text-sm font-semibold text-primary">{{ currentTrackInfo?.title || 'No track' }}</p>
        <p class="truncate text-xs text-tertiary">{{ currentTrackInfo?.artist || '' }}</p>
      </div>

      <!-- Listen Along button -->
      <button
        v-if="musicStatus?.playing && !isOwnProfile"
        aria-label="Listen along"
        class="shrink-0 rounded-full bg-accent px-4 py-2 text-xs font-bold text-black transition-all duration-200 hover:bg-accent-hover hover:scale-105 active:scale-95 focus-visible:outline-2 focus-visible:outline-primary"
        @click="$emit('listenAlong', currentTrackInfo!.id)"
      >
        Listen
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Music, PauseCircle } from 'lucide-vue-next'
import { computed } from 'vue'

const props = defineProps<{
  musicStatus?: {
    playing: boolean
    currentTrack?: {
      id: string
      title: string
      artist: string
      coverUrl: string
    } | null
    current_track_id?: string | null
    updated_at?: string | null
  } | null
  isOwnProfile?: boolean
  lastPlayedLabel?: string
}>()

defineEmits<{
  listenAlong: [trackId: string]
}>()

const currentTrackInfo = computed(() => props.musicStatus?.currentTrack || null)
</script>

<style scoped>
@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin-slow {
  animation: spin-slow 4s linear infinite;
}
.animate-paused {
  animation-play-state: paused !important;
}

@keyframes equalizer {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(2.5); }
}
.animate-equalizer {
  animation: equalizer 0.6s ease-in-out infinite;
  transform-origin: bottom;
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 0.8; }
}

@media (prefers-reduced-motion: reduce) {
  .animate-spin-slow { animation: none; }
  .animate-equalizer { animation: none; }
  .animate-paused { animation: none; }
  .motion-safe\:animate-pulse { animation: none; }
}
</style>
