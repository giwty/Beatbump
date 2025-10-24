package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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

	for i := 0; i < len(playerResponse.StreamingData.AdaptiveFormats); i++ {
		format := &playerResponse.StreamingData.AdaptiveFormats[i]
		/*if !strings.Contains(format.MimeType, "audio") {
			continue
		}*/
		streamUrl := format.URL

		format.URL = strings.Clone(streamUrl)
	}

	return c.JSON(http.StatusOK, playerResponse)
}



func callPlayerAPI(clientInfo api.ClientInfo, videoId string, playlistId string, companionBaseURL string, companionAPIKey *string) ([]byte, error) {

	responseBytes, err := api.Player(videoId, playlistId, clientInfo, nil, companionBaseURL, companionAPIKey)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error building API request: %s", err))
	}

	return responseBytes, err
}
