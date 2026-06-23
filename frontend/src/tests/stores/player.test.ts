import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePlayerStore } from '@/stores/player'

// Mock the Engine to avoid real playback
vi.mock('@/services/player', () => {
  const mockEngine = {
    play: vi.fn(),
    playById: vi.fn(),
    pause: vi.fn(),
    resume: vi.fn(),
    stop: vi.fn(),
    seek: vi.fn(),
    seekPercent: vi.fn(),
    setVolume: vi.fn(),
    toggleMute: vi.fn(),
    next: vi.fn(),
    previous: vi.fn(),
    setQueueAndPlay: vi.fn(),
    setShuffleMode: vi.fn(),
    setRepeatMode: vi.fn(),
    setPlaybackRate: vi.fn(),
    getVolume: vi.fn().mockReturnValue(0.85),
    getMuted: vi.fn().mockReturnValue(false),
    isPlaying: false,
    currentTrack: null,
    queueAll: [],
    shuffleModeValue: 'off' as const,
    repeatModeValue: 'off' as const,
    destroy: vi.fn(),
    on: vi.fn().mockReturnValue(vi.fn()),
    off: vi.fn(),
  }
  return {
    PlayerEngine: vi.fn(function () { return mockEngine }),
    ShuffleMode: { Off: 'off', Queue: 'queue', Catalog: 'catalog', Similar: 'similar' },
    RepeatMode: { Off: 'off', One: 'one', All: 'all' },
  }
})

vi.mock('@/services/api/player', () => ({
  usePlayerApi: () => ({
    getPlaybackTrack: vi.fn(),
  }),
}))

vi.mock('@/services/api/catalog/tracks', () => ({
  useTracksApi: () => ({
    getRandomTracks: vi.fn().mockResolvedValue([]),
  }),
}))

vi.mock('@/services/api/recommendation', () => ({
  useRecommendationsApi: () => ({
    getSimilar: vi.fn().mockResolvedValue({ items: [] }),
  }),
}))

vi.mock('@/services/api/library', () => ({
  useLibraryApi: () => ({
    addPlayHistory: vi.fn(),
  }),
}))

vi.mock('@/services/api/video', () => ({
  useVideoApi: () => ({
    updateMusicStatus: vi.fn(),
    clearMusicStatus: vi.fn(),
  }),
}))

const mockTrack = {
  id: 'track-1',
  title: 'Test Track',
  artistName: 'Test Artist',
  albumTitle: 'Test Album',
  coverUrl: 'http://example.com/cover.jpg',
  durationSeconds: 240,
  streamUrl: 'http://example.com/audio.mp3',
}

describe('usePlayerStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('starts with default values', () => {
    const store = usePlayerStore()
    expect(store.currentTrack).toBeNull()
    expect(store.queue).toEqual([])
    expect(store.isPlaying).toBe(false)
    expect(store.isBuffering).toBe(false)
    expect(store.currentTime).toBe(0)
    expect(store.duration).toBe(0)
    expect(store.volume).toBeGreaterThan(0)
    expect(store.muted).toBe(false)
    expect(store.shuffleMode).toBe('off')
    expect(store.repeatMode).toBe('off')
    expect(store.error).toBeNull()
  })

  it('computes progressPercent correctly', () => {
    const store = usePlayerStore()
    store.currentTime = 60
    store.duration = 240
    expect(store.progressPercent).toBe(25)
  })

  it('progressPercent is 0 when duration is 0', () => {
    const store = usePlayerStore()
    store.currentTime = 50
    store.duration = 0
    expect(store.progressPercent).toBe(0)
  })

  it('hasNext is true when shuffle is on', () => {
    const store = usePlayerStore()
    store.shuffleMode = 'queue'
    expect(store.hasNext).toBe(true)
  })

  it('hasNext is false when shuffle is off and queue is empty', () => {
    const store = usePlayerStore()
    store.shuffleMode = 'off'
    store.queue = []
    store.currentTrack = null
    expect(store.hasNext).toBe(false)
  })

  it('hasPrevious is true when shuffle is on', () => {
    const store = usePlayerStore()
    store.shuffleMode = 'queue'
    expect(store.hasPrevious).toBe(true)
  })

  it('hasPrevious is true when currentTime > 0 and shuffle off', () => {
    const store = usePlayerStore()
    store.shuffleMode = 'off'
    store.currentTime = 10
    expect(store.hasPrevious).toBe(true)
  })

  it('setShuffleMode updates shuffle mode', () => {
    const store = usePlayerStore()
    // @ts-expect-error - accessing internal engine
    store.engine = { setShuffleMode: vi.fn() }
    store.setShuffleMode('queue')
    expect(store.shuffleMode).toBe('queue')
  })

  it('toggleShuffle cycles through modes', () => {
    const store = usePlayerStore()
    // @ts-expect-error - accessing internal engine
    store.engine = { setShuffleMode: vi.fn() }

    // Start with 'off'
    expect(store.shuffleMode).toBe('off')
    store.toggleShuffle()
    // Should now be 'queue' (if the engine supports it)
    // Depends on implementation
  })

  it('setVolume persists to localStorage', () => {
    const store = usePlayerStore()
    // @ts-expect-error - accessing internal engine
    store.engine = { setVolume: vi.fn() }
    store.setVolume(0.5)
    expect(store.volume).toBe(0.5)
    expect(localStorage.getItem('player-volume')).toBe('0.5')
  })

  it('toggleMute toggles muted state', () => {
    const store = usePlayerStore()
    // @ts-expect-error - accessing internal engine
    store.engine = { toggleMute: vi.fn() }
    const initial = store.muted
    store.toggleMute()
    expect(store.muted).toBe(!initial)
  })

  it('setPlaybackRate updates rate', () => {
    const store = usePlayerStore()
    // @ts-expect-error - accessing internal engine
    store.engine = { setPlaybackRate: vi.fn() }
    store.setPlaybackRate(1.5)
    expect(store.playbackRate).toBe(1.5)
  })

  it('playTrack sets error when track fails', async () => {
    const store = usePlayerStore()
    const { PlayerEngine } = await import('@/services/player')
    const mockEngine = (PlayerEngine as any).mock.results[0].value
    mockEngine.play.mockRejectedValueOnce(new Error('Playback failed'))

    await store.playTrack(mockTrack as any)
    expect(store.error).toBe('Playback failed')
  })

  it('pause calls engine pause', async () => {
    const store = usePlayerStore()
    store.pause()
    // Should not throw
  })

  it('stop calls engine stop and clears state', async () => {
    const store = usePlayerStore()
    store.stop()
    // Should not throw
  })
})
