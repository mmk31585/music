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
        component: () => import('@/pages/admin/PageAdminDashboard.vue'),
        meta: { title: 'Dashboard' },
      },
      {
        path: 'catalog',
        name: 'admin.catalog',
        component: () => import('@/pages/admin/PageAdminCatalog.vue'),
        meta: { title: 'Catalog' },
      },
      {
        path: 'tracks',
        name: 'admin.tracks',
        component: () => import('@/pages/admin/PageAdminTracks.vue'),
        meta: { title: 'Tracks' },
      },
      {
        path: 'tracks/:id',
        name: 'admin.track.detail',
        component: () => import('@/pages/admin/PageAdminTrackDetail.vue'),
        meta: { title: 'Track Detail' },
      },
      {
        path: 'artists',
        name: 'admin.artists',
        component: () => import('@/pages/admin/AdminArtistsPage.vue'),
        meta: { title: 'Artists' },
      },
      {
        path: 'artists/:id',
        name: 'admin.artist.detail',
        component: () => import('@/pages/admin/PageAdminArtistDetail.vue'),
        meta: { title: 'Artist Detail' },
      },
      {
        path: 'albums',
        name: 'admin.albums',
        component: () => import('@/pages/admin/AdminAlbumsPage.vue'),
        meta: { title: 'Albums' },
      },
      {
        path: 'albums/:id',
        name: 'admin.album.detail',
        component: () => import('@/pages/admin/PageAdminAlbumDetail.vue'),
        meta: { title: 'Album Detail' },
      },
      {
        path: 'genres',
        name: 'admin.genres',
        component: () => import('@/pages/admin/AdminGenresPage.vue'),
        meta: { title: 'Genres' },
      },
      {
        path: 'users',
        name: 'admin.users',
        component: () => import('@/pages/admin/PageAdminUsers.vue'),
        meta: { title: 'Users' },
      },
      {
        path: 'media',
        name: 'admin.media',
        component: () => import('@/pages/admin/PageAdminMedia.vue'),
        meta: { title: 'Media' },
      },
      {
        path: 'videos',
        name: 'admin.videos',
        component: () => import('@/pages/admin/PageAdminVideos.vue'),
        meta: { title: 'Videos' },
      },
      {
        path: 'video-upload',
        name: 'admin.video-upload',
        component: () => import('@/pages/admin/PageAdminVideoUpload.vue'),
        meta: { title: 'Video Upload' },
      },
      {
        path: 'import',
        name: 'admin.import',
        component: () => import('@/pages/admin/PageAdminImport.vue'),
        meta: { title: 'Import from Internet' },
      },
      {
        path: 'import/artist',
        name: 'admin.import.artist',
        component: () => import('@/pages/admin/PageAdminImportArtist.vue'),
        meta: { title: 'Import by Artist' },
      },
      {
        path: 'ingestion',
        name: 'admin.ingestion',
        component: () => import('@/pages/admin/PageAdminIngestion.vue'),
        meta: { title: 'Music Ingestion' },
      },
      {
        path: 'ingestion/review/:id',
        name: 'admin.ingestion.review',
        component: () => import('@/pages/admin/PageAdminIngestionReview.vue'),
        meta: { title: 'Review Draft' },
      },
      {
        path: 'moderation',
        name: 'admin.moderation',
        component: () => import('@/pages/admin/PageAdminModeration.vue'),
        meta: { title: 'Moderation' },
      },
      {
        path: 'subscriptions',
        name: 'admin.subscriptions',
        component: () => import('@/pages/admin/PageAdminSubscriptions.vue'),
        meta: { title: 'Subscriptions' },
      },
      {
        path: 'contributions',
        name: 'admin.contributions',
        component: () => import('@/pages/admin/PageAdminContributions.vue'),
        meta: { title: 'Contributions' },
      },
    ],
  },
]
