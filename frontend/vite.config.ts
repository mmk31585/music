import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import Components from 'unplugin-vue-components/vite'
import { PrimeVueResolver } from '@primevue/auto-import-resolver'

// https://vite.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    host: '0.0.0.0',
     proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        ws: true,
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      output: {
        manualChunks(id: string) {
          // Vendor chunk: Vue ecosystem
          if (id.includes('node_modules/vue') ||
              id.includes('node_modules/pinia') ||
              id.includes('node_modules/vue-router') ||
              id.includes('node_modules/@vue')) {
            return 'vendor-vue'
          }
          // Heavy PrimeVue components (DataTable, Password) → separate chunk
          // These are only used in admin/auth pages, not in the main app shell
          if (id.includes('node_modules/primevue/datatable') ||
              id.includes('node_modules/primevue/column') ||
              id.includes('node_modules/primevue/paginator') ||
              id.includes('node_modules/primevue/password')) {
            return 'vendor-primevue-heavy'
          }
          // Vendor chunk: PrimeVue UI library (lightweight components)
          if (id.includes('node_modules/primevue') ||
              id.includes('node_modules/primeicons')) {
            return 'vendor-primevue'
          }
          // Vendor chunk: Other heavy deps
          if (id.includes('node_modules/lucide-vue-next') ||
              id.includes('node_modules/axios') ||
              id.includes('node_modules/zod') ||
              id.includes('node_modules/@vueuse') ||
              id.includes('node_modules/lodash')) {
            return 'vendor-libs'
          }
          // Shared core: composables, stores, and API types/enums/routes
          // These are tightly coupled (composables import stores import API types;
          // API routes import composables), so they must live in ONE chunk
          // to avoid Rollup circular-chunk duplication warnings.
          if (id.includes('/composables/') ||
              id.includes('/stores/') ||
              (id.includes('/services/api/') &&
               (id.endsWith('/enums.ts') || id.endsWith('/types.ts') || id.endsWith('/api.ts') || id.endsWith('routes.ts')))) {
            return 'shared-core'
          }
        },
      },
    },
  },
  plugins: [
    vue(),
    Components({
      resolvers: [PrimeVueResolver()],
    }),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
