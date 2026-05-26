import type { RouteLocationNormalized } from 'vue-router'

export function checkMaintenanceGuard(to: RouteLocationNormalized) {
  // Prevent redirect loop
  if (to.path.startsWith('/maintenance')) return

  if (to.path.startsWith('/api') || to.path.startsWith('/_')) return

  return null
}
