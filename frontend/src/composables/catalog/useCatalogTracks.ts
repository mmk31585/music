import { onMounted, ref } from 'vue'
import {  useTracksApi } from '@/services/api/catalog'
import type { Track } from '@/services/api/catalog'

export function useCatalogTracks() {
  const api = useTracksApi()

  const tracks = ref<Track[]>([])
  const loading = ref(false)
  const error = ref<any>(null)

  async function fetchTracks() {
    loading.value = true
    error.value = null

    try {
      const response = await api.getTracks()

      tracks.value = Array.isArray(response) ? response : []
    } catch (err) {
      error.value = err
      console.error('Failed to fetch tracks:', err)
      tracks.value = []
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    void fetchTracks()
  })

  return {
    tracks,
    loading,
    error,
    fetchTracks,
  }
}
