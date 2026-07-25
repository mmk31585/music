import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { axiosClient } from '@/plugins'

export const useMaintenanceStore = defineStore('useMaintenance', () => {
  const inMaintenance = ref(false)
  const checked = ref(false)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isMaintenance = computed(() => inMaintenance.value)
  const isChecked = computed(() => checked.value)

  /**
   * Fetch maintenance status from the public endpoint.
   * Sets inMaintenance and checked flags.
   */
  async function checkStatus() {
    if (checked.value) return
    loading.value = true
    error.value = null
    try {
      const res = await axiosClient.get('/status/maintenance')
      inMaintenance.value = res.data?.data?.maintenance === true
      checked.value = true
    } catch (err: any) {
      // If the endpoint doesn't exist (older backend), silently degrade
      if (err?.response?.status === 404) {
        inMaintenance.value = false
        checked.value = true
        return
      }
      error.value = err?.message || 'Failed to check maintenance status'
      inMaintenance.value = false
      checked.value = true
    } finally {
      loading.value = false
    }
  }

  /**
   * Toggle maintenance mode (admin only).
   */
  async function toggleMaintenance(enabled: boolean) {
    loading.value = true
    error.value = null
    try {
      const res = await axiosClient.post('/admin/dashboard/maintenance', { enabled })
      inMaintenance.value = res.data?.data?.maintenance === true
    } catch (err: any) {
      error.value = err?.response?.data?.message || err?.message || 'Failed to toggle maintenance'
      throw err
    } finally {
      loading.value = false
    }
  }

  function setMaintenance(value: boolean) {
    inMaintenance.value = value
  }

  return {
    inMaintenance,
    checked,
    loading,
    error,
    isMaintenance,
    isChecked,
    checkStatus,
    toggleMaintenance,
    setMaintenance,
  }
})
