import { ref } from 'vue'
import { useSocialApi } from '@/services/api/social'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useToast } from 'primevue/usetoast'
import type { ActivityFeedItem } from '@/services/api/social/types'
import type { RecommendationTrack } from '@/services/api/recommendation/types'

export function useDiscover() {
  const socialApi = useSocialApi()
  const recsApi = useRecommendationsApi()
  const toast = useToast()

  const trending = ref<RecommendationTrack[]>([])
  const forYou = ref<RecommendationTrack[]>([])
  const popular = ref<RecommendationTrack[]>([])
  const recent = ref<RecommendationTrack[]>([])
  const feed = ref<ActivityFeedItem[]>([])
  const loading = ref(false)
  const error = ref<unknown>(null)

  async function fetchDiscover() {
    loading.value = true
    error.value = null

    try {
      const [popularData, forYouData, recentData, feedData] = await Promise.all([
        recsApi.getPopular({ limit: 10 }).catch(() => null),
        recsApi.getForYou({ limit: 10 }).catch(() => null),
        recsApi.getRecent({ limit: 10 }).catch(() => null),
        socialApi.getFeed({ limit: 20, types: 'upload,follow' }).catch(() => null),
      ])

      if (popularData) {
        popular.value = popularData.items
      }

      if (forYouData) {
        forYou.value = forYouData.items
      }

      if (recentData) {
        recent.value = recentData.items
      }

      if (feedData) {
        feed.value = feedData.items
      }
    } catch (err: unknown) {
      error.value = err
      const msg = err instanceof Error ? err.message : 'Failed to load discover data'
      toast.add({ severity: 'error', summary: 'Discover Error', detail: msg, life: 5000 })
    } finally {
      loading.value = false
    }
  }

  return {
    trending,
    forYou,
    popular,
    recent,
    feed,
    loading,
    error,
    fetchDiscover,
  }
}
