<template>
  <NowPlayingBarRoot
    @toggle-queue="$emit('toggle-queue')"
    @toggle-fullscreen="$emit('toggle-fullscreen')"
    @toggle-lyrics="$emit('toggle-lyrics')"
    @toggle-mobile-sheet="$emit('toggle-mobile-sheet')"
    @navigate-artist="$emit('navigate-artist', $event)"
  />
</template>

<script setup lang="ts">
import NowPlayingBarRoot from './NowPlayingBarRoot.vue'

defineEmits<{
  'toggle-queue': []
  'toggle-fullscreen': []
  'toggle-lyrics': []
  'toggle-mobile-sheet': []
  'navigate-artist': [artistId: string]
}>()
</script>