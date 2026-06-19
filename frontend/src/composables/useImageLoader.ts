import { ref, watch } from 'vue'

interface ImageLoaderOptions {
  src: string | null | undefined
  placeholder?: string
  lazy?: boolean
}

export function useImageLoader(options: ImageLoaderOptions) {
  const { lazy = true } = options
  const isLoaded = ref(false)
  const hasError = ref(false)
  const currentSrc = ref(options.placeholder || '')
  const isLoading = ref(false)
  let observer: IntersectionObserver | null = null
  let shouldLoad = !lazy

  function loadImage(src: string) {
    if (!src) {
      hasError.value = true
      isLoading.value = false
      return
    }

    isLoading.value = true
    const img = new Image()

    img.onload = () => {
      isLoaded.value = true
      hasError.value = false
      currentSrc.value = src
      isLoading.value = false
    }

    img.onerror = () => {
      hasError.value = true
      isLoading.value = false
    }

    img.src = src
  }

  function triggerLoad() {
    if (shouldLoad || !options.src) return
    shouldLoad = true
    if (options.src) {
      loadImage(options.src)
    }
  }

  function bindLazy(el: HTMLElement | null) {
    if (!el || !lazy || shouldLoad) return

    observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0]
        if (entry?.isIntersecting) {
          triggerLoad()
          observer?.disconnect()
          observer = null
        }
      },
      { rootMargin: '200px 0px' },
    )

    observer.observe(el)
  }

  function unbindLazy() {
    observer?.disconnect()
    observer = null
  }

  watch(
    () => options.src,
    (newSrc) => {
      if (newSrc) {
        if (shouldLoad) {
          loadImage(newSrc)
        }
      } else {
        hasError.value = true
        isLoading.value = false
      }
    },
  )

  if (options.src && !lazy) {
    loadImage(options.src)
  } else if (options.src) {
    isLoading.value = true
  }

  return {
    isLoaded,
    hasError,
    currentSrc,
    isLoading,
    loadImage,
    triggerLoad,
    bindLazy,
    unbindLazy,
  }
}
