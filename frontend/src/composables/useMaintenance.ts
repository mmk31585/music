import { safeLocalStorage } from '@/services/storage'
import { useMaintenanceStore } from '@/stores/maintenance'
import { computed, ref } from 'vue'

export function useMaintenance() {
  const secret = ref<string | null>(
    safeLocalStorage.getItem<string>('maintenance_secret'),
  )
  const getSecret = computed(() => secret.value)

  const check = async () => {
    const maintenanceStore = useMaintenanceStore()
    const s = getSecret.value

    if (import.meta.env.VITE_IN_MAINTENANCE_MODE === 'true') {
      const bypassCode = import.meta.env.VITE_IN_MAINTENANCE_MODE_CODE as string | undefined
      if (!(s && bypassCode && bypassCode === s)) {
        maintenanceStore.inMaintenance = true
        maintenanceStore.checked = true
        return
      }
    }
    maintenanceStore.inMaintenance = false
    maintenanceStore.checked = true
  }

  return { check, getSecret }
}
