import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { LyricsApiRoutes } from './enums'
import {
  LyricsSchema,
  type Lyrics,
  type CreateLyricsPayload,
  type UpdateLyricsPayload,
} from './types'

export const useLyricsApi = () => {
  const getTrackLyrics = async (
    trackId: string | number,
    lang?: string,
    config?: UseRequestConfig<Lyrics>,
  ) => {
    const query = lang ? `?lang=${encodeURIComponent(lang)}` : ''

    return useRequest<Lyrics>(
      `${LyricsApiRoutes.GET_TRACK_LYRICS.replace(':trackId', String(trackId))}${query}`,
      { method: 'GET' },
      {
        schema: LyricsSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getLyricsByTrackID = async (
    trackId: string | number,
    config?: UseRequestConfig<Lyrics[]>,
  ) => {
    return useRequest<Lyrics, true>(
      LyricsApiRoutes.GET_TRACK_LYRICS_ALL.replace(':trackId', String(trackId)),
      { method: 'GET' },
      {
        schema: LyricsSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateLyrics = async (
    payload: CreateLyricsPayload,
    config?: UseRequestConfig<Lyrics>,
  ) => {
    return useRequest<Lyrics>(
      LyricsApiRoutes.ADMIN_CREATE,
      {
        method: 'POST',
        data: payload,
      },
      {
        schema: LyricsSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUpdateLyrics = async (
    id: string | number,
    payload: UpdateLyricsPayload,
    config?: UseRequestConfig<Lyrics>,
  ) => {
    return useRequest<Lyrics>(
      LyricsApiRoutes.ADMIN_UPDATE.replace(':id', String(id)),
      {
        method: 'PUT',
        data: payload,
      },
      {
        schema: LyricsSchema,
        silent: false,
        ...config,
      },
    )
  }

  const fetchLrcLyrics = async (trackId: string | number, config?: UseRequestConfig<Lyrics>) => {
    return useRequest<Lyrics>(
      LyricsApiRoutes.ADMIN_FETCH_LRC.replace(':trackId', String(trackId)),
      { method: 'POST' },
      {
        schema: LyricsSchema,
        silent: false,
        ...config,
      },
    )
  }

  /**
   * Fetches lyrics using a pipeline:
   * 1. Tries LRCLIB (internet)
   * 2. If not found, falls back to AI generation via ML service
   *
   * Response sources:
   * - `lrclib` — lyrics were found on LRCLIB and saved
   * - `ai` — AI generation enqueued (poll track lyrics after ~30s)
   * - `none` — neither source had results
   */
  const fetchOrGenerateLyrics = async (
    trackId: string | number,
    config?: UseRequestConfig,
  ) => {
    return useRequest<{
      success: boolean
      source: 'lrclib' | 'openrouter' | 'ai' | 'none'
      data?: Lyrics
      job_id?: string
      track_id?: string
      message?: string
    }>(
      LyricsApiRoutes.ADMIN_FETCH_OR_GENERATE.replace(':trackId', String(trackId)),
      { method: 'POST' },
      {
        silent: false,
        ...config,
      },
    )
  }

  /**
   * Polls the AI lyrics generation status for a track.
   * Returns progress info while AI is working, or completed lyrics when done.
   */
  const aiStatus = async (trackId: string | number, config?: UseRequestConfig) => {
    return useRequest<{
      success: boolean
      status: 'queued' | 'downloading' | 'transcribing' | 'formatting' | 'callback_pending' | 'completed' | 'failed' | 'not_found' | 'unavailable'
      source?: string
      lyrics?: Lyrics
      job_id?: string
      track_id?: string
      message?: string
      progress_pct?: number
    }>(
      LyricsApiRoutes.ADMIN_AI_STATUS.replace(':trackId', String(trackId)),
      { method: 'GET' },
      {
        silent: true,
        ...config,
      },
    )
  }

  const adminDeleteLyrics = async (id: string | number, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      LyricsApiRoutes.ADMIN_DELETE.replace(':id', String(id)),
      { method: 'DELETE' },
      {
        silent: false,
        ...config,
      },
    )
  }

  /**
   * Syncs plain text lyrics to LRC format using OpenRouter AI.
   * Sends plain text, receives synced LRC lyrics.
   */
  const syncWithAI = async (
    payload: {
      track_id: string
      plain_text: string
      track_title?: string
      artist_name?: string
    },
    config?: UseRequestConfig,
  ) => {
    return useRequest<{
      success: boolean
      source: 'openrouter' | 'none'
      data?: Lyrics
      message?: string
    }>(
      LyricsApiRoutes.ADMIN_SYNC_AI,
      { method: 'POST', data: payload, timeout: 180000 }, // 3min for slow AI models
      {
        silent: false,
        ...config,
      },
    )
  }

  /**
   * Reviews and fixes existing LRC lyrics using OpenRouter AI.
   * Sends existing LRC, receives fixed LRC content.
   */
  const reviewWithAI = async (
    payload: {
      track_id: string
      existing_lrc: string
      track_title?: string
      artist_name?: string
    },
    config?: UseRequestConfig,
  ) => {
    return useRequest<{
      success: boolean
      source: 'openrouter' | 'none'
      language?: string
      data?: {
        content: string
        type: string
      }
      message?: string
    }>(
      LyricsApiRoutes.ADMIN_REVIEW_AI,
      { method: 'POST', data: payload, timeout: 180000 }, // 3min for slow AI models
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    getTrackLyrics,
    getLyricsByTrackID,
    adminCreateLyrics,
    fetchLrcLyrics,
    fetchOrGenerateLyrics,
    syncWithAI,
    reviewWithAI,
    aiStatus,
    adminUpdateLyrics,
    adminDeleteLyrics,
  }
}
