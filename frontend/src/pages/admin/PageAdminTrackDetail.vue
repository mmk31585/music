<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <Button
      label="Back to Tracks"
      icon="pi pi-arrow-left"
      severity="secondary"
      text
      size="small"
      class="mb-4"
      @click="router.push({ name: 'admin.tracks' })"
    />

    <div v-if="loading" class="space-y-6">
      <div class="flex items-start gap-6">
        <div class="aspect-square w-36 animate-pulse rounded-2xl bg-white/6 lg:w-44" />
        <div class="flex-1 space-y-3">
          <div class="h-8 w-64 animate-pulse rounded bg-white/6" />
          <div class="h-4 w-40 animate-pulse rounded bg-white/4" />
          <div class="h-4 w-56 animate-pulse rounded bg-white/4" />
        </div>
      </div>
    </div>

    <div v-else-if="error" class="text-center">
      <AdminEmptyState
        icon="pi pi-exclamation-triangle"
        title="Failed to load track"
        :description="error"
      >
        <template #action>
          <Button label="Retry" icon="pi pi-refresh" @click="loadTrack" />
        </template>
      </AdminEmptyState>
    </div>

    <template v-else-if="track">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-start">
        <div class="shrink-0">
          <div class="relative aspect-square w-36 overflow-hidden rounded-2xl bg-white/6 shadow-lg lg:w-44">
            <img
              v-if="track.cover_url"
              :src="track.cover_url"
              :alt="track.title"
              class="h-full w-full object-cover"
              @error="($event.target as HTMLImageElement).style.display='none'"
            />
            <div v-else class="flex h-full w-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-headphones text-4xl text-slate-600" />
            </div>
          </div>

          <div class="mt-3 space-y-1 text-center text-xs text-slate-500 lg:text-left">
            <div v-if="track.play_count" class="flex items-center gap-1">
              <i aria-hidden="true" class="pi pi-play" />
              {{ formatPlayCount(track.play_count) }} plays
            </div>
            <div v-if="track.duration_seconds" class="flex items-center gap-1">
              <i aria-hidden="true" class="pi pi-clock" />
              {{ formatDuration(track.duration_seconds) }}
            </div>
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <h1 class="text-2xl font-bold text-white lg:text-3xl">{{ track.title }}</h1>

          <div class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-slate-400">
            <span v-if="artist" class="flex items-center gap-1">
              <i aria-hidden="true" class="pi pi-user" />
              <button
                class="cursor-pointer text-teal-400 transition-colors hover:text-teal-300"
                @click="router.push({ name: 'admin.artist.detail', params: { id: artist.id } })"
              >
                {{ artist.name }}
              </button>
            </span>
            <button
              v-if="parentAlbum"
              class="flex items-center gap-1 text-teal-400 transition-colors hover:text-teal-300"
              @click="router.push({ name: 'admin.album.detail', params: { id: parentAlbum.id } })"
            >
              <i aria-hidden="true" class="pi pi-book" />
              {{ parentAlbum.title }}
            </button>
            <span v-if="track.genres.length > 0">
              <i aria-hidden="true" class="pi pi-tag mr-1" />
              {{ track.genres.map(g => g.name).join(', ') }}
            </span>
          </div>

          <div class="mt-6 flex flex-wrap gap-2">
            <Button
              label="Edit Track"
              icon="pi pi-pencil"
              size="small"
              severity="secondary"
              @click="openEdit"
            />
            <Button
              label="Fetch LRC"
              icon="pi pi-cloud-download"
              size="small"
              severity="secondary"
              :loading="fetchingLrc"
              class="text-teal-400!"
              @click="handleFetchLRC"
            />
            <Button
              label="Fetch All"
              icon="pi pi-sync"
              size="small"
              severity="secondary"
              :loading="enriching"
              class="text-purple-400!"
              @click="handleEnrichTrack"
            />
            <Button
              label="Delete"
              icon="pi pi-trash"
              size="small"
              severity="danger"
              class="text-red-400!"
              @click="openDeleteConfirm"
            />
          </div>
        </div>
      </div>

      <div v-if="trackInOtherAlbums.length > 0" class="mt-10">
        <h2 class="mb-4 text-lg font-semibold text-white">Also appears on</h2>
        <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          <div
            v-for="alb in trackInOtherAlbums"
            :key="alb.id"
            class="group cursor-pointer"
            @click="router.push({ name: 'admin.album.detail', params: { id: alb.id } })"
          >
            <div class="relative aspect-square overflow-hidden rounded-xl border border-white/6 bg-white/4">
              <img
                v-if="alb.cover_url"
                :src="alb.cover_url"
                :alt="alb.title"
                class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                @error="($event.target as HTMLImageElement).style.display='none'"
              />
              <div v-else class="flex h-full w-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-image text-3xl text-slate-700" />
              </div>
            </div>
            <p class="mt-2 truncate text-sm font-medium text-white">{{ alb.title }}</p>
          </div>
        </div>
      </div>
    </template>

    <TrackFormDialog
      v-model:visible="showForm"
      mode="edit"
      :track="editingTrack"
      :loading="saving"
      :artists-options="artistOptions"
      :albums-options="albumOptions"
      :genres-options="genreOptions"
      @submit="handleEditSubmit"
      @cancel="closeForm"
    />

    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete track"
      :item-name="track?.title"
      :deleting="deleting"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import AdminEmptyState from '@/components/admin/AdminEmptyState.vue'
import TrackFormDialog from '@/components/admin/TrackFormDialog.vue'
import type { TrackFormPayload } from '@/components/admin/TrackFormDialog.vue'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'
import { useTracksApi, type TrackArtistRequest } from '@/services/api/catalog/tracks'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { useAlbumsApi } from '@/services/api/catalog/albums'
import { useLyricsApi } from '@/services/api/lyrics'
import { useAdminArtists } from '@/composables/admin/useAdminArtists'
import { useAdminAlbums } from '@/composables/admin/useAdminAlbums'
import { useAdminGenres } from '@/composables/admin/useAdminGenres'
import type { Track } from '@/services/api/catalog/tracks'
import type { Artist } from '@/services/api/catalog/artists'
import type { Album } from '@/services/api/catalog/albums'

type CatalogId = string | number

type CatalogOption = {
  id: CatalogId
  name: string
  slug?: string
  image_url?: string | null
  avatar_url?: string | null
  cover_url?: string | null
}

const router = useRouter()
const route = useRoute()
const toast = useToast()
const tracksApi = useTracksApi()
const artistsApi = useArtistsApi()
const albumsApi = useAlbumsApi()
const lyricsApi = useLyricsApi()

const { artists, fetchArtists } = useAdminArtists()
const { albums, fetchAlbums } = useAdminAlbums()
const { genres, fetchGenres } = useAdminGenres()

const track = ref<Track | null>(null)
const artist = ref<Artist | null>(null)
const allAlbums = ref<Album[]>([])
const loading = ref(true)
const error = ref('')

const showForm = ref(false)
const editingTrack = ref<Track | null>(null)
const saving = ref(false)
const showDelete = ref(false)
const deleting = ref(false)
const fetchingLrc = ref(false)
const enriching = ref(false)

const parentAlbum = computed(() =>
  track.value?.album_id
    ? allAlbums.value.find(a => String(a.id) === String(track.value!.album_id)) ?? null
    : null
)

const trackInOtherAlbums = computed(() =>
  allAlbums.value.filter(a => String(a.id) !== String(track.value?.album_id))
)

const artistOptions = computed(() =>
  normalizeCatalogOptions(artists.value, 'artist').filter(Boolean) as CatalogOption[]
)

const albumOptions = computed(() =>
  normalizeCatalogOptions(albums.value, 'album').filter(Boolean) as CatalogOption[]
)

const genreOptions = computed(() =>
  normalizeCatalogOptions(genres.value, 'genre').filter(Boolean) as CatalogOption[]
)

onMounted(() => {
  loadTrack()
  Promise.allSettled([fetchArtists(), fetchAlbums(), fetchGenres()])
})

async function loadTrack() {
  const id = route.params.id as string
  loading.value = true
  error.value = ''
  try {
    track.value = await tracksApi.getTrack(id)
    const [artistData, albumsData] = await Promise.all([
      track.value.artist_id ? artistsApi.getArtist(String(track.value.artist_id)) : Promise.resolve(null),
      albumsApi.getAlbums(),
    ])
    artist.value = artistData
    allAlbums.value = albumsData
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Failed to load track.'
  } finally {
    loading.value = false
  }
}

async function openEdit() {
  // Start with the track data from the detail page
  const augmented: Record<string, unknown> = { ...track.value }

  // Fetch lyrics — they aren't included in getTrack response
  try {
    const lyrics = await lyricsApi.getTrackLyrics(augmented.id as string, undefined, { silent: true })
    if (lyrics?.content) {
      augmented.lyrics = lyrics.content
      augmented.lyrics_language = lyrics.language || 'en'
      augmented.lyrics_type = lyrics.type === 'synced' ? 'synced' : 'plain'
    }
  } catch {
    // No lyrics stored yet — that's fine
  }

  editingTrack.value = augmented as Track
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  editingTrack.value = null
}

function openDeleteConfirm() {
  showDelete.value = true
}

async function handleEditSubmit(payload: TrackFormPayload) {
  if (track.value!) return
  saving.value = true
  try {
    // Build artists array from credits/artist_ids/featured_artist_ids
    let artists: TrackArtistRequest[] | undefined
    if (payload.credits?.length) {
      artists = payload.credits.map((credit, index) => ({
        artistId: credit.artist_id,
        role: credit.role,
        position: index,
      }))
    } else if (payload.artists?.length) {
      artists = payload.artists.map((artist, index) => ({
        artistId: artist.artist_id,
        role: artist.role || (index === 0 ? 'primary' : 'featured'),
        position: artist.position ?? index,
      }))
    } else {
      const primaryId = payload.artist_ids?.[0] ?? payload.artist_id
      if (primaryId) {
        artists = [{ artistId: primaryId, role: 'primary', position: 0 }]
        // Add featured artists
        if (payload.featured_artist_ids?.length) {
          artists.push(
            ...payload.featured_artist_ids.map((id, index) => ({
              artistId: id,
              role: 'featured' as const,
              position: index + 1,
            })),
          )
        }
      }
    }

    // Build the update payload — only send fields that exist in TrackUpdatePayload
    const updatePayload = {
      title: payload.title,
      artistId: payload.artist_ids?.[0] ?? payload.artist_id ?? null,
      artists,
      albumId: payload.album_id ?? null,
      durationSeconds: payload.duration_seconds ?? null,
      audioUrl: payload.audio_url ?? null,
      coverUrl: payload.cover_url ?? null,
      genreIds: payload.genre_ids ?? [],
      trackNumber: payload.track_number ?? null,
      explicit: payload.explicit ?? false,
      isPublic: true,
    }

    const updated = await tracksApi.adminUpdateTrack(track.value.id, updatePayload)
    track.value = updated

    // Sync lyrics — use existing lyric's ID if available, otherwise create fresh
    const content = payload.lyrics?.trim()
    if (content) {
      const language = payload.lyrics_language || 'en'
      const type = payload.lyrics_type || 'plain'

      // Find if there's already a lyric row for this track (of any language)
      let existingId: string | number | null = null
      try {
        const allLyrics = await lyricsApi.getLyricsByTrackID(track.value.id, { silent: true })
        // Prefer a lyric with the target language; fall back to the first one
        const target = allLyrics.find((l) => l.language === language) ?? allLyrics[0]
        existingId = target?.id ?? null
      } catch {
        // no existing lyrics
      }

      if (existingId) {
        // Update the existing row with new language, content, and type
        await lyricsApi.adminUpdateLyrics(existingId, {
          track_id: track.value.id,
          content,
          language,
          type,
        })
      } else {
        await lyricsApi.adminCreateLyrics({ track_id: track.value.id, content, language, type })
      }
    } else if (content === '') {
      // Lyrics cleared — delete all existing lyrics
      try {
        const allLyrics = await lyricsApi.getLyricsByTrackID(track.value.id, { silent: true })
        for (const l of allLyrics) {
          if (l?.id) await lyricsApi.adminDeleteLyrics(l.id)
        }
      } catch { /* nothing to delete */ }
    }

    showForm.value = false
    toast.add({ severity: 'success', summary: 'Track updated', life: 2500 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Update failed', detail: err instanceof Error ? err.message : 'Unknown error', life: 3000 })
  } finally {
    saving.value = false
  }
}

async function handleFetchLRC() {
  if (track.value!) return
  fetchingLrc.value = true
  try {
    // Use fetch-or-generate pipeline: LRCLIB first, then AI fallback
    const result = await lyricsApi.fetchOrGenerateLyrics(String(track.value.id))
    if (result.source === 'lrclib') {
      toast.add({ severity: 'success', summary: 'Lyrics found on LRCLIB and saved', life: 3000 })
    } else if (result.source === 'ai') {
      toast.add({
        severity: 'info',
        summary: 'AI lyrics generation started',
        detail: 'No internet lyrics found. AI is generating from audio... This may take a minute.',
        life: 5000,
      })
      // Poll until done
      let attempts = 0
      const maxAttempts = 30
      const poll = async (): Promise<void> => {
        attempts++
        const status = await lyricsApi.aiStatus(String(track.value!.id))
        if (status.status === 'completed') {
          toast.add({ severity: 'success', summary: 'AI lyrics generated and saved!', life: 3000 })
          await loadTrack()
          return
        }
        if (status.status === 'failed') {
          toast.add({
            severity: 'warn',
            summary: 'AI lyrics generation failed',
            detail: status.message || 'Unknown error',
            life: 5000,
          })
          return
        }
        if (attempts < maxAttempts) {
          setTimeout(poll, 3000)
        } else {
          toast.add({ severity: 'warn', summary: 'AI lyrics generation still in progress', detail: 'Check back later', life: 5000 })
        }
      }
      await poll()
    } else {
      toast.add({ severity: 'warn', summary: 'No lyrics found', detail: result.message || 'Neither LRCLIB nor AI had results', life: 3000 })
    }
    await loadTrack()
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Failed to fetch lyrics',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  } finally {
    fetchingLrc.value = false
  }
}

async function handleEnrichTrack() {
  if (track.value!) return
  enriching.value = true
  try {
    await tracksApi.adminEnrichTrack(track.value.id)
    toast.add({ severity: 'success', summary: 'Track enrichment complete', life: 3000 })
    await loadTrack()
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Enrichment failed',
      detail: err instanceof Error ? err.message : 'Unknown error',
      life: 3000,
    })
  } finally {
    enriching.value = false
  }
}

async function handleDelete() {
  if (track.value!) return
  deleting.value = true
  try {
    await tracksApi.adminDeleteTrack(track.value.id)
    toast.add({ severity: 'success', summary: 'Track deleted', life: 2500 })
    router.push({ name: 'admin.tracks' })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Delete failed', detail: err instanceof Error ? err.message : 'Unknown error', life: 3000 })
  } finally {
    deleting.value = false
  }
}

function normalizeCatalogOptions(items: Array<Record<string, unknown>> | undefined | null, _type?: string) {
  if (Array.isArray!(items)) return []
  return items.map(item => {
    const id = item.id ?? item.artist_id ?? item.album_id ?? item.genre_id
    const name = item.name ?? item.title ?? item.artist_name ?? item.album_title ?? item.genre_name
    if (id === undefined || id === null || name!) return null
    return { id, name, slug: item.slug, image_url: item.image_url ?? item.avatar_url ?? item.cover_url ?? null, avatar_url: item.avatar_url ?? item.image_url ?? null, cover_url: item.cover_url ?? item.image_url ?? null } as CatalogOption
  })
}

function formatPlayCount(count: number): string {
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>
