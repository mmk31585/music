export enum SocialApiRoutes {
   
  FOLLOW = '/social/follow/:userId',
   
  UNFOLLOW = '/social/unfollow/:userId',
  FOLLOWERS = '/social/followers/:userId',
  FOLLOWING = '/social/following/:userId',
  IS_FOLLOWING = '/social/is-following/:userId',
  FEED = '/social/feed',

  // Listening Parties
  PARTIES = '/social/parties',
  PARTY = '/social/parties/:id',
  PARTY_STATUS = '/social/parties/:id/status',
  PARTY_JOIN = '/social/parties/:id/join',
  PARTY_LEAVE = '/social/parties/:id/leave',

  // Live Rooms
  ROOMS = '/social/rooms',
  ROOM = '/social/rooms/:id',
  ROOM_JOIN = '/social/rooms/:id/join',
  ROOM_LEAVE = '/social/rooms/:id/leave',
  ROOM_PARTICIPANTS = '/social/rooms/:id/participants',
  ROOM_QUEUE = '/social/rooms/:id/queue',

  // Music Clubs (Phase 5)
  CLUBS = '/social/clubs',
  CLUB = '/social/clubs/:id',
  CLUB_BROWSE = '/social/clubs/browse',
  CLUB_DETAIL = '/social/clubs/:id/detail',
  CLUB_JOIN = '/social/clubs/:id/join',
  CLUB_LEAVE = '/social/clubs/:id/leave',
  CLUB_LAUNCH_PARTY = '/social/clubs/:id/launch-party',
  CLUB_MEMBERS = '/social/clubs/:id/members',
  CLUB_POSTS = '/social/clubs/:id/posts',

  // Discussions
  DISCUSSIONS = '/social/discussions',
  DISCUSSION_REPLIES = '/social/discussions/:id/replies',

  // Club Discussions (Phase 6)
  CLUB_DISCUSSIONS = '/social/clubs/:clubId/discussions',
  CLUB_DISCUSSION = '/social/discussions/:id',
  CLUB_DISCUSSION_REPLIES = '/social/discussions/:id/replies',
  CLUB_DISCUSSION_REPLY_DELETE = '/social/discussions/:id/replies/:replyId',

  // Track Ratings
  RATINGS = '/social/ratings',
  TRACK_RATINGS = '/social/ratings/:trackId',
}
