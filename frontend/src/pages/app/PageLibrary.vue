<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section class="relative overflow-hidden rounded-[2rem] border border-white/[0.06] bg-[#0C0C14] p-10 text-white">
      <div class="absolute -top-20 -right-20 h-60 w-60 rounded-full bg-[#1db954]/10 blur-3xl" />
      <div class="relative">
        <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Your music</p>
        <h1 class="mt-2 text-4xl font-black md:text-6xl">Library</h1>
        <p class="mt-3 max-w-2xl text-sm text-white/50">
          Saved songs, favorite artists, albums, and playlists will appear here.
        </p>
      </div>
    </section>

    <div v-if="loading" class="mt-10 space-y-6">
      <div v-for="i in 3" :key="i">
        <div class="mb-4 h-6 w-32 animate-pulse rounded bg-white/[0.06]" />
        <div class="flex gap-4">
          <div v-for="j in 4" :key="j" class="h-44 w-40 shrink-0 animate-pulse rounded-2xl bg-white/[0.06]" />
        </div>
      </div>
    </div>

    <template v-else>
      <section v-if="likedTracks.length" class="mt-10" aria-live="polite">
        <div class="mb-5 flex items-end justify-between gap-4">
          <div>
            <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Songs</p>
            <h2 class="mt-1 text-xl font-bold text-white md:text-2xl">Liked Tracks</h2>
          </div>
          <span class="text-xs tabular-nums text-white/30">{{ likedTracks.length }} tracks</span>
        </div>
        <div class="overflow-hidden rounded-2xl border border-white/[0.06] bg-white/[0.02]">
          <TrackRow
            v-for="(item, index) in likedTrackRows"
            :key="item.id"
            :track="item"
            :index="index"
            :queue="likedTrackRows"
          />
        </div>
      </section>

      <section v-if="likedAlbums.length" class="mt-12">
        <div class="mb-5 flex items-end justify-between gap-4">
          <div>
            <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Albums</p>
            <h2 class="mt-1 text-xl font-bold text-white md:text-2xl">Saved Albums</h2>
          </div>
          <span class="text-xs tabular-nums text-white/30">{{ likedAlbums.length }} albums</span>
        </div>
        <HomeCarousel>
          <AlbumCard
            v-for="album in likedAlbums"
            :key="album.album_id"
            :album="libraryAlbumToAlbum(album)"
          />
        </HomeCarousel>
      </section>

      <section v-if="followedArtists.length" class="mt-12">
        <div class="mb-5 flex items-end justify-between gap-4">
          <div>
            <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Artists</p>
            <h2 class="mt-1 text-xl font-bold text-white md:text-2xl">Followed Artists</h2>
          </div>
          <span class="text-xs tabular-nums text-white/30">{{ followedArtists.length }} artists</span>
        </div>
        <HomeCarousel>
          <RouterLink
            v-for="artist in followedArtists"
            :key="artist.artist_id"
            :to="`/artist/${artist.artist_id}`"
            class="group block w-40 shrink-0 space-y-3"
          >
            <div class="mx-auto h-36 w-36 overflow-hidden rounded-full bg-white/[0.06] ring-1 ring-white/10 transition group-hover:ring-[#1db954]/30">
              <img
                v-if="artist.cover_url"
                :src="artist.cover_url"
                :alt="artist.name"
                loading="lazy"
                class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-user text-3xl text-slate-500" />
              </div>
            </div>
            <div class="space-y-0.5 text-center">
              <p class="truncate text-sm font-bold text-white">{{ artist.name }}</p>
              <p class="text-xs text-slate-500">Artist</p>
            </div>
          </RouterLink>
        </HomeCarousel>
      </section>

      <section v-if="!hasAnyData" class="mt-12 grid gap-5 md:grid-cols-3">
        <div
          v-for="item in placeholderItems"
          :key="item.title"
          class="rounded-2xl border border-white/[0.06] bg-white/[0.02] p-6"
        >
          <div
            class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.06] text-lg"
          >
            <i aria-hidden="true" :class="[item.icon, 'text-white/50']" />
          </div>
          <h2 class="mt-4 text-base font-bold text-white">{{ item.title }}</h2>
          <p class="mt-1 text-sm text-white/40">{{ item.description }}</p>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { TrackRow, HomeCarousel, AlbumCard } from '@/components/music'
import { useLibraryApi } from '@/services/api/library'
import type { LibraryAlbum, LibraryArtist, LibraryTrack } from '@/services/api/library'
import type { Album } from '@/services/api/catalog/albums'

const libraryApi = useLibraryApi()

const loading = ref(false)
const likedTracks = ref<LibraryTrack[]>([])
const likedAlbums = ref<LibraryAlbum[]>([])
const followedArtists = ref<LibraryArtist[]>([])

const hasAnyData = computed(() =>
  likedTracks.value.length > 0 || likedAlbums.value.length > 0 || followedArtists.value.length > 0
)

const likedTrackRows = computed(() =>
  likedTracks.value.map((item) => ({
    id: item.track_id,
    title: item.title,
    artist_name: item.artist_name,
    album_title: item.album_title,
    cover_url: item.cover_url ?? null,
    duration_seconds: item.duration_seconds,
  })),
)

function libraryAlbumToAlbum(item: LibraryAlbum): Album {
  return {
    id: item.album_id,
    title: item.title,
    cover_url: item.cover_url ?? null,
    artist_id: item.artist_id ?? null,
    artist_name: item.artist_name ?? null,
    release_date: item.release_date ?? null,
    track_count: 0,
  }
}

const placeholderItems = [
  {
    title: 'Liked Songs',
    description: 'Tracks you liked will be collected here.',
    icon: 'pi pi-heart',
  },
  {
    title: 'Saved Albums',
    description: 'Albums saved to your library.',
    icon: 'pi pi-images',
  },
  {
    title: 'Followed Artists',
    description: 'Artists you follow and listen to often.',
    icon: 'pi pi-users',
  },
]

onMounted(async () => {
  loading.value = true
  try {
    const [tracks, albums, artists] = await Promise.all([
      libraryApi.getLikedTracks().catch(() => []),
      libraryApi.getLikedAlbums().catch(() => []),
      libraryApi.getFollowedArtists().catch(() => []),
    ])
    likedTracks.value = Array.isArray(tracks) ? tracks : []
    likedAlbums.value = Array.isArray(albums) ? albums : []
    followedArtists.value = Array.isArray(artists) ? artists : []
  } catch (err) {
    console.error('Failed to fetch library:', err)
  } finally {
    loading.value = false
  }
})
</script>
