import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { LibraryApiRoutes } from './enums'
import {
  AddPlayHistoryPayloadSchema,
  LibraryAlbumLikePayloadSchema,
  LibraryAlbumSchema,
  LibraryArtistFollowPayloadSchema,
  LibraryArtistSchema,
  LibraryTrackLikePayloadSchema,
  LibraryTrackSchema,
  PlayHistoryItemSchema,
  type AddPlayHistoryPayload,
  type LibraryAlbum,
  type LibraryAlbumLikePayload,
  type LibraryArtist,
  type LibraryArtistFollowPayload,
  type LibraryTrack,
  type LibraryTrackLikePayload,
  type PlayHistoryItem,
} from './types'

export const useLibraryApi = () => {
  const getLikedTracks = async (config?: UseRequestConfig<LibraryTrack[]>) => {
    return useRequest<LibraryTrack, true>(
      LibraryApiRoutes.TRACKS,
      { method: 'GET' },
      {
        schema: LibraryTrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const likeTrack = async (payload: LibraryTrackLikePayload, config?: UseRequestConfig<void>) => {
    LibraryTrackLikePayloadSchema.parse(payload)

    return useRequest<void>(
      LibraryApiRoutes.LIKE_TRACK,
      {
        method: 'POST',
        data: payload,
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const unlikeTrack = async (trackId: string | number, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      LibraryApiRoutes.UNLIKE_TRACK.replace(':trackId', String(trackId)),
      {
        method: 'DELETE',
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const getLikedAlbums = async (config?: UseRequestConfig<LibraryAlbum[]>) => {
    return useRequest<LibraryAlbum, true>(
      LibraryApiRoutes.ALBUMS,
      { method: 'GET' },
      {
        schema: LibraryAlbumSchema,
        silent: true,
        ...config,
      },
    )
  }

  const likeAlbum = async (payload: LibraryAlbumLikePayload, config?: UseRequestConfig<void>) => {
    LibraryAlbumLikePayloadSchema.parse(payload)

    return useRequest<void>(
      LibraryApiRoutes.LIKE_ALBUM,
      {
        method: 'POST',
        data: payload,
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const unlikeAlbum = async (albumId: string | number, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      LibraryApiRoutes.UNLIKE_ALBUM.replace(':albumId', String(albumId)),
      {
        method: 'DELETE',
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const getFollowedArtists = async (config?: UseRequestConfig<LibraryArtist[]>) => {
    return useRequest<LibraryArtist, true>(
      LibraryApiRoutes.ARTISTS,
      { method: 'GET' },
      {
        schema: LibraryArtistSchema,
        silent: true,
        ...config,
      },
    )
  }

  const followArtist = async (
    payload: LibraryArtistFollowPayload,
    config?: UseRequestConfig<void>,
  ) => {
    LibraryArtistFollowPayloadSchema.parse(payload)

    return useRequest<void>(
      LibraryApiRoutes.FOLLOW_ARTIST,
      {
        method: 'POST',
        data: payload,
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const unfollowArtist = async (artistId: string | number, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      LibraryApiRoutes.UNFOLLOW_ARTIST.replace(':artistId', String(artistId)),
      {
        method: 'DELETE',
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const addPlayHistory = async (
    payload: AddPlayHistoryPayload,
    config?: UseRequestConfig<void>,
  ) => {
    AddPlayHistoryPayloadSchema.parse(payload)

    return useRequest<void>(
      LibraryApiRoutes.HISTORY,
      {
        method: 'POST',
        data: payload,
      },
      {
        silent: true,
        ...config,
      },
    )
  }

  const getPlayHistory = async (config?: UseRequestConfig<PlayHistoryItem[]>) => {
    return useRequest<PlayHistoryItem, true>(
      LibraryApiRoutes.HISTORY,
      { method: 'GET' },
      {
        schema: PlayHistoryItemSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getRecentlyPlayed = async (config?: UseRequestConfig<PlayHistoryItem[]>) => {
    return useRequest<PlayHistoryItem, true>(
      LibraryApiRoutes.RECENTLY_PLAYED,
      { method: 'GET' },
      {
        schema: PlayHistoryItemSchema,
        silent: true,
        ...config,
      },
    )
  }

  return {
    getLikedTracks,
    likeTrack,
    unlikeTrack,

    getLikedAlbums,
    likeAlbum,
    unlikeAlbum,

    getFollowedArtists,
    followArtist,
    unfollowArtist,

    addPlayHistory,
    getPlayHistory,
    getRecentlyPlayed,
  }
}
