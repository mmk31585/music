package catalog

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
	"music/internal/modules/catalog/album"
	"music/internal/modules/catalog/artist"
	"music/internal/modules/catalog/genre"
	"music/internal/modules/catalog/track"
)

type Handlers struct {
	Artist *artist.Handler
	Album  *album.Handler
	Track  *track.Handler
	Genre  *genre.Handler
}

func RegisterPublicRoutes(rg *gin.RouterGroup, h Handlers, mw ...gin.HandlerFunc) {
	cg := rg.Group("/catalog")
	if len(mw) > 0 {
		cg.Use(mw...)
	}
	{
		cg.GET("/artists", h.Artist.List)
		cg.GET("/artists/:artistID", h.Artist.Get)
		cg.GET("/artists/:artistID/overview", h.Artist.Overview)
		cg.GET("/artists/:artistID/tracks", h.Artist.Tracks)
		cg.GET("/artists/:artistID/albums", h.Artist.Albums)
		cg.GET("/artists/:artistID/singles", h.Artist.Singles)
		cg.GET("/artists/:artistID/appears-on", h.Artist.AppearsOn)
		cg.GET("/artists/:artistID/top-tracks", h.Artist.TopTracks)
		cg.GET("/artists/:artistID/related", h.Artist.Related)

		cg.GET("/albums", h.Album.List)
		cg.GET("/albums/:albumID", h.Album.Get)
		cg.GET("/albums/:albumID/tracks", h.Album.Tracks)
		cg.GET("/albums/:albumID/artists", h.Album.Artists)

		cg.GET("/tracks", h.Track.ListPublic)
		cg.GET("/tracks/:trackID", h.Track.Get)
		cg.GET("/tracks/:trackID/credits", h.Track.Credits)
		cg.GET("/tracks/:trackID/artists", h.Track.Artists)

		cg.GET("/genres", h.Genre.List)
		cg.GET("/genres/:genreID", h.Genre.Get)
	}
}

func RegisterAdminRoutes(rg *gin.RouterGroup, h Handlers, authMW gin.HandlerFunc) {
	adminCatalog := rg.Group("/admin/catalog")
	adminCatalog.Use(authMW)
	adminCatalog.Use(auth.RequireRole("admin"))

	artists := adminCatalog.Group("/artist")
	{
		artists.GET("", h.Artist.List)
		artists.POST("", h.Artist.Create)
		artists.PATCH("/:artistID", h.Artist.Update)
		artists.DELETE("/:artistID", h.Artist.Delete)
		artists.PUT("/:artistID/related", h.Artist.ReplaceRelated)
		artists.PUT("/:artistID/top-tracks", h.Artist.ReplaceTopTracks)

	}

	albums := adminCatalog.Group("/albums")
	{
		albums.GET("", h.Album.List)
		albums.POST("", h.Album.Create)
		albums.PATCH("/:albumID", h.Album.Update)
		albums.DELETE("/:albumID", h.Album.Delete)
		albums.PUT("/:albumID/artists", h.Album.ReplaceArtists)
	}

	tracks := adminCatalog.Group("/tracks")
	{
		tracks.GET("", h.Track.ListAdmin)
		tracks.POST("", h.Track.Create)
		tracks.PATCH("/:trackID", h.Track.Update)
		tracks.DELETE("/:trackID", h.Track.Delete)
		tracks.PUT("/:trackID/credits", h.Track.ReplaceCredits)
		tracks.PUT("/:trackID/artists", h.Track.ReplaceArtists)
	}

	genres := adminCatalog.Group("/genres")
	{
		genres.GET("", h.Genre.List)
		genres.POST("", h.Genre.Create)
		genres.PATCH("/:genreID", h.Genre.Update)
		genres.DELETE("/:genreID", h.Genre.Delete)
	}
}
