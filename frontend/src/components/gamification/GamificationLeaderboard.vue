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
            activeTab === tab.key ? 'bg-spotify text-black' : 'text-white/40 hover:text-white/60'
          "
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- User rank (shown when user is not in top entries) -->
    <div
      v-if="userRank && !userInEntries"
      class="mb-6 rounded-2xl border border-spotify/20 bg-spotify/5 px-5 py-4"
    >
      <p class="text-sm text-white/60">Your Rank</p>
      <p class="text-3xl font-black text-white">#{{ userRank }}</p>
    </div>

    <!-- Friends empty state guidance -->
    <AppEmptyState
      v-if="activeTab === 'friends' && !props.entries.length"
      icon="pi pi-users"
      title="No friends yet"
      description="Follow creators and friends to see their rankings here"
    >
      <RouterLink
        to="/search"
        class="mt-3 inline-block rounded-full bg-spotify px-5 py-2 text-sm font-bold text-black transition hover:bg-spotify-hover"
      >
        Discover Artists
      </RouterLink>
    </AppEmptyState>

    <template v-else-if="!props.entries.length">
      <div class="py-8 text-center text-sm text-white/30">
        No rankings available yet.
      </div>
    </template>

    <div v-else class="space-y-2">
      <div
        v-for="entry in props.entries"
        :key="entry.user_id"
        class="group flex items-center gap-3 rounded-xl px-4 py-3 transition-all duration-200"
        :class="entry.rank <= 3 ? 'bg-white/5' : 'hover:bg-white/3'"
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
import { ref, computed, watch } from 'vue'
import { onImgError } from '@/utils/helpers'
import type { LeaderboardEntry } from '@/services/api/gamification'

const props = defineProps<{
  entries: LeaderboardEntry[]
  type?: string
  userRank?: number | null
  userId?: string | null
  onTypeChange?: (type: string) => void
}>()

const tabs = [
  { key: 'all', label: 'All Time' },
  { key: 'monthly', label: 'Monthly' },
  { key: 'weekly', label: 'Weekly' },
  { key: 'streams', label: 'Streams' },
  { key: 'friends', label: 'Friends' },
]

const activeTab = ref(props.type ?? 'all')

const userInEntries = computed(() => {
  if (!props.userId) return false
  return props.entries.some((e) => e.user_id === props.userId)
})

watch(() => props.type, (val) => {
  if (val) activeTab.value = val
})

watch(activeTab, (tab) => {
  props.onTypeChange?.(tab)
})

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
