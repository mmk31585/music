import type { RouteRecordRaw } from 'vue-router'
import { authRoutes } from './auth'
import { appRoutes } from './app'
import { adminRoutes } from './admin'

const routes: RouteRecordRaw[] = [
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
