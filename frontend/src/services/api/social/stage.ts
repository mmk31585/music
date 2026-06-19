import { useRequest } from '@/composables/useRequest'

export interface StageMember {
  room_id: string
  user_id: string
  role: 'host' | 'speaker' | 'listener'
  joined_at: string
  muted: boolean
}

export interface HandRaiseData {
  id: string
  room_id: string
  user_id: string
  status: string
  created_at: string
}

export interface StageStateResponse {
  host: StageMember
  speakers: StageMember[]
  pending_requests: HandRaiseData[]
}

export interface StageState {
  host: StageMember & { username?: string; avatar_url?: string }
  speakers: (StageMember & { username?: string; avatar_url?: string })[]
  pending_hand_raises: (HandRaiseData & { username?: string; avatar_url?: string })[]
  my_role: 'host' | 'speaker' | 'listener'
  my_hand_raised: boolean
}

export async function getStageState(roomId: string): Promise<StageStateResponse> {
  return useRequest<StageStateResponse>(
    `/social/rooms/${roomId}/stage`,
    { method: 'GET' },
    { silent: true },
  )
}

export async function raiseHand(roomId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/raise-hand`,
    { method: 'POST' },
  )
}

export async function lowerHand(roomId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/lower-hand`,
    { method: 'POST' },
  )
}

export async function approveHand(roomId: string, userId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/${userId}/approve`,
    { method: 'POST' },
  )
}

export async function denyHand(roomId: string, userId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/${userId}/deny`,
    { method: 'POST' },
  )
}

export async function removeFromStage(roomId: string, userId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/${userId}`,
    { method: 'DELETE' },
  )
}

export async function leaveStage(roomId: string): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/leave`,
    { method: 'POST' },
  )
}

export async function toggleMute(roomId: string, userId: string, muted: boolean): Promise<void> {
  return useRequest<void>(
    `/social/rooms/${roomId}/stage/${userId}/mute`,
    { method: 'POST', data: { muted } },
  )
}
