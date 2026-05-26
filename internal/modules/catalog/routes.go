package catalog

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
)

func RegisterPublicRoutes(rg *gin.RouterGroup, h *Handler) {
	catalog := rg.Group("/catalog")
	{
		// Artists
		catalog.GET("/artists", h.ListArtists)
		catalog.GET("/artists/:artistID", h.GetArtist)

		// Albums
		catalog.GET("/albums", h.ListAlbums)
		catalog.GET("/albums/:albumID", h.GetAlbum)

		// Tracks
		catalog.GET("/tracks", h.ListTracks)
		catalog.GET("/tracks/:trackID", h.GetTrack)

		// Genres
		catalog.GET("/genres", h.ListGenres)
		catalog.GET("/genres/:genreID", h.GetGenre)

		// Search
		catalog.GET("/search", h.Search)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	adminCatalog := rg.Group("/admin/catalog")
	adminCatalog.Use(authMW)
	adminCatalog.Use(auth.RequireRole("admin")) // ensure this returns gin.HandlerFunc

	// Artists admin
	artists := adminCatalog.Group("/artists")
	{
		artists.POST("/", h.RegisterArtist)
		artists.PATCH("/:artistID", h.UpdateArtist)
		artists.DELETE("/:artistID", h.DeleteArtist)
	}

	// Albums admin
	albums := adminCatalog.Group("/albums")
	{
		albums.POST("/", h.RegisterAlbum)
		albums.PATCH("/:albumID", h.UpdateAlbum)
		albums.DELETE("/:albumID", h.DeleteAlbum)
	}

	// Tracks admin
	tracks := adminCatalog.Group("/tracks")
	{
		tracks.GET("/", h.ListAdminTracks)
		tracks.POST("/upload", h.RegisterTrackWithAudio)
		tracks.POST("/", h.RegisterTrack)
		tracks.PATCH("/:trackID", h.UpdateTrack)
		tracks.DELETE("/:trackID", h.DeleteTrack)
	}

	// Genres admin
	genres := adminCatalog.Group("/genres")
	{
		genres.POST("/", h.CreateGenre)
		genres.PATCH("/:genreID", h.UpdateGenre)
		genres.DELETE("/:genreID", h.DeleteGenre)
	}
}
