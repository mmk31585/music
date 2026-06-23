<template>
  <div
    class="group rounded-2xl bg-white/6 p-5 ring-1 ring-white/10 transition-all duration-300
           hover:bg-white/8 hover:ring-white/15 focus-within:ring-2 focus-within:ring-spotify"
    role="article"
    :aria-label="`Live room: ${room.title}`"
  >
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2">
          <h3 class="truncate text-sm font-bold text-white">{{ room.title }}</h3>
          <span class="inline-flex items-center gap-1 rounded-full bg-red-500/10 px-2 py-0.5 text-[9px] font-semibold text-red-400 uppercase">
            <span class="h-1.5 w-1.5 rounded-full bg-red-500 motion-safe:animate-pulse" aria-hidden="true" />
            Live
          </span>
        </div>
        <p v-if="room.description" class="mt-1 line-clamp-2 text-xs text-white/40 leading-relaxed">
          {{ room.description }}
        </p>
      </div>
    </div>

    <!-- Listener info -->
    <div class="mt-4 flex items-center gap-2">
      <div class="flex -space-x-1.5" :aria-label="`${room.listener_count || 0} listeners`">
        <div
          v-for="(avatar, i) in roomAvatars.slice(0, 4)"
          :key="i"
          class="h-6 w-6 overflow-hidden rounded-full border-2 border-surface-base transition group-hover:border-white/20"
        >
          <img :src="avatar" alt="" class="h-full w-full object-cover" />
        </div>
        <div
          v-if="(room.listener_count || 0) > 4"
          class="flex h-6 w-6 items-center justify-center rounded-full border-2 border-surface-base bg-white/10 text-[8px] font-bold text-white/50"
        >
          +{{ (room.listener_count || 0) - 4 }}
        </div>
      </div>
      <span class="text-[10px] text-white/30 tabular-nums">{{ room.listener_count || 0 }} listening</span>
    </div>

    <!-- Host info -->
    <div v-if="hostName" class="mt-3 flex items-center gap-1.5 text-[10px] text-white/30">
      <i aria-hidden="true" class="pi pi-user text-[10px] shrink-0" />
      <span class="truncate">{{ hostName }}</span>
    </div>

    <!-- Now playing badge -->
    <div v-if="room.current_track_id" class="mt-1 flex items-center gap-1.5 text-[10px] text-white/25">
      <i aria-hidden="true" class="pi pi-music text-[10px] shrink-0" />
      <span>Now playing</span>
    </div>

    <div class="mt-4 flex items-center gap-2">
      <button
        aria-label="Listen live"
        class="flex-1 rounded-lg bg-red-500/10 py-2.5 text-xs font-bold text-red-400 transition
               hover:bg-red-500/20 active:scale-[0.98] focus-visible:outline-2 focus-visible:outline-red-400"
        @click="$emit('join', room.id)"
      >
        <span class="flex items-center justify-center gap-1.5">
          <i aria-hidden="true" class="pi pi-megaphone text-[10px]" />
          Listen Live
        </span>
      </button>
      <button
        aria-label="Share room"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-white/30 transition
               hover:bg-white/6 hover:text-white/60 focus-visible:outline-2 focus-visible:outline-[#1db954]"
        @click="$emit('share', room.id)"
      >
        <i aria-hidden="true" class="pi pi-share-alt text-sm" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LiveRoom } from '@/services/api/social'

const props = defineProps<{
  room: LiveRoom
  hostName?: string
}>()

defineEmits<{
  join: [id: string]
  share: [id: string]
}>()

const roomAvatars = computed(() => {
  return (props.room as any).avatars || []
})
</script>
