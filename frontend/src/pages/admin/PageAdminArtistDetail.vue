<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <Button
      label="Back to Artists"
      icon="pi pi-arrow-left"
      severity="secondary"
      text
      size="small"
      class="mb-4"
      @click="router.push({ name: 'admin.artists' })"
    />

    <div v-if="loading" class="space-y-6">
      <div class="flex items-start gap-6">
        <div class="h-40 w-40 animate-pulse rounded-2xl bg-white/6" />
        <div class="flex-1 space-y-3">
          <div class="h-8 w-48 animate-pulse rounded bg-white/6" />
          <div class="h-4 w-72 animate-pulse rounded bg-white/4" />
          <div class="h-20 w-full animate-pulse rounded bg-white/4" />
        </div>
      </div>
    </div>

    <div v-else-if="error" class="text-center">
      <AdminEmptyState
        icon="pi pi-exclamation-triangle"
        title="Failed to load artist"
        :description="error"
      >
        <template #action>
          <Button label="Retry" icon="pi pi-refresh" @click="loadArtist" />
        </template>
      </AdminEmptyState>
    </div>

    <template v-else-if="artist">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-start">
        <div class="shrink-0">
          <div class="relative h-40 w-40 overflow-hidden rounded-2xl bg-white/6 shadow-lg lg:h-56 lg:w-56">
            <img
              v-if="artist.image_url"
              :src="artist.image_url"
              :alt="artist.name"
              class="h-full w-full object-cover"
              @error="($event.target as HTMLImageElement).style.display='none'"
            />
            <div v-else class="flex h-full w-full items-center justify-center">
              <User aria-hidden="true" class="text-5xl text-slate-600"  />
            </div>
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-3">
            <h1 class="text-2xl font-bold text-white lg:text-3xl">{{ artist.name }}</h1>
            <BadgeCheck aria-hidden="true" v-if="artist.is_verified" class="text-emerald-400" v-tooltip.top="'Verified'"  />
          </div>

          <div v-if="artist.monthly_listeners" class="mt-2 text-sm text-slate-400">
            {{ formatListeners(artist.monthly_listeners) }} monthly listeners
          </div>

          <div v-if="artist.bio" class="mt-4 rounded-xl border border-white/6 bg-white/2 p-4">
            <p class="text-xs font-medium text-slate-400">Biography</p>
            <p class="mt-2 text-sm leading-relaxed text-slate-300">{{ artist.bio }}</p>
          </div>

          <div class="mt-6 flex flex-wrap gap-2">
            <Button
              label="Enrich"
              icon="pi pi-refresh"
              size="small"
              severity="info"
              :loading="enriching"
              class="rounded-xl! bg-amber-500/10! text-amber-400! hover:bg-amber-500/20!"
              @click="handleEnrich"
            />
            <Button
              label="Edit Artist"
              icon="pi pi-pencil"
              size="small"
              severity="secondary"
              @click="openEdit"
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

      <div class="mt-10">
        <h2 class="mb-4 text-lg font-semibold text-white">Albums</h2>
        <div v-if="loadingAlbums" class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          <div v-for="i in 5" :key="i" class="space-y-3">
            <div class="aspect-square animate-pulse rounded-xl bg-white/6" />
            <div class="h-4 w-3/4 animate-pulse rounded bg-white/6" />
          </div>
        </div>
        <AdminEmptyState
          v-else-if="albums.length === 0"
          icon="pi pi-book"
          title="No albums"
          description="This artist has no albums yet."
          compact
        />
        <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          <div
            v-for="album in albums"
            :key="album.id"
            class="group cursor-pointer"
            @click="router.push({ name: 'admin.album.detail', params: { id: album.id } })"
          >
            <div class="relative aspect-square overflow-hidden rounded-xl border border-white/6 bg-white/4">
              <img
                v-if="album.cover_url"
                :src="album.cover_url"
                :alt="album.title"
                class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                @error="($event.target as HTMLImageElement).style.display='none'"
              />
              <div v-else class="flex h-full w-full items-center justify-center">
                <Image aria-hidden="true" class="text-3xl text-slate-700"  />
              </div>
            </div>
            <p class="mt-2 truncate text-sm font-medium text-white">{{ album.title }}</p>
          </div>
        </div>
      </div>

      <div class="mt-10">
        <h2 class="mb-4 text-lg font-semibold text-white">Tracks</h2>
        <div v-if="loadingTracks" class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-12 animate-pulse rounded-lg bg-white/4" />
        </div>
        <AdminEmptyState
          v-else-if="tracks.length === 0"
          icon="pi pi-music"
          title="No tracks"
          description="This artist has no tracks yet."
          compact
        />
        <div v-else class="overflow-hidden rounded-xl border border-white/6">
          <div
            v-for="(track, i) in tracks"
            :key="track.id"
            class="flex items-center gap-3 px-4 py-2 transition-colors hover:bg-white/2"
            :class="i < tracks.length - 1 ? 'border-b border-white/4' : ''"
          >
            <button
              type="button"
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-slate-500 transition-all hover:bg-spotify/20 hover:text-spotify disabled:opacity-30"
              :disabled="loadingTrackId === String(track.id)"
              :aria-label="'Play ' + track.title"
              :title="isTrackPlaying(track) ? 'Now playing' : 'Play track'"
              @click.stop="handlePlayTrack(track)"
            >
              <Loader2
                v-if="loadingTrackId === String(track.id)"
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
            <span class="w-5 text-center text-xs text-slate-500">{{ i + 1 }}</span>
            <div
              class="min-w-0 flex-1 cursor-pointer"
              @click="router.push({ name: 'admin.track.detail', params: { id: track.id } })"
            >
              <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
              <p v-if="track.duration_seconds" class="text-xs text-slate-500">{{ formatDuration(track.duration_seconds) }}</p>
            </div>
            <ChevronRight
              aria-hidden="true"
              class="cursor-pointer text-xs text-slate-600"
              @click="router.push({ name: 'admin.track.detail', params: { id: track.id } })"
             />
          </div>
        </div>
      </div>
    </template>

    <ArtistFormDialog
      v-model="showForm"
      :artist="editingArtist"
      :saving="saving"
      @submit="handleEditSubmit"
    />

    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete artist"
      :item-name="artist?.name"
      :deleting="deleting"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { BadgeCheck, ChevronRight, Image, Loader2, User } from 'lucide-vue-next'
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import AdminEmptyState from '@/components/admin/AdminEmptyState.vue'
import ArtistFormDialog from '@/components/admin/ArtistFormDialog.vue'
import type { ArtistFormPayload } from '@/composables/admin/useAdminArtists'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { useAlbumsApi } from '@/services/api/catalog/albums'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { usePlayer } from '@/composables/player'
import { buildPlaybackTrack as factoryBuildPlaybackTrack } from '@/factories/playbackTrack'
import type { Artist } from '@/services/api/catalog/artists'
import type { Album } from '@/services/api/catalog/albums'
import type { Track } from '@/services/api/catalog/tracks'
import { formatDuration } from '@/utils/format'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const artistsApi = useArtistsApi()
const albumsApi = useAlbumsApi()
const tracksApi = useTracksApi()

const artist = ref<Artist | null>(null)
const albums = ref<Album[]>([])
const tracks = ref<Track[]>([])
const loading = ref(true)
const error = ref('')
const loadingAlbums = ref(false)
const loadingTracks = ref(false)

const showForm = ref(false)
const editingArtist = ref<Artist | null>(null)
const saving = ref(false)
const showDelete = ref(false)
const deleting = ref(false)
const enriching = ref(false)
const loadingTrackId = ref<string | null>(null)
const player = usePlayer()

function buildPlaybackTrack(track: Track) {
  return factoryBuildPlaybackTrack({
    id: String(track.id),
    title: track.title || 'Untitled',
    artist_name: artist.value?.name || 'Unknown artist',
    cover_url: track.cover_url || null,
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

function getTrackPlayButtonIcon(track: Track): string {
  if (isTrackPlaying(track) && player.isPlaying.value) return 'pi pi-pause-fill'
  return 'pi pi-play-fill'
}

onMounted(loadArtist)

async function loadArtist() {
  const id = route.params.id as string
  loading.value = true
  error.value = ''
  try {
    const [artistData, allAlbums, allTracks] = await Promise.all([
      artistsApi.getArtist(id),
      albumsApi.getAlbums(),
      tracksApi.getTracks(),
    ])
    artist.value = artistData
    albums.value = allAlbums.filter((a: Record<string, unknown>) => String(a.artist_id) === id)
    tracks.value = allTracks.filter((t: Record<string, unknown>) => String(t.artist_id) === id)
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Failed to load artist.'
  } finally {
    loading.value = false
  }
}

function openEdit() {
  editingArtist.value = artist.value
  showForm.value = true
}

function openDeleteConfirm() {
  showDelete.value = true
}

async function handleEnrich() {
  if (!artist.value) return
  enriching.value = true
  try {
    await artistsApi.adminEnrichArtist(artist.value.id)
    toast.add({ severity: 'success', summary: 'Artist enriched', detail: 'Bio and image fetched from external sources', life: 3000 })
    await loadArtist()
  } catch {
    toast.add({ severity: 'error', summary: 'Enrich failed', detail: 'Could not enrich artist', life: 4000 })
  } finally {
    enriching.value = false
  }
}

async function handleEditSubmit(payload: ArtistFormPayload) {
  if (!artist.value) return
  saving.value = true
  try {
    const updated = await artistsApi.adminUpdateArtist(artist.value.id, payload)
    artist.value = updated
    showForm.value = false
    toast.add({ severity: 'success', summary: 'Artist updated', life: 2500 })
  } catch {
    toast.add({ severity: 'error', summary: 'Update failed', life: 3000 })
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!artist.value) return
  deleting.value = true
  try {
    await artistsApi.adminDeleteArtist(artist.value.id)
    toast.add({ severity: 'success', summary: 'Artist deleted', life: 2500 })
    router.push({ name: 'admin.artists' })
  } catch {
    toast.add({ severity: 'error', summary: 'Delete failed', life: 3000 })
  } finally {
    deleting.value = false
  }
}

function formatListeners(count: number): string {
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}

</script>
