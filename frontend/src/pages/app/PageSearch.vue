<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section class="mb-8">
      <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">Discover</p>
      <h1 class="mt-2 text-4xl font-black text-white md:text-5xl">Search</h1>
      <p class="mt-2 text-sm text-slate-400">Find tracks, artists, and albums.</p>
    </section>

    <div class="relative mb-8">
      <i class="pi pi-search absolute left-4 top-1/2 -translate-y-1/2 text-sm text-slate-500" />
      <input
        v-model="query"
        type="text"
        placeholder="Search tracks, artists, albums..."
        class="w-full rounded-2xl border border-white/[0.08] bg-white/[0.04] py-4 pl-11 pr-4 text-sm text-white outline-none transition placeholder:text-slate-600 focus:border-[#1db954]/40 focus:bg-white/[0.06]"
        @input="onInput"
        @keydown.enter="doSearch"
      />
      <button
        v-if="query"
        type="button"
        class="absolute right-4 top-1/2 -translate-y-1/2 flex h-6 w-6 items-center justify-center rounded-full bg-white/10 text-xs text-slate-400 hover:bg-white/20"
        @click="clearSearch"
      >
        <i class="pi pi-times" />
      </button>
    </div>

    <div v-if="searching" class="space-y-3">
      <div v-for="i in 6" :key="i" class="flex items-center gap-4 rounded-2xl bg-white/[0.02] p-3">
        <div class="h-14 w-14 animate-pulse rounded-xl bg-white/[0.06]" />
        <div class="flex-1 space-y-2">
          <div class="h-4 w-40 animate-pulse rounded bg-white/[0.06]" />
          <div class="h-3 w-24 animate-pulse rounded bg-white/[0.04]" />
        </div>
      </div>
    </div>

    <template v-else-if="hasSearched">
      <div v-if="hasNoResults" class="flex flex-col items-center gap-4 py-16 text-center">
        <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
          <i class="pi pi-search text-3xl text-slate-600" />
        </div>
        <p class="text-sm text-slate-400">No results for "<span class="font-medium text-white">{{ lastQuery }}</span>"</p>
        <p class="text-xs text-slate-500">Try a different search term</p>
      </div>

      <div v-else class="space-y-10">
        <section v-if="results.tracks.length">
          <div class="mb-4 flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10">
              <i class="pi pi-play-circle text-xs text-emerald-400" />
            </div>
            <h2 class="text-lg font-bold text-white">Tracks</h2>
            <span class="text-xs text-slate-500">{{ results.tracks.length }}</span>
          </div>
          <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
            <div
              v-for="(track, i) in results.tracks"
              :key="track.id"
              class="group flex cursor-pointer items-center gap-3 px-4 py-3 transition hover:bg-white/[0.03]"
              @click="playTrack(track)"
            >
              <span class="w-6 text-center text-xs text-slate-600">{{ i + 1 }}</span>
              <div class="h-11 w-11 shrink-0 overflow-hidden rounded-lg bg-white/[0.04]">
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i class="pi pi-music text-xs text-slate-600" />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p class="truncate text-xs text-slate-500">{{ track.artist_name || 'Unknown artist' }}</p>
              </div>
              <span class="text-xs text-slate-600 tabular-nums">{{ formatDuration(track.duration_seconds) }}</span>
            </div>
          </div>
        </section>

        <section v-if="results.artists.length">
          <div class="mb-4 flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10">
              <i class="pi pi-users text-xs text-blue-400" />
            </div>
            <h2 class="text-lg font-bold text-white">Artists</h2>
            <span class="text-xs text-slate-500">{{ results.artists.length }}</span>
          </div>
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
            <RouterLink
              v-for="artist in results.artists"
              :key="artist.id"
              :to="`/artist/${artist.id}`"
              class="group rounded-2xl border border-white/[0.06] bg-white/[0.02] p-4 text-center transition hover:bg-white/[0.05]"
            >
              <div class="mx-auto h-20 w-20 overflow-hidden rounded-full bg-white/[0.06]">
                <img
                  v-if="artist.image_url"
                  :src="artist.image_url"
                  :alt="artist.name"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i class="pi pi-user text-xl text-slate-500" />
                </div>
              </div>
              <p class="mt-3 truncate text-sm font-medium text-white">{{ artist.name }}</p>
              <p class="text-xs text-slate-500">Artist</p>
            </RouterLink>
          </div>
        </section>

        <section v-if="results.albums.length">
          <div class="mb-4 flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-500/10">
              <i class="pi pi-book text-xs text-purple-400" />
            </div>
            <h2 class="text-lg font-bold text-white">Albums</h2>
            <span class="text-xs text-slate-500">{{ results.albums.length }}</span>
          </div>
          <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <RouterLink
              v-for="album in results.albums"
              :key="album.id"
              :to="`/album/${album.id}`"
              class="group"
            >
              <div class="relative aspect-square overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50">
                <img
                  v-if="album.cover_url"
                  :src="album.cover_url"
                  :alt="album.title"
                  class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i class="pi pi-compact-disc text-3xl text-slate-500" />
                </div>
              </div>
              <p class="mt-2 truncate text-sm font-medium text-white">{{ album.title }}</p>
              <p class="truncate text-xs text-slate-400">{{ album.artist_name || '—' }}</p>
            </RouterLink>
          </div>
        </section>
      </div>
    </template>

    <div v-else class="flex flex-col items-center gap-4 py-16 text-center">
      <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
        <i class="pi pi-search text-3xl text-slate-600" />
      </div>
      <p class="text-sm text-slate-400">Type to search tracks, artists, and albums</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSearchApi } from '@/services/api/catalog/search'
import type { SearchResult } from '@/services/api/catalog/search'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'

const searchApi = useSearchApi()
const player = usePlayer()
const playerApi = usePlayerApi()

const query = ref('')
const lastQuery = ref('')
const results = ref<SearchResult>({ tracks: [], artists: [], albums: [] })
const searching = ref(false)
const hasSearched = ref(false)
let debounce: ReturnType<typeof setTimeout> | null = null

const hasNoResults = ref(false)

function onInput() {
  if (debounce) clearTimeout(debounce)
  debounce = setTimeout(() => {
    doSearch()
  }, 350)
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return

  lastQuery.value = q
  searching.value = true
  hasSearched.value = true
  hasNoResults.value = false

  try {
    const res = await searchApi.searchCatalog({ query: q, limit: 20 })
    results.value = res
    hasNoResults.value = !res.tracks.length && !res.artists.length && !res.albums.length
  } catch {
    hasNoResults.value = true
    results.value = { tracks: [], artists: [], albums: [] }
  } finally {
    searching.value = false
  }
}

function clearSearch() {
  query.value = ''
  results.value = { tracks: [], artists: [], albums: [] }
  hasSearched.value = false
  hasNoResults.value = false
}

function playTrack(track: any) {
  void player.setQueueAndPlay(
    [{
      id: String(track.id),
      title: track.title,
      artistName: track.artist_name || 'Unknown',
      albumTitle: track.album_title || null,
      coverUrl: track.cover_url || null,
      durationSeconds: track.duration_seconds ?? null,
      streamUrl: playerApi.getTrackStreamUrl(String(track.id)),
    }],
    0,
  )
}

function formatDuration(value?: number | null): string {
  if (!value) return '—'
  const mins = Math.floor(value / 60)
  const secs = value % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}
</script>