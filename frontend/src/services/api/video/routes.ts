import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { apiReplaceParams } from '@/utils/api-replace-params'
import { VideoApiRoutes } from './enums'
import {
  ExploreResponseSchema,
  VideoJobStatusSchema,
  UploadEditResponseSchema,
  type ExploreResponse,
  type MusicStatusResponse,
  type CreateEditPayload,
  type UploadEditResponse,
  type LikeVisibilityPayload,
  type LikedTrackItem,
  type VideoJobStatus,
  type VideoItem,
} from './types'
import type { AxiosRequestConfig } from 'axios'

export const useVideoApi = () => {
  // ── Admin methods ────────────────────────────────────────────────────
  const adminGetVideos = async (
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<ExploreResponse>,
  ) => {
    return useRequest<ExploreResponse>(
      '/admin/videos',
      { method: 'GET', params },
      { schema: ExploreResponseSchema, silent: true, ...config },
    )
  }

  const adminUpdateVideo = async (
    id: string,
    payload: {
      title?: string
      description?: string
      type?: string
      is_public?: boolean
      is_approved?: boolean
      status?: string
    },
    config?: UseRequestConfig<VideoItem>,
  ) => {
    return useRequest<VideoItem>(
      `/admin/videos/${id}`,
      { method: 'PUT', data: payload },
      config,
    )
  }

  const adminDeleteVideo = async (
    id: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      `/admin/videos/${id}`,
      { method: 'DELETE' },
      config,
    )
  }

  const adminApproveVideo = async (
    id: string,
    config?: UseRequestConfig<VideoItem>,
  ) => {
    return useRequest<VideoItem>(
      `/admin/videos/${id}/approve`,
      { method: 'POST' },
      config,
    )
  }

  // ── Admin comment methods ────────────────────────────────────────────
  interface AdminCommentItem {
    id: string
    video_id: string
    user_id: string
    author_name: string
    content: string
    created_at: string
    updated_at: string
  }

  const adminGetVideoComments = async (
    videoId: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: AdminCommentItem[]; total: number }>,
  ) => {
    return useRequest<{ items: AdminCommentItem[]; total: number }>(
      `/admin/videos/${videoId}/comments`,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const adminUpdateComment = async (
    commentId: string,
    payload: { content: string },
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      `/admin/comments/${commentId}`,
      { method: 'PUT', data: payload },
      config,
    )
  }

  const adminDeleteComment = async (
    commentId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      `/admin/comments/${commentId}`,
      { method: 'DELETE' },
      config,
    )
  }

  // ── Public / user methods ────────────────────────────────────────────
  const getExploreVideos = async (
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<ExploreResponse>,
  ) => {
    return useRequest<ExploreResponse>(
      VideoApiRoutes.EXPLORE,
      { method: 'GET', params },
      { schema: ExploreResponseSchema, silent: true, ...config },
    )
  }

  const getTrackVideos = async (
    trackId: string,
    config?: UseRequestConfig<ExploreResponse>,
  ) => {
    return useRequest<ExploreResponse>(
      VideoApiRoutes.TRACK_VIDEOS.replace(':trackId', trackId),
      { method: 'GET' },
      { schema: ExploreResponseSchema, silent: true, ...config },
    )
  }

  const getVideoJobStatus = async (
    jobId: string,
    config?: UseRequestConfig<VideoJobStatus>,
  ) => {
    return useRequest<VideoJobStatus>(
      VideoApiRoutes.VIDEO_JOB.replace(':jobId', jobId),
      { method: 'GET' },
      { schema: VideoJobStatusSchema, silent: true, ...config },
    )
  }

  const likeVideo = async (
    videoId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      VideoApiRoutes.VIDEO_LIKE.replace(':id', videoId),
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const unlikeVideo = async (
    videoId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      VideoApiRoutes.VIDEO_LIKE.replace(':id', videoId),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  const viewVideo = async (
    videoId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      VideoApiRoutes.VIDEO_VIEW.replace(':id', videoId),
      { method: 'POST' },
      config,
    )
  }

  const uploadEdit = async (
    payload: CreateEditPayload,
    config?: UseRequestConfig<UploadEditResponse>,
  ) => {
    return useRequest<UploadEditResponse>(
      VideoApiRoutes.UPLOAD_EDIT,
      { method: 'POST', data: payload },
      { schema: UploadEditResponseSchema, silent: false, ...config },
    )
  }

  const uploadVideoFile = async (
    file: File,
    metadata: {
      track_id: string | number
      title?: string
      description?: string
      track_start_ms?: number
      track_end_ms?: number
    },
    config?: UseRequestConfig<UploadEditResponse>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const formData = new FormData()
    formData.append('video', file)
    formData.append('track_id', String(metadata.track_id))
    if (metadata.title) formData.append('title', metadata.title)
    if (metadata.description) formData.append('description', metadata.description)
    if (metadata.track_start_ms !== undefined) formData.append('track_start_ms', String(metadata.track_start_ms))
    if (metadata.track_end_ms !== undefined) formData.append('track_end_ms', String(metadata.track_end_ms))

    return useRequest<UploadEditResponse>(
      VideoApiRoutes.UPLOAD_EDIT,
      {
        method: 'POST',
        data: formData,
        ...requestConfig,
      },
      { schema: UploadEditResponseSchema, silent: false, ...config },
    )
  }

  const getMusicStatus = async (
    userId: string,
    config?: UseRequestConfig<MusicStatusResponse>,
  ) => {
    return useRequest<MusicStatusResponse>(
      `/users/${userId}/music-status`,
      { method: 'GET' },
      config,
    )
  }

  const updateMusicStatus = async (
    payload: { track_id: string },
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      VideoApiRoutes.MUSIC_STATUS,
      { method: 'PUT', data: payload },
      { silent: true, ...config },
    )
  }

  const clearMusicStatus = async (config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      VideoApiRoutes.MUSIC_STATUS,
      { method: 'DELETE' },
      { silent: true, ...config },
    )
  }

  const setLikeVisibility = async (
    trackId: string,
    payload: LikeVisibilityPayload,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      VideoApiRoutes.LIKE_VISIBILITY.replace(':id', trackId),
      { method: 'PUT', data: payload },
      { silent: false, ...config },
    )
  }

  const getUserVideos = async (
    userId: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<ExploreResponse>,
  ) => {
    return useRequest<ExploreResponse>(
      VideoApiRoutes.USER_VIDEOS.replace(':userId', userId),
      { method: 'GET', params },
      config,
    )
  }

  const getLikedTracksPublic = async (
    userId: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: LikedTrackItem[] }>,
  ) => {
    return useRequest<{ items: LikedTrackItem[] }>(
      VideoApiRoutes.LIKED_TRACKS.replace(':userId', userId),
      { method: 'GET', params },
      config,
    )
  }

  return {
    // Admin
    adminGetVideos,
    adminUpdateVideo,
    adminDeleteVideo,
    adminApproveVideo,
    adminGetVideoComments,
    adminUpdateComment,
    adminDeleteComment,
    // Public / user
    getExploreVideos,
    getTrackVideos,
    getVideoJobStatus,
    likeVideo,
    unlikeVideo,
    viewVideo,
    uploadEdit,
    uploadVideoFile,
    getMusicStatus,
    updateMusicStatus,
    clearMusicStatus,
    setLikeVisibility,
    getUserVideos,
    getLikedTracksPublic,
  }
}
