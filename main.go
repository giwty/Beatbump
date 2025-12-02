package main

import (
	"beatbump-server/backend/api"
	"beatbump-server/backend/api/downloader"
	"beatbump-server/backend/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db.InitDB()
	downloader.StartWorker()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
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

	// Download & Settings
	e.GET("/api/v1/download/playlist", api.DownloadPlaylistHandler)
	e.GET("/api/v1/download/song", api.DownloadSongMixHandler)
	e.GET("/api/v1/downloads", api.GetDownloadsHandler)
	e.GET("/api/v1/downloads/:taskId/tracks", api.GetTaskTracksHandler)
	e.POST("/api/v1/downloads/:taskId/pause", api.PauseTaskHandler)
	e.POST("/api/v1/downloads/:taskId/resume", api.ResumeTaskHandler)
	e.POST("/api/v1/downloads/:taskId/retry", api.RetryTaskHandler)
	e.GET("/api/v1/stream/:taskId/:videoId", api.StreamTrackHandler)
	e.GET("/api/v1/settings", api.GetSettingsHandler)
	e.POST("/api/v1/settings", api.UpdateSettingsHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
