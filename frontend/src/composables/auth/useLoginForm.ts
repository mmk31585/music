import { reactive, ref } from 'vue'
import { useAuth } from './useAuth.ts'

export function useLoginForm() {
  const { login } = useAuth()

  const loading = ref(false)
  const apiError = ref('') // <-- add API error state

  const form = reactive({
    email: '',
    password: '',
  })

  const errors = reactive({
    email: '',
    password: '',
  })

  function validate() {
    errors.email = ''
    errors.password = ''
    apiError.value = '' // clear previous API error

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
      await login(form.email, form.password)
    } catch (err) {
      const errMsg = (err as Record<string, unknown>)?.message as string | undefined
      // Backend validation errors carry a more specific message inside .data
      const errData = (err as Record<string, unknown>)?.data as Record<string, unknown> | undefined
      apiError.value = (errData?.message as string) || errMsg || 'Login failed'
    } finally {
      loading.value = false
    }
  }

  return {
    form,
    errors,
    apiError, // expose for template
    loading,
    onSubmit,
  }
}
