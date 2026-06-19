<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <Button
      label="Back to Albums"
      icon="pi pi-arrow-left"
      severity="secondary"
      text
      size="small"
      class="mb-4"
      @click="router.push({ name: 'admin.albums' })"
    />

    <div v-if="loading" class="space-y-6">
      <div class="flex items-start gap-6">
        <div class="aspect-square w-48 animate-pulse rounded-2xl bg-white/[0.06]" />
        <div class="flex-1 space-y-3">
          <div class="h-8 w-56 animate-pulse rounded bg-white/[0.06]" />
          <div class="h-4 w-40 animate-pulse rounded bg-white/[0.04]" />
          <div class="h-4 w-64 animate-pulse rounded bg-white/[0.04]" />
        </div>
      </div>
    </div>

    <div v-else-if="error" class="text-center">
      <AdminEmptyState
        icon="pi pi-exclamation-triangle"
        title="Failed to load album"
        :description="error"
      >
        <template #action>
          <Button label="Retry" icon="pi pi-refresh" @click="loadAlbum" />
        </template>
      </AdminEmptyState>
    </div>

    <template v-else-if="album">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-start">
        <div class="shrink-0">
          <div class="relative aspect-square w-48 overflow-hidden rounded-2xl bg-white/[0.06] shadow-lg lg:w-56">
            <img
              v-if="album.cover_url"
              :src="album.cover_url"
              :alt="album.title"
              class="h-full w-full object-cover"
              @error="($event.target as HTMLImageElement).style.display='none'"
            />
            <div v-else class="flex h-full w-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-image text-5xl text-slate-600" />
            </div>
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <h1 class="text-2xl font-bold text-white lg:text-3xl">{{ album.title }}</h1>

          <div class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-slate-400">
            <button
              v-if="artist"
              class="cursor-pointer text-teal-400 transition-colors hover:text-teal-300"
              @click="router.push({ name: 'admin.artist.detail', params: { id: artist.id } })"
            >
              {{ artist.name }}
            </button>
            <span v-if="album.release_date">
              <i aria-hidden="true" class="pi pi-calendar mr-1" />
              {{ album.release_date }}
            </span>
            <span v-if="album.track_count">
              <i aria-hidden="true" class="pi pi-music mr-1" />
              {{ album.track_count }} tracks
            </span>
          </div>

          <div class="mt-6 flex flex-wrap gap-2">
            <Button
              label="Fetch Cover"
              icon="pi pi-refresh"
              size="small"
              severity="info"
              :loading="enriching"
              class="!rounded-xl !bg-amber-500/10 !text-amber-400 hover:!bg-amber-500/20"
              @click="handleEnrich"
            />
            <Button
              label="Edit Album"
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
              class="!text-red-400"
              @click="openDeleteConfirm"
            />
          </div>
        </div>
      </div>

      <div class="mt-10">
        <h2 class="mb-4 text-lg font-semibold text-white">Tracks</h2>
        <div v-if="loadingTracks" class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-12 animate-pulse rounded-lg bg-white/[0.04]" />
        </div>
        <AdminEmptyState
          v-else-if="albumTracks.length === 0"
          icon="pi pi-music"
          title="No tracks"
          description="This album has no tracks yet."
          compact
        />
        <div v-else class="overflow-hidden rounded-xl border border-white/[0.06]">
          <div
            v-for="(track, i) in albumTracks"
            :key="track.id"
            class="flex cursor-pointer items-center gap-4 px-4 py-3 transition-colors hover:bg-white/[0.02]"
            :class="i < albumTracks.length - 1 ? 'border-b border-white/[0.04]' : ''"
            @click="router.push({ name: 'admin.track.detail', params: { id: track.id } })"
          >
            <span class="w-6 text-center text-xs text-slate-500">{{ track.track_number || i + 1 }}</span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
              <p v-if="track.artist_name" class="truncate text-xs text-slate-500">{{ track.artist_name }}</p>
            </div>
            <span v-if="track.duration_seconds" class="text-xs text-slate-500">{{ formatDuration(track.duration_seconds) }}</span>
            <i aria-hidden="true" class="pi pi-chevron-right text-xs text-slate-600" />
          </div>
        </div>
      </div>
    </template>

    <AlbumFormDialog
      v-model="showForm"
      :album="editingAlbum"
      :saving="saving"
      @submit="handleEditSubmit"
    />

    <AdminDeleteConfirm
      v-model="showDelete"
      title="Delete album"
      :item-name="album?.title"
      :deleting="deleting"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Button from 'primevue/button'
import { useToast } from 'primevue/usetoast'
import AdminEmptyState from '@/components/admin/AdminEmptyState.vue'
import AlbumFormDialog from '@/components/admin/AlbumFormDialog.vue'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'
import { useAlbumsApi } from '@/services/api/catalog/albums'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { useTracksApi } from '@/services/api/catalog/tracks'
import type { Album } from '@/services/api/catalog/albums'
import type { Artist } from '@/services/api/catalog/artists'
import type { Track } from '@/services/api/catalog/tracks'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const albumsApi = useAlbumsApi()
const artistsApi = useArtistsApi()
const tracksApi = useTracksApi()

const album = ref<Album | null>(null)
const artist = ref<Artist | null>(null)
const allTracks = ref<Track[]>([])
const loading = ref(true)
const error = ref('')
const loadingTracks = ref(false)

const showForm = ref(false)
const editingAlbum = ref<Album | null>(null)
const saving = ref(false)
const showDelete = ref(false)
const deleting = ref(false)
const enriching = ref(false)

const albumTracks = computed(() =>
  allTracks.value.filter(t => String(t.album_id) === route.params.id)
)

onMounted(loadAlbum)

async function loadAlbum() {
  const id = route.params.id as string
  loading.value = true
  error.value = ''
  try {
    album.value = await albumsApi.getAlbum(id)
    const [artistData, tracksData] = await Promise.all([
      album.value.artist_id ? artistsApi.getArtist(String(album.value.artist_id)) : Promise.resolve(null),
      tracksApi.getTracks(),
    ])
    artist.value = artistData
    allTracks.value = tracksData
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Failed to load album.'
  } finally {
    loading.value = false
  }
}

function openEdit() {
  editingAlbum.value = album.value
  showForm.value = true
}

function openDeleteConfirm() {
  showDelete.value = true
}

async function handleEnrich() {
  if (!album.value) return
  enriching.value = true
  try {
    await albumsApi.adminEnrichAlbum(album.value.id)
    toast.add({ severity: 'success', summary: 'Album enriched', detail: 'Cover art fetched from external sources', life: 3000 })
    await loadAlbum()
  } catch {
    toast.add({ severity: 'error', summary: 'Enrich failed', detail: 'Could not fetch album cover', life: 4000 })
  } finally {
    enriching.value = false
  }
}

async function handleEditSubmit(payload: Record<string, unknown>) {
  if (!album.value) return
  saving.value = true
  try {
    const updated = await albumsApi.adminUpdateAlbum(album.value.id, payload)
    album.value = updated
    showForm.value = false
    toast.add({ severity: 'success', summary: 'Album updated', life: 2500 })
  } catch {
    toast.add({ severity: 'error', summary: 'Update failed', life: 3000 })
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!album.value) return
  deleting.value = true
  try {
    await albumsApi.adminDeleteAlbum(album.value.id)
    toast.add({ severity: 'success', summary: 'Album deleted', life: 2500 })
    router.push({ name: 'admin.albums' })
  } catch {
    toast.add({ severity: 'error', summary: 'Delete failed', life: 3000 })
  } finally {
    deleting.value = false
  }
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>
