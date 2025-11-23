package downloader

import (
	"beatbump-server/backend/_youtube"
	yt_api "beatbump-server/backend/_youtube/api"
	"beatbump-server/backend/api"
	"beatbump-server/backend/db"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type TrackInfo struct {
	VideoID      string
	Title        string
	Artist       string
	Album        string
	ThumbnailURL string
}

func DownloadPlaylist(playlistID string, taskID int) {
	log.Printf("Starting processing for playlist %s (Task %d)", playlistID, taskID)
	err := db.UpdateTaskStatus(taskID, "processing")
	if err != nil {
		log.Printf("Failed to update task status: %v", err)
		return
	}

	downloadPath, err := db.GetSetting("download_path")
	if err != nil || downloadPath == "" {
		downloadPath = "downloads"
	}

	// Get playlist name from task payload if available
	var playlistName string
	task, err := db.GetTask(taskID)
	if err == nil && task.Payload != "" {
		var payloadMap map[string]interface{}
		if err := json.Unmarshal([]byte(task.Payload), &payloadMap); err == nil {
			if name, ok := payloadMap["playlistName"].(string); ok {
				playlistName = name
			}
		}
	}

	// Create playlist folder
	playlistFolder := playlistID
	if playlistName != "" {
		playlistFolder = sanitizeFilename(playlistName)
	}
	fullDownloadPath := filepath.Join(downloadPath, playlistFolder)

	// Ensure download directory exists
	if _, err := os.Stat(fullDownloadPath); os.IsNotExist(err) {
		os.MkdirAll(fullDownloadPath, 0755)
	}

	// Phase 1: Populate Tracks if not already populated
	hasTracks, err := db.HasTaskTracks(taskID)
	if err != nil {
		log.Printf("Failed to check task tracks: %v", err)
		db.UpdateTaskStatus(taskID, "failed")
		return
	}

	if !hasTracks {
		log.Printf("Phase 1: Populating tracks for task %d", taskID)
		tracks, err := fetchPlaylistTracks(playlistID)
		if err != nil {
			log.Printf("Failed to fetch playlist tracks: %v", err)
			db.UpdateTaskStatus(taskID, "failed")
			return
		}

		for _, track := range tracks {
			if track.VideoID == "" {
				continue
			}
			err := db.AddTaskTrack(taskID, track.VideoID, track.Title, track.Artist, track.Album, track.ThumbnailURL)
			if err != nil {
				log.Printf("Failed to add track %s to task: %v", track.Title, err)
			}
		}
		log.Printf("Phase 1 Complete: Populated %d tracks", len(tracks))
	}

	// Phase 2: Process Downloads
	// Fetch all tracks again to get their current status
	existingTracks, err := db.GetTaskTracks(taskID)
	if err != nil {
		log.Printf("Failed to get task tracks: %v", err)
		db.UpdateTaskStatus(taskID, "failed")
		return
	}

	log.Printf("Phase 2: Processing %d tracks for task %d", len(existingTracks), taskID)

	allCompleted := true
	for _, track := range existingTracks {
		// Check if task is paused
		currentTask, err := db.GetTask(taskID)
		if err != nil {
			log.Printf("Failed to get task status: %v", err)
			return
		}
		if currentTask.Status == "paused" {
			log.Printf("Task %d is paused, stopping download", taskID)
			break
		}

		// Skip if already completed
		if track.Status == "completed" {
			continue
		}

		// Update status to in_progress
		db.UpdateTaskTrackStatus(taskID, track.VideoID, "in_progress")

		trackInfo := TrackInfo{
			VideoID: track.VideoID,
			Title:   track.Title,
			Artist:  track.Artist,
			Album:   track.Album,
		}

		relativePath, err := downloadTrack(trackInfo, fullDownloadPath, playlistFolder)
		if err != nil {
			log.Printf("Failed to download track %s: %v", track.Title, err)
			db.UpdateTaskTrackStatus(taskID, track.VideoID, "failed")
			allCompleted = false
			// Continue to next track
		} else {
			db.MarkTrackCompleted(taskID, track.VideoID, relativePath)
		}

		// Sleep to avoid rate limiting (2-8 seconds)
		sleepDuration := time.Duration(rand.Intn(7)+2) * time.Second
		log.Printf("Sleeping for %v...", sleepDuration)
		time.Sleep(sleepDuration)
	}

	if allCompleted {
		db.UpdateTaskStatus(taskID, "completed")
		log.Printf("Completed download for playlist %s", playlistID)
	} else {
		// Mark as completed even if some failed, so it doesn't get stuck in processing loop forever.
		// Users can retry later if we implement a retry mechanism.
		db.UpdateTaskStatus(taskID, "completed")
		log.Printf("Finished processing playlist %s (some tracks may have failed)", playlistID)
	}

	// Generate Metadata files (M3U and NFO)
	// Fetch fresh tracks to ensure we have file paths for completed ones
	finalTracks, err := db.GetTaskTracks(taskID)
	if err == nil {
		if err := generateM3U(playlistName, finalTracks, fullDownloadPath); err != nil {
			log.Printf("Failed to generate M3U: %v", err)
		}
		if err := generateNFO(playlistName, fullDownloadPath); err != nil {
			log.Printf("Failed to generate NFO: %v", err)
		}
	}
}

func fetchPlaylistTracks(playlistID string) ([]TrackInfo, error) {
	playlistResponse, err := api.GetPlaylist(playlistID, "", "")
	if err != nil {
		return nil, err
	}

	var tracks []TrackInfo

	for _, item := range playlistResponse.Tracks {
		title := item.Title

		var videoId string
		if item.VideoId != nil {
			videoId = *item.VideoId
		}

		var artist string
		if len(item.ArtistInfo.Artist) > 0 {
			artist = item.ArtistInfo.Artist[0].Text
		}

		var album string
		if item.Album != nil {
			album = item.Album.Text
		}

		var thumbnailURL string
		if len(item.Thumbnails) > 0 {
			// Get the last thumbnail (usually highest quality)
			thumbnailURL = item.Thumbnails[len(item.Thumbnails)-1].URL
		}

		if videoId != "" {
			tracks = append(tracks, TrackInfo{
				VideoID:      videoId,
				Title:        title,
				Artist:       artist,
				Album:        album,
				ThumbnailURL: thumbnailURL,
			})
		}
	}

	return tracks, nil
}

func downloadTrack(track TrackInfo, downloadPath, playlistFolder string) (string, error) {
	log.Printf("Downloading %s - %s", track.Artist, track.Title)

	// Get Stream URL
	// Use yt_api.Player instead of api.Player because api is now the beatbump api package
	responseBytes, err := yt_api.Player(track.VideoID, "", yt_api.IOS_MUSIC, nil, "", nil)
	if err != nil {
		return "", err
	}

	var playerResponse _youtube.PlayerResponse
	err = json.Unmarshal(responseBytes, &playerResponse)
	if err != nil {
		return "", err
	}

	if playerResponse.PlayabilityStatus.Status != "OK" {
		return "", fmt.Errorf("not playable: %s", playerResponse.PlayabilityStatus.Status)
	}

	var streamUrl string
	// Find best audio format
	// Prefer high bitrate audio
	bestBitrate := 0
	for _, format := range playerResponse.StreamingData.AdaptiveFormats {
		if strings.HasPrefix(format.MimeType, "audio") {
			if format.Bitrate > bestBitrate {
				bestBitrate = format.Bitrate
				streamUrl = format.URL
			}
		}
	}

	if streamUrl == "" {
		return "", fmt.Errorf("no audio stream found")
	}

	// Sanitize filename
	filename := fmt.Sprintf("%s - %s.mp3", track.Artist, track.Title)
	filename = sanitizeFilename(filename)
	filePath := filepath.Join(downloadPath, filename)

	// Use ffmpeg to download and convert
	cmd := exec.Command("ffmpeg", "-y", "-i", streamUrl,
		"-metadata", "title="+track.Title,
		"-metadata", "artist="+track.Artist,
		"-metadata", "album="+track.Album,
		"-c:a", "libmp3lame", "-q:a", "2",
		filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	// Return relative path including playlist folder
	return filepath.Join(playlistFolder, filename), nil
}

func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return -1
		}
		return r
	}, name)
}

func generateM3U(playlistName string, tracks []db.TaskTrack, outputDir string) error {
	filename := filepath.Join(outputDir, "playlist.m3u")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header
	if _, err := f.WriteString("#EXTM3U\n"); err != nil {
		return err
	}

	for _, track := range tracks {
		if track.Status == "completed" && track.FilePath != "" {
			// FilePath in DB is relative to downloads root (e.g. "Playlist/Song.mp3")
			// We need just the filename for the M3U since it's in the same folder
			baseName := filepath.Base(track.FilePath)

			// Write EXTINF
			duration := -1 // Unknown duration
			title := fmt.Sprintf("%s - %s", track.Artist, track.Title)
			if _, err := f.WriteString(fmt.Sprintf("#EXTINF:%d,%s\n", duration, title)); err != nil {
				return err
			}
			// Write filename
			if _, err := f.WriteString(baseName + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateNFO(playlistName string, outputDir string) error {
	filename := filepath.Join(outputDir, "album.nfo")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<album>
  <title>%s</title>
  <artist>Various Artists</artist>
  <albumartist>Various Artists</albumartist>
  <compilation>true</compilation>
</album>`, playlistName)

	_, err = f.WriteString(content)
	return err
}
