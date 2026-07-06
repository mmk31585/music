export enum AlbumApiRoutes {
  // Public routes
  LIST = '/catalog/albums',
  GET = '/catalog/albums/:albumId',
  TRACKS = '/catalog/albums/:albumId/tracks',
  ARTISTS = '/catalog/albums/:albumId/artists',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/albums',
  ADMIN_CREATE = '/admin/catalog/albums',
  ADMIN_UPDATE = '/admin/catalog/albums/:albumId',
  ADMIN_DELETE = '/admin/catalog/albums/:albumId',
  ADMIN_ENRICH = '/admin/catalog/albums/:albumId/enrich',
  ADMIN_REPLACE_ARTISTS = '/admin/catalog/albums/:albumId/artists',
}
