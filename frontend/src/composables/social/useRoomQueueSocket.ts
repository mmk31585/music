import { ref, readonly, onMounted, onUnmounted } from 'vue'
import { wsClient } from '@/services/socket'
import { getQueueState, suggestTrack, castVote, removeVote, reportTrackEnded } from '@/services/api/social/room-queue'
import type { QueueState } from '@/services/api/social/room-queue'

export function useRoomQueueSocket(roomId: string) {
  const queueState = ref<QueueState | null>(null)
  const isConnected = ref(false)
  // Track the playback context so we only report track-ended for this room
  const playbackContext = ref<{ type: 'room'; roomId: string } | null>(null)

  let unsubQueueUpdated: (() => void) | null = null
  let unsubTrackChanged: (() => void) | null = null

  async function fetchState() {
    try {
      const state = await getQueueState(roomId)
      queueState.value = state
    } catch {
      // keep previous state on error
    }
  }

  function setupSocket() {
    wsClient.connect()
    wsClient.subscribe(`room:${roomId}`)

    unsubQueueUpdated = wsClient.on('room.queue_updated', () => {
      fetchState()
    })

    unsubTrackChanged = wsClient.on('room.track_changed', () => {
      fetchState()
    })

    isConnected.value = true
  }

  function teardownSocket() {
    if (unsubQueueUpdated) {
      unsubQueueUpdated()
      unsubQueueUpdated = null
    }
    if (unsubTrackChanged) {
      unsubTrackChanged()
      unsubTrackChanged = null
    }
    wsClient.unsubscribe(`room:${roomId}`)
    isConnected.value = false
  }

  onMounted(() => {
    fetchState()
    setupSocket()
  })

  onUnmounted(() => {
    teardownSocket()
  })

  // Watch the player store for track-ended events.
  // When the currently playing audio track ends AND this room is
  // the active playback context, report it to the backend.
  // We use the Pinia `$subscribe`-like pattern via watcher on
  // the player store's currentTime/duration to detect end.
  // Simpler: expose a manual end-of-track callback and let the
  // page wire it up. The actual ended detection happens at the
  // page level (PagePartyDetail) via audioEngine 'ended' event.

  return {
    queueState: readonly(queueState),
    isConnected: readonly(isConnected),
    playbackContext,
    refresh: fetchState,
    suggest: (trackId: string) => suggestTrack(roomId, trackId),
    vote: (candidateId: string) => castVote(roomId, candidateId),
    unvote: (candidateId: string) => removeVote(roomId, candidateId),
    reportEnded: (trackId: string) => reportTrackEnded(roomId, trackId),
  }
}
