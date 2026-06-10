<template>
  <a href="#main-content" class="skip-link">Skip to main content</a>

  <div class="min-h-screen bg-transparent text-white">
    <div class="flex min-h-screen">
      <MusicSidebar />

      <main id="main-content" class="min-w-0 flex-1">
        <MusicAppHeader
          :page-title="pageTitle"
          :unread-count="unreadCount"
          @toggle-mobile="mobileOpen = true"
          @toggle-search="searchOpen = true"
        />

        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <KeepAlive :max="3">
              <component :is="Component" />
            </KeepAlive>
          </Transition>
        </RouterView>
      </main>
    </div>

    <Transition name="fade">
      <div
        v-if="mobileOpen"
        class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm lg:hidden"
        :aria-hidden="!mobileOpen"
        @click="mobileOpen = false"
      >
        <div class="flex h-full w-80 max-w-[85vw] flex-col bg-black p-4" @click.stop>
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
              aria-label="Close navigation menu"
              class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white"
              @click="mobileOpen = false"
            >
              <i class="pi pi-times" />
            </button>
          </div>

          <nav role="navigation" aria-label="Mobile navigation" class="space-y-1">
            <RouterLink
              v-for="item in navItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-white"
              :class="activeNavBase === item.to ? 'bg-white/[0.10] text-white' : ''"
              @click="mobileOpen = false"
            >
              <i :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </RouterLink>
          </nav>

          <div class="mt-auto space-y-2 border-t border-white/10 pt-4">
            <template v-if="store.isAuthenticated">
              <div class="flex items-center gap-3 rounded-xl px-4 py-2">
                <div
                  class="flex h-9 w-9 items-center justify-center rounded-full bg-[#1db954]/20 text-sm font-bold text-[#1db954]"
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
                class="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-sm font-bold text-slate-400 transition hover:bg-white/[0.08] hover:text-red-400"
                @click="handleLogout"
              >
                <i class="pi pi-sign-out text-lg" />
                <span>Log out</span>
              </button>
            </template>

            <template v-else>
              <RouterLink
                to="/auth/login"
                class="flex w-full items-center justify-center gap-2 rounded-xl bg-[#1db954] px-4 py-3 text-sm font-bold text-black transition hover:bg-[#1ed760]"
                @click="mobileOpen = false"
              >
                <i class="pi pi-sign-in" />
                <span>Log in</span>
              </RouterLink>

              <RouterLink
                to="/auth/register"
                class="flex w-full items-center justify-center gap-2 rounded-xl border border-white/15 px-4 py-3 text-sm font-bold text-white transition hover:bg-white/[0.08]"
                @click="mobileOpen = false"
              >
                <span>Sign up</span>
              </RouterLink>
            </template>
          </div>
        </div>
      </div>
    </Transition>

    <NowPlayingBar
      @toggle-fullscreen="fullscreenOpen = !fullscreenOpen"
      @toggle-queue="showQueue = !showQueue"
      @toggle-lyrics="onToggleLyrics"
      @toggle-mobile-sheet="mobileSheetOpen = !mobileSheetOpen"
    />
    <SearchOverlay v-model:visible="searchOpen" />
    <FullscreenPlayer v-if="!ffEnabled" v-model:visible="fullscreenOpen" :initial-tab="playerInitialTab" />
    <ExpandedPlayer v-else v-model:visible="fullscreenOpen" />
    <QueuePanel v-model:visible="showQueue" />
    <MobileBottomSheet v-model:visible="mobileSheetOpen" @open-fullscreen="fullscreenOpen = true" />
    <KeyboardShortcuts v-model:visible="showShortcuts" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  MusicSidebar,
  NowPlayingBar,
  SearchOverlay,
  FullscreenPlayer,
  ExpandedPlayer,
  QueuePanel,
  MobileBottomSheet,
  KeyboardShortcuts,
} from '@/components/music'
import MusicAppHeader from '@/components/layouts/MusicAppHeader.vue'
import { useUserAuthStore, useFeatureFlagsStore } from '@/stores'
import type { FeatureFlagKey } from '@/services/api/feature-flags'
import { useAuth } from '@/composables/auth/useAuth'
import { client } from '@/composables'
import { wsClient } from '@/services/socket/client'
import type { NotificationResponse } from '@/services/api/notification/routes'

const allNavItems = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'Search', icon: 'pi pi-search', to: '/search', featureFlag: 'search' },
  { label: 'Recommendations', icon: 'pi pi-star', to: '/recommendations', featureFlag: 'recommendation' },
  { label: 'Library', icon: 'pi pi-bookmark', to: '/library' },
  { label: 'Playlists', icon: 'pi pi-list', to: '/playlists' },
  { label: 'Recently Played', icon: 'pi pi-history', to: '/recently-played' },
]

const route = useRoute()
const store = useUserAuthStore()
const ff = useFeatureFlagsStore()
const { logout } = useAuth()
const mobileOpen = ref(false)
const searchOpen = ref(false)
const fullscreenOpen = ref(false)
const showQueue = ref(false)
const mobileSheetOpen = ref(false)
const showShortcuts = ref(false)
const playerInitialTab = ref<'now-playing' | 'queue' | 'lyrics'>('now-playing')
const unreadCount = ref(0)
const ffEnabled = computed(() => ff.isEnabled('redesignedPlayer'))
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
  } else {
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
    if (e.key === '?' && !e.ctrlKey && !e.metaKey && !e.altKey) {
      showShortcuts.value = !showShortcuts.value
    }
  })
})

onUnmounted(() => {
  stopPolling()
  document.removeEventListener('visibilitychange', onVisibilityChange)
  unsubNotif?.()
})

async function fetchUnreadCount() {
  try {
    const res = await client.get('/notifications', { params: { limit: 1 } })
    unreadCount.value = res.data?.data?.unreadCount ?? 0
  } catch {
    /* ignore */
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

const navItems = computed(() =>
  allNavItems.filter(
    (item) => !item.featureFlag || ff.isEnabled(item.featureFlag as FeatureFlagKey),
  ),
)

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
  const match = allNavItems.find((item) => {
    if (item.to === '/') return route.path === '/'
    return route.path.startsWith(item.to)
  })
  return match?.label || 'Music'
})

function onToggleLyrics() {
  playerInitialTab.value = 'lyrics'
  fullscreenOpen.value = true
}

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
