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
	social.GET("/discussions/:id/replies", handler.GetDiscussionReplies)
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

		// Music Clubs
		protected.POST("/clubs", handler.CreateClub)
		protected.POST("/clubs/:id/join", handler.JoinClub)
		protected.POST("/clubs/:id/leave", handler.LeaveClub)
		protected.POST("/clubs/:id/posts", handler.CreateClubPost)

		// Discussions
		protected.POST("/discussions", handler.CreateDiscussion)

		// Track Ratings
		protected.POST("/ratings", handler.CreateRating)
	}
}
