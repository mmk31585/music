import { ref, computed } from 'vue'
import {
  useTracksApi,
  type Track,
} from '@/services/api/catalog/tracks'
import { useRecommendationsApi, type RecommendationTrack } from '@/services/api/recommendation'
import { useLibraryApi } from '@/services/api/library'
import { useLyricsApi, type Lyrics } from '@/services/api/lyrics'

export function useTrack(id: string | number) {
  const tracksApi = useTracksApi()
  const recsApi = useRecommendationsApi()
  const libraryApi = useLibraryApi()
  const lyricsApi = useLyricsApi()

  const track = ref<Track | null>(null)
  const similarTracks = ref<RecommendationTrack[]>([])
  const lyrics = ref<Lyrics | null>(null)
  interface TrackArtist { id: string | number; role?: string; name: string; artistId?: string | number }
  interface TrackCredit { id: string | number; role?: string; name?: string; artistName?: string; creditType?: string }
  const trackArtists = ref<TrackArtist[]>([])
  const trackCredits = ref<TrackCredit[]>([])
  const isLiked = ref(false)
  const loading = ref(false)
  const error = ref<any>(null)

  const mainArtist = computed(() => {
    if (!trackArtists.value.length) return null
    return (
      trackArtists.value.find((a: TrackArtist) => a.role === 'main' || a.role === 'primary') ??
      trackArtists.value[0]!
    )
  })

  const featuredArtists = computed(() => {
    return trackArtists.value.filter((a: TrackArtist) => a.role && !['main', 'primary'].includes(a.role))
  })

  const genreList = computed(() => {
    return track.value?.genres ?? []
  })

  async function fetchTrack() {
    loading.value = true
    error.value = null

    try {
      const [trackData, likedTracks] = await Promise.all([
        tracksApi.getTrack(id),
        libraryApi.getLikedTracks().catch(() => []),
      ])

      if (!trackData) {
        throw new Error('Track not found')
      }

      track.value = trackData

      isLiked.value = Array.isArray(likedTracks)
        ? likedTracks.some((t) => t.track_id === String(id))
        : false

      trackArtists.value = (trackData.artists || []).map((a: any) => ({
        id: a.artist_id,
        artistId: a.artist_id ?? undefined,
        name: a.name,
        role: a.role,
      }))
      trackCredits.value = []

      recsApi
        .getSimilar(String(id), { limit: 8 })
        .then((res) => {
          if (res?.items) similarTracks.value = res.items.slice(0, 8)
        })
        .catch(() => {})

      lyricsApi
        .getTrackLyrics(id)
        .then((l) => {
          lyrics.value = l
        })
        .catch(() => {})
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  async function toggleLike() {
    try {
      if (isLiked.value) {
        await libraryApi.unlikeTrack(String(id))
        isLiked.value = false
      } else {
        await libraryApi.likeTrack({ track_id: String(id) })
        isLiked.value = true
      }
    } catch {
      // silent
    }
  }

  return {
    track,
    similarTracks,
    lyrics,
    trackArtists,
    trackCredits,
    mainArtist,
    featuredArtists,
    genreList,
    isLiked,
    loading,
    error,
    fetchTrack,
    toggleLike,
  }
}
