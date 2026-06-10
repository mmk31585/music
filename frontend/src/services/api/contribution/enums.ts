export enum ContributionApiRoutes {
  CREATE = '/contributions',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  LIST_MY = '/contributions',
  LIST_PENDING = '/contributions/pending',
  GET_BY_ID = '/contributions/:id',
  REVIEW = '/contributions/:id/review',
  HISTORY = '/contributions/:id/history',
  LEADERBOARD = '/contributions/leaderboard',
  BY_TARGET = '/content/:targetType/:targetId/contributions',
  VERSIONS = '/content/:targetType/:targetId/versions',
  APPLY = '/admin/contributions/:id/apply',
  ADMIN_LIST = '/admin/contributions',
}
