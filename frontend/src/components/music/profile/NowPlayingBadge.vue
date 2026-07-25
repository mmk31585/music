<template>
  <div v-if="musicStatus?.playing" class="glass-strong rounded-full px-3 py-1.5 flex items-center gap-2">
    <template v-if="musicStatus.currentTrack">
      <img
        v-if="musicStatus.currentTrack.coverUrl"
        :src="musicStatus.currentTrack.coverUrl"
        alt=""
        class="w-8 h-8 rounded-lg object-cover shrink-0"
        loading="lazy"
      />
      <div v-else class="w-8 h-8 rounded-lg bg-white/8 flex items-center justify-center shrink-0">
        <Music aria-hidden="true" class="text-xs text-white/40"  />
      </div>
      <div class="flex items-center gap-2 min-w-0">
        <span class="flex gap-0.5 items-end h-4" aria-hidden="true">
          <span class="eq-bar eq-bar-1 w-0.5 rounded-full bg-primary" />
          <span class="eq-bar eq-bar-2 w-0.5 rounded-full bg-primary" />
          <span class="eq-bar eq-bar-3 w-0.5 rounded-full bg-primary" />
        </span>
        <span class="text-[11px] text-white/70 truncate max-w-[200px]">
          <span class="font-semibold text-white/90">NOW PLAYING</span>
          &middot; {{ musicStatus.currentTrack.title }} &middot; {{ musicStatus.currentTrack.artist }}
        </span>
      </div>
      <button
        v-if="showListenAlong"
        aria-label="Listen along"
        class="rounded-full bg-primary px-3 py-1 text-[11px] font-bold text-black transition hover:bg-primary-hover active:scale-95"
        @click="$emit('listen-along', musicStatus.currentTrack!.id)"
      >
        گوش بده
      </button>
    </template>
    <template v-else>
      <span class="flex gap-0.5 items-end h-4" aria-hidden="true">
        <span class="eq-bar eq-bar-1 w-0.5 rounded-full bg-primary" />
        <span class="eq-bar eq-bar-2 w-0.5 rounded-full bg-primary" />
        <span class="eq-bar eq-bar-3 w-0.5 rounded-full bg-primary" />
      </span>
      <span class="text-[11px] text-white/70">در حال گوش دادن...</span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { Music } from 'lucide-vue-next'
defineProps<{
  musicStatus?: {
    playing: boolean
    currentTrack?: {
      id: string
      title: string
      artist: string
      coverUrl: string
    } | null
  } | null
  showListenAlong?: boolean
}>()

defineEmits<{
  'listen-along': [trackId: string]
}>()
</script>

<style scoped>
.glass-strong {
  background: rgba(18,18,18,0.8);
  backdrop-filter: blur(28px);
}

@keyframes eq-bar {
  0%, 100% { height: 4px; }
  50% { height: 16px; }
}
.eq-bar-1 { animation: eq-bar 0.8s ease-in-out infinite; }
.eq-bar-2 { animation: eq-bar 0.8s ease-in-out infinite 0.15s; }
.eq-bar-3 { animation: eq-bar 0.8s ease-in-out infinite 0.3s; }

@media (prefers-reduced-motion: reduce) {
  .eq-bar-1, .eq-bar-2, .eq-bar-3 { animation: none; height: 8px; }
}
</style>
