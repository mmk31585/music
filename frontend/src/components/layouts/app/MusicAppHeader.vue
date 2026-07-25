<template>
  <header
    ref="headerRef"
    role="banner"
    class="sticky top-0 z-30 border-b border-border-default backdrop-blur-xl transition-all duration-300"
    :class="[isScrolled ? 'bg-surface-elevated shadow-lg' : 'bg-transparent']"
  >
    <div class="flex h-[72px] items-center justify-between px-6 md:px-8 lg:px-10">
      <div class="flex items-center gap-3">
        <button
          type="button"
          aria-label="Open navigation menu"
          class="flex h-10 w-10 items-center justify-center rounded-xl bg-surface-overlay text-primary transition hover:bg-surface-active lg:hidden"
          @click="$emit('toggle-mobile')"
        >
          <Menu aria-hidden="true" class="text-sm"  />
        </button>
        <h1 class="text-xl font-bold tracking-tight text-primary">{{ pageTitle }}</h1>
      </div>

      <div class="hidden flex-1 items-center justify-center gap-2 px-8 md:flex">
        <RouterLink
          v-if="!isHome"
          to="/"
          aria-label="Go to home"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-secondary transition-all duration-150 hover:bg-surface-overlay hover:text-primary active:scale-90"
        >
          <ChevronLeft :size="18" />
        </RouterLink>
        <button
          type="button"
          class="group flex w-full max-w-md items-center gap-3 rounded-xl border border-border-default bg-surface-overlay px-4 py-2.5 text-sm text-muted transition-all duration-200 hover:border-border-strong hover:bg-surface-active hover:text-secondary focus-visible:border-accent focus-visible:outline-none"
          @click="$emit('toggle-search')"
        >
          <Search aria-hidden="true" class="text-accent text-sm transition-colors group-hover:text-accent"  />
          <span>Search songs, artists, albums...</span>
          <kbd class="ml-auto rounded-lg border border-border-default bg-surface-active px-2 py-0.5 text-[11px] font-medium text-muted">Ctrl+K</kbd>
        </button>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          aria-label="Search"
          class="flex h-10 w-10 items-center justify-center rounded-xl bg-surface-overlay text-secondary transition hover:bg-surface-active hover:text-primary md:hidden"
          @click="$emit('toggle-search')"
        >
          <Search aria-hidden="true" class="text-sm"  />
        </button>

        <template v-if="isAuthenticated">
          <RouterLink
            to="/notifications"
            aria-label="Notifications"
            class="relative flex h-10 w-10 items-center justify-center rounded-xl bg-surface-overlay text-secondary transition hover:bg-surface-active hover:text-primary"
          >
            <Bell aria-hidden="true" class="text-sm"  />
            <span
              v-if="unreadCount > 0"
              class="absolute -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-accent px-1 text-[10px] font-bold text-accent-text"
              :class="isRTL ? '-left-0.5' : '-right-0.5'"
            >
              {{ unreadCount > 99 ? '99+' : unreadCount }}
            </span>
          </RouterLink>

          <div ref="userMenuContainer" class="relative">
            <button
              type="button"
              ref="userButtonRef"
              aria-label="User menu"
              aria-haspopup="true"
              :aria-expanded="showUserMenu"
              class="flex h-9 w-9 items-center justify-center overflow-hidden rounded-xl bg-accent text-sm font-bold text-accent-text ring-2 ring-transparent transition-all duration-200 hover:bg-accent-hover focus-visible:ring-accent/50"
              :class="{ 'ring-accent/40': showUserMenu }"
              @click="showUserMenu = !showUserMenu"
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

            <Transition name="menu-pop">
              <ProfileMenu
                v-if="showUserMenu"
                :user="user"
                :display-name="displayName"
                :initials="initials"
                @close="showUserMenu = false"
              />
            </Transition>
          </div>
        </template>

        <template v-else>
          <RouterLink
            to="/auth/login"
            class="rounded-xl bg-accent px-5 py-2 text-sm font-bold text-accent-text transition hover:bg-accent-hover"
          >
            Log in
          </RouterLink>
          <RouterLink
            to="/auth/register"
            class="hidden rounded-xl border border-border-strong px-5 py-2 text-sm font-bold text-primary transition hover:bg-surface-overlay md:inline-flex"
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
import { useRoute } from 'vue-router'
import { Bell, ChevronLeft, Menu, Search } from 'lucide-vue-next'
import { useRTL } from '@/composables'
import { useUserAuthStore } from '@/stores'
import ProfileMenu from './ProfileMenu.vue'

const { isRTL } = useRTL()
const route = useRoute()

defineEmits<{
  'toggle-mobile': []
  'toggle-search': []
}>()

defineProps<{
  pageTitle: string
  unreadCount: number
}>()

const isHome = computed(() => route.path === '/')

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

function onAvatarError(e: Event) {
  const img = e.currentTarget as HTMLImageElement
  img.style.display = 'none'
}

function onClickOutside(e: MouseEvent) {
  if (
    showUserMenu.value &&
    userMenuContainer.value &&
    !userMenuContainer.value.contains(e.target as Node)
  ) {
    showUserMenu.value = false
  }
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape' && showUserMenu.value) {
    showUserMenu.value = false
    userButtonRef.value?.focus()
  }
}

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
      isScrolled.value = !entry!.isIntersecting
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
.menu-pop-enter-active {
  transition: opacity 0.15s ease-out, transform 0.15s cubic-bezier(0.16, 1, 0.3, 1);
}
.menu-pop-leave-active {
  transition: opacity 0.1s ease-in, transform 0.1s ease-in;
}
.menu-pop-enter-from {
  opacity: 0;
  transform: translateY(-4px) scale(0.96);
}
.menu-pop-leave-to {
  opacity: 0;
  transform: translateY(-2px) scale(0.97);
}
</style>
