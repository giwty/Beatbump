package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	// Playlist ID provided by the user
	videoId := "_fW2rw8SwoA"
	paramsMap := map[string]string{}
	responseBytes, err := api.Next(videoId, "", api.WebMusic, paramsMap)

	var nextResponse _youtube.NextResponse
	err = json.Unmarshal(responseBytes, &nextResponse)
	
	parsedResponse := ParseNextBody(nextResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, parsedResponse)

}
