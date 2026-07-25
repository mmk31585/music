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
        component: () => import(/* webpackChunkName: "auth-login" */ '@/pages/auth/PageLogin.vue'),
        meta: {
          title: 'Login',
        },
      },
      {
        path: 'register',
        name: 'auth.register',
        component: () => import(/* webpackChunkName: "auth-register" */ '@/pages/auth/PageRegister.vue'),
        meta: {
          title: 'Register',
        },
      },
      {
        path: 'forgot-password',
        name: 'auth.forgot-password',
        component: () => import(/* webpackChunkName: "auth-forgot-password" */ '@/pages/auth/PageForgotPassword.vue'),
        meta: {
          title: 'Forgot Password',
        },
      },
    ],
  },
]
