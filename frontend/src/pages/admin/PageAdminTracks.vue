<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Tracks"
      description="Create, edit, and delete tracks in the catalog."
    >
      <template #actions>
        <Button
          label="Enrich All"
          icon="pi pi-sync"
          size="small"
          severity="secondary"
          :loading="enrichingAll"
          class="text-purple-400! rounded-xl!"
          @click="handleEnrichAll"
        />
        <Button
          label="Add track"
          icon="pi pi-plus"
          size="small"
          class="rounded-xl! bg-emerald-500! px-4! text-black! hover:bg-emerald-400!"
          @click="openCreate"
        />
      </template>
    </AdminSectionHeader>

    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-wrap items-center gap-3">
        <span class="text-sm text-slate-500 tabular-nums">
          {{ filteredTracks.length }} track{{ filteredTracks.length !== 1 ? 's' : '' }}
        </span>

        <div class="hidden h-4 w-px bg-white/10 sm:block" />

        <InputText
          v-model="searchQuery"
          placeholder="Search tracks, artists, albums, genres..."
          aria-label="Search tracks"
          class="h-9! w-full! rounded-lg! border-white/8! bg-white/3! text-sm! text-white! placeholder:text-slate-600! sm:w-80!"
        />
      </div>

      <div class="flex items-center gap-2">
        <Button
          :icon="viewMode === 'table' ? 'pi pi-list' : 'pi pi-th-large'"
          text
          rounded
          size="small"
          class="text-slate-400!"
          v-tooltip.top="viewMode === 'table' ? 'Switch to cards' : 'Switch to table'"
          @click="viewMode = viewMode === 'table' ? 'card' : 'table'"
        />
      </div>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
      <div v-if="loading" class="divide-y divide-white/4">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 px-5 py-4">
          <div class="h-10 w-10 animate-pulse rounded-lg bg-white/6" />
          <div class="flex-1 space-y-2">
            <div class="h-4 w-40 animate-pulse rounded bg-white/6" />
            <div class="h-3 w-56 animate-pulse rounded bg-white/4" />
          </div>
          <div class="h-4 w-12 animate-pulse rounded bg-white/4" />
        </div>
      </div>

      <AdminEmptyState
        v-else-if="filteredTracks.length === 0 && !searchQuery"
        icon="pi pi-play-circle"
        title="No tracks yet"
        description="Upload your first track to get started."
      >
        <template #action>
          <Button
            label="Add track"
            icon="pi pi-plus"
            size="small"
            class="rounded-xl! bg-emerald-500! px-4! text-black! hover:bg-emerald-400!"
            @click="openCreate"
          />
        </template>
      </AdminEmptyState>

      <AdminEmptyState
        v-else-if="filteredTracks.length === 0 && searchQuery"
        icon="pi pi-search"
        title="No results found"
        :description="`No tracks matching &quot;${searchQuery}&quot;`"
      />

      <template v-else-if="viewMode === 'table'">
        <div class="overflow-x-auto">
          <table class="w-full min-w-225">
            <thead>
              <tr
                class="border-b border-white/6 text-left text-xs tracking-wider text-slate-500 uppercase"
              >
                <th class="px-5 py-3 font-medium">
                  <Play aria-hidden="true" class="text-xs"  />
                </th>
                <th class="px-5 py-3 font-medium">Title</th>
                <th class="px-5 py-3 font-medium">Artists</th>
                <th class="hidden px-5 py-3 font-medium lg:table-cell">Album</th>
                <th class="hidden px-5 py-3 font-medium xl:table-cell">Genres</th>
                <th class="px-5 py-3 text-right font-medium">
                  <Clock aria-hidden="true" class="text-xs"  />
                </th>
                <th class="px-5 py-3" />
              </tr>
            </thead>

            <tbody class="divide-y divide-white/4">
              <tr
                v-for="(track, index) in filteredTracks"
                :key="track.id"
                class="group transition-colors hover:bg-white/2"
              >
                <td class="px-5 py-3.5">
                  <button
                    type="button"
                    class="flex h-8 w-8 items-center justify-center rounded-lg text-slate-500 transition-all hover:bg-spotify/20 hover:text-spotify disabled:opacity-30"
                    :disabled="isTrackLoading(track)"
                    :aria-label="'Play ' + track.title"
                    :title="isTrackPlaying(track) ? 'Now playing' : 'Play track'"
                    @click="handlePlayTrack(track)"
                  >
                    <Loader2
                      v-if="isTrackLoading(track)"
                      aria-hidden="true"
                      class="text-sm animate-spin"
                     />
                    <i
                      v-else
                      aria-hidden="true"
                      :class="getTrackPlayButtonIcon(track)"
                      class="text-sm"
                    />
                  </button>
                </td>

                <td class="px-5 py-3.5">
                  <div class="flex min-w-0 items-center gap-3">
                    <div class="h-10 w-10 shrink-0 overflow-hidden rounded-md bg-white/4">
                      <img
                        v-if="getTrackCoverUrl(track)"
                        :src="getTrackCoverUrl(track)"
                        :alt="track.title"
                        class="h-full w-full object-cover"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                      />

                      <div v-else class="flex h-full w-full items-center justify-center">
                        <Headphones aria-hidden="true" class="text-xs text-slate-700"  />
                      </div>
                    </div>

                    <div class="min-w-0">
                      <p class="truncate text-sm font-medium text-white">
                        {{ track.title }}
                      </p>

                      <div class="mt-1 flex flex-wrap items-center gap-1.5">
                        <span
                          v-if="isExplicit(track)"
                          class="rounded-sm bg-white/10 px-1.5 py-0.5 text-[10px] font-semibold text-slate-300"
                        >
                          E
                        </span>

                        <span v-if="getTrackIsrc(track)" class="text-[11px] text-slate-600">
                          {{ getTrackIsrc(track) }}
                        </span>
                      </div>
                    </div>
                  </div>
                </td>

                <td class="px-5 py-3.5">
                  <div class="min-w-0">
                    <p class="truncate text-sm text-slate-300">
                      {{ getTrackPrimaryArtistDisplay(track) }}
                    </p>

                    <p
                      v-if="getTrackFeaturedArtistDisplay(track)"
                      class="mt-0.5 truncate text-xs text-slate-500"
                    >
                      feat. {{ getTrackFeaturedArtistDisplay(track) }}
                    </p>
                  </div>
                </td>

                <td class="hidden px-5 py-3.5 lg:table-cell">
                  <span class="truncate text-sm text-slate-500">
                    {{ getTrackAlbumTitle(track) }}
                  </span>
                </td>

                <td class="hidden px-5 py-3.5 xl:table-cell">
                  <div class="flex max-w-xs flex-wrap gap-1.5">
                    <span
                      v-for="genre in getTrackGenreNames(track).slice(0, 3)"
                      :key="genre"
                      class="rounded-full border border-white/6 bg-white/3 px-2 py-0.5 text-xs text-slate-400"
                    >
                      {{ genre }}
                    </span>

                    <span
                      v-if="getTrackGenreNames(track).length > 3"
                      class="rounded-full border border-white/6 bg-white/3 px-2 py-0.5 text-xs text-slate-500"
                    >
                      +{{ getTrackGenreNames(track).length - 3 }}
                    </span>

                    <span
                      v-if="getTrackGenreNames(track).length === 0"
                      class="text-sm text-slate-600"
                    >
                      —
                    </span>
                  </div>
                </td>

                <td class="px-5 py-3.5 text-right text-sm text-slate-500 tabular-nums">
                  {{ formatDuration(track.duration_seconds) }}
                </td>

                <td class="px-5 py-3.5">
                  <div
                    class="flex items-center justify-end gap-1 opacity-0 transition-opacity group-hover:opacity-100"
                  >
                    <Button
                      icon="pi pi-pencil"
                      text
                      rounded
                      size="small"
                      class="h-7! w-7! text-slate-400! hover:text-emerald-400!"
                      @click="openEdit(track)"
                    />

                    <Button
                      icon="pi pi-trash"
                      text
                      rounded
                      size="small"
                      class="h-7! w-7! text-slate-400! hover:text-red-400!"
                      @click="openDeleteConfirm(track)"
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>

      <div v-else class="grid grid-cols-2 gap-4 p-5 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
        <div v-for="track in filteredTracks" :key="track.id" class="group">
          <div
            class="relative aspect-square overflow-hidden rounded-xl border border-white/6 bg-white/4"
          >
            <img
              v-if="getTrackCoverUrl(track)"
              :src="getTrackCoverUrl(track)"
              :alt="track.title"
              class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
              @error="($event.target as HTMLImageElement).style.display = 'none'"
            />

            <div v-else class="flex h-full w-full items-center justify-center">
              <Headphones aria-hidden="true" class="text-3xl text-slate-700"  />
            </div>

            <div class="absolute top-2 left-2 flex flex-wrap gap-1">
              <span
                v-if="isExplicit(track)"
                class="rounded-sm bg-black/60 px-1.5 py-0.5 text-[10px] font-semibold text-white backdrop-blur-xs"
              >
                E
              </span>
            </div>

            <div
              class="absolute inset-0 flex items-center justify-center gap-2 bg-black/40 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <button
                type="button"
                class="flex h-12 w-12 items-center justify-center rounded-full bg-spotify text-black shadow-xl transition-all hover:scale-110 active:scale-90 disabled:opacity-40"
                :disabled="isTrackLoading(track)"
                aria-label="Play track"
                @click="handlePlayTrack(track)"
              >
                <Loader2
                  v-if="isTrackLoading(track)"
                  aria-hidden="true"
                  class="text-lg animate-spin"
                 />
                <i
                  v-else
                  aria-hidden="true"
                  :class="getTrackPlayButtonIcon(track)"
                  class="ml-0.5 text-lg"
                />
              </button>
            </div>

            <div
              class="absolute right-2 top-2 flex gap-1 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <Button
                icon="pi pi-pencil"
                rounded
                size="small"
                class="h-7! w-7! bg-white/20! text-white! backdrop-blur-xs! hover:bg-white/30!"
                @click="openEdit(track)"
              />

              <Button
                icon="pi pi-trash"
                rounded
                size="small"
                class="h-7! w-7! bg-white/20! text-white! backdrop-blur-xs! hover:bg-red-500/60!"
                @click="openDeleteConfirm(track)"
              />
            </div>
          </div>

          <div class="mt-2 min-w-0">
            <p class="truncate text-sm font-medium text-white">
              {{ track.title }}
            </p>

            <p class="mt-0.5 truncate text-xs text-slate-500">
              {{ getTrackArtistDisplay(track) }}
            </p>

            <p class="mt-0.5 truncate text-xs text-slate-600">
              <span v-if="getTrackAlbumTitle(track) !== '—'">
                {{ getTrackAlbumTitle(track) }}
              </span>

              <span v-if="getTrackAlbumTitle(track) !== '—' && track.duration_seconds"> · </span>

              <span v-if="track.duration_seconds">
                {{ formatDuration(track.duration_seconds) }}
              </span>
            </p>

            <div v-if="getTrackGenreNames(track).length" class="mt-1.5 flex flex-wrap gap-1">
              <span
                v-for="genre in getTrackGenreNames(track).slice(0, 2)"
                :key="genre"
                class="rounded-full bg-white/4 px-2 py-0.5 text-[10px] text-slate-500"
              >
                {{ genre }}
              </span>

              <span
                v-if="getTrackGenreNames(track).length > 2"
                class="rounded-full bg-white/4 px-2 py-0.5 text-[10px] text-slate-600"
              >
                +{{ getTrackGenreNames(track).length - 2 }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <TrackFormDialog
      v-model:visible="showForm"
      :mode="selectedTrack ? 'edit' : 'create'"
      :track="selectedTrack"
      :loading="saving"
      :artists-options="artistOptions"
      :albums-options="albumOptions"
      :genres-options="genreOptions"
      @submit="handleSubmit"
      @cancel="closeForm"
      @create-artist="handleCreateArtist"
      @create-album="handleCreateAlbum"
    />

    <ArtistFormDialog
      v-model="showArtistForm"
      :saving="artistSaving"
      :prefill-name="artistPrefill"
      @submit="handleArtistSubmit"
    />

    <AlbumFormDialog
      v-model="showAlbumForm"
      :saving="albumSaving"
      :prefill-title="albumPrefill"
      @submit="handleAlbumSubmit"
    />

    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete track"
      :item-name="deleteTarget?.title"
      :deleting="deleting"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { Clock, Headphones, Loader2, Play } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'

import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import AdminEmptyState from '@/components/admin/AdminEmptyState.vue'
import TrackFormDialog from '@/components/admin/TrackFormDialog.vue'
import ArtistFormDialog from '@/components/admin/ArtistFormDialog.vue'
import AlbumFormDialog from '@/components/admin/AlbumFormDialog.vue'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'

import { useAdminTracks, type TrackFormPayload } from '@/composables/admin/useAdminTracks'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { useLyricsApi } from '@/services/api/lyrics'
import { useAdminArtists } from '@/composables/admin/useAdminArtists'
import { useAdminAlbums } from '@/composables/admin/useAdminAlbums'
import { useAdminGenres } from '@/composables/admin/useAdminGenres'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { buildPlaybackTrack as factoryBuildPlaybackTrack } from '@/factories/playbackTrack'

import type { Track } from '@/services/api/catalog/tracks'
import type { ArtistFormPayload } from '@/composables/admin/useAdminArtists'
import type { AlbumFormPayload } from '@/composables/admin/useAdminAlbums'
import { formatDuration } from '@/utils/format'

type CatalogId = string | number

type CatalogOption = {
  id: CatalogId
  name: string
  slug?: string
  image_url?: string | null
  avatar_url?: string | null
  cover_url?: string | null
}

/**
 * Admin-specific Track type that extends the base Track with all the
 * extra optional fields the admin API response may include.
 * Replace AnyTrack with this proper interface for type safety.
 */
interface AdminTrack extends Track {
  // Cover art fallback shapes
  cover?: { url?: string | null }
  cover_media?: { url?: string | null }
  cover_media_url?: string | null
  media?: { cover_url?: string | null }
  album_cover_url?: string | null

  // Album-related shapes
  album?: { title?: string; name?: string }
  album_name?: string | null

  // Artist shapes
  artist?: { name?: string }
  primary_artist_name?: string | null
  primary_artists?: Array<Record<string, unknown>>
  featured_artists?: Array<Record<string, unknown>>
  artists?: Array<{
    artist_id: string | number | null
    name: string
    role?: string
    is_primary?: boolean
    artist?: Record<string, unknown>
  }>
  credits?: Array<{
    role?: string
    name?: string
    artist_name?: string
    artist?: Record<string, unknown>
  }>

  // Genre shapes
  genre_names?: string[]
  genre_name?: string | null
  genre?: string | null

  // ISRC / explicit
  isrc?: string | null
  is_explicit?: boolean

  // Lyrics
  lyrics?: string | null
  lyrics_language?: string | null
  lyrics_type?: 'plain' | 'synced' | null

  // Search / metadata
  album_artist?: string | { name?: string }
  album_artist_name?: string | null
  composer?: string | null
  label?: string | null
  language?: string | null
  release_date?: string | null
}

const toast = useToast()

const { tracks, loading, saving, deleting, fetchTracks, createTrack, updateTrack, deleteTrack } =
  useAdminTracks()

/**
 * Adjust these destructurings if your composables expose different names.
 *
 * Expected:
 * - artists: array of artist records
 * - albums: array of album records
 * - genres: array of genre records
 * - fetchArtists/fetchAlbums/fetchGenres: loaders
 */
const {
  artists: adminArtistsList,
  fetchArtists,
  createArtist,
  saving: artistSaving,
  error: artistError,
} = useAdminArtists()

const {
  albums: adminAlbumsList,
  fetchAlbums,
  createAlbum,
  saving: albumSaving,
  error: albumError,
} = useAdminAlbums()

const { genres, fetchGenres } = useAdminGenres()

const tracksApi = useTracksApi()
const { getTrackLyrics } = useLyricsApi()
const player = usePlayer()
const playerApi = usePlayerApi()
const loadingTrackId = ref<string | null>(null)
const searchQuery = ref('')
const viewMode = ref<'table' | 'card'>('table')
const showForm = ref(false)
const showDelete = ref(false)
const enrichingAll = ref(false)
const selectedTrack = ref<Track | null>(null)
const deleteTarget = ref<Track | null>(null)
const showArtistForm = ref(false)
const showAlbumForm = ref(false)
const artistPrefill = ref('')
const albumPrefill = ref('')

function buildPlaybackTrack(track: Track) {
  return factoryBuildPlaybackTrack({
    id: String(track.id),
    title: track.title || 'Untitled',
    artist_name: getTrackPrimaryArtistDisplay(track) || 'Unknown artist',
    cover_url: getTrackCoverUrl(track),
    duration_seconds: track.duration_seconds ?? null,
  })
}

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

function isTrackLoading(track: Track): boolean {
  return loadingTrackId.value === String(track.id)
}

function getTrackPlayButtonIcon(track: Track): string {
  if (isTrackPlaying(track) && player.isPlaying.value) return 'pi pi-pause'
  return 'pi pi-play'
}

const artistOptions = computed(() => {
  return normalizeCatalogOptions(adminArtistsList.value, 'artist').filter(Boolean) as CatalogOption[]
})

const albumOptions = computed(() => {
  return normalizeCatalogOptions(adminAlbumsList.value, 'album').filter(Boolean) as CatalogOption[]
})

const genreOptions = computed(() => {
  return normalizeCatalogOptions(genres.value, 'genre').filter(Boolean) as CatalogOption[]
})

const filteredTracks = computed(() => {
  const q = normalizeSearch(searchQuery.value)

  if (!q) return tracks.value

  return tracks.value.filter((track) => {
    return getTrackSearchText(track).includes(q)
  })
})

onMounted(async () => {
  await Promise.allSettled([fetchTracks(), fetchArtists(), fetchAlbums(), fetchGenres()])
})

function openCreate() {
  selectedTrack.value = null
  showForm.value = true
}

async function openEdit(track: Track) {
  // Start with the track data from the list
  const augmented: AdminTrack = { ...track }

  // Fetch lyrics from the backend — they aren't included in the tracks list
  try {
    const lyrics = await getTrackLyrics(track.id, undefined, { silent: true })
    if (lyrics?.content) {
      augmented.lyrics = lyrics.content
      augmented.lyrics_language = lyrics.language || 'en'
      augmented.lyrics_type = lyrics.type === 'synced' ? 'synced' : 'plain'
    }
  } catch (err) {
    console.error('Failed to fetch lyrics:', err)
  }

  selectedTrack.value = augmented as Track
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  selectedTrack.value = null
}

function openDeleteConfirm(track: Track) {
  deleteTarget.value = track
  showDelete.value = true
}

async function handleSubmit(payload: TrackFormPayload) {
  try {
    if (selectedTrack.value) {
      await updateTrack(selectedTrack.value.id, payload, selectedTrack.value)

      toast.add({
        severity: 'success',
        summary: 'Track updated',
        detail: payload.lyrics ? 'Track and lyrics were updated.' : 'Track updated successfully.',
        life: 2500,
      })
    } else {
      await createTrack(payload)

      toast.add({
        severity: 'success',
        summary: 'Track created',
        detail: payload.lyrics ? 'Track and lyrics were created.' : 'Track created successfully.',
        life: 2500,
      })
    }

    showForm.value = false
    selectedTrack.value = null

    await fetchTracks()
  } catch (err) {
    console.error('Track operation failed:', err)

    toast.add({
      severity: 'error',
      summary: selectedTrack.value ? 'Update failed' : 'Create failed',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  }
}

function handleCreateArtist(name: string) {
  artistPrefill.value = name
  showArtistForm.value = true
}

function handleCreateAlbum(title: string) {
  albumPrefill.value = title
  showAlbumForm.value = true
}

async function handleArtistSubmit(payload: ArtistFormPayload) {
  try {
    await createArtist(payload)
    showArtistForm.value = false
    artistPrefill.value = ''

    toast.add({
      severity: 'success',
      summary: 'Artist created',
      life: 2000,
    })

    await fetchArtists()
  } catch (err) {
    console.error('Create artist failed:', err)
    toast.add({
      severity: 'error',
      summary: 'Create artist failed',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  }
}

async function handleAlbumSubmit(payload: AlbumFormPayload) {
  try {
    await createAlbum(payload)
    showAlbumForm.value = false
    albumPrefill.value = ''

    toast.add({
      severity: 'success',
      summary: 'Album created',
      life: 2000,
    })

    await fetchAlbums()
  } catch (err) {
    console.error('Create album failed:', err)
    toast.add({
      severity: 'error',
      summary: 'Create album failed',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  }
}

async function handleDelete() {
  if (!deleteTarget.value) return

  try {
    await deleteTrack(deleteTarget.value.id)

    toast.add({
      severity: 'success',
      summary: 'Track deleted',
      life: 2500,
    })

    showDelete.value = false
    deleteTarget.value = null

    await fetchTracks()
  } catch (err) {
    console.error('Delete track failed:', err)

    toast.add({
      severity: 'error',
      summary: 'Delete failed',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  }
}

async function handleEnrichAll() {
  enrichingAll.value = true
  try {
    await tracksApi.adminEnrichAllTracks()
    toast.add({
      severity: 'success',
      summary: 'Enrichment started',
      detail: 'All tracks are being enriched in the background.',
      life: 4000,
    })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Enrichment failed to start',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  } finally {
    enrichingAll.value = false
  }
}

function normalizeCatalogOptions(
  items: Array<Record<string, unknown>> | undefined | null,
  _type?: string,
) {
  if (!Array.isArray(items)) return []

  return (items ?? [])
    .map((item) => {
      const id = item.id ?? item.artist_id ?? item.album_id ?? item.genre_id

      const name =
        item.name ?? item.title ?? item.artist_name ?? item.album_title ?? item.genre_name

      if (id === undefined || id === null || !name) return null

      const opt: CatalogOption = {
        id: id as CatalogId,
        name: name as string,
        slug: (item.slug as string | undefined) ?? undefined,
        image_url: (item.image_url as string | null | undefined) ?? null,
        avatar_url: (item.avatar_url as string | null | undefined) ?? null,
        cover_url: (item.cover_url as string | null | undefined) ?? null,
      }
      return opt
    })
}

function normalizeSearch(value: string) {
  return value.toLowerCase().trim().replace(/\s+/g, ' ')
}

function compactStrings(values: (string | null | undefined)[]) {
  return values
    .filter((value) => value !== null && value !== undefined && String(value).trim().length > 0)
    .map((value) => String(value).trim())
}

function uniqStrings(values: string[]) {
  const seen = new Set<string>()

  return values.filter((value) => {
    const key = normalizeSearch(value)

    if (seen.has(key)) return false

    seen.add(key)
    return true
  })
}

function getTrackCoverUrl(track: Track) {
  const t = track as AdminTrack

  return (
    t.cover_url ??
    t.cover?.url ??
    t.cover_media?.url ??
    t.cover_media_url ??
    t.media?.cover_url ??
    t.album_cover_url ??
    undefined
  )
}

function getTrackAlbumTitle(track: Track) {
  const t = track as AdminTrack

  return t.album_title ?? t.album?.title ?? t.album?.name ?? t.album_name ?? '—'
}

function getTrackIsrc(track: Track) {
  const t = track as AdminTrack
  return t.isrc ?? null
}

function isExplicit(track: Track) {
  const t = track as AdminTrack
  return Boolean(t.explicit ?? t.is_explicit)
}

function getTrackArtistNamesByRole(track: Track, roles: string[]) {
  const t = track as AdminTrack

  const normalizedRoles = roles.map((role) => role.toLowerCase())

  const fromArtists = Array.isArray(t.artists)
    ? t.artists
        .filter((item: Record<string, unknown>) => {
          const role = String(item.role ?? '').toLowerCase()

          if (normalizedRoles.includes('primary')) {
            if (item.is_primary === true) return true
          }

          return normalizedRoles.includes(role)
        })
        .map((item: Record<string, unknown>) => {
          const artist = item.artist as Record<string, unknown> | undefined
          return String(item.name ?? item.artist_name ?? artist?.name ?? artist?.title ?? '')
        })
    : []

  const fromCredits = Array.isArray(t.credits)
    ? t.credits
        .filter((item: Record<string, unknown>) => normalizedRoles.includes(String(item.role ?? '').toLowerCase()))
        .map((item: Record<string, unknown>) => {
          const artist = item.artist as Record<string, unknown> | undefined
          return String(item.name ?? item.artist_name ?? artist?.name ?? artist?.title ?? '')
        })
    : []

  return uniqStrings(compactStrings([...fromArtists, ...fromCredits]))
}

function getTrackPrimaryArtistNames(track: Track) {
  const t = track as AdminTrack

  const names = getTrackArtistNamesByRole(track, ['primary'])

  if (names.length) return names

  return uniqStrings(
    compactStrings([
      t.artist_name,
      t.primary_artist_name,
      String(t.artist?.name ?? ''),
      ...(Array.isArray(t.primary_artists)
        ? t.primary_artists.map((item: Record<string, unknown>) => {
            const artist = item.artist as Record<string, unknown> | undefined
            return String(item.name ?? item.artist_name ?? artist?.name ?? '')
          })
        : []),
    ]),
  )
}

function getTrackFeaturedArtistNames(track: Track) {
  const t = track as AdminTrack

  const names = getTrackArtistNamesByRole(track, ['featured'])

  if (names.length) return names

  return uniqStrings(
    compactStrings([
      ...(Array.isArray(t.featured_artists)
        ? t.featured_artists.map((item: Record<string, unknown>) => {
            const artist = item.artist as Record<string, unknown> | undefined
            return String(item.name ?? item.artist_name ?? artist?.name ?? '')
          })
        : []),
    ]),
  )
}

function getTrackPrimaryArtistDisplay(track: Track) {
  const names = getTrackPrimaryArtistNames(track)
  return names.length ? names.join(', ') : 'Unknown artist'
}

function getTrackFeaturedArtistDisplay(track: Track) {
  return getTrackFeaturedArtistNames(track).join(', ')
}

function getTrackArtistDisplay(track: Track) {
  const primary = getTrackPrimaryArtistDisplay(track)
  const featured = getTrackFeaturedArtistDisplay(track)

  if (!featured) return primary

  return `${primary} feat. ${featured}`
}

function getTrackGenreNames(track: Track) {
  const t = track as AdminTrack

  const fromGenres = Array.isArray(t.genres)
    ? t.genres.map((item: Record<string, unknown>) => {
        const genre = item.genre as Record<string, unknown> | undefined
        return String(item.name ?? item.genre_name ?? genre?.name ?? '')
      })
    : []

  const fromGenreNames = Array.isArray(t.genre_names) ? t.genre_names : []

  const legacy = compactStrings([t.genre_name, String(t.genre ?? '')])

  return uniqStrings(compactStrings([...fromGenres, ...fromGenreNames, ...legacy]))
}

function getTrackCreditNames(track: Track) {
  const t = track as AdminTrack

  if (!Array.isArray(t.credits)) return []

  return uniqStrings(
    compactStrings(
      t.credits.map((item: Record<string, unknown>) => {
        const artist = item.artist as Record<string, unknown> | undefined
        return String(item.name ?? item.artist_name ?? artist?.name ?? '')
      }),
    ),
  )
}

function getTrackSearchText(track: Track) {
  const t = track as AdminTrack

  const parts = compactStrings([
    t.title,
    getTrackAlbumTitle(track),
    getTrackArtistDisplay(track),
    getTrackGenreNames(track).join(' '),
    getTrackCreditNames(track).join(' '),
    t.album_artist_name,
    typeof t.album_artist === 'object' && t.album_artist ? t.album_artist.name : t.album_artist,
    t.composer,
    t.isrc,
    t.label,
    t.language,
    t.release_date,
  ])

  return normalizeSearch(parts.join(' '))
}

// TODO MEDIUM: formatDuration treats 0 as falsy — 0-second tracks show '—'.
</script>
