<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import PwaInstallPrompt from '@/components/common/pwa/PwaInstallPrompt.vue'
import PwaUpdateToast from '@/components/common/pwa/PwaUpdateToast.vue'
import { useRoute, useRouter } from 'vue-router'
import { ThemeProvider } from '@/components/layouts'
import PageProgressBar from '@/components/widgets/page-progressbar/PageProgressBar.vue'
import { registerRouter, registerToast } from '@/composables'
import { useToast } from 'primevue/usetoast'
import { useMaintenanceStore, useUserAuthStore, usePlayerStore } from '@/stores'
import { useFeatureFlags } from '@/composables/useFeatureFlags'
import { useLibraryApi } from '@/services/api/library'
import { pendingUnlikeTrackId } from '@/composables/player/useTrackLike'
import { pendingPlaylistRemoveData } from '@/composables/catalog/usePlaylistDetail'
import { usePlaylistsApi } from '@/services/api/playlist'

const route = useRoute()
const maintenanceStore = useMaintenanceStore()
const authStore = useUserAuthStore()
const playerStore = usePlayerStore()
const toast = useToast()

const trackAnnouncement = ref('')
watch(() => playerStore.currentTrack?.title, (title, oldTitle) => {
  if (title && title !== oldTitle) {
    const artist = playerStore.currentTrack?.artistName ?? ''
    trackAnnouncement.value = `Now playing: ${title} by ${artist}`
  }
})

const layout = computed(() => {
  return route?.meta?.layout ?? 'layout-empty'
})

registerToast(useToast())
registerRouter(useRouter())

watch(() => authStore.restoreError, (hasError) => {
  if (hasError) {
    toast.add({
      severity: 'warn',
      summary: 'Session expired',
      detail: 'Your session has expired. Please log in again.',
      life: 5000,
    })
  }
})

const libraryApi = useLibraryApi()
const undoToast = useToast()
async function undoUnlike() {
  const trackId = pendingUnlikeTrackId.value
  if (!trackId) return
  pendingUnlikeTrackId.value = null
  try {
    await libraryApi.likeTrack({ track_id: trackId })
    undoToast.add({ severity: 'success', summary: 'Track re-added to library', life: 3000 })
  } catch {
    undoToast.add({ severity: 'error', summary: 'Failed to undo', life: 3000 })
  }
}

async function undoPlaylistRemove() {
  const data = pendingPlaylistRemoveData.value
  if (!data) return
  pendingPlaylistRemoveData.value = null
  try {
    const api = usePlaylistsApi()
    await api.addTrack(data.playlistId, { track_id: data.trackId })
    undoToast.add({ severity: 'success', summary: `Track re-added to ${data.playlistName}`, life: 3000 })
  } catch {
    undoToast.add({ severity: 'error', summary: 'Failed to undo', life: 3000 })
  }
}

onMounted(() => {
  maintenanceStore.checkStatus()
  useFeatureFlags().init()
})
</script>

<template>
  <PageProgressBar />
  <Toast />
  <Toast group="undo">
    <template #message="{ message }">
      <div class="flex items-center gap-3 px-2 py-1">
        <span class="text-sm font-medium text-primary">{{ message.summary }}</span>
        <button
          type="button"
          class="rounded-full bg-surface-active px-3 py-1 text-xs font-bold text-primary transition hover:bg-surface-hover"
          @click="undoUnlike"
        >Undo</button>
      </div>
    </template>
  </Toast>
  <Toast group="playlist-undo">
    <template #message="{ message }">
      <div class="flex items-center gap-3 px-2 py-1">
        <span class="text-sm font-medium text-primary">{{ message.summary }}</span>
        <button
          type="button"
          class="rounded-full bg-surface-active px-3 py-1 text-xs font-bold text-primary transition hover:bg-surface-hover"
          @click="undoPlaylistRemove"
        >Undo</button>
      </div>
    </template>
  </Toast>

  <div
    id="route-announcer"
    aria-live="polite"
    aria-atomic="true"
    class="sr-only"
  />

  <div
    aria-live="polite"
    aria-atomic="true"
    class="sr-only"
  >{{ trackAnnouncement }}</div>

  <ThemeProvider>
    <PwaInstallPrompt />
    <PwaUpdateToast />

    <component :is="layout">
      <router-view />
    </component>
  </ThemeProvider>
</template>
