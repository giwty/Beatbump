package api

import (
	"beatbump-server/backend/_youtube"
	"strings"
)

type Carousel struct {
	Header struct { // Use pointer for optional header
		Title      string  `json:"title"`
		Subheading *string `json:"subheading,omitempty"`
		BrowseId   *string `json:"browseId,omitempty"`
		Params     *string `json:"params"`
	} `json:"header,omitempty"`
	Contents   []IListItemRenderer `json:"items"`
	Categories *[]ExploreCategory  `json:"categories"`
}

type MusicShelf struct {
	Header struct { // Use pointer for optional header
		Title string `json:"title"`
	} `json:"header,omitempty"`
	Contents []IListItemRenderer `json:"contents"`
}
type Artist struct {
	PageType string `json:"pageType,omitempty"`
	Text     string `json:"text,omitempty"`
	BrowseId string `json:"browseId,omitempty"`
}

// ArtistInfo represents the artist information structure
type ArtistInfo struct {
	Artist []Artist `json:"artist"`
}

// Thumbnail represents information about a thumbnail
type Thumbnail struct {
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Placeholder string `json:"placeholder"`
}

type ContinuationData struct {
	Continuation        string `json:"continuation"`
	ClickTrackingParams string `json:"clickTrackingParams"`
}

type Endpoint struct {
	BrowseId string  `json:"browseId"`
	PageType string  `json:"pageType"`
	Params   *string `json:"params,omitempty"`
}

// IListItemRenderer represents an item in a list
type IListItemRenderer struct {
	Album               *Artist    `json:"album,omitempty"`
	Subtitle            []Artist   `json:"subtitle"`
	ArtistInfo          ArtistInfo `json:"artistInfo"`
	Explicit            bool       `json:"explicit"`
	Title               string     `json:"title"`
	AspectRatio         *string    `json:"aspectRatio,omitempty"`
	PlayerParams        string     `json:"playerParams,omitempty"`
	PlaylistSetVideoId  string     `json:"playlistSetVideoId,omitempty"`
	ClickTrackingParams string     `json:"clickTrackingParams,omitempty"`
	Endpoint            *Endpoint  `json:"endpoint,omitempty"`
	MusicVideoType      string     `json:"musicVideoType,omitempty"`
	Params              string     `json:"params,omitempty"`
	Index               int        `json:"index,omitempty"`
	Length              string     `json:"length,omitempty"`
	VideoId             *string    `json:"videoId,omitempty"`
	PlaylistId          *string    `json:"playlistId,omitempty"`
	LoggingContext      struct {
		VssLoggingContext VssLoggingContext `json:"vssLoggingContext"`
	} `json:"loggingContext,omitempty"`
	Thumbnails []Thumbnail `json:"thumbnails"`
	Type       string      `json:"type,omitempty"`
	BrowseId   string      `json:"browseId,omitempty"`
}

// VssLoggingContext represents the logging context
type VssLoggingContext struct {
	SerializedContextData string `json:"serializedContextData"`
}

func stringPtr(s string) *string {
	return &s
}

func parseMusicResponsiveListItemRenderer(itemRender _youtube.MusicResponsiveListItemRenderer) IListItemRenderer {

	thumbnails := itemRender.Thumbnail.MusicThumbnailRenderer.Thumbnail.Thumbnails
	ItemThumbnails := make([]Thumbnail, 0, len(thumbnails))
	for _, thumbnail := range thumbnails {
		ItemThumbnails = append(ItemThumbnails, Thumbnail{
			URL:    thumbnail.Url,
			Width:  thumbnail.Width,
			Height: thumbnail.Height,
			//Placeholder: thumbnail.Placeholder,
		})
	}

	params := itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchPlaylistEndpoint.Params

	if params == "" {
		params = itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchPlaylistEndpoint.PlaylistId
	}

	item := IListItemRenderer{
		Thumbnails: ItemThumbnails,
		BrowseId:   itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId,
		Params:     itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchPlaylistEndpoint.Params,
		//VideoId:             &itemRender.NavigationEndpoint.WatchEndpoint.VideoId,
		ClickTrackingParams: itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.ClickTrackingParams,
		//LoggingContext:      itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchEndpoint.LoggingContext,
		//	Endpoint:endpoint
	}

	if itemRender.FlexColumnDisplayStyle != "" {
		item.AspectRatio = &itemRender.FlexColumnDisplayStyle
	}

	if len(itemRender.FlexColumns) > 0 {
		subtitleRuns := itemRender.FlexColumns[1].MusicResponsiveListItemFlexColumnRenderer.Text.Runs
		itemArtists := make([]Artist, 0, len(subtitleRuns))
		for _, subtitleItem := range subtitleRuns {
			/*if subtitleItem.NavigationEndpoint.BrowseEndpoint.BrowseId == "" {
				continue
			}*/
			artist := Artist{
				Text:     subtitleItem.Text,
				BrowseId: subtitleItem.NavigationEndpoint.BrowseEndpoint.BrowseId,
				PageType: subtitleItem.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType,
			}

			if strings.Contains(artist.PageType, "ALBUM") {
				item.Album = &artist
			}

			if len(subtitleRuns) > 0 && (subtitleItem.Text == "Artist" || subtitleItem.Text == "Album") {
				artist.PageType = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType
				artist.BrowseId = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId
			}
			itemArtists = append(itemArtists, artist)
		}

		item.Subtitle = itemArtists
		if len(itemArtists) != 0 {
			item.ArtistInfo = ArtistInfo{Artist: []Artist{itemArtists[0]}}
		}
		if len(itemRender.FlexColumns[0].MusicResponsiveListItemFlexColumnRenderer.Text.Runs) > 0 {
			flexCol0 := itemRender.FlexColumns[0].MusicResponsiveListItemFlexColumnRenderer.Text.Runs[0]
			item.Title = flexCol0.Text
			item.MusicVideoType = flexCol0.NavigationEndpoint.WatchEndpoint.WatchEndpointMusicSupportedConfigs.WatchEndpointMusicConfig.MusicVideoType
			if flexCol0.NavigationEndpoint.WatchEndpoint.VideoId != "" {
				item.VideoId = &flexCol0.NavigationEndpoint.WatchEndpoint.VideoId
			}
			if flexCol0.NavigationEndpoint.WatchEndpoint.PlaylistId != "" {
				item.PlaylistId = &flexCol0.NavigationEndpoint.WatchEndpoint.PlaylistId
			}
			item.PlayerParams = flexCol0.NavigationEndpoint.WatchEndpoint.PlayerParams
			if item.PlayerParams == "" {
				item.PlayerParams = flexCol0.NavigationEndpoint.WatchEndpoint.Params
			}
		}

		if itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchEndpoint.LoggingContext.VssLoggingContext.SerializedContextData != "" {
			item.LoggingContext.VssLoggingContext = itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchEndpoint.LoggingContext.VssLoggingContext
		} else {
			item.LoggingContext.VssLoggingContext = itemRender.NavigationEndpoint.WatchEndpoint.LoggingContext.VssLoggingContext
		}

		if (item.PlaylistId == nil || *item.PlaylistId == "") && len(itemRender.Menu.MenuRenderer.Items) > 0 {
			item.PlaylistId = &itemRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchPlaylistEndpoint.PlaylistId
		}

		if (item.PlaylistId == nil || *item.PlaylistId == "") && len(itemRender.Menu.MenuRenderer.Items) > 0 {
			item.PlaylistId = &itemRender.Menu.MenuRenderer.Items[0].MenuNavigationItemRenderer.NavigationEndpoint.WatchEndpoint.PlaylistId
		}

		if item.PlaylistId == nil || *item.PlaylistId == "" {
			item.PlaylistId = &itemRender.NavigationEndpoint.WatchEndpoint.PlaylistId
		}

		if itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId != "" ||
			itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType != "" {
			endpoint := Endpoint{}
			endpoint.BrowseId = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseId
			endpoint.PageType = itemRender.NavigationEndpoint.BrowseEndpoint.BrowseEndpointContextSupportedConfigs.BrowseEndpointContextMusicConfig.PageType
			item.Endpoint = &endpoint
		}

		//Explicit:       "badges" in itemRender,
	}

	if len(itemRender.FixedColumns) != 0 {
		item.Length = itemRender.FixedColumns[0].MusicResponsiveListItemFixedColumnRenderer.Text.Runs[0].Text
	}

	item.PlaylistSetVideoId = itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchEndpoint.PlaylistSetVideoId

	/*


		if playlistSetVideoID {
			Item.PlaylistSetVideoID = itemRender.PlaylistSetVideoID ||
				itemRender.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchEndpoint.PlaylistSetVideoID
			Item.PlaylistID = playlistID
		}
	*/

	return item
}
