<template>
  <div
    class="relative h-full w-full overflow-hidden rounded-2xl bg-gradient-to-br from-[#1db954]/20 to-[#a855f7]/20"
  >
    <!-- 2×2 Grid -->
    <div
      v-if="covers.length > 0"
      class="grid h-full w-full"
      :class="gridClass"
    >
      <div
        v-for="(cover, index) in gridCovers"
        :key="index"
        class="relative overflow-hidden"
        :class="cellClass(index)"
      >
        <img
          v-if="cover"
          :src="cover"
          :alt="`Track ${index + 1}`"
          class="h-full w-full object-cover transition duration-500 hover:scale-110"
          loading="lazy"
          @error="onCoverError"
        />
        <div
          v-else
          class="flex h-full items-center justify-center bg-white/[0.03]"
        >
          <i aria-hidden="true" class="pi pi-music text-lg text-white/20" />
        </div>

        <!-- Overlays for partial grids -->
        <div
          v-if="index === 0 && covers.length === 1"
          class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent"
        />
        <div
          v-if="index === 1 && covers.length === 2"
          class="pointer-events-none absolute inset-0 bg-gradient-to-l from-black/40 to-transparent"
        />
      </div>
    </div>

    <!-- Track count badge -->
    <div
      v-if="trackCount !== undefined"
      class="absolute right-2 bottom-2 z-10 rounded-full bg-black/70 px-2 py-0.5 text-[10px] font-bold text-white/80 backdrop-blur-sm"
    >
      {{ trackCount }} {{ trackCount === 1 ? 'track' : 'tracks' }}
    </div>

    <!-- Empty state when no tracks -->
    <div
      v-if="covers.length === 0"
      class="flex h-full w-full flex-col items-center justify-center gap-2"
    >
      <i aria-hidden="true" class="pi pi-list text-3xl text-white/30" />
      <span class="text-xs font-medium text-white/40">No tracks</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    covers: (string | null | undefined)[]
    trackCount?: number
  }>(),
  {
    trackCount: undefined,
  },
)

const gridCovers = computed(() => {
  const valid = props.covers.filter(Boolean) as string[]
  if (valid.length === 0) return []
  if (valid.length === 1) return [valid[0], null, null, null]
  if (valid.length === 2) return [valid[0], valid[1], null, null]
  if (valid.length === 3) return [valid[0], valid[1], valid[2], null]
  return valid.slice(0, 4)
})

const gridClass = computed(() => {
  const count = gridCovers.value.filter((c) => c !== null).length
  if (count <= 1) return 'grid-cols-1 grid-rows-1'
  if (count === 2) return 'grid-cols-2 grid-rows-1'
  return 'grid-cols-2 grid-rows-2'
})

function cellClass(index: number): string {
  const count = gridCovers.value.filter((c) => c !== null).length
  if (count === 1) return 'rounded-2xl'
  if (count === 2) {
    if (index === 0) return 'rounded-l-2xl'
    if (index === 1) return 'rounded-r-2xl'
  }
  if (count === 3) {
    if (index === 0) return 'rounded-tl-2xl'
    if (index === 1) return 'rounded-tr-2xl'
    if (index === 2) return 'rounded-bl-2xl'
  }
  // 4 covers or full grid
  if (index === 0) return 'rounded-tl-2xl'
  if (index === 1) return 'rounded-tr-2xl'
  if (index === 2) return 'rounded-bl-2xl'
  if (index === 3) return 'rounded-br-2xl'
  return ''
}

function onCoverError(e: Event) {
  const img = e.currentTarget as HTMLImageElement
  img.style.display = 'none'
}
</script>
