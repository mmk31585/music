<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="relative overflow-hidden rounded-[2rem] border border-white/10 bg-[#121212] p-8 text-white shadow-2xl"
    >
      <div class="absolute -top-20 -right-20 h-72 w-72 rounded-full bg-[#1db954]/30 blur-3xl" />
      <div class="absolute right-1/3 bottom-0 h-44 w-44 rounded-full bg-emerald-400/20 blur-3xl" />

      <div class="relative z-10 max-w-3xl">
        <p class="text-sm font-bold tracking-[0.35em] text-[#1db954] uppercase">Welcome back</p>

        <h1 class="mt-4 text-4xl leading-tight font-black md:text-6xl">
          Discover your next favorite track.
        </h1>

        <p class="mt-5 max-w-2xl text-base leading-7 text-slate-300">
          Stream tracks, explore artists, search albums, and build the perfect vibe.
        </p>

        <div class="mt-8 flex flex-wrap gap-3">
          <RouterLink
            to="/search"
            class="rounded-full bg-[#1db954] px-6 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
          >
            Explore Music
          </RouterLink>

          <button
            type="button"
            class="rounded-full border border-white/15 bg-white/10 px-6 py-3 text-sm font-bold text-white backdrop-blur transition hover:bg-white/15"
            @click="playFirstTrack"
          >
            Play First Track
          </button>
        </div>
      </div>
    </section>

    <section class="mt-10">
      <div class="mb-5 flex items-end justify-between gap-4">
        <div>
          <p class="text-sm font-semibold text-[#1db954]">Public catalog</p>
          <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Tracks</h2>
        </div>

        <span class="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-slate-300">
          {{ tracks.length }} results
        </span>
      </div>

      <div v-if="loading" class="space-y-3">
        <div v-for="i in 8" :key="i" class="h-[68px] animate-pulse rounded-2xl bg-white/[0.06]" />
      </div>

      <div
        v-else-if="tracks.length === 0"
        class="rounded-3xl border border-white/10 bg-white/[0.04] px-6 py-16 text-center"
      >
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <i class="pi pi-music" />
        </div>

        <h3 class="mt-5 text-xl font-bold text-white">No tracks yet</h3>

        <p class="mt-2 text-sm text-slate-400">
          Upload tracks from the admin panel to see them here.
        </p>
      </div>

      <div
        v-else
        class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur"
      >
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
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'

const { tracks, loading } = useCatalogTracks()
const player = usePlayer()
const playerApi = usePlayerApi()

async function playFirstTrack() {
  const firstTrack = tracks.value[0] as any
  if (!firstTrack) return

  const id = String(firstTrack.id)

  const playbackTrack: PlaybackTrack = {
    id,
    title: firstTrack.title || 'Untitled',
    artistName:
      firstTrack.artist_name ||
      firstTrack.artistName ||
      firstTrack.artist?.name ||
      firstTrack.artist ||
      'Unknown artist',
    albumTitle: firstTrack.album_title || firstTrack.albumTitle || firstTrack.album?.title || null,
    coverUrl: firstTrack.cover_url || firstTrack.coverUrl || firstTrack.cover || null,
    durationSeconds:
      firstTrack.duration_seconds ?? firstTrack.durationSeconds ?? firstTrack.duration ?? null,
    streamUrl: playerApi.getTrackStreamUrl(id),
  }

  await player.toggleTrack(playbackTrack)
}
</script>
