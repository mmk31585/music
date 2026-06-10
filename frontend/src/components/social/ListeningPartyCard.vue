<template>
  <div
    class="group overflow-hidden rounded-xl border border-white/5 bg-white/[0.03] transition-all duration-300 hover:border-white/10 hover:bg-white/[0.06]"
  >
    <div class="p-5">
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold text-white">{{ party.title }}</p>
          <p v-if="hostName" class="mt-0.5 text-xs text-white/30">{{ hostName }}</p>
          <p v-if="party.description" class="mt-1 line-clamp-2 text-xs text-white/40">
            {{ party.description }}
          </p>
        </div>
        <div
          class="flex shrink-0 items-center gap-1.5 rounded-full px-2.5 py-1 text-[10px] font-semibold tracking-wider uppercase"
          :class="statusClass"
        >
          <span
            v-if="party.status === 'active'"
            class="h-1.5 w-1.5 animate-pulse rounded-full bg-green-500"
          />
          {{ party.status }}
        </div>
      </div>

      <div class="mt-4 flex items-center gap-4 text-xs text-white/30">
        <span class="flex items-center gap-1">👥 {{ party.participant_count }}</span>
        <span v-if="party.current_track_id && currentTrackName" class="truncate" :title="currentTrackName">
          🎵 {{ currentTrackName }}
        </span>
        <span v-else-if="party.current_track_id" class="flex items-center gap-1">🎵 Now playing</span>
        <span v-else class="text-white/20">No track</span>
      </div>

      <div class="mt-3 flex items-center gap-2">
        <button
          class="flex-1 rounded-lg bg-[#1db954]/10 py-2 text-xs font-semibold text-[#1db954] transition hover:bg-[#1db954]/20"
          @click="$emit('join', party.id)"
        >
          Join Party
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ListeningParty } from '@/services/api/social'

const props = defineProps<{
  party: ListeningParty
  hostName?: string
  currentTrackName?: string
}>()

defineEmits<{
  join: [id: string]
}>()

const statusClass = computed(() => {
  switch (props.party.status) {
    case 'active': return 'bg-green-500/10 text-green-400'
    case 'paused': return 'bg-amber-500/10 text-amber-400'
    case 'ended': return 'bg-white/10 text-white/30'
    default: return 'bg-white/10 text-white/30'
  }
})
</script>
