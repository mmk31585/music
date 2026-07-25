import type { RouteLocationNormalized } from 'vue-router'
import { useMaintenanceStore } from '@/stores/maintenance'

/**
 * Redirect to /maintenance when the app is in maintenance mode.
 * Returns a redirect object, or null to continue.
 */
export async function checkMaintenanceGuard(to: RouteLocationNormalized) {
  // Prevent redirect loop
  if (to.path.startsWith('/maintenance')) return null

  // Skip API/internal paths
  if (to.path.startsWith('/api') || to.path.startsWith('/_')) return null

  const maintenance = useMaintenanceStore()

  // Ensure maintenance status has been checked
  if (!maintenance.isChecked) {
    await maintenance.checkStatus()
  }

  if (maintenance.isMaintenance) {
    return { name: 'maintenance' }
  }

  return null
}
