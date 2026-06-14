import { computed, ref } from 'vue'
import { usePlayer } from './usePlayer'

export function usePlayerControls() {
  const player = usePlayer()

  const playIcon = computed(() => {
    if (player.isBuffering.value || player.isLoadingTrack.value) return 'pi pi-spin pi-spinner'
    return player.isPlaying.value ? 'pi pi-pause' : 'pi pi-play'
  })

  const volumeIcon = computed(() => {
    if (player.muted.value || player.volume.value === 0) return 'pi pi-volume-off'
    if (player.volume.value < 0.5) return 'pi pi-volume-down'
    return 'pi pi-volume-up'
  })

  async function togglePlayPause() {
    if (player.isPlaying.value) {
      player.pause()
      return
    }

    await player.resume()
  }

  const sleepTimerMinutes = ref(0)
  let sleepTimerId: ReturnType<typeof setTimeout> | null = null

  function clearSleepTimer() {
    if (sleepTimerId !== null) {
      clearTimeout(sleepTimerId)
      sleepTimerId = null
    }
    sleepTimerMinutes.value = 0
  }

  const speedLabel = computed(() => {
    const rate = player.playbackRate
    if (rate === 1) return 'Normal'
    return `${rate}x`
  })

  return {
    ...player,
    playIcon,
    volumeIcon,
    togglePlayPause,
    speedLabel,
  }
}
