<template>
  <div class="group relative overflow-hidden rounded-2xl bg-surface-overlay/50 transition-all duration-200 hover:bg-surface-active/80">
    <div
      v-if="albumCover"
      class="pointer-events-none absolute inset-0 scale-110"
      aria-hidden="true"
    >
      <img
        :src="albumCover"
        alt=""
        class="h-full w-full object-cover opacity-[0.08]"
        style="filter: blur(40px) saturate(1.5)"
      />
    </div>

    <div class="relative z-10 p-5">
      <div class="flex items-start gap-4">
        <div class="relative shrink-0">
          <div
            class="h-24 w-24 overflow-hidden rounded-2xl shadow-[0_8px_32px_rgba(0,0,0,0.5)] ring-1 ring-border-default transition-all duration-500"
            :class="isPlaying ? 'scale-100' : 'scale-95 opacity-80'"
          >
            <img
              v-if="albumCover"
              :src="albumCover"
              :alt="track.title"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-primary/30 to-accent/30">
              <Headphones :size="22" class="text-tertiary" />
            </div>
          </div>
          <div
            v-if="isPlaying"
            class="absolute -bottom-1 -right-1 flex h-6 w-6 items-center justify-center rounded-full bg-primary shadow-lg"
          >
            <Music :size="10" class="text-black" />
          </div>
        </div>

        <div class="min-w-0 flex-1 pt-1">
          <p class="truncate text-base font-bold text-primary leading-tight">{{ track.title }}</p>
          <p class="mt-0.5 truncate text-sm text-secondary">{{ track.artistName }}</p>
          <p v-if="track.albumTitle" class="mt-0.5 truncate text-xs text-tertiary">{{ track.albumTitle }}</p>

          <div class="mt-3 flex items-center gap-1">
            <button
              type="button"
              aria-label="Like"
              class="flex h-8 w-8 items-center justify-center rounded-full text-secondary transition-all duration-150 hover:bg-surface-active hover:text-primary active:scale-90"
              @click="$emit('toggle-like')"
            >
              <Heart
                v-if="liked"
                :size="14"
                class="text-accent"
                fill="currentColor"
              />
              <Heart v-else :size="14" />
            </button>
            <button
              type="button"
              aria-label="Fullscreen"
              class="flex h-8 w-8 items-center justify-center rounded-full text-secondary transition-all duration-150 hover:bg-surface-active hover:text-primary active:scale-90"
              @click="$emit('toggle-fullscreen')"
            >
              <Maximize2 :size="14" />
            </button>
          </div>
        </div>
      </div>

      <div class="mt-5">
        <div
          role="slider"
          tabindex="0"
          :aria-valuemin="0"
          :aria-valuemax="1"
          :aria-valuenow="progressPercent / 100"
          aria-label="Seek"
          class="group/progress relative h-1 cursor-pointer rounded-full transition-all duration-150 hover:h-1.5"
          :class="isPlaying ? 'bg-surface-active' : 'bg-surface-overlay'"
          @click="onSeek"
          @keydown.enter="onSeek"
          @keydown.space.prevent="onSeek"
        >
          <div
            v-if="isPlaying && accentColor"
            class="absolute inset-0 rounded-full opacity-20 blur-xs transition-opacity"
            :style="{ background: accentColor }"
          />
          <div
            class="relative h-full rounded-full transition-all duration-100"
            :style="{ width: `${progressPercent}%`, background: accentColor || 'var(--color-primary)', boxShadow: accentColor ? `0 0 6px ${accentColor}44, 0 0 12px ${accentColor}22` : undefined }"
          />
        </div>
        <div class="mt-1.5 flex items-center justify-between">
          <span class="text-[10px] font-mono tabular-nums text-tertiary">{{ formatTime(currentTime) }}</span>
          <span class="text-[10px] font-mono tabular-nums text-tertiary">{{ formatTime(duration) }}</span>
        </div>
      </div>

      <div class="mt-4 flex items-center justify-center gap-2">
        <button
          type="button"
          aria-label="Shuffle"
          class="flex h-9 w-9 items-center justify-center rounded-full transition-all duration-150 active:scale-90"
          :class="shuffleMode !== 'off' ? 'text-accent' : 'text-secondary hover:bg-surface-active hover:text-primary'"
          @click="$emit('toggle-shuffle')"
        >
          <Shuffle :size="16" />
        </button>
        <button
          type="button"
          aria-label="Previous track"
          class="flex h-9 w-9 items-center justify-center rounded-full text-secondary transition-all duration-150 hover:bg-surface-active hover:text-primary active:scale-90 disabled:opacity-30"
          :disabled="!hasPrevious"
          @click="$emit('prev')"
        >
          <SkipBack :size="16" />
        </button>
        <button
          type="button"
          aria-label="Play or pause"
          class="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-black shadow-lg transition-all duration-150 hover:bg-primary/90 active:scale-90 disabled:opacity-40"
          :disabled="!track"
          @click="$emit('toggle-play')"
        >
          <Loader2 v-if="isLoading || isBuffering" :size="16" class="animate-spin" />
          <Play v-else-if="!isPlaying" :size="16" class="ml-0.5" />
          <Pause v-else :size="16" />
        </button>
        <button
          type="button"
          aria-label="Next track"
          class="flex h-9 w-9 items-center justify-center rounded-full text-secondary transition-all duration-150 hover:bg-surface-active hover:text-primary active:scale-90 disabled:opacity-30"
          :disabled="!hasNext"
          @click="$emit('next')"
        >
          <SkipForward :size="16" />
        </button>
        <button
          type="button"
          aria-label="Repeat"
          class="flex h-9 w-9 items-center justify-center rounded-full transition-all duration-150 active:scale-90"
          :class="repeatMode !== 'off' ? 'text-accent' : 'text-secondary hover:bg-surface-active hover:text-primary'"
          @click="$emit('toggle-repeat')"
        >
          <Repeat :size="16" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Headphones, Heart, Loader2, Maximize2, Music, Pause, Play, Repeat, Shuffle, SkipBack, SkipForward } from 'lucide-vue-next'
import type { PlaybackTrack } from '@/services/api/player/types'

defineProps<{
  track: PlaybackTrack
  isPlaying: boolean
  isLoading: boolean
  isBuffering: boolean
  liked: boolean
  currentTime?: number | null
  duration?: number | null
  progressPercent: number
  shuffleMode: string
  repeatMode: string
  hasNext: boolean
  hasPrevious: boolean
  albumCover?: string | null
  accentColor?: string | null
}>()

const emit = defineEmits<{
  'toggle-like': []
  'toggle-fullscreen': []
  'toggle-play': []
  'toggle-shuffle': []
  'toggle-repeat': []
  'prev': []
  'next': []
  'seek': [percent: number]
}>()

function onSeek(e: MouseEvent | KeyboardEvent) {
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = Math.max(0, Math.min(100, ((e as MouseEvent).clientX - rect.left) / rect.width * 100))
  emit('seek', pct)
}

function formatTime(seconds?: number | null) {
  if (typeof seconds !== 'number' || !isFinite(seconds)) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>
