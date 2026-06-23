import { useToast } from 'primevue/usetoast'
import { translateMessage } from '@/utils/message-translations'

const DEFAULT_LIFETIME = 4000
const ERROR_LIFETIME = 6000

export function useAppToast() {
  const toast = useToast()

  function success(detail: string, summary?: string) {
    toast.add({ severity: 'success', summary: summary || 'Success', detail: translateMessage(detail) ?? detail, life: DEFAULT_LIFETIME })
  }

  function error(detail: string, summary?: string) {
    toast.add({ severity: 'error', summary: summary || 'Error', detail: translateMessage(detail) ?? detail, life: ERROR_LIFETIME })
  }

  function info(detail: string, summary?: string) {
    toast.add({ severity: 'info', summary: summary || 'Info', detail: translateMessage(detail) ?? detail, life: DEFAULT_LIFETIME })
  }

  function warn(detail: string, summary?: string) {
    toast.add({ severity: 'warn', summary: summary || 'Warning', detail: translateMessage(detail) ?? detail, life: DEFAULT_LIFETIME })
  }

  function apiError(err: any, fallback?: string) {
    const e = err as Record<string, any>
    const response = e.response as Record<string, any> | undefined
    const data = response?.data as Record<string, any> | undefined
    const message =
      data?.error ||
      data?.message ||
      e?.message ||
      fallback ||
      'Something went wrong'
    error(message, 'خطا در درخواست')
  }

  return { success, error, info, warn, apiError, toast }
}

export type AppToast = ReturnType<typeof useAppToast>
