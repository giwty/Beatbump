package api

import (
	"beatbump-server/backend/db"
	"fmt"
	"net/http"
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

	payload := ""
	if playlistName != "" {
		payload = fmt.Sprintf(`{"playlistName": "%s"}`, playlistName)
	}

	companionAPIKey := c.Request().Header.Get("x-companion-api-key")
	companionBaseURL := c.Request().Header.Get("x-companion-base-url")

	err := db.AddGroupTask("playlist_download", playlistID, payload, companionAPIKey, companionBaseURL, "user")
	if err != nil {
		c.Logger().Errorf("Failed to add group task: %v", err)
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
	downloadPath, _ := db.GetSetting("download_path")
	ongoingListeningEnabled, _ := db.GetSetting("ongoing_listening_enabled")
	return c.JSON(http.StatusOK, map[string]string{
		"downloadPath":            downloadPath,
		"ongoingListeningEnabled": ongoingListeningEnabled,
	})
}

func UpdateSettingsHandler(c echo.Context) error {
	companionAPIKey := c.Request().Header.Get("x-companion-api-key")
	companionBaseURL := c.Request().Header.Get("x-companion-base-url")

	if companionBaseURL == "" || companionAPIKey == "" {
		return c.String(http.StatusBadRequest, "Missing companion API configuration headers")
	}

	var req SettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if req.DownloadPath != "" {
		err := db.SetSetting("download_path", req.DownloadPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update download path")
		}
	}

	if req.OngoingListeningEnabled != "" {
		err := db.SetSetting("ongoing_listening_enabled", req.OngoingListeningEnabled)
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
