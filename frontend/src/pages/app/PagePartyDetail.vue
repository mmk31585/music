<template>
  <div class="mx-auto max-w-4xl space-y-6 px-4 pt-20 pb-24 md:px-8">
    <button
      class="inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="goBack"
    >
      &larr; Back to Social
    </button>

    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" class="h-48" />
    </div>

    <div v-else-if="error" class="rounded-2xl bg-white/[0.03] p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="party">
      <!-- Party header -->
      <div class="glass-strong rounded-2xl p-6">
        <div class="flex items-start gap-4">
          <div class="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-[#1db954]/10 text-2xl">
            🎉
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-3">
              <h1 class="truncate text-2xl font-black text-white">{{ party.title }}</h1>
              <span
                class="rounded-full px-3 py-1 text-xs font-semibold"
                :class="statusClass"
              >
                {{ party.status }}
              </span>
            </div>
            <p v-if="party.description" class="mt-1.5 line-clamp-2 text-sm text-white/40">
              {{ party.description }}
            </p>
            <div class="mt-3 flex items-center gap-4 text-xs text-white/30">
              <span>{{ party.participant_count }} participants</span>
              <span>Hosted by {{ userName(party.host_id) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Host Controls -->
      <div v-if="party.host_id === auth.user?.id" class="glass-strong rounded-2xl p-6">
        <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">Host Controls</h2>
        <div class="flex flex-wrap gap-3">
          <button
            class="rounded-xl bg-green-500/10 px-4 py-2 text-sm font-semibold text-green-400 transition hover:bg-green-500/20"
            @click="updateStatus('active')"
          >Play</button>
          <button
            class="rounded-xl bg-yellow-500/10 px-4 py-2 text-sm font-semibold text-yellow-400 transition hover:bg-yellow-500/20"
            @click="updateStatus('paused')"
          >Pause</button>
          <button
            class="rounded-xl bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-400 transition hover:bg-red-500/20"
            @click="updateStatus('ended')"
          >End</button>
        </div>
      </div>

      <!-- Now Playing Hero -->
      <RoomNowPlayingHero
        :now-playing="queueState?.now_playing ?? null"
        :is-loading="!queueState"
        :is-playing="isPlayingTrack"
        @toggle-play="handleTogglePlay"
      />

      <!-- Candidate Queue -->
      <RoomQueueList
        :candidates="queueState?.candidates ?? []"
        @vote="handleVote"
        @unvote="handleUnvote"
        @suggest-clicked="showTrackPicker = true"
      />

      <!-- Join/Leave -->
      <button
        class="w-full rounded-xl py-3 text-sm font-bold transition"
        :class="isParticipant
          ? 'border border-red-500/20 text-red-400 hover:bg-red-500/10'
          : 'bg-[#1db954] text-black hover:bg-[#1db954]/90'"
        @click="isParticipant ? handleLeave() : handleJoin()"
      >
        {{ isParticipant ? 'Leave Party' : 'Join Party' }}
      </button>

      <!-- Track picker -->
      <TrackPickerDialog
        :visible="showTrackPicker"
        title="Choose a track for the party"
        @update:visible="showTrackPicker = false"
        @select="handleSuggest"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
/*
 * MANUAL TEST — RUN WITH TWO BROWSER WINDOWS:
 * 1. Start backend (go run cmd/api/main.go) + frontend (npm run dev)
 * 2. Log in as User A in window 1, User B in window 2
 * 3. User A creates a Listening Party
 * 4. User B joins the same party
 * 5. Verify: both windows show the same now playing (or "waiting" state)
 * 6. User A clicks "+ پیشنهاد آهنگ", picks a track
 * 7. Verify: track appears in User B's queue list in real time
 * 8. User B votes on it (▲)
 * 9. Verify: vote count updates in both windows instantly
 * 10. Wait for current track to finish (or manually trigger track-ended)
 * 11. Verify: the voted track becomes Now Playing in both windows
 * 12. Verify: it's removed from the candidate queue
 *
 * FLIP transition: we use TransitionGroup with `queue-promote` name on
 * the candidate list. When a candidate wins and becomes now_playing,
 * the track's card is removed from the list (candidates shrinks) and
 * the hero shows the new track. The fade-out/fade-in approach is used
 * for reliability — the candidate fades out (200ms), then the hero
 * fades in (200ms, 100ms delay). prefers-reduced-motion disables all.
 */
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerApi } from '@/services/api/player'
import { usePlayerStore } from '@/stores/player'
import { useUserApi } from '@/services/api/users'
import { useUserAuthStore } from '@/stores'
import { TrackPickerDialog, RoomNowPlayingHero, RoomQueueList } from '@/components/social'
import { useRoomQueueSocket } from '@/composables/social/useRoomQueueSocket'
import type { ListeningParty } from '@/services/api/social'
import type { Track } from '@/services/api/catalog/tracks'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const playerApi = usePlayerApi()
const playerStore = usePlayerStore()
const auth = useUserAuthStore()

const partyId = route.params.id as string

const loading = ref(true)
const error = ref('')
const party = ref<ListeningParty | null>(null)
const currentTrack = ref<PlaybackTrack | null>(null)
const userNames = ref<Record<string, string>>({})

const isParticipant = ref(false)
const showTrackPicker = ref(false)

// Room queue socket composable
const roomQueue = useRoomQueueSocket(partyId)
const { queueState, playbackContext } = roomQueue

/* ---- Track from player store for play/pause ---- */
const isPlayingTrack = computed(() => {
  return playerStore.isPlaying && playerStore.currentTrack?.id === currentTrack.value?.id
})

/* ---- Re-import PlaybackTrack type inline ---- */
interface PlaybackTrack {
  id: string
  title: string
  artistName: string
  albumTitle?: string | null
  coverUrl?: string | null
  durationSeconds?: number | null
  streamUrl: string
}
/* ---- end PlaybackTrack type ---- */

function handleTogglePlay() {
  if (isPlayingTrack.value) {
    playerStore.pause()
  } else if (currentTrack.value && party.value?.current_track_id) {
    if (playerStore.currentTrack?.id === party.value.current_track_id) {
      playerStore.resume()
    } else {
      playerStore.playTrackById(party.value.current_track_id)
    }
  }
}

/* ---- Track suggestion & voting ---- */
function handleSuggest(track: Track) {
  showTrackPicker.value = false
  const trackId = String(track.id)
  roomQueue.suggest(trackId)
  // Also reflect the track as now playing if host selected it
  if (party.value && party.value.host_id === auth.user?.id) {
    updateStatus('active', trackId)
    party.value.current_track_id = trackId
    currentTrack.value = {
      id: trackId,
      title: track.title,
      artistName: track.artist_name || 'Unknown',
      albumTitle: track.album_title || null,
      coverUrl: track.cover_url || null,
      durationSeconds: track.duration_seconds ?? null,
      streamUrl: playerApi.getTrackStreamUrl(trackId),
    }
  }
}

async function handleVote(candidateId: string) {
  try {
    await roomQueue.vote(candidateId)
  } catch { /* ignore */ }
}

async function handleUnvote(candidateId: string) {
  try {
    await roomQueue.unvote(candidateId)
  } catch { /* ignore */ }
}

/* ---- Existing party logic ---- */

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  try {
    const profile = await useUserApi().getPublicUserProfile(userId) as Record<string, any>
    const p = profile
    userNames.value[userId] = String(p.full_name || p.username || userId.slice(0, 8))
  } catch {
    userNames.value[userId] = userId.slice(0, 8)
  }
}

function userName(userId: string): string {
  return userNames.value[userId] || userId?.slice(0, 8) || ''
}

const statusClass = computed(() => {
  if (!party.value) return ''
  switch (party.value.status) {
    case 'active': return 'bg-green-500/10 text-green-400'
    case 'paused': return 'bg-yellow-500/10 text-yellow-400'
    case 'ended': return 'bg-white/10 text-white/40'
    default: return 'bg-white/10 text-white/40'
  }
})

async function loadParty() {
  try {
    const partyData = await api.getParty(partyId)
    party.value = partyData
    await fetchUserName(partyData.host_id)

    if (partyData.current_track_id) {
      try {
        currentTrack.value = await playerApi.getPlaybackTrack(partyData.current_track_id)
      } catch { /* ignore */ }
    }

    isParticipant.value = false
  } catch (e: any) {
    error.value = e instanceof Error ? e.message : 'Failed to load party'
  } finally {
    loading.value = false
  }
}

async function updateStatus(status: string, trackId?: string) {
  try {
    await api.updatePartyStatus(partyId, status, trackId)
    if (party.value) party.value.status = status as ListeningParty['status']
  } catch { /* ignore */ }
}

async function handleJoin() {
  try {
    await api.joinParty(partyId)
    isParticipant.value = true
  } catch { /* ignore */ }
}

async function handleLeave() {
  try {
    await api.leaveParty(partyId)
    isParticipant.value = false
  } catch { /* ignore */ }
}

function goBack() {
  router.push({ name: 'social' })
}

/* ---- Track-ended detection ---- */
// Watch for audio engine "ended" events via the player store.
// When the currentParty track finishes and this room is our
// playback context, report it to the backend so the next
// candidate advances.
let unsubscribeEnded: (() => void) | null = null

onMounted(() => {
  loadParty()
  // Set playback context when this page mounts
  playbackContext.value = { type: 'room', roomId: partyId }

  // Detect track ended via player store's currentTime + duration watcher
  const stopWatch = watch(
    () => [playerStore.currentTime, playerStore.duration, playerStore.isPlaying] as const,
    ([time, dur, playing]) => {
      if (!playing && dur > 0 && time >= dur - 1 && currentTrack.value) {
        // Track ended naturally
        const trackedId = currentTrack.value.id
        roomQueue.reportEnded(trackedId)
      }
    },
  )
  unsubscribeEnded = stopWatch
})

onUnmounted(() => {
  playbackContext.value = null
  if (unsubscribeEnded) {
    unsubscribeEnded()
    unsubscribeEnded = null
  }
})
</script>
