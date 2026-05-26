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
      // The store's login expects a LoginPayload object
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

  async function register(payload: { email: string; password: string; name?: string }) {
    // Use the raw API because the store doesn't have a register action
    const response = await api.register(payload)
    // The register response should match the shape expected by setSession
    // Assuming response contains { access_token, user, refresh_token? }
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
  }

  async function logout() {
    try {
      // The store's logout already calls api.logout and resets the session
      await store.logout()
    } catch (error) {
      // Even if the API call fails, we still want to clear the local session
      store.$reset()
    } finally {
      await router.push({ name: 'auth.login' })
    }
  }

  async function fetchMe() {
    // The store's me method returns the user and updates the store
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
