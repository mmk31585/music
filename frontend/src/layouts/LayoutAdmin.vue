<template>
  <div class="h-screen overflow-hidden bg-black text-white" :dir="dir">
    <div class="flex h-full">
      <!-- Desktop sidebar -->
      <AdminSidebar
        :collapsed="sidebarCollapsed"
        class="hidden shrink-0 lg:flex"
        @toggle="sidebarCollapsed = !sidebarCollapsed"
      />

      <!-- Mobile overlay backdrop -->
      <Transition name="fade">
        <div
          v-if="mobileOpen"
          class="fixed inset-0 z-40 bg-black/60 backdrop-blur-xs lg:hidden"
          role="button"
          tabindex="0"
          @click="mobileOpen = false"
          @keydown.enter="mobileOpen = false"
          @keydown.space.prevent="mobileOpen = false"
        />
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
          class="flex-1 overflow-y-auto bg-linear-to-b from-[#151515] to-black"
          :class="mainPadding"
        >
          <router-view />
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
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRTL } from '@/composables'
import { AdminSidebar, AdminTopbar } from '@/components/admin'
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
      showShortcuts.value = showShortcuts.value!
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
