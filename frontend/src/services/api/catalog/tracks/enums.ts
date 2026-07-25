export enum TrackApiRoutes {
  // Public routes
  LIST = '/catalog/tracks',
  RANDOM = '/catalog/tracks/random',
  GET = '/catalog/tracks/:trackId',
  CREDITS = '/catalog/tracks/:trackId/credits',
  TRACK_ARTISTS = '/catalog/tracks/:trackId/artists',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/tracks',
  ADMIN_CREATE = '/admin/catalog/tracks',
  ADMIN_UPDATE = '/admin/catalog/tracks/:trackId',
  ADMIN_DELETE = '/admin/catalog/tracks/:trackId',
  ADMIN_ENRICH = '/admin/catalog/tracks/:trackId/enrich',
  ADMIN_ENRICH_ALL = '/admin/catalog/tracks/enrich-all',
  ADMIN_UPLOAD = '/admin/media/upload',
  ADMIN_REPLACE_CREDITS = '/admin/catalog/tracks/:trackId/credits',
  ADMIN_REPLACE_ARTISTS = '/admin/catalog/tracks/:trackId/artists',
}
