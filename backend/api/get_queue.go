package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetQueueHandler(c echo.Context) error {
	url := c.Request().URL
	query := url.Query()

	playlistId := query.Get("playlistId")
	videoId := query.Get("videoIds")

	//{"context":{"client":{"clientName":"WEB_REMIX","clientVersion":"1.20230501.01.00","visitorData":"CgtybFRXTFh0anpCVSi45tyvBjIKCgJVUxIEGgAgGg=="},"request":{"useSsl":true},"user":{"lockedSafetyMode":false},"captionParams":{},"clickTracking":{"clickTrackingParams":"CCYQxyAiEwjT0qbK8_uEAxUR1pcIHVXXCtI="}},"enablePersistentPlaylistPanel":true,"isAudioOnly":true,"tunerSettingValue":"AUTOMIX_SETTING_NORMAL","videoId":"YL5VIiyTVso","playlistSetVideoId":"289F4A46DF0A30D2","playlistId":"RDAMVMgQudqB29sT8","params":"gAQBiAQB"}
	if videoId == "" && playlistId == "" {
		return errors.New("missing required param: videoId")
	}
	responseBytes, err := api.GetQueue(videoId, playlistId, api.WebMusic)

	var nextResponse _youtube.NextResponse
	err = json.Unmarshal(responseBytes, &nextResponse)
	if err != nil {
		return errors.New(fmt.Sprintf("Error building API request: %s", err))
	}

	parsedResponse := ParseNextBody(nextResponse)
	return c.JSON(http.StatusOK, parsedResponse)

}
