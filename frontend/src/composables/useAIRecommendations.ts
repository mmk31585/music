import { ref, computed } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useAIApi } from '@/services/api/ai'
import { usePlayerApi } from '@/services/api/player'
import { usePlayer } from '@/composables/player'
import type { RecommendationTrack, ListeningStats, DiscoverWeeklyResponse } from '@/services/api/recommendation/types'
import type { AIPlaylistResponse, TrackMeta } from '@/services/api/ai/types'
import type { PlaybackTrack } from '@/services/api/player/types'

interface PlayQueueItem {
  id: string
  title: string
  artist?: string | null
  artist_name?: string | null
  album?: string | null
  album_title?: string | null
  cover_url?: string | null
  duration?: number | null
  duration_seconds?: number | null
}

export function useAIRecommendations() {
  const recsApi = useRecommendationsApi()
  const aiApi = useAIApi()
  const playerApi = usePlayerApi()
  const player = usePlayer()

  const popularTracks = ref<RecommendationTrack[]>([])
  const forYouTracks = ref<RecommendationTrack[]>([])
  const discoverWeekly = ref<DiscoverWeeklyResponse | null>(null)
  const listeningStats = ref<ListeningStats | null>(null)
  const aiSimilarTracks = ref<TrackMeta[]>([])
  const moodPlaylist = ref<AIPlaylistResponse | null>(null)

  const loadingPopular = ref(false)
  const loadingForYou = ref(false)
  const loadingDiscoverWeekly = ref(false)
  const loadingStats = ref(false)
  const loadingSimilar = ref(false)
  const loadingMoodPlaylist = ref(false)
  const loadingPersonalized = ref(false)

  async function fetchPopular(limit = 10): Promise<void> {
    loadingPopular.value = true
    try {
      const res = await recsApi.getPopular({ limit })
      popularTracks.value = res.items
    } catch {
      popularTracks.value = []
    } finally {
      loadingPopular.value = false
    }
  }

  async function fetchForYou(limit = 10): Promise<void> {
    loadingForYou.value = true
    try {
      const res = await recsApi.getForYou({ limit })
      forYouTracks.value = res.items
    } catch {
      forYouTracks.value = []
    } finally {
      loadingForYou.value = false
    }
  }

  async function fetchDiscoverWeekly(): Promise<void> {
    loadingDiscoverWeekly.value = true
    try {
      const res = await recsApi.getDiscoverWeekly()
      discoverWeekly.value = res
    } catch {
      discoverWeekly.value = null
    } finally {
      loadingDiscoverWeekly.value = false
    }
  }

  async function fetchListeningStats(period = 'month'): Promise<void> {
    loadingStats.value = true
    try {
      const res = await recsApi.getListeningStats({ period })
      listeningStats.value = res
    } catch {
      listeningStats.value = null
    } finally {
      loadingStats.value = false
    }
  }

  async function fetchSimilarByMood(trackId: string, mood?: string): Promise<void> {
    loadingSimilar.value = true
    try {
      const res = await aiApi.similarByMood(trackId, mood)
      aiSimilarTracks.value = res
    } catch {
      aiSimilarTracks.value = []
    } finally {
      loadingSimilar.value = false
    }
  }

  async function fetchMoodPlaylist(mood: string, limit = 20): Promise<void> {
    loadingMoodPlaylist.value = true
    try {
      const res = await aiApi.generatePlaylist({ prompt: `mood: ${mood}`, mood, limit })
      moodPlaylist.value = res
    } catch {
      moodPlaylist.value = null
    } finally {
      loadingMoodPlaylist.value = false
    }
  }

  const personalizedTracks = ref<RecommendationTrack[]>([])

  async function fetchPersonalized(limit = 10): Promise<void> {
    loadingPersonalized.value = true
    try {
      const res = await recsApi.getPersonalized({ limit })
      personalizedTracks.value = res.items
    } catch {
      personalizedTracks.value = []
    } finally {
      loadingPersonalized.value = false
    }
  }

  function formatDuration(seconds?: number | null): string {
    if (seconds == null || seconds <= 0) return '0:00'
    const m = Math.floor(seconds / 60)
    const s = Math.floor(seconds % 60)
    return `${m}:${s.toString().padStart(2, '0')}`
  }

  function getEnergyColor(energy: number): string {
    if (energy >= 0.8) return 'text-red-400'
    if (energy >= 0.6) return 'text-orange-400'
    if (energy >= 0.4) return 'text-yellow-400'
    if (energy >= 0.2) return 'text-lime-400'
    return 'text-blue-400'
  }

  function getMoodGradient(mood: string): string {
    const gradients: Record<string, string> = {
      energetic: 'from-orange-500/20 via-red-500/10 to-transparent',
      happy: 'from-yellow-400/20 via-amber-500/10 to-transparent',
      chill: 'from-cyan-400/20 via-teal-500/10 to-transparent',
      calm: 'from-blue-400/20 via-indigo-500/10 to-transparent',
      sad: 'from-indigo-400/20 via-violet-500/10 to-transparent',
      focus: 'from-emerald-400/20 via-green-500/10 to-transparent',
      romantic: 'from-pink-400/20 via-rose-500/10 to-transparent',
      intense: 'from-red-500/20 via-orange-600/10 to-transparent',
      confident: 'from-purple-400/20 via-fuchsia-500/10 to-transparent',
      sleep: 'from-slate-400/20 via-blue-500/10 to-transparent',
    }
    return gradients[mood] || 'from-[#1db954]/20 via-emerald-500/10 to-transparent'
  }

  function buildPlayQueue(tracks: PlayQueueItem[], startIndex = 0): void {
    const playbackTracks: PlaybackTrack[] = tracks.map((t) => ({
      id: t.id,
      title: t.title,
      artistName: t.artist ?? t.artist_name ?? 'Unknown artist',
      albumTitle: t.album ?? t.album_title,
      coverUrl: t.cover_url,
      durationSeconds: t.duration ?? t.duration_seconds,
      streamUrl: playerApi.getTrackStreamUrl(t.id),
    }))
    player.setQueueAndPlay(playbackTracks, startIndex)
  }

  function getStreamUrl(trackId: string): string {
    return playerApi.getTrackStreamUrl(trackId)
  }

  const dominantMood = computed(() => {
    if (!listeningStats.value) return 'chill'
    const genres = listeningStats.value.top_genres
    if (!genres.length) return 'chill'
    const top = genres[0]!.genre_name.toLowerCase()
    if (/rock|metal|dubstep|drum/.test(top)) return 'energetic'
    if (/pop|dance|edm|house/.test(top)) return 'happy'
    if (/jazz|soul|r&b|ambient/.test(top)) return 'chill'
    if (/classical|piano|acoustic/.test(top)) return 'calm'
    if (/blues|folk|country/.test(top)) return 'sad'
    if (/lo[- ]?fi|study|focus/.test(top)) return 'focus'
    return 'chill'
  })

  const discoveryScore = computed(() => {
    if (!listeningStats.value) return '—'
    return `${Math.round(listeningStats.value.discovery_score)}%`
  })

  return {
    popularTracks,
    forYouTracks,
    personalizedTracks,
    discoverWeekly,
    listeningStats,
    aiSimilarTracks,
    moodPlaylist,
    loadingPopular,
    loadingForYou,
    loadingDiscoverWeekly,
    loadingStats,
    loadingSimilar,
    loadingMoodPlaylist,
    loadingPersonalized,
    fetchPopular,
    fetchForYou,
    fetchDiscoverWeekly,
    fetchListeningStats,
    fetchSimilarByMood,
    fetchMoodPlaylist,
    fetchPersonalized,
    formatDuration,
    getEnergyColor,
    getMoodGradient,
    buildPlayQueue,
    getStreamUrl,
    dominantMood,
    discoveryScore,
  }
}

export function usePlayFromRecommendation() {
  const player = usePlayer()
  const playerApi = usePlayerApi()

  function playFromRecommendation(tracks: PlayQueueItem[], startIndex = 0): void {
    const playbackTracks: PlaybackTrack[] = tracks.map((t) => ({
      id: t.id,
      title: t.title,
      artistName: t.artist ?? t.artist_name ?? 'Unknown artist',
      albumTitle: t.album ?? t.album_title,
      coverUrl: t.cover_url,
      durationSeconds: t.duration ?? t.duration_seconds,
      streamUrl: playerApi.getTrackStreamUrl(t.id),
    }))
    player.setQueueAndPlay(playbackTracks, startIndex)
  }

  return { playFromRecommendation }
}
