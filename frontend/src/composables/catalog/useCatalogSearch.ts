/**
 * Persian/Arabic character normalization map for typo-tolerant search.
 * Maps visually similar characters to a canonical form.
 */
const PERSIAN_NORMALIZE: Record<string, string> = {
  // Yeh variants → standard ی
  'ي': 'ی', 'ى': 'ی', 'ێ': 'ی', 'ۍ': 'ی',
  // Kaf variants → standard ک
  'ك': 'ک', 'ڪ': 'ک', 'ګ': 'ک',
  // Heh variants → standard ه
  'ة': 'ه', 'ۀ': 'ه',
  // Other common Persian substitutions
  'إ': 'ا', 'أ': 'ا', 'آ': 'ا',
  'ؤ': 'و',
  'ئ': 'ی',
  'ٱ': 'ا',
  'ے': 'ی',
  'ۓ': 'ی',
}

/** Normalize a Persian search query for typo-tolerant matching. */
function normalizePersianQuery(q: string): string {
  let normalized = ''
  for (const ch of q) {
    normalized += PERSIAN_NORMALIZE[ch] ?? ch
  }
  return normalized.trim()
}

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
    const raw = query.value.trim()
    // Normalize Persian characters for typo-tolerant search
    const q = normalizePersianQuery(raw)

    if (!q) {
      results.value = {
        tracks: [],
        artists: [],
        albums: [],
        playlists: [],
      }
      return results.value
    }

    // Keep the original query in the input but send normalized version to API
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
