import { createRouter, createWebHistory } from 'vue-router'
import type { RouteLocationNormalizedGeneric } from 'vue-router'
import routes from './routes'
import { useUserAuthStore } from '@/stores'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

function getHomeRoute(isAdmin: boolean) {
  return isAdmin ? { name: 'admin.dashboard' } : { name: 'app.home' }
}

router.beforeEach((to: RouteLocationNormalizedGeneric) => {
  const auth = useUserAuthStore()

  const isAuthenticated = auth.isAuthenticated
  const isAdmin = auth.isAdmin

  // Logged-in users should not visit guest pages like login/register
  if (to.meta.guestOnly && isAuthenticated) {
    return getHomeRoute(isAdmin)
  }

  // Require authentication
  if (to.meta.requiresAuth && !isAuthenticated) {
    return {
      name: 'auth.login',
      query: {
        redirect: to.fullPath,
      },
    }
  }

  // Guest mode: protect routes that require an account
  if (!isAuthenticated) {
    const guestRestricted = ['/library', '/playlists', '/profile', '/settings', '/subscriptions']
    if (guestRestricted.some((path) => to.path.startsWith(path))) {
      return {
        name: 'auth.login',
        query: {
          redirect: to.fullPath,
        },
      }
    }
  }

  // Require admin role
  if (to.meta.requiresRole === 'admin') {
    if (!isAuthenticated) {
      return {
        name: 'auth.login',
        query: {
          redirect: to.fullPath,
        },
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

  return true
})

router.afterEach((to) => {
  document.title = typeof to.meta.title === 'string' ? to.meta.title : 'Music App'
})

export default router
