import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './assets/css/main.css'

import App from './App.vue'
import router from './router'

import LayoutEmpty from './components/layouts/LayoutEmpty.vue'

import PrimeVue from 'primevue/config'
import { AppPreset, primeLocale } from '@/utils'
import ToastService from 'primevue/toastservice'
import { Buffer } from 'buffer'
  ;(globalThis as Record<string, any>).Buffer = Buffer


const app = createApp(App)
const pinia = createPinia()


app.use(pinia)

void (async () => {

  app
    .use(router)
    .use(PrimeVue, {
      ripple: true,
      theme: {
        rtl: true,
        preset: AppPreset,
        options: {
          darkModeSelector: '.app-dark',
          cssLayer: {
            name: 'primevue',
            order: 'theme, base, primevue',
          },
        },
        locale: primeLocale,
      },
    })
    .use(ToastService)

  app.component('layout-empty', LayoutEmpty)
  app.mount('#app')
})()
