<template>
  <div
    class="group grid grid-cols-[48px_1fr_auto] items-center gap-4 rounded-xl px-3 py-2.5 transition-all duration-200 hover:bg-surface-active"
    :class="isCurrent ? 'bg-surface-active shadow-[inset_3px_0_0_var(--accent)]' : ''"
    @contextmenu.prevent="openContextMenu"
  >
    <button
      type="button"
      class="relative flex h-11 w-11 items-center justify-center overflow-hidden rounded-xl bg-surface-active text-primary transition-all duration-200 hover:scale-105 hover:bg-accent hover:text-black disabled:cursor-wait disabled:opacity-70"
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
        <Loader2 aria-hidden="true" v-if="loadingThisTrack" class="text-sm animate-spin"  />

        <span
          v-else-if="isCurrent && player.isPlaying.value"
          class="flex h-4 items-end gap-0.5"
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
      <RouterLink
        :to="`/track/${trackId}`"
        class="truncate text-sm font-semibold transition hover:underline"
        :class="isCurrent ? 'text-accent' : 'text-primary'"
        @click.stop
      >
        {{ title }}
      </RouterLink>

      <div class="mt-0.5 truncate text-xs text-secondary">
        {{ artistName }}
      </div>
    </div>

    <div class="flex items-center gap-4 text-xs text-secondary">
      <span v-if="durationLabel" class="tabular-nums">
        {{ durationLabel }}
      </span>

      <button
        type="button"
        aria-label="More options"
        class="hidden rounded-full p-2 text-secondary transition group-hover:block hover:bg-surface-active hover:text-primary"
        @click.stop="openContextMenu"
      >
        <MoreHorizontal aria-hidden="true" class=""  />
      </button>
    </div>
  </div>

  <ContextMenu
    v-model:visible="menuVisible"
    :sections="menuSections"
    :position="{ x: menuX, y: menuY }"
  />
</template>

<script setup lang="ts">
import { Loader2, MoreHorizontal } from 'lucide-vue-next'
import { computed, inject, ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { buildPlaybackTrack } from '@/factories/playbackTrack'
import ContextMenu from '@/components/common/ContextMenu.vue'
import type { ContextMenuSection } from '@/types/context-menu'

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
  album_id?: string | number | null
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
  if (Number.isFinite!(total) || total <= 0) return ''

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

const openRadioFromTrack = inject<(trackId: string, seedLabel?: string) => void>('openRadio', () => {})

const queueTrackInput = () => ({
  id: trackId.value,
  title: title.value,
  artist_name: String(artistName.value),
  album_title: props.track.albumTitle || props.track.album_title || props.track.album?.title || null,
  cover_url: coverUrl.value,
  duration_seconds: durationSeconds.value,
})

const menuSections = computed<ContextMenuSection[]>(() => [
  {
    id: 'playback',
    label: 'PLAYBACK',
    items: [
      {
        id: 'play-now',
        label: 'Play Now',
        icon: 'Play',
        action: () => handlePlay(),
      },
      {
        id: 'play-next',
        label: 'Play Next',
        icon: 'SkipForward',
        action: () => player.playNextInQueue(buildPlaybackTrack(queueTrackInput())),
      },
      {
        id: 'add-to-queue',
        label: 'Add to Queue',
        icon: 'ListMusic',
        action: () => player.addToQueue(buildPlaybackTrack(queueTrackInput())),
      },
      {
        id: 'start-radio',
        label: 'Start Radio',
        icon: 'Radio',
        separator: true,
        action: () => {
          openRadioFromTrack(trackId.value, `${title.value} • ${artistName.value}`)
        },
      },
    ],
  },
  {
    id: 'navigate',
    label: 'GO TO',
    items: [
      {
        id: 'go-to-track',
        label: 'Go to Track',
        icon: 'Music2',
        separator: true,
        action: () => router.push(`/track/${trackId.value}`),
      },
      {
        id: 'go-to-artist',
        label: 'Go to Artist',
        icon: 'UserRound',
        hidden: !props.track.artist_id,
        action: () => {
          const artistId = props.track.artist_id
          if (artistId) router.push(`/artist/${artistId}`)
        },
      },
      {
        id: 'go-to-album',
        label: 'Go to Album',
        icon: 'Disc3',
        action: () => {
          const albumId = props.track.album?.id || props.track.album_id
          if (albumId) router.push(`/album/${albumId}`)
        },
      },
    ],
  },
])

function openContextMenu(e: MouseEvent) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  menuVisible.value = true
}

const buttonIcon = computed(() => {
  if (isCurrent.value && player.isPlaying.value) return 'pi pi-pause'
  return 'pi pi-play'
})

async function handlePlay() {
  await player.toggleTrack(buildPlaybackTrack({
    id: trackId.value,
    title: title.value,
    artist_name: String(artistName.value),
    album_title: props.track.albumTitle || props.track.album_title || props.track.album?.title || null,
    cover_url: coverUrl.value,
    duration_seconds: durationSeconds.value,
  }))
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
