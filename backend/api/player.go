package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"beatbump-server/backend/_youtube/api/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
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
	authObj := (c.(*auth.AuthContext)).AuthContext
	var potRequired bool = false
	//authorization := headers.Get("Authorization")
	if videoId == "" {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Missing required param: videoId"))
	}

	var responseBytes []byte
	var err error
	if authObj.AuthType == auth.AUTH_TYPE_OAUTH {
		responseBytes, err = callPlayerAPI(api.IOS_MUSIC, videoId, playlistId, authObj)
	} else if authObj.AuthType == auth.AUTH_TYPE_COOKIES {
		//default cookies - web_creator / tv
		responseBytes, err = callPlayerAPI(api.TV, videoId, playlistId, authObj)
	} else {
		// no auth clients - tv / ios
		potRequired = true
		responseBytes, err = callPlayerAPI(api.TV_SIMPLY, videoId, playlistId, authObj)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var playerResponse _youtube.PlayerResponse
	err = json.Unmarshal(responseBytes, &playerResponse)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to parse reeponse: "+err.Error()+" (using:"+getAuthTypeString(authObj)+")")
	}

	if playerResponse.PlayabilityStatus.Status != "OK" {
		return c.JSON(http.StatusInternalServerError, "Playability status is not OK: "+playerResponse.PlayabilityStatus.Status+" (using:"+getAuthTypeString(authObj)+")")
	}

	if len(playerResponse.StreamingData.AdaptiveFormats) == 0 {
		return c.JSON(http.StatusInternalServerError, "Playability status is not OK: "+playerResponse.PlayabilityStatus.Status+" (using:"+getAuthTypeString(authObj)+")")
	}

	for i := 0; i < len(playerResponse.StreamingData.AdaptiveFormats); i++ {
		format := &playerResponse.StreamingData.AdaptiveFormats[i]
		/*if !strings.Contains(format.MimeType, "audio") {
			continue
		}*/
		streamUrl := format.URL

		if streamUrl == "" {
			signatureCipher := format.SignatureCipher
			streamUrl, err = api.DecipherSignatureCipher(videoId, signatureCipher)
		} else {
			resultUrl, err := url.Parse(streamUrl)
			if err != nil {
				continue
			}
			resultUrl, err = api.DecipherNsig(resultUrl, videoId)
			streamUrl = resultUrl.String()
		}

		if potRequired {
			poToken := api.GetPoToken()
			if poToken != nil {
				streamUrl += "&pot=" + poToken.Potoken
			}
		}

		if err != nil {
			fmt.Println("failed to decrypt - %s", err)
			continue
		}
		format.URL = strings.Clone(streamUrl)
	}

	return c.JSON(http.StatusOK, playerResponse)
}

func getAuthTypeString(authObj *auth.Auth) string {
	if authObj.AuthType == auth.AUTH_TYPE_OAUTH {
		return "OAuth"
	}
	if authObj.AuthType == auth.AUTH_TYPE_COOKIES {
		return "Cookies"
	}
	return "No Auth configured"
}

func callPlayerAPI(clientInfo api.ClientInfo, videoId string, playlistId string, authObj *auth.Auth) ([]byte, error) {

	responseBytes, err := api.Player(videoId, playlistId, clientInfo, nil, authObj)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error building API request: %s", err))
	}

	return responseBytes, err
}
