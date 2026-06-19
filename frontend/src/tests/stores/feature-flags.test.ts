import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import type { FeatureFlags } from '@/services/api/feature-flags'
import { useFeatureFlagsStore } from '@/stores/feature-flags'

describe('useFeatureFlagsStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with null flags and not fetched', () => {
    const store = useFeatureFlagsStore()
    expect(store.flags).toBeNull()
    expect(store.fetched).toBe(false)
    expect(store.loading).toBe(false)
    expect(store.isLoaded).toBe(false)
  })

  it('isEnabled returns true by default when flags not loaded', () => {
    const store = useFeatureFlagsStore()
    expect(store.isEnabled('social')).toBe(true)
    expect(store.isEnabled('ai')).toBe(true)
  })

  it('isEnabled returns false after setFlags with disabled feature', () => {
    const store = useFeatureFlagsStore()
    store.setFlags({
      analytics: true,
      recommendation: true,
      search: true,
      social: false,
      reactions: true,
      creator: true,
      moderation: true,
      ai: true,
      contribution: true,
      gamification: true,
      tips: true,
      subscription: true,
      notification: true,
      redesignedPlayer: true,
    })
    expect(store.isEnabled('social')).toBe(false)
    expect(store.isEnabled('analytics')).toBe(true)
    expect(store.fetched).toBe(true)
    expect(store.isLoaded).toBe(true)
  })

  it('$reset clears all state', () => {
    const store = useFeatureFlagsStore()
    store.setFlags({
      analytics: false, recommendation: false, search: false,
      social: false, reactions: false, creator: false,
      moderation: false, ai: false, contribution: false,
      gamification: false, tips: false, subscription: false,
      notification: false, redesignedPlayer: false,
    })
    expect(store.fetched).toBe(true)
    store.$reset()
    expect(store.flags).toBeNull()
    expect(store.fetched).toBe(false)
    expect(store.loading).toBe(false)
  })

  it('handles all feature flag keys', () => {
    const keys = [
      'analytics', 'recommendation', 'search', 'social',
      'reactions', 'creator', 'moderation', 'ai',
      'contribution', 'gamification', 'tips', 'subscription',
      'notification', 'redesignedPlayer',
    ] as const
    const store = useFeatureFlagsStore()
    const enabled: Record<string, boolean> = {}
    keys.forEach((k) => { enabled[k] = true })
    enabled.social = false
    enabled.ai = false
    store.setFlags(enabled as any as FeatureFlags)
    keys.forEach((k) => {
      if (k === 'social' || k === 'ai') {
        expect(store.isEnabled(k)).toBe(false)
      } else {
        expect(store.isEnabled(k)).toBe(true)
      }
    })
  })
})
