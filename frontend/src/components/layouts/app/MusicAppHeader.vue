<template>
  <header
    ref="headerRef"
    role="banner"
    class="sticky top-0 z-30 border-b border-white/10 backdrop-blur-xl transition-all duration-300"
    :class="[
      isScrolled
        ? 'bg-black/70 shadow-lg shadow-black/20'
        : 'bg-black/30',
    ]"
  >
    <div class="flex h-20 items-center justify-between px-4 md:px-6 lg:px-8">
      <!-- ── Left: Hamburger + Title ── -->
      <div class="flex items-center gap-3">
        <button
          type="button"
          aria-label="Open navigation menu"
          class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/15 lg:hidden"
          @click="$emit('toggle-mobile')"
        >
          <i aria-hidden="true" class="pi pi-bars" />
        </button>

        <div>
          <p class="text-xs font-semibold tracking-[0.25em] text-spotify uppercase">Music</p>
          <h1 class="text-lg font-black text-white md:text-xl">
            {{ pageTitle }}
          </h1>
        </div>
      </div>

      <!-- ── Right: Actions ── -->
      <div class="flex items-center gap-2 md:gap-3">
        <!-- Search -->
        <button
          type="button"
          aria-label="Main navigation"
          class="hidden items-center gap-2 rounded-full bg-white/10 px-4 py-2 text-sm font-bold text-white transition hover:bg-white/15 md:inline-flex"
          @click="$emit('toggle-search')"
        >
          <i aria-hidden="true" class="pi pi-search" />
          Search
          <kbd
            class="rounded-md border border-white/10 bg-white/10 px-1.5 py-0.5 text-[10px] text-slate-400"
          >Ctrl+K</kbd>
        </button>

        <!-- ════════════════════════════════════════ -->
        <!-- AUTHENTICATED                           -->
        <!-- ════════════════════════════════════════ -->
        <template v-if="isAuthenticated">
          <!-- Notifications -->
          <RouterLink
            to="/notifications"
            aria-label="Notifications"
            class="relative flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/15"
          >
            <i aria-hidden="true" class="pi pi-bell" />
            <span
              v-if="unreadCount > 0"
              class="absolute -top-0.5 -right-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-spotify px-1 text-[10px] font-bold text-black"
            >
              {{ unreadCount > 99 ? '99+' : unreadCount }}
            </span>
          </RouterLink>

          <!-- User avatar + dropdown -->
          <div ref="userMenuContainer" class="relative">
            <button
              type="button"
              ref="userButtonRef"
              aria-label="User menu"
              aria-haspopup="true"
              :aria-expanded="showUserMenu"
              class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-full bg-spotify text-sm font-bold text-black ring-2 ring-transparent transition-all duration-200 hover:bg-spotify-hover hover:ring-spotify/30 focus-visible:ring-spotify/50"
              :class="{ 'ring-spotify/40': showUserMenu }"
              @click="toggleUserMenu"
            >
              <img
                v-if="user?.avatarUrl"
                :src="user.avatarUrl"
                :alt="displayName"
                class="h-full w-full object-cover"
                @error="onAvatarError"
              />
              <span v-else>{{ initials }}</span>
            </button>

            <!-- Dropdown Menu -->
            <Transition name="user-menu">
              <div
                v-if="showUserMenu"
                class="absolute left-0 top-full z-50 mt-3 w-64 origin-top-right overflow-hidden rounded-2xl border border-white/8 bg-[#1a1a2e]/95 shadow-2xl shadow-black/40 backdrop-blur-2xl"
                @click.stop
              >
                <!-- User info header -->
                <div class="border-b border-white/6 px-4 py-4">
                  <div class="flex items-center gap-3">
                    <div
                      class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-spotify text-sm font-bold text-black"
                    >
                      <img
                        v-if="user?.avatarUrl"
                        :src="user.avatarUrl"
                        :alt="displayName"
                        class="h-full w-full rounded-full object-cover"
                        @error="onAvatarError"
                      />
                      <span v-else>{{ initials }}</span>
                    </div>
                    <div class="min-w-0">
                      <p class="truncate text-sm font-bold text-white">{{ displayName }}</p>
                      <p class="truncate text-xs text-white/40">{{ user?.email || '' }}</p>
                    </div>
                  </div>
                  <div
                    v-if="user?.role"
                    class="mt-2 flex items-center gap-1.5"
                  >
                    <span
                      class="rounded-full bg-spotify/10 px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider text-spotify"
                    >{{ user.role }}</span>
                  </div>
                </div>

                <!-- Menu items -->
                <div class="py-1.5">
                  <RouterLink
                    v-for="item in menuItems"
                    :key="item.to"
                    :to="item.to"
                    class="flex items-center gap-3 px-4 py-2.5 text-sm text-white/70 transition hover:bg-white/6 hover:text-white"
                    @click="closeUserMenu"
                  >
                    <span
                      class="flex h-7 w-7 items-center justify-center rounded-lg"
                      :class="item.iconBg"
                    >
                      <i aria-hidden="true" :class="item.icon" class="text-xs" />
                    </span>
                    <span>{{ item.label }}</span>
                    <span
                      v-if="item.badge"
                      class="ml-auto rounded-full bg-spotify/15 px-2 py-0.5 text-[10px] font-bold text-spotify"
                    >
                      {{ item.badge }}
                    </span>
                  </RouterLink>
                </div>

                <!-- Divider + Logout -->
                <div class="border-t border-white/6 py-1.5">
                  <button
                    type="button"
                    class="flex w-full items-center gap-3 px-4 py-2.5 text-sm text-red-400/70 transition hover:bg-red-500/10 hover:text-red-400"
                    @click="handleLogout"
                  >
                    <span class="flex h-7 w-7 items-center justify-center rounded-lg bg-red-500/10">
                      <i aria-hidden="true" class="pi pi-sign-out text-xs" />
                    </span>
                    <span>Sign out</span>
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </template>

        <!-- ════════════════════════════════════════ -->
        <!-- GUEST                                   -->
        <!-- ════════════════════════════════════════ -->
        <template v-else>
          <RouterLink
            to="/auth/login"
            class="rounded-full bg-spotify px-5 py-2 text-sm font-bold text-black transition hover:bg-spotify-hover"
          >
            Log in
          </RouterLink>

          <RouterLink
            to="/auth/register"
            class="hidden rounded-full border border-white/15 px-5 py-2 text-sm font-bold text-white transition hover:bg-white/8 md:inline-flex"
          >
            Sign up
          </RouterLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserAuthStore, usePlayerStore } from '@/stores'
import { wsClient } from '@/services/socket/client'

defineEmits<{
  'toggle-mobile': []
  'toggle-search': []
}>()

defineProps<{
  pageTitle: string
  unreadCount: number
}>()

const router = useRouter()
const store = useUserAuthStore()
const headerRef = ref<HTMLElement | null>(null)
const userButtonRef = ref<HTMLElement | null>(null)
const userMenuContainer = ref<HTMLElement | null>(null)
const isScrolled = ref(false)
const showUserMenu = ref(false)

const isAuthenticated = computed(() => store.isAuthenticated)
const user = computed(() => store.user)

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

const menuItems = computed<Array<{ label: string; icon: string; iconBg: string; to: string; badge?: string }>>(() => {
  const items: Array<{ label: string; icon: string; iconBg: string; to: string; badge?: string }> = [
    {
      label: 'Profile',
      icon: 'pi pi-user',
      iconBg: 'bg-white/6 text-white/60',
      to: '/profile',
    },
    {
      label: 'Settings',
      icon: 'pi pi-cog',
      iconBg: 'bg-white/6 text-white/60',
      to: '/settings',
    },
    {
      label: 'Listening Stats',
      icon: 'pi pi-chart-bar',
      iconBg: 'bg-white/6 text-white/60',
      to: '/stats',
    },
  ]

  if (store.isAdmin) {
    items.push({
      label: 'Admin Panel',
      icon: 'pi pi-shield',
      iconBg: 'bg-amber-500/10 text-amber-400',
      to: '/admin',
      badge: 'Admin',
    })
  }

  return items
})

function toggleUserMenu() {
  showUserMenu.value = !showUserMenu.value
}

function closeUserMenu() {
  showUserMenu.value = false
}

function onAvatarError(e: Event) {
  const img = e.currentTarget as HTMLImageElement
  img.style.display = 'none'
}

async function handleLogout() {
  closeUserMenu()
  await store.logout()
  wsClient.disconnect()
  usePlayerStore().$reset()
  router.push('/')
}

// ── Click outside to close ──
function onClickOutside(e: MouseEvent) {
  if (
    showUserMenu.value &&
    userMenuContainer.value &&
    userMenuContainer.value.contains!(e.target as Node)
  ) {
    closeUserMenu()
  }
}

// ── Escape key to close ──
function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape' && showUserMenu.value) {
    closeUserMenu()
    userButtonRef.value?.focus()
  }
}

// ── Scroll observer for header bg ──
let scrollObserver: IntersectionObserver | null = null

onMounted(() => {
  document.addEventListener('click', onClickOutside)
  document.addEventListener('keydown', onKeyDown)

  if (!headerRef.value) return
  const sentinel = document.createElement('div')
  sentinel.style.position = 'absolute'
  sentinel.style.top = '0'
  sentinel.style.left = '0'
  sentinel.style.width = '1px'
  sentinel.style.height = '1px'
  sentinel.style.pointerEvents = 'none'
  document.body.prepend(sentinel)

  scrollObserver = new IntersectionObserver(
    ([entry]) => {
      isScrolled.value = entry!.isIntersecting!
    },
    { threshold: 0 },
  )
  scrollObserver.observe(sentinel)
})

onUnmounted(() => {
  document.removeEventListener('click', onClickOutside)
  document.removeEventListener('keydown', onKeyDown)
  scrollObserver?.disconnect()
})
</script>

<style scoped>
/* ── Dropdown animation ── */
.user-menu-enter-active {
  transition: opacity 0.2s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.user-menu-leave-active {
  transition: opacity 0.12s ease-in,
              transform 0.12s ease-in;
}
.user-menu-enter-from {
  opacity: 0;
  transform: translateY(-6px) scale(0.96);
}
.user-menu-leave-to {
  opacity: 0;
  transform: translateY(-4px) scale(0.96);
}
</style>
