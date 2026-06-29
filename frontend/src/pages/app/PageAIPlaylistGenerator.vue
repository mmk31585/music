<template>
  <div class="relative mx-auto min-h-screen w-full pb-36">
    <!-- ── Ambient purple/green aurora ── -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1" />
      <div class="aurora-spot-2" />
    </div>

    <div class="relative z-10 px-4 pt-8 md:px-6 lg:px-8">
      <div class="mx-auto max-w-7xl">
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
          <h1 class="text-4xl font-black text-white md:text-5xl">Playlist Generator</h1>
          <p class="mt-3 max-w-xl text-base text-white/50">
            Describe the vibe — a feeling, an activity, a scene — and let AI curate the perfect playlist.
          </p>
        </div>

        <!-- ════════════════════════════════════ -->
        <!-- MAIN CONTENT — Canvas Layout         -->
        <!-- ════════════════════════════════════ -->
        <div class="grid gap-8 lg:grid-cols-5">
          <!-- ──── LEFT: Controls ──── -->
          <div class="space-y-6 lg:col-span-2">
            <!-- Prompt Card -->
            <div class="rounded-2xl border border-white/6 bg-white/2 p-6 backdrop-blur-xs">
              <h2 class="mb-5 text-sm font-bold text-white">What are you in the mood for?</h2>

              <!-- Prompt textarea -->
              <div class="mb-5">
                <label class="mb-2 block text-xs font-medium text-white/40">Describe the vibe</label>
                <Textarea
                  v-model="prompt"
                  placeholder="e.g. Upbeat electronic music for a sunrise drive along the coast..."
                  :auto-resize="true"
                  rows="3"
                  class="w-full"
                />
              </div>

              <!-- Mood chips -->
              <div class="mb-5">
                <label class="mb-2 block text-xs font-medium text-white/40">Mood</label>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="mood in moodOptions"
                    :key="mood.value"
                    type="button"
                    class="group relative flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-xs font-medium transition-all duration-300"
                    :class="
                      selectedMood === mood.value
                        ? 'border-spotify/40 bg-spotify/15 text-white'
                        : 'border-white/6 bg-white/3 text-white/40 hover:border-white/12 hover:text-white'
                    "
                    @click="selectedMood = selectedMood === mood.value ? '' : mood.value"
                  >
                    <i aria-hidden="true" :class="mood.icon" class="text-[10px]" />
                    {{ mood.label }}
                  </button>
                </div>
              </div>

              <!-- Activity chips -->
              <div class="mb-5">
                <label class="mb-2 block text-xs font-medium text-white/40">Activity</label>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="activity in activityOptions"
                    :key="activity.value"
                    type="button"
                    class="group relative flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-xs font-medium transition-all duration-300"
                    :class="
                      selectedActivity === activity.value
                        ? 'border-aurora-purple/40 bg-aurora-purple/15 text-white'
                        : 'border-white/6 bg-white/3 text-white/40 hover:border-white/12 hover:text-white'
                    "
                    @click="selectedActivity = selectedActivity === activity.value ? '' : activity.value"
                  >
                    <i aria-hidden="true" :class="activity.icon" class="text-[10px]" />
                    {{ activity.label }}
                  </button>
                </div>
              </div>

              <!-- Genre + Track limit -->
              <div class="mb-5 grid grid-cols-2 gap-4">
                <div>
                  <label class="mb-2 block text-xs font-medium text-white/40">Genre</label>
                  <InputText
                    v-model="genre"
                    placeholder="Any genre"
                    class="w-full"
                  />
                </div>
                <div>
                  <label class="mb-2 block text-xs font-medium text-white/40">Track limit</label>
                  <Select v-model="limit" :options="[10, 20, 30, 50]" class="w-full" />
                </div>
              </div>

              <!-- Generate button -->
              <Button
                label="Generate Playlist"
                icon="pi pi-sparkles"
                severity="success"
                class="mt-2 w-full"
                :loading="generating"
                :disabled="generating || (!prompt && !selectedMood && !selectedActivity)"
                @click="handleGenerate"
              />
            </div>

            <!-- Generation History -->
            <div
              v-if="history.length"
              class="rounded-2xl border border-white/6 bg-white/2 p-6 backdrop-blur-xs"
            >
              <h3 class="mb-3 flex items-center gap-2 text-[10px] font-bold tracking-[0.2em] text-white/40 uppercase">
                <i aria-hidden="true" class="pi pi-history" />
                Recent Generations
              </h3>
              <div class="space-y-2">
                <button
                  v-for="(item, i) in history"
                  :key="i"
                  type="button"
                  class="group flex w-full items-center gap-3 rounded-xl bg-white/3 px-4 py-3 text-left transition hover:bg-white/6"
                  @click="restoreHistory(item)"
                >
                  <div
                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-linear-to-br from-spotify/20 to-aurora-purple/20"
                  >
                    <i aria-hidden="true" class="pi pi-sparkles text-[10px] text-spotify" />
                  </div>
                  <div class="min-w-0 flex-1">
                    <div class="truncate text-sm font-medium text-white group-hover:text-spotify">
                      {{ item.name }}
                    </div>
                    <div class="text-[10px] text-white/30">
                      {{ item.tracks }} tracks · {{ item.date }}
                    </div>
                  </div>
                </button>
              </div>
            </div>
          </div>

          <!-- ──── RIGHT: Results ──── -->
          <div class="lg:col-span-3">
            <!-- Generating State -->
            <div
              v-if="generating"
              class="flex flex-col items-center justify-center rounded-2xl border border-white/6 bg-linear-to-br from-spotify/5 via-aurora-purple/5 to-black/40 px-6 py-24 text-center backdrop-blur-xs"
            >
              <div class="relative mb-6">
                <div
                  class="flex h-20 w-20 items-center justify-center rounded-full bg-linear-to-br from-spotify/20 to-aurora-purple/20"
                >
                  <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-spotify" />
                </div>
                <div
                  class="absolute -inset-2 animate-ping rounded-full border border-spotify/20"
                />
              </div>
              <p class="text-lg font-bold text-white">Generating your perfect playlist...</p>
              <p class="mt-2 text-sm text-white/40">
                AI is analyzing your preferences and curating tracks
              </p>
              <div class="mt-6 flex items-center gap-2 text-xs text-white/20">
                <span class="flex h-1.5 w-1.5 rounded-full bg-spotify/50" />
                <span>Analyzing taste profile</span>
                <span class="flex h-1.5 w-1.5 rounded-full bg-spotify/30" />
                <span>Matching mood</span>
                <span class="flex h-1.5 w-1.5 rounded-full bg-white/10" />
                <span>Curating tracks</span>
              </div>
            </div>

            <!-- Result -->
            <div v-else-if="result" class="space-y-5">
              <!-- Playlist Header -->
              <div
                class="relative overflow-hidden rounded-2xl border border-white/6 bg-linear-to-br from-spotify/10 via-aurora-purple/5 to-black/40 p-8 backdrop-blur-2xl md:p-10"
              >
                <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-2xl">
                  <div
                    class="absolute -top-1/3 -left-1/4 h-64 w-64 rounded-full bg-spotify/10 blur-[100px]"
                  />
                  <div
                    class="absolute -bottom-1/3 -right-1/4 h-64 w-64 rounded-full bg-aurora-purple/10 blur-[100px]"
                  />
                </div>

                <div class="relative flex flex-col gap-6 md:flex-row md:items-start">
                  <!-- Playlist art (first track cover) -->
                  <div
                    class="h-40 w-40 shrink-0 overflow-hidden rounded-2xl bg-white/5 shadow-2xl ring-1 ring-white/6"
                  >
                    <img
                      v-if="result.tracks[0]?.cover_url"
                      :src="result.tracks[0].cover_url"
                      :alt="result.name"
                      class="h-full w-full object-cover"
                      @error="onImgError"
                    />
                    <div v-else class="flex h-full items-center justify-center bg-linear-to-br from-spotify/20 to-aurora-purple/20">
                      <i aria-hidden="true" class="pi pi-sparkles text-4xl text-white/40" />
                    </div>
                  </div>

                  <div class="flex-1">
                    <div class="mb-2 flex items-center gap-2">
                      <span
                        class="inline-flex items-center gap-1 rounded-full border border-spotify/20 bg-spotify/10 px-2.5 py-0.5 text-[9px] font-bold tracking-[0.15em] text-spotify uppercase"
                      >
                        <i aria-hidden="true" class="pi pi-sparkles text-[8px]" />
                        AI Generated
                      </span>
                    </div>
                    <h2 class="text-2xl font-black text-white md:text-3xl">{{ result.name }}</h2>
                    <p class="mt-2 text-sm text-white/50">{{ result.description }}</p>
                    <div class="mt-4 flex flex-wrap items-center gap-4">
                      <div class="flex items-center gap-1.5 text-xs text-white/40">
                        <i aria-hidden="true" class="pi pi-music" />
                        {{ result.tracks.length }} tracks
                      </div>
                      <div class="flex items-center gap-1.5 text-xs text-white/40">
                        <i aria-hidden="true" class="pi pi-calendar" />
                        {{ new Date(result.generated_at).toLocaleDateString() }}
                      </div>
                      <div class="flex items-center gap-1.5 text-xs text-white/40">
                        <i aria-hidden="true" class="pi pi-clock" />
                        {{ totalDuration }}
                      </div>
                    </div>
                    <div class="mt-5 flex gap-3">
                      <Button
                        label="Play All"
                        icon="pi pi-play"
                        severity="success"
                        size="small"
                        :disabled="!result.tracks.length"
                        @click="playAll"
                      />
                      <Button
                        label="New Playlist"
                        icon="pi pi-refresh"
                        severity="secondary"
                        size="small"
                        text
                        @click="reset"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <!-- Track List -->
              <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2 backdrop-blur-xs">
                <template v-if="result.tracks.length">
                  <div
                    v-for="(track, index) in result.tracks"
                    :key="track.id"
                    role="button"
                    tabindex="0"
                    class="group flex items-center gap-3 px-4 py-2.5 transition hover:bg-white/4"
                    :style="{ animationDelay: `${index * 40}ms` }"
                    @click="playTrack(index)"
                    @keydown.enter="playTrack(index)"
                  >
                    <span class="flex w-6 items-center justify-center">
                      <span class="text-xs font-bold text-white/20 group-hover:hidden">{{ index + 1 }}</span>
                      <i aria-hidden="true" class="pi pi-play-fill hidden text-xs text-white group-hover:block" />
                    </span>

                    <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/5 ring-1 ring-white/6">
                      <img
                        v-if="track.cover_url"
                        :src="track.cover_url"
                        :alt="track.title"
                        loading="lazy"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div v-else class="flex h-full items-center justify-center">
                        <i aria-hidden="true" class="pi pi-music text-xs text-white/20" />
                      </div>
                      <button
                        type="button"
                        aria-label="Play track"
                        class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
                      >
                        <i aria-hidden="true" class="pi pi-play-fill text-xs text-white" />
                      </button>
                    </div>

                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                      <p v-if="track.artist" class="truncate text-xs text-white/40">
                        {{ track.artist }}
                      </p>
                    </div>

                    <span class="shrink-0 text-xs text-white/30 tabular-nums">
                      {{ formatTime(track.duration) }}
                    </span>
                  </div>
                </template>
                <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
                  <i aria-hidden="true" class="pi pi-inbox text-4xl text-white/20" />
                  <p class="text-sm text-white/40">No tracks in generated playlist</p>
                </div>
              </div>
            </div>

            <!-- Empty / Ready State -->
            <div
              v-else
              class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-white/6 bg-linear-to-br from-white/1 to-white/2 px-6 py-24 text-center backdrop-blur-xs"
            >
              <div
                class="mb-6 flex h-20 w-20 items-center justify-center rounded-2xl bg-linear-to-br from-spotify/20 to-aurora-purple/20"
              >
                <i aria-hidden="true" class="pi pi-sparkles text-3xl text-spotify" />
              </div>
              <h3 class="text-xl font-bold text-white">Ready when you are</h3>
              <p class="mt-2 max-w-md text-sm leading-relaxed text-white/40">
                Describe the vibe on the left — a mood, an activity, a genre.
                AI will craft a custom playlist that matches your moment.
              </p>
              <div class="mt-8 flex flex-wrap justify-center gap-3">
                <div
                  v-for="suggestion in quickSuggestions"
                  :key="suggestion"
                  class="cursor-pointer rounded-full border border-white/6 bg-white/3 px-4 py-2 text-xs text-white/40 transition hover:border-white/12 hover:text-white"
                  @click="prompt = suggestion"
                >
                  {{ suggestion }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAIApi } from '@/services/api/ai'
import { MOOD_OPTIONS, ACTIVITY_OPTIONS } from '@/services/api/ai/types'
import type { AIPlaylistResponse } from '@/services/api/ai/types'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'
import { useToast } from 'primevue/usetoast'

const aiApi = useAIApi()
const player = usePlayer()
const playerApi = usePlayerApi()
const toast = useToast()

const prompt = ref('')
const selectedMood = ref('')
const selectedActivity = ref('')
const genre = ref('')
const limit = ref(20)
const generating = ref(false)
const result = ref<AIPlaylistResponse | null>(null)
const moodOptions = [...MOOD_OPTIONS]
const activityOptions = [...ACTIVITY_OPTIONS]

const history = ref<{ name: string; tracks: number; date: string; data: AIPlaylistResponse }[]>([])

const quickSuggestions = [
  'Late night jazz for a rainy city',
  'Upbeat workout mix with heavy bass',
  'Chill lo-fi for deep focus session',
  'Sunset beach vibes with acoustic',
  'Dark electronic for a cyberpunk night',
]

const totalDuration = computed(() => {
  if (!result.value?.tracks.length) return '0 min'
  const total = result.value.tracks.reduce((acc, t) => acc + (t.duration || 0), 0)
  const mins = Math.floor(total / 60)
  return `${mins} min`
})

function formatTime(seconds?: number) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

async function handleGenerate() {
  generating.value = true
  try {
    const response = await aiApi.generatePlaylist({
      prompt: prompt.value,
      mood: selectedMood.value || undefined,
      activity: selectedActivity.value || undefined,
      genre: genre.value || undefined,
      limit: limit.value,
    })
    if (response) {
      result.value = response
      history.value.unshift({
        name: response.name,
        tracks: response.tracks.length,
        date: new Date().toLocaleDateString(),
        data: { ...response },
      })
      if (history.value.length > 10) history.value.pop()
    }
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Generation failed',
      detail: 'Could not generate playlist. Try again.',
      life: 3000,
    })
  } finally {
    generating.value = false
  }
}

function playTrack(index: number) {
  if (!result.value) return
  const queue = result.value.tracks.map((t) => ({
    id: t.id,
    title: t.title,
    artistName: t.artist || 'Unknown',
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration ?? null,
    streamUrl: playerApi.getTrackStreamUrl(t.id),
  }))
  player.setQueueAndPlay(queue, index)
}

function playAll() {
  if (result.value?.tracks.length) playTrack(0)
}

function restoreHistory(item: { data: AIPlaylistResponse }) {
  result.value = item.data
}

function reset() {
  result.value = null
  prompt.value = ''
  selectedMood.value = ''
  selectedActivity.value = ''
  genre.value = ''
  limit.value = 20
}
</script>
