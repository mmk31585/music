import { checkAuthGuard, checkLoginGuard, checkMaintenanceGuard } from './index'
import type { RouteLocationNormalized, RouteLocationNormalizedLoaded } from 'vue-router'

/**
 * In each hook, priority is important(from upper middleware to lower one)
 */
export default {
  beforeEach: [
    // async (to: RouteLocationNormalized) => await startPageLoading(to),
    (to: RouteLocationNormalized) => checkMaintenanceGuard(to),
    (to: RouteLocationNormalized) => checkLoginGuard(to),
    (to: RouteLocationNormalized, from: RouteLocationNormalizedLoaded, result: object | null) =>
      checkAuthGuard(to, from, result),
  ],

  // beforeResolve: [() => endPageLoading()],

  afterEach: [],
}
