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
  },
]
export default routes
