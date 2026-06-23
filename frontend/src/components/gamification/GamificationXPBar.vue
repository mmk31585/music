<template>
  <div v-if="profile" class="glass-strong rounded-2xl p-6 md:p-8">
    <div class="flex items-center gap-5">
      <!-- Avatar -->
      <div class="relative h-16 w-16 shrink-0 md:h-20 md:w-20">
        <div class="h-full w-full overflow-hidden rounded-2xl bg-white/10 ring-2 ring-spotify/30">
          <img
            v-if="profile.avatar_url"
            :src="profile.avatar_url"
            :alt="profile.username"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div
            v-else
            class="flex h-full w-full items-center justify-center text-xl font-bold text-white/40"
          >
            {{ profile.username.charAt(0).toUpperCase() }}
          </div>
        </div>
        <!-- Level badge -->
        <div
          class="absolute -right-1 -bottom-1 flex h-7 w-7 items-center justify-center rounded-full bg-spotify text-[10px] font-bold text-black shadow-lg md:h-8 md:w-8 md:text-xs"
        >
          {{ profile.level }}
        </div>
      </div>

      <!-- Info -->
      <div class="min-w-0 flex-1">
        <p class="text-lg font-bold text-white md:text-xl">{{ profile.username }}</p>
        <p class="text-sm text-white/50">{{ profile.title }} · {{ profile.title_persian }}</p>
        <p class="mt-0.5 text-xs text-white/30">
          Rank #{{ profile.rank }} · {{ profile.total_xp.toLocaleString() }} XP
        </p>
      </div>
    </div>

    <!-- XP Progress bar -->
    <div class="mt-5">
      <div class="flex items-center justify-between text-xs text-white/40">
        <span>Level {{ profile.level }}</span>
        <span
          >{{ profile.current_xp.toLocaleString() }} /
          {{ profile.next_level_xp.toLocaleString() }} XP</span
        >
      </div>
      <div class="mt-2 h-2 overflow-hidden rounded-full bg-white/5">
        <div
          class="h-full rounded-full bg-linear-to-r from-spotify to-aurora-blue transition-all duration-1000 ease-out"
          :style="{ width: `${progressPercent}%` }"
        />
      </div>
      <p class="mt-1 text-[10px] text-white/20">
        {{ xpToNextLevel.toLocaleString() }} XP to next level
      </p>
    </div>
  </div>

  <div v-else class="glass-strong flex items-center gap-4 rounded-2xl p-6">
    <div class="h-16 w-16 animate-pulse rounded-2xl bg-white/5" />
    <div class="flex-1 space-y-2">
      <div class="h-5 w-32 animate-pulse rounded bg-white/5" />
      <div class="h-4 w-48 animate-pulse rounded bg-white/5" />
      <div class="h-2 w-full animate-pulse rounded bg-white/5" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onImgError } from '@/utils/helpers'
import type { GamificationProfile } from '@/services/api/gamification'

const props = defineProps<{
  profile: GamificationProfile | null
}>()

const progressPercent = computed(() => {
  if (!props.profile) return 0
  const current = props.profile.current_xp
  const next = props.profile.next_level_xp
  const prev = props.profile.level <= 1 ? 0 : xpForLevel(props.profile.level - 1)
  const range = next - prev
  if (range <= 0) return 100
  return Math.min(100, ((current - prev) / range) * 100)
})

const xpToNextLevel = computed(() => {
  if (!props.profile) return 0
  return Math.max(0, props.profile.next_level_xp - props.profile.current_xp)
})

function xpForLevel(level: number): number {
  return Math.floor(Math.pow(level, 2.5) * 100)
}
</script>
