<template>
  <RouterLink :to="`/artist/${artist.id}`" class="group block w-40 shrink-0 space-y-3">
    <div
      class="mx-auto h-36 w-36 overflow-hidden rounded-full bg-white/[0.06] ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50"
    >
      <img
        v-if="artist.image_url"
        :src="artist.image_url"
        :alt="artist.name"
        loading="lazy"
        class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
        @error="onImgError"
      />
      <div v-else class="flex h-full items-center justify-center">
        <i aria-hidden="true" class="pi pi-user text-3xl text-slate-500" />
      </div>
    </div>
    <div class="space-y-0.5 text-center">
      <p class="truncate text-sm font-bold text-white">
        {{ artist.name }}
      </p>
      <p class="text-xs text-slate-500">Artist</p>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { onImgError } from '@/utils/helpers'
import type { Artist } from '@/services/api/catalog/artists'

defineProps<{
  artist: Artist
}>()
</script>
