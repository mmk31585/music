<template>
  <div
    class="group overflow-hidden rounded-xl border border-white/5 bg-white/[0.03] transition-all duration-300 hover:border-white/10 hover:bg-white/[0.06]"
  >
    <div class="p-5">
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold text-white">{{ room.title }}</p>
          <p v-if="room.description" class="mt-1 line-clamp-2 text-xs text-white/40">
            {{ room.description }}
          </p>
        </div>
        <div class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-red-500/10">
          <span class="h-2 w-2 animate-pulse rounded-full bg-red-500" />
        </div>
      </div>

      <div class="mt-4 flex items-center gap-4 text-xs text-white/30">
        <span class="flex items-center gap-1">🎤 {{ room.listener_count }} listening</span>
        <span v-if="room.current_track_id" class="flex items-center gap-1">🎵 Now playing</span>
      </div>

      <div class="mt-3 flex items-center gap-2">
        <button
          class="flex-1 rounded-lg bg-red-500/10 py-2 text-xs font-semibold text-red-400 transition hover:bg-red-500/20"
          @click="$emit('join', room.id)"
        >
          Listen Live
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LiveRoom } from '@/services/api/social'

defineProps<{
  room: LiveRoom
}>()

defineEmits<{
  join: [id: string]
}>()
</script>
