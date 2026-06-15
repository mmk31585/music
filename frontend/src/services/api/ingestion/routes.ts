import { z } from 'zod'
import { useRequest } from '@/composables/useRequest'
import type { AxiosRequestConfig } from 'axios'
import type { UseRequestConfig } from '@/plugins/client/types'
import { IngestionApiRoutes } from './enums'
import {
  UploadResponseSchema,
  DraftDetailResponseSchema,
  ListDraftsResponseSchema,
  DraftListItemSchema,
  EnrichmentResultSchema,
  SaveFinalMetadataRequestSchema,
  ArtistSearchResultSchema,
  AlbumSearchResultSchema,
  type UploadResponse,
  type DraftDetailResponse,
  type ListDraftsResponse,
  type DraftListItem,
  type EnrichmentResult,
  type SaveFinalMetadataRequest,
  type ArtistSearchResult,
  type AlbumSearchResult,
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
        schema: ListDraftsResponseSchema,
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
    config?: UseRequestConfig<any>,
  ) => {
    return useRequest<any>(
      IngestionApiRoutes.ADMIN_ENRICH.replace(':id', id),
      { method: 'POST' },
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

  return {
    uploadAudio,
    listDrafts,
    getDraftDetail,
    enrichDraft,
    getDraftSuggestions,
    saveFinalMetadata,
    rejectDraft,
    searchArtists,
    searchAlbums,
  }
}
