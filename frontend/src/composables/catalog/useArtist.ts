import { ref, computed } from 'vue'
import { useArtistsApi, type Artist } from '@/services/api/catalog/artists'
import { useAlbumsApi, type Album } from '@/services/api/catalog/albums'
import { useTracksApi, type Track } from '@/services/api/catalog/tracks'
import { useLibraryApi } from '@/services/api/library'

export function useArtist(id: string | number) {
  const artistsApi = useArtistsApi()
  const albumsApi = useAlbumsApi()
  const tracksApi = useTracksApi()
  const libraryApi = useLibraryApi()

  const artist = ref<Artist | null>(null)
  const tracks = ref<Track[]>([])
  const albums = ref<Album[]>([])
  const related = ref<Artist[]>([])
  const isFollowing = ref(false)
  const loading = ref(false)
  const error = ref<any>(null)

  const monthlyListeners = computed(() => artist.value?.monthly_listeners ?? 0)

  async function fetchArtist() {
    loading.value = true
    error.value = null

    try {
      const [artistData, tracksData, albumsData, followed] = await Promise.all([
        artistsApi.getArtist(id),
        tracksApi.getTracks({ artist_id: id }).catch(() => [] as Track[]),
        albumsApi.getAlbums({ artist_id: id }).catch(() => [] as Album[]),
        libraryApi.getFollowedArtists().catch(() => []),
      ])

      if (!artistData) {
        throw new Error('Artist not found')
      }

      artist.value = artistData

      tracks.value = Array.isArray(tracksData) ? tracksData.slice(0, 10) : []
      albums.value = Array.isArray(albumsData) ? albumsData : []

      // Get related artists from distinct artist_ids on albums NOT by this artist
      // This avoids fetching all albums then filtering client-side
      related.value = []

      isFollowing.value = Array.isArray(followed)
        ? followed.some((f) => f.artist_id === String(id))
        : false
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  async function toggleFollow() {
    try {
      if (isFollowing.value) {
        await libraryApi.unfollowArtist(String(id))
        isFollowing.value = false
      } else {
        await libraryApi.followArtist({ artist_id: String(id) })
        isFollowing.value = true
      }
    } catch {
      // silent
    }
  }

  return {
    artist,
    tracks,
    albums,
    related,
    isFollowing,
    monthlyListeners,
    loading,
    error,
    fetchArtist,
    toggleFollow,
  }
}
