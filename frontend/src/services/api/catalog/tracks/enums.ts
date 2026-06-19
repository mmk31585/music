export enum TrackApiRoutes {
  // Public routes
  LIST = '/catalog/tracks',
  RANDOM = '/catalog/tracks/random',
  GET = '/catalog/tracks/:trackId',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/tracks',
  ADMIN_CREATE = '/admin/catalog/tracks',
  ADMIN_UPDATE = '/admin/catalog/tracks/:trackId',
  ADMIN_DELETE = '/admin/catalog/tracks/:trackId',
  ADMIN_ENRICH = '/admin/catalog/tracks/:trackId/enrich',
  ADMIN_ENRICH_ALL = '/admin/catalog/tracks/enrich-all',
  ADMIN_UPLOAD = '/admin/media/upload',
}
