<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Overview"
      title="Admin dashboard"
      description="A central place to manage uploaded media and catalog content."
    >
      <template #actions>
        <RouterLink
          to="/admin/media"
          class="rounded-full bg-[#1db954] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90"
        >
          Upload media
        </RouterLink>
      </template>
    </AdminSectionHeader>

    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard
        label="Tracks"
        :value="trackCount"
        hint="Loaded from catalog API"
        icon="pi pi-play-circle"
      />
      <AdminStatCard
        label="Artists"
        :value="artistCount"
        hint="Loaded from catalog API"
        icon="pi pi-users"
      />
      <AdminStatCard
        label="Albums"
        :value="albumCount"
        hint="Loaded from catalog API"
        icon="pi pi-book"
      />
      <AdminStatCard
        label="Genres"
        :value="genreCount"
        hint="Loaded from catalog API"
        icon="pi pi-tags"
      />
    </section>

    <section class="mt-8">
      <CatalogQuickActions />
    </section>

    <section class="mt-8 grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
        <h2 class="text-xl font-bold text-white">Recent tracks</h2>
        <p class="mt-2 text-sm text-slate-400">
          Quick glance at tracks available in the public catalog.
        </p>

        <div class="mt-6 space-y-3">
          <div
            v-if="loading"
            v-for="i in 5"
            :key="i"
            class="h-16 animate-pulse rounded-2xl bg-white/5"
          />

          <div
            v-else
            v-for="track in recentTracks"
            :key="track.id"
            class="flex items-center justify-between rounded-2xl border border-white/10 bg-black/20 px-4 py-4"
          >
            <div class="min-w-0">
              <p class="truncate font-medium text-white">{{ track.title }}</p>
              <p class="truncate text-sm text-slate-400">
                {{ track.artist_name || 'Unknown artist' }}
              </p>
            </div>

            <div class="text-sm text-slate-500">
              {{ formatDuration(track.duration_seconds) }}
            </div>
          </div>
        </div>
      </div>

      <div class="rounded-3xl border border-white/10 bg-white/5 p-6">
        <h2 class="text-xl font-bold text-white">Admin notes</h2>
        <ul class="mt-4 space-y-3 text-sm text-slate-400">
          <li>• Use the media upload page to upload audio files.</li>
          <li>• Connect catalog create/update endpoints when backend is ready.</li>
          <li>• Add table pages for tracks, artists, albums, and genres later.</li>
          <li>• Protect all admin routes with role-based guards.</li>
        </ul>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useCatalogApi } from '@/services/api/catalog'
import type { Track, Artist, Album, Genre } from '@/services/api/catalog'
import { AdminSectionHeader, AdminStatCard, CatalogQuickActions } from '@/components/admin'

const api = useCatalogApi()

const loading = ref(false)
const tracks = ref<Track[]>([])
const artists = ref<Artist[]>([])
const albums = ref<Album[]>([])
const genres = ref<Genre[]>([])

const trackCount = computed(() => tracks.value.length)
const artistCount = computed(() => artists.value.length)
const albumCount = computed(() => albums.value.length)
const genreCount = computed(() => genres.value.length)
const recentTracks = computed(() => tracks.value.slice(0, 5))

function formatDuration(value?: number | null) {
  if (!value) return '--:--'
  const mins = Math.floor(value / 60)
  const secs = value % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}

async function fetchDashboard() {
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

onMounted(fetchDashboard)
</script>
