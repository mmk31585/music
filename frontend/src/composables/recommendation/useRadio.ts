import { ref, watch } from 'vue'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { useRadioApi, type RadioSession } from '@/services/api/recommendation/radio'
import { useAppToast } from '@/composables/useAppToast'
import type { PlaybackTrack } from '@/services/api/player'

const REFILL_THRESHOLD = 3
const BATCH_SIZE = 10

export function useRadio() {
  const currentSession = ref<RadioSession | null>(null)
  const isRadioActive = ref(false)
  const isLoadingBatch = ref(false)
  const radioSeedLabel = ref('')

  const player = usePlayer()
  const playerApi = usePlayerApi()
  const radioApi = useRadioApi()
  const toast = useAppToast()

  function mapToPlaybackTrack(t: any): PlaybackTrack {
    return {
      id: t.id,
      title: t.title,
      artistName: t.artist_name || 'Unknown',
      albumTitle: t.album_title || null,
      coverUrl: t.cover_url || null,
      durationSeconds: t.duration_seconds ?? null,
      streamUrl: playerApi.getTrackStreamUrl(t.id),
    }
  }

  async function startFromTrack(trackId: string, trackName?: string) {
    try {
      isLoadingBatch.value = true
      const result = await radioApi.startRadio(trackId)
      currentSession.value = { session_id: result.session_id, seed_track_id: trackId }
      isRadioActive.value = true
      radioSeedLabel.value = trackName || ''

      const tracks: PlaybackTrack[] = (result.tracks || []).map(mapToPlaybackTrack)

      if (tracks.length > 0) {
        player.updateQueue(tracks)
      }
    } catch (err: any) {
      toast.error(err?.message || 'Failed to start radio')
    } finally {
      isLoadingBatch.value = false
    }
  }

  async function fetchNextBatch() {
    if (!currentSession.value || isLoadingBatch.value) return

    isLoadingBatch.value = true
    try {
      const result = await radioApi.getNextRadioBatch(currentSession.value.session_id, BATCH_SIZE)
      if (result.tracks && result.tracks.length > 0) {
        const newTracks: PlaybackTrack[] = result.tracks.map(mapToPlaybackTrack)
        player.updateQueue([...player.queue.value, ...newTracks])
      }
    } catch {
      // silent — keep current queue intact
    } finally {
      isLoadingBatch.value = false
    }
  }

  async function endRadio() {
    if (currentSession.value) {
      try {
        await radioApi.endRadio(currentSession.value.session_id)
      } catch {
        // silent
      }
    }
    currentSession.value = null
    isRadioActive.value = false
    radioSeedLabel.value = ''
  }

  // Auto-refill: watch the queue length and fetch more when running low
  watch(
    () => player.queue.value.length,
    (length) => {
      if (isRadioActive.value && length <= REFILL_THRESHOLD && !isLoadingBatch.value) {
        fetchNextBatch()
      }
    },
  )

  return {
    currentSession,
    isRadioActive,
    isLoadingBatch,
    radioSeedLabel,
    startFromTrack,
    fetchNextBatch,
    endRadio,
  }
}
