<template>
  <div class="rounded-2xl bg-white/[0.04] p-6 ring-1 ring-white/[0.07]">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xs font-bold uppercase tracking-wider text-white/30">
        صف پیشنهادی
      </h2>
      <button
        class="inline-flex items-center gap-1.5 rounded-xl bg-[#1db954]/10 px-4 py-2 text-sm font-semibold text-[#1db954] transition hover:bg-[#1db954]/20"
        @click="$emit('suggest-clicked')"
      >
        <i aria-hidden="true" class="pi pi-plus text-xs" />
        پیشنهاد آهنگ
      </button>
    </div>

    <!-- Empty state -->
    <div
      v-if="!candidates.length"
      class="flex flex-col items-center gap-4 py-12 text-sm text-white/30"
    >
      <i aria-hidden="true" class="pi pi-music text-4xl text-white/20" />
      <p>هنوز کسی آهنگی پیشنهاد نداده. اولین نفر باش!</p>
      <button
        class="inline-flex items-center gap-1.5 rounded-xl bg-[#1db954] px-5 py-2.5 text-sm font-bold text-black transition hover:bg-[#1db954]/90"
        @click="$emit('suggest-clicked')"
      >
        <i aria-hidden="true" class="pi pi-plus text-xs" />
        پیشنهاد آهنگ
      </button>
    </div>

    <!-- Candidate list -->
    <TransitionGroup
      v-else
      name="queue-promote"
      tag="div"
      class="space-y-2"
    >
      <div
        v-for="candidate in candidates"
        :key="candidate.id"
        class="group flex items-center gap-3 rounded-xl px-3 py-2.5 ring-1 ring-transparent transition-all hover:ring-white/10"
      >
        <div class="relative h-12 w-12 shrink-0 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="candidate.track.cover_url"
            :src="candidate.track.cover_url"
            :alt="candidate.track.title"
            loading="lazy"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full items-center justify-center">
            <i aria-hidden="true" class="pi pi-headphones text-sm text-white/30" />
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium text-white">{{ candidate.track.title }}</p>
          <p class="truncate text-xs text-white/40">
            پیشنهاد {{ candidate.suggested_by.username || candidate.suggested_by.id.slice(0, 8) }}
          </p>
        </div>

        <button
          class="flex shrink-0 items-center gap-1.5 rounded-lg px-3 py-2 text-sm font-semibold transition hover:bg-white/10"
          :class="candidate.has_voted ? 'text-[#1db954]' : 'text-white/40'"
          @click="toggleVote(candidate)"
        >
          <i
            class="text-lg transition-transform duration-150"
            :class="candidate.has_voted ? 'pi pi-caret-up' : 'pi pi-caret-up'"
            :style="{ transform: candidate.has_voted ? 'scale(1.1)' : 'scale(1)' }"
          />
          <span
            class="tabular-nums"
            :class="{ 'vote-pulse': candidate.vote_count > 0 }"
          >
            {{ candidate.vote_count }}
          </span>
        </button>
      </div>
    </TransitionGroup>

  </div>
</template>

<script setup lang="ts">
import type { QueueCandidate } from '@/services/api/social/room-queue'

defineProps<{
  candidates: readonly QueueCandidate[]
}>()

const emit = defineEmits<{
  vote: [candidateId: string]
  unvote: [candidateId: string]
  'suggest-clicked': []
}>()

function toggleVote(candidate: QueueCandidate) {
  if (candidate.has_voted) {
    emit('unvote', candidate.id)
  } else {
    emit('vote', candidate.id)
  }
}
</script>

<style scoped>
.vote-pulse {
  transition: transform 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.queue-promote-leave-active {
  transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
}

.queue-promote-enter-active {
  transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
  transition-delay: 100ms;
}

.queue-promote-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.queue-promote-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

@media (prefers-reduced-motion: reduce) {
  .queue-promote-enter-active,
  .queue-promote-leave-active {
    transition: none !important;
  }
  .queue-promote-enter-from,
  .queue-promote-leave-to {
    opacity: 1 !important;
    transform: none !important;
  }
}
</style>
