<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="relative overflow-hidden rounded-[2rem] border border-white/[0.06] bg-[#0C0C14] p-10 text-white"
    >
      <div class="absolute -top-20 -right-20 h-60 w-60 rounded-full bg-[#1db954]/10 blur-3xl" />
      <div class="relative">
        <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Collections</p>
        <h1 class="mt-2 text-4xl font-black md:text-6xl">Playlists</h1>
        <p class="mt-3 max-w-2xl text-sm text-white/50">
          Create, organize, and play custom music collections.
        </p>
      </div>
    </section>

    <section class="mt-10">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-bold text-white">My Playlists</h2>
        <button
          type="button"
          class="rounded-full bg-[#1db954] px-6 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
          @click="showCreate = true"
        >
          <i aria-hidden="true" class="pi pi-plus mr-2" />
          Create Playlist
        </button>
      </div>

      <div v-if="loading" class="mt-6 space-y-3">
        <div v-for="i in 4" :key="i" class="h-20 animate-pulse rounded-2xl bg-white/[0.06]" />
      </div>

      <div
        v-else-if="playlists.length === 0"
        class="mt-6 rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center"
      >
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <i aria-hidden="true" class="pi pi-list" />
        </div>
        <h2 class="mt-5 text-xl font-black text-white">No playlists yet</h2>
        <p class="mt-2 text-sm text-slate-400">Start by creating your first playlist.</p>
        <button
          type="button"
          class="mt-5 rounded-full bg-[#1db954] px-6 py-3 text-sm font-bold text-black transition hover:bg-[#1ed760]"
          @click="showCreate = true"
        >
          <i aria-hidden="true" class="pi pi-plus mr-2" />
          Create Playlist
        </button>
      </div>

      <div v-else class="mt-6 grid gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4" aria-live="polite">
        <RouterLink
          v-for="playlist in playlists"
          :key="playlist.id"
          :to="`/playlist/${playlist.id}`"
          class="group rounded-2xl border border-white/10 bg-white/[0.04] p-4 transition hover:-translate-y-0.5 hover:bg-white/[0.08]"
        >
          <div class="relative mb-3 aspect-square overflow-hidden rounded-xl bg-white/10">
            <img
              v-if="playlist.cover_url"
              :src="playlist.cover_url"
              :alt="playlist.name"
              loading="lazy"
              class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
            />
            <div v-else class="flex h-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-list text-3xl text-slate-500" />
            </div>
            <div class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100">
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl">
                <i aria-hidden="true" class="pi pi-play-fill text-lg" />
              </div>
            </div>
          </div>
          <p class="truncate font-semibold text-white">{{ playlist.name }}</p>
          <p class="text-xs text-slate-400">{{ playlist.track_count }} tracks</p>
        </RouterLink>
      </div>
    </section>

    <Dialog
      v-model:visible="showCreate"
      header="Create Playlist"
      :modal="true"
      class="w-full max-w-md rounded-2xl bg-[#121212]"
    >
      <div class="space-y-4 p-4">
        <div>
          <label class="mb-1 block text-sm font-semibold text-slate-300">Name</label>
          <InputText
            v-model="newName"
            placeholder="My awesome playlist"
            class="w-full rounded-xl border border-white/10 bg-white/[0.05] px-4 py-3 text-white placeholder:text-slate-500"
          />
        </div>
        <div>
          <label class="mb-1 block text-sm font-semibold text-slate-300">Description</label>
          <Textarea
            v-model="newDescription"
            placeholder="Optional description"
            rows="3"
            class="w-full rounded-xl border border-white/10 bg-white/[0.05] px-4 py-3 text-white placeholder:text-slate-500"
          />
        </div>
        <div class="flex items-center gap-2">
          <Checkbox v-model="newIsPublic" :binary="true" input-id="public" />
          <label for="public" class="text-sm text-slate-300">Public playlist</label>
        </div>
        <div class="flex justify-end gap-3 pt-2">
          <Button
            label="Cancel"
            class="rounded-xl bg-white/10 px-5 py-2 text-sm font-semibold text-white hover:bg-white/20"
            @click="showCreate = false"
          />
          <Button
            label="Create"
            class="rounded-xl bg-[#1db954] px-5 py-2 text-sm font-semibold text-black hover:bg-[#1ed760]"
            :disabled="!newName.trim() || creating"
            @click="handleCreate"
          />
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useRouter } from 'vue-router'
import { usePlaylistsApi } from '@/services/api/playlist'
import type { PlaylistListItem } from '@/services/api/playlist'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'

const playlistsApi = usePlaylistsApi()
const router = useRouter()

const playlists = ref<PlaylistListItem[]>([])
const loading = ref(false)
const showCreate = ref(false)
const creating = ref(false)
const newName = ref('')
const newDescription = ref('')
const newIsPublic = ref(true)

async function fetchPlaylists() {
  loading.value = true
  try {
    const data = await playlistsApi.getMyPlaylists()
    playlists.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Failed to fetch playlists:', err)
    playlists.value = []
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!newName.value.trim() || creating.value) return
  creating.value = true
  try {
    const result = await playlistsApi.createPlaylist({
      name: newName.value.trim(),
      description: newDescription.value.trim() || undefined,
      is_public: newIsPublic.value,
    })
    showCreate.value = false
    newName.value = ''
    newDescription.value = ''
    if (result?.id) {
      await router.push(`/playlist/${result.id}`)
    } else {
      await fetchPlaylists()
    }
  } catch (err) {
    const toast = useToast()
    toast.add({ severity: 'error', summary: 'Failed to create playlist', detail: 'Please try again later.', life: 4000 })
    console.error('Failed to create playlist:', err)
  } finally {
    creating.value = false
  }
}

onMounted(() => {
  void fetchPlaylists()
})
</script>
