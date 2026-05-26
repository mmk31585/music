import { ref } from 'vue'
import { useCatalogApi } from '@/services/api/catalog'
import type { Track, TrackFormPayload } from '@/services/api/catalog'

export function useAdminTracks() {
  const api = useCatalogApi()

  const tracks = ref<Track[]>([])
  const loading = ref(false)

  async function fetchTracks() {
    loading.value = true
    try {
      const response = await api.adminGetTracks()
      tracks.value = response
    } finally {
      loading.value = false
    }
  }

  async function createTrack(payload: TrackFormPayload) {
    const created = payload.audioFile
      ? await api.adminUploadTrackWithAudio(payload, payload.audioFile)
      : await api.adminCreateTrack(payload)

    // Optimistic update
    tracks.value = [created, ...tracks.value]
    return created
  }

  async function updateTrack(id: string | number, payload: TrackFormPayload) {
    const updated = await api.adminUpdateTrack(id, payload)
    tracks.value = tracks.value.map((t) => (t.id === id ? updated : t))
    return updated
  }

  async function deleteTrack(id: string | number) {
    await api.adminDeleteTrack(id)
    tracks.value = tracks.value.filter((t) => t.id !== id)
  }

  return {
    tracks,
    loading,
    fetchTracks,
    createTrack,
    updateTrack,
    deleteTrack,
  }
}
