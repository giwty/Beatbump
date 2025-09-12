package api

import (
	"beatbump-server/backend/_youtube/api/auth"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const URL_BASE = "https://music.youtube.com/youtubei/v1/"

const DEBUG = false

func Browse(browseId string, pageType PageType, params string,
	visitorData *string, itct *string, ctoken *string, client ClientInfo, authObj *auth.Auth) ([]byte, error) {

	urlAddress := URL_BASE + "browse" + "?prettyPrint=false"
	innertubeContext := prepareInnertubeContext(client, visitorData)
	//playlistUrl := "https://music.youtube.com/playlist?list=" + ""
	//innertubeContext.Client.OriginalUrl = &playlistUrl
	data := innertubeRequest{
		Context: innertubeContext,
	}
	if itct != nil && ctoken != nil {
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

	if params != "" {
		data.Params = params
	}
	resp, err := callAPI(urlAddress, data, client, authObj)
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

	resp, err := callAPI(url, data, client, nil)
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

	resp, err := callAPI(url, data, client, nil)
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

	resp, err := callAPI(url, data, client, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func Next(videoId string, playlistId string, client ClientInfo, params Params, authObj *auth.Auth) ([]byte, error) {

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

	resp, err := callAPI(url, data, client, authObj)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func Player(videoId string, playlistId string, client ClientInfo, params Params, authObj *auth.Auth) ([]byte, error) {

	playerUrl := URL_BASE + "player" + "?prettyPrint=false"

	innertubeContext := prepareInnertubeContext(client, nil)
	state, err := GetPlayerInfo(videoId, authObj)

	str := ""
	if err == nil {
		str = strconv.FormatUint(state.SignatureTimestamp, 10)
	}

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
				SignatureTimestamp: str,
				HTML5Preference:    "HTML5_PREF_WANTS",
				//Referer:            "https://www.youtube.com/watch?v=" + videoId,
			},
		},
	}

	resp, err := callAPI(playerUrl, data, client, authObj)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func DownloadWebpage(urlAddress string, clientInfo ClientInfo, authbj *auth.Auth) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, urlAddress, nil)
	if err != nil {
		return nil, err
	}
	return doRequest(clientInfo, authbj, req, nil)
}

func callAPI(urlAddress string, requestPayload innertubeRequest, clientInfo ClientInfo, authbj *auth.Auth) ([]byte, error) {

	req, err := http.NewRequest(http.MethodPost, urlAddress, nil)

	if err != nil {
		return nil, err
	}
	return doRequest(clientInfo, authbj, req, &requestPayload)
}

func doRequest(clientInfo ClientInfo, authbj *auth.Auth, req *http.Request, requestPayload *innertubeRequest) ([]byte, error) {

	urlAddress := req.URL.String()
	client := getHttpClient(authbj)

	if strings.Contains(urlAddress, "youtubei/v1/player") {
		if authbj != nil && authbj.AuthType == auth.AUTH_TYPE_COOKIES && authbj.CookieHeader != "" {
			u, _ := url.Parse("https://www.youtube.com/")
			for _, cookie := range client.Jar.Cookies(u) {

				if cookie.Name == "__Secure-3PAPISID" || cookie.Name == "SAPISID" {
					authorization := getAuthorization(cookie.Value)
					req.Header.Set("Authorization", authorization)

					/*if requestPayload != nil {
						info, _ := GetPlayerInfo("", authbj)
						//contextMap := info.InnerTubeContext.Context.(map[string]interface{})
						clinetAttr := info.InnerTubeContext["client"].(map[string]interface{})
						req.Header.Set("X-Goog-Visitor-Id", clinetAttr["visitorData"].(string))
					}*/

					req.Header.Set("X-Goog-AuthUser", "0")
				}
			}
		} else if authbj != nil && authbj.AuthType == auth.AUTH_TYPE_OAUTH && authbj.OauthToken != nil {
			req.Header.Set("Authorization", authbj.OauthToken.TokenType+" "+authbj.OauthToken.AccessToken)
			// Get the current Unix timestamp as an int64
			timestamp := time.Now().Unix()

			// Convert the timestamp to a string
			timestampString := strconv.FormatInt(timestamp, 10)
			req.Header.Set("X-Goog-Request-Time", timestampString)
			//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0 Cobalt/Version")
		} else {
			poToken := GetPoToken()
			if poToken != nil {
				requestPayload.ServiceIntegrityDimensions = &ServiceIntegrityDimensions{
					PoToken: poToken.Potoken,
				}
				requestPayload.Context.Client.VisitorData = poToken.VisitorData
			}
		}
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

	if DEBUG {
		res, _ := httputil.DumpRequest(req, true)
		fmt.Print(string(res))
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("got non 200 response %d", resp.StatusCode)
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	respBytes, err := io.ReadAll(reader)

	return respBytes, nil
}

func getHttpClient(authObj *auth.Auth) http.Client {

	myDialer := net.Dialer{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return myDialer.DialContext(ctx, "tcp4", addr)
	}
	client := http.Client{
		Transport: transport,
	}

	if authObj != nil && authObj.AuthType == auth.AUTH_TYPE_COOKIES && authObj.CookieHeader != "" {
		cookies, err := parseCookieString(authObj.CookieHeader)
		if err != nil {
			fmt.Println("Error parsing cookies:", err)
			return client
		}
		var cookiesList []*http.Cookie
		for key, value := range cookies {

			cookie := &http.Cookie{
				Name:    key,
				Value:   value,
				Path:    "/",
				Domain:  ".youtube.com",
				Expires: time.Now().Add(5 * time.Minute),
				Secure:  true,
			}
			cookiesList = append(cookiesList, cookie)
		}

		jar, err := cookiejar.New(nil)
		u, _ := url.Parse("https://www.youtube.com/")
		jar.SetCookies(u, cookiesList)
		client.Jar = jar
	}
	return client
}

func prepareInnertubeContext(clientInfo ClientInfo, visitorData *string) inntertubeContext {
	client := innertubeClient{
		HL:            "en",
		GL:            "US",
		ClientName:    clientInfo.ClientName,
		ClientVersion: clientInfo.ClientVersion,
		TimeZone:      "UTC",
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
		User: map[string]string{
			"lockedSafetyMode": "false",
		},
		Request: map[string]string{
			"useSsl": "true",
		},
	}
}

func strPtr(s string) *string {
	return &s
}

// ParseCookieString parses a cookie string into a map of key-value pairs.
func parseCookieString(cookieStr string) (map[string]string, error) {
	cookies := make(map[string]string)

	// Split the cookie string by ';' to get individual components
	pairs := strings.Split(cookieStr, ";")

	for _, pair := range pairs {
		// Trim any leading/trailing whitespace
		pair = strings.TrimSpace(pair)

		// Split the pair by '=' to separate name and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			cookies[parts[0]] = parts[1]
		} else {
			// If there's no '=', treat the pair as a flag
			cookies[parts[0]] = ""
		}
	}

	return cookies, nil
}

func getAuthorization(auth string) string {
	ts := time.Now().Unix()
	hash := sha1.Sum([]byte(fmt.Sprintf("%d %s %s", ts, auth, "https://music.youtube.com")))
	return fmt.Sprintf("SAPISIDHASH %d_%x", ts, hash[:])
}
