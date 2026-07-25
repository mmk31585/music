import { ref, computed, onUnmounted } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useAlbumsApi } from '@/services/api/catalog/albums'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { useUserAuthStore } from '@/stores'
import type { RecommendationTrack, HomeFeedSection } from '@/services/api/recommendation'
import type { Album } from '@/services/api/catalog/albums'
import type { Artist } from '@/services/api/catalog/artists'

export function useHomeFeed() {
  const recsApi = useRecommendationsApi()
  const albumsApi = useAlbumsApi()
  const artistsApi = useArtistsApi()
  const auth = useUserAuthStore()

  const sections = ref<HomeFeedSection[]>([])
  const albums = ref<Album[]>([])
  const artists = ref<Artist[]>([])
  const loading = ref(false)
  const error = ref<unknown>(null)
  const sectionErrors = ref<Record<string, boolean>>({})
  const sectionLoading = ref(false)

  function isUnauthorized(err: unknown): boolean {
    if (err && typeof err === 'object' && 'response' in err) {
      const resp = (err as any).response
      return resp?.status === 401
    }
    return false
  }

  const hasData = computed(() =>
    sections.value.length > 0 ||
    albums.value.length > 0 ||
    artists.value.length > 0
  )

  const personalized = ref<RecommendationTrack[]>([])
  const sectionMap = computed(() => {
    const map: Record<string, HomeFeedSection> = {}
    for (const s of sections.value) {
      map[s.id] = s
    }
    return map
  })

  const recentPlays = computed(() => sectionMap.value.recently_played?.items ?? [])
  const popular = computed(() => sectionMap.value.trending?.items ?? [])
  const forYou = computed(() => sectionMap.value.for_you?.items ?? [])
  const fromYourArtists = computed(() => sectionMap.value.from_your_artists?.items ?? [])
  const becauseOfSections = computed(() =>
    sections.value.filter(s => s.id.startsWith('because_of_'))
  )
  const yourGenres = computed(() => sectionMap.value.your_genres?.items ?? [])

  async function fetchHomeFeed() {
    loading.value = true
    error.value = null

    try {
      const results = await Promise.allSettled([
        auth.isAuthenticated
          ? recsApi.getHomeFeed().catch(() => null)
          : Promise.resolve(null),
        auth.isAuthenticated
          ? recsApi.getPersonalized({ limit: 10 }).catch(() => null)
          : Promise.resolve(null),
        albumsApi.getAlbums().catch(() => [] as Album[]),
        artistsApi.getArtists().catch(() => [] as Artist[]),
      ])

      const homeResult = results[0].status === 'fulfilled' ? results[0].value : null
      sections.value = homeResult?.sections ?? []

      const personalizedResult = results[1].status === 'fulfilled' ? results[1].value : null
      personalized.value = personalizedResult?.items ?? []

      albums.value = results[2].status === 'fulfilled' ? results[2].value ?? [] : []
      artists.value = results[3].status === 'fulfilled' ? results[3].value ?? [] : []

      const homeRejected = results[0].status === 'rejected' && isUnauthorized(results[0].reason)
      const personalizedRejected = results[1].status === 'rejected' && isUnauthorized(results[1].reason)
      if (homeRejected || personalizedRejected) {
        personalized.value = []
      }

      sectionErrors.value = {}
      if (results[0].status === 'rejected') sectionErrors.value.recs = true
      if (results[1].status === 'rejected') sectionErrors.value.personalized = true
      if (results[2].status === 'rejected') sectionErrors.value.albums = true
      if (results[3].status === 'rejected') sectionErrors.value.artists = true
    } catch (err: unknown) {
      error.value = err
    } finally {
      loading.value = false
      sectionLoading.value = false
    }
  }

  let refreshTimer: ReturnType<typeof setInterval> | null = null
  function startPeriodicRefresh(intervalMs = 5 * 60 * 1000) {
    stopPeriodicRefresh()
    refreshTimer = setInterval(() => {
      fetchHomeFeed()
    }, intervalMs)
  }
  function stopPeriodicRefresh() {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }
  onUnmounted(stopPeriodicRefresh)

  return {
    sections,
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
    sectionLoading,
    error,
    hasData,
    sectionErrors,
    startPeriodicRefresh,
    stopPeriodicRefresh,
    fetchHomeFeed,
  }
}
