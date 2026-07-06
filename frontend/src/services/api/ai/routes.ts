import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { AIApiRoutes } from './enums'
import {
  MoodResponseSchema,
  AIPlaylistResponseSchema,
  EmbeddingResponseSchema,
  type MoodResponse,
  type AIPlaylistResponse,
  type EmbeddingResponse,
  type GeneratePlaylistPayload,
  type AnalyzeMoodPayload,
  type EmbeddingPayload,
  type TrackMeta,
} from './types'

export const useAIApi = () => {
  const generatePlaylist = async (
    payload: GeneratePlaylistPayload,
    config?: UseRequestConfig<AIPlaylistResponse>,
  ) => {
    return useRequest<AIPlaylistResponse>(
      AIApiRoutes.GENERATE_PLAYLIST,
      { method: 'POST', data: payload },
      { schema: AIPlaylistResponseSchema, silent: false, ...config },
    )
  }

  const analyzeMood = async (
    payload: AnalyzeMoodPayload,
    config?: UseRequestConfig<MoodResponse>,
  ) => {
    return useRequest<MoodResponse>(
      AIApiRoutes.ANALYZE_MOOD,
      { method: 'POST', data: payload },
      { schema: MoodResponseSchema, silent: false, ...config },
    )
  }

  const getMood = async (trackId: string, config?: UseRequestConfig<MoodResponse>) => {
    return useRequest<MoodResponse>(
      AIApiRoutes.GET_MOOD.replace(':trackId', trackId),
      { method: 'GET' },
      { schema: MoodResponseSchema, silent: true, ...config },
    )
  }

  const generateEmbedding = async (payload: EmbeddingPayload, config?: UseRequestConfig<EmbeddingResponse[]>) => {
    return useRequest<EmbeddingResponse[]>(
      AIApiRoutes.GENERATE_EMBEDDING,
      { method: 'POST', data: payload },
      { silent: false, schema: EmbeddingResponseSchema, ...config },
    )
  }

  const similarByMood = async (
    trackId: string,
    mood?: string,
    config?: UseRequestConfig<TrackMeta[]>,
  ) => {
    const params: Record<string, string> = {}
    if (mood) params.mood = mood
    return useRequest<TrackMeta[]>(
      AIApiRoutes.SIMILAR_BY_MOOD.replace(':trackId', trackId),
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const similarByEmbedding = async (trackId: string, config?: UseRequestConfig<TrackMeta[]>) => {
    return useRequest<TrackMeta[]>(
      AIApiRoutes.SIMILAR_BY_EMBEDDING.replace(':trackId', trackId),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  return {
    generatePlaylist,
    analyzeMood,
    getMood,
    generateEmbedding,
    similarByMood,
    similarByEmbedding,
  }
}
