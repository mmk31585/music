package social

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	social := rg.Group("/social")

	// Public read-only endpoints (no auth required)
	social.GET("/parties", handler.ListActiveParties)
	social.GET("/parties/:id", handler.GetParty)
	social.GET("/rooms", handler.ListActiveRooms)
	social.GET("/rooms/:id", handler.GetRoom)
	social.GET("/rooms/:id/participants", handler.GetRoomParticipants)
	social.GET("/rooms/:id/queue", handler.GetRoomQueue)
	social.GET("/clubs", handler.ListClubs)
	social.GET("/clubs/:id", handler.GetClub)
	social.GET("/clubs/:id/members", handler.GetClubMembers)
	social.GET("/clubs/:id/posts", handler.GetClubPosts)
	social.GET("/discussions", handler.GetDiscussions)
	social.GET("/followers/:userId", handler.GetFollowers)
	social.GET("/following/:userId", handler.GetFollowing)
	social.GET("/following", handler.GetFollowing)
	social.GET("/ratings/:trackId", handler.GetTrackRatings)

	// Protected endpoints (require auth)
	protected := social.Group("")
	protected.Use(authMW)
	{
		protected.GET("/feed", handler.GetFeed)
		protected.POST("/follow/:userId", handler.Follow)
		protected.DELETE("/follow/:userId", handler.Unfollow)
		protected.GET("/is-following/:userId", handler.IsFollowing)

		// Listening Parties
		protected.POST("/parties", handler.CreateParty)
		protected.PUT("/parties/:id/status", handler.UpdatePartyStatus)
		protected.POST("/parties/:id/join", handler.JoinParty)
		protected.POST("/parties/:id/leave", handler.LeaveParty)

		// Live Rooms
		protected.POST("/rooms", handler.CreateRoom)
		protected.POST("/rooms/:id/join", handler.JoinRoom)
		protected.POST("/rooms/:id/leave", handler.LeaveRoom)
		protected.POST("/rooms/:id/queue", handler.AddToRoomQueue)

		// Democratic Voting Queue
		protected.POST("/rooms/:id/queue/suggest", handler.SuggestTrack)
		protected.POST("/rooms/:id/queue/:candidateId/vote", handler.CastVote)
		protected.DELETE("/rooms/:id/queue/:candidateId/vote", handler.RemoveVote)
		protected.GET("/rooms/:id/queue/state", handler.GetQueueState)
		protected.POST("/rooms/:id/track-ended", handler.TrackEnded)

		// Stage & Raise-Hand
		protected.POST("/rooms/:id/stage/raise-hand", handler.RaiseHand)
		protected.POST("/rooms/:id/stage/lower-hand", handler.LowerHand)
		protected.POST("/rooms/:id/stage/:userId/approve", handler.ApproveHand)
		protected.POST("/rooms/:id/stage/:userId/deny", handler.DenyHand)
		protected.DELETE("/rooms/:id/stage/:userId", handler.RemoveFromStage)
		protected.POST("/rooms/:id/stage/leave", handler.LeaveStage)
		protected.POST("/rooms/:id/stage/:userId/mute", handler.ToggleMute)
		protected.GET("/rooms/:id/stage", handler.GetStageState)

		// Music Clubs (Phase 5)
		protected.POST("/clubs", handler.CreateClub)
		protected.GET("/clubs/browse", handler.ListClubsWithGenre)
		protected.GET("/clubs/:id/detail", handler.GetClubDetail)
		protected.POST("/clubs/:id/join", handler.JoinClub)
		protected.POST("/clubs/:id/leave", handler.LeaveClub)
		protected.POST("/clubs/:id/launch-party", handler.LaunchPartyFromClub)
		protected.POST("/clubs/:id/posts", handler.CreateClubPost)

		// Discussions
		protected.POST("/discussions", handler.CreateDiscussion)

		// Club Discussions (Phase 6)
		protected.POST("/clubs/:id/discussions", handler.CreateClubDiscussion)
		protected.GET("/clubs/:id/discussions", handler.ListClubDiscussions)
		protected.GET("/discussions/:id", handler.GetClubDiscussion)
		protected.GET("/discussions/:id/replies", handler.GetClubDiscussionReplies)
		protected.POST("/discussions/:id/replies", handler.CreateClubDiscussionReply)
		protected.DELETE("/discussions/:id", handler.DeleteClubDiscussion)
		protected.DELETE("/discussions/:id/replies/:replyId", handler.DeleteClubDiscussionReply)

		// Track Ratings
		protected.POST("/ratings", handler.CreateRating)
	}
}
