package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const URL_BASE = "https://music.youtube.com/youtubei/v1/"

const DEBUG = false

var companionBaseURL string
var companionAPIKey string

func init() {
	companionBaseURL = os.Getenv("COMPANION_URL")
	companionAPIKey = os.Getenv("COMPANION_SECRET_KEY")
}

func Browse(browseId string, pageType PageType, params string,
	visitorData *string, itct *string, ctoken *string, client ClientInfo) ([]byte, error) {

	urlAddress := URL_BASE + "browse" + "?prettyPrint=false"
	innertubeContext := prepareInnertubeContext(client, visitorData)

	data := innertubeRequest{
		Context: innertubeContext,
	}
	if (itct == nil || *itct == "") && ctoken != nil {
		data = innertubeRequest{
			//RequestAttributes: additionalRequestAttributes,
			Continuation: ctoken,
			Context:  innertubeContext,
			//ContentCheckOK: true,
			//RacyCheckOk:    true,
			BrowseEndpointContextMusicConfig: &BrowseEndpointContextMusicConfig{
				PageType: string(pageType),
			},
		}
	} else if itct != nil && ctoken != nil {
		urlAddress = handleContinuation(urlAddress, *itct, *ctoken)
	} else {

		data = innertubeRequest{
			//RequestAttributes: additionalRequestAttributes,
			BrowseID: browseId,
			Context:  innertubeContext,
			//ContentCheckOK: true,
			//RacyCheckOk:    true,
			BrowseEndpointContextMusicConfig: &BrowseEndpointContextMusicConfig{
				PageType: string(pageType),
			},
		}
	}

	if ctoken != nil {
		data.Continuation = ctoken
	}

	if params != "" {
		data.Params = params
	}
	resp, err := callAPI(urlAddress, data, client)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func handleContinuation(url string, itct string, ctoken string) string {
	url += "&itct=" + itct
	url += "&continuation=" + ctoken
	url += "&ctoken=" + ctoken
	url += "&type=next"
	return url
}

func GetSearchSuggestions(query string, client ClientInfo) ([]byte, error) {
	innertubeContext := prepareInnertubeContext(client, nil)

	url := URL_BASE + "music/get_search_suggestions" + "?prettyPrint=false"

	data := innertubeRequest{
		//RequestAttributes: additionalRequestAttributes,
		Context: innertubeContext,
		Input:   strPtr(query),
	}

	resp, err := callAPI(url, data, client)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Search(query string, filter string, itct *string, ctoken *string, client ClientInfo) ([]byte, error) {
	innertubeContext := prepareInnertubeContext(client, nil)

	url := URL_BASE + "search" + "?prettyPrint=false"

	if itct != nil && ctoken != nil {
		url = handleContinuation(url, *itct, *ctoken)
	}

	data := innertubeRequest{
		//RequestAttributes: additionalRequestAttributes,
		Context:  innertubeContext,
		BrowseID: "",
		Query:    query,
		/*user: {
			lockedSafetyMode: locals.preferences.Restricted,
		},*/
		//Continuation: continuationMap,
		//ContentCheckOK: true,
		//RacyCheckOk:    true,
		//Params: reqParams,
		/*PlaybackContext: &playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				// SignatureTimestamp: sts,
				HTML5Preference: "HTML5_PREF_WANTS",
			},
		},*/
	}

	if filter != "" {
		data.Params = filter
	}

	resp, err := callAPI(url, data, client)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetQueue(videoId string, playlistId string, client ClientInfo) ([]byte, error) {

	url := URL_BASE + "/music/get_queue" + "?prettyPrint=false"

	innertubeContext := prepareInnertubeContext(client, nil)

	//reqParams, err := createRequestParams(params)

	data := innertubeRequest{
		//RequestAttributes: additionalRequestAttributes,
		VideoID:        videoId,
		Context:        innertubeContext,
		ContentCheckOK: true,
		RacyCheckOk:    true,
		PlaylistId:     playlistId,
		/*EnablePersistentPlaylistPanel: true,
		  IsAudioOnly:                   true,
		  TunerSettingValue:             "AUTOMIX_SETTING_NORMAL",*/
		/*PlaybackContext: &playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				// SignatureTimestamp: sts,
				HTML5Preference: "HTML5_PREF_WANTS",
			},
		},*/
	}

	resp, err := callAPI(url, data, client)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func Next(videoId string, playlistId string, client ClientInfo, params Params) ([]byte, error) {

	url := URL_BASE + "next" + "?prettyPrint=false"

	innertubeContext := prepareInnertubeContext(client, strPtr(params["visitorData"]))

	//reqParams, err := createRequestParams(params)

	data := innertubeRequest{
		//RequestAttributes: additionalRequestAttributes,
		VideoID: videoId,
		Context: innertubeContext,
		// ContentCheckOK:                true,
		// RacyCheckOk:                   true,
		PlaylistId:                    playlistId,
		EnablePersistentPlaylistPanel: true,
		IsAudioOnly:                   true,
		TunerSettingValue:             "AUTOMIX_SETTING_NORMAL",
		//PlaylistSetVideoId:            params["playlistSetVideoId"],
		Params: "wAEB",

		/*PlaybackContext: &playbackContext{
		    ContentPlaybackContext: contentPlaybackContext{
		        // SignatureTimestamp: sts,
		        HTML5Preference: "HTML5_PREF_WANTS",
		    },
		},*/
	}

	resp, err := callAPI(url, data, client)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func Player(videoId string, playlistId string, client ClientInfo, params Params) ([]byte, error) {
		
	if companionBaseURL == "" {
		return nil, errors.New("Missing companion base URL")
	}
	//if companion url ends with /, remove it
	if companionBaseURL[len(companionBaseURL)-1] == '/' {
		companionBaseURL = companionBaseURL[:len(companionBaseURL)-1]
	}
	playerUrl := companionBaseURL + "/companion/youtubei/v1/player"
	innertubeContext := prepareInnertubeContext(client, nil)

	data := innertubeRequest{
		//RequestAttributes: additionalRequestAttributes,
		VideoID:        videoId,
		Context:        innertubeContext, //innertubeContext,
		ContentCheckOK: true,
		RacyCheckOk:    true,
		//Params:         reqParams,
		PlaylistId: playlistId,

		PlaybackContext: &playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				//SignatureTimestamp: str,
				HTML5Preference: "HTML5_PREF_WANTS",
				//Referer:            "https://www.youtube.com/watch?v=" + videoId,
			},
		},
	}

	resp, err := callAPI(playerUrl, data, client)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func DownloadWebpage(urlAddress string, clientInfo ClientInfo) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, urlAddress, nil)
	if err != nil {
		return nil, err
	}
	return doRequest(clientInfo, req, nil)
}

func callAPI(urlAddress string, requestPayload innertubeRequest, clientInfo ClientInfo) ([]byte, error) {

	req, err := http.NewRequest(http.MethodPost, urlAddress, nil)

	if err != nil {
		return nil, err
	}
	return doRequest(clientInfo, req, &requestPayload)
}

func doRequest(clientInfo ClientInfo, req *http.Request, requestPayload *innertubeRequest) ([]byte, error) {

	client := getHttpClient()
	urlAddress := req.URL.String()

	if strings.Contains(urlAddress, "companion") {
		req.Header.Set("Authorization", "Bearer "+companionAPIKey)
	}

	if clientInfo.userAgent != "" {
		req.Header.Set("User-Agent", clientInfo.userAgent)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.19 Safari/537.36")
	}
	if req.Method == http.MethodPost {
		req.Header.Set("X-Youtube-Client-Name", clientInfo.ClientId)
		req.Header.Set("X-Youtube-Client-Version", clientInfo.ClientVersion)
		req.Header.Set("Origin", "https://music.youtube.com")
		req.Header.Set("X-Origin", "https://music.youtube.com")
		req.Header.Set("Content-Type", "application/json")

		payload, err := json.Marshal(requestPayload)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewBuffer(payload))
	}
	req.Header.Set("Accept-Language", "en-us,en;q=0.5")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9")
	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("referer", "https://music.youtube.com")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	respBytes, err := io.ReadAll(reader)

	if resp.StatusCode != http.StatusOK {
		log.Printf("API call failed with status %d \n  %s", resp.StatusCode, string(respBytes))
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(dump))
		log.Println(string(respBytes))
		return nil, errors.New(resp.Status)
	}

	return respBytes, nil
}

func getHttpClient() http.Client {

	myDialer := net.Dialer{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return myDialer.DialContext(ctx, "tcp4", addr)
	}
	client := http.Client{
		Transport: transport,
	}

	return client
}

func prepareInnertubeContext(clientInfo ClientInfo, visitorData *string) inntertubeContext {
	client := innertubeClient{
		//	HL:            "en",
		//	GL:            "US",
		ClientName:    clientInfo.ClientName,
		ClientVersion: clientInfo.ClientVersion,
		//	TimeZone:      "UTC",
	}
	if clientInfo.DeviceModel != "" {
		client.DeviceModel = clientInfo.DeviceModel
	}
	if clientInfo.DeviceMake != "" {
		client.DeviceMake = clientInfo.DeviceMake
	}
	if clientInfo.OsVersion != "" {
		client.OsVersion = clientInfo.OsVersion
	}
	if clientInfo.OsName != "" {
		client.OsName = clientInfo.OsName
	}
	if clientInfo.DeviceMake != "" {
		client.DeviceMake = clientInfo.DeviceMake
	}
	if clientInfo.AndroidSdkVersion != 0 {
		client.AndroidSDKVersion = clientInfo.AndroidSdkVersion
	}
	if visitorData != nil {
		escape := url.QueryEscape(*visitorData)
		client.VisitorData = escape
	}
	return inntertubeContext{
		Client: client,
		User:   map[string]string{
			//	"lockedSafetyMode": "false",
		},
		//Request: map[string]string{
		//	"useSsl": "true",
		//},
	}
}

func strPtr(s string) *string {
	return &s
}


