import { ref, computed } from 'vue'
import { useUserAuthStore } from '@/stores'
import { useSocialApi } from '@/services/api/social'
import { useReactionsApi } from '@/services/api/reactions'
import { useToast } from 'primevue/usetoast'

export function useUserProfile(userId?: string) {
  const auth = useUserAuthStore()
  const socialApi = useSocialApi()
  const reactionsApi = useReactionsApi()
  const toast = useToast()

  const targetUserId = computed(() => String(userId || auth.user?.id || ''))
  const isOwnProfile = computed(() => !userId || String(userId) === String(auth.user?.id))

  const followers = ref<any[]>([])
  const following = ref<any[]>([])
  const followerCount = ref(0)
  const followingCount = ref(0)
  const isFollowing = ref(false)
  const feed = ref<any[]>([])
  const likedTracks = ref<any[]>([])
  const likedAlbums = ref<any[]>([])
  const feedHasMore = ref(false)
  const loading = ref(false)
  const error = ref<unknown>(null)

  const displayName = computed(() => {
    return auth.user?.displayName || auth.user?.username || auth.user?.name || 'User'
  })

  async function fetchProfile() {
    if (!targetUserId.value) return
    loading.value = true
    error.value = null

    try {
      const uid = targetUserId.value
      const [followersData, followingData, feedData, likedTracksData, likedAlbumsData] =
        await Promise.all([
          socialApi.getFollowers(uid).catch(() => null),
          socialApi.getFollowing(uid).catch(() => null),
          socialApi.getFeed({ limit: 10 }).catch(() => null),
          reactionsApi.getLikedTracks({ limit: 10 }).catch(() => null),
          reactionsApi.getLikedAlbums({ limit: 10 }).catch(() => null),
        ])

      if (followersData) {
        followers.value = followersData.items
        followerCount.value = followersData.total_count
      }

      if (followingData) {
        following.value = followingData.items
        followingCount.value = followingData.total_count
      }

      if (feedData) {
        feed.value = feedData.items
        feedHasMore.value = feedData.pagination.has_more
      }

      if (likedTracksData) likedTracks.value = likedTracksData.items
      if (likedAlbumsData) likedAlbums.value = likedAlbumsData.items

      if (!isOwnProfile.value) {
        const followingCheck = await socialApi.isFollowing(uid).catch(() => null)
        if (followingCheck) isFollowing.value = followingCheck.is_following
      }
    } catch (err) {
      error.value = err
      const msg = err instanceof Error ? err.message : 'Failed to load profile'
      toast.add({ severity: 'error', summary: 'Profile Error', detail: msg, life: 5000 })
    } finally {
      loading.value = false
    }
  }

  async function toggleFollow() {
    if (!targetUserId.value || isOwnProfile.value) return
    try {
      const uid = targetUserId.value
      if (isFollowing.value) {
        await socialApi.unfollow(uid)
        isFollowing.value = false
        followerCount.value = Math.max(0, followerCount.value - 1)
      } else {
        await socialApi.follow(uid)
        isFollowing.value = true
        followerCount.value += 1
      }
    } catch {
      /* silent */
    }
  }

  return {
    targetUserId,
    isOwnProfile,
    displayName,
    followers,
    following,
    followerCount,
    followingCount,
    isFollowing,
    feed,
    likedTracks,
    likedAlbums,
    feedHasMore,
    loading,
    error,
    fetchProfile,
    toggleFollow,
  }
}
