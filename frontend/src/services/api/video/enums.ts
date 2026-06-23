export enum VideoApiRoutes {
  EXPLORE = '/explore/videos',
  TRACK_VIDEOS = '/tracks/:trackId/videos',
  VIDEO_LIKE = '/videos/:id/like',
  VIDEO_VIEW = '/videos/:id/view',
  UPLOAD_EDIT = '/videos',
  VIDEO_JOB = '/video/jobs/:jobId',
  MUSIC_STATUS = '/users/me/music-status',
  LIKE_VISIBILITY = '/tracks/:id/like/visibility',
  USER_VIDEOS = '/users/:userId/videos',
  LIKED_TRACKS = '/users/:userId/liked-tracks',
}
