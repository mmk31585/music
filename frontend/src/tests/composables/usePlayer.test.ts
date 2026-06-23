import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePlayer } from '@/composables/player/usePlayer'

// Mock the player engine
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

describe('usePlayer', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('returns player state and actions', () => {
    const player = usePlayer()
    expect(player).toHaveProperty('currentTrack')
    expect(player).toHaveProperty('queue')
    expect(player).toHaveProperty('isPlaying')
    expect(player).toHaveProperty('playTrack')
    expect(player).toHaveProperty('resume')
    expect(player).toHaveProperty('pause')
    expect(player).toHaveProperty('seek')
    expect(player).toHaveProperty('playNext')
    expect(player).toHaveProperty('playPrevious')
    expect(player).toHaveProperty('setShuffleMode')
    expect(player).toHaveProperty('setVolume')
    expect(player).toHaveProperty('toggleShuffle')
    expect(player).toHaveProperty('toggleRepeat')
  })

  it('currentTrack is null initially', () => {
    const player = usePlayer()
    expect(player.currentTrack.value).toBeNull()
  })

  it('isPlaying is false initially', () => {
    const player = usePlayer()
    expect(player.isPlaying.value).toBe(false)
  })

  it('shuffleMode starts as off', () => {
    const player = usePlayer()
    expect(player.shuffleMode.value).toBe('off')
  })

  it('repeatMode starts as off', () => {
    const player = usePlayer()
    expect(player.repeatMode.value).toBe('off')
  })

  it('playTrack delegates to store', () => {
    const player = usePlayer()
    // Should not throw
    expect(() => player.playTrack(mockTrack)).not.toThrow()
  })

  it('pause delegates to store', () => {
    const player = usePlayer()
    expect(() => player.pause()).not.toThrow()
  })

  it('setShuffleMode updates mode', () => {
    const player = usePlayer()
    player.setShuffleMode('queue')
    expect(player.shuffleMode.value).toBe('queue')
  })
})
