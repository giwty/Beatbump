package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func ExploreEndpointHandler(c echo.Context) error {
	category := c.Param("category")
	browseID := "FEmusic_moods_and_genres"
	if category != "" {
		browseID = "FEmusic_moods_and_genres_category"
	}
	var responseBytes []byte
	var err error
	responseBytes, err = api.Browse(browseID, api.PageType_MusicPageTypePlaylist, category, nil, nil, nil, api.WebMusic)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	if category == "" {

		var exploreResponse _youtube.Explore
		err = json.Unmarshal(responseBytes, &exploreResponse)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
		}

		r := parseExplore(exploreResponse)
		return c.JSON(http.StatusOK, r)
	} else {
		var homeResponse _youtube.HomeResponse
		err = json.Unmarshal(responseBytes, &homeResponse)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
		}

		r := ParseHome(homeResponse)

		return c.JSON(http.StatusOK, r)
	}
}

func parseExplore(response _youtube.Explore) []map[string]interface{} {

	result := []map[string]interface{}{}

	if len(response.Contents.SingleColumnBrowseResultsRenderer.Tabs) == 0 {
		return result
	}
	if len(response.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) == 0 {
		return result
	}

	for _, section := range response.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents {
		items := section.GridRenderer.Items
		header := section.GridRenderer.Header
		sectionItems := []ExploreCategory{}

		for _, item := range items {

			itemToAdd := ExploreCategory{
				Text:  item.MusicNavigationButtonRenderer.ButtonText.Runs[0].Text,
				Color: "#" + strconv.FormatInt(item.MusicNavigationButtonRenderer.Solid.LeftStripeColor&0xffffff, 16),
				Endpoint: Endpoint{
					Params:   stringPtr(item.MusicNavigationButtonRenderer.ClickCommand.BrowseEndpoint.Params),
					BrowseId: item.MusicNavigationButtonRenderer.ClickCommand.BrowseEndpoint.BrowseID,
				},
			}
			sectionItems = append(sectionItems, itemToAdd)
		}
		result = append(result, map[string]interface{}{
			"title":   header.GridHeaderRenderer.Title.Runs[0].Text,
			"section": sectionItems,
		})

	}
	return result
}

type ExploreCategory struct {
	Text     string   `json:"text"`
	Color    string   `json:"color"`
	Endpoint Endpoint `json:"endpoint"`
}
