import type { RouteLocationNormalized, RouteLocationNormalizedLoaded } from 'vue-router'
import { useUserAuthStore } from '@/stores'
// import { endPageLoading } from '@/router/middleware/page-progressbar.ts'

export function checkAuthGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalizedLoaded,
  result: object | null,
) {
  const store = useUserAuthStore()

  if (null !== result) {
    return result
  }

  if (
    to.name !== 'dashboard' &&
    to.meta?.placeTag &&
    to.meta?.placePermission
  ) {
    // endPageLoading()
    return {
      name: 'app.home',
    }
  }

  return null
}
