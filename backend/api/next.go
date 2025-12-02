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

type NextEndpointResponse struct {
	Results             []Item  `json:"results"`
	Continuation        string  `json:"continuation"`
	ClickTrackingParams string  `json:"clickTrackingParams"`
	CurrentMixID        string  `json:"currentMixId"`
	VisitorData         string  `json:"visitorData"`
	Related             Related `json:"related"`
}

type Related struct {
	BrowseID                              string `json:"browseId"`
	BrowseEndpointContextSupportedConfigs struct {
		BrowseEndpointContextMusicConfig struct {
			PageType string `json:"pageType"`
		} `json:"browseEndpointContextMusicConfig"`
	} `json:"browseEndpointContextSupportedConfigs"`
}

func NextEndpointHandler(c echo.Context) error {
	url := c.Request().URL
	query := url.Query()

	clickTracking := query.Get("clickTracking")
	continuation := query.Get("continuation")
	loggingContext := query.Get("loggingContext")
	params := query.Get("params")
	playlistId := query.Get("playlistId")
	playlistSetVideoId := query.Get("playlistSetVideoId")
	videoId := query.Get("videoId")
	//index := query.Get("index")
	visitorData := query.Get("visitorData")
	paramsMap := map[string]string{}

	paramsMap["visitorData"] = visitorData
	paramsMap["params"] = params
	paramsMap["playlistSetVideoId"] = playlistSetVideoId
	paramsMap["loggingContext"] = loggingContext
	paramsMap["clickTracking"] = clickTracking
	paramsMap["continuation"] = continuation
	//{"context":{"client":{"clientName":"WEB_REMIX","clientVersion":"1.20230501.01.00","visitorData":"CgtybFRXTFh0anpCVSi45tyvBjIKCgJVUxIEGgAgGg=="},"request":{"useSsl":true},"user":{"lockedSafetyMode":false},"captionParams":{},"clickTracking":{"clickTrackingParams":"CCYQxyAiEwjT0qbK8_uEAxUR1pcIHVXXCtI="}},"enablePersistentPlaylistPanel":true,"isAudioOnly":true,"tunerSettingValue":"AUTOMIX_SETTING_NORMAL","videoId":"YL5VIiyTVso","playlistSetVideoId":"289F4A46DF0A30D2","playlistId":"RDAMVMgQudqB29sT8","params":"gAQBiAQB"}
	if videoId == "" && playlistId == "" && playlistSetVideoId == "" {
		return errors.New("missing required param: videoId")
	}

	responseBytes, err := api.Next(videoId, playlistId, api.WebMusic, paramsMap)

	var nextResponse _youtube.NextResponse
	err = json.Unmarshal(responseBytes, &nextResponse)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing response: %s", err))
	}

	if continuation == "" {
		parsedResponse := ParseNextBody(nextResponse)
		return c.JSON(http.StatusOK, parsedResponse)
	} /*else {
		parsedResponse = parseContinuationNextBody(nextResponse)
	}*/

	return c.JSON(http.StatusOK, struct{}{})
}

func ParseNextBody(nextResponse _youtube.NextResponse) NextEndpointResponse {
	tabs := nextResponse.Contents.SingleColumnMusicWatchNextResultsRenderer.TabbedRenderer.WatchNextTabbedResultsRenderer.Tabs
	if len(tabs) < 3 {
		return NextEndpointResponse{}
	}

	related := tabs[2].TabRenderer.Endpoint.BrowseEndpoint
	contents := tabs[0].TabRenderer.Content.MusicQueueRenderer.Content.PlaylistPanelRenderer.Contents

	if len(contents) == 0 {
		return NextEndpointResponse{}
	}
	//nextRadioContinuationData := tabs[0].TabRenderer.Content.MusicQueueRenderer.Content.PlaylistPanelRenderer.c
	//clickTrackingParams := nextResponse.c
	//continuation := nextResponse.
	//watchEndpoint := nextResponse.CurrentVideoEndpoint.WatchEndpoint
	visitorData := nextResponse.ResponseContext.VisitorData
	mixId := nextResponse.CurrentVideoEndpoint.WatchEndpoint.PlaylistId

	var results []Item
	for _, content := range contents {
		if len(content.PlaylistPanelVideoRenderer.Title.Runs) == 0 {
			continue
		}
		item := parsePlaylistPanelVideoRenderer(content.PlaylistPanelVideoRenderer)
		results = append(results, item)
	}

	return NextEndpointResponse{
		Results:     results,
		VisitorData: visitorData,
		Related: Related{
			BrowseEndpointContextSupportedConfigs: related.BrowseEndpointContextSupportedConfigs,
			BrowseID:                              related.BrowseId,
		},
		CurrentMixID: mixId,
		//	ClickTrackingParams:
	}
}

func parsePlaylistPanelVideoRenderer(playlistVideoRender _youtube.PlaylistPanelVideoRenderer) Item {
	var item Item = Item{}

	if len(playlistVideoRender.LongBylineText.Runs) != 0 {
		itemArtists := make([]Artist, 0, len(playlistVideoRender.LongBylineText.Runs))
		for _, sub := range playlistVideoRender.LongBylineText.Runs {
			itemArtists = append(itemArtists, Artist{
				Text:     sub.Text,
				BrowseId: sub.NavigationEndpoint.BrowseEndpoint.BrowseId,
				PageType: sub.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType,
			})
		}
		item.Subtitle = itemArtists
		artist := Artist{Text: playlistVideoRender.LongBylineText.Runs[0].Text,
			BrowseId: playlistVideoRender.LongBylineText.Runs[0].NavigationEndpoint.BrowseEndpoint.BrowseId,
			PageType: playlistVideoRender.LongBylineText.Runs[0].NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType,
		}
		item.ArtistInfo = ArtistInfo{Artist: []Artist{artist}}

	}

	thumbnails := playlistVideoRender.Thumbnail.Thumbnails
	itemThumbnails := make([]Thumbnail, 0, len(thumbnails))
	for _, thumbnail := range thumbnails {
		itemThumbnails = append(itemThumbnails, Thumbnail{
			URL:    thumbnail.Url,
			Width:  thumbnail.Width,
			Height: thumbnail.Height,
		})
	}
	item.Thumbnails = itemThumbnails

	item.VideoID = playlistVideoRender.NavigationEndpoint.WatchEndpoint.VideoId
	item.PlaylistID = playlistVideoRender.NavigationEndpoint.WatchEndpoint.PlaylistId
	item.PlaylistSetVideoID = playlistVideoRender.NavigationEndpoint.WatchEndpoint.PlaylistSetVideoId
	if item.PlaylistSetVideoID == "" {
		item.PlaylistSetVideoID = playlistVideoRender.PlaylistSetVideoId
	}
	item.PlayerParams = playlistVideoRender.NavigationEndpoint.WatchEndpoint.PlayerParams
	item.Itct = playlistVideoRender.NavigationEndpoint.WatchEndpoint.Params
	item.Index = playlistVideoRender.NavigationEndpoint.WatchEndpoint.Index
	if len(playlistVideoRender.Title.Runs) > 0 {
		item.Title = playlistVideoRender.Title.Runs[0].Text
	}

	if len(playlistVideoRender.Menu.MenuRenderer.Items) > 0 {
		item.AutoMixList = playlistVideoRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchEndpoint.PlaylistId
	}

	if len(playlistVideoRender.LengthText.Runs) > 0 {
		item.Length = playlistVideoRender.LengthText.Runs[0].Text
	}

	item.ClickTrackingParams = playlistVideoRender.NavigationEndpoint.ClickTrackingParams

	//todo: filter out albums

	return item
}

type Item struct {
	Subtitle            []Artist    `json:"subtitle"`
	Thumbnails          []Thumbnail `json:"thumbnails"`
	ArtistInfo          ArtistInfo  `json:"artistInfo"`
	VideoID             string      `json:"videoId"`
	PlaylistID          string      `json:"playlistId"`
	PlaylistSetVideoID  string      `json:"playlistSetVideoId"`
	PlayerParams        string      `json:"playerParams"`
	Itct                string      `json:"itct"`
	Index               int         `json:"index"`
	Title               string      `json:"title"`
	AutoMixList         string      `json:"autoMixList"`
	Length              string      `json:"length"`
	ClickTrackingParams string      `json:"clickTrackingParams"`
}

