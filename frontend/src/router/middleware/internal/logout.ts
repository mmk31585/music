import { useUserAuthStore, usePlayerStore } from '@/stores'
import { isValidInternalRedirectLink } from '@/utils'
import { wsClient } from '@/services/socket/client'
import type {
  NavigationGuardNext,
  RouteLocationNormalized,
  RouteLocationNormalizedLoaded,
} from 'vue-router'

export const logoutMiddleware = async (
  to: RouteLocationNormalized,
  from: RouteLocationNormalizedLoaded,
  next: NavigationGuardNext,
) => {
  const store = useUserAuthStore()

  let route = null

  // Check if the user is NOT authenticated — nothing to do
  if (!store.user || !store.token) {
    next(from)
    return
  }

  // Logout the user
  await store.logout().then(() => {
    // Clean up WebSocket and player state
    wsClient.disconnect()
    usePlayerStore().$reset()

    if (from.meta.requiresAuth) {
      // If the original route requires authentication, redirect to log in
      route = { name: 'auth.login', query: {} }

      // @ts-expect-error It may have redirect in query
      if ('redirect' in to.query && isValidInternalRedirectLink(to.query?.redirect)) {
        route.query = { redirect: to.query.redirect }
      }
    }
  })

  if (route) {
    location.reload()
  }

  return next(from)
}
