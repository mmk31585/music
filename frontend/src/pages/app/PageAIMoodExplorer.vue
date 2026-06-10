<template>
  <div class="mx-auto w-full max-w-5xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div class="mb-8">
      <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">AI-Powered</p>
      <h1 class="mt-2 text-3xl font-black text-white md:text-4xl">Mood Explorer</h1>
      <p class="mt-2 text-sm text-slate-400">
        Browse tracks by their mood profile — energy, valence, and more.
      </p>
    </div>

    <div class="mb-8">
      <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-5">
        <button
          v-for="mood in moodOptions"
          :key="mood.value"
          type="button"
          class="flex flex-col items-center gap-2 rounded-2xl border px-4 py-5 text-center transition"
          :class="
            selectedMood === mood.value
              ? 'border-[#1db954]/50 bg-[#1db954]/10 text-[#1db954]'
              : 'border-white/[0.06] bg-white/[0.03] text-slate-300 hover:bg-white/[0.06]'
          "
          @click="toggleMood(mood.value)"
        >
          <div
            class="flex h-12 w-12 items-center justify-center rounded-xl text-lg"
            :class="selectedMood === mood.value ? 'bg-[#1db954]/20' : 'bg-white/[0.06]'"
          >
            <i :class="mood.icon" />
          </div>
          <span class="text-xs font-bold">{{ mood.label }}</span>
        </button>
      </div>
    </div>

    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-sm font-bold text-white">
        {{ selectedMood ? `${getMoodLabel(selectedMood)} Tracks` : 'Select a mood to browse' }}
      </h2>
      <span v-if="filteredTracks.length" class="text-xs text-slate-500"
        >{{ filteredTracks.length }} tracks</span
      >
    </div>

    <div v-if="loading" class="space-y-2">
      <SkeletonLoader v-for="i in 6" :key="i" variant="track" />
    </div>

    <div
      v-else-if="filteredTracks.length"
      class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]"
    >
      <div
        v-for="(track, index) in filteredTracks"
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
          <p v-if="track.artist" class="truncate text-xs text-slate-400">{{ track.artist }}</p>
        </div>

        <div v-if="track.mood" class="hidden items-center gap-1 sm:flex">
          <div class="flex h-1.5 w-12 overflow-hidden rounded-full bg-white/10">
            <div
              class="h-full rounded-full transition-all"
              :class="getEnergyColor(track.mood.energy)"
              :style="{ width: `${track.mood.energy * 100}%` }"
            />
          </div>
          <span class="text-[10px] text-slate-500">{{ Math.round(track.mood.energy * 100) }}%</span>
        </div>

        <span class="text-xs text-slate-500">{{ formatTime(track.duration) }}</span>
      </div>
    </div>

    <div
      v-else-if="!selectedMood"
      class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.08] bg-white/[0.02] px-6 py-20 text-center"
    >
      <div
        class="flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-[#1db954]/20 to-blue-500/20"
      >
        <i class="pi pi-heart text-xl text-[#1db954]" />
      </div>
      <h3 class="text-lg font-bold text-white">Pick a mood</h3>
      <p class="max-w-xs text-sm text-slate-400">
        Choose a mood above to discover tracks that match the vibe.
      </p>
    </div>

    <div
      v-else
      class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.08] bg-white/[0.02] px-6 py-20 text-center"
    >
      <div
        class="flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-[#1db954]/20 to-amber-500/20"
      >
        <i class="pi pi-inbox text-xl text-slate-400" />
      </div>
      <h3 class="text-lg font-bold text-white">No tracks found</h3>
      <p class="max-w-xs text-sm text-slate-400">
        No tracks match this mood. Try a different mood.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { SkeletonLoader } from '@/components/common'
import { useAIApi } from '@/services/api/ai'
import { MOOD_OPTIONS } from '@/services/api/ai/types'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'
import { useToast } from 'primevue/usetoast'

const aiApi = useAIApi()
const player = usePlayer()
const playerApi = usePlayerApi()
const toast = useToast()

const selectedMood = ref('')
const loading = ref(false)
const tracks = ref<any[]>([])

const moodOptions = MOOD_OPTIONS

const filteredTracks = computed(() => tracks.value)

function getMoodLabel(value: string) {
  return MOOD_OPTIONS.find((m) => m.value === value)?.label || value
}

function getEnergyColor(energy: number) {
  if (energy > 0.7) return 'bg-green-400'
  if (energy > 0.4) return 'bg-yellow-400'
  return 'bg-blue-400'
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
  const queue = tracks.value.map((t) => ({
    id: t.id,
    title: t.title,
    artistName: t.artist || 'Unknown',
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration ?? null,
    streamUrl: playerApi.getTrackStreamUrl(t.id),
  }))
  player.setQueueAndPlay(queue, index)
}

function formatTime(seconds?: number) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>
