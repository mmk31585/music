<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Catalog overview"
      description="Browse all public catalog entities and manage content."
    >
      <template #actions>
        <Button
          label="Refresh"
          icon="pi pi-refresh"
          text
          size="small"
          :loading="loading"
          class="text-slate-400! hover:text-white!"
          @click="fetchCatalog"
        />
      </template>
    </AdminSectionHeader>

    <!-- Stats Row -->
    <section class="mb-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard
        label="Tracks"
        :value="tracks.length"
        icon="pi pi-play-circle"
        color="emerald"
        :loading="loading"
      />
      <AdminStatCard
        label="Artists"
        :value="artists.length"
        icon="pi pi-users"
        color="blue"
        :loading="loading"
      />
      <AdminStatCard
        label="Albums"
        :value="albums.length"
        icon="pi pi-book"
        color="purple"
        :loading="loading"
      />
      <AdminStatCard
        label="Genres"
        :value="genres.length"
        icon="pi pi-tags"
        color="amber"
        :loading="loading"
      />
    </section>

    <!-- Content Grid -->
    <section class="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <!-- Tracks Panel -->
      <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
        <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10">
              <i aria-hidden="true" class="pi pi-play-circle text-xs text-emerald-400" />
            </div>
            <h2 class="text-base font-semibold text-white">Tracks</h2>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs tabular-nums text-slate-500">{{ tracks.length }}</span>
            <RouterLink
              to="/admin/tracks"
              class="text-xs font-medium text-emerald-400 transition-colors hover:text-emerald-300"
            >
              View all →
            </RouterLink>
          </div>
        </div>

        <div v-if="loading" class="divide-y divide-white/4">
          <div v-for="i in 6" :key="i" class="flex items-center gap-3 px-5 py-3.5">
            <div class="h-9 w-9 animate-pulse rounded-md bg-white/6" />
            <div class="flex-1 space-y-1.5">
              <div class="h-3.5 w-32 animate-pulse rounded bg-white/6" />
              <div class="h-3 w-44 animate-pulse rounded bg-white/4" />
            </div>
          </div>
        </div>

        <div v-else-if="tracks.length === 0" class="py-12 text-center">
          <i aria-hidden="true" class="pi pi-play-circle text-2xl text-slate-700" />
          <p class="mt-2 text-sm text-slate-500">No tracks found</p>
        </div>

        <div v-else class="max-h-125 divide-y divide-white/4 overflow-y-auto">
          <div
            v-for="(track, i) in tracks"
            :key="track.id"
            class="group flex items-center gap-3 px-5 py-3 transition-colors hover:bg-white/2"
          >
            <button
              type="button"
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-slate-500 transition-all hover:bg-spotify/20 hover:text-spotify disabled:opacity-30"
              :disabled="loadingTrackId === String(track.id)"
              :aria-label="'Play ' + track.title"
              :title="isTrackPlaying(track) ? 'Now playing' : 'Play track'"
              @click.stop="handlePlayTrack(track)"
            >
              <i
                v-if="loadingTrackId === String(track.id)"
                aria-hidden="true"
                class="pi pi-spin pi-spinner text-sm"
              />
              <i
                v-else
                aria-hidden="true"
                :class="getTrackPlayButtonIcon(track)"
                class="text-sm"
              />
            </button>

            <span class="w-5 text-center text-xs tabular-nums text-slate-600">{{ i + 1 }}</span>

            <div class="h-9 w-9 shrink-0 overflow-hidden rounded-md bg-white/4">
              <img
                v-if="track.cover_url"
                :src="track.cover_url"
                :alt="track.title"
                class="h-full w-full object-cover"
                @error="($event.target as HTMLImageElement).style.display = 'none'"
              />
              <div v-else class="flex h-full w-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-xs text-slate-700" />
              </div>
            </div>

            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
              <p class="truncate text-xs text-slate-500">
                {{ track.artist_name || 'Unknown artist' }}
                <span v-if="track.album_title" class="text-slate-700">
                  · {{ track.album_title }}
                </span>
              </p>
            </div>

            <span class="text-xs tabular-nums text-slate-600">
              {{ formatDuration(track.duration_seconds) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Right Column -->
      <div class="space-y-6">
        <!-- Artists Panel -->
        <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
          <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
            <div class="flex items-center gap-3">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10">
                <i aria-hidden="true" class="pi pi-users text-xs text-blue-400" />
              </div>
              <h2 class="text-base font-semibold text-white">Artists</h2>
            </div>
            <RouterLink
              to="/admin/artists"
              class="text-xs font-medium text-emerald-400 transition-colors hover:text-emerald-300"
            >
              Manage →
            </RouterLink>
          </div>

          <div v-if="loading" class="divide-y divide-white/4">
            <div v-for="i in 4" :key="i" class="flex items-center gap-3 px-5 py-3">
              <div class="h-8 w-8 animate-pulse rounded-full bg-white/6" />
              <div class="h-3.5 w-28 animate-pulse rounded bg-white/6" />
            </div>
          </div>

          <div v-else-if="artists.length === 0" class="py-8 text-center">
            <p class="text-sm text-slate-500">No artists</p>
          </div>

          <div v-else class="max-h-60 divide-y divide-white/4 overflow-y-auto">
            <div
              v-for="artist in artists"
              :key="artist.id"
              class="flex items-center gap-3 px-5 py-3 transition-colors hover:bg-white/2"
            >
              <div
                class="flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden rounded-full bg-white/6"
              >
                <img
                  v-if="artist.image_url && artist.image_url.length > 5"
                  :src="artist.image_url"
                  :alt="artist.name"
                  class="h-full w-full object-cover"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <span v-else class="text-xs font-semibold text-slate-500">
                  {{ artist.name.charAt(0).toUpperCase() }}
                </span>
              </div>
              <p class="truncate text-sm text-white">{{ artist.name }}</p>
              <i
                v-if="artist.is_verified"
                class="pi pi-verified text-xs text-emerald-400"
              />
            </div>
          </div>
        </div>

        <!-- Albums Panel -->
        <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
          <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
            <div class="flex items-center gap-3">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-500/10">
                <i aria-hidden="true" class="pi pi-book text-xs text-purple-400" />
              </div>
              <h2 class="text-base font-semibold text-white">Albums</h2>
            </div>
            <RouterLink
              to="/admin/albums"
              class="text-xs font-medium text-emerald-400 transition-colors hover:text-emerald-300"
            >
              Manage →
            </RouterLink>
          </div>

          <div v-if="loading" class="grid grid-cols-3 gap-3 p-5">
            <div v-for="i in 3" :key="i" class="space-y-2">
              <div class="aspect-square animate-pulse rounded-lg bg-white/6" />
              <div class="h-3 w-3/4 animate-pulse rounded bg-white/4" />
            </div>
          </div>

          <div v-else-if="albums.length === 0" class="py-8 text-center">
            <p class="text-sm text-slate-500">No albums</p>
          </div>

          <div v-else class="grid grid-cols-3 gap-3 p-5">
            <div v-for="album in albums.slice(0, 6)" :key="album.id" class="min-w-0">
              <div
                class="aspect-square overflow-hidden rounded-lg border border-white/6 bg-white/4"
              >
                <img
                  v-if="album.cover_url"
                  :src="album.cover_url"
                  :alt="album.title"
                  class="h-full w-full object-cover"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <div v-else class="flex h-full w-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-image text-lg text-slate-700" />
                </div>
              </div>
              <p class="mt-1.5 truncate text-xs font-medium text-white">{{ album.title }}</p>
              <p class="truncate text-xs text-slate-500">{{ album.artist_name || '—' }}</p>
            </div>
          </div>
        </div>

        <!-- Genres Panel -->
        <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
          <div class="flex items-center justify-between border-b border-white/6 px-5 py-4">
            <div class="flex items-center gap-3">
              <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-500/10">
                <i aria-hidden="true" class="pi pi-tags text-xs text-amber-400" />
              </div>
              <h2 class="text-base font-semibold text-white">Genres</h2>
            </div>
            <RouterLink
              to="/admin/genres"
              class="text-xs font-medium text-emerald-400 transition-colors hover:text-emerald-300"
            >
              Manage →
            </RouterLink>
          </div>

          <div v-if="loading" class="flex flex-wrap gap-2 p-5">
            <div
              v-for="i in 6"
              :key="i"
              class="h-8 animate-pulse rounded-full bg-white/6"
              :style="{ width: `${50 + Math.random() * 50}px` }"
            />
          </div>

          <div v-else-if="genres.length === 0" class="py-8 text-center">
            <p class="text-sm text-slate-500">No genres</p>
          </div>

          <div v-else class="flex flex-wrap gap-2 p-5">
            <span
              v-for="genre in genres"
              :key="genre.id"
              class="rounded-full border border-white/8 bg-white/3 px-3.5 py-1.5 text-xs text-slate-300"
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
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import AdminStatCard from '@/components/admin/AdminStatCard.vue'
import { useTracksApi, type Track } from '@/services/api/catalog/tracks'
import { useArtistsApi, type Artist } from '@/services/api/catalog/artists'
import { useAlbumsApi, type Album } from '@/services/api/catalog/albums'
import { useGenresApi, type Genre } from '@/services/api/catalog/genres'
import { usePlayer } from '@/composables/player'
import { buildPlaybackTrack } from '@/factories/playbackTrack'
import { formatDuration } from '@/utils/format'

const { getTracks } = useTracksApi()
const { getArtists } = useArtistsApi()
const { getAlbums } = useAlbumsApi()
const { getGenres } = useGenresApi()

const loading = ref(false)
const loadingTrackId = ref<string | null>(null)
const tracks = ref<Track[]>([])
const artists = ref<Artist[]>([])
const albums = ref<Album[]>([])
const genres = ref<Genre[]>([])
const player = usePlayer()

async function handlePlayTrack(track: Track) {
  loadingTrackId.value = String(track.id)
  try {
    await player.toggleTrack(buildPlaybackTrack(track))
  } finally {
    loadingTrackId.value = null
  }
}

function isTrackPlaying(track: Track): boolean {
  return player.currentTrack.value?.id === String(track.id)
}

function getTrackPlayButtonIcon(track: Track): string {
  if (isTrackPlaying(track) && player.isPlaying.value) return 'pi pi-pause-fill'
  return 'pi pi-play-fill'
}

// TODO HIGH: Fetching ALL catalog items just for preview counts is wasteful.
// Use a dedicated stats endpoint when available. Also, Promise.all will crash all on single failure — use Promise.allSettled instead.
async function fetchCatalog() {
  loading.value = true
  try {
    const [tracksRes, artistsRes, albumsRes, genresRes] = await Promise.all([
      getTracks(),
      getArtists(),
      getAlbums(),
      getGenres(),
    ])
    tracks.value = tracksRes
    artists.value = artistsRes
    albums.value = albumsRes
    genres.value = genresRes
  } finally {
    loading.value = false
  }
}

// TODO MEDIUM: formatDuration treats 0 as falsy — 0-second tracks show '—' instead of '0:00'.
onMounted(fetchCatalog)
</script>
