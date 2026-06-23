<template>
  <button
    type="button"
    class="group relative flex h-16 w-full cursor-pointer items-center gap-3 overflow-hidden rounded-xl border border-white/[0.06] bg-white/[0.04] p-2 text-right transition-all duration-200 hover:bg-white/[0.08] focus-visible:ring-2 focus-visible:ring-[#1db954] focus-visible:ring-offset-2 focus-visible:outline-none"
    :class="{ 'border-l-[3px] border-l-[#1db954]': isPlaying }"
    :style="{ transitionDelay: `${delay}ms` }"
    @click="$emit('play', item)"
  >
    <div class="h-12 w-12 shrink-0 overflow-hidden rounded-lg bg-white/10">
      <img
        v-if="item.cover_url || item.coverUrl"
        :src="coverSrc"
        :alt="item.title || item.track_title || ''"
        class="h-full w-full object-cover"
        loading="lazy"
      />
      <div v-else class="flex h-full items-center justify-center">
        <i aria-hidden="true" class="pi pi-music text-lg text-white/30" />
      </div>
    </div>
    <div class="min-w-0 flex-1">
      <p
        class="truncate text-sm font-semibold"
        :class="isPlaying ? 'text-[#1db954]' : 'text-white'"
      >
        {{ item.title || item.track_title || 'بدون عنوان' }}
      </p>
      <p class="truncate text-xs text-white/60">
        {{ item.artist_name || item.artistName || '' }}
      </p>
    </div>
    <div
      class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[#1db954] text-black opacity-0 transition-all duration-200 group-hover:opacity-100"
    >
      <i aria-hidden="true" class="pi pi-play-fill text-sm" />
    </div>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface QuickPlayItem {
  cover_url?: string | null
  coverUrl?: string | null
  title?: string | null
  track_title?: string | null
  artist_name?: string | null
  artistName?: string | null
}

const props = defineProps<{
  item: QuickPlayItem
  isPlaying?: boolean
  delay?: number
}>()

defineEmits<{
  play: [item: QuickPlayItem]
}>()

const coverSrc = computed((): string | undefined => (props.item.cover_url || props.item.coverUrl) ?? undefined)
</script>
