<template>
  <!-- Glass border at 0.10 for 3:1 non-text contrast (WCAG 1.4.11) -->
  <div class="rounded-2xl bg-surface-overlay/60 p-6 ring-1 ring-border-default">
    <h2 class="mb-4 text-xs font-bold uppercase tracking-wider text-muted">Stage</h2>

    <div class="flex flex-col items-center gap-6">
      <!-- Host avatar (always shown, largest) -->
      <div class="flex flex-col items-center gap-2">
        <div
          class="relative h-20 w-20 overflow-hidden rounded-full ring-2"
          :class="stageState.host?.user_id ? 'ring-accent' : 'ring-border-default'"
        >
          <img
            v-if="stageState.host?.avatar_url"
            :src="stageState.host?.avatar_url"
            :alt="stageState.host?.username || 'Host'"
            class="h-full w-full object-cover"
          />
          <!-- Lucide Crown SVG icon (replaces 👑 emoji — emoji-as-icon anti-pattern) -->
          <div v-else class="flex h-full w-full items-center justify-center bg-surface-active">
            <svg class="h-7 w-7 text-warning" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M2 4l3 12h14l3-12-6 7-4-7-4 7-6-7z" />
              <path d="M3 20h18" />
            </svg>
          </div>
          <!-- Speaking pulse ring (stub) -->
          <div
            v-if="speakingUserIds.has(stageState.host?.user_id ?? '')"
            class="absolute inset-0 animate-pulse rounded-full ring-2 ring-accent ring-offset-2 ring-offset-transparent"
          />
        </div>
        <p class="text-sm font-semibold text-primary">{{ stageState.host?.username || 'Host' }}</p>
        <span class="rounded-full bg-warning-subtle px-2.5 py-0.5 text-[10px] font-semibold text-warning">Host</span>
      </div>

      <!-- Speakers grid -->
      <div v-if="stageState.speakers?.length" class="flex flex-wrap justify-center gap-4">
        <div
          v-for="speaker in stageState.speakers"
          :key="speaker.user_id"
          :data-speaker-id="speaker.user_id"
          class="group relative flex flex-col items-center gap-2"
          @mouseenter="hoveredSpeaker = speaker.user_id"
          @mouseleave="hoveredSpeaker = null"
          @focusin="hoveredSpeaker = speaker.user_id"
          @focusout="onSpeakerFocusOut(speaker.user_id)"
        >
          <div
            class="relative h-16 w-16 overflow-hidden rounded-full ring-2"
            :class="speakingUserIds.has(speaker.user_id) ? 'ring-accent' : 'ring-border-default'"
          >
            <img
              v-if="speaker?.avatar_url"
              :src="speaker.avatar_url"
              :alt="speaker?.username || 'Speaker'"
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-surface-active text-lg font-bold text-primary">
              {{ (speaker?.username || '?').charAt(0).toUpperCase() }}
            </div>

            <!-- Muted badge -->
            <div
              v-if="speaker.muted"
              class="absolute bottom-0 right-0 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-[10px] text-primary"
            >
              <MicOff aria-hidden="true" class=""  />
            </div>

            <!-- Speaking pulse ring (stub — TODO: wire to real audio levels) -->
            <div
              v-if="speakingUserIds.has(speaker.user_id)"
              class="absolute inset-0 animate-pulse rounded-full ring-2 ring-accent ring-offset-2 ring-offset-transparent"
              role="status"
              aria-label="Currently speaking"
            />
          </div>

          <p class="truncate text-xs font-medium text-primary">{{ speaker?.username || 'Unknown' }}</p>

          <!-- Host actions: visible on hover OR when any child is focused (keyboard accessible) -->
          <div
            v-if="showHostActions && hoveredSpeaker === speaker.user_id"
            class="absolute -bottom-12 left-1/2 z-10 flex -translate-x-1/2 gap-1 rounded-xl bg-surface-overlay px-2 py-1.5 shadow-lg ring-1 ring-border-default"
            role="toolbar"
            :aria-label="`Actions for ${speaker?.username || 'speaker'}`"
          >
            <button
              class="rounded-lg px-3 py-1.5 text-[10px] text-primary/60 transition hover:bg-surface-active hover:text-primary focus-visible:outline-2 focus-visible:outline-accent"
              @click="$emit('toggle-mute', { userId: speaker.user_id, muted: !speaker.muted })"
            >
              {{ speaker.muted ? 'رفع بی‌صوتی' : 'بی‌صدا کردن' }}
            </button>
            <button
              class="rounded-lg px-3 py-1.5 text-[10px] text-danger transition hover:bg-danger-subtle focus-visible:outline-2 focus-visible:outline-red-400"
              @click="$emit('remove-speaker', speaker.user_id)"
            >
              حذف از استیج
            </button>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="py-4 text-center text-sm text-muted">
        هنوز کسی روی استیج نیست
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { MicOff } from 'lucide-vue-next'
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

/**
 * On focusout, check if focus moved to a child button inside this speaker's host-actions.
 * If so, keep hoveredSpeaker set so the toolbar stays visible.
 * If focus left the entire speaker container, clear hoveredSpeaker.
 */
function onSpeakerFocusOut(speakerId: string) {
  // Use requestAnimationFrame to let the browser update document.activeElement first
  requestAnimationFrame(() => {
    const focused = document.activeElement
    const container = focused?.closest(`[data-speaker-id="${speakerId}"]`)
    if (!container) {
      hoveredSpeaker.value = null
    }
  })
}
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
