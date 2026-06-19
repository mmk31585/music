<template>
  <aside
    class="fixed left-0 top-0 z-40 flex h-screen flex-col border-r border-white/10 bg-[#0a0a0a] transition-all duration-300 lg:static"
    :class="collapsed ? 'w-[68px]' : 'w-72'"
  >
    <!-- Brand -->
    <div class="flex shrink-0 items-center gap-3 px-3 pt-4 pb-3" :class="collapsed ? 'justify-center' : 'px-4 pt-4 pb-3'">
      <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[#1db954] font-bold text-black text-sm">
        A
      </div>
      <div v-if="!collapsed" class="overflow-hidden">
        <h1 class="text-base font-bold text-white truncate">Admin Panel</h1>
        <p class="text-[10px] text-slate-500 truncate">Catalog & media control</p>
      </div>
    </div>

    <!-- Toggle button (desktop) -->
    <button
      class="mx-2 flex items-center justify-center rounded-lg py-2 text-slate-500 hover:bg-white/5 hover:text-white transition-colors hidden lg:flex"
      @click="$emit('toggle')"
    >
      <i aria-hidden="true" class="pi text-xs" :class="collapsed ? 'pi-chevron-right' : 'pi-chevron-left'" />
    </button>

    <!-- Nav items -->
    <nav class="mt-2 flex-1 overflow-y-auto scroll-bar px-2">
      <template v-for="(section, sIdx) in navSections" :key="sIdx">
        <!-- Section label -->
        <p
          v-if="!collapsed && section.label"
          class="mb-1 mt-4 px-3 text-[10px] font-semibold tracking-wider text-slate-600 uppercase"
        >
          {{ section.label }}
        </p>

        <div class="space-y-0.5">
          <RouterLink
            v-for="item in section.items"
            :key="item.to"
            :to="item.to"
            :class="navItemClass(item)"
            @click="$emit('close')"
          >
            <i aria-hidden="true" :class="item.icon" class="text-base" />
            <span v-if="!collapsed" class="truncate text-sm">{{ item.label }}</span>
          </RouterLink>
        </div>
      </template>
    </nav>

    <!-- Bottom actions -->
    <div class="mt-auto shrink-0 border-t border-white/5 px-2 py-3">
      <RouterLink
        to="/"
        class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-slate-500 transition hover:bg-white/5 hover:text-white"
        :class="collapsed ? 'justify-center' : ''"
        @click="$emit('close')"
      >
        <i aria-hidden="true" class="pi pi-home" />
        <span v-if="!collapsed" class="text-sm">Back to app</span>
      </RouterLink>

      <button
        type="button"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-slate-500 transition hover:bg-white/5 hover:text-white"
        :class="collapsed ? 'justify-center' : ''"
        @click="handleLogout"
      >
        <i aria-hidden="true" class="pi pi-sign-out" />
        <span v-if="!collapsed" class="text-sm">Logout</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuth } from '@/composables/auth/useAuth'

const props = defineProps<{
  collapsed: boolean
}>()

defineEmits<{
  toggle: []
  close: []
}>()

const route = useRoute()
const { logout } = useAuth()

interface NavItem {
  label: string
  icon: string
  to: string
  exact?: boolean
}

interface NavSection {
  label?: string
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    label: 'Overview',
    items: [
      { label: 'Dashboard', icon: 'pi pi-th-large', to: '/admin', exact: true },
    ],
  },
  {
    label: 'Catalog',
    items: [
      { label: 'Catalog', icon: 'pi pi-database', to: '/admin/catalog' },
      { label: 'Tracks', icon: 'pi pi-play-circle', to: '/admin/tracks' },
      { label: 'Artists', icon: 'pi pi-users', to: '/admin/artists' },
      { label: 'Albums', icon: 'pi pi-images', to: '/admin/albums' },
      { label: 'Genres', icon: 'pi pi-tags', to: '/admin/genres' },
    ],
  },
  {
    label: 'Management',
    items: [
      { label: 'Users', icon: 'pi pi-user', to: '/admin/users' },
      { label: 'Media', icon: 'pi pi-upload', to: '/admin/media' },
      { label: 'Import from Internet', icon: 'pi pi-globe', to: '/admin/import' },
      { label: 'Import by Artist', icon: 'pi pi-user-plus', to: '/admin/import/artist' },
      { label: 'Music Ingestion', icon: 'pi pi-cloud-upload', to: '/admin/ingestion' },
      { label: 'Moderation', icon: 'pi pi-shield', to: '/admin/moderation' },
    ],
  },
  {
    label: 'Finance',
    items: [
      { label: 'Subscriptions', icon: 'pi pi-credit-card', to: '/admin/subscriptions' },
      { label: 'Contributions', icon: 'pi pi-heart', to: '/admin/contributions' },
    ],
  },
]

function isActive(item: NavItem): boolean {
  if (item.exact) {
    return route.path === item.to
  }
  return route.path.startsWith(item.to)
}

function navItemClass(item: NavItem): string {
  const active = isActive(item)
  return [
    'flex items-center gap-3 rounded-xl px-3 py-2.5 transition',
    props.collapsed && 'justify-center',
    active
      ? 'bg-[#1db954]/15 text-[#1db954]'
      : 'text-slate-400 hover:bg-white/5 hover:text-white',
  ].filter(Boolean).join(' ')
}

function handleLogout() {
  void logout()
}
</script>
