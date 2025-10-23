package api

import (
	"beatbump-server/backend/_youtube"
	"beatbump-server/backend/_youtube/api"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AlbumEndpointHandler(c echo.Context) error {

	query := c.Request().URL.Query()
	browseID := query.Get("browseId")
	//pt := query.Get("pt")
	//	ctoken := query.Get("ctoken")
	//referrer := query.Get("ref")
	//visitorData := query.Get("visitorData")

	var responseBytes []byte
	var err error
	responseBytes, err = api.Browse(browseID, api.PageType_MusicPageTypeAlbum, "", nil, nil, nil, api.WebMusic)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	var albumResponse _youtube.AlbumResponse
	err = json.Unmarshal(responseBytes, &albumResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error building API request: %s", err))
	}

	r := parseAlbum(albumResponse)

	return c.JSON(http.StatusOK, r)

}

func parseAlbum(playlistResponse _youtube.AlbumResponse) map[string]interface{} {

	//header := map[string]interface{}{}
	var items []IListItemRenderer = make([]IListItemRenderer, 0)
	var response = make(map[string]interface{}, 0)

	section := playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Contents
	/*musicShelf := Carousel{}

	  //section.MusicTastebuilderShelfRenderer.PrimaryText.Runs[0].Text
	  /*if len(section.MusicDescriptionShelfRenderer.Header.Runs) != 0 {
	  	description.Header = section.MusicDescriptionShelfRenderer.Header.Runs[0].Text
	  }
	  if len(section.MusicDescriptionShelfRenderer.Description.Runs) != 0 {
	  	description.Description = section.MusicDescriptionShelfRenderer.Description.Runs[0].Text
	  }
	  header := section.MusicCarouselShelfRenderer.Header

	  if len(header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs) != 0 {
	  	musicShelf.Header.Title = header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs[0].Text
	  }
	  if len(section.MusicCarouselShelfRenderer.Header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs) != 0 {
	  	musicShelf.Header.Subheading = stringPtr(header.MusicCarouselShelfBasicHeaderRenderer.Title.Runs[0].Text)
	  }*/

	//musicShelf.Header.BrowseId = stringPtr(section.MusicCarouselShelfRenderer.Header.MusicCarouselShelfBasicHeaderRenderer.MoreContentButton.ButtonRenderer.NavigationEndpoint.WatchPlaylistEndpoint.PlaylistID)

	//contents := make([]IListItemRenderer, 0, len(section[0].MusicCarouselShelfRenderer.Contents))
	for _, carousel := range section[0].MusicShelfRenderer.Contents {
		var item IListItemRenderer

		item = parseMusicResponsiveListItemRenderer(carousel.MusicResponsiveListItemRenderer)

		items = append(items, item)
	}

	header := playlistResponse.Contents.TwoColumnBrowseResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicResponsiveHeaderRenderer
	var releaseInfo = make(map[string]interface{}, 0)
	//	response.Tracks = tracks
	//header["playlistId"] = musicPlaylistShelfRenderer.PlaylistId
	subtitle := []map[string]interface{}{}
	subtitle = append(subtitle, map[string]interface{}{
		"year":   header.SecondSubtitle.Runs[2].Text,
		"tracks": header.SecondSubtitle.Runs[0].Text,
		"length": header.Subtitle.Runs[2].Text,
	})

	artists := []map[string]interface{}{}
	if len(header.StraplineTextOne.Runs) != 0 {
		for _, t := range header.StraplineTextOne.Runs {
			artists = append(artists, map[string]interface{}{
				"name":      t.Text,
				"channelId": t.NavigationEndpoint.BrowseEndpoint.BrowseID,
			})
		}
	}

	thumbnails := []Thumbnail{}
	if len(header.Thumbnail.MusicThumbnailRenderer.Thumbnail.Thumbnails) != 0 {
		for _, t := range header.Thumbnail.MusicThumbnailRenderer.Thumbnail.Thumbnails {
			thumbnails = append(thumbnails, Thumbnail{
				URL:    t.URL,
				Width:  t.Width,
				Height: t.Height,
			})
		}
	}

	playlistId := ""
	if len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Contents) != 0 &&
		len(playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Contents[0].MusicShelfRenderer.Contents) != 0 {
		playlistId = playlistResponse.Contents.TwoColumnBrowseResultsRenderer.SecondaryContents.SectionListRenderer.Contents[0].MusicShelfRenderer.Contents[0].MusicResponsiveListItemRenderer.Overlay.MusicItemThumbnailOverlayRenderer.Content.MusicPlayButtonRenderer.PlayNavigationEndpoint.WatchEndpoint.PlaylistId
	}

	releaseInfo["thumbnails"] = thumbnails
	releaseInfo["artist"] = artists
	releaseInfo["title"] = header.Title.Runs[0].Text
	releaseInfo["subtitles"] = subtitle
	releaseInfo["playlistId"] = playlistId
	releaseInfo["autoMixId"] = playlistId
	response["items"] = items
	response["releaseInfo"] = releaseInfo

	return map[string]interface{}{
		"items": response,
	}

	/*return {
	  items: promise,
	  	id: browseId,
	  		path,
	  };*/
}

/*import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BrowseRequest struct {
	Context struct {
		Client struct {
			ClientName           string `json:"clientName"`
			ClientVersion        string `json:"clientVersion"`
			DeviceMake           string `json:"deviceMake"`
			Platform             string `json:"platform"`
			DeviceModel          string `json:"deviceModel"`
			ExperimentIds        []interface{} `json:"experimentIds"`
			ExperimentsToken     string `json:"experimentsToken"`
			OsName               string `json:"osName"`
			OsVersion            string `json:"osVersion"`
			LocationInfo struct {
				LocationPermissionAuthorizationStatus string `json:"locationPermissionAuthorizationStatus"`
			} `json:"locationInfo"`
			MusicAppInfo struct {
				MusicActivityMasterSwitch      string `json:"musicActivityMasterSwitch"`
				MusicLocationMasterSwitch      string `json:"musicLocationMasterSwitch"`
				PwaInstallabilityStatus        string `json:"pwaInstallabilityStatus"`
			} `json:"musicAppInfo"`
			UtcOffsetMinutes int `json:"utcOffsetMinutes"`
		} `json:"client"`
		Capabilities map[string]interface{} `json:"capabilities"`
		Request struct {
			InternalExperimentFlags []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"internalExperimentFlags"`
			SessionIndex map[string]interface{} `json:"sessionIndex"`
		} `json:"request"`
		User struct {
			LockedSafetyMode bool `json:"lockedSafetyMode"`
		} `json:"user"`
		ActivePlayers map[string]interface{} `json:"activePlayers"`
	} `json:"context"`
	BrowseEndpointContextMusicConfig struct {
		BrowseEndpointContextMusicConfig struct {
			PageType string `json:"pageType"`
		} `json:"browseEndpointContextMusicConfig"`
	} `json:"browseEndpointContextMusicConfig"`
	BrowseID string `json:"browseId"`
}

func MusicBrowseHandler(c echo.Context) error {
	browseID := c.QueryParam("browseId")

	requestBody := BrowseRequest{
		Context: struct {
			Client struct {
				ClientName           string `json:"clientName"`
				ClientVersion        string `json:"clientVersion"`
				DeviceMake           string `json:"deviceMake"`
				Platform             string `json:"platform"`
				DeviceModel          string `json:"deviceModel"`
				ExperimentIds        []interface{} `json:"experimentIds"`
				ExperimentsToken     string `json:"experimentsToken"`
				OsName               string `json:"osName"`
				OsVersion            string `json:"osVersion"`
				LocationInfo struct {
					LocationPermissionAuthorizationStatus string `json:"locationPermissionAuthorizationStatus"`
				} `json:"locationInfo"`
				MusicAppInfo struct {
					MusicActivityMasterSwitch      string `json:"musicActivityMasterSwitch"`
					MusicLocationMasterSwitch      string `json:"musicLocationMasterSwitch"`
					PwaInstallabilityStatus        string `json:"pwaInstallabilityStatus"`
				} `json:"musicAppInfo"`
				UtcOffsetMinutes int `json:"utcOffsetMinutes"`
			} `json:"client"`
			Capabilities map[string]interface{} `json:"capabilities"`
			Request struct {
				InternalExperimentFlags []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"internalExperimentFlags"`
				SessionIndex map[string]interface{} `json:"sessionIndex"`
			} `json:"request"`
			User struct {
				LockedSafetyMode bool `json:"lockedSafetyMode"`
			} `json:"user"`
			ActivePlayers map[string]interface{} `json:"activePlayers"`
		}{
			Client: struct {
				ClientName           string `json:"clientName"`
				ClientVersion        string `json:"clientVersion"`
				DeviceMake           string `json:"deviceMake"`
				Platform             string `json:"platform"`
				DeviceModel          string `json:"deviceModel"`
				ExperimentIds        []interface{} `json:"experimentIds"`
				ExperimentsToken     string `json:"experimentsToken"`
				OsName               string `json:"osName"`
				OsVersion            string `json:"osVersion"`
				LocationInfo struct {
					LocationPermissionAuthorizationStatus string `json:"locationPermissionAuthorizationStatus"`
				} `json:"locationInfo"`
				MusicAppInfo struct {
					MusicActivityMasterSwitch      string `json:"musicActivityMasterSwitch"`
					MusicLocationMasterSwitch      string `json:"musicLocationMasterSwitch"`
					PwaInstallabilityStatus        string `json:"pwaInstallabilityStatus"`
				} `json:"musicAppInfo"`
				UtcOffsetMinutes int `json:"utcOffsetMinutes"`
			}{
				ClientName:           "WEB_REMIX",
				ClientVersion:        "0.1",
				DeviceMake:           "google",
				Platform:             "DESKTOP",
				DeviceModel:          "bot",
				ExperimentIds:        []interface{}{},
				ExperimentsToken:     "",
				OsName:               "Googlebot",
				OsVersion:            "2.1",
				LocationInfo: struct {
					LocationPermissionAuthorizationStatus string `json:"locationPermissionAuthorizationStatus"`
				}{
					LocationPermissionAuthorizationStatus: "LOCATION_PERMISSION_AUTHORIZATION_STATUS_UNSUPPORTED",
				},
				MusicAppInfo: struct {
					MusicActivityMasterSwitch      string `json:"musicActivityMasterSwitch"`
					MusicLocationMasterSwitch      string `json:"musicLocationMasterSwitch"`
					PwaInstallabilityStatus        string `json:"pwaInstallabilityStatus"`
				}{
					MusicActivityMasterSwitch:      "MUSIC_ACTIVITY_MASTER_SWITCH_INDETERMINATE",
					MusicLocationMasterSwitch:      "MUSIC_LOCATION_MASTER_SWITCH_INDETERMINATE",
					PwaInstallabilityStatus:        "PWA_INSTALLABILITY_STATUS_UNKNOWN",
				},
				UtcOffsetMinutes: -new Date().getTimezoneOffset(),
			},
			Capabilities: map[string]interface{}{},
			Request: struct {
				InternalExperimentFlags []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"internalExperimentFlags"`
				SessionIndex map[string]interface{} `json:"sessionIndex"`
			}{
				InternalExperimentFlags: []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					{
						Key:   "force_music_enable_outertube_tastebuilder_browse",
						Value: "true",
					},
					{
						Key:   "force_music_enable_outertube_playlist_detail_browse",
						Value: "true",
					},
					{
						Key:   "force_music_enable_outertube_search_suggestions",
						Value: "true",
					},
				},
				SessionIndex: map[string]interface{}{},
			},
			User: struct {
				LockedSafetyMode bool `json:"lockedSafetyMode"`
			}{
				LockedSafetyMode: false,
			},
			ActivePlayers: map[string]interface{}{},
		},
		BrowseEndpointContextMusicConfig: struct {
			BrowseEndpointContextMusicConfig struct {
				PageType string `json:"pageType"`
			} `json:"browseEndpointContextMusicConfig"`
		}{
			BrowseEndpointContextMusicConfig: struct {
				PageType string `json:"pageType"`
			}{
				PageType: "MUSIC_PAGE_TYPE_ARTIST",
			},
		},
		BrowseID: browseID,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error encoding request body")
	}

	url := "https://music.youtube.com/youtubei/v1/browse?key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30&prettyPrint=false"
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error making request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Request failed with status: %d - %s", resp.StatusCode, resp.Status))
	}

	return c.JSON(http.StatusOK, string(body))
}
*/
