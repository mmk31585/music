<template>
  <section v-if="artists.length" class="mt-12">
    <div class="mb-5">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Spotlight</p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">Featured Artists</h2>
    </div>
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <RouterLink
        v-for="artist in artists"
        :key="artist.id"
        :to="`/artist/${artist.id}`"
        class="group relative overflow-hidden rounded-2xl border border-white/[4%] bg-white/[1%] p-5 transition-all hover:-translate-y-1 hover:border-white/[10%] hover:bg-white/[4%]"
      >
        <div class="flex flex-col items-center gap-4 text-center">
          <div class="relative h-28 w-28 overflow-hidden rounded-full bg-white/10 ring-2 ring-white/[6%] transition group-hover:ring-spotify/40">
            <img
              v-if="artist.image_url"
              :src="artist.image_url"
              :alt="artist.name"
              class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
              loading="lazy"
            />
            <div v-else class="flex h-full items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-8 w-8 text-white/30"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
            </div>
          </div>
          <div class="space-y-1">
            <p class="truncate text-base font-bold text-white">{{ artist.name }}</p>
            <p v-if="artist.monthly_listeners" class="text-xs text-white/40">
              {{ formatListeners(artist.monthly_listeners) }} monthly listeners
            </p>
          </div>
          <button
            type="button"
            class="mt-1 inline-flex h-9 cursor-pointer items-center gap-2 rounded-full bg-white/10 px-5 text-xs font-semibold text-white/80 backdrop-blur-xs transition hover:bg-white/20 hover:text-white"
            @click.prevent="$emit('playArtist', artist)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3 w-3"><polygon points="5 3 19 12 5 21 5 3"/></svg>
            Play
          </button>
        </div>
      </RouterLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Artist } from '@/services/api/catalog/artists/types'

defineProps<{
  artists: Artist[]
}>()

defineEmits<{
  playArtist: [artist: Artist]
}>()

function formatListeners(count: number): string {
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}
</script>
