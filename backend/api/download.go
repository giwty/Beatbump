package api

import (
	"beatbump-server/backend/db"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DownloadPlaylistRequest struct {
	PlaylistID   string `json:"playlistId"`
	PlaylistName string `json:"playlistName"`
}
type SettingsRequest struct {
	DownloadPath            string `json:"downloadPath"`
	OngoingListeningEnabled string `json:"ongoingListeningEnabled"`
}

func DownloadPlaylistHandler(c echo.Context) error {
	playlistID := c.QueryParam("playlistId")
	playlistName := c.QueryParam("playlistName")

	if playlistID == "" {
		return c.String(http.StatusBadRequest, "Playlist ID is required")
	}

	err := db.AddGroupTask(db.TaskTypePlaylistDownload, playlistID, playlistName, "user", -1)
	if err != nil {
		c.Logger().Errorf("Failed to add group task: %v", err)
		if len(err.Error()) > 24 && err.Error()[:24] == "UNIQUE constraint failed" {
			return c.JSON(http.StatusConflict, map[string]string{"status": "already_queued", "message": "Task already queued"})
		}
		return c.String(http.StatusInternalServerError, "Failed to create task")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "queued"})
}

func DownloadSongMixHandler(c echo.Context) error {
	videoID := c.QueryParam("videoId")
	limitStr := c.QueryParam("limit")

	if videoID == "" {
		return c.String(http.StatusBadRequest, "Video ID is required")
	}

	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid limit")
		}
	}
	if limit < 0 {
		limit = 0
	}
	if limit > 500 {
		limit = 500
	}

	// Use videoID as the referenceID
	referenceID := "songmix:" + videoID

	err := db.AddGroupTask(db.TaskTypeSongMixDownload, referenceID, "Song mix "+videoID, "user", limit)
	if err != nil {
		c.Logger().Errorf("Failed to add song mix task: %v", err)
		if len(err.Error()) > 24 && err.Error()[:24] == "UNIQUE constraint failed" {
			return c.JSON(http.StatusConflict, map[string]string{"status": "already_queued", "message": "Task already queued"})
		}
		return c.String(http.StatusInternalServerError, "Failed to create task")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "queued"})
}

func GetDownloadsHandler(c echo.Context) error {
	tasks, err := db.GetAllGroupTasks()
	if err != nil {
		c.Logger().Errorf("Failed to fetch group tasks: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
	}
	// Return empty array if nil to avoid null in JSON
	if tasks == nil {
		tasks = []db.GroupTask{}
	}
	c.Logger().Infof("Returning %d tasks", len(tasks))
	return c.JSON(http.StatusOK, tasks)
}

func GetTaskTracksHandler(c echo.Context) error {
	taskIDStr := c.Param("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task ID")
	}

	songs, err := db.GetSongTasks(taskID)
	if err != nil {
		c.Logger().Errorf("Failed to fetch song tasks: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch song tasks")
	}
	if songs == nil {
		songs = []db.SongTask{}
	}
	return c.JSON(http.StatusOK, songs)
}

func GetSettingsHandler(c echo.Context) error {
	downloadPath, _ := db.GetSetting(db.DownloadPathSetting)
	ongoingListeningEnabled, _ := db.GetSetting(db.OngoingListeningEnabledSetting)
	return c.JSON(http.StatusOK, map[string]string{
		"downloadPath":            downloadPath,
		"ongoingListeningEnabled": ongoingListeningEnabled,
	})
}

func UpdateSettingsHandler(c echo.Context) error {

	var req SettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if req.DownloadPath != "" {
		// Validate that the path exists and is a directory
		info, err := os.Stat(req.DownloadPath)
		if err != nil {
			if os.IsNotExist(err) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Directory does not exist"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate directory"})
		}
		if !info.IsDir() {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Path is not a directory"})
		}

		err = db.SetSetting(db.DownloadPathSetting, req.DownloadPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update download path")
		}
	}

	if req.OngoingListeningEnabled != "" {
		if req.OngoingListeningEnabled == "true" {
			// Ensure we have a valid download path before enabling
			currentPath, _ := db.GetSetting(db.DownloadPathSetting)
			// If we are updating both, use the new one, otherwise use the stored one
			if req.DownloadPath != "" {
				currentPath = req.DownloadPath
			}

			if currentPath == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Download path must be set first"})
			}

			info, err := os.Stat(currentPath)
			if err != nil || !info.IsDir() {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Valid download path required"})
			}
		}

		err := db.SetSetting(db.OngoingListeningEnabledSetting, req.OngoingListeningEnabled)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update ongoing listening setting")
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func PauseTaskHandler(c echo.Context) error {
	taskIDStr := c.Param("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task ID")
	}

	err = db.UpdateGroupTaskStatus(taskID, "paused")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to pause task")
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "paused"})
}

func ResumeTaskHandler(c echo.Context) error {
	taskIDStr := c.Param("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task ID")
	}

	err = db.UpdateGroupTaskStatus(taskID, "pending")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to resume task")
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "resumed"})
}

func RetryTaskHandler(c echo.Context) error {
	taskIDStr := c.Param("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task ID")
	}

	err = db.RetryGroupTask(taskID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retry task")
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "retrying"})
}
