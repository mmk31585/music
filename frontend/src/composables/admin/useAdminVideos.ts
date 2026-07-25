import { ref } from 'vue'
import { useVideoApi } from '@/services/api/video'
import type { VideoItem } from '@/services/api/video/types'

export function useAdminVideos() {
  const {
    adminGetVideos,
    adminUpdateVideo,
    adminDeleteVideo,
    adminApproveVideo,
  } = useVideoApi()

  const videos = ref<VideoItem[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<any>(null)
  const hasMore = ref(true)
  const offset = ref(0)
  const limit = 50

  async function fetchVideos(reset = true) {
    if (reset) {
      offset.value = 0
      hasMore.value = true
    }
    loading.value = true
    error.value = null

    try {
      const res = await adminGetVideos({ limit, offset: offset.value })
      if (reset) {
        videos.value = res.items || []
      } else {
        videos.value.push(...(res.items || []))
      }
      hasMore.value = (res.items?.length || 0) >= limit
      offset.value += res.items?.length || 0
      return res
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function loadMore() {
    if (loading.value || !hasMore.value) return
    await fetchVideos(false)
  }

  async function updateVideo(
    id: string,
    payload: {
      title?: string
      description?: string
      type?: string
      is_public?: boolean
      is_approved?: boolean
      status?: string
    },
  ) {
    saving.value = true
    error.value = null

    try {
      const updated = await adminUpdateVideo(id, payload)
      videos.value = videos.value.map((v) =>
        String(v.id) === String(id) ? updated : v,
      )
      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function deleteVideo(id: string) {
    deleting.value = true
    error.value = null

    try {
      await adminDeleteVideo(id)
      videos.value = videos.value.filter((v) => String(v.id) !== String(id))
    } catch (err) {
      error.value = err
      throw err
    } finally {
      deleting.value = false
    }
  }

  async function approveVideo(id: string) {
    saving.value = true
    error.value = null

    try {
      const updated = await adminApproveVideo(id)
      videos.value = videos.value.map((v) =>
        String(v.id) === String(id) ? updated : v,
      )
      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  function resetVideos() {
    videos.value = []
    offset.value = 0
    hasMore.value = true
    error.value = null
  }

  return {
    videos,
    loading,
    saving,
    deleting,
    error,
    hasMore,
    fetchVideos,
    loadMore,
    updateVideo,
    deleteVideo,
    approveVideo,
    resetVideos,
  }
}
