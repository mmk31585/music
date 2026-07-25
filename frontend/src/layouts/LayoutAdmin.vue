<template>
  <div class="h-screen overflow-hidden bg-[var(--bg-base)] text-[var(--text-primary)]" :dir="dir">
    <!-- Skip link for keyboard users -->
    <a href="#main-content" class="skip-link">Skip to main content</a>

    <!-- Offline banner -->
    <div
      v-if="!isOnline"
      role="alert"
      class="fixed top-0 left-0 right-0 z-9999 flex items-center justify-center gap-2 bg-red-600/90 px-4 py-2 text-sm font-medium text-white backdrop-blur-xs"
      style="padding-top: max(0.5rem, env(safe-area-inset-top, 0.5rem))"
    >
      <Wifi class="text-xs" aria-hidden="true"  />
      <span>You are offline. Some features may be unavailable.</span>
    </div>

    <div class="flex h-full">
      <!-- Desktop sidebar -->
      <AdminSidebar
        :collapsed="sidebarCollapsed"
        class="hidden shrink-0 lg:flex"
        @toggle="sidebarCollapsed = !sidebarCollapsed"
      />

      <!-- Mobile overlay backdrop -->
      <Transition name="fade">
        <button
          type="button"
          v-if="mobileOpen"
          class="fixed inset-0 z-40 bg-[var(--bg-overlay)] backdrop-blur-xs lg:hidden"
          aria-label="Close sidebar"
          @click="mobileOpen = false"
        ></button>
      </Transition>

      <!-- Mobile sidebar drawer -->
      <Transition name="slide">
        <AdminSidebar
          v-if="mobileOpen"
          :collapsed="false"
          class="lg:hidden!"
          @close="mobileOpen = false"
        />
      </Transition>

      <!-- Main content area (scrolls independently) -->
      <div class="flex min-w-0 flex-1 flex-col overflow-hidden">
        <AdminTopbar
          :collapsed="sidebarCollapsed"
          @toggle-mobile="mobileOpen = !mobileOpen"
          @toggle-collapse="sidebarCollapsed = !sidebarCollapsed"
        />
        <main
          id="main-content"
          class="flex-1 overflow-y-auto bg-linear-to-b from-[var(--surface-1)] to-[var(--surface-0)]"
          :class="mainPadding"
        >
          <ErrorBoundary>
            <router-view />
          </ErrorBoundary>
        </main>
      </div>
    </div>

    <!-- ── Player Components ── -->
    <NowPlayingBar
      @toggle-fullscreen="fullscreenOpen = !fullscreenOpen"
      @toggle-queue="showQueue = !showQueue"
      @toggle-lyrics="onToggleLyrics"
      @toggle-mobile-sheet="mobileSheetOpen = !mobileSheetOpen"
    />
    <FullscreenPlayer v-if="!ffEnabled" v-model:visible="fullscreenOpen" :initial-tab="playerInitialTab" />
    <ExpandedPlayer v-else v-model:visible="fullscreenOpen" />
    <QueuePanel v-model:visible="showQueue" />
    <MobileBottomSheet v-model:visible="mobileSheetOpen" @open-fullscreen="fullscreenOpen = true" />
    <KeyboardShortcuts v-model:visible="showShortcuts" />
  </div>
</template>

<script setup lang="ts">
import { Wifi } from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRTL, useOnlineStatus } from '@/composables'
import { AdminSidebar, AdminTopbar } from '@/components/admin'
import ErrorBoundary from '@/components/common/ErrorBoundary.vue'
import {
  NowPlayingBar,
  FullscreenPlayer,
  ExpandedPlayer,
  QueuePanel,
  MobileBottomSheet,
  KeyboardShortcuts,
} from '@/components/music'
import { usePlayerStore, useFeatureFlagsStore } from '@/stores'

const { dir } = useRTL()
const { isOnline } = useOnlineStatus()
const sidebarCollapsed = ref(false)
const mobileOpen = ref(false)
const fullscreenOpen = ref(false)
const showQueue = ref(false)
const mobileSheetOpen = ref(false)
const showShortcuts = ref(false)
const playerInitialTab = ref<'now-playing' | 'queue' | 'lyrics'>('now-playing')

const playerStore = usePlayerStore()
const ff = useFeatureFlagsStore()
const ffEnabled = computed(() => ff.isEnabled('redesignedPlayer'))

const barCollapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')

function onBarCollapse(e: Event) {
  barCollapsed.value = (e as CustomEvent).detail
}

onMounted(() => {
  window.addEventListener('playerbar-collapse', onBarCollapse)

  // Keyboard shortcuts
  document.addEventListener('keydown', (e) => {
    const tag = (e.target as HTMLElement)?.tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return
    if (e.key === '?' && e.ctrlKey! && e.metaKey! && e.altKey!) {
      showShortcuts.value = !showShortcuts.value
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('playerbar-collapse', onBarCollapse)
})

const mainPadding = computed(() => {
  if (playerStore.currentTrack) {
    return barCollapsed.value ? 'pb-20 lg:pb-16' : 'pb-32 lg:pb-28'
  }
  return 'pb-0'
})

function onToggleLyrics() {
  playerInitialTab.value = 'lyrics'
  fullscreenOpen.value = true
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.25s ease;
}
/* LTR: slide in from left. RTL: slide in from right */
html:not([dir="rtl"]) .slide-enter-from,
html:not([dir="rtl"]) .slide-leave-to {
  transform: translateX(-100%);
}
html[dir="rtl"] .slide-enter-from,
html[dir="rtl"] .slide-leave-to {
  transform: translateX(100%);
}
</style>
