import { ref, watch } from 'vue'
import { useLibraryApi } from '@/services/api/library'
import { useUserAuthStore } from '@/stores/user-auth'

/**
 * Manages the liked state for a track, syncing with the backend library API.
 * Pass a reactive trackId ref to automatically check liked status on change.
 */
export function useTrackLike(trackIdRef: import('vue').Ref<string | undefined>) {
  const libraryApi = useLibraryApi()
  const auth = useUserAuthStore()

  const liked = ref(false)
  const loading = ref(false)
  const likedTrackIds = ref<Set<string>>(new Set())

  // Fetch all liked track IDs once on mount
  async function fetchLikedTracks() {
    if (!auth.isAuthenticated) return
    try {
      const tracks = await libraryApi.getLikedTracks()
      likedTrackIds.value = new Set((tracks || []).map((t: any) => String(t.track_id)))
    } catch {
      // silently fail
    }
  }

  // Check if current track is liked
  function checkLiked(trackId: string | undefined) {
    if (!trackId || !auth.isAuthenticated) {
      liked.value = false
      return
    }
    liked.value = likedTrackIds.value.has(trackId)
  }

  // Watch track changes
  watch(trackIdRef, (id) => checkLiked(id), { immediate: true })

  async function toggleLike() {
    const trackId = trackIdRef.value
    if (!trackId || !auth.isAuthenticated) return

    loading.value = true
    try {
      if (liked.value) {
        await libraryApi.unlikeTrack(trackId)
        likedTrackIds.value.delete(trackId)
        liked.value = false
      } else {
        await libraryApi.likeTrack({ track_id: trackId })
        likedTrackIds.value.add(trackId)
        liked.value = true
      }
    } catch {
      // revert on failure
      liked.value = !liked.value
    } finally {
      loading.value = false
    }
  }

  // Initialize
  fetchLikedTracks()

  return {
    liked,
    loading,
    toggleLike,
    fetchLikedTracks,
  }
}
