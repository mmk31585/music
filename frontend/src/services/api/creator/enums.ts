export enum CreatorApiRoutes {
  OVERVIEW = '/creator/overview',
  DAILY_STATS = '/creator/daily',
  TRACK_STATS = '/creator/tracks',
  REFRESH = '/creator/refresh',
  IS_CREATOR = '/creator/check',

  EARNINGS = '/creator/earnings',
  PAYOUTS = '/creator/earnings/payouts',
  PAYOUT_METHODS = '/creator/earnings/methods',

  AUDIENCE = '/creator/audience',

  CONTENT = '/creator/content',
  UPDATE_TRACK = '/creator/tracks/:trackId',
  UPDATE_ALBUM = '/creator/albums/:albumId',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  DELETE_TRACK = '/creator/tracks/:trackId',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  DELETE_ALBUM = '/creator/albums/:albumId',
}
