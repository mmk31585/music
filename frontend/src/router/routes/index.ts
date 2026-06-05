import type { RouteRecordRaw } from 'vue-router'
import { authRoutes } from './auth'
import { appRoutes } from './app'
import { adminRoutes } from './admin'

const routes: RouteRecordRaw[] = [
  ...authRoutes,
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
