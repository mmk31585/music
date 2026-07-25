import { ref } from 'vue'
import {
  usePlaylistsApi,
  type PlaylistDetail,
  type PlaylistTrackItem,
} from '@/services/api/playlist'
import { useToast } from 'primevue/usetoast'

export const pendingPlaylistRemoveData = ref<{ playlistId: string; trackId: string; playlistName: string } | null>(null)

export function usePlaylistDetail(id: string) {
  const playlistsApi = usePlaylistsApi()
  const toast = useToast()

  const playlist = ref<PlaylistDetail | null>(null)
  const tracks = ref<PlaylistTrackItem[]>([])
  const loading = ref(false)
  const error = ref<any>(null)

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

  const pendingRemoveTrackId = ref<string | null>(null)

  async function removeTrack(trackId: string) {
    pendingRemoveTrackId.value = trackId
    pendingPlaylistRemoveData.value = { playlistId: id, trackId, playlistName: playlist.value?.name || 'Playlist' }
    try {
      await playlistsApi.removeTrack(id, trackId)
      tracks.value = tracks.value.filter((t) => t.track_id !== trackId)

      toast.add({
        severity: 'success',
        summary: 'Track removed from playlist',
        group: 'playlist-undo',
        life: 6000,
      })
    } catch {
      toast.add({
        severity: 'error',
        summary: 'Failed to remove track',
        life: 3000,
      })
    }
  }

  async function undoRemoveTrack() {
    const trackId = pendingRemoveTrackId.value
    if (!trackId) return
    pendingRemoveTrackId.value = null
    pendingPlaylistRemoveData.value = null
    try {
      await playlistsApi.addTrack(id, { track_id: trackId })
      await fetchPlaylist()
      toast.add({ severity: 'success', summary: 'Track re-added', life: 3000 })
    } catch {
      toast.add({ severity: 'error', summary: 'Failed to undo', life: 3000 })
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
    undoRemoveTrack,
    pendingRemoveTrackId,
    deletePlaylist,
  }
}
