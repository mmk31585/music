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
            class="aurora-spot-1 -top-40 -left-40 bg-accent-subtle"
            style="animation-duration: 25s"
          />
          <div
            class="aurora-spot-2 -right-40 -bottom-40 bg-aurora-blue/10"
            style="animation-duration: 30s"
          />
        </div>

        <div class="relative z-10 flex items-center justify-between px-5 pt-5 pb-3">
          <button
            type="button"
            class="spring flex h-10 w-10 items-center justify-center rounded-full text-secondary backdrop-blur-xs transition-all hover:bg-surface-active hover:text-primary"
            aria-label="بستن رادیو"
            title="رادیو متوقف می‌شه"
            @click="close"
          >
            <ChevronDown aria-hidden="true" class="text-lg"  />
          </button>
          <div class="glass flex items-center gap-2 rounded-full px-4 py-2 text-xs text-secondary">
            <span class="flex h-2 w-2 animate-pulse rounded-full bg-accent" />
            <span class="font-semibold tracking-wider uppercase">Radio</span>
            <span v-if="displaySeedLabel" class="text-muted">· {{ displaySeedLabel }}</span>
          </div>
          <button
            type="button"
            class="spring flex h-10 w-10 items-center justify-center rounded-full text-secondary backdrop-blur-xs transition-all hover:bg-surface-active hover:text-primary"
            :disabled="!canSkip"
            aria-label="آهنگ بعدی"
            @click="skipTrack"
          >
            <Forward aria-hidden="true" class="text-lg"  />
          </button>
        </div>

        <div class="relative z-10 flex flex-1 overflow-hidden">
          <div class="flex w-full flex-col gap-6 px-6 pb-6 lg:flex-row lg:gap-8 lg:px-10">
            <div class="flex flex-col items-center justify-center lg:w-1/2">
              <div class="relative h-48 w-48 md:h-64 md:w-64 lg:h-80 lg:w-80">
                <div
                  class="h-full w-full overflow-hidden rounded-3xl shadow-2xl ring-1 shadow-bg-overlay/40 ring-border-default"
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
                    class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-surface-raised"
                  >
                    <Music aria-hidden="true" class="text-5xl text-muted"  />
                  </div>
                </div>
              </div>

              <div class="mt-5 text-center">
                <p class="text-xl font-bold text-primary md:text-2xl">{{ title }}</p>
                <p
                  class="mt-1 cursor-pointer text-sm text-secondary transition-colors hover:text-accent"
                >
                  {{ artistName }}
                </p>
              </div>

              <div class="mt-5 flex items-center gap-5">
          <button
            type="button"
            class="spring flex h-12 w-12 items-center justify-center rounded-full text-secondary transition-all hover:bg-surface-active hover:text-primary disabled:opacity-20"
            :disabled="!hasPrevious"
            aria-label="Previous track"
            @click="playPrevious"
          >
            <SkipBack aria-hidden="true" class="text-xl"  />
          </button>

          <button
            type="button"
            class="glow-green spring relative flex h-16 w-16 items-center justify-center rounded-full bg-white text-black shadow-2xl transition-all hover:scale-105 hover:bg-accent hover:text-primary disabled:opacity-40"
            :class="{ '!bg-accent !text-primary': isPlaying }"
            :disabled="!currentTrack || isLoadingTrack"
            :aria-label="isLoadingTrack || isBuffering ? 'در حال بارگذاری' : isPlaying ? 'توقف' : 'پخش'"
            @click="togglePlayPause"
          >
            <Loader2 v-if="isLoadingTrack || isBuffering" aria-hidden="true" class="animate-spin text-xl" />
            <Pause v-else-if="isPlaying" aria-hidden="true" class="text-xl" />
            <Play v-else aria-hidden="true" class="text-xl ms-0.5" />
            <div
              v-if="isPlaying"
              class="absolute -inset-2 animate-ping rounded-full border-2 border-spotify/30"
            />
          </button>

          <button
            type="button"
            class="spring flex h-12 w-12 items-center justify-center rounded-full text-secondary transition-all hover:bg-surface-active hover:text-primary"
            aria-label="Next track"
            @click="skipTrack"
          >
                  <SkipForward aria-hidden="true" class="text-xl"  />
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
                <div class="mt-1 flex justify-between text-[11px] text-tertiary tabular-nums">
                  <span>{{ currentTimeLabel }}</span>
                  <span>{{ durationLabel }}</span>
                </div>
              </div>
            </div>

            <div
              class="flex flex-1 flex-col overflow-hidden lg:border-l lg:border-border-subtle lg:pl-8"
            >
              <div class="flex items-center gap-4 border-b border-border-subtle pb-3">
                <button
                  type="button"
                  class="spring text-xs font-semibold tracking-wider uppercase transition-all"
                  :class="
                    subTab === 'upcoming' ? 'text-primary' : 'text-tertiary hover:text-primary/70'
                  "
                  @click="subTab = 'upcoming'"
                >
                  بعدی
                </button>
                <button
                  type="button"
                  class="spring text-xs font-semibold tracking-wider uppercase transition-all"
                  :class="subTab === 'history' ? 'text-primary' : 'text-tertiary hover:text-primary/70'"
                  @click="subTab = 'history'"
                >
                  تاریخچه
                </button>
              </div>

              <div class="mt-4 flex-1 scrollbar-thin space-y-1 overflow-y-auto pr-2">
                <div v-show="subTab === 'upcoming'">
                  <div
                    v-for="(track, idx) in upcomingList"
                    :key="track.id"
                    role="button"
                    tabindex="0"
                    class="group spring flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2 transition-all hover:bg-surface-overlay"
                    @click="playUpcoming(idx)"
                    @keydown.enter="playUpcoming(idx)"
                    @keydown.space.prevent="playUpcoming(idx)"
                  >
                    <span class="w-5 text-center text-xs text-muted tabular-nums">{{
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
                        class="flex h-full w-full items-center justify-center bg-surface-active"
                      >
                        <Music aria-hidden="true" class="text-xs text-muted"  />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p
                        class="truncate text-sm font-medium text-primary/80 transition-colors group-hover:text-primary"
                      >
                        {{ track.title }}
                      </p>
                      <p class="truncate text-xs text-tertiary">{{ track.artistName }}</p>
                    </div>
                    <div class="text-accent opacity-0 group-hover:opacity-100">
                      <Play aria-hidden="true" class="text-sm"  />
                    </div>
                  </div>
                  <div v-if="upcomingList.length === 0 && isLoadingBatch" class="flex flex-col items-center gap-3 py-12 text-center">
                    <div class="flex items-end gap-1 h-8">
                      <div class="w-1 rounded-full bg-accent/60" style="animation: eq-bar 0.8s ease-in-out infinite; height: 8px;" />
                      <div class="w-1 rounded-full bg-accent/80" style="animation: eq-bar 0.8s ease-in-out infinite 0.15s; height: 16px;" />
                      <div class="w-1 rounded-full bg-accent" style="animation: eq-bar 0.8s ease-in-out infinite 0.3s; height: 12px;" />
                      <div class="w-1 rounded-full bg-accent/80" style="animation: eq-bar 0.8s ease-in-out infinite 0.45s; height: 20px;" />
                      <div class="w-1 rounded-full bg-accent/60" style="animation: eq-bar 0.8s ease-in-out infinite 0.6s; height: 10px;" />
                    </div>
                    <p class="text-xs text-muted">در حال ساخت رادیو بر اساس سلیقه شما...</p>
                  </div>
                  <div
                    v-else-if="upcomingList.length === 0 && !isLoadingBatch"
                    class="flex flex-col items-center gap-2 py-8 text-muted"
                  >
                    <Activity aria-hidden="true" class="text-2xl"  />
                    <p class="text-xs">در حال ساخت رادیو...</p>
                  </div>
                </div>

                <div v-show="subTab === 'history'">
                  <div
                    v-for="(track, idx) in historyList"
                    :key="track.id + '-' + idx"
                    class="flex items-center gap-3 rounded-xl px-3 py-2 text-tertiary"
                  >
                    <span class="w-5 text-center text-xs text-muted tabular-nums">{{
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
                        class="flex h-full w-full items-center justify-center bg-surface-active"
                      >
                        <Music aria-hidden="true" class="text-xs text-muted"  />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm text-secondary">{{ track.title }}</p>
                      <p class="truncate text-xs text-muted">{{ track.artistName }}</p>
                    </div>
                  </div>
                  <div
                    v-if="historyList.length === 0"
                    class="flex flex-col items-center gap-2 py-8 text-muted"
                  >
                    <History aria-hidden="true" class="text-2xl"  />
                    <p class="text-xs">هنوز تاریخچه‌ای نیست</p>
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
import { Activity, ChevronDown, Forward, History, Loader2, Music, Pause, Play, SkipBack, SkipForward } from 'lucide-vue-next'
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

const title = computed(() => currentTrack.value?.title || 'آهنگی انتخاب نشده')
const artistName = computed(() => currentTrack.value?.artistName || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')
const { palette: albumPalette } = useAlbumColors(coverUrl)

const displaySeedLabel = computed(() => {
  return radio.radioSeedLabel.value || props.seedLabel || currentTrack.value?.title || ''
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
    return { background: 'var(--bg-gradient)' }
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
    return
  }
  const current = currentTrack.value
  if (current && !radio.isRadioActive.value) {
    await radio.startFromTrack(current.id, current.title)
    return
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

@keyframes eq-bar {
  0%, 100% { transform: scaleY(0.4); }
  50% { transform: scaleY(1); }
}

@media (prefers-reduced-motion: reduce) {
  [style*="eq-bar"] { animation: none !important; }
}
</style>
