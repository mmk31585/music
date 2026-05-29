import type { PlaybackTrack } from '@/services/api/player'

export function updateMediaSession(
  track: PlaybackTrack,
  handlers: {
    play: () => void
    pause: () => void
    next?: () => void
    previous?: () => void
    seek?: (time: number) => void
  },
) {
  if (!('mediaSession' in navigator)) return

  navigator.mediaSession.metadata = new MediaMetadata({
    title: track.title,
    artist: track.artistName,
    album: track.albumTitle || '',
    artwork: track.coverUrl
      ? [
          {
            src: track.coverUrl,
            sizes: '512x512',
            type: 'image/jpeg',
          },
        ]
      : [],
  })

  navigator.mediaSession.setActionHandler('play', handlers.play)
  navigator.mediaSession.setActionHandler('pause', handlers.pause)

  if (handlers.next) {
    navigator.mediaSession.setActionHandler('nexttrack', handlers.next)
  }

  if (handlers.previous) {
    navigator.mediaSession.setActionHandler('previoustrack', handlers.previous)
  }

  if (handlers.seek) {
    navigator.mediaSession.setActionHandler('seekto', (details) => {
      if (typeof details.seekTime === 'number') {
        handlers.seek?.(details.seekTime)
      }
    })
  }
}

export function setMediaSessionPlaybackState(state: MediaSessionPlaybackState) {
  if (!('mediaSession' in navigator)) return

  navigator.mediaSession.playbackState = state
}
