import type { RouteLocationNormalized } from 'vue-router'
import { useUserAuthStore } from '@/stores'
import { usePermissionsApi } from '@/services/api/permissions'
import { useGlobalToast } from '@/composables/useRequest'

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
      useGlobalToast()?.add({ severity: 'warn', summary: 'Access denied', detail: 'This area is for administrators only', life: 4000 })
      return { name: 'app.home' }
    }
  }

	// Redirect admins from normal app landing pages to admin dashboard
	if (to.meta.redirectIfAdmin && isAuthenticated && isAdmin) {
		return { name: 'admin.dashboard' }
	}

	// Require specific permission
	if (to.meta.requiresPermission && isAuthenticated) {
		try {
			const access = await usePermissionsApi().getMyAccess()
			if (!access.permissions.includes(to.meta.requiresPermission as string)) {
				useGlobalToast()?.add({ severity: 'warn', summary: 'Access denied', detail: 'You don\'t have permission to access this page', life: 4000 })
				return { name: 'app.home' }
			}
		} catch {
			return { name: 'auth.login', query: { redirect: to.fullPath } }
		}
	}

	return null
}
