import { useRequest } from '@/composables/useRequest'
import type {
  SearchResult,
  ImportResponse,
  ProgressResponse,
  ArtistSearchResponse,
  BatchImportItem,
  BatchImportResponse,
  BatchProgressResponse,
} from './types'

const BASE = '/admin/import'

export const useImportApi = () => {
  const search = async (q: string) => {
    return useRequest<SearchResult, true>(
      `${BASE}/search?q=${encodeURIComponent(q)}`,
      { method: 'GET' },
      { silent: true },
    )
  }

  const importTrack = async (track: SearchResult | { title: string; artist: string; source?: string; duration?: number; external_ids?: Record<string, string> }) => {
    const data: Record<string, string | number | Record<string, string>> = {}
    if ('url' in track && track.url) {
      data.url = track.url
    }
    if (track.source) {
      data.source = track.source
    }
    if (track.title) data.title = track.title
    if (track.artist) data.artist = track.artist
    if (track.duration) data.duration = track.duration
    if ('isrc' in track && track.isrc && typeof track.isrc === 'string') data.isrc = track.isrc
    if (track.external_ids && Object.keys(track.external_ids).length > 0) {
      data.external_ids = track.external_ids
    }

    return useRequest<ImportResponse>(
      `${BASE}/import`,
      { method: 'POST', data: data as Record<string, unknown> },
      { silent: false },
    )
  }

  const getProgress = async (jobId: string) => {
    return useRequest<ProgressResponse>(
      `${BASE}/${jobId}/progress`,
      { method: 'GET' },
      { silent: true },
    )
  }

  // ── Artist Discography Search ─────────────────────────────────

  const searchArtist = async (name: string) => {
    return useRequest<ArtistSearchResponse>(
      `${BASE}/artist?name=${encodeURIComponent(name)}`,
      { method: 'GET' },
      { silent: true },
    )
  }

  // ── Batch Import ─────────────────────────────────────────────

  const batchImport = async (tracks: BatchImportItem[]) => {
    return useRequest<BatchImportResponse>(
      `${BASE}/batch`,
      { method: 'POST', data: { tracks } },
      { silent: false },
    )
  }

  const getBatchProgress = async (batchId: string) => {
    return useRequest<BatchProgressResponse>(
      `${BASE}/batch/${batchId}/progress`,
      { method: 'GET' },
      { silent: true },
    )
  }

  return { search, importTrack, getProgress, searchArtist, batchImport, getBatchProgress }
}
