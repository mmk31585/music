<template>
  <div :key="String(route.params.id)" class="mx-auto w-full max-w-7xl px-4 pt-4 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="hero" />
      <div class="space-y-3">
        <SkeletonLoader v-for="i in 5" :key="i" variant="track" />
      </div>
    </div>

    <AppEmptyState
      v-else-if="error"
      variant="error"
      icon="pi pi-exclamation-circle"
      title="Artist not found"
      description="This artist may have been removed or the link is invalid."
    >
      <template #action>
        <RouterLink
          to="/"
          class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-5 py-2.5 text-sm font-bold text-black transition hover:bg-[#1ed760]"
        >
          <i aria-hidden="true" class="pi pi-home" />
          Go home
        </RouterLink>
      </template>
    </AppEmptyState>

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
      <section class="mt-14" aria-live="polite">
        <template v-if="displayedTracks.length">
          <div class="flex items-baseline justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Top tracks</p>
              <h2 class="mt-1 text-xl font-bold text-white md:text-2xl">Popular</h2>
            </div>
            <span class="text-xs tabular-nums text-white/25">{{ tracks.length }} tracks</span>
          </div>
          <div class="mt-5">
            <TrackList :tracks="displayedTracks" />
            <button
              v-if="tracks.length > 5"
              type="button"
              class="mt-3 text-sm font-medium text-white/40 transition hover:text-white"
              @click="showAllTracks = !showAllTracks"
            >
              {{ showAllTracks ? 'Show less' : `Show all (${tracks.length} tracks)` }}
            </button>
          </div>
        </template>
        <AppEmptyState
          v-else
          icon="pi pi-music"
          title="No tracks found"
          description="No tracks found for this artist"
        />
      </section>

      <!-- Section divider -->
      <div v-if="albums.length || related.length || artist.bio" class="relative mt-14">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-white/[0.06]" />
        </div>
      </div>

      <!-- Albums -->
      <section class="mt-14">
        <template v-if="albums.length">
          <div class="flex items-baseline justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Discography</p>
              <h2 class="mt-1 text-xl font-bold text-white md:text-2xl">Albums</h2>
            </div>
            <span class="text-xs tabular-nums text-white/25">{{ albums.length }} album{{ albums.length === 1 ? '' : 's' }}</span>
          </div>
          <div class="mt-5">
            <HomeCarousel>
              <AlbumCard v-for="album in albums" :key="String(album.id)" :album="album" />
            </HomeCarousel>
          </div>
        </template>
        <AppEmptyState
          v-else
          icon="pi pi-inbox"
          title="No albums yet"
          description="Albums will appear here once they're released."
        />
      </section>

      <!-- Related -->
      <section class="mt-14">
        <template v-if="related.length">
          <div class="flex items-baseline justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">You might also like</p>
              <h2 class="mt-1 text-xl font-bold text-white md:text-2xl">Related Artists</h2>
            </div>
          </div>
          <div class="mt-5">
            <HomeCarousel>
              <ArtistCard v-for="a in related" :key="String(a.id)" :artist="a" />
            </HomeCarousel>
          </div>
        </template>
        <AppEmptyState
          v-else
          icon="pi pi-inbox"
          title="No related artists"
          description="Check back later for similar artist suggestions."
        />
      </section>

      <!-- Bio -->
      <section v-if="artist.bio" class="mt-14">
        <div class="relative mb-8">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-white/[0.06]" />
          </div>
          <div class="relative flex justify-center">
            <span class="bg-[#0a0a0a] px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
              Biography
            </span>
          </div>
        </div>
        <div class="mx-auto max-w-3xl rounded-2xl border border-white/[0.04] bg-white/[0.02] p-8">
          <p class="text-sm leading-relaxed whitespace-pre-line text-white/50">
            {{ bioExpanded ? artist.bio : truncateBio(artist.bio) }}
          </p>
          <button
            v-if="artist.bio.length > 300"
            type="button"
            class="mt-3 text-sm font-medium text-[#1db954] transition hover:underline"
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
import { SkeletonLoader, AppEmptyState } from '@/components/common'
import { useArtist } from '@/composables/catalog/useArtist'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import {
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
