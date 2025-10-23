package api

type innertubeRequest struct {
	ActivePlayers                    []map[string]string               `json:"activePlayers,omitempty"`
	VideoID                          string                            `json:"videoId,omitempty"`
	PlaylistId                       string                            `json:"playlistId,omitempty"`
	PLayerParams                     string                            `json:"playerParams,omitempty"`
	BrowseID                         string                            `json:"browseId,omitempty"`
	Query                            string                            `json:"query,omitempty"`
	Input                            *string                           `json:"input,omitempty"`
	Continuation                     map[string]string                 `json:"continuation,omitempty"`
	Context                          inntertubeContext                 `json:"context"`
	PlaybackContext                  *playbackContext                  `json:"playbackContext,omitempty"`
	ContentCheckOK                   bool                              `json:"contentCheckOk,omitempty"`
	RacyCheckOk                      bool                              `json:"racyCheckOk,omitempty"`
	Params                           string                            `json:"params,omitempty"`
	EnablePersistentPlaylistPanel    bool                              `json:"enablePersistentPlaylistPanel"`
	IsAudioOnly                      bool                              `json:"isAudioOnly"`
	TunerSettingValue                string                            `json:"tunerSettingValue"`
	PlaylistSetVideoId               string                            `json:"playlistSetVideoId,omitempty"`
	BrowseEndpointContextMusicConfig *BrowseEndpointContextMusicConfig `json:"browseEndpointContextMusicConfig,omitempty"`
	User                             *User                             `json:"user,omitempty"`
	ClickTracking                    *ClickTracking                    `json:"clickTracking,omitempty"`
	ServiceIntegrityDimensions       *ServiceIntegrityDimensions       `json:"serviceIntegrityDimensions,omitempty"`
}
type BrowseEndpointContextMusicConfig struct {
	PageType string `json:"pageType"`
}
type ServiceIntegrityDimensions struct {
	PoToken string `json:"poToken"`
}

type ClickTracking struct {
	ClickTrackingParams string `json:"clickTrackingParams"`
}

type User struct {
	LockedSafetyMode bool `json:"lockedSafetyMode,omitempty"`
}

type SearchRequest struct {
	Query  string `json:"query,omitempty"`
	Filter string `json:"filter,omitempty"`
	Itct   string `json:"itct,omitempty"`
	Ctoken string `json:"ctoken,omitempty"`
}

type playbackContext struct {
	ContentPlaybackContext contentPlaybackContext `json:"contentPlaybackContext"`
}

type contentPlaybackContext struct {
	SignatureTimestamp string `json:"signatureTimestamp"`
	HTML5Preference    string `json:"html5Preference"`
	Referer            string `json:"referer,omitempty"`
}

type inntertubeContext struct {
	Client  innertubeClient   `json:"client"`
	Request map[string]string `json:"request",omitempty`
	User    map[string]string `json:"user"`
	//CaptionParams string          `json:"caption_params"`
}

type innertubeClient struct {
	HL                 string  `json:"hl,omitempty"`
	GL                 string  `json:"gl,omitempty"`
	ClientName         string  `json:"clientName"`
	ClientVersion      string  `json:"clientVersion"`
	AndroidSDKVersion  int     `json:"androidSDKVersion,omitempty"`
	UserAgent          string  `json:"userAgent,omitempty"`
	TimeZone           string  `json:"timeZone,omitempty"`
	UTCOffset          int     `json:"utcOffsetMinutes,omitempty"`
	VisitorData        string  `json:"visitorData,omitempty"`
	OriginalUrl        *string `json:"originalUrl,omitempty"`
	DeviceExperimentId string  `json:"deviceExperimentId,omitempty"`
	DeviceMake         string  `json:"deviceMake",omitempty`
	BrowserName        *string `json:"BrowserName,omitempty"`
	BrowserVersion     *string `json:"browserVersion,omitempty"`
	Platform           string  `json:"platform,omitempty"`
	DeviceModel        string  `json:"deviceModel",omitempty`
	OsName             string  `json:"osName,omitempty"`
	OsVersion          string  `json:"osVersion,omitempty"`
	ScreenPixelDensity string  `json:"screenPixelDensity,omitempty"`
	UserInterfaceTheme string  `json:"userInterfaceTheme,omitempty"`
	AcceptHeader       string  `json:"acceptHeader,omitempty"`
}

type Params map[string]string

type ClientInfo struct {
	ClientId          string
	ClientName        string
	ClientKey         string
	ClientVersion     string
	userAgent         string
	AndroidSdkVersion int
	DeviceModel       string
	DeviceMake        string
	OsName            string
	OsVersion         string
	Platform          string
}

var (
	MWebClient = ClientInfo{
		ClientName:    "MWEB",
		ClientVersion: "2.20241202.07.00",
		//ClientKey:     "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
		userAgent: "Mozilla/5.0 (iPad; CPU OS 16_7_10 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1,gzip(gfe)",
	}

	// WebClient, better to use Android client but go ahead.
	WebClient = ClientInfo{
		ClientName:    "WEB",
		ClientVersion: "2.20250312.04.00",
		ClientId:      "1",
		//ClientKey:     "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.19 Safari/537.36",
	}

	// AndroidClient, download go brrrrrr.
	AndroidClient = ClientInfo{
		ClientName:    "ANDROID",
		ClientVersion: "19.09.37",
		ClientKey:     "AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w",
		userAgent:     "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip",

		AndroidSdkVersion: 30,
	}

	IOS = ClientInfo{
		ClientName:    "IOS",
		ClientVersion: "20.10.4",
		//ClientKey:     "AIzaSyB-63vPrdThhKuerbB2N_l7Kwwcxj6yUAc", // seems like same ClientKey works for both clients
		DeviceModel: "iPhone16,2",
		OsVersion:   "18.3.2.22D82",
		userAgent:   "com.google.ios.youtube/20.10.4 (iPhone16,2; U; CPU iOS 18_3_2 like Mac OS X;)",
	}

	IOS_MUSIC = ClientInfo{
		ClientName:    "IOS_MUSIC",
		ClientVersion: "6.33",
		DeviceMake:    "Apple",
		DeviceModel:   "iPhone16,2",
		//OsName:        "iPhone",
		//OsVersion:     "18.3.2.22D82",
		//ClientKey:     "AIzaSyB-63vPrdThhKuerbB2N_l7Kwwcxj6yUAc", // seems like same ClientKey works for both clients
		userAgent: "com.google.ios.youtubemusic/6.33 (iPhone16,2; U; CPU iOS 18_1_0 like Mac OS X;)",
	}

	WEB_CREATOR = ClientInfo{
		ClientName:    "WEB_CREATOR",
		ClientVersion: "1.20250312.03.01",
		ClientId:      "62",
		//ClientKey:     "AIzaSyB-63vPrdThhKuerbB2N_l7Kwwcxj6yUAc", // seems like same ClientKey works for both clients
		//DeviceModel: "iPhone16,2",
		//userAgent:   "com.google.ios.youtubemusic/7.27.0 (iPhone16,2; U; CPU iOS 18_1_0 like Mac OS X;)",
	}

	ANDROID_MUSIC = ClientInfo{
		ClientName:    "ANDROID_MUSIC",
		ClientVersion: " 7.27.52",
		//ClientId:      "62",
				userAgent:     "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip",
		//ClientKey:     "AIzaSyB-63vPrdThhKuerbB2N_l7Kwwcxj6yUAc", // seems like same ClientKey works for both clients
		//DeviceModel: "iPhone16,2",
		//userAgent:   "com.google.ios.youtubemusic/7.27.0 (iPhone16,2; U; CPU iOS 18_1_0 like Mac OS X;)",
	}

	ANDROID_VR = ClientInfo{
		ClientName:        "ANDROID_VR",
		ClientVersion:     "1.62.27",
		ClientId:          "28",
		DeviceModel:       "Quest 3",
		DeviceMake:        "Oculus",
		AndroidSdkVersion: 32,
		OsName:            "Android",
		OsVersion:         "12L",
		//	ClientKey: "AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30", // seems like same ClientKey works for both clients
		userAgent: "com.google.android.apps.youtube.vr.oculus/1.62.27 (Linux; U; Android 12L; eureka-user Build/SQ3A.220605.009.A1) gzip",
	}

	WebMusic = ClientInfo{
		ClientName:    "WEB_REMIX",
		ClientId:      "67",
		ClientVersion: "1.20250310.01.00",
		OsName:        "Windows",
		OsVersion:     "10.0",
		Platform:      "DESKTOP",
		//	ClientKey: "AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30", // seems like same ClientKey works for both clients
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.19 Safari/537.36",
	}

	TV = ClientInfo{
		ClientName:    "TVHTML5",
		ClientVersion: "7.20240925.00.00",
		ClientId:      "7",
		userAgent: "Mozilla/5.0 (ChromiumStylePlatform) Cobalt/25.lts.30.1034943-gold (unlike Gecko), Unknown_TV_Unknown_0/Unknown (Unknown, Unknown)",
		//OsName:        "Macintosh",
		//OsVersion:     "10.15",
		//Platform:      "DESKTOP",
		//	ClientKey: "AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30", // seems like same ClientKey works for both clients
		//userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0,gzip(gfe)",
	}

	TV_SIMPLY = ClientInfo{
		ClientName:    "TVHTML5_SIMPLY",
		ClientVersion: "1.0",
		ClientId:      "75",
	}


	TVHTML5_SIMPLY_EMBEDDED_PLAYER = ClientInfo{
		ClientName:    "TVHTML5_SIMPLY_EMBEDDED_PLAYER",
		ClientVersion: "2.0",
		ClientId:      "85",
		//OsName:        "Macintosh",
		//OsVersion:     "10.15",
		//Platform:      "DESKTOP",
		//	ClientKey: "AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30", // seems like same ClientKey works for both clients
		//userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0,gzip(gfe)",
	}

	TVHTML5_AUDIO = ClientInfo{
		ClientName:    "TVHTML5_AUDIO",
		ClientVersion: "2.0",
		//OsName:        "Macintosh",
		//OsVersion:     "10.15",
		//Platform:      "DESKTOP",
		//	ClientKey: "AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30", // seems like same ClientKey works for both clients
		//userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0,gzip(gfe)",
	}

	// EmbeddedClient, not really tested.
	EmbeddedClient = ClientInfo{
		ClientName:    "WEB_EMBEDDED_PLAYER",
		ClientVersion: "1.20220731.00.00",
		ClientKey:     "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8", // seems like same ClientKey works for both clients
		userAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
)

type PageType string

var PageType_MusicPageTypeAlbum PageType = "MUSIC_PAGE_TYPE_ALBUM"
var PageType_MusicPageTypeArtist PageType = "MUSIC_PAGE_TYPE_ARTIST"
var PageType_MusicPageTypePlaylist PageType = "MUSIC_PAGE_TYPE_PLAYLIST"
var PageType_MusicPageTypeTrackLyrics PageType = "MUSIC_PAGE_TYPE_TRACK_LYRICS"
var PageType_MusicPageTypeTrackRelated PageType = "MUSIC_PAGE_TYPE_TRACK_RELATED"
var PageType_MusicPageTypeUserChannel PageType = "MUSIC_PAGE_TYPE_USER_CHANNEL"
