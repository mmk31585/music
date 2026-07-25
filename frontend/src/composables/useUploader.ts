import { ref } from 'vue'
import { useUploadsApi } from '@/services/api/uploads'
import { useReputationApi } from '@/services/api/reputation'
import type { Draft, UploadSlotStatus } from '@/services/api/uploads/types'

export function useUploader() {
  const drafts = ref<Draft[]>([])
  const slots = ref<UploadSlotStatus | null>(null)
  const loading = ref(false)

  async function fetchMyDrafts() {
    loading.value = true
    try {
      const result = await useUploadsApi().getMyDrafts()
      if (Array.isArray(result)) {
        drafts.value = result
      } else if (result && typeof result === 'object' && 'items' in result) {
        drafts.value = (result as any).items ?? []
      } else {
        drafts.value = result ? [result] : []
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchSlots() {
    try {
      slots.value = await useUploadsApi().getUploadSlots()
    } catch {
      slots.value = { available: false, used: 0, max: 0 }
    }
  }

  async function submitUpload(filename: string, file?: File, clubId?: number): Promise<Draft | null> {
    const draft = await useUploadsApi().createDraft(filename, clubId)
    if (file && draft?.id) {
      await useUploadsApi().uploadDraftFile(draft.id, file)
    }
    await fetchSlots()
    await fetchMyDrafts()
    return draft
  }

  async function canUserUpload(userId: string): Promise<boolean> {
    try {
      const rep = await useReputationApi().getUserReputation(userId)
      return rep.uploadSlots > 0
    } catch {
      return false
    }
  }

  return {
    drafts,
    slots,
    loading,
    fetchMyDrafts,
    fetchSlots,
    submitUpload,
    canUserUpload,
  }
}
