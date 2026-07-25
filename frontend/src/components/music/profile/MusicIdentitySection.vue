<template>
  <div v-if="hasData" class="rounded-2xl bg-surface-overlay/50 ring-1 ring-border-subtle p-5">
    <h3 class="text-xs uppercase tracking-widest text-muted mb-4">سلیقه موسیقی</h3>

    <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
      <!-- Top Genres -->
      <div v-if="topGenres && topGenres.length > 0">
        <h4 class="text-[10px] uppercase tracking-wider text-muted mb-3">ژانرهای برتر</h4>
        <div class="space-y-2.5">
          <div v-for="g in topGenres" :key="g.name" class="flex items-center gap-3">
            <span class="text-xs text-secondary w-20 truncate shrink-0">{{ g.name }}</span>
            <div class="flex-1 h-1.5 rounded-full bg-surface-active/80 overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-700"
                :style="{ width: g.percent + '%', background: g.color }"
              />
            </div>
            <span class="text-[10px] text-muted w-8 text-right tabular-nums shrink-0">{{ g.percent }}%</span>
          </div>
        </div>
      </div>

      <!-- Top Artists -->
      <div v-if="topArtists && topArtists.length > 0">
        <h4 class="text-[10px] uppercase tracking-wider text-muted mb-3">هنرمندان برتر</h4>
        <div class="space-y-2">
          <div
            v-for="a in topArtists.slice(0, 5)"
            :key="a.id"
            class="flex items-center gap-3 group cursor-pointer"
            @click="navigateToArtist(a.id)"
          >
            <div class="h-9 w-9 shrink-0 overflow-hidden rounded-full ring-2 ring-border-default">
              <img
                v-if="a.avatarUrl"
                :src="a.avatarUrl"
                :alt="a.name"
                class="h-full w-full object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-full w-full items-center justify-center bg-surface-active/80 text-xs font-bold text-tertiary">
                {{ a.name.charAt(0).toUpperCase() }}
              </div>
            </div>
            <span class="text-sm text-secondary truncate group-hover:text-primary transition">{{ a.name }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Streak badge -->
    <div v-if="listeningStreak && listeningStreak >= 3" class="mt-4">
      <span class="inline-flex items-center gap-1.5 rounded-full bg-orange-500/15 border border-orange-500/20 px-3 py-1 text-xs text-orange-300">
        {{ listeningStreak }} روز پشت سر هم
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{
  topGenres?: Array<{ name: string; percent: number; color: string }>
  topArtists?: Array<{ id: string; name: string; avatarUrl?: string }>
  listeningStreak?: number
  recentMood?: string
}>()

const hasData = computed(() =>
  (props.topGenres && props.topGenres.length > 0) ||
  (props.topArtists && props.topArtists.length > 0)
)

const router = useRouter()

function navigateToArtist(id: string) {
  router.push(`/artists/${id}`)
}
</script>
