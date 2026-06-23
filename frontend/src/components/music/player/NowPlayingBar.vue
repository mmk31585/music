<template>
  <!-- DESKTOP BAR -->
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      ref="barRef"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 bottom-0 z-50 hidden md:block select-none"
    >
      <div
        class="relative flex flex-col overflow-visible transition-all duration-300 ease-out"
        :class="collapsed ? 'pb-0' : ''"
        style="background: #08080A; backdrop-filter: blur(32px); -webkit-backdrop-filter: blur(32px);"
      >
        <!-- Background blur gets its own overflow clip so it doesn't clip popup menus -->
        <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-[inherit]">
          <div
            v-if="currentTrack.coverUrl"
            class="absolute inset-0 scale-110"
            aria-hidden="true"
          >
            <img
              :src="currentTrack.coverUrl"
              alt=""
              class="h-full w-full object-cover"
              :class="collapsed ? 'opacity-[0.15]' : 'opacity-[0.35]'"
              style="filter: blur(60px) saturate(1.5)"
            />
          </div>
          <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-black/80 via-black/50 to-transparent" />
        </div>

        <!-- ── MINI COLLAPSED BAR ── -->
        <!-- Thin progress line at top -->
        <div class="relative z-10 h-0.5 bg-white/5">
          <div class="h-full rounded-full transition-[width] duration-100" :style="progressStyle" />
        </div>
        <div v-if="collapsed" class="relative z-10 flex items-center gap-3 px-4 h-14">
          <div role="button" tabindex="0" aria-label="Open fullscreen player" class="relative shrink-0 cursor-pointer" @click="emit('toggle-fullscreen')" @keydown.enter="emit('toggle-fullscreen')" @keydown.space.prevent="emit('toggle-fullscreen')">
            <div class="h-10 w-10 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10">
              <img
                v-if="currentTrack.coverUrl"
                :src="currentTrack.coverUrl"
                :alt="currentTrack.title"
                class="h-full w-full object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
                <i aria-hidden="true" class="pi pi-headphones text-sm text-slate-500" />
              </div>
            </div>
          </div>

          <div role="button" tabindex="0" aria-label="Open fullscreen player" class="min-w-0 flex-1 cursor-pointer" @click="emit('toggle-fullscreen')" @keydown.enter="emit('toggle-fullscreen')" @keydown.space.prevent="emit('toggle-fullscreen')">
            <p class="truncate text-sm font-bold text-white leading-tight">{{ currentTrack.title }}</p>
            <p class="truncate text-xs text-white/50 leading-tight">{{ currentTrack.artistName }}</p>
          </div>

          <div class="flex items-center gap-2">
            <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full shadow-lg disabled:opacity-40 hover:scale-110 active:scale-90 transition-all duration-200"
                :style="{ background: progressColor }"
                :disabled="currentTrack! || isLoadingTrack"
                :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
                @click="isPlaying ? togglePlayPause() : proceed()"
              >
                <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm text-white" />
                <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="ml-0.5 text-sm text-white" />
              </button>
            </GuestPlayGate>

            <button
              type="button"
              class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              :aria-label="collapsed ? 'Expand player' : 'Collapse player'"
              @click="toggleCollapsed"
            >
              <i aria-hidden="true" :class="collapsed ? 'pi pi-chevron-up' : 'pi pi-chevron-down'" class="text-sm" />
            </button>
          </div>
        </div>

        <!-- ── FULL EXPANDED BAR ── -->
        <template v-if="!collapsed">
          <div class="relative z-10 flex items-center gap-4 px-6 pt-3">
            <div class="flex min-w-0 w-[25%] items-center gap-3">
              <div role="button" tabindex="0" aria-label="Open fullscreen player" class="relative shrink-0 cursor-pointer" @click="emit('toggle-fullscreen')" @keydown.enter="emit('toggle-fullscreen')" @keydown.space.prevent="emit('toggle-fullscreen')">
                <div
                  class="h-14 w-14 overflow-hidden rounded-[18px] shadow-[0_16px_32px_rgba(0,0,0,0.5)] ring-1 ring-white/10 transition-all duration-700"
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
                    <i aria-hidden="true" class="pi pi-headphones text-lg text-slate-500" />
                  </div>
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <p class="truncate text-sm font-bold text-white leading-tight">{{ currentTrack.title }}</p>
                  <button
                    type="button"
                    class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90"
                    :class="liked ? 'text-aurora-pink' : 'text-white/30 hover:text-white/60'"
                    :aria-label="liked ? 'Unlike' : 'Like'"
                    @click.stop="toggleLike"
                  >
                    <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xs" />
                  </button>
                  <button
                    type="button"
                    class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90 text-white/30 hover:text-white/60"
                    aria-label="Add to playlist"
                    @click.stop="showAddToPlaylist = true"
                  >
                    <i aria-hidden="true" class="pi pi-plus text-xs" />
                  </button>
                </div>
                <p class="truncate text-xs text-white/50 mt-0.5 leading-tight">{{ currentTrack.artistName }}</p>
                <div v-if="isPlaying" class="flex items-center gap-1 mt-1">
                  <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 0ms" />
                  <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 150ms" />
                  <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 300ms" />
                </div>
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center gap-1.5">
              <div class="relative">
                <button
                  ref="shuffleBtnRef"
                  type="button"
                  class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
                  :class="shuffleMode !== 'off' ? 'text-aurora-purple' : 'text-white/40 hover:text-white hover:bg-white/10'"
                  :disabled="currentTrack!"
                  :aria-label="shuffleMode === 'queue' ? 'Shuffle queue' : shuffleMode === 'catalog' ? 'Random catalog tracks' : shuffleMode === 'similar' ? 'Similar tracks' : 'Shuffle off'"
                  @click="showShuffleMenu = !showShuffleMenu"
                >
                  <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
                  <span
                    v-if="shuffleMode !== 'off'"
                    class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-aurora-purple text-[8px] font-bold text-white"
                  >{{ shuffleMode === 'queue' ? 'Q' : shuffleMode === 'catalog' ? 'R' : 'S' }}</span>
                </button>

                <!-- Shuffle mode selector popup -->
                <Transition name="fade">
                  <div
                    v-if="showShuffleMenu"
                    class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-[60] min-w-[150px] rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
                    style="backdrop-filter: blur(24px);"
                  >
                    <button
                      v-for="mode in shuffleModes"
                      :key="mode.value"
                      type="button"
                      class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8"
                      :class="shuffleMode === mode.value ? 'text-aurora-purple bg-white/6' : 'text-slate-400 hover:text-white'"
                      @click="setShuffleMode(mode.value)"
                    >
                      <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                      <span class="flex-1 text-left">{{ mode.label }}</span>
                      <span v-if="shuffleMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-purple" />
                    </button>
                  </div>
                </Transition>
              </div>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 transition-all duration-200"
                :disabled="hasPrevious!"
                aria-label="Previous track"
                @click="playPrevious"
              >
                <i aria-hidden="true" class="pi pi-step-backward text-base" />
              </button>

              <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
                <button
                  type="button"
                  class="flex h-11 w-11 items-center justify-center rounded-full shadow-xl disabled:opacity-40 hover:scale-110 active:scale-95 transition-all duration-200"
                  :style="{ background: progressColor }"
                  :disabled="currentTrack! || isLoadingTrack"
                  :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
                  @click="isPlaying ? togglePlayPause() : proceed()"
                >
                  <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-base text-white" />
                  <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="ml-0.5 text-base text-white" />
                </button>
              </GuestPlayGate>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 transition-all duration-200"
                :disabled="hasNext!"
                aria-label="Next track"
                @click="playNext"
              >
                <i aria-hidden="true" class="pi pi-step-forward text-base" />
              </button>

              <div class="relative">
                <button
                  ref="repeatBtnRef"
                  type="button"
                  class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
                  :class="repeatMode !== 'off' ? 'text-aurora-pink' : 'text-white/40 hover:text-white hover:bg-white/10'"
                  :disabled="currentTrack!"
                  :aria-label="repeatMode === 'off' ? 'Repeat off' : repeatMode === 'all' ? 'Repeat all' : 'Repeat one'"
                  @click="showRepeatMenu = !showRepeatMenu"
                >
                  <i aria-hidden="true" class="pi pi-refresh text-sm" />
                  <span
                    v-if="repeatMode === 'one'"
                    class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-aurora-pink text-[8px] font-bold text-black"
                  >1</span>
                </button>

                <!-- Repeat mode selector popup -->
                <Transition name="fade">
                  <div
                    v-if="showRepeatMenu"
                    class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-[60] min-w-[130px] rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
                    style="backdrop-filter: blur(24px);"
                  >
                    <button
                      v-for="mode in repeatModes"
                      :key="mode.value"
                      type="button"
                      class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8"
                      :class="repeatMode === mode.value ? 'text-aurora-pink bg-white/6' : 'text-slate-400 hover:text-white'"
                      @click="setRepeatMode(mode.value)"
                    >
                      <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                      <span class="flex-1 text-left">{{ mode.label }}</span>
                      <span v-if="repeatMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-pink" />
                    </button>
                  </div>
                </Transition>
              </div>
            </div>

            <div class="flex w-[25%] items-center justify-end gap-1">
              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :disabled="currentTrack!"
                aria-label="Open queue"
                @click="emit('toggle-queue')"
              >
                <i aria-hidden="true" class="pi pi-list text-sm" />
              </button>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :disabled="currentTrack!"
                aria-label="Toggle lyrics"
                @click="emit('toggle-lyrics')"
              >
                <i aria-hidden="true" class="pi pi-align-left text-sm" />
              </button>

              <div class="mx-1 h-6 w-px bg-white/10" />

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :aria-label="muted ? 'Unmute' : 'Mute'"
                @click="toggleMute"
              >
                <i aria-hidden="true" :class="volumeIcon" class="text-sm" />
              </button>
              <div class="w-20" dir="ltr">
                <Slider
                  :model-value="muted ? 0 : volume"
                  @update:model-value="onVolume"
                  :min="0"
                  :max="1"
                  :step="0.01"
                  aria-label="Volume"
                />
              </div>

              <div class="mx-1 h-6 w-px bg-white/10" />

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :disabled="currentTrack!"
                aria-label="Fullscreen"
                @click="emit('toggle-fullscreen')"
              >
                <i aria-hidden="true" class="pi pi-arrow-up-right-and-arrow-down-left-from-center text-sm" />
              </button>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :aria-label="collapsed ? 'Expand player' : 'Collapse player'"
                @click="toggleCollapsed"
              >
                <i aria-hidden="true" :class="collapsed ? 'pi pi-chevron-up' : 'pi pi-chevron-down'" class="text-sm" />
              </button>

              <div class="relative overflow-menu-container">
                <button
                  type="button"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                  :disabled="currentTrack!"
                  aria-label="More options"
                  @click.stop="showOverflow = !showOverflow"
                >
                  <i aria-hidden="true" class="pi pi-ellipsis-h text-sm" />
                </button>
                <PlayerOverflowMenu v-if="showOverflow" @close="showOverflow = false" />
              </div>
            </div>
          </div>

          <div class="relative z-10 flex items-center gap-3 px-6 pb-3 pt-1">
            <span class="text-[11px] text-white/50 font-mono tabular-nums w-10 text-right">{{ formatTime(currentTime) }}</span>
            <div
              role="button"
              tabindex="0"
              class="group relative flex-1 h-1.5 cursor-pointer rounded-full transition-all duration-150"
              :class="isPlaying ? 'bg-white/15' : 'bg-white/10'"
              dir="ltr"
              @click="onSeekClick"
              @keydown.enter="onSeekClick"
              @keydown.space.prevent="onSeekClick"
            >
              <!-- Base track glow -->
              <div
                v-if="isPlaying"
                class="absolute inset-0 rounded-full opacity-20 blur-xs transition-opacity duration-500"
                :style="{ background: progressColor }"
              />
              <!-- Fill -->
              <div
                class="relative h-full rounded-full"
                :class="isPlaying ? 'progress-bar-fill' : ''"
                :style="progressStyle"
              />
              <!-- Thumb dot -->
              <div
                class="absolute top-1/2 -translate-y-1/2 size-4 rounded-full shadow-xl opacity-0 group-hover:opacity-100 scale-0 group-hover:scale-100 transition-all duration-300 ease-spring"
                :style="{
                  left: `calc(${progressPercent}% - 8px)`,
                  background: progressColor,
                  boxShadow: `0 0 0 3px ${progressColor}22, 0 4px 12px rgba(0,0,0,0.5)`,
                }"
              />
            </div>
            <span class="text-[11px] text-white/50 font-mono tabular-nums w-10">{{ formatTime(duration) }}</span>
          </div>

          <div
            v-if="playbackError"
            class="relative z-10 flex items-center justify-center gap-2 bg-red-500/10 px-5 py-1.5 text-xs text-red-400"
          >
            <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
            <span>{{ playbackError }}</span>
          </div>
        </template>
      </div>
    </div>
  </Transition>

  <!-- MOBILE BAR - floats above bottom nav -->
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      key="mobile-bar"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 z-[45] md:hidden select-none"
      style="bottom: calc(4.5rem + env(safe-area-inset-bottom, 0px))"
    >
      <button
        type="button"
        class="relative mx-3 flex w-[calc(100%-1.5rem)] items-center gap-3 overflow-hidden rounded-2xl text-left shadow-2xl transition-all duration-500 active:scale-[0.98]"
        :style="{ background: coverUrl ? `${palette.dark}dd` : 'rgba(8, 8, 10, 0.92)', backdropFilter: 'blur(24px)', WebkitBackdropFilter: 'blur(24px)' }"
        @click="emit('toggle-fullscreen')"
      >
        <div
          v-if="currentTrack.coverUrl"
          class="absolute inset-0 scale-110 bg-cover bg-center blur-2xl opacity-[0.25] transition-all duration-700"
          :style="{ backgroundImage: `url(${currentTrack.coverUrl})` }"
          aria-hidden="true"
        />
        <div class="absolute inset-0 bg-linear-to-r from-black/80 via-black/50 to-black/80" />

        <div class="relative shrink-0">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10">
            <img
              v-if="currentTrack.coverUrl"
              :src="currentTrack.coverUrl"
              :alt="currentTrack.title"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
              <i aria-hidden="true" class="pi pi-headphones text-sm text-slate-500" />
            </div>
          </div>
        </div>

        <div class="relative min-w-0 flex-1">
          <p class="truncate text-sm font-bold text-white">{{ currentTrack.title }}</p>
          <p class="truncate text-xs text-white/50">{{ currentTrack.artistName }}</p>
        </div>

        <div class="relative flex items-center gap-0.5 pr-1" @click.stop>
          <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-black shadow-lg disabled:opacity-40 active:scale-90 transition-transform"
              :disabled="currentTrack! || isLoadingTrack"
              :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
              @click="isPlaying ? togglePlayPause() : proceed()"
            >
              <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
              <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="ml-0.5 text-sm" />
            </button>
          </GuestPlayGate>
          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-white/60 hover:text-white active:scale-90 transition-transform"
            :disabled="hasNext!"
            aria-label="Next track"
            @click.stop="playNext"
          >
            <i aria-hidden="true" class="pi pi-step-forward text-sm" />
          </button>
        </div>
      </button>

      <div
        v-if="playbackError"
        class="mx-3 mt-1 flex items-center justify-center gap-2 rounded-xl bg-red-500/10 px-4 py-1.5 text-xs text-red-400"
      >
        <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
        <span>{{ playbackError }}</span>
      </div>
    </div>
  </Transition>

  <AddToPlaylistDialog
    :visible="showAddToPlaylist"
    :track-id="currentTrack?.id || ''"
    :track-title="currentTrack?.title"
    @update:visible="showAddToPlaylist = $event"
  />
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, onUnmounted, ref } from 'vue'

import { usePlayerControls, useTrackLike } from '@/composables/player'
import { usePlayerStore } from '@/stores/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import AddToPlaylistDialog from './AddToPlaylistDialog.vue'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'
import GuestPlayGate from '@/components/common/GuestPlayGate.vue'

const emit = defineEmits<{
  'toggle-queue': []
  'toggle-fullscreen': []
  'toggle-lyrics': []
  'toggle-mobile-sheet': []
}>()

const showOverflow = ref(false)
const showAddToPlaylist = ref(false)
const collapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')
const showShuffleMenu = ref(false)
const shuffleBtnRef = ref<HTMLElement | null>(null)
const showRepeatMenu = ref(false)
const repeatBtnRef = ref<HTMLElement | null>(null)

const shuffleModes = [
  { value: 'off' as const, label: 'Off', icon: 'pi pi-ban' },
  { value: 'queue' as const, label: 'Shuffle Queue', icon: 'pi pi-sort-alt' },
  { value: 'catalog' as const, label: 'Random Catalog', icon: 'pi pi-globe' },
  { value: 'similar' as const, label: 'Similar Tracks', icon: 'pi pi-star' },
]

const repeatModes = [
  { value: 'off' as const, label: 'No Repeat', icon: 'pi pi-refresh' },
  { value: 'all' as const, label: 'Repeat All', icon: 'pi pi-sync' },
  { value: 'one' as const, label: 'Repeat One', icon: 'pi pi-undo' },
]

function setShuffleMode(mode: 'off' | 'queue' | 'catalog' | 'similar') {
  pc.setShuffleMode(mode)
  showShuffleMenu.value = false
}

function setRepeatMode(mode: 'off' | 'all' | 'one') {
  usePlayerStore().repeatMode = mode
  showRepeatMenu.value = false
}

// Close popup menus on outside click
onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
})
onUnmounted(() => {
  document.removeEventListener('click', handleOutsideClick)
})
function handleOutsideClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (showShuffleMenu.value && shuffleBtnRef.value && shuffleBtnRef.value.contains!(target)) {
    showShuffleMenu.value = false
  }
  if (showRepeatMenu.value && repeatBtnRef.value && repeatBtnRef.value.contains!(target)) {
    showRepeatMenu.value = false
  }
}

function toggleCollapsed() {
  collapsed.value = collapsed.value!
  localStorage.setItem('player-bar-collapsed', String(collapsed.value))
  window.dispatchEvent(new CustomEvent('playerbar-collapse', { detail: collapsed.value }))
}

const pc = usePlayerControls()

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  volume,
  muted,
  currentTime,
  duration,
  error: playbackError,
  hasNext,
  hasPrevious,
  shuffleMode,
  repeatMode,
  volumeIcon,
  togglePlayPause,
  toggleMute,
  setVolume,
  progressPercent,
  seekPercent,
  playNext,
  playPrevious,
} = pc

const trackId = computed(() => currentTrack.value?.id)
const { liked, toggleLike } = useTrackLike(trackId)

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressColor = computed(() => {
  return coverUrl.value && palette.value.vibrant ? palette.value.vibrant : '#1db954'
})

const progressGlow = computed(() => {
  const c = progressColor.value
  return `0 0 8px ${c}66, 0 0 20px ${c}33`
})

const progressStyle = computed(() => ({
  width: `${progressPercent.value}%`,
  background: progressColor.value,
  boxShadow: progressGlow.value,
  transition: 'width 100ms linear, background 0.5s ease, box-shadow 0.3s ease',
}))

function formatTime(s: number) {
  if (s! || isFinite!(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

const barRef = ref<HTMLElement | null>(null)

function onVolume(val: number) {
  setVolume(val)
}

function onSeekClick(e: MouseEvent | KeyboardEvent) {
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = (((e as MouseEvent).clientX - rect.left) / rect.width) * 100
  seekPercent(pct)
}

function onOverflowClickOutside(e: MouseEvent) {
  if (!(e.target as HTMLElement).closest('.overflow-menu-container')) {
    showOverflow.value = false
  }
}

onMounted(() => document.addEventListener('click', onOverflowClickOutside))
onBeforeUnmount(() => document.removeEventListener('click', onOverflowClickOutside))
</script>

<style scoped>
.bar-slide-enter-active {
  transition: transform 400ms cubic-bezier(0.19, 1, 0.22, 1);
}
.bar-slide-leave-active {
  transition: transform 300ms ease-in;
}
.bar-slide-enter-from,
.bar-slide-leave-to {
  transform: translateY(100%);
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
.animate-bounce {
  animation: bounce 0.6s ease-in-out infinite;
}

/* ── Progress Bar ─────────────────────────────── */

/* Subtle pulse glow on the progress fill when playing */
@keyframes progress-glow {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}
.progress-bar-fill::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: inherit;
  filter: blur(6px);
  opacity: 0.35;
  animation: progress-glow 2s ease-in-out infinite;
  z-index: -1;
}

/* Ease-spring utility for the thumb dot */
.ease-spring {
  transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* Time labels hover highlight */
.group:hover .time-highlight {
  color: rgba(255, 255, 255, 0.7);
}

@media (prefers-reduced-motion: reduce) {
  .bar-slide-enter-active,
  .bar-slide-leave-active { transition: none; }
  .bar-slide-enter-from,
  .bar-slide-leave-to { transform: none; }
  .animate-bounce { animation: none; }
  .progress-bar-fill::after { animation: none; opacity: 0; }
}
</style>
