export interface FeatureFlags {
  analytics: boolean
  recommendation: boolean
  search: boolean
  social: boolean
  reactions: boolean
  creator: boolean
  moderation: boolean
  ai: boolean
  contribution: boolean
  gamification: boolean
  tips: boolean
  subscription: boolean
  notification: boolean
  redesignedPlayer: boolean
}

export type FeatureFlagKey = keyof FeatureFlags
