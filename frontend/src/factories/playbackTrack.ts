import type { PlaybackTrack } from '@/services/api/player/types'
import { apiReplaceParams } from '@/utils/api-replace-params'
import { PlayerApiRoutes } from '@/services/api/player/enums'

const STREAM_ORIGIN = String(
  import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
).replace(/\/$/, '')

function buildStreamUrl(path: string) {
  return `${STREAM_ORIGIN}${path}`
}

function getTrackStreamUrl(id: string) {
  return buildStreamUrl(apiReplaceParams(PlayerApiRoutes.STREAM_TRACK, { id }))
}

/**
 * Input shape accepted by the factory.
 * Supports the common API response patterns (snake_case, camelCase, mixed).
 */
export interface PlaybackTrackInput {
  id: string | number
  title?: string | null
  /** Common aliases: artist, artist_name, artistName */
  artist_name?: string | null
  artist?: string | null
  artistName?: string | null
  /** Common aliases: album_title, albumTitle */
  album_title?: string | null
  albumTitle?: string | null
  /** Common aliases: cover_url, coverUrl */
  cover_url?: string | null
  coverUrl?: string | null
  /** Common aliases: duration_seconds, durationSeconds, duration */
  duration_seconds?: number | null
  durationSeconds?: number | null
  duration?: number | null
  /** Common aliases: audio_url, streamUrl */
  audio_url?: string | null
  streamUrl?: string | null
  /** Unused fields are ignored (allows passing raw API data) */
  [key: string]: unknown
}

/**
 * Build a validated PlaybackTrack from various API response shapes.
 *
 * Normalizes snake_case / camelCase / mixed input into the PlaybackTrack schema.
 * Falls back to `getTrackStreamUrl(id)` when no stream URL is provided.
 */
export function buildPlaybackTrack(input: PlaybackTrackInput): PlaybackTrack {
  const id = String(input.id)

  const streamUrl =
    input.streamUrl ||
    input.audio_url ||
    getTrackStreamUrl(id)

  return {
    id,
    title: input.title ?? 'Unknown Track',
    artistName:
      input.artistName ||
      input.artist_name ||
      input.artist ||
      'Unknown artist',
    albumTitle:
      input.albumTitle ||
      input.album_title ||
      null,
    coverUrl:
      input.coverUrl ||
      input.cover_url ||
      null,
    durationSeconds:
      input.durationSeconds ??
      input.duration_seconds ??
      input.duration ??
      null,
    streamUrl,
  }
}

/**
 * Convenience wrapper: map an array of raw API items to PlaybackTrack[].
 */
export function mapToPlaybackTracks(items: PlaybackTrackInput[]): PlaybackTrack[] {
  return items.map(buildPlaybackTrack)
}
