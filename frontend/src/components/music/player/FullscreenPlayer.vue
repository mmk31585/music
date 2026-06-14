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
        <!-- Aurora + blurred art background -->
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
            class="absolute inset-0 transition-all duration-1000"
            :style="auroraGradient"
          />
        </div>

        <!-- HEADER (56px) -->
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
              @click="onOverflowToggle"
              aria-label="More options"
              :aria-expanded="showOverflow"
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

        <!-- TAB CONTENT -->
        <div class="relative z-10 flex flex-1 flex-col overflow-hidden">

          <!-- ─── NOW PLAYING TAB ─── -->
          <div
            v-if="activeTab === 'now-playing'"
            class="flex flex-1 flex-col items-center overflow-y-auto px-4 pb-4"
          >
            <div class="flex w-full max-w-md flex-1 flex-col items-center justify-center gap-5">
              <!-- Album Art -->
              <div class="relative" :class="coverSizeClass">
                <div
                  class="relative h-full w-full overflow-hidden rounded-2xl shadow-[0_24px_80px_rgba(0,0,0,0.7)] transition-all duration-1000"
                  :class="{
                    'shadow-[0_0_60px_rgba(29,185,84,0.2)]': isPlaying
                  }"
                >
                  <img
                    v-if="coverUrl"
                    :src="coverUrl"
                    :alt="title"
                    loading="lazy"
                    class="h-full w-full object-cover transition-all duration-300"
                    :class="artScaleClass"
                    @error="onImgError"
                  />
                  <div
                    v-else
                    class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212]"
                  >
                    <i class="pi pi-music text-5xl text-white/30" />
                  </div>
                </div>
                <!-- Pulse ring when playing -->
                <div
                  v-if="isPlaying"
                  class="pointer-events-none absolute -inset-4 animate-[pulse-ring_2s_ease-in-out_infinite] rounded-full border border-[#1db954]/30"
                />
              </div>

              <!-- Track Info -->
              <div class="flex w-full max-w-md flex-col items-center gap-1 text-center">
                <h2 class="max-w-full truncate text-[22px] font-bold text-white" dir="auto">
                  {{ title }}
                </h2>
                <div class="flex items-center gap-3">
                  <p class="text-base text-[rgba(255,255,255,0.6)]" dir="auto">
                    {{ artistName }}
                  </p>
                  <button
                    type="button"
                    :aria-label="isLiked ? 'Unlike' : 'Like'"
                    class="flex items-center justify-center transition-all active:scale-90"
                    :class="isLiked ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)] hover:text-white'"
                    @click="toggleLike"
                  >
                    <i :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-lg" />
                  </button>
                </div>
              </div>

              <!-- Seekbar -->
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

              <!-- Controls -->
              <div class="flex items-center gap-5">
                <button
                  type="button"
                  :aria-label="shuffleMode ? 'Shuffle on' : 'Shuffle off'"
                  class="flex h-8 w-8 items-center justify-center rounded-full transition-all hover:bg-white/10 hover:text-white active:scale-90 disabled:opacity-30"
                  :class="shuffleMode ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
                  :disabled="!currentTrack"
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
                  <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-2xl" />
                  <i v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-2xl" />
                  <!-- Expanding ring pulse when playing -->
                  <span
                    v-if="isPlaying"
                    class="absolute inset-0 animate-[play-ring_2s_ease-out_infinite] rounded-full border-2 border-[#1db954]/40"
                  />
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
                  class="relative flex h-8 w-8 items-center justify-center rounded-full transition-all hover:bg-white/10 hover:text-white active:scale-90 disabled:opacity-30"
                  :class="repeatMode !== 'off' ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
                  :disabled="!currentTrack"
                  @click="toggleRepeat"
                >
                  <i class="pi pi-refresh text-sm" />
                  <span
                    v-if="repeatMode === 'one'"
                    class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#1db954] text-[9px] font-bold text-black"
                  >1</span>
                </button>
              </div>

              <!-- Secondary controls -->
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

                <button
                  type="button"
                  class="flex h-8 items-center gap-1.5 rounded-full bg-white/[0.04] px-3 text-xs font-bold transition-all hover:bg-white/10 hover:text-white active:scale-95"
                  :class="{ 'bg-[#1db954]/10 text-[#1db954]': sleepTimerMinutes > 0 }"
                  @click="sleepMenuOpen = !sleepMenuOpen"
                >
                  <i class="pi pi-clock text-[10px]" />
                  <span v-if="sleepTimerMinutes > 0" class="tabular-nums">{{ sleepTimerMinutes }}m</span>
                </button>

                <button
                  type="button"
                  class="flex h-8 w-8 items-center justify-center rounded-full transition-all hover:text-white active:scale-90"
                  :aria-label="crossfadeDuration > 0 ? `Crossfade ${crossfadeDuration}s` : 'Crossfade off'"
                  @click="cycleCrossfade"
                >
                  <i class="pi pi-arrows-alt text-xs" />
                  <span v-if="crossfadeDuration > 0" class="ml-0.5 text-[10px] tabular-nums">{{ crossfadeDuration }}s</span>
                </button>

                <button
                  type="button"
                  class="flex h-8 w-8 items-center justify-center rounded-full transition-all hover:text-white active:scale-90"
                  aria-label="Picture in Picture"
                  @click="onTogglePiP"
                >
                  <i class="pi pi-window-maximize text-xs" />
                </button>
              </div>
            </div>

            <!-- Up Next -->
            <div v-if="upNextTracks.length" class="w-full max-w-md mt-4">
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
                    <p class="truncate text-sm font-medium text-white/80" dir="auto">{{ track.title }}</p>
                    <p class="truncate text-xs text-white/40" dir="auto">{{ track.artistName }}</p>
                  </div>
                  <span class="text-xs text-white/30 tabular-nums">{{ formatTime(track.durationSeconds) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- ─── QUEUE TAB ─── -->
          <div
            v-if="activeTab === 'queue'"
            class="flex flex-1 flex-col overflow-hidden"
          >
            <div class="flex-1 overflow-y-auto px-4 py-4 md:px-6">
              <!-- Now Playing -->
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
                    <p class="truncate text-sm font-bold text-white" dir="auto">{{ currentTrack.title }}</p>
                    <p class="truncate text-xs text-white/40" dir="auto">{{ currentTrack.artistName }}</p>
                  </div>
                  <i class="pi pi-waveform text-lg text-[#1db954]" />
                </div>
              </div>

              <!-- Next Up -->
              <div>
                <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">Next Up</p>
                <div v-if="queueTracks.length" class="space-y-1">
                  <div
                    v-for="(track, idx) in queueTracks"
                    :key="track.id"
                    :draggable="queueTracks.length > 1"
                    class="group flex cursor-grab items-center gap-3 rounded-xl px-3 py-2.5 transition-all hover:bg-white/[0.06]"
                    @dragstart="onDragStart(idx)"
                    @dragover.prevent="onDragOver(idx)"
                    @dragend="onDragEnd"
                    @click="playQueueItem(idx)"
                  >
                    <span class="cursor-grab text-[rgba(255,255,255,0.2)] hover:text-white/50" @click.stop>
                      <i class="pi pi-bars text-xs" />
                    </span>
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
                      <p class="truncate text-sm font-medium text-white/80" dir="auto">{{ track.title }}</p>
                      <p class="truncate text-xs text-white/40" dir="auto">{{ track.artistName }}</p>
                    </div>
                    <span class="text-xs text-white/30 tabular-nums">{{ formatTime(track.durationSeconds) }}</span>
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
                <button
                  v-if="queueTracks.length"
                  type="button"
                  class="mt-4 text-xs font-medium text-white/30 transition-colors hover:text-red-400"
                  @click="clearQueue"
                >
                  Clear Queue
                </button>
              </div>
            </div>
          </div>

          <!-- ─── LYRICS TAB ─── -->
          <div
            v-if="activeTab === 'lyrics'"
            class="relative flex flex-1 flex-col overflow-hidden"
          >
            <div class="pointer-events-none absolute inset-0 bg-black/30" />
            <div class="relative z-10 flex-1 overflow-y-auto px-4 py-6 md:px-8">
              <!-- Loading -->
              <div
                v-if="lyricsLoading"
                class="flex flex-1 items-center justify-center"
              >
                <div class="w-3/4 space-y-4">
                  <div
                    v-for="i in 6"
                    :key="i"
                    class="h-5 animate-pulse rounded bg-white/[0.06]"
                    :style="{ width: `${55 + ((i * 7) % 30)}%` }"
                  />
                </div>
              </div>

              <!-- Error / no lyrics -->
              <div
                v-else-if="lyricsError || !lyricsContent"
                class="flex flex-1 flex-col items-center justify-center gap-3 text-center py-16"
              >
                <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
                  <i class="pi pi-align-left text-3xl text-white/15" />
                </div>
                <p class="text-sm text-white/25">No lyrics available</p>
                <p class="text-xs text-white/15">Lyrics will appear here when available</p>
              </div>

              <!-- Synced lyrics (karaoke mode) -->
              <div
                v-else
                ref="lyricsContainer"
                class="flex flex-col items-center justify-center py-8"
              >
                <div
                  v-for="(line, idx) in parsedLines"
                  :key="idx"
                  ref="lyricLinesRef"
                  class="cursor-pointer px-4 py-2.5 text-center leading-relaxed transition-all duration-500"
                  :class="{
                    'text-[20px] font-semibold text-white drop-shadow-[0_0_20px_rgba(255,255,255,0.1)]': idx === activeLineIndex,
                    'text-[17px] text-[rgba(255,255,255,0.35)]': idx !== activeLineIndex && !isPastLine(idx),
                    'text-[18px] text-[rgba(255,255,255,0.35)]': isPastLine(idx),
                  }"
                  @click="seek(line.timeSeconds)"
                >
                  <!-- Karaoke word-by-word highlighting -->
                  <template v-if="karaokeMode && idx === activeLineIndex && line.words">
                    <span
                      v-for="(word, wIdx) in line.words"
                      :key="wIdx"
                      class="transition-all duration-150"
                      :class="wIdx <= (activeWordMap.get(idx) ?? -1) ? 'text-[#1db954]' : 'text-white/40'"
                    >
                      {{ word.text }}<template v-if="wIdx < line.words.length - 1">&nbsp;</template>
                    </span>
                  </template>
                  <template v-else>
                    <span
                      :class="idx === activeLineIndex ? 'text-white' : ''"
                    >
                      {{ line.text }}
                    </span>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB BAR -->
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

        <!-- Sleep timer menu -->
        <Transition name="fade">
          <div
            v-if="sleepMenuOpen"
            class="absolute bottom-20 left-1/2 z-20 -translate-x-1/2"
          >
            <div class="glass-darker flex items-center gap-1.5 rounded-full px-3 py-1.5 shadow-lg">
              <button
                v-for="opt in sleepOptions"
                :key="opt.value"
                type="button"
                class="rounded-full px-3 py-1.5 text-xs font-medium text-white/50 transition-all hover:bg-white/10 hover:text-white"
                :class="{ '!bg-[#1db954]/15 !text-[#1db954]': sleepTimerMinutes === opt.value }"
                @click="setTimer(opt.value)"
              >
                {{ opt.label }}
              </button>
              <button
                type="button"
                class="ml-1 rounded-full px-2 py-1.5 text-xs text-white/30 transition-colors hover:text-white/60"
                @click="sleepMenuOpen = false"
              >
                <i class="pi pi-times" />
              </button>
            </div>
          </div>
        </Transition>

        <!-- Overflow menu -->
        <div v-if="showOverflow" class="absolute top-14 right-4 z-20">
          <PlayerOverflowMenu
            @close="showOverflow = false"
            @toggle-pip="onTogglePiP"
          />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import { usePlayerControls, usePlayer } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useLyricsApi, type Lyrics } from '@/services/api/lyrics'
import { useRecommendationsApi, type RecommendationTrack } from '@/services/api/recommendation'
import type { PlaybackTrack } from '@/services/api/player'
import type { ParsedLine } from '@/composables/lyrics'
import { parseLRCLines, parsePlainLines } from '@/composables/lyrics'
import { onImgError } from '@/utils/helpers'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'

const props = withDefaults(defineProps<{
  visible: boolean
  initialTab?: 'now-playing' | 'queue' | 'lyrics'
}>(), {
  initialTab: 'now-playing',
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const activeTab = ref<'now-playing' | 'queue' | 'lyrics'>(props.initialTab)

watch(() => props.initialTab, (tab) => {
  activeTab.value = tab
})

const tabs = [
  { key: 'now-playing' as const, label: 'Now Playing' },
  { key: 'queue' as const, label: 'Queue' },
  { key: 'lyrics' as const, label: 'Lyrics' },
]

const tabIndicatorStyle = computed(() => {
  const idx = tabs.findIndex((t) => t.key === activeTab.value)
  const width = 100 / tabs.length
  return {
    width: `${width}%`,
    left: `${idx * width}%`,
  }
})

// --- Player controls ---
const pc = usePlayerControls()
const currentTrack = pc.currentTrack
const isPlaying = pc.isPlaying
const isBuffering = pc.isBuffering
const isLoadingTrack = pc.isLoadingTrack
const currentTime = pc.currentTime
const duration = pc.duration
const progressPercent = pc.progressPercent
const volume = pc.volume
const muted = pc.muted
const shuffleMode = pc.shuffleMode
const repeatMode = pc.repeatMode
const playbackRate = pc.playbackRate
const sleepTimerMinutes = pc.sleepTimerMinutes
const crossfadeDuration = pc.crossfadeDuration
const hasNext = pc.hasNext
const hasPrevious = pc.hasPrevious
const volumeIcon = pc.volumeIcon
const speedLabel = pc.speedLabel

const togglePlayPause = pc.togglePlayPause
const toggleMute = pc.toggleMute
const seekPercent = pc.seekPercent
const seek = pc.seek
const playNext = pc.playNext
const playPrevious = pc.playPrevious
const toggleShuffle = pc.toggleShuffle
const toggleRepeat = pc.toggleRepeat
const setPlaybackRate = pc.setPlaybackRate
const setSleepTimer = pc.setSleepTimer
const clearSleepTimer = pc.clearSleepTimer

const sleepOptions = [
  { value: 5, label: '5m' },
  { value: 15, label: '15m' },
  { value: 30, label: '30m' },
  { value: 45, label: '45m' },
  { value: 60, label: '1h' },
]

const sleepMenuOpen = ref(false)
const showOverflow = ref(false)
const karaokeMode = ref(true)

function setTimer(minutes: number) {
  if (minutes === 0) clearSleepTimer()
  else setSleepTimer(minutes)
  sleepMenuOpen.value = false
}

const speedOptions = [0.5, 0.75, 1, 1.25, 1.5, 2]
function cycleSpeed() {
  const idx = speedOptions.indexOf(playbackRate)
  const nextIdx = (idx + 1) % speedOptions.length
  setPlaybackRate(speedOptions[nextIdx]!)
}

function cycleCrossfade() {
  const p = usePlayer()
  p.crossfadeDuration = p.crossfadeDuration > 0 ? 0 : 5
}

function onOverflowToggle() {
  showOverflow.value = !showOverflow.value
}

function onTogglePiP() {
  // PiP handled by parent via FloatingMiniPlayerHost
}

// --- Queue ---
const queueTracks = computed(() => pc.queue.value as PlaybackTrack[])

const queueIndex = computed(() => {
  if (!currentTrack.value) return -1
  return queueTracks.value.findIndex((t) => t.id === currentTrack.value?.id)
})

const upNextTracks = computed(() => {
  const idx = queueIndex.value
  if (idx < 0) return queueTracks.value.slice(0, 3)
  return queueTracks.value.slice(idx + 1, idx + 4)
})

const dragIndex = ref<number | null>(null)

function onDragStart(index: number) {
  dragIndex.value = index
}

function onDragOver(index: number) {
  if (dragIndex.value === null || dragIndex.value === index) return
  const items = [...queueTracks.value]
  const [moved] = items.splice(dragIndex.value, 1)
  if (!moved) return
  items.splice(index, 0, moved)
  pc.updateQueue(items)
  dragIndex.value = index
}

function onDragEnd() {
  dragIndex.value = null
}

function playQueueItem(index: number) {
  const track = queueTracks.value[index]
  if (!track || track.id === currentTrack.value?.id) return
  pc.setQueueAndPlay(queueTracks.value, index)
}

function removeFromQueue(index: number) {
  const newQueue = [...pc.queue.value]
  newQueue.splice(index, 1)
  pc.updateQueue(newQueue)
}

function clearQueue() {
  pc.updateQueue([])
}

// --- Recommendations ---
const recsApi = useRecommendationsApi()
const recs = shallowRef<RecommendationTrack[]>([])
const recsLoading = ref(false)

async function fetchRecs() {
  if (!currentTrack.value?.id) return
  recsLoading.value = true
  try {
    const res = await recsApi.getSimilar(currentTrack.value.id, { limit: 10 })
    recs.value = res?.items ?? []
  } catch {
    recs.value = []
  } finally {
    recsLoading.value = false
  }
}

// --- Lyrics ---
const lyricsApi = useLyricsApi()
const lyricsData = ref<Lyrics | null>(null)
const lyricsLoading = ref(false)
const lyricsError = ref(false)

watch(
  () => currentTrack.value?.id,
  (id) => {
    if (id) {
      void fetchRecs()
      void fetchLyrics(id)
    }
  },
  { immediate: true },
)

const lyricsContent = computed(() => lyricsData.value?.content || '')
const lyricsType = computed(() => lyricsData.value?.type || 'plain')

const parsedLines = ref<ParsedLine[]>([])

watch(
  [lyricsType, lyricsContent],
  () => {
    if (!lyricsContent.value) {
      parsedLines.value = []
      return
    }
    parsedLines.value =
      lyricsType.value === 'lrc'
        ? parseLRCLines(lyricsContent.value)
        : parsePlainLines(lyricsContent.value)
  },
  { immediate: true },
)

const activeLineIndex = computed(() => {
  const t = currentTime.value
  const lines = parsedLines.value
  for (let i = lines.length - 1; i >= 0; i--) {
    if (t >= lines[i]!.timeSeconds) return i
  }
  return -1
})

function isPastLine(idx: number) {
  return activeLineIndex.value >= 0 && idx < activeLineIndex.value
}

const activeWordMap = computed(() => {
  const map = new Map<number, number>()
  if (!currentTrack.value) return map
  const t = currentTime.value
  for (let i = 0; i < parsedLines.value.length; i++) {
    const line = parsedLines.value[i]!
    if (!line.words || line.words.length === 0) {
      map.set(i, -1)
      continue
    }
    let found = -1
    for (let w = 0; w < line.words.length; w++) {
      const word = line.words[w]!
      const nextWord = line.words[w + 1]
      if (!word) continue
      const start = word.timeSeconds >= 0 ? word.timeSeconds : line.timeSeconds
      const end = nextWord?.timeSeconds ?? line.timeSeconds + 4
      if (t >= start && t < end) {
        found = w
        break
      }
    }
    map.set(i, found)
  }
  return map
})

const lyricsContainer = ref<HTMLElement | null>(null)
const lyricLinesRef = ref<HTMLElement[]>([])

let autoScrollTimer: ReturnType<typeof setTimeout> | null = null

watch(activeLineIndex, (idx) => {
  if (autoScrollTimer) clearTimeout(autoScrollTimer)
  autoScrollTimer = setTimeout(() => {
    if (idx < 0 || !lyricsContainer.value) return
    const target = lyricLinesRef.value[idx]
    target?.scrollIntoView({ block: 'center', behavior: 'smooth' })
  }, 80)
})

async function fetchLyrics(trackId: string) {
  lyricsLoading.value = true
  lyricsError.value = false
  lyricsData.value = null
  try {
    const res = await lyricsApi.getTrackLyrics(trackId)
    lyricsData.value = (res ?? null) as Lyrics | null
  } catch {
    lyricsError.value = true
  } finally {
    lyricsLoading.value = false
  }
}

// --- Computed values ---
const title = computed(() => currentTrack.value?.title || 'No track')
const artistName = computed(() => currentTrack.value?.artistName || '')
const albumTitle = computed(() => currentTrack.value?.albumTitle || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')

const { palette: albumPalette } = useAlbumColors(coverUrl)

const accentColor = computed(() => albumPalette.value.vibrant || '#1db954')

const isLiked = ref(false)
function toggleLike() {
  isLiked.value = !isLiked.value
}

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

const progressPct = computed(() => {
  if (!duration.value) return 0
  return (currentTime.value / duration.value) * 100
})

const repeatTitle = computed(() => {
  if (repeatMode === 'off') return 'Repeat: off'
  if (repeatMode === 'all') return 'Repeat: all'
  return 'Repeat: one'
})

// --- Seek ---
const progressRef = ref<HTMLDivElement | null>(null)
const hoverPos = ref<number | null>(null)

function seekFromEvent(e: MouseEvent) {
  const rect = progressRef.value?.getBoundingClientRect()
  if (!rect) return
  const pct = (e.clientX - rect.left) / rect.width
  if (duration.value) {
    seek(pct * duration.value)
  }
}

function onProgressHover(e: MouseEvent) {
  const rect = progressRef.value?.getBoundingClientRect()
  if (!rect) return
  hoverPos.value = ((e.clientX - rect.left) / rect.width) * 100
}

const hoverTimeLabel = computed(() => {
  if (hoverPos.value === null || !duration.value) return '0:00'
  return formatTime((hoverPos.value / 100) * duration.value)
})

function onVolume(e: Event) {
  pc.setVolume(Number((e.target as HTMLInputElement).value))
}

// --- Background aurora ---
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

// --- Album art ---
const coverSizeClass = 'h-[280px] w-[280px] md:h-[320px] md:w-[320px]'

const artScaleClass = computed(() => {
  if (!currentTrack.value?.id) return ''
  return 'animate-[art-pop_300ms_cubic-bezier(0.34,1.56,0.64,1)]'
})

// --- Navigation ---
const rootEl = ref<HTMLElement | null>(null)

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

// --- Lifecycle ---
onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
  window.addEventListener('popstate', close)
})

onBeforeUnmount(() => {
  window.removeEventListener('popstate', close)
  if (autoScrollTimer) clearTimeout(autoScrollTimer)
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

.fade-enter-active,
.fade-leave-active {
  transition: opacity 200ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes pulse-ring {
  0% { transform: scale(1); opacity: 0.4; }
  50% { transform: scale(1.08); opacity: 0.15; }
  100% { transform: scale(1); opacity: 0.4; }
}

@keyframes play-ring {
  0% { transform: scale(1); opacity: 0.5; }
  100% { transform: scale(1.4); opacity: 0; }
}

@keyframes art-pop {
  from { transform: scale(0.97); opacity: 0.7; }
  to { transform: scale(1); opacity: 1; }
}

@media (prefers-reduced-motion: reduce) {
  .fullscreen-slide-enter-active,
  .fullscreen-slide-leave-active,
  .fade-enter-active,
  .fade-leave-active {
    transition: none;
  }
  .fullscreen-slide-enter-from,
  .fullscreen-slide-leave-to {
    transform: none;
  }
  [class*="pulse-ring"],
  [class*="play-ring"] {
    animation-name: none;
  }
}
</style>
