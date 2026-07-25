<script setup lang="ts">
import { computed, ref } from 'vue'
import { useFeatureFlagsStore } from '@/stores'
import { useKeyboardShortcuts } from '@/composables/player'
import NowPlayingBar from './NowPlayingBar.vue'
import FullscreenPlayer from './FullscreenPlayer.vue'
import ExpandedPlayer from './ExpandedPlayer.vue'
import QueuePanel from './QueuePanel.vue'
import MobileBottomSheet from './MobileBottomSheet.vue'
import KeyboardShortcuts from './KeyboardShortcuts.vue'

const showQueue = ref(false)
const fullscreenOpen = ref(false)
const mobileSheetOpen = ref(false)
const showShortcuts = ref(false)
const playerInitialTab = ref<'now-playing' | 'queue' | 'lyrics'>('now-playing')
const ff = useFeatureFlagsStore()
const ffEnabled = computed(() => ff.isEnabled('redesignedPlayer'))

function onToggleLyrics() {
  playerInitialTab.value = 'lyrics'
  fullscreenOpen.value = true
}

// Wire global keyboard shortcuts to actual player actions
useKeyboardShortcuts({
  onToggleFullscreen: () => { fullscreenOpen.value = !fullscreenOpen.value },
  onToggleQueue: () => { showQueue.value = !showQueue.value },
  onToggleLyrics: () => { onToggleLyrics() },
  onToggleShortcuts: () => { showShortcuts.value = !showShortcuts.value },
  onToggleMobileSheet: () => { mobileSheetOpen.value = !mobileSheetOpen.value },
})
</script>

<template>
  <NowPlayingBar
    @toggle-fullscreen="fullscreenOpen = !fullscreenOpen"
    @toggle-queue="showQueue = !showQueue"
    @toggle-lyrics="onToggleLyrics"
    @toggle-mobile-sheet="mobileSheetOpen = !mobileSheetOpen"
  />
  <FullscreenPlayer
    v-if="!ffEnabled"
    v-model:visible="fullscreenOpen"
    :initial-tab="playerInitialTab"
  />
  <ExpandedPlayer v-else v-model:visible="fullscreenOpen" />
  <QueuePanel v-model:visible="showQueue" />
  <MobileBottomSheet
    v-model:visible="mobileSheetOpen"
    @open-fullscreen="fullscreenOpen = true"
  />
  <KeyboardShortcuts v-model:visible="showShortcuts" />
</template>
