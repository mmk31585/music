export enum ArtistApiRoutes {
  // Public routes
  LIST = '/catalog/artists',
  GET = '/catalog/artists/:artistId',
  OVERVIEW = '/catalog/artists/:artistId/overview',
  TRACKS = '/catalog/artists/:artistId/tracks',
  ALBUMS = '/catalog/artists/:artistId/albums',
  SINGLES = '/catalog/artists/:artistId/singles',
  APPEARS_ON = '/catalog/artists/:artistId/appears-on',
  TOP_TRACKS = '/catalog/artists/:artistId/top-tracks',
  RELATED = '/catalog/artists/:artistId/related',

  // Admin routes
  ADMIN_LIST = '/admin/catalog/artist',
  ADMIN_CREATE = '/admin/catalog/artist',
  ADMIN_UPDATE = '/admin/catalog/artist/:artistId',
  ADMIN_DELETE = '/admin/catalog/artist/:artistId',
  ADMIN_ENRICH = '/admin/catalog/artist/:artistId/enrich',
  ADMIN_REPLACE_RELATED = '/admin/catalog/artist/:artistId/related',
  ADMIN_REPLACE_TOP_TRACKS = '/admin/catalog/artist/:artistId/top-tracks',
}
