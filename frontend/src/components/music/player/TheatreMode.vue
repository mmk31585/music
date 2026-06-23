<template>
  <Teleport to="body">
    <Transition name="theatre-enter">
      <div
        v-if="visible"
        ref="rootEl"
        class="fixed inset-0 z-[200] flex select-none"
        :style="bgStyle"
        @keydown="onKeydown"
        tabindex="0"
      >
        <!-- Dimmed overlay -->
        <div class="pointer-events-none absolute inset-0 bg-black/40 backdrop-blur-sm" />

        <div class="relative z-10 flex w-full flex-col">
          <!-- Top bar -->
          <div class="flex items-center justify-between px-5 pt-5 pb-3">
            <button
              type="button"
              class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 backdrop-blur-sm transition-all hover:bg-white/10 hover:text-white"
              aria-label="Close"
              @click="close"
            >
              <i aria-hidden="true" class="pi pi-chevron-down text-lg" />
            </button>
            <div class="glass flex items-center gap-2 rounded-full px-4 py-2 text-xs text-white/50">
              <span
                class="flex h-2 w-2 rounded-full bg-[#1db954]"
                :class="{ 'glow-spread': isPlaying }"
              />
              <span class="font-semibold tracking-wider uppercase">Theatre Mode</span>
            </div>
              <button
                type="button"
                class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 backdrop-blur-sm transition-all hover:bg-white/10 hover:text-white"
                :class="{ '!bg-[#1db954]/15 !text-[#1db954]': karaokeMode }"
                aria-label="Toggle karaoke"
                @click="karaokeMode = !karaokeMode"
              >
                <i aria-hidden="true" class="pi pi-file text-sm" />
            </button>
          </div>

          <!-- Main: Album art left, lyrics right -->
          <div class="flex flex-1 items-center gap-8 overflow-hidden px-8 pb-8">
            <!-- LEFT: Album art -->
            <div class="w-80 shrink-0 md:w-96 lg:w-[440px]">
              <div
                class="relative aspect-square overflow-hidden rounded-3xl shadow-2xl ring-1 ring-white/10"
              >
                <img
                  v-if="coverUrl"
                  :src="coverUrl"
                  :alt="title"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  :class="{ 'vinyl-spin': isPlaying, 'vinyl-spin-paused': !isPlaying }"
                  @error="onImgError"
                />
                <div
                  v-else
                  class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212]"
                >
                  <i aria-hidden="true" class="pi pi-music text-5xl text-white/30" />
                </div>
              </div>

              <!-- Track info under album art -->
              <div class="mt-5 text-center">
                <p class="truncate text-xl font-bold text-white">{{ title }}</p>
                <p
                  class="mt-1 cursor-pointer truncate text-sm text-white/50 transition-colors hover:text-[#1db954]"
                >
                  {{ artistName }}
                </p>
                <p v-if="albumTitle" class="mt-0.5 truncate text-xs text-white/30">
                  {{ albumTitle }}
                </p>
              </div>

              <!-- Playback controls -->
              <div class="mt-5 flex items-center justify-center gap-4">
                <button
                  type="button"
                  class="spring relative flex h-10 w-10 items-center justify-center rounded-full text-white/40 transition-all hover:bg-white/10 hover:text-white"
                  :class="{ '!text-[#1db954]': shuffleMode !== 'off' }"
                  :disabled="!currentTrack"
                  aria-label="Shuffle"
                  @click="toggleShuffle"
                >
                  <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
                  <span
                    v-if="shuffleMode !== 'off'"
                    class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-[#a855f7] text-[8px] font-bold text-white"
                  >{{ shuffleMode === 'queue' ? 'Q' : shuffleMode === 'catalog' ? 'R' : 'S' }}</span>
                </button>

                <button
                  type="button"
                  class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 transition-all hover:bg-white/10 hover:text-white disabled:opacity-20"
                  :disabled="!hasPrevious"
                  aria-label="Previous track"
                  @click="playPrevious"
                >
                  <i aria-hidden="true" class="pi pi-step-backward text-lg" />
                </button>

                <button
                  type="button"
                  class="glow-green spring relative flex h-14 w-14 items-center justify-center rounded-full bg-white text-black shadow-2xl transition-all hover:scale-105 hover:bg-[#1db954] hover:text-white disabled:opacity-40"
                  :class="{ '!bg-[#1db954] !text-white': isPlaying }"
                  :disabled="!currentTrack || isLoadingTrack"
                  :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
                  @click="togglePlayPause"
                >
                  <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-lg" />
                  <i
                    v-else
                    :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'"
                    class="text-lg"
                  />
                  <div
                    v-if="isPlaying"
                    class="absolute -inset-2 animate-ping rounded-full border-2 border-[#1db954]/30"
                  />
                </button>

                <button
                  type="button"
                  class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 transition-all hover:bg-white/10 hover:text-white disabled:opacity-20"
                  :disabled="!hasNext"
                  aria-label="Next track"
                  @click="playNext"
                >
                  <i aria-hidden="true" class="pi pi-step-forward text-lg" />
                </button>

                <button
                  type="button"
                  class="spring relative flex h-10 w-10 items-center justify-center rounded-full text-white/40 transition-all hover:bg-white/10 hover:text-white"
                  :class="{ '!text-[#1db954]': repeatMode !== 'off' }"
                  :disabled="!currentTrack"
                  aria-label="Repeat"
                  @click="toggleRepeat"
                >
                  <i aria-hidden="true" class="pi pi-refresh text-sm" />
                  <span
                    v-if="repeatMode === 'one'"
                    class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#1db954] text-[9px] font-bold text-black"
                    >1</span
                  >
                </button>
              </div>

              <!-- Seekbar -->
              <div class="mt-4 w-full">
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="0.1"
                  class="fullscreen-range w-full"
                  :style="progressStyle"
                  :value="progressPercent"
                  :disabled="!currentTrack"
                  @input="onSeek"
                />
                <div class="mt-1 flex justify-between text-[11px] text-white/40 tabular-nums">
                  <span>{{ currentTimeLabel }}</span>
                  <span>{{ durationLabel }}</span>
                </div>
              </div>
            </div>

            <!-- RIGHT: Lyrics -->
            <div class="flex h-full flex-1 flex-col overflow-hidden">
              <div class="flex items-center gap-3 border-b border-white/[0.06] pb-3">
                <h3 class="text-sm font-semibold tracking-wider text-white/60 uppercase">Lyrics</h3>
                <span
                  v-if="lyricsLanguage"
                  class="rounded bg-white/5 px-2 py-0.5 text-[10px] font-medium text-white/30 uppercase"
                  >{{ lyricsLanguage }}</span
                >
              </div>

              <div class="relative flex-1">
                <KaraokeLyrics
                  :content="lyricsContent"
                  :type="lyricsType"
                  :language="lyricsLanguage ?? undefined"
                  :current-time="currentTime"
                  :loading="lyricsLoading"
                  :karaoke="karaokeMode"
                  @seek="seekTo"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import KaraokeLyrics from './KaraokeLyrics.vue'
import { useLyricsApi, type Lyrics } from '@/services/api/lyrics'
import { onImgError } from '@/utils/helpers'

defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const rootEl = ref<HTMLElement | null>(null)
const karaokeMode = ref(true)

const pc = usePlayerControls()
const currentTrack = pc.currentTrack
const isPlaying = pc.isPlaying
const isBuffering = pc.isBuffering
const isLoadingTrack = pc.isLoadingTrack
const currentTime = pc.currentTime
const duration = pc.duration
const progressPercent = pc.progressPercent
const hasNext = pc.hasNext
const hasPrevious = pc.hasPrevious
const shuffleMode = pc.shuffleMode
const repeatMode = pc.repeatMode

const togglePlayPause = pc.togglePlayPause
const seekPercent = pc.seekPercent
const seek = pc.seek
const playNext = pc.playNext
const playPrevious = pc.playPrevious
const toggleShuffle = pc.toggleShuffle
const toggleRepeat = pc.toggleRepeat

const title = computed(() => currentTrack.value?.title || 'No track')
const artistName = computed(() => currentTrack.value?.artistName || '')
const albumTitle = computed(() => currentTrack.value?.albumTitle || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')
const { palette: albumPalette } = useAlbumColors(coverUrl)

const bgStyle = computed(() => {
  const p = albumPalette.value
  if (!coverUrl.value) {
    return { background: 'linear-gradient(135deg, #0a0a0a 0%, #121212 100%)' }
  }
  return {
    background: `
      radial-gradient(ellipse at 30% 30%, ${p.vibrant}30 0%, transparent 60%),
      radial-gradient(ellipse at 70% 70%, ${p.light}20 0%, transparent 50%),
      linear-gradient(135deg, ${p.dark} 0%, ${p.dominant}aa 50%, ${p.muted} 100%)
    `,
  }
})

function fmtTime(s: number) {
  const total = Math.max(0, Math.floor(Number(s) || 0))
  const m = Math.floor(total / 60)
  const sec = total % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}

const currentTimeLabel = computed(() => fmtTime(currentTime.value))
const durationLabel = computed(() =>
  fmtTime(duration.value || currentTrack.value?.durationSeconds || 0),
)
const progressStyle = computed(() => ({ '--range-progress': `${progressPercent.value}%` }))

function onSeek(e: Event) {
  seekPercent(Number((e.target as HTMLInputElement).value))
}

function close() {
  emit('update:visible', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
  if (e.key === ' ' && e.target === e.currentTarget) {
    e.preventDefault()
    togglePlayPause()
  }
}

onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
})

// Lyrics
const lyricsApi = useLyricsApi()
const lyricsData = ref<Lyrics | null>(null)
const lyricsLoading = ref(false)
const lyricsError = ref(false)
const lyricsLanguage = ref<string | null>(null)

const lyricsContent = computed(() => lyricsData.value?.content || '')
const lyricsType = computed(() => lyricsData.value?.type || 'plain')

async function fetchLyrics(trackId: string) {
  lyricsLoading.value = true
  lyricsError.value = false
  lyricsData.value = null
  lyricsLanguage.value = null
  try {
    const res = await lyricsApi.getTrackLyrics(trackId)
    lyricsData.value = (res ?? null) as Lyrics | null
    lyricsLanguage.value = (res as any as { language?: string })?.language || null
  } catch {
    lyricsError.value = true
  } finally {
    lyricsLoading.value = false
  }
}

function seekTo(seconds: number) {
  seek(seconds)
}

watch(
  () => currentTrack.value?.id,
  (id) => {
    if (id) {
      void fetchLyrics(id)
    }
  },
  { immediate: true },
)

onBeforeUnmount(() => {
})
</script>

<style scoped>
.theatre-enter-enter-active {
  transition:
    opacity 250ms ease,
    transform 250ms ease;
}
.theatre-enter-leave-active {
  transition:
    opacity 150ms ease,
    transform 150ms ease;
}
.theatre-enter-enter-from {
  opacity: 0;
  transform: scale(0.96);
}
.theatre-enter-leave-to {
  opacity: 0;
  transform: scale(0.96);
}
.theatre-lyric-line {
  scroll-margin: 40vh;
}
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
