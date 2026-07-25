<template>
  <RouterLink :to="`/artist/${artist.id}`" class="group block w-40 shrink-0 space-y-3">
    <div
      class="mx-auto h-36 w-36 overflow-hidden rounded-full bg-surface-overlay ring-1 ring-border-default transition group-hover:ring-accent/50"
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
        <User aria-hidden="true" class="text-3xl text-tertiary"  />
      </div>
    </div>
    <div class="space-y-0.5 text-center">
      <p class="truncate text-sm font-bold text-primary">
        {{ artist.name }}
      </p>
      <p class="text-xs text-tertiary">Artist</p>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { User } from 'lucide-vue-next'
import { onImgError } from '@/utils/helpers'
import type { Artist } from '@/services/api/catalog/artists'

defineProps<{
  artist: Artist
}>()
</script>
