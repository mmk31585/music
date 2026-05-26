import { computed, ref } from 'vue'
import { type CaptchaProps, useAuthApi } from '@/services/api'
import { useIntervalFn } from '@vueuse/core'

export function useCaptcha() {
  const captcha = ref<CaptchaProps | null>(null)
  const loading = ref<boolean>(true)

  const isLoading = computed<boolean>(() => loading.value)
  const getCaptcha = computed<CaptchaProps | null>(() => captcha.value)
  const delay = ref(60000)

  async function fetch() {
    loading.value = true

    captcha.value = await useAuthApi().captcha()
    captcha.value.expires_in = captcha.value.expires_in - 1

    delay.value = captcha.value.expires_in * 1000

    loading.value = false
  }

  const control = useIntervalFn(
    async () => {
      await fetch()
    },
    delay,
    {
      immediate: false,
    },
  )

  fetch().then(() => {
    control.resume()
  })

  return {
    getCaptcha,
    loading: isLoading,
    fetch,
    fetchControl: control,
  }
}
