<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <section
      class="rounded-3xl bg-gradient-to-r from-[#1db954] via-[#169c46] to-[#0f0f0f] p-8 text-white"
    >
      <p class="text-sm tracking-[0.3em] text-white/70 uppercase">Welcome</p>
      <h1 class="mt-3 text-4xl font-extrabold md:text-5xl">Discover your next favorite track</h1>
      <p class="mt-4 max-w-2xl text-white/85">
        Explore public catalog tracks, artists, albums, and genres.
      </p>
    </section>

    <section class="mt-8">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-2xl font-bold">Tracks</h2>
        <span class="text-sm text-slate-400">{{ tracks.length }} results</span>
      </div>

      <div v-if="loading" class="space-y-3">
        <div v-for="i in 8" :key="i" class="h-16 animate-pulse rounded-2xl bg-white/5" />
      </div>

      <div
        v-else-if="tracks.length === 0"
        class="rounded-2xl border border-white/10 bg-white/5 px-5 py-8 text-center text-slate-300"
      >
        No public tracks yet.
      </div>

      <div v-else class="space-y-2">
        <TrackRow
          v-for="(track, index) in tracks"
          :key="track.id"
          :track="track"
          :index="index"
          :queue="tracks"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { TrackRow } from '@/components/music'
import { useCatalogTracks } from '@/composables/catalog/useCatalogTracks'

const { tracks, loading } = useCatalogTracks()
</script>
