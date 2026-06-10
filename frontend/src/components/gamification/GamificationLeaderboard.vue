<template>
  <div class="glass-strong rounded-2xl p-6 md:p-8">
    <div class="mb-6 flex items-center justify-between">
      <h3 class="text-lg font-bold text-white">Leaderboard</h3>

      <div class="flex gap-1 rounded-lg bg-white/5 p-0.5">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="rounded-md px-3 py-1 text-xs font-medium transition-all duration-200"
          :class="
            activeTab === tab.key ? 'bg-[#1db954] text-black' : 'text-white/40 hover:text-white/60'
          "
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div v-if="!filteredEntries.length" class="py-8 text-center text-sm text-white/30">
      No rankings available yet.
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="(entry, i) in filteredEntries"
        :key="entry.user_id"
        class="group flex items-center gap-3 rounded-xl px-4 py-3 transition-all duration-200"
        :class="entry.rank <= 3 ? 'bg-white/5' : 'hover:bg-white/[0.03]'"
      >
        <!-- Rank -->
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center font-bold"
          :class="rankClass(entry.rank)"
        >
          <span v-if="entry.rank <= 3" class="text-sm">{{ rankEmoji(entry.rank) }}</span>
          <span v-else class="text-xs text-white/30">{{ entry.rank }}</span>
        </div>

        <!-- Avatar -->
        <div class="h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="entry.avatar_url"
            :src="entry.avatar_url"
            :alt="entry.username"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div
            v-else
            class="flex h-full w-full items-center justify-center text-xs font-bold text-white/30"
          >
            {{ entry.username.charAt(0).toUpperCase() }}
          </div>
        </div>

        <!-- Name -->
        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold text-white">{{ entry.username }}</p>
        </div>

        <!-- Score -->
        <div class="shrink-0 text-right">
          <p class="text-sm font-bold text-white/60">{{ entry.score.toLocaleString() }}</p>
          <p class="text-[10px] text-white/20">XP</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onImgError } from '@/utils/helpers'
import type { LeaderboardEntry } from '@/services/api/gamification'

const props = defineProps<{
  entries: LeaderboardEntry[]
}>()

const tabs = [
  { key: 'all', label: 'All Time' },
  { key: 'monthly', label: 'Monthly' },
  { key: 'weekly', label: 'Weekly' },
  { key: 'daily', label: 'Daily' },
]

const activeTab = ref('all')

const filteredEntries = computed(() => {
  return props.entries
})

function onTabChange(key: string) {
  activeTab.value = key
}

function rankClass(rank: number): string {
  if (rank === 1) return 'text-amber-400'
  if (rank === 2) return 'text-slate-300'
  if (rank === 3) return 'text-amber-600'
  return 'text-white/20'
}

function rankEmoji(rank: number): string {
  if (rank === 1) return '🥇'
  if (rank === 2) return '🥈'
  if (rank === 3) return '🥉'
  return ''
}
</script>
