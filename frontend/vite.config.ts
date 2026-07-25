import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import Components from 'unplugin-vue-components/vite'
import { PrimeVueResolver } from '@primevue/auto-import-resolver'
import { VitePWA } from 'vite-plugin-pwa'

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
        chunkFileNames: 'assets/[name]-[hash:8].js',
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
    VitePWA({
      registerType: 'prompt',
      includeAssets: ['icons/icon.svg'],
      manifest: {
        name: 'Muse — Persian Music Platform',
        short_name: 'Muse',
        description: 'Discover and stream Persian music. Create playlists, follow artists, and connect with the community.',
        display: 'standalone',
        orientation: 'portrait',
        start_url: '/',
        scope: '/',
        theme_color: '#050505',
        background_color: '#050505',
        lang: 'fa-IR',
        dir: 'rtl',
        categories: ['music', 'entertainment', 'social'],
        icons: [
          {
            src: 'icons/icon-192.png',
            sizes: '192x192',
            type: 'image/png',
            purpose: 'any',
          },
          {
            src: 'icons/icon-512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any',
          },
          {
            src: 'icons/icon-maskable-512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable',
          },
          {
            src: 'icons/icon.svg',
            sizes: 'any',
            type: 'image/svg+xml',
            purpose: 'any',
          },
        ],
        shortcuts: [
          {
            name: 'My Profile',
            short_name: 'Profile',
            description: 'View your profile',
            url: '/profile',
            icons: [{ src: 'icons/icon.svg', sizes: 'any' }],
          },
          {
            name: 'Liked Tracks',
            short_name: 'Tracks',
            description: 'View your liked tracks',
            url: '/library',
            icons: [{ src: 'icons/icon.svg', sizes: 'any' }],
          },
        ],
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,woff2}'],
        navigateFallback: '/index.html',
        navigateFallbackDenylist: [/^\/api\/v1\//, /^\/uploads\//],
        runtimeCaching: [
          {
            urlPattern: /\.(?:png|jpg|jpeg|gif|webp|avif|svg)$/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'muse-images',
              expiration: {
                maxEntries: 200,
                maxAgeSeconds: 30 * 24 * 60 * 60,
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
          {
            urlPattern: /^\/api\/v1\/(?!.*(?:stream|audio|hls|m3u8|ts|mp3|aac|flac|ogg|wav)).*$/,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'muse-api',
              networkTimeoutSeconds: 4,
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 24 * 60 * 60,
              },
            },
          },
          {
            urlPattern: /^\/api\/v1\/(?:gamification|recommendations\/stats)/,
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'muse-stats',
              expiration: {
                maxAgeSeconds: 60 * 60,
              },
            },
          },
        ],
      },
      devOptions: {
        enabled: true,
        type: 'module',
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
