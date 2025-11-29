


package downloader

import (
	"testing"
)

func TestGetPlaylist(t *testing.T) {
	// Playlist ID provided by the user
	playlistID := "VLPLciALUPf8sIJlvH350O8PWh4UbSbOffSJ"

	// Initial call to GetPlaylist
	response, err := fetchPlaylistTracks(playlistID)
	if err != nil {
		t.Errorf("Error fetching playlist: %v", err)
	}
	t.Logf("Playlist response: %v", response)
}	