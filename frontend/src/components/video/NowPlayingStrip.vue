<template>
  <Transition name="slide-up">
    <div
      v-if="track"
      class="flex items-center gap-3 rounded-2xl bg-black/70 px-4 py-3 backdrop-blur-xs"
    >
      <!-- Album art -->
      <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-xl">
        <AppImage
          :src="track.cover_url"
          :alt="track.title"
          class="h-full w-full object-cover"
          fallback-icon="pi pi-music"
          icon-size="0.875rem"
        />
        <div class="absolute inset-0 rounded-xl ring-1 ring-inset ring-white/10" />
      </div>

      <!-- Track info -->
      <div class="min-w-0 flex-1">
        <p class="truncate text-sm font-semibold text-white">
          {{ track.title }}
        </p>
        <p class="truncate text-xs text-white/50">
          {{ track.artist_name || 'Unknown artist' }}
        </p>
      </div>

      <!-- Play button -->
      <button
        type="button"
        class="inline-flex shrink-0 items-center gap-1.5 rounded-xl bg-white/10 px-3.5 py-2 text-xs font-bold text-white transition hover:bg-spotify hover:text-black"
        @click="$emit('play', track)"
      >
        <i aria-hidden="true" class="pi pi-play-fill text-sm" />
        <span>پخش آهنگ</span>
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import type { TrackSummary } from '@/services/api/video/types'
import AppImage from '@/components/common/AppImage.vue'

defineProps<{
  track: TrackSummary | null
}>()

defineEmits<{
  play: [track: TrackSummary]
}>()
</script>

<style scoped>
.slide-up-enter-active {
  transition: transform 200ms ease-out, opacity 200ms ease-out;
}
.slide-up-leave-active {
  transition: transform 150ms ease-in, opacity 150ms ease-in;
}
.slide-up-enter-from {
  transform: translateY(100%);
  opacity: 0;
}
.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0;
}
</style>
