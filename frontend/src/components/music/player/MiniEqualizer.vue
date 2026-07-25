<template>
  <div class="flex items-end gap-[3px]" role="img" aria-label="Audio equalizer">
    <span
      v-for="i in 3"
      :key="i"
      class="w-[3px] rounded-full bg-spotify transition-all"
      :class="[`eq-bar-${i}`, isPlaying ? 'is-playing' : 'is-paused']"
    />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isPlaying: boolean
}>()
</script>

<style scoped>
.eq-bar-1 { animation-duration: 0.6s; animation-delay: 0s; }
.eq-bar-2 { animation-duration: 0.4s; animation-delay: 0.1s; }
.eq-bar-3 { animation-duration: 0.5s; animation-delay: 0.2s; }

.is-paused {
  height: 4px;
  animation-name: none;
}

.is-playing {
  animation-name: mini-equalizer;
  animation-timing-function: ease-in-out;
  animation-iteration-count: infinite;
  animation-direction: alternate;
}

@keyframes mini-equalizer {
  0%, 100% { height: 4px; }
  50% { height: 16px; }
}

@media (prefers-reduced-motion: reduce) {
  .is-playing {
    animation-name: none;
    height: 4px;
  }
}
</style>
