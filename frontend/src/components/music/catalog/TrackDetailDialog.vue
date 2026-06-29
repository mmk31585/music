<template>
  <Dialog
    :visible="visible"
    modal
    :draggable="false"
    :style="{ maxWidth: '480px', width: '92vw' }"
    :pt="{
      root: 'border-none',
      mask: 'backdrop-blur-xs',
      header: 'border-b border-white/5 p-0',
      title: 'text-white text-sm font-bold',
      content: 'p-0',
    }"
    @update:visible="emit('update:visible', $event)"
  >
    <template #header>
      <div class="flex items-center gap-2 px-1" />
    </template>

    <div class="flex flex-col">
      <!-- Cover art header -->
      <div class="relative h-48 w-full overflow-hidden bg-white/4">
        <img
          v-if="coverUrl"
          :src="coverUrl"
          :alt="title"
          class="h-full w-full object-cover"
          @error="onImgError"
        />
        <div v-else class="flex h-full items-center justify-center">
          <i aria-hidden="true" class="pi pi-compact-disc text-4xl text-slate-500" />
        </div>

        <!-- Dark gradient overlay for text readability -->
        <div class="absolute inset-0 bg-linear-to-t from-surface-base via-surface-base/60 to-transparent" />

        <!-- Track info overlaying cover -->
        <div class="absolute bottom-4 left-5 right-5">
          <h2 class="text-xl font-black text-white drop-shadow-lg">{{ title }}</h2>
          <p v-if="artistName" class="mt-1 text-sm font-medium text-white/80 drop-shadow-md">
            {{ artistName }}
          </p>
        </div>
      </div>

      <!-- Body -->
      <div class="flex flex-col gap-5 p-5">
        <!-- Album + Duration + Explicit -->
        <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-white/50">
          <RouterLink
            v-if="albumId && albumTitle"
            :to="`/album/${albumId}`"
            class="flex items-center gap-1.5 font-medium text-white/60 transition hover:text-white"
            @click="emit('update:visible', false)"
          >
            <i aria-hidden="true" class="pi pi-book text-[10px]" />
            {{ albumTitle }}
          </RouterLink>

          <span v-if="durationSeconds" class="tabular-nums">
            <i aria-hidden="true" class="pi pi-clock mr-1 text-[10px]" />
            {{ formatDuration(durationSeconds) }}
          </span>

          <span
            v-if="explicit"
            class="inline-flex h-4 w-4 items-center justify-center rounded bg-white/15 text-[9px] font-bold tracking-wide text-white"
          >E</span>

          <span v-if="playCount && playCount > 0" class="text-white/40">
            {{ formatCount(playCount) }} plays
          </span>
        </div>

        <!-- Genres -->
        <div v-if="genres.length > 0" class="flex flex-wrap gap-1.5">
          <span
            v-for="g in genres"
            :key="g.id ?? g.name"
            class="rounded-full border border-white/10 bg-white/6 px-2.5 py-0.5 text-[11px] font-medium text-slate-300"
          >
            {{ g.name || g.id || '' }}
          </span>
        </div>

        <!-- Divider -->
        <div class="border-t border-white/6" />

        <!-- Action buttons -->
        <div class="flex flex-wrap gap-2">
          <button
            type="button"
            class="glow-green inline-flex items-center gap-2 rounded-full bg-spotify px-5 py-2 text-sm font-bold text-black transition hover:bg-spotify-hover"
            :class="{ 'bg-white/20! text-white! shadow-none!': isCurrentTrack }"
            @click="handlePlay"
          >
            <i aria-hidden="true" :class="isCurrentTrack && player.isPlaying.value ? 'pi pi-pause-fill' : 'pi pi-play-fill'" />
            {{ isCurrentTrack && player.isPlaying.value ? 'Pause' : 'Play' }}
          </button>

          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/4 px-4 py-2 text-sm font-bold text-white/80 transition hover:border-white/30 hover:bg-white/8 hover:text-white"
            @click="addToQueue"
          >
            <i aria-hidden="true" class="pi pi-list text-xs" />
            Queue
          </button>

          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/3 px-4 py-2 text-sm font-bold transition"
            :class="liked ? 'border-spotify/30 text-spotify' : 'text-white/60 hover:bg-white/8 hover:text-white'"
            @click="emit('toggle-like')"
          >
            <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" />
            {{ liked ? 'Liked' : 'Like' }}
          </button>

          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/3 px-4 py-2 text-sm font-bold text-white/60 transition hover:bg-white/8 hover:text-white"
            @click="shareTrack"
          >
            <i aria-hidden="true" class="pi pi-share-alt" />
            Share
          </button>
        </div>

        <!-- Quick links -->
        <div class="flex items-center gap-4 text-xs text-white/40">
          <RouterLink
            v-if="artistId"
            :to="`/artist/${artistId}`"
            class="flex items-center gap-1 transition hover:text-white"
            @click="emit('update:visible', false)"
          >
            <i aria-hidden="true" class="pi pi-user" />
            View artist
          </RouterLink>
          <RouterLink
            v-if="albumId"
            :to="`/album/${albumId}`"
            class="flex items-center gap-1 transition hover:text-white"
            @click="emit('update:visible', false)"
          >
            <i aria-hidden="true" class="pi pi-book" />
            View album
          </RouterLink>
        </div>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import { useSocialShare } from '@/composables/social'
import { onImgError } from '@/utils/helpers'
import { formatDuration } from '@/utils/format'

interface GenreItem {
  id?: string | number
  name?: string
  [key: string]: unknown
}

const props = withDefaults(defineProps<{
  visible: boolean
  trackId: string
  title: string
  artistName?: string | null
  artistId?: string | number | null
  albumTitle?: string | null
  albumId?: string | number | null
  coverUrl?: string | null
  durationSeconds?: number | null
  explicit?: boolean
  playCount?: number | null
  genres?: GenreItem[]
  liked?: boolean
}>(), {
  genres: () => [],
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'toggle-like': []
}>()

const player = usePlayer()
const playerApi = usePlayerApi()
const { copyLink } = useSocialShare()

const isCurrentTrack = computed(() => player.currentTrack.value?.id === props.trackId)

function buildPlaybackTrack(): PlaybackTrack {
  return {
    id: props.trackId,
    title: props.title,
    artistName: props.artistName || 'Unknown',
    albumTitle: props.albumTitle || null,
    coverUrl: props.coverUrl || null,
    durationSeconds: props.durationSeconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(props.trackId),
  }
}

function handlePlay() {
  const pb = buildPlaybackTrack()
  if (isCurrentTrack.value && player.isPlaying.value) {
    player.pause()
  } else {
    player.playTrack(pb)
  }
  emit('update:visible', false)
}

function addToQueue() {
  const pb = buildPlaybackTrack()
  player.updateQueue([...player.queue.value, pb])
  emit('update:visible', false)
}

function shareTrack() {
  copyLink({
    id: props.trackId,
    title: props.title,
    type: 'track',
    artistName: props.artistName,
    coverUrl: props.coverUrl,
  })
}

function formatCount(count?: number | null): string {
  if (!count) return '0'
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}
</script>
