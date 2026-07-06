import type { RouteLocationNormalized } from 'vue-router'
import { useUserAuthStore } from '@/stores'

/**
 * Returns null to continue, or a route-location object to redirect.
 *
 * Rules:
 * 1. Admin-only routes → redirect non-admins to home.
 * 2. redirectIfAdmin meta → redirect admins to admin dashboard.
 */
export async function checkAuthGuard(to: RouteLocationNormalized) {
  const auth = useUserAuthStore()
  // Wait for auth restore to complete before checking auth state
  await auth.ready()
  const isAuthenticated = auth.isAuthenticated
  const isAdmin = auth.isAdmin

  // Require admin role
  if (to.meta.requiresRole === 'admin') {
    if (!isAuthenticated) {
      return {
        name: 'auth.login',
        query: { redirect: to.fullPath },
      }
    }
    if (!isAdmin) {
      return { name: 'app.home' }
    }
  }

  // Redirect admins from normal app landing pages to admin dashboard
  if (to.meta.redirectIfAdmin && isAuthenticated && isAdmin) {
    return { name: 'admin.dashboard' }
  }

  return null
}
