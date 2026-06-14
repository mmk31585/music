import { safeLocalStorage } from '@/services/storage'
import { useMaintenanceStore } from '@/stores/maintenance'
import { computed, ref } from 'vue'

export function useMaintenance() {
  const secret = ref<string | null>(
    (safeLocalStorage.getItem('maintenance_secret') as string | null) || null,
  )
  const getSecret = computed(() => secret.value)

  const check = async () => {
    const maintenanceState = useMaintenanceStore()
    const s = getSecret.value

    if (import.meta.env.VITE_IN_MAINTENANCE_MODE === 'true') {
      const bypassCode = import.meta.env.VITE_IN_MAINTENANCE_MODE_CODE as string | undefined
      if (!(s && bypassCode && bypassCode === s)) {
        maintenanceState.state = { inMaintenance: true, checked: true }
        return
      }
    }
    maintenanceState.state = { inMaintenance: false, checked: true }
  }

  return { check, getSecret }
}
