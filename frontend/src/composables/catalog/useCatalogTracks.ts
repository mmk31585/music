import { onMounted, ref } from 'vue'
import { useCatalogApi } from '@/services/api/catalog'
import type { Track } from '@/services/api/catalog'

export function useCatalogTracks() {
  const api = useCatalogApi()

  const tracks = ref<Track[]>([])
  const loading = ref(false)

  async function fetchTracks() {
    loading.value = true
    try {
      const response = await api.getTracks()
      tracks.value = response
    } finally {
      loading.value = false
    }
  }

  onMounted(fetchTracks)

  return {
    tracks,
    loading,
    fetchTracks,
  }
}
