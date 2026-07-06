import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './assets/css/main.css'

import App from './App.vue'
import router from './router'

import LayoutEmpty from './components/layouts/LayoutEmpty.vue'

import PrimeVue from 'primevue/config'
import { AppPreset } from '@/utils'
import ToastService from 'primevue/toastservice'
import { Buffer } from 'buffer'
;(globalThis as Record<string, any>).Buffer = Buffer

import i18n from '@/locales'
import { getPrimeLocale, useLocaleStore } from '@/stores/locale'

// Preload local fonts (paths resolved by Vite for content-hashed filenames)
const fontUrls = [
  new URL('@/assets/fonts/iranyekanweblight.woff2', import.meta.url).href,
  new URL('@/assets/fonts/iranyekanwebregular.woff2', import.meta.url).href,
  new URL('@/assets/fonts/iranyekanwebbold.woff2', import.meta.url).href,
]
for (const href of fontUrls) {
  const link = document.createElement('link')
  link.rel = 'preload'
  link.as = 'font'
  link.type = 'font/woff2'
  link.href = href
  link.crossOrigin = 'anonymous'
  document.head.appendChild(link)
}

const app = createApp(App)
const pinia = createPinia()

// Expose Pinia globally so the PiP controller (which mounts a separate
// Vue app into the Document Picture-in-Picture window) can share the
// same reactive store instances.
;(window as any).__PINIA__ = pinia

app.use(pinia)
app.use(i18n)

void (async () => {
  // ── Initialize locale from saved preference ────────────────────────
  const localeStore = useLocaleStore()
  const initialLocale = localeStore.locale

  app
    .use(router)
    .use(PrimeVue, {
      ripple: true,
      theme: {
        rtl: initialLocale === 'fa',
        preset: AppPreset,
        options: {
          darkModeSelector: '.app-dark',
          cssLayer: {
            name: 'primevue',
            order: 'theme, base, primevue',
          },
        },
        locale: getPrimeLocale(initialLocale),
      },
    })
    .use(ToastService)

  // ── React to locale changes ────────────────────────────────────────
  // Update PrimeVue locale when the user switches languages.
  window.addEventListener('locale-change', ((e: CustomEvent) => {
    const { locale } = e.detail
    const primevue = app.config.globalProperties.$primevue
    if (primevue) {
      primevue.config.locale = getPrimeLocale(locale)
      primevue.config.rtl = locale === 'fa'
    }
  }) as EventListener)

  app.component('layout-empty', LayoutEmpty)
  app.mount('#app')
})()
