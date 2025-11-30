package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// Rate limiter: 20 calls per minute (1 call every 3 seconds)
var itunesRateLimiter = rate.NewLimiter(rate.Every(3*time.Second), 1)

type iTunesResponse struct {
	ResultCount int          `json:"resultCount"`
	Results     []iTunesItem `json:"results"`
}

type iTunesItem struct {
	ArtistName       string `json:"artistName"`
	TrackName        string `json:"trackName"`
	CollectionName   string `json:"collectionName"`
	ArtworkUrl100    string `json:"artworkUrl100"`
	PrimaryGenreName string `json:"primaryGenreName"`
	ReleaseDate      string `json:"releaseDate"` // ISO 8601 format: 2005-03-01T08:00:00Z
}

// FetchMetadata attempts to find better metadata for a song using the iTunes Search API.
// It respects rate limits and returns the original metadata if no match is found or an error occurs.
func FetchMetadata(artist, title string) (*AudioMetadata, error) {
	// Wait for rate limiter
	if err := itunesRateLimiter.Wait(context.Background()); err != nil {
		return nil, err
	}

	// Clean up search terms
	cleanTitle := cleanString(title)
	cleanArtist := cleanString(artist)

	// Construct query
	query := fmt.Sprintf("%s %s", cleanArtist, cleanTitle)
	encodedQuery := url.QueryEscape(query)
	// Detect country based on script to improve search results
	countryCode := detectStoreCountry(query)
	countryParam := ""
	if countryCode != "" {
		countryParam = fmt.Sprintf("&country=%s", countryCode)
	}

	apiURL := fmt.Sprintf("https://itunes.apple.com/search?term=%s&entity=song&limit=1%s", encodedQuery, countryParam)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes API returned status: %d", resp.StatusCode)
	}

	var result iTunesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.ResultCount == 0 || len(result.Results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	item := result.Results[0]

	// Parse Year from ReleaseDate
	year := ""
	if len(item.ReleaseDate) >= 4 {
		year = item.ReleaseDate[:4]
	}

	// Get high-res artwork (replace 100x100 with 600x600)
	artworkURL := strings.Replace(item.ArtworkUrl100, "100x100bb", "600x600bb", 1)

	return &AudioMetadata{
		Title:      item.TrackName,
		Artist:     item.ArtistName,
		Album:      item.CollectionName,
		ArtworkURL: artworkURL,
		Genre:      item.PrimaryGenreName,
		Year:       year,
	}, nil
}

// cleanString removes common noise from YouTube titles
func cleanString(s string) string {
	// Remove things in brackets/parentheses like (Official Video), [Lyrics], etc.
	re := regexp.MustCompile(`(?i)(\(|\[)(official|video|audio|lyrics|hq|hd|4k|music video).*?(\)|\])`)
	s = re.ReplaceAllString(s, "")

	// Remove "ft.", "feat."
	reFeat := regexp.MustCompile(`(?i)\s(ft\.|feat\.|featuring)\s.*`)
	s = reFeat.ReplaceAllString(s, "")

	return strings.TrimSpace(s)
}

// detectStoreCountry returns an iTunes country code based on the script detected in the string.
// Returns empty string if no specific script is detected (defaults to US/Global).
func detectStoreCountry(s string) string {
	for _, r := range s {
		switch {
		case r >= 0x0590 && r <= 0x05FF: // Hebrew
			return "IL"
		case r >= 0x0400 && r <= 0x04FF: // Cyrillic
			return "RU"
		case r >= 0x0600 && r <= 0x06FF: // Arabic
			return "EG"
		case r >= 0x3040 && r <= 0x309F: // Hiragana
			return "JP"
		case r >= 0x30A0 && r <= 0x30FF: // Katakana
			return "JP"
		case r >= 0x4E00 && r <= 0x9FFF: // CJK Unified Ideographs (Kanji)
			// CJK is shared, but JP is a good default for music metadata in this context.
			// Could be refined if needed.
			return "JP"
		}
	}
	return ""
}
