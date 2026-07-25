import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useUserAuthStore } from '@/stores'
import { useAuthApi } from '@/services/api/auth'
import { wsClient } from '@/services/socket/client'
import { usePlayerStore } from '@/stores/player'

export function useAuth() {
  const router = useRouter()
  const route = useRoute()
  const toast = useToast()
  const store = useUserAuthStore()
  const api = useAuthApi()
  const player = usePlayerStore()

  async function login(email: string, password: string, options?: { suppressToast?: boolean; rememberMe?: boolean }) {
    try {
      await store.login({ email, password }, { rememberMe: options?.rememberMe })

      // Use ?redirect= query param if present, otherwise go to home
      const redirect = route.query.redirect as string | undefined
      if (redirect && redirect.startsWith('/')) {
        await router.push(redirect)
      } else {
        await router.push({ name: 'app.home' })
      }

      toast.add({
        severity: 'success',
        summary: 'Welcome back',
        detail: 'Logged in successfully',
        life: 2500,
      })
    } catch (error) {
      // Only show toast if not suppressed (caller can handle their own display)
      if (!options?.suppressToast) {
        const axiosErr = error as Record<string, unknown>
        const status = (axiosErr?.response as Record<string, unknown> | undefined)?.status
        let msg: string
        if (status === 429) {
          msg = 'Too many attempts. Please try again later.'
        } else {
          msg =
            axiosErr?.message as string ??
            (error instanceof Error ? error.message : 'Please try again')
        }
        toast.add({
          severity: 'error',
          summary: 'Login failed',
          detail: msg,
          life: 5000,
        })
      }
      throw error
    }
  }

  async function register(payload: {
    displayName: string
    username: string
    email: string
    password: string
  }) {
    try {
      const response = await api.register(payload)

      store.setSession({
        access_token: response.access_token,
        refresh_token: response.refresh_token,
        user: response.user,
      })

      await router.push({ name: 'onboarding.genres' })

      toast.add({
        severity: 'success',
        summary: 'Account created',
        detail: 'Welcome to your music app',
        life: 2500,
      })
    } catch (error) {
      const msg =
        (error as Record<string, unknown>)?.message as string ??
        (error instanceof Error ? error.message : 'Please try again')
      toast.add({
        severity: 'error',
        summary: 'Registration failed',
        detail: msg,
        life: 5000,
      })
      throw error
    }
  }

  async function logout() {
    try {
      await store.logout()
    } catch {
      store.$reset()
    } finally {
      // Disconnect WebSocket — the old token in its URL is now invalid.
      wsClient.disconnect()
      // Reset player state so the next user starts clean.
      player.$reset()
      await router.push({ name: 'auth.login' })
    }
  }

  async function fetchMe() {
    return await store.me()
  }

  return {
    store,
    login,
    register,
    logout,
    fetchMe,
  }
}
