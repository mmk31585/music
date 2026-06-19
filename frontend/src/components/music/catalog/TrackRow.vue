<template>
  <div
    class="group grid grid-cols-[48px_1fr_auto] items-center gap-4 rounded-xl px-3 py-2.5 transition-all duration-200 hover:bg-white/[0.08]"
    :class="isCurrent ? 'bg-white/[0.10] shadow-[inset_3px_0_0_#1db954]' : ''"
    @contextmenu.prevent="openContextMenu"
  >
    <button
      type="button"
      class="relative flex h-11 w-11 items-center justify-center overflow-hidden rounded-xl bg-white/10 text-white transition-all duration-200 hover:scale-105 hover:bg-[#1db954] hover:text-black disabled:cursor-wait disabled:opacity-70"
      :disabled="loadingThisTrack"
      @click="handlePlay"
    >
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="title"
        class="absolute inset-0 h-full w-full object-cover opacity-60 transition group-hover:opacity-35"
        loading="lazy"
        @error="onImgError"
      />

      <span class="relative z-10 flex items-center justify-center">
        <i aria-hidden="true" v-if="loadingThisTrack" class="pi pi-spin pi-spinner text-sm" />

        <span
          v-else-if="isCurrent && player.isPlaying.value"
          class="flex h-4 items-end gap-[2px]"
          aria-label="Playing"
        >
          <span class="eq-bar h-2" />
          <span class="eq-bar animation-delay-150 h-4" />
          <span class="eq-bar animation-delay-300 h-3" />
        </span>

        <i aria-hidden="true" v-else :class="buttonIcon" class="text-sm" />
      </span>
    </button>

    <div class="min-w-0">
      <div
        class="truncate text-sm font-semibold transition"
        :class="isCurrent ? 'text-[#1db954]' : 'text-white'"
      >
        {{ title }}
      </div>

      <div class="mt-0.5 truncate text-xs text-slate-400">
        {{ artistName }}
      </div>
    </div>

    <div class="flex items-center gap-4 text-xs text-slate-400">
      <span v-if="durationLabel" class="tabular-nums">
        {{ durationLabel }}
      </span>

      <button
        type="button"
        aria-label="More options"
        class="hidden rounded-full p-2 text-slate-400 transition group-hover:block hover:bg-white/10 hover:text-white"
        @click.stop="openContextMenu"
      >
        <i aria-hidden="true" class="pi pi-ellipsis-h" />
      </button>
    </div>
  </div>

  <AppContextMenu
    :visible="menuVisible"
    :items="menuItems"
    :position="{ x: menuX, y: menuY }"
    @close="menuVisible = false"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import AppContextMenu from '@/components/common/AppContextMenu.vue'

interface TrackRowTrack {
  id: string | number
  title?: string | null
  artistName?: string | null
  artist_name?: string | null
  artist_id?: string | number | null
  artist?: { name?: string | null } | null
  artists?: { name?: string | null }[] | null
  coverUrl?: string | null
  cover_url?: string | null
  cover?: string | null
  album?: { title?: string | null; coverUrl?: string | null; cover_url?: string | null; id?: string | number } | null
  albumTitle?: string | null
  album_title?: string | null
  durationSeconds?: number | null
  duration_seconds?: number | null
  duration?: number | null
}

const props = defineProps<{
  track: TrackRowTrack
  index?: number
  queue?: TrackRowTrack[]
}>()

const player = usePlayer()
const playerApi = usePlayerApi()

const title = computed(() => props.track.title || 'Untitled')

const artistName = computed(() => {
  return (
    props.track.artistName ||
    props.track.artist_name ||
    props.track.artist?.name ||
    (Array.isArray(props.track.artists) && props.track.artists[0]?.name) ||
    props.track.artist ||
    'Unknown artist'
  )
})

const coverUrl = computed(() => {
  return (
    props.track.coverUrl ||
    props.track.cover_url ||
    props.track.cover ||
    props.track.album?.coverUrl ||
    props.track.album?.cover_url ||
    null
  )
})

const durationSeconds = computed(() => {
  return props.track.durationSeconds ?? props.track.duration_seconds ?? props.track.duration ?? null
})

const durationLabel = computed(() => {
  const total = Number(durationSeconds.value)
  if (!Number.isFinite(total) || total <= 0) return ''

  const minutes = Math.floor(total / 60)
  const seconds = Math.floor(total % 60)

  return `${minutes}:${String(seconds).padStart(2, '0')}`
})

const trackId = computed(() => String(props.track.id))

const isCurrent = computed(() => player.currentTrack.value?.id === trackId.value)

const loadingThisTrack = computed(() => {
  return isCurrent.value && (player.isLoadingTrack.value || player.isBuffering.value)
})

const router = useRouter()

const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const menuItems = [
  { label: 'Play Now', icon: 'pi pi-play', action: () => handlePlay() },
  { label: 'Play next', icon: 'pi pi-step-forward', action: () => {
    if (player.queue.value.length > 0) {
      const items = [...player.queue.value]
      items.splice(0, 0, buildPlaybackTrack())
      player.updateQueue(items)
    } else {
      player.updateQueue([buildPlaybackTrack()])
    }
  }},
  { label: 'Add to queue', icon: 'pi pi-list', action: () => player.updateQueue([...player.queue.value, buildPlaybackTrack()]) },
  { label: 'Go to artist', icon: 'pi pi-user', action: () => {
    const artistId = props.track.artist_id
    if (artistId) router.push(`/artist/${artistId}`)
  }},
  { label: 'Go to album', icon: 'pi pi-book', action: () => {
    if (props.track.album?.id) router.push(`/album/${props.track.album.id}`)
  }},
]

function openContextMenu(e: MouseEvent) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  menuVisible.value = true
}

const buttonIcon = computed(() => {
  if (isCurrent.value && player.isPlaying.value) return 'pi pi-pause'
  return 'pi pi-play'
})

function buildPlaybackTrack(): PlaybackTrack {
  const id = trackId.value
  return {
    id,
    title: title.value,
    artistName: String(artistName.value),
    albumTitle:
      props.track.albumTitle || props.track.album_title || props.track.album?.title || null,
    coverUrl: coverUrl.value,
    durationSeconds: durationSeconds.value,
    streamUrl: playerApi.getTrackStreamUrl(id),
  }
}

async function handlePlay() {
  await player.toggleTrack(buildPlaybackTrack())
}
</script>

<style scoped>
.eq-bar {
  width: 3px;
  border-radius: 999px;
  background: currentColor;
  animation: equalizer 850ms ease-in-out infinite alternate;
  will-change: transform, opacity;
}

.animation-delay-150 {
  animation-delay: 150ms;
}

.animation-delay-300 {
  animation-delay: 300ms;
}

@keyframes equalizer {
  from { transform: scaleY(0.45); opacity: 0.6; }
  to { transform: scaleY(1); opacity: 1; }
}
</style>
