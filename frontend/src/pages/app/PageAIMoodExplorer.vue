<template>
  <div class="relative mx-auto min-h-screen w-full pb-36">
    <!-- ── Ambient Mood Background ── -->
    <div class="pointer-events-none fixed inset-0 transition-all duration-1000" aria-hidden="true">
      <div
        class="absolute inset-0 transition-opacity duration-1000"
        :class="moodBgClass"
      />
      <div
        class="absolute -top-1/3 -end-1/4 h-125 w-125 rounded-full opacity-30 blur-[150px] transition-all duration-1000"
        :class="moodSpotClass"
      />
      <div
        class="absolute -bottom-1/3 -start-1/4 h-100 w-100 rounded-full opacity-20 blur-[120px] transition-all duration-1000"
        :class="moodSpotClass2"
      />
    </div>

    <div class="relative z-10 px-4 pt-8 md:px-6 lg:px-8">
      <div class="mx-auto max-w-6xl">
        <!-- ════════════════════════════════════ -->
        <!-- HEADER                               -->
        <!-- ════════════════════════════════════ -->
        <div class="mb-10">
          <div class="mb-3 flex items-center gap-2">
            <span
              class="inline-flex items-center gap-1.5 rounded-full border border-spotify/20 bg-spotify/10 px-3 py-1 text-[10px] font-bold tracking-[0.2em] text-spotify uppercase"
            >
              <i aria-hidden="true" class="pi pi-sparkles text-[10px]" />
              AI-Powered
            </span>
          </div>
          <h1 class="text-4xl font-black text-white md:text-5xl">Mood Explorer</h1>
          <p class="mt-3 max-w-xl text-base text-white/50">
            Pick a feeling and let AI find tracks that match the vibe — by energy, valence, and emotional profile.
          </p>
        </div>

        <!-- ════════════════════════════════════ -->
        <!-- MOOD GRID — Gradient Cards           -->
        <!-- ════════════════════════════════════ -->
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-5">
          <button
            v-for="mood in moodOptions"
            :key="mood.value"
            type="button"
            class="group relative overflow-hidden rounded-2xl border p-5 text-center transition-all duration-500"
            :class="[
              selectedMood === mood.value
                ? moodSelectedBorder(mood.value)
                : 'border-white/6 bg-white/2 hover:border-white/12',
            ]"
            @click="toggleMood(mood.value)"
          >
            <!-- Hover/Active gradient -->
            <div
              class="pointer-events-none absolute inset-0 bg-linear-to-br opacity-0 transition-opacity duration-500 group-hover:opacity-100"
              :class="getMoodGradient(mood.value)"
            />
            <!-- Selected glow -->
            <div
              v-if="selectedMood === mood.value"
              class="pointer-events-none absolute inset-0 opacity-20 blur-2xl transition-opacity duration-500"
              :class="moodGlowClass(mood.value)"
            />

            <div class="relative">
              <div
                class="mx-auto mb-3 flex h-14 w-14 items-center justify-center rounded-xl text-xl transition-all duration-500"
                :class="
                  selectedMood === mood.value
                    ? moodIconActiveBg(mood.value)
                    : 'bg-white/6 text-white/40 group-hover:bg-white/10'
                "
              >
                <i aria-hidden="true" :class="mood.icon" />
              </div>
              <span
                class="text-sm font-bold transition-colors duration-300"
                :class="selectedMood === mood.value ? 'text-white' : 'text-white/60 group-hover:text-white'"
              >
                {{ mood.label }}
              </span>
              <p v-if="moodDescriptions[mood.value]" class="mt-1 text-[10px] leading-relaxed text-white/30">
                {{ moodDescriptions[mood.value] }}
              </p>
            </div>
          </button>
        </div>

        <!-- ════════════════════════════════════ -->
        <!-- RESULTS SECTION                      -->
        <!-- ════════════════════════════════════ -->
        <div class="mt-8">
          <!-- Section header -->
          <div class="mb-4 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <h2 class="text-sm font-bold text-white">
                <template v-if="selectedMood">
                  <span :class="moodTextClass(selectedMood)">{{ getMoodLabel(selectedMood) }}</span>
                  Tracks
                </template>
                <template v-else>
                  Pick a mood to explore
                </template>
              </h2>
              <span
                v-if="tracks.length && selectedMood"
                class="rounded-full bg-white/6 px-2 py-0.5 text-[10px] font-medium text-white/40"
              >
                {{ tracks.length }} tracks
              </span>
            </div>
          </div>

          <!-- Loading state -->
          <div v-if="loading" class="space-y-2">
            <div class="flex items-center gap-3 rounded-2xl bg-white/2 px-5 py-4">
              <i aria-hidden="true" class="pi pi-spin pi-sparkles text-spotify" />
              <div>
                <span class="text-sm font-medium text-white">AI is analyzing your mood</span>
                <p class="text-xs text-white/40">Matching tracks by energy and emotional profile...</p>
              </div>
            </div>
            <div v-for="i in 5" :key="i" class="flex animate-pulse items-center gap-3 rounded-2xl bg-white/2 px-4 py-3">
              <div class="h-10 w-10 shrink-0 rounded-lg bg-white/6" />
              <div class="flex-1 space-y-2">
                <div class="h-3 w-3/4 rounded bg-white/6" />
                <div class="h-2 w-1/2 rounded bg-white/4" />
              </div>
              <div class="h-2 w-14 rounded bg-white/4" />
              <div class="h-3 w-12 rounded bg-white/4" />
            </div>
          </div>

          <!-- Track list -->
          <div
            v-else-if="tracks.length"
            class="overflow-hidden rounded-2xl border border-white/6 bg-white/2 backdrop-blur-xs"
          >
            <div
              v-for="(track, index) in tracks"
              :key="String(track.id)"
              role="button"
              tabindex="0"
              class="group flex items-center gap-3 px-4 py-2.5 transition hover:bg-white/4"
              :style="{ animationDelay: `${index * 50}ms` }"
              @click="playTrack(index)"
              @keydown.enter="playTrack(index)"
            >
              <!-- Number -->
              <span class="flex w-6 items-center justify-center">
                <span class="text-xs font-bold text-white/20 group-hover:hidden">{{ index + 1 }}</span>
                <i aria-hidden="true" class="pi pi-play-fill hidden text-xs text-white group-hover:block" />
              </span>

              <!-- Cover -->
              <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/5 ring-1 ring-white/6">
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url as string"
                  :alt="String(track.title ?? '')"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-white/20" />
                </div>
              </div>

              <!-- Info -->
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p v-if="track.artist" class="truncate text-xs text-white/40">{{ track.artist }}</p>
              </div>

              <!-- Mood energy bar -->
              <div v-if="track.mood" class="hidden items-center gap-1.5 sm:flex">
                <div class="flex h-1.5 w-14 overflow-hidden rounded-full bg-white/10">
                  <div
                    class="h-full rounded-full transition-all duration-500"
                    :class="getEnergyBarColor((track.mood as Record<string, unknown>).energy as number || 0)"
                    :style="{ width: `${Math.min(((track.mood as Record<string, unknown>).energy as number || 0) * 100, 100)}%` }"
                  />
                </div>
                <span class="w-8 text-right text-[10px] text-white/30 tabular-nums">
                  {{ Math.round(((track.mood as Record<string, unknown>).energy as number || 0) * 100) }}%
                </span>
              </div>

              <!-- Mood tags -->
              <div
                v-if="track.mood && (track.mood as Record<string, unknown>).mood_tags"
                class="hidden items-center gap-1 lg:flex"
              >
                <span
                  v-for="tag in ((track.mood as Record<string, unknown>).mood_tags as Array<{name: string}>)?.slice(0, 2)"
                  :key="tag.name"
                  class="rounded-full bg-white/6 px-2 py-0.5 text-[9px] font-medium text-white/40"
                >
                  {{ tag.name }}
                </span>
              </div>

              <!-- Duration -->
              <span class="shrink-0 text-xs text-white/30 tabular-nums">
                {{ formatTime(track.duration as number | undefined) }}
              </span>
            </div>
          </div>

          <!-- Empty / Pick a mood state -->
          <div
            v-else-if="!selectedMood"
            class="flex flex-col items-center gap-4 rounded-2xl border border-dashed border-white/6 bg-white/2 px-6 py-20 text-center"
          >
            <div
              class="flex h-16 w-16 items-center justify-center rounded-2xl bg-linear-to-br from-spotify/20 to-blue-500/20"
            >
              <i aria-hidden="true" class="pi pi-heart text-2xl text-spotify" />
            </div>
            <h3 class="text-lg font-bold text-white">What's your mood?</h3>
            <p class="max-w-xs text-sm text-white/40">
              Choose a mood above and AI will find tracks that match the feeling.
            </p>
          </div>

          <!-- No results state -->
          <div
            v-else
            class="flex flex-col items-center gap-4 rounded-2xl border border-dashed border-white/6 bg-white/2 px-6 py-20 text-center"
          >
            <div
              class="flex h-16 w-16 items-center justify-center rounded-2xl bg-linear-to-br from-white/5 to-white/2"
            >
              <i aria-hidden="true" class="pi pi-inbox text-2xl text-white/20" />
            </div>
            <h3 class="text-lg font-bold text-white">No tracks found</h3>
            <p class="max-w-xs text-sm text-white/40">
              No tracks match this mood yet. Try a different mood.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAIApi } from '@/services/api/ai'
import { MOOD_OPTIONS } from '@/services/api/ai/types'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { mapToPlaybackTracks } from '@/factories/playbackTrack'
import { useToast } from 'primevue/usetoast'

const aiApi = useAIApi()
const player = usePlayer()
const toast = useToast()

const selectedMood = ref('')
const loading = ref(false)
const tracks = ref<Record<string, unknown>[]>([])

const moodOptions = MOOD_OPTIONS

const moodDescriptions: Record<string, string> = {
  energetic: 'High energy, fast tempo',
  happy: 'Upbeat, positive vibes',
  chill: 'Laid back, easy going',
  calm: 'Peaceful, slow tempo',
  sad: 'Melancholic, emotional',
  focus: 'Concentration flow',
  romantic: 'Warm, intimate mood',
  intense: 'Powerful, dramatic',
  confident: 'Bold, self-assured',
  sleep: 'Gentle, restful sounds',
}

function getMoodGradient(mood: string): string {
  const gradients: Record<string, string> = {
    energetic: 'from-orange-500/15 via-red-500/10 to-transparent',
    happy: 'from-yellow-400/15 via-amber-500/10 to-transparent',
    chill: 'from-cyan-400/15 via-teal-500/10 to-transparent',
    calm: 'from-blue-400/15 via-indigo-500/10 to-transparent',
    sad: 'from-indigo-400/15 via-violet-500/10 to-transparent',
    focus: 'from-emerald-400/15 via-green-500/10 to-transparent',
    romantic: 'from-pink-400/15 via-rose-500/10 to-transparent',
    intense: 'from-red-500/15 via-orange-600/10 to-transparent',
    confident: 'from-purple-400/15 via-fuchsia-500/10 to-transparent',
    sleep: 'from-slate-400/15 via-blue-500/10 to-transparent',
  }
  return gradients[mood] || 'from-spotify/15 via-emerald-500/10 to-transparent'
}

function moodBgClass(): string {
  if (!selectedMood.value) return 'bg-black/0'
  const bg: Record<string, string> = {
    energetic: 'bg-linear-to-b from-surface-base via-orange-950/20 to-black',
    happy: 'bg-linear-to-b from-surface-base via-amber-950/15 to-black',
    chill: 'bg-linear-to-b from-surface-base via-teal-950/20 to-black',
    calm: 'bg-linear-to-b from-surface-base via-blue-950/20 to-black',
    sad: 'bg-linear-to-b from-surface-base via-indigo-950/20 to-black',
    focus: 'bg-linear-to-b from-surface-base via-emerald-950/20 to-black',
    romantic: 'bg-linear-to-b from-surface-base via-rose-950/20 to-black',
    intense: 'bg-linear-to-b from-surface-base via-red-950/20 to-black',
    confident: 'bg-linear-to-b from-surface-base via-purple-950/20 to-black',
    sleep: 'bg-linear-to-b from-surface-base via-slate-950/20 to-black',
  }
  return bg[selectedMood.value] || 'bg-black/0'
}

function moodSpotClass(): string {
  const spots: Record<string, string> = {
    energetic: 'bg-orange-500',
    happy: 'bg-amber-400',
    chill: 'bg-teal-400',
    calm: 'bg-blue-400',
    sad: 'bg-indigo-400',
    focus: 'bg-emerald-400',
    romantic: 'bg-pink-400',
    intense: 'bg-red-500',
    confident: 'bg-purple-400',
    sleep: 'bg-slate-400',
  }
  return spots[selectedMood.value] || 'bg-transparent'
}

function moodSpotClass2(): string {
  const spots: Record<string, string> = {
    energetic: 'bg-red-600',
    happy: 'bg-yellow-500',
    chill: 'bg-cyan-400',
    calm: 'bg-indigo-500',
    sad: 'bg-violet-500',
    focus: 'bg-green-400',
    romantic: 'bg-rose-400',
    intense: 'bg-orange-600',
    confident: 'bg-fuchsia-400',
    sleep: 'bg-blue-500',
  }
  return spots[selectedMood.value] || 'bg-transparent'
}

function moodGlowClass(mood: string): string {
  const glows: Record<string, string> = {
    energetic: 'bg-orange-500',
    happy: 'bg-amber-400',
    chill: 'bg-teal-400',
    calm: 'bg-blue-400',
    sad: 'bg-indigo-400',
    focus: 'bg-emerald-400',
    romantic: 'bg-pink-400',
    intense: 'bg-red-500',
    confident: 'bg-purple-400',
    sleep: 'bg-slate-400',
  }
  return glows[mood] || 'bg-transparent'
}

function moodIconActiveBg(mood: string): string {
  const bg: Record<string, string> = {
    energetic: 'bg-orange-500/20 text-orange-400',
    happy: 'bg-amber-400/20 text-amber-300',
    chill: 'bg-teal-400/20 text-teal-300',
    calm: 'bg-blue-400/20 text-blue-300',
    sad: 'bg-indigo-400/20 text-indigo-300',
    focus: 'bg-emerald-400/20 text-emerald-300',
    romantic: 'bg-pink-400/20 text-pink-300',
    intense: 'bg-red-500/20 text-red-400',
    confident: 'bg-purple-400/20 text-purple-300',
    sleep: 'bg-slate-400/20 text-slate-300',
  }
  return bg[mood] || 'bg-spotify/20 text-spotify'
}

function moodSelectedBorder(mood: string): string {
  const borders: Record<string, string> = {
    energetic: 'border-orange-500/40 bg-orange-500/10',
    happy: 'border-amber-400/40 bg-amber-400/10',
    chill: 'border-teal-400/40 bg-teal-400/10',
    calm: 'border-blue-400/40 bg-blue-400/10',
    sad: 'border-indigo-400/40 bg-indigo-400/10',
    focus: 'border-emerald-400/40 bg-emerald-400/10',
    romantic: 'border-pink-400/40 bg-pink-400/10',
    intense: 'border-red-500/40 bg-red-500/10',
    confident: 'border-purple-400/40 bg-purple-400/10',
    sleep: 'border-slate-400/40 bg-slate-400/10',
  }
  return borders[mood] || 'border-spotify/40 bg-spotify/10'
}

function moodTextClass(mood: string): string {
  const txt: Record<string, string> = {
    energetic: 'text-orange-400',
    happy: 'text-amber-300',
    chill: 'text-teal-300',
    calm: 'text-blue-300',
    sad: 'text-indigo-300',
    focus: 'text-emerald-300',
    romantic: 'text-pink-300',
    intense: 'text-red-400',
    confident: 'text-purple-300',
    sleep: 'text-slate-300',
  }
  return txt[mood] || 'text-spotify'
}

function getEnergyBarColor(energy: number): string {
  if (energy > 0.7) return 'bg-green-400'
  if (energy > 0.4) return 'bg-yellow-400'
  return 'bg-blue-400'
}

function getMoodLabel(value: string) {
  return MOOD_OPTIONS.find((m) => m.value === value)?.label || value
}

async function toggleMood(mood: string) {
  if (selectedMood.value === mood) {
    selectedMood.value = ''
    tracks.value = []
    return
  }
  selectedMood.value = mood
  loading.value = true
  try {
    const res = await aiApi.generatePlaylist({
      prompt: '',
      mood,
      limit: 30,
    })
    if (res?.tracks) {
      tracks.value = res.tracks
    }
  } catch {
    toast.add({ severity: 'error', summary: 'Could not load tracks for this mood', life: 2000 })
    tracks.value = []
  } finally {
    loading.value = false
  }
}

function playTrack(index: number) {
  if (!tracks.value.length) return
  const queue = mapToPlaybackTracks(tracks.value)
  player.setQueueAndPlay(queue, index)
}

function formatTime(seconds?: number) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>
