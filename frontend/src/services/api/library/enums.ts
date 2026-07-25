export enum LibraryApiRoutes {
  TRACKS = '/library/tracks',
  LIKE_TRACK = '/library/tracks/like',
  UNLIKE_TRACK = '/library/tracks/:trackId/like',

  ALBUMS = '/library/albums',
  LIKE_ALBUM = '/library/albums/like',
  UNLIKE_ALBUM = '/library/albums/:albumId/like',

  ARTISTS = '/library/artists',
  FOLLOW_ARTIST = '/library/artists/follow',
  UNFOLLOW_ARTIST = '/library/artists/:artistId/follow',

  HISTORY = '/library/history',
  RECENTLY_PLAYED = '/library/recently-played',
}
