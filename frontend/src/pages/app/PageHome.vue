<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <HomeHero
      v-if="heroItems.length"
      :items="heroItems"
      @play="handleHeroPlay"
      @add-to-library="handleHeroAddToLibrary"
    />

    <div v-if="loading" class="mt-6 space-y-8">
      <div v-for="s in 4" :key="s">
        <div class="mb-4 h-6 w-40 animate-pulse rounded bg-white/[0.06]" />
        <div class="flex gap-4">
          <div v-for="i in 5" :key="i" class="h-36 w-28 shrink-0 animate-pulse rounded-2xl bg-white/[0.06] sm:w-32 md:w-36" />
        </div>
      </div>
    </div>

    <template v-else-if="!hasData">
      <section class="relative mt-6 overflow-hidden rounded-[2rem] border border-white/10 bg-[#121212] p-8 text-white shadow-2xl">
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
            <RouterLink
              to="/discover"
              class="rounded-full border border-white/15 bg-white/10 px-6 py-3 text-sm font-bold text-white backdrop-blur transition hover:bg-white/15"
            >
              Discover
            </RouterLink>
          </div>
        </div>
      </section>

      <section class="mt-10">
        <HomeSectionHeader title="Start exploring" subtitle="Upload tracks from the admin panel to see them here." />
        <div class="rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center">
          <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white">
            <i class="pi pi-music" />
          </div>
          <h3 class="mt-5 text-xl font-bold text-white">No tracks yet</h3>
          <p class="mt-2 text-sm text-slate-400">
            Upload tracks from the admin panel to see them here.
          </p>
        </div>
      </section>
    </template>

    <template v-else>
      <section v-if="recentPlays.length" class="mt-6">
        <HomeSectionHeader title="Recently played" subtitle="Jump back in" see-all-route="/recently-played" />
        <HomeCarousel>
          <HomeTrackCard
            v-for="(item, i) in recentPlays"
            :key="item.track_id"
            :item="item"
            :delay="i * 30"
            :is-playing="isCurrentlyPlaying(item)"
            @play="handlePlay"
          />
        </HomeCarousel>
      </section>

      <section v-if="popular.length" class="mt-10">
        <HomeSectionHeader title="Popular right now" subtitle="Trending across the catalog" see-all-route="/recommendations/popular" />
        <HomeCarousel>
          <HomeTrackCard
            v-for="(item, i) in popular"
            :key="item.id"
            :item="item"
            :delay="i * 30"
            :badge="getScoreBadge(item)"
            :is-playing="isCurrentlyPlaying(item)"
            @play="handlePlay"
          />
        </HomeCarousel>
      </section>

      <section v-if="forYou.length" class="mt-10">
        <HomeSectionHeader title="Made for you" subtitle="Personalized picks" see-all-route="/recommendations/for-you" />
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          <HomeQuickPlayCard
            v-for="(item, i) in forYou.slice(0, 10)"
            :key="item.id"
            :item="item"
            :delay="i * 20"
            :is-playing="isCurrentlyPlaying(item)"
            @play="handlePlay"
          />
        </div>
      </section>

      <section v-if="recent.length" class="mt-10">
        <HomeSectionHeader title="New arrivals" subtitle="Latest additions" see-all-route="/recommendations/recent" />
        <HomeCarousel>
          <HomeTrackCard
            v-for="(item, i) in recent"
            :key="item.id"
            :item="item"
            :delay="i * 30"
            :is-playing="isCurrentlyPlaying(item)"
            @play="handlePlay"
          />
        </HomeCarousel>
      </section>

      <section v-if="albums.length" class="mt-10">
        <HomeSectionHeader title="Albums" subtitle="Explore full collections" see-all-route="/search" />
        <HomeCarousel>
          <RouterLink
            v-for="album in albums"
            :key="album.id"
            :to="`/album/${album.id}`"
            class="group block w-44 shrink-0 space-y-3"
          >
            <div class="relative aspect-square overflow-hidden rounded-2xl bg-white/[0.06] shadow-lg ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50">
              <img
                v-if="album.cover_url"
                :src="album.cover_url"
                :alt="album.title"
                loading="lazy"
                class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i class="pi pi-compact-disc text-3xl text-slate-500" />
              </div>
              <div class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100">
                <div class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl">
                  <i class="pi pi-play-fill text-lg" />
                </div>
              </div>
            </div>
            <div class="space-y-0.5 px-1">
              <p class="truncate text-sm font-bold text-white">{{ album.title }}</p>
              <p v-if="album.artist_name" class="truncate text-xs text-slate-400">{{ album.artist_name }}</p>
              <p class="text-xs text-slate-500">{{ album.release_date?.slice(0, 4) || '' }}{{ album.track_count ? ` • ${album.track_count} tracks` : '' }}</p>
            </div>
          </RouterLink>
        </HomeCarousel>
      </section>

      <section v-if="artists.length" class="mt-10">
        <HomeSectionHeader title="Featured artists" subtitle="Meet the creators" />
        <HomeCarousel>
          <RouterLink
            v-for="artist in artists"
            :key="artist.id"
            :to="`/artist/${artist.id}`"
            class="group block w-40 shrink-0 space-y-3"
          >
            <div class="mx-auto h-36 w-36 overflow-hidden rounded-full bg-white/[0.06] ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50">
              <img
                v-if="artist.image_url"
                :src="artist.image_url"
                :alt="artist.name"
                loading="lazy"
                class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i class="pi pi-user text-3xl text-slate-500" />
              </div>
            </div>
            <div class="space-y-0.5 text-center">
              <p class="truncate text-sm font-bold text-white">{{ artist.name }}</p>
              <p class="text-xs text-slate-500">Artist</p>
            </div>
          </RouterLink>
        </HomeCarousel>
      </section>

      <section class="mt-10">
        <div class="grid gap-4 md:grid-cols-3">
          <RouterLink
            to="/recommendations"
            class="group rounded-2xl border border-white/10 bg-gradient-to-br from-[#1db954]/20 to-transparent p-5 transition hover:-translate-y-0.5 hover:border-[#1db954]/30"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-[#1db954]/30 text-[#1db954]">
              <i class="pi pi-star text-xl" />
            </div>
            <h3 class="mt-4 text-base font-bold text-white">Recommendations</h3>
            <p class="mt-1 text-sm text-slate-400">Discover popular, recent, and personalized tracks.</p>
          </RouterLink>
          <RouterLink
            to="/discover"
            class="group rounded-2xl border border-white/10 bg-gradient-to-br from-blue-500/20 to-transparent p-5 transition hover:-translate-y-0.5 hover:border-blue-500/30"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-500/30 text-blue-400">
              <i class="pi pi-compass text-xl" />
            </div>
            <h3 class="mt-4 text-base font-bold text-white">Discover</h3>
            <p class="mt-1 text-sm text-slate-400">Explore trending tracks and new releases.</p>
          </RouterLink>
          <RouterLink
            to="/ai/mood-explorer"
            class="group rounded-2xl border border-white/10 bg-gradient-to-br from-purple-500/20 to-transparent p-5 transition hover:-translate-y-0.5 hover:border-purple-500/30"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-purple-500/30 text-purple-400">
              <i class="pi pi-magic text-xl" />
            </div>
            <h3 class="mt-4 text-base font-bold text-white">Mood Explorer</h3>
            <p class="mt-1 text-sm text-slate-400">Find music that matches your vibe.</p>
          </RouterLink>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useHomeFeed } from '@/composables/useHomeFeed'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'
import {
  HomeHero,
  HomeCarousel,
  HomeTrackCard,
  HomeQuickPlayCard,
} from '@/components/music'
import HomeSectionHeader from '@/components/music/home/HomeSectionHeader.vue'
import type { HeroItem } from '@/components/music/home/HomeHero.vue'

const {
  popular, forYou, recent, albums, artists, recentPlays,
  loading, hasData, fetchHomeFeed,
} = useHomeFeed()

const player = usePlayer()
const playerApi = usePlayerApi()

const heroItems = computed<HeroItem[]>(() => {
  const tracks = popular.value.slice(0, 5)
  return tracks.map((t, i) => ({
    id: String(t.id),
    title: t.title || 'Untitled',
    subtitle: t.artist_name || 'Unknown artist',
    image: t.cover_url || '',
    badge: i === 0 ? 'Trending' : i === 1 ? 'Popular' : 'Featured',
    badgeVariant: i === 0 ? 'green' : 'purple',
    type: 'album' as const,
  }))
})

function getScoreBadge(item: any): string | undefined {
  if (item.score >= 90) return '🔥 Hot'
  if (item.score >= 75) return 'Trending'
  return undefined
}

function isCurrentlyPlaying(item: any): boolean {
  const currentId = player.currentTrack.value?.id
  if (!currentId) return false
  return currentId === String(item.id) || currentId === String(item.track_id)
}

function buildPlaybackTrack(item: any): PlaybackTrack {
  const id = String(item.id || item.track_id)
  return {
    id,
    title: item.title || 'Untitled',
    artistName: item.artist_name || item.artistName || 'Unknown artist',
    albumTitle: item.album_title || item.albumTitle || null,
    coverUrl: item.cover_url || item.coverUrl || null,
    durationSeconds: item.duration_seconds ?? item.durationSeconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(id),
  }
}

async function handlePlay(item: any) {
  await player.toggleTrack(buildPlaybackTrack(item))
}

function handleHeroPlay(item: HeroItem) {
  const track = popular.value.find((t) => String(t.id) === item.id)
  if (track) {
    void handlePlay(track)
  }
}

function handleHeroAddToLibrary(_item: HeroItem) {
  // Placeholder - library add would go here
}

onMounted(() => {
  void fetchHomeFeed()
})
</script>
