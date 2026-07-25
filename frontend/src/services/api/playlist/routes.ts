import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { PlaylistApiRoutes } from './enums'
import {
  PlaylistListItemSchema,
  PlaylistDetailSchema,
  CollaboratorResponseSchema,
  type PlaylistListItem,
  type PlaylistDetail,
  type CollaboratorResponse,
  type CreatePlaylistPayload,
  type AddTrackPayload,
} from './types'

export const usePlaylistsApi = () => {
  const listPublicPlaylists = async (config?: UseRequestConfig<PlaylistListItem[]>) => {
    return useRequest<PlaylistListItem, true>(
      PlaylistApiRoutes.LIST_PUBLIC,
      { method: 'GET' },
      { schema: PlaylistListItemSchema, silent: true, ...config },
    )
  }

  const getMyPlaylists = async (config?: UseRequestConfig<PlaylistListItem[]>) => {
    return useRequest<PlaylistListItem, true>(
      PlaylistApiRoutes.LIST_MY,
      { method: 'GET' },
      { schema: PlaylistListItemSchema, silent: true, ...config },
    )
  }

  const getPlaylist = async (id: string, config?: UseRequestConfig<PlaylistDetail>) => {
    return useRequest<PlaylistDetail>(
      PlaylistApiRoutes.GET.replace(':playlistId', id),
      { method: 'GET' },
      { schema: PlaylistDetailSchema, silent: true, ...config },
    )
  }

  const createPlaylist = async (
    payload: CreatePlaylistPayload,
    config?: UseRequestConfig<PlaylistDetail>,
  ) => {
    return useRequest<PlaylistDetail>(
      PlaylistApiRoutes.CREATE,
      { method: 'POST', data: payload },
      { schema: PlaylistDetailSchema, silent: false, ...config },
    )
  }

  /**
   * Update playlist metadata.
   * Backend accepts: name, description, cover_url, is_public.
   * Collaborative toggle uses the dedicated setCollaborative endpoint.
   */
  const updatePlaylist = async (
    id: string,
    payload: Partial<CreatePlaylistPayload>,
    config?: UseRequestConfig<PlaylistDetail>,
  ) => {
    return useRequest<PlaylistDetail>(
      PlaylistApiRoutes.UPDATE.replace(':playlistId', id),
      { method: 'PUT', data: payload },
      { schema: PlaylistDetailSchema, silent: false, ...config },
    )
  }

  const deletePlaylist = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      PlaylistApiRoutes.DELETE.replace(':playlistId', id),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  const addTrack = async (
    playlistId: string,
    payload: AddTrackPayload,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.ADD_TRACK.replace(':playlistId', playlistId),
      { method: 'POST', data: payload },
      { silent: false, ...config },
    )
  }

  const removeTrack = async (
    playlistId: string,
    trackId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.REMOVE_TRACK.replace(':playlistId', playlistId).replace(
        ':trackId',
        trackId,
      ),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  /**
   * Reorder a single track to a new position.
   * Backend expects { track_id, new_position } (1-indexed).
   */
  const reorderTracks = async (
    playlistId: string,
    trackId: string,
    newPosition: number,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.REORDER_TRACKS.replace(':playlistId', playlistId),
      { method: 'PUT', data: { track_id: trackId, new_position: newPosition } },
      { silent: false, ...config },
    )
  }

  /**
   * Toggle collaborative mode.
   * Backend expects { collaborative: boolean }.
   */
  const setCollaborative = async (
    playlistId: string,
    collaborative: boolean,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.SET_COLLABORATIVE.replace(':playlistId', playlistId),
      { method: 'PUT', data: { collaborative } },
      { silent: false, ...config },
    )
  }

  /**
   * List collaborators for a collaborative playlist.
   * Backend returns: { success: true, data: [{ user_id, added_at, is_creator }] }
   */
  const listCollaborators = async (
    playlistId: string,
    config?: UseRequestConfig<CollaboratorResponse[]>,
  ) => {
    return useRequest<CollaboratorResponse, true>(
      PlaylistApiRoutes.LIST_COLLABORATORS.replace(':playlistId', playlistId),
      { method: 'GET' },
      { schema: CollaboratorResponseSchema, silent: true, ...config },
    )
  }

  return {
    listPublicPlaylists,
    getMyPlaylists,
    getPlaylist,
    createPlaylist,
    updatePlaylist,
    deletePlaylist,
    addTrack,
    removeTrack,
    reorderTracks,
    setCollaborative,
    listCollaborators,
  }
}
