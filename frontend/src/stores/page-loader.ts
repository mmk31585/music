import { defineStore } from 'pinia'
import { ref } from 'vue'

export const usePageLoaderStore = defineStore('usePageLoader', () => {
  const loading = ref(false)

  function setLoading(boolean: boolean) {
    loading.value = boolean
  }

  function $reset() {
    loading.value = false
  }

  return {
    loading,
    setLoading,
    $reset,
  }
})
