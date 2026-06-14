<template>
  <div :key="String(route.params.id)" class="mx-auto w-full max-w-5xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <div class="flex gap-6">
        <SkeletonLoader variant="card" class="w-64 shrink-0" />
        <div class="flex-1 space-y-3">
          <SkeletonLoader variant="lines" :lines="3" />
        </div>
      </div>
      <div class="space-y-2">
        <SkeletonLoader v-for="i in 6" :key="i" variant="track" />
      </div>
    </div>

    <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
      <i class="pi pi-exclamation-circle text-4xl text-slate-500" />
      <h2 class="text-xl font-bold text-white">Album not found</h2>
      <RouterLink to="/" class="text-sm font-medium text-[#1db954] underline underline-offset-2">
        Go home
      </RouterLink>
    </div>

    <template v-else-if="album">
      <!-- Header with dynamic background -->
      <div
        class="relative overflow-hidden rounded-[2rem] border border-white/[0.06] p-8 md:p-10"
        :style="headerBg"
      >
        <div class="aurora-spot-1 -top-40 -left-40 bg-[#1db954]/10" />
        <div class="aurora-spot-2 -right-40 -bottom-40 bg-[#60a5fa]/8" />

        <div class="relative z-10 flex flex-col gap-6 md:flex-row md:items-end">
          <div
            class="vinyl-glow h-64 w-64 shrink-0 overflow-hidden rounded-2xl bg-white/[0.06] shadow-2xl ring-1 ring-white/10"
          >
            <img
              v-if="album.cover_url"
              :src="album.cover_url"
              :alt="album.title"
              loading="lazy"
              class="h-full w-full object-cover"
              @error="onImgError"
            />
            <div v-else class="flex h-full items-center justify-center">
              <i class="pi pi-compact-disc text-5xl text-slate-500" />
            </div>
          </div>

          <div class="flex-1">
            <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">Album</p>
            <h1 class="mt-2 text-3xl font-black text-white md:text-5xl">
              {{ album.title }}
            </h1>

            <div class="mt-3 flex flex-wrap items-center gap-2 text-sm text-slate-400">
              <RouterLink
                v-if="artist"
                :to="`/artist/${artist.id}`"
                class="font-bold text-white underline underline-offset-2 transition hover:text-[#1db954]"
              >
                {{ artist.name }}
              </RouterLink>
              <template v-else-if="album.artist_name">
                <span class="font-bold text-white">{{ album.artist_name }}</span>
              </template>
              <span v-if="album.release_date"> • {{ album.release_date.slice(0, 4) }}</span>
              <span> • {{ tracks.length }} {{ tracks.length === 1 ? 'track' : 'tracks' }}</span>
              <span v-if="totalDuration > 0"> • {{ formattedDuration }}</span>
            </div>

            <div class="mt-6 flex flex-wrap items-center gap-3">
              <button
                type="button"
                class="glow-green inline-flex items-center gap-2 rounded-full bg-[#1db954] px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
                @click="playAll"
              >
                <i class="pi pi-play-fill" />
                Play
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-6 py-3 text-sm font-bold text-white backdrop-blur transition hover:bg-white/15"
                :class="{ 'border-[#1db954]/50 text-[#1db954]': isLiked }"
                @click="toggleLike"
              >
                <i :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" />
                {{ isLiked ? 'Saved' : 'Save' }}
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-6 py-3 text-sm font-bold text-white backdrop-blur transition hover:bg-white/15"
                @click="shareAlbum"
              >
                <i class="pi pi-share-alt" />
                Share
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tracklist -->
      <section class="mt-10">
        <TrackList :tracks="tracks" empty-message="No tracks in this album" />
      </section>

      <!-- Credits / Artists -->
      <section v-if="albumArtists.length > 1" class="mt-10">
        <h2 class="mb-4 text-lg font-bold text-white">Credits</h2>
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="aa in albumArtists"
            :key="aa.id"
            class="flex items-center gap-3 rounded-xl border border-white/[0.06] bg-white/[0.03] px-4 py-3 transition hover:bg-white/[0.06]"
          >
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[#1db954]/20 to-purple-500/20 text-sm font-bold text-white"
            >
              {{ aa.name.charAt(0).toUpperCase() }}
            </div>
            <div class="min-w-0">
              <RouterLink
                :to="`/artist/${aa.id}`"
                class="block truncate text-sm font-semibold text-white transition hover:text-[#1db954] hover:underline"
              >
                {{ aa.name }}
              </RouterLink>
              <p v-if="aa.role" class="text-xs text-slate-500 capitalize">{{ aa.role }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Related albums -->
      <section v-if="relatedAlbums.length > 0" class="mt-12">
        <HomeSection title="More by" :eyebrow="artist?.name ?? mainArtist?.name ?? 'Artist'">
          <template #action>
            <RouterLink
              v-if="artist"
              :to="`/artist/${artist.id}`"
              class="text-xs font-bold text-slate-400 transition hover:text-white"
            >
              View artist
            </RouterLink>
          </template>
        </HomeSection>
        <div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          <RouterLink
            v-for="ra in relatedAlbums"
            :key="ra.id"
            :to="`/album/${ra.id}`"
            class="group block"
          >
            <div
              class="relative mb-2 aspect-square overflow-hidden rounded-xl bg-white/[0.06] shadow-lg ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50"
            >
              <img
                v-if="ra.cover_url"
                :src="ra.cover_url"
                :alt="ra.title"
                loading="lazy"
                class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i class="pi pi-compact-disc text-2xl text-slate-500" />
              </div>
              <div
                class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
              >
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954]/90 text-black shadow-xl"
                >
                  <i class="pi pi-play-fill text-lg" />
                </div>
              </div>
            </div>
            <p class="truncate text-sm font-semibold text-white">{{ ra.title }}</p>
            <p class="truncate text-xs text-slate-500">{{ ra.artist_name || 'Artist' }}</p>
          </RouterLink>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useAlbum } from '@/composables/catalog/useAlbum'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { HomeSection, TrackList } from '@/components/music'
import { onImgError } from '@/utils/helpers'
import { useSocialShare } from '@/composables/social'

const route = useRoute()
const albumId = String(route.params.id)
const {
  album,
  tracks,
  artist,
  albumArtists,
  relatedAlbums,
  mainArtist,
  featuredArtists,
  isLiked,
  totalDuration,
  formattedDuration,
  loading,
  error,
  fetchAlbum,
  toggleLike,
} = useAlbum(albumId) as any

const coverUrl = computed(() => album.value?.cover_url || null)
const { palette } = useAlbumColors(coverUrl)

const headerBg = computed(() => {
  const p = palette.value
  if (!coverUrl.value) return { background: 'linear-gradient(135deg, #0a0a0a 0%, #121212 100%)' }
  return {
    background: `linear-gradient(180deg, ${p.dark} 0%, ${p.dominant}99 40%, ${p.muted} 100%)`,
  }
})

const player = usePlayer()
const playerApi = usePlayerApi()

onMounted(() => {
  fetchAlbum()
})

function playAll() {
  if (!tracks.value.length) return
  const queue = tracks.value.map((t: any) => ({
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

const { copyLink } = useSocialShare()
function shareAlbum() {
  if (!album.value) return
  copyLink({
    id: album.value.id,
    title: album.value.title,
    type: 'album',
    artistName: artist.value?.name || album.value.artist_name,
    coverUrl: album.value.cover_url,
  })
}
</script>
