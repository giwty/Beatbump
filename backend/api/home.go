package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func HomeEndpointHandler(c echo.Context) error {
	query := c.Request().URL.Query()
	browseID := "FEmusic_home"
	itct := query.Get("itct")
	ctoken := query.Get("ctoken")
	params := query.Get("params")
	//referrer := query.Get("ref")
	visitorData := query.Get("visitorData")
	var responseBytes []byte
	var err error
	if ctoken != "" && itct != "" {
		responseBytes, err = api.Browse(browseID, api.PageType_MusicPageTypePlaylist, params, &visitorData, &itct, &ctoken, api.WebMusic)
	} else {
		responseBytes, err = api.Browse(browseID, api.PageType_MusicPageTypePlaylist, params, &visitorData, nil, nil, api.WebMusic)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	var homeResponse _youtube.HomeResponse
	err = json.Unmarshal(responseBytes, &homeResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	r := ParseHome(homeResponse)

	return c.JSON(http.StatusOK, r)
}

func ParseHome(homeResponse _youtube.HomeResponse) interface{} {
	var carouselResponse []Carousel = make([]Carousel, 0)
	var response = make(map[string]interface{}, 0)

	if len(homeResponse.Header.MusicHeaderRenderer.Title.Runs) != 0 {
		response["header"] = homeResponse.Header.MusicHeaderRenderer.Title.Runs[0].Text
	}
	//var description Description = Description{}
	if len(homeResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) != 0 {
		for _, section := range homeResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents {
			musicShelf := Carousel{}
			contents := make([]IListItemRenderer, 0)
			categories := make([]ExploreCategory, 0)
			if section.MusicCarouselShelfRenderer != nil {

				//section.MusicTastebuilderShelfRenderer.PrimaryText.Runs[0].Text
				/*if len(section.MusicDescriptionShelfRenderer.Header.Runs) != 0 {
				  	description.Header = section.MusicDescriptionShelfRenderer.Header.Runs[0].Text
				  }
				  if len(section.MusicDescriptionShelfRenderer.Description.Runs) != 0 {
				  	description.Description = section.MusicDescriptionShelfRenderer.Description.Runs[0].Text
				  }*/
				header := section.MusicCarouselShelfRenderer.Header

				if len(header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs) != 0 {
					musicShelf.Header.Title = header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs[0].Text
				}
				if len(section.MusicCarouselShelfRenderer.Header.MusicCarouselShelfBasicHeaderRenderer.Strapline.Runs) != 0 {
					musicShelf.Header.Subheading = stringPtr(header.MusicCarouselShelfBasicHeaderRenderer.Strapline.Runs[0].Text)
				}

				if header.MusicCarouselShelfBasicHeaderRenderer.MoreContentButton.ButtonRenderer.NavigationEndpoint.BrowseEndpoint.Params != "" {
					musicShelf.Header.Params = stringPtr(header.MusicCarouselShelfBasicHeaderRenderer.MoreContentButton.ButtonRenderer.NavigationEndpoint.BrowseEndpoint.Params)
				}

				musicShelf.Header.BrowseId = stringPtr(header.MusicCarouselShelfBasicHeaderRenderer.MoreContentButton.ButtonRenderer.NavigationEndpoint.BrowseEndpoint.BrowseID)

				if strings.Contains(musicShelf.Header.Title, "episodes") {
					continue
				}

				for _, carousel := range section.MusicCarouselShelfRenderer.Contents {
					var item *IListItemRenderer = nil
					var category *ExploreCategory = nil
					if carousel.MusicResponsiveListItemRenderer != nil {
						itemA := parseMusicResponsiveListItemRenderer(*carousel.MusicResponsiveListItemRenderer)
						item = &itemA
					} else if carousel.MusicTwoRowItemRenderer != nil {
						itemA := parseMusicTwoRowItemRenderer(*carousel.MusicTwoRowItemRenderer)
						item = &itemA
					} else if carousel.MusicMultiRowListItemRenderer != nil {
						itemA := parseMusicMultiRowItemRenderer(*carousel.MusicMultiRowListItemRenderer)
						item = &itemA
					} else if carousel.MusicNavigationButtonRenderer != nil {
						category = &ExploreCategory{
							Text:  carousel.MusicNavigationButtonRenderer.ButtonText.Runs[0].Text,
							Color: "#" + strconv.FormatInt(carousel.MusicNavigationButtonRenderer.Solid.LeftStripeColor&0xffffff, 16),
							Endpoint: Endpoint{
								Params:   stringPtr(carousel.MusicNavigationButtonRenderer.ClickCommand.BrowseEndpoint.Params),
								BrowseId: carousel.MusicNavigationButtonRenderer.ClickCommand.BrowseEndpoint.BrowseID,
							},
						}
					}
					if category != nil {
						categories = append(categories, *category)
					}
					if item != nil {
						contents = append(contents, *item)
					}
				}

			} else if section.GridRenderer != nil {
				for _, carousel := range section.GridRenderer.Items {
					if carousel.MusicTwoRowItemRenderer != nil {
						item := parseMusicTwoRowItemRenderer(*carousel.MusicTwoRowItemRenderer)
						contents = append(contents, item)
					}
				}

			} else if section.MusicShelfRenderer != nil {

				if len(section.MusicShelfRenderer.Title.Runs) != 0 {
					musicShelf.Header.Title = section.MusicShelfRenderer.Title.Runs[0].Text
				}

				if section.MusicShelfRenderer.BottomEndpoint.BrowseEndpoint.BrowseId != "" {
					musicShelf.Header.BrowseId = stringPtr(section.MusicShelfRenderer.BottomEndpoint.BrowseEndpoint.BrowseId)
				}
				for _, carousel := range section.MusicShelfRenderer.Contents {

					item := parseMusicResponsiveListItemRenderer(carousel.MusicResponsiveListItemRenderer)

					contents = append(contents, item)
				}

			}
			if len(contents) != 0 {
				musicShelf.Contents = contents
				carouselResponse = append(carouselResponse, musicShelf)
			}
			if len(categories) != 0 {
				musicShelf.Categories = &categories
				carouselResponse = append(carouselResponse, musicShelf)
			}
		}

		chips := homeResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Header.ChipCloudRenderer.Chips
		var chipsResponse []map[string]interface{} = make([]map[string]interface{}, 0)
		for _, chip := range chips {
			text := chip.ChipCloudChipRenderer.Text.Runs[0].Text
			browseEndpoint := chip.ChipCloudChipRenderer.NavigationEndpoint.BrowseEndpoint
			ctoken := chip.ChipCloudChipRenderer.TrackingParams
			/*chipsMap["text"] = text
			  chipsMap["browseEndpoint"] = browseEndpoint
			  chipsMap["ctoken"] = ctoken*/
			chipsResponse = append(chipsResponse, map[string]interface{}{
				"text":           text,
				"browseEndpoint": browseEndpoint,
				"ctoken":         ctoken,
			})
		}
		response["chips"] = chipsResponse
		response["visitorData"] = homeResponse.ResponseContext.VisitorData
		if len(homeResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Continuations) != 0 {
			response["continuations"] = homeResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Continuations[0].NextContinuationData
		}
	}
	if len(homeResponse.ContinuationContents.SectionListContinuation.Contents) != 0 {

		for _, section := range homeResponse.ContinuationContents.SectionListContinuation.Contents {

			musicShelf := Carousel{}

			//section.MusicTastebuilderShelfRenderer.PrimaryText.Runs[0].Text
			/*if len(section.MusicDescriptionShelfRenderer.Header.Runs) != 0 {
			  	description.Header = section.MusicDescriptionShelfRenderer.Header.Runs[0].Text
			  }
			  if len(section.MusicDescriptionShelfRenderer.Description.Runs) != 0 {
			  	description.Description = section.MusicDescriptionShelfRenderer.Description.Runs[0].Text
			  }*/
			header := section.MusicCarouselShelfRenderer.Header
			if len(header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs) != 0 {
				musicShelf.Header.Title = header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs[0].Text
			}
			contents := make([]IListItemRenderer, 0, len(section.MusicCarouselShelfRenderer.Contents))
			for _, carousel := range section.MusicCarouselShelfRenderer.Contents {
				var item IListItemRenderer
				if carousel.MusicResponsiveListItemRenderer != nil {
					item = parseMusicResponsiveListItemRenderer(*carousel.MusicResponsiveListItemRenderer)
				} else {
					item = parseMusicTwoRowItemRenderer(*carousel.MusicTwoRowItemRenderer)
				}

				contents = append(contents, item)
			}
			if len(contents) != 0 {
				musicShelf.Contents = contents
				carouselResponse = append(carouselResponse, musicShelf)
			}
		}
		if len(homeResponse.ContinuationContents.SectionListContinuation.Contents) != 0 && len(homeResponse.ContinuationContents.SectionListContinuation.Continuations) != 0 {
			if homeResponse.ContinuationContents.SectionListContinuation.Continuations[0].NextContinuationData.Continuation != "" {

				response["continuations"] = homeResponse.ContinuationContents.SectionListContinuation.Continuations[0].NextContinuationData
			}
		} else {
			response["continuations"] = struct{}{}
		}
	}

	headerThumbnail := homeResponse.Background.MusicThumbnailRenderer.Thumbnail.Thumbnails
	if len(headerThumbnail) != 0 {
		response["headerThumbnail"] = headerThumbnail
	}
	response["carousels"] = carouselResponse

	return response

}
