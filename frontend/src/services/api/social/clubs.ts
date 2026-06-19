import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { SocialApiRoutes } from './enums'
import type {
  MusicClub,
  MusicClubMember,
  MusicClubPost,
  ClubDetailResponse,
  ListeningParty,
} from './types'

export interface CreateClubPayload {
  name: string
  slug?: string
  description?: string
  genre?: string
  cover_url?: string
  is_public?: boolean
  max_members?: number
}

export interface LaunchPartyPayload {
  title?: string
  description?: string
  is_public?: boolean
}

export const createClub = async (
  data: CreateClubPayload,
  config?: UseRequestConfig<MusicClub>,
) => {
  return useRequest<MusicClub>(SocialApiRoutes.CLUBS, { method: 'POST', data }, config)
}

export const listClubs = async (
  params?: { limit?: number; offset?: number; genre?: string },
  config?: UseRequestConfig<MusicClub[]>,
) => {
  return useRequest<MusicClub[]>(
    SocialApiRoutes.CLUB_BROWSE,
    { method: 'GET', params },
    { silent: true, ...config },
  )
}

export const getClubDetail = async (id: string, config?: UseRequestConfig<ClubDetailResponse>) => {
  return useRequest<ClubDetailResponse>(
    SocialApiRoutes.CLUB_DETAIL.replace(':id', id),
    { method: 'GET' },
    { silent: true, ...config },
  )
}

export const joinClub = async (id: string, config?: UseRequestConfig<void>) => {
  return useRequest<void>(SocialApiRoutes.CLUB_JOIN.replace(':id', id), { method: 'POST' }, config)
}

export const leaveClub = async (id: string, config?: UseRequestConfig<void>) => {
  return useRequest<void>(
    SocialApiRoutes.CLUB_LEAVE.replace(':id', id),
    { method: 'DELETE' },
    config,
  )
}

export const launchParty = async (
  id: string,
  data?: LaunchPartyPayload,
  config?: UseRequestConfig<ListeningParty>,
) => {
  return useRequest<ListeningParty>(
    SocialApiRoutes.CLUB_LAUNCH_PARTY.replace(':id', id),
    { method: 'POST', data },
    config,
  )
}

export const getClubMembers = async (id: string, config?: UseRequestConfig<MusicClubMember[]>) => {
  return useRequest<MusicClubMember[]>(
    SocialApiRoutes.CLUB_MEMBERS.replace(':id', id),
    { method: 'GET' },
    { silent: true, ...config },
  )
}

export const getClubPosts = async (
  id: string,
  params?: { limit?: number; offset?: number },
  config?: UseRequestConfig<MusicClubPost[]>,
) => {
  return useRequest<MusicClubPost[]>(
    SocialApiRoutes.CLUB_POSTS.replace(':id', id),
    { method: 'GET', params },
    { silent: true, ...config },
  )
}
