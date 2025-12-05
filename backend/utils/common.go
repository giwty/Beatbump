package utils

import (
	"beatbump-server/backend/db"
	"fmt"
	"path/filepath"
	"strings"
)

func SanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return -1
		}
		return r
	}, name)
}

func ResolveDownloadDirectory(groupTask *db.GroupTask, downloadPath string) (string, string, string) {
	playlistName := groupTask.PlaylistName
	if groupTask.Type == db.TaskTypeOngoingDownload {
		if playlistName == "" {
			// Use safe format for folder name (no colons)
			// Format: Songs YYYY-MM-DD_HH-MM
			playlistName = fmt.Sprintf("Songs %s", groupTask.CreatedAt.Format("2006-01-02_15-04"))
		}

		// For ongoing downloads, place them in a dedicated subfolder
		// We sanitize the playlist name first, then join with the parent folder
		safePlaylistName := SanitizeFilename(playlistName)
		playlistFolder := filepath.Join("Ongoing Listening", safePlaylistName)
		fullDownloadPath := filepath.Join(downloadPath, playlistFolder)

		return fullDownloadPath, playlistFolder, playlistName
	}

	if groupTask.Type == db.TaskTypeSongMixDownload {
		// Format: {Song Title}-{MaxTracks}
		// If MaxTracks is 0, maybe just {Song Title}?
		// The requirement said: "All the songs will be stored under the same folder with the song title then dash and the number of related songs"
		folderName := fmt.Sprintf("%s mix-(%d songs)", playlistName, groupTask.MaxTracks)
		playlistFolder := SanitizeFilename(folderName)
		fullDownloadPath := filepath.Join(downloadPath, playlistFolder)
		return fullDownloadPath, playlistFolder, playlistName
	}

	playlistFolder := SanitizeFilename(playlistName)
	fullDownloadPath := filepath.Join(downloadPath, playlistFolder)
	return fullDownloadPath, playlistFolder, playlistName
}
