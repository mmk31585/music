<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section class="rounded-[2rem] border border-white/10 bg-white/[0.05] p-8 text-white">
      <p class="text-sm font-bold tracking-[0.35em] text-[#1db954] uppercase">Your music</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">Library</h1>
      <p class="mt-4 max-w-2xl text-slate-300">
        Saved songs, favorite artists, albums, and playlists will appear here.
      </p>
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
      <section v-if="likedTracks.length" class="mt-10">
        <div class="mb-4 flex items-end justify-between">
          <div>
            <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Songs</p>
            <h2 class="mt-1 text-2xl font-black text-white">Liked Tracks</h2>
          </div>
          <span class="text-xs text-slate-400">{{ likedTracks.length }} tracks</span>
        </div>
        <div class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur">
          <TrackRow
            v-for="(item, index) in likedTrackRows"
            :key="item.id"
            :track="item"
            :index="index"
            :queue="likedTrackRows"
          />
        </div>
      </section>

      <section v-if="likedAlbums.length" class="mt-10">
        <div class="mb-4 flex items-end justify-between">
          <div>
            <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Albums</p>
            <h2 class="mt-1 text-2xl font-black text-white">Saved Albums</h2>
          </div>
          <span class="text-xs text-slate-400">{{ likedAlbums.length }} albums</span>
        </div>
        <HomeCarousel>
          <AlbumCard
            v-for="album in likedAlbums"
            :key="album.album_id"
            :album="libraryAlbumToAlbum(album)"
          />
        </HomeCarousel>
      </section>

      <section v-if="followedArtists.length" class="mt-10">
        <div class="mb-4 flex items-end justify-between">
          <div>
            <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Artists</p>
            <h2 class="mt-1 text-2xl font-black text-white">Followed Artists</h2>
          </div>
          <span class="text-xs text-slate-400">{{ followedArtists.length }} artists</span>
        </div>
        <HomeCarousel>
          <RouterLink
            v-for="artist in followedArtists"
            :key="artist.artist_id"
            :to="`/artist/${artist.artist_id}`"
            class="group block w-40 shrink-0 space-y-3"
          >
            <div class="mx-auto h-36 w-36 overflow-hidden rounded-full bg-white/[0.06] ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50">
              <img
                v-if="artist.cover_url"
                :src="artist.cover_url"
                :alt="artist.name"
                loading="lazy"
                class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
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

      <section v-if="!hasAnyData" class="mt-10 grid gap-5 md:grid-cols-3">
        <div
          v-for="item in placeholderItems"
          :key="item.title"
          class="rounded-3xl border border-white/10 bg-black/20 p-6"
        >
          <div
            class="flex h-14 w-14 items-center justify-center rounded-2xl text-2xl"
            :class="item.iconClass"
          >
            <i :class="item.icon" />
          </div>
          <h2 class="mt-5 text-xl font-black text-white">{{ item.title }}</h2>
          <p class="mt-2 text-sm text-slate-400">{{ item.description }}</p>
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
    iconClass: 'bg-pink-500/20 text-pink-300',
  },
  {
    title: 'Saved Albums',
    description: 'Albums saved to your library.',
    icon: 'pi pi-images',
    iconClass: 'bg-purple-500/20 text-purple-300',
  },
  {
    title: 'Followed Artists',
    description: 'Artists you follow and listen to often.',
    icon: 'pi pi-users',
    iconClass: 'bg-[#1db954]/20 text-[#1db954]',
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
  } catch {
    // silent
  } finally {
    loading.value = false
  }
})
</script>
