<template>
  <aside
    class="hidden h-screen w-72 shrink-0 border-r border-white/10 bg-black/40 p-4 backdrop-blur-xl lg:block"
  >
    <RouterLink to="/" class="flex items-center gap-3 rounded-2xl px-3 py-4">
      <div
        class="flex h-11 w-11 items-center justify-center rounded-2xl bg-[#1db954] text-xl text-black"
      >
        <i class="pi pi-volume-up" />
      </div>

      <div>
        <div class="text-lg font-black text-white">Music App</div>
        <div class="text-xs font-medium text-slate-400">Stream everything</div>
      </div>
    </RouterLink>

    <nav class="mt-6 space-y-1">
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-white"
        :class="isActive(item.to) ? 'bg-white/[0.10] text-white' : ''"
      >
        <i :class="item.icon" class="text-lg" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <div class="mt-8 rounded-3xl border border-white/10 bg-white/[0.05] p-5">
      <div
        class="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#1db954]/20 text-[#1db954]"
      >
        <i class="pi pi-sparkles" />
      </div>

      <h3 class="mt-4 text-sm font-black text-white">Discover weekly</h3>

      <p class="mt-2 text-xs leading-5 text-slate-400">
        Fresh tracks and personalized recommendations are waiting for you.
      </p>

      <RouterLink
        to="/recommendations"
        class="mt-4 inline-flex rounded-full bg-white px-4 py-2 text-xs font-black text-black transition hover:scale-105"
      >
        Explore
      </RouterLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const navItems = [
  {
    label: 'Home',
    icon: 'pi pi-home',
    to: '/',
  },
  {
    label: 'Search',
    icon: 'pi pi-search',
    to: '/search',
  },
  {
    label: 'Recommendations',
    icon: 'pi pi-star',
    to: '/recommendations',
  },
  {
    label: 'Library',
    icon: 'pi pi-bookmark',
    to: '/library',
  },
  {
    label: 'Playlists',
    icon: 'pi pi-list',
    to: '/playlists',
  },
  {
    label: 'Recently Played',
    icon: 'pi pi-history',
    to: '/recently-played',
  },
]

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>
