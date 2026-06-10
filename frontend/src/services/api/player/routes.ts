import { useRequest } from '@/composables/useRequest'
import { apiReplaceParams } from '@/utils/api-replace-params'
import { PlayerApiRoutes } from './enums'
import { PlaybackTrackSchema, type PlaybackTrack } from './types'

function getApiOrigin() {
  return String(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1').replace(/\/$/, '')
}

function buildStreamUrl(path: string) {
  const origin = getApiOrigin()
console.log(origin)
console.log(path)

  return `${origin}${path}`
}

export const usePlayerApi = () => {
  const getPlaybackTrack = async (id: string) => {
    return useRequest<PlaybackTrack>(
      apiReplaceParams(PlayerApiRoutes.GET_PLAYBACK_TRACK, { id }),
      {
        method: 'GET',
      },
      {
        schema: PlaybackTrackSchema,
      },
    )
  }

  const getTrackStreamUrl = (id: string) => {
    return buildStreamUrl(apiReplaceParams(PlayerApiRoutes.STREAM_TRACK, { id }))
  }

  return {
    getPlaybackTrack,
    getTrackStreamUrl,
  }
}

export type { PlaybackTrack }
