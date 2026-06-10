<template>
  <div class="glass-strong rounded-2xl p-6 md:p-8">
    <h3 class="mb-6 text-lg font-bold text-white">Top Contributors</h3>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 5" :key="i" class="flex items-center gap-3">
        <div class="h-8 w-8 animate-pulse rounded-full bg-white/5" />
        <div class="h-8 w-8 animate-pulse rounded-full bg-white/5" />
        <div class="flex-1 space-y-1">
          <div class="h-4 w-24 animate-pulse rounded bg-white/5" />
          <div class="h-3 w-16 animate-pulse rounded bg-white/5" />
        </div>
      </div>
    </div>

    <div
      v-else-if="contributors.length === 0"
      class="flex flex-col items-center gap-3 py-8 text-center"
    >
      <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-white/5">
        <i class="pi pi-trophy text-xl text-white/15" />
      </div>
      <p class="text-sm text-white/25">No contributors yet. Be the first!</p>
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="(c, idx) in contributors"
        :key="c.user_id"
        class="group spring flex items-center gap-3 rounded-xl px-3 py-2.5 transition-all hover:bg-white/[0.04]"
      >
        <!-- Rank badge -->
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-bold"
          :class="rankClass(idx)"
        >
          <template v-if="idx === 0">
            <i class="pi pi-star-fill text-sm text-[#f59e0b]" />
          </template>
          <template v-else-if="idx === 1">
            <i class="pi pi-star-fill text-sm text-[#94a3b8]" />
          </template>
          <template v-else-if="idx === 2">
            <i class="pi pi-star-fill text-sm text-[#cd7f32]" />
          </template>
          <template v-else>
            {{ c.rank }}
          </template>
        </div>

        <!-- Avatar -->
        <div class="h-9 w-9 shrink-0 overflow-hidden rounded-full bg-white/10">
          <img
            v-if="c.avatar_url"
            :src="c.avatar_url"
            :alt="c.username"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div
            v-else
            class="flex h-full w-full items-center justify-center text-xs font-bold text-white/30"
          >
            {{ c.username.charAt(0).toUpperCase() }}
          </div>
        </div>

        <!-- Info -->
        <div class="min-w-0 flex-1">
          <p
            class="truncate text-sm font-medium text-white/80 transition-colors group-hover:text-white"
          >
            {{ c.username }}
          </p>
          <p class="text-xs text-white/30">
            {{ c.approved_count }} approved · {{ c.total_contributions }} total
          </p>
        </div>

        <!-- XP -->
        <div class="text-right">
          <p class="text-sm font-bold text-[#1db954]">{{ c.xp_earned.toLocaleString() }}</p>
          <p class="text-[10px] text-white/20">XP</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onImgError } from '@/utils/helpers'
import type { ContributorStats } from '@/services/api/contribution'

defineProps<{
  contributors: ContributorStats[]
  loading?: boolean
}>()

function rankClass(idx: number) {
  if (idx === 0) return 'bg-[#f59e0b]/15'
  if (idx === 1) return 'bg-[#94a3b8]/15'
  if (idx === 2) return 'bg-[#cd7f32]/15'
  return 'bg-white/5 text-white/30'
}
</script>
