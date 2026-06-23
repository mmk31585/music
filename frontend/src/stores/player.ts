import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { PlaybackTrack } from '@/services/api/player/types'
import { usePlayerApi } from '@/services/api/player/routes'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useLibraryApi } from '@/services/api/library'
import { useVideoApi } from '@/services/api/video'
import { useUserAuthStore } from '@/stores/user-auth'
import {
  PlayerEngine,
  type ShuffleMode,
  type RepeatMode,
} from '@/services/player'

const QUEUE_STORAGE_KEY = 'player-queue-track-ids'

interface QueuePersistData {
  trackIds: string[]
  hydratedAt: number
}

export const usePlayerStore = defineStore('player', () => {
  const playerApi = usePlayerApi()
  const tracksApi = useTracksApi()
  const recsApi = useRecommendationsApi()
  const libraryApi = useLibraryApi()

const currentTrack = ref<PlaybackTrack | null>(null)
const queue = ref<PlaybackTrack[]>([])

const isPlaying = ref(false)
const isBuffering = ref(false)
const isLoadingTrack = ref(false)

const currentTime = ref(0)
const duration = ref(0)

const volume = ref(0.85)
const muted = ref(false)

const shuffleMode = ref<ShuffleMode>('off')
const repeatMode = ref<RepeatMode>('off')
  const playbackRate = ref(1)
  const sleepTimerMinutes = ref(0)
  const crossfadeDuration = ref(0)
  type AudioQuality = 'auto' | 'low' | 'medium' | 'high' | 'lossless'
  const audioQuality = ref<AudioQuality>('auto')
  let sleepTimerId: ReturnType<typeof setTimeout> | null = null

  const error = ref<string | null>(null)

  // Consecutive-failure guard for the player error handler.
  // Defined at store top level so both initialize() (which sets up the
  // engine.on('error') listener) and resetFailureGuard() (called from
  // playTrack / toggleTrack / etc.) share the same closure scope.
  let consecutiveFailures = 0
  const maxConsecutiveFailures = 3
  let playbackStopped = false

  const progressPercent = computed(() => {
    if (!duration.value) return 0
    return Math.min(100, Math.max(0, (currentTime.value / duration.value) * 100))
  })

  const _currentTrackIndex = computed(() => {
    if (!currentTrack.value) return -1
    return queue.value.findIndex((t) => t.id === currentTrack.value?.id)
  })

  const hasNext = computed(() => {
    if (shuffleMode.value !== 'off') return true
    if (_currentTrackIndex.value < 0) return false
    return _currentTrackIndex.value + 1 < queue.value.length
  })

  const hasPrevious = computed(() => {
    if (shuffleMode.value !== 'off') return true
    return currentTime.value > 0
  })

  let engine: PlayerEngine | null = null
  let initialized = false
  const unsubs: (() => void)[] = []

  // ── Music Status (Now Playing) ─────────────────────────────────────
  const videoApiRef = useVideoApi()
  let musicStatusTimer: ReturnType<typeof setTimeout> | null = null
  let pendingStatusTrackId: string | null = null
  let _beforeUnloadHandler: (() => void) | null = null

  /**
   * Register the beforeunload handler — must be called from a component's
   * onMounted (or from a composable used in a component) so it's properly
   * tied to the component lifecycle.
   */
  function registerBeforeUnload() {
    if (_beforeUnloadHandler || typeof window === 'undefined') return
    _beforeUnloadHandler = () => {
      if (musicStatusTimer && pendingStatusTrackId) {
        clearTimeout(musicStatusTimer)
        try {
          const auth = useUserAuthStore()
          if (auth.isAuthenticated && !auth.isGuest) {
            const privacy = localStorage.getItem('music-status-privacy') || 'public'
            if (privacy !== 'private') {
              navigator.sendBeacon(
                '/api/v1/users/me/music-status',
                JSON.stringify({ track_id: pendingStatusTrackId }),
              )
            }
          }
        } catch {
          // Best-effort flush
        }
      }
    }
    window.addEventListener('beforeunload', _beforeUnloadHandler)
  }

  function unregisterBeforeUnload() {
    if (_beforeUnloadHandler && typeof window !== 'undefined') {
      window.removeEventListener('beforeunload', _beforeUnloadHandler)
      _beforeUnloadHandler = null
    }
  }

  async function updateMusicStatus(trackId?: string) {
    // Debounce: only send at most 1 update per 3 seconds
    if (musicStatusTimer) {
      clearTimeout(musicStatusTimer)
    }

    pendingStatusTrackId = trackId || null

    musicStatusTimer = setTimeout(async () => {
      musicStatusTimer = null
      pendingStatusTrackId = null

      try {
        const auth = useUserAuthStore()
        if (!auth.isAuthenticated || auth.isGuest) return

        // Check cached privacy setting
        const privacy = localStorage.getItem('music-status-privacy') || 'public'
        if (privacy === 'private') return

        if (trackId) {
          await videoApiRef.updateMusicStatus({ track_id: trackId })
        } else {
          await videoApiRef.clearMusicStatus()
        }
      } catch {
        // Silent fail — music status is best-effort
      }
    }, 3000)
  }

  /** Helper: map raw API track data to PlaybackTrack shape */
  function mapToPlaybackTrack(t: {
    id: string | number
    title?: string
    artist_name?: string | null
    album_title?: string | null
    cover_url?: string | null
    duration_seconds?: number | null
    audio_url?: string | null
  }): PlaybackTrack {
    return {
      id: String(t.id),
      title: t.title ?? 'Unknown Track',
      artistName: t.artist_name ?? 'Unknown artist',
      albumTitle: t.album_title ?? null,
      coverUrl: t.cover_url ?? null,
      durationSeconds: t.duration_seconds ?? null,
      streamUrl: t.audio_url ?? '',
    }
  }

  function mapItemToPlaybackTrack(item: {
    id: string
    title?: string
    artist_name?: string | null
    album_title?: string | null
    cover_url?: string | null
    duration_seconds?: number | null
    audio_url?: string | null
  }): PlaybackTrack {
    return mapToPlaybackTrack(item)
  }

  function initialize() {
    if (initialized) return
    initialized = true

    engine = new PlayerEngine(undefined, undefined, undefined, {
      fetchTrack: (id: string) => playerApi.getPlaybackTrack(id),
      fetchRandomTracks: (limit: number) =>
        tracksApi.getRandomTracks({ limit }).then((tracks) =>
          tracks.map(mapToPlaybackTrack),
        ),
      fetchSimilarTracks: (trackId: string, limit: number) =>
        recsApi.getSimilar(trackId, { limit }).then((res) =>
          (res.items || []).map(mapItemToPlaybackTrack),
        ),
      onPlayHistory: (trackId: string, dur: number) => {
        libraryApi.addPlayHistory({ track_id: trackId, duration: dur })
      },
    })

    engine.setVolume(volume.value)
    if (muted.value) engine.toggleMute()

    // Restore persisted queue (async, best-effort)
    restorePersistedQueue().then((restored) => {
      if (restored.length > 0) {
        queue.value = restored
        engine?.updateQueue(restored)
      }
    })

    unsubs.push(
      engine.on('trackchange', (track) => {
        currentTrack.value = track ? { ...track } : null
        updateMusicStatus(track?.id)
      }),
    )

    unsubs.push(
      engine.on('playstate', (state) => {
        isPlaying.value = state === 'playing'
      }),
    )

    unsubs.push(
      engine.on('timeupdate', (payload) => {
        currentTime.value = payload.currentTime
        duration.value = payload.duration
      }),
    )

    unsubs.push(
      engine.on('buffering', (buffering) => {
        isBuffering.value = buffering
      }),
    )

    unsubs.push(
      engine.on('queuechange', (q) => {
        queue.value = q
        persistQueue(q)
      }),
    )

    unsubs.push(
      engine.on('error', (msg) => {
        // Ignore errors after we've already stopped — stale audio-element
        // events (404 responses already in-flight) can keep firing.
        if (playbackStopped) return

        error.value = msg
        isBuffering.value = false
        isPlaying.value = false

        // Increment failure counter BEFORE any action so recursive
        // errors from next() don't bypass the guard.
        consecutiveFailures++

        if (consecutiveFailures >= maxConsecutiveFailures) {
          playbackStopped = true
          if (import.meta.env.DEV) {
            console.warn('[Player] Too many consecutive failures — stopping playback')
          }
          error.value = 'Playback unavailable — the media server may be offline'
          engine?.stop()
          return
        }

        // Skip to the next track in queue (don't retry — if the backend
        // is down every retry will fail and the queue advances anyway)
        if (import.meta.env.DEV) {
          console.warn(`[Player] Skipping track due to error: ${currentTrack.value?.title ?? msg}`)
        }
        engine?.next().catch(() => {})
      }),
    )

    // Volume/mute: listen for audio engine changes and sync to store + localStorage
    unsubs.push(
      engine.on('volumechange', (payload) => {
        volume.value = payload.volume
        muted.value = payload.muted
        localStorage.setItem('player-volume', String(payload.volume))
        localStorage.setItem('player-muted', String(payload.muted))
      }),
    )
  }

  /** Extract error message from unknown error */
  function getErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof Error) return err.message || fallback
    if (typeof err === 'string') return err
    if (err && typeof err === 'object' && 'message' in err) return String((err as { message: unknown }).message) || fallback
    return fallback
  }

  // ── Queue persistence (track IDs only) ──────────────────────────────
  function persistQueue(tracks: PlaybackTrack[]) {
    try {
      const data: QueuePersistData = {
        trackIds: tracks.map((t) => t.id),
        hydratedAt: Date.now(),
      }
      localStorage.setItem(QUEUE_STORAGE_KEY, JSON.stringify(data))
    } catch {
      // Storage full or blocked — silently skip
    }
  }

  /** Try to restore a previously persisted queue by fetching fresh metadata. */
  async function restorePersistedQueue(): Promise<PlaybackTrack[]> {
    try {
      const raw = localStorage.getItem(QUEUE_STORAGE_KEY)
      if (!raw) return []
      const data = JSON.parse(raw) as QueuePersistData
      if (!Array.isArray(data.trackIds) || data.trackIds.length === 0) return []

      // Fetch fresh track data (handle missing tracks gracefully)
      const results = await Promise.allSettled(
        data.trackIds.map((id) => playerApi.getPlaybackTrack(id)),
      )

      const restored: PlaybackTrack[] = []
      for (const result of results) {
        if (result.status === 'fulfilled' && result.value) {
          restored.push(result.value)
        }
        // Silently skip tracks that no longer exist
      }

      return restored
    } catch {
      return []
    }
  }

  /** Reset the consecutive-failure guard so new playback can start fresh. */
  function resetFailureGuard() {
    consecutiveFailures = 0
    playbackStopped = false
  }

  async function playTrack(track: PlaybackTrack) {
    initialize()
    if (!engine) return

    resetFailureGuard()
    error.value = null
    isLoadingTrack.value = true

    try {
      await engine.play(track)
    } catch (err: unknown) {
      error.value = getErrorMessage(err, 'Could not play track')
    } finally {
      isLoadingTrack.value = false
    }
  }

  async function playTrackById(id: string) {
    initialize()
    if (!engine) return

    resetFailureGuard()
    isLoadingTrack.value = true
    error.value = null

    try {
      await engine.playById(id)
    } catch (err: unknown) {
      error.value = getErrorMessage(err, 'Could not load track')
    } finally {
      isLoadingTrack.value = false
    }
  }

  async function setQueueAndPlay(tracks: PlaybackTrack[], startIndex = 0) {
    initialize()
    if (!engine) return
    if (!tracks.length) return

    resetFailureGuard()
    shuffleMode.value = engine.shuffleModeValue

    await engine.setQueueAndPlay(tracks, startIndex)
  }

  async function toggleTrack(track: PlaybackTrack) {
    initialize()
    if (!engine) return

    if (currentTrack.value?.id === track.id) {
      if (isPlaying.value) {
        engine.pause()
      } else {
        await engine.resume()
      }
      return
    }

    resetFailureGuard()
    await engine.play(track)
  }

  async function resume() {
    if (!engine) return

    error.value = null

    try {
      await engine.resume()
    } catch (err: unknown) {
      error.value = getErrorMessage(err, 'Could not resume playback')
    }
  }

  function pause() {
    engine?.pause()
  }

  function stop() {
    engine?.stop()
    currentTime.value = 0
  }

  function seek(seconds: number) {
    engine?.seek(seconds)
    currentTime.value = seconds
  }

  function seekPercent(percent: number) {
    if (!duration.value) return
    const safePercent = Math.max(0, Math.min(100, percent))
    seek((safePercent / 100) * duration.value)
  }

  function setVolume(value: number) {
    volume.value = value
    muted.value = false
    localStorage.setItem('player-volume', String(value))
    localStorage.setItem('player-muted', 'false')
    engine?.setVolume(value)
  }

  function toggleMute() {
    muted.value = !muted.value
    localStorage.setItem('player-muted', String(muted.value))
    engine?.toggleMute()
  }

  async function playNext() {
    initialize()
    if (!engine) return

    shuffleMode.value = engine.shuffleModeValue
    repeatMode.value = engine.repeatModeValue

    await engine.next()
  }

  async function playPrevious() {
    initialize()
    if (!engine) return

    shuffleMode.value = engine.shuffleModeValue
    repeatMode.value = engine.repeatModeValue

    await engine.previous()
  }

  function setShuffleMode(mode: ShuffleMode) {
    engine?.setShuffleMode(mode)
    shuffleMode.value = mode
  }

  function toggleShuffle() {
    engine?.toggleShuffle()
    if (engine) shuffleMode.value = engine.shuffleModeValue
  }

  function toggleRepeat() {
    engine?.toggleRepeat()
    if (engine) repeatMode.value = engine.repeatModeValue
  }

  function setPlaybackRate(rate: number) {
    playbackRate.value = rate
    engine?.setPlaybackRate(rate)
  }

  function updateQueue(newQueue: PlaybackTrack[]) {
    queue.value = newQueue
    engine?.updateQueue(newQueue)
    if (engine) shuffleMode.value = engine.shuffleModeValue
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
    registerBeforeUnload,
    unregisterBeforeUnload,
    restorePersistedQueue,
  }
})
