<template>
  <Teleport to="body">
    <Transition name="fullscreen-slide">
      <div
        v-if="visible"
        ref="rootEl"
        class="fixed inset-0 z-50 flex flex-col overflow-hidden text-white select-none outline-none"
        :style="resolvedBgStyle"
        @keydown="onKeydown"
        tabindex="0"
      >
        <!-- Aurora background layer -->
        <div class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
          <img
            v-if="coverUrl"
            :src="coverUrl"
            alt=""
            class="absolute h-[120%] w-[120%] -top-[10%] -left-[10%] object-cover opacity-[0.15]"
            style="filter: blur(80px) saturate(1.2)"
          />
          <div class="absolute inset-0 bg-black/55" />
          <div
            class="absolute inset-0"
            :style="auroraGradient"
          />
        </div>

        <!-- Header -->
        <div class="relative z-10 flex h-14 items-center justify-between px-4 md:px-6">
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-2 focus-visible:ring-[#1db954]"
            @click="close"
            aria-label="Close"
          >
            <i class="pi pi-chevron-down text-lg" />
          </button>

          <div class="text-[13px] font-semibold uppercase tracking-widest text-[rgba(255,255,255,0.6)]">
            NOW PLAYING
          </div>

          <div class="flex items-center gap-1">
            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-2 focus-visible:ring-[#1db954]"
              @click="showOverflow = !showOverflow"
              aria-label="More options"
            >
              <i class="pi pi-ellipsis-h text-sm" />
            </button>
            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-2 focus-visible:ring-[#1db954]"
              :class="{ '!text-[#1db954]': activeTab === 'queue' }"
              @click="activeTab = 'queue'"
              aria-label="Toggle queue"
            >
              <i class="pi pi-list text-sm" />
            </button>
          </div>
        </div>

        <!-- Content -->
        <div class="relative z-10 flex flex-1 flex-col overflow-hidden">
          <div
            v-if="activeTab === 'now-playing'"
            class="flex flex-1 flex-col items-center justify-center gap-5 px-4 pb-4 overflow-y-auto"
          >
            <div class="relative h-[280px] w-[280px] md:h-[320px] md:w-[320px]">
              <div
                class="relative h-full w-full overflow-hidden rounded-2xl shadow-[0_24px_80px_rgba(0,0,0,0.7)] transition-all duration-1000"
                :class="{ 'shadow-[0_0_60px_rgba(29,185,84,0.2)]': isPlaying }"
              >
                <img
                  v-if="coverUrl"
                  :src="coverUrl"
                  :alt="title"
                  loading="lazy"
                  class="h-full w-full object-cover transition-all duration-300"
                  :class="{ 'animate-[art-pop_300ms_cubic-bezier(0.34,1.56,0.64,1)]': !!currentTrack }"
                  @error="onImgError"
                />
                <div
                  v-else
                  class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212]"
                >
                  <i class="pi pi-music text-5xl text-white/30" />
                </div>
              </div>
            </div>

            <div class="flex w-full max-w-md flex-col items-center gap-1 text-center">
              <h2 class="max-w-full truncate text-[22px] font-bold text-white">
                {{ title }}
              </h2>
              <div class="flex items-center gap-3">
                <p class="text-base text-[rgba(255,255,255,0.6)]">
                  {{ artistName }}
                </p>
                <button
                  type="button"
                  :aria-label="isLiked ? 'Unlike' : 'Like'"
                  class="flex items-center justify-center transition-all active:scale-90"
                  :class="isLiked ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)] hover:text-white'"
                  @click="isLiked = !isLiked"
                >
                  <i :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-lg" />
                </button>
              </div>
            </div>

            <div class="w-full max-w-md">
              <div
                ref="progressRef"
                class="group/seeks relative flex h-5 cursor-pointer items-center"
                @click="seekFromEvent"
                @mousemove="onProgressHover"
                @mouseleave="hoverPos = null"
              >
                <div class="h-1 w-full rounded-full bg-white/[0.15] transition-all duration-150 group-hover/seeks:h-1.5">
                  <div
                    class="relative h-full rounded-full bg-[#1db954] transition-all duration-100"
                    :style="{ width: `${Math.min(progressPct, 100)}%` }"
                  >
                    <div
                      class="absolute -right-1.5 -top-1.5 h-3 w-3 scale-0 rounded-full bg-white shadow-lg transition-transform duration-150 group-hover/seeks:scale-100"
                    />
                  </div>
                </div>
                <div
                  v-if="hoverPos !== null"
                  class="pointer-events-none absolute -top-7 rounded-md bg-black/80 px-2 py-1 text-xs tabular-nums text-white"
                  :style="{ left: `${hoverPos}%` }"
                >
                  {{ hoverTimeLabel }}
                </div>
              </div>
              <div class="mt-1 flex justify-between text-xs tabular-nums text-[rgba(255,255,255,0.35)]">
                <span>{{ currentTimeLabel }}</span>
                <span>{{ durationLabel }}</span>
              </div>
            </div>

            <div class="flex items-center gap-5">
              <button
                type="button"
                :aria-label="shuffleMode ? 'Shuffle on' : 'Shuffle off'"
                class="flex h-8 w-8 items-center justify-center rounded-full transition-all hover:bg-white/10 hover:text-white active:scale-90"
                :class="shuffleMode ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
                @click="toggleShuffle"
              >
                <i class="pi pi-sort-alt text-sm" />
              </button>

              <button
                type="button"
                class="flex h-10 w-10 items-center justify-center rounded-full text-white/60 transition-all hover:bg-white/10 hover:text-white active:scale-90 disabled:opacity-20"
                :disabled="!hasPrevious"
                @click="playPrevious"
                aria-label="Previous track"
              >
                <i class="pi pi-step-backward text-xl" />
              </button>

              <button
                type="button"
                :aria-label="isPlaying ? 'Pause' : 'Play'"
                class="relative flex h-16 w-16 items-center justify-center rounded-full bg-[#1db954] text-white shadow-lg transition-all active:scale-95 disabled:opacity-40"
                :style="{ transitionTimingFunction: 'var(--ease-spring)', transitionDuration: '150ms' }"
                :disabled="!currentTrack || isLoadingTrack"
                @click="togglePlayPause"
              >
                <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-xl" />
                <i v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-xl" />
              </button>

              <button
                type="button"
                class="flex h-10 w-10 items-center justify-center rounded-full text-white/60 transition-all hover:bg-white/10 hover:text-white active:scale-90 disabled:opacity-20"
                :disabled="!hasNext"
                @click="playNext"
                aria-label="Next track"
              >
                <i class="pi pi-step-forward text-xl" />
              </button>

              <button
                type="button"
                :aria-label="repeatTitle"
                class="relative flex h-8 w-8 items-center justify-center rounded-full transition-all hover:bg-white/10 hover:text-white active:scale-90"
                :class="repeatMode !== 'off' ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
                @click="toggleRepeat"
              >
                <i class="pi pi-refresh text-sm" />
                <span
                  v-if="repeatMode === 'one'"
                  class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#1db954] text-[9px] font-bold text-black"
                >1</span>
              </button>
            </div>

            <div class="flex items-center gap-4 text-white/40">
              <button
                type="button"
                class="flex h-8 w-8 items-center justify-center rounded-full transition-all hover:text-white active:scale-90"
                :aria-label="muted ? 'Unmute' : 'Mute'"
                @click="toggleMute"
              >
                <i :class="volumeIcon" class="text-sm" />
              </button>
              <div class="group/vol relative flex items-center">
                <div class="relative h-1 w-24 overflow-hidden rounded-full bg-white/[0.15]">
                  <div
                    class="h-full rounded-full bg-white/60"
                    :style="{ width: `${muted ? 0 : Number(volume) * 100}%` }"
                  />
                </div>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  class="absolute inset-0 cursor-pointer opacity-0"
                  :value="muted ? 0 : volume"
                  @input="onVolume"
                />
              </div>

              <button
                type="button"
                class="flex h-8 items-center gap-1.5 rounded-full bg-white/[0.04] px-3 text-xs font-bold transition-all hover:bg-white/10 hover:text-white active:scale-95"
                :class="{ 'bg-[#1db954]/10 text-[#1db954]': playbackRate !== 1 }"
                @click="cycleSpeed"
              >
                <i class="pi pi-forward text-[10px]" />
                <span class="tabular-nums">{{ speedLabel }}</span>
              </button>
            </div>

            <div v-if="upNextTracks.length" class="w-full max-w-md">
              <p class="mb-2 px-1 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">Up Next</p>
              <div class="space-y-1">
                <div
                  v-for="(track, idx) in upNextTracks"
                  :key="track.id"
                  class="group flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2 transition-all hover:bg-white/[0.06]"
                  @click="playQueueItem(queueIndex + 1 + idx)"
                >
                  <div class="h-8 w-8 shrink-0 overflow-hidden rounded-lg bg-white/10">
                    <img
                      v-if="track.coverUrl"
                      :src="track.coverUrl"
                      :alt="track.title"
                      loading="lazy"
                      class="h-full w-full object-cover"
                      @error="onImgError"
                    />
                    <div v-else class="flex h-full items-center justify-center">
                      <i class="pi pi-music text-xs text-white/30" />
                    </div>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-medium text-white/80">{{ track.title }}</p>
                    <p class="truncate text-xs text-white/40">{{ track.artistName }}</p>
                  </div>
                  <span class="text-xs text-white/30 tabular-nums">{{ formatTime(track.durationSeconds) }}</span>
                </div>
              </div>
            </div>
          </div>

          <div
            v-if="activeTab === 'queue'"
            class="flex flex-1 flex-col overflow-hidden"
          >
            <div class="flex-1 overflow-y-auto px-4 py-4 md:px-6">
              <div v-if="currentTrack" class="mb-6">
                <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">Now Playing</p>
                <div class="flex items-center gap-3 rounded-xl border-l-2 border-[#1db954] bg-white/[0.03] px-4 py-3">
                  <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                    <img
                      v-if="currentTrack.coverUrl"
                      :src="currentTrack.coverUrl"
                      :alt="currentTrack.title"
                      loading="lazy"
                      class="h-full w-full object-cover"
                      @error="onImgError"
                    />
                    <div v-else class="flex h-full items-center justify-center">
                      <i class="pi pi-music text-xs text-white/30" />
                    </div>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-bold text-white">{{ currentTrack.title }}</p>
                    <p class="truncate text-xs text-white/40">{{ currentTrack.artistName }}</p>
                  </div>
                  <i class="pi pi-waveform text-lg text-[#1db954]" />
                </div>
              </div>

              <div>
                <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">Next Up</p>
                <div v-if="queueTracks.length" class="space-y-1">
                  <div
                    v-for="(track, idx) in queueTracks"
                    :key="track.id"
                    class="group flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 transition-all hover:bg-white/[0.06]"
                    @click="playQueueItem(idx)"
                  >
                    <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                      <img
                        v-if="track.coverUrl"
                        :src="track.coverUrl"
                        :alt="track.title"
                        loading="lazy"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div v-else class="flex h-full items-center justify-center">
                        <i class="pi pi-music text-xs text-white/30" />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm font-medium text-white/80">{{ track.title }}</p>
                      <p class="truncate text-xs text-white/40">{{ track.artistName }}</p>
                    </div>
                    <button
                      type="button"
                      class="flex h-7 w-7 items-center justify-center rounded-full text-white/20 opacity-0 transition-all hover:bg-white/10 hover:text-white/60 group-hover:opacity-100"
                      @click.stop="removeFromQueue(idx)"
                      aria-label="Remove from queue"
                    >
                      <i class="pi pi-times text-xs" />
                    </button>
                  </div>
                </div>
                <div
                  v-else
                  class="flex flex-col items-center justify-center py-16 text-center"
                >
                  <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-white/5">
                    <i class="pi pi-list text-lg text-white/20" />
                  </div>
                  <p class="text-sm text-white/25">Your queue is empty</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Tab bar -->
        <div class="relative z-10 flex justify-center border-t border-white/[0.06] bg-black/20 backdrop-blur-sm">
          <div class="relative flex gap-8">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              type="button"
              class="relative py-3 text-sm font-medium transition-colors"
              :class="activeTab === tab.key ? 'text-white' : 'text-[rgba(255,255,255,0.35)]'"
              @click="activeTab = tab.key"
            >
              {{ tab.label }}
            </button>
            <div
              class="absolute bottom-0 h-0.5 bg-[#1db954] transition-all duration-200"
              :style="tabIndicatorStyle"
              style="transition-timing-function: cubic-bezier(0.19, 1, 0.22, 1)"
            />
          </div>
        </div>

        <!-- Overflow menu -->
        <div v-if="showOverflow" class="absolute top-14 right-4 z-20">
          <PlayerOverflowMenu
            @close="showOverflow = false"
            @toggle-pip="emit('toggle-pip')"
          />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { usePlayer } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { onImgError } from '@/utils/helpers'
import type { PlaybackTrack } from '@/services/api/player'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'

const visible = defineModel<boolean>('visible', { default: false })

const emit = defineEmits<{
  'toggle-pip': []
}>()

const player = usePlayer()

const currentTrack = computed(() => player.currentTrack.value)
const isPlaying = computed(() => player.isPlaying.value)
const isBuffering = computed(() => player.isBuffering.value)
const isLoadingTrack = computed(() => player.isLoadingTrack.value)
const title = computed(() => currentTrack.value?.title ?? 'No track')
const artistName = computed(() => currentTrack.value?.artistName ?? 'Unknown')
const coverUrl = computed(() => currentTrack.value?.coverUrl)
const currentTime = computed(() => player.currentTime.value)
const duration = computed(() => player.duration.value)
const volume = computed(() => player.volume.value)
const muted = computed(() => player.muted.value)
const shuffleMode = computed(() => player.shuffleMode.value)
const repeatMode = computed(() => player.repeatMode.value)
const playbackRate = computed(() => player.playbackRate.value)
const hasNext = computed(() => player.hasNext.value)
const hasPrevious = computed(() => player.hasPrevious.value)

const { palette } = useAlbumColors(coverUrl)

const accentColor = computed(() => palette.value.vibrant || '#1db954')

const activeTab = ref<'now-playing' | 'queue'>('now-playing')

const tabs = [
  { key: 'now-playing' as const, label: 'Now Playing' },
  { key: 'queue' as const, label: 'Queue' },
]

const tabIndicatorStyle = computed(() => {
  const idx = tabs.findIndex((t) => t.key === activeTab.value)
  const width = 100 / tabs.length
  return {
    width: `${width}%`,
    left: `${idx * width}%`,
  }
})

const rootEl = ref<HTMLElement | null>(null)
const progressRef = ref<HTMLDivElement | null>(null)
const hoverPos = ref<number | null>(null)
const showOverflow = ref(false)
const isLiked = ref(false)

// Locally scoped player methods for template
const toggleShuffle = player.toggleShuffle
const toggleRepeat = player.toggleRepeat
const playPrevious = player.playPrevious
const playNext = player.playNext
const toggleMute = player.toggleMute

const queueTracks = computed(() => player.queue.value as PlaybackTrack[])

const queueIndex = computed(() => {
  if (!currentTrack.value) return -1
  return queueTracks.value.findIndex((t) => t.id === currentTrack.value?.id)
})

const upNextTracks = computed(() => {
  const idx = queueIndex.value
  if (idx < 0) return queueTracks.value.slice(0, 3)
  return queueTracks.value.slice(idx + 1, idx + 4)
})

const progressPct = computed(() => {
  if (!duration.value) return 0
  return (currentTime.value / duration.value) * 100
})

function formatTime(s?: number | null): string {
  if (!s) return '0:00'
  const total = Math.max(0, Math.floor(Number(s) || 0))
  const m = Math.floor(total / 60)
  const sec = total % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}

const currentTimeLabel = computed(() => formatTime(currentTime.value))
const durationLabel = computed(() =>
  formatTime(duration.value || currentTrack.value?.durationSeconds || 0),
)
const hoverTimeLabel = computed(() => {
  if (hoverPos.value === null || !duration.value) return '0:00'
  return formatTime((hoverPos.value / 100) * duration.value)
})

const repeatTitle = computed(() => {
  if (repeatMode.value === 'off') return 'Repeat: off'
  if (repeatMode.value === 'all') return 'Repeat: all'
  return 'Repeat: one'
})

const volumeIcon = computed(() => {
  if (muted.value || volume.value === 0) return 'pi pi-volume-off'
  if (volume.value < 0.5) return 'pi pi-volume-down'
  return 'pi pi-volume-up'
})

const speedLabel = computed(() => `${playbackRate.value}x`)

function togglePlayPause() {
  if (isPlaying.value) {
    player.pause()
  } else {
    player.resume()
  }
}

function seekFromEvent(e: MouseEvent) {
  const rect = progressRef.value?.getBoundingClientRect()
  if (!rect) return
  const pct = (e.clientX - rect.left) / rect.width
  if (duration.value) {
    player.seek(pct * duration.value)
  }
}

function onProgressHover(e: MouseEvent) {
  const rect = progressRef.value?.getBoundingClientRect()
  if (!rect) return
  hoverPos.value = ((e.clientX - rect.left) / rect.width) * 100
}

function onVolume(e: Event) {
  player.setVolume(Number((e.target as HTMLInputElement).value))
}

function playQueueItem(index: number) {
  const track = queueTracks.value[index]
  if (!track || track.id === currentTrack.value?.id) return
  player.setQueueAndPlay(queueTracks.value, index)
}

function removeFromQueue(index: number) {
  const newQueue = [...player.queue.value]
  newQueue.splice(index, 1)
  player.updateQueue(newQueue)
}

const speedOptions = [0.5, 0.75, 1, 1.25, 1.5, 2]
function cycleSpeed() {
  const idx = speedOptions.indexOf(playbackRate.value)
  const nextIdx = (idx + 1) % speedOptions.length
  player.setPlaybackRate(speedOptions[nextIdx]!)
}

function close() {
  visible.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
  if (e.key === ' ') { e.preventDefault(); togglePlayPause() }
  if (e.key === 'ArrowLeft') player.playPrevious()
  if (e.key === 'ArrowRight') player.playNext()
}

function parseHexColor(hex: string): string {
  const clean = hex.replace('#', '')
  if (clean.length === 6) {
    const r = parseInt(clean.substring(0, 2), 16)
    const g = parseInt(clean.substring(2, 4), 16)
    const b = parseInt(clean.substring(4, 6), 16)
    if (!isNaN(r) && !isNaN(g) && !isNaN(b)) return `${r},${g},${b}`
  }
  return '29,185,84'
}

const auroraGradient = computed(() => {
  const p = accentColor.value || '#1db954'
  if (!coverUrl.value) return {}
  return {
    background: `
      radial-gradient(ellipse 80% 60% at 50% 0%,
        rgba(${parseHexColor(p)}, 0.25) 0%,
        transparent 70%)
    `,
    transition: 'background 1s ease',
  }
})

// Override bgStyle from useAlbumColors to use our aurora spec
const resolvedBgStyle = computed(() => {
  if (!coverUrl.value) {
    return { background: '#0a0a0a' }
  }
  return {
    background:
      `radial-gradient(ellipse 80% 60% at 50% 0%, rgba(${parseHexColor(accentColor.value || '#1db954')}, 0.25) 0%, transparent 70%), #0a0a0a`,
    transition: 'background 1s ease',
  }
})

onMounted(() => {
  rootEl.value?.focus()
  window.addEventListener('popstate', close)
})

onBeforeUnmount(() => {
  window.removeEventListener('popstate', close)
})
</script>

<style scoped>
.fullscreen-slide-enter-active {
  transition: transform 400ms cubic-bezier(0.19, 1, 0.22, 1);
}
.fullscreen-slide-leave-active {
  transition: transform 300ms ease-in;
}
.fullscreen-slide-enter-from {
  transform: translateY(100%);
}
.fullscreen-slide-leave-to {
  transform: translateY(100%);
}

@keyframes art-pop {
  from { transform: scale(0.97); opacity: 0.7; }
  to { transform: scale(1); opacity: 1; }
}

@media (prefers-reduced-motion: reduce) {
  .fullscreen-slide-enter-active,
  .fullscreen-slide-leave-active {
    transition: none;
  }
  .fullscreen-slide-enter-from,
  .fullscreen-slide-leave-to {
    transform: none;
  }
}
</style>
