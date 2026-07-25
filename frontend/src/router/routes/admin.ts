import type { RouteRecordRaw } from 'vue-router'
import LayoutAdmin from '@/layouts/LayoutAdmin.vue'

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin',
    component: LayoutAdmin,
    meta: {
      requiresAuth: true,
      requiresRole: 'admin',
    },
    children: [
      {
        path: '',
        name: 'admin.dashboard',
        component: () => import(/* webpackChunkName: "admin-dashboard" */ '@/pages/admin/PageAdminDashboard.vue'),
        meta: { title: 'Dashboard' },
      },
      {
        path: 'catalog',
        name: 'admin.catalog',
        component: () => import(/* webpackChunkName: "admin-catalog" */ '@/pages/admin/PageAdminCatalog.vue'),
        meta: { title: 'Catalog' },
      },
      {
        path: 'tracks',
        name: 'admin.tracks',
        component: () => import(/* webpackChunkName: "admin-tracks" */ '@/pages/admin/PageAdminTracks.vue'),
        meta: { title: 'Tracks' },
      },
      {
        path: 'tracks/:id',
        name: 'admin.track.detail',
        component: () => import(/* webpackChunkName: "admin-track-detail" */ '@/pages/admin/PageAdminTrackDetail.vue'),
        meta: { title: 'Track Detail' },
      },
      {
        path: 'artists',
        name: 'admin.artists',
        component: () => import(/* webpackChunkName: "admin-artists" */ '@/pages/admin/AdminArtistsPage.vue'),
        meta: { title: 'Artists' },
      },
      {
        path: 'artists/:id',
        name: 'admin.artist.detail',
        component: () => import(/* webpackChunkName: "admin-artist-detail" */ '@/pages/admin/PageAdminArtistDetail.vue'),
        meta: { title: 'Artist Detail' },
      },
      {
        path: 'albums',
        name: 'admin.albums',
        component: () => import(/* webpackChunkName: "admin-albums" */ '@/pages/admin/AdminAlbumsPage.vue'),
        meta: { title: 'Albums' },
      },
      {
        path: 'albums/:id',
        name: 'admin.album.detail',
        component: () => import(/* webpackChunkName: "admin-album-detail" */ '@/pages/admin/PageAdminAlbumDetail.vue'),
        meta: { title: 'Album Detail' },
      },
      {
        path: 'genres',
        name: 'admin.genres',
        component: () => import(/* webpackChunkName: "admin-genres" */ '@/pages/admin/AdminGenresPage.vue'),
        meta: { title: 'Genres' },
      },
      {
        path: 'users',
        name: 'admin.users',
        component: () => import(/* webpackChunkName: "admin-users" */ '@/pages/admin/PageAdminUsers.vue'),
        meta: { title: 'Users' },
      },
      {
        path: 'media',
        name: 'admin.media',
        component: () => import(/* webpackChunkName: "admin-media" */ '@/pages/admin/PageAdminMedia.vue'),
        meta: { title: 'Media' },
      },
      {
        path: 'videos',
        name: 'admin.videos',
        component: () => import(/* webpackChunkName: "admin-videos" */ '@/pages/admin/PageAdminVideos.vue'),
        meta: { title: 'Videos' },
      },
      {
        path: 'video-upload',
        name: 'admin.video-upload',
        component: () => import(/* webpackChunkName: "admin-video-upload" */ '@/pages/admin/PageAdminVideoUpload.vue'),
        meta: { title: 'Video Upload' },
      },
      {
        path: 'import',
        name: 'admin.import',
        component: () => import(/* webpackChunkName: "admin-import" */ '@/pages/admin/PageAdminImport.vue'),
        meta: { title: 'Import from Internet' },
      },
      {
        path: 'import/artist',
        name: 'admin.import.artist',
        component: () => import(/* webpackChunkName: "admin-import-artist" */ '@/pages/admin/PageAdminImportArtist.vue'),
        meta: { title: 'Import by Artist' },
      },
      {
        path: 'ingestion',
        name: 'admin.ingestion',
        component: () => import(/* webpackChunkName: "admin-ingestion" */ '@/pages/admin/PageAdminIngestion.vue'),
        meta: { title: 'Music Ingestion' },
      },
      {
        path: 'ingestion/review/:id',
        name: 'admin.ingestion.review',
        component: () => import(/* webpackChunkName: "admin-ingestion-review" */ '@/pages/admin/PageAdminIngestionReview.vue'),
        meta: { title: 'Review Draft' },
      },
      {
        path: 'moderation',
        name: 'admin.moderation',
        component: () => import(/* webpackChunkName: "admin-moderation" */ '@/pages/admin/PageAdminModeration.vue'),
        meta: { title: 'Moderation' },
      },
      {
        path: 'subscriptions',
        name: 'admin.subscriptions',
        component: () => import(/* webpackChunkName: "admin-subscriptions" */ '@/pages/admin/PageAdminSubscriptions.vue'),
        meta: { title: 'Subscriptions' },
      },
      {
        path: 'contributions',
        name: 'admin.contributions',
        component: () => import(/* webpackChunkName: "admin-contributions" */ '@/pages/admin/PageAdminContributions.vue'),
        meta: { title: 'Contributions' },
      },
      {
        path: 'analytics',
        name: 'admin.analytics',
        component: () => import(/* webpackChunkName: "admin-analytics" */ '@/pages/admin/PageAdminAnalytics.vue'),
        meta: { title: 'Analytics' },
      },
      {
        path: 'cache',
        name: 'admin.cache',
        component: () => import(/* webpackChunkName: "admin-cache" */ '@/pages/admin/PageAdminCache.vue'),
        meta: { title: 'Cache Management' },
      },
      {
        path: 'permissions',
        name: 'admin.permissions',
        component: () => import(/* webpackChunkName: "admin-permissions" */ '@/pages/admin/PageAdminPermissions.vue'),
        meta: { title: 'Permissions' },
      },
      {
        path: 'settings',
        name: 'admin.settings',
        component: () => import(/* webpackChunkName: "admin-settings" */ '@/pages/admin/PageAdminSettings.vue'),
        meta: { title: 'Settings' },
      },
    ],
  },
]
