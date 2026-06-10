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
  const error = ref<unknown>(null)

  const monthlyListeners = computed(() => artist.value?.monthly_listeners ?? 0)

  async function fetchArtist() {
    loading.value = true
    error.value = null

    try {
      const [artistData, tracksData, albumsData, followed] = await Promise.all([
        artistsApi.getArtist(id),
        tracksApi.getTracks({ limit: 100 }).catch(() => [] as Track[]),
        albumsApi.getAlbums({ limit: 100 }).catch(() => [] as Album[]),
        libraryApi.getFollowedArtists({ limit: 50 }).catch(() => []),
      ])

      artist.value = artistData

      const artistTracks = Array.isArray(tracksData)
        ? tracksData.filter((t) => String(t.artist_id) === String(id))
        : []

      const artistAlbums = Array.isArray(albumsData)
        ? albumsData.filter((a) => String(a.artist_id) === String(id))
        : []

      tracks.value = artistTracks.slice(0, 10)
      albums.value = artistAlbums

      const otherArtistIds = [
        ...new Set(
          (Array.isArray(albumsData) ? albumsData : [])
            .filter((a) => String(a.artist_id) !== String(id))
            .map((a) => a.artist_id)
            .filter(Boolean),
        ),
      ].slice(0, 6)

      const relatedArtists = await Promise.all(
        otherArtistIds.map((aid) => artistsApi.getArtist(aid!).catch(() => null)),
      )
      related.value = relatedArtists.filter((a): a is Artist => a !== null)

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
