import { ref } from 'vue'

const STORAGE_KEY = 'muse_recent_searches'
const MAX_ITEMS = 5

export interface SearchEntry {
  query: string
  timestamp: number
}

function loadFromStorage(): SearchEntry[] {
  if (typeof window === 'undefined') return []
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as SearchEntry[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function saveToStorage(entries: SearchEntry[]) {
  if (typeof window === 'undefined') return
  localStorage.setItem(STORAGE_KEY, JSON.stringify(entries))
}

const recentSearches = ref<SearchEntry[]>(loadFromStorage())

export function useRecentSearches() {
  function addRecentSearch(query: string) {
    const trimmed = query.trim()
    if (!trimmed) return

    const filtered = recentSearches.value.filter(
      (entry) => entry.query.toLowerCase() !== trimmed.toLowerCase()
    )

    recentSearches.value = [
      { query: trimmed, timestamp: Date.now() },
      ...filtered,
    ].slice(0, MAX_ITEMS)

    saveToStorage(recentSearches.value)
  }

  function removeRecentSearch(query: string) {
    recentSearches.value = recentSearches.value.filter(
      (entry) => entry.query.toLowerCase() !== query.toLowerCase()
    )
    saveToStorage(recentSearches.value)
  }

  function clearRecentSearches() {
    recentSearches.value = []
    saveToStorage(recentSearches.value)
  }

  return {
    recentSearches,
    addRecentSearch,
    removeRecentSearch,
    clearRecentSearches,
  }
}
