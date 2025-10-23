package main

import (
	"beatbump-server/backend/api"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	companionBaseURL := os.Getenv("COMPANION_BASE_URL")
	companionAPIKey := os.Getenv("COMPANION_API_KEY")

	if companionBaseURL == "" || companionAPIKey == "" {
		panic("COMPANION_BASE_URL and COMPANION_API_KEY must be set")
	}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(LoggingMiddleware)
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "./build",
		Browse:     true,
		IgnoreBase: true,
		HTML5:      true,
	}))

	e.GET("/api/v1/search.json", api.SearchEndpointHandler)
	e.GET("/api/v1/player.json", api.PlayerEndpointHandler)
	e.GET("/api/v1/playlist.json", api.PlaylistEndpointHandler)
	e.GET("/api/v1/next.json", api.NextEndpointHandler)
	e.GET("/api/v1/related.json", api.RelatedEndpointHandler)
	e.GET("/api/v1/main.json", api.AlbumEndpointHandler)
	e.GET("/api/v1/get_queue.json", api.GetQueueHandler)
	e.GET("/api/v1/get_search_suggestions.json", api.GetSearchSuggstionsHandler)

	e.GET("/api/v1/home.json", api.HomeEndpointHandler)
	e.GET("/api/v1/explore/:category", api.ExploreEndpointHandler)
	e.GET("/api/v1/explore", api.ExploreEndpointHandler)
	e.GET("/api/v1/trending", api.TrendingEndpointHandler)
	e.GET("/api/v1/trending/:browseId", api.TrendingEndpointHandler)

	e.GET("/api/v1/artist/:artistId", api.ArtistEndpointHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
