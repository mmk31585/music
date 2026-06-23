import { ref, readonly, onMounted, onUnmounted } from 'vue'
import { wsClient } from '@/services/socket'
import { getQueueState, suggestTrack, castVote, removeVote, reportTrackEnded } from '@/services/api/social/room-queue'
import type { QueueState } from '@/services/api/social/room-queue'

export function useRoomQueueSocket(roomId: string) {
  const queueState = ref<QueueState | null>(null)
  const isConnected = ref(false)
  // Track the playback context so we only report track-ended for this room
  const playbackContext = ref<{ type: 'room'; roomId: string } | null>(null)
  const partyStatus = ref<string | null>(null)

  let unsubQueueUpdated: (() => void) | null = null
  let unsubTrackChanged: (() => void) | null = null
  let unsubPartyStatus: (() => void) | null = null

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
    // Also subscribe to party channel for party-level events (status, position, etc.)
    wsClient.subscribe(`party:${roomId}`)

    unsubQueueUpdated = wsClient.on('room.queue_updated', () => {
      fetchState()
    })

    unsubTrackChanged = wsClient.on('room.track_changed', () => {
      fetchState()
    })

    unsubPartyStatus = wsClient.on('party.status_changed', (payload: { status: string }) => {
      partyStatus.value = payload.status
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
    if (unsubPartyStatus) {
      unsubPartyStatus()
      unsubPartyStatus = null
    }
    wsClient.unsubscribe(`room:${roomId}`)
    wsClient.unsubscribe(`party:${roomId}`)
    isConnected.value = false
  }

  onMounted(() => {
    fetchState()
    setupSocket()
  })

  onUnmounted(() => {
    teardownSocket()
  })

  return {
    queueState: readonly(queueState),
    isConnected: readonly(isConnected),
    playbackContext,
    partyStatus: readonly(partyStatus),
    refresh: fetchState,
    suggest: (trackId: string) => suggestTrack(roomId, trackId),
    vote: (candidateId: string) => castVote(roomId, candidateId),
    unvote: (candidateId: string) => removeVote(roomId, candidateId),
    reportEnded: (trackId: string) => reportTrackEnded(roomId, trackId),
  }
}
