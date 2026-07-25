<template>
  <div
    class="flex min-w-0 w-[25%] items-center gap-3 pb-1"
    :dir="dir"
  >
    <button
      type="button"
      aria-label="Open fullscreen player"
      class="relative shrink-0"
      @click="$emit('open-fullscreen')"
    >
      <div
        class="overflow-hidden rounded-[18px] shadow-[0_16px_32px_rgba(0,0,0,0.5)] ring-1 ring-border-default transition-all duration-700 ease-out-expo"
        :class="[
          isPlaying ? 'scale-100' : 'scale-95 opacity-80',
        ]"
        :style="artworkStyle"
      >
        <img
          v-if="coverUrl"
          :src="coverUrl"
          :alt="title"
          class="h-14 w-14 md:h-16 md:w-16 object-cover"
          loading="lazy"
        />
        <div v-else class="flex h-14 w-14 md:h-16 md:w-16 items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
          <Music class="h-8 w-8 text-tertiary" aria-hidden="true" />
        </div>
        <!-- Dynamic glow from album color -->
        <div
          v-if="isPlaying && progressColor"
          class="absolute inset-0 rounded-[18px] opacity-0 transition-opacity duration-500"
          :style="{ boxShadow: `0 0 40px ${progressColor}80, 0 0 80px ${progressColor}40` }"
        />
      </div>
    </button>

    <div class="min-w-0 flex-1">
      <div class="flex items-center gap-2">
        <p class="truncate text-sm font-bold text-primary leading-tight">{{ title }}</p>
        <!-- Multiple artists with individual links -->
        <div v-if="artists && artists.length > 1" class="flex items-center gap-1 text-xs text-primary/50">
          <template v-for="(artist, i) in artists" :key="artist.id || i">
            <button
              v-if="artist.id"
              type="button"
              class="flex items-center gap-0.5 transition-colors hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-spotify/60 rounded px-0.5"
              @click.stop="$emit('navigate-artist', artist.id)"
              :aria-label="`View ${artist.name}`"
            >
              {{ artist.name }}
              <span v-if="i < artists.length - 1" class="text-primary/30">•</span>
            </button>
            <span v-else class="flex items-center gap-0.5">
              {{ artist.name }}
              <span v-if="i < artists.length - 1" class="text-primary/30">•</span>
            </span>
          </template>
        </div>
        <button
          v-else-if="artistId"
          type="button"
          class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90"
          :class="liked ? 'text-aurora-pink' : 'text-primary/30 hover:text-primary/60'"
          :aria-label="liked ? 'Unlike' : 'Like'"
          @click.stop="$emit('toggle-like')"
        >
          <Heart
            :class="['text-xs', liked ? 'fill-current' : '']"
            :stroke-width="liked ? 0 : 2"
            aria-hidden="true"
          />
        </button>
        <button
          type="button"
          class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90 text-primary/30 hover:text-primary/60"
          aria-label="Add to playlist"
          @click.stop="$emit('add-to-playlist')"
        >
          <Plus class="size-5" aria-hidden="true" />
        </button>
      </div>
      <p v-if="artistName && (!artists || artists.length <= 1)" class="truncate text-xs text-primary/50 mt-0.5 leading-tight">{{ artistName }}</p>
      <p v-else-if="artists && artists.length > 1" class="truncate text-xs text-primary/40 mt-0.5 leading-tight">
        {{ artists.slice(1).map(a => a.name).join(' • ') }}
      </p>
      <!-- Playing indicator - animated equalizer bars -->
      <div v-if="isPlaying" class="flex items-center gap-1 mt-1" role="status" aria-live="polite" aria-label="Now playing">
        <span class="size-1 rounded-full bg-surface-overlay animate-equalizer equalizer-bar-1" />
        <span class="size-1 rounded-full bg-surface-overlay animate-equalizer equalizer-bar-2" style="animation-delay: 0.15s" />
        <span class="size-1 rounded-full bg-surface-overlay animate-equalizer equalizer-bar-3" style="animation-delay: 0.3s" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Heart, Music, Plus } from 'lucide-vue-next'
import { useRTL } from '@/composables/useRTL'

interface Artist {
  id?: string
  name: string
}

const props = defineProps<{
  title: string
  artistName: string
  artistId?: string
  artists?: Artist[]
  coverUrl?: string
  isPlaying: boolean
  liked: boolean
  progressColor?: string
}>()

defineEmits<{
  'open-fullscreen': []
  'toggle-like': []
  'add-to-playlist': []
  'navigate-artist': [artistId: string]
}>()

const { dir } = useRTL()

const artworkStyle = computed(() => {
  if (!props.progressColor) return {}
  return {
    boxShadow: `
      0 16px 32px rgba(0,0,0,0.5),
      0 0 0 1px rgba(255,255,255,0.1),
      0 0 40px ${props.progressColor}40,
      0 0 80px ${props.progressColor}20
    `
  }
})
</script>

<style scoped>
/* Equalizer bar animations */
.animate-equalizer {
  animation: eq-wave 600ms ease-in-out infinite alternate;
  transform-origin: bottom;
}

.equalizer-bar-1 {
  animation-delay: 0ms;
}

.equalizer-bar-2 {
  animation-delay: 150ms;
}

.equalizer-bar-3 {
  animation-delay: 300ms;
}

@keyframes eq-wave {
  0% {
    transform: scaleY(0.4);
    opacity: 0.5;
  }
  100% {
    transform: scaleY(1);
    opacity: 1;
  }
}

</style>