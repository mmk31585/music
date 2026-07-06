<template>
  <div class="relative mx-auto min-h-screen w-full pb-36">
    <!-- ── Ambient Aurora ── -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1" />
      <div class="aurora-spot-2" />
    </div>

    <div class="relative z-10 px-4 pt-8 md:px-6 lg:px-8">
      <div class="mx-auto max-w-6xl">
        <!-- Loading -->
        <div v-if="loading" class="space-y-6">
          <div class="flex gap-6">
            <div class="shimmer h-56 w-56 shrink-0 rounded-2xl" />
            <div class="flex-1 space-y-3">
              <div class="shimmer h-4 w-24 rounded-lg" />
              <div class="shimmer h-8 w-64 rounded-lg" />
              <div class="shimmer h-4 w-48 rounded-lg" />
            </div>
          </div>
          <div class="space-y-2">
            <div v-for="i in 5" :key="i" class="shimmer h-15 rounded-2xl" />
          </div>
        </div>

        <!-- Error -->
        <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
          <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/4">
            <i aria-hidden="true" class="pi pi-exclamation-circle text-3xl text-white/20" />
          </div>
          <h2 class="text-xl font-bold text-white">Playlist not found</h2>
          <RouterLink to="/library" class="text-sm font-medium text-spotify hover:underline">
            Back to library
          </RouterLink>
        </div>

        <!-- ── Content ── -->
        <template v-else-if="playlist">
          <div class="flex flex-col gap-8 md:flex-row md:items-end">
            <!-- Cover -->
          <div class="group relative h-56 w-56 shrink-0 overflow-hidden rounded-2xl shadow-2xl ring-1 ring-white/6">
            <!-- Custom cover image -->
            <img
              v-if="playlist.cover_url"
              :src="playlist.cover_url"
              :alt="playlist.name"
              loading="lazy"
              class="h-full w-full object-cover"
              @error="onImgError"
            />
            <!-- Default: 4-track grid collage -->
            <PlaylistCoverGrid
              v-else
              :covers="trackCoverUrls"
              :track-count="tracks.length"
              class="h-full w-full"
            />

            <!-- Edit cover overlay (owner only) -->
            <div
              v-if="isOwner"
              class="absolute inset-0 flex cursor-pointer items-center justify-center bg-black/50 opacity-0 backdrop-blur-xs transition-opacity duration-300 group-hover:opacity-100"
              @click="triggerCoverUpload"
            >
              <div class="flex flex-col items-center gap-1.5 text-white">
                <div class="flex h-10 w-10 items-center justify-center rounded-full bg-white/20 backdrop-blur-xs">
                  <i aria-hidden="true" class="pi pi-camera text-lg" />
                </div>
                <span class="text-xs font-bold">{{ playlist.cover_url ? 'Change cover' : 'Add cover' }}</span>
              </div>
            </div>

            <!-- Hidden file input -->
            <input
              ref="coverFileInput"
              type="file"
              accept="image/*"
              class="hidden"
              @change="onCoverFileSelected"
            />
          </div>

            <!-- Info -->
            <div class="flex-1">
              <div class="flex items-center gap-2">
                <p class="text-[10px] font-bold tracking-[0.25em] text-spotify uppercase">Playlist</p>
                <span
                  v-if="playlist.is_collaborative"
                  class="rounded-full bg-blue-500/10 px-2 py-0.5 text-[10px] font-medium text-blue-400"
                >
                  <i aria-hidden="true" class="pi pi-users mr-1 text-[8px]" />Collaborative
                </span>
              </div>
              <h1 class="mt-2 text-3xl font-black text-white md:text-5xl">{{ playlist.name }}</h1>

              <div class="mt-3 flex flex-wrap items-center gap-2 text-sm text-white/40">
                <span>{{ tracks.length }} {{ tracks.length === 1 ? 'track' : 'tracks' }}</span>
                <span>• Updated {{ timeAgo }}</span>
                <span v-if="collaborators.length" class="text-white/30">
                  • {{ collaborators.length }}
                  {{ collaborators.length === 1 ? 'collaborator' : 'collaborators' }}
                </span>
              </div>

              <p v-if="playlist.description" class="mt-2 max-w-lg text-sm leading-relaxed text-white/40">
                {{ playlist.description }}
              </p>

              <div class="mt-6 flex flex-wrap items-center gap-3">
                <button
                  type="button"
                  :disabled="!tracks.length"
                  class="inline-flex items-center gap-2 rounded-full bg-spotify px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-spotify-hover disabled:opacity-40 disabled:hover:scale-100"
                  @click="playAll"
                >
                  <i aria-hidden="true" class="pi pi-play-fill" /> Play
                </button>

                <button
                  v-if="isOwner"
                  type="button"
                  class="inline-flex items-center gap-2 rounded-full border px-6 py-3 text-sm font-bold transition"
                  :class="
                    playlist.is_collaborative
                      ? 'border-blue-500/50 bg-blue-500/10 text-blue-400'
                      : 'border-white/6 bg-white/4 text-white/60 hover:bg-white/8 hover:text-white'
                  "
                  @click="toggleCollaborative"
                >
                  <i aria-hidden="true" class="pi pi-users text-xs" />
                  {{ playlist.is_collaborative ? 'Collaborative' : 'Make Collaborative' }}
                </button>

                <button
                  v-if="isOwner"
                  type="button"
                  class="inline-flex items-center gap-2 rounded-full border border-red-500/20 px-6 py-3 text-sm font-bold text-red-400 transition hover:bg-red-500/10"
                  @click="deletePlaylist"
                >
                  <i aria-hidden="true" class="pi pi-trash text-xs" />
                  Delete
                </button>
              </div>

              <!-- Collaborators -->
              <div v-if="isOwner && collaborators.length" class="mt-4 flex flex-wrap gap-2">
                <div
                  v-for="c in collaborators"
                  :key="String(c.user_id)"
                  class="flex items-center gap-2 rounded-full bg-white/4 px-3 py-1.5 text-xs text-white/50"
                >
                  <i aria-hidden="true" class="pi pi-user text-[10px]" />
                  <span>{{ c.is_creator ? 'You' : `User #${String(c.user_id)}` }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- ── Tracks Section ── -->
          <section class="mt-10">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold text-white">Tracks</h2>
              <div v-if="tracks.length" class="flex items-center gap-2">
                <button
                  v-if="isOwner || isCollaborator"
                  type="button"
                  class="inline-flex items-center gap-1.5 rounded-full border border-white/6 bg-white/4 px-4 py-1.5 text-xs font-medium text-white/60 transition hover:bg-white/8 hover:text-white"
                  @click="showAddTrack = true"
                >
                  <i aria-hidden="true" class="pi pi-plus text-[10px]" />
                  Add Track
                </button>
              </div>
            </div>

            <div
              v-if="tracks.length"
              class="overflow-hidden rounded-2xl border border-white/6 bg-white/2 backdrop-blur-xs"
              aria-live="polite"
            >
              <div
                v-for="(item, index) in tracks"
                :key="item.playlist_track_id"
                :draggable="isOwner || isCollaborator"
                class="group flex items-center gap-3 px-4 py-2.5 transition"
                :class="{
                  'hover:bg-white/4': true,
                  'opacity-50': dragIndex === index,
                  'border-t-2 border-spotify': dropTargetIndex === index,
                }"
                @dragstart="onDragStart(index)"
                @dragover="onDragOver(index)"
                @dragleave="onDragLeave(index)"
                @drop="onDrop(index)"
                @dragend="onDragEnd"
              >
                <!-- Drag handle -->
                <span
                  v-if="isOwner || isCollaborator"
                  class="flex w-6 cursor-grab items-center justify-center text-white/20 active:cursor-grabbing"
                >
                  <i aria-hidden="true" class="pi pi-bars text-xs opacity-0 transition group-hover:opacity-100" />
                </span>

                <!-- Number / Play -->
                <span class="flex w-6 items-center justify-center">
                  <span class="text-xs font-bold text-white/20 group-hover:hidden">{{ index + 1 }}</span>
                  <i aria-hidden="true" class="pi pi-play-fill hidden text-xs text-white group-hover:block" />
                </span>

                <!-- Cover -->
                <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/5 ring-1 ring-white/6">
                  <img
                    v-if="item.cover_url"
                    :src="item.cover_url"
                    :alt="item.title"
                    loading="lazy"
                    class="h-full w-full object-cover"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i aria-hidden="true" class="pi pi-music text-xs text-white/20" />
                  </div>
                  <button
                    type="button"
                    aria-label="Play track"
                    class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
                    @click="playTrack(index)"
                  >
                    <i aria-hidden="true" class="pi pi-play-fill text-xs text-white" />
                  </button>
                </div>

                <!-- Info -->
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-white">{{ item.title }}</p>
                  <p class="truncate text-xs text-white/40">{{ item.artist_name }}</p>
                </div>

                <!-- Duration -->
                <span class="shrink-0 text-xs text-white/30 tabular-nums">
                  {{ formatTime(item.duration_seconds) }}
                </span>

                <!-- Remove -->
                <button
                  v-if="isOwner || isCollaborator"
                  type="button"
                  aria-label="Remove from playlist"
                  class="shrink-0 rounded-full p-1.5 text-white/20 opacity-0 transition group-hover:opacity-100 hover:bg-white/6 hover:text-red-400"
                  @click="removeTrack(item.track_id)"
                >
                  <i aria-hidden="true" class="pi pi-times text-xs" />
                </button>
              </div>
            </div>

            <div
              v-else
              class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/6 bg-white/2 px-6 py-16 text-center"
            >
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/4">
                <i aria-hidden="true" class="pi pi-list text-xl text-white/20" />
              </div>
              <h3 class="text-base font-bold text-white">Empty playlist</h3>
              <p class="text-sm text-white/40">Add tracks from search or your library</p>
              <button
                v-if="isOwner || isCollaborator"
                type="button"
                class="rounded-full bg-spotify px-6 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover"
                @click="showAddTrack = true"
              >
                <i aria-hidden="true" class="pi pi-plus mr-1 text-xs" />
                Add Track
              </button>
            </div>
          </section>
        </template>
      </div>
    </div>

    <!-- ════════════════════════════════════════ -->
    <!-- ADD TRACK DIALOG                         -->
    <!-- ════════════════════════════════════════ -->
    <Dialog
      :visible="showAddTrack"
      modal
      :draggable="false"
      :style="{ maxWidth: '500px', width: '90vw' }"
      :pt="{
        root: 'border-none',
        mask: 'backdrop-blur-xs bg-black/60',
        header: 'border-b border-white/5',
        title: 'text-white text-sm font-bold',
        content: 'p-0',
      }"
      @update:visible="showAddTrack = $event"
    >
      <template #header>
        <div class="flex items-center gap-2 px-1">
          <i aria-hidden="true" class="pi pi-search text-sm text-spotify" />
          <span>Add Track</span>
        </div>
      </template>

      <div class="flex flex-col max-h-[70vh]">
        <div class="shrink-0 p-3">
          <div class="relative">
            <i aria-hidden="true" class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-white/20" />
            <input
              v-model="addTrackQuery"
              type="text"
              placeholder="Search tracks..."
              aria-label="Search tracks to add"
              class="w-full rounded-xl border border-white/6 bg-white/3 py-2.5 pl-9 pr-3 text-sm text-white outline-hidden transition placeholder:text-white/20 focus:border-spotify/30 focus:bg-white/6"
              @input="onAddTrackSearch"
            />
          </div>
        </div>

        <div class="flex-1 overflow-y-auto px-1">
          <div v-if="addTrackLoading" class="flex items-center justify-center py-12">
            <i aria-hidden="true" class="pi pi-spin pi-spinner text-lg text-white/20" />
          </div>

          <div
            v-else-if="addTrackResults.length === 0 && addTrackQuery"
            class="flex flex-col items-center gap-2 py-12 text-center"
          >
            <i aria-hidden="true" class="pi pi-search text-2xl text-white/20" />
            <p class="text-sm text-white/40">No tracks found</p>
          </div>

          <div v-else-if="!addTrackQuery" class="flex flex-col items-center gap-2 py-12 text-center">
            <i aria-hidden="true" class="pi pi-music text-2xl text-white/20" />
            <p class="text-sm text-white/40">Type to search tracks</p>
          </div>

          <button
            v-for="t in addTrackResults"
            :key="String(t.id)"
            type="button"
            :disabled="addingTrackId === String(t.id)"
            class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left transition hover:bg-white/4 disabled:opacity-50"
            @click="addSelectedTrack(String(t.id))"
          >
            <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/5 ring-1 ring-white/6">
              <img
                v-if="t.cover_url"
                :src="t.cover_url as string | undefined"
                :alt="String(t.title ?? '')"
                class="h-full w-full object-cover"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-xs text-white/20" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">{{ t.title }}</p>
              <p class="truncate text-xs text-white/40">{{ t.artist_name }}</p>
            </div>
            <span class="text-xs text-white/30">{{ formatDuration(t.duration_seconds as number) }}</span>
            <i v-if="addingTrackId === t.id" aria-hidden="true" class="pi pi-spin pi-spinner text-xs text-spotify" />
            <i v-else aria-hidden="true" class="pi pi-plus text-xs text-white/30" />
          </button>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useUserAuthStore } from '@/stores'
import { PlaylistCoverGrid } from '@/components/music'
import { usePlaylistDetail } from '@/composables/catalog/usePlaylistDetail'
import { usePlaylistsApi } from '@/services/api/playlist'
import { useMediaApi } from '@/services/api/media'
import { usePlayer } from '@/composables/player'
import { useSearchApi } from '@/services/api/catalog/search'
import { mapToPlaybackTracks } from '@/factories/playbackTrack'
import { onImgError } from '@/utils/helpers'
import { useCollaborativePlaylist } from '@/composables/useCollaborativePlaylist'
import { wsClient } from '@/services/socket'
import { formatDuration } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const playlistId = String(route.params.id)

const {
  playlist,
  tracks,
  loading,
  error,
  fetchPlaylist,
  removeTrack,
  deletePlaylist: remove,
} = usePlaylistDetail(playlistId)
const player = usePlayer()
const auth = useUserAuthStore()
const playlistsApi = usePlaylistsApi()

const collaborators = ref<Record<string, unknown>[]>([])
const isCollaborator = ref(false)
const dragIndex = ref(-1)
const dropTargetIndex = ref(-1)

function onDragStart(index: number) {
  dragIndex.value = index
}

function onDragOver(index: number) {
  if (dragIndex.value === index) return
  dropTargetIndex.value = index
}

function onDragLeave(index: number) {
  if (dropTargetIndex.value === index) dropTargetIndex.value = -1
}

async function onDrop(targetIndex: number) {
  if (dragIndex.value < 0 || dragIndex.value === targetIndex) {
    onDragEnd()
    return
  }

  const reordered = [...tracks.value]
  const [moved] = reordered.splice(dragIndex.value, 1)
  reordered.splice(targetIndex, 0, moved!)

  tracks.value = reordered

  try {
    // Backend expects { track_id, new_position } (1-indexed position)
    await playlistsApi.reorderTracks(playlistId, moved!.track_id, targetIndex + 1)
  } catch {
    await fetchPlaylist()
    toast.add({ severity: 'error', summary: 'Failed to reorder', life: 3000 })
  }
  onDragEnd()
}

function onDragEnd() {
  dragIndex.value = -1
  dropTargetIndex.value = -1
}

const mediaApi = useMediaApi()
const coverFileInput = ref<HTMLInputElement | null>(null)
const uploadingCover = ref(false)

const trackCoverUrls = computed(() => {
  return tracks.value.map((t) => t.cover_url)
})

function triggerCoverUpload() {
  coverFileInput.value?.click()
}

async function onCoverFileSelected(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || !playlist.value) return

  uploadingCover.value = true
  try {
    const uploaded = await mediaApi.uploadAdminMedia('playlistCover', file, undefined, {
      onUploadProgress(progressEvent) {
        // Could add progress bar here
      },
    })
    const url = uploaded.url || uploaded.file_url || uploaded.fileUrl || uploaded.path
    if (!url) throw new Error('No URL returned from upload')

    await playlistsApi.updatePlaylist(playlistId, { cover_url: url })
    playlist.value.cover_url = url
    toast.add({ severity: 'success', summary: 'Cover updated', life: 2000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to upload cover', life: 3000 })
  } finally {
    uploadingCover.value = false
    input.value = '' // Reset file input
  }
}

const isOwner = computed(() => {
  if (!playlist.value) return false
  return playlist.value.user_id === String(auth.user?.id)
})

const timeAgo = computed(() => {
  if (!playlist.value?.updated_at) return ''
  const diff = Date.now() - new Date(playlist.value.updated_at).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
})

const collabHelper = useCollaborativePlaylist(
  playlistId,
  () => { fetchPlaylist() },
  () => { fetchPlaylist() },
  () => { fetchPlaylist() },
  () => { fetchPlaylist() },
)

async function toggleCollaborative() {
  if (!playlist.value) return
  try {
    const newVal = !playlist.value.is_collaborative
    await playlistsApi.setCollaborative(playlistId, newVal)
    playlist.value.is_collaborative = newVal
    toast.add({ severity: 'success', summary: newVal ? 'Collaborative mode on' : 'Collaborative mode off', life: 2000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to update', life: 3000 })
  }
}

// TODO: Backend route POST/DELETE /playlists/:id/collaborators not yet registered.
// Re-enable when backend exposes these routes.
// async function removeCollab(userId: string) {
//   try {
//     await playlistsApi.removeCollaborator(playlistId, userId)
//     collaborators.value = collaborators.value.filter((c) => String(c.user_id) !== userId)
//     toast.add({ severity: 'success', summary: 'Collaborator removed', life: 2000 })
//   } catch {
//     toast.add({ severity: 'error', summary: 'Failed to remove', life: 3000 })
//   }
// }

onMounted(async () => {
  await fetchPlaylist()
  wsClient.connect()
  collabHelper.setup()
  const collabData = await playlistsApi.listCollaborators(playlistId).catch(() => null)
  if (collabData) {
    collaborators.value = Array.isArray(collabData) ? collabData : []
    isCollaborator.value = collaborators.value.some(
      (c: Record<string, unknown>) => String(c.user_id) === String(auth.user?.id) &&       c.is_creator,
    )
  }
})

onUnmounted(() => {
  collabHelper.teardown()
})

// ── Add Track ──
const showAddTrack = ref(false)
const addTrackQuery = ref('')
const addTrackLoading = ref(false)
const addTrackResults = ref<Record<string, unknown>[]>([])
const addingTrackId = ref<string | null>(null)
let addTrackDebounce: ReturnType<typeof setTimeout> | null = null
const searchApi = useSearchApi()

function onAddTrackSearch() {
  if (addTrackDebounce) clearTimeout(addTrackDebounce)
  const q = addTrackQuery.value.trim()
  if (!q) {
    addTrackResults.value = []
    return
  }
  addTrackDebounce = setTimeout(async () => {
    addTrackLoading.value = true
    try {
      const res = await searchApi.searchCatalog({ query: q, type: 'tracks', limit: 10 })
      addTrackResults.value = res?.tracks || []
    } catch {
      addTrackResults.value = []
    } finally {
      addTrackLoading.value = false
    }
  }, 300)
}

async function addSelectedTrack(trackId: string) {
  addingTrackId.value = trackId
  try {
    await playlistsApi.addTrack(playlistId, { track_id: trackId })
    toast.add({ severity: 'success', summary: 'Track added to playlist', life: 2500 })
    showAddTrack.value = false
    addTrackQuery.value = ''
    addTrackResults.value = []
    await fetchPlaylist()
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to add track', life: 3000 })
  } finally {
    addingTrackId.value = null
  }
}

function playAll() {
  if (tracks.value.length) playTrack(0)
}

function playTrack(index: number) {
  const queue = mapToPlaybackTracks(tracks.value)
  player.setQueueAndPlay(queue, index)
}

async function deletePlaylist() {
  const deleted = await remove()
  if (deleted) {
    toast.add({ severity: 'success', summary: 'Playlist deleted', life: 2000 })
    router.push('/library')
  }
}

function formatTime(seconds?: number | null) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>
