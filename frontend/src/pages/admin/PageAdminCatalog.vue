<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Catalog management"
      description="Review public catalog entities while preparing create/edit/delete actions."
    />

    <section class="grid gap-6 xl:grid-cols-2">
      <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xl font-bold text-white">Tracks</h2>
          <span class="text-sm text-slate-400">{{ tracks.length }}</span>
        </div>

        <div v-if="loading" class="space-y-3">
          <div v-for="i in 6" :key="i" class="h-16 animate-pulse rounded-2xl bg-white/5" />
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="track in tracks"
            :key="track.id"
            class="rounded-2xl border border-white/10 bg-black/20 px-4 py-4"
          >
            <p class="font-medium text-white">{{ track.title }}</p>
            <p class="mt-1 text-sm text-slate-400">
              {{ track.artist_name || 'Unknown artist' }}
              <span v-if="track.album_title">• {{ track.album_title }}</span>
            </p>
          </div>
        </div>
      </div>

      <div class="space-y-6">
        <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-xl font-bold text-white">Artists</h2>
            <span class="text-sm text-slate-400">{{ artists.length }}</span>
          </div>

          <div class="space-y-3">
            <div
              v-for="artist in artists"
              :key="artist.id"
              class="rounded-2xl border border-white/10 bg-black/20 px-4 py-4"
            >
              <p class="font-medium text-white">{{ artist.name }}</p>
            </div>
          </div>
        </div>

        <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-xl font-bold text-white">Albums</h2>
            <span class="text-sm text-slate-400">{{ albums.length }}</span>
          </div>

          <div class="space-y-3">
            <div
              v-for="album in albums"
              :key="album.id"
              class="rounded-2xl border border-white/10 bg-black/20 px-4 py-4"
            >
              <p class="font-medium text-white">{{ album.title }}</p>
              <p class="mt-1 text-sm text-slate-400">
                {{ album.artist_name || 'Unknown artist' }}
              </p>
            </div>
          </div>
        </div>

        <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-xl font-bold text-white">Genres</h2>
            <span class="text-sm text-slate-400">{{ genres.length }}</span>
          </div>

          <div class="flex flex-wrap gap-3">
            <span
              v-for="genre in genres"
              :key="genre.id"
              class="rounded-full border border-white/10 bg-black/20 px-4 py-2 text-sm text-slate-300"
            >
              {{ genre.name }}
            </span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useCatalogApi } from '@/services/api/catalog'
import type { Track, Artist, Album, Genre } from '@/services/api/catalog'
import { AdminSectionHeader } from '@/components/admin'

const api = useCatalogApi()

const loading = ref(false)
const tracks = ref<Track[]>([])
const artists = ref<Artist[]>([])
const albums = ref<Album[]>([])
const genres = ref<Genre[]>([])

async function fetchCatalog() {
  loading.value = true
  try {
    const [tracksRes, artistsRes, albumsRes, genresRes] = await Promise.all([
      api.getTracks(),
      api.getArtists(),
      api.getAlbums(),
      api.getGenres(),
    ])

    tracks.value = tracksRes
    artists.value = artistsRes
    albums.value = albumsRes
    genres.value = genresRes
  } finally {
    loading.value = false
  }
}

onMounted(fetchCatalog)
</script>
