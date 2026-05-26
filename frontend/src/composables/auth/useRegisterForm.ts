import { reactive, ref } from 'vue'
import { useAuth } from './useAuth.ts'

export function useRegisterForm() {
  const { register } = useAuth()

  const loading = ref(false)
  const apiError = ref('')

  const form = reactive({
    name: '',
    email: '',
    password: '',
  })

  const errors = reactive({
    name: '',
    email: '',
    password: '',
  })

  function validate() {
    errors.name = ''
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

    if (form.password.trim().length < 6) {
      errors.password = 'Password must be at least 6 characters'
      valid = false
    }

    return valid
  }

  async function onSubmit() {
    if (!validate()) return

    loading.value = true
    try {
      await register({
        name: form.name || undefined,
        email: form.email,
        password: form.password,
      })
    } catch (err) {
      apiError.value = err instanceof Error ? err.message : 'Registration failed'
    } finally {
      loading.value = false
    }
  }

  return {
    form,
    errors,
    apiError,
    loading,
    onSubmit,
  }
}
