import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

const mockList = vi.fn()

vi.mock('@/services/api/feature-flags', () => ({
  useFeatureFlagsApi: () => ({ list: mockList }),
  FeatureFlagsApiRoutes: { ALL: '/features' },
}))

describe('useFeatureFlags', () => {
  beforeEach(async () => {
    setActivePinia(createPinia())
    mockList.mockReset()
    vi.resetModules()
  })

  it('isEnabled returns true for unknown keys before init', async () => {
    const { useFeatureFlags } = await import('@/composables/useFeatureFlags')
    const ff = useFeatureFlags()
    expect(ff.isEnabled('nonexistent')).toBe(true)
  })

  it('init fetches flags and stores them', async () => {
    mockList.mockResolvedValue({
      analytics: true,
      recommendation: true,
      search: true,
      social: false,
      reactions: true,
      creator: false,
      moderation: true,
      ai: true,
      contribution: true,
      gamification: true,
      tips: true,
      subscription: true,
      notification: true,
    })

    const { useFeatureFlags } = await import('@/composables/useFeatureFlags')
    const ff = useFeatureFlags()
    await ff.init()

    expect(ff.store.fetched).toBe(true)
    expect(ff.isEnabled('social')).toBe(false)
    expect(ff.isEnabled('creator')).toBe(false)
    expect(ff.isEnabled('analytics')).toBe(true)
  })

  it('init does not fetch twice if already fetched', async () => {
    mockList.mockResolvedValue({
      analytics: true, recommendation: true, search: true,
      social: true, reactions: true, creator: true,
      moderation: true, ai: true, contribution: true,
      gamification: true, tips: true, subscription: true, notification: true,
    })

    const { useFeatureFlags } = await import('@/composables/useFeatureFlags')
    const ff = useFeatureFlags()
    await ff.init()
    expect(mockList).toHaveBeenCalledTimes(1)

    await ff.init()
    expect(mockList).toHaveBeenCalledTimes(1)
  })

  it('init handles API errors gracefully', async () => {
    mockList.mockRejectedValue(new Error('Network error'))

    const { useFeatureFlags } = await import('@/composables/useFeatureFlags')
    const ff = useFeatureFlags()
    await ff.init()

    expect(ff.store.fetched).toBe(true)
  })
})
