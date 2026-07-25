import { computed, watch } from 'vue'
import type { RouteLocationNormalizedLoaded } from 'vue-router'

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

const allNavItems = [...browseItems, ...libraryItems, ...socialItems, ...moreItems]

export function usePageMeta(route: RouteLocationNormalizedLoaded) {
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

  watch(pageTitle, (title) => {
    document.title = title ? `${title} — Muse` : 'Muse'
  }, { immediate: true })

  return {
    pageTitle,
    activeNavBase,
    browseItems,
    libraryItems,
    socialItems,
    moreItems,
  }
}
