<template>
  <aside
    class="hidden h-screen shrink-0 border-r border-white/10 bg-black/20 backdrop-blur-2xl xl:flex xl:flex-col transition-all duration-300 ease-out z-30"
    :class="collapsed ? 'w-14 items-center' : 'w-88'"
    style="backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);"
  >
    <div class="flex h-full flex-col overflow-hidden" :class="collapsed ? 'items-center' : ''">
      <!-- ── Pane Header / Collapse Toggle ── -->
      <div class="flex shrink-0 items-center px-5 py-4" :class="collapsed ? 'flex-col gap-4' : 'gap-3 border-b border-white/5 w-full'">
        <button
          type="button"
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg transition hover:bg-white/10"
          :class="currentTrack ? 'bg-spotify/20' : 'bg-white/5'"
          :title="collapsed ? 'Show now playing panel' : 'Hide now playing panel'"
          @click="toggleCollapsed"
        >
          <i aria-hidden="true" :class="[collapsed ? 'pi pi-chevron-left' : 'pi pi-chevron-right', currentTrack ? 'text-spotify' : 'text-slate-500', 'pi text-xs']" />
        </button>
        <template v-if="!collapsed">
          <div class="flex h-6 w-6 items-center justify-center rounded-lg bg-spotify/20">
            <i aria-hidden="true" class="pi pi-waveform text-[10px] text-spotify" />
          </div>
          <h2 class="text-xs font-bold tracking-[0.2em] text-slate-400 uppercase">Now Playing</h2>
        </template>
        <template v-else>
          <div v-if="currentTrack" class="flex flex-col items-center gap-2">
            <div class="h-10 w-10 overflow-hidden rounded-xl ring-1 ring-white/10">
              <img v-if="currentTrack.coverUrl" :src="currentTrack.coverUrl" :alt="currentTrack.title" class="h-full w-full object-cover" />
              <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
                <i aria-hidden="true" class="pi pi-headphones text-sm text-slate-500" />
              </div>
            </div>
            <span v-if="isPlaying" class="h-1.5 w-1.5 rounded-full bg-spotify animate-pulse" />
          </div>
        </template>
      </div>

      <!-- ── Scrollable Content (hidden when collapsed) ── -->
      <template v-if="!collapsed">
      <div class="pane-scroll flex-1 space-y-3 overflow-y-auto p-4 pb-32 scroll-smooth" style="scrollbar-width: thin; scrollbar-color: rgba(255,255,255,0.08) transparent;">
        <!-- ─── Now Playing Card ─── -->
        <div v-if="currentTrack" class="group relative overflow-hidden rounded-2xl border border-white/10 bg-white/4 transition-all duration-300 hover:border-white/20 hover:bg-white/6">
          <!-- Background blur -->
          <div v-if="currentTrack.coverUrl" class="pointer-events-none absolute inset-0 scale-110" aria-hidden="true">
            <img
              :src="currentTrack.coverUrl"
              alt=""
              class="h-full w-full object-cover opacity-[0.15]"
              style="filter: blur(40px) saturate(1.5)"
            />
          </div>

          <div class="relative z-10 p-4">
            <!-- Album Art + Info -->
            <div class="flex items-start gap-4">
              <div class="relative shrink-0">
                <div
                  class="h-20 w-20 overflow-hidden rounded-2xl shadow-[0_8px_24px_rgba(0,0,0,0.4)] ring-1 ring-white/10 transition-all duration-700"
                  :class="isPlaying ? 'scale-100' : 'scale-95 opacity-80'"
                >
                  <img
                    v-if="currentTrack.coverUrl"
                    :src="currentTrack.coverUrl"
                    :alt="currentTrack.title"
                    class="h-full w-full object-cover"
                    loading="lazy"
                  />
                  <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
                    <i aria-hidden="true" class="pi pi-headphones text-xl text-slate-500" />
                  </div>
                </div>
                <!-- Now Playing Badge -->
                <div
                  v-if="isPlaying"
                  class="absolute -bottom-1 -right-1 flex h-6 w-6 items-center justify-center rounded-full bg-spotify shadow-lg"
                >
                  <i aria-hidden="true" class="pi pi-waveform text-xs text-black" />
                </div>
              </div>

              <div class="min-w-0 flex-1 pt-1">
                <p class="truncate text-base font-bold text-white leading-tight">{{ currentTrack.title }}</p>
                <p class="mt-0.5 truncate text-sm text-slate-400">{{ currentTrack.artistName }}</p>
                <p v-if="currentTrack.albumTitle" class="mt-0.5 truncate text-xs text-slate-500">{{ currentTrack.albumTitle }}</p>

                <!-- Mini Controls -->
                <div class="mt-3 flex items-center gap-1">
                  <button
                    type="button"
                    aria-label="Like"
                    class="flex h-8 w-8 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
                    @click="toggleLike"
                  >
                    <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill text-pink-400' : 'pi pi-heart'" class="text-xs" />
                  </button>

                  <button
                    type="button"
                    aria-label="Fullscreen"
                    class="flex h-8 w-8 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
                    @click="$emit('toggle-fullscreen')"
                  >
                    <i aria-hidden="true" class="pi pi-arrow-up-right-and-arrow-down-left-from-center text-xs" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Progress -->
            <div class="mt-4">
              <div role="button" tabindex="0" class="group relative h-1 cursor-pointer rounded-full transition-all hover:h-1.5" :class="isPlaying ? 'bg-white/15' : 'bg-white/10'" @click="onSeekClick" @keydown.enter="onSeekClick" @keydown.space.prevent="onSeekClick">
                <!-- Base glow -->
                <div v-if="isPlaying" class="absolute inset-0 rounded-full opacity-20 blur-xs transition-opacity" :style="{ background: progressColor }" />
                <!-- Fill -->
                <div class="relative h-full rounded-full" :style="{ width: `${progressPercent}%`, background: progressColor, boxShadow: progressGlow }" />
              </div>
              <div class="mt-1.5 flex items-center justify-between">
                <span class="text-[10px] font-mono tabular-nums text-slate-500">{{ formatTime(currentTime) }}</span>
                <span class="text-[10px] font-mono tabular-nums text-slate-500">{{ formatTime(duration) }}</span>
              </div>
            </div>

            <!-- Playback Controls -->
            <div class="mt-3 flex items-center justify-center gap-2">
              <button
                type="button"
                aria-label="Shuffle"
                class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
                @click="toggleShuffle"
                :class="shuffleMode !== 'off' ? 'text-aurora-purple' : ''"
              >
                <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
              </button>
              <button
                type="button"
                aria-label="Previous track"
                class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white disabled:opacity-30"
                :disabled="hasPrevious!"
                @click="playPrevious"
              >
                <i aria-hidden="true" class="pi pi-step-backward text-sm" />
              </button>
              <button
                type="button"
                aria-label="Play or pause"
                class="flex h-10 w-10 items-center justify-center rounded-full bg-white/90 text-black shadow-lg transition hover:bg-white active:scale-95 disabled:opacity-40"
                :disabled="currentTrack!"
                @click="togglePlayPause"
              >
                <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
                <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="ml-0.5 text-sm" />
              </button>
              <button
                type="button"
                aria-label="Next track"
                class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white disabled:opacity-30"
                :disabled="hasNext!"
                @click="playNext"
              >
                <i aria-hidden="true" class="pi pi-step-forward text-sm" />
              </button>
              <button
                type="button"
                aria-label="Repeat"
                class="flex h-9 w-9 items-center justify-center rounded-full transition"
                :class="repeatMode !== 'off' ? 'text-aurora-pink' : 'text-slate-400 hover:bg-white/10 hover:text-white'"
                @click="toggleRepeat"
              >
                <i aria-hidden="true" class="pi pi-refresh text-sm" />
              </button>
            </div>
          </div>
        </div>

        <!-- Empty Now Playing -->
        <div v-else class="flex flex-col items-center justify-center gap-3 rounded-2xl border border-dashed border-white/10 bg-white/2 px-6 py-12 text-center">
          <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5">
            <i aria-hidden="true" class="pi pi-music text-2xl text-slate-500" />
          </div>
          <div>
            <p class="text-sm font-bold text-slate-400">No track playing</p>
            <p class="mt-1 text-xs text-slate-500">Select a track to start listening</p>
          </div>
        </div>

        <!-- ─── Queue Section ─── -->
        <div v-if="queue.length > 0" class="overflow-hidden rounded-2xl border border-white/10 bg-white/4">
          <div class="flex items-center justify-between border-b border-white/5 px-4 py-3">
            <h3 class="flex items-center gap-2 text-xs font-bold tracking-[0.15em] text-slate-500 uppercase">
              <i aria-hidden="true" class="pi pi-list text-[10px]" />
              Up Next
              <span class="flex h-4 min-w-4 items-center justify-center rounded-full bg-white/10 px-1.5 text-[9px] font-bold text-slate-400">{{ queue.length }}</span>
            </h3>
            <button
              type="button"
              class="text-[10px] font-bold text-spotify transition hover:text-spotify-hover"
              @click="clearQueue"
            >
              Clear
            </button>
          </div>

          <div class="divide-y divide-white/3">
            <div
              v-for="(track, idx) in visibleQueue"
              :key="track.id"
              role="button"
              tabindex="0"
              class="flex cursor-pointer items-center gap-3 px-4 py-2.5 transition hover:bg-white/4"
              @click="playFromQueue(idx)"
              @keydown.enter="playFromQueue(idx)"
              @keydown.space.prevent="playFromQueue(idx)"
            >
              <div class="h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-white/10">
                <img
                  v-if="track.coverUrl"
                  :src="track.coverUrl"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-slate-500" />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p class="truncate text-xs text-slate-500">{{ track.artistName }}</p>
              </div>
              <span class="shrink-0 text-[10px] font-mono tabular-nums text-slate-500">{{ formatTime(track.durationSeconds) }}</span>
            </div>
          </div>

          <!-- Show more toggle -->
          <button
            v-if="queue.length > maxVisibleQueue"
            type="button"
            class="flex w-full items-center justify-center gap-1.5 border-t border-white/3 px-4 py-2.5 text-[11px] font-bold text-slate-400 transition hover:bg-white/4 hover:text-white"
            @click="showAllQueue = !showAllQueue"
          >
            <i aria-hidden="true" :class="showAllQueue ? 'pi pi-chevron-up' : 'pi pi-chevron-down'" class="text-[10px]" />
            {{ showAllQueue ? 'Show less' : `Show all (${queue.length})` }}
          </button>
        </div>

        <!-- ─── Suggestions Section ─── -->
        <div v-if="suggestions.length > 0" class="overflow-hidden rounded-2xl border border-white/10 bg-white/4">
          <div class="flex items-center justify-between border-b border-white/5 px-4 py-3">
            <h3 class="flex items-center gap-2 text-xs font-bold tracking-[0.15em] text-slate-500 uppercase">
              <i aria-hidden="true" class="pi pi-star text-[10px]" />
              Suggestions
            </h3>
          </div>

          <div class="divide-y divide-white/3">
            <div
              v-for="track in suggestions"
              :key="track.id"
              role="button"
              tabindex="0"
              class="flex cursor-pointer items-center gap-3 px-4 py-2.5 transition hover:bg-white/4"
              @click="playSuggestion(track)"
              @keydown.enter="playSuggestion(track)"
              @keydown.space.prevent="playSuggestion(track)"
            >
              <div class="h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-white/10">
                <img
                  v-if="track.coverUrl"
                  :src="track.coverUrl"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-slate-500" />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p class="truncate text-xs text-slate-500">{{ track.artistName }}</p>
              </div>
              <button
                type="button"
                aria-label="Play suggestion"
                class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-slate-500 transition hover:bg-white/10 hover:text-spotify"
                @click.stop="playSuggestion(track)"
              >
                <i aria-hidden="true" class="pi pi-play text-xs" />
              </button>
            </div>
          </div>
        </div>

        <!-- Player bar spacer -->
        <div class="h-2" />

        <!-- ─── Quick Actions ─── -->
        <div class="rounded-2xl border border-white/10 bg-black/60 px-4 py-3 backdrop-blur-xl">
          <div class="flex items-center justify-between">
            <button
              type="button"
              class="flex items-center gap-2 rounded-xl px-3 py-2 text-xs font-bold text-slate-400 transition hover:bg-white/10 hover:text-white"
              @click="$emit('toggle-queue-overlay')"
            >
              <i aria-hidden="true" class="pi pi-ellipsis-v text-xs" />
              <span>Full Queue</span>
            </button>
            <div class="flex items-center gap-1">
              <button
                type="button"
                aria-label="Toggle mute"
                class="flex h-7 w-7 items-center justify-center rounded-full text-slate-500 transition hover:bg-white/10 hover:text-white"
                @click="toggleMute"
              >
                <i aria-hidden="true" :class="volumeIcon" class="text-xs" />
              </button>
              <div class="w-16" dir="ltr">
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  class="player-range h-1 w-full cursor-pointer appearance-none rounded-full bg-white/10 outline-hidden"
                  :value="muted ? 0 : volume"
                  @input="onVolume"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      </template>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { usePlayerStore } from '@/stores/player'
import { usePlayerControls, useTrackLike } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { queueManager } from '@/services/player/queue-manager'
import { onImgError } from '@/utils/helpers'
import type { RecommendationTrack } from '@/services/api/recommendation/types'

const playerStore = usePlayerStore()
const pc = usePlayerControls()

const collapsed = ref(localStorage.getItem('rightpane-collapsed') === 'true')
function toggleCollapsed() {
  collapsed.value = collapsed.value!
  localStorage.setItem('rightpane-collapsed', String(collapsed.value))
  window.dispatchEvent(new CustomEvent('rightpane-collapse', { detail: collapsed.value }))
}

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  volume,
  muted,
  currentTime,
  duration,
  progressPercent,
  shuffleMode,
  repeatMode,
  volumeIcon,
  hasNext,
  hasPrevious,
  togglePlayPause,
  toggleMute,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
} = pc

const queue = computed(() => playerStore.queue)
const trackId = computed(() => currentTrack.value?.id)
const { liked, toggleLike } = useTrackLike(trackId)

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressColor = computed(() => {
  return coverUrl.value && palette.value.vibrant ? palette.value.vibrant : '#1db954'
})

const progressGlow = computed(() => {
  const c = progressColor.value
  return `0 0 6px ${c}44, 0 0 12px ${c}22`
})

const showAllQueue = ref(false)
const maxVisibleQueue = 5

const visibleQueue = computed(() => {
  if (showAllQueue.value) return queue.value
  return queue.value.slice(0, maxVisibleQueue)
})

interface SuggestionItem {
  id: string
  title: string
  artistName: string
  coverUrl?: string | null
  durationSeconds?: number | null
}

const suggestions = ref<SuggestionItem[]>([])

// Fetch suggestions when a track is playing
watch(currentTrack, async (track) => {
  suggestions.value = []
  if (track?.id) {
    try {
      const recsApi = useRecommendationsApi()
      const result = await recsApi.getSimilar(track.id, { limit: 3 })
      const items = result?.items ?? []
      suggestions.value = items.map((item: RecommendationTrack) => ({
        id: item.id,
        title: item.title,
        artistName: item.artist_name || 'Unknown',
        coverUrl: item.cover_url,
        durationSeconds: item.duration_seconds,
      }))
    } catch (err) {
      console.error('Failed to fetch suggestions:', err)
    }
  }
}, { immediate: true })

function playFromQueue(index: number) {
  const track = queue.value[index]
  if (track) {
    playerStore.playTrack(track)
  }
}

function playSuggestion(track: SuggestionItem) {
  playerStore.playTrackById(track.id)
}

function clearQueue() {
  queueManager.clear()
  playerStore.updateQueue([])
}

function onSeekClick(e: MouseEvent | KeyboardEvent) {
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = (((e as MouseEvent).clientX - rect.left) / rect.width) * 100
  pc.seekPercent(pct)
}

function onVolume(e: Event) {
  pc.setVolume(Number((e.target as HTMLInputElement).value))
}

function formatTime(seconds?: number | null) {
  if (seconds! || isFinite!(seconds)) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>

<style scoped>
/* Scroll container bottom fade */
.pane-scroll {
  position: relative;
}
.pane-scroll::after {
  content: '';
  position: sticky;
  bottom: 0;
  left: 0;
  right: 0;
  display: block;
  height: 2rem;
  background: linear-gradient(to top, rgba(0,0,0,0.4), transparent);
  pointer-events: none;
  z-index: 5;
  margin-top: -2rem;
}

.player-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: white;
  border: none;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s;
}
.player-range:hover::-webkit-slider-thumb {
  opacity: 1;
}
.player-range::-moz-range-thumb {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: white;
  border: none;
  cursor: pointer;
  opacity: 0;
}
.player-range:hover::-moz-range-thumb {
  opacity: 1;
}
</style>
