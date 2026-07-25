<template>
  <div class="flex items-center gap-1" :dir="dir">
    <button
      type="button"
      class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
      :aria-label="muted ? 'Unmute' : 'Mute'"
      @click="$emit('mute')"
    >
      <VolumeX v-if="muted || volume === 0" class="h-4 w-4" aria-hidden="true" />
      <Volume2 v-else class="h-4 w-4" aria-hidden="true" />
    </button>

    <div class="w-20 relative">
      <input
        type="range"
        class="player-range"
        :min="0"
        :max="1"
        :step="0.01"
        :value="muted ? 0 : volume"
        @input="onVolume"
        @change="onVolume"
        aria-label="Volume"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Volume2, VolumeX } from 'lucide-vue-next'
import { useRTL } from '@/composables/useRTL'

interface Props {
  volume: number
  muted: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  volume: [value: number]
  mute: []
}>()

const { dir } = useRTL()

function onVolume(e: Event) {
  const val = (e.target as HTMLInputElement).valueAsNumber
  emit('volume', val)
}
</script>

<style scoped>
/* Uses .player-range styles from main.css */
</style>