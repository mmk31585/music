<template>
  <div
    class="group w-28 shrink-0 cursor-pointer space-y-2 sm:w-32 md:w-36"
    :style="{ transitionDelay: `${delay}ms` }"
    @click="$emit('play', item)"
  >
    <div
      class="relative aspect-square overflow-hidden rounded-xl bg-white/[0.06] ring-1 ring-white/10 transition-all duration-300 group-hover:shadow-[0_8px_32px_rgba(0,0,0,0.4)] group-hover:ring-[#1db954]/40"
      :style="{ borderRadius: radius + 'px' }"
    >
      <img
        v-if="item.cover_url || item.coverUrl"
        :src="item.cover_url || item.coverUrl"
        :alt="altText"
        class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
        loading="lazy"
      />
      <div v-else class="flex h-full items-center justify-center">
        <i class="pi pi-music text-2xl text-white/30" />
      </div>
      <div
        class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl transition-transform group-hover:scale-110"
        >
          <i class="pi pi-play-fill text-sm" />
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

const props = withDefaults(defineProps<{
  item: Record<string, any>
  isPlaying?: boolean
  delay?: number
  radius?: number
  badge?: string
  reason?: string
}>(), {
  radius: 12,
})

defineEmits<{
  play: [item: Record<string, any>]
}>()

const artistName = computed(() => {
  return props.item.artist_name || props.item.artistName || ''
})

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
