import type { RouteRecordRaw } from 'vue-router'
import LayoutAuth from '@/layouts/LayoutAuth.vue'

export const authRoutes: RouteRecordRaw[] = [
  {
    path: '/auth',
    component: LayoutAuth,
    meta: {
      guestOnly: true,
    },
    children: [
      {
        path: 'login',
        name: 'auth.login',
        component: () => import('@/pages/auth/PageLogin.vue'),
        meta: {
          title: 'Login',
        },
      },
      {
        path: 'register',
        name: 'auth.register',
        component: () => import('@/pages/auth/PageRegister.vue'),
        meta: {
          title: 'Register',
        },
      },
    ],
  },
]
