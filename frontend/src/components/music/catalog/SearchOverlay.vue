<template>
  <Teleport to="body">
    <Transition name="search-fade">
      <div
        v-if="visible"
        class="fixed inset-0 z-[100] flex items-start justify-center bg-black/80 pt-16 backdrop-blur-sm md:pt-24"
        @click.self="close"
      >
        <div
          ref="panelRef"
          class="mx-4 w-full max-w-2xl overflow-hidden rounded-2xl bg-gradient-to-b from-[#1a1a2e] to-[#121212] shadow-2xl ring-1 ring-white/10"
        >
          <!-- Search input -->
          <div class="relative flex items-center border-b border-white/10 px-4">
            <i class="pi pi-search text-sm text-slate-400" />
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="Search tracks, artists, albums, playlists..."
              class="flex-1 bg-transparent px-3 py-4 text-sm text-white outline-none placeholder:text-slate-500"
              @keydown="onKeydown"
              @input="onInput"
            />
            <button
              v-if="query"
              type="button"
              class="spring mr-2 flex h-6 w-6 items-center justify-center rounded-full bg-white/10 text-xs text-slate-400 transition hover:bg-white/20"
              @click="clearQuery"
            >
              <i class="pi pi-times" />
            </button>
            <kbd
              class="hidden rounded-md border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] text-slate-500 md:inline-block"
            >
              ESC
            </kbd>
          </div>

          <!-- Recent / Trending (empty query) -->
          <div v-if="!query && !searching" class="max-h-[60vh] space-y-4 overflow-y-auto p-4">
            <div v-if="recentSearches.length">
              <p
                class="mb-2 flex items-center justify-between text-xs font-bold tracking-wider text-slate-400 uppercase"
              >
                <span>Recent</span>
                <button
                  type="button"
                  class="text-[10px] text-[#1db954] hover:underline"
                  @click="clearRecent"
                >
                  Clear
                </button>
              </p>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="term in recentSearches"
                  :key="term"
                  type="button"
                  class="spring flex items-center gap-2 rounded-full bg-white/[0.08] px-4 py-2 text-sm text-white transition-all hover:scale-105 hover:bg-white/[0.12]"
                  @click="query = term; doSearch()"
                >
                  <i class="pi pi-history text-xs text-slate-500" />
                  {{ term }}
                </button>
              </div>
            </div>

            <div>
              <p class="mb-2 text-xs font-bold tracking-wider text-slate-400 uppercase">
                Suggestions
              </p>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="suggestion in suggestions"
                  :key="suggestion"
                  type="button"
                  class="spring rounded-full bg-white/[0.06] px-4 py-2 text-sm text-slate-300 transition-all hover:scale-105 hover:bg-white/[0.10] hover:text-white"
                  @click="query = suggestion; doSearch()"
                >
                  {{ suggestion }}
                </button>
              </div>
            </div>

            <!-- Keyboard shortcut hint -->
            <div class="flex items-center gap-4 text-[11px] text-slate-600">
              <span
                ><kbd class="rounded border border-white/10 px-1.5 py-0.5 text-[10px]">↑↓</kbd>
                Navigate</span
              >
              <span
                ><kbd class="rounded border border-white/10 px-1.5 py-0.5 text-[10px]">↩</kbd>
                Select</span
              >
              <span
                ><kbd class="rounded border border-white/10 px-1.5 py-0.5 text-[10px]">Esc</kbd>
                Close</span
              >
            </div>
          </div>

          <!-- Searching indicator -->
          <div
            v-if="searching"
            class="flex items-center justify-center gap-3 p-12 text-sm text-slate-400"
          >
            <i class="pi pi-spin pi-spinner" />
            Searching...
          </div>

          <!-- Results -->
          <div v-if="query && !searching" class="max-h-[60vh] overflow-y-auto p-2">
            <div v-if="noResults" class="flex flex-col items-center gap-4 p-12 text-center">
              <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
                <i class="pi pi-search text-3xl text-slate-600" />
              </div>
              <p class="text-sm text-slate-400">
                No results for "<span class="font-medium text-white">{{ query }}</span
                >"
              </p>
              <p class="text-xs text-slate-500">Try a different search term</p>
            </div>

            <template v-else>
              <!-- Top Result (first track) -->
              <div v-if="results.tracks?.length" class="mb-4 px-2">
                <p class="mb-2 text-xs font-bold tracking-wider text-slate-400 uppercase">
                  Top Result
                </p>
                <div
                  class="group spring flex cursor-pointer items-center gap-4 rounded-xl bg-white/[0.04] p-3 transition-all hover:bg-white/[0.08]"
                  @click="selectTrack(results.tracks[0])"
                >
                  <div class="h-16 w-16 shrink-0 overflow-hidden rounded-xl bg-white/10 shadow-lg">
                    <img
                      v-if="results.tracks[0].cover_url"
                      :src="results.tracks[0].cover_url"
                      :alt="results.tracks[0].title"
                      loading="lazy"
                      class="h-full w-full object-cover"
                      @error="onImgError"
                    />
                    <div v-else class="flex h-full items-center justify-center">
                      <i class="pi pi-music text-lg text-slate-500" />
                    </div>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-base font-bold text-white">
                      {{ results.tracks[0].title }}
                    </p>
                    <p class="truncate text-sm text-slate-400">
                      {{ results.tracks[0].artist_name || 'Unknown' }}
                    </p>
                  </div>
                  <div
                    class="spring flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954]/0 text-white opacity-0 transition-all group-hover:bg-[#1db954] group-hover:opacity-100"
                  >
                    <i class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>

              <!-- Tracks -->
              <div v-if="results.tracks?.length" class="mb-3">
                <p class="mb-2 px-2 text-xs font-bold tracking-wider text-slate-400 uppercase">
                  Tracks
                </p>
                <div
                  v-for="(item, i) in results.tracks.slice(0, 5)"
                  :key="item.id"
                  :ref="(el) => setItemRef('track', i, el)"
                  class="spring flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 transition-all"
                  :class="
                    highlightedIndex === `track-${i}` ? 'bg-white/[0.10]' : 'hover:bg-white/[0.06]'
                  "
                  @click="selectTrack(item)"
                  @mouseenter="highlightedIndex = `track-${i}`"
                >
                  <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                    <img
                      v-if="item.cover_url"
                      :src="item.cover_url"
                      :alt="item.title"
                      loading="lazy"
                      class="h-full w-full object-cover"
                      @error="onImgError"
                    />
                    <div v-else class="flex h-full items-center justify-center">
                      <i class="pi pi-music text-xs text-slate-500" />
                    </div>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-medium text-white">{{ item.title }}</p>
                    <p class="truncate text-xs text-slate-400">
                      {{ item.artist_name || 'Unknown' }}
                    </p>
                  </div>
                  <span v-if="item.duration_seconds" class="text-xs text-slate-500 tabular-nums">{{
                    fmtDuration(item.duration_seconds)
                  }}</span>
                </div>
              </div>

              <!-- Artists -->
              <div v-if="results.artists?.length" class="mb-3">
                <p class="mb-2 px-2 text-xs font-bold tracking-wider text-slate-400 uppercase">
                  Artists
                </p>
                <div class="space-y-1">
                  <RouterLink
                    v-for="(item, i) in results.artists.slice(0, 3)"
                    :key="item.id"
                    :to="`/artist/${item.id}`"
                    :ref="(el) => setItemRef('artist', i, el)"
                    class="spring flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 transition-all"
                    :class="
                      highlightedIndex === `artist-${i}`
                        ? 'bg-white/[0.10]'
                        : 'hover:bg-white/[0.06]'
                    "
                    @click="close"
                    @mouseenter="highlightedIndex = `artist-${i}`"
                  >
                    <div
                      class="h-12 w-12 shrink-0 overflow-hidden rounded-full bg-white/10 ring-2 ring-white/10"
                    >
                      <img
                        v-if="item.cover_url"
                        :src="item.cover_url"
                        :alt="item.name"
                        loading="lazy"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div v-else class="flex h-full items-center justify-center">
                        <i class="pi pi-user text-sm text-slate-500" />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm font-medium text-white">{{ item.name }}</p>
                      <p class="truncate text-xs text-slate-400">
                        {{ item.genre ? `Artist · ${item.genre}` : 'Artist' }}
                      </p>
                    </div>
                    <i class="pi pi-chevron-left text-xs text-slate-500" />
                  </RouterLink>
                </div>
              </div>

              <!-- Albums -->
              <div v-if="results.albums?.length" class="mb-3">
                <p class="mb-2 px-2 text-xs font-bold tracking-wider text-slate-400 uppercase">
                  Albums
                </p>
                <div class="space-y-1">
                  <RouterLink
                    v-for="(item, i) in results.albums.slice(0, 3)"
                    :key="item.id"
                    :to="`/album/${item.id}`"
                    :ref="(el) => setItemRef('album', i, el)"
                    class="spring flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 transition-all"
                    :class="
                      highlightedIndex === `album-${i}`
                        ? 'bg-white/[0.10]'
                        : 'hover:bg-white/[0.06]'
                    "
                    @click="close"
                    @mouseenter="highlightedIndex = `album-${i}`"
                  >
                    <div class="h-12 w-12 shrink-0 overflow-hidden rounded-lg bg-white/10">
                      <img
                        v-if="item.cover_url"
                        :src="item.cover_url"
                        :alt="item.title"
                        loading="lazy"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div v-else class="flex h-full items-center justify-center">
                        <i class="pi pi-compact-disc text-sm text-slate-500" />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm font-medium text-white">{{ item.title }}</p>
                      <p class="truncate text-xs text-slate-400">
                        {{ item.artist_name || 'Album' }}
                      </p>
                    </div>
                    <i class="pi pi-chevron-left text-xs text-slate-500" />
                  </RouterLink>
                </div>
              </div>

              <!-- Playlists -->
              <div v-if="results.playlists?.length" class="mb-3">
                <p class="mb-2 px-2 text-xs font-bold tracking-wider text-slate-400 uppercase">
                  Playlists
                </p>
                <div class="space-y-1">
                  <RouterLink
                    v-for="(item, i) in results.playlists.slice(0, 3)"
                    :key="item.id"
                    :to="`/playlists/${item.id}`"
                    :ref="(el) => setItemRef('playlist', i, el)"
                    class="spring flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 transition-all"
                    :class="
                      highlightedIndex === `playlist-${i}`
                        ? 'bg-white/[0.10]'
                        : 'hover:bg-white/[0.06]'
                    "
                    @click="close"
                    @mouseenter="highlightedIndex = `playlist-${i}`"
                  >
                    <div
                      class="h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-gradient-to-br from-purple-500/20 to-purple-500/5"
                    >
                      <img
                        v-if="item.cover_url"
                        :src="item.cover_url"
                        :alt="item.name"
                        loading="lazy"
                        class="h-full w-full object-cover"
                        @error="onImgError"
                      />
                      <div v-else class="flex h-full items-center justify-center">
                        <i class="pi pi-list text-sm text-purple-400" />
                      </div>
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="truncate text-sm font-medium text-white">{{ item.name }}</p>
                      <p class="truncate text-xs text-slate-400">
                        {{ item.description || 'Playlist' }}
                      </p>
                    </div>
                    <i class="pi pi-chevron-left text-xs text-slate-500" />
                  </RouterLink>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useSearchApi } from '@/services/api/catalog/search'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { usePlayerApi } from '@/services/api/player'

const router = useRouter()
const player = usePlayer()
const playerApi = usePlayerApi()
const searchApi = useSearchApi()

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const _visible = ref(false)
const query = ref('')
const searching = ref(false)
const highlightedIndex = ref<string | null>(null)
const results = ref<{ tracks?: any[]; artists?: any[]; albums?: any[]; playlists?: any[] }>({})
const recentSearches = ref<string[]>([])
const noResults = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | null = null
const itemRefs: Record<string, HTMLElement> = {}

const inputRef = ref<HTMLInputElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)

const suggestions = [
  'Popular',
  'New releases',
  'Classical',
  'Electronic',
  'Top 50',
  'Trending',
  'Chill',
  'Workout',
  'Focus',
]

watch(
  () => props.visible,
  (v) => {
    _visible.value = v
    if (v) {
      query.value = ''
      results.value = {}
      noResults.value = false
      highlightedIndex.value = null
      loadRecent()
      nextTick(() => inputRef.value?.focus())
      window.history.pushState(null, '')
    }
  },
)

watch(_visible, (v) => {
  emit('update:visible', v)
})

function close() {
  _visible.value = false
}
function clearQuery() {
  query.value = ''
  results.value = {}
  noResults.value = false
  highlightedIndex.value = null
  nextTick(() => inputRef.value?.focus())
}

function loadRecent() {
  try {
    const raw = localStorage.getItem('music_recent_searches')
    recentSearches.value = raw ? JSON.parse(raw) : []
  } catch {
    recentSearches.value = []
  }
}

function saveRecent(term: string) {
  try {
    let list: string[] = JSON.parse(localStorage.getItem('music_recent_searches') || '[]')
    list = [term, ...list.filter((t) => t !== term)].slice(0, 8)
    localStorage.setItem('music_recent_searches', JSON.stringify(list))
    recentSearches.value = list
  } catch {}
}

function clearRecent() {
  localStorage.removeItem('music_recent_searches')
  recentSearches.value = []
}

function setItemRef(group: string, index: number, el: any) {
  if (el && el instanceof Element) itemRefs[`${group}-${index}`] = el as HTMLElement
}

let abortController: AbortController | null = null
async function doSearch() {
  const term = query.value.trim()
  if (!term) {
    results.value = {}
    noResults.value = false
    searching.value = false
    return
  }

  abortController?.abort()
  abortController = new AbortController()

  searching.value = true
  noResults.value = false

  try {
    const res = await searchApi.search({ q: term, limit: 20 }, undefined, { signal: abortController.signal })
    const data = res as any
    results.value = {
      tracks: data.tracks || [],
      artists: data.artists || [],
      albums: data.albums || [],
      playlists: data.playlists || [],
    }
    noResults.value =
      !results.value.tracks?.length &&
      !results.value.artists?.length &&
      !results.value.albums?.length &&
      !results.value.playlists?.length
    saveRecent(term)
  } catch (err) {
    if ((err as any)?.name === 'AbortError' || (err as any)?.code === 'ERR_CANCELED') return
    results.value = {}
    noResults.value = true
  } finally {
    searching.value = false
  }
}

function onInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  highlightedIndex.value = null
  debounceTimer = setTimeout(doSearch, 250)
}

function getFlatItems(): { group: string; index: number; el?: HTMLElement }[] {
  const flat: { group: string; index: number; el?: HTMLElement }[] = []
  for (const group of ['track', 'artist', 'album', 'playlist'] as const) {
    const key = group === 'playlist' ? 'playlists' : `${group}s`
    const items = (results.value as any)[key] || []
    items.forEach((_: any, i: number) =>
      flat.push({ group, index: i, el: itemRefs[`${group}-${i}`] }),
    )
  }
  return flat
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close()
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    moveHighlight(1)
    return
  }
  if (e.key === 'ArrowUp') {
    e.preventDefault()
    moveHighlight(-1)
    return
  }
  if (e.key === 'Enter') {
    e.preventDefault()
    activateHighlight()
    return
  }
}

function moveHighlight(dir: number) {
  const flat = getFlatItems()
  if (!flat.length) return
  const currentIdx = flat.findIndex(
    (f) => f.group && f.index && highlightedIndex.value === `${f.group}-${f.index}`,
  )
  let next = currentIdx === -1 ? (dir > 0 ? 0 : flat.length - 1) : currentIdx + dir
  if (next < 0) next = flat.length - 1
  if (next >= flat.length) next = 0
  const target = flat[next]
  if (!target) return
  highlightedIndex.value = `${target.group}-${target.index}`
  target.el?.scrollIntoView?.({ block: 'nearest' })
}

function activateHighlight() {
  if (!highlightedIndex.value) return
  const [group, indexStr] = highlightedIndex.value.split('-')
  const i = Number(indexStr)
  const key = group === 'track' ? 'tracks' : group === 'playlist' ? 'playlists' : `${group}s`
  const items = (results.value as any)[key] || []
  const item = items[i]
  if (!item) return
  if (group === 'track') selectTrack(item)
  else if (group === 'artist') {
    router.push(`/artist/${item.id}`)
    close()
  } else if (group === 'album') {
    router.push(`/album/${item.id}`)
    close()
  } else if (group === 'playlist') {
    router.push(`/playlists/${item.id}`)
    close()
  }
}

function selectTrack(track: any) {
  const playable = {
    id: String(track.id),
    title: track.title,
    artistName: track.artist_name || 'Unknown',
    albumTitle: track.album_title || null,
    coverUrl: track.cover_url || null,
    durationSeconds: track.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(track.id)),
  }
  player.playTrack(playable)
  if (track.artist_name) saveRecent(`${track.artist_name} - ${track.title}`)
  else saveRecent(track.title)
  close()
}

function fmtDuration(s: number) {
  const m = Math.floor(s / 60)
  const sec = s % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}

function onKeybind(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    _visible.value = !_visible.value
    if (_visible.value) nextTick(() => inputRef.value?.focus())
    return
  }
  if (
    e.key === '/' &&
    !['INPUT', 'TEXTAREA', 'SELECT'].includes((e.target as HTMLElement)?.tagName || '')
  ) {
    if (!_visible.value) {
      e.preventDefault()
      _visible.value = true
      nextTick(() => inputRef.value?.focus())
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeybind)
  window.addEventListener('popstate', close)
  loadRecent()
})

onUnmounted(() => {
  abortController?.abort()
  document.removeEventListener('keydown', onKeybind)
  window.removeEventListener('popstate', close)
})
</script>

<style scoped>
.search-fade-enter-active {
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}
.search-fade-leave-active {
  transition:
    opacity 120ms ease,
    transform 120ms ease;
}
.search-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}
.search-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
