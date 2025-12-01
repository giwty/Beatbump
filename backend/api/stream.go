package api

import (
	"beatbump-server/backend/db"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func StreamTrackHandler(c echo.Context) error {
	taskIdParam := c.Param("taskId")
	videoId := c.Param("videoId")

	taskId, err := strconv.Atoi(taskIdParam)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid taskId")
	}

	// Get SongTask to find the file path
	songTask, err := db.GetSongTask(taskId, videoId)
	if err != nil {
		return c.String(http.StatusNotFound, "Song not found")
	}

	if songTask.FilePath == "" {
		return c.String(http.StatusNotFound, "File not found for this song")
	}

	// Get Download Path Setting
	downloadPath, _ := db.GetSetting(db.DownloadPathSetting)

	// Construct full path
	fullPath := filepath.Join(downloadPath, songTask.FilePath)

	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		return c.String(http.StatusNotFound, "File not found on disk")
	}
	defer file.Close()

	// Get file info for ServeContent
	fileInfo, err := file.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not get file info")
	}

	// ServeContent handles Range requests automatically
	http.ServeContent(c.Response(), c.Request(), fileInfo.Name(), fileInfo.ModTime(), file)
	return nil
}
