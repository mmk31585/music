import type { Component } from 'vue'

export interface NavItem {
  id: string
  label: string
  icon: Component
  to?: string
  badge?: string | number
  badgeVariant?: 'default' | 'new' | 'count'
  section: NavSectionId
  permissions?: 'auth' | 'admin'
  future?: boolean
  children?: NavItem[]
}

export type NavSectionId =
  | 'discover'
  | 'your-music'
  | 'community'
  | 'ai'
  | 'creator'
  | 'settings'

export interface NavSection {
  id: NavSectionId
  label: string
  items: NavItem[]
}
