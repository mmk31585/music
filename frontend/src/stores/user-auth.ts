import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import Cookie from 'js-cookie'

import { safeLocalStorage } from '@/services/storage'
import { useAuthApi } from '@/services/api/auth/routes'
import type { User, LoginPayload, AuthResponse } from '@/services/api/auth/types'
import type { UseRequestConfig } from '@/plugins/client/types'

interface TokenState {
  access_token: string
  refresh_token?: string
}

function createSafeNamespace<T>(key: string) {
  return {
    get: (): T | null => safeLocalStorage.getItem<T>(key),
    set: (value: T): void => safeLocalStorage.setItem<T>(key, value),
    remove: (): void => safeLocalStorage.removeItem(key),
  }
}

export const useUserAuthStore = defineStore('auth', () => {
  const cookieName = 'auth_cookie'
  const tokenCookieName = 'token_cookie'
  const deviceCookieName = 'device_cookie'

  const storage = createSafeNamespace<User>('user_auth')
  const tokenStorage = createSafeNamespace<TokenState>('user_token')
  const deviceStorage = createSafeNamespace<string>('device_id')

  const user = ref<User | null>(storage.get())
  const state = ref<TokenState | null>(tokenStorage.get())
  const deviceId = ref<string | null>(deviceStorage.get())
  const loading = ref(false)

  const isAuthenticated = computed(() => Boolean(state.value?.access_token))
  const isGuest = computed(() => user.value === null && state.value?.access_token == null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const token = computed(() => state.value?.access_token ?? null)
  const device = computed(() => deviceId.value ?? null)
  const isLoading = computed(() => loading.value)

  function persistToken(data: TokenState) {
    tokenStorage.set(data)
    Cookie.set(tokenCookieName, JSON.stringify(data), {
      secure: true,
      sameSite: 'lax',
    })
  }

  function persistUser(data: User) {
    storage.set(data)
    Cookie.set(cookieName, JSON.stringify(data), {
      secure: true,
      sameSite: 'lax',
    })
  }

  function persistDeviceId(id: string) {
    deviceStorage.set(id)
    Cookie.set(deviceCookieName, JSON.stringify(id), {
      secure: true,
      sameSite: 'lax',
    })
  }

  function setRefreshToken(refreshToken: string) {
    if (state.value) {
      state.value = { ...state.value, refresh_token: refreshToken }
      persistToken(state.value)
    }
  }

  function clearToken() {
    state.value = null
    tokenStorage.remove()
    Cookie.remove(tokenCookieName)
  }

  function clearUser() {
    user.value = null
    storage.remove()
    Cookie.remove(cookieName)
  }

  function setToken(newToken: string) {
    if (!state.value) {
      state.value = {
        access_token: newToken,
      }
    } else {
      state.value = {
        ...state.value,
        access_token: newToken,
      }
    }

    persistToken(state.value)
  }

  function setSession(payload: { access_token: string; refresh_token?: string; user: User }) {
    const tokenState: TokenState = {
      access_token: payload.access_token,
      refresh_token: payload.refresh_token,
    }

    state.value = tokenState
    user.value = payload.user

    persistToken(tokenState)
    persistUser(payload.user)
  }

  async function login(payload: LoginPayload, testing = false): Promise<User> {
    loading.value = true

    try {
      const config = testing
        ? { headers: { 'X-Testing': 'true' } } as UseRequestConfig<AuthResponse>
        : {} as UseRequestConfig<AuthResponse>

      const response = await useAuthApi().login(payload, config)

      if (!response.access_token) {
        throw new Error('Login failed: no access token received')
      }

      if (!response.user) {
        throw new Error('Login failed: no user data received')
      }

      setSession({
        access_token: response.access_token,
        refresh_token: response.refresh_token,
        user: response.user,
      })

      return response.user
    } finally {
      loading.value = false
    }
  }

  async function me(config: UseRequestConfig<any> = {}) {
    const userData = await useAuthApi().me(config)

    user.value = userData
    persistUser(userData)

    return userData
  }

  async function logout(): Promise<void> {
    const refreshToken = state.value?.refresh_token

    try {
      if (refreshToken) {
        await useAuthApi().logout({
          refreshToken,
        })
      }
    } catch (error) {
      console.warn('Logout request failed, clearing local session anyway:', error)
    } finally {
      $reset()
    }
  }

  function createDeviceId() {
    if (!device.value) {
      const id = uuidv4()

      deviceId.value = id
      persistDeviceId(id)
    }
  }

  let _readyPromise: Promise<void> | null = null
  let _readyResolve: (() => void) | null = null

  async function restore(): Promise<void> {
    const storedToken = tokenStorage.get()
    const storedUser = storage.get()
    const storedDevice = deviceStorage.get()

    if (storedToken) {
      state.value = storedToken
    }

    if (storedUser) {
      user.value = storedUser
    }

    if (storedDevice) {
      deviceId.value = storedDevice
    } else {
      createDeviceId()
    }

    if (token.value) {
      loading.value = true

      try {
        await me()
      } catch (error) {
        console.warn('Failed to restore auth session:', error)
        $reset()
      } finally {
        loading.value = false
      }
    }

    _readyResolve?.()
  }

  /** Returns a promise that resolves when the initial auth restore completes. */
  function ready(): Promise<void> {
    return _readyPromise ?? Promise.resolve()
  }

  function $reset(): void {
    clearUser()
    clearToken()
  }

  // Set up the ready promise
  _readyPromise = new Promise<void>((resolve) => {
    _readyResolve = resolve
  })

  // Fire restore asynchronously (don't block app mount)
  restore()

  return {
    user,
    entity: state,
    deviceId,

    loading: isLoading,
    isAuthenticated,
    isGuest,
    isAdmin,
    token,
    device,

    login,
    logout,
    me,
    restore,
    ready,
    setToken,
    setRefreshToken,
    setSession,
    clearToken,
    clearUser,
    $reset,
  }
})
