<template>
  <div class="relative">
    <a href="#main-content" class="skip-link">Skip to main content</a>

    <div
      v-if="!isOnline"
      role="alert"
      class="fixed top-0 start-0 end-0 z-9999 flex items-center justify-center gap-2 bg-red-600/90 px-4 py-2 text-sm font-medium text-white backdrop-blur-xs"
      style="padding-top: max(0.5rem, env(safe-area-inset-top, 0.5rem))"
    >
      <Wifi class="text-xs" aria-hidden="true"  />
      <span>You are offline. Some features may be unavailable.</span>
    </div>

    <div
      class="flex h-screen overflow-hidden bg-transparent"
      :dir="rtlDir"
    >
      <MusicSidebar />

      <div class="flex min-w-0 flex-1 flex-col overflow-hidden">
        <MusicAppHeader
          :page-title="pageTitle"
          :unread-count="unreadCount"
          @toggle-mobile="mobileOpen = true"
          @toggle-search="searchOpen = true"
        />

        <main
          id="main-content"
          class="flex-1 overflow-y-auto scroll-smooth scroll-bar"
          :class="mainPadding"
        >
          <ErrorBoundary>
            <RouterView v-slot="{ Component }">
              <Transition name="page" mode="out-in">
                <KeepAlive :max="3">
                  <component :is="Component" />
                </KeepAlive>
              </Transition>
            </RouterView>
          </ErrorBoundary>
        </main>
      </div>

      <MusicRightPane
        @toggle-fullscreen="fullscreenOpen = !fullscreenOpen"
        @toggle-queue-overlay="showQueue = !showQueue"
      />

      <MobileNavPanel v-model="mobileOpen" />

      <SearchOverlay v-model:visible="searchOpen" />
      <RadioMode v-model:visible="radioVisible" :seed-id="radioSeedId" :seed-label="radioSeedLabel" />
      <PlayerRegion />
      <MobileBottomNav />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Wifi } from 'lucide-vue-next'
import { ref, computed, onMounted, provide } from 'vue'
import { useRoute } from 'vue-router'
import {
  MusicSidebar,
  SearchOverlay,
  RadioMode,
  PlayerRegion,
} from '@/components/music'
import MusicRightPane from '@/components/music/layout/MusicRightPane.vue'
import { MobileBottomNav, MusicAppHeader, MobileNavPanel } from '@/components/layouts'
import { usePlayerStore } from '@/stores'
import { useOnlineStatus, useRTL, usePageMeta, useNotificationCount } from '@/composables'
import ErrorBoundary from '@/components/common/ErrorBoundary.vue'

const route = useRoute()
const playerStore = usePlayerStore()

const mobileOpen = ref(false)
const fullscreenOpen = ref(false)
const showQueue = ref(false)
const searchOpen = ref(false)

const radioVisible = ref(false)
const radioSeedId = ref('')
const radioSeedLabel = ref('')

function openRadio(trackId: string, seedLabel?: string) {
  radioSeedId.value = trackId
  radioSeedLabel.value = seedLabel || ''
  radioVisible.value = true
}
provide('openRadio', openRadio)

const { dir: rtlDir } = useRTL()
const { isOnline } = useOnlineStatus()
const { pageTitle } = usePageMeta(route)
const { unreadCount } = useNotificationCount()

const barCollapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')

onMounted(() => {
  const handler = (e: Event) => {
    barCollapsed.value = (e as CustomEvent).detail
  }
  window.addEventListener('playerbar-collapse', handler)
})

const mainPadding = computed(() => {
  if (playerStore.currentTrack) {
    return barCollapsed.value ? 'pb-20 lg:pb-16' : 'pb-32 lg:pb-28'
  }
  return 'pb-24 lg:pb-0'
})
</script>
