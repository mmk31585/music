import { storeToRefs } from 'pinia'
import { usePlayerStore } from '@/stores/player'

export function usePlayer() {
  const player = usePlayerStore()
  const state = storeToRefs(player)

  player.initialize()

  return {
    ...state,

    playTrack: player.playTrack,
    playTrackById: player.playTrackById,
    setQueueAndPlay: player.setQueueAndPlay,
    toggleTrack: player.toggleTrack,

    resume: player.resume,
    pause: player.pause,
    stop: player.stop,

    seek: player.seek,
    seekPercent: player.seekPercent,

    setVolume: player.setVolume,
    toggleMute: player.toggleMute,

    playNext: player.playNext,
    playPrevious: player.playPrevious,
  }
}
