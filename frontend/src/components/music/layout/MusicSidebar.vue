<template>
  <nav aria-label="Main navigation"
    class="hidden h-screen shrink-0 border-r border-white/10 bg-black/40 backdrop-blur-2xl lg:block lg:overflow-y-auto transition-all duration-300 ease-out z-30"
    :class="collapsed ? 'w-14' : 'w-56'"
    style="backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px); scrollbar-width: thin; scrollbar-color: rgba(255,255,255,0.06) transparent;"
  >
    <!-- Toggle button -->
    <div class="flex items-center justify-between px-3 py-3" :class="collapsed ? 'flex-col gap-3' : ''">
      <RouterLink v-if="!collapsed" to="/" class="flex items-center gap-3 rounded-2xl px-1 py-2">
        <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl bg-spotify text-xl text-black">
          <i aria-hidden="true" class="pi pi-volume-up" />
        </div>
        <div>
          <div class="text-lg font-black text-white">Music App</div>
        </div>
      </RouterLink>
      <button
        type="button"
        class="flex items-center justify-center rounded-lg text-slate-500 hover:text-white hover:bg-white/10 transition-all duration-200"
        :class="collapsed ? 'h-10 w-10' : 'h-8 w-8 shrink-0'"
          :aria-label="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          :title="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          @click="toggleCollapsed"
      >
        <i aria-hidden="true" :class="collapsed ? 'pi pi-chevron-right' : 'pi pi-chevron-left'" class="text-xs" />
      </button>
    </div>

    <template v-if="!collapsed">
      <nav aria-label="Browse" class="mt-2 space-y-1">
        <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">Browse</p>
        <RouterLink
          v-for="item in mainNav" :key="item.to" :to="item.to"
          class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
          :class="isActive(item.to) ? 'bg-white/10 text-white' : ''"
          :aria-current="isActive(item.to) ? 'page' : undefined"
        >
          <i aria-hidden="true" :class="item.icon" class="text-lg" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
      <nav aria-label="Library" class="mt-6 space-y-1">
        <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">Library</p>
        <RouterLink
          v-for="item in libraryNav" :key="item.to" :to="item.to"
          class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
          :class="isActive(item.to) ? 'bg-white/10 text-white' : ''"
          :aria-current="isActive(item.to) ? 'page' : undefined"
        >
          <i aria-hidden="true" :class="item.icon" class="text-lg" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
      <nav aria-label="Social" class="mt-6 space-y-1">
        <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">Social</p>
        <RouterLink
          v-for="item in socialNav" :key="item.to" :to="item.to"
          class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
          :class="isActive(item.to) ? 'bg-white/10 text-white' : ''"
          :aria-current="isActive(item.to) ? 'page' : undefined"
        >
          <i aria-hidden="true" :class="item.icon" class="text-lg" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
      <nav aria-label="More" class="mt-6 space-y-1">
        <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">More</p>
        <RouterLink
          v-for="item in moreNav" :key="item.to" :to="item.to"
          class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
        :class="isActive(item.to) ? 'bg-white/10 text-white' : ''"
        :aria-current="isActive(item.to) ? 'page' : undefined"
      >
          <i aria-hidden="true" :class="item.icon" class="text-lg" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
    </template>
    <!-- Collapsed nav: icons only -->
    <nav v-else aria-label="Main navigation" class="mt-2 flex flex-col items-center gap-1 px-1">
      <RouterLink
        v-for="item in allNav" :key="item.to" :to="item.to"
        class="flex h-10 w-10 items-center justify-center rounded-xl text-lg transition hover:bg-white/8"
        :class="isActive(item.to) ? 'text-white bg-white/10' : 'text-slate-400 hover:text-white'"
        :title="item.label"
        :aria-current="isActive(item.to) ? 'page' : undefined"
      >
        <i aria-hidden="true" :class="item.icon" />
      </RouterLink>
    </nav>

    <!-- Spacer to prevent content from being hidden behind player bar -->
    <div v-if="playerStore.currentTrack" class="h-32" />
    <div v-else class="h-4" />
  </nav>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { usePlayerStore } from '@/stores/player'
import { useUserAuthStore } from '@/stores'

const route = useRoute()
const playerStore = usePlayerStore()
const authStore = useUserAuthStore()

const collapsed = ref(localStorage.getItem('sidebar-collapsed') === 'true')
function toggleCollapsed() {
  collapsed.value = collapsed.value!
  localStorage.setItem('sidebar-collapsed', String(collapsed.value))
  window.dispatchEvent(new CustomEvent('sidebar-collapse', { detail: collapsed.value }))
}

const mainNav = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'Search', icon: 'pi pi-search', to: '/search' },
  { label: 'Recommendations', icon: 'pi pi-star', to: '/recommendations' },
  { label: 'Music Videos', icon: 'pi pi-video', to: '/videos' },
]

const libraryNav = [
  { label: 'Library', icon: 'pi pi-bookmark', to: '/library' },
]

const socialNav = [
  { label: 'Explore', icon: 'pi pi-compass', to: '/explore' },
  { label: 'Social Hub', icon: 'pi pi-users', to: '/social' },
  { label: 'Notifications', icon: 'pi pi-bell', to: '/notifications' },
]

const moreNav = computed(() => {
  const items = [
    { label: 'Profile', icon: 'pi pi-user', to: '/profile' },
    { label: 'Settings', icon: 'pi pi-cog', to: '/settings' },
    { label: 'AI Mood Explorer', icon: 'pi pi-magic', to: '/ai/mood-explorer' },
    { label: 'AI Playlist Generator', icon: 'pi pi-sync', to: '/ai/playlist-generator' },
    { label: 'Subscription', icon: 'pi pi-credit-card', to: '/subscription' },
    { label: 'Gamification', icon: 'pi pi-trophy', to: '/gamification' },
    { label: 'Contributions', icon: 'pi pi-cloud-upload', to: '/contributions' },
    { label: 'Creator Dashboard', icon: 'pi pi-chart-bar', to: '/creator-dashboard' },
  ]
  // Admin users get a direct link to the admin panel
  if (authStore.isAdmin) {
    items.push({ label: 'Admin Panel', icon: 'pi pi-shield', to: '/admin' })
  }
  return items
})

const allNav = computed(() => [...mainNav, ...libraryNav, ...socialNav, ...moreNav.value])

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>
