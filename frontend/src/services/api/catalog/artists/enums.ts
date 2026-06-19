export enum ArtistApiRoutes {
  // Public routes
  LIST = '/catalog/artists',
  GET = '/catalog/artists/:artistId',

  ADMIN_LIST = '/admin/catalog/artist',
  ADMIN_CREATE = '/admin/catalog/artist',
  ADMIN_UPDATE = '/admin/catalog/artist/:artistId',
  ADMIN_DELETE = '/admin/catalog/artist/:artistId',
  ADMIN_ENRICH = '/admin/catalog/artist/:artistId/enrich',
}
