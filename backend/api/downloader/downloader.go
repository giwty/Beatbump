package downloader

import (
	"beatbump-server/backend/_youtube"
	yt_api "beatbump-server/backend/_youtube/api"
	"beatbump-server/backend/api"
	"beatbump-server/backend/db"
	"beatbump-server/backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

type TrackInfo struct {
	VideoID      string
	Title        string
	Artist       string
	Album        string
	ThumbnailURL string
}

func PopulateGroupTask(playlistID string, groupTaskID int) {
	log.Printf("Populating songs for group task %d (Playlist %s)", groupTaskID, playlistID)

	// Phase 1: Populate Tracks if not already populated
	existingSongs, err := db.GetSongTasks(groupTaskID)
	if err != nil {
		log.Printf("Failed to get song tasks: %v", err)
		db.UpdateGroupTaskStatus(groupTaskID, db.TaskStatusFailed)
		return
	}

	if len(existingSongs) == 0 {
		tracks, err := fetchPlaylistTracks(playlistID)
		if err != nil {
			log.Printf("Failed to fetch playlist tracks: %v", err)
			db.UpdateGroupTaskStatus(groupTaskID, db.TaskStatusFailed)
			return
		}

		for _, track := range tracks {
			if track.VideoID == "" {
				continue
			}
			err := db.AddSongTask(groupTaskID, track.VideoID, track.Title, track.Artist, track.Album, track.ThumbnailURL)
			if err != nil {
				log.Printf("Failed to add song %s to task: %v", track.Title, err)
			}
		}
		log.Printf("Populated %d songs for group task %d", len(tracks), groupTaskID)
	} else {
		log.Printf("Group task %d already has %d songs", groupTaskID, len(existingSongs))
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

func downloadTrack(track TrackInfo, downloadPath, playlistFolder, companionBaseURL string, companionAPIKey *string) (string, error) {
	log.Printf("Downloading %s - %s", track.Artist, track.Title)

	// 1. Get Stream Info
	streamUrl, contentLength, err := getStreamInfo(track.VideoID, companionBaseURL, companionAPIKey)
	if err != nil {
		return "", err
	}

	// 2. Build filename and path (always .m4a)
	baseFilename := fmt.Sprintf("%s - %s", track.Artist, track.Title)
	filename := sanitizeFilename(baseFilename + ".m4a")
	finalFilePath := filepath.Join(downloadPath, filename)

	// 3. Create output file
	downloadTarget, err := os.Create(finalFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	// 4. Perform Download
	err = performDownload(streamUrl, downloadTarget, contentLength)
	downloadTarget.Close() // Close immediately after download

	if err != nil {
		os.Remove(finalFilePath)
		return "", err
	}

	log.Printf("Finished downloading %s - %s", track.Artist, track.Title)

	// Return absolute path to the downloaded file
	return finalFilePath, nil
}

func getStreamInfo(videoID, companionBaseURL string, companionAPIKey *string) (string, int64, error) {
	responseBytes, err := yt_api.Player(videoID, "", yt_api.IOS_MUSIC, nil, companionBaseURL, companionAPIKey)
	if err != nil {
		return "", 0, err
	}

	var playerResponse _youtube.PlayerResponse
	err = json.Unmarshal(responseBytes, &playerResponse)
	if err != nil {
		return "", 0, err
	}

	if playerResponse.PlayabilityStatus.Status != "OK" {
		return "", 0, fmt.Errorf("not playable: %s", playerResponse.PlayabilityStatus.Status)
	}

	var streamUrl string
	var contentLength int64
	bestBitrate := 0

	for _, format := range playerResponse.StreamingData.AdaptiveFormats {
		if strings.HasPrefix(format.MimeType, "audio") {
			if format.Bitrate > bestBitrate {
				bestBitrate = format.Bitrate
				streamUrl = format.URL
				contentLength, _ = strconv.ParseInt(format.ContentLength, 10, 64)
			}
		}
	}

	if streamUrl == "" {
		return "", 0, fmt.Errorf("no audio stream found")
	}

	return streamUrl, contentLength, nil
}

func performDownload(streamUrl string, output *os.File, contentLength int64) error {
	req, err := http.NewRequest(http.MethodGet, streamUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	r, w := io.Pipe()

	if contentLength == 0 {
		resp, err := http.Get(streamUrl)
		if err != nil {
			return fmt.Errorf("failed to download stream: %v", err)
		}

		go func() {
			defer resp.Body.Close()
			_, err := io.Copy(w, resp.Body)
			if err == nil {
				w.Close()
			} else {
				w.CloseWithError(err)
			}
		}()
	} else {
		downloadChunked(req, w, contentLength)
	}
	defer r.Close()

	_, err = io.Copy(output, r)
	if err != nil {
		return fmt.Errorf("failed to save stream: %v", err)
	}
	return nil
}

func convertTrack(track TrackInfo, inputM4aPath, downloadPath string) (string, error) {
	log.Printf("Converting %s - %s to MP3", track.Artist, track.Title)

	baseFilename := fmt.Sprintf("%s - %s", track.Artist, track.Title)
	mp3Filename := sanitizeFilename(baseFilename + ".mp3")
	mp3FilePath := filepath.Join(downloadPath, mp3Filename)

	// 1. Fetch Enhanced Metadata (Best Effort)
	enhancedMeta, err := utils.FetchMetadata(track.Artist, track.Title)
	if err != nil {
		log.Printf("Metadata fetch failed for %s - %s: %v", track.Artist, track.Title, err)
		// Fallback to basic metadata
		enhancedMeta = &utils.AudioMetadata{
			Title:  track.Title,
			Artist: track.Artist,
			Album:  track.Album,
		}
	} else {
		log.Printf("Enhanced metadata found: %s - %s (%s)", enhancedMeta.Artist, enhancedMeta.Title, enhancedMeta.Album)
	}

	// 2. Determine Cover Art
	var coverPath string

	// Priority: iTunes Artwork > YouTube Thumbnail
	artworkURL := enhancedMeta.ArtworkURL
	if artworkURL == "" {
		artworkURL = track.ThumbnailURL
	}

	if artworkURL != "" {
		cf, err := os.CreateTemp("", "beatbump_cover_*.jpg")
		if err == nil {
			coverPath = cf.Name()
			cf.Close()
			if err := downloadFile(artworkURL, coverPath); err != nil {
				os.Remove(coverPath)
				coverPath = ""
			} else {
				defer os.Remove(coverPath)
			}
		}
	}

	// 3. Convert
	err = utils.ConvertToMp3(inputM4aPath, mp3FilePath, coverPath, *enhancedMeta)
	if err != nil {
		log.Printf("Conversion failed for %s: %v", track.Title, err)
		return "", fmt.Errorf("conversion failed: %v", err)
	}

	// 4. Remove original .m4a file after successful conversion
	os.Remove(inputM4aPath)
	log.Printf("Conversion complete: %s", mp3Filename)

	return mp3FilePath, nil
}

const (
	Size1Kb  = 1024
	Size1Mb  = Size1Kb * 1024
	Size10Mb = Size1Mb * 10
)

func downloadChunked(req *http.Request, w *io.PipeWriter, contentLength int64) {
	chunks := getChunks(contentLength, Size10Mb)

	currentChunk := atomic.Uint32{}

	for i := 0; i < 5; i++ {
		go func() {
			for {
				chunkIndex := int(currentChunk.Add(1)) - 1
				if chunkIndex >= len(chunks) {
					// no more chunks
					return
				}

				chunk := &chunks[chunkIndex]
				err := downloadChunk(req.Clone(context.Background()), chunk)
				close(chunk.data)

				if err != nil {
					w.CloseWithError(err)
					return
				}
			}
		}()
	}

	go func() {
		// copy chunks into the PipeWriter
		for i := 0; i < len(chunks); i++ {
			data := <-chunks[i].data
			_, err := io.Copy(w, bytes.NewBuffer(data))
			if err != nil {
				w.CloseWithError(err)
			}
		}

		// everything succeeded
		w.Close()
	}()
}

func downloadChunk(req *http.Request, chunk *chunk) error {
	q := req.URL.Query()
	q.Set("range", fmt.Sprintf("%d-%d", chunk.start, chunk.end))
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	expected := int(chunk.end-chunk.start) + 1
	data, err := io.ReadAll(resp.Body)
	n := len(data)

	if err != nil {
		return err
	}

	if n != expected {
		return fmt.Errorf("chunk at offset %d has invalid size: expected=%d actual=%d", chunk.start, expected, n)
	}

	chunk.data <- data

	return nil
}

type chunk struct {
	start int64
	end   int64
	data  chan []byte
}

func getChunks(totalSize, chunkSize int64) []chunk {
	var chunks []chunk

	for start := int64(0); start < totalSize; start += chunkSize {
		end := chunkSize + start - 1
		if end > totalSize-1 {
			end = totalSize - 1
		}

		chunks = append(chunks, chunk{start, end, make(chan []byte, 1)})
	}

	return chunks
}

func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return -1
		}
		return r
	}, name)
}

func generateM3U(playlistName string, tracks []db.SongTask, outputDir string) error {
	filename := filepath.Join(outputDir, "playlist.m3u8")
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
			// FilePath in DB is relative to downloads root (e.g. "Playlist/Song.m4a")
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

func HandleSongTask(track *db.SongTask) {
	// Update status to in_progress
	db.UpdateSongTaskStatus(int(track.GroupTaskID), track.VideoID, db.TaskStatusProcessing)

	// Fetch parent group task to get keys and playlist name
	groupTask, err := db.GetGroupTask(int(track.GroupTaskID))
	if err != nil {
		log.Printf("Failed to get group task %d: %v", track.GroupTaskID, err)
		db.UpdateSongTaskStatus(int(track.GroupTaskID), track.VideoID, db.TaskStatusFailed)
		return
	}

	fullDownloadPath, playlistFolder, playlistName := resolveDownloadDirectory(groupTask)

	// Ensure download directory exists
	if _, err := os.Stat(fullDownloadPath); os.IsNotExist(err) {
		os.MkdirAll(fullDownloadPath, 0755)
	}

	// Check for free disk space (require at least 50MB)
	freeSpace, err := utils.GetFreeDiskSpace(fullDownloadPath)
	if err != nil {
		log.Printf("Failed to check disk space: %v", err)
		// Proceeding anyway as it might be a permission issue or unsupported OS
	} else {
		const minSpace = 50 * 1024 * 1024 // 50MB
		if freeSpace < minSpace {
			log.Printf("Insufficient disk space: %d bytes available, %d required", freeSpace, minSpace)
			db.UpdateSongTaskStatus(int(track.GroupTaskID), track.VideoID, db.TaskStatusFailed)
			return
		}
	}

	trackInfo := TrackInfo{
		VideoID:      track.VideoID,
		Title:        track.Title,
		Artist:       track.Artist,
		Album:        track.Album,
		ThumbnailURL: track.ThumbnailURL,
	}

	// Step 1: Download the track (always as .m4a)
	absolutePath, err := downloadTrack(trackInfo, fullDownloadPath, playlistFolder, groupTask.CompanionBaseURL, &groupTask.CompanionAPIKey)
	if err != nil {
		log.Printf("Failed to download track %s: %v", track.Title, err)
		db.UpdateSongTaskStatus(int(track.GroupTaskID), track.VideoID, db.TaskStatusFailed)
		return
	}

	// Step 2: Convert to MP3 if FFmpeg is available
	finalPath := absolutePath
	if utils.IsFFmpegAvailable() {
		convertedPath, err := convertTrack(trackInfo, absolutePath, fullDownloadPath)
		if err != nil {
			log.Printf("Conversion failed for %s: %v. Keeping .m4a file.", track.Title, err)
			// Continue with .m4a file (don't fail the task)
		} else {
			finalPath = convertedPath
		}
	}

	// Step 3: Finalize the task
	relativePath := filepath.Join(playlistFolder, filepath.Base(finalPath))
	finalizeTask(track, relativePath, playlistName, fullDownloadPath)
}

func resolveDownloadDirectory(groupTask *db.GroupTask) (string, string, string) {
	downloadPath, _ := db.GetSetting(db.DownloadPathSetting)

	playlistName := groupTask.PlaylistName
	if groupTask.Type == db.TaskTypeOngoingDownload && playlistName == "" {
		// Use safe format for folder name (no colons)
		playlistName = fmt.Sprintf("Songs %s", groupTask.CreatedAt.Format("2006-01-02_15-04"))
	}

	playlistFolder := sanitizeFilename(playlistName)
	fullDownloadPath := filepath.Join(downloadPath, playlistFolder)
	return fullDownloadPath, playlistFolder, playlistName
}

func finalizeTask(track *db.SongTask, relativePath, playlistName, fullDownloadPath string) {
	db.MarkSongTaskCompleted(int(track.GroupTaskID), track.VideoID, relativePath)

	// Check if all songs in the group are completed
	completed, err := db.CheckGroupCompletion(int(track.GroupTaskID))
	if err == nil && completed {
		log.Printf("All songs completed for group task %d. Finalizing...", track.GroupTaskID)
		db.UpdateGroupTaskStatus(int(track.GroupTaskID), db.TaskStatusCompleted)

		// Generate Metadata
		finalSongs, err := db.GetSongTasks(int(track.GroupTaskID))
		if err == nil {
			if err := generateM3U(playlistName, finalSongs, fullDownloadPath); err != nil {
				log.Printf("Failed to generate M3U: %v", err)
			}
			if err := generateNFO(playlistName, fullDownloadPath); err != nil {
				log.Printf("Failed to generate NFO: %v", err)
			}
		}
	}
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
