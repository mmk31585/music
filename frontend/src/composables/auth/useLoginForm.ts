import { reactive, ref } from 'vue'
import { useAuth } from './useAuth.ts'
import type { ApiResponseProps } from '@/plugins/client/types'

export function useLoginForm() {
  const { login } = useAuth()

  const loading = ref(false)
  const apiError = ref('')
  const rememberMe = ref(true)

  const form = reactive({
    email: '',
    password: '',
  })

  const errors = reactive({
    email: '',
    password: '',
  })

  function validateField(field: 'email' | 'password') {
    errors[field] = ''
    if (field === 'email' && !form.email.trim()) {
      errors.email = 'Email is required'
    }
    if (field === 'password' && !form.password.trim()) {
      errors.password = 'Password is required'
    }
  }

  function validate() {
    errors.email = ''
    errors.password = ''
    apiError.value = ''

    let valid = true

    if (!form.email.trim()) {
      errors.email = 'Email is required'
      valid = false
    }

    if (!form.password.trim()) {
      errors.password = 'Password is required'
      valid = false
    }

    return valid
  }

  async function onSubmit() {
    if (!validate()) return

    loading.value = true
    apiError.value = ''
    try {
      await login(form.email, form.password, { suppressToast: true, rememberMe: rememberMe.value })
    } catch (err: unknown) {
      // Check for 429 rate limit
      const axiosErr = err as Record<string, unknown>
      const status = (axiosErr?.response as Record<string, unknown> | undefined)?.status
      if (status === 429) {
        apiError.value = 'Too many attempts. Please try again later.'
      } else if (err instanceof Error) {
        apiError.value = err.message
      } else if (err && typeof err === 'object') {
        const payload = err as ApiResponseProps
        apiError.value = (payload.data as Record<string, unknown> | undefined)?.message as string
          || payload.message
          || 'Login failed'
      } else {
        apiError.value = 'Login failed'
      }
    } finally {
      loading.value = false
    }
  }

  return {
    form,
    errors,
    apiError,
    loading,
    rememberMe,
    validateField,
    onSubmit,
  }
}
