package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"strings"
)

var searchFilters = map[string]string{
	"all":                 "",
	"songs":               "EgWKAQIIAWoKEAMQBBAJEAoQBQ%3D%3D",
	"videos":              "EgWKAQIQAWoKEAMQBBAJEAoQBQ%3D%3D",
	"albums":              "EgWKAQIYAWoKEAMQBBAJEAoQBQ%3D%3D",
	"artists":             "EgWKAQIgAWoKEAMQBBAJEAoQBQ%3D%3D",
	"community_playlists": "EgeKAQQoAEABagwQDhAKEAkQAxAEEAU%3D",
	"featured_playlists":  "EgeKAQQoADgBagwQDhAKEAMQBBAJEAU%3D",
	"all_playlists":       "EgWKAQIoAWoKEAMQBBAKEAUQCQ%3D%3D",
}

func SearchEndpointHandler(c echo.Context) error {
	urlQuery := c.Request().URL.Query()
	query := urlQuery.Get("q")
	filter := urlQuery.Get("filter")
	filterId := ""
	if filter != "all" && filter != "" {
		filterId = searchFilters[filter]
	}

	ctoken := urlQuery.Get("ctoken")
	itct := urlQuery.Get("itct")

	if query == "" && itct == "" && ctoken == "" {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Missing required params"))
	}

	queryUnescape, err := url.QueryUnescape(query)

	var responseBytes []byte
	if itct != "" && ctoken != "" {
		responseBytes, err = api.Search(queryUnescape, filterId, &itct, &ctoken, api.WebMusic)
	} else {
		responseBytes, err = api.Search(queryUnescape, filterId, nil, nil, api.WebMusic)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	var searchResponse _youtube.SearchResponse
	err = json.Unmarshal(responseBytes, &searchResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	var regularResponse []MusicShelf
	var continuationResponse []IListItemRenderer
	var continuation _youtube.NextContinuationData
	var responseType *string = nil
	// continuation mode
	if len(searchResponse.ContinuationContents.MusicShelfContinuation.Continuations) != 0 {
		searchContinuationContent := searchResponse.ContinuationContents.MusicShelfContinuation
		continuationResponse, err = parseContinuationResponse(searchContinuationContent.Content, filter)
		if len(searchContinuationContent.Continuations) == 1 {
			continuation = searchContinuationContent.Continuations[0].NextContinuationData
		}

		responseType = stringPtr("next")

	} else if len(searchResponse.Content.TabbedSearchResultsRenderer.Tabs) != 0 {
		searchContent := searchResponse.Content.TabbedSearchResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.SectionListRendererContents
		regularResponse, err = parseResponse(searchContent)
		if len(searchContent) == 1 && len(searchContent[0].MusicShelfRenderer.Continuations) != 0 {
			continuation = searchContent[0].MusicShelfRenderer.Continuations[0].NextContinuationData
		}
	} else {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	if continuationResponse != nil {
		r := struct {
			ContinuationResults []IListItemRenderer            `json:"results"`
			Response            _youtube.SearchResponse        `json:"response"`
			Continuation        *_youtube.NextContinuationData `json:"continuation,omitempty"`
			Type                *string                        `json:"type,omitempty"`
		}{
			ContinuationResults: continuationResponse,
			Response:            searchResponse,
			Continuation:        &continuation,
			Type:                responseType,
		}
		return c.JSON(http.StatusOK, r)
	} else {
		r := struct {
			Results      []MusicShelf                   `json:"results"`
			Response     _youtube.SearchResponse        `json:"response"`
			Continuation *_youtube.NextContinuationData `json:"continuation,omitempty"`
			Type         *string                        `json:"type,omitempty"`
		}{
			Results:      regularResponse,
			Response:     searchResponse,
			Continuation: &continuation,
			Type:         responseType,
		}
		return c.JSON(http.StatusOK, r)
	}

}

func parseContinuationResponse(content []_youtube.MusicShelfContinuationContent, filter string) ([]IListItemRenderer, error) {

	response := make([]IListItemRenderer, 0, len(content))
	for _, entry := range content {
		item := parseMusicResponsiveListItemRenderer(entry.MusicResponsiveListItemRenderer)
		item.Type = filter
		response = append(response, item)
	}

	return response, nil
}

func parseResponse(content []_youtube.SectionListRendererContents) ([]MusicShelf, error) {
	response := make([]MusicShelf, 0, len(content))
	for _, shelf := range content {
		currShelf := MusicShelf{}
		if shelf.MusicShelfRenderer == nil {
			continue
		}

		title := ""
		if len(shelf.MusicShelfRenderer.Title.Runs) != 0{
			title = shelf.MusicShelfRenderer.Title.Runs[0].Text
		}
		
		currShelf.Header.Title = title
		currShelf.Contents = make([]IListItemRenderer, 0, len(shelf.MusicShelfRenderer.Contents))
		for _, entry := range shelf.MusicShelfRenderer.Contents {
			item := parseMusicResponsiveListItemRenderer(entry.MusicResponsiveListItemRenderer)

			entryTitle := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
			item.Type = entryTitle
			if entryTitle == "top_result" && item.Endpoint != nil {
				if strings.Contains(item.Endpoint.PageType, "SINGLE") ||
					strings.Contains(item.Endpoint.PageType, "ALBUM") {
					item.Type = "albums"
				}
			}

			currShelf.Contents = append(currShelf.Contents, item)
		}

		response = append(response, currShelf)
	}
	return response, nil
}
