<template>
  <RouterLink :to="`/album/${album.id}`" class="group block w-44 shrink-0 space-y-3">
    <div
      class="relative aspect-square overflow-hidden rounded-2xl bg-surface-overlay shadow-lg ring-1 ring-border-default transition group-hover:ring-accent/50"
    >
      <img
        v-if="album.cover_url"
        :src="album.cover_url"
        :alt="album.title"
        loading="lazy"
        class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
        @error="onImgError"
      />
      <div v-else class="flex h-full items-center justify-center">
        <Disc3 aria-hidden="true" class="text-3xl text-tertiary"  />
      </div>
      <div
        class="absolute inset-0 flex items-center justify-center bg-bg-overlay/30 opacity-0 transition group-hover:opacity-100"
      >
        <div
          class="flex h-12 w-12 items-center justify-center rounded-full bg-accent text-black shadow-xl"
        >
          <Play aria-hidden="true" class="text-lg"  />
        </div>
      </div>
    </div>
    <div class="space-y-0.5 px-1">
      <p class="truncate text-sm font-bold text-primary">
        {{ album.title }}
      </p>
      <p v-if="album.artist_name" class="truncate text-xs text-secondary">
        {{ album.artist_name }}
      </p>
      <p class="text-xs text-tertiary">
        {{ album.release_date?.slice(0, 4) || ''
        }}{{ album.track_count ? ` • ${album.track_count} tracks` : '' }}
      </p>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { Disc3, Play } from 'lucide-vue-next'
import { onImgError } from '@/utils/helpers'
import type { Album } from '@/services/api/catalog/albums'

defineProps<{
  album: Album
}>()
</script>
