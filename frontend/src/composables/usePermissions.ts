import { usePermissionsStore } from '@/stores/permissions'

export function usePermissions() {
  const store = usePermissionsStore()

  function can(permSlug: string): boolean {
    return store.hasPermission(permSlug)
  }

  function canAny(...permSlugs: string[]): boolean {
    return store.hasAnyPermission(...permSlugs)
  }

  function isLevel(minLevel: number): boolean {
    return store.isAtLeast(minLevel)
  }

  const canUpload = () => can('upload_track')
  const canReview = () => canAny('review_uploads', 'manage_users')
  const canModerate = () => canAny('moderate_comments', 'moderate_users', 'moderate_clubs')
  const isAdmin = () => can('manage_system')
  const isModerator = () => can('review_uploads')

  return {
    can,
    canAny,
    isLevel,
    canUpload,
    canReview,
    canModerate,
    isAdmin,
    isModerator,
    store,
  }
}
