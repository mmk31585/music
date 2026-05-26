import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const usePageLoaderStore = defineStore('usePageLoader', () => {
  const loading = ref(false)

  const isLoading = computed(() => {
    return loading.value
  })

  function setLoading(boolean: boolean) {
    loading.value = boolean
  }

  function $reset() {
    loading.value = false
  }

  return {
    loading,
    isLoading,
    setLoading,
    //
    $reset,
  }
})
