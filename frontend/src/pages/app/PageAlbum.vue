<template>
  <div :key="String(route.params.id)" class="mx-auto w-full max-w-6xl px-4 pt-8 pb-36 md:px-8 lg:px-10">
    <button
      type="button"
      class="mb-6 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-white/50 transition hover:bg-white/6 hover:text-white"
      @click="goBack"
    >
      <i aria-hidden="true" class="pi pi-arrow-left text-xs" />
      Back
    </button>

    <div v-if="loading" class="grid gap-8 md:grid-cols-[300px_1fr]">
      <SkeletonLoader variant="card" class="w-full" />
      <div class="space-y-4">
        <SkeletonLoader variant="lines" :lines="4" class="max-w-md" />
        <SkeletonLoader variant="lines" :lines="2" class="max-w-xs" />
      </div>
      <div class="col-span-full space-y-2">
        <SkeletonLoader v-for="i in 8" :key="i" variant="track" />
      </div>
    </div>

    <AppEmptyState
      v-else-if="error"
      variant="error"
      icon="pi pi-exclamation-circle"
      title="Album not found"
      description="This album may have been removed or the link is invalid."
    >
      <template #action>
        <RouterLink
          to="/"
          class="inline-flex items-center gap-2 rounded-full bg-spotify px-5 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover"
        >
          <i aria-hidden="true" class="pi pi-home" />
          Go home
        </RouterLink>
      </template>
    </AppEmptyState>

    <template v-else-if="album">
      <!-- Ambient gradient background from cover art -->
      <div
        class="pointer-events-none fixed inset-0 transition-all duration-1000"
        :style="ambientBg"
        aria-hidden="true"
      />

      <!-- Subtle groove texture overlay -->
      <div
        class="pointer-events-none fixed inset-0 opacity-[0.015]"
        style="background-image: repeating-radial-gradient(circle at 50% 50%, transparent 0, transparent 2px, rgba(255,255,255,0.04) 2px, rgba(255,255,255,0.04) 3px); background-size: 6px 6px;"
        aria-hidden="true"
      />

      <div class="relative z-10">
        <!-- Hero -->
        <div class="flex flex-col gap-10 md:flex-row md:items-end">
          <!-- Cover art with physical-object treatment -->
          <div class="group relative shrink-0">
            <div
              class="relative h-75 w-75 overflow-hidden rounded-2xl bg-white/6 shadow-2xl ring-1 ring-white/10 transition-all duration-500 group-hover:shadow-xl group-hover:-rotate-1 group-hover:scale-[1.02]"
              :style="coverGlowStyle"
            >
              <img
                v-if="album.cover_url"
                :src="album.cover_url"
                :alt="album.title"
                loading="eager"
                class="h-full w-full object-cover"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-compact-disc text-5xl text-slate-500" />
              </div>
            </div>
            <!-- Sleeve frame accent -->
            <div
              class="pointer-events-none absolute inset-0 rounded-2xl ring-1 ring-inset ring-white/4"
              aria-hidden="true"
            />
          </div>

          <!-- Album info -->
          <div class="flex-1">
            <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">
              Album
            </p>
            <h1 class="mt-3 text-4xl font-black leading-[1.1] text-white md:text-5xl lg:text-6xl">
              {{ album.title }}
            </h1>

            <div class="mt-4 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm">
              <template v-if="artist">
                <RouterLink
                  :to="`/artist/${artist.id}`"
                  class="font-bold text-white underline underline-offset-4 decoration-white/20 transition hover:text-spotify hover:decoration-[#1db954]"
                >
                  {{ artist.name }}
                </RouterLink>
              </template>
              <template v-else-if="album.artist_name">
                <span class="font-bold text-white">{{ album.artist_name }}</span>
              </template>
              <span v-if="album.release_date" class="text-white/30">·</span>
              <span v-if="album.release_date" class="text-white/50">{{ album.release_date.slice(0, 4) }}</span>
              <span v-if="album.genre" class="text-white/30">·</span>
              <span v-if="album.genre" class="text-white/50">{{ album.genre }}</span>
            </div>

            <!-- Stats -->
            <div class="mt-6 flex items-center gap-8">
              <div class="flex items-baseline gap-2">
                <span class="text-2xl font-black tabular-nums text-white">{{ tracks.length }}</span>
                <span class="text-xs text-white/40">{{ tracks.length === 1 ? 'track' : 'tracks' }}</span>
              </div>
              <div v-if="totalDuration" class="flex items-baseline gap-2.5">
                <span class="text-2xl font-black tabular-nums text-white">{{ formattedDuration }}</span>
                <span class="text-xs text-white/40">total</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="mt-8 flex flex-wrap items-center gap-3">
              <button
                type="button"
                class="glow-green inline-flex items-center gap-2.5 rounded-full bg-spotify px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-spotify-hover"
                @click="playAll"
              >
                <i aria-hidden="true" class="pi pi-play-fill" />
                Play all
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/4 px-6 py-3 text-sm font-bold text-white/80 transition hover:border-white/30 hover:bg-white/8 hover:text-white"
                @click="shuffleAll"
              >
                <i aria-hidden="true" class="pi pi-sort-alt" />
                Shuffle
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2.5 rounded-full border border-white/10 bg-white/3 px-5 py-3 text-sm font-bold text-white/60 transition hover:bg-white/8 hover:text-white"
                :class="{ 'border-spotify/30 text-spotify': isLiked }"
                @click="toggleLike"
              >
                <i aria-hidden="true" :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" />
                {{ isLiked ? 'Saved' : 'Save' }}
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2.5 rounded-full border border-white/10 bg-white/3 px-5 py-3 text-sm font-bold text-white/60 transition hover:bg-white/8 hover:text-white"
                @click="shareAlbum"
              >
                <i aria-hidden="true" class="pi pi-share-alt" />
                Share
              </button>
            </div>
          </div>
        </div>

        <!-- Section divider -->
        <div class="relative my-14">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-white/6" />
          </div>
          <div class="relative flex justify-center">
            <span class="bg-surface-base px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
              Tracklist
            </span>
          </div>
        </div>

        <!-- Tracklist -->
        <section aria-live="polite">
          <div class="space-y-1">
            <div
              v-for="(track, index) in tracks"
              :key="String(track.id)"
              role="button"
              tabindex="0"
              class="group flex cursor-pointer items-center gap-4 rounded-2xl px-4 py-3 transition-all duration-200 hover:bg-white/4"
              :class="isCurrentTrack(track) ? 'bg-white/6 ring-1 ring-inset ring-spotify/15' : ''"
              @click="playTrack(track, Number(index))"
              @keydown.enter="playTrack(track, Number(index))"
              @keydown.space.prevent="playTrack(track, Number(index))"
            >
              <!-- Track number / Play icon -->
              <span
                class="flex w-8 items-center justify-center text-center text-sm tabular-nums text-white/20"
                :class="isCurrentTrack(track) ? 'hidden' : 'group-hover:hidden'"
              >
                {{ String(Number(index) + 1).padStart(2, '0') }}
              </span>
              <span
                class="hidden w-8 items-center justify-center group-hover:flex"
                :class="isCurrentTrack(track) ? 'flex!' : ''"
              >
                <template v-if="isCurrentTrack(track)">
                  <span class="flex h-4 items-end gap-0.5">
                    <span class="eq-bar h-2" />
                    <span class="eq-bar animation-delay-150 h-4" />
                    <span class="eq-bar animation-delay-300 h-3" />
                  </span>
                </template>
                <i aria-hidden="true" v-else class="pi pi-play-fill text-xs text-white" />
              </span>

              <!-- Thumbnail -->
              <div class="relative h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10 ring-1 ring-white/4">
                <img
                  v-if="track.cover_url || album.cover_url"
                  :src="track.cover_url ?? album.cover_url ?? undefined"
                  :alt="track.title"
                  class="absolute inset-0 h-full w-full object-cover"
                  loading="lazy"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                </div>
              </div>

              <!-- Track info -->
              <div class="min-w-0 flex-1">
                <p
                  class="truncate text-sm font-semibold transition"
                  :class="isCurrentTrack(track) ? 'text-spotify' : 'text-white'"
                >
                  {{ track.title }}
                </p>
                <p class="mt-0.5 truncate text-xs text-white/40">
                  {{ track.artist_name || artist?.name || '' }}
                </p>
              </div>

              <!-- Duration -->
              <span class="text-xs tabular-nums text-white/25 transition group-hover:text-white/50">
                {{ formatDuration(track.duration_seconds) }}
              </span>
            </div>
          </div>

          <AppEmptyState
            v-if="!tracks.length"
            icon="pi pi-music"
            title="No tracks in this album"
            description="Tracks will appear here once they're added."
            bordered
          />
        </section>

        <!-- Credits -->
        <section v-if="albumArtists.length > 1" class="mt-14">
          <div class="relative mb-8">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/6" />
            </div>
            <div class="relative flex justify-center">
              <span class="bg-surface-base px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
                Credits
              </span>
            </div>
          </div>
          <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="aa in albumArtists"
              :key="aa.id"
              class="group flex items-center gap-4 rounded-2xl border border-white/4 bg-white/2 px-5 py-4 transition hover:border-white/8 hover:bg-white/4"
            >
              <div
                class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-linear-to-br from-white/8 to-white/2 text-sm font-bold text-white/70 ring-1 ring-white/4"
              >
                {{ aa.name?.charAt(0).toUpperCase() ?? '' }}
              </div>
              <div class="min-w-0">
                <RouterLink
                  v-if="aa.role !== 'main' && aa.role !== 'primary'"
                  :to="`/artist/${aa.id}`"
                  class="block truncate text-sm font-semibold text-white transition group-hover:text-spotify"
                >
                  {{ aa.name }}
                </RouterLink>
                <span v-else class="block truncate text-sm font-semibold text-white">{{ aa.name }}</span>
                <p v-if="aa.role && !['main', 'primary'].includes(aa.role)" class="mt-0.5 text-xs text-white/40 capitalize">
                  {{ aa.role }}
                </p>
                <p v-else-if="aa.role === 'main'" class="mt-0.5 text-xs text-white/30">Main artist</p>
              </div>
            </div>
          </div>
        </section>

        <!-- Related albums -->
        <section v-if="relatedAlbums.length > 0" class="mt-14">
          <div class="mb-8 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">More from</p>
              <h2 class="mt-1 text-xl font-bold text-white">{{ artist?.name ?? mainArtist?.name ?? 'Artist' }}</h2>
            </div>
            <RouterLink
              v-if="artist"
              :to="`/artist/${artist.id}`"
              class="text-xs font-medium text-white/30 transition hover:text-white"
            >
              View artist
              <i aria-hidden="true" class="pi pi-chevron-left ml-0.5 text-[10px]" />
            </RouterLink>
          </div>
          <div class="grid grid-cols-2 gap-5 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <RouterLink
              v-for="ra in relatedAlbums"
              :key="ra.id"
              :to="`/album/${ra.id}`"
              class="group block"
            >
              <div
                class="relative mb-3 aspect-square overflow-hidden rounded-2xl bg-white/6 shadow-lg ring-1 ring-white/10 transition-all duration-300 group-hover:shadow-lg group-hover:-translate-y-1 group-hover:ring-spotify/30"
              >
                <img
                  v-if="ra.cover_url"
                  :src="ra.cover_url"
                  :alt="ra.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-compact-disc text-2xl text-slate-500" />
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-spotify/90 text-black shadow-xl backdrop-blur-xs transition-transform group-hover:scale-110"
                  >
                    <i aria-hidden="true" class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>
              <p class="truncate text-sm font-semibold text-white">{{ ra.title }}</p>
              <p class="mt-0.5 truncate text-xs text-white/40">{{ ra.artist_name || 'Artist' }}</p>
            </RouterLink>
          </div>
        </section>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SkeletonLoader, AppEmptyState } from '@/components/common'
import { useAlbum } from '@/composables/catalog/useAlbum'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
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
  isLiked,
  totalDuration,
  formattedDuration,
  loading,
  error,
  fetchAlbum,
  toggleLike,
} = useAlbum(albumId)

const coverUrl = computed(() => album.value?.cover_url || null)
const { palette } = useAlbumColors(coverUrl)

const accentColor = computed(() => palette.value.vibrant || '#1db954')

const ambientBg = computed(() => {
  if (coverUrl.value!) return { background: '#06060A' }
  const c = accentColor.value
  return {
    background: `
      radial-gradient(ellipse 80% 50% at 50% 0%, ${c}1A 0%, transparent 70%),
      radial-gradient(ellipse 60% 40% at 100% 100%, ${c}0D 0%, transparent 50%),
      #06060A
    `,
  }
})

const coverGlowStyle = computed(() => {
  if (coverUrl.value!) return {}
  const c = accentColor.value
  return {
    boxShadow: `0 0 40px ${c}40, 0 0 80px ${c}20, 0 0 120px ${c}10`,
    transition: 'box-shadow 0.6s ease',
  }
})

const player = usePlayer()
const playerApi = usePlayerApi()
const router = useRouter()

function goBack() {
  router.back()
}

onMounted(() => {
  fetchAlbum()
})

function isCurrentTrack(track: Record<string, unknown>): boolean {
  return player.currentTrack.value?.id === String(track.id)
}

function formatDuration(seconds: number | null | undefined): string {
  if (seconds!) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

function playTrack(track: Record<string, unknown>, index: number) {
  if (tracks.value.length!) return
  const queue = buildQueue()
  player.setQueueAndPlay(queue, index)
}

function playAll() {
  if (tracks.value.length!) return
  const queue = buildQueue()
  player.setQueueAndPlay(queue, 0)
}

function shuffleAll() {
  if (tracks.value.length!) return
  const queue = buildQueue()
  const shuffled = [...queue].sort(() => Math.random() - 0.5)
  player.setQueueAndPlay(shuffled, 0)
}

function buildQueue() {
  if (tracks.value.length!) return []
  return tracks.value.map((t: Record<string, unknown>) => ({
    id: String(t.id),
    title: t.title as string,
    artistName: (t.artist_name as string) || artist?.value?.name || 'Unknown',
    albumTitle: (t.album_title as string) || album.value?.title || null,
    coverUrl: (t.cover_url as string) || album.value?.cover_url || null,
    durationSeconds: (t.duration_seconds as number) ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
  }))
}

const { copyLink } = useSocialShare()
function shareAlbum() {
  if (album.value!) return
  copyLink({
    id: album.value.id,
    title: album.value.title,
    type: 'album',
    artistName: artist.value?.name || album.value?.artist_name,
    coverUrl: album.value.cover_url,
  })
}
</script>

<style scoped>
.eq-bar {
  width: 3px;
  border-radius: 999px;
  background: currentColor;
  animation: equalizer 850ms ease-in-out infinite alternate;
}
.animation-delay-150 { animation-delay: 150ms; }
.animation-delay-300 { animation-delay: 300ms; }
@keyframes equalizer {
  from { transform: scaleY(0.45); opacity: 0.6; }
  to { transform: scaleY(1); opacity: 1; }
}
</style>
