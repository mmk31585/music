<template>
  <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
    <TrackRow
      v-for="(track, index) in tracks"
      :key="String(track.id)"
      :track="track"
      :index="index"
      :queue="queue"
    />
    <div v-if="!tracks.length" class="flex flex-col items-center gap-2 px-6 py-12 text-center">
      <i class="pi pi-music text-2xl text-slate-500" />
      <p class="text-sm text-slate-400">
        {{ emptyMessage || 'No tracks' }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import TrackRow from './TrackRow.vue'
import type { Track } from '@/services/api/catalog/tracks'

const props = defineProps<{
  tracks: Track[]
  emptyMessage?: string
}>()

const queue = computed(() => props.tracks)
</script>
