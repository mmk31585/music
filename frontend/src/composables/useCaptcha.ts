import { ref } from 'vue'

export function useCaptcha() {
  const loading = ref(false)
  const getCaptcha = ref(null)

  async function fetch() {
    // Captcha not implemented on backend
  }

  return { getCaptcha, loading, fetch }
}
