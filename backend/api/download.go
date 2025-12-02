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
	title := c.QueryParam("title")
	limitStr := c.QueryParam("limit")

	if videoID == "" || title == "" {
		return c.String(http.StatusBadRequest, "Video ID and Title are required")
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

	// For song downloads, we use the videoID as the referenceID (or maybe "song:<videoID>" to avoid collision if needed, but videoID should be unique enough for now if we treat it as a playlist of 1+related)
	// Actually, to allow downloading the same song multiple times with different limits or just re-downloading, we might want a unique reference ID or just rely on the fact that if it's pending/processing we don't duplicate.
	// But the user might want to download "Song A" + 5 related, and later "Song A" + 0 related.
	// The current DB constraint is `ReferenceID string 'gorm:"uniqueIndex"'`.
	// So we should probably make the reference ID unique per request if we want to allow multiples, OR just use "song:<videoID>" and enforce one active task per song.
	// Given the playlist logic, let's stick to "song:<videoID>" for now to prevent spamming.
	referenceID := "song:" + videoID

	err := db.AddGroupTask(db.TaskTypeSongMixDownload, referenceID, title, "user", limit)
	if err != nil {
		c.Logger().Errorf("Failed to add song task: %v", err)
		if len(err.Error()) > 24 && err.Error()[:24] == "UNIQUE constraint failed" {
			// If it failed, it might be because there's already a task.
			// For songs, maybe we should just return success if it's already there?
			return c.JSON(http.StatusConflict, map[string]string{"status": "already_queued", "message": "Task already queued"})
		}
		return c.String(http.StatusInternalServerError, "Failed to create task")
	}

	// Add the initial song task immediately
	// We need the group task ID. AddGroupTask doesn't return it.
	// We can fetch it by ReferenceID.
	groupTask, err := db.GetGroupTaskByReferenceID(referenceID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve created task")
	}

	// We need artist/album/thumbnail. These are not passed in the query params usually (or maybe they are?).
	// The frontend should pass them if possible, or we fetch them.
	// Let's assume the frontend passes them in the body or we fetch them.
	// For now, let's try to get them from query params or leave them empty and let the worker fetch/fill them?
	// The worker `PopulateGroupTask` usually fetches playlist tracks.
	// For a single song, we might want to pass metadata.
	// Let's check what the frontend sends. The plan said "Connect UI to backend API".
	// I should probably accept a JSON body for `DownloadSongHandler` to get full metadata.

	// Let's update the handler to read from JSON body if possible, or query params.
	// But `DownloadPlaylistHandler` uses query params.
	// Let's stick to query params for consistency if possible, or switch to JSON.
	// JSON is better for more fields.
	// Let's check `DownloadPlaylistRequest` struct in `download.go`. It exists but `DownloadPlaylistHandler` uses `c.QueryParam`.
	// I will use `c.QueryParam` for now to match style, but add artist/image params.

	artist := c.QueryParam("artist")
	album := c.QueryParam("album")
	thumbnailURL := c.QueryParam("thumbnailUrl")

	err = db.AddSongTask(int(groupTask.ID), videoID, title, artist, album, thumbnailURL)
	if err != nil {
		c.Logger().Errorf("Failed to add initial song task: %v", err)
		// Don't fail the request, the worker might retry or we can handle it.
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
