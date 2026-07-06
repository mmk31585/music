import type { RouteLocationNormalized } from 'vue-router'
import { useUserAuthStore } from '@/stores'

/**
 * Returns null to continue, or a route-location object to redirect.
 *
 * Rules:
 * 1. Guest-only pages (login/register) → redirect authenticated users to home.
 * 2. Auth-required pages → redirect unauthenticated users to login.
 */
export async function checkLoginGuard(to: RouteLocationNormalized) {
  const auth = useUserAuthStore()
  // Wait for auth restore to complete before checking auth state
  await auth.ready()

  // Logged-in users should not visit guest pages like login/register
  if (to.meta.guestOnly && auth.isAuthenticated) {
    const isAdmin = auth.isAdmin
    return isAdmin ? { name: 'admin.dashboard' } : { name: 'app.home' }
  }

  // Require authentication
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return {
      name: 'auth.login',
      query: { redirect: to.fullPath },
    }
  }

  return null
}
