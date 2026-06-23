import { axiosClient, createRequestWrapper } from '@/plugins'
import type { RefreshToken, RequestHooks } from '@/plugins/client/types'
import { useUserAuthStore } from '@/stores/user-auth'
import { useMaintenance } from '@/composables/useMaintenance'
import { type Router } from 'vue-router'
import { useAuthApi } from '@/services/api/auth/routes'
import type { ToastServiceMethods } from 'primevue/toastservice'
import type { AxiosRequestConfig } from 'axios'

let router: Router | null = null
let toast: ToastServiceMethods | null = null
let _initialized = false
const _pendingInit: Array<() => void> = []

export function registerRouter(r: Router) {
  router = r
  checkInit()
}

export function registerToast(t: ToastServiceMethods) {
  toast = t
  checkInit()
}

function checkInit() {
  if (_initialized) return
  if (router && toast) {
    _initialized = true
    _pendingInit.splice(0).forEach((fn) => fn())
  }
}

/** Returns a promise that resolves when both router and toast are registered. */
export function useRequestReady(): Promise<void> {
  if (_initialized) return Promise.resolve()
  return new Promise((resolve) => {
    _pendingInit.push(resolve)
  })
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
    await router?.push({ name: 'auth.login' })
  },

  // ---------------- Refresh token ----------------
  refreshToken(): Promise<RefreshToken> {
    const store = useUserAuthStore()
    const token = store.entity?.refresh_token
    if (!token) {
      return Promise.reject(new Error('No refresh token available'))
    }
    return useAuthApi().refresh(token).then((res) => {
      if (res?.refresh_token) {
        store.setRefreshToken(res.refresh_token)
      }
      return res
    })
  },
  refreshTokenUrlRejecter(config: AxiosRequestConfig): boolean {
    const store = useUserAuthStore()
    // Don't try to refresh when:
    // 1. The failing request is itself the refresh endpoint (avoid loops), OR
    // 2. The user has no stored token (first-time 401, e.g. bad login credentials)
    return !!config.url?.includes('/auth/refresh') || !store.token
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
