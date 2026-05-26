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
    console.log(form)
      await login(form.email, form.password)
      console.log(form)
    } catch (err) {
      apiError.value = err instanceof Error ? err.message : 'Login failed'
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
