<template>
  <div class="relative overflow-hidden rounded-[2rem]">
    <div
      class="absolute inset-0 bg-gradient-to-b from-[#1db954]/20 via-[#121212]/60 to-[#121212]"
    />
    <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />

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
          class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212] text-5xl text-white/40"
        >
          <i class="pi pi-user" />
        </div>
      </div>

      <div class="flex-1">
        <div
          v-if="artist?.is_verified"
          class="mb-2 flex items-center justify-center gap-1.5 md:justify-start"
        >
          <i class="pi pi-verified text-sm text-[#1db954]" />
          <span class="text-xs font-medium text-[#1db954]">Verified Artist</span>
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
            class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
            @click="$emit('playAll')"
          >
            <i class="pi pi-play-fill" />
            Play
          </button>

          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-8 py-3 text-sm font-bold text-white backdrop-blur transition hover:bg-white/15"
            @click="$emit('shuffle')"
          >
            <i class="pi pi-shuffle" />
            Shuffle
          </button>

          <button
            type="button"
            :class="[
              'inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-bold transition',
              isFollowing
                ? 'border border-[#1db954]/50 bg-[#1db954]/10 text-[#1db954] hover:bg-[#1db954]/20'
                : 'border border-white/15 bg-white/10 text-white hover:bg-white/15',
            ]"
            @click="$emit('toggleFollow')"
          >
            <i :class="isFollowing ? 'pi pi-check' : 'pi pi-plus'" class="text-xs" />
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
