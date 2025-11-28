package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlaylist(t *testing.T) {
	// Playlist ID provided by the user
	playlistID := "VLPLciALUPf8sIJlvH350O8PWh4UbSbOffSJ"

	// Initial call to GetPlaylist
	response, err := GetPlaylist(playlistID, "", "")
	if !assert.NoError(t, err) {
		return
	}

	totalTracks := len(response.Tracks)
	assert.NotEmpty(t, response.Tracks, "Playlist should have tracks")
	t.Logf("Fetched %d tracks in initial request", len(response.Tracks))

	// Loop for continuations
	token, itct := extractContinuationInfo(response.Continuations)

	pageCount := 1
	maxPages := 10 // Limit to avoid infinite loops

	for token != "" && pageCount < maxPages {
		t.Logf("Fetching page %d with token: %s", pageCount+1, token)

		contResponse, err := GetPlaylist(playlistID, token, itct)
		if !assert.NoError(t, err) {
			break
		}

		fetched := len(contResponse.Tracks)
		totalTracks += fetched
		t.Logf("Fetched %d tracks in page %d", fetched, pageCount+1)

		if fetched == 0 {
			t.Log("No more tracks returned, stopping.")
			break
		}

		token, itct = extractContinuationInfo(contResponse.Continuations)
		pageCount++
	}

	t.Logf("Total tracks fetched: %d", totalTracks)
}

func extractContinuationInfo(continuations interface{}) (string, string) {
	if continuations == nil {
		return "", ""
	}
	importJson, _ := json.Marshal(continuations)
	var contMap map[string]interface{}
	if err := json.Unmarshal(importJson, &contMap); err != nil {
		return "", ""
	}

	var token, itct string

	if t, ok := contMap["token"].(string); ok {
		token = t
	} else if nextContData, ok := contMap["nextContinuationData"].(map[string]interface{}); ok {
		if t, ok := nextContData["continuation"].(string); ok {
			token = t
		}
		if c, ok := nextContData["clickTrackingParams"].(string); ok {
			itct = c
		}
	} else if continuationsArr, ok := contMap["continuations"].([]interface{}); ok {
		if len(continuationsArr) > 0 {
			if firstCont, ok := continuationsArr[0].(map[string]interface{}); ok {
				if t, ok := firstCont["continuation"].(string); ok {
					token = t
				}
				if c, ok := firstCont["clickTrackingParams"].(string); ok {
					itct = c
				}
			}
		}
	} else if t, ok := contMap["continuation"].(string); ok {
		token = t
		if c, ok := contMap["clickTrackingParams"].(string); ok {
			itct = c
		}
	}

	return token, itct
}
