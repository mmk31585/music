import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { PlaylistApiRoutes } from './enums'
import {
  PlaylistListItemSchema,
  PlaylistDetailSchema,
  type PlaylistListItem,
  type PlaylistDetail,
  type CreatePlaylistPayload,
  type AddTrackPayload,
} from './types'

export const usePlaylistsApi = () => {
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

  const updatePlaylist = async (
    id: string,
    payload: Partial<CreatePlaylistPayload & { collaborative: boolean }>,
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

  const reorderTracks = async (
    playlistId: string,
    trackIds: string[],
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.REORDER_TRACKS.replace(':playlistId', playlistId),
      { method: 'PUT', data: { track_ids: trackIds } },
      { silent: false, ...config },
    )
  }

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

  const addCollaborator = async (
    playlistId: string,
    userId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.ADD_COLLABORATOR.replace(':playlistId', playlistId),
      { method: 'POST', data: { user_id: userId } },
      { silent: false, ...config },
    )
  }

  const removeCollaborator = async (
    playlistId: string,
    userId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      PlaylistApiRoutes.REMOVE_COLLABORATOR.replace(':playlistId', playlistId).replace(
        ':userId',
        userId,
      ),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  const listCollaborators = async (playlistId: string, config?: UseRequestConfig<any>) => {
    return useRequest<any>(
      PlaylistApiRoutes.LIST_COLLABORATORS.replace(':playlistId', playlistId),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  return {
    getMyPlaylists,
    getPlaylist,
    createPlaylist,
    updatePlaylist,
    deletePlaylist,
    addTrack,
    removeTrack,
    reorderTracks,
    setCollaborative,
    addCollaborator,
    removeCollaborator,
    listCollaborators,
  }
}
