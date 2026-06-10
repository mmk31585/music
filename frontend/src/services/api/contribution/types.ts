export type ContributionType =
  | 'lyrics'
  | 'translation'
  | 'credits'
  | 'metadata'
  | 'album_art'
  | 'bio'
export type TargetType = 'track' | 'album' | 'artist'
export type ContributionStatus = 'pending' | 'approved' | 'rejected' | 'needs_review'
export type ReviewAction = 'approve' | 'reject'

export interface Contribution {
  id: string
  contributor_id: string
  contributor_username?: string | null
  contributor_avatar_url?: string | null
  contribution_type: ContributionType
  target_type: TargetType
  target_id: string
  locale: string | null
  data: any
  summary: string | null
  status: ContributionStatus
  ai_verdict: string | null
  ai_confidence: number | null
  version: number
  is_minor: boolean
  xp_awarded: number
  created_at: string
  decided_at: string | null
  applied_at: string | null
  moderator_note?: string | null
}

export interface ContributionHistoryItem {
  id: string
  data: any
  previous: any | null
  changed_by: string
  change_type: string
  created_at: string
}

export interface ContentVersion {
  id: string
  version: number
  data: any
  applied_by: string
  created_at: string
}

export interface ContributorStats {
  user_id: string
  username: string
  avatar_url: string
  total_contributions: number
  approved_count: number
  pending_count: number
  rejected_count: number
  xp_earned: number
  rank: number
}

export interface CreateContributionPayload {
  contribution_type: ContributionType
  target_type: TargetType
  target_id: string
  locale?: string
  data: any
  summary?: string
  is_minor?: boolean
}

export interface ReviewContributionPayload {
  action: ReviewAction
  note?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  meta: {
    page: number
    page_size: number
    total: number
  }
}
