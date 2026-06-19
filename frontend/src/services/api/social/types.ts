import { z } from 'zod'

export const UserFollowSchema = z.object({
  follower_id: z.string(),
  followed_id: z.string(),
  created_at: z.string(),
})

export const ActivityFeedItemSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  type: z.string(),
  target_id: z.string(),
  target_type: z.string(),
  metadata: z.record(z.string(), z.unknown()).nullable().optional(),
  created_at: z.string(),
  user_display_name: z.string(),
  user_avatar_url: z.string().nullable().optional(),
  target_name: z.string().nullable().optional(),
  target_image_url: z.string().nullable().optional(),
})

export const FollowersResponseSchema = z.object({
  items: z.array(UserFollowSchema).catch([]),
  total_count: z.number(),
  limit: z.number(),
  offset: z.number(),
})

export const ActivityFeedResponseSchema = z.object({
  items: z.array(ActivityFeedItemSchema).catch([]),
  pagination: z.object({
    limit: z.number(),
    offset: z.number(),
    count: z.number(),
    has_more: z.boolean(),
  }),
})

// --- Social Feature Types (no Zod - simpler) ---

export interface ListeningParty {
  id: string
  host_id: string
  title: string
  description?: string
  cover_url?: string
  is_public: boolean
  status: 'active' | 'paused' | 'ended'
  current_track_id?: string
  current_position_ms: number
  started_at: string
  ended_at?: string
  created_at: string
  participant_count: number
}

export interface LiveRoom {
  id: string
  host_id: string
  title: string
  description?: string
  cover_url?: string
  is_public: boolean
  status: 'live' | 'ended'
  current_track_id?: string
  listener_count: number
  created_at: string
  ended_at?: string
}

export interface LiveRoomParticipant {
  id: string
  room_id: string
  user_id: string
  role: 'host' | 'co_host' | 'speaker' | 'listener'
  joined_at: string
  left_at?: string
  is_active: boolean
}

export interface LiveRoomQueueItem {
  id: string
  room_id: string
  track_id: string
  added_by: string
  position: number
  status: string
  played_at?: string
  created_at: string
}

export interface MusicClub {
  id: string
  name: string
  slug: string
  description?: string
  cover_url?: string
  genre?: string
  playlist_id?: string
  created_by: string
  is_public: boolean
  max_members: number
  member_count: number
  created_at: string
  updated_at: string
}

export interface MusicClubMember {
  id: string
  club_id: string
  user_id: string
  role: 'admin' | 'moderator' | 'member'
  joined_at: string
}

export interface MusicClubPost {
  id: string
  club_id: string
  user_id: string
  content: string
  created_at: string
  updated_at: string
}

export interface ClubDetailResponse {
  club: MusicClub
  is_member: boolean
  member_role: string
  members: MusicClubMember[]
  posts: MusicClubPost[]
  post_count: number
  track_count: number
}

export interface ClubMember {
  id: string
  club_id: string
  user_id: string
  role: 'admin' | 'moderator' | 'member'
  joined_at: string
}

export interface LaunchPartyPayload {
  title?: string
  description?: string
  is_public?: boolean
}

export interface Discussion {
  id: string
  user_id: string
  target_type: string
  target_id: string
  content: string
  parent_id?: string
  created_at: string
  updated_at: string
}

// Club Discussions (Phase 6)
export interface ClubDiscussion {
  id: string
  club_id: string
  author_id: string
  title: string
  body: string
  reply_count: number
  created_at: string
  updated_at: string
}

export interface ClubDiscussionReply {
  id: string
  discussion_id: string
  author_id: string
  body: string
  created_at: string
}

export interface CreateClubDiscussionRequest {
  title: string
  body: string
}

export interface CreateDiscussionReplyRequest {
  body: string
}

export interface TrackRating {
  id: string
  user_id: string
  track_id: string
  rating: number
  review?: string
  created_at: string
  updated_at: string
}

// Request types
export interface CreatePartyRequest {
  title: string
  description?: string
  is_public?: boolean
  track_id?: string
}

export interface CreateRoomRequest {
  title: string
  description?: string
  is_public?: boolean
}

export interface CreateClubRequest {
  name: string
  slug?: string
  description?: string
  genre?: string
  cover_url?: string
  is_public?: boolean
  max_members?: number
}

export interface CreateDiscussionRequest {
  target_type: string
  target_id: string
  content: string
  parent_id?: string
}

export interface CreateRatingRequest {
  track_id: string
  rating: number
  review?: string
}

export type UserFollow = z.infer<typeof UserFollowSchema>
export type ActivityFeedItem = z.infer<typeof ActivityFeedItemSchema>
export type FollowersResponse = z.infer<typeof FollowersResponseSchema>
export type ActivityFeedResponse = z.infer<typeof ActivityFeedResponseSchema>
