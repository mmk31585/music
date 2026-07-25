export enum ReactionsApiRoutes {
  REACT = '/reactions',
  REMOVE = '/reactions/:targetType/:targetId',
  USER_REACTION = '/reactions/:targetType/:targetId/mine',
  COUNTS = '/reactions/:targetType/:targetId/counts',
  BULK_COUNTS = '/reactions/bulk-counts',
  LIKED_TRACKS = '/reactions/tracks',
  LIKED_ALBUMS = '/reactions/albums',
}
