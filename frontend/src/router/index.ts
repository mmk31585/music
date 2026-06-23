import { createRouter, createWebHistory } from 'vue-router'
import type { RouteLocationNormalized } from 'vue-router'
import routes from './routes'
import { checkMaintenanceGuard } from './middleware/maintenance-guard'
import { checkLoginGuard } from './middleware/login-guard'
import { checkAuthGuard } from './middleware/auth-guard'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

/**
 * Composes middleware functions into a single beforeEach chain.
 * Each middleware returns either a redirect object or null.
 * The first non-null result short-circuits and becomes the redirect.
 */
function chainMiddleware(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  middlewares: Array<(to: RouteLocationNormalized) => ReturnType<typeof checkLoginGuard>>,
) {
  for (const middleware of middlewares) {
    const result = middleware(to)
    if (result !== null && result !== undefined) return result
  }
  return true
}

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized) => {
  return chainMiddleware(to, from, [
    checkMaintenanceGuard,
    checkLoginGuard,
    checkAuthGuard,
  ])
})

router.afterEach((to) => {
  document.title = typeof to.meta.title === 'string' ? to.meta.title : 'Music App'
})

export default router
