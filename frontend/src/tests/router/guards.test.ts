import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { Router, RouteRecordRaw } from 'vue-router'
import { useUserAuthStore } from '@/stores/user-auth'

function createTestRouter(routes: RouteRecordRaw[]): Router {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', name: 'app.home', component: { template: '<div>Home</div>' }, meta: {} },
      { path: '/login', name: 'auth.login', component: { template: '<div>Login</div>' }, meta: { guestOnly: true } },
      { path: '/register', name: 'auth.register', component: { template: '<div>Register</div>' }, meta: { guestOnly: true } },
      { path: '/admin/dashboard', name: 'admin.dashboard', component: { template: '<div>Admin</div>' }, meta: { requiresRole: 'admin' } },
      { path: '/profile', name: 'app.profile', component: { template: '<div>Profile</div>' }, meta: { requiresAuth: true } },
      ...routes,
    ],
  })

  // Router guards from index.ts
  router.beforeEach((to) => {
    const auth = useUserAuthStore()
    const isAuthenticated = auth.isAuthenticated
    const isAdmin = auth.isAdmin

    if (to.meta.guestOnly && isAuthenticated) {
      return isAdmin ? { name: 'admin.dashboard' } : { name: 'app.home' }
    }

    if (to.meta.requiresAuth && !isAuthenticated) {
      return { name: 'auth.login', query: { redirect: to.fullPath } }
    }

    if (to.meta.requiresRole === 'admin') {
      if (!isAuthenticated) {
        return { name: 'auth.login', query: { redirect: to.fullPath } }
      }
      if (!isAdmin) {
        return { name: 'app.home' }
      }
    }

    return true
  })

  return router
}

describe('Router Guards', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('allows unauthenticated access to login page', async () => {
    const router = createTestRouter([])
    router.push('/login')
    await router.isReady()
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('redirects authenticated user from login to home', async () => {
    const router = createTestRouter([])
    const auth = useUserAuthStore()
    auth.setToken('some-token')
    auth.$patch({ user: { id: 'user-1', role: 'listener' } as any })

    router.push('/login')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('app.home')
  })

  it('redirects unauthenticated user from protected route to login', async () => {
    const router = createTestRouter([])
    router.push('/profile')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('auth.login')
  })

  it('redirects non-admin from admin route', async () => {
    const router = createTestRouter([])
    const auth = useUserAuthStore()
    auth.setToken('listener-token')
    auth.$patch({ user: { id: 'user-1', role: 'listener' } as any })

    router.push('/admin/dashboard')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('app.home')
  })

  it('allows admin access to admin routes', async () => {
    const router = createTestRouter([])
    const auth = useUserAuthStore()
    auth.setToken('admin-token')
    auth.$patch({ user: { id: 'admin-1', role: 'admin' } as any })

    router.push('/admin/dashboard')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('admin.dashboard')
  })

  it('redirects admin from guest-only pages to admin dashboard', async () => {
    const router = createTestRouter([])
    const auth = useUserAuthStore()
    auth.setToken('admin-token')
    auth.$patch({ user: { id: 'admin-1', role: 'admin' } as any })

    router.push('/login')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('admin.dashboard')
  })

  it('allows unauthenticated access to public routes', async () => {
    const router = createTestRouter([
      { path: '/about', name: 'app.about', component: { template: '<div>About</div>' } },
    ])
    router.push('/about')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('app.about')
  })

  it('redirects unauthenticated user from admin route to login', async () => {
    const router = createTestRouter([])
    router.push('/admin/dashboard')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('auth.login')
  })
})
