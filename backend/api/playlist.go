package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PlaylistEndpointParams struct {
	BrowseID string `json:"browseId"`
}

type PlaylistEndpointContinuation struct {
	Ctoken       string `json:"ctoken"`
	Continuation string `json:"continuation"`
	Itct         string `json:"itct"`
	Type         string `json:"type"`
}

type NextContinuationData struct {
	Continuations []PlaylistEndpointContinuation `json:"continuations"`
}

type MusicDetailHeaderRenderer struct {
	Subtitle       []string `json:"subtitle"`
	SecondSubtitle []string `json:"secondSubtitle"`
	Description    string   `json:"description"`
	Thumbnails     []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
	Title      string `json:"title"`
	PlaylistID string `json:"playlistId"`
}

type PlaylistAPIResponse struct {
	//Original              _youtube.PlaylistResponse `json:"_original"`
	Continuations         interface{}            `json:"continuations"`
	Tracks                []IListItemRenderer    `json:"tracks"`
	VisitorData           string                 `json:"visitorData"`
	CarouselContinuations interface{}            `json:"carouselContinuations"`
	Header                map[string]interface{} `json:"header"`
}

func PlaylistEndpointHandler(c echo.Context) error {
	query := c.Request().URL.Query()
	browseID := query.Get("list")
	itct := query.Get("itct")
	ctoken := query.Get("ctoken")

	r, err := GetPlaylist(browseID, ctoken, itct)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, r)
}
func GetPlaylist(browseID string, ctoken string, itct string) (PlaylistAPIResponse, error) {
	var responseBytes []byte
	var err error
	if ctoken != "" && itct != "" {
		responseBytes, err = api.Browse(browseID, api.PageType_MusicPageTypePlaylist, "", nil, &itct, &ctoken, api.WebMusic)
	} else {
		responseBytes, err = api.Browse(browseID, api.PageType_MusicPageTypePlaylist, "", nil, nil, nil, api.WebMusic)
	}

	if err != nil {
		return PlaylistAPIResponse{}, fmt.Errorf("error building API request: %w", err)
	}

	playlistResponse := _youtube.PlaylistResponse{}
	err = json.Unmarshal(responseBytes, &playlistResponse)
	if err != nil {
		return PlaylistAPIResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	r := parsePlaylist(playlistResponse)
	return r, nil
}

func parsePlaylist(playlistResponse _youtube.PlaylistResponse) PlaylistAPIResponse {

	response := PlaylistAPIResponse{}
	header := map[string]interface{}{}
	var musicPlaylistShelfRenderer *_youtube.MusicPlaylistShelfRenderer = nil
	if len(playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs) != 0 {
		if len(playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) != 0 {
			musicPlaylistShelfRenderer = playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicPlaylistShelfRenderer
		}
		_carouselContinuation := playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Continuations
		if len(_carouselContinuation) != 0 {
			//continuationData := _carouselContinuation[0].NextContinuationData
			response.CarouselContinuations = _carouselContinuation[0].NextContinuationData
			response.Continuations = _carouselContinuation[0].NextContinuationData
		}
	} else if len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Contents) != 0 {
		musicPlaylistShelfRenderer = playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Contents[0].MusicPlaylistShelfRenderer
		_carouselContinuation := playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Continuations
		if len(_carouselContinuation) != 0 {
			//continuationData := _carouselContinuation[0].NextContinuationData
			response.CarouselContinuations = _carouselContinuation[0].NextContinuationData
			response.Continuations = _carouselContinuation[0].NextContinuationData
		}
	}

	if musicPlaylistShelfRenderer != nil {
		continuations := musicPlaylistShelfRenderer.Continuations
		if len(continuations) != 0 {
			//continuationData := continuations[0].NextContinuationData
			response.Continuations = continuations[0].NextContinuationData
		}
	}

	if musicPlaylistShelfRenderer != nil {
		tracks := []IListItemRenderer{}
		for _, musicResponsiveListItemRenderer := range musicPlaylistShelfRenderer.Contents {
			item := parseMusicResponsiveListItemRenderer(musicResponsiveListItemRenderer.MusicResponsiveListItemRenderer)
			tracks = append(tracks, item)
		}
		response.Tracks = tracks
		header["playlistId"] = musicPlaylistShelfRenderer.PlaylistId
	} else {
		if len(playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs) != 0 && len(playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) != 0 {
			gridRenderer := playlistResponse.Contents.SingleColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].GridRenderer
			if gridRenderer != nil {
				tracks := []IListItemRenderer{}
				for _, musicResponsiveListItemRenderer := range gridRenderer.Items {
					item := parseMusicTwoRowItemRenderer(*musicResponsiveListItemRenderer.MusicTwoRowItemRenderer)
					tracks = append(tracks, item)
				}
				response.Tracks = tracks
				header["playlistId"] = gridRenderer.TrackingParams
			}
		}
	}

	values := []string{}

	if len(playlistResponse.Header.MusicDetailHeaderRenderer.Subtitle.Runs) != 0 {
		for _, el := range playlistResponse.Header.MusicDetailHeaderRenderer.Subtitle.Runs {
			values = append(values, el.Text)
		}
	} else if len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs) != 0 {
		for _, el := range playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicResponsiveHeaderRenderer.Subtitle.Runs {
			values = append(values, el.Text)
		}
	}
	header["subtitle"] = values

	values = []string{}
	if len(playlistResponse.Header.MusicDetailHeaderRenderer.SecondSubtitle.Runs) != 0 {
		for _, el := range playlistResponse.Header.MusicDetailHeaderRenderer.SecondSubtitle.Runs {
			values = append(values, el.Text)
		}
	} else if playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs != nil && len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) != 0 {
		for _, el := range playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicResponsiveHeaderRenderer.SecondSubtitle.Runs {
			values = append(values, el.Text)
		}
	}
	header["secondSubtitle"] = values

	values = []string{}
	if len(playlistResponse.Header.MusicDetailHeaderRenderer.Title.Runs) != 0 {
		for _, el := range playlistResponse.Header.MusicDetailHeaderRenderer.Title.Runs {
			values = append(values, el.Text)
		}
	} else if playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs != nil && len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) != 0 {
		for _, el := range playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicResponsiveHeaderRenderer.Title.Runs {
			values = append(values, el.Text)
		}
	} else if len(playlistResponse.Header.MusicHeaderRenderer.Title.Runs) != 0 {
		for _, el := range playlistResponse.Header.MusicHeaderRenderer.Title.Runs {
			values = append(values, el.Text)
		}
	}
	header["title"] = values

	tvalues := []Thumbnail{}
	if len(playlistResponse.Header.MusicDetailHeaderRenderer.Thumbnail.CroppedSquareThumbnailRenderer.Thumbnail.Thumbnails) != 0 {
		for _, t := range playlistResponse.Header.MusicDetailHeaderRenderer.Thumbnail.CroppedSquareThumbnailRenderer.Thumbnail.Thumbnails {
			tvalues = append(tvalues, Thumbnail{
				URL:    t.Url,
				Width:  t.Width,
				Height: t.Height,
			})
		}
	} else if playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs != nil && len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents) != 0 {

		for _, t := range playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicResponsiveHeaderRenderer.Thumbnail.MusicThumbnailRenderer.Thumbnail.Thumbnails {
			tvalues = append(tvalues, Thumbnail{
				URL:    t.URL,
				Width:  t.Width,
				Height: t.Height,
			})
		}
	}

	if len(tvalues) != 0 {
		header["thumbnails"] = tvalues
	}

	response.Header = header
	visitorData := playlistResponse.ResponseContext.VisitorData
	response.VisitorData = visitorData

	return response
}
