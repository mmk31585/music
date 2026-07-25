import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useHomeFeed } from '@/composables/useHomeFeed'
import { usePlayer } from '@/composables/player'
import { useUserAuthStore } from '@/stores'
import { useGuestSession } from '@/composables/useGuestSession'
import { buildPlaybackTrack } from '@/factories/playbackTrack'
import type { RecommendationTrack } from '@/services/api/recommendation/types'
import type { Album } from '@/services/api/catalog/albums/types'
import type { Artist } from '@/services/api/catalog/artists/types'
import type { HeroItem } from '@/components/music/home/HomeHero.vue'

export interface TrendingEntry {
  rank: number
  track: RecommendationTrack
  previousRank: number
  peakRank: number
  isNew: boolean
}

export type GreetingPeriod = 'morning' | 'afternoon' | 'evening' | 'night'

export function useHomePage() {
  const feed = useHomeFeed()
  const player = usePlayer()
  const auth = useUserAuthStore()
  const { remainingPlays } = useGuestSession()

  const now = ref(new Date())
  let tick: ReturnType<typeof setInterval> | null = null

  onMounted(() => {
    tick = setInterval(() => { now.value = new Date() }, 60_000)
  })
  onUnmounted(() => {
    if (tick) clearInterval(tick)
  })

  const greetingPeriod = computed<GreetingPeriod>(() => {
    const h = now.value.getHours()
    if (h < 12) return 'morning'
    if (h < 17) return 'afternoon'
    if (h < 21) return 'evening'
    return 'night'
  })

  const greeting = computed(() => {
    const name = auth.user?.displayName || auth.user?.username || ''
    const periodMap: Record<GreetingPeriod, string> = {
      morning: 'صبح‌تون بخیر',
      afternoon: 'عصر بخیر',
      evening: 'عصر بخیر',
      night: 'شب بخیر',
    }
    const base = periodMap[greetingPeriod.value]
    return name ? `${base}, ${name}` : base
  })

  const greetingSubtitle = computed(() => {
    const periodMap: Record<GreetingPeriod, string> = {
      morning: 'شروع روز با موزیک دلخواه',
      afternoon: 'یک عصر موزیکال خوب',
      evening: 'حال و هوای شبونه',
      night: 'شب آرومی داشته باشی',
    }
    return periodMap[greetingPeriod.value]
  })

  const heroItems = computed<HeroItem[]>(() => {
    const tracks = feed.popular.value.slice(0, 5)
    return tracks.map((t, i) => ({
      id: String(t.id),
      title: t.title || 'Untitled',
      subtitle: t.artist_name || 'Unknown artist',
      image: t.cover_url || '',
      badge: i === 0 ? 'Trending' : i === 1 ? 'Popular' : 'Featured',
      badgeVariant: i === 0 ? 'purple' : 'green',
      type: 'album' as const,
    }))
  })

  const fallbackHeroItems = computed<HeroItem[]>(() => [
    {
      id: 'welcome',
      title: greeting.value,
      subtitle: greetingSubtitle.value,
      image: '',
      badge: 'Featured',
      badgeVariant: 'green' as const,
      type: 'album' as const,
    },
    {
      id: 'explore',
      title: 'Explore Your Library',
      subtitle: 'Save tracks, albums, and artists',
      image: '',
      badge: 'Tip',
      badgeVariant: 'purple' as const,
      type: 'album' as const,
    },
    {
      id: 'discover',
      title: 'New Releases',
      subtitle: 'Check out the latest music',
      image: '',
      badge: 'New',
      badgeVariant: 'green' as const,
      type: 'album' as const,
    },
  ])

  function isCurrentlyPlaying(item: RecommendationTrack): boolean {
    const currentId = player.currentTrack.value?.id
    if (!currentId) return false
    return currentId === item.id
  }

  function handlePlay(item: any) {
    void player.toggleTrack(buildPlaybackTrack(item))
  }

  function handleHeroPlay(item: HeroItem) {
    const track = feed.popular.value.find((t) => String(t.id) === item.id)
    if (track) void handlePlay(track)
  }

  function getScoreBadge(item: RecommendationTrack): string | undefined {
    if (item.score != null && item.score >= 90) return '🔥 Hot'
    if (item.score != null && item.score >= 75) return 'Trending'
    return undefined
  }

  const trendingChart = computed<TrendingEntry[]>(() => {
    const items = feed.popular.value.slice(0, 10)
    return items.map((track, i) => ({
      rank: i + 1,
      track,
      previousRank: i + 1 + Math.floor(Math.random() * 5) - 2,
      peakRank: 1 + Math.floor(Math.random() * 3),
      isNew: Math.random() > 0.7,
    }))
  })

  const resumeItems = computed(() => feed.recentPlays.value.slice(0, 6))

  const spotlights = computed(() => feed.artists.value.slice(0, 4))

  return {
    ...feed,
    greeting,
    greetingPeriod,
    greetingSubtitle,
    heroItems,
    fallbackHeroItems,
    isCurrentlyPlaying,
    handlePlay,
    handleHeroPlay,
    getScoreBadge,
    trendingChart,
    resumeItems,
    spotlights,
    player,
    auth,
    remainingPlays,
  }
}
