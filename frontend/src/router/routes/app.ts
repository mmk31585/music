import type { RouteRecordRaw } from 'vue-router'
import LayoutMusicApp from '@/layouts/LayoutMusicApp.vue'

export const appRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: LayoutMusicApp,
    children: [
      {
        path: '',
        name: 'app.home',
        component: () => import('@/pages/app/PageHome.vue'),
        meta: { title: 'Home' },
      },
      {
        path: 'search',
        name: 'search',
        component: () => import('@/pages/app/PageSearch.vue'),
        meta: { title: 'Search' },
      },
      {
        path: 'discover',
        name: 'discover',
        component: () => import('@/pages/app/PageDiscover.vue'),
        meta: { title: 'Discover' },
      },
      {
        path: 'recommendations',
        name: 'recommendations',
        component: () => import('@/pages/app/recommendations/PageRecommendationsHub.vue'),
        meta: { title: 'Recommendations' },
      },
      {
        path: 'recommendations/for-you',
        name: 'recommendations.for-you',
        component: () => import('@/pages/app/recommendations/PageRecommendationsForYou.vue'),
        meta: { title: 'For You' },
      },
      {
        path: 'recommendations/popular',
        name: 'recommendations.popular',
        component: () => import('@/pages/app/recommendations/PageRecommendationsPopular.vue'),
        meta: { title: 'Popular' },
      },
      {
        path: 'recommendations/best',
        name: 'recommendations.best',
        component: () => import('@/pages/app/recommendations/PageRecommendationsBest.vue'),
        meta: { title: 'Best' },
      },
      {
        path: 'recommendations/recent',
        name: 'recommendations.recent',
        component: () => import('@/pages/app/recommendations/PageRecommendationsRecent.vue'),
        meta: { title: 'Recent' },
      },
      {
        path: 'library',
        name: 'library',
        component: () => import('@/pages/app/PageLibrary.vue'),
        meta: { title: 'Library', requiresAuth: true },
      },
      {
        path: 'playlists',
        name: 'playlists',
        component: () => import('@/pages/app/PagePlaylists.vue'),
        meta: { title: 'Playlists', requiresAuth: true },
      },
      {
        path: 'playlist/:id',
        name: 'playlist.detail',
        component: () => import('@/pages/app/PagePlaylistDetail.vue'),
        meta: { title: 'Playlist' },
      },
      {
        path: 'recently-played',
        name: 'recently-played',
        component: () => import('@/pages/app/PageRecentlyPlayed.vue'),
        meta: { title: 'Recently Played' },
      },
      {
        path: 'track/:id',
        name: 'track.detail',
        component: () => import('@/pages/app/PageTrack.vue'),
        meta: { title: 'Track' },
      },
      {
        path: 'album/:id',
        name: 'album.detail',
        component: () => import('@/pages/app/PageAlbum.vue'),
        meta: { title: 'Album' },
      },
      {
        path: 'artist/:id',
        name: 'artist.detail',
        component: () => import('@/pages/app/PageArtist.vue'),
        meta: { title: 'Artist' },
      },
      {
        path: 'notifications',
        name: 'notifications',
        component: () => import('@/pages/app/PageNotifications.vue'),
        meta: { title: 'Notifications', requiresAuth: true },
      },
      {
        path: 'social',
        name: 'social',
        component: () => import('@/pages/app/PageSocial.vue'),
        meta: { title: 'Social Hub', requiresAuth: true },
      },
      {
        path: 'social/party/:id',
        name: 'social.party',
        component: () => import('@/pages/app/PagePartyDetail.vue'),
        meta: { title: 'Listening Party' },
      },
      {
        path: 'social/room/:id',
        name: 'social.room',
        component: () => import('@/pages/app/PageRoomLive.vue'),
        meta: { title: 'Live Room' },
      },
      {
        path: 'social/clubs/browse',
        name: 'social.clubs.browse',
        component: () => import('@/pages/app/PageClubsBrowse.vue'),
        meta: { title: 'Browse Clubs', requiresAuth: true },
      },
      {
        path: 'social/club/:id',
        name: 'social.club',
        component: () => import('@/pages/app/PageClubDetail.vue'),
        meta: { title: 'Music Club' },
      },
      {
        path: 'profile',
        redirect: { name: 'user.profile' },
      },
      {
        path: 'user/:id?',
        name: 'user.profile',
        component: () => import('@/pages/app/PageUserProfile.vue'),
        meta: { title: 'Profile' },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/pages/app/PageUserSettings.vue'),
        meta: { title: 'Settings', requiresAuth: true },
      },
      {
        path: 'subscription',
        name: 'subscription',
        component: () => import('@/pages/app/PageSubscriptions.vue'),
        meta: { title: 'Subscription', requiresAuth: true },
      },
      {
        path: 'contributions',
        name: 'contributions',
        component: () => import('@/pages/app/PageContributions.vue'),
        meta: { title: 'Contributions', requiresAuth: true },
      },
      {
        path: 'creator-dashboard',
        name: 'creator.dashboard',
        component: () => import('@/pages/app/PageCreatorDashboard.vue'),
        meta: { title: 'Creator Dashboard', requiresAuth: true },
      },
      {
        path: 'gamification',
        name: 'gamification',
        component: () => import('@/pages/app/PageGamification.vue'),
        meta: { title: 'Gamification', requiresAuth: true },
      },
      {
        path: 'ai/mood-explorer',
        name: 'ai.mood-explorer',
        component: () => import('@/pages/app/PageAIMoodExplorer.vue'),
        meta: { title: 'Mood Explorer' },
      },
      {
        path: 'ai/playlist-generator',
        name: 'ai.playlist-generator',
        component: () => import('@/pages/app/PageAIPlaylistGenerator.vue'),
        meta: { title: 'AI Playlist Generator', requiresAuth: true },
      },
    ],
  },
]
