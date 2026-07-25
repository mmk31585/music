<template>
  <section v-if="items.length" class="mt-12">
    <div class="mb-5">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Community</p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">What friends are listening</h2>
    </div>
    <div class="flex gap-4 overflow-x-auto pb-2 scrollbar-none">
      <div
        v-for="(item, i) in items.slice(0, 6)"
        :key="i"
        class="w-44 shrink-0 rounded-2xl border border-white/[4%] bg-white/[2%] p-4 backdrop-blur-xs"
        @contextmenu.prevent="openContextMenu($event, item, i)"
      >
        <div class="flex items-center gap-2">
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-linear-to-br from-spotify to-aurora-blue text-[11px] font-bold text-white">
            {{ getInitials(item) }}
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-semibold text-white">{{ item.username || item.user?.username || 'User' }}</p>
            <p class="text-[11px] text-white/50">{{ item.timeAgo || 'Now' }}</p>
          </div>
        </div>
        <div class="mt-3 flex justify-center">
          <div class="h-16 w-16 overflow-hidden rounded-xl bg-white/10 ring-1 ring-white/10">
            <img
              v-if="item.cover_url || item.track?.cover_url"
              :src="item.cover_url || item.track?.cover_url || ''"
              alt=""
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div v-else class="flex h-full items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-white/30"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
            </div>
          </div>
        </div>
        <p class="mt-2 truncate text-center text-sm font-semibold text-white">{{ item.track_title || item.track?.title || 'Track' }}</p>
        <p class="truncate text-center text-xs text-white/50">{{ item.artist_name || item.track?.artist_name || '' }}</p>
      </div>
    </div>
    <ContextMenu
      v-model:visible="menuVisible"
      :sections="sections"
      :header="header"
      :accent-color="accentColor"
      :position="{ x: menuX, y: menuY }"
    />
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { TrackContextItem } from '@/composables/useTrackContextMenu'
import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
import ContextMenu from '@/components/common/ContextMenu.vue'

interface ActivityItem {
  id?: string | number
  username?: string
  timeAgo?: string
  cover_url?: string | null
  track_title?: string | null
  artist_name?: string | null
  user?: {
    username?: string
  }
  track?: {
    id?: string | number
    cover_url?: string | null
    title?: string | null
    artist_name?: string | null
  }
}

const props = defineProps<{
  items: ActivityItem[]
}>()

function getInitials(item: ActivityItem) {
  const name = item.username || item.user?.username || '?'
  return name.slice(0, 2).toUpperCase()
}

// ── Context menu ──────────────────────────────────────────────────
const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTrack = ref<TrackContextItem | null>(null)

function openContextMenu(e: MouseEvent, item: ActivityItem, index: number) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  // Map ActivityItem to TrackContextItem shape
  contextTrack.value = {
    id: item.track?.id ?? item.id ?? index,
    title: item.track_title ?? item.track?.title ?? null,
    artist_name: item.artist_name ?? item.track?.artist_name ?? null,
    cover_url: item.cover_url ?? item.track?.cover_url ?? null,
  }
  menuVisible.value = true
}

const { sections, header, accentColor } = useTrackContextMenu(
  computed(() => contextTrack.value),
)
</script>

<style scoped>
.scrollbar-none {
  scrollbar-width: none;
}
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
</style>
