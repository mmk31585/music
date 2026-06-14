import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { PlaybackTrack } from '@/services/api/player'
import { usePlayerApi } from '@/services/api/player'
import {
  audioEngine,
  preloadManager,
  queueManager,
  setMediaSessionPlaybackState,
  updateMediaSession,
} from '@/services/player'

export const usePlayerStore = defineStore('player', () => {
  const playerApi = usePlayerApi()

  const currentTrack = ref<PlaybackTrack | null>(null)
  const queue = ref<PlaybackTrack[]>([])

  const isPlaying = ref(false)
  const isBuffering = ref(false)
  const isLoadingTrack = ref(false)

  const currentTime = ref(0)
  const duration = ref(0)

  const volume = ref(Number(localStorage.getItem('player-volume') || 0.85))
  const muted = ref(localStorage.getItem('player-muted') === 'true')

  const shuffleMode = ref(false)
  type RepeatMode = 'off' | 'one' | 'all'
  const repeatMode = ref<RepeatMode>('off')
  const playbackRate = ref(1)
  const sleepTimerMinutes = ref(0)
  const crossfadeDuration = ref(0)
  type AudioQuality = 'auto' | 'low' | 'medium' | 'high' | 'lossless'
  const audioQuality = ref<AudioQuality>('auto')
  let sleepTimerId: ReturnType<typeof setTimeout> | null = null

  const error = ref<string | null>(null)

  const progressPercent = computed(() => {
    if (!duration.value) return 0
    return Math.min(100, Math.max(0, (currentTime.value / duration.value) * 100))
  })

  const hasNext = computed(() => Boolean(queueManager.getNext()))
  const hasPrevious = computed(() => Boolean(queueManager.getPrevious()))

  let initialized = false

  function initialize() {
    if (initialized) return
    initialized = true

    audioEngine.setVolume(volume.value)
    audioEngine.setMuted(muted.value)

    audioEngine.on('play', () => {
      isPlaying.value = true
      setMediaSessionPlaybackState('playing')
    })

    audioEngine.on('pause', () => {
      isPlaying.value = false
      setMediaSessionPlaybackState('paused')
    })

    audioEngine.on('waiting', () => {
      isBuffering.value = true
    })

    audioEngine.on('playing', () => {
      isBuffering.value = false
    })

    audioEngine.on('canplay', () => {
      isBuffering.value = false
    })

    audioEngine.on('timeupdate', (payload) => {
      currentTime.value = payload.currentTime
      duration.value = payload.duration || duration.value
    })

    audioEngine.on('durationchange', (payload) => {
      if (payload.duration) {
        duration.value = payload.duration
      }
    })

    audioEngine.on('volumechange', (payload) => {
      volume.value = payload.volume
      muted.value = payload.muted

      localStorage.setItem('player-volume', String(payload.volume))
      localStorage.setItem('player-muted', String(payload.muted))
    })

    audioEngine.on('ended', async () => {
      await playNext()
    })

    audioEngine.on('error', (err) => {
      error.value = err.message
      isBuffering.value = false
      isPlaying.value = false
      setMediaSessionPlaybackState('none')
    })
  }

  async function playTrack(track: PlaybackTrack) {
    initialize()

    error.value = null
    isLoadingTrack.value = true
    isBuffering.value = true

    try {
      currentTrack.value = track
      queueManager.setCurrent(track)
      queue.value = queueManager.all()

      duration.value = track.durationSeconds || 0
      currentTime.value = 0

      updateMediaSession(track, {
        play: resume,
        pause,
        next: playNext,
        previous: playPrevious,
        seek,
      })

      await audioEngine.play(track.streamUrl)

      const nextTrack = queueManager.getNext()
      if (nextTrack) {
        preloadManager.preload(nextTrack.streamUrl)
      }
    } catch (err: any) {
      error.value = err?.message || 'Could not play track'
    } finally {
      isLoadingTrack.value = false
      isBuffering.value = false
    }
  }

  async function playTrackById(id: string) {
    initialize()

    isLoadingTrack.value = true
    error.value = null

    try {
      const track = await playerApi.getPlaybackTrack(id)
      await playTrack(track)
    } catch (err: any) {
      error.value = err?.message || 'Could not load track'
    } finally {
      isLoadingTrack.value = false
    }
  }

  async function setQueueAndPlay(tracks: PlaybackTrack[], startIndex = 0) {
    initialize()

    if (!tracks.length) return

    queueManager.setQueue(tracks, startIndex)
    queue.value = queueManager.all()

    const track = tracks[startIndex]
    if (track) {
      await playTrack(track)
    }
  }

  async function toggleTrack(track: PlaybackTrack) {
    initialize()

    if (currentTrack.value?.id === track.id) {
      if (isPlaying.value) {
        pause()
      } else {
        await resume()
      }

      return
    }

    await playTrack(track)
  }

  async function resume() {
    initialize()

    if (!currentTrack.value) return

    error.value = null

    try {
      await audioEngine.play()
    } catch (err: any) {
      error.value = err?.message || 'Could not resume playback'
    }
  }

  function pause() {
    audioEngine.pause()
  }

  function stop() {
    audioEngine.stop()
    isPlaying.value = false
    currentTime.value = 0
  }

  function seek(seconds: number) {
    audioEngine.seek(seconds)
    currentTime.value = seconds
  }

  function seekPercent(percent: number) {
    if (!duration.value) return

    const safePercent = Math.max(0, Math.min(100, percent))
    seek((safePercent / 100) * duration.value)
  }

  function setVolume(value: number) {
    audioEngine.setVolume(value)
  }

  function toggleMute() {
    audioEngine.setMuted(!muted.value)
  }

  async function playNext() {
    const nextTrack = queueManager.next()

    if (!nextTrack) {
      pause()
      currentTime.value = 0
      return
    }

    await playTrack(nextTrack)
  }

  async function playPrevious() {
    if (currentTime.value > 4) {
      seek(0)
      return
    }

    const previousTrack = queueManager.previous()

    if (!previousTrack) {
      seek(0)
      return
    }

    await playTrack(previousTrack)
  }

  function toggleShuffle() {
    shuffleMode.value = !shuffleMode.value
  }

  function toggleRepeat() {
    if (repeatMode.value === 'off') repeatMode.value = 'all'
    else if (repeatMode.value === 'all') repeatMode.value = 'one'
    else repeatMode.value = 'off'
  }

  function setPlaybackRate(rate: number) {
    playbackRate.value = rate
    audioEngine.setPlaybackRate(rate)
  }

  function updateQueue(newQueue: PlaybackTrack[]) {
    queue.value = newQueue
    queueManager.replaceAll(newQueue)
  }

  function setSleepTimer(minutes: number) {
    clearSleepTimer()
    sleepTimerMinutes.value = minutes
    sleepTimerId = setTimeout(() => {
      pause()
      sleepTimerMinutes.value = 0
    }, minutes * 60 * 1000)
  }

  function clearSleepTimer() {
    if (sleepTimerId !== null) {
      clearTimeout(sleepTimerId)
      sleepTimerId = null
    }
    sleepTimerMinutes.value = 0
  }

  return {
    currentTrack,
    queue,

    isPlaying,
    isBuffering,
    isLoadingTrack,

    currentTime,
    duration,
    progressPercent,

    volume,
    muted,

    shuffleMode,
    repeatMode,
    playbackRate,
    sleepTimerMinutes,
    crossfadeDuration,
    audioQuality,

    error,

    hasNext,
    hasPrevious,

    initialize,
    playTrack,
    playTrackById,
    setQueueAndPlay,
    toggleTrack,
    resume,
    pause,
    stop,
    seek,
    seekPercent,
    setVolume,
    toggleMute,
    playNext,
    playPrevious,
    toggleShuffle,
    toggleRepeat,
    setPlaybackRate,
    updateQueue,
    setSleepTimer,
    clearSleepTimer,
  }
})
