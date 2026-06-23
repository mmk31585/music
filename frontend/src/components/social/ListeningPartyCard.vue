<template>
  <div
    class="group rounded-2xl bg-white/[0.06] p-5 ring-1 ring-white/[0.10] transition-all duration-300
           hover:bg-white/[0.08] hover:ring-white/[0.15] focus-within:ring-2 focus-within:ring-[#1db954]"
    role="article"
    :aria-label="`Party: ${party.title}`"
  >
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2">
          <h3 class="truncate text-sm font-bold text-white">{{ party.title }}</h3>
          <span
            v-if="party.status === 'active'"
            class="inline-flex items-center gap-1 rounded-full bg-green-500/10 px-2 py-0.5 text-[9px] font-semibold text-green-400 uppercase"
          >
            <span class="h-1.5 w-1.5 rounded-full bg-green-500 motion-safe:animate-pulse" aria-hidden="true" />
            Live
          </span>
        </div>
        <p v-if="hostName" class="mt-0.5 text-xs text-white/30">{{ hostName }}</p>
        <p v-if="party.description" class="mt-1 line-clamp-2 text-xs text-white/40 leading-relaxed">
          {{ party.description }}
        </p>
      </div>

      <!-- Cover art thumbnail if there's a current track -->
      <div
        v-if="party.currentTrackCover || party.current_track_id"
        class="h-14 w-14 shrink-0 overflow-hidden rounded-xl ring-1 ring-white/10"
      >
        <img
          v-if="party.currentTrackCover"
          :src="party.currentTrackCover"
          alt=""
          class="h-full w-full object-cover"
        />
        <div v-else class="flex h-full w-full items-center justify-center bg-white/[0.04]">
          <i aria-hidden="true" class="pi pi-music text-sm text-white/20" />
        </div>
      </div>
    </div>

    <!-- Participant avatars row -->
    <div class="mt-4 flex items-center gap-2">
      <div class="flex -space-x-1.5" :aria-label="`${party.participant_count || 0} participants`">
        <div
          v-for="(avatar, i) in partyAvatars.slice(0, 4)"
          :key="i"
          class="h-6 w-6 overflow-hidden rounded-full border-2 border-[#0A0A0A] transition group-hover:border-white/20"
        >
          <img :src="avatar" alt="" class="h-full w-full object-cover" />
        </div>
        <div
          v-if="(party.participant_count || 0) > 4"
          class="flex h-6 w-6 items-center justify-center rounded-full border-2 border-[#0A0A0A] bg-white/10 text-[8px] font-bold text-white/50"
        >
          +{{ (party.participant_count || 0) - 4 }}
        </div>
      </div>
      <span class="text-[10px] text-white/30 tabular-nums">{{ party.participant_count || 0 }} listening</span>
    </div>

    <!-- Now playing mini info -->
    <div v-if="currentTrackName" class="mt-3 flex items-center gap-1.5 text-[10px] text-white/30">
      <i aria-hidden="true" class="pi pi-music text-[10px] shrink-0" />
      <span class="truncate">{{ currentTrackName }}</span>
    </div>

    <div class="mt-3 flex items-center gap-2">
      <button
        aria-label="Join party"
        class="flex-1 rounded-lg bg-[#1db954]/10 py-2.5 text-xs font-bold text-[#1db954] transition
               hover:bg-[#1db954]/20 active:scale-[0.98] focus-visible:outline-2 focus-visible:outline-[#1db954]"
        @click="$emit('join', party.id)"
      >
        <span class="flex items-center justify-center gap-1.5">
          <i aria-hidden="true" class="pi pi-headphones text-[10px]" />
          Join Party
        </span>
      </button>
      <button
        aria-label="Share party"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-white/30 transition
               hover:bg-white/[0.06] hover:text-white/60 focus-visible:outline-2 focus-visible:outline-[#1db954]"
        @click="$emit('share', party.id)"
      >
        <i aria-hidden="true" class="pi pi-share-alt text-sm" />
      </button>
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
  share: [id: string]
}>()

const partyAvatars = computed(() => {
  // Use provided avatars or generate initials-based placeholders
  return (props.party as any).avatars || []
})
</script>
