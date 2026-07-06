<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageProgressBar from '@/components/widgets/page-progressbar/PageProgressBar.vue'
import { registerRouter, registerToast } from '@/composables'
import { useToast } from 'primevue/usetoast'
import { useMaintenanceStore, useUserAuthStore, usePlayerStore } from '@/stores'
import { useFeatureFlags } from '@/composables/useFeatureFlags'
import { useLibraryApi } from '@/services/api/library'
import { pendingUnlikeTrackId } from '@/composables/player/useTrackLike'

const route = useRoute()
const maintenanceStore = useMaintenanceStore()
const authStore = useUserAuthStore()
const playerStore = usePlayerStore()
const toast = useToast()

// ── Track-change live region for screen readers (WCAG 4.1.3) ──────
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

// Injection for request wrapper working
registerToast(useToast())
registerRouter(useRouter())

// Watch for auth restore failures and show toast
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

// ── Undo unlike ───────────────────────────────────────────────
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

// Check maintenance mode and feature flags at app startup
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
        <span class="text-sm font-medium">{{ message.summary }}</span>
        <button
          type="button"
          class="rounded-full bg-white/10 px-3 py-1 text-xs font-bold text-white transition hover:bg-white/20"
          @click="undoUnlike"
        >Undo</button>
      </div>
    </template>
  </Toast>

  <!-- Route announcer for screen readers (WCAG 4.1.3, visually hidden) -->
  <div
    id="route-announcer"
    aria-live="polite"
    aria-atomic="true"
    class="sr-only"
  />

  <!-- Track-change announcer for screen readers (WCAG 4.1.3, visually hidden) -->
  <div
    aria-live="polite"
    aria-atomic="true"
    class="sr-only"
  >{{ trackAnnouncement }}</div>

  <component :is="layout">
    <router-view />
  </component>
</template>
