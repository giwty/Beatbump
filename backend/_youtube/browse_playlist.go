package _youtube

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PlaylistResponse struct {
	ResponseContext struct {
		VisitorData           string `json:"visitorData"`
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
	} `json:"responseContext"`
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			SecondaryContents struct {
				SectionListRenderer struct {
					Contents      []PlaylistContents `json:"contents"`
					Continuations []struct {
						NextContinuationData struct {
							Continuation        string `json:"continuation"`
							ClickTrackingParams string `json:"clickTrackingParams"`
						} `json:"nextContinuationData"`
					} `json:"continuations"`
					TrackingParams string `json:"trackingParams"`
				} `json:"sectionListRenderer"`
			} `json:"secondaryContents"`
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								MusicResponsiveHeaderRenderer struct {
									Thumbnail struct {
										MusicThumbnailRenderer struct {
											Thumbnail struct {
												Thumbnails []struct {
													URL    string `json:"url"`
													Width  int    `json:"width"`
													Height int    `json:"height"`
												} `json:"thumbnails"`
											} `json:"thumbnail"`
											ThumbnailCrop  string `json:"thumbnailCrop"`
											ThumbnailScale string `json:"thumbnailScale"`
											TrackingParams string `json:"trackingParams"`
										} `json:"musicThumbnailRenderer"`
									} `json:"thumbnail"`
									Buttons []struct {
										ToggleButtonRenderer struct {
											IsToggled   bool `json:"isToggled"`
											IsDisabled  bool `json:"isDisabled"`
											DefaultIcon struct {
												IconType string `json:"iconType"`
											} `json:"defaultIcon"`
											ToggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"toggledIcon"`
											TrackingParams            string `json:"trackingParams"`
											DefaultNavigationEndpoint struct {
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
											} `json:"defaultNavigationEndpoint"`
											AccessibilityData struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"accessibilityData"`
											ToggledAccessibilityData struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"toggledAccessibilityData"`
										} `json:"toggleButtonRenderer,omitempty"`
										MusicPlayButtonRenderer struct {
											PlayNavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												WatchEndpoint       struct {
													VideoID        string `json:"videoId"`
													PlaylistID     string `json:"playlistId"`
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
											IconLoadingColor      int `json:"iconLoadingColor"`
											ActiveScaleFactor     int `json:"activeScaleFactor"`
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
										} `json:"musicPlayButtonRenderer,omitempty"`
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
															PlaylistID string `json:"playlistId"`
															Params     string `json:"params"`
														} `json:"watchPlaylistEndpoint"`
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
																PlaylistID   string `json:"playlistId"`
																OnEmptyQueue struct {
																	ClickTrackingParams string `json:"clickTrackingParams"`
																	WatchEndpoint       struct {
																		PlaylistID string `json:"playlistId"`
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
																PlaylistID string `json:"playlistId"`
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
										} `json:"menuRenderer,omitempty"`
									} `json:"buttons"`
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Subtitle struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"subtitle"`
									TrackingParams   string `json:"trackingParams"`
									StraplineTextOne struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"straplineTextOne"`
									Description struct {
										MusicDescriptionShelfRenderer struct {
											Description struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"description"`
											MoreButton struct {
												ToggleButtonRenderer struct {
													IsToggled   bool `json:"isToggled"`
													IsDisabled  bool `json:"isDisabled"`
													DefaultIcon struct {
														IconType string `json:"iconType"`
													} `json:"defaultIcon"`
													DefaultText struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"defaultText"`
													ToggledIcon struct {
														IconType string `json:"iconType"`
													} `json:"toggledIcon"`
													ToggledText struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"toggledText"`
													TrackingParams string `json:"trackingParams"`
												} `json:"toggleButtonRenderer"`
											} `json:"moreButton"`
											TrackingParams string `json:"trackingParams"`
											ShelfStyle     string `json:"shelfStyle"`
										} `json:"musicDescriptionShelfRenderer"`
									} `json:"description"`
									SecondSubtitle struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"secondSubtitle"`
								} `json:"musicResponsiveHeaderRenderer"`
							} `json:"contents"`
							TrackingParams string `json:"trackingParams"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
					TrackingParams string `json:"trackingParams"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
		SingleColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents      []PlaylistContents `json:"contents"`
							Continuations []struct {
								NextContinuationData struct {
									Continuation        string `json:"continuation"`
									ClickTrackingParams string `json:"clickTrackingParams"`
								} `json:"nextContinuationData"`
							} `json:"continuations"`
							TrackingParams string `json:"trackingParams"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
					TrackingParams string `json:"trackingParams"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"singleColumnBrowseResultsRenderer"`
	} `json:"contents"`
	OnResponseReceivedActions []struct {
		ClickTrackingParams           string `json:"clickTrackingParams"`
		AppendContinuationItemsAction struct {
			ContinuationItems []struct {
				MusicResponsiveListItemRenderer MusicResponsiveListItemRenderer `json:"musicResponsiveListItemRenderer,omitempty"` 	 
				ContinuationItemRenderer        *ContinuationItemRenderer       `json:"continuationItemRenderer,omitempty"`
			} `json:"continuationItems"`
			TargetID string `json:"targetId"`
		} `json:"appendContinuationItemsAction"`
	} `json:"onResponseReceivedActions"`
	ContinuationContents struct {
		SectionListContinuation struct {
			Contents []struct {
				MusicCarouselShelfRenderer struct {
					Header struct {
						MusicCarouselShelfBasicHeaderRenderer struct {
							Title struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"title"`
							AccessibilityData struct {
								AccessibilityData struct {
									Label string `json:"label"`
								} `json:"accessibilityData"`
							} `json:"accessibilityData"`
							HeaderStyle    string `json:"headerStyle"`
							TrackingParams string `json:"trackingParams"`
						} `json:"musicCarouselShelfBasicHeaderRenderer"`
					} `json:"header"`
					Contents []struct {
						MusicTwoRowItemRenderer *MusicTwoRowItemRenderer `json:"musicTwoRowItemRenderer,omitempty"`
					} `json:"contents"`
					TrackingParams string `json:"trackingParams"`
					ItemSize       string `json:"itemSize"`
				} `json:"musicCarouselShelfRenderer"`
			} `json:"contents"`
			TrackingParams string `json:"trackingParams"`
		} `json:"sectionListContinuation"`
	} `json:"continuationContents"`
	Header struct {
		MusicHeaderRenderer struct {
			Title struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"title"`
			TrackingParams string `json:"trackingParams"`
		} `json:"musicHeaderRenderer"`
		MusicDetailHeaderRenderer struct {
			Title struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"title"`
			Subtitle struct {
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
					} `json:"navigationEndpoint,omitempty"`
				} `json:"runs"`
			} `json:"subtitle"`
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
										PlaylistId   string `json:"playlistId"`
										OnEmptyQueue struct {
											ClickTrackingParams string `json:"clickTrackingParams"`
											WatchEndpoint       struct {
												PlaylistId string `json:"playlistId"`
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
						ButtonRenderer struct {
							Style string `json:"style"`
							Size  string `json:"size"`
							Text  struct {
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
								} `json:"watchPlaylistEndpoint"`
							} `json:"navigationEndpoint"`
							Accessibility struct {
								Label string `json:"label"`
							} `json:"accessibility"`
							TrackingParams    string `json:"trackingParams"`
							AccessibilityData struct {
								AccessibilityData struct {
									Label string `json:"label"`
								} `json:"accessibilityData"`
							} `json:"accessibilityData"`
						} `json:"buttonRenderer,omitempty"`
						ToggleButtonRenderer struct {
							IsToggled   bool `json:"isToggled"`
							IsDisabled  bool `json:"isDisabled"`
							DefaultIcon struct {
								IconType string `json:"iconType"`
							} `json:"defaultIcon"`
							DefaultText struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
								Accessibility struct {
									AccessibilityData struct {
										Label string `json:"label"`
									} `json:"accessibilityData"`
								} `json:"accessibility"`
							} `json:"defaultText"`
							ToggledIcon struct {
								IconType string `json:"iconType"`
							} `json:"toggledIcon"`
							ToggledText struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
								Accessibility struct {
									AccessibilityData struct {
										Label string `json:"label"`
									} `json:"accessibilityData"`
								} `json:"accessibility"`
							} `json:"toggledText"`
							TrackingParams            string `json:"trackingParams"`
							DefaultNavigationEndpoint struct {
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
							} `json:"defaultNavigationEndpoint"`
						} `json:"toggleButtonRenderer,omitempty"`
					} `json:"topLevelButtons"`
					Accessibility struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibility"`
				} `json:"menuRenderer"`
			} `json:"menu"`
			Thumbnail struct {
				CroppedSquareThumbnailRenderer struct {
					Thumbnail struct {
						Thumbnails []struct {
							Url    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"thumbnails"`
					} `json:"thumbnail"`
					TrackingParams string `json:"trackingParams"`
				} `json:"croppedSquareThumbnailRenderer"`
			} `json:"thumbnail"`
			TrackingParams string `json:"trackingParams"`
			MoreButton     struct {
				ToggleButtonRenderer struct {
					IsToggled   bool `json:"isToggled"`
					IsDisabled  bool `json:"isDisabled"`
					DefaultIcon struct {
						IconType string `json:"iconType"`
					} `json:"defaultIcon"`
					DefaultText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"defaultText"`
					ToggledIcon struct {
						IconType string `json:"iconType"`
					} `json:"toggledIcon"`
					ToggledText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"toggledText"`
					TrackingParams string `json:"trackingParams"`
				} `json:"toggleButtonRenderer"`
			} `json:"moreButton"`
			SecondSubtitle struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"secondSubtitle"`
		} `json:"musicDetailHeaderRenderer"`
	} `json:"header"`
	TrackingParams string `json:"trackingParams"`
}

type ContinuationItemRenderer struct {
	Trigger              string `json:"trigger"`
	ContinuationEndpoint struct {
		ClickTrackingParams string `json:"clickTrackingParams"`
		ContinuationCommand struct {
			Token   string `json:"token"`
			Request string `json:"request"`
		} `json:"continuationCommand"`
	} `json:"continuationEndpoint"`
}

type MusicPlaylistShelfRenderer struct {
	PlaylistId string `json:"playlistId"`
	Contents   []struct {
		MusicResponsiveListItemRenderer MusicResponsiveListItemRenderer `json:"musicResponsiveListItemRenderer"`
		ContinuationItemRenderer        *ContinuationItemRenderer       `json:"continuationItemRenderer,omitempty"`
	} `json:"contents"`
	Continuations []struct {
		NextContinuationData struct {
			Continuation        string `json:"continuation"`
			ClickTrackingParams string `json:"clickTrackingParams"`
		} `json:"nextContinuationData"`
	} `json:"continuations"`
	CollapsedItemCount      int    `json:"collapsedItemCount"`
	TrackingParams          string `json:"trackingParams"`
	ContentsMultiSelectable bool   `json:"contentsMultiSelectable"`
}

type PlaylistContents struct {
	MusicPlaylistShelfRenderer *MusicPlaylistShelfRenderer `json:"musicPlaylistShelfRenderer,omitempty"`
	GridRenderer               *GridRenderer               `json:"gridRenderer,omitempty"`
}

func (c *PlaylistContents) UnmarshalJSON(data []byte) error {

	if strings.Contains(string(data[0:80]), "musicPlaylistShelfRenderer") {
		musicResponsiveListRenderer := map[string]MusicPlaylistShelfRenderer{}
		err := json.Unmarshal(data, &musicResponsiveListRenderer)
		if err != nil {
			return err
		}
		renderer := musicResponsiveListRenderer["musicPlaylistShelfRenderer"]
		c.MusicPlaylistShelfRenderer = &renderer
		return nil
	} else if strings.Contains(string(data[0:80]), "gridRenderer") {
		musicTwoRowRenderer := map[string]GridRenderer{}
		err := json.Unmarshal(data, &musicTwoRowRenderer)
		if err != nil {
			return err
		}
		renderer := musicTwoRowRenderer["gridRenderer"]
		c.GridRenderer = &renderer
		return nil
	} else {
		return fmt.Errorf("unknown type: %v", string(data[0:50]))
	}
}
