export enum PlaylistApiRoutes {
  LIST_PUBLIC = '/playlists',
  GET = '/playlists/:playlistId',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  CREATE = '/playlists',
  LIST_MY = '/playlists/me',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  UPDATE = '/playlists/:playlistId',
  // eslint-disable-next-line @typescript-eslint/no-duplicate-enum-values
  DELETE = '/playlists/:playlistId',
  ADD_TRACK = '/playlists/:playlistId/tracks',
  REMOVE_TRACK = '/playlists/:playlistId/tracks/:trackId',
  REORDER_TRACKS = '/playlists/:playlistId/tracks/reorder',
  SET_COLLABORATIVE = '/playlists/:playlistId/collaborative',
   
  LIST_COLLABORATORS = '/playlists/:playlistId/collaborators',
}
