export interface GamificationProfile {
  user_id: string
  username: string
  avatar_url: string
  level: number
  current_xp: number
  next_level_xp: number
  total_xp: number
  title: string
  title_persian: string
  rank: number
  badges: UserBadge[]
  challenges: UserChallenge[]
}

export interface Badge {
  id: string
  name: string
  description: string
  icon_url: string
  category: string
  rarity: string
  criteria: string
  xp_reward: number
  is_hidden: boolean
}

export interface UserBadge {
  user_id: string
  badge_id: string
  earned_at: string
  is_displayed: boolean
  badge?: Badge
}

export interface DailyChallenge {
  id: string
  title: string
  description: string
  challenge_type: string
  target_count: number
  xp_reward: number
  badge_reward_id: string | null
  is_active: boolean
  valid_from: string
  valid_until: string
}

export interface UserChallenge {
  user_id: string
  challenge_id: string
  progress: number
  is_completed: boolean
  completed_at: string | null
  challenge?: DailyChallenge
}

export interface LeaderboardEntry {
  rank: number
  user_id: string
  username: string
  avatar_url: string
  score: number
}

export interface LeaderboardResponse {
  type: string
  entries: LeaderboardEntry[]
}

export interface BadgesResponse {
  all: Badge[]
  mine: UserBadge[]
}

export interface ChallengesResponse {
  challenges: DailyChallenge[]
  progress: UserChallenge[]
  earned_xp: number
}
