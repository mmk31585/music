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
      <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
        Browse
      </p>

      <RouterLink
        v-for="item in mainNav"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-white"
        :class="isActive(item.to) ? 'bg-white/[0.10] text-white' : ''"
      >
        <i :class="item.icon" class="text-lg" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <nav class="mt-6 space-y-1">
      <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
        Library
      </p>

      <RouterLink
        v-for="item in libraryNav"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-white"
        :class="isActive(item.to) ? 'bg-white/[0.10] text-white' : ''"
      >
        <i :class="item.icon" class="text-lg" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <nav class="mt-6 space-y-1">
      <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
        Social
      </p>

      <RouterLink
        v-for="item in socialNav"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-white"
        :class="isActive(item.to) ? 'bg-white/[0.10] text-white' : ''"
      >
        <i :class="item.icon" class="text-lg" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <nav class="mt-6 space-y-1">
      <p class="px-4 pb-1 pt-2 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
        More
      </p>

      <RouterLink
        v-for="item in moreNav"
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

const mainNav = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'Discover', icon: 'pi pi-compass', to: '/discover' },
  { label: 'Search', icon: 'pi pi-search', to: '/search' },
  { label: 'Recommendations', icon: 'pi pi-star', to: '/recommendations' },
]

const libraryNav = [
  { label: 'Library', icon: 'pi pi-bookmark', to: '/library' },
  { label: 'Playlists', icon: 'pi pi-list', to: '/playlists' },
  { label: 'Recently Played', icon: 'pi pi-history', to: '/recently-played' },
]

const socialNav = [
  { label: 'Social Hub', icon: 'pi pi-users', to: '/social' },
  { label: 'Notifications', icon: 'pi pi-bell', to: '/notifications' },
]

const moreNav = [
  { label: 'Profile', icon: 'pi pi-user', to: '/profile' },
  { label: 'Settings', icon: 'pi pi-cog', to: '/settings' },
  { label: 'AI Mood Explorer', icon: 'pi pi-magic', to: '/ai/mood-explorer' },
  { label: 'AI Playlist Generator', icon: 'pi pi-sync', to: '/ai/playlist-generator' },
  { label: 'Subscription', icon: 'pi pi-credit-card', to: '/subscription' },
  { label: 'Gamification', icon: 'pi pi-trophy', to: '/gamification' },
  { label: 'Contributions', icon: 'pi pi-cloud-upload', to: '/contributions' },
  { label: 'Creator Dashboard', icon: 'pi pi-chart-bar', to: '/creator-dashboard' },
]

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>
