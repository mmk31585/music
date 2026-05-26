import type { RouteLocationNormalized } from 'vue-router'
import { useUserAuthStore } from '@/stores'

export function checkLoginGuard(to: RouteLocationNormalized) {
  const store = useUserAuthStore()

  if (
    to.matched.some((record) => record.meta.requiresAuth) &&
    to.name !== 'auth.login' &&
    to.name !== 'logout'
  ) {
    if (!store.token) {
      store.$reset()

      return {
        name: 'auth.login',
        query: { redirect: to.fullPath },
      }
    }
  }

  if (to.name === 'auth.login') {
    if (store.user && store.token) {
      return {
        name: 'app.home',
      }
    }
  }

  return null
}
