export interface BreadcrumbItem {
  name: string
  link?: string | (() => string)
  params?: string[]
}

export interface AppRouteMeta {
  requiresAuth?: boolean
  guestOnly?: boolean // ← add
  requiresRole?: string // ← add (e.g. 'admin')
  titleAppearance?: boolean
  title?: string
  noNeedRouteWaiting?: boolean
  breadcrumb?: BreadcrumbItem
  layout?: string
}
