import type { PlaybackTrack } from '@/services/api/player'
import { AudioEngine, audioEngine } from './audio-engine'
import { QueueManager, queueManager } from './queue-manager'
import { PreloadManager, preloadManager } from './preload-manager'
import { updateMediaSession, setMediaSessionPlaybackState } from './media-session'

export type ShuffleMode = 'off' | 'queue' | 'catalog' | 'similar'
export type RepeatMode = 'off' | 'one' | 'all'

export interface PlayerEngineProviders {
  fetchTrack?: (id: string) => Promise<PlaybackTrack>
  fetchRandomTracks?: (limit: number) => Promise<PlaybackTrack[]>
  fetchSimilarTracks?: (trackId: string, limit: number) => Promise<PlaybackTrack[]>
  onPlayHistory?: (trackId: string, duration: number) => void
}

export interface PlayerEngineTimeupdate {
  currentTime: number
  duration: number
}

type EventPayloads = {
  trackchange: PlaybackTrack | null
  playstate: 'playing' | 'paused' | 'stopped'
  timeupdate: PlayerEngineTimeupdate
  buffering: boolean
  queuechange: PlaybackTrack[]
  error: string
}

type EventName = keyof EventPayloads
type Listener<K extends EventName> = (payload: EventPayloads[K]) => void

export class PlayerEngine {
  private audio: AudioEngine
  private queue: QueueManager
  private preload: PreloadManager
  private providers: PlayerEngineProviders

  private shuffleOrder: number[] = []
  private shufflePosition = -1
  private shuffleMode: ShuffleMode = 'off'
  private repeatMode: RepeatMode = 'off'

  private _currentTrack: PlaybackTrack | null = null

  private listeners = new Map<EventName, Set<Listener<any>>>()
  private unsubs: (() => void)[] = []

  constructor(audio?: AudioEngine, queue?: QueueManager, preload?: PreloadManager, providers?: PlayerEngineProviders) {
    this.audio = audio ?? audioEngine
    this.queue = queue ?? queueManager
    this.preload = preload ?? preloadManager
    this.providers = providers ?? {}

    this.bindAudioEvents()
  }

  private bindAudioEvents() {
    this.unsubs.push(
      this.audio.on('play', () => {
        setMediaSessionPlaybackState('playing')
        this.emit('playstate', 'playing')
      }),
    )

    this.unsubs.push(
      this.audio.on('pause', () => {
        setMediaSessionPlaybackState('paused')
        this.emit('playstate', 'paused')
      }),
    )

    this.unsubs.push(
      this.audio.on('waiting', () => {
        this.emit('buffering', true)
      }),
    )

    this.unsubs.push(
      this.audio.on('playing', () => {
        this.emit('buffering', false)
      }),
    )

    this.unsubs.push(
      this.audio.on('canplay', () => {
        this.emit('buffering', false)
      }),
    )

    this.unsubs.push(
      this.audio.on('timeupdate', (payload) => {
        this.emit('timeupdate', {
          currentTime: payload.currentTime,
          duration: payload.duration,
        })
      }),
    )

    this.unsubs.push(
      this.audio.on('durationchange', (payload) => {
        if (payload.duration) {
          this.emit('timeupdate', {
            currentTime: this.audio.currentTime,
            duration: payload.duration,
          })
        }
      }),
    )

    this.unsubs.push(
      this.audio.on('ended', async () => {
        if (this.repeatMode === 'one') {
          this.audio.seek(0)
          await this.audio.play()
          return
        }
        await this.next()
      }),
    )

    this.unsubs.push(
      this.audio.on('error', async (err) => {
        const message = err.message || 'Playback failed'
        setMediaSessionPlaybackState('none')
        this.emit('error', message)
        this.emit('playstate', 'paused')
        // NOTE: Do NOT auto-advance the queue here. The store's error
        // handler owns the skip-or-stop decision so it can count consecutive
        // failures and halt playback when the backend is unreachable.
        // If we call next() here the store loses control and we get an
        // infinite retry loop when every track in the queue is broken.
      }),
    )
  }

  private emit<K extends EventName>(event: K, payload: EventPayloads[K]) {
    this.listeners.get(event)?.forEach((listener) => {
      listener(payload)
    })
  }

  on<K extends EventName>(event: K, listener: Listener<K>): () => void {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event)!.add(listener)
    return () => {
      this.listeners.get(event)?.delete(listener)
    }
  }

  private setCurrentTrack(track: PlaybackTrack | null) {
    this._currentTrack = track
    this.emit('trackchange', track)
  }

  get currentTrack(): PlaybackTrack | null {
    return this._currentTrack
  }

  get currentTime(): number {
    return this.audio.currentTime
  }

  get duration(): number {
    return this.audio.duration
  }

  get isPlaying(): boolean {
    return !this.audio.paused
  }

  get queueAll(): PlaybackTrack[] {
    return this.queue.all()
  }

  private buildShuffleOrder() {
    const q = this.queue.all()
    const indices = q.map((_, i) => i)
    for (let i = indices.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1))
      ;[indices[i]!, indices[j]!] = [indices[j]!, indices[i]!]
    }
    this.shuffleOrder = indices
    this.shufflePosition = -1
  }

  private async doPlay(track: PlaybackTrack) {
    this.emit('buffering', true)

    this.setCurrentTrack({
      ...track,
      coverUrl: track.coverUrl ?? undefined,
    })

    this.queue.setCurrent(track)
    this.emit('queuechange', this.queue.all())

    updateMediaSession(track, {
      play: () => this.resume(),
      pause: () => this.pause(),
      next: () => this.next(),
      previous: () => this.previous(),
      seek: (t) => this.seek(t),
    })

    await this.audio.play(track.streamUrl)

    this.providers.onPlayHistory?.(track.id, track.durationSeconds ?? 0)

    const nextTrack = this.queue.getNext()
    if (nextTrack) {
      this.preload.preload(nextTrack.streamUrl)
    }

    this.emit('buffering', false)
  }

  async play(track: PlaybackTrack) {
    await this.doPlay(track)
  }

  async playById(id: string) {
    if (!this.providers.fetchTrack) {
      this.emit('error', 'Cannot fetch track: no provider configured')
      return
    }

    this.emit('buffering', true)

    try {
      const track = await this.providers.fetchTrack(id)
      await this.doPlay(track)
    } catch (err: any) {
      const message = err instanceof Error ? err.message : String(err)
      this.emit('error', message || 'Could not load track')
    } finally {
      this.emit('buffering', false)
    }
  }

  async setQueueAndPlay(tracks: PlaybackTrack[], startIndex = 0) {
    if (!tracks.length) return

    this.queue.setQueue(tracks, startIndex)
    this.emit('queuechange', this.queue.all())

    if (this.shuffleMode === 'queue') {
      this.buildShuffleOrder()
    }

    const track = tracks[startIndex]
    if (track) {
      await this.doPlay(track)
    }
  }

  async resume() {
    await this.audio.play()
  }

  pause() {
    this.audio.pause()
  }

  stop() {
    this.audio.stop()
    this.setCurrentTrack(null)
    this.emit('playstate', 'stopped')
  }

  seek(seconds: number) {
    this.audio.seek(seconds)
  }

  async next() {
    if (this.shuffleMode === 'queue') {
      await this.nextShuffleQueue()
      return
    }

    if (this.shuffleMode === 'catalog') {
      await this.nextCatalog()
      return
    }

    if (this.shuffleMode === 'similar') {
      await this.nextSimilar()
      return
    }

    await this.nextSequential()
  }

  private async nextShuffleQueue() {
    const q = this.queue.all()
    this.shufflePosition++

    const currentId = this._currentTrack?.id
    while (this.shufflePosition < this.shuffleOrder.length) {
      const nextRealIdx = this.shuffleOrder[this.shufflePosition]!
      const candidate = q[nextRealIdx]
      if (candidate && candidate.id !== currentId) break
      this.shufflePosition++
    }

    if (this.shufflePosition < this.shuffleOrder.length) {
      const nextRealIdx: number = this.shuffleOrder[this.shufflePosition]!
      const nextTrack = q[nextRealIdx]
      if (nextTrack) {
        this.queue.setCurrent(nextTrack)
        await this.doPlay(nextTrack)
        return
      }
    }

    if (this.repeatMode === 'all') {
      this.buildShuffleOrder()
      if (this.shuffleOrder.length > 0) {
        this.shufflePosition = 0
        const nextRealIdx = this.shuffleOrder[0]!
        const nextTrack = q[nextRealIdx]
        if (nextTrack) {
          this.queue.setCurrent(nextTrack)
          await this.doPlay(nextTrack)
          return
        }
      }
    }

    this.audio.stop()
  }

  private async nextCatalog() {
    if (!this.providers.fetchRandomTracks) return
    try {
      const randomTracks = await this.providers.fetchRandomTracks(20)
      if (randomTracks && randomTracks.length > 0) {
        const pick = randomTracks[Math.floor(Math.random() * randomTracks.length)]!
        await this.playById(pick.id)
        return
      }
    } catch { /* fall through to sequential */ }
    await this.nextSequential()
  }

  private async nextSimilar() {
    if (!this.providers.fetchSimilarTracks || !this._currentTrack) return
    try {
      const similar = await this.providers.fetchSimilarTracks(this._currentTrack.id, 10)
      if (similar && similar.length > 0) {
        const pick = similar[Math.floor(Math.random() * similar.length)]!
        await this.playById(pick.id)
        return
      }
    } catch { /* fall through to sequential */ }
    await this.nextSequential()
  }

  private async nextSequential() {
    const nextTrack = this.queue.next()

    if (!nextTrack) {
      if (this.repeatMode === 'all') {
        this.queue.setQueue(this.queue.all(), 0)
        const firstTrack = this.queue.next()
        if (firstTrack) {
          await this.doPlay(firstTrack)
          return
        }
      }
      this.audio.stop()
      return
    }

    await this.doPlay(nextTrack)
  }

  async previous() {
    if (this.audio.currentTime > 4) {
      this.audio.seek(0)
      return
    }

    if (this.shuffleMode === 'queue') {
      await this.prevShuffleQueue()
      return
    }

    if (this.shuffleMode === 'catalog') {
      await this.prevCatalog()
      return
    }

    if (this.shuffleMode === 'similar') {
      await this.prevSimilar()
      return
    }

    await this.prevSequential()
  }

  private async prevShuffleQueue() {
    this.shufflePosition--
    if (this.shufflePosition < 0 && this.repeatMode === 'all') {
      this.shufflePosition = this.shuffleOrder.length - 1
    }
    if (this.shufflePosition >= 0) {
      const q = this.queue.all()
      const prevRealIdx = this.shuffleOrder[this.shufflePosition]!
      const prevTrack = q[prevRealIdx]
      if (prevTrack) {
        this.queue.setCurrent(prevTrack)
        await this.doPlay(prevTrack)
        return
      }
    }
    this.audio.seek(0)
  }

  private async prevCatalog() {
    if (!this.providers.fetchRandomTracks) return
    try {
      const randomTracks = await this.providers.fetchRandomTracks(20)
      if (randomTracks && randomTracks.length > 0) {
        const pick = randomTracks[Math.floor(Math.random() * randomTracks.length)]!
        await this.playById(pick.id)
        return
      }
    } catch { /* fall through */ }
    this.audio.seek(0)
  }

  private async prevSimilar() {
    if (!this.providers.fetchSimilarTracks || !this._currentTrack) return
    try {
      const similar = await this.providers.fetchSimilarTracks(this._currentTrack.id, 10)
      if (similar && similar.length > 0) {
        const pick = similar[Math.floor(Math.random() * similar.length)]!
        await this.playById(pick.id)
        return
      }
    } catch { /* fall through */ }
    this.audio.seek(0)
  }

  private async prevSequential() {
    const previousTrack = this.queue.previous()
    if (!previousTrack) {
      this.audio.seek(0)
      return
    }
    await this.doPlay(previousTrack)
  }

  setVolume(value: number) {
    this.audio.setVolume(value)
  }

  toggleMute() {
    this.audio.setMuted(!this.audio.element.muted)
  }

  get volume(): number {
    return this.audio.element.volume
  }

  get muted(): boolean {
    return this.audio.element.muted
  }

  setShuffleMode(mode: ShuffleMode) {
    this.shuffleMode = mode
    this.shuffleOrder = []
    this.shufflePosition = -1
    if (mode === 'queue') {
      this.buildShuffleOrder()
    }
  }

  toggleShuffle() {
    const modes: ShuffleMode[] = ['off', 'queue', 'catalog', 'similar']
    const idx = modes.indexOf(this.shuffleMode)
    this.setShuffleMode(modes[(idx + 1) % modes.length]!)
  }

  get shuffleModeValue(): ShuffleMode {
    return this.shuffleMode
  }

  toggleRepeat() {
    if (this.repeatMode === 'off') this.repeatMode = 'all'
    else if (this.repeatMode === 'all') this.repeatMode = 'one'
    else this.repeatMode = 'off'
  }

  get repeatModeValue(): RepeatMode {
    return this.repeatMode
  }

  setPlaybackRate(rate: number) {
    this.audio.setPlaybackRate(rate)
  }

  updateQueue(newQueue: PlaybackTrack[]) {
    this.queue.replaceAll(newQueue)
    this.emit('queuechange', this.queue.all())
    if (this.shuffleMode === 'queue') {
      this.buildShuffleOrder()
    }
  }

  getAnalyserNode() {
    return this.audio.getAnalyserNode()
  }

  destroy() {
    this.unsubs.forEach((fn) => fn())
    this.unsubs = []
    this.listeners.clear()
    this.audio.stop()
  }
}
