import { ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useLibraryApi } from '@/services/api/library'
import { useUserAuthStore } from '@/stores/user-auth'

const GUEST_FREE_ACTION_KEY = 'guest_free_action_used'

function isGuestFreeActionUsed(): boolean {
  return localStorage.getItem(GUEST_FREE_ACTION_KEY) === 'true'
}

function markGuestFreeActionUsed(): void {
  localStorage.setItem(GUEST_FREE_ACTION_KEY, 'true')
}

/**
 * Singleton set of liked track IDs shared across all instances.
 * Only the first caller fetches from the API; subsequent instances
 * read the cached set and stay in sync.
 */
let _fetched = false
const _likedTrackIds = ref<Set<string>>(new Set())

/** Track ID pending undo from an unlike action (consumed by App.vue undo toast). */
export const pendingUnlikeTrackId = ref<string | null>(null)

/**
 * Manages the liked state for a track, syncing with the backend library API.
 * Pass a reactive trackId ref to automatically check liked status on change.
 *
 * All instances share the same liked-track cache, so liking a track in one
 * component immediately reflects in all others.
 */
export function useTrackLike(trackIdRef: import('vue').Ref<string | undefined>) {
  const libraryApi = useLibraryApi()
  const auth = useUserAuthStore()

  const liked = ref(false)
  const loading = ref(false)

  // Fetch all liked track IDs once — singleton across all callers
  if (!_fetched) {
    _fetched = true
    fetchLikedTracks()
  }

  async function fetchLikedTracks() {
    if (!auth.isAuthenticated) return
    try {
      const tracks = await libraryApi.getLikedTracks()
      _likedTrackIds.value = new Set((tracks || []).map((t: any) => String(t.track_id)))
    } catch {
      // silently fail
    }
  }

  // Check if current track is liked
  function checkLiked(trackId: string | undefined) {
    if (!trackId) {
      liked.value = false
      return
    }
    liked.value = _likedTrackIds.value.has(trackId)
  }

  // Watch track changes
  watch(trackIdRef, (id) => checkLiked(id), { immediate: true })

  const toast = useToast()

  async function toggleLike() {
    const trackId = trackIdRef.value
    if (!trackId) return

    if (!auth.isAuthenticated) {
      if (isGuestFreeActionUsed()) {
        toast.add({
          severity: 'info',
          summary: 'Sign up to keep using this feature',
          life: 4000,
        })
        return
      }
      markGuestFreeActionUsed()
      // Local-only like/unlike without API call
      if (liked.value) {
        _likedTrackIds.value.delete(trackId)
        liked.value = false
      } else {
        _likedTrackIds.value.add(trackId)
        liked.value = true
      }
      return
    }

    loading.value = true
    try {
      if (liked.value) {
        await libraryApi.unlikeTrack(trackId)
        _likedTrackIds.value.delete(trackId)
        liked.value = false
        pendingUnlikeTrackId.value = trackId
        toast.add({
          severity: 'success',
          summary: 'Track removed from library',
          group: 'undo',
          life: 6000,
        })
      } else {
        await libraryApi.likeTrack({ track_id: trackId })
        _likedTrackIds.value.add(trackId)
        liked.value = true
      }
    } catch {
      // revert on failure
      liked.value = !liked.value
    } finally {
      loading.value = false
    }
  }

  return {
    liked,
    loading,
    toggleLike,
    fetchLikedTracks,
  }
}
