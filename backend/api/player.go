package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"beatbump-server/backend/db"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type PlayerAPIResponse struct {
	// Define your response structure based on the actual data fields
}

func PlayerEndpointHandler(c echo.Context) error {
	requestUrl := c.Request().URL
	query := requestUrl.Query()
	videoId := query.Get("videoId")
	playlistId := query.Get("playlistId")
	//playerParams := query.Get("playerParams")
	if videoId == "" {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Missing required param: videoId"))
	}

	var responseBytes []byte
	var err error

	companionAPIKey := c.Request().Header.Get("x-companion-api-key")
	companionBaseURL := c.Request().Header.Get("x-companion-base-url")

	if companionBaseURL == "" || companionAPIKey == "" {
		return c.String(http.StatusInternalServerError, "Missing companion API configuration headers")
	}

	responseBytes, err = callPlayerAPI(api.IOS_MUSIC, videoId, playlistId, companionBaseURL, &companionAPIKey)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var playerResponse _youtube.PlayerResponse
	err = json.Unmarshal(responseBytes, &playerResponse)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to parse reeponse: "+err.Error())
	}

	if playerResponse.PlayabilityStatus.Status != "OK" {
		return c.JSON(http.StatusInternalServerError, "Playability status is not OK: "+playerResponse.PlayabilityStatus.Status)
	}

	if len(playerResponse.StreamingData.AdaptiveFormats) == 0 {
		return c.JSON(http.StatusInternalServerError, "Playability status is not OK: "+playerResponse.PlayabilityStatus.Status)
	}

	//258/251/22/256/140/250/18/249/139
	for i := 0; i < len(playerResponse.StreamingData.AdaptiveFormats); i++ {
		format := &playerResponse.StreamingData.AdaptiveFormats[i]
		/*if !strings.Contains(format.MimeType, "audio") {
			continue
		}*/
		streamUrl := format.URL

		format.URL = strings.Clone(streamUrl)
	}

	// Ongoing Listening Logic
	handleTrackdownloadTask(playlistId, companionAPIKey, companionBaseURL, playerResponse, videoId)

	return c.JSON(http.StatusOK, playerResponse)
}

func handleTrackdownloadTask(playlistId string, companionAPIKey string, companionBaseURL string, playerResponse _youtube.PlayerResponse, videoId string) {
	go func() {
		enabled, _ := db.GetSetting(db.OngoingListeningEnabledSetting)
		if enabled == "true" {
			// Check if task exists
			var refID string
			var playlistName string
			var task *db.GroupTask

			// Treat "RD" (Radio) playlists as part of the ongoing session, not as a separate playlist
			// Also handle "undefined" string which can come from frontend
			if playlistId != "" && playlistId != "undefined" &&
				!(strings.HasPrefix(playlistId, "RDEM") || strings.HasPrefix(playlistId, "RDAMVM") || strings.HasPrefix(playlistId, "RDAT")) {
				refID = "ongoing:playlist:" + playlistId
				// Try to get task to see if we already have the name
				t, err := db.GetGroupTaskByReferenceID(refID)
				if err == nil && t != nil {
					task = t
				} else {
					if strings.HasPrefix(playlistId, "R") && !strings.HasPrefix(playlistId, "VL") {
						playlistId = "VL" + playlistId
					}
					// Task doesn't exist, fetch playlist name
					pl, err := GetPlaylist(playlistId, "", "")
					if err == nil {
						if title, ok := pl.Header["title"].([]string); ok && len(title) > 0 {
							playlistName = title[0]
						} else {
							playlistName = "Playlist " + playlistId
						}
					} else {
						playlistName = "Playlist " + playlistId
					}
				}
			} else {
				// Try to find an active session task (updated within last 30 mins)
				log.Println("Checking for active session task...")
				activeTask, err := db.GetActiveSessionGroupTask(30 * time.Minute)
				if err == nil && activeTask != nil {
					log.Printf("Found active session task: %d (Ref: %s)", activeTask.ID, activeTask.ReferenceID)
					task = activeTask
					refID = task.ReferenceID
				} else {
					if err != nil {
						log.Printf("Error getting active session task: %v", err)
					} else {
						log.Println("No active session task found.")
					}
					// Create new session
					timestamp := time.Now().Format("2006-01-02 15:04")
					refID = fmt.Sprintf("ongoing:songs:%d", time.Now().Unix())
					playlistName = fmt.Sprintf("Listening Session %s", timestamp)
					log.Printf("Creating new session: %s (Ref: %s)", playlistName, refID)
				}
			}

			if task == nil {
				// Create task
				db.AddGroupTask(db.TaskTypeOngoingDownload, refID, playlistName, companionAPIKey, companionBaseURL, db.TaskSourceSystem)
				task, _ = db.GetGroupTaskByReferenceID(refID)
			}

			if task != nil {
				// Add track
				details := playerResponse.VideoDetails
				title := details.Title
				artist := details.Author
				thumbnail := ""
				if len(details.Thumbnail.Thumbnails) > 0 {
					thumbnail = details.Thumbnail.Thumbnails[len(details.Thumbnail.Thumbnails)-1].URL
				}

				db.AddSongTask(int(task.ID), videoId, title, artist, "", thumbnail)
			}
		}
	}()
}

func callPlayerAPI(clientInfo api.ClientInfo, videoId string, playlistId string, companionBaseURL string, companionAPIKey *string) ([]byte, error) {

	responseBytes, err := api.Player(videoId, playlistId, clientInfo, nil, companionBaseURL, companionAPIKey)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error building API request: %s", err))
	}

	return responseBytes, err
}
