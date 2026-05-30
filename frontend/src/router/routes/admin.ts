import type { RouteRecordRaw } from 'vue-router'
import LayoutAdmin from '@/layouts/LayoutAdmin.vue'
import { useUserAuthStore } from '@/stores'

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin',
    component: LayoutAdmin,
    meta: { requiresAuth: true, requiresRole: 'admin' },
    children: [
      {
        path: '',
        name: 'admin.dashboard',
        component: () => import('@/pages/admin/PageAdminDashboard.vue'),
      },
      {
        path: 'media',
        name: 'admin.media',
        component: () => import('@/pages/admin/PageAdminMediaUpload.vue'),
      },
      {
        path: 'catalog',
        name: 'admin.catalog',
        component: () => import('@/pages/admin/PageAdminCatalog.vue'),
      },
      {
        path: 'catalog/tracks',
        name: 'admin.catalog.tracks',
        component: () => import('@/pages/admin/PageAdminTracks.vue'),
      },
      {
        path: '/admin/artist',
        name: 'AdminArtists',
        component: () => import('@/pages/admin/AdminArtistsPage.vue'),
      },
      {
        path: '/admin/albums',
        name: 'AdminAlbums',
        component: () => import('@/pages/admin/AdminAlbumsPage.vue'),
      },
      {
        path: '/admin/genres',
        name: 'AdminGenres',
        component: () => import('@/pages/admin/AdminGenresPage.vue'),
      },
    ],
  },
]

export function setupRouteGuards(router: import('vue-router').Router) {
  router.beforeEach(async (to) => {
    const auth = useUserAuthStore()

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
      return { name: 'auth.login' }
    }

    if (to.meta.guestOnly && auth.isAuthenticated) {
      return { name: 'app.home' }
    }

    if (to.meta.requiresRole === 'admin' && !auth.isAdmin) {
      return { name: 'app.home' }
    }

    return true
  })
}
