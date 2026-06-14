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
        dir="rtl"
      >
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

        <div class="relative z-10 flex h-14 items-center justify-between px-4 md:px-6">
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-2 focus-visible:ring-[#1db954]"
            @click="close"
            aria-label="بستن"
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
              aria-label="بیشتر"
              :aria-expanded="showOverflow"
            >
              <i class="pi pi-ellipsis-h text-sm" />
            </button>
          </div>
        </div>

        <div class="relative z-10 flex flex-1 flex-col overflow-hidden">
          <div
            v-if="activeTab === 'now-playing'"
            class="flex flex-1 flex-col items-center overflow-y-auto px-4 pb-4"
          >
            <div class="flex w-full max-w-md flex-1 flex-col items-center justify-center gap-5">
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
                <div
                  v-if="isPlaying"
                  class="pointer-events-none absolute -inset-4 animate-[pulse-ring_2s_ease-in-out_infinite] rounded-full border border-[#1db954]/30"
                />
              </div>

              <div class="flex w-full max-w-md flex-col items-center gap-1 text-center">
                <h2 class="max-w-full truncate text-[22px] font-bold text-white" dir="auto">
                  {{ title }}
                </h2>
                <p class="text-base text-[rgba(255,255,255,0.6)]" dir="auto">
                  {{ artistName }}
                </p>
              </div>

              <div v-if="isPlaying" class="h-12 w-full max-w-md">
                <VisualizerSystem
                  :bar-count="55"
                  mode="spectrum"
                  :color-palette="palette"
                  :is-playing="isPlaying"
                  :sensitivity="0.7"
                  :glow-intensity="0.3"
                  :rounded-bars="true"
                  class="h-full w-full"
                />
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
                  aria-label="قبلی"
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
                  aria-label="بعدی"
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
                  :aria-label="muted ? 'باز کردن صدا' : 'بی صدا'"
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
                  <span v-if="sleepTimerMinutes > 0" class="tabular-nums">{{ sleepTimerLabel }}</span>
                </button>

                <button
                  type="button"
                  class="flex h-8 items-center gap-1.5 rounded-full bg-white/[0.04] px-3 text-xs font-bold transition-all hover:bg-white/10 hover:text-white active:scale-95"
                  :class="{ 'bg-[#1db954]/10 text-[#1db954]': crossfadeDuration > 0 }"
                  @click="cycleCrossfade"
                >
                  <i class="pi pi-arrows-alt text-[10px]" />
                  <span v-if="crossfadeDuration > 0" class="tabular-nums">{{ crossfadeDuration }}s</span>
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

            <div v-if="upNextTracks.length" class="w-full max-w-md mt-4">
              <p class="mb-2 px-1 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">در ادامه</p>
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

          <div
            v-if="activeTab === 'lyrics'"
            class="relative flex flex-1 flex-col overflow-hidden"
          >
            <div class="pointer-events-none absolute inset-0 bg-black/30" />
            <div class="relative z-10 flex-1 overflow-y-auto px-4 py-6 md:px-8">
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

              <div
                v-else-if="lyricsError || (!lyricsContent && !lyricsLoading)"
                class="flex flex-1 flex-col items-center justify-center gap-4 text-center py-16"
              >
                <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
                  <i class="pi pi-align-left text-3xl text-white/15" />
                </div>
                <p class="text-sm text-white/40">متن این آهنگ هنوز اضافه نشده</p>
                <p class="text-xs text-white/20">لطفاً بعداً مراجعه کنید یا در تکمیل متن مشارکت کنید</p>
                <button
                  type="button"
                  class="mt-2 rounded-full bg-white/10 px-5 py-2 text-xs font-medium text-white/60 transition-all hover:bg-white/20 hover:text-white"
                  @click="goToContributions"
                >
                  کمک به ترجمه متن
                </button>
              </div>

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
                  @click="seekLyrics(line.timeSeconds)"
                >
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
                    <span :class="idx === activeLineIndex ? 'text-white' : ''">{{ line.text }}</span>
                  </template>
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
                <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">در حال پخش</p>
                <div class="flex items-center gap-3 rounded-xl border-r-2 border-[#1db954] bg-white/[0.03] px-4 py-3">
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

              <div>
                <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-[rgba(255,255,255,0.35)]">در صف پخش</p>
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
                      aria-label="حذف از صف"
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
                  <p class="text-sm text-white/25">صف پخش خالی است</p>
                </div>
              </div>
            </div>
          </div>
        </div>

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
                class="mr-1 rounded-full px-2 py-1.5 text-xs text-white/30 transition-colors hover:text-white/60"
                @click="sleepMenuOpen = false"
              >
                <i class="pi pi-times" />
              </button>
            </div>
          </div>
        </Transition>

        <Transition name="fade">
          <div
            v-if="showOverflow"
            class="absolute top-14 left-4 z-20"
            @click.self="showOverflow = false"
          >
            <div
              class="min-w-[220px] origin-top-right rounded-2xl border border-white/10 bg-[#121212]/95 p-2 shadow-2xl backdrop-blur-xl"
            >
              <button
                type="button"
                class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white"
                @click="shareTelegram"
              >
                <i class="pi pi-send text-base text-white/40" />
                <span>اشتراک‌گذاری در تلگرام</span>
              </button>

              <button
                type="button"
                class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white"
                @click="addToPlaylist"
              >
                <i class="pi pi-plus-circle text-base text-white/40" />
                <span>افزودن به پلی‌لیست</span>
              </button>

              <button
                type="button"
                class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white"
                @click="goToAlbum"
              >
                <i class="pi pi-disc text-base text-white/40" />
                <span>رفتن به صفحه آلبوم</span>
              </button>

              <button
                type="button"
                class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white"
                @click="goToArtist"
              >
                <i class="pi pi-user text-base text-white/40" />
                <span>رفتن به صفحه خواننده</span>
              </button>

              <div class="my-1 border-t border-white/5" />

              <button
                type="button"
                class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white"
                @click="switchToClassic"
              >
                <i class="pi pi-external-link text-base text-white/40" />
                <span>نمایش کلاسیک</span>
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayer } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useLyricsApi, type Lyrics } from '@/services/api/lyrics'
import { useRecommendationsApi, type RecommendationTrack } from '@/services/api/recommendation'
import type { PlaybackTrack } from '@/services/api/player'
import type { ParsedLine } from '@/composables/lyrics'
import { parseLRCLines, parsePlainLines } from '@/composables/lyrics'
import { onImgError } from '@/utils/helpers'
import VisualizerSystem from './VisualizerSystem.vue'

const visible = defineModel<boolean>('visible', { default: false })

const emit = defineEmits<{
  'toggle-pip': []
}>()

const router = useRouter()
const player = usePlayer()

const { playNext, playPrevious, toggleMute } = player

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
const hasNext = computed(() => player.hasNext.value)
const hasPrevious = computed(() => player.hasPrevious.value)

const { palette } = useAlbumColors(coverUrl)

const accentColor = computed(() => palette.value.vibrant || '#1db954')

const activeTab = ref<'now-playing' | 'lyrics' | 'queue'>('now-playing')

const tabs = [
  { key: 'now-playing' as const, label: 'Now Playing' },
  { key: 'lyrics' as const, label: 'Lyrics' },
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
const sleepMenuOpen = ref(false)
const karaokeMode = ref(true)

const shuffleMode = ref(false)
const repeatMode = ref<'off' | 'all' | 'one'>('off')
const playbackRate = ref(1)

function toggleShuffle() {
  shuffleMode.value = !shuffleMode.value
}

function toggleRepeat() {
  if (repeatMode.value === 'off') repeatMode.value = 'all'
  else if (repeatMode.value === 'all') repeatMode.value = 'one'
  else repeatMode.value = 'off'
}

const repeatTitle = computed(() => {
  if (repeatMode.value === 'off') return 'Repeat: off'
  if (repeatMode.value === 'all') return 'Repeat: all'
  return 'Repeat: one'
})

const speedOptions = [0.5, 0.75, 1, 1.25, 1.5, 2]
const speedLabel = computed(() => `${playbackRate.value}x`)

function cycleSpeed() {
  const idx = speedOptions.indexOf(playbackRate.value)
  const nextIdx = (idx + 1) % speedOptions.length
  playbackRate.value = speedOptions[nextIdx]!
}

const sleepTimerMinutes = ref(0)
let sleepTimerHandle: ReturnType<typeof setTimeout> | null = null
let sleepCountdownHandle: ReturnType<typeof setInterval> | null = null
const sleepRemainingSeconds = ref(0)

const sleepTimerLabel = computed(() => {
  if (sleepRemainingSeconds.value > 0) {
    const m = Math.floor(sleepRemainingSeconds.value / 60)
    const s = sleepRemainingSeconds.value % 60
    return `${m}:${String(s).padStart(2, '0')}`
  }
  return `${sleepTimerMinutes.value}m`
})

const sleepOptions = [
  { value: 5, label: '۵ دقیقه' },
  { value: 15, label: '۱۵ دقیقه' },
  { value: 30, label: '۳۰ دقیقه' },
  { value: 45, label: '۴۵ دقیقه' },
  { value: 60, label: '۱ ساعت' },
]

function setTimer(minutes: number) {
  clearSleepTimer()
  if (minutes === 0) return
  sleepTimerMinutes.value = minutes
  sleepRemainingSeconds.value = minutes * 60
  sleepTimerHandle = setTimeout(() => {
    player.pause()
    clearSleepTimer()
  }, minutes * 60 * 1000)
  sleepCountdownHandle = setInterval(() => {
    if (sleepRemainingSeconds.value > 0) {
      sleepRemainingSeconds.value--
    } else {
      clearSleepTimer()
    }
  }, 1000)
  sleepMenuOpen.value = false
}

function clearSleepTimer() {
  if (sleepTimerHandle) {
    clearTimeout(sleepTimerHandle)
    sleepTimerHandle = null
  }
  if (sleepCountdownHandle) {
    clearInterval(sleepCountdownHandle)
    sleepCountdownHandle = null
  }
  sleepTimerMinutes.value = 0
  sleepRemainingSeconds.value = 0
}

const crossfadeDuration = ref(0)

function cycleCrossfade() {
  crossfadeDuration.value = crossfadeDuration.value > 0 ? 0 : 5
}

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
  player.updateQueue(items)
  dragIndex.value = index
}

function onDragEnd() {
  dragIndex.value = null
}

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

const volumeIcon = computed(() => {
  if (muted.value || volume.value === 0) return 'pi pi-volume-off'
  if (volume.value < 0.5) return 'pi pi-volume-down'
  return 'pi pi-volume-up'
})

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

function onTogglePiP() {
  emit('toggle-pip')
}

function close() {
  visible.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
  if (e.key === ' ' && e.target === e.currentTarget) { e.preventDefault(); togglePlayPause() }
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

const lyricsApi = useLyricsApi()
const lyricsData = ref<Lyrics | null>(null)
const lyricsLoading = ref(false)
const lyricsError = ref(false)

watch(
  () => currentTrack.value?.id,
  (id) => {
    lyricsData.value = null
    if (id) {
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

function seekLyrics(seconds: number) {
  if (duration.value) {
    player.seek(seconds)
  }
}

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

watch(
  () => currentTrack.value?.id,
  (id) => {
    if (id) void fetchRecs()
  },
  { immediate: false },
)

function shareTelegram() {
  if (currentTrack.value) {
    const text = `${currentTrack.value.title} - ${currentTrack.value.artistName}\n${window.location.origin}/track/${currentTrack.value.id}`
    window.open(`https://t.me/share/url?url=${encodeURIComponent(window.location.origin + '/track/' + currentTrack.value.id)}&text=${encodeURIComponent(text)}`, '_blank')
  }
  showOverflow.value = false
}

function addToPlaylist() {
  showOverflow.value = false
}

function goToAlbum() {
  if (currentTrack.value?.albumTitle) {
    router.push(`/album/${currentTrack.value.id}`)
  }
  showOverflow.value = false
}

function goToArtist() {
  if (currentTrack.value?.artistName) {
    router.push(`/artist/${currentTrack.value.id}`)
  }
  showOverflow.value = false
}

function goToContributions() {
  router.push('/contributions')
}

function switchToClassic() {
  localStorage.setItem('player-redesigned-override', 'false')
  showOverflow.value = false
  close()
  window.location.reload()
}

onMounted(() => {
  nextTick().then(() => rootEl.value?.focus())
  window.addEventListener('popstate', close)
})

onBeforeUnmount(() => {
  window.removeEventListener('popstate', close)
  if (autoScrollTimer) clearTimeout(autoScrollTimer)
  clearSleepTimer()
})
</script>

<style scoped>
.fullscreen-slide-enter-active {
  transition: transform 400ms cubic-bezier(0.19, 1, 0.22, 1), opacity 350ms ease;
}
.fullscreen-slide-leave-active {
  transition: transform 300ms ease-in, opacity 200ms ease;
}
.fullscreen-slide-enter-from {
  transform: translateY(100%);
  opacity: 0;
}
.fullscreen-slide-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 200ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes art-pop {
  from { transform: scale(0.97); opacity: 0.7; }
  to { transform: scale(1); opacity: 1; }
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
    opacity: 1;
  }
  [class*="pulse-ring"],
  [class*="play-ring"] {
    animation-name: none;
  }
}
</style>
