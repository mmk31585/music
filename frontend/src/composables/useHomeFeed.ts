import { ref, computed } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useAlbumsApi } from '@/services/api/catalog/albums'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { useLibraryApi } from '@/services/api/library'
import { useUserAuthStore } from '@/stores'
import type { RecommendationTrack } from '@/services/api/recommendation'
import type { Album } from '@/services/api/catalog/albums'
import type { Artist } from '@/services/api/catalog/artists'

export function useHomeFeed() {
  const recsApi = useRecommendationsApi()
  const albumsApi = useAlbumsApi()
  const artistsApi = useArtistsApi()
  const libraryApi = useLibraryApi()
  const auth = useUserAuthStore()

  const popular = ref<RecommendationTrack[]>([])
  const forYou = ref<RecommendationTrack[]>([])
  const recent = ref<RecommendationTrack[]>([])
  const albums = ref<Album[]>([])
  const artists = ref<Artist[]>([])
  const recentPlays = ref<any[]>([])
  const loading = ref(false)
  const error = ref<any>(null)

  const hasData = computed(() =>
    popular.value.length > 0 ||
    forYou.value.length > 0 ||
    recent.value.length > 0 ||
    albums.value.length > 0 ||
    artists.value.length > 0
  )

  async function fetchHomeFeed() {
    loading.value = true
    error.value = null

    try {
      const results = await Promise.allSettled([
        recsApi.getPopular({ limit: 10 }).catch(() => null),
        recsApi.getForYou({ limit: 10 }).catch(() => null),
        recsApi.getRecent({ limit: 10 }).catch(() => null),
        albumsApi.getAlbums().catch(() => [] as Album[]),
        artistsApi.getArtists().catch(() => [] as Artist[]),
        auth.isAuthenticated
          ? libraryApi.getRecentlyPlayed().catch(() => [])
          : Promise.resolve([]),
      ])

      const popResult = results[0].status === 'fulfilled' ? results[0].value : null
      popular.value = popResult?.items ?? []

      const forYouResult = results[1].status === 'fulfilled' ? results[1].value : null
      forYou.value = forYouResult?.items ?? []

      const recentResult = results[2].status === 'fulfilled' ? results[2].value : null
      recent.value = recentResult?.items ?? []

      albums.value = results[3].status === 'fulfilled' ? results[3].value ?? [] : []
      artists.value = results[4].status === 'fulfilled' ? results[4].value ?? [] : []
      recentPlays.value = results[5].status === 'fulfilled' ? results[5].value ?? [] : []
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  return {
    popular,
    forYou,
    recent,
    albums,
    artists,
    recentPlays,
    loading,
    error,
    hasData,
    fetchHomeFeed,
  }
}
