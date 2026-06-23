<template>
  <Teleport to="body">
    <Transition name="radio-enter">
      <div
        v-if="visible"
        ref="rootEl"
        class="fixed inset-0 z-[200] flex flex-col select-none"
        :style="bgStyle"
        @keydown="onKeydown"
        tabindex="0"
      >
        <div class="pointer-events-none absolute inset-0 overflow-hidden">
          <div
            class="aurora-spot-1 -top-40 -left-40 bg-[#1db954]/15"
            style="animation-duration: 25s"
          />
          <div
            class="aurora-spot-2 -right-40 -bottom-40 bg-[#60a5fa]/10"
            style="animation-duration: 30s"
          />
        </div>

        <div class="relative z-10 flex items-center justify-between px-5 pt-5 pb-3">
          <button
            type="button"
            class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 backdrop-blur-sm transition-all hover:bg-white/10 hover:text-white"
            aria-label="Close"
            @click="close"
          >
            <i aria-hidden="true" class="pi pi-chevron-down text-lg" />
          </button>
          <div class="glass flex items-center gap-2 rounded-full px-4 py-2 text-xs text-white/50">
            <span class="flex h-2 w-2 animate-pulse rounded-full bg-[#1db954]" />
            <span class="font-semibold tracking-wider uppercase">Radio</span>
            <span v-if="displaySeedLabel" class="text-white/30">· {{ displaySeedLabel }}</span>
          </div>
          <button
            type="button"
            class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 backdrop-blur-sm transition-all hover:bg-white/10 hover:text-white"
            :disabled="!canSkip"
            aria-label="Skip track"
            @click="skipTrack"
          >
            <i aria-hidden="true" class="pi pi-forward text-lg" />
          </button>
        </div>

        <div class="relative z-10 flex flex-1 overflow-hidden">
          <div class="flex w-full flex-col gap-6 px-6 pb-6 lg:flex-row lg:gap-8 lg:px-10">
            <div class="flex flex-col items-center justify-center lg:w-1/2">
              <div class="relative h-48 w-48 md:h-64 md:w-64 lg:h-80 lg:w-80">
                <div
                  class="h-full w-full overflow-hidden rounded-3xl shadow-2xl ring-1 shadow-black/40 ring-white/10"
                >
                  <img
                    v-if="coverUrl"
                    :src="coverUrl"
                    :alt="title"
                    class="h-full w-full object-cover"
                    :class="{ 'vinyl-spin': isPlaying, 'vinyl-spin-paused': !isPlaying }"
                    @error="onImgError"
                  />
                  <div
                    v-else
                    class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212]"
                  >
                    <i aria-hidden="true" class="pi pi-music text-5xl text-white/20" />
                  </div>
                </div>
              </div>

              <div class="mt-5 text-center">
                <p class="text-xl font-bold text-white md:text-2xl">{{ title }}</p>
                <p
                  class="mt-1 cursor-pointer text-sm text-white/50 transition-colors hover:text-[#1db954]"
                >
                  {{ artistName }}
                </p>
              </div>

              <div class="mt-5 flex items-center gap-5">
          <button
            type="button"
            class="spring flex h-12 w-12 items-center justify-center rounded-full text-white/60 transition-all hover:bg-white/10 hover:text-white disabled:opacity-20"
            :disabled="!hasPrevious"
            aria-label="Previous track"
            @click="playPrevious"
          >
            <i aria-hidden="true" class="pi pi-step-backward text-xl" />
          </button>

          <button
            type="button"
            class="glow-green spring relative flex h-16 w-16 items-center justify-center rounded-full bg-white text-black shadow-2xl transition-all hover:scale-105 hover:bg-[#1db954] hover:text-white disabled:opacity-40"
            :class="{ '!bg-[#1db954] !text-white': isPlaying }"
            :disabled="!currentTrack || isLoadingTrack"
            :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
            @click="togglePlayPause"
          >
            <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-xl" />
            <i
              v-else
              :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'"
              class="text-xl"
            />
            <div
              v-if="isPlaying"
              class="absolute -inset-2 animate-ping rounded-full border-2 border-[#1db954]/30"
            />
          </button>

          <button
            type="button"
            class="spring flex h-12 w-12 items-center justify-center rounded-full text-white/60 transition-all hover:bg-white/10 hover:text-white"
            aria-label="Next track"
            @click="skipTrack"
          >
                  <i aria-hidden="true" class="pi pi-step-forward text-xl" />
                </button>
              </div>

              <div class="mt-4 w-full max-w-sm">
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

            <div
              class="flex flex-1 flex-col overflow-hidden lg:border-l lg:border-white/[0.06] lg:pl-8"
            >
              <div class="flex items-center gap-4 border-b border-white/[0.06] pb-3">
                <button
                  type="button"
                  class="spring text-xs font-semibold tracking-wider uppercase transition-all"
                  :class="
                    subTab === 'upcoming' ? 'text-white' : 'text-white/40 hover:text-white/70'
                  "
                  @click="subTab = 'upcoming'"
                >
                  Up Next
                </button>
                <button
                  type="button"
                  class="spring text-xs font-semibold tracking-wider uppercase transition-all"
                  :class="subTab === 'history' ? 'text-white' : 'text-white/40 hover:text-white/70'"
                  @click="subTab = 'history'"
                >
                  History
                </button>
              </div>

              <div class="mt-4 flex-1 scrollbar-thin space-y-1 overflow-y-auto pr-2">
                <div v-show="subTab === 'upcoming'">
                  <div
                    v-for="(track, idx) in upcomingList"
                    :key="track.id"
                    role="button"
                    tabindex="0"
                    class="group spring flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2 transition-all hover:bg-white/[0.06]"
                    @click="playUpcoming(idx)"
                    @keydown.enter="playUpcoming(idx)"
                    @keydown.space.prevent="playUpcoming(idx)"
                  >
                    <span class="w-5 text-center text-xs text-white/20 tabular-nums">{{
                      idx + 1
                    }}</span>
                    <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg">
                      <img
                        v-if="track.coverUrl"
                        :src="track.coverUrl"
                        :alt="track.title"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div
                        v-else
                        class="flex h-full w-full items-center justify-center bg-white/10"
                      >
                        <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p
                        class="truncate text-sm font-medium text-white/80 transition-colors group-hover:text-white"
                      >
                        {{ track.title }}
                      </p>
                      <p class="truncate text-xs text-white/40">{{ track.artistName }}</p>
                    </div>
                    <div class="text-[#1db954] opacity-0 group-hover:opacity-100">
                      <i aria-hidden="true" class="pi pi-play-fill text-sm" />
                    </div>
                  </div>
                  <div v-if="isLoadingBatch" class="flex justify-center py-4">
                    <i aria-hidden="true" class="pi pi-spin pi-spinner text-white/30" />
                  </div>
                  <div
                    v-else-if="upcomingList.length === 0"
                    class="flex flex-col items-center gap-2 py-8 text-white/20"
                  >
                    <i aria-hidden="true" class="pi pi-wave-pulse text-2xl" />
                    <p class="text-xs">Generating your radio...</p>
                  </div>
                </div>

                <div v-show="subTab === 'history'">
                  <div
                    v-for="(track, idx) in historyList"
                    :key="track.id + '-' + idx"
                    class="flex items-center gap-3 rounded-xl px-3 py-2 text-white/40"
                  >
                    <span class="w-5 text-center text-xs text-white/15 tabular-nums">{{
                      historyList.length - idx
                    }}</span>
                    <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg opacity-60">
                      <img
                        v-if="track.coverUrl"
                        :src="track.coverUrl"
                        :alt="track.title"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div
                        v-else
                        class="flex h-full w-full items-center justify-center bg-white/10"
                      >
                        <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm text-white/50">{{ track.title }}</p>
                      <p class="truncate text-xs text-white/25">{{ track.artistName }}</p>
                    </div>
                  </div>
                  <div
                    v-if="historyList.length === 0"
                    class="flex flex-col items-center gap-2 py-8 text-white/20"
                  >
                    <i aria-hidden="true" class="pi pi-history text-2xl" />
                    <p class="text-xs">No history yet</p>
                  </div>
                </div>
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
import { useRadio } from '@/composables/recommendation/useRadio'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { onImgError } from '@/utils/helpers'
import type { PlaybackTrack } from '@/services/api/player'

const props = defineProps<{
  visible: boolean
  seedType?: 'track' | 'artist' | 'genre'
  seedId?: string
  seedLabel?: string
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const rootEl = ref<HTMLElement | null>(null)
const subTab = ref<'upcoming' | 'history'>('upcoming')

const pc = usePlayerControls()
const radio = useRadio()

const currentTrack = pc.currentTrack
const isPlaying = pc.isPlaying
const isBuffering = pc.isBuffering
const isLoadingTrack = pc.isLoadingTrack
const currentTime = pc.currentTime
const duration = pc.duration
const progressPercent = pc.progressPercent
const hasNext = pc.hasNext
const hasPrevious = pc.hasPrevious

const togglePlayPause = pc.togglePlayPause
const seekPercent = pc.seekPercent
const playNext = pc.playNext
const playPrevious = pc.playPrevious

const title = computed(() => currentTrack.value?.title || 'No track')
const artistName = computed(() => currentTrack.value?.artistName || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')
const { palette: albumPalette } = useAlbumColors(coverUrl)

const displaySeedLabel = computed(() => {
  return radio.radioSeedLabel.value || props.seedLabel || ''
})

const upcomingList = computed<PlaybackTrack[]>(() => {
  // Show upcoming tracks from the queue, skipping the current track
  const q = pc.queue.value
  const currentId = currentTrack.value?.id
  const currentIdx = q.findIndex((t: PlaybackTrack) => t.id === currentId)
  if (currentIdx >= 0) {
    return q.slice(currentIdx + 1)
  }
  return q.slice(1)
})

const historyList = computed<PlaybackTrack[]>(() => {
  const q = pc.queue.value
  const currentId = currentTrack.value?.id
  const currentIdx = q.findIndex((t: PlaybackTrack) => t.id === currentId)
  if (currentIdx > 0) {
    return q.slice(0, currentIdx).reverse()
  }
  return []
})

const isLoadingBatch = computed(() => radio.isLoadingBatch.value)

const bgStyle = computed(() => {
  const p = albumPalette.value
  if (!coverUrl.value) {
    return { background: 'linear-gradient(135deg, #0a0a0a 0%, #121212 100%)' }
  }
  return {
    background: `
      radial-gradient(ellipse at 20% 50%, ${p.vibrant}22 0%, transparent 60%),
      radial-gradient(ellipse at 80% 20%, ${p.light}15 0%, transparent 50%),
      linear-gradient(180deg, ${p.dark} 0%, ${p.dominant}99 50%, ${p.muted} 100%)
    `,
  }
})

const canSkip = computed(() => upcomingList.value.length > 0 || hasNext.value)

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
  radio.endRadio()
  emit('update:visible', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
  if (e.key === ' ' && e.target === e.currentTarget) {
    e.preventDefault()
    togglePlayPause()
  }
}

function skipTrack() {
  playNext()
}

function playUpcoming(index: number) {
  const track = upcomingList.value[index]
  if (!track) return
  pc.setQueueAndPlay(pc.queue.value, pc.queue.value.findIndex((t: PlaybackTrack) => t.id === track.id))
}

async function startRadioSession() {
  if (props.seedId) {
    await radio.startFromTrack(props.seedId, props.seedLabel)
  }
}

watch(
  () => props.visible,
  async (v) => {
    if (v) {
      await nextTick()
      await startRadioSession()
    }
  },
)

watch(
  () => props.seedId,
  async (seedId) => {
    if (seedId && props.visible) {
      await nextTick()
      await startRadioSession()
    }
  },
)

onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
})

onBeforeUnmount(() => {
  radio.endRadio()
})
</script>

<style scoped>
.radio-enter-enter-active {
  transition:
    opacity 250ms ease,
    transform 250ms ease;
}
.radio-enter-leave-active {
  transition:
    opacity 150ms ease,
    transform 150ms ease;
}
.radio-enter-enter-from {
  opacity: 0;
  transform: scale(0.96);
}
.radio-enter-leave-to {
  opacity: 0;
  transform: scale(0.96);
}
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
