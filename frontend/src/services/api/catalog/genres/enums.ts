export enum GenreApiRoutes {
  // Public routes
  LIST = '/catalog/genres',
  GET = '/catalog/genres/:genreId',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/genres',
  ADMIN_CREATE = '/admin/catalog/genres/create',
  ADMIN_UPDATE = '/admin/catalog/genres/:genreId',
  ADMIN_DELETE = '/admin/catalog/genres/:genreId/delete',
}
