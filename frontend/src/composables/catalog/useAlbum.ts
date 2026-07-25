import { ref, computed } from 'vue'
import { useAlbumsApi, type Album } from '@/services/api/catalog/albums'
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
  interface AlbumArtist { id: string | number; role?: string; name: string }
  const albumArtists = ref<AlbumArtist[]>([])
  const relatedAlbums = ref<Album[]>([])
  const isLiked = ref(false)
  const loading = ref(false)
  const error = ref<any>(null)

  const mainArtist = computed(() => {
    return (
      albumArtists.value.find((a: AlbumArtist) => a.role === 'main' || a.role === 'primary') ??
      albumArtists.value[0] ??
      null
    )
  })

  const featuredArtists = computed(() => {
    return albumArtists.value.filter((a: AlbumArtist) => a.role && !['main', 'primary'].includes(a.role))
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

      if (!albumData) {
        throw new Error('Album not found')
      }

      album.value = albumData

      const [tracksData, likedAlbums] = await Promise.all([
        tracksApi.getTracks({ album_id: id }).catch(() => [] as Track[]),
        libraryApi.getLikedAlbums().catch(() => []),
      ])

      albumArtists.value = []

      tracks.value = Array.isArray(tracksData) ? tracksData : []

      const artistId = mainArtist.value?.id ?? albumData.artist_id
      if (artistId) {
        const [artistData, relatedAlbumsData] = await Promise.all([
          artistsApi.getArtist(artistId).catch(() => null),
          albumsApi.getAlbums({ artist_id: artistId }).catch(() => [] as Album[]),
        ])
        artist.value = artistData
        relatedAlbums.value = Array.isArray(relatedAlbumsData)
          ? relatedAlbumsData.filter((a) => String(a.id) !== String(id))
          : []
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
