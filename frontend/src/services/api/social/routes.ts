import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { SocialApiRoutes } from './enums'
import {
  FollowersResponseSchema,
  ActivityFeedResponseSchema,
  type FollowersResponse,
  type ActivityFeedResponse,
  type ListeningParty,
  type LiveRoom,
  type LiveRoomParticipant,
  type LiveRoomQueueItem,
  type MusicClub,
  type MusicClubMember,
  type MusicClubPost,
  type ClubDetailResponse,
  type Discussion,
  type ClubDiscussion,
  type ClubDiscussionReply,
  type TrackRating,
  type CreatePartyRequest,
  type CreateRoomRequest,
  type CreateClubRequest,
  type CreateDiscussionRequest,
  type CreateClubDiscussionRequest,
  type CreateRatingRequest,
  type LaunchPartyPayload,
} from './types'

export const useSocialApi = () => {
  const follow = async (userId: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.FOLLOW.replace(':userId', userId),
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const unfollow = async (userId: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.UNFOLLOW.replace(':userId', userId),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  const getFollowers = async (
    userId: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<FollowersResponse>,
  ) => {
    return useRequest<FollowersResponse>(
      SocialApiRoutes.FOLLOWERS.replace(':userId', userId),
      { method: 'GET', params },
      { schema: FollowersResponseSchema, silent: true, ...config },
    )
  }

  const getFollowing = async (
    userId: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<FollowersResponse>,
  ) => {
    return useRequest<FollowersResponse>(
      SocialApiRoutes.FOLLOWING.replace(':userId', userId),
      { method: 'GET', params },
      { schema: FollowersResponseSchema, silent: true, ...config },
    )
  }

  const isFollowing = async (
    userId: string,
    config?: UseRequestConfig<{ is_following: boolean }>,
  ) => {
    return useRequest<{ is_following: boolean }>(
      SocialApiRoutes.IS_FOLLOWING.replace(':userId', userId),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getFeed = async (
    params?: { user_id?: string; limit?: number; offset?: number; types?: string },
    config?: UseRequestConfig<ActivityFeedResponse>,
  ) => {
    return useRequest<ActivityFeedResponse>(
      SocialApiRoutes.FEED,
      { method: 'GET', params },
      { schema: ActivityFeedResponseSchema, silent: true, ...config },
    )
  }

  // --- Listening Parties ---

  const listParties = async (
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<ListeningParty[]>,
  ) => {
    return useRequest<ListeningParty[]>(
      SocialApiRoutes.PARTIES,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const createParty = async (
    data: CreatePartyRequest,
    config?: UseRequestConfig<ListeningParty>,
  ) => {
    return useRequest<ListeningParty>(SocialApiRoutes.PARTIES, { method: 'POST', data }, config)
  }

  const getParty = async (id: string, config?: UseRequestConfig<ListeningParty>) => {
    return useRequest<ListeningParty>(
      SocialApiRoutes.PARTY.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const updatePartyStatus = async (id: string, status: string, trackId?: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.PARTY_STATUS.replace(':id', id),
      { method: 'PUT', data: trackId ? { status, track_id: trackId } : { status } },
      config,
    )
  }

  const joinParty = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.PARTY_JOIN.replace(':id', id),
      { method: 'POST' },
      config,
    )
  }

  const leaveParty = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.PARTY_LEAVE.replace(':id', id),
      { method: 'POST' },
      config,
    )
  }

  // --- Live Rooms ---

  const listRooms = async (
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<LiveRoom[]>,
  ) => {
    return useRequest<LiveRoom[]>(
      SocialApiRoutes.ROOMS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const createRoom = async (data: CreateRoomRequest, config?: UseRequestConfig<LiveRoom>) => {
    return useRequest<LiveRoom>(SocialApiRoutes.ROOMS, { method: 'POST', data }, config)
  }

  const getRoom = async (id: string, config?: UseRequestConfig<LiveRoom>) => {
    return useRequest<LiveRoom>(
      SocialApiRoutes.ROOM.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const joinRoom = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.ROOM_JOIN.replace(':id', id),
      { method: 'POST' },
      config,
    )
  }

  const leaveRoom = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.ROOM_LEAVE.replace(':id', id),
      { method: 'POST' },
      config,
    )
  }

  const getRoomParticipants = async (
    id: string,
    config?: UseRequestConfig<LiveRoomParticipant[]>,
  ) => {
    return useRequest<LiveRoomParticipant[]>(
      SocialApiRoutes.ROOM_PARTICIPANTS.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const addToRoomQueue = async (id: string, trackId: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.ROOM_QUEUE.replace(':id', id),
      { method: 'POST', data: { track_id: trackId } },
      config,
    )
  }

  const getRoomQueue = async (id: string, config?: UseRequestConfig<LiveRoomQueueItem[]>) => {
    return useRequest<LiveRoomQueueItem[]>(
      SocialApiRoutes.ROOM_QUEUE.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  // --- Music Clubs (Phase 5) ---

  const listClubs = async (
    params?: { limit?: number; offset?: number; genre?: string },
    config?: UseRequestConfig<MusicClub[]>,
  ) => {
    return useRequest<MusicClub[]>(
      SocialApiRoutes.CLUBS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const listClubsWithGenre = async (
    params?: { genre?: string; limit?: number; offset?: number },
    config?: UseRequestConfig<MusicClub[]>,
  ) => {
    return useRequest<MusicClub[]>(
      SocialApiRoutes.CLUB_BROWSE,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const createClub = async (data: CreateClubRequest, config?: UseRequestConfig<MusicClub>) => {
    return useRequest<MusicClub>(SocialApiRoutes.CLUBS, { method: 'POST', data }, config)
  }

  const getClub = async (id: string, config?: UseRequestConfig<MusicClub>) => {
    return useRequest<MusicClub>(
      SocialApiRoutes.CLUB.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getClubDetail = async (id: string, config?: UseRequestConfig<ClubDetailResponse>) => {
    return useRequest<ClubDetailResponse>(
      SocialApiRoutes.CLUB_DETAIL.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const joinClub = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.CLUB_JOIN.replace(':id', id),
      { method: 'POST' },
      config,
    )
  }

  const leaveClub = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.CLUB_LEAVE.replace(':id', id),
      { method: 'POST' },
      config,
    )
  }

  const launchParty = async (
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

  const getClubMembers = async (id: string, config?: UseRequestConfig<MusicClubMember[]>) => {
    return useRequest<MusicClubMember[]>(
      SocialApiRoutes.CLUB_MEMBERS.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getClubPosts = async (
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

  const createClubPost = async (
    id: string,
    content: string,
    config?: UseRequestConfig<MusicClubPost>,
  ) => {
    return useRequest<MusicClubPost>(
      SocialApiRoutes.CLUB_POSTS.replace(':id', id),
      { method: 'POST', data: { content } },
      config,
    )
  }

  // --- Discussions ---

  const getDiscussions = async (
    params: { target_type: string; target_id: string; limit?: number; offset?: number },
    config?: UseRequestConfig<Discussion[]>,
  ) => {
    return useRequest<Discussion[]>(
      SocialApiRoutes.DISCUSSIONS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const createDiscussion = async (
    data: CreateDiscussionRequest,
    config?: UseRequestConfig<Discussion>,
  ) => {
    return useRequest<Discussion>(SocialApiRoutes.DISCUSSIONS, { method: 'POST', data }, config)
  }

  const getDiscussionReplies = async (id: string, config?: UseRequestConfig<Discussion[]>) => {
    return useRequest<Discussion[]>(
      SocialApiRoutes.DISCUSSION_REPLIES.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  // --- Club Discussions (Phase 6) ---

  const listClubDiscussions = async (
    clubId: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<ClubDiscussion[]>,
  ) => {
    return useRequest<ClubDiscussion[]>(
      SocialApiRoutes.CLUB_DISCUSSIONS.replace(':clubId', clubId),
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const createClubDiscussion = async (
    clubId: string,
    data: CreateClubDiscussionRequest,
    config?: UseRequestConfig<ClubDiscussion>,
  ) => {
    return useRequest<ClubDiscussion>(
      SocialApiRoutes.CLUB_DISCUSSIONS.replace(':clubId', clubId),
      { method: 'POST', data },
      config,
    )
  }

  const getClubDiscussion = async (id: string, config?: UseRequestConfig<ClubDiscussion>) => {
    return useRequest<ClubDiscussion>(
      SocialApiRoutes.CLUB_DISCUSSION.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getClubDiscussionReplies = async (
    id: string,
    config?: UseRequestConfig<ClubDiscussionReply[]>,
  ) => {
    return useRequest<ClubDiscussionReply[]>(
      SocialApiRoutes.CLUB_DISCUSSION_REPLIES.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const createClubDiscussionReply = async (
    discussionId: string,
    data: { body: string },
    config?: UseRequestConfig<ClubDiscussionReply>,
  ) => {
    return useRequest<ClubDiscussionReply>(
      SocialApiRoutes.CLUB_DISCUSSION_REPLIES.replace(':id', discussionId),
      { method: 'POST', data },
      config,
    )
  }

  const deleteClubDiscussion = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      SocialApiRoutes.CLUB_DISCUSSION.replace(':id', id),
      { method: 'DELETE' },
      config,
    )
  }

  const deleteClubDiscussionReply = async (
    discussionId: string,
    replyId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      SocialApiRoutes.CLUB_DISCUSSION_REPLY_DELETE.replace(':id', discussionId).replace(':replyId', replyId),
      { method: 'DELETE' },
      config,
    )
  }

  // --- Track Ratings ---

  const createRating = async (
    data: CreateRatingRequest,
    config?: UseRequestConfig<TrackRating>,
  ) => {
    return useRequest<TrackRating>(SocialApiRoutes.RATINGS, { method: 'POST', data }, config)
  }

  const getTrackRatings = async (trackId: string, config?: UseRequestConfig<TrackRating[]>) => {
    return useRequest<TrackRating[]>(
      SocialApiRoutes.TRACK_RATINGS.replace(':trackId', trackId),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  return {
    follow,
    unfollow,
    getFollowers,
    getFollowing,
    isFollowing,
    getFeed,
    listParties,
    createParty,
    getParty,
    updatePartyStatus,
    joinParty,
    leaveParty,
    listRooms,
    createRoom,
    getRoom,
    joinRoom,
    leaveRoom,
    getRoomParticipants,
    addToRoomQueue,
    getRoomQueue,
    listClubs,
    listClubsWithGenre,
    createClub,
    getClub,
    getClubDetail,
    joinClub,
    leaveClub,
    launchParty,
    getClubMembers,
    getClubPosts,
    createClubPost,
    getDiscussions,
    createDiscussion,
    getDiscussionReplies,
    listClubDiscussions,
    createClubDiscussion,
    getClubDiscussion,
    getClubDiscussionReplies,
    createClubDiscussionReply,
    deleteClubDiscussion,
    deleteClubDiscussionReply,
    createRating,
    getTrackRatings,
  }
}
