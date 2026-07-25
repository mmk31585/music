import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { PlaybackTrack } from '@/services/api/player/types'
import { buildPlaybackTrack, mapToPlaybackTracks } from '@/factories/playbackTrack'
import { usePlayerApi } from '@/services/api/player/routes'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useLibraryApi } from '@/services/api/library'
import { useVideoApi } from '@/services/api/video'
import { useHistoryApi } from '@/services/api/history'
import { useUserAuthStore } from '@/stores/user-auth'
import {
  PlayerEngine,
  type ShuffleMode,
  type RepeatMode,
} from '@/services/player'
import { audioEngine } from '@/services/player/audio-engine'
import { queueManager } from '@/services/player/queue-manager'

const QUEUE_STORAGE_KEY = 'player-queue-track-ids'
const QUEUE_META_KEY = 'player-queue-metadata'
const SESSION_STORAGE_KEY = 'player-session-state'

interface QueuePersistData {
  trackIds: string[]
  hydratedAt: number
}

/** Serialisable subset of PlaybackTrack for cache fallback */
interface CachedTrackMeta {
  id: string
  title: string
  artistName: string
  albumTitle?: string | null
  coverUrl?: string | null
  durationSeconds?: number | null
  streamUrl: string
}

interface SessionPersistData {
  currentTrackId: string | null
  currentTime: number
  updatedAt: number
  isPlaying: boolean
  shuffleMode: ShuffleMode
  repeatMode: RepeatMode
  currentIndex: number
}

export const usePlayerStore = defineStore('player', () => {
  const playerApi = usePlayerApi()
  const tracksApi = useTracksApi()
  const recsApi = useRecommendationsApi()
  const libraryApi = useLibraryApi()
  const historyApi = useHistoryApi()

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

  // Play tracking: reports plays to backend for gamification (XP, challenges, badges)
  let playSessionId = crypto.randomUUID()
  let lastReportedTrackId: string | null = null
  let playStartTime = 0

  /** Report a completed play to the backend (fire-and-forget). */
  function reportPlay(track: PlaybackTrack, durationSec: number, completed: boolean) {
    if (!track?.id) return
    // Don't double-report the same track consecutively
    if (track.id === lastReportedTrackId && completed) return
    lastReportedTrackId = track.id

    historyApi.recordPlay({
      track_id: track.id,
      duration: Math.floor(durationSec),
      completed,
      session_id: playSessionId,
      track_duration_ms: track.durationSeconds ? track.durationSeconds * 1000 : undefined,
    }).catch(() => {})
  }

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

  /**
   * The actual next track that will play, accounting for shuffle order.
   * - shuffle='off': next sequential track
   * - shuffle='queue': next index from the internal shuffleOrder
   * - shuffle='catalog'|'similar': null (unpredictable API fetch)
   *
   * Accesses store refs as explicit deps so Vue re-evaluates when
   * shuffle mode, queue, or current track changes — engine's internal
   * state is non-reactive so the computed needs reactive anchors.
   */
  const nextUpTrack = computed<PlaybackTrack | null>(() => {
    if (!engine) return null
    // Reactive deps: re-evaluate when these store values change
    void shuffleMode.value
    void queue.value.length
    void currentTrack.value?.id
    void repeatMode.value
    return engine.peekNextTrack()
  })

  /**
   * Queue items in the order they will actually be played.
   * - 'queue' shuffle mode: ordered by internal shuffleOrder (remaining after current)
   * - 'catalog'/'similar'/off: returns the raw queue
   *
   * Always returns a new array reference so downstream watchers fire correctly.
   */
  const remainingShuffledQueue = computed<PlaybackTrack[]>(() => {
    void shuffleMode.value
    void queue.value.length
    void repeatMode.value
    if (!engine) return [...queue.value]
    return engine.getRemainingShuffledQueue()
  })

  let engine: PlayerEngine | null = null
  let initialized = false
  const unsubs: (() => void)[] = []

  // ── Music Status (Now Playing) ─────────────────────────────────────
  const videoApiRef = useVideoApi()
  let musicStatusTimer: ReturnType<typeof setTimeout> | null = null
  let pendingStatusTrackId: string | null = null
  let _beforeUnloadHandler: (() => void) | null = null
  let _pageHideHandler: (() => void) | null = null

  /**
   * Register the pagehide handler — must be called from a component's
   * onMounted (or from a composable used in a component) so it's properly
   * tied to the component lifecycle.
   *
   * Uses `pagehide` instead of `beforeunload` because `beforeunload`
   * prevents the browser's back/forward cache (bfcache), degrading
   * navigation performance. The `pagehide` event fires in all the same
   * scenarios but is bfcache-friendly.
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
    window.addEventListener('pagehide', _beforeUnloadHandler)
  }

  function unregisterBeforeUnload() {
    if (_beforeUnloadHandler && typeof window !== 'undefined') {
      window.removeEventListener('pagehide', _beforeUnloadHandler)
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

  /** @deprecated Use buildPlaybackTrack from @/factories/playbackTrack instead */
  function mapItemToPlaybackTrack(item: {
    id: string | number
    title?: string
    artist_name?: string | null
    album_title?: string | null
    cover_url?: string | null
    duration_seconds?: number | null
    audio_url?: string | null
  }): PlaybackTrack {
    return buildPlaybackTrack(item)
  }

  function initialize() {
    if (initialized) return
    initialized = true

    engine = new PlayerEngine(undefined, undefined, undefined, {
      fetchTrack: (id: string) => playerApi.getPlaybackTrack(id),
      fetchRandomTracks: (limit: number) =>
        tracksApi.getRandomTracks({ limit }).then((tracks) =>
          mapToPlaybackTracks(tracks),
        ),
      fetchSimilarTracks: (trackId: string, limit: number) =>
        recsApi.getSimilar(trackId, { limit }).then((res) =>
          mapToPlaybackTracks(res.items || []),
        ),
      onPlayHistory: (trackId: string, dur: number) => {
        libraryApi.addPlayHistory({ track_id: trackId, duration: dur }).catch(() => { /* silent fail */ })
      },
    })

    engine.setVolume(volume.value)
    if (muted.value) engine.toggleMute()
    if (crossfadeDuration.value > 0) {
      engine.setCrossfadeDuration(crossfadeDuration.value)
    }

    // Restore persisted queue + session (async, best-effort)
    restorePersistedQueue().then(async (restored) => {
      if (restored.length > 0) {
        queue.value = restored
        engine?.updateQueue(restored)
      }

      // Try to auto-resume playback from saved session
      const sessionRestored = await restoreSession()

      // If session had a current track not in the queue, insert it at position 0
      // (This handles the case where the queue was saved but current track was elsewhere)
      if (!sessionRestored) {
        try {
          const raw = localStorage.getItem(SESSION_STORAGE_KEY)
          if (raw) {
            const data = JSON.parse(raw) as SessionPersistData
            if (data.currentTrackId && queue.value.length > 0) {
              // Check if current track is already in queue
              const inQueue = queue.value.some(t => t.id === data.currentTrackId)
              if (!inQueue) {
                // Try to fetch and prepend
                const track = await playerApi.getPlaybackTrack(data.currentTrackId)
                if (track) {
                  queue.value = [track, ...queue.value]
                  engine?.updateQueue(queue.value)
                  currentTrack.value = track
                  currentTime.value = data.currentTime
                }
              } else {
                // Already in queue — find its index and seek
                currentTime.value = data.currentTime
              }
            }
          }
        } catch {
          // Best-effort
        }
      }
    })

    // Persist session on page hide (covers refresh, tab switch, navigate away)
    if (typeof window !== 'undefined') {
      _pageHideHandler = () => {
        if (currentTrack.value) {
          persistSession()
          // Report current track play on page hide
          if (playStartTime > 0) {
            const playedSec = (Date.now() - playStartTime) / 1000
            reportPlay(currentTrack.value, playedSec, false)
          }
        }
      }
      window.addEventListener('pagehide', _pageHideHandler)
    }

    unsubs.push(
      engine.on('trackchange', (track) => {
        // Report the previous track as played when switching to a new one
        if (currentTrack.value && playStartTime > 0) {
          const playedSec = (Date.now() - playStartTime) / 1000
          const durationSec = currentTrack.value.durationSeconds || 0
          const completed = durationSec > 0 && playedSec >= durationSec * 0.8
          reportPlay(currentTrack.value, playedSec, completed)
        }

        currentTrack.value = track ? { ...track } : null
        playStartTime = track ? Date.now() : 0
        updateMusicStatus(track?.id)
      }),
    )

    unsubs.push(
      engine.on('playstate', (state) => {
        isPlaying.value = state === 'playing'
        // Start/stop periodic session persistence
        if (isPlaying.value) {
          startSessionPersist()
        } else {
          stopSessionPersist()
          // Save final state on pause/stop
          if (currentTrack.value) persistSession()
        }
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

        // Exponential backoff: 1s, 2s, then skip
        const backoffMs = Math.min(1000 * Math.pow(2, consecutiveFailures - 1), 4000)
        if (consecutiveFailures < maxConsecutiveFailures - 1 && currentTrack.value?.streamUrl) {
          if (import.meta.env.DEV) {
            console.warn(`[Player] Retrying track in ${backoffMs}ms: ${currentTrack.value.title}`)
          }
          setTimeout(() => {
            engine?.play(currentTrack.value!).catch(() => {})
          }, backoffMs)
          return
        }

        // Last retry failed — skip to the next track
        if (import.meta.env.DEV) {
          console.warn(`[Player] Skipping track due to error: ${currentTrack.value?.title ?? msg}`)
        }
        engine?.next().catch(() => {})
      }),
    )

    // Volume/mute: listen for audio engine changes and sync to store + localStorage
    unsubs.push(
      audioEngine.on('volumechange', (payload) => {
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

  // ── Queue persistence ──────────────────────────────────────────────
  /** Persist both track IDs (primary) and full metadata (cache fallback). */
  function persistQueue(tracks: PlaybackTrack[]) {
    try {
      const data: QueuePersistData = {
        trackIds: tracks.map((t) => t.id),
        hydratedAt: Date.now(),
      }
      localStorage.setItem(QUEUE_STORAGE_KEY, JSON.stringify(data))
      // Cache full metadata as fallback when API is unavailable
      const meta: CachedTrackMeta[] = tracks.map((t) => ({
        id: t.id,
        title: t.title,
        artistName: t.artistName,
        albumTitle: t.albumTitle,
        coverUrl: t.coverUrl,
        durationSeconds: t.durationSeconds,
        streamUrl: t.streamUrl,
      }))
      localStorage.setItem(QUEUE_META_KEY, JSON.stringify(meta))
    } catch {
      // Storage full or blocked — silently skip
    }
  }

  /** Try to restore a previously persisted queue, falling back to cached metadata. */
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
      const cacheRaw = localStorage.getItem(QUEUE_META_KEY)
      let cache: CachedTrackMeta[] = []
      if (cacheRaw) {
        try { cache = JSON.parse(cacheRaw) as CachedTrackMeta[] } catch { /* ignore */ }
      }

      for (let i = 0; i < results.length; i++) {
        const result = results[i]!
        if (result.status === 'fulfilled' && result.value) {
          restored.push(result.value)
        } else if (cache[i]) {
          // Fall back to cached metadata when API fails
          restored.push(cache[i] as PlaybackTrack)
        }
        // Silently skip tracks that no longer exist in either source
      }

      return restored
    } catch {
      return []
    }
  }

  // ── Session persistence (current track, time, queue position, playback state) ──────
  function persistSession() {
    try {
      const data: SessionPersistData = {
        currentTrackId: currentTrack.value?.id ?? null,
        currentTime: currentTime.value,
        updatedAt: Date.now(),
        isPlaying: isPlaying.value,
        shuffleMode: shuffleMode.value,
        repeatMode: repeatMode.value,
        currentIndex: queue.value.findIndex((t) => t.id === currentTrack.value?.id),
      }
      localStorage.setItem(SESSION_STORAGE_KEY, JSON.stringify(data))
    } catch {
      // Silently skip
    }
  }

  /** Save current playback position periodically (every 5s while playing). */
  let _sessionPersistTimer: ReturnType<typeof setInterval> | null = null
  function startSessionPersist() {
    stopSessionPersist()
    _sessionPersistTimer = setInterval(() => {
      if (isPlaying.value && currentTrack.value) {
        persistSession()
      }
    }, 5000)
  }
  function stopSessionPersist() {
    if (_sessionPersistTimer) {
      clearInterval(_sessionPersistTimer)
      _sessionPersistTimer = null
    }
  }

  /** Try to restore session state and auto-resume playback. */
  async function restoreSession(): Promise<boolean> {
    try {
      const raw = localStorage.getItem(SESSION_STORAGE_KEY)
      if (!raw) return false
      const data = JSON.parse(raw) as SessionPersistData
      if (!data.currentTrackId) return false

      // Check freshness: if session is older than 1 hour, don't resume
      const age = Date.now() - data.updatedAt
      if (age > 3600_000) {
        localStorage.removeItem(SESSION_STORAGE_KEY)
        return false
      }

      // Restore mode preferences
      shuffleMode.value = data.shuffleMode || 'off'
      repeatMode.value = data.repeatMode || 'off'
      // Sync modes to the engine's internal state
      if (engine) {
        engine.setShuffleMode(shuffleMode.value)
        // Sync repeat mode: cycle engine to the correct mode
        while (engine.repeatModeValue !== repeatMode.value) {
          engine.toggleRepeat()
        }
      }

      // Restore playback from the persisted track
      try {
        const track = await playerApi.getPlaybackTrack(data.currentTrackId)
        if (!track) return false

        // Preserve queue position when restoring within the existing queue
        const existingIdx = queue.value.findIndex((t) => t.id === data.currentTrackId)
        if (existingIdx >= 0) {
          // Track is already in the restored queue — seek to position within it
          currentTrack.value = track
          currentTime.value = data.currentTime
          if (engine) {
            // Move queue position to the correct index
            const q = queueManager.all()
            const idx = data.currentIndex >= 0 && data.currentIndex < q.length ? data.currentIndex : 0
            if (idx > 0) {
              queueManager.setQueue(q, idx)
              engine?.updateQueue(q)
            }
          }
        } else {
          // Track not in queue — create single-item queue
          queue.value = [track]
          engine?.updateQueue([track])
          currentTrack.value = track
          currentTime.value = data.currentTime
        }

        // Auto-resume if was playing
        if (data.isPlaying && engine) {
          await engine.play(track)
          if (data.currentTime > 0) {
            engine.seek(data.currentTime)
          }
        }
        return true
      } catch {
        return false
      }
    } catch {
      return false
    }
  }

  /** Reset all player state to defaults — used on logout / full teardown. */
  function $reset() {
    stop()
    clearSleepTimer()
    stopSessionPersist()
    // Clean up pagehide listener
    if (_pageHideHandler && typeof window !== 'undefined') {
      window.removeEventListener('pagehide', _pageHideHandler)
      _pageHideHandler = null
    }
    unregisterBeforeUnload()
    if (musicStatusTimer) {
      clearTimeout(musicStatusTimer)
      musicStatusTimer = null
    }
    pendingStatusTrackId = null
    unsubs.forEach((fn) => fn())
    unsubs.length = 0
    initialized = false
    engine = null

    currentTrack.value = null
    queue.value = []
    isPlaying.value = false
    isBuffering.value = false
    isLoadingTrack.value = false
    currentTime.value = 0
    duration.value = 0
    volume.value = 0.85
    muted.value = false
    shuffleMode.value = 'off'
    repeatMode.value = 'off'
    playbackRate.value = 1
    sleepTimerMinutes.value = 0
    crossfadeDuration.value = 0
    audioQuality.value = 'auto'
    error.value = null
    consecutiveFailures = 0
    playbackStopped = false
  }

  /** Reset the consecutive-failure guard so new playback can start fresh. */
  function resetFailureGuard() {
    consecutiveFailures = 0
    playbackStopped = false
  }

  /** Manually retry the current track after a playback error. */
  async function retryCurrentTrack() {
    if (!engine || !currentTrack.value) return

    resetFailureGuard()
    error.value = null
    isLoadingTrack.value = true

    try {
      await engine.play(currentTrack.value)
    } catch (err: unknown) {
      error.value = getErrorMessage(err, 'Could not retry track')
    } finally {
      isLoadingTrack.value = false
    }
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

  /** Append tracks to the queue without replacing existing content. */
  async function appendQueueAndPlay(tracks: PlaybackTrack[], playIndex?: number) {
    initialize()
    if (!engine) return
    if (!tracks.length) return

    resetFailureGuard()
    shuffleMode.value = engine.shuffleModeValue

    await engine.appendQueueAndPlay(tracks, playIndex)
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

  function addToQueue(track: PlaybackTrack) {
    queueManager.addToQueue(track)
    queue.value = queueManager.all()
    engine?.updateQueue(queue.value)
  }

  function playNextInQueue(track: PlaybackTrack) {
    queueManager.playNext(track)
    queue.value = queueManager.all()
    engine?.updateQueue(queue.value)
  }

  function setCrossfadeDuration(seconds: number) {
    crossfadeDuration.value = Math.max(0, seconds)
    engine?.setCrossfadeDuration(crossfadeDuration.value)
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

    nextUpTrack,
    remainingShuffledQueue,
    hasNext,
    hasPrevious,

    initialize,
    playTrack,
    playTrackById,
    setQueueAndPlay,
    appendQueueAndPlay,
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
    addToQueue,
    playNextInQueue,
    retryCurrentTrack,
    setCrossfadeDuration,
    setSleepTimer,
    clearSleepTimer,
    registerBeforeUnload,
    unregisterBeforeUnload,
    restorePersistedQueue,
    persistSession,
    restoreSession,
    $reset,
  }
})
