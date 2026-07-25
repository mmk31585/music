<template>
  <div class="relative">
    <button
      ref="queueBtnRef"
      type="button"
      class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
      :class="upcomingCount > 0
        ? 'text-primary hover:bg-surface-active active:scale-95'
        : 'text-tertiary hover:text-primary hover:bg-surface-active active:scale-95'"
      :disabled="!currentTrack"
      :aria-label="`Queue — ${upcomingCount} upcoming`"
      @click="showQueuePreview = !showQueuePreview"
      @keydown.escape="showQueuePreview = false"
    >
      <ListMusic class="h-4 w-4" aria-hidden="true" />
      <span
        v-if="upcomingCount > 0"
        class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-accent text-[8px] font-bold text-black"
      >
        {{ upcomingCount > 9 ? '9+' : upcomingCount }}
      </span>
    </button>

    <!-- Up Next mini-preview popup -->
    <Transition name="fade-scale">
      <div
        v-if="showQueuePreview"
        class="absolute bottom-full right-0 mb-2 z-99 w-72 origin-bottom-right rounded-2xl border border-border-default p-2 shadow-floating glass-strong"
      >
        <!-- Now Playing -->
        <div class="mb-2 px-2 pt-1">
          <p class="text-[10px] font-semibold tracking-wider text-muted uppercase">{{ $t('player.now_playing') }}</p>
          <div class="mt-1.5 flex items-center gap-2.5">
            <div class="h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-surface-active ring-1 ring-border-subtle">
              <img
                v-if="currentTrack?.coverUrl"
                :src="currentTrack.coverUrl"
                :alt="currentTrack.title"
                class="h-full w-full object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-full items-center justify-center">
                <Music class="h-4 w-4 text-muted" aria-hidden="true" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-bold text-primary">{{ currentTrack?.title }}</p>
              <p class="truncate text-xs text-tertiary">{{ currentTrack?.artistName }}</p>
            </div>
            <div v-if="isPlaying" class="flex items-center gap-1">
              <span class="size-1 rounded-full bg-accent animate-equalizer equalizer-bar-1" />
              <span class="size-1 rounded-full bg-accent animate-equalizer equalizer-bar-2" style="animation-delay: 0.15s" />
              <span class="size-1 rounded-full bg-accent animate-equalizer equalizer-bar-3" style="animation-delay: 0.3s" />
            </div>
          </div>
        </div>

        <div class="mx-2 my-1.5 border-t border-border-subtle" />

        <!-- Next Up -->
        <div v-if="nextTrack" class="px-2 pb-2">
          <div class="flex items-center gap-2 text-[10px] font-semibold tracking-wider text-muted uppercase mb-2">
            <ChevronDown class="h-3 w-3" aria-hidden="true" />
            {{ $t('player.up_next') }}
            <span v-if="upcomingCount > 1" class="h-3.5 w-3.5 rounded-full bg-surface-active/80 flex items-center justify-center text-[8px] font-bold text-tertiary">{{ upcomingCount }}</span>
            <Shuffle v-if="shuffleMode !== 'off'" class="h-3 w-3 text-aurora-purple ml-auto" aria-hidden="true" title="Shuffle is on — next track from shuffle order" />
          </div>
          <div
            class="group flex items-center gap-2.5 rounded-xl px-2.5 py-2 transition-all duration-200 cursor-pointer ring-1 ring-border-subtle bg-surface-overlay hover:bg-surface-active"
            @click="playNextTrack"
          >
            <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-surface-active shadow-sm ring-1 ring-border-subtle">
              <img
                v-if="nextTrack.coverUrl"
                :src="nextTrack.coverUrl"
                :alt="nextTrack.title"
                class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
                loading="lazy"
              />
              <div v-else class="flex h-full items-center justify-center">
                <Music class="h-4 w-4 text-muted" aria-hidden="true" />
              </div>
              <div class="absolute inset-0 flex items-center justify-center bg-bg-overlay/0 transition-all duration-200 group-hover:bg-bg-overlay/30">
                <Play class="h-5 w-5 text-primary opacity-0 transition-all duration-200 group-hover:opacity-100 drop-shadow-lg" aria-hidden="true" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-primary group-hover:text-primary transition-colors">
                {{ nextTrack.title }}
              </p>
              <p class="truncate text-xs text-tertiary">{{ nextTrack.artistName }}</p>
            </div>
            <span class="text-[10px] font-mono tabular-nums shrink-0 text-muted group-hover:text-primary/50 transition-colors">
              {{ formatTime(nextTrack.durationSeconds ?? 0) }}
            </span>
          </div>
        </div>

        <!-- Shuffle catalog/similar: can't predict next track -->
        <div v-else-if="shuffleMode === 'catalog' || shuffleMode === 'similar'" class="px-2 pb-2">
          <div class="flex items-center gap-2 text-[10px] font-semibold tracking-wider text-muted uppercase mb-2">
            <Shuffle class="h-3 w-3 text-aurora-purple" aria-hidden="true" />
            {{ $t('player.up_next') }}
            <span class="text-[8px] text-aurora-purple/60 font-normal">— random track</span>
          </div>
          <div
            class="flex items-center gap-2.5 rounded-xl px-2.5 py-2 transition-all duration-200 cursor-pointer ring-1 ring-border-subtle bg-surface-overlay hover:bg-surface-active group"
            @click="playNextTrack"
          >
            <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-aurora-purple/10 ring-1 ring-aurora-purple/20">
              <Shuffle class="h-4 w-4 text-aurora-purple" aria-hidden="true" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium text-secondary group-hover:text-primary transition-colors">Skip to random track</p>
              <p class="text-xs text-muted">Next track selected from {{ shuffleMode === 'catalog' ? 'full catalog' : 'similar tracks' }}</p>
            </div>
          </div>
        </div>

        <div v-else class="flex flex-col items-center gap-1.5 px-2 pb-3 pt-2 text-center">
          <ListMusic class="h-6 w-6 text-muted" aria-hidden="true" />
          <p class="text-xs text-muted">{{ $t('player.queue_empty') }}</p>
        </div>

        <!-- View full queue -->
        <button
          type="button"
          class="flex w-full items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-xs font-medium text-tertiary transition-all hover:bg-surface-overlay hover:text-secondary"
          @click="openFullQueue"
        >
          {{ $t('player.view_full_queue') }}
          <ChevronRight class="h-3 w-3" aria-hidden="true" />
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { ChevronDown, ChevronRight, ListMusic, Music, Play, Shuffle } from 'lucide-vue-next'
import type { PlaybackTrack } from '@/services/api/player'

interface Props {
  currentTrack: PlaybackTrack | null
  nextTrack: PlaybackTrack | null
  upcomingCount: number
  shuffleMode: 'off' | 'queue' | 'catalog' | 'similar'
  isPlaying: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'play-next': []
  'open-queue': []
}>()

const showQueuePreview = ref(false)
const queueBtnRef = ref<HTMLElement | null>(null)

watch(showQueuePreview, async (v) => {
  if (v) await nextTick()
})

function playNextTrack() {
  emit('play-next')
  showQueuePreview.value = false
}

function openFullQueue() {
  showQueuePreview.value = false
  emit('open-queue')
}

// Close on outside click
onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
})
onUnmounted(() => {
  document.removeEventListener('click', handleOutsideClick)
})
function handleOutsideClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (showQueuePreview.value && queueBtnRef.value && !queueBtnRef.value.contains(target)) {
    showQueuePreview.value = false
  }
}

function formatTime(s: number) {
  if (!isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}
</script>

<style scoped>
.fade-scale-enter-active {
  transition: opacity 200ms cubic-bezier(0.34, 1.56, 0.64, 1), transform 200ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
.fade-scale-leave-active {
  transition: opacity 100ms ease, transform 100ms ease;
}
.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(4px);
}
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

@media (prefers-reduced-motion: reduce) {
  .fade-scale-enter-active,
  .fade-scale-leave-active {
    transition: none;
  }
  .fade-scale-enter-from,
  .fade-scale-leave-to {
    opacity: 1;
    transform: none;
  }
}

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