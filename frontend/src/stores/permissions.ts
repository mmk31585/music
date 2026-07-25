import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { usePermissionsApi } from '@/services/api/permissions'
import type { UserAccess } from '@/services/api/permissions/types'

export const usePermissionsStore = defineStore('permissions', () => {
  const access = ref<UserAccess | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const roleSlug = computed(() => access.value?.roleSlug ?? null)
  const roleLevel = computed(() => access.value?.roleLevel ?? 0)
  const permissions = computed(() => access.value?.permissions ?? [])

  const hasPermission = computed(() => (slug: string) => permissions.value.includes(slug))
  const hasAnyPermission = computed(() => (...slugs: string[]) => slugs.some((s) => permissions.value.includes(s)))

  const isAtLeast = computed(() => (minLevel: number) => roleLevel.value >= minLevel)

  async function fetchMyAccess() {
    loading.value = true
    error.value = null
    try {
      access.value = await usePermissionsApi().getMyAccess()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch permissions'
      access.value = null
    } finally {
      loading.value = false
    }
  }

  function $reset() {
    access.value = null
    loading.value = false
    error.value = null
  }

  return {
    access,
    loading,
    error,
    roleSlug,
    roleLevel,
    permissions,
    hasPermission,
    hasAnyPermission,
    isAtLeast,
    fetchMyAccess,
    $reset,
  }
})
