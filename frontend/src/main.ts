import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './assets/css/main.css'

import App from './App.vue'
import router from './router'

import LayoutEmpty from './components/layouts/LayoutEmpty.vue'
import { useUserAuthStore } from '@/stores'

import PrimeVue from 'primevue/config'
import { IndigoPreset, primeLocale } from '@/utils'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'


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
        preset: IndigoPreset,
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
    .use(ConfirmationService)

  app.component('layout-empty', LayoutEmpty)
  app.mount('#app')
})()
