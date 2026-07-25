<template>
  <div
    class="select-none"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
  >
    <!-- Pull-to-refresh indicator -->
    <div
      v-if="pulling || refreshing"
      class="flex justify-center"
      :style="{ height: `${pullDistance || (refreshing ? 40 : 0)}px`, overflow: 'hidden', transition: refreshing ? 'height 0.3s ease' : 'none' }"
    >
      <div
        v-if="refreshing"
        class="mt-3 h-5 w-5 animate-spin rounded-full border-2 border-spotify border-t-transparent"
      />
      <div
        v-else
        class="mt-3 h-5 w-5 rounded-full border-2 border-white/20"
        :style="{
          transform: `rotate(${pullDistance * 3}deg)`,
          transition: 'transform 0.1s ease',
        }"
      />
    </div>

    <HomeSkeleton v-if="loading" />

    <template v-else-if="error">
      <div class="flex flex-col items-center justify-center gap-4 py-20">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-12 w-12 text-white/20"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        <p class="text-white/50">Could not load content</p>
        <button
          type="button"
          class="cursor-pointer rounded-full bg-white/10 px-6 py-2 text-sm font-medium text-white transition hover:bg-white/20"
          @click="fetchHomeFeed"
        >
          Try again
        </button>
      </div>
    </template>

    <template v-else>
      <HomeHero
        :greeting="greeting"
        :subtitle="greetingSubtitle"
        :period-label="periodLabel"
        :art-src="heroArtSrc"
        @explore="router.push('/search')"
        @discover="router.push('/search')"
      />

      <ContinueListening
        v-if="auth.isAuthenticated && resumeItems.length"
        :items="resumeItems"
        @play="handlePlay"
      />

      <MoodDiscovery @select="handleMoodSelect" />

      <QuickActionsGrid />

      <MiniWidgets />

      <TrendingChart
        v-if="auth.isAuthenticated && trendingChart.length"
        :entries="trendingChart"
        @play="handlePlay"
      />

      <ArtistSpotlight
        v-if="spotlights.length"
        :artists="spotlights"
        @play-artist="handlePlayArtist"
      />

      <AIDiscoveries
        v-if="personalized.length"
        :tracks="personalized.slice(0, 8)"
        @play="handlePlay"
      />

      <EditorialDiscover
        v-if="forYou.length"
        :items="forYou.slice(0, 6)"
        @play="handlePlay"
      />

      <PersianHighlights
        v-if="albums.length"
        :items="albums.slice(0, 8)"
      />

      <CommunityActivity
        v-if="auth.isAuthenticated"
        :items="[]"
      />

      <!-- Spacer for player bar -->
      <div class="h-24" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useHomePage } from '@/composables/useHomePage'
import HomeHero from '@/pages/app/home/HomeHero.vue'
import ContinueListening from '@/pages/app/home/ContinueListening.vue'
import MoodDiscovery from '@/pages/app/home/MoodDiscovery.vue'
import QuickActionsGrid from '@/pages/app/home/QuickActionsGrid.vue'
import MiniWidgets from '@/pages/app/home/MiniWidgets.vue'
import TrendingChart from '@/pages/app/home/TrendingChart.vue'
import ArtistSpotlight from '@/pages/app/home/ArtistSpotlight.vue'
import AIDiscoveries from '@/pages/app/home/AIDiscoveries.vue'
import EditorialDiscover from '@/pages/app/home/EditorialDiscover.vue'
import PersianHighlights from '@/pages/app/home/PersianHighlights.vue'
import CommunityActivity from '@/pages/app/home/CommunityActivity.vue'
import HomeSkeleton from '@/pages/app/home/HomeSkeleton.vue'

const {
  greeting,
  greetingPeriod,
  greetingSubtitle,
  heroItems,
  fallbackHeroItems,
  resumeItems,
  trendingChart,
  spotlights,
  personalized,
  forYou,
  albums,
  loading,
  error,
  isCurrentlyPlaying,
  handlePlay,
  getScoreBadge,
  fetchHomeFeed,
  player,
  auth,
  remainingPlays,
  popular,
} = useHomePage()

const router = useRouter()
const toast = useToast()

const periodLabel = computed(() => {
  const map: Record<string, string> = {
    morning: '☀️ Morning',
    afternoon: '🌤 Afternoon',
    evening: '🌅 Evening',
    night: '🌙 Night',
  }
  return map[greetingPeriod.value] || 'Featured'
})

const heroArtSrc = computed(() => {
  if (popular.value.length) {
    return popular.value[0]?.cover_url || undefined
  }
  return undefined
})

function handleMoodSelect(moodId: string) {
  router.push(`/ai/mood-explorer?mood=${moodId}`)
}

function handlePlayArtist(artist: any) {
  const track = popular.value.find((t) => t.artist_id === String(artist.id))
  if (track) handlePlay(track)
}

// ── Pull-to-refresh ──────────────────────────────────────────────
const touchStartY = ref(0)
const pulling = ref(false)
const pullDistance = ref(0)
const refreshing = ref(false)
const PULL_THRESHOLD = 80

function onTouchStart(e: TouchEvent) {
  if (window.scrollY > 10) return
  touchStartY.value = e.touches[0]?.clientY ?? 0
  pulling.value = true
}

function onTouchMove(e: TouchEvent) {
  if (!pulling.value) return
  const delta = (e.touches[0]?.clientY ?? 0) - touchStartY.value
  if (delta > 0) {
    pullDistance.value = Math.min(delta * 0.5, 120)
  }
}

function onTouchEnd() {
  if (!pulling.value) return
  pulling.value = false
  const distance = pullDistance.value
  pullDistance.value = 0
  if (distance >= PULL_THRESHOLD) {
    refreshing.value = true
    void fetchHomeFeed().finally(() => {
      refreshing.value = false
    })
  }
}

onMounted(() => {
  void fetchHomeFeed()
  const route = router.currentRoute.value
  if (route.query.onboarded === 'true') {
    toast.add({
      severity: 'success',
      summary: 'Welcome to Muse!',
      detail: 'Your music preferences are set. Start exploring!',
      life: 5000,
    })
    router.replace({ query: {} })
  }
})
</script>
