<template>
  <nav
    aria-label="Main navigation"
    class="fixed inset-x-0 bottom-0 z-40 border-t border-white/[0.06] bg-black/80 backdrop-blur-2xl lg:hidden"
    style="padding-bottom: max(0.5rem, env(safe-area-inset-bottom, 0.5rem))"
  >
    <div class="flex items-center justify-around px-2 pt-1">
      <RouterLink
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        class="flex flex-col items-center gap-0.5 rounded-xl px-3 py-1.5 text-[10px] font-medium transition-all"
        :class="isActive(tab.to) ? 'text-white' : 'text-slate-500 hover:text-slate-300'"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-xl transition-all duration-200"
          :class="isActive(tab.to) ? 'bg-[#1db954]/15 text-[#1db954]' : 'text-slate-500'"
        >
          <i aria-hidden="true" :class="tab.icon" class="text-lg" />
        </div>
        <span :class="isActive(tab.to) ? 'font-semibold' : ''">{{ tab.label }}</span>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const tabs = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'Search', icon: 'pi pi-search', to: '/search' },
  { label: 'Library', icon: 'pi pi-bookmark', to: '/library' },
  { label: 'Social', icon: 'pi pi-users', to: '/social' },
  { label: 'Profile', icon: 'pi pi-user', to: '/profile' },
]

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>
