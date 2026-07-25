<template>
  <div
    class="group rounded-2xl bg-surface-overlay/60 p-5 ring-1 ring-border-default transition-all duration-300
           hover:bg-surface-active/80 hover:ring-border-strong focus-within:ring-2 focus-within:ring-spotify"
    role="article"
    :aria-label="`Live room: ${room.title}`"
  >
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2">
          <h3 class="truncate text-sm font-bold text-primary">{{ room.title }}</h3>
          <span class="inline-flex items-center gap-1 rounded-full bg-danger-subtle px-2 py-0.5 text-[9px] font-semibold text-danger uppercase">
            <span class="h-1.5 w-1.5 rounded-full bg-red-500 motion-safe:animate-pulse" aria-hidden="true" />
            Live
          </span>
        </div>
        <p v-if="room.description" class="mt-1 line-clamp-2 text-xs text-primary/40 leading-relaxed">
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
          class="h-6 w-6 overflow-hidden rounded-full border-2 border-surface-base transition group-hover:border-border-strong"
        >
          <img :src="avatar" alt="" class="h-full w-full object-cover" />
        </div>
        <div
          v-if="(room.listener_count || 0) > 4"
          class="flex h-6 w-6 items-center justify-center rounded-full border-2 border-surface-base bg-surface-active text-[8px] font-bold text-primary/50"
        >
          +{{ (room.listener_count || 0) - 4 }}
        </div>
      </div>
      <span class="text-[10px] text-primary/30 tabular-nums">{{ room.listener_count || 0 }} listening</span>
    </div>

    <!-- Host info -->
    <div v-if="hostName" class="mt-3 flex items-center gap-1.5 text-[10px] text-primary/30">
      <User aria-hidden="true" class="text-[10px] shrink-0"  />
      <span class="truncate">{{ hostName }}</span>
    </div>

    <!-- Now playing badge -->
    <div v-if="room.current_track_id" class="mt-1 flex items-center gap-1.5 text-[10px] text-primary/25">
      <Music aria-hidden="true" class="text-[10px] shrink-0"  />
      <span>Now playing</span>
    </div>

    <div class="mt-4 flex items-center gap-2">
      <button
        aria-label="Listen live"
        class="flex-1 rounded-lg bg-danger-subtle py-2.5 text-xs font-bold text-danger transition
               hover:bg-danger-subtle active:scale-[0.98] focus-visible:outline-2 focus-visible:outline-red-400"
        @click="$emit('join', room.id)"
      >
        <span class="flex items-center justify-center gap-1.5">
          <Megaphone aria-hidden="true" class="text-[10px]"  />
          Listen Live
        </span>
      </button>
      <button
        aria-label="Share room"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-primary/30 transition
               hover:bg-surface-overlay/60 hover:text-primary/60 focus-visible:outline-2 focus-visible:outline-accent"
        @click="$emit('share', room.id)"
      >
        <Share2 aria-hidden="true" class="text-sm"  />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Megaphone, Music, Share2, User } from 'lucide-vue-next'
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
