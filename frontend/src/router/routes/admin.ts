import type { RouteRecordRaw } from 'vue-router'
import LayoutAdmin from '@/layouts/LayoutAdmin.vue'

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin',
    component: LayoutAdmin,
    meta: {
      requiresAuth: true,
      requiresRole: 'admin',
    },
    children: [
      {
        path: '',
        name: 'admin.dashboard',
        component: () => import('@/pages/admin/PageAdminDashboard.vue'),
        meta: {
          title: 'Admin Dashboard',
        },
      },
      {
        path: 'catalog',
        name: 'admin.catalog',
        component: () => import('@/pages/admin/PageAdminCatalog.vue'),
        meta: {
          title: 'Catalog',
        },
      },
      {
        path: 'tracks',
        name: 'admin.tracks',
        component: () => import('@/pages/admin/PageAdminTracks.vue'),
        meta: {
          title: 'Tracks',
        },
      },
      {
        path: 'artists',
        name: 'admin.artists',
        component: () => import('@/pages/admin/AdminArtistsPage.vue'),
        meta: {
          title: 'Artists',
        },
      },
      {
        path: 'albums',
        name: 'admin.albums',
        component: () => import('@/pages/admin/AdminAlbumsPage.vue'),
        meta: {
          title: 'Albums',
        },
      },
      {
        path: 'genres',
        name: 'admin.genres',
        component: () => import('@/pages/admin/AdminGenresPage.vue'),
        meta: {
          title: 'Genres',
        },
      },
    ],
  },
]
