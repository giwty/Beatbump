package _youtube

import (
	"encoding/json"
	"strings"
)

type MusicShelfContinuationContent struct {
	MusicResponsiveListItemRenderer MusicResponsiveListItemRenderer `json:"musicResponsiveListItemRenderer"`
}

type SearchResponse struct {
	ResponseContext struct {
		VisitorData           string `json:"visitorData"`
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
		MaxAgeSeconds int `json:"maxAgeSeconds"`
	} `json:"responseContext"`
	Content struct {
		TabbedSearchResultsRenderer TabbedSearchResultsRenderer `json:"tabbedSearchResultsRenderer"`
	} `json:"contents,omitempty"`
	ContinuationContents struct {
		MusicShelfContinuation struct {
			Content        []MusicShelfContinuationContent `json:"contents"`
			TrackingParams string                          `json:"trackingParams"`
			Continuations  []struct {
				NextContinuationData struct {
					Continuation        string `json:"continuation"`
					ClickTrackingParams string `json:"clickTrackingParams"`
				} `json:"nextContinuationData"`
			} `json:"continuations"`
			ShelfDivider struct {
				MusicShelfDividerRenderer struct {
					Hidden bool `json:"hidden"`
				} `json:"musicShelfDividerRenderer"`
			} `json:"shelfDivider"`
			AutoReloadWhenEmpty bool `json:"autoReloadWhenEmpty"`
		} `json:"musicShelfContinuation"`
	} `json:"continuationContents,omitempty"`
	Continuation struct {
		Continuation        string `json:"continuation"`
		ClickTrackingParams string `json:"clickTrackingParams"`
	} `json:"continuation,omitempty"`
	TrackingParams string `json:"trackingParams"`
}

type Runs struct {
	Text               string `json:"text"`
	NavigationEndpoint struct {
		ClickTrackingParams string `json:"clickTrackingParams"`
		BrowseEndpoint      struct {
			BrowseId                              string `json:"browseId"`
			BrowseEndpointContextSupportedConfigs struct {
				BrowseEndpointContextMusicConfig struct {
					PageType string `json:"pageType"`
				} `json:"browseEndpointContextMusicConfig"`
			} `json:"browseEndpointContextSupportedConfigs"`
		} `json:"browseEndpoint,omitempty"`
		WatchEndpoint struct {
			VideoId            string `json:"videoId"`
			PlaylistId         string `json:"playlistId"`
			Index              int    `json:"index"`
			Params             string `json:"params"`
			PlayerParams       string `json:"playerParams"`
			PlaylistSetVideoId string `json:"playlistSetVideoId"`
			LoggingContext     struct {
				VssLoggingContext struct {
					SerializedContextData string `json:"serializedContextData"`
				} `json:"vssLoggingContext"`
			} `json:"loggingContext"`
			WatchEndpointMusicSupportedConfigs struct {
				WatchEndpointMusicConfig struct {
					HasPersistentPlaylistPanel bool   `json:"hasPersistentPlaylistPanel"`
					MusicVideoType             string `json:"musicVideoType"`
				} `json:"watchEndpointMusicConfig"`
			} `json:"watchEndpointMusicSupportedConfigs"`
		} `json:"watchEndpoint"`
	} `json:"navigationEndpoint,omitempty"`
}
type Test struct {
	MusicResponsiveListItemRenderer struct {
		TrackingParams string `json:"trackingParams"`
		Thumbnail      struct {
			MusicThumbnailRenderer struct {
				Thumbnail struct {
					Thumbnails []struct {
						Url    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"thumbnails"`
				} `json:"thumbnail"`
				ThumbnailCrop  string `json:"thumbnailCrop"`
				ThumbnailScale string `json:"thumbnailScale"`
				TrackingParams string `json:"trackingParams"`
			} `json:"musicThumbnailRenderer"`
		} `json:"thumbnail"`
		Overlay struct {
			MusicItemThumbnailOverlayRenderer struct {
				Background struct {
					VerticalGradient struct {
						GradientLayerColors []string `json:"gradientLayerColors"`
					} `json:"verticalGradient"`
				} `json:"background"`
				Content struct {
					MusicPlayButtonRenderer struct {
						PlayNavigationEndpoint struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							WatchEndpoint       struct {
								VideoId            string `json:"videoId"`
								PlaylistId         string `json:"playlistId"`
								PlaylistSetVideoId string `json:"playlistSetVideoId"`
								LoggingContext     struct {
									VssLoggingContext struct {
										SerializedContextData string `json:"serializedContextData"`
									} `json:"vssLoggingContext"`
								} `json:"loggingContext"`
								WatchEndpointMusicSupportedConfigs struct {
									WatchEndpointMusicConfig struct {
										MusicVideoType string `json:"musicVideoType"`
									} `json:"watchEndpointMusicConfig"`
								} `json:"watchEndpointMusicSupportedConfigs"`
							} `json:"watchEndpoint"`
						} `json:"playNavigationEndpoint"`
						TrackingParams string `json:"trackingParams"`
						PlayIcon       struct {
							IconType string `json:"iconType"`
						} `json:"playIcon"`
						PauseIcon struct {
							IconType string `json:"iconType"`
						} `json:"pauseIcon"`
						IconColor             int64 `json:"iconColor"`
						BackgroundColor       int   `json:"backgroundColor"`
						ActiveBackgroundColor int   `json:"activeBackgroundColor"`
						LoadingIndicatorColor int64 `json:"loadingIndicatorColor"`
						PlayingIcon           struct {
							IconType string `json:"iconType"`
						} `json:"playingIcon"`
						IconLoadingColor      int    `json:"iconLoadingColor"`
						ActiveScaleFactor     int    `json:"activeScaleFactor"`
						ButtonSize            string `json:"buttonSize"`
						RippleTarget          string `json:"rippleTarget"`
						AccessibilityPlayData struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"accessibilityPlayData"`
						AccessibilityPauseData struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"accessibilityPauseData"`
					} `json:"musicPlayButtonRenderer"`
				} `json:"content"`
				ContentPosition string `json:"contentPosition"`
				DisplayStyle    string `json:"displayStyle"`
			} `json:"musicItemThumbnailOverlayRenderer"`
		} `json:"overlay"`
		FlexColumns []struct {
			MusicResponsiveListItemFlexColumnRenderer struct {
				Text struct {
					Runs []struct {
						Text               string `json:"text"`
						NavigationEndpoint struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							WatchEndpoint       struct {
								VideoId        string `json:"videoId"`
								PlaylistId     string `json:"playlistId"`
								PlayerParams   string `json:"playerParams"`
								LoggingContext struct {
									VssLoggingContext struct {
										SerializedContextData string `json:"serializedContextData"`
									} `json:"vssLoggingContext"`
								} `json:"loggingContext"`
								WatchEndpointMusicSupportedConfigs struct {
									WatchEndpointMusicConfig struct {
										MusicVideoType string `json:"musicVideoType"`
									} `json:"watchEndpointMusicConfig"`
								} `json:"watchEndpointMusicSupportedConfigs"`
							} `json:"watchEndpoint,omitempty"`
							BrowseEndpoint struct {
								BrowseId                              string `json:"browseId"`
								BrowseEndpointContextSupportedConfigs struct {
									BrowseEndpointContextMusicConfig struct {
										PageType string `json:"pageType"`
									} `json:"browseEndpointContextMusicConfig"`
								} `json:"browseEndpointContextSupportedConfigs"`
							} `json:"browseEndpoint,omitempty"`
						} `json:"navigationEndpoint"`
					} `json:"runs"`
				} `json:"text"`
				DisplayPriority string `json:"displayPriority"`
			} `json:"musicResponsiveListItemFlexColumnRenderer"`
		} `json:"flexColumns"`
		FixedColumns []struct {
			MusicResponsiveListItemFixedColumnRenderer struct {
				Text struct {
					Runs []struct {
						Text string `json:"text"`
					} `json:"runs"`
				} `json:"text"`
				DisplayPriority string `json:"displayPriority"`
				Size            string `json:"size"`
			} `json:"musicResponsiveListItemFixedColumnRenderer"`
		} `json:"fixedColumns"`
		Menu struct {
			MenuRenderer struct {
				Items []struct {
					MenuNavigationItemRenderer struct {
						Text struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"text"`
						Icon struct {
							IconType string `json:"iconType"`
						} `json:"icon"`
						NavigationEndpoint struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							WatchEndpoint       struct {
								VideoId        string `json:"videoId"`
								PlaylistId     string `json:"playlistId"`
								Params         string `json:"params"`
								LoggingContext struct {
									VssLoggingContext struct {
										SerializedContextData string `json:"serializedContextData"`
									} `json:"vssLoggingContext"`
								} `json:"loggingContext"`
								WatchEndpointMusicSupportedConfigs struct {
									WatchEndpointMusicConfig struct {
										MusicVideoType string `json:"musicVideoType"`
									} `json:"watchEndpointMusicConfig"`
								} `json:"watchEndpointMusicSupportedConfigs"`
							} `json:"watchEndpoint,omitempty"`
							ModalEndpoint struct {
								Modal struct {
									ModalWithTitleAndButtonRenderer struct {
										Title struct {
											Runs []struct {
												Text string `json:"text"`
											} `json:"runs"`
										} `json:"title"`
										Content struct {
											Runs []struct {
												Text string `json:"text"`
											} `json:"runs"`
										} `json:"content"`
										Button struct {
											ButtonRenderer struct {
												Style      string `json:"style"`
												IsDisabled bool   `json:"isDisabled"`
												Text       struct {
													Runs []struct {
														Text string `json:"text"`
													} `json:"runs"`
												} `json:"text"`
												NavigationEndpoint struct {
													ClickTrackingParams string `json:"clickTrackingParams"`
													SignInEndpoint      struct {
														Hack bool `json:"hack"`
													} `json:"signInEndpoint"`
												} `json:"navigationEndpoint"`
												TrackingParams string `json:"trackingParams"`
											} `json:"buttonRenderer"`
										} `json:"button"`
									} `json:"modalWithTitleAndButtonRenderer"`
								} `json:"modal"`
							} `json:"modalEndpoint,omitempty"`
							BrowseEndpoint struct {
								BrowseId                              string `json:"browseId"`
								BrowseEndpointContextSupportedConfigs struct {
									BrowseEndpointContextMusicConfig struct {
										PageType string `json:"pageType"`
									} `json:"browseEndpointContextMusicConfig"`
								} `json:"browseEndpointContextSupportedConfigs"`
							} `json:"browseEndpoint,omitempty"`
							ShareEntityEndpoint struct {
								SerializedShareEntity string `json:"serializedShareEntity"`
								SharePanelType        string `json:"sharePanelType"`
							} `json:"shareEntityEndpoint,omitempty"`
						} `json:"navigationEndpoint"`
						TrackingParams string `json:"trackingParams"`
					} `json:"menuNavigationItemRenderer,omitempty"`
					MenuServiceItemRenderer struct {
						Text struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"text"`
						Icon struct {
							IconType string `json:"iconType"`
						} `json:"icon"`
						ServiceEndpoint struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							QueueAddEndpoint    struct {
								QueueTarget struct {
									VideoId      string `json:"videoId"`
									OnEmptyQueue struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										WatchEndpoint       struct {
											VideoId string `json:"videoId"`
										} `json:"watchEndpoint"`
									} `json:"onEmptyQueue"`
								} `json:"queueTarget"`
								QueueInsertPosition string `json:"queueInsertPosition"`
								Commands            []struct {
									ClickTrackingParams string `json:"clickTrackingParams"`
									AddToToastAction    struct {
										Item struct {
											NotificationTextRenderer struct {
												SuccessResponseText struct {
													Runs []struct {
														Text string `json:"text"`
													} `json:"runs"`
												} `json:"successResponseText"`
												TrackingParams string `json:"trackingParams"`
											} `json:"notificationTextRenderer"`
										} `json:"item"`
									} `json:"addToToastAction"`
								} `json:"commands"`
							} `json:"queueAddEndpoint"`
						} `json:"serviceEndpoint"`
						TrackingParams string `json:"trackingParams"`
					} `json:"menuServiceItemRenderer,omitempty"`
				} `json:"items"`
				TrackingParams  string `json:"trackingParams"`
				TopLevelButtons []struct {
					LikeButtonRenderer struct {
						Target struct {
							VideoId string `json:"videoId"`
						} `json:"target"`
						LikeStatus                string `json:"likeStatus"`
						TrackingParams            string `json:"trackingParams"`
						LikesAllowed              bool   `json:"likesAllowed"`
						DislikeNavigationEndpoint struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							ModalEndpoint       struct {
								Modal struct {
									ModalWithTitleAndButtonRenderer struct {
										Title struct {
											Runs []struct {
												Text string `json:"text"`
											} `json:"runs"`
										} `json:"title"`
										Content struct {
											Runs []struct {
												Text string `json:"text"`
											} `json:"runs"`
										} `json:"content"`
										Button struct {
											ButtonRenderer struct {
												Style      string `json:"style"`
												IsDisabled bool   `json:"isDisabled"`
												Text       struct {
													Runs []struct {
														Text string `json:"text"`
													} `json:"runs"`
												} `json:"text"`
												NavigationEndpoint struct {
													ClickTrackingParams string `json:"clickTrackingParams"`
													SignInEndpoint      struct {
														Hack bool `json:"hack"`
													} `json:"signInEndpoint"`
												} `json:"navigationEndpoint"`
												TrackingParams string `json:"trackingParams"`
											} `json:"buttonRenderer"`
										} `json:"button"`
									} `json:"modalWithTitleAndButtonRenderer"`
								} `json:"modal"`
							} `json:"modalEndpoint"`
						} `json:"dislikeNavigationEndpoint"`
						LikeCommand struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							ModalEndpoint       struct {
								Modal struct {
									ModalWithTitleAndButtonRenderer struct {
										Title struct {
											Runs []struct {
												Text string `json:"text"`
											} `json:"runs"`
										} `json:"title"`
										Content struct {
											Runs []struct {
												Text string `json:"text"`
											} `json:"runs"`
										} `json:"content"`
										Button struct {
											ButtonRenderer struct {
												Style      string `json:"style"`
												IsDisabled bool   `json:"isDisabled"`
												Text       struct {
													Runs []struct {
														Text string `json:"text"`
													} `json:"runs"`
												} `json:"text"`
												NavigationEndpoint struct {
													ClickTrackingParams string `json:"clickTrackingParams"`
													SignInEndpoint      struct {
														Hack bool `json:"hack"`
													} `json:"signInEndpoint"`
												} `json:"navigationEndpoint"`
												TrackingParams string `json:"trackingParams"`
											} `json:"buttonRenderer"`
										} `json:"button"`
									} `json:"modalWithTitleAndButtonRenderer"`
								} `json:"modal"`
							} `json:"modalEndpoint"`
						} `json:"likeCommand"`
					} `json:"likeButtonRenderer"`
				} `json:"topLevelButtons"`
				Accessibility struct {
					AccessibilityData struct {
						Label string `json:"label"`
					} `json:"accessibilityData"`
				} `json:"accessibility"`
			} `json:"menuRenderer"`
		} `json:"menu"`
		PlaylistItemData struct {
			PlaylistSetVideoId string `json:"playlistSetVideoId"`
			VideoId            string `json:"videoId"`
		} `json:"playlistItemData"`
		MultiSelectCheckbox struct {
			CheckboxRenderer struct {
				OnSelectionChangeCommand struct {
					ClickTrackingParams           string `json:"clickTrackingParams"`
					UpdateMultiSelectStateCommand struct {
						MultiSelectParams string `json:"multiSelectParams"`
						MultiSelectItem   string `json:"multiSelectItem"`
					} `json:"updateMultiSelectStateCommand"`
				} `json:"onSelectionChangeCommand"`
				CheckedState   string `json:"checkedState"`
				TrackingParams string `json:"trackingParams"`
			} `json:"checkboxRenderer"`
		} `json:"multiSelectCheckbox"`
	} `json:"musicResponsiveListItemRenderer"`
}

type MusicResponsiveListItemRenderer struct {
	TrackingParams string `json:"trackingParams"`
	Thumbnail      struct {
		MusicThumbnailRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					Url    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			ThumbnailCrop  string `json:"thumbnailCrop"`
			ThumbnailScale string `json:"thumbnailScale"`
			TrackingParams string `json:"trackingParams"`
		} `json:"musicThumbnailRenderer"`
	} `json:"thumbnail"`
	Overlay struct {
		MusicItemThumbnailOverlayRenderer struct {
			Background struct {
				VerticalGradient struct {
					GradientLayerColors []string `json:"gradientLayerColors"`
				} `json:"verticalGradient"`
			} `json:"background"`
			Content struct {
				MusicPlayButtonRenderer struct {
					PlayNavigationEndpoint struct {
						ClickTrackingParams   string `json:"clickTrackingParams"`
						WatchPlaylistEndpoint struct {
							PlaylistId string `json:"playlistId"`
							Params     string `json:"params,omitempty"`
						} `json:"watchPlaylistEndpoint,omitempty"`
						WatchEndpoint struct {
							VideoId            string `json:"videoId"`
							PlaylistId         string `json:"playlistId"`
							PlaylistSetVideoId string `json:"playlistSetVideoId"`
							LoggingContext     struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointMusicSupportedConfigs struct {
								WatchEndpointMusicConfig struct {
									MusicVideoType string `json:"musicVideoType"`
								} `json:"watchEndpointMusicConfig"`
							} `json:"watchEndpointMusicSupportedConfigs"`
							Params string `json:"params,omitempty"`
						} `json:"watchEndpoint,omitempty"`
					} `json:"playNavigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
					PlayIcon       struct {
						IconType string `json:"iconType"`
					} `json:"playIcon"`
					PauseIcon struct {
						IconType string `json:"iconType"`
					} `json:"pauseIcon"`
					IconColor             int64 `json:"iconColor"`
					BackgroundColor       int   `json:"backgroundColor"`
					ActiveBackgroundColor int   `json:"activeBackgroundColor"`
					LoadingIndicatorColor int64 `json:"loadingIndicatorColor"`
					PlayingIcon           struct {
						IconType string `json:"iconType"`
					} `json:"playingIcon"`
					IconLoadingColor      int    `json:"iconLoadingColor"`
					ActiveScaleFactor     int    `json:"activeScaleFactor"`
					ButtonSize            string `json:"buttonSize"`
					RippleTarget          string `json:"rippleTarget"`
					AccessibilityPlayData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityPlayData"`
					AccessibilityPauseData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityPauseData"`
				} `json:"musicPlayButtonRenderer"`
			} `json:"content"`
			ContentPosition string `json:"contentPosition"`
			DisplayStyle    string `json:"displayStyle"`
		} `json:"musicItemThumbnailOverlayRenderer"`
	} `json:"overlay"`
	FixedColumns []struct {
		MusicResponsiveListItemFixedColumnRenderer struct {
			Text struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"text"`
			DisplayPriority string `json:"displayPriority"`
			Size            string `json:"size"`
		} `json:"musicResponsiveListItemFixedColumnRenderer"`
	} `json:"fixedColumns"`
	FlexColumns []struct {
		MusicResponsiveListItemFlexColumnRenderer struct {
			Text struct {
				Runs []Runs `json:"runs"`
			} `json:"text"`
			DisplayPriority string `json:"displayPriority"`
		} `json:"musicResponsiveListItemFlexColumnRenderer"`
	} `json:"flexColumns"`
	Menu struct {
		MenuRenderer struct {
			Items []struct {
				MenuNavigationItemRenderer struct {
					Text struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					NavigationEndpoint struct {
						ClickTrackingParams   string `json:"clickTrackingParams"`
						WatchPlaylistEndpoint struct {
							PlaylistId string `json:"playlistId"`
							Params     string `json:"params"`
						} `json:"watchPlaylistEndpoint,omitempty"`
						ModalEndpoint struct {
							Modal struct {
								ModalWithTitleAndButtonRenderer struct {
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Content struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"content"`
									Button struct {
										ButtonRenderer struct {
											Style      string `json:"style"`
											IsDisabled bool   `json:"isDisabled"`
											Text       struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												SignInEndpoint      struct {
													Hack bool `json:"hack"`
												} `json:"signInEndpoint"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"buttonRenderer"`
									} `json:"button"`
								} `json:"modalWithTitleAndButtonRenderer"`
							} `json:"modal"`
						} `json:"modalEndpoint,omitempty"`
						BrowseEndpoint struct {
							BrowseId                              string `json:"browseId"`
							BrowseEndpointContextSupportedConfigs struct {
								BrowseEndpointContextMusicConfig struct {
									PageType string `json:"pageType"`
								} `json:"browseEndpointContextMusicConfig"`
							} `json:"browseEndpointContextSupportedConfigs"`
						} `json:"browseEndpoint,omitempty"`
						ShareEntityEndpoint struct {
							SerializedShareEntity string `json:"serializedShareEntity"`
							SharePanelType        string `json:"sharePanelType"`
						} `json:"shareEntityEndpoint,omitempty"`
						WatchEndpoint struct {
							VideoId        string `json:"videoId"`
							PlaylistId     string `json:"playlistId"`
							Params         string `json:"params"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointMusicSupportedConfigs struct {
								WatchEndpointMusicConfig struct {
									MusicVideoType string `json:"musicVideoType"`
								} `json:"watchEndpointMusicConfig"`
							} `json:"watchEndpointMusicSupportedConfigs"`
						} `json:"watchEndpoint,omitempty"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"menuNavigationItemRenderer,omitempty"`
				MenuServiceItemRenderer struct {
					Text struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						QueueAddEndpoint    struct {
							QueueTarget struct {
								PlaylistId   string `json:"playlistId"`
								OnEmptyQueue struct {
									ClickTrackingParams string `json:"clickTrackingParams"`
									WatchEndpoint       struct {
										PlaylistId string `json:"playlistId,omitempty"`
										VideoId    string `json:"videoId,omitempty"`
									} `json:"watchEndpoint"`
								} `json:"onEmptyQueue"`
								VideoId string `json:"videoId,omitempty"`
							} `json:"queueTarget"`
							QueueInsertPosition string `json:"queueInsertPosition"`
							Commands            []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								AddToToastAction    struct {
									Item struct {
										NotificationTextRenderer struct {
											SuccessResponseText struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"successResponseText"`
											TrackingParams string `json:"trackingParams"`
										} `json:"notificationTextRenderer"`
									} `json:"item"`
								} `json:"addToToastAction"`
							} `json:"commands"`
						} `json:"queueAddEndpoint"`
					} `json:"serviceEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"menuServiceItemRenderer,omitempty"`
				ToggleMenuServiceItemRenderer struct {
					DefaultText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"defaultText"`
					DefaultIcon struct {
						IconType string `json:"iconType"`
					} `json:"defaultIcon"`
					DefaultServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						ModalEndpoint       struct {
							Modal struct {
								ModalWithTitleAndButtonRenderer struct {
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Content struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"content"`
									Button struct {
										ButtonRenderer struct {
											Style      string `json:"style"`
											IsDisabled bool   `json:"isDisabled"`
											Text       struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												SignInEndpoint      struct {
													Hack bool `json:"hack"`
												} `json:"signInEndpoint"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"buttonRenderer"`
									} `json:"button"`
								} `json:"modalWithTitleAndButtonRenderer"`
							} `json:"modal"`
						} `json:"modalEndpoint"`
					} `json:"defaultServiceEndpoint"`
					ToggledText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"toggledText"`
					ToggledIcon struct {
						IconType string `json:"iconType"`
					} `json:"toggledIcon"`
					ToggledServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						LikeEndpoint        struct {
							Status string `json:"status"`
							Target struct {
								PlaylistId string `json:"playlistId"`
							} `json:"target"`
						} `json:"likeEndpoint"`
					} `json:"toggledServiceEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"toggleMenuServiceItemRenderer,omitempty"`
			} `json:"items"`
			TrackingParams string `json:"trackingParams"`
			Accessibility  struct {
				AccessibilityData struct {
					Label string `json:"label"`
				} `json:"accessibilityData"`
			} `json:"accessibility"`
			TopLevelButtons []struct {
				LikeButtonRenderer struct {
					Target struct {
						VideoId string `json:"videoId"`
					} `json:"target"`
					LikeStatus                string `json:"likeStatus"`
					TrackingParams            string `json:"trackingParams"`
					LikesAllowed              bool   `json:"likesAllowed"`
					DislikeNavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						ModalEndpoint       struct {
							Modal struct {
								ModalWithTitleAndButtonRenderer struct {
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Content struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"content"`
									Button struct {
										ButtonRenderer struct {
											Style      string `json:"style"`
											IsDisabled bool   `json:"isDisabled"`
											Text       struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												SignInEndpoint      struct {
													Hack bool `json:"hack"`
												} `json:"signInEndpoint"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"buttonRenderer"`
									} `json:"button"`
								} `json:"modalWithTitleAndButtonRenderer"`
							} `json:"modal"`
						} `json:"modalEndpoint"`
					} `json:"dislikeNavigationEndpoint"`
					LikeCommand struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						ModalEndpoint       struct {
							Modal struct {
								ModalWithTitleAndButtonRenderer struct {
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Content struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"content"`
									Button struct {
										ButtonRenderer struct {
											Style      string `json:"style"`
											IsDisabled bool   `json:"isDisabled"`
											Text       struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												SignInEndpoint      struct {
													Hack bool `json:"hack"`
												} `json:"signInEndpoint"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"buttonRenderer"`
									} `json:"button"`
								} `json:"modalWithTitleAndButtonRenderer"`
							} `json:"modal"`
						} `json:"modalEndpoint"`
					} `json:"likeCommand"`
				} `json:"likeButtonRenderer"`
			} `json:"topLevelButtons,omitempty"`
		} `json:"menuRenderer"`
	} `json:"menu"`
	FlexColumnDisplayStyle string `json:"flexColumnDisplayStyle"`
	NavigationEndpoint     struct {
		ClickTrackingParams string `json:"clickTrackingParams"`
		WatchEndpoint       struct {
			VideoId        string `json:"videoId"`
			PlaylistId     string `json:"playlistId"`
			LoggingContext struct {
				VssLoggingContext struct {
					SerializedContextData string `json:"serializedContextData"`
				} `json:"vssLoggingContext"`
			} `json:"loggingContext"`
			WatchEndpointMusicSupportedConfigs struct {
				WatchEndpointMusicConfig struct {
					MusicVideoType string `json:"musicVideoType"`
				} `json:"watchEndpointMusicConfig"`
			} `json:"watchEndpointMusicSupportedConfigs"`
		} `json:"watchEndpoint,omitempty"`
		BrowseEndpoint struct {
			BrowseId                              string `json:"browseId"`
			BrowseEndpointContextSupportedConfigs struct {
				BrowseEndpointContextMusicConfig struct {
					PageType string `json:"pageType"`
				} `json:"browseEndpointContextMusicConfig"`
			} `json:"browseEndpointContextSupportedConfigs"`
		} `json:"browseEndpoint"`
	} `json:"navigationEndpoint"`
	ItemHeight       string `json:"itemHeight"`
	PlaylistItemData struct {
		PlaylistSetVideoID string `json:"playlistSetVideoId"`
		VideoID            string `json:"videoId"`
		VoteSortValue      int    `json:"voteSortValue"`
	} `json:"playlistItemData,omitempty"`
	Badges []struct {
		MusicInlineBadgeRenderer struct {
			TrackingParams string `json:"trackingParams"`
			Icon           struct {
				IconType string `json:"iconType"`
			} `json:"icon"`
			AccessibilityData struct {
				AccessibilityData struct {
					Label string `json:"label"`
				} `json:"accessibilityData"`
			} `json:"accessibilityData"`
		} `json:"musicInlineBadgeRenderer"`
	} `json:"badges,omitempty"`
}

type MusicCardShelfRenderer struct {
	TrackingParams string `json:"trackingParams"`
	Thumbnail      struct {
		MusicThumbnailRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					Url    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			ThumbnailCrop  string `json:"thumbnailCrop"`
			ThumbnailScale string `json:"thumbnailScale"`
			TrackingParams string `json:"trackingParams"`
		} `json:"musicThumbnailRenderer"`
	} `json:"thumbnail"`
	Title struct {
		Runs []struct {
			Text               string `json:"text"`
			NavigationEndpoint struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				BrowseEndpoint      struct {
					BrowseId                              string `json:"browseId"`
					BrowseEndpointContextSupportedConfigs struct {
						BrowseEndpointContextMusicConfig struct {
							PageType string `json:"pageType"`
						} `json:"browseEndpointContextMusicConfig"`
					} `json:"browseEndpointContextSupportedConfigs"`
				} `json:"browseEndpoint"`
			} `json:"navigationEndpoint"`
		} `json:"runs"`
	} `json:"title"`
	Subtitle struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"subtitle"`
	Contents []struct {
		MusicResponsiveListItemRenderer MusicResponsiveListItemRenderer `json:"musicResponsiveListItemRenderer"`
	} `json:"contents"`
	Buttons []struct {
		ButtonRenderer struct {
			Style string `json:"style"`
			Size  string `json:"size,omitempty"`
			Text  struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"text"`
			Icon struct {
				IconType string `json:"iconType"`
			} `json:"icon"`
			Accessibility struct {
				Label string `json:"label"`
			} `json:"accessibility"`
			TrackingParams    string `json:"trackingParams"`
			AccessibilityData struct {
				AccessibilityData struct {
					Label string `json:"label"`
				} `json:"accessibilityData"`
			} `json:"accessibilityData,omitempty"`
			Command struct {
				ClickTrackingParams   string `json:"clickTrackingParams"`
				WatchPlaylistEndpoint struct {
					PlaylistId string `json:"playlistId"`
					Params     string `json:"params"`
				} `json:"watchPlaylistEndpoint"`
			} `json:"command"`
		} `json:"buttonRenderer"`
	} `json:"buttons"`
	Menu struct {
		MenuRenderer struct {
			Items []struct {
				MenuNavigationItemRenderer struct {
					Text struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					NavigationEndpoint struct {
						ClickTrackingParams   string `json:"clickTrackingParams"`
						WatchPlaylistEndpoint struct {
							PlaylistId string `json:"playlistId"`
							Params     string `json:"params"`
						} `json:"watchPlaylistEndpoint,omitempty"`
						ShareEntityEndpoint struct {
							SerializedShareEntity string `json:"serializedShareEntity"`
							SharePanelType        string `json:"sharePanelType"`
						} `json:"shareEntityEndpoint,omitempty"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"menuNavigationItemRenderer,omitempty"`
				ToggleMenuServiceItemRenderer struct {
					DefaultText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"defaultText"`
					DefaultIcon struct {
						IconType string `json:"iconType"`
					} `json:"defaultIcon"`
					DefaultServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						ModalEndpoint       struct {
							Modal struct {
								ModalWithTitleAndButtonRenderer struct {
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Content struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"content"`
									Button struct {
										ButtonRenderer struct {
											Style      string `json:"style"`
											IsDisabled bool   `json:"isDisabled"`
											Text       struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												SignInEndpoint      struct {
													Hack bool `json:"hack"`
												} `json:"signInEndpoint"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"buttonRenderer"`
									} `json:"button"`
								} `json:"modalWithTitleAndButtonRenderer"`
							} `json:"modal"`
						} `json:"modalEndpoint"`
					} `json:"defaultServiceEndpoint"`
					ToggledText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"toggledText"`
					ToggledIcon struct {
						IconType string `json:"iconType"`
					} `json:"toggledIcon"`
					TrackingParams string `json:"trackingParams"`
				} `json:"toggleMenuServiceItemRenderer,omitempty"`
			} `json:"items"`
			TrackingParams string `json:"trackingParams"`
			Accessibility  struct {
				AccessibilityData struct {
					Label string `json:"label"`
				} `json:"accessibilityData"`
			} `json:"accessibility"`
		} `json:"menuRenderer"`
	} `json:"menu"`
	OnTap struct {
		ClickTrackingParams string `json:"clickTrackingParams"`
		BrowseEndpoint      struct {
			BrowseId                              string `json:"browseId"`
			BrowseEndpointContextSupportedConfigs struct {
				BrowseEndpointContextMusicConfig struct {
					PageType string `json:"pageType"`
				} `json:"browseEndpointContextMusicConfig"`
			} `json:"browseEndpointContextSupportedConfigs"`
		} `json:"browseEndpoint"`
	} `json:"onTap"`
	Header struct {
		MusicCardShelfHeaderBasicRenderer struct {
			Title struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"title"`
			TrackingParams string `json:"trackingParams"`
		} `json:"musicCardShelfHeaderBasicRenderer"`
	} `json:"header"`
	EndIcon struct {
		IconType string `json:"iconType"`
	} `json:"endIcon"`
}

type MusicShelfRenderer struct {
	Title struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	Contents []struct {
		MusicResponsiveListItemRenderer MusicResponsiveListItemRenderer `json:"musicResponsiveListItemRenderer"`
	} `json:"contents"`
	TrackingParams string `json:"trackingParams"`
	BottomText     struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"bottomText"`
	BottomEndpoint struct {
		ClickTrackingParams string `json:"clickTrackingParams"`
		SearchEndpoint      struct {
			Query  string `json:"query"`
			Params string `json:"params"`
		} `json:"searchEndpoint"`
		BrowseEndpoint struct {
			BrowseId                              string `json:"browseId"`
			Params                                string `json:"params"`
			BrowseEndpointContextSupportedConfigs struct {
				BrowseEndpointContextMusicConfig struct {
					PageType string `json:"pageType"`
				} `json:"browseEndpointContextMusicConfig"`
			} `json:"browseEndpointContextSupportedConfigs"`
		} `json:"browseEndpoint"`
	} `json:"bottomEndpoint"`
	Continuations []struct {
		NextContinuationData NextContinuationData `json:"nextContinuationData"`
	} `json:"continuations"`
	ShelfDivider struct {
		MusicShelfDividerRenderer struct {
			Hidden bool `json:"hidden"`
		} `json:"musicShelfDividerRenderer"`
	} `json:"shelfDivider"`
}

type NextContinuationData struct {
	Continuation        string `json:"continuation"`
	ClickTrackingParams string `json:"clickTrackingParams"`
}

type ItemSectionRenderer struct {
	Contents []struct {
		ShowingResultsForRenderer struct {
			ShowingResultsFor struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"showingResultsFor"`
			CorrectedQuery struct {
				Runs []struct {
					Text    string `json:"text"`
					Italics bool   `json:"italics,omitempty"`
				} `json:"runs"`
			} `json:"correctedQuery"`
			CorrectedQueryEndpoint struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				SearchEndpoint      struct {
					Query string `json:"query"`
				} `json:"searchEndpoint"`
			} `json:"correctedQueryEndpoint"`
			SearchInsteadFor struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"searchInsteadFor"`
			OriginalQuery struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"originalQuery"`
			OriginalQueryEndpoint struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				SearchEndpoint      struct {
					Query  string `json:"query"`
					Params string `json:"params"`
				} `json:"searchEndpoint"`
			} `json:"originalQueryEndpoint"`
			TrackingParams string `json:"trackingParams"`
		} `json:"showingResultsForRenderer"`
	} `json:"contents"`
	TrackingParams string `json:"trackingParams"`
}

type SectionListRendererContents struct {
	MusicCardShelfRenderer *MusicCardShelfRenderer `json:"musicCardShelfRenderer,omitempty"`
	MusicShelfRenderer     *MusicShelfRenderer     `json:"musicShelfRenderer,omitempty"`
	ItemSectionRenderer    *ItemSectionRenderer    `json:"itemSectionRenderer,omitempty"`
}

func (c *SectionListRendererContents) UnmarshalJSON(data []byte) error {

	if strings.Contains(string(data[0:50]), "musicCardShelfRenderer") {
		/*musicCardShelfRenderer := map[string]MusicCardShelfRenderer{}
		err := json.Unmarshal(data, &musicCardShelfRenderer)
		if err != nil {
			return err
		}
		c.MusicCardShelfRenderer = &(musicCardShelfRenderer["musicCardShelfRenderer"])*/
		return nil
	} else if strings.Contains(string(data[0:50]), "musicShelfRenderer") {
		musicShelfRenderer := map[string]MusicShelfRenderer{}
		err := json.Unmarshal(data, &musicShelfRenderer)
		if err != nil {
			return err
		}
		renderer := musicShelfRenderer["musicShelfRenderer"]
		c.MusicShelfRenderer = &renderer
		return nil
	} else if strings.Contains(string(data[0:50]), "itemSectionRenderer") {
		/*musicShelfRenderer := map[string]MusicShelfRenderer{}
		err := json.Unmarshal(data, &musicShelfRenderer)
		if err != nil {
			return err
		}
		renderer := musicShelfRenderer["musicShelfRenderer"]
		c.MusicShelfRenderer = &renderer*/
		return nil
	} else {
		//return fmt.Errorf("unknown type: %v", string(data[0:50]))
		return nil
	}
}

type TabbedSearchResultsRenderer struct {
	Tabs []struct {
		TabRenderer struct {
			Title    string `json:"title"`
			Selected bool   `json:"selected"`
			Content  struct {
				SectionListRenderer struct {
					SectionListRendererContents []SectionListRendererContents `json:"contents"`
					TrackingParams              string                        `json:"trackingParams"`
					Header                      struct {
						ChipCloudRenderer struct {
							Chips []struct {
								ChipCloudChipRenderer struct {
									Style struct {
										StyleType string `json:"styleType"`
									} `json:"style"`
									Text struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										SearchEndpoint      struct {
											Query  string `json:"query"`
											Params string `json:"params"`
										} `json:"searchEndpoint"`
									} `json:"navigationEndpoint"`
									TrackingParams    string `json:"trackingParams"`
									AccessibilityData struct {
										AccessibilityData struct {
											Label string `json:"label"`
										} `json:"accessibilityData"`
									} `json:"accessibilityData"`
									IsSelected bool   `json:"isSelected"`
									UniqueId   string `json:"uniqueId"`
								} `json:"chipCloudChipRenderer"`
							} `json:"chips"`
							CollapsedRowCount    int    `json:"collapsedRowCount"`
							TrackingParams       string `json:"trackingParams"`
							HorizontalScrollable bool   `json:"horizontalScrollable"`
						} `json:"chipCloudRenderer"`
					} `json:"header"`
				} `json:"sectionListRenderer"`
			} `json:"content"`
			TabIdentifier  string `json:"tabIdentifier"`
			TrackingParams string `json:"trackingParams"`
		} `json:"tabRenderer"`
	} `json:"tabs"`
}
