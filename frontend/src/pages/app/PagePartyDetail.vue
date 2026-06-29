<template>
  <div class="mx-auto max-w-4xl space-y-6 px-4 pt-20 pb-24 md:px-8">
    <button
      class="inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="goBack"
      aria-label="Back to Social"
    >
      &larr; Back to Social
    </button>

    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" class="h-48" />
    </div>

    <div v-else-if="error" class="rounded-2xl bg-white/3 p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="party">
      <!-- Party header -->
      <div class="glass-strong rounded-2xl p-6">
        <div class="flex items-start gap-4">
          <div class="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-spotify/10 text-2xl">
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

      <!-- Now Playing Hero -->
      <RoomNowPlayingHero
        :now-playing="queueState?.now_playing ?? null"
        :is-loading="(!queueState && !party?.current_track_id) || isTrackTransitioning"
        :is-playing="isPlayingTrack"
        @toggle-play="handleTogglePlay"
      />

      <!-- Playback error banner -->
      <div
        v-if="playerStore.error"
        class="rounded-xl bg-red-500/10 px-4 py-3 ring-1 ring-red-500/20"
      >
        <div class="flex items-center justify-between gap-3">
          <div class="flex items-center gap-2 text-sm text-red-400">
            <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
            <span>{{ playerStore.error }}</span>
          </div>
          <button
            class="inline-flex items-center gap-1 rounded-lg bg-red-500/20 px-3 py-1.5 text-xs font-semibold text-red-300 transition hover:bg-red-500/30 active:scale-95"
            @click="retryPlayback"
          >
            <i aria-hidden="true" class="pi pi-refresh text-xs" />
            Retry
          </button>
        </div>
      </div>

      <!-- Host: Skip + Manual controls -->
      <div v-if="isHost" class="flex flex-wrap gap-3">
        <button
          class="inline-flex items-center gap-1.5 rounded-xl bg-white/5 px-4 py-2 text-sm font-semibold text-white/60 transition hover:bg-white/10 hover:text-white"
          :disabled="!nowPlayingTrackId || isTrackTransitioning"
          @click="skipTrack"
        >
          <i aria-hidden="true" class="pi pi-forward text-xs" />
          Skip
        </button>
        <button
          class="rounded-xl bg-green-500/10 px-4 py-2 text-sm font-semibold text-green-400 transition hover:bg-green-500/20"
          @click="updateStatus('active')"
        >
          <i aria-hidden="true" class="pi pi-play text-xs" />
          Resume Playback
        </button>
        <button
          class="rounded-xl bg-yellow-500/10 px-4 py-2 text-sm font-semibold text-yellow-400 transition hover:bg-yellow-500/20"
          @click="updateStatus('paused')"
        >
          <i aria-hidden="true" class="pi pi-pause text-xs" />
          Pause Playback
        </button>
        <button
          class="rounded-xl bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-400 transition hover:bg-red-500/20"
          @click="updateStatus('ended')"
        >
          <i aria-hidden="true" class="pi pi-stop text-xs" />
          End Party
        </button>
      </div>

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
          : 'bg-spotify text-black hover:bg-spotify/90'"
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
import { usePlayerStore } from '@/stores/player'
import { useUserApi } from '@/services/api/users'
import { useUserAuthStore } from '@/stores'
import { TrackPickerDialog, RoomNowPlayingHero, RoomQueueList } from '@/components/social'
import { useRoomQueueSocket } from '@/composables/social/useRoomQueueSocket'
import { useAppToast } from '@/composables/useAppToast'
import { audioEngine } from '@/services/player'
import type { ListeningParty } from '@/services/api/social'
import type { Track } from '@/services/api/catalog/tracks'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const playerStore = usePlayerStore()
const auth = useUserAuthStore()
const toast = useAppToast()

const partyId = String(route.params.id)

const loading = ref(true)
const error = ref('')
const party = ref<ListeningParty | null>(null)
const userNames = ref<Record<string, string>>({})

const isParticipant = ref(false)
const showTrackPicker = ref(false)
const isTrackTransitioning = ref(false)

// Room queue socket composable
const roomQueue = useRoomQueueSocket(partyId)
const { queueState, playbackContext, partyStatus } = roomQueue

/* ---- Now Playing track driven by queue state ---- */
const nowPlayingTrackId = computed(() => queueState.value?.now_playing?.track?.id ?? null)

const isPlayingTrack = computed(() => {
  return playerStore.isPlaying && playerStore.currentTrack?.id === nowPlayingTrackId.value
})

const isHost = computed(() => party.value?.host_id === auth.user?.id)

/* ---- Play / Pause ---- */
async function handleTogglePlay() {
  const trackId = nowPlayingTrackId.value
  if (!trackId) return

  // Clear any previous error before retry
  playerStore.error = null

  if (isPlayingTrack.value) {
    playerStore.pause()
  } else if (playerStore.currentTrack?.id === trackId) {
    playerStore.resume()
  } else {
    isTrackTransitioning.value = true
    try {
      await playerStore.playTrackById(trackId)
    } finally {
      isTrackTransitioning.value = false
    }
  }
}

/** Retry playback after an error */
async function retryPlayback() {
  const trackId = nowPlayingTrackId.value
  if (!trackId) return
  playerStore.error = null
  isTrackTransitioning.value = true
  try {
    await playerStore.playTrackById(trackId)
  } finally {
    isTrackTransitioning.value = false
  }
}

/** Force-play a specific track ID (called when now_playing changes) */
async function playNowPlayingTrack(trackId: string) {
  if (playerStore.currentTrack?.id === trackId && playerStore.isPlaying) return
  isTrackTransitioning.value = true
  try {
    await playerStore.playTrackById(trackId)
  } finally {
    isTrackTransitioning.value = false
  }
}

/** Skip the current track (host only) */
async function skipTrack() {
  const trackId = nowPlayingTrackId.value
  if (!trackId || !isHost.value) return
  try {
    await roomQueue.reportEnded(trackId)
  } catch {
    toast.error('Failed to skip track')
  }
}

/* ---- Auto-play when queue advances or initial state loads ---- */
watch(
  () => queueState.value?.now_playing?.track?.id,
  (newTrackId, oldTrackId) => {
    // Only auto-play if we're still a participant in the party
    if (newTrackId && newTrackId !== oldTrackId && isParticipant.value) {
      playNowPlayingTrack(newTrackId)
    }
  },
)

/* ---- Clear transitioning state when player settles after any load attempt ---- */
watch(
  () => playerStore.isLoadingTrack,
  (loading) => {
    if (!loading) {
      setTimeout(() => {
        isTrackTransitioning.value = false
      }, 0)
    }
  },
)

/* ---- Restore playback when party loads and has a current_track_id ---- */
watch(
  [() => queueState.value, () => party.value],
  ([qs, p]) => {
    if (qs && p?.current_track_id && !qs.now_playing) {
      // Party has a track but queue system hasn't picked it up yet.
      // Suggest it to the queue so it becomes now_playing on next advance.
      roomQueue.suggest(p.current_track_id)
    }
  },
  { once: true },
)

/* ---- Track suggestion & voting ---- */
function handleSuggest(track: Track) {
  showTrackPicker.value = false
  const trackId = String(track.id)
  roomQueue.suggest(trackId)
}

async function handleVote(candidateId: string) {
  try {
    await roomQueue.vote(candidateId)
  } catch (err) {
    console.error('Failed to vote:', err)
  }
}

async function handleUnvote(candidateId: string) {
  try {
    await roomQueue.unvote(candidateId)
  } catch (err) {
    console.error('Failed to unvote:', err)
  }
}

/* ---- Existing party logic ---- */

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  try {
    const profile = await useUserApi().getPublicUserProfile(userId) as Record<string, unknown>
    const p = profile
    userNames.value[userId] = String(p.full_name || p.username || userId.slice(0, 8))
  } catch (err) {
    console.error('Failed to fetch user name:', err)
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

    isParticipant.value = false
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load party'
  } finally {
    loading.value = false
  }
}

/* ---- When both party + queue are loaded, try to restore playback ---- */
watch(
  [() => party.value, () => queueState.value],
  ([p, qs]) => {
    if (!p || !qs) return

    // If the party has an active track and nothing is playing yet, play it.
    if (p.current_track_id && !playerStore.currentTrack && !qs.now_playing) {
      playerStore.playTrackById(p.current_track_id)
        .catch(() => {
          roomQueue.suggest(String(p.current_track_id))
        })
    }
  },
  { once: true },
)

async function updateStatus(status: string, trackId?: string) {
  try {
    await api.updatePartyStatus(partyId, status, trackId)
    if (party.value) party.value.status = status as ListeningParty['status']
    // Start playback immediately for the host (WebSocket broadcast will
    // handle other participants via the partyStatus watcher below)
    if (status === 'active' && nowPlayingTrackId.value) {
      playNowPlayingTrack(nowPlayingTrackId.value)
    } else if (status === 'paused' && isPlayingTrack.value) {
      playerStore.pause()
    }
  } catch (err) {
    console.error('Failed to update status:', err)
  }
}

/* ---- React to party status changes from WebSocket ---- */
watch(partyStatus, (status) => {
  if (!status || !isParticipant.value) return
  // Sync local party status for UI badge
  if (party.value) {
    party.value.status = status as ListeningParty['status']
  }
  if (status === 'active' && nowPlayingTrackId.value) {
    playNowPlayingTrack(nowPlayingTrackId.value)
  } else if (status === 'paused' && isPlayingTrack.value) {
    playerStore.pause()
  }
})

async function handleJoin() {
  try {
    await api.joinParty(partyId)
    isParticipant.value = true
  } catch (err) {
    console.error('Failed to join party:', err)
  }
}

async function handleLeave() {
  try {
    await api.leaveParty(partyId)
    isParticipant.value = false
    playbackContext.value = null
  } catch (err) {
    console.error('Failed to leave party:', err)
  }
}

function goBack() {
  router.push({ name: 'social' })
}

/* ---- Track-ended detection ---- */
// Listen for the native audio engine "ended" event and report it
// to the party backend so the queue advances to the next track.
let unsubscribeEnded: (() => void) | null = null

onMounted(() => {
  loadParty()
  // Set playback context when this page mounts
  playbackContext.value = { type: 'room', roomId: partyId }

  // Use the native audio "ended" event instead of a time-based watcher.
  // The time-based watcher (currentTime + duration) had false positives
  // when the user PAUSED near the end of a track — those would wrongly
  // call reportEnded and advance the room queue.
  unsubscribeEnded = audioEngine.on('ended', () => {
    if (nowPlayingTrackId.value) {
      roomQueue.reportEnded(nowPlayingTrackId.value)
    }
  })
})

onUnmounted(() => {
  playbackContext.value = null
  if (unsubscribeEnded) {
    unsubscribeEnded()
    unsubscribeEnded = null
  }
})
</script>
