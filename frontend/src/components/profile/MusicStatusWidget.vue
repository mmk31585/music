<template>
  <div
    class="overflow-hidden rounded-2xl border border-white/6 bg-white/2 backdrop-blur-xs"
    dir="rtl"
  >
    <!-- Header -->
    <div class="flex items-center justify-between px-4 pt-4 pb-2">
      <div class="flex items-center gap-2">
        <span class="relative flex h-2 w-2">
          <span
            v-if="isPlayingAny"
            class="absolute inline-flex h-full w-full animate-ping rounded-full bg-spotify opacity-75"
          />
          <span
            :class="isPlayingAny ? 'bg-spotify' : 'bg-slate-500'"
            class="relative inline-flex h-2 w-2 rounded-full"
          />
        </span>
        <span class="text-xs font-bold text-white/50">
          {{ isOwnProfile ? 'در حال گوش دادن' : 'در حال پخش' }}
        </span>
      </div>
    </div>

    <!-- Body -->
    <div class="px-4 pb-4">
      <!-- No music playing -->
      <div v-if="!isPlayingAny && isOwnProfile" class="py-3 text-sm text-white/30">
        الان موزیکی پخش نمی‌شه
      </div>

      <!-- Playing track (owner or visible visitor) -->
      <div v-else-if="currentTrackData" class="flex items-center gap-3 rounded-xl px-3 py-3 backdrop-blur-xs">
        <!-- Album art -->
        <div class="h-11 w-11 shrink-0 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="currentTrackData.coverUrl || currentTrackData.cover_url"
            :src="(currentTrackData.coverUrl || currentTrackData.cover_url)!"
            :alt="currentTrackData.title"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full items-center justify-center">
            <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
          </div>
        </div>

        <!-- Track info -->
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-white">
            {{ currentTrackData.title }}
          </p>
          <p class="mt-0.5 truncate text-xs text-white/40">
            {{ currentTrackData.artistName || currentTrackData.artist_name || '' }}
          </p>
        </div>

        <!-- Play button (visitor) -->
        <button
          v-if="!isOwnProfile"
          type="button"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-xs transition hover:bg-spotify hover:text-black"
          @click="playTrack"
          aria-label="Play track"
        >
          <i aria-hidden="true" class="pi pi-play-fill text-sm" />
        </button>
      </div>

      <!-- Privacy controls (owner only) -->
      <div v-if="isOwnProfile" class="mt-2 flex items-center justify-between">
        <span class="text-[10px] font-medium text-white/30">نمایش برای:</span>
        <div class="flex gap-1.5">
          <button
            v-for="opt in privacyOptions"
            :key="opt.value"
            type="button"
            class="rounded-full px-2.5 py-1 text-[10px] font-bold transition"
            :class="privacy === opt.value
              ? 'bg-spotify/20 text-spotify'
              : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'"
            @click="setPrivacy(opt.value)"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { usePlayerStore } from '@/stores/player'
import { useVideoApi } from '@/services/api/video'
import { useTracksApi } from '@/services/api/catalog/tracks'

const props = defineProps<{
  isOwnProfile: boolean
  musicStatus?: {
    playing: boolean
    current_track_id?: string | null
    updated_at?: string | null
  } | null
}>()

const playerStore = usePlayerStore()
const videoApi = useVideoApi()
const tracksApi = useTracksApi()

type PrivacyLevel = 'public' | 'followers' | 'private'

const privacyOptions = [
  { value: 'public' as PrivacyLevel, label: 'همه' },
  { value: 'followers' as PrivacyLevel, label: 'فالوورها' },
  { value: 'private' as PrivacyLevel, label: 'خصوصی' },
]

const privacy = ref<PrivacyLevel>(
  (localStorage.getItem('music-status-privacy') as PrivacyLevel) || 'public',
)
const isPlayingAny = ref(false)
const currentTrackData = ref<{
  id: string
  title: string
  artistName?: string | null
  artist_name?: string | null
  coverUrl?: string | null
  cover_url?: string | null
} | null>(null)

// ── Lifecycle ───────────────────────────────────────────────────────────

onMounted(() => {
  if (props.isOwnProfile) {
    isPlayingAny.value = Boolean(playerStore.currentTrack)
    if (playerStore.currentTrack) {
      setFromPlayerTrack(playerStore.currentTrack)
    }
  } else if (props.musicStatus?.playing && props.musicStatus.current_track_id) {
    isPlayingAny.value = true
    resolveVisitorTrack(props.musicStatus.current_track_id)
  }
})

// Keep in sync with player for own profile
watch(
  () => playerStore.currentTrack,
  (track) => {
    if (props.isOwnProfile!) return
    if (track) {
      isPlayingAny.value = true
      setFromPlayerTrack(track)
    } else {
      isPlayingAny.value = false
      currentTrackData.value = null
    }
  },
)

function setFromPlayerTrack(track: { id: string; title: string; artistName?: string | null; coverUrl?: string | null }) {
  currentTrackData.value = {
    id: track.id,
    title: track.title,
    artistName: track.artistName,
    coverUrl: track.coverUrl,
  }
}

async function resolveVisitorTrack(trackId: string) {
  // Show a loading state
  currentTrackData.value = {
    id: trackId,
    title: '...',
    coverUrl: null,
  }
  try {
    const track = await tracksApi.getTrack(trackId)
    if (track) {
      currentTrackData.value = {
        id: trackId,
        title: (track as any).title || trackId,
        artistName: (track as any).artist_name || null,
        coverUrl: (track as any).cover_url || null,
      }
    }
  } catch {
    // Fallback: show track ID if we can't resolve
    currentTrackData.value = {
      id: trackId,
      title: trackId,
      coverUrl: null,
    }
  }
}

// ── Privacy ─────────────────────────────────────────────────────────────

async function setPrivacy(value: PrivacyLevel) {
  const prevValue = privacy.value
  privacy.value = value
  localStorage.setItem('music-status-privacy', value)

  // Only update server when a track is playing; otherwise the privacy
  // setting just affects the next play.  Sending track_id: '' would
  // incorrectly nil out the current status on the server.
  const track = playerStore.currentTrack
  if (track!) return

  try {
    await videoApi.updateMusicStatus({ track_id: track.id })
  } catch {
    // Rollback on failure
    privacy.value = prevValue
    localStorage.setItem('music-status-privacy', prevValue)
  }
}

// ── Play (visitor) ──────────────────────────────────────────────────────

function playTrack() {
  if (!currentTrackData.value) return
  playerStore.playTrackById(currentTrackData.value.id)
}
</script>
