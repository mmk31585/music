<template>
  <Dialog
    :visible="visible"
    modal
    :draggable="false"
    aria-labelledby="add-to-playlist-title"
    :style="{ maxWidth: '420px', width: '90vw' }"
    :pt="{
      root: 'border-none',
      mask: 'backdrop-blur-xs',
      header: 'border-b border-white/5',
      title: 'text-white text-sm font-bold',
      content: 'p-0',
    }"
    @update:visible="emit('update:visible', $event)"
  >
    <template #header>
      <div class="flex items-center gap-2 px-1">
        <List aria-hidden="true" class="text-sm text-spotify"  />
        <span id="add-to-playlist-title">Add to Playlist</span>
      </div>
    </template>

    <div class="flex flex-col max-h-[60vh]">
      <div class="shrink-0 p-3">
        <div class="relative">
          <Search
            aria-hidden="true"
            class="absolute left-3 top-1/2 -translate-y-1/2 text-xs text-white/30"
           />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Filter playlists..."
            aria-label="Filter playlists"
            class="w-full rounded-xl border border-white/10 bg-white/5 py-2.5 pl-9 pr-3 text-sm text-white outline-hidden transition placeholder:text-white/30 focus:border-spotify/50"
          />
        </div>
      </div>

      <div class="flex-1 overflow-y-auto px-1">
        <div v-if="loading" class="flex items-center justify-center py-12">
          <Loader2 aria-hidden="true" class="text-lg text-white/30 animate-spin"  />
        </div>

        <div v-else-if="filteredPlaylists.length === 0" class="flex flex-col items-center gap-2 py-12 text-center">
          <Inbox aria-hidden="true" class="text-2xl text-white/20"  />
          <p class="text-sm text-white/40">No playlists found</p>
        </div>

        <button
          v-for="p in filteredPlaylists"
          :key="p.id"
          type="button"
          class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left transition hover:bg-white/5"
          @click="addTo(p.id)"
        >
          <div class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-white/10">
            <img
              v-if="p.cover_url"
              :src="p.cover_url"
              :alt="p.name"
              class="h-full w-full object-cover"
            />
            <List v-else aria-hidden="true" class="text-sm text-white/30"  />
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium text-white">{{ p.name }}</p>
            <p class="text-xs text-white/40">{{ p.track_count || 0 }} tracks</p>
          </div>
          <Loader2
            v-if="addingId === p.id"
            aria-hidden="true"
            class="text-xs text-spotify animate-spin"
           />
          <Plus
            v-else
            aria-hidden="true"
            class="text-xs text-white/30 transition group-hover:text-white"
           />
        </button>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { Inbox, List, Loader2, Plus, Search } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { usePlaylistsApi, type PlaylistListItem } from '@/services/api/playlist'
import { useToast } from 'primevue/usetoast'

const props = defineProps<{
  visible: boolean
  trackId: string
  trackTitle?: string
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'added': [playlistId: string]
}>()

const toast = useToast()
const playlistsApi = usePlaylistsApi()

const loading = ref(false)
const searchQuery = ref('')
const addingId = ref<string | null>(null)
const playlists = ref<PlaylistListItem[]>([])

const filteredPlaylists = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return playlists.value
  return playlists.value.filter(
    (p) => p.name.toLowerCase().includes(q),
  )
})

async function fetch() {
  loading.value = true
  try {
    const data = await playlistsApi.getMyPlaylists()
    playlists.value = data || []
  } catch {
    playlists.value = []
  } finally {
    loading.value = false
  }
}

async function addTo(playlistId: string) {
  addingId.value = playlistId
  try {
    // Duplicate check: fetch playlist detail to see if track already exists
    const detail = await playlistsApi.getPlaylist(playlistId)
    const alreadyIn = (detail.tracks || []).some((t) => t.track_id === props.trackId)
    if (alreadyIn) {
      toast.add({
        severity: 'info',
        summary: 'Already in this playlist',
        life: 2500,
      })
      addingId.value = null
      return
    }

    await playlistsApi.addTrack(playlistId, { track_id: props.trackId })
    toast.add({
      severity: 'success',
      summary: props.trackTitle
        ? `"${props.trackTitle}" added to playlist`
        : 'Track added to playlist',
      life: 2500,
    })
    emit('added', playlistId)
    emit('update:visible', false)
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Failed to add track',
      life: 3000,
    })
  } finally {
    addingId.value = null
  }
}

watch(() => props.visible, (v) => {
  if (v) fetch()
})

onMounted(() => {
  if (props.visible) fetch()
})
</script>
