import { useAuthApi, type MaintenanceRequest } from '@/services/api'
import { safeLocalStorage } from '@/services/storage'
import { useMaintenanceStore } from '@/stores/maintenance'
import { computed, type Ref, ref } from 'vue'

export function useMaintenance() {
  const secret: Ref<string | null> = ref(
    (safeLocalStorage.getItem('maintenance_secret') as string | null) || null,
  )
  const getSecret = computed(() => {
    return secret.value
  })

  function setSecret(newSecret: string | null) {
    secret.value = newSecret
    safeLocalStorage.setItem('maintenance_secret', newSecret)
  }

  function removeSecret() {
    secret.value = ''
    safeLocalStorage.setItem('maintenance_secret', '')
  }

  const check = async () => {
    const maintenanceState = useMaintenanceStore()
    const secret = getSecret.value

    if (import.meta.env.VITE_IN_MAINTENANCE_MODE === 'true') {
      const bypassCode = import.meta.env.VITE_IN_MAINTENANCE_MODE_CODE
      if (!(secret && bypassCode && bypassCode === secret)) {
        maintenanceState.state = { inMaintenance: true, checked: true }
        return
      }
    }
    try {
      const data = {} as MaintenanceRequest
      if (secret) {
        data.maintenance_secret = secret
      }
      const res = await useAuthApi().maintenance(data)
      if (!res?.in_maintenance_mode) {
        removeSecret()
        maintenanceState.state = { inMaintenance: false, checked: true }
        return
      }
      if (secret) {
        setSecret(secret)
        maintenanceState.state = { inMaintenance: false, checked: true }
        return
      }
      maintenanceState.state = { inMaintenance: true, checked: true }
    } catch (e) {
      console.warn('Maintenance check failed', e)
      removeSecret()
      maintenanceState.state = { inMaintenance: false, checked: true }
    }
  }
  return {
    check,
    getSecret,
  }
}
