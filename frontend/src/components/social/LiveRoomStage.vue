<template>
  <div class="rounded-2xl bg-white/[0.04] p-6 ring-1 ring-white/[0.07]">
    <h2 class="mb-4 text-xs font-bold uppercase tracking-wider text-white/30">Stage</h2>

    <div class="flex flex-col items-center gap-6">
      <!-- Host avatar (always shown, largest) -->
      <div class="flex flex-col items-center gap-2">
        <div
          class="relative h-20 w-20 overflow-hidden rounded-full ring-2"
          :class="stageState.host?.user_id ? 'ring-[#1db954]' : 'ring-white/10'"
        >
          <img
            v-if="stageState.host?.avatar_url"
            :src="stageState.host?.avatar_url"
            :alt="stageState.host?.username || 'Host'"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full w-full items-center justify-center bg-white/10 text-2xl">
            👑
          </div>
          <!-- Speaking pulse ring (stub) -->
          <div
            v-if="speakingUserIds.has(stageState.host?.user_id ?? '')"
            class="absolute inset-0 animate-pulse rounded-full ring-2 ring-[#1db954] ring-offset-2 ring-offset-transparent"
          />
        </div>
        <p class="text-sm font-semibold text-white">{{ stageState.host?.username || 'Host' }}</p>
        <span class="rounded-full bg-yellow-500/10 px-2.5 py-0.5 text-[10px] font-semibold text-yellow-400">Host</span>
      </div>

      <!-- Speakers grid -->
      <div v-if="stageState.speakers?.length" class="flex flex-wrap justify-center gap-4">
        <div
          v-for="speaker in stageState.speakers"
          :key="speaker.user_id"
          class="group relative flex flex-col items-center gap-2"
          @mouseenter="hoveredSpeaker = speaker.user_id"
          @mouseleave="hoveredSpeaker = null"
        >
          <div
            class="relative h-16 w-16 overflow-hidden rounded-full ring-2"
            :class="speakingUserIds.has(speaker.user_id) ? 'ring-[#1db954]' : 'ring-white/10'"
          >
            <img
              v-if="speaker?.avatar_url"
              :src="speaker.avatar_url"
              :alt="speaker?.username || 'Speaker'"
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-white/10 text-lg font-bold text-white">
              {{ (speaker?.username || '?').charAt(0).toUpperCase() }}
            </div>

            <!-- Muted badge -->
            <div
              v-if="speaker.muted"
              class="absolute bottom-0 right-0 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-[10px] text-white"
            >
              <i aria-hidden="true" class="pi pi-microphone-off" />
            </div>

            <!-- Speaking pulse ring (stub — TODO: wire to real audio levels) -->
            <div
              v-if="speakingUserIds.has(speaker.user_id)"
              class="absolute inset-0 animate-pulse rounded-full ring-2 ring-[#1db954] ring-offset-2 ring-offset-transparent"
            />
          </div>

          <p class="truncate text-xs font-medium text-white">{{ speaker?.username || 'Unknown' }}</p>

          <!-- Host actions on hover (only visible to host) -->
          <div
            v-if="showHostActions && hoveredSpeaker === speaker.user_id"
            class="absolute -bottom-12 left-1/2 z-10 flex -translate-x-1/2 gap-1 rounded-xl bg-[#1a1a1a] px-2 py-1.5 shadow-lg ring-1 ring-white/10"
          >
            <button
              class="rounded-lg px-2 py-1 text-[10px] text-white/60 transition hover:bg-white/10 hover:text-white"
              @click="$emit('toggle-mute', { userId: speaker.user_id, muted: !speaker.muted })"
            >
              {{ speaker.muted ? 'رفع بی‌صوتی' : 'بی‌صدا کردن' }}
            </button>
            <button
              class="rounded-lg px-2 py-1 text-[10px] text-red-400 transition hover:bg-red-500/10"
              @click="$emit('remove-speaker', speaker.user_id)"
            >
              حذف از استیج
            </button>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="py-4 text-center text-sm text-white/30">
        هنوز کسی روی استیج نیست
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

/*
 * TODO: speakingUserIds is a stub. Wire to real audio level detection
 * once WebRTC / SFU integration is added. The prop accepts a Set<string>
 * of user_ids that are currently producing audio.
 */
defineProps<{
  stageState: {
    host: { user_id: string; username?: string; avatar_url?: string } | null
    speakers: readonly { user_id: string; username?: string; avatar_url?: string; muted: boolean }[]
  }
  speakingUserIds: Set<string>
  showHostActions: boolean
}>()

defineEmits<{
  'remove-speaker': [userId: string]
  'toggle-mute': [payload: { userId: string; muted: boolean }]
}>()

const hoveredSpeaker = ref<string | null>(null)
</script>

<style scoped>
.animate-pulse {
  animation: pulse-ring 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes pulse-ring {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.6;
  }
}

@media (prefers-reduced-motion: reduce) {
  .animate-pulse {
    animation: none !important;
  }
}
</style>
