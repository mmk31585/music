package catalog

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/identity"
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

func RegisterPublicRoutes(rg *gin.RouterGroup, h Handlers) {
	cg := rg.Group("/catalog")
	{
		cg.GET("/artists", h.Artist.List)
		cg.GET("/artists/:artistID", h.Artist.Get)

		cg.GET("/albums", h.Album.List)
		cg.GET("/albums/:albumID", h.Album.Get)

		cg.GET("/tracks", h.Track.ListPublic)
		cg.GET("/tracks/:trackID", h.Track.Get)

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
	}

	albums := adminCatalog.Group("/albums")
	{
		albums.GET("", h.Album.List)
		albums.POST("", h.Album.Create)
		albums.PATCH("/:albumID", h.Album.Update)
		albums.DELETE("/:albumID", h.Album.Delete)
	}

	tracks := adminCatalog.Group("/tracks")
	{
		tracks.GET("", h.Track.ListAdmin)
		tracks.POST("", h.Track.Create)
		tracks.PATCH("/:trackID", h.Track.Update)
		tracks.DELETE("/:trackID", h.Track.Delete)
	}

	genres := adminCatalog.Group("/genres")
	{
		genres.GET("", h.Genre.List)
		genres.POST("", h.Genre.Create)
		genres.PATCH("/:genreID", h.Genre.Update)
		genres.DELETE("/:genreID", h.Genre.Delete)
	}
}
