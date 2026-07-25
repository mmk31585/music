<template>
  <div class="relative overflow-hidden rounded-2xl">
    <div
      class="absolute inset-0 bg-linear-to-b from-spotify/20 via-surface-raised/60 to-surface-raised"
    />
    <div class="absolute inset-0 bg-linear-to-t from-black/60 via-transparent to-transparent" />

    <div
      class="relative z-10 flex flex-col items-center gap-6 px-6 pt-16 pb-8 text-center md:flex-row md:items-end md:gap-8 md:pt-8 md:text-left"
    >
      <div
        class="h-48 w-48 shrink-0 overflow-hidden rounded-full border-4 border-white/10 shadow-2xl md:h-56 md:w-56"
      >
        <img
          v-if="artist?.image_url"
          :src="artist.image_url"
          :alt="artist?.name"
          loading="lazy"
          class="h-full w-full object-cover"
          @error="onImgError"
        />
        <div
          v-else
          class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-surface-raised text-5xl text-white/40"
        >
          <User aria-hidden="true" class=""  />
        </div>
      </div>

      <div class="flex-1">
        <div
          v-if="artist?.is_verified"
          class="mb-2 flex items-center justify-center gap-1.5 md:justify-start"
        >
          <BadgeCheck aria-hidden="true" class="text-sm text-spotify"  />
          <span class="text-xs font-medium text-spotify">Verified Artist</span>
        </div>

        <h1 class="text-4xl leading-tight font-black text-white md:text-6xl">
          {{ artist?.name || 'Artist' }}
        </h1>

        <div
          class="mt-3 flex flex-wrap items-center justify-center gap-4 text-sm text-slate-400 md:justify-start"
        >
          <span class="font-semibold text-white">
            {{ formatNumber(monthlyListeners) }}
          </span>
          monthly listeners
        </div>

        <div class="mt-6 flex flex-wrap items-center justify-center gap-3 md:justify-start">
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-full bg-spotify px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-spotify-hover"
            @click="$emit('playAll')"
          >
            <Play aria-hidden="true" class=""  />
            Play
          </button>

          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-8 py-3 text-sm font-bold text-white backdrop-blur-xs transition hover:bg-white/15"
            @click="$emit('shuffle')"
          >
            <Shuffle aria-hidden="true" class=""  />
            Shuffle
          </button>

          <button
            type="button"
            :class="[
              'inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-bold transition',
              isFollowing
                ? 'border border-spotify/50 bg-spotify/10 text-spotify hover:bg-spotify/20'
                : 'border border-white/15 bg-white/10 text-white hover:bg-white/15',
            ]"
            @click="$emit('toggleFollow')"
          >
            <component :is="isFollowing ? Check : Plus"<i aria-hidden="true"  class="text-xs" /> />
            {{ isFollowing ? 'Following' : 'Follow' }}
          </button>
        </div>

        <div v-if="artist?.bio" class="mt-4 max-w-lg text-sm leading-relaxed text-slate-400">
          {{ truncateBio(artist.bio) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BadgeCheck, Check, Play, Plus, Shuffle, User } from 'lucide-vue-next'
import { onImgError } from '@/utils/helpers'
import type { Artist } from '@/services/api/catalog/artists'

defineProps<{
  artist: Artist | null
  monthlyListeners: number
  isFollowing: boolean
}>()

defineEmits<{
  playAll: []
  shuffle: []
  toggleFollow: []
}>()

function formatNumber(n: number) {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`
  return String(n)
}

function truncateBio(bio: string) {
  return bio.length > 200 ? bio.slice(0, 200) + '...' : bio
}
</script>
