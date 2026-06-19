<template>
  <div :key="String(route.params.id)" class="mx-auto w-full max-w-4xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="hero" />
    </div>

    <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
      <i aria-hidden="true" class="pi pi-exclamation-circle text-4xl text-slate-500" />
      <h2 class="text-xl font-bold text-white">Track not found</h2>
      <RouterLink to="/" class="text-sm font-medium text-[#1db954] underline underline-offset-2">
        Go home
      </RouterLink>
    </div>

    <template v-else-if="track">
      <div class="flex flex-col items-center">
        <!-- Cover art with spinning vinyl effect -->
        <div class="group relative">
          <div
            class="relative h-72 w-72 overflow-hidden rounded-full bg-white/[0.06] shadow-2xl ring-1 ring-white/10 md:h-80 md:w-80"
            :class="{ 'animate-spin-slow': isPlaying }"
          >
            <img
              v-if="track.cover_url"
              :src="track.cover_url"
              :alt="track.title"
              loading="lazy"
              class="h-full w-full object-cover"
              @error="onImgError"
            />
            <div v-else class="flex h-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-compact-disc text-6xl text-slate-500" />
            </div>

            <!-- Center pin -->
            <div class="absolute inset-0 flex items-center justify-center">
              <div
                class="h-10 w-10 rounded-full bg-black/60 shadow-lg ring-2 ring-white/20 backdrop-blur-sm"
              />
            </div>
          </div>

          <!-- Equalizer ring -->
          <div v-if="isPlaying" class="absolute -inset-3">
            <svg class="h-full w-full -rotate-90" viewBox="0 0 100 100">
              <circle
                v-for="(bar, i) in 12"
                :key="i"
                cx="50"
                cy="50"
                :r="44 + (i % 3) * 2"
                fill="none"
                :stroke="'#1db954'"
                :stroke-width="1.5"
                :opacity="0.15 + (i % 4) * 0.2"
                :stroke-dasharray="`${6 + (i % 3) * 4} ${12 - (i % 3) * 2}`"
                class="equalizer-ring"
                :style="{ animationDelay: `${i * 0.08}s` }"
              />
            </svg>
          </div>
        </div>

        <!-- Track Info -->
        <div class="mt-8 text-center">
          <div class="flex items-center justify-center gap-2">
            <h1 class="text-3xl font-black text-white">{{ track.title }}</h1>
            <span
              v-if="track.explicit"
              class="inline-flex h-5 w-5 items-center justify-center rounded bg-white/20 text-[10px] font-bold tracking-wide text-white"
              title="Explicit"
              >E</span
            >
          </div>
          <div class="mt-2 flex items-center justify-center gap-2 text-sm text-slate-400">
            <RouterLink
              v-if="track.artist_name"
              :to="`/artist/${track.artist_id}`"
              class="font-bold text-white underline underline-offset-2 transition hover:text-[#1db954]"
            >
              {{ track.artist_name }}
            </RouterLink>
            <span v-if="track.album_title"> • </span>
            <RouterLink
              v-if="track.album_title"
              :to="`/album/${track.album_id}`"
              class="transition hover:text-white"
            >
              {{ track.album_title }}
            </RouterLink>
          </div>

          <!-- Genres -->
          <div
            v-if="genreList.length > 0"
            class="mt-3 flex flex-wrap items-center justify-center gap-1.5"
          >
            <span
              v-for="g in genreList"
              :key="g.id"
              class="rounded-full border border-white/10 bg-white/[0.06] px-3 py-0.5 text-xs font-medium text-slate-300 transition hover:border-[#1db954]/30 hover:text-white"
            >
              {{ g.name }}
            </span>
          </div>

          <!-- Play count -->
          <p v-if="track.play_count > 0" class="mt-3 text-xs text-slate-500">
            {{ formatPlayCount(track.play_count) }} plays
          </p>
        </div>

        <!-- Main Controls -->
        <div class="mt-8 flex items-center gap-6">
          <button
            type="button"
            aria-label="Previous track"
            class="flex h-10 w-10 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
            @click="player.playPrevious"
          >
            <i aria-hidden="true" class="pi pi-step-backward text-lg" />
          </button>

          <button
            type="button"
            aria-label="Toggle play"
            class="relative flex h-16 w-16 items-center justify-center rounded-full bg-white text-black shadow-2xl transition hover:scale-105 hover:bg-[#1db954] hover:text-white"
            @click="togglePlay"
          >
            <i
              :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'"
              class="ml-0.5 text-2xl"
            />
          </button>

          <button
            type="button"
            aria-label="Next track"
            class="flex h-10 w-10 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
            @click="player.playNext"
          >
            <i aria-hidden="true" class="pi pi-step-forward text-lg" />
          </button>
        </div>

        <!-- Progress bar -->
        <div class="mt-6 flex w-full max-w-md items-center gap-3">
          <span class="w-10 text-right text-xs font-medium text-slate-500 tabular-nums">
            {{ formatTime(player.currentTime.value) }}
          </span>
          <input
            type="range"
            min="0"
            max="100"
            step="0.1"
            aria-label="Seek"
            class="player-range flex-1"
            :style="{ '--range-progress': `${Number(player.progressPercent.value || 0)}%` }"
            :value="player.progressPercent.value"
            @input="onSeek"
          />
          <span class="w-10 text-xs font-medium text-slate-500 tabular-nums">
            {{ formatTime(player.duration.value || track.duration_seconds || 0) }}
          </span>
        </div>

        <!-- Secondary Controls -->
        <div class="mt-6 flex items-center gap-6">
          <button
            type="button"
            class="flex items-center gap-2 text-sm font-medium transition"
            :class="isLiked ? 'text-[#1db954]' : 'text-slate-400 hover:text-white'"
            @click="toggleLike"
          >
            <i aria-hidden="true" :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-lg" />
            {{ isLiked ? 'Liked' : 'Like' }}
          </button>

          <button
            type="button"
            class="flex items-center gap-2 text-sm font-medium text-slate-400 transition hover:text-white"
            @click="player.toggleShuffle"
          >
            <i
              class="pi pi-sort-alt text-lg"
              :class="{ 'text-[#1db954]': player.shuffleMode }"
            />
            Shuffle
          </button>

          <button
            type="button"
            class="relative flex items-center gap-2 text-sm font-medium transition"
            :class="
              player.repeatMode !== 'off'
                ? 'text-[#1db954]'
                : 'text-slate-400 hover:text-white'
            "
            @click="player.toggleRepeat"
          >
            <i aria-hidden="true" class="pi pi-refresh text-lg" />
            <span
              v-if="player.repeatMode === 'one'"
              class="absolute -top-1 -right-3 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-[#1db954] text-[8px] font-bold text-black"
              >1</span
            >
            Repeat
          </button>

          <button
            type="button"
            class="flex items-center gap-2 text-sm font-medium text-slate-400 transition hover:text-white"
            @click="showQueue = true"
          >
            <i aria-hidden="true" class="pi pi-list text-lg" />
            Queue
          </button>
        </div>
      </div>

      <!-- Waveform -->
      <section class="mx-auto mt-16 w-full max-w-lg">
        <div class="relative h-32 overflow-hidden rounded-2xl bg-white/[0.03]">
          <div class="flex h-full items-end justify-center gap-[3px] px-4 pb-3">
            <div
              v-for="i in 80"
              :key="i"
              class="waveform-bar w-[3px] rounded-full"
              :class="{
                'bg-[#1db954]': isPlaying && i > 30 && i < 50,
                'bg-white/20': !(isPlaying && i > 30 && i < 50),
              }"
              :style="{
                height: `${getWaveHeight(i)}%`,
                animationDelay: isPlaying ? `${i * 0.04}s` : '0s',
              }"
            />
          </div>
          <!-- Center play overlay on hover -->
          <div
            class="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 transition hover:opacity-100"
          >
            <div
              class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954]/90 text-black shadow-xl backdrop-blur-sm"
            >
              <i aria-hidden="true" class="pi pi-play-fill text-lg" />
            </div>
          </div>
        </div>
      </section>

      <!-- Lyrics -->
      <section class="mt-10">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-bold text-white">Lyrics</h2>
          <button
            v-if="lyrics && lyrics.content"
            type="button"
            class="flex items-center gap-1.5 rounded-full bg-white/10 px-3 py-1.5 text-xs font-bold transition hover:bg-white/15"
            :class="karaokeActive ? 'bg-[#1db954]/15 text-[#1db954]' : 'text-white/60'"
            @click="karaokeActive = !karaokeActive"
          >
            <i aria-hidden="true" class="pi pi-mic text-[10px]" />
            Karaoke
          </button>
        </div>
        <div
          v-if="karaokeActive && lyrics?.content"
          class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.03] backdrop-blur-xl"
          style="height: 400px"
        >
          <KaraokeLyrics
            :content="lyrics.content"
            :type="lyrics.type || 'plain'"
            :current-time="player.currentTime.value"
            :loading="loading"
            :karaoke="true"
            @seek="player.seek"
          />
        </div>
        <div
          v-else
          class="rounded-2xl border border-white/[0.06] bg-white/[0.03] p-6 backdrop-blur-xl"
        >
          <LyricsDisplay :lyrics="lyrics" :loading="loading" :error="!lyrics && !loading" />
        </div>
      </section>

      <!-- Credits -->
      <section v-if="trackArtists.length > 0 || trackCredits.length > 0" class="mt-14">
        <div class="relative mb-8">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-white/[0.06]" />
          </div>
          <div class="relative flex justify-center">
            <span class="bg-[#0A0A0F] px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
              Credits
            </span>
          </div>
        </div>
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <!-- Artists -->
          <div v-if="trackArtists.length > 0">
            <h3 class="mb-4 text-xs font-bold tracking-[0.2em] text-white/30 uppercase">Artists</h3>
            <div class="space-y-2">
              <div
                v-for="a in trackArtists"
                :key="a.artistId"
                class="group flex items-center gap-3 rounded-2xl border border-white/[0.04] bg-white/[0.02] px-4 py-3 transition hover:border-white/[0.08] hover:bg-white/[0.04]"
              >
                <div
                  class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-white/[0.08] to-white/[0.02] text-sm font-bold text-white/70 ring-1 ring-white/[0.04]"
                >
                  {{ a.name.charAt(0).toUpperCase() }}
                </div>
                <div>
                  <RouterLink
                    :to="`/artist/${a.artistId}`"
                    class="text-sm font-semibold text-white transition group-hover:text-[#1db954]"
                  >
                    {{ a.name }}
                  </RouterLink>
                  <p
                    v-if="a.role && !['main', 'primary'].includes(a.role)"
                    class="mt-0.5 text-xs text-white/40 capitalize"
                  >
                    {{ a.role }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Credits -->
          <div v-if="trackCredits.length > 0">
            <h3 class="mb-4 text-xs font-bold tracking-[0.2em] text-white/30 uppercase">Production</h3>
            <div class="space-y-2">
              <div
                v-for="c in trackCredits"
                :key="c.id"
                class="group flex items-center gap-3 rounded-2xl border border-white/[0.04] bg-white/[0.02] px-4 py-3 transition hover:border-white/[0.08] hover:bg-white/[0.04]"
              >
                <div
                  class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-white/[0.08] to-white/[0.02] text-sm font-bold text-white/70 ring-1 ring-white/[0.04]"
                >
                  {{ c.artistName.charAt(0).toUpperCase() }}
                </div>
                <div>
                  <p class="text-sm font-semibold text-white">{{ c.artistName }}</p>
                  <p class="mt-0.5 text-xs text-white/40 capitalize">{{ c.creditType }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Section divider -->
      <div v-if="trackArtists.length > 0 || trackCredits.length > 0 || similarTracks.length" class="relative mt-14">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-white/[0.06]" />
        </div>
      </div>

      <!-- Similar Tracks -->
      <section v-if="similarTracks.length" class="mt-14" aria-live="polite">
        <div class="relative mb-8">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-white/[0.06]" />
          </div>
          <div class="relative flex justify-center">
            <span class="bg-[#0A0A0F] px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
              You might like
            </span>
          </div>
        </div>
        <div class="space-y-1">
            <div
              v-for="(st, index) in similarTracks"
              :key="st.id"
              role="button"
              tabindex="0"
              class="group flex cursor-pointer items-center gap-4 rounded-2xl px-4 py-3 transition-all duration-200 hover:bg-white/[0.04]"
              @click="playSimilar(st, index)"
              @keydown.enter="playSimilar(st, index)"
              @keydown.space.prevent="playSimilar(st, index)"
            >
            <span
              class="flex w-8 items-center justify-center text-center text-sm tabular-nums text-white/20 group-hover:hidden"
            >
              {{ String(index + 1).padStart(2, '0') }}
            </span>
            <span class="hidden w-8 items-center justify-center group-hover:flex">
              <i aria-hidden="true" class="pi pi-play-fill text-xs text-white" />
            </span>
            <div class="relative h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10 ring-1 ring-white/[0.04]">
              <img
                v-if="st.cover_url"
                :src="st.cover_url"
                :alt="st.title"
                loading="lazy"
                class="h-full w-full object-cover"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-white">{{ st.title }}</p>
              <p class="mt-0.5 truncate text-xs text-white/40">{{ st.artist_name }}</p>
            </div>
            <span class="text-xs tabular-nums text-white/25 group-hover:text-white/50">{{
              formatTime(st.duration_seconds)
            }}</span>
          </div>
        </div>
      </section>
    </template>

    <!-- Queue Panel -->
    <QueuePanel v-model:visible="showQueue" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useTrack } from '@/composables/catalog/useTrack'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'
import { LyricsDisplay, QueuePanel, KaraokeLyrics } from '@/components/music'

const route = useRoute()
const trackId = String(route.params.id)
const {
  track,
  similarTracks,
  lyrics,
  trackArtists: _trackArtists,
  trackCredits: _trackCredits,
  genreList,
  isLiked,
  loading,
  error,
  fetchTrack,
  toggleLike,
} = useTrack(trackId)

const trackArtists = _trackArtists as Record<string, unknown>[]
const trackCredits = _trackCredits as Record<string, unknown>[]
const player = usePlayer()
const playerApi = usePlayerApi()
const showQueue = ref(false)
const karaokeActive = ref(false)

const isPlaying = player.isPlaying

onMounted(() => {
  fetchTrack()
})

function togglePlay() {
  if (!track.value) return
  const pb: PlaybackTrack = {
    id: trackId,
    title: track.value.title,
    artistName: track.value.artist_name || 'Unknown',
    albumTitle: track.value.album_title || null,
    coverUrl: track.value.cover_url || null,
    durationSeconds: track.value.duration_seconds ?? null,
    streamUrl: track.value.audio_url || playerApi.getTrackStreamUrl(trackId),
  }

  if (isPlaying.value && player.currentTrack.value?.id === trackId) {
    player.pause()
  } else {
    player.playTrack(pb)
  }
}

function playSimilar(
  st: {
    id: string
    title: string
    artist_name?: string | null
    cover_url?: string | null
    duration_seconds?: number | null
  },
  index: number,
) {
  const queue: PlaybackTrack[] = similarTracks.value.map((t) => ({
    id: String(t.id),
    title: t.title,
    artistName: t.artist_name || 'Unknown',
    albumTitle: t.album_title || null,
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration_seconds ?? null,
    streamUrl: t.audio_url || playerApi.getTrackStreamUrl(String(t.id)),
  }))

  player.setQueueAndPlay(queue, index)
}

function onSeek(event: Event) {
  const target = event.target as HTMLInputElement
  player.seekPercent(Number(target.value))
}

function formatTime(seconds?: number | null) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

function formatPlayCount(count: number) {
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}

const waveHeights = Array.from({ length: 80 }, () => Math.random())
function getWaveHeight(i: number) {
  return 15 + waveHeights[i]! * 60
}
</script>

<style scoped>
@keyframes spin-slow {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin-slow {
  animation: spin-slow 6s linear infinite;
}

@keyframes eq-ring-pulse {
  0%,
  100% {
    opacity: 0.2;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.02);
  }
}

.equalizer-ring {
  animation: eq-ring-pulse 1.2s ease-in-out infinite alternate;
}

.player-range {
  --range-progress: 0%;
  width: 100%;
  height: 18px;
  cursor: pointer;
  appearance: none;
  background: transparent;
}

.player-range::-webkit-slider-runnable-track {
  height: 4px;
  border-radius: 999px;
  background: linear-gradient(
    to right,
    #1db954 0%,
    #1db954 var(--range-progress),
    rgba(255, 255, 255, 0.12) var(--range-progress),
    rgba(255, 255, 255, 0.12) 100%
  );
}

.player-range::-webkit-slider-thumb {
  width: 14px;
  height: 14px;
  margin-top: -5px;
  border-radius: 999px;
  appearance: none;
  background: #fff;
  box-shadow: 0 0 16px rgba(29, 185, 84, 0.6);
  opacity: 0;
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.player-range:hover::-webkit-slider-thumb,
.player-range:active::-webkit-slider-thumb {
  opacity: 1;
}

.player-range:active::-webkit-slider-thumb {
  transform: scale(1.25);
}

.player-range::-moz-range-track {
  height: 4px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
}

.player-range::-moz-range-progress {
  height: 4px;
  border-radius: 999px;
  background: #1db954;
}

.player-range::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border: 0;
  border-radius: 999px;
  background: #fff;
  box-shadow: 0 0 16px rgba(29, 185, 84, 0.6);
  opacity: 0;
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.player-range:hover::-moz-range-thumb,
.player-range:active::-moz-range-thumb {
  opacity: 1;
}

.waveform-bar {
  animation: wave-pulse 800ms ease-in-out infinite alternate;
}

.waveform-bar:nth-child(even) {
  animation-delay: 0.2s;
}

.waveform-bar:nth-child(3n) {
  animation-delay: 0.4s;
}

@keyframes wave-pulse {
  0% {
    transform: scaleY(0.6);
    opacity: 0.5;
  }
  100% {
    transform: scaleY(1);
    opacity: 1;
  }
}
</style>
