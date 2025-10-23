package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RelatedEndpointHandler(c echo.Context) error {

	url := c.Request().URL
	query := url.Query()
	browseId := query.Get("browseId")
	if browseId == "" {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Missing required param: browseId"))
	}

	responseBytes, err := api.Browse(browseId, api.PageType_MusicPageTypeTrackRelated, "", nil, nil, nil, api.WebMusic)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	var relatedResponse _youtube.RelatedResponse
	err = json.Unmarshal(responseBytes, &relatedResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}
	if len(relatedResponse.Contents.SectionListRenderer.Contents) == 0 {
		return nil
	}
	var carouselResponse []Carousel = make([]Carousel, 0)
	var description Description = Description{}
	for _, section := range relatedResponse.Contents.SectionListRenderer.Contents {
		musicShelf := Carousel{}
		if len(section.MusicDescriptionShelfRenderer.Header.Runs) != 0 {
			description.Header = section.MusicDescriptionShelfRenderer.Header.Runs[0].Text
		}
		if len(section.MusicDescriptionShelfRenderer.Description.Runs) != 0 {
			description.Description = section.MusicDescriptionShelfRenderer.Description.Runs[0].Text
		}
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
			//

			contents = append(contents, item)
		}
		if len(contents) != 0 {
			musicShelf.Contents = contents
			carouselResponse = append(carouselResponse, musicShelf)
		}
	}

	r := struct {
		Carousels   []Carousel  `json:"carousels"`
		Description Description `json:"description"`
	}{
		Carousels:   carouselResponse,
		Description: description,
	}

	return c.JSON(http.StatusOK, r)
}

type Description struct {
	Header      string `json:"header"`
	Description string `json:"description"`
}

func parseMusicMultiRowItemRenderer(itemRender _youtube.MusicMultiRowListItemRenderer) IListItemRenderer {
	item := IListItemRenderer{}
	thumbnails := itemRender.Thumbnail.MusicThumbnailRenderer.Thumbnail.Thumbnails
	itemThumbnails := make([]Thumbnail, 0, len(thumbnails))
	for _, thumbnail := range thumbnails {
		itemThumbnails = append(itemThumbnails, Thumbnail{
			URL:    thumbnail.URL,
			Width:  thumbnail.Width,
			Height: thumbnail.Height,
			//Placeholder: thumbnail.Placeholder,
		})
	}

	/*playlistId := itemRender.TrackingParams[0]..WatchEndpoint.PlaylistId
	  if playlistId == "" {
	  	playlistId = itemRender.ThumbnailOverlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchPlaylistEndpoint.PlaylistId
	  }*/

	item.Title = itemRender.Title.Runs[0].Text
	item.Thumbnails = itemThumbnails
	//	item.AspectRatio = &itemRender.AspectRatio
	item.VideoId = &itemRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchEndpoint.VideoID
	item.PlaylistId = &itemRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchEndpoint.PlaylistID
	item.MusicVideoType = itemRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchEndpoint.WatchEndpointMusicSupportedConfigs.WatchEndpointMusicConfig.MusicVideoType
	item.PlayerParams = itemRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchEndpoint.Params

	/*if itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId != "" ||
		itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType != "" {
		endpoint := Endpoint{}
		endpoint.BrowseId = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId
		endpoint.PageType = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType
		item.Endpoint = &endpoint
	}*/
	itemArtists := make([]Artist, 0, len(itemRender.Subtitle.Runs))
	for _, subtitleItem := range itemRender.Subtitle.Runs {
		/*		if subtitleItem.NavigationEndpoint == nil {
				continue
			}*/
		itemArtists = append(itemArtists, Artist{
			Text: subtitleItem.Text,
			//BrowseId: subtitleItem.Text.BrowseEndpoint.BrowseId,
			//PageType: subtitleItem.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType,
		})
	}
	item.Subtitle = itemArtists

	return item
}

func parseMusicTwoRowItemRenderer(itemRender _youtube.MusicTwoRowItemRenderer) IListItemRenderer {
	item := IListItemRenderer{}
	itemThumbnails := parseThumbnails(itemRender)

	playlistId := itemRender.NavigationEndpoint.WatchEndpoint.PlaylistId
	if playlistId == "" {
		playlistId = itemRender.ThumbnailOverlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchPlaylistEndpoint.PlaylistId
	}

	item.Title = itemRender.Title.Runs[0].Text
	item.Thumbnails = itemThumbnails
	item.AspectRatio = &itemRender.AspectRatio
	item.VideoId = &itemRender.NavigationEndpoint.WatchEndpoint.VideoId
	item.PlaylistId = &playlistId
	item.MusicVideoType = itemRender.NavigationEndpoint.WatchEndpoint.WatchEndpointMusicSupportedConfigs.WatchEndpointMusicConfig.MusicVideoType
	item.PlayerParams = itemRender.NavigationEndpoint.WatchEndpoint.Params

	if itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId != "" ||
		itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType != "" {
		endpoint := Endpoint{}
		endpoint.BrowseId = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId
		endpoint.PageType = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType
		item.Endpoint = &endpoint
	}
	itemArtists := make([]Artist, 0, len(itemRender.Subtitle.Runs))
	for _, subtitleItem := range itemRender.Subtitle.Runs {
		/*		if subtitleItem.NavigationEndpoint == nil {
				continue
			}*/
		itemArtists = append(itemArtists, Artist{
			Text:     subtitleItem.Text,
			BrowseId: subtitleItem.NavigationEndpoint.BrowseEndpoint.BrowseId,
			PageType: subtitleItem.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType,
		})
	}
	item.Subtitle = itemArtists

	return item

}

func parseThumbnails(itemRender _youtube.MusicTwoRowItemRenderer) []Thumbnail {
	thumbnails := itemRender.ThumbnailRenderer.MusicThumbnailRenderer.Thumbnail.Thumbnails
	itemThumbnails := make([]Thumbnail, 0, len(thumbnails))
	for _, thumbnail := range thumbnails {
		itemThumbnails = append(itemThumbnails, Thumbnail{
			URL:    thumbnail.Url,
			Width:  thumbnail.Width,
			Height: thumbnail.Height,
			//Placeholder: thumbnail.Placeholder,
		})
	}
	return itemThumbnails
}
