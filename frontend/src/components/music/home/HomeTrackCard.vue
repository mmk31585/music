<template>
  <div
    role="button"
    tabindex="0"
    class="group w-44 shrink-0 cursor-pointer space-y-2"
    :style="{ transitionDelay: `${delay}ms` }"
    @click="$emit('play', item)"
    @keydown.enter="$emit('play', item)"
    @keydown.space.prevent="$emit('play', item)"
  >
    <div
      class="relative aspect-square overflow-hidden rounded-xl bg-white/[0.06] ring-1 ring-white/10 transition-all duration-300 hover:-translate-y-0.5 hover:shadow-[0_8px_32px_rgba(0,0,0,0.5)] hover:ring-[#1db954]/40"
      :style="{ borderRadius: radius + 'px' }"
    >
      <img
        v-if="item.cover_url || item.coverUrl"
        :src="coverSrc"
        :alt="altText"
        class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
        loading="lazy"
      />
      <div v-else class="flex h-full items-center justify-center">
        <i aria-hidden="true" class="pi pi-music text-2xl text-white/30" />
      </div>
      <div
        class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl transition-transform group-hover:scale-110"
        >
          <i aria-hidden="true" class="pi pi-play-fill text-sm" />
        </div>
      </div>
      <div
        v-if="isPlaying"
        class="absolute right-2 bottom-2 flex h-5 items-end gap-0.5 rounded-full bg-black/60 px-1.5 py-1"
      >
        <span class="equalizer-bar-small h-2 w-0.5 rounded-full bg-[#1db954]" />
        <span class="equalizer-bar-small animation-delay-150 h-3 w-0.5 rounded-full bg-[#1db954]" />
        <span class="equalizer-bar-small animation-delay-300 h-2.5 w-0.5 rounded-full bg-[#1db954]" />
      </div>
      <div
        v-if="badge"
        class="absolute top-2 left-2 rounded-full bg-[#1db954]/90 px-2 py-0.5 text-[10px] font-bold text-black"
      >
        {{ badge }}
      </div>
      <div
        v-if="reason"
        class="absolute top-2 left-2 rounded-full bg-white/10 px-2 py-0.5 text-[10px] font-medium text-white backdrop-blur"
      >
        {{ reason }}
      </div>
    </div>
    <div class="space-y-0.5 px-0.5">
      <p class="truncate text-sm font-semibold text-white">
        {{ item.title || item.track_title || 'بدون عنوان' }}
      </p>
      <p class="truncate text-xs text-white/60">
        {{ artistName }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface TrackCardItem {
  cover_url?: string | null
  coverUrl?: string | null
  title?: string | null
  track_title?: string | null
  artist_name?: string | null
  artistName?: string | null
}

const props = withDefaults(defineProps<{
  item: TrackCardItem
  isPlaying?: boolean
  delay?: number
  radius?: number
  badge?: string
  reason?: string
}>(), {
  radius: 12,
})

defineEmits<{
  play: [item: TrackCardItem]
}>()

const artistName = computed(() => {
  return props.item.artist_name || props.item.artistName || ''
})

const coverSrc = computed((): string | undefined => (props.item.cover_url || props.item.coverUrl) ?? undefined)

const altText = computed(() => {
  const title = props.item.title || props.item.track_title || ''
  const artist = artistName.value
  return artist ? `${title} - ${artist}` : title
})
</script>

<style scoped>
.equalizer-bar-small {
  animation: eqBeat 600ms ease-in-out infinite alternate;
  transform-origin: bottom;
}
.animation-delay-150 { animation-delay: 150ms; }
.animation-delay-300 { animation-delay: 300ms; }
@keyframes eqBeat {
  from { transform: scaleY(0.4); opacity: 0.5; }
  to   { transform: scaleY(1);   opacity: 1; }
}
</style>
