import type { FeatureFlags, FeatureFlagKey } from '@/services/api/feature-flags'
import { useFeatureFlagsStore } from '@/stores/feature-flags'
import { useFeatureFlagsApi } from '@/services/api/feature-flags'

let initPromise: Promise<void> | null = null

export function useFeatureFlags() {
  const store = useFeatureFlagsStore()
  const api = useFeatureFlagsApi()

  async function init() {
    if (store.fetched) return
    if (initPromise) return initPromise
    initPromise = (async () => {
      store.loading = true
      try {
        const res = await api.list()
        if (res) {
          store.setFlags(res as unknown as FeatureFlags)
        }
      } catch {
        store.setFlags({} as FeatureFlags)
      } finally {
        store.loading = false
      }
    })()
    return initPromise
  }

  function isEnabled(key: string): boolean {
    return store.isEnabled(key as FeatureFlagKey)
  }

  return {
    store,
    init,
    isEnabled,
  }
}
