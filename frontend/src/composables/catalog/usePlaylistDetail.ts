import { ref } from 'vue'
import {
  usePlaylistsApi,
  type PlaylistDetail,
  type PlaylistTrackItem,
} from '@/services/api/playlist'
import { useToast } from 'primevue/usetoast'

export function usePlaylistDetail(id: string) {
  const playlistsApi = usePlaylistsApi()
  const toast = useToast()

  const playlist = ref<PlaylistDetail | null>(null)
  const tracks = ref<PlaylistTrackItem[]>([])
  const loading = ref(false)
  const error = ref<unknown>(null)

  async function fetchPlaylist() {
    loading.value = true
    error.value = null

    try {
      const data = await playlistsApi.getPlaylist(id)
      playlist.value = data
      tracks.value = data.tracks ?? []
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  async function removeTrack(trackId: string) {
    try {
      await playlistsApi.removeTrack(id, trackId)
      tracks.value = tracks.value.filter((t) => t.track_id !== trackId)

      toast.add({
        severity: 'success',
        summary: 'Track removed',
        detail: 'Track removed from playlist',
        life: 2000,
      })
    } catch {
      toast.add({
        severity: 'error',
        summary: 'Failed to remove track',
        life: 3000,
      })
    }
  }

  async function deletePlaylist() {
    try {
      await playlistsApi.deletePlaylist(id)
      return true
    } catch {
      toast.add({
        severity: 'error',
        summary: 'Failed to delete playlist',
        life: 3000,
      })
      return false
    }
  }

  return {
    playlist,
    tracks,
    loading,
    error,
    fetchPlaylist,
    removeTrack,
    deletePlaylist,
  }
}
