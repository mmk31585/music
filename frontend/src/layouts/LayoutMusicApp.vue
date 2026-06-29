<template>
  <div class="relative">
    <a href="#main-content" class="skip-link">Skip to main content</a>

    <!-- Offline banner -->
    <div
      v-if="!isOnline"
      role="alert"
      class="fixed top-0 left-0 right-0 z-9999 flex items-center justify-center gap-2 bg-red-600/90 px-4 py-2 text-sm font-medium text-white backdrop-blur-xs"
      style="padding-top: max(0.5rem, env(safe-area-inset-top, 0.5rem))"
    >
      <i class="pi pi-wifi text-xs" aria-hidden="true" />
      <span>You are offline. Some features may be unavailable.</span>
    </div>

    <div
      class="flex h-screen overflow-hidden bg-transparent text-white"
      :dir="rtlDir"
    >
    <!-- ── right Sidebar ── -->
    <MusicSidebar />

    <!-- ── Main Content Area (scrolls independently) ── -->
    <div class="flex min-w-0 flex-1 flex-col overflow-hidden">
      <MusicAppHeader
        :page-title="pageTitle"
        :unread-count="unreadCount"
        @toggle-mobile="mobileOpen = true"
        @toggle-search="searchOpen = true"
      />

      <main
        id="main-content"
        class="flex-1 overflow-y-auto scroll-smooth scroll-bar"
        :class="mainPadding"
      >
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <KeepAlive :max="3">
              <component :is="Component" />
            </KeepAlive>
          </Transition>
        </RouterView>
      </main>
    </div>

    <!-- ── Right Sticky Pane (persistent, categorized) ── -->
    <MusicRightPane
      @toggle-fullscreen="fullscreenOpen = !fullscreenOpen"
      @toggle-queue-overlay="showQueue = !showQueue"
    />

    <Transition name="fade">
      <div
        v-if="mobileOpen"
        class="fixed inset-0 z-50 bg-black/70 backdrop-blur-xs lg:hidden"
        :aria-hidden="!mobileOpen"
        role="button"
        tabindex="0"
        @click="mobileOpen = false"
        @keydown.enter="mobileOpen = false"
        @keydown.space.prevent="mobileOpen = false"
      >
        <div class="flex h-full w-80 max-w-[85vw] flex-col bg-black p-4" @click.stop>
          <div class="mb-4 flex items-center justify-between">
            <RouterLink to="/" class="flex items-center gap-3" @click="mobileOpen = false">
              <div
                class="flex h-10 w-10 items-center justify-center rounded-2xl bg-spotify text-black"
              >
                <i aria-hidden="true" class="pi pi-volume-up" />
              </div>
              <span class="font-black text-white">Music App</span>
            </RouterLink>

            <button
              type="button"
              aria-label="Close navigation menu"
              class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white"
              @click="mobileOpen = false"
            >
              <i aria-hidden="true" class="pi pi-times" />
            </button>
          </div>

          <nav role="navigation" aria-label="Mobile navigation" class="flex-1 space-y-1 overflow-y-auto">
            <p class="px-4 pb-1 pt-4 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
              Browse
            </p>
            <RouterLink
              v-for="item in browseItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
              :class="activeNavBase === item.to ? 'bg-white/10 text-white' : ''"
              @click="mobileOpen = false"
            >
              <i aria-hidden="true" :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </RouterLink>

            <p class="px-4 pb-1 pt-4 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
              Library
            </p>
            <RouterLink
              v-for="item in libraryItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
              :class="activeNavBase === item.to ? 'bg-white/10 text-white' : ''"
              @click="mobileOpen = false"
            >
              <i aria-hidden="true" :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </RouterLink>

            <p class="px-4 pb-1 pt-4 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
              Social
            </p>
            <RouterLink
              v-for="item in socialItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
              :class="activeNavBase === item.to ? 'bg-white/10 text-white' : ''"
              @click="mobileOpen = false"
            >
              <i aria-hidden="true" :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </RouterLink>

            <p class="px-4 pb-1 pt-4 text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">
              More
            </p>
            <RouterLink
              v-for="item in moreItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-white"
              :class="activeNavBase === item.to ? 'bg-white/10 text-white' : ''"
              @click="mobileOpen = false"
            >
              <i aria-hidden="true" :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </RouterLink>
          </nav>

          <div class="mt-auto space-y-2 border-t border-white/10 pt-4">
            <template v-if="store.isAuthenticated">
              <div class="flex items-center gap-3 rounded-xl px-4 py-2">
                <div
                  class="flex h-9 w-9 items-center justify-center rounded-full bg-spotify/20 text-sm font-bold text-spotify"
                >
                  {{ initials }}
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-bold text-white">{{ displayName }}</p>
                  <p class="text-xs text-slate-500">{{ isAdmin ? 'Admin' : 'Listener' }}</p>
                </div>
              </div>

              <button
                type="button"
                class="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/8 hover:text-red-400"
                @click="handleLogout"
              >
                <i aria-hidden="true" class="pi pi-sign-out text-lg" />
                <span>Log out</span>
              </button>
            </template>

            <template v-else>
              <RouterLink
                to="/auth/login"
                class="flex w-full items-center justify-center gap-2 rounded-xl bg-spotify px-4 py-3 text-sm font-bold text-black transition hover:bg-spotify-hover"
                @click="mobileOpen = false"
              >
                <i aria-hidden="true" class="pi pi-sign-in" />
                <span>Log in</span>
              </RouterLink>

              <RouterLink
                to="/auth/register"
                class="flex w-full items-center justify-center gap-2 rounded-xl border border-white/15 px-4 py-3 text-sm font-bold text-white transition hover:bg-white/8"
                @click="mobileOpen = false"
              >
                <span>Sign up</span>
              </RouterLink>
            </template>
          </div>
        </div>
      </div>
    </Transition>

    <SearchOverlay v-model:visible="searchOpen" />
    <RadioMode v-model:visible="radioVisible" :seed-id="radioSeedId" :seed-label="radioSeedLabel" />
    <PlayerRegion />
    <MobileBottomNav />
  </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, provide, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  MusicSidebar,
  SearchOverlay,
  RadioMode,
  PlayerRegion,
} from '@/components/music'
import MusicRightPane from '@/components/music/layout/MusicRightPane.vue'
import { MobileBottomNav, MusicAppHeader  } from '@/components/layouts'
import { useUserAuthStore, usePlayerStore } from '@/stores'
import { useAuth } from '@/composables/auth/useAuth'
import { client, useRTL } from '@/composables'
import { wsClient } from '@/services/socket/client'
import axios from 'axios'
import type { NotificationResponse } from '@/services/api/notification/routes'

const browseItems = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'Search', icon: 'pi pi-search', to: '/search' },
  { label: 'Recommendations', icon: 'pi pi-star', to: '/recommendations' },
]

const libraryItems = [
  { label: 'Library', icon: 'pi pi-bookmark', to: '/library' },
  { label: 'Playlists', icon: 'pi pi-list', to: '/playlists' },
  { label: 'Recently Played', icon: 'pi pi-history', to: '/recently-played' },
]

const socialItems = [
  { label: 'Explore', icon: 'pi pi-compass', to: '/explore' },
  { label: 'Social Hub', icon: 'pi pi-users', to: '/social' },
  { label: 'Notifications', icon: 'pi pi-bell', to: '/notifications' },
]

const moreItems = [
  { label: 'Profile', icon: 'pi pi-user', to: '/profile' },
  { label: 'Settings', icon: 'pi pi-cog', to: '/settings' },
  { label: 'AI Mood Explorer', icon: 'pi pi-magic', to: '/ai/mood-explorer' },
  { label: 'AI Playlist Generator', icon: 'pi pi-sync', to: '/ai/playlist-generator' },
  { label: 'Subscription', icon: 'pi pi-credit-card', to: '/subscription' },
  { label: 'Gamification', icon: 'pi pi-trophy', to: '/gamification' },
  { label: 'Contributions', icon: 'pi pi-cloud-upload', to: '/contributions' },
  { label: 'Creator Dashboard', icon: 'pi pi-chart-bar', to: '/creator-dashboard' },
]

const pageTitleMap: Record<string, string> = {
  '/': 'Home',
  '/explore': 'Explore',
  '/search': 'Search',
  '/recommendations': 'Recommendations',
  '/library': 'Library',
  '/playlists': 'Playlists',
  '/recently-played': 'Recently Played',
  '/notifications': 'Notifications',
  '/social': 'Social Hub',
  '/profile': 'Profile',
  '/settings': 'Settings',
  '/subscription': 'Subscription',
  '/contributions': 'Contributions',
  '/creator-dashboard': 'Creator Dashboard',
  '/gamification': 'Gamification',
  '/ai/mood-explorer': 'AI Mood Explorer',
  '/ai/playlist-generator': 'AI Playlist Generator',
}

const route = useRoute()
const store = useUserAuthStore()
const playerStore = usePlayerStore()
const { logout } = useAuth()
const mobileOpen = ref(false)
const fullscreenOpen = ref(false)
const showQueue = ref(false)
const searchOpen = ref(false)

const radioVisible = ref(false)
const radioSeedId = ref('')
const radioSeedLabel = ref('')

function openRadio(trackId: string, seedLabel?: string) {
  radioSeedId.value = trackId
  radioSeedLabel.value = seedLabel || ''
  radioVisible.value = true
}

provide('openRadio', openRadio)

const { dir: rtlDir } = useRTL()

const barCollapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')
const isOnline = ref(navigator.onLine)

onMounted(() => {
  const handler = (e: Event) => {
    barCollapsed.value = (e as CustomEvent).detail
  }
  window.addEventListener('playerbar-collapse', handler)
})

const mainPadding = computed(() => {
  if (playerStore.currentTrack) {
    return barCollapsed.value ? 'pb-20 lg:pb-16' : 'pb-32 lg:pb-28'
  }
  return 'pb-24 lg:pb-0'
})
const showShortcuts = ref(false)
const unreadCount = ref(0)
let unreadInterval: ReturnType<typeof setInterval> | null = null
let unsubNotif: (() => void) | null = null

function startPolling() {
  fetchUnreadCount()
  unreadInterval = setInterval(fetchUnreadCount, 30000)
}

function stopPolling() {
  if (unreadInterval) {
    clearInterval(unreadInterval)
    unreadInterval = null
  }
}

function onVisibilityChange() {
  if (document.hidden) {
    stopPolling()
  } else if (store.isAuthenticated) {
    startPolling()
  }
}

onMounted(() => {
  if (store.isAuthenticated) {
    startPolling()
    unsubNotif = wsClient.on('notification', (msg) => {
      const n = msg.payload as NotificationResponse
      if (n && n.id && !n.isRead) {
        unreadCount.value++
      }
    })
    wsClient.connect()
  }
  document.addEventListener('visibilitychange', onVisibilityChange)
  document.addEventListener('keydown', (e) => {
    const tag = (e.target as HTMLElement)?.tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return
    if (e.key === '?' && e.ctrlKey! && e.metaKey! && e.altKey!) {
      showShortcuts.value = !showShortcuts.value
    }
  })
  const updateOnline = () => { isOnline.value = navigator.onLine }
  window.addEventListener('online', updateOnline)
  window.addEventListener('offline', updateOnline)
})

onUnmounted(() => {
  stopPolling()
  document.removeEventListener('visibilitychange', onVisibilityChange)
  unsubNotif?.()
})

async function fetchUnreadCount() {
  if (!store.isAuthenticated) return
  try {
    const res = await client.get('/notifications', { params: { limit: 1 } })
    unreadCount.value = res.data?.data?.unreadCount ?? 0
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.status === 401) {
      // User is not authenticated, silently ignore
      return
    }
    console.error('Failed to fetch unread count:', err)
  }
}

const isAdmin = computed(() => store.isAdmin)
const displayName = computed(
  () => store.user?.displayName || store.user?.username || store.user?.name || 'User',
)
const initials = computed(() => {
  const name = displayName.value || '?'
  const words = name.split(/\s+/).filter(Boolean)
  return words.length >= 2
    ? (words[0]![0]! + words[words.length - 1]![0]!).toUpperCase()
    : name.slice(0, 2).toUpperCase()
})

const allNavItems = [...browseItems, ...libraryItems, ...socialItems, ...moreItems]

const activeNavBase = computed(() => {
  const path = route.path
  for (const item of allNavItems) {
    if (item.to === '/' ? path === '/' : path.startsWith(item.to)) {
      return item.to
    }
  }
  return null
})

const pageTitle = computed(() => {
  const path = route.path
  let bestMatch = 'Music'
  let bestLength = 0

  for (const [prefix, title] of Object.entries(pageTitleMap)) {
    if (path.startsWith(prefix) && prefix.length > bestLength) {
      bestMatch = title
      bestLength = prefix.length
    }
  }

  if (path.startsWith('/playlist/')) return 'Playlist'
  if (path.startsWith('/track/')) return 'Track'
  if (path.startsWith('/album/')) return 'Album'
  if (path.startsWith('/artist/')) return 'Artist'
  if (path.startsWith('/user/')) return 'Profile'
  if (path.startsWith('/social/party/')) return 'Listening Party'
  if (path.startsWith('/social/room/')) return 'Live Room'
  if (path.startsWith('/social/club/')) return 'Music Club'
  if (path.startsWith('/recommendations/')) return 'Recommendations'

  return bestMatch
})

// Update document title when page changes
watch(pageTitle, (title) => {
  document.title = title ? `${title} — Muse` : 'Muse'
}, { immediate: true })

function handleLogout() {
  mobileOpen.value = false
  void logout()
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
