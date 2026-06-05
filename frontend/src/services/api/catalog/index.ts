// Re-export all modules
export * from './common'
export * from './artists'
export * from './albums'
export * from './genres'
export * from './tracks'
export * from './search'

// Combined composable for backward compatibility
// You can use this during migration, then gradually switch to individual composables
import { useArtistsApi } from './artists'
import { useAlbumsApi } from './albums'
import { useGenresApi } from './genres'
import { useTracksApi } from './tracks'
import { useSearchApi } from './search'

/**
 * @deprecated Use individual API composables instead:
 * - useArtistsApi()
 * - useAlbumsApi()
 * - useGenresApi()
 * - useTracksApi()
 * - useSearchApi()
 */
export const useCatalogApi = () => {
  const artists = useArtistsApi()
  const albums = useAlbumsApi()
  const genres = useGenresApi()
  const tracks = useTracksApi()
  const search = useSearchApi()

  return {
    // Artists
    getArtists: artists.getArtists,
    adminGetArtists: artists.adminGetArtists,
    adminCreateArtist: artists.adminCreateArtist,
    adminUpdateArtist: artists.adminUpdateArtist,
    adminDeleteArtist: artists.adminDeleteArtist,

    // Albums
    getAlbums: albums.getAlbums,
    adminGetAlbums: albums.adminGetAlbums,
    adminCreateAlbum: albums.adminCreateAlbum,
    adminUpdateAlbum: albums.adminUpdateAlbum,
    adminDeleteAlbum: albums.adminDeleteAlbum,

    // Genres
    getGenres: genres.getGenres,
    adminGetGenres: genres.adminGetGenres,
    adminCreateGenre: genres.adminCreateGenre,
    adminUpdateGenre: genres.adminUpdateGenre,
    adminDeleteGenre: genres.adminDeleteGenre,

    // Tracks
    getTracks: tracks.getTracks,
    getTrack: tracks.getTrack,
    adminGetTracks: tracks.adminGetTracks,
    adminCreateTrack: tracks.adminCreateTrack,
    adminUploadTrackWithAudio: tracks.adminUploadTrackWithAudio,
    adminUpdateTrack: tracks.adminUpdateTrack,
    adminDeleteTrack: tracks.adminDeleteTrack,

    // Search
    searchCatalog: search.searchCatalog,
  }
}
