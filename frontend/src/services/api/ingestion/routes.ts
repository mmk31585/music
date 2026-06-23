import { z } from 'zod'
import { useRequest } from '@/composables/useRequest'
import type { AxiosRequestConfig } from 'axios'
import type { UseRequestConfig } from '@/plugins/client/types'
import { IngestionApiRoutes } from './enums'
import {
  DraftListItemSchema,
  UploadResponseSchema,
  DraftDetailResponseSchema,
  ListDraftsResponseSchema,
  EnrichmentResultSchema,
  ArtistSearchResultSchema,
  AlbumSearchResultSchema,
  FinalizeResultSchema,
  IngestionStatsSchema,
  IngestionConfigSchema,
  type UploadResponse,
  type DraftDetailResponse,
  type ListDraftsResponse,
  type EnrichmentResult,
  type SaveFinalMetadataRequest,
  type ArtistSearchResult,
  type AlbumSearchResult,
  type FinalizeResult,
  type IngestionStats,
  type IngestionConfig,
} from './types'

export const useIngestionApi = () => {
  const uploadAudio = async (
    file: File,
    config?: UseRequestConfig<UploadResponse>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const formData = new FormData()
    formData.append('audio', file)

    return useRequest<UploadResponse>(
      IngestionApiRoutes.ADMIN_UPLOAD,
      {
        method: 'POST',
        data: formData,
        ...requestConfig,
      },
      {
        schema: UploadResponseSchema,
        silent: false,
        ...config,
      },
    )
  }

  const listDrafts = async (
    params?: { status?: string; page?: number; limit?: number },
    config?: UseRequestConfig<ListDraftsResponse>,
  ) => {
    const query = new URLSearchParams()
    if (params?.status) query.set('status', params.status)
    if (params?.page) query.set('page', String(params.page))
    if (params?.limit) query.set('limit', String(params.limit))
    const qs = query.toString()

    return useRequest<ListDraftsResponse>(
      `${IngestionApiRoutes.ADMIN_DRAFTS}${qs ? `?${qs}` : ''}`,
      { method: 'GET' },
      {
        // Use per-item schema so request-factory's pagination handler validates correctly
        schema: DraftListItemSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getDraftDetail = async (
    id: string,
    config?: UseRequestConfig<DraftDetailResponse>,
  ) => {
    return useRequest<DraftDetailResponse>(
      `${IngestionApiRoutes.ADMIN_DRAFTS}/${id}`,
      { method: 'GET' },
      {
        schema: DraftDetailResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const enrichDraft = async (
    id: string,
    scope?: string[],
    config?: UseRequestConfig<any>,
  ) => {
    const body = scope && scope.length > 0 ? { scope } : undefined
    return useRequest<any>(
      IngestionApiRoutes.ADMIN_ENRICH.replace(':id', id),
      { method: 'POST', data: body },
      { silent: false, ...config },
    )
  }

  const getDraftSuggestions = async (
    id: string,
    config?: UseRequestConfig<EnrichmentResult>,
  ) => {
    return useRequest<EnrichmentResult>(
      IngestionApiRoutes.ADMIN_SUGGESTIONS.replace(':id', id),
      { method: 'GET' },
      {
        schema: EnrichmentResultSchema,
        silent: true,
        ...config,
      },
    )
  }

  const updateDraftMetadata = async (
    id: string,
    data: { title?: string; artist?: string; album?: string },
    config?: UseRequestConfig<any>,
  ) => {
    return useRequest<any>(
      IngestionApiRoutes.ADMIN_UPDATE_METADATA.replace(':id', id),
      { method: 'PATCH', data },
      { silent: false, ...config },
    )
  }

  const saveFinalMetadata = async (
    id: string,
    data: SaveFinalMetadataRequest,
    config?: UseRequestConfig<any>,
  ) => {
    return useRequest<any>(
      IngestionApiRoutes.ADMIN_FINAL_METADATA.replace(':id', id),
      { method: 'PATCH', data },
      { silent: false, ...config },
    )
  }

  const rejectDraft = async (
    id: string,
    reason?: string,
    config?: UseRequestConfig<any>,
  ) => {
    return useRequest<any>(
      `${IngestionApiRoutes.ADMIN_REJECT.replace(':id', id)}${reason ? `?reason=${encodeURIComponent(reason)}` : ''}`,
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const searchArtists = async (
    q: string,
    config?: UseRequestConfig<ArtistSearchResult[]>,
  ) => {
    return useRequest<ArtistSearchResult[]>(
      `${IngestionApiRoutes.ADMIN_ARTISTS_SEARCH}?q=${encodeURIComponent(q)}`,
      { method: 'GET' },
      { schema: z.array(ArtistSearchResultSchema), silent: true, ...config },
    )
  }

  const searchAlbums = async (
    q: string,
    config?: UseRequestConfig<AlbumSearchResult[]>,
  ) => {
    return useRequest<AlbumSearchResult[]>(
      `${IngestionApiRoutes.ADMIN_ALBUMS_SEARCH}?q=${encodeURIComponent(q)}`,
      { method: 'GET' },
      { schema: z.array(AlbumSearchResultSchema), silent: true, ...config },
    )
  }

  const finalizeDraft = async (
    id: string,
    config?: UseRequestConfig<FinalizeResult>,
  ) => {
    return useRequest<FinalizeResult>(
      IngestionApiRoutes.ADMIN_FINALIZE.replace(':id', id),
      { method: 'POST' },
      { schema: FinalizeResultSchema, silent: false, ...config },
    )
  }

  const uploadDraftImage = async (
    draftId: string,
    entity: string,
    file: File,
    config?: UseRequestConfig<{ url: string }>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const formData = new FormData()
    formData.append('image', file)

    return useRequest<{ url: string }>(
      IngestionApiRoutes.ADMIN_UPLOAD_IMAGE.replace(':id', draftId).replace(':entity', entity),
      {
        method: 'POST',
        data: formData,
        ...requestConfig,
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  const getIngestionStats = async (
    config?: UseRequestConfig<IngestionStats>,
  ) => {
    return useRequest<IngestionStats>(
      IngestionApiRoutes.ADMIN_STATS,
      { method: 'GET' },
      { schema: IngestionStatsSchema, silent: true, ...config },
    )
  }

  const getIngestionConfig = async (
    config?: UseRequestConfig<IngestionConfig>,
  ) => {
    return useRequest<IngestionConfig>(
      IngestionApiRoutes.ADMIN_CONFIG,
      { method: 'GET' },
      { schema: IngestionConfigSchema, silent: true, ...config },
    )
  }

  const triggerCleanup = async (
    config?: UseRequestConfig<any>,
  ) => {
    return useRequest<any>(
      IngestionApiRoutes.ADMIN_CLEANUP,
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  return {
    uploadAudio,
    listDrafts,
    getDraftDetail,
    enrichDraft,
    getDraftSuggestions,
    saveFinalMetadata,
    updateDraftMetadata,
    rejectDraft,
    searchArtists,
    searchAlbums,
    finalizeDraft,
    getIngestionStats,
    getIngestionConfig,
    triggerCleanup,
    uploadDraftImage,
  }
}
