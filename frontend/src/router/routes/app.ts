import type { RouteRecordRaw } from 'vue-router'
import LayoutMusicApp from '@/layouts/LayoutMusicApp.vue'

export const appRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: LayoutMusicApp,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'app.home',
        component: () => import('@/pages/app/PageHome.vue'),
      },
      {
        path: 'search',
        name: 'app.search',
        component: () => import('@/pages/app/PageSearch.vue'),
      },
      // {
      //   path: 'tracks/:trackId',
      //   name: 'app.track.details',
      //   component: () => import('@/pages/app/PageTrackDetails.vue'),
      //   props: true,
      // },
      // {
      //   path: 'albums/:albumId',
      //   name: 'app.album.details',
      //   component: () => import('@/pages/app/PageAlbumDetails.vue'),
      //   props: true,
      // },
      // {
      //   path: 'artist/:artistId',
      //   name: 'app.artist.details',
      //   component: () => import('@/pages/app/PageArtistDetails.vue'),
      //   props: true,
      // },
      // {
      //   path: 'genres/:genreId',
      //   name: 'app.genre.details',
      //   component: () => import('@/pages/app/PageGenreDetails.vue'),
      //   props: true,
      // },
    ],
  },
]
