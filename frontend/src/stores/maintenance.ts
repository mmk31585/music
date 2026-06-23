import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface MaintenanceState {
  inMaintenance: boolean
  checked: boolean
}

export const useMaintenanceStore = defineStore('useMaintenance', () => {
  const state = ref({
    inMaintenance: false,
    checked: false,
  })
  const isMaintenance = computed(() => state.value.inMaintenance)
  const isChecked = computed(() => state.value.checked)

  return {
    state,
    isChecked,
    isMaintenance,
  }
})
