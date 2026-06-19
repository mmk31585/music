import type { RouteLocationNormalized, RouteLocationNormalizedLoaded } from 'vue-router'

export function checkAuthGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalizedLoaded,
  result: object | null,
) {
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
