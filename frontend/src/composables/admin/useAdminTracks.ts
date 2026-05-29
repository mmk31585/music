import { ref } from 'vue'
import { useCatalogApi } from '@/services/api/catalog'
import type { Track, TrackFormPayload } from '@/services/api/catalog'

export function useAdminTracks() {
  const api = useCatalogApi()

  const tracks = ref<Track[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<unknown>(null)

  async function fetchTracks() {
    loading.value = true
    error.value = null

    try {
      const response = await api.adminGetTracks()
      tracks.value = response
      return response
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createTrack(payload: TrackFormPayload) {
    saving.value = true
    error.value = null

    try {
      const created = payload.audioFile
        ? await api.adminUploadTrackWithAudio(payload, payload.audioFile)
        : await api.adminCreateTrack(payload)

      tracks.value = [created, ...tracks.value]

      return created
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function updateTrack(id: string | number, payload: TrackFormPayload) {
    saving.value = true
    error.value = null

    try {
      const updated = await api.adminUpdateTrack(id, payload)

      tracks.value = tracks.value.map((track) =>
        String(track.id) === String(id) ? updated : track,
      )

      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function deleteTrack(id: string | number) {
    deleting.value = true
    error.value = null

    try {
      await api.adminDeleteTrack(id)

      tracks.value = tracks.value.filter((track) => String(track.id) !== String(id))
    } catch (err) {
      error.value = err
      throw err
    } finally {
      deleting.value = false
    }
  }

  function setTracks(value: Track[]) {
    tracks.value = value
  }

  function resetTracks() {
    tracks.value = []
    error.value = null
  }

  return {
    tracks,
    loading,
    saving,
    deleting,
    error,

    fetchTracks,
    createTrack,
    updateTrack,
    deleteTrack,

    setTracks,
    resetTracks,
  }
}
