import { useRequest } from '@/composables/useRequest'

export interface TrackSummary {
  id: string
  title: string
  artist_name: string
  album_title?: string | null
  cover_url?: string | null
  duration_seconds?: number | null
}

export interface QueueCandidate {
  id: string
  room_id: string
  track: TrackSummary
  suggested_by: { id: string; username: string; avatar_url?: string }
  vote_count: number
  has_voted: boolean
  created_at: string
}

export interface RoomNowPlaying {
  track: TrackSummary
  started_at: string
  suggested_by?: { id: string; username: string }
  source: 'vote' | 'autofill'
}

export interface QueueState {
  now_playing: RoomNowPlaying | null
  candidates: QueueCandidate[]
}

export async function getQueueState(roomId: string): Promise<QueueState> {
  return useRequest<QueueState>(
    `/social/rooms/${roomId}/queue/state`,
    { method: 'GET' },
    { silent: true },
  )
}

export async function suggestTrack(roomId: string, trackId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/queue/suggest`,
    { method: 'POST', data: { track_id: trackId } },
  )
}

export async function castVote(roomId: string, candidateId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/queue/${candidateId}/vote`,
    { method: 'POST' },
  )
}

export async function removeVote(roomId: string, candidateId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/queue/${candidateId}/vote`,
    { method: 'DELETE' },
  )
}

export async function reportTrackEnded(roomId: string, trackId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/track-ended`,
    { method: 'POST', data: { track_id: trackId } },
  )
}
