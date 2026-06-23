import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserAuthStore } from '@/stores/user-auth'

vi.mock('@/services/api/auth', () => ({
  useAuthApi: () => ({
    login: vi.fn().mockResolvedValue({
      access_token: 'access-123',
      refresh_token: 'refresh-456',
      user: {
        id: 'user-1',
        email: 'test@example.com',
        username: 'testuser',
        display_name: 'Test User',
        role: 'listener',
      },
    }),
    refresh: vi.fn().mockResolvedValue({
      access_token: 'new-access-789',
      refresh_token: 'new-refresh-012',
    }),
    register: vi.fn().mockResolvedValue({
      access_token: 'access-123',
      refresh_token: 'refresh-456',
      user: { id: 'user-1', email: 'test@example.com', username: 'testuser' },
    }),
    me: vi.fn().mockResolvedValue({
      id: 'user-1',
      email: 'test@example.com',
      username: 'testuser',
      display_name: 'Test User',
      role: 'listener',
    }),
    logout: vi.fn().mockResolvedValue({}),
  }),
}))

vi.mock('uuid', () => ({
  v4: () => 'mocked-uuid-1234',
}))

vi.mock('@/services/storage', () => ({
  safeLocalStorage: {
    getItem: vi.fn().mockReturnValue(null),
    setItem: vi.fn(),
    removeItem: vi.fn(),
  },
}))

// Mock js-cookie
vi.mock('js-cookie', () => ({
  default: {
    set: vi.fn(),
    get: vi.fn().mockReturnValue(null),
    remove: vi.fn(),
  },
}))

describe('useUserAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts as guest with null values', () => {
    const store = useUserAuthStore()
    expect(store.user).toBeNull()
    expect(store.state).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.isGuest).toBe(true)
    expect(store.isAdmin).toBe(false)
    expect(store.token).toBeNull()
    expect(store.loading).toBe(false)
  })

  it('setToken updates access token', () => {
    const store = useUserAuthStore()
    store.setToken('new-token-abc')
    expect(store.token).toBe('new-token-abc')
    expect(store.isAuthenticated).toBe(true)
  })

  it('setRefreshToken updates refresh token', () => {
    const store = useUserAuthStore()
    store.setToken('access-token')
    store.setRefreshToken('refresh-token-xyz')
    expect(store.state?.refresh_token).toBe('refresh-token-xyz')
  })

  it('setSession sets user and tokens', () => {
    const store = useUserAuthStore()
    store.setSession({
      access_token: 'access-123',
      refresh_token: 'refresh-456',
      user: { id: 'user-1', email: 'test@example.com', username: 'testuser' } as any,
    })
    expect(store.isAuthenticated).toBe(true)
    expect(store.token).toBe('access-123')
    expect(store.user?.email).toBe('test@example.com')
    expect(store.isGuest).toBe(false)
  })

  it('isAdmin returns true for admin role', () => {
    const store = useUserAuthStore()
    store.setSession({
      access_token: 'admin-token',
      user: { id: 'admin-1', email: 'admin@example.com', username: 'admin', role: 'admin' } as any,
    })
    expect(store.isAdmin).toBe(true)
  })

  it('clearToken removes tokens', () => {
    const store = useUserAuthStore()
    store.setToken('some-token')
    expect(store.isAuthenticated).toBe(true)
    store.clearToken()
    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBeNull()
  })

  it('clearUser removes user data', () => {
    const store = useUserAuthStore()
    store.setSession({
      access_token: 'token',
      user: { id: 'user-1', email: 'test@example.com', username: 'testuser' } as any,
    })
    expect(store.user).not.toBeNull()
    store.clearUser()
    expect(store.user).toBeNull()
  })

  it('$reset clears everything', () => {
    const store = useUserAuthStore()
    store.setSession({
      access_token: 'token',
      refresh_token: 'refresh',
      user: { id: 'user-1', email: 'test@example.com', username: 'testuser' } as any,
    })
    expect(store.isAuthenticated).toBe(true)
    store.$reset()
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
    expect(store.state).toBeNull()
  })

  it('login sets session on success', async () => {
    const store = useUserAuthStore()
    const result = await store.login({ email: 'test@example.com', password: 'pass123' } as any, true)
    expect(result).toBeDefined()
    expect(store.isAuthenticated).toBe(true)
  })

  it('login throws error when no access token', async () => {
    // Override mock for this test
    const authModule = await import('@/services/api/auth')
    ;(authModule.useAuthApi().login as any).mockResolvedValueOnce({})

    const store = useUserAuthStore()
    await expect(
      store.login({ email: 'test@example.com', password: 'pass123' } as any, true),
    ).rejects.toThrow('no access token')
  })
})
