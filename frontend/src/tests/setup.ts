import { config } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'

// Mock localStorage for test environment (jsdom doesn't provide it by default)
if (typeof window !== 'undefined' && !window.localStorage) {
  const store = new Map<string, string>()
  ;(window as Record<string, unknown>).localStorage = {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => store.set(key, value),
    removeItem: (key: string) => store.delete(key),
    clear: () => store.clear(),
    get length() { return store.size },
    key: (index: number) => [...store.keys()][index] ?? null,
  } as Storage
}

const router = createRouter({
  history: createWebHistory(),
  routes: [],
})

config.global.plugins = [router]
