<template>
  <section class="mt-12">
    <div class="mb-5">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Persian scene</p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">Highlights from Iran</h2>
    </div>
    <div class="flex gap-4 overflow-x-auto pb-2 scrollbar-none">
      <RouterLink
        v-for="(item, i) in items"
        :key="item.id"
        :to="`/album/${item.id}`"
        class="group w-40 shrink-0 space-y-2"
      >
        <div class="relative aspect-square overflow-hidden rounded-xl bg-white/6 ring-1 ring-white/10 transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg hover:ring-rose-500/30">
          <img
            v-if="item.cover_url"
            :src="item.cover_url"
            :alt="item.title || ''"
            class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
            loading="lazy"
          />
          <div v-else class="flex h-full items-center justify-center bg-linear-to-br from-rose-500/20 to-amber-500/20">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-8 w-8 text-rose-300/40"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
          </div>
          <div class="absolute right-2 bottom-2 rounded-full bg-rose-500/70 px-2 py-0.5 text-[10px] font-bold text-white backdrop-blur-xs">
            {{ tags[i % tags.length] }}
          </div>
        </div>
        <p class="truncate text-sm font-semibold text-white">{{ item.title || 'Untitled' }}</p>
        <p class="truncate text-xs text-white/50">{{ item.artist_name || 'Unknown artist' }}</p>
      </RouterLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Album } from '@/services/api/catalog/albums/types'

defineProps<{
  items: Album[]
}>()

const tags = ['تازه', 'محبوب', 'پیشنهاد', 'ویژه']
</script>

<style scoped>
.scrollbar-none {
  scrollbar-width: none;
}
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
</style>
