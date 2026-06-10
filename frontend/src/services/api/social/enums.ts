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

  // Music Clubs
  CLUBS = '/social/clubs',
  CLUB = '/social/clubs/:id',
  CLUB_JOIN = '/social/clubs/:id/join',
  CLUB_LEAVE = '/social/clubs/:id/leave',
  CLUB_MEMBERS = '/social/clubs/:id/members',
  CLUB_POSTS = '/social/clubs/:id/posts',

  // Discussions
  DISCUSSIONS = '/social/discussions',
  DISCUSSION_REPLIES = '/social/discussions/:id/replies',

  // Track Ratings
  RATINGS = '/social/ratings',
  TRACK_RATINGS = '/social/ratings/:trackId',
}
