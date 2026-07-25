import { createRouter, createWebHistory } from 'vue-router'
import type { RouteLocationNormalized } from 'vue-router'
import routes from './routes'
import { checkMaintenanceGuard } from './middleware/maintenance-guard'
import { checkLoginGuard } from './middleware/login-guard'
import { checkAuthGuard } from './middleware/auth-guard'
import { usePageLoaderStore } from '@/stores/page-loader'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

/**
 * Composes middleware functions into a single async beforeEach chain.
 * Each middleware returns either a redirect object or null.
 * The first non-null result short-circuits and becomes the redirect.
 */
async function chainMiddleware(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  middlewares: Array<(to: RouteLocationNormalized) => ReturnType<typeof checkLoginGuard>>,
) {
  for (const middleware of middlewares) {
    const result = await middleware(to)
    if (result !== null && result !== undefined) return result
  }
  return true
}

router.beforeEach(async (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
  const loader = usePageLoaderStore()
  if (to.meta.noNeedRouteWaiting !== true) loader.setLoading(true)
  return chainMiddleware(to, from, [
    checkMaintenanceGuard,
    checkLoginGuard,
    checkAuthGuard,
  ])
})

router.afterEach((to) => {
  usePageLoaderStore().setLoading(false)
  const title = typeof to.meta.title === 'string' ? to.meta.title : 'Music App'
  document.title = title

  // Announce route change to screen readers
  const announcer = document.getElementById('route-announcer')
  if (announcer) {
    announcer.textContent = `Navigated to ${title}`
  }

  // Reset focus to main content or skip link after navigation
  // Use nextTick to wait for DOM update
  requestAnimationFrame(() => {
    const main = document.getElementById('main-content')
    if (main) {
      main.setAttribute('tabindex', '-1')
      main.focus({ preventScroll: true })
      // Remove tabindex after focus so it doesn't remain in tab order
      setTimeout(() => main.removeAttribute('tabindex'), 100)
    }
  })
})

export default router
