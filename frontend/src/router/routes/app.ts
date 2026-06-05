import type { RouteRecordRaw } from 'vue-router'
import LayoutMusicApp from '@/layouts/LayoutMusicApp.vue'

export const appRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: LayoutMusicApp,
    children: [
      {
        path: '',
        name: 'app.home',
        component: () => import('@/pages/app/PageHome.vue'),
      },
      {
        path: 'search',
        name: 'search',
        component: () => import('@/pages/app/PageSearch.vue'),
      },
      {
        path: 'recommendations',
        name: 'recommendations',
        component: () => import('@/pages/app/PageRecommendations.vue'),
      },
      {
        path: 'library',
        name: 'library',
        component: () => import('@/pages/app/PageLibrary.vue'),
      },
      {
        path: 'playlists',
        name: 'playlists',
        component: () => import('@/pages/app/PagePlaylists.vue'),
      },
      {
        path: 'recently-played',
        name: 'recently-played',
        component: () => import('@/pages/app/PageRecentlyPlayed.vue'),
      },
    ],
  },
]
