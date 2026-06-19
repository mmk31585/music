import { computed, ref } from 'vue'
import { useSearchApi, type SearchResult } from '@/services/api/catalog/search'
import type { PlaylistListItem } from '@/services/api/playlist'

type SearchResultWithPlaylists = SearchResult & { playlists: PlaylistListItem[] }

export function useCatalogSearch() {
  const { searchCatalog } = useSearchApi()

  type SearchType = 'all' | 'artists' | 'albums' | 'tracks'
  const query = ref('')
  const type = ref<SearchType>('all')
  const results = ref<SearchResultWithPlaylists>({
    tracks: [],
    artists: [],
    albums: [],
    playlists: [],
  })

  const loading = ref(false)
  const error = ref<any>(null)

  const hasQuery = computed(() => query.value.trim().length > 0)
  const totalResults = computed(
    () =>
      results.value.tracks.length +
      results.value.artists.length +
      results.value.albums.length +
      results.value.playlists.length,
  )
  const hasResults = computed(() => totalResults.value > 0)

  async function runSearch() {
    const q = query.value.trim()

    if (!q) {
      results.value = {
        tracks: [],
        artists: [],
        albums: [],
        playlists: [],
      }
      return results.value
    }

    loading.value = true
    error.value = null

    try {
      const response = await searchCatalog({
        query: q,
        type: type.value,
        limit: 30,
      })

      results.value = response as SearchResultWithPlaylists
      return response
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearSearch() {
    query.value = ''
    results.value = {
      tracks: [],
      artists: [],
      albums: [],
      playlists: [],
    }
    error.value = null
  }

  return {
    query,
    type,
    results,
    loading,
    error,
    hasQuery,
    totalResults,
    hasResults,
    runSearch,
    clearSearch,
  }
}
