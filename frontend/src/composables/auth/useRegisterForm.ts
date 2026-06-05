import { reactive, ref } from 'vue'
import { useAuth } from './useAuth'

type RegisterErrors = {
  displayName: string
  username: string
  email: string
  password: string
}

type BackendValidationError = {
  success?: boolean
  code?: string
  message?: string
  details?: Record<string, string>
}

export function useRegisterForm() {
  const { register } = useAuth()

  const loading = ref(false)
  const apiError = ref('')

  const form = reactive({
    displayName: '',
    username: '',
    email: '',
    password: '',
  })

  const errors = reactive<RegisterErrors>({
    displayName: '',
    username: '',
    email: '',
    password: '',
  })

  function resetErrors() {
    errors.displayName = ''
    errors.username = ''
    errors.email = ''
    errors.password = ''
    apiError.value = ''
  }

  function validateEmail(email: string) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
  }

  function validateUsername(username: string) {
    return /^[a-zA-Z0-9._-]{3,20}$/.test(username)
  }

  function validate() {
    resetErrors()

    let valid = true

    if (!form.displayName.trim()) {
      errors.displayName = 'Display name is required'
      valid = false
    } else if (form.displayName.trim().length < 2) {
      errors.displayName = 'Display name must be at least 2 characters'
      valid = false
    }

    if (!form.username.trim()) {
      errors.username = 'Username is required'
      valid = false
    } else if (!validateUsername(form.username.trim())) {
      errors.username =
        'Username must be 3-20 characters and contain only letters, numbers, dot, dash, or underscore'
      valid = false
    }

    if (!form.email.trim()) {
      errors.email = 'Email is required'
      valid = false
    } else if (!validateEmail(form.email.trim())) {
      errors.email = 'Please enter a valid email address'
      valid = false
    }

    if (!form.password.trim()) {
      errors.password = 'Password is required'
      valid = false
    } else if (form.password.trim().length < 6) {
      errors.password = 'Password must be at least 6 characters'
      valid = false
    }

    return valid
  }

  function applyBackendErrors(error: unknown) {
    const err = error as {
      response?: {
        data?: BackendValidationError
      }
      message?: string
    }

    const data = err?.response?.data

    if (data?.code === 'VALIDATION_ERROR' && data.details) {
      if (data.details.displayName) {
        errors.displayName = data.details.displayName
      }
      if (data.details.username) {
        errors.username = data.details.username
      }
      if (data.details.email) {
        errors.email = data.details.email
      }
      if (data.details.password) {
        errors.password = data.details.password
      }

      apiError.value = data.message || 'Please fix the highlighted fields'
      return
    }

    apiError.value = err?.message || 'Registration failed'
  }

  async function onSubmit() {
    if (!validate()) return

    loading.value = true
    apiError.value = ''

    try {
      await register({
        displayName: form.displayName.trim(),
        username: form.username.trim(),
        email: form.email.trim(),
        password: form.password,
      })
    } catch (error) {
      applyBackendErrors(error)
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
