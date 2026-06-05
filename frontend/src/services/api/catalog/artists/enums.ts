export enum ArtistApiRoutes {
  // Public routes
  LIST = '/catalog/artists',
  GET = '/catalog/artists/:artistId',

  // Admin routes - NOTE: Backend uses singular "/artist", fix this!
  ADMIN_LIST = '/admin/catalog/artist',
  ADMIN_CREATE = '/admin/catalog/artists',
  ADMIN_UPDATE = '/admin/catalog/artists/:artistId',
  ADMIN_DELETE = '/admin/catalog/artists/:artistId',
}
