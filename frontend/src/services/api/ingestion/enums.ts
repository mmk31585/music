export enum IngestionApiRoutes {
  ADMIN_UPLOAD = '/admin/ingestion/upload',
  ADMIN_DRAFTS = '/admin/ingestion/drafts',
  ADMIN_ENRICH = '/admin/ingestion/drafts/:id/enrich',
  ADMIN_SUGGESTIONS = '/admin/ingestion/drafts/:id/suggestions',
  ADMIN_FINAL_METADATA = '/admin/ingestion/drafts/:id/final-metadata',
  ADMIN_REJECT = '/admin/ingestion/drafts/:id/reject',
  ADMIN_FINALIZE = '/admin/ingestion/drafts/:id/finalize',
  ADMIN_ARTISTS_SEARCH = '/admin/catalog/artists/search',
  ADMIN_ALBUMS_SEARCH = '/admin/catalog/albums/search',
  ADMIN_STATS = '/admin/ingestion/stats',
  ADMIN_CLEANUP = '/admin/ingestion/cleanup',
  ADMIN_CONFIG = '/admin/ingestion/config',
}
