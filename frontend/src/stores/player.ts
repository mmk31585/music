import { computed, ref, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import type { Track } from '@/services/api/catalog'

type RepeatMode = 'off' | 'one' | 'all'

type PersistedPlayerState = {
  queue: Track[]
  currentTrackId: string | number | null
  currentTime: number
  volume: number
  repeatMode: RepeatMode
  shuffle: boolean
}

type AudioCleanup = () => void

const STORAGE_KEY = 'music_player_state'

function readPersistedState(): PersistedPlayerState | null {
  if (typeof window === 'undefined') return null

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    return raw ? (JSON.parse(raw) as PersistedPlayerState) : null
  } catch {
    window.localStorage.removeItem(STORAGE_KEY)
    return null
  }
}

function writePersistedState(state: PersistedPlayerState) {
  if (typeof window === 'undefined') return

  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  } catch {
    // Playback persistence is nice-to-have; quota/privacy failures should not break audio.
  }
}

function clamp(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value))
}

export const usePlayerStore = defineStore('player', () => {
  const persisted = readPersistedState()

  const audio = shallowRef<HTMLAudioElement | null>(null)
  const audioCleanup = shallowRef<AudioCleanup | null>(null)
  const queue = ref<Track[]>(persisted?.queue ?? [])
  const currentIndex = ref(
    persisted?.currentTrackId
      ? queue.value.findIndex((track) => track.id === persisted.currentTrackId)
      : -1,
  )
  const isPlaying = ref(false)
  const isBuffering = ref(false)
  const currentTime = ref(persisted?.currentTime ?? 0)
  const duration = ref(0)
  const volume = ref(persisted?.volume ?? 0.8)
  const repeatMode = ref<RepeatMode>(persisted?.repeatMode ?? 'off')
  const shuffle = ref(persisted?.shuffle ?? false)
  const error = ref('')

  if (currentIndex.value < 0 && queue.value.length > 0) {
    currentIndex.value = 0
  }

  const currentTrack = computed(() => {
    if (currentIndex.value < 0) return null
    return queue.value[currentIndex.value] ?? null
  })

  const canPlayCurrent = computed(() => !!currentTrack.value?.audio_url)
  const hasNext = computed(
    () =>
      (shuffle.value && queue.value.length > 1) ||
      currentIndex.value < queue.value.length - 1 ||
      (repeatMode.value === 'all' && queue.value.length > 1),
  )
  const hasPrevious = computed(
    () =>
      currentTime.value > 3 ||
      currentIndex.value > 0 ||
      (repeatMode.value === 'all' && queue.value.length > 1),
  )

  function persist() {
    writePersistedState({
      queue: queue.value,
      currentTrackId: currentTrack.value?.id ?? null,
      currentTime: currentTime.value,
      volume: volume.value,
      repeatMode: repeatMode.value,
      shuffle: shuffle.value,
    })
  }

  function addAudioListener<K extends keyof HTMLMediaElementEventMap>(
    element: HTMLAudioElement,
    event: K,
    listener: (event: HTMLMediaElementEventMap[K]) => void,
  ) {
    element.addEventListener(event, listener)
    return () => element.removeEventListener(event, listener)
  }

  function disposeAudio() {
    const element = audio.value

    audioCleanup.value?.()
    audioCleanup.value = null

    if (element) {
      element.pause()
      element.removeAttribute('src')
      element.load()
    }

    audio.value = null
    isPlaying.value = false
    isBuffering.value = false
  }

  function ensureAudio() {
    if (audio.value || typeof Audio === 'undefined') return audio.value

    const element = new Audio()
    element.preload = 'metadata'
    element.volume = volume.value
    const cleanup: AudioCleanup[] = []

    cleanup.push(
      addAudioListener(element, 'loadstart', () => {
        isBuffering.value = true
        error.value = ''
      }),
    )
    cleanup.push(
      addAudioListener(element, 'loadedmetadata', () => {
        duration.value = Number.isFinite(element.duration) ? element.duration : 0
        if (currentTime.value > 0 && currentTime.value < duration.value) {
          element.currentTime = currentTime.value
        }
        isBuffering.value = false
      }),
    )
    cleanup.push(
      addAudioListener(element, 'waiting', () => {
        isBuffering.value = true
      }),
    )
    cleanup.push(
      addAudioListener(element, 'playing', () => {
        isBuffering.value = false
        isPlaying.value = true
      }),
    )
    cleanup.push(
      addAudioListener(element, 'pause', () => {
        isPlaying.value = false
        persist()
      }),
    )
    cleanup.push(
      addAudioListener(element, 'timeupdate', () => {
        currentTime.value = element.currentTime
        if (Math.floor(element.currentTime) % 5 === 0) {
          persist()
        }
      }),
    )
    cleanup.push(
      addAudioListener(element, 'durationchange', () => {
        duration.value = Number.isFinite(element.duration) ? element.duration : 0
      }),
    )
    cleanup.push(
      addAudioListener(element, 'ended', () => {
        if (repeatMode.value === 'one') {
          void playCurrent(true)
          return
        }

        if (hasNext.value) {
          void playNext()
        } else {
          isPlaying.value = false
          currentTime.value = 0
          persist()
        }
      }),
    )
    cleanup.push(
      addAudioListener(element, 'error', () => {
        error.value = 'This track could not be played.'
        isBuffering.value = false
        isPlaying.value = false
        persist()
      }),
    )

    audioCleanup.value = () => {
      for (const remove of cleanup) {
        remove()
      }
    }
    audio.value = element
    loadCurrent(false)

    return element
  }

  function loadCurrent(keepTime: boolean) {
    const element = ensureAudio()
    const track = currentTrack.value

    if (!element || !track?.audio_url) {
      duration.value = track?.duration_seconds ?? 0
      return
    }

    if (element.src !== track.audio_url) {
      element.src = track.audio_url
      element.load()
      duration.value = track.duration_seconds ?? 0
      if (!keepTime) currentTime.value = 0
    }
  }

  async function playCurrent(keepTime = false) {
    const element = ensureAudio()
    const track = currentTrack.value

    if (!element || !track) return
    if (!track.audio_url) {
      error.value = 'This track has no audio file yet.'
      return
    }

    loadCurrent(keepTime)
    error.value = ''

    try {
      await element.play()
      isPlaying.value = true
      persist()
    } catch {
      error.value = 'Playback was blocked. Press play again.'
      isPlaying.value = false
    }
  }

  function setQueue(tracks: Track[], startTrack?: Track) {
    queue.value = tracks.filter((track) => !!track.audio_url)
    currentIndex.value = startTrack
      ? queue.value.findIndex((track) => track.id === startTrack.id)
      : currentIndex.value

    if (currentIndex.value < 0 && queue.value.length > 0) {
      currentIndex.value = 0
    }

    persist()
  }

  async function playTrack(track: Track, tracks: Track[] = queue.value) {
    setQueue(tracks.length > 0 ? tracks : [track], track)
    currentTime.value = 0
    loadCurrent(false)
    await playCurrent(false)
  }

  async function togglePlay() {
    const element = ensureAudio()
    if (!element) return

    if (isPlaying.value) {
      element.pause()
      return
    }

    await playCurrent(true)
  }

  async function playNext() {
    if (queue.value.length === 0) return

    if (shuffle.value && queue.value.length > 1) {
      let nextIndex = currentIndex.value
      while (nextIndex === currentIndex.value) {
        nextIndex = Math.floor(Math.random() * queue.value.length)
      }
      currentIndex.value = nextIndex
    } else if (currentIndex.value < queue.value.length - 1) {
      currentIndex.value += 1
    } else if (repeatMode.value === 'all') {
      currentIndex.value = 0
    } else {
      return
    }

    currentTime.value = 0
    loadCurrent(false)
    await playCurrent(false)
  }

  async function playPrevious() {
    const element = ensureAudio()
    if (!element || queue.value.length === 0) return

    if (currentTime.value > 3) {
      seek(0)
      return
    }

    if (currentIndex.value > 0) {
      currentIndex.value -= 1
    } else if (repeatMode.value === 'all') {
      currentIndex.value = queue.value.length - 1
    }

    currentTime.value = 0
    loadCurrent(false)
    await playCurrent(false)
  }

  function seek(seconds: number) {
    const element = ensureAudio()
    const nextTime = clamp(seconds, 0, duration.value || currentTrack.value?.duration_seconds || 0)

    currentTime.value = nextTime
    if (element) {
      element.currentTime = nextTime
    }
    persist()
  }

  function setVolume(nextVolume: number) {
    volume.value = clamp(nextVolume, 0, 1)
    const element = ensureAudio()
    if (element) {
      element.volume = volume.value
    }
    persist()
  }

  function toggleShuffle() {
    shuffle.value = !shuffle.value
    persist()
  }

  function cycleRepeat() {
    repeatMode.value =
      repeatMode.value === 'off' ? 'all' : repeatMode.value === 'all' ? 'one' : 'off'
    persist()
  }

  function restore() {
    ensureAudio()
  }

  function dispose() {
    persist()
    disposeAudio()
  }

  return {
    queue,
    currentTrack,
    currentIndex,
    isPlaying,
    isBuffering,
    currentTime,
    duration,
    volume,
    repeatMode,
    shuffle,
    error,
    canPlayCurrent,
    hasNext,
    hasPrevious,
    restore,
    dispose,
    setQueue,
    playTrack,
    togglePlay,
    playNext,
    playPrevious,
    seek,
    setVolume,
    toggleShuffle,
    cycleRepeat,
  }
})
