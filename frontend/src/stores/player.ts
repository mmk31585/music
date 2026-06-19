import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { PlaybackTrack } from '@/services/api/player'
import { usePlayerApi } from '@/services/api/player'
import { useLibraryApi } from '@/services/api/library'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { useRecommendationsApi } from '@/services/api/recommendation'
import {
  audioEngine,
  preloadManager,
  queueManager,
  setMediaSessionPlaybackState,
  updateMediaSession,
} from '@/services/player'
import { useAppToast } from '@/composables/useAppToast'

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

  type ShuffleMode = 'off' | 'queue' | 'catalog' | 'similar'
  const shuffleMode = ref<ShuffleMode>('off')
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

  const hasNext = computed(() => {
    if (shuffleMode.value !== 'off') return true
    return Boolean(queueManager.getNext())
  })
  const hasPrevious = computed(() => {
    if (shuffleMode.value !== 'off') return true
    return Boolean(queueManager.getPrevious())
  })

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
      if (repeatMode.value === 'one') {
        currentTime.value = 0
        audioEngine.seek(0)
        await audioEngine.play()
        return
      }
      await playNext()
    })

    audioEngine.on('error', async (err) => {
      error.value = err.message
      isBuffering.value = false
      isPlaying.value = false
      setMediaSessionPlaybackState('none')
      try {
        const toast = useAppToast()
        toast.error(err.message || 'Playback failed')
      } catch { /* ignore */ }
      const next = queueManager.getNext()
      if (next) {
        await playNext()
      }
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

      const libraryApi = useLibraryApi()
      libraryApi.addPlayHistory({ track_id: track.id, duration: track.durationSeconds })

      const nextTrack = queueManager.getNext()
      if (nextTrack) {
        preloadManager.preload(nextTrack.streamUrl)
      }
    } catch (err: any) {
      const message = err instanceof Error ? err.message : String(err)
      error.value = message || 'Could not play track'
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
      const message = err instanceof Error ? err.message : String(err)
      error.value = message || 'Could not load track'
    } finally {
      isLoadingTrack.value = false
    }
  }

  async function setQueueAndPlay(tracks: PlaybackTrack[], startIndex = 0) {
    initialize()

    if (!tracks.length) return

    queueManager.setQueue(tracks, startIndex)
    queue.value = queueManager.all()

    if (shuffleMode.value === 'queue') {
      buildShuffleOrder()
    }

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
      const message = err instanceof Error ? err.message : String(err)
      error.value = message || 'Could not resume playback'
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
    if (shuffleMode.value === 'queue') {
      const q = queueManager.all()
      // advance position; if at -1 (first click after build) this becomes 0
      shufflePosition.value++

      // skip the track that's already playing if it happens to be next in order
      const currentId = currentTrack.value?.id
      while (shufflePosition.value < shuffleOrder.value.length) {
        const nextRealIdx = shuffleOrder.value[shufflePosition.value]
        const candidate = q[nextRealIdx]
        if (candidate && candidate.id !== currentId) break
        shufflePosition.value++
      }

      if (shufflePosition.value < shuffleOrder.value.length) {
        const nextRealIdx = shuffleOrder.value[shufflePosition.value]
        const nextTrack = q[nextRealIdx]
        if (nextTrack) {
          queueManager.setCurrent(nextTrack)
          await playTrack(nextTrack)
          return
        }
      }

      // exhausted shuffled order
      if (repeatMode.value === 'all') {
        buildShuffleOrder()
        if (shuffleOrder.value.length > 0) {
          shufflePosition.value = 0
          const nextRealIdx = shuffleOrder.value[0]
          const nextTrack = q[nextRealIdx]
          if (nextTrack) {
            queueManager.setCurrent(nextTrack)
            await playTrack(nextTrack)
            return
          }
        }
      }

      pause()
      currentTime.value = 0
      return
    }

    if (shuffleMode.value === 'catalog') {
      const tracksApi = useTracksApi()
      try {
        const randomTracks = await tracksApi.getRandomTracks({ limit: 20 })
        if (randomTracks && randomTracks.length > 0) {
          const pick = randomTracks[Math.floor(Math.random() * randomTracks.length)]
          await playTrackById(pick.id)
          return
        }
      } catch {
        // fall through to sequential
      }
    }

    if (shuffleMode.value === 'similar' && currentTrack.value) {
      const recsApi = useRecommendationsApi()
      try {
        const similar = await recsApi.getSimilar(currentTrack.value.id, { limit: 10 })
        const tracks: any[] = (similar as any)?.items ?? []
        if (tracks.length > 0) {
          const pick = tracks[Math.floor(Math.random() * tracks.length)]
          await playTrackById(pick.id)
          return
        }
      } catch {
        // fall through to sequential
      }
    }

    // Default sequential
    const nextTrack = queueManager.next()

    if (!nextTrack) {
      if (repeatMode.value === 'all') {
        // loop back to beginning of queue
        queueManager.setQueue(queueManager.all(), 0)
        const firstTrack = queueManager.next()
        if (firstTrack) {
          await playTrack(firstTrack)
          return
        }
      }
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

    if (shuffleMode.value === 'queue') {
      shufflePosition.value--

      // if we ended up at -1 and repeat is all, wrap to end
      if (shufflePosition.value < 0 && repeatMode.value === 'all') {
        shufflePosition.value = shuffleOrder.value.length - 1
      }

      if (shufflePosition.value >= 0) {
        const q = queueManager.all()
        const prevRealIdx = shuffleOrder.value[shufflePosition.value]
        const prevTrack = q[prevRealIdx]
        if (prevTrack) {
          queueManager.setCurrent(prevTrack)
          await playTrack(prevTrack)
          return
        }
      }
      seek(0)
      return
    }

    if (shuffleMode.value === 'catalog') {
      const tracksApi = useTracksApi()
      try {
        const randomTracks = await tracksApi.getRandomTracks({ limit: 20 })
        if (randomTracks && randomTracks.length > 0) {
          const pick = randomTracks[Math.floor(Math.random() * randomTracks.length)]
          await playTrackById(pick.id)
          return
        }
      } catch { /* fall through */ }
      seek(0)
      return
    }

    if (shuffleMode.value === 'similar' && currentTrack.value) {
      const recsApi = useRecommendationsApi()
      try {
        const similar = await recsApi.getSimilar(currentTrack.value.id, { limit: 10 })
        const tracks: any[] = (similar as any)?.items ?? []
        if (tracks.length > 0) {
          const pick = tracks[Math.floor(Math.random() * tracks.length)]
          await playTrackById(pick.id)
          return
        }
      } catch { /* fall through */ }
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

  const shuffleOrder = ref<number[]>([])
  const shufflePosition = ref(-1)

  function buildShuffleOrder() {
    const q = queueManager.all()
    const indices = q.map((_, i) => i)
    // Fisher-Yates
    for (let i = indices.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1))
      ;[indices[i], indices[j]] = [indices[j], indices[i]]
    }
    shuffleOrder.value = indices
    // Start before the first track — next click advances to shuffleOrder[0]
    shufflePosition.value = -1
  }

  function setShuffleMode(mode: ShuffleMode) {
    shuffleMode.value = mode
    if (mode === 'queue') {
      buildShuffleOrder()
    }
  }

  function toggleShuffle() {
    const modes: ShuffleMode[] = ['off', 'queue', 'catalog', 'similar']
    const idx = modes.indexOf(shuffleMode.value)
    setShuffleMode(modes[(idx + 1) % modes.length])
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
    if (shuffleMode.value === 'queue') {
      buildShuffleOrder()
    }
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
    setShuffleMode,
    toggleShuffle,
    toggleRepeat,
    setPlaybackRate,
    updateQueue,
    setSleepTimer,
    clearSleepTimer,
  }
})
