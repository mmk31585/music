<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Tracks"
      description="Create, edit, and delete tracks in the catalog."
    >
      <template #actions>
        <Button
          label="Add track"
          icon="pi pi-plus"
          size="small"
          class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
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
          class="!h-9 !w-full !rounded-lg !border-white/[0.08] !bg-white/[0.03] !text-sm !text-white placeholder:!text-slate-600 sm:!w-80"
        />
      </div>

      <div class="flex items-center gap-2">
        <Button
          :icon="viewMode === 'table' ? 'pi pi-list' : 'pi pi-th-large'"
          text
          rounded
          size="small"
          class="!text-slate-400"
          v-tooltip.top="viewMode === 'table' ? 'Switch to cards' : 'Switch to table'"
          @click="viewMode = viewMode === 'table' ? 'card' : 'table'"
        />
      </div>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
      <div v-if="loading" class="divide-y divide-white/[0.04]">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 px-5 py-4">
          <div class="h-10 w-10 animate-pulse rounded-lg bg-white/[0.06]" />
          <div class="flex-1 space-y-2">
            <div class="h-4 w-40 animate-pulse rounded bg-white/[0.06]" />
            <div class="h-3 w-56 animate-pulse rounded bg-white/[0.04]" />
          </div>
          <div class="h-4 w-12 animate-pulse rounded bg-white/[0.04]" />
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
            class="!rounded-xl !bg-emerald-500 !px-4 !text-black hover:!bg-emerald-400"
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
          <table class="w-full min-w-[900px]">
            <thead>
              <tr
                class="border-b border-white/[0.06] text-left text-xs tracking-wider text-slate-500 uppercase"
              >
                <th class="px-5 py-3 font-medium">#</th>
                <th class="px-5 py-3 font-medium">Title</th>
                <th class="px-5 py-3 font-medium">Artists</th>
                <th class="hidden px-5 py-3 font-medium lg:table-cell">Album</th>
                <th class="hidden px-5 py-3 font-medium xl:table-cell">Genres</th>
                <th class="px-5 py-3 text-right font-medium">
                  <i aria-hidden="true" class="pi pi-clock text-xs" />
                </th>
                <th class="px-5 py-3" />
              </tr>
            </thead>

            <tbody class="divide-y divide-white/[0.04]">
              <tr
                v-for="(track, index) in filteredTracks"
                :key="track.id"
                class="group transition-colors hover:bg-white/[0.02]"
              >
                <td class="px-5 py-3.5 text-sm text-slate-600 tabular-nums">
                  {{ index + 1 }}
                </td>

                <td class="px-5 py-3.5">
                  <div class="flex min-w-0 items-center gap-3">
                    <div class="h-10 w-10 shrink-0 overflow-hidden rounded-md bg-white/[0.04]">
                      <img
                        v-if="getTrackCoverUrl(track)"
                        :src="getTrackCoverUrl(track)"
                        :alt="track.title"
                        class="h-full w-full object-cover"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                      />

                      <div v-else class="flex h-full w-full items-center justify-center">
                        <i aria-hidden="true" class="pi pi-music text-xs text-slate-700" />
                      </div>
                    </div>

                    <div class="min-w-0">
                      <p class="truncate text-sm font-medium text-white">
                        {{ track.title }}
                      </p>

                      <div class="mt-1 flex flex-wrap items-center gap-1.5">
                        <span
                          v-if="isExplicit(track)"
                          class="rounded bg-white/10 px-1.5 py-0.5 text-[10px] font-semibold text-slate-300"
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
                      class="rounded-full border border-white/[0.06] bg-white/[0.03] px-2 py-0.5 text-xs text-slate-400"
                    >
                      {{ genre }}
                    </span>

                    <span
                      v-if="getTrackGenreNames(track).length > 3"
                      class="rounded-full border border-white/[0.06] bg-white/[0.03] px-2 py-0.5 text-xs text-slate-500"
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
                      class="!h-7 !w-7 !text-slate-400 hover:!text-emerald-400"
                      @click="openEdit(track)"
                    />

                    <Button
                      icon="pi pi-trash"
                      text
                      rounded
                      size="small"
                      class="!h-7 !w-7 !text-slate-400 hover:!text-red-400"
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
            class="relative aspect-square overflow-hidden rounded-xl border border-white/[0.06] bg-white/[0.04]"
          >
            <img
              v-if="getTrackCoverUrl(track)"
              :src="getTrackCoverUrl(track)"
              :alt="track.title"
              class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
              @error="($event.target as HTMLImageElement).style.display = 'none'"
            />

            <div v-else class="flex h-full w-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-music text-3xl text-slate-700" />
            </div>

            <div class="absolute top-2 left-2 flex flex-wrap gap-1">
              <span
                v-if="isExplicit(track)"
                class="rounded bg-black/60 px-1.5 py-0.5 text-[10px] font-semibold text-white backdrop-blur-sm"
              >
                E
              </span>
            </div>

            <div
              class="absolute inset-0 flex items-end justify-end gap-1 bg-gradient-to-t from-black/70 via-transparent p-2.5 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <Button
                icon="pi pi-pencil"
                rounded
                size="small"
                class="!h-7 !w-7 !bg-white/20 !text-white !backdrop-blur-sm hover:!bg-white/30"
                @click="openEdit(track)"
              />

              <Button
                icon="pi pi-trash"
                rounded
                size="small"
                class="!h-7 !w-7 !bg-white/20 !text-white !backdrop-blur-sm hover:!bg-red-500/60"
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
                class="rounded-full bg-white/[0.04] px-2 py-0.5 text-[10px] text-slate-500"
              >
                {{ genre }}
              </span>

              <span
                v-if="getTrackGenreNames(track).length > 2"
                class="rounded-full bg-white/[0.04] px-2 py-0.5 text-[10px] text-slate-600"
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
import { computed, onMounted, ref } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import { useToast } from 'primevue/usetoast'

import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import AdminEmptyState from '@/components/admin/AdminEmptyState.vue'
import TrackFormDialog from '@/components/admin/TrackFormDialog.vue'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'

import { useAdminTracks, type TrackFormPayload } from '@/composables/admin/useAdminTracks'
import { useAdminArtists } from '@/composables/admin/useAdminArtists'
import { useAdminAlbums } from '@/composables/admin/useAdminAlbums'
import { useAdminGenres } from '@/composables/admin/useAdminGenres'

import type { Track } from '@/services/api/catalog/tracks'

type CatalogId = string | number

type CatalogOption = {
  id: CatalogId
  name: string
  slug?: string
  image_url?: string | null
  avatar_url?: string | null
  cover_url?: string | null
}

// TODO HIGH: AnyTrack defeats TypeScript safety. The Track type should match actual API shape,
// or normalizer functions should live in the API service layer. 10+ getter functions add 150+ fragile lines.
type AnyTrack = Track & Record<string, any>

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
const { artists, fetchArtists } = useAdminArtists()

const { albums, fetchAlbums } = useAdminAlbums()

const { genres, fetchGenres } = useAdminGenres()

const searchQuery = ref('')
const viewMode = ref<'table' | 'card'>('table')
const showForm = ref(false)
const showDelete = ref(false)
const selectedTrack = ref<Track | null>(null)
const deleteTarget = ref<Track | null>(null)

const artistOptions = computed(() => {
  return normalizeCatalogOptions(artists.value, 'artist').filter(Boolean) as CatalogOption[]
})

const albumOptions = computed(() => {
  return normalizeCatalogOptions(albums.value, 'album').filter(Boolean) as CatalogOption[]
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

function openEdit(track: Track) {
  selectedTrack.value = track
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
      await updateTrack(selectedTrack.value.id, payload)

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

function normalizeCatalogOptions(
  items: any[] | undefined | null,
) {
  if (!Array.isArray(items)) return []

  return items
    .map((item) => {
      const id = item.id ?? item.artist_id ?? item.album_id ?? item.genre_id

      const name =
        item.name ?? item.title ?? item.artist_name ?? item.album_title ?? item.genre_name

      if (id === undefined || id === null || !name) return null

      const opt: CatalogOption = {
        id,
        name,
        slug: item.slug,
        image_url: item.image_url ?? item.avatar_url ?? item.cover_url ?? null,
        avatar_url: item.avatar_url ?? item.image_url ?? null,
        cover_url: item.cover_url ?? item.image_url ?? null,
      }
      return opt
    })
}

function normalizeSearch(value: string) {
  return value.toLowerCase().trim().replace(/\s+/g, ' ')
}

function compactStrings(values: Array<any>) {
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
  const t = track as AnyTrack

  return (
    t.cover_url ??
    t.cover?.url ??
    t.cover_media?.url ??
    t.cover_media_url ??
    t.media?.cover_url ??
    t.album_cover_url ??
    null
  )
}

function getTrackAlbumTitle(track: Track) {
  const t = track as AnyTrack

  return t.album_title ?? t.album?.title ?? t.album?.name ?? t.album_name ?? '—'
}

function getTrackIsrc(track: Track) {
  const t = track as AnyTrack
  return t.isrc ?? null
}

function isExplicit(track: Track) {
  const t = track as AnyTrack
  return Boolean(t.explicit ?? t.is_explicit)
}

function getTrackArtistNamesByRole(track: Track, roles: string[]) {
  const t = track as AnyTrack

  const normalizedRoles = roles.map((role) => role.toLowerCase())

  const fromArtists = Array.isArray(t.artists)
    ? t.artists
        .filter((item: any) => {
          const it = item as Record<string, any>
          const role = String(it.role ?? '').toLowerCase()

          if (normalizedRoles.includes('primary')) {
            if (it.is_primary === true) return true
          }

          return normalizedRoles.includes(role)
        })
        .map((item: any) => {
          const it = item as Record<string, any>
          return it.name ?? it.artist_name ?? (it.artist as Record<string, any>)?.name ?? (it.artist as Record<string, any>)?.title
        })
    : []

  const fromCredits = Array.isArray(t.credits)
    ? t.credits
        .filter((item: any) => normalizedRoles.includes(String((item as Record<string, any>).role ?? '').toLowerCase()))
        .map((item: any) => {
          const it = item as Record<string, any>
          return it.name ?? it.artist_name ?? (it.artist as Record<string, any>)?.name ?? (it.artist as Record<string, any>)?.title
        })
    : []

  return uniqStrings(compactStrings([...fromArtists, ...fromCredits]))
}

function getTrackPrimaryArtistNames(track: Track) {
  const t = track as AnyTrack

  const names = getTrackArtistNamesByRole(track, ['primary'])

  if (names.length) return names

  return uniqStrings(
    compactStrings([
      t.artist_name,
      t.primary_artist_name,
      t.artist?.name,
      ...(Array.isArray(t.primary_artists)
        ? t.primary_artists.map((item: any) => {
            const it = item as Record<string, any>
            return it.name ?? it.artist_name ?? (it.artist as Record<string, any>)?.name
          })
        : []),
    ]),
  )
}

function getTrackFeaturedArtistNames(track: Track) {
  const t = track as AnyTrack

  const names = getTrackArtistNamesByRole(track, ['featured'])

  if (names.length) return names

  return uniqStrings(
    compactStrings([
      ...(Array.isArray(t.featured_artists)
        ? t.featured_artists.map((item: any) => {
            const it = item as Record<string, any>
            return it.name ?? it.artist_name ?? (it.artist as Record<string, any>)?.name
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
  const t = track as AnyTrack

  const fromGenres = Array.isArray(t.genres)
    ? t.genres.map((item: any) => {
        const it = item as Record<string, any>
        return it.name ?? it.genre_name ?? (it.genre as Record<string, any>)?.name
      })
    : []

  const fromGenreNames = Array.isArray(t.genre_names) ? t.genre_names : []

  const legacy = compactStrings([t.genre_name, t.genre])

  return uniqStrings(compactStrings([...fromGenres, ...fromGenreNames, ...legacy]))
}

function getTrackCreditNames(track: Track) {
  const t = track as AnyTrack

  if (!Array.isArray(t.credits)) return []

  return uniqStrings(
    compactStrings(
      t.credits.map((item: any) => {
        const it = item as Record<string, any>
        return it.name ?? it.artist_name ?? (it.artist as Record<string, any>)?.name
      }),
    ),
  )
}

function getTrackSearchText(track: Track) {
  const t = track as AnyTrack

  const parts = compactStrings([
    t.title,
    getTrackAlbumTitle(track),
    getTrackArtistDisplay(track),
    getTrackGenreNames(track).join(' '),
    getTrackCreditNames(track).join(' '),
    t.album_artist,
    t.album_artist_name,
    t.album_artist?.name,
    t.composer,
    t.isrc,
    t.label,
    t.language,
    t.release_date,
  ])

  return normalizeSearch(parts.join(' '))
}

// TODO MEDIUM: formatDuration treats 0 as falsy — 0-second tracks show '—'.
function formatDuration(value?: number | null): string {
  if (value === null || value === undefined || value <= 0) return '—'

  const mins = Math.floor(value / 60)
  const secs = value % 60

  return `${mins}:${String(secs).padStart(2, '0')}`
}
</script>
