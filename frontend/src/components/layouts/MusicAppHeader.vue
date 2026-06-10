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
      <div class="flex items-center gap-3">
        <button
          type="button"
          aria-label="Open navigation menu"
          class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/15 lg:hidden"
          @click="$emit('toggle-mobile')"
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
        <button
          type="button"
          aria-label="Main navigation"
          class="hidden items-center gap-2 rounded-full bg-white/10 px-4 py-2 text-sm font-bold text-white transition hover:bg-white/15 md:inline-flex"
          @click="$emit('toggle-search')"
        >
          <i class="pi pi-search" />
          Search
          <kbd class="rounded-md border border-white/10 bg-white/10 px-1.5 py-0.5 text-[10px] text-slate-400">Ctrl+K</kbd>
        </button>

        <template v-if="isAuthenticated">
          <RouterLink
            to="/notifications"
            aria-label="Notifications"
            class="relative flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/15"
          >
            <i class="pi pi-bell" />
            <span
              v-if="unreadCount > 0"
              class="absolute -top-0.5 -right-0.5 flex h-4 min-w-[16px] items-center justify-center rounded-full bg-[#1db954] px-1 text-[10px] font-bold text-black"
            >
              {{ unreadCount > 99 ? '99+' : unreadCount }}
            </span>
          </RouterLink>

          <button
            type="button"
            aria-label="User menu"
            class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-sm font-bold text-black transition hover:bg-[#1ed760]"
            :title="displayName"
          >
            {{ initials }}
          </button>
        </template>

        <template v-else>
          <RouterLink
            to="/auth/login"
            class="rounded-full bg-[#1db954] px-5 py-2 text-sm font-bold text-black transition hover:bg-[#1ed760]"
          >
            Log in
          </RouterLink>

          <RouterLink
            to="/auth/register"
            class="hidden rounded-full border border-white/15 px-5 py-2 text-sm font-bold text-white transition hover:bg-white/[0.08] md:inline-flex"
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
import { useUserAuthStore } from '@/stores'

defineEmits<{
  'toggle-mobile': []
  'toggle-search': []
}>()

defineProps<{
  pageTitle: string
  unreadCount: number
}>()

const store = useUserAuthStore()
const headerRef = ref<HTMLElement | null>(null)
const isScrolled = ref(false)

const isAuthenticated = computed(() => store.isAuthenticated)
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

let scrollObserver: IntersectionObserver | null = null

onMounted(() => {
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
  scrollObserver?.disconnect()
})
</script>
