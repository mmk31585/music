export enum LyricsApiRoutes {
  GET_TRACK_LYRICS = 'tracks/:trackId/lyrics',
  GET_TRACK_LYRICS_ALL = 'tracks/:trackId/lyrics/all',

  ADMIN_CREATE = 'admin/lyrics',
  ADMIN_FETCH_LRC = 'admin/lyrics/fetch/:trackId',
  ADMIN_FETCH_OR_GENERATE = 'admin/lyrics/fetch-or-generate/:trackId',
  ADMIN_AI_STATUS = 'admin/lyrics/ai-status/:trackId',
  ADMIN_UPDATE = 'admin/lyrics/:id',
  ADMIN_DELETE = 'admin/lyrics/:id',
}
