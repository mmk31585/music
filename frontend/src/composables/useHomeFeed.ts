import { ref, computed } from 'vue'
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
    } catch (err: unknown) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

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
    error,
    hasData,
    fetchHomeFeed,
  }
}
