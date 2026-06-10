<template>
  <div class="glass-strong rounded-2xl p-6 md:p-8">
    <div class="mb-6 flex items-center justify-between">
      <h3 class="text-lg font-bold text-white">
        Badges ({{ userBadges.length }}/{{ allBadges.length }})
      </h3>
    </div>

    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
      <div
        v-for="badge in allBadges"
        :key="badge.id"
        class="group relative flex flex-col items-center rounded-xl p-4 text-center transition-all duration-300"
        :class="earnedBadgeMap[badge.id] ? 'bg-white/5' : 'bg-white/[0.02] opacity-40 saturate-0'"
      >
        <!-- Badge icon -->
        <div
          class="mb-2 flex h-14 w-14 items-center justify-center rounded-xl text-2xl transition-transform duration-300 group-hover:scale-110"
          :class="rarityClass(badge.rarity)"
        >
          <img
            v-if="badge.icon_url"
            :src="badge.icon_url"
            :alt="badge.name"
            class="h-10 w-10 object-contain"
            @error="onImgError"
          />
          <span v-else>🏅</span>
        </div>

        <!-- Name -->
        <p class="text-xs font-semibold text-white/80">{{ badge.name }}</p>
        <p class="mt-0.5 text-[10px] leading-tight text-white/30">{{ badge.description }}</p>

        <!-- Rarity indicator -->
        <span
          class="mt-1.5 rounded-full px-2 py-0.5 text-[9px] font-medium tracking-wider uppercase"
          :class="rarityBadge(badge.rarity)"
        >
          {{ badge.rarity }}
        </span>

        <!-- Earned stamp overlay -->
        <div
          v-if="earnedBadgeMap[badge.id]"
          class="absolute top-1.5 right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-[#1db954] text-[10px] text-black shadow-lg"
        >
          ✓
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onImgError } from '@/utils/helpers'
import type { Badge, UserBadge } from '@/services/api/gamification'

const props = defineProps<{
  allBadges: Badge[]
  userBadges: UserBadge[]
}>()

const earnedBadgeMap = computed(() => {
  const map: Record<string, boolean> = {}
  for (const ub of props.userBadges) {
    map[ub.badge_id] = true
  }
  return map
})

const rarityColors: Record<string, string> = {
  common: 'bg-white/5 text-white/50',
  uncommon: 'bg-green-500/10 text-green-400',
  rare: 'bg-blue-500/10 text-blue-400',
  epic: 'bg-purple-500/10 text-purple-400',
  legendary: 'bg-amber-500/10 text-amber-400',
  special: 'bg-pink-500/10 text-pink-400',
}

function rarityClass(rarity: string): string {
  return rarityColors[rarity] ?? rarityColors.common!
}

const rarityBadgeColors: Record<string, string> = {
  common: 'bg-white/5 text-white/30',
  uncommon: 'bg-green-500/10 text-green-400',
  rare: 'bg-blue-500/10 text-blue-400',
  epic: 'bg-purple-500/10 text-purple-300',
  legendary: 'bg-amber-500/10 text-amber-300',
  special: 'bg-pink-500/10 text-pink-300',
}

function rarityBadge(rarity: string): string {
  return rarityBadgeColors[rarity] ?? rarityBadgeColors.common!
}
</script>
