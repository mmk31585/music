<template>
  <div
    class="group overflow-hidden rounded-2xl bg-surface-overlay/60 ring-1 ring-border-default transition-all duration-300
           hover:bg-surface-active/80 hover:ring-border-strong focus-within:ring-2 focus-within:ring-spotify"
    role="article"
    :aria-label="`Club: ${club.name}`"
  >
    <!-- Cover image header -->
    <div
      v-if="club.cover_url"
      class="h-28 bg-cover bg-center transition-transform duration-500 group-hover:scale-105"
      :style="{ backgroundImage: `url(${club.cover_url})` }"
    />
    <div
      v-else
      class="flex h-20 items-center justify-center bg-linear-to-br from-purple-500/10 to-purple-500/5"
    >
      <Building2 aria-hidden="true" class="text-2xl text-accent/30"  />
    </div>

    <div class="p-5">
      <div class="flex items-start gap-3">
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-bold text-primary">{{ club.name }}</h3>
          <p v-if="creatorName" class="mt-0.5 text-xs text-primary/30">{{ creatorName }}</p>
          <p v-if="club.description" class="mt-1 line-clamp-2 text-xs text-primary/40 leading-relaxed">
            {{ club.description }}
          </p>

          <!-- Meta row -->
          <div class="mt-3 flex items-center gap-3 text-[11px] text-primary/30">
            <span class="flex items-center gap-1">
              <Users aria-hidden="true" class="text-[10px]"  />
              {{ club.member_count || 0 }} / {{ club.max_members || '∞' }}
            </span>
            <span
              class="rounded-full px-2 py-0.5 text-[10px] font-medium"
              :class="club.is_public ? 'bg-green-500/10 text-green-400' : 'bg-amber-500/10 text-amber-400'"
            >
              {{ club.is_public ? 'Public' : 'Private' }}
            </span>
          </div>

          <!-- Member capacity bar -->
          <div class="mt-2 h-1 w-full overflow-hidden rounded-full bg-surface-overlay">
            <div
              class="h-full rounded-full bg-purple-500/40 transition-all duration-500"
              :style="{ width: `${capacityPercent}%` }"
            />
          </div>

          <!-- Recent member avatars -->
          <div v-if="clubAvatars.length" class="mt-3 flex -space-x-1.5">
            <div
              v-for="(avatar, i) in clubAvatars.slice(0, 5)"
              :key="i"
              class="h-5 w-5 overflow-hidden rounded-full border-2 border-surface-base"
            >
              <img :src="avatar" alt="" class="h-full w-full object-cover" />
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4 flex items-center gap-2">
        <button
          aria-label="Join club"
          class="flex-1 rounded-lg bg-accent-subtle py-2.5 text-xs font-bold text-accent transition
                 hover:bg-purple-500/20 active:scale-[0.98] focus-visible:outline-2 focus-visible:outline-purple-400"
          @click="$emit('join', club.id)"
        >
          <span class="flex items-center justify-center gap-1.5">
            <Plus aria-hidden="true" class="text-[10px]"  />
            Join Club
          </span>
        </button>
        <button
          aria-label="Share club"
          class="flex h-9 w-9 items-center justify-center rounded-lg text-primary/30 transition
                 hover:bg-surface-overlay/60 hover:text-primary/60 focus-visible:outline-2 focus-visible:outline-accent"
          @click="$emit('share', club.id)"
        >
          <Share2 aria-hidden="true" class="text-sm"  />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Building2, Plus, Share2, Users } from 'lucide-vue-next'
import { computed } from 'vue'
import type { MusicClub } from '@/services/api/social'

const props = defineProps<{
  club: MusicClub
  creatorName?: string
}>()

defineEmits<{
  join: [id: string]
  share: [id: string]
}>()

const capacityPercent = computed(() => {
  if (!props.club.max_members) return 0
  return Math.min(100, ((props.club.member_count || 0) / props.club.max_members) * 100)
})

const clubAvatars = computed(() => {
  return (props.club as any).avatars || []
})
</script>
