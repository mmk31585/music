import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserAuthStore } from '@/stores'
import { navSections } from './navigation'
import type { NavSection } from './types'

export function useNavigation() {
  const route = useRoute()
  const authStore = useUserAuthStore()

  const visibleSections = computed<NavSection[]>(() => {
    return navSections
      .map((section) => ({
        ...section,
        items: section.items.filter((item) => {
          if (item.permissions === 'admin' && !authStore.isAdmin) return false
          if (item.permissions === 'auth' && !authStore.isAuthenticated) return false
          return true
        }),
      }))
      .filter((section) => section.items.length > 0)
  })

  const activeNavId = computed(() => {
    const path = route.path
    for (const section of navSections) {
      for (const item of section.items) {
        if (!item.to) continue
        if (item.to === '/' ? path === '/' : path.startsWith(item.to)) {
          return item.id
        }
      }
    }
    return null
  })

  return { visibleSections, activeNavId }
}
