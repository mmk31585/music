<template>
  <div :key="route.params.id" class="mx-auto w-full max-w-7xl px-4 pt-4 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="hero" />
      <div class="space-y-3">
        <SkeletonLoader v-for="i in 5" :key="i" variant="track" />
      </div>
    </div>

    <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
      <i class="pi pi-exclamation-circle text-4xl text-slate-500" />
      <h2 class="text-xl font-bold text-white">Artist not found</h2>
      <RouterLink to="/" class="text-sm font-medium text-[#1db954] underline underline-offset-2">
        Go home
      </RouterLink>
    </div>

    <template v-else-if="artist">
      <div
        class="relative overflow-hidden rounded-[2rem] border border-white/[0.06]"
        :style="headerBg"
      >
        <div class="aurora-spot-1 -top-60 -left-40 bg-[#1db954]/10" />
        <div class="aurora-spot-2 -right-40 -bottom-40 bg-[#60a5fa]/8" />
        <div class="relative z-10">
          <ArtistHero
            :artist="artist"
            :monthly-listeners="monthlyListeners"
            :is-following="isFollowing"
            @play-all="playAll"
            @shuffle="shuffleAll"
            @toggle-follow="toggleFollow"
          />
        </div>
      </div>

      <!-- Top Tracks -->
      <section class="mt-10">
        <template v-if="displayedTracks.length">
          <HomeSection title="Popular" eyebrow="Top tracks" />
          <div class="mt-4">
            <TrackList :tracks="displayedTracks" />
            <button
              v-if="tracks.length > 5"
              type="button"
              class="mt-3 text-sm font-medium text-slate-400 transition hover:text-white"
              @click="showAllTracks = !showAllTracks"
            >
              {{ showAllTracks ? 'Show less' : `Show all (${tracks.length} tracks)` }}
            </button>
          </div>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-music text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No tracks found for this artist</p>
        </div>
      </section>

      <!-- Albums -->
      <section class="mt-10">
        <template v-if="albums.length">
          <HomeSection title="Albums" eyebrow="Discography">
            <template #action>
              <span class="text-xs text-slate-500">{{ albums.length }} albums</span>
            </template>
          </HomeSection>
          <div class="mt-4">
            <HomeCarousel>
              <AlbumCard v-for="album in albums" :key="String(album.id)" :album="album" />
            </HomeCarousel>
          </div>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-inbox text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No albums yet</p>
        </div>
      </section>

      <!-- Related -->
      <section class="mt-10">
        <template v-if="related.length">
          <HomeSection title="Related Artists" eyebrow="You might also like" />
          <div class="mt-4">
            <HomeCarousel>
              <ArtistCard v-for="a in related" :key="String(a.id)" :artist="a" />
            </HomeCarousel>
          </div>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-inbox text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No related artists</p>
        </div>
      </section>

      <!-- Bio -->
      <section v-if="artist.bio" class="mt-10">
        <HomeSection title="About" eyebrow="Biography" />
        <div class="mt-4 max-w-3xl rounded-2xl bg-white/[0.03] p-6">
          <p class="text-sm leading-relaxed whitespace-pre-line text-slate-400">
            {{ bioExpanded ? artist.bio : truncateBio(artist.bio) }}
          </p>
          <button
            v-if="artist.bio.length > 300"
            type="button"
            class="mt-2 text-sm font-medium text-[#1db954] hover:underline"
            @click="bioExpanded = !bioExpanded"
          >
            {{ bioExpanded ? 'Show less' : 'Show more' }}
          </button>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useArtist } from '@/composables/catalog/useArtist'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import {
  HomeSection,
  HomeCarousel,
  TrackList,
  ArtistHero,
  ArtistCard,
  AlbumCard,
} from '@/components/music'

const route = useRoute()
const artistId = String(route.params.id)
const {
  artist,
  tracks,
  albums,
  related,
  isFollowing,
  monthlyListeners,
  loading,
  error,
  fetchArtist,
  toggleFollow,
} = useArtist(artistId)

const artistImageUrl = computed(() => artist.value?.image_url || null)
const { palette } = useAlbumColors(artistImageUrl)

const headerBg = computed(() => {
  const p = palette.value
  if (!artistImageUrl.value)
    return { background: 'linear-gradient(135deg, #0a0a0a 0%, #121212 100%)' }
  return {
    background: `linear-gradient(180deg, ${p.dark} 0%, ${p.dominant}99 40%, ${p.muted} 100%)`,
  }
})

const player = usePlayer()
const playerApi = usePlayerApi()
const showAllTracks = ref(false)
const bioExpanded = ref(false)

const displayedTracks = computed(() =>
  showAllTracks.value ? tracks.value : tracks.value.slice(0, 5),
)

function truncateBio(bio: string) {
  return bio.length > 300 ? bio.slice(0, 300) + '...' : bio
}

onMounted(() => {
  fetchArtist()
})

function playAll() {
  if (!tracks.value.length) return
  const queue = tracks.value.map((t) => ({
    id: String(t.id),
    title: t.title,
    artistName: t.artist_name || 'Unknown',
    albumTitle: t.album_title || null,
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
  }))
  player.setQueueAndPlay(queue, 0)
}

function shuffleAll() {
  const shuffled = [...tracks.value].sort(() => Math.random() - 0.5)
  if (!shuffled.length) return
  const queue = shuffled.map((t) => ({
    id: String(t.id),
    title: t.title,
    artistName: t.artist_name || 'Unknown',
    albumTitle: t.album_title || null,
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
  }))
  player.setQueueAndPlay(queue, 0)
}
</script>
