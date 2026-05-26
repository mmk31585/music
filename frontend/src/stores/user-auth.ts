import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { safeLocalStorage } from '@/services/storage'
import { AuthResponseSchema, useAuthApi, UserSchema } from '@/services/api'  // assuming this still works
import Cookie from 'js-cookie'
import type { UseRequestConfig } from '@/plugins/client/types'
import type { User, LoginPayload } from '@/services/api'  // adjust import

// Local token shape (simpler than old AuthTokenProps)
interface TokenState {
  access_token: string
  refresh_token?: string
}

function createSafeNamespace<T>(key: string) {
  return {
    get: (): string | T | null => safeLocalStorage.getItem<T>(key),
    set: (value: T): void => safeLocalStorage.setItem<T>(key, value),
    remove: (): void => safeLocalStorage.removeItem(key),
  }
}

// Helper to extract token & user from a flexible API response
function normalizeAuthResponse(raw: unknown): { user?: User; token?: TokenState } {
  const parsed = AuthResponseSchema.parse(raw)  // you'll need to import this schema
  // Prefer top-level fields, then data.*
  const user = parsed.user ?? parsed.data?.user
  const accessToken = parsed.access_token ?? parsed.data?.access_token ?? parsed.token ?? parsed.data?.token
  const refreshToken = parsed.refresh_token ?? parsed.data?.refresh_token
  return {
    user: user as User | undefined,
    token: accessToken ? { access_token: accessToken, refresh_token: refreshToken } : undefined,
  }
}

export const useUserAuthStore = defineStore('auth', () => {
  const cookieName = 'auth_cookie'
  const storage = createSafeNamespace<User>('user_auth')
  const tokenCookieName = 'token_cookie'
  const tokenStorage = createSafeNamespace<TokenState>('user_token')
  const deviceCookieName = 'device_cookie'
  const deviceStorage = createSafeNamespace<string>('device_id')

  const user = ref(storage.get())
  const state = ref(tokenStorage.get()) // holds { access_token, refresh_token? }
  const deviceId = ref(deviceStorage.get())
  const loading = ref(false)

  const isAuthenticated = computed(() => !!state.value?.access_token)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const token = computed(() => state.value?.access_token ?? null)
  const device = computed(() => deviceId.value ?? null)
  const isLoading = computed(() => loading.value)

  function persistToken(data: TokenState) {
    tokenStorage.set(data)
    Cookie.set(tokenCookieName, JSON.stringify(data), { secure: true, sameSite: 'lax' })
  }

  function persistUser(data: User) {
    storage.set(data)
    Cookie.set(cookieName, JSON.stringify(data), { secure: true, sameSite: 'lax' })
  }

  function persistDeviceId(id: string) {
    deviceStorage.set(id)
    Cookie.set(deviceCookieName, JSON.stringify(id), { secure: true, sameSite: 'lax' })
  }

  function clearToken() {
    tokenStorage.remove()
    Cookie.remove(tokenCookieName)
  }

  function clearUser() {
    user.value = null
    state.value = null
  }

  function setToken(newToken: string) {
    if (!state.value) return
    state.value = {
      ...state.value,
      access_token: newToken,
    }
    persistToken(state.value)
  }
  function setSession(payload: { access_token: string; refresh_token?: string; user: User }) {
    const tokenState: TokenState = {
      access_token: payload.access_token,
      refresh_token: payload.refresh_token,
    }
    state.value = tokenState
    persistToken(tokenState)
    user.value = payload.user
    persistUser(payload.user)
  }
  async function login(payload: LoginPayload, testing = false): Promise<User> {
    let options = {}
    if (testing) {
      options = { headers: { 'X-Testing': 'true' } }
    }

    const rawResponse = await useAuthApi().login(payload, {}, options)
    const { user: userData, token: tokenData } = normalizeAuthResponse(rawResponse)

    if (!tokenData?.access_token) throw new Error('Login failed: no access token received')
    if (!userData) throw new Error('Login failed: no user data received')

    state.value = tokenData
    persistToken(tokenData)
    user.value = userData
    persistUser(userData)

    return userData
  }

  async function me(config: UseRequestConfig<User> = {}) {
    const rawUser = await useAuthApi().me(config)
    // Assuming me() already returns a valid User object
    const userData = UserSchema.parse(rawUser) // validate shape
    user.value = userData
    persistUser(userData)
    return userData
  }

  function logout(): Promise<void> {
    return useAuthApi().logout({
      success: () => $reset(),
    })
  }

  function createDeviceId() {
    if (!device.value) {
      deviceId.value = uuidv4()
      persistDeviceId(deviceId.value)
    }
  }

  function restore(): void {
    const storedToken = tokenStorage.get()
    if (storedToken) state.value = storedToken

    const storedDevice = deviceStorage.get()
    if (storedDevice) {
      deviceId.value = storedDevice
    } else {
      createDeviceId()
    }

    if (token.value) {
      loading.value = true
      me().finally(() => (loading.value = false))
    }
  }

  function $reset(): void {
    clearUser()
    clearToken()
  }

  restore()

  return {
    // state
    user,
    entity: state,
    deviceId,
    loading: isLoading,
    // computed
    isAuthenticated,
    isAdmin,
    token,
    device,
    // actions
    login,
    logout,
    me,
    restore,
    setToken,
    setSession,
    clearToken,
    clearUser,
    $reset,
  }
})
