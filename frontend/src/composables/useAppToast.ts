import { useToast } from 'primevue/usetoast'

type ToastSeverity = 'success' | 'error' | 'info' | 'warn'

const DEFAULT_LIFETIME = 4000
const ERROR_LIFETIME = 6000

export function useAppToast() {
  const toast = useToast()

  function success(detail: string, summary?: string) {
    toast.add({ severity: 'success', summary: summary || 'Success', detail, life: DEFAULT_LIFETIME })
  }

  function error(detail: string, summary?: string) {
    toast.add({ severity: 'error', summary: summary || 'Error', detail, life: ERROR_LIFETIME })
  }

  function info(detail: string, summary?: string) {
    toast.add({ severity: 'info', summary: summary || 'Info', detail, life: DEFAULT_LIFETIME })
  }

  function warn(detail: string, summary?: string) {
    toast.add({ severity: 'warn', summary: summary || 'Warning', detail, life: DEFAULT_LIFETIME })
  }

  function apiError(err: unknown, fallback?: string) {
    const message =
      (err as any)?.response?.data?.error ||
      (err as any)?.response?.data?.message ||
      (err as any)?.message ||
      fallback ||
      'Something went wrong'
    error(message, 'Request failed')
  }

  return { success, error, info, warn, apiError, toast }
}

export type AppToast = ReturnType<typeof useAppToast>
