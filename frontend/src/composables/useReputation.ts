import { ref, computed } from 'vue'
import { useReputationApi } from '@/services/api/reputation'
import { usePermissionsApi } from '@/services/api/permissions'
import type { ReputationSummary } from '@/services/api/reputation/types'

export function useReputation(userId?: string | number) {
  const resolvedId = computed(() => userId != null ? String(userId) : undefined)
  const summary = ref<ReputationSummary | null>(null)
  const loading = ref(false)
  const error = ref<unknown>(null)

  const tier = computed(() => summary.value?.tier ?? 'newcomer')
  const tierLabel = computed(() => summary.value?.tierLabel ?? 'Newcomer')
  const trustScore = computed(() => summary.value?.trustScore ?? 0)
  const acceptedContributions = computed(() => summary.value?.acceptedContributions ?? 0)
  const uploadSlots = computed(() => summary.value?.uploadSlots ?? 0)
  const canAutoPublish = computed(() => summary.value?.autoPublish ?? false)
  const canReview = computed(() => summary.value?.canReview ?? false)

  async function fetchReputation() {
    const id = resolvedId.value
    if (!id) return
    loading.value = true
    error.value = null
    try {
      summary.value = await useReputationApi().getUserReputation(id)
    } catch (err) {
      error.value = err
      summary.value = null
    } finally {
      loading.value = false
    }
  }

  async function checkCanUpload(): Promise<boolean> {
    try {
      const access = await usePermissionsApi().getMyAccess()
      return access.permissions.includes('upload_track')
    } catch {
      return false
    }
  }

  return {
    summary,
    loading,
    error,
    tier,
    tierLabel,
    trustScore,
    acceptedContributions,
    uploadSlots,
    canAutoPublish,
    canReview,
    fetchReputation,
    checkCanUpload,
  }
}
