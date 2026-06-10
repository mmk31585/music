import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { Router, RouteLocationNormalizedGeneric, RouteRecordRaw } from 'vue-router'

const mockIsEnabled = vi.fn()
const mockInit = vi.fn()
let mockFetched = false

vi.mock('@/composables/useFeatureFlags', () => ({
  useFeatureFlags: () => ({
    store: {
      flags: null,
      loading: false,
      fetched: mockFetched,
      isLoaded: false,
      isEnabled: mockIsEnabled,
      setFlags: () => {},
      $reset: () => {},
      $patch: () => {},
      $subscribe: () => () => {},
      $onAction: () => () => {},
      $dispose: () => {},
    },
    init: mockInit,
    isEnabled: mockIsEnabled,
  }),
}))

function getRouteMeta(route: RouteLocationNormalizedGeneric) {
  let meta = route.meta
  if (!meta.featureFlag) {
    for (const record of route.matched) {
      if (record.meta?.featureFlag) {
        meta = record.meta
        break
      }
    }
  }
  return meta
}

function installGuard(router: Router) {
  router.beforeEach(async (to) => {
    const meta = getRouteMeta(to)
    const featureFlag = meta.featureFlag as string | undefined
    if (featureFlag) {
      const { useFeatureFlags } = await import('@/composables/useFeatureFlags')
      const ff = useFeatureFlags()
      if (!ff.store.fetched) {
        await ff.init()
      }
      if (!ff.store.isEnabled(featureFlag as import('@/services/api/feature-flags').FeatureFlagKey)) {
        return { name: 'app.home' }
      }
    }
    return true
  })
}

async function createTestRouter(routes: RouteRecordRaw[]): Promise<Router> {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', name: 'app.home', component: { template: '<div>Home</div>' } },
      ...routes,
    ],
  })
  installGuard(router)
  router.push('/')
  await router.isReady()
  return router
}

describe('featureFlag router guard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockIsEnabled.mockReset()
    mockInit.mockReset()
    mockFetched = false
  })

  it('allows navigation when feature is enabled', async () => {
    mockFetched = true
    mockIsEnabled.mockReturnValue(true)

    const router = await createTestRouter([
      {
        path: '/social',
        name: 'app.social',
        component: { template: '<div>Social</div>' },
        meta: { featureFlag: 'social' },
      },
    ])

    await router.push('/social')
    expect(router.currentRoute.value.name).toBe('app.social')
  })

  it('redirects to home when feature is disabled', async () => {
    mockFetched = true
    mockIsEnabled.mockReturnValue(false)

    const router = await createTestRouter([
      {
        path: '/social',
        name: 'app.social',
        component: { template: '<div>Social</div>' },
        meta: { featureFlag: 'social' },
      },
    ])

    await router.push('/social')
    expect(router.currentRoute.value.name).toBe('app.home')
  })

  it('calls init when store not fetched yet', async () => {
    mockFetched = false
    mockInit.mockImplementation(() => {
      mockFetched = true
    })
    mockIsEnabled.mockReturnValue(false)

    const router = await createTestRouter([
      {
        path: '/ai',
        name: 'app.ai',
        component: { template: '<div>AI</div>' },
        meta: { featureFlag: 'ai' },
      },
    ])

    await router.push('/ai')
    expect(mockInit).toHaveBeenCalledOnce()
    expect(router.currentRoute.value.name).toBe('app.home')
  })

  it('allows navigation when no featureFlag meta is set', async () => {
    mockFetched = true

    const router = await createTestRouter([
      {
        path: '/about',
        name: 'app.about',
        component: { template: '<div>About</div>' },
        meta: { title: 'About' },
      },
    ])

    await router.push('/about')
    expect(router.currentRoute.value.name).toBe('app.about')
    expect(mockIsEnabled).not.toHaveBeenCalled()
  })

  it('inherits featureFlag from parent route meta', async () => {
    mockFetched = true
    mockIsEnabled.mockReturnValue(false)

    const router = await createTestRouter([
      {
        path: '/recommendations',
        component: { template: '<div><router-view /></div>' },
        meta: { featureFlag: 'recommendation' },
        children: [
          {
            path: 'for-you',
            name: 'app.recommendations.for-you',
            component: { template: '<div>For You</div>' },
          },
        ],
      },
    ])

    await router.push('/recommendations/for-you')
    expect(router.currentRoute.value.name).toBe('app.home')
    expect(mockIsEnabled).toHaveBeenCalledWith('recommendation')
  })

  it('allows child route when parent feature is enabled', async () => {
    mockFetched = true
    mockIsEnabled.mockReturnValue(true)

    const router = await createTestRouter([
      {
        path: '/recommendations',
        component: { template: '<div><router-view /></div>' },
        meta: { featureFlag: 'recommendation' },
        children: [
          {
            path: 'for-you',
            name: 'app.recommendations.for-you',
            component: { template: '<div>For You</div>' },
          },
        ],
      },
    ])

    await router.push('/recommendations/for-you')
    expect(router.currentRoute.value.name).toBe('app.recommendations.for-you')
  })
})
