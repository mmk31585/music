type AudioEngineEventMap = {
  play: void
  pause: void
  ended: void
  waiting: void
  playing: void
  canplay: void
  loadedmetadata: void
  timeupdate: {
    currentTime: number
    duration: number
  }
  durationchange: {
    duration: number
  }
  volumechange: {
    volume: number
    muted: boolean
  }
  error: Error
}

type AudioEngineEventName = keyof AudioEngineEventMap
type AudioEngineListener<K extends AudioEngineEventName> = (payload: AudioEngineEventMap[K]) => void

class AudioEngine {
  private audio: HTMLAudioElement
  private listeners = new Map<AudioEngineEventName, Set<Function>>()
  private animationFrameId: number | null = null
  private lastProgressEmit = 0
  private audioContext: AudioContext | null = null
  private analyser: AnalyserNode | null = null
  private sourceNode: MediaElementAudioSourceNode | null = null

  constructor() {
    this.audio = new Audio()
    this.audio.preload = 'metadata'
    this.audio.crossOrigin = 'anonymous'

    this.bindNativeEvents()
  }

  private ensureAudioContext() {
    if (this.audioContext) return
    this.audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
    this.sourceNode = this.audioContext.createMediaElementSource(this.audio)
    this.analyser = this.audioContext.createAnalyser()
    this.analyser.fftSize = 256
    this.analyser.smoothingTimeConstant = 0.8
    this.sourceNode.connect(this.analyser)
    this.analyser.connect(this.audioContext.destination)
  }

  get element() {
    return this.audio
  }

  get currentTime() {
    return this.audio.currentTime || 0
  }

  get duration() {
    return Number.isFinite(this.audio.duration) ? this.audio.duration : 0
  }

  get paused() {
    return this.audio.paused
  }

  get src() {
    return this.audio.src
  }

  on<K extends AudioEngineEventName>(event: K, listener: AudioEngineListener<K>) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }

    this.listeners.get(event)?.add(listener)

    return () => {
      this.listeners.get(event)?.delete(listener)
    }
  }

  private emit<K extends AudioEngineEventName>(event: K, payload: AudioEngineEventMap[K]) {
    this.listeners.get(event)?.forEach((listener) => {
      listener(payload)
    })
  }

  private bindNativeEvents() {
    this.audio.addEventListener('play', () => {
      this.startProgressLoop()
      this.emit('play', undefined)
    })

    this.audio.addEventListener('pause', () => {
      this.stopProgressLoop()
      this.emit('pause', undefined)
    })

    this.audio.addEventListener('ended', () => {
      this.stopProgressLoop()
      this.emit('ended', undefined)
    })

    this.audio.addEventListener('waiting', () => {
      this.emit('waiting', undefined)
    })

    this.audio.addEventListener('playing', () => {
      this.emit('playing', undefined)
    })

    this.audio.addEventListener('canplay', () => {
      this.emit('canplay', undefined)
    })

    this.audio.addEventListener('loadedmetadata', () => {
      this.emit('loadedmetadata', undefined)
      this.emit('durationchange', {
        duration: this.duration,
      })
    })

    this.audio.addEventListener('durationchange', () => {
      this.emit('durationchange', {
        duration: this.duration,
      })
    })

    this.audio.addEventListener('volumechange', () => {
      this.emit('volumechange', {
        volume: this.audio.volume,
        muted: this.audio.muted,
      })
    })

    this.audio.addEventListener('error', () => {
      const mediaError = this.audio.error
      let message = 'Audio playback error'
      if (mediaError) {
        switch (mediaError.code) {
          case MediaError.MEDIA_ERR_ABORTED:
            message = 'Playback was aborted'
            break
          case MediaError.MEDIA_ERR_NETWORK:
            message = 'Network error — check your connection'
            break
          case MediaError.MEDIA_ERR_DECODE:
            message = 'Could not decode audio'
            break
          case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
            message = 'Audio format not supported or stream unavailable'
            break
        }
      }
      const error = new Error(message)
      this.emit('error', error)
    })
  }

  async load(src: string) {
    if (this.audio.src === src) return

    this.audio.pause()
    this.audio.src = src
    this.audio.load()
  }

  async play(src?: string) {
    if (src && this.audio.src !== src) {
      await this.load(src)
    }

    await this.audio.play()
  }

  pause() {
    this.audio.pause()
  }

  stop() {
    this.audio.pause()
    this.audio.currentTime = 0
  }

  seek(seconds: number) {
    if (!Number.isFinite(seconds)) return

    this.audio.currentTime = Math.max(0, Math.min(seconds, this.duration || seconds))

    this.emit('timeupdate', {
      currentTime: this.currentTime,
      duration: this.duration,
    })
  }

  setVolume(volume: number) {
    this.audio.volume = Math.max(0, Math.min(1, volume))
  }

  setMuted(muted: boolean) {
    this.audio.muted = muted
  }

  private startProgressLoop() {
    this.stopProgressLoop()

    const tick = (timestamp: number) => {
      // Throttle progress updates to avoid rerendering everything 60 times/sec.
      if (timestamp - this.lastProgressEmit >= 250) {
        this.lastProgressEmit = timestamp

        this.emit('timeupdate', {
          currentTime: this.currentTime,
          duration: this.duration,
        })
      }

      this.animationFrameId = window.requestAnimationFrame(tick)
    }

    this.animationFrameId = window.requestAnimationFrame(tick)
  }

  private stopProgressLoop() {
    if (this.animationFrameId !== null) {
      window.cancelAnimationFrame(this.animationFrameId)
      this.animationFrameId = null
    }
  }

  setPlaybackRate(rate: number) {
    this.audio.playbackRate = rate
  }

  getAnalyserNode(): AnalyserNode | null {
    try {
      this.ensureAudioContext()
    } catch {
      return null
    }
    return this.analyser
  }
}

export type AudioQuality = 'auto' | 'low' | 'medium' | 'high' | 'lossless'

export const audioEngine = new AudioEngine()
export type { AudioEngine }
