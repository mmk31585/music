<template>
  <div class="relative mx-auto w-full max-w-7xl px-4 pt-6 pb-36 md:px-6 lg:px-8">
    <!-- Aurora gradient background -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1 -top-40 -left-40 bg-purple-600/15" />
      <div class="aurora-spot-2 -top-60 -right-40 bg-blue-500/10" />
      <div
        class="aurora-spot-1 top-20 left-1/3 bg-violet-500/10"
        style="animation-delay: -8s; width: 350px; height: 350px"
      />
      <div class="aurora-flow absolute inset-0 opacity-30" />
    </div>

    <div class="relative z-10">
      <HomeHero
        v-if="heroItems.length"
        :items="heroItems"
        @play="handleHeroPlay"
        @add-to-library="handleHeroAddToLibrary"
      />

      <div v-if="loading" class="mt-8 space-y-10">
        <div v-for="s in 4" :key="s">
          <SkeletonLoader variant="lines" :lines="1" class="w-40" />
          <div class="mt-5 flex gap-4">
            <SkeletonLoader v-for="i in 5" :key="i" variant="card" class="w-44 shrink-0" />
          </div>
        </div>
      </div>
      <template v-else-if="!hasData">
        <section
          class="relative mt-8 overflow-hidden rounded-2xl border border-white/6 bg-[#0C0C14] text-white shadow-2xl"
          dir="ltr"
        >
          <!-- Background layers -->
          <div class="absolute inset-0">
            <div
              class="absolute -top-24 -right-24 h-80 w-80 rounded-full bg-spotify/15 blur-[100px]"
            />
            <div
              class="absolute -bottom-16 -left-16 h-64 w-64 rounded-full bg-purple-600/10 blur-[80px]"
            />
            <div
              class="bg-blue-600/08 absolute top-1/2 left-1/2 h-48 w-48 -translate-x-1/2 -translate-y-1/2 rounded-full blur-[60px]"
            />
            <!-- Subtle grid -->
            <div
              class="absolute inset-0 opacity-[0.025]"
              style="
                background-image:
                  linear-gradient(#fff 1px, transparent 1px),
                  linear-gradient(90deg, #fff 1px, transparent 1px);
                background-size: 32px 32px;
              "
            />
          </div>

          <!-- Top accent line -->
          <div
            class="absolute top-0 right-0 left-0 h-px bg-linear-to-r from-transparent via-spotify/40 to-transparent"
          />

          <div class="relative z-10 px-8 py-12 md:px-12 md:py-16">
            <!-- Icon cluster -->
            <div class="mb-8 flex items-center gap-3">
              <div
                class="flex h-12 w-12 items-center justify-center rounded-2xl bg-spotify/10 ring-1 ring-spotify/20"
              >
                <i class="pi pi-headphones text-lg text-spotify" aria-hidden="true" />
              </div>
              <div class="flex gap-2">
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-xl bg-white/4 ring-1 ring-white/6"
                >
                  <i class="pi pi-bell text-xs text-white/30" aria-hidden="true" />
                </div>
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-xl bg-white/4 ring-1 ring-white/6"
                >
                  <i class="pi pi-star text-xs text-white/30" aria-hidden="true" />
                </div>
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-xl bg-white/4 ring-1 ring-white/6"
                >
                  <i class="pi pi-heart text-xs text-white/30" aria-hidden="true" />
                </div>
              </div>
            </div>

            <!-- Text -->
            <div class="max-w-lg">
              <p class="text-[10px] font-bold tracking-[0.35em] text-spotify/60 uppercase">
                Ready to listen
              </p>
              <h1 class="mt-3 text-3xl leading-tight font-black md:text-5xl">
                Your soundtrack<br />
                <span class="text-white/25">starts right here.</span>
              </h1>
              <p class="mt-4 text-sm leading-relaxed text-white/40">
                Discover new artists, save albums you love, and build playlists that match every
                mood.
              </p>
            </div>

            <!-- Divider -->
            <div class="my-8 h-px bg-linear-to-r from-white/6 to-transparent" />

            <!-- CTAs -->
            <div class="flex flex-wrap gap-3">
              <RouterLink
                to="/search"
                class="inline-flex items-center gap-2 rounded-full bg-spotify px-6 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-spotify-hover active:scale-95"
              >
                <i class="pi pi-search" aria-hidden="true" />
                Explore Music
              </RouterLink>
              <RouterLink
                to="/discover"
                class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/4 px-6 py-3 text-sm font-bold text-white/70 transition hover:border-white/20 hover:bg-white/7 hover:text-white active:scale-95"
              >
                <i class="pi pi-compass" aria-hidden="true" />
                Discover
              </RouterLink>
            </div>
          </div>
        </section>
      </template>

      <template v-else>
        <section v-if="personalized.length" class="mt-6">
          <HomeSectionHeader
            title="Based on Your Taste"
            eyebrow="Personalized picks"
            subtitle="Powered by your listening history"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in personalized"
              :key="item.id"
              :item="item"
              :delay="i * 30"
              :is-playing="isCurrentlyPlaying(item)"
              @play="handlePlay"
            />
          </HomeCarousel>
        </section>

        <section v-if="recentPlays.length" class="mt-6">
          <HomeSectionHeader
            title="Recently played"
            eyebrow="Jump back in"
            subtitle="Continue where you left off"
            see-all-route="/recently-played"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in recentPlays"
              :key="item.id"
              :item="item"
              :delay="i * 30"
              :is-playing="isCurrentlyPlaying(item)"
              @play="handlePlay"
            />
          </HomeCarousel>
        </section>

        <section v-if="popular.length" class="mt-12">
          <HomeSectionHeader
            title="Trending"
            eyebrow="Popular right now"
            see-all-route="/recommendations/popular"
          />
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

        <section v-if="forYou.length" class="mt-12">
          <HomeSectionHeader
            title="For You"
            eyebrow="Personalized picks"
            see-all-route="/recommendations/for-you"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in forYou"
              :key="item.id"
              :item="item"
              :delay="i * 30"
              :is-playing="isCurrentlyPlaying(item)"
              @play="handlePlay"
            />
          </HomeCarousel>
        </section>

        <section v-if="fromYourArtists.length" class="mt-12">
          <HomeSectionHeader
            title="More from your favorite artists"
            eyebrow="Based on your listening"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in fromYourArtists"
              :key="item.id"
              :item="item"
              :delay="i * 30"
              :is-playing="isCurrentlyPlaying(item)"
              @play="handlePlay"
            />
          </HomeCarousel>
        </section>

        <section v-for="sec in becauseOfSections" :key="sec.id" class="mt-12">
          <HomeSectionHeader :title="sec.title" eyebrow="Similar tracks" />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in sec.items"
              :key="item.id"
              :item="item"
              :delay="i * 30"
              :is-playing="isCurrentlyPlaying(item)"
              @play="handlePlay"
            />
          </HomeCarousel>
        </section>

        <section v-if="yourGenres.length" class="mt-12">
          <HomeSectionHeader title="Popular in your genres" eyebrow="From genres you listen to" />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in yourGenres"
              :key="item.id"
              :item="item"
              :delay="i * 30"
              :is-playing="isCurrentlyPlaying(item)"
              @play="handlePlay"
            />
          </HomeCarousel>
        </section>

        <section v-if="albums.length" class="mt-12">
          <HomeSectionHeader title="Albums" eyebrow="Collections" see-all-route="/search" />
          <HomeCarousel>
            <RouterLink
              v-for="album in albums"
              :key="album.id"
              :to="`/album/${album.id}`"
              class="group block w-44 shrink-0 space-y-3"
            >
              <div
                class="relative aspect-square overflow-hidden rounded-xl bg-white/6 shadow-lg ring-1 ring-white/10 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:ring-spotify/30"
              >
                <img
                  v-if="album.cover_url"
                  :src="album.cover_url"
                  :alt="album.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-compact-disc text-3xl text-slate-500" />
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-spotify text-black shadow-xl backdrop-blur-xs transition-transform group-hover:scale-110"
                  >
                    <i aria-hidden="true" class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>
              <div class="space-y-0.5 px-1">
                <p class="truncate text-sm font-semibold text-white">{{ album.title }}</p>
                <p v-if="album.artist_name" class="truncate text-xs text-white/40">
                  {{ album.artist_name }}
                </p>
                <p class="text-xs text-white/30">
                  {{ album.release_date?.slice(0, 4) || ''
                  }}{{ album.track_count ? ` \u2022 ${album.track_count} tracks` : '' }}
                </p>
              </div>
            </RouterLink>
          </HomeCarousel>
        </section>

        <section v-if="artists.length" class="mt-12">
          <HomeSectionHeader title="Featured artists" eyebrow="Meet the creators" />
          <HomeCarousel>
            <RouterLink
              v-for="artist in artists"
              :key="artist.id"
              :to="`/artist/${artist.id}`"
              class="group block w-40 shrink-0 space-y-3"
            >
              <div
                class="mx-auto h-36 w-36 overflow-hidden rounded-full bg-white/6 ring-1 ring-white/10 transition group-hover:ring-spotify/30"
              >
                <img
                  v-if="artist.image_url"
                  :src="artist.image_url"
                  :alt="artist.name"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-user text-3xl text-slate-500" />
                </div>
              </div>
              <div class="space-y-0.5 text-center">
                <p class="truncate text-sm font-semibold text-white">{{ artist.name }}</p>
                <p class="text-xs text-white/40">Artist</p>
              </div>
            </RouterLink>
          </HomeCarousel>
        </section>

        <section class="mt-12">
          <div class="grid gap-4 md:grid-cols-3">
            <RouterLink
              to="/recommendations"
              class="group rounded-2xl border border-white/6 bg-white/2 p-5 transition hover:-translate-y-0.5 hover:border-white/12 hover:bg-white/4"
            >
              <div
                class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/6 text-white/50"
              >
                <i aria-hidden="true" class="pi pi-star text-xl" />
              </div>
              <h3 class="mt-4 text-sm font-bold text-white">Recommendations</h3>
              <p class="mt-1 text-sm text-white/40">
                Discover popular, recent, and personalized tracks.
              </p>
            </RouterLink>
            <RouterLink
              to="/discover"
              class="group rounded-2xl border border-white/6 bg-white/2 p-5 transition hover:-translate-y-0.5 hover:border-white/12 hover:bg-white/4"
            >
              <div
                class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/6 text-white/50"
              >
                <i aria-hidden="true" class="pi pi-compass text-xl" />
              </div>
              <h3 class="mt-4 text-sm font-bold text-white">Discover</h3>
              <p class="mt-1 text-sm text-white/40">Explore trending tracks and new releases.</p>
            </RouterLink>
            <RouterLink
              to="/ai/mood-explorer"
              class="group rounded-2xl border border-white/6 bg-white/2 p-5 transition hover:-translate-y-0.5 hover:border-white/12 hover:bg-white/4"
            >
              <div
                class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/6 text-white/50"
              >
                <i aria-hidden="true" class="pi pi-magic text-xl" />
              </div>
              <h3 class="mt-4 text-sm font-bold text-white">Mood Explorer</h3>
              <p class="mt-1 text-sm text-white/40">Find music that matches your vibe.</p>
            </RouterLink>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useHomeFeed } from '@/composables/useHomeFeed'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import type { RecommendationTrack } from '@/services/api/recommendation'
import { onImgError } from '@/utils/helpers'
import { HomeHero, HomeCarousel, HomeTrackCard } from '@/components/music'
import HomeSectionHeader from '@/components/music/home/HomeSectionHeader.vue'
import type { HeroItem } from '@/components/music/home/HomeHero.vue'

const {
  popular,
  forYou,
  personalized,
  recentPlays,
  fromYourArtists,
  becauseOfSections,
  yourGenres,
  albums,
  artists,
  loading,
  hasData,
  fetchHomeFeed,
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

function getScoreBadge(item: RecommendationTrack): string | undefined {
  if (item.score != null && item.score >= 90) return '🔥 Hot'
  if (item.score != null && item.score >= 75) return 'Trending'
  return undefined
}

function isCurrentlyPlaying(item: RecommendationTrack): boolean {
  const currentId = player.currentTrack.value?.id
  if (currentId!) return false
  return currentId === item.id
}

function buildPlaybackTrack(item: RecommendationTrack): PlaybackTrack {
  return {
    id: item.id,
    title: item.title || 'Untitled',
    artistName: item.artist_name || 'Unknown artist',
    albumTitle: item.album_title || null,
    coverUrl: item.cover_url || null,
    durationSeconds: item.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(item.id),
  }
}

function handlePlay(item: any) {
  void player.toggleTrack(buildPlaybackTrack(item))
}

function handleHeroPlay(item: HeroItem) {
  const track = popular.value.find((t) => String(t.id) === item.id)
  if (track) {
    void handlePlay(track)
  }
}

const toast = useToast()
function handleHeroAddToLibrary() {
  toast.add({
    severity: 'info',
    summary: 'Coming Soon',
    detail: 'Library management is coming soon!',
    life: 3000,
  })
}

onMounted(() => {
  void fetchHomeFeed()
})
</script>
