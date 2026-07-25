import { wsClient } from '@/services/socket'

interface PlaylistEvent {
  playlist_id: string
  track_id?: string
  user_id: string
  position?: number
}

export function useCollaborativePlaylist(
  playlistId: string,
  onTrackAdded?: (trackId: string, userId: string) => void,
  onTrackRemoved?: (trackId: string, userId: string) => void,
  onTrackReordered?: (trackId: string, position: number, userId: string) => void,
  onPlaylistUpdated?: (userId: string) => void,
) {
  let cleanupFns: (() => void)[] = []

  function setup() {
    if (!playlistId) return

    wsClient.connect()
    wsClient.subscribe(`playlist:${playlistId}`)

    if (onTrackAdded) {
      cleanupFns.push(
        wsClient.on('playlist.track_added', (msg: any) => {
          const ev = (msg as { payload: PlaylistEvent }).payload
          if (ev.playlist_id === playlistId) {
            onTrackAdded(ev.track_id!, ev.user_id)
          }
        }),
      )
    }

    if (onTrackRemoved) {
      cleanupFns.push(
        wsClient.on('playlist.track_removed', (msg: any) => {
          const ev = (msg as { payload: PlaylistEvent }).payload
          if (ev.playlist_id === playlistId) {
            onTrackRemoved(ev.track_id!, ev.user_id)
          }
        }),
      )
    }

    if (onTrackReordered) {
      cleanupFns.push(
        wsClient.on('playlist.track_reordered', (msg: any) => {
          const ev = (msg as { payload: PlaylistEvent }).payload
          if (ev.playlist_id === playlistId) {
            onTrackReordered(ev.track_id!, ev.position!, ev.user_id)
          }
        }),
      )
    }

    if (onPlaylistUpdated) {
      cleanupFns.push(
        wsClient.on('playlist.updated', (msg: any) => {
          const ev = (msg as { payload: PlaylistEvent }).payload
          if (ev.playlist_id === playlistId) {
            onPlaylistUpdated(ev.user_id)
          }
        }),
      )
    }
  }

  function teardown() {
    if (playlistId) {
      wsClient.unsubscribe(`playlist:${playlistId}`)
    }
    cleanupFns.forEach((fn) => fn())
    cleanupFns = []
  }

  return { setup, teardown }
}
