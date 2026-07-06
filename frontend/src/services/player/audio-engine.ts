export type VolumeChangePayload = {
  volume: number
  muted: boolean
}

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
  volumechange: VolumeChangePayload
  error: Error
}

type AudioEngineEventName = keyof AudioEngineEventMap
type AudioEngineListener<K extends AudioEngineEventName> = (payload: AudioEngineEventMap[K]) => void

type NativeEventName =
  | 'play' | 'pause' | 'ended' | 'waiting' | 'playing' | 'canplay'
  | 'loadedmetadata' | 'durationchange' | 'volumechange' | 'error'

class AudioEngine {
  private audio: HTMLAudioElement
  private listeners: { [K in AudioEngineEventName]?: Set<AudioEngineListener<K>> } = {} as { [K in AudioEngineEventName]?: Set<AudioEngineListener<K>> }
  private animationFrameId: number | null = null
  private lastProgressEmit = 0
  private audioContext: AudioContext | null = null
  private analyser: AnalyserNode | null = null
  private sourceNode: MediaElementAudioSourceNode | null = null
  private gainNode: GainNode | null = null
  private crossfadeDuration = 0
  private nativeHandlers = new Map<NativeEventName, EventListener>()
  private disposed = false

  constructor() {
    this.audio = new Audio()
    this.audio.preload = 'metadata'
    this.audio.crossOrigin = 'anonymous'

    this.bindNativeEvents()
  }

  private ensureAudioContext() {
    if (this.audioContext) return
    this.audioContext = new (window.AudioContext || (window as any as { webkitAudioContext: typeof AudioContext }).webkitAudioContext)()
    this.sourceNode = this.audioContext.createMediaElementSource(this.audio)
    this.analyser = this.audioContext.createAnalyser()
    this.analyser.fftSize = 256
    this.analyser.smoothingTimeConstant = 0.8
    this.gainNode = this.audioContext.createGain()
    this.gainNode.gain.value = 1
    this.sourceNode.connect(this.analyser)
    this.analyser.connect(this.gainNode)
    this.gainNode.connect(this.audioContext.destination)
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
    const set = (this.listeners[event] ??= new Set() as any) as Set<AudioEngineListener<K>>
    set.add(listener)

    return () => {
      set.delete(listener)
    }
  }

  private emit<K extends AudioEngineEventName>(event: K, payload: AudioEngineEventMap[K]) {
    (this.listeners[event] as Set<AudioEngineListener<K>> | undefined)?.forEach((listener) => {
      listener(payload)
    })
  }

  private bindNativeEvents() {
    const onPlay: EventListener = () => {
      this.startProgressLoop()
      this.emit('play', undefined)
    }
    this.nativeHandlers.set('play', onPlay)
    this.audio.addEventListener('play', onPlay)

    const onPause: EventListener = () => {
      this.stopProgressLoop()
      this.emit('pause', undefined)
    }
    this.nativeHandlers.set('pause', onPause)
    this.audio.addEventListener('pause', onPause)

    const onEnded: EventListener = () => {
      this.stopProgressLoop()
      this.emit('ended', undefined)
    }
    this.nativeHandlers.set('ended', onEnded)
    this.audio.addEventListener('ended', onEnded)

    const onWaiting: EventListener = () => {
      this.emit('waiting', undefined)
    }
    this.nativeHandlers.set('waiting', onWaiting)
    this.audio.addEventListener('waiting', onWaiting)

    const onPlaying: EventListener = () => {
      this.emit('playing', undefined)
    }
    this.nativeHandlers.set('playing', onPlaying)
    this.audio.addEventListener('playing', onPlaying)

    const onCanplay: EventListener = () => {
      this.emit('canplay', undefined)
    }
    this.nativeHandlers.set('canplay', onCanplay)
    this.audio.addEventListener('canplay', onCanplay)

    const onLoadedmetadata: EventListener = () => {
      this.emit('loadedmetadata', undefined)
      this.emit('durationchange', {
        duration: this.duration,
      })
    }
    this.nativeHandlers.set('loadedmetadata', onLoadedmetadata)
    this.audio.addEventListener('loadedmetadata', onLoadedmetadata)

    const onDurationchange: EventListener = () => {
      this.emit('durationchange', {
        duration: this.duration,
      })
    }
    this.nativeHandlers.set('durationchange', onDurationchange)
    this.audio.addEventListener('durationchange', onDurationchange)

    const onVolumechange: EventListener = () => {
      this.emit('volumechange', {
        volume: this.audio.volume,
        muted: this.audio.muted,
      })
    }
    this.nativeHandlers.set('volumechange', onVolumechange)
    this.audio.addEventListener('volumechange', onVolumechange)

    const onError: EventListener = () => {
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
    }
    this.nativeHandlers.set('error', onError)
    this.audio.addEventListener('error', onError)
  }

  /** Remove all native event listeners and custom listeners, then stop playback. */
  dispose() {
    if (this.disposed) return
    this.disposed = true

    this.stop()
    this.stopProgressLoop()

    // Remove all native event listeners
    for (const [event, handler] of this.nativeHandlers) {
      this.audio.removeEventListener(event, handler)
    }
    this.nativeHandlers.clear()

    // Clear all custom listeners
    this.listeners = {} as { [K in AudioEngineEventName]?: Set<AudioEngineListener<K>> }

    // Clean up AudioContext
    if (this.audioContext) {
      this.audioContext.close().catch(() => {})
      this.audioContext = null
      this.analyser = null
      this.sourceNode = null
    }
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

    try {
      await this.audio.play()
    } catch {
      // HTMLAudioElement.play() rejects with a DOMException when the media
      // resource is not suitable (e.g. 404). The native 'error' event already
      // fired and will be handled upstream — swallow the rejection here to
      // prevent "Uncaught (in promise) DOMException" from hitting Vue's
      // global error handler.
    }
  }

  /** Set the crossfade duration in seconds (0 = disabled). */
  setCrossfadeDuration(seconds: number) {
    this.crossfadeDuration = Math.max(0, seconds)
  }

  /**
   * Fade the current audio from its current gain to target gain over `duration` seconds.
   * Resolves when the fade completes or immediately if no AudioContext is available.
   */
  fadeTo(targetGain: number, duration: number): Promise<void> {
    return new Promise((resolve) => {
      try {
        this.ensureAudioContext()
      } catch {
        resolve()
        return
      }
      if (!this.gainNode || !this.audioContext || duration <= 0) {
        if (this.gainNode) this.gainNode.gain.value = targetGain
        resolve()
        return
      }

      const currentTime = this.audioContext.currentTime
      this.gainNode.gain.cancelScheduledValues(currentTime)
      this.gainNode.gain.setValueAtTime(this.gainNode.gain.value, currentTime)
      this.gainNode.gain.linearRampToValueAtTime(targetGain, currentTime + duration)

      // Resolve after the ramp completes
      setTimeout(resolve, duration * 1000)
    })
  }

  /** Quick fade out (uses crossfadeDuration or 300ms default). */
  async fadeOut(duration?: number) {
    const dur = duration ?? this.crossfadeDuration || 0.3
    await this.fadeTo(0, dur)
  }

  /** Fade in from 0 to 1 over given duration. */
  async fadeIn(duration?: number) {
    const dur = duration ?? this.crossfadeDuration || 0.3
    // Reset gain to 0 before fading in
    if (this.gainNode) this.gainNode.gain.value = 0
    await this.fadeTo(1, dur)
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

export { AudioEngine }
export const audioEngine = new AudioEngine()
export type { AudioEngine as AudioEngineType }
