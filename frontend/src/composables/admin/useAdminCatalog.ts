import { computed } from 'vue'
import { useAdminTracks } from './useAdminTracks'
import { useAdminArtists } from './useAdminArtists'
import { useAdminAlbums } from './useAdminAlbums'
import { useAdminGenres } from './useAdminGenres'

export function useAdminCatalog() {
  const tracks = useAdminTracks()
  const artists = useAdminArtists()
  const albums = useAdminAlbums()
  const genres = useAdminGenres()

  const loading = computed(() => {
    return (
      tracks.loading.value ||
      artists.loading.value ||
      albums.loading.value ||
      genres.loading.value
    )
  })

  const saving = computed(() => {
    return (
      tracks.saving.value ||
      artists.saving.value ||
      albums.saving.value ||
      genres.saving.value
    )
  })

  const deleting = computed(() => {
    return (
      tracks.deleting.value ||
      artists.deleting.value ||
      albums.deleting.value ||
      genres.deleting.value
    )
  })

  async function fetchAll() {
    await Promise.all([
      tracks.fetchTracks(),
      artists.fetchArtists(),
      albums.fetchAlbums(),
      genres.fetchGenres(),
    ])
  }

  function resetAll() {
    tracks.resetTracks()
    artists.resetArtists()
    albums.resetAlbums()
    genres.resetGenres()
  }

  return {
    tracks,
    artists,
    albums,
    genres,

    loading,
    saving,
    deleting,

    fetchAll,
    resetAll,
  }
}
