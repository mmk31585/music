<template>
  <aside
    class="hidden h-screen shrink-0 border-e border-border-subtle bg-surface-raised/40 backdrop-blur-2xl xl:flex xl:flex-col transition-all duration-300 ease-out-expo z-30"
    :class="collapsed ? 'w-14 items-center' : 'w-88'"
  >
    <div class="flex h-full flex-col overflow-hidden" :class="collapsed ? 'items-center' : ''">
      <div class="flex shrink-0 items-center px-5 py-4" :class="collapsed ? 'flex-col gap-4' : 'gap-3 border-b border-border-subtle w-full'">
        <button
          type="button"
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg transition-all duration-150 hover:bg-surface-active active:scale-90"
          :class="currentTrack ? 'bg-accent-subtle' : 'bg-surface-overlay'"
          :title="collapsed ? 'Show now playing panel' : 'Hide now playing panel'"
          @click="toggleCollapsed"
        >
          <ChevronRight
            v-if="collapsed"
            :size="14"
            :class="currentTrack ? 'text-accent' : 'text-tertiary'"
          />
          <ChevronLeft
            v-else
            :size="14"
            :class="currentTrack ? 'text-accent' : 'text-tertiary'"
          />
        </button>
        <template v-if="!collapsed">
          <div class="flex h-6 w-6 items-center justify-center rounded-lg bg-accent-subtle">
            <Music :size="12" class="text-accent" />
          </div>
          <h2 class="text-xs font-bold tracking-[0.2em] text-secondary uppercase">Now Playing</h2>
        </template>
        <template v-else>
          <div v-if="currentTrack" class="flex flex-col items-center gap-2">
            <div class="h-10 w-10 overflow-hidden rounded-xl ring-1 ring-border-default">
              <img v-if="currentTrack.coverUrl" :src="currentTrack.coverUrl" :alt="currentTrack.title" class="h-full w-full object-cover" />
              <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-accent/30 to-accent-hover/30">
                <Headphones :size="16" class="text-tertiary" />
              </div>
            </div>
            <span v-if="isPlaying" class="h-1.5 w-1.5 rounded-full bg-accent animate-pulse" />
          </div>
        </template>
      </div>

      <template v-if="!collapsed">
        <div class="flex-1 space-y-5 overflow-y-auto p-4 pb-32 scroll-smooth scroll-bar">
          <RightPaneNowPlaying
            v-if="currentTrack"
            :track="currentTrack"
            :is-playing="isPlaying"
            :is-loading="isLoadingTrack"
            :is-buffering="isBuffering"
            :liked="liked"
            :current-time="currentTime"
            :duration="duration"
            :progress-percent="progressPercent"
            :shuffle-mode="shuffleMode"
            :repeat-mode="repeatMode"
            :has-next="hasNext"
            :has-previous="hasPrevious"
            :album-cover="albumCover"
            :accent-color="accentColor"
            @toggle-like="toggleLike"
            @toggle-fullscreen="$emit('toggle-fullscreen')"
            @toggle-play="togglePlayPause"
            @toggle-shuffle="toggleShuffle"
            @toggle-repeat="toggleRepeat"
            @prev="playPrevious"
            @next="playNext"
            @seek="onSeek"
          />
          <RightPaneNowPlayingEmpty v-else />

          <RightPaneQueue
            :items="queue"
            @play="playFromQueue"
            @clear="clearQueue"
          />

          <RightPaneSuggestions
            :items="suggestions"
            @play="playSuggestion"
          />

          <div class="flex items-center justify-center pt-2">
            <button
              type="button"
              class="flex items-center gap-2 rounded-xl px-4 py-2 text-[11px] font-bold text-tertiary transition-all duration-150 hover:bg-surface-hover hover:text-primary active:scale-95"
              @click="$emit('toggle-queue-overlay')"
            >
              <MoreHorizontal :size="12" />
              <span>Show Queue &rarr;</span>
            </button>
          </div>

          <div class="h-2" />
        </div>
      </template>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ChevronLeft, ChevronRight, Headphones, MoreHorizontal, Music } from 'lucide-vue-next'
import { usePlayerStore } from '@/stores/player'
import { usePlayerControls, useTrackLike } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { queueManager } from '@/services/player/queue-manager'
import type { RecommendationTrack } from '@/services/api/recommendation/types'
import RightPaneNowPlaying from './RightPaneNowPlaying.vue'
import RightPaneNowPlayingEmpty from './RightPaneNowPlayingEmpty.vue'
import RightPaneQueue from './RightPaneQueue.vue'
import RightPaneSuggestions from './RightPaneSuggestions.vue'
import type { SuggestionItem } from './RightPaneSuggestions.vue'

defineEmits<{
  'toggle-fullscreen': []
  'toggle-queue-overlay': []
}>()

const playerStore = usePlayerStore()
const pc = usePlayerControls()

const collapsed = ref(localStorage.getItem('rightpane-collapsed') === 'true')
function toggleCollapsed() {
  collapsed.value = !collapsed.value
  localStorage.setItem('rightpane-collapsed', String(collapsed.value))
  window.dispatchEvent(new CustomEvent('rightpane-collapse', { detail: collapsed.value }))
}

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  currentTime,
  duration,
  progressPercent,
  shuffleMode,
  repeatMode,
  hasNext,
  hasPrevious,
  togglePlayPause,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
  seekPercent,
} = pc

const queue = computed(() => playerStore.queue)
const trackId = computed(() => currentTrack.value?.id)
const { liked, toggleLike } = useTrackLike(trackId)

const albumCover = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(albumCover)

const accentColor = computed(() => {
  return albumCover.value && palette.value.vibrant ? palette.value.vibrant : null
})

const suggestions = ref<SuggestionItem[]>([])

watch(currentTrack, async (track) => {
  suggestions.value = []
  if (track?.id) {
    try {
      const recsApi = useRecommendationsApi()
      const result = await recsApi.getSimilar(track.id, { limit: 3 })
      const items = result?.items ?? []
      suggestions.value = items.map((item: RecommendationTrack) => ({
        id: item.id,
        title: item.title,
        artistName: item.artist_name || 'Unknown',
        coverUrl: item.cover_url,
        durationSeconds: item.duration_seconds,
      }))
    } catch (err) {
      console.error('Failed to fetch suggestions:', err)
    }
  }
}, { immediate: true })

function playFromQueue(index: number) {
  const track = queue.value[index]
  if (track) {
    playerStore.playTrack(track)
  }
}

function playSuggestion(id: string) {
  playerStore.playTrackById(id)
}

function clearQueue() {
  queueManager.clear()
  playerStore.updateQueue([])
}

function onSeek(pct: number) {
  seekPercent(pct)
}
</script>
