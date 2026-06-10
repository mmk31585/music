import { ref, computed } from 'vue'
import { useAlbumsApi, type Album, type AlbumArtist } from '@/services/api/catalog/albums'
import { useArtistsApi, type Artist } from '@/services/api/catalog/artists'
import { useTracksApi, type Track } from '@/services/api/catalog/tracks'
import { useLibraryApi } from '@/services/api/library'

export function useAlbum(id: string | number) {
  const albumsApi = useAlbumsApi()
  const artistsApi = useArtistsApi()
  const tracksApi = useTracksApi()
  const libraryApi = useLibraryApi()

  const album = ref<Album | null>(null)
  const tracks = ref<Track[]>([])
  const artist = ref<Artist | null>(null)
  const albumArtists = ref<AlbumArtist[]>([])
  const relatedAlbums = ref<Album[]>([])
  const isLiked = ref(false)
  const loading = ref(false)
  const error = ref<unknown>(null)

  const mainArtist = computed(() => {
    return (
      albumArtists.value.find((a) => a.role === 'main' || a.role === 'primary') ??
      albumArtists.value[0] ??
      null
    )
  })

  const featuredArtists = computed(() => {
    return albumArtists.value.filter((a) => a.role && !['main', 'primary'].includes(a.role))
  })

  const totalDuration = computed(() => {
    return tracks.value.reduce((sum, t) => sum + (t.duration_seconds ?? 0), 0)
  })

  const formattedDuration = computed(() => {
    const min = Math.floor(totalDuration.value / 60)
    return `${min} min`
  })

  async function fetchAlbum() {
    loading.value = true
    error.value = null

    try {
      const albumData = await albumsApi.getAlbum(id)
      album.value = albumData

      const [allTracks, likedAlbums, artistsData, allAlbums] = await Promise.all([
        tracksApi.getTracks({ limit: 100 }).catch(() => [] as Track[]),
        libraryApi.getLikedAlbums({ limit: 50 }).catch(() => []),
        albumsApi.getAlbumArtists(id).catch(() => [] as AlbumArtist[]),
        albumsApi.getAlbums({ limit: 100 }).catch(() => [] as Album[]),
      ])

      albumArtists.value = artistsData

      const albumTracks = Array.isArray(allTracks)
        ? allTracks.filter((t) => String(t.album_id) === String(id))
        : []

      tracks.value = albumTracks

      const artistId = mainArtist.value?.id ?? albumData.artist_id
      if (artistId) {
        artistsApi
          .getArtist(artistId)
          .then((a) => {
            artist.value = a
          })
          .catch(() => {})

        if (Array.isArray(allAlbums)) {
          relatedAlbums.value = allAlbums.filter(
            (a) => String(a.id) !== String(id) && String(a.artist_id) === String(artistId),
          )
        }
      }

      isLiked.value = Array.isArray(likedAlbums)
        ? likedAlbums.some((a) => a.album_id === String(id))
        : false
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  async function toggleLike() {
    try {
      if (isLiked.value) {
        await libraryApi.unlikeAlbum(String(id))
        isLiked.value = false
      } else {
        await libraryApi.likeAlbum({ album_id: String(id) })
        isLiked.value = true
      }
    } catch {
      // silent
    }
  }

  return {
    album,
    tracks,
    artist,
    albumArtists,
    relatedAlbums,
    mainArtist,
    featuredArtists,
    isLiked,
    totalDuration,
    formattedDuration,
    loading,
    error,
    fetchAlbum,
    toggleLike,
  }
}
