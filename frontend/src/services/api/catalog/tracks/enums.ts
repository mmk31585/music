export enum TrackApiRoutes {
  // Public routes
  LIST = '/catalog/tracks',
  GET = '/catalog/tracks/:trackId',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/tracks',
  ADMIN_CREATE = '/admin/catalog/tracks/create',
  ADMIN_UPDATE = '/admin/catalog/tracks/:trackId',
  ADMIN_DELETE = '/admin/catalog/tracks/:trackId/delete',
  ADMIN_UPLOAD = '/admin/media/upload',
}
