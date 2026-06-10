export enum ModerationApiRoutes {
  REPORT = '/moderation/report',
  PENDING = '/moderation/pending',
  RESOLVE = '/moderation/resolve/:id',
  BY_STATUS = '/moderation/status/:status',
  REPORT_DETAIL = '/moderation/reports/:id',
  BULK = '/moderation/bulk',
  FLAG = '/moderation/flag',
  FLAGS = '/moderation/flags',
  STATS = '/moderation/stats',
  ACTIONS = '/moderation/actions',
}
