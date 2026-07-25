<template>
  <div class="glass-strong rounded-2xl p-6 md:p-8">
    <div class="mb-6 flex items-center justify-between">
      <h3 class="text-lg font-bold text-white">Challenges</h3>
      <span class="rounded-full bg-spotify/10 px-2.5 py-0.5 text-xs font-medium text-spotify">
        +{{ earnedXP }} XP earned
      </span>
    </div>

    <div v-if="!challenges.length" class="py-8 text-center text-sm text-white/30">
      No active challenges right now.
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="item in combined"
        :key="item.challenge.id"
        class="overflow-hidden rounded-xl border border-white/5 transition-all duration-300"
        :class="item.progress?.is_completed ? 'bg-spotify/5' : 'bg-white/3'"
      >
        <div class="p-4">
          <div class="flex items-start justify-between gap-4">
            <div class="flex items-start gap-3 min-w-0 flex-1">
              <!-- Challenge type icon -->
              <div class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white/5 text-lg">
                {{ challengeIcon(item.challenge.challenge_type) }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold text-white">{{ item.challenge.title }}</p>
                <p class="mt-0.5 text-xs text-white/40">{{ item.challenge.description }}</p>
                <!-- Time remaining -->
                <p class="mt-1 text-[10px] text-white/20">
                  {{ timeRemaining(item.challenge.valid_until) }}
                </p>
              </div>
            </div>
            <div class="shrink-0 text-right">
              <span
                class="inline-flex items-center gap-1 rounded-full bg-amber-500/10 px-2.5 py-0.5 text-xs font-medium text-amber-400"
              >
                +{{ item.challenge.xp_reward }} XP
              </span>
            </div>
          </div>

          <!-- Progress bar -->
          <div class="mt-3">
            <div class="flex items-center justify-between text-[11px] text-white/30">
              <span>{{ item.progress?.progress || 0 }} / {{ item.challenge.target_count }}</span>
              <span v-if="item.progress?.is_completed" class="text-spotify font-medium">Completed!</span>
            </div>
            <div class="mt-1.5 h-2 overflow-hidden rounded-full bg-white/5">
              <div
                class="h-full rounded-full transition-all duration-700"
                :class="item.progress?.is_completed ? 'bg-spotify' : 'bg-linear-to-r from-amber-400 to-orange-500'"
                :style="{ width: `${completionPercent(item)}%` }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { DailyChallenge, UserChallenge } from '@/services/api/gamification'

const props = defineProps<{
  challenges: DailyChallenge[]
  progress: UserChallenge[]
  earnedXP: number
}>()

interface ChallengeItem {
  challenge: DailyChallenge
  progress: UserChallenge | undefined
}

const combined = computed<ChallengeItem[]>(() => {
  return props.challenges.map((ch) => ({
    challenge: ch,
    progress: props.progress.find((p) => p.challenge_id === ch.id),
  }))
})

const CHALLENGE_ICONS: Record<string, string> = {
  streams: '🎵',
  likes: '❤️',
  shares: '🔗',
  upload: '📤',
  publish: '📢',
  contribute: '✏️',
}

function challengeIcon(type: string): string {
  return CHALLENGE_ICONS[type] ?? '🎯'
}

function completionPercent(item: ChallengeItem): number {
  if (!item.challenge.target_count) return 0
  const progress = item.progress?.progress || 0
  return Math.min(100, (progress / item.challenge.target_count) * 100)
}

function timeRemaining(dateStr: string): string {
  if (!dateStr) return ''
  const diff = new Date(dateStr).getTime() - Date.now()
  if (diff <= 0) return 'Expired'
  const hours = Math.floor(diff / (1000 * 60 * 60))
  if (hours < 1) return 'Ending soon'
  if (hours < 24) return `${hours}h remaining`
  const days = Math.floor(hours / 24)
  return `${days}d remaining`
}
</script>
