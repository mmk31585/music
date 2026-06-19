<template>
  <div class="mx-auto w-full max-w-5xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <div class="flex gap-6">
        <SkeletonLoader variant="card" class="w-56 shrink-0" />
        <div class="flex-1 space-y-3">
          <SkeletonLoader variant="lines" :lines="3" />
        </div>
      </div>
      <div class="space-y-2">
        <SkeletonLoader v-for="i in 6" :key="i" variant="track" />
      </div>
    </div>

    <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
      <i aria-hidden="true" class="pi pi-exclamation-circle text-4xl text-slate-500" />
      <h2 class="text-xl font-bold text-white">Playlist not found</h2>
      <RouterLink
        to="/playlists"
        class="text-sm font-medium text-[#1db954] underline underline-offset-2"
      >
        Back to playlists
      </RouterLink>
    </div>

    <template v-else-if="playlist">
      <div class="flex flex-col gap-6 md:flex-row md:items-end">
        <div
          class="relative h-56 w-56 shrink-0 overflow-hidden rounded-2xl bg-white/[0.06] shadow-2xl ring-1 ring-white/10"
        >
          <img
            v-if="playlist.cover_url"
            :src="playlist.cover_url"
            :alt="playlist.name"
            loading="lazy"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div
            v-else
            class="flex h-full items-center justify-center bg-gradient-to-br from-[#1db954]/20 to-[#121212]"
          >
            <i aria-hidden="true" class="pi pi-list text-4xl text-slate-500" />
          </div>
        </div>

        <div class="flex-1">
          <div class="flex items-center gap-2">
            <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">Playlist</p>
            <span
              v-if="playlist.is_collaborative"
              class="rounded-full bg-blue-500/10 px-2 py-0.5 text-xs font-medium text-blue-400"
            >
              <i aria-hidden="true" class="pi pi-users mr-1 text-[10px]" />Collaborative
            </span>
          </div>
          <h1 class="mt-2 text-3xl font-black text-white md:text-5xl">{{ playlist.name }}</h1>

          <div class="mt-3 flex flex-wrap items-center gap-2 text-sm text-slate-400">
            <span>{{ tracks.length }} {{ tracks.length === 1 ? 'track' : 'tracks' }}</span>
            <span> • Updated {{ timeAgo }}</span>
            <span v-if="collaborators.length" class="text-slate-500">
              • {{ collaborators.length }}
              {{ collaborators.length === 1 ? 'collaborator' : 'collaborators' }}
            </span>
          </div>

          <p
            v-if="playlist.description"
            class="mt-2 max-w-lg text-sm leading-relaxed text-slate-400"
          >
            {{ playlist.description }}
          </p>

          <div class="mt-6 flex flex-wrap items-center gap-3">
            <button
              type="button"
              :disabled="!tracks.length"
              class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760] disabled:opacity-40 disabled:hover:scale-100"
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
                  : 'border-white/15 bg-white/10 text-white hover:bg-white/15'
              "
              @click="toggleCollaborative"
            >
              <i aria-hidden="true" class="pi pi-users text-xs" />
              {{ playlist.is_collaborative ? 'Collaborative' : 'Make Collaborative' }}
            </button>
          </div>

          <!-- Collaborators (owner only) -->
          <div v-if="isOwner && collaborators.length" class="mt-4 flex flex-wrap gap-2">
            <div
              v-for="c in collaborators"
              :key="c.user_id"
              class="flex items-center gap-2 rounded-full bg-white/5 px-3 py-1.5 text-xs text-slate-300"
            >
              <i aria-hidden="true" class="pi pi-user text-[10px]" />
              <span>{{ c.is_creator ? 'You' : `User #${c.user_id}` }}</span>
              <button
                v-if="!c.is_creator"
                type="button"
                class="text-slate-500 transition hover:text-red-400"
                @click="removeCollab(c.user_id)"
              >
                <i aria-hidden="true" class="pi pi-times text-[10px]" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <section class="mt-10">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-bold text-white">Tracks</h2>
          <div v-if="tracks.length" class="flex items-center gap-2">
            <button
              v-if="isOwner || isCollaborator"
              type="button"
              class="rounded-full border border-red-500/30 px-4 py-1.5 text-xs font-medium text-red-400 transition hover:bg-red-500/10"
              @click="deletePlaylist"
            >
              Delete playlist
            </button>
          </div>
        </div>

        <div
          v-if="tracks.length"
          class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]"
        >
          <div
            v-for="(item, index) in tracks"
            :key="item.playlist_track_id"
            :draggable="isOwner || isCollaborator"
            class="group flex items-center gap-3 px-4 py-2 transition"
            :class="{
              'hover:bg-white/[0.06]': true,
              'opacity-50': dragIndex === index,
              'border-t-2 border-[#1db954]': dropTargetIndex === index,
            }"
            @dragstart="onDragStart(index)"
            @dragover="onDragOver(index)"
            @dragleave="onDragLeave(index)"
            @drop="onDrop(index)"
            @dragend="onDragEnd"
          >
            <!-- Drag handle for owners/collaborators -->
            <span
              v-if="isOwner || isCollaborator"
              class="flex w-6 cursor-grab items-center justify-center text-slate-500 active:cursor-grabbing"
            >
              <i aria-hidden="true" class="pi pi-bars text-xs opacity-0 transition group-hover:opacity-100" />
            </span>
            <span class="w-6 text-right text-xs text-slate-500">{{ index + 1 }}</span>

            <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
              <img
                v-if="item.cover_url"
                :src="item.cover_url"
                :alt="item.title"
                loading="lazy"
                class="h-full w-full object-cover"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-xs text-slate-500" />
              </div>
              <button
                type="button"
                class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                @click="playTrack(index)"
              >
                <i aria-hidden="true" class="pi pi-play-fill text-xs text-white" />
              </button>
            </div>

            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">{{ item.title }}</p>
              <p class="truncate text-xs text-slate-400">{{ item.artist_name }}</p>
            </div>

            <span class="text-xs text-slate-500">{{ formatTime(item.duration_seconds) }}</span>

            <button
              v-if="isOwner || isCollaborator"
              type="button"
              class="shrink-0 rounded-full p-1.5 text-slate-500 opacity-0 transition group-hover:opacity-100 hover:bg-white/10 hover:text-red-400"
              title="Remove from playlist"
              @click="removeTrack(item.track_id)"
            >
              <i aria-hidden="true" class="pi pi-times text-xs" />
            </button>
          </div>
        </div>

        <div
          v-else
          class="flex flex-col items-center gap-3 rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-16 text-center"
        >
          <i aria-hidden="true" class="pi pi-list text-3xl text-slate-500" />
          <h3 class="text-lg font-bold text-white">Empty playlist</h3>
          <p class="text-sm text-slate-400">Add tracks from the search or your library</p>
          <RouterLink
            to="/search"
            class="rounded-full bg-[#1db954] px-6 py-2 text-sm font-bold text-black transition hover:bg-[#1ed760]"
          >
            Search tracks
          </RouterLink>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useUserAuthStore } from '@/stores'
import { usePlaylistDetail } from '@/composables/catalog/usePlaylistDetail'
import { usePlaylistsApi } from '@/services/api/playlist'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'
import { useCollaborativePlaylist } from '@/composables/useCollaborativePlaylist'
import { wsClient } from '@/services/socket'
import { useToast } from 'primevue/usetoast'

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
const playerApi = usePlayerApi()
const auth = useUserAuthStore()
const playlistsApi = usePlaylistsApi()

const collaborators = ref<any[]>([])
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
  if (dropTargetIndex.value === index) {
    dropTargetIndex.value = -1
  }
}

async function onDrop(targetIndex: number) {
  if (dragIndex.value < 0 || dragIndex.value === targetIndex) {
    onDragEnd()
    return
  }

  const reordered = [...tracks.value]
  const [moved] = reordered.splice(dragIndex.value, 1)
  reordered.splice(targetIndex, 0, moved!)

  const newOrder = reordered.map((t) => t.track_id)

  // Optimistic update
  tracks.value = reordered

  try {
    await playlistsApi.reorderTracks(playlistId, newOrder)
    toast.add({ severity: 'success', summary: 'Tracks reordered', life: 2000 })
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
  () => {
    fetchPlaylist()
  },
  () => {
    fetchPlaylist()
  },
  () => {
    fetchPlaylist()
  },
  () => {
    fetchPlaylist()
  },
)

async function toggleCollaborative() {
  if (!playlist.value) return
  try {
    const newVal = !playlist.value.is_collaborative
    await playlistsApi.setCollaborative(playlistId, newVal)
    playlist.value.is_collaborative = newVal
    toast.add({
      severity: 'success',
      summary: newVal ? 'Collaborative mode on' : 'Collaborative mode off',
      life: 2000,
    })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to update', life: 3000 })
  }
}

async function removeCollab(userId: string) {
  try {
    await playlistsApi.removeCollaborator(playlistId, userId)
    collaborators.value = collaborators.value.filter((c) => c.user_id !== userId)
    toast.add({ severity: 'success', summary: 'Collaborator removed', life: 2000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to remove', life: 3000 })
  }
}

onMounted(async () => {
  await fetchPlaylist()
  wsClient.connect()
  collabHelper.setup()
  const collabData = await playlistsApi.listCollaborators(playlistId).catch(() => null)
  if (collabData) {
    collaborators.value = collabData
    isCollaborator.value = collabData.some(
      (c: any) => (c as { user_id: string; is_creator: boolean }).user_id === String(auth.user?.id) && !(c as { user_id: string; is_creator: boolean }).is_creator,
    )
  }
})

onUnmounted(() => {
  collabHelper.teardown()
})

function playAll() {
  if (tracks.value.length) playTrack(0)
}

function playTrack(index: number) {
  const queue = tracks.value.map((t) => ({
    id: String(t.track_id),
    title: t.title,
    artistName: t.artist_name || 'Unknown',
    albumTitle: t.album_title || null,
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.track_id)),
  }))
  player.setQueueAndPlay(queue, index)
}

async function deletePlaylist() {
  const deleted = await remove()
  if (deleted) {
    toast.add({ severity: 'success', summary: 'Playlist deleted', life: 2000 })
    router.push('/playlists')
  }
}

function formatTime(seconds?: number | null) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
</script>
