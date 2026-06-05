<template>
  <div class="min-h-screen bg-transparent text-white">
    <div class="flex min-h-screen">
      <MusicSidebar />

      <main class="min-w-0 flex-1">
        <header
          class="sticky top-0 z-30 flex h-20 items-center justify-between border-b border-white/10 bg-black/30 px-4 backdrop-blur-xl md:px-6 lg:px-8"
        >
          <div class="flex items-center gap-3">
            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/15 lg:hidden"
              @click="mobileOpen = true"
            >
              <i class="pi pi-bars" />
            </button>

            <div>
              <p class="text-xs font-semibold tracking-[0.25em] text-[#1db954] uppercase">Music</p>
              <h1 class="text-lg font-black text-white md:text-xl">
                {{ pageTitle }}
              </h1>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <RouterLink
              to="/search"
              class="hidden rounded-full bg-white/10 px-4 py-2 text-sm font-bold text-white transition hover:bg-white/15 md:inline-flex"
            >
              <i class="pi pi-search mr-2" />
              Search
            </RouterLink>

            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-black"
            >
              <i class="pi pi-user" />
            </button>
          </div>
        </header>

        <RouterView />
      </main>
    </div>

    <Transition name="fade">
      <div
        v-if="mobileOpen"
        class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm lg:hidden"
        @click="mobileOpen = false"
      >
        <div class="h-full w-80 max-w-[85vw] bg-black p-4" @click.stop>
          <div class="mb-4 flex items-center justify-between">
            <RouterLink to="/" class="flex items-center gap-3" @click="mobileOpen = false">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-2xl bg-[#1db954] text-black"
              >
                <i class="pi pi-volume-up" />
              </div>

              <span class="font-black text-white">Music App</span>
            </RouterLink>

            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white"
              @click="mobileOpen = false"
            >
              <i class="pi pi-times" />
            </button>
          </div>

          <nav class="space-y-1">
            <RouterLink
              v-for="item in navItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-white"
              :class="isActive(item.to) ? 'bg-white/[0.10] text-white' : ''"
              @click="mobileOpen = false"
            >
              <i :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </RouterLink>
          </nav>
        </div>
      </div>
    </Transition>

    <NowPlayingBar />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { MusicSidebar, NowPlayingBar } from '@/components/music'

const route = useRoute()
const mobileOpen = ref(false)

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

const pageTitle = computed(() => {
  const match = navItems.find((item) => {
    if (item.to === '/') return route.path === '/'
    return route.path.startsWith(item.to)
  })

  return match?.label || 'Music'
})

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 160ms ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
