export enum LyricsApiRoutes {
  GET_TRACK_LYRICS = 'tracks/:trackId/lyrics',
  GET_TRACK_LYRICS_ALL = 'tracks/:trackId/lyrics/all',

  ADMIN_CREATE = 'admin/lyrics',
  ADMIN_FETCH_LRC = 'admin/lyrics/fetch/:trackId',
  ADMIN_UPDATE = 'admin/lyrics/:id',
  ADMIN_DELETE = 'admin/lyrics/:id',
}
