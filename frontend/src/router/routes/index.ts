import type { RouteRecordRaw } from 'vue-router'
import { authRoutes } from './auth'
import { appRoutes } from './app'
import { adminRoutes } from './admin'

const routes: RouteRecordRaw[] = [
  {
    path: '/maintenance',
    name: 'maintenance',
    component: () => import('@/pages/errors/PageMaintenance.vue'),
    meta: {
      title: 'Under Maintenance',
      layout: 'layout-empty',
    },
  },
  ...authRoutes,
  {
    path: '/onboarding',
    meta: { layout: 'layout-empty' },
    children: [
      {
        path: 'genres',
        name: 'onboarding.genres',
        component: () => import('@/pages/onboarding/PageGenreOnboarding.vue'),
        meta: { title: 'Genre Onboarding', requiresAuth: true },
      },
    ],
  },
  ...appRoutes,
  ...adminRoutes,
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/pages/errors/PageNotFound.vue'),
    meta: {
      title: 'Page Not Found',
      layout: 'layout-empty',
    },
  },
]

export default routes
