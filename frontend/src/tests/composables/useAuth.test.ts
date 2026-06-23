import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuth } from '@/composables/auth/useAuth'

const mockRouterPush = vi.fn()
const mockToastAdd = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockRouterPush,
  }),
  useRoute: () => ({}),
}))

vi.mock('primevue/usetoast', () => ({
  useToast: () => ({
    add: mockToastAdd,
  }),
}))

vi.mock('@/services/api/auth', () => ({
  useAuthApi: () => ({
    register: vi.fn().mockResolvedValue({
      access_token: 'access-123',
      refresh_token: 'refresh-456',
      user: { id: 'user-1', email: 'test@example.com', username: 'testuser', display_name: 'Test User' },
    }),
  }),
}))

vi.mock('@/services/storage', () => ({
  safeLocalStorage: {
    getItem: vi.fn().mockReturnValue(null),
    setItem: vi.fn(),
    removeItem: vi.fn(),
  },
}))

vi.mock('js-cookie', () => ({
  default: {
    set: vi.fn(),
    get: vi.fn().mockReturnValue(null),
    remove: vi.fn(),
  },
}))

vi.mock('uuid', () => ({
  v4: () => 'mocked-uuid-1234',
}))

describe('useAuth', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockRouterPush.mockReset()
    mockToastAdd.mockReset()
  })

  it('login calls store.login and navigates to home', async () => {
    // Mock the store's login to succeed
    const store = useAuth().store
    vi.spyOn(store, 'login').mockResolvedValue({ id: 'user-1' } as any)

    const auth = useAuth()
    await auth.login('test@example.com', 'password123')

    expect(store.login).toHaveBeenCalledWith({ email: 'test@example.com', password: 'password123' })
    expect(mockRouterPush).toHaveBeenCalledWith({ name: 'app.home' })
    expect(mockToastAdd).toHaveBeenCalled()
  })

  it('login shows error toast on failure', async () => {
    const store = useAuth().store
    vi.spyOn(store, 'login').mockRejectedValue(new Error('Invalid credentials'))

    const auth = useAuth()
    await expect(auth.login('test@example.com', 'wrong')).rejects.toThrow('Invalid credentials')
    expect(mockToastAdd).toHaveBeenCalledWith(
      expect.objectContaining({ severity: 'error' }),
    )
  })

  it('register creates account and navigates to onboarding', async () => {
    const auth = useAuth()
    const store = auth.store
    vi.spyOn(store, 'setSession').mockImplementation(() => {})

    await auth.register({
      displayName: 'New User',
      username: 'newuser',
      email: 'new@example.com',
      password: 'SecurePass1!',
    })

    expect(store.setSession).toHaveBeenCalled()
    expect(mockRouterPush).toHaveBeenCalledWith({ name: 'onboarding.genres' })
    expect(mockToastAdd).toHaveBeenCalledWith(
      expect.objectContaining({ severity: 'success' }),
    )
  })

  it('register shows error toast on failure', async () => {
    const authModule = await import('@/services/api/auth')
    ;(authModule.useAuthApi().register as any).mockRejectedValueOnce(
      new Error('Email already exists'),
    )

    const auth = useAuth()
    const store = auth.store
    vi.spyOn(store, 'setSession')

    await expect(
      auth.register({
        displayName: 'Fail User',
        username: 'failuser',
        email: 'exists@example.com',
        password: 'pass123',
      }),
    ).rejects.toThrow('Email already exists')

    expect(store.setSession).not.toHaveBeenCalled()
    expect(mockToastAdd).toHaveBeenCalledWith(
      expect.objectContaining({ severity: 'error' }),
    )
  })

  it('logout navigates to login page', async () => {
    const auth = useAuth()
    const store = auth.store
    vi.spyOn(store, 'logout').mockResolvedValue(undefined as any)
    vi.spyOn(store, '$reset').mockImplementation(() => {})

    await auth.logout()
    expect(mockRouterPush).toHaveBeenCalledWith({ name: 'auth.login' })
  })

  it('logout resets store on failure', async () => {
    const auth = useAuth()
    const store = auth.store
    vi.spyOn(store, 'logout').mockRejectedValue(new Error('Network error'))
    vi.spyOn(store, '$reset').mockImplementation(() => {})

    await auth.logout()
    expect(store.$reset).toHaveBeenCalled()
    expect(mockRouterPush).toHaveBeenCalledWith({ name: 'auth.login' })
  })
})
