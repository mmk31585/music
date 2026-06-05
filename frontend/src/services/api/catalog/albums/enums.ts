export enum AlbumApiRoutes {
  // Public routes
  LIST = '/catalog/albums',
  GET = '/catalog/albums/:albumId',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/albums',
  ADMIN_CREATE = '/admin/catalog/albums',
  ADMIN_UPDATE = '/admin/catalog/albums/:albumId',
  ADMIN_DELETE = '/admin/catalog/albums/:albumId',
}
