<template>
  <div class="mx-auto w-full max-w-5xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div class="mb-8">
      <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">AI-Powered</p>
      <h1 class="mt-2 text-3xl font-black text-white md:text-4xl">Playlist Generator</h1>
      <p class="mt-2 text-sm text-slate-400">
        Describe the vibe you want — mood, activity, genre — and let AI curate the perfect playlist.
      </p>
    </div>

    <div class="grid gap-8 lg:grid-cols-5">
      <div class="space-y-6 lg:col-span-2">
        <div class="rounded-2xl border border-white/[0.06] bg-white/[0.03] p-6">
          <h2 class="mb-5 text-sm font-bold text-white">What are you in the mood for?</h2>

          <div class="space-y-4">
            <div>
              <label class="mb-2 block text-xs font-medium text-slate-400">Prompt</label>
              <Textarea
                v-model="prompt"
                placeholder="e.g. Upbeat electronic music for a morning run"
                :auto-resize="true"
                rows="3"
                class="w-full"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="mb-2 block text-xs font-medium text-slate-400">Mood</label>
                <Select
                  v-model="selectedMood"
                  :options="moodOptions"
                  option-label="label"
                  option-value="value"
                  placeholder="Any mood"
                  class="w-full"
                  show-clear
                >
                  <template #value="slotProps">
                    <div v-if="slotProps.value" class="flex items-center gap-2">
                      <i :class="getMoodIcon(slotProps.value)" class="text-xs" />
                      <span>{{ getMoodLabel(slotProps.value) }}</span>
                    </div>
                    <span v-else class="text-slate-400">Any mood</span>
                  </template>
                  <template #option="slotProps">
                    <div class="flex items-center gap-2">
                      <i :class="slotProps.option.icon" class="text-xs" />
                      <span>{{ slotProps.option.label }}</span>
                    </div>
                  </template>
                </Select>
              </div>

              <div>
                <label class="mb-2 block text-xs font-medium text-slate-400">Activity</label>
                <Select
                  v-model="selectedActivity"
                  :options="activityOptions"
                  option-label="label"
                  option-value="value"
                  placeholder="Any activity"
                  class="w-full"
                  show-clear
                >
                  <template #value="slotProps">
                    <div v-if="slotProps.value" class="flex items-center gap-2">
                      <i :class="getActivityIcon(slotProps.value)" class="text-xs" />
                      <span>{{ getActivityLabel(slotProps.value) }}</span>
                    </div>
                    <span v-else class="text-slate-400">Any activity</span>
                  </template>
                  <template #option="slotProps">
                    <div class="flex items-center gap-2">
                      <i :class="slotProps.option.icon" class="text-xs" />
                      <span>{{ slotProps.option.label }}</span>
                    </div>
                  </template>
                </Select>
              </div>
            </div>

            <div>
              <label class="mb-2 block text-xs font-medium text-slate-400">Genre (optional)</label>
              <InputText v-model="genre" placeholder="e.g. rock, jazz, electronic" class="w-full" />
            </div>

            <div>
              <label class="mb-2 block text-xs font-medium text-slate-400">Track limit</label>
              <Select v-model="limit" :options="[10, 20, 30, 50]" class="w-full" />
            </div>

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
        </div>

        <div
          v-if="history.length"
          class="rounded-2xl border border-white/[0.06] bg-white/[0.03] p-6"
        >
          <h3 class="text-xs font-bold tracking-wider text-slate-400 uppercase">
            Recent Generations
          </h3>
          <div class="mt-3 space-y-2">
            <button
              v-for="(item, i) in history"
              :key="i"
              type="button"
              class="w-full rounded-xl bg-white/[0.04] px-4 py-3 text-left text-xs text-slate-300 transition hover:bg-white/[0.08]"
              @click="restoreHistory(item)"
            >
              <div class="font-medium text-white">{{ item.name }}</div>
              <div class="mt-0.5 text-slate-500">{{ item.tracks }} tracks • {{ item.date }}</div>
            </button>
          </div>
        </div>
      </div>

      <div class="lg:col-span-3">
        <div
          v-if="generating"
          class="flex flex-col items-center justify-center rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-24 text-center"
        >
          <i class="pi pi-spin pi-spinner text-3xl text-[#1db954]" />
          <p class="mt-4 text-sm font-medium text-white">Generating your perfect playlist...</p>
          <p class="mt-1 text-xs text-slate-400">AI is curating tracks based on your preferences</p>
        </div>

        <div v-else-if="result" class="space-y-4">
          <div
            class="flex items-center justify-between rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-4"
          >
            <div>
              <h2 class="text-lg font-bold text-white">{{ result.name }}</h2>
              <p class="text-xs text-slate-400">{{ result.description }}</p>
              <p class="mt-1 text-xs text-slate-500">
                {{ result.tracks.length }} tracks • Generated
                {{ new Date(result.generated_at).toLocaleDateString() }}
              </p>
            </div>
            <div class="flex gap-2">
              <Button
                label="Play All"
                icon="pi pi-play"
                severity="success"
                size="small"
                :disabled="!result.tracks.length"
                @click="playAll"
              />
              <Button
                label="New"
                icon="pi pi-refresh"
                severity="secondary"
                size="small"
                text
                @click="reset"
              />
            </div>
          </div>

          <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
            <template v-if="result.tracks.length">
              <div
                v-for="(track, index) in result.tracks"
                :key="track.id"
                class="group flex items-center gap-3 px-4 py-2 transition hover:bg-white/[0.06]"
              >
                <span class="w-6 text-right text-xs text-slate-500">{{ index + 1 }}</span>

                <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                  <img
                    v-if="track.cover_url"
                    :src="track.cover_url"
                    :alt="track.title"
                    loading="lazy"
                    class="h-full w-full object-cover"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i class="pi pi-music text-xs text-slate-500" />
                  </div>
                  <button
                    type="button"
                    class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                    @click="playTrack(index)"
                  >
                    <i class="pi pi-play-fill text-xs text-white" />
                  </button>
                </div>

                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                  <p v-if="track.artist" class="truncate text-xs text-slate-400">
                    {{ track.artist }}
                  </p>
                </div>

                <span class="text-xs text-slate-500">{{ formatTime(track.duration) }}</span>
              </div>
            </template>
            <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
              <i class="pi pi-inbox text-4xl text-slate-500" />
              <p class="text-sm text-slate-400">No tracks in generated playlist</p>
            </div>
          </div>
        </div>

        <div
          v-else
          class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-white/[0.08] bg-white/[0.02] px-6 py-24 text-center"
        >
          <div
            class="flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-[#1db954]/20 to-purple-500/20"
          >
            <i class="pi pi-sparkles text-2xl text-[#1db954]" />
          </div>
          <h3 class="mt-6 text-lg font-bold text-white">Ready when you are</h3>
          <p class="mt-2 max-w-xs text-sm text-slate-400">
            Describe the vibe, pick a mood or activity, and let AI create a custom playlist for you.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
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
const moodOptions = ref([...MOOD_OPTIONS])
const activityOptions = ref([...ACTIVITY_OPTIONS])

const history = ref<{ name: string; tracks: number; date: string; data: AIPlaylistResponse }[]>([])

function getMoodIcon(value: string) {
  return MOOD_OPTIONS.find((m) => m.value === value)?.icon || ''
}
function getMoodLabel(value: string) {
  return MOOD_OPTIONS.find((m) => m.value === value)?.label || value
}
function getActivityIcon(value: string) {
  return ACTIVITY_OPTIONS.find((a) => a.value === value)?.icon || ''
}
function getActivityLabel(value: string) {
  return ACTIVITY_OPTIONS.find((a) => a.value === value)?.label || value
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

function formatTime(seconds?: number) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>
