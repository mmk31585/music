import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useUserAuthStore } from '@/stores'
import { useAuthApi } from '@/services/api/auth'

export function useAuth() {
  const router = useRouter()
  const toast = useToast()
  const store = useUserAuthStore()
  const api = useAuthApi()

  async function login(email: string, password: string) {
    try {
      await store.login({ email, password })
      await router.push({ name: 'app.home' })
      toast.add({
        severity: 'success',
        summary: 'Welcome back',
        detail: 'Logged in successfully',
        life: 2500,
      })
    } catch (error) {
      toast.add({
        severity: 'error',
        summary: 'Login failed',
        detail: error instanceof Error ? error.message : 'Please try again',
        life: 5000,
      })
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

      await router.push({ name: 'app.home' })

      toast.add({
        severity: 'success',
        summary: 'Account created',
        detail: 'Welcome to your music app',
        life: 2500,
      })
    } catch (error) {
      toast.add({
        severity: 'error',
        summary: 'Registration failed',
        detail: error instanceof Error ? error.message : 'Please try again',
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
