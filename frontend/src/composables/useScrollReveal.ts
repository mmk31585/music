import { ref, onMounted, onUnmounted, watch } from 'vue'

interface ScrollRevealOptions {
  threshold?: number
  rootMargin?: string
  triggerOnce?: boolean
  immediate?: boolean
}

export function useScrollReveal(
  target: () => HTMLElement | null | undefined,
  options: ScrollRevealOptions = {},
) {
  const { threshold = 0.1, rootMargin = '0px 0px -40px 0px', triggerOnce = true, immediate } = options
  const isVisible = ref(false)
  const hasTriggered = ref(false)

  let observer: IntersectionObserver | null = null

  function observe() {
    const el = target()
    if (!el || observer) return

    observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0]
        if (!entry) return
        if (entry.isIntersecting) {
          isVisible.value = true
          if (triggerOnce) {
            hasTriggered.value = true
            observer?.disconnect()
            observer = null
          }
        } else if (!triggerOnce) {
          isVisible.value = false
        }
      },
      { threshold, rootMargin },
    )

    observer.observe(el)
  }

  function disconnect() {
    observer?.disconnect()
    observer = null
  }

  onMounted(() => {
    if (immediate) {
      isVisible.value = true
      hasTriggered.value = true
      return
    }
    observe()
  })

  onUnmounted(() => {
    disconnect()
  })

  watch(
    () => target(),
    (el) => {
      if (el && !hasTriggered.value && !immediate) {
        disconnect()
        observe()
      }
    },
  )

  return { isVisible, hasTriggered }
}
