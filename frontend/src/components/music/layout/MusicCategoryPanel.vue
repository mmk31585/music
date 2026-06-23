<template>
  <aside
    class="hidden h-screen w-56 shrink-0 border-x border-white/4 bg-black/20 backdrop-blur-2xl lg:block"
    style="backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);"
  >
    <div class="flex h-full flex-col">
      <!-- ── Panel Header ── -->
      <div class="flex shrink-0 items-center gap-3 border-b border-white/4 px-4 py-4">
        <div
          class="flex h-8 w-8 items-center justify-center rounded-xl"
          :class="categoryMeta.iconBg"
        >
          <i aria-hidden="true" :class="categoryMeta.icon" class="text-sm" />
        </div>
        <div>
          <p class="text-[10px] font-bold tracking-[0.2em] text-slate-500 uppercase">Category</p>
          <p class="text-sm font-bold text-white">{{ categoryMeta.label }}</p>
        </div>
      </div>

      <!-- ── Quick Filters ── -->
      <div v-if="quickFilters.length > 0" class="shrink-0 border-b border-white/4 px-3 py-3">
        <div class="flex flex-wrap gap-1.5">
          <button
            v-for="filter in quickFilters"
            :key="filter.label"
            type="button"
            class="rounded-lg px-3 py-1.5 text-[11px] font-bold transition"
            :class="filter.active ? 'bg-white/15 text-white' : 'bg-white/4 text-slate-400 hover:bg-white/10 hover:text-white'"
            @click="filter.action?.()"
          >
            {{ filter.label }}
          </button>
        </div>
      </div>

      <!-- ── Navigation Items ── -->
      <nav class="flex-1 space-y-0.5 overflow-y-auto px-3 pb-24 pt-3" style="scrollbar-width: thin; scrollbar-color: rgba(255,255,255,0.06) transparent;">
        <template v-for="section in navSections" :key="section.label">
          <p class="px-2 pb-1 pt-3 text-[9px] font-bold tracking-[0.2em] text-slate-600 uppercase">
            {{ section.label }}
          </p>
          <RouterLink
            v-for="item in section.items"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-2.5 rounded-xl px-3 py-2.5 text-sm font-medium transition"
            :class="isActive(item.to) ? 'bg-white/8 text-white' : 'text-slate-400 hover:bg-white/4 hover:text-white'"
            :aria-current="isActive(item.to) ? 'page' : undefined"
          >
            <i v-if="item.icon" aria-hidden="true" :class="item.icon" class="text-xs" />
            <span>{{ item.label }}</span>
            <span
              v-if="item.badge"
              class="ml-auto flex h-4 min-w-4 items-center justify-center rounded-full bg-white/10 px-1.5 text-[9px] font-bold text-slate-400"
            >
              {{ item.badge }}
            </span>
          </RouterLink>
        </template>
      </nav>

      <!-- ── Bottom Context ── -->
      <div class="shrink-0 border-t border-white/4 px-4 py-3">
        <div class="flex items-center gap-2 rounded-xl bg-white/3 px-3 py-2">
          <div class="flex h-7 w-7 items-center justify-center rounded-lg bg-linear-to-br from-spotify/30 to-aurora-purple/30">
            <i aria-hidden="true" class="pi pi-sparkles text-[10px] text-white/70" />
          </div>
          <p class="text-[10px] font-medium text-slate-500 leading-tight">
            {{ contextHint }}
          </p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserAuthStore } from '@/stores'

interface NavItem {
  label: string
  to: string
  icon?: string
  badge?: string | number
}

interface NavSection {
  label: string
  items: NavItem[]
}

interface QuickFilter {
  label: string
  active: boolean
  action?: () => void
}

const route = useRoute()
const router = useRouter()
const store = useUserAuthStore()

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}

// ── Category Detection ──
type Category = 'home' | 'search' | 'recommendations' | 'library' | 'social' | 'ai' | 'more'

const activeCategory = computed<Category>(() => {
  const path = route.path
  if (path === '/' || path.startsWith('/track/') || path.startsWith('/album/') || path.startsWith('/artist/')) return 'home'
  if (path.startsWith('/search')) return 'search'
  if (path.startsWith('/recommendations')) return 'recommendations'
  if (path.startsWith('/library') || path.startsWith('/playlist') || path.startsWith('/recently-played')) return 'library'
  if (path.startsWith('/social') || path.startsWith('/notifications')) return 'social'
  if (path.startsWith('/ai/')) return 'ai'
  return 'more'
})

const categoryMeta = computed(() => {
  const map: Record<Category, { label: string; icon: string; iconBg: string }> = {
    home: { label: 'Browse', icon: 'pi pi-home', iconBg: 'bg-spotify/20 text-spotify' },
    search: { label: 'Search', icon: 'pi pi-search', iconBg: 'bg-spotify/20 text-spotify' },
    recommendations: { label: 'For You', icon: 'pi pi-star', iconBg: 'bg-aurora-pink/20 text-aurora-pink' },
    library: { label: 'Library', icon: 'pi pi-bookmark', iconBg: 'bg-aurora-blue/20 text-aurora-blue' },
    social: { label: 'Community', icon: 'pi pi-users', iconBg: 'bg-amber-500/20 text-amber-500' },
    ai: { label: 'AI', icon: 'pi pi-magic', iconBg: 'bg-linear-to-br from-spotify/20 to-aurora-purple/20 text-white' },
    more: { label: 'More', icon: 'pi pi-ellipsis-h', iconBg: 'bg-white/10 text-white' },
  }
  return map[activeCategory.value]
})

// ── Quick Filters ──
const quickFilters = computed<QuickFilter[]>(() => {
  const path = route.path
  switch (activeCategory.value) {
    case 'home':
      return [
        { label: 'For You', active: path === '/', action: () => router.push('/') },
        { label: 'Trending', active: false, action: () => {} },
        { label: 'New', active: false, action: () => {} },
      ]
    case 'search':
      return [
        { label: 'All', active: path === '/search', action: () => router.push('/search') },
        { label: 'Songs', active: false },
        { label: 'Artists', active: false },
        { label: 'Albums', active: false },
      ]
    case 'recommendations':
      return [
        { label: 'All', active: path === '/recommendations' || path.startsWith('/recommendations/for-you'), action: () => router.push('/recommendations') },
        { label: 'Popular', active: path.startsWith('/recommendations/popular'), action: () => router.push('/recommendations/popular') },
        { label: 'Best', active: path.startsWith('/recommendations/best'), action: () => router.push('/recommendations/best') },
        { label: 'Recent', active: path.startsWith('/recommendations/recent'), action: () => router.push('/recommendations/recent') },
      ]
    case 'library':
      return [
        { label: 'Tracks', active: path === '/library' },
        { label: 'Albums', active: false },
        { label: 'Artists', active: false },
        { label: 'Playlists', active: path.startsWith('/playlist') },
      ]
    case 'social':
      return [
        { label: 'Feed', active: path === '/social' },
        { label: 'Parties', active: path.startsWith('/social/party') },
        { label: 'Rooms', active: path.startsWith('/social/room') },
        { label: 'Clubs', active: path.startsWith('/social/club') || path.startsWith('/social/clubs') },
      ]
    default:
      return []
  }
})

// ── Navigation Sections ──
const navSections = computed<NavSection[]>(() => {
  const cat = activeCategory.value
  const sections: NavSection[] = []

  switch (cat) {
    case 'home':
      sections.push({
        label: 'Main',
        items: [
          { label: 'Home', to: '/', icon: 'pi pi-home' },
          { label: 'Search', to: '/search', icon: 'pi pi-search' },
        ],
      })
      sections.push({
        label: 'Quick Access',
        items: [
          { label: 'Library', to: '/library', icon: 'pi pi-bookmark' },
          { label: 'Playlists', to: '/playlists', icon: 'pi pi-list' },
          { label: 'Recently Played', to: '/recently-played', icon: 'pi pi-history' },
        ],
      })
      break

    case 'search':
      sections.push({
        label: 'Browse',
        items: [
          { label: 'Genres', to: '/search', icon: 'pi pi-tag' },
          { label: 'Moods', to: '/ai/mood-explorer', icon: 'pi pi-heart' },
          { label: 'New Releases', to: '/recommendations/recent', icon: 'pi pi-star' },
          { label: 'Popular', to: '/recommendations/popular', icon: 'pi pi-chart-bar' },
        ],
      })
      sections.push({
        label: 'Discover',
        items: [
          { label: 'Made For You', to: '/search', icon: 'pi pi-user' },
          { label: 'Trending', to: '/search', icon: 'pi pi-fire' },
          { label: 'Viral Hits', to: '/search', icon: 'pi pi-bolt' },
        ],
      })
      break

    case 'recommendations':
      sections.push({
        label: 'Curated',
        items: [
          { label: 'For You', to: '/recommendations', icon: 'pi pi-user' },
          { label: 'Popular', to: '/recommendations/popular', icon: 'pi pi-fire' },
          { label: 'Best of All Time', to: '/recommendations/best', icon: 'pi pi-trophy' },
          { label: 'Recent', to: '/recommendations/recent', icon: 'pi pi-clock' },
        ],
      })
      break

    case 'library':
      sections.push({
        label: 'Collection',
        items: [
          { label: 'Tracks', to: '/library', icon: 'pi pi-music' },
          { label: 'Albums', to: '/library/albums', icon: 'pi pi-book' },
          { label: 'Artists', to: '/library/artists', icon: 'pi pi-users' },
          { label: 'Playlists', to: '/playlists', icon: 'pi pi-list' },
        ],
      })
      if (store.isAuthenticated) {
        sections.push({
          label: 'History',
          items: [
            { label: 'Recently Played', to: '/recently-played', icon: 'pi pi-history' },
          ],
        })
      }
      break

    case 'social':
      sections.push({
        label: 'Social',
        items: [
          { label: 'Feed', to: '/social', icon: 'pi pi-clock' },
          { label: 'Notifications', to: '/notifications', icon: 'pi pi-bell', badge: '3' },
          { label: 'Messages', to: '/social/messages', icon: 'pi pi-comments' },
        ],
      })
      sections.push({
        label: 'Live',
        items: [
          { label: 'Listening Parties', to: '/social/party', icon: 'pi pi-users' },
          { label: 'Live Rooms', to: '/social/room', icon: 'pi pi-video' },
          { label: 'Music Clubs', to: '/social/clubs/browse', icon: 'pi pi-building' },
        ],
      })
      break

    case 'ai':
      sections.push({
        label: 'AI Features',
        items: [
          { label: 'Mood Explorer', to: '/ai/mood-explorer', icon: 'pi pi-magic' },
          { label: 'Playlist Generator', to: '/ai/playlist-generator', icon: 'pi pi-sync' },
        ],
      })
      break

    case 'more':
      sections.push({
        label: 'Account',
        items: [
          { label: 'Profile', to: '/profile', icon: 'pi pi-user' },
          { label: 'Settings', to: '/settings', icon: 'pi pi-cog' },
          { label: 'Subscription', to: '/subscription', icon: 'pi pi-credit-card' },
        ],
      })
      sections.push({
        label: 'Community',
        items: [
          { label: 'Gamification', to: '/gamification', icon: 'pi pi-trophy' },
          { label: 'Contributions', to: '/contributions', icon: 'pi pi-cloud-upload' },
          { label: 'Creator Dashboard', to: '/creator-dashboard', icon: 'pi pi-chart-bar' },
        ],
      })
      break
  }

  return sections
})

// ── Context Hint ──
const contextHint = computed(() => {
  const map: Record<Category, string> = {
    home: 'Explore your music world',
    search: 'Search songs, artists, albums, and more',
    recommendations: 'Handpicked just for you',
    library: 'Your personal collection',
    social: 'Connect through music',
    ai: 'Powered by intelligence',
    more: 'Everything else',
  }
  return map[activeCategory.value]
})
</script>
