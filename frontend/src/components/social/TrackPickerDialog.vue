<template>
  <Teleport to="body">
    <Transition name="track-picker-fade">
      <div
        v-if="visible"
        class="fixed inset-0 z-[100] flex items-start justify-center bg-black/70 pt-16 backdrop-blur-sm md:pt-24"
        @click.self="close"
      >
        <div
          class="mx-4 w-full max-w-xl overflow-hidden rounded-2xl bg-gradient-to-b from-[#1a1a2e] to-[#121212] shadow-2xl ring-1 ring-white/10"
        >
          <div class="flex items-center justify-between border-b border-white/10 px-4 py-3">
            <h2 class="text-sm font-bold text-white">{{ title }}</h2>
            <button
              type="button"
              class="flex h-7 w-7 items-center justify-center rounded-full bg-white/10 text-xs text-slate-400 transition hover:bg-white/20"
              @click="close"
            >
              <i class="pi pi-times" />
            </button>
          </div>

          <div class="relative flex items-center border-b border-white/10 px-4">
            <i class="pi pi-search text-sm text-slate-400" />
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="Search tracks..."
              class="flex-1 bg-transparent px-3 py-3 text-sm text-white outline-none placeholder:text-slate-500"
              @input="onInput"
              @keydown="onKeydown"
            />
            <i v-if="searching" class="pi pi-spin pi-spinner text-xs text-slate-400" />
          </div>

          <div class="max-h-72 overflow-y-auto p-2">
            <div v-if="!query" class="flex items-center justify-center py-12 text-xs text-slate-500">
              <i class="pi pi-headphones mr-2" /> Type to search for tracks
            </div>

            <div v-else-if="searching" class="flex items-center justify-center py-12 text-xs text-slate-400">
              <div class="flex flex-col items-center gap-2">
                <div class="flex gap-1">
                  <div class="h-2 w-2 animate-bounce rounded-full bg-slate-500 [animation-delay:0ms]" />
                  <div class="h-2 w-2 animate-bounce rounded-full bg-slate-500 [animation-delay:150ms]" />
                  <div class="h-2 w-2 animate-bounce rounded-full bg-slate-500 [animation-delay:300ms]" />
                </div>
                <span>Searching...</span>
              </div>
            </div>

            <div v-else-if="!results.length" class="flex items-center justify-center py-12 text-xs text-slate-500">
              <i class="pi pi-info-circle mr-2" /> No tracks found for "{{ query }}"
            </div>

            <div v-else class="space-y-1">
              <div
                v-for="(track, i) in results"
                :key="track.id"
                role="option"
                :aria-selected="focusedIdx === i"
                class="group flex cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 transition-all"
                :class="focusedIdx === i ? 'bg-white/[0.12] ring-1 ring-white/20' : 'hover:bg-white/[0.08]'"
                @click="selectTrack(track)"
                @mouseenter="focusedIdx = i"
              >
                <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                  <img
                    v-if="track.cover_url"
                    :src="track.cover_url"
                    :alt="track.title"
                    loading="lazy"
                    class="h-full w-full object-cover"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i class="pi pi-headphones text-xs text-slate-500" />
                  </div>
                  <div class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 transition group-hover:opacity-100">
                    <i class="pi pi-plus text-xs text-white" />
                  </div>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                  <p class="truncate text-xs text-slate-400">{{ track.artist_name || 'Unknown' }}</p>
                </div>
                <span v-if="track.duration_seconds" class="shrink-0 text-xs text-slate-500 tabular-nums">
                  {{ fmtDuration(track.duration_seconds) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted } from 'vue'
import { useSearchApi } from '@/services/api/catalog/search'
import type { Track } from '@/services/api/catalog/tracks'

const props = defineProps<{
  visible: boolean
  title?: string
}>()

const emit = defineEmits<{
  select: [track: Track]
  'update:visible': [value: boolean]
}>()

const searchApi = useSearchApi()

const query = ref('')
const results = ref<Track[]>([])
const searching = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | null = null
let abortController: AbortController | null = null
const inputRef = ref<HTMLInputElement | null>(null)
const focusedIdx = ref(-1)

watch(() => props.visible, (v) => {
  if (v) {
    focusedIdx.value = -1
    nextTick(() => {
      inputRef.value?.focus()
      inputRef.value?.select()
    })
  } else {
    query.value = ''
    results.value = []
  }
})

function onInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(doSearch, 250)
  focusedIdx.value = -1
}

async function doSearch() {
  const term = query.value.trim()
  if (!term) {
    results.value = []
    searching.value = false
    return
  }

  abortController?.abort()
  abortController = new AbortController()

  searching.value = true
  try {
    const res = await searchApi.searchCatalog({ query: term, type: 'tracks', limit: 10 }, { signal: abortController.signal } as any)
    results.value = res.tracks ?? []
  } catch (err) {
    if ((err as any)?.name === 'AbortError' || (err as any)?.code === 'ERR_CANCELED') return
    results.value = []
  } finally {
    searching.value = false
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (focusedIdx.value < results.value.length - 1) focusedIdx.value++
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (focusedIdx.value > 0) focusedIdx.value--
  } else if (e.key === 'Enter' && focusedIdx.value >= 0) {
    e.preventDefault()
    const track = results.value[focusedIdx.value]
    if (track) selectTrack(track)
  } else if (e.key === 'Escape') {
    close()
  }
}

function selectTrack(track: Track) {
  emit('select', track)
  query.value = ''
  results.value = []
}

onUnmounted(() => {
  abortController?.abort()
})

function close() {
  emit('update:visible', false)
}

function fmtDuration(s: number) {
  const m = Math.floor(s / 60)
  const sec = s % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}
</script>

<style scoped>
.track-picker-fade-enter-active {
  transition: opacity 160ms ease, transform 160ms ease;
}
.track-picker-fade-leave-active {
  transition: opacity 120ms ease, transform 120ms ease;
}
.track-picker-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}
.track-picker-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}
</style>
