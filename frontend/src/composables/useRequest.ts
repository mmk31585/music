import { axiosClient, createRequestWrapper } from '@/plugins'
import type { RefreshToken, RequestHooks } from '@/plugins/client/types'
import { useUserAuthStore } from '@/stores'
import { useMaintenance } from '@/composables'
import { type Router } from 'vue-router'
import { useAuthApi } from '@/services/api'
import type { ToastServiceMethods } from 'primevue/toastservice'
import type { AxiosRequestConfig } from 'axios'

let router: Router | null = null
let toast: ToastServiceMethods | null = null

export function registerRouter(r: Router) {
  router = r
}

export function registerToast(t: ToastServiceMethods) {
  toast = t
}

/**
 * This MUST implement by real ui route, toast, store and storage
 */
const UIHooks: RequestHooks = {
  // ---------------- Toast ----------------
  toast: {
    success(msg: string) {
      toast?.add({ summary: msg, severity: 'success', life: 4000 })
    },
    error(msg: string) {
      toast?.add({ summary: msg, severity: 'error', life: 6000 })
    },
    info(msg: string) {
      toast?.add({ summary: msg, severity: 'info', life: 5000 })
    },
    warning(msg: string) {
      toast?.add({ summary: msg, severity: 'warn', life: 5000 })
    },
  },

  // ---------------- Maintenance ----------------
  getMaintenanceSecrets() {
    return useMaintenance().getSecret.value
  },

  // ---------------- Auth ----------------
  getAuthToken() {
    return useUserAuthStore().token
  },
  setAuthToken(token: string) {
    useUserAuthStore().setToken(token)
  },
  clearAuthToken() {
    useUserAuthStore().clearToken()
  },
  resetAuthStore() {
    useUserAuthStore().clearUser()
  },
  async redirectToLogin() {
    await router?.push({ name: 'login' })
  },

  // ---------------- Refresh token ----------------
  refreshToken(): Promise<RefreshToken> {
    const store = useUserAuthStore()
    const token = store.entity?.refresh_token
    return useAuthApi().refresh(token)
  },
  refreshTokenUrlRejecter(config: AxiosRequestConfig): boolean {
    return !!config.url?.includes('/auth/refresh')
  },

  // ---------------- Extras ----------------
  extraHeaders(): Record<string, string> {
    const user = useUserAuthStore()
    const device = user.device

    if (device) {
      return {
        'X-Device-Id': device,
      }
    }

    return {}
  },
}

export const { useRequest, client } = createRequestWrapper(axiosClient, UIHooks)
