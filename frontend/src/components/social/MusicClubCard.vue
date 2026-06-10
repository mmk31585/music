<template>
  <div
    class="group overflow-hidden rounded-xl border border-white/5 bg-white/[0.03] transition-all duration-300 hover:border-white/10 hover:bg-white/[0.06]"
  >
    <div
      v-if="club.cover_url"
      class="h-24 bg-cover bg-center"
      :style="{ backgroundImage: `url(${club.cover_url})` }"
    />
    <div class="p-5">
      <div class="flex items-start gap-3">
        <div
          v-if="!club.cover_url"
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-purple-500/10 text-lg"
        >
          🏛
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold text-white">{{ club.name }}</p>
          <p v-if="creatorName" class="mt-0.5 text-xs text-white/30">{{ creatorName }}</p>
          <p v-if="club.description" class="mt-1 line-clamp-2 text-xs text-white/40">
            {{ club.description }}
          </p>
          <div class="mt-2 flex items-center gap-3 text-[11px] text-white/30">
            <span class="flex items-center gap-1">
              👥 {{ club.member_count }} / {{ club.max_members }}
            </span>
            <span
              class="rounded-full px-2 py-0.5 text-[10px] font-medium"
              :class="club.is_public ? 'bg-green-500/10 text-green-400' : 'bg-amber-500/10 text-amber-400'"
            >
              {{ club.is_public ? 'Public' : 'Private' }}
            </span>
          </div>
          <!-- Member bar -->
          <div class="mt-2 h-1 w-full overflow-hidden rounded-full bg-white/5">
            <div
              class="h-full rounded-full bg-purple-500/40 transition-all"
              :style="{ width: `${Math.min(100, (club.member_count / club.max_members) * 100)}%` }"
            />
          </div>
        </div>
      </div>

      <div class="mt-3 flex items-center gap-2">
        <button
          class="flex-1 rounded-lg bg-purple-500/10 py-2 text-xs font-semibold text-purple-400 transition hover:bg-purple-500/20"
          @click="$emit('join', club.id)"
        >
          Join Club
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MusicClub } from '@/services/api/social'

defineProps<{
  club: MusicClub
  creatorName?: string
}>()

defineEmits<{
  join: [id: string]
}>()
</script>
