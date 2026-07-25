import { ref, onUnmounted } from 'vue'
import { useIngestionApi } from '@/services/api/ingestion/routes'
import type { DraftDetailResponse } from '@/services/api/ingestion/types'

export interface PollingOptions {
  intervalMs?: number
  maxAttempts?: number
  onComplete?: (detail: DraftDetailResponse) => void
  onTimeout?: () => void
  onError?: (err: unknown) => void
}

export function useEnrichmentPolling() {
  const ingestionApi = useIngestionApi()
  const polling = ref(false)
  const attempts = ref(0)

  let timer: ReturnType<typeof setInterval> | null = null
  const opts: Required<PollingOptions> = {
    intervalMs: 2000,
    maxAttempts: 30,
    onComplete: () => {},
    onTimeout: () => {},
    onError: () => {},
  }

  function start(draftId: string, options?: PollingOptions) {
    stop()
    if (options?.intervalMs) opts.intervalMs = options.intervalMs
    if (options?.maxAttempts) opts.maxAttempts = options.maxAttempts
    if (options?.onComplete) opts.onComplete = options.onComplete
    if (options?.onTimeout) opts.onTimeout = options.onTimeout
    if (options?.onError) opts.onError = options.onError

    polling.value = true
    attempts.value = 0

    timer = setInterval(async () => {
      attempts.value++
      try {
        const detail = await ingestionApi.getDraftDetail(draftId)
        if (detail && detail.status !== 'enriching' && detail.status !== 'refetching') {
          stop()
          opts.onComplete(detail)
          return
        }
      } catch (err) {
        opts.onError(err)
      }
      if (attempts.value >= opts.maxAttempts) {
        stop()
        opts.onTimeout()
      }
    }, opts.intervalMs)
  }

  function stop() {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    polling.value = false
    attempts.value = 0
  }

  onUnmounted(stop)

  return { start, stop, polling, attempts }
}
