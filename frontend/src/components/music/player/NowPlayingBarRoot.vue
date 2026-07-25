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
        :style="barBackgroundStyle"
      >
        <!-- Background with album art blur -->
        <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-[inherit]">
          <div
            v-if="coverUrl"
            class="absolute inset-0 scale-110"
            aria-hidden="true"
          >
            <img
              :src="coverUrl"
              alt=""
              class="h-full w-full object-cover transition-opacity duration-300 ease-out"
              :class="collapsed ? 'opacity-[0.15]' : 'opacity-[0.35]'"
              style="filter: blur(60px) saturate(1.5)"
            />
          </div>
          <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-black/80 via-black/50 to-transparent" />
        </div>

        <!-- ── MINI COLLAPSED BAR ── -->
        <div class="relative z-10 h-0.5 bg-white/5">
          <div class="h-full rounded-full transition-[width] duration-100" :style="progressStyle" />
        </div>
        
        <div v-if="collapsed" class="relative z-10 flex items-center gap-3 px-4 h-14">
          <button
            type="button"
            aria-label="Open fullscreen player"
            class="relative shrink-0"
            @click="$emit('toggle-fullscreen')"
          >
            <div class="h-10 w-10 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10">
              <img
                v-if="coverUrl"
                :src="coverUrl"
                :alt="title"
                class="h-full w-full object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
                <Music class="h-5 w-5 text-slate-500" aria-hidden="true" />
              </div>
            </div>
          </button>

          <button
            type="button"
            aria-label="Open fullscreen player"
            class="min-w-0 flex-1 text-left"
            @click="$emit('toggle-fullscreen')"
          >
            <p class="truncate text-sm font-bold text-white leading-tight">{{ title }}</p>
            <p class="truncate text-xs text-white/60 leading-tight">{{ artistName }}</p>
          </button>

          <div class="flex items-center gap-2">
            <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full shadow-lg disabled:opacity-40 hover:scale-110 active:scale-90 transition-all duration-200"
                :style="{ background: progressColor }"
                :disabled="!currentTrack || isLoadingTrack"
                :aria-label="playAriaLabel"
                @click="isPlaying ? togglePlayPause() : proceed()"
              >
                <Loader v-if="isLoadingTrack || isBuffering" class="h-4 w-4 animate-spin-slow text-white" aria-hidden="true" />
                <Pause v-else-if="isPlaying" class="h-4 w-4 text-white" aria-hidden="true" />
                <Play v-else class="h-4 w-4 text-white ml-0.5" aria-hidden="true" />
              </button>
            </GuestPlayGate>

            <button
              type="button"
              class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              :aria-label="collapsed ? 'Expand player' : 'Collapse player'"
              @click="toggleCollapsed"
            >
              <ChevronUp v-if="collapsed" class="h-4 w-4" aria-hidden="true" />
              <ChevronDown v-else class="h-4 w-4" aria-hidden="true" />
            </button>
          </div>
        </div>

<!-- ── FULL EXPANDED BAR ── -->
        <Transition name="expand">
          <div v-if="!collapsed" class="overflow-hidden">
            <div class="relative z-10 flex items-center gap-4 px-6 pt-2">
              <!-- Track Info - Left Section (25%) -->
              <TrackInfo
                :title="title"
                :artist-name="artistName"
                :artist-id="(currentTrack as any)?.artistId"
                :artists="(currentTrack as any)?.artists"
                :cover-url="coverUrl || undefined"
                :is-playing="isPlaying"
                :liked="liked"
                :progress-color="progressColor"
                @open-fullscreen="$emit('toggle-fullscreen')"
                @toggle-like="toggleLike"
                @add-to-playlist="showAddToPlaylist = true"
                @navigate-artist="$emit('navigate-artist', $event)"
              />
              <div
                v-if="radio.isRadioActive.value && radio.radioSeedLabel.value"
                class="flex items-center gap-1.5 rounded-full bg-accent/15 px-2 py-0.5 text-[10px] text-accent mt-1"
              >
                <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-accent" />
                <span class="font-semibold uppercase tracking-wider">Radio</span>
                <span class="text-accent/60">· {{ radio.radioSeedLabel.value }}</span>
              </div>

              <!-- Playback Controls - Center Section (50%, flex-1) -->
              <div class="flex-1 flex justify-center">
                <PlaybackControls
                  :current-track="currentTrack"
                  :is-playing="isPlaying"
                  :is-buffering="isBuffering"
                  :is-loading-track="isLoadingTrack"
                  :has-next="hasNext"
                  :has-previous="hasPrevious"
                  :shuffle-mode="shuffleMode"
                  :repeat-mode="repeatMode"
                  :progress-color="progressColor"
                  :current-time="currentTime"
                  :duration="duration"
                  :progress-percent="progressPercent"
                  :toggle-play-pause="togglePlayPause"
                  :play-next="playNext"
                  :play-previous="playPrevious"
                  :set-shuffle-mode="setShuffleMode"
                  :set-repeat-mode="setRepeatMode"
                />
              </div>

              <!-- Right Side Actions - Right Section (25%) -->
              <div class="flex w-[25%] items-center justify-end gap-1">
                <QueuePreview
                  :current-track="currentTrack"
                  :next-track="nextTrack"
                  :upcoming-count="upcomingCount"
                  :shuffle-mode="shuffleMode"
                  :is-playing="isPlaying"
                  @play-next="playNext"
                  @open-queue="$emit('toggle-queue')"
                />

                <button
                  type="button"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
                  :disabled="!currentTrack"
                  aria-label="Toggle lyrics"
                  @click="$emit('toggle-lyrics')"
                >
                  <FileText class="h-4 w-4" aria-hidden="true" />
                </button>

                <button
                  v-if="currentTrack"
                  type="button"
                  class="flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
                  :class="radio.isRadioActive.value
                    ? 'text-accent bg-accent/10 hover:bg-accent/20'
                    : 'text-white/40 hover:text-white hover:bg-white/10'"
                  :aria-label="radio.isRadioActive.value ? 'رادیو فعال است' : 'شروع رادیو'"
                  :title="radio.isRadioActive.value && radio.radioSeedLabel.value ? 'رادیو در حال پخش · ' + radio.radioSeedLabel.value : 'شروع رادیو از این آهنگ'"
                  @click="handleRadioClick"
                >
                  <Radio aria-hidden="true" class="h-4 w-4" />
                </button>

                <div class="mx-1 h-6 w-px bg-white/10" />

                <VolumeControl
                  :volume="volume"
                  :muted="muted"
                  @volume="onVolume"
                  @mute="toggleMute"
                />

                <div class="mx-1 h-6 w-px bg-white/10" />

                <button
                  type="button"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
                  :disabled="!currentTrack"
                  aria-label="Fullscreen"
                  @click="$emit('toggle-fullscreen')"
                >
                  <Expand class="h-4 w-4" aria-hidden="true" />
                </button>

                <button
                  type="button"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
                  :aria-label="collapsed ? 'Expand player' : 'Collapse player'"
                  @click="toggleCollapsed"
                >
                  <ChevronUp v-if="collapsed" class="h-4 w-4" aria-hidden="true" />
                  <ChevronDown v-else class="h-4 w-4" aria-hidden="true" />
                </button>

                <div class="relative">
                  <button
                    type="button"
                    class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
                    :disabled="!currentTrack"
                    aria-label="More options"
                    @click.stop="showOverflow = !showOverflow"
                  >
                    <MoreHorizontal class="h-4 w-4" aria-hidden="true" />
                  </button>
                  <ContextMenu
                    v-model:visible="showOverflow"
                    :sections="overflowSections"
                    :header="overflowHeader"
                    :accent-color="progressColor"
                    :anchor-el="overflowAnchorEl"
                    @close="showOverflow = false"
                  />
                </div>
              </div>
            </div>

            <!-- Playback Error -->
            <div
              v-if="playbackError"
              class="relative z-10 flex items-center justify-center gap-2 bg-red-500/10 px-5 py-1.5 text-xs text-red-400"
            >
              <AlertCircle class="h-3 w-3" aria-hidden="true" />
              <span>{{ playbackError }}</span>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </Transition>

  <!-- MOBILE BAR -->
  <MobilePlayer
    v-if="currentTrack"
    :current-track="currentTrack"
    :is-playing="isPlaying"
    :is-buffering="isBuffering"
    :is-loading-track="isLoadingTrack"
    :has-next="hasNext"
    :playback-error="playbackError"
    :cover-url="coverUrl"
    :title="title"
    :artist-name="artistName"
    :toggle-play-pause="togglePlayPause"
    @open-fullscreen="$emit('toggle-fullscreen')"
    @play-next="playNext"
  />

  <!-- Dialogs -->
  <AddToPlaylistDialog
    :visible="showAddToPlaylist"
    :track-id="currentTrack?.id || ''"
    :track-title="currentTrack?.title"
    @update:visible="showAddToPlaylist = $event"
  />
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onBeforeUnmount, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { AlertCircle, Ban, ChevronDown, ChevronUp, Clock, Copy, Disc3, Download, Expand, FileText, Flag, Globe, Info, Loader, Monitor, MoreHorizontal, Music, Pause, Play, PlusCircle, Radio, Repeat, RotateCcw, Settings2, Share2, Shuffle, SkipForward, Sparkles, Star, UserRound } from 'lucide-vue-next'
import { useRadio } from '@/composables/recommendation/useRadio'
import { usePlayerControls, useTrackLike } from '@/composables/player'
import { usePlayerStore } from '@/stores/player'
import { usePlayerPiPController } from '@/composables/usePlayerPiPController'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useAppToast } from '@/composables/useAppToast'
import type { ContextMenuSection, ContextMenuHeader } from '@/types/context-menu'
import AddToPlaylistDialog from './AddToPlaylistDialog.vue'
import TrackInfo from './TrackInfo.vue'
import PlaybackControls from './PlaybackControls.vue'
import VolumeControl from './VolumeControl.vue'
import QueuePreview from './QueuePreview.vue'
import ContextMenu from '@/components/common/ContextMenu.vue'
import MobilePlayer from './MobilePlayer.vue'
import GuestPlayGate from '@/components/common/GuestPlayGate.vue'
import type { PlaybackTrack } from '@/services/api/player'

const toast = useAppToast()
const router = useRouter()

const emit = defineEmits<{
  'toggle-queue': []
  'toggle-fullscreen': []
  'toggle-lyrics': []
  'toggle-mobile-sheet': []
  'navigate-artist': [artistId: string]
}>()

// State
const showOverflow = ref(false)
const showAddToPlaylist = ref(false)
const collapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')
const showShuffleMenu = ref(false)
const shuffleBtnRef = ref<HTMLElement | null>(null)
const shuffleMenuRef = ref<HTMLElement | null>(null)
const shuffleItemRefs = ref<HTMLElement[]>([])
const showRepeatMenu = ref(false)
const repeatBtnRef = ref<HTMLElement | null>(null)
const repeatMenuRef = ref<HTMLElement | null>(null)
const repeatItemRefs = ref<HTMLElement[]>([])
const showQueuePreview = ref(false)
const queueBtnRef = ref<HTMLElement | null>(null)

// Menu focus management
watch(showShuffleMenu, async (v) => {
  if (v) { await nextTick(); shuffleItemRefs.value[0]?.focus() }
})
watch(showRepeatMenu, async (v) => {
  if (v) { await nextTick(); repeatItemRefs.value[0]?.focus() }
})

const shuffleModes = [
  { value: 'off' as const, label: 'Off', icon: Ban },
  { value: 'queue' as const, label: 'Shuffle Queue', icon: Shuffle },
  { value: 'catalog' as const, label: 'Random Catalog', icon: Globe },
  { value: 'similar' as const, label: 'Similar Tracks', icon: Star },
]

const repeatModes = [
  { value: 'off' as const, label: 'No Repeat', icon: RotateCcw },
  { value: 'all' as const, label: 'Repeat All', icon: Repeat },
  { value: 'one' as const, label: 'Repeat One', icon: RotateCcw },
]

function setShuffleMode(mode: 'off' | 'queue' | 'catalog' | 'similar') {
  pc.setShuffleMode(mode)
  showShuffleMenu.value = false
}

function setRepeatMode(mode: 'off' | 'all' | 'one') {
  usePlayerStore().repeatMode = mode
  showRepeatMenu.value = false
}

// Player controls
const pc = usePlayerControls()
const pip = usePlayerPiPController()

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
  progressPercent,
  seekPercent,
  playNext,
  playPrevious,
} = pc

const togglePlayPause = () => pc.togglePlayPause()
const toggleMute = () => pc.toggleMute()

const trackId = computed(() => currentTrack.value?.id)
const { liked, toggleLike } = useTrackLike(trackId)

// Queue next-up preview
const playerStore = usePlayerStore()
const queueTracks = computed(() => playerStore.queue as PlaybackTrack[])
const currentTrackIndex = computed(() => {
  if (!currentTrack.value) return -1
  return queueTracks.value.findIndex((t) => t.id === currentTrack.value?.id)
})
const nextTrack = computed(() => {
  const engineNext = playerStore.nextUpTrack
  if (engineNext) return engineNext
  const idx = currentTrackIndex.value
  if (idx < 0) return queueTracks.value[0] ?? null
  return queueTracks.value[idx + 1] ?? null
})
const upcomingCount = computed(() => {
  const idx = currentTrackIndex.value
  if (idx < 0) return queueTracks.value.length
  if (shuffleMode.value === 'catalog' || shuffleMode.value === 'similar') {
    return queueTracks.value.length - idx - 1
  }
  return Math.max(0, queueTracks.value.length - idx - 1)
})

// Album colors
const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressColor = computed(() => {
  return coverUrl.value && palette.value.vibrant ? palette.value.vibrant : 'var(--p-500)'
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

const title = computed(() => currentTrack.value?.title || 'Unknown Track')
const artistName = computed(() => currentTrack.value?.artistName || 'Unknown Artist')

const barBackgroundStyle = computed(() => {
  return {
    background: 'var(--s-950)',
    backdropFilter: 'blur(32px)',
    WebkitBackdropFilter: 'blur(32px)',
  }
})

const playAriaLabel = computed(() => {
  if (isLoadingTrack.value || isBuffering.value) return 'Loading'
  return isPlaying.value ? 'Pause' : 'Play'
})

function toggleCollapsed() {
  collapsed.value = !collapsed.value
  localStorage.setItem('player-bar-collapsed', String(collapsed.value))
  window.dispatchEvent(new CustomEvent('playerbar-collapse', { detail: collapsed.value }))
}

function playNextTrack() {
  playNext()
  showQueuePreview.value = false
}

function openFullQueue() {
  showQueuePreview.value = false
  emit('toggle-queue')
}

function handleTogglePiP() {
  pip.toggle().catch(() => {
    // PiP not supported or user denied — silently fail
  })
}

// ── Overflow menu sections ───────────────────────────────────────

const overflowAnchorEl = ref<HTMLElement | null>(null)

const {
  sleepTimerMinutes,
  crossfadeDuration,
  audioQuality,
  setSleepTimer,
  clearSleepTimer,
} = pc

const qualityOptions = ['auto', 'low', 'medium', 'high', 'lossless'] as const
type AudioQuality = typeof qualityOptions[number]

const qualityLabel = computed(() => {
  const q = audioQuality.value
  if (q === 'lossless') return 'Lossless'
  if (q === 'high') return 'High'
  if (q === 'medium') return 'Medium'
  if (q === 'low') return 'Low'
  return 'Auto'
})

function cycleQuality() {
  const idx = qualityOptions.indexOf(audioQuality.value as AudioQuality)
  const nextIdx = (idx + 1) % qualityOptions.length
  audioQuality.value = qualityOptions[nextIdx]!
}

function cycleSleepTimer() {
  const current = sleepTimerMinutes.value
  if (current === 0) { setSleepTimer(15); return }
  if (current === 15) { setSleepTimer(30); return }
  if (current === 30) { setSleepTimer(60); return }
  clearSleepTimer()
}

const sleepLabel = computed(() => {
  const m = sleepTimerMinutes.value
  if (m <= 0) return 'Off'
  return `${m}m`
})

function cycleCrossfade() {
  const current = crossfadeDuration.value
  const store = usePlayerStore()
  if (current === 0) { store.setCrossfadeDuration(3); return }
  if (current === 3) { store.setCrossfadeDuration(5); return }
  if (current === 5) { store.setCrossfadeDuration(8); return }
  if (current === 8) { store.setCrossfadeDuration(12); return }
  store.setCrossfadeDuration(0)
}

function copyTrackLink() {
  if (!currentTrack.value) return
  const url = `${window.location.origin}/track/${currentTrack.value.id}`
  navigator.clipboard?.writeText(url).then(() => {
    toast.success('Track link copied to clipboard')
  }).catch(() => {
    toast.error('Could not copy link')
  })
  showOverflow.value = false
}

function goToAlbumFromOverflow() {
  const ct = currentTrack.value as any
  if (ct?.albumId) {
    router.push(`/album/${ct.albumId}`)
  }
  showOverflow.value = false
}

function goToArtistFromOverflow() {
  const ct = currentTrack.value as any
  if (ct?.artistId) {
    router.push(`/artist/${ct.artistId}`)
  }
  showOverflow.value = false
}

const overflowSections = computed<ContextMenuSection[]>(() => {
  if (!currentTrack.value) return []

  return [
    {
      id: 'playback',
      label: 'PLAYBACK',
      items: [
        {
          id: 'audio-quality',
          label: 'Audio Quality',
          icon: 'Settings2',
          badge: qualityLabel.value,
          shortcut: 'Q',
          action: cycleQuality,
        },
        {
          id: 'sleep-timer',
          label: 'Sleep Timer',
          icon: 'Clock',
          badge: sleepLabel.value,
          action: cycleSleepTimer,
        },
        {
          id: 'crossfade',
          label: 'Crossfade',
          icon: 'Sparkles',
          badge: crossfadeDuration.value > 0 ? `${crossfadeDuration.value}s` : 'Off',
          action: cycleCrossfade,
        },
        {
          id: 'pip',
          label: 'Picture in Picture',
          icon: 'Monitor',
          shortcut: 'P',
          action: handleTogglePiP,
        },
      ],
    },
    {
      id: 'library',
      label: 'LIBRARY',
      items: [
        {
          id: 'like',
          label: liked.value ? 'Saved to Library' : 'Save to Library',
          icon: 'Heart',
          checked: liked.value,
          shortcut: 'S',
          action: toggleLike,
        },
        {
          id: 'add-to-playlist',
          label: 'Add to Playlist',
          icon: 'PlusCircle',
          shortcut: 'A',
          action: () => { showAddToPlaylist.value = true },
        },
      ],
    },
    {
      id: 'navigate',
      label: 'GO TO',
      items: [
        {
          id: 'go-to-album',
          label: 'Go to Album',
          icon: 'Disc3',
          action: goToAlbumFromOverflow,
        },
        {
          id: 'go-to-artist',
          label: 'Go to Artist',
          icon: 'UserRound',
          action: goToArtistFromOverflow,
        },
      ],
    },
    {
      id: 'share',
      label: 'SHARE',
      items: [
        {
          id: 'copy-link',
          label: 'Copy Link',
          icon: 'Link2',
          shortcut: 'Ctrl+C',
          separator: true,
          action: copyTrackLink,
        },
        {
          id: 'share',
          label: 'Share Track',
          icon: 'Share2',
          action: () => {
            if (!currentTrack.value) return
            const url = `${window.location.origin}/track/${currentTrack.value.id}`
            if (navigator.share) {
              navigator.share({ title: currentTrack.value.title, url }).catch(() => {})
            } else {
              copyTrackLink()
            }
            showOverflow.value = false
          },
        },
      ],
    },
    {
      id: 'advanced',
      label: 'ADVANCED',
      items: [
        {
          id: 'download',
          label: 'Download',
          icon: 'Download',
          hidden: true, // premium check
          separator: true,
          action: () => { showOverflow.value = false },
        },
        {
          id: 'track-info',
          label: 'Track Information',
          icon: 'Info',
          shortcut: 'I',
          action: () => {
            if (currentTrack.value?.id) router.push(`/track/${currentTrack.value.id}`)
            showOverflow.value = false
          },
        },
        {
          id: 'report',
          label: 'Report Issue',
          icon: 'Flag',
          danger: true,
          action: () => {
            toast.info('Issue reporting will be available in a future update')
            showOverflow.value = false
          },
        },
      ],
    },
  ]
})

const overflowHeader = computed<ContextMenuHeader | undefined>(() => {
  if (!currentTrack.value) return undefined
  return {
    coverUrl: currentTrack.value.coverUrl ?? undefined,
    title: currentTrack.value.title ?? 'Unknown Track',
    artistName: currentTrack.value.artistName ?? 'Unknown Artist',
    artistId: (currentTrack.value as any)?.artistId,
    duration: currentTrack.value.durationSeconds ?? undefined,
  }
})

function onVolume(val: number | number[]) {
  pc.setVolume(typeof val === 'number' ? val : val[0] ?? 0)
}

function seekRelative(seconds: number) {
  const newTime = Math.max(0, Math.min(duration.value, currentTime.value + seconds))
  const pct = duration.value > 0 ? (newTime / duration.value) * 100 : 0
  seekPercent(pct)
}

function formatTime(s: number) {
  if (!isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
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
  if (showShuffleMenu.value && shuffleBtnRef.value && !shuffleBtnRef.value.contains(target)) {
    showShuffleMenu.value = false
  }
  if (showRepeatMenu.value && repeatBtnRef.value && !repeatBtnRef.value.contains(target)) {
    showRepeatMenu.value = false
  }
  if (showQueuePreview.value && queueBtnRef.value && !queueBtnRef.value.contains(target)) {
    showQueuePreview.value = false
  }
}

// ── Radio integration ──
const radio = useRadio()
const openRadio = inject<((trackId: string, seedLabel?: string) => void) | null>('openRadio', null)

function handleRadioClick() {
  if (!currentTrack.value?.id) return
  if (openRadio) {
    openRadio(currentTrack.value.id, currentTrack.value.title)
  }
}

// Overflow menu outside click — handled by ContextMenu component
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

.expand-enter-active {
  transition: opacity 300ms ease, transform 300ms ease;
}
.expand-leave-active {
  transition: opacity 200ms ease, transform 200ms ease;
}
.expand-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.expand-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.animate-spin-slow {
  animation: spin-slow 1.5s linear infinite;
}
@keyframes spin-slow {
  to { transform: rotate(360deg); }
}

@media (prefers-reduced-motion: reduce) {
  .bar-slide-enter-active,
  .bar-slide-leave-active,
  .expand-enter-active,
  .expand-leave-active { transition: none; }
  .bar-slide-enter-from,
  .bar-slide-leave-to,
  .expand-enter-from,
  .expand-leave-to,
  .expand-enter-to,
  .expand-leave-from { opacity: 1; transform: none; }
  .animate-spin-slow { animation: none; }
}
</style>