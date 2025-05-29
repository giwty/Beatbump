package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type TokenInfo struct {
	Updated     int    `json:"updated"`
	Potoken     string `json:"poToken"`
	VisitorData string `json:"visitorData"`
}

func (t *TokenInfo) String() string {
	// Convert the TokenInfo struct to JSON string
	tokenJSON, err := json.Marshal(t)
	if err != nil {
		log.Printf("Error marshalling TokenInfo to JSON: %v", err)
		return ""
	}
	return string(tokenJSON)
}

type PotokenExtractor struct {
	updateInterval  float64
	ongoingUpdate   sync.Mutex
	updateRequested *sync.Cond
	tokenInfo       *TokenInfo
}

func NewPotokenExtractor(updateInterval float64) *PotokenExtractor {
	extractor := &PotokenExtractor{
		updateInterval: updateInterval,
	}
	extractor.updateRequested = sync.NewCond(&sync.Mutex{})
	return extractor
}

func (e *PotokenExtractor) Get() *TokenInfo {
	return e.tokenInfo
}

func (e *PotokenExtractor) RunOnce() (*TokenInfo, error) {
	err := e.update()
	if err != nil {
		return nil, err
	}
	return e.Get(), nil
}

func (e *PotokenExtractor) Run() {
	// Continuously run the extractor with a scheduled update interval
	for {
		select {
		case <-time.After(time.Duration(e.updateInterval) * time.Second):
			log.Println("Initiating scheduled update.")
			_ = e.update()
		}
	}
}

func (e *PotokenExtractor) RequestUpdate() bool {
	// Request an immediate update
	e.ongoingUpdate.Lock()
	defer e.ongoingUpdate.Unlock()

	if e.updateRequested != nil {
		log.Println("Force update has already been requested")
		return false
	}

	e.updateRequested.Signal()
	log.Println("Force update requested")
	return true
}

func (e *PotokenExtractor) update() error {
	// Simulate async update with a timeout
	timeout := time.After(30 * time.Second)
	done := make(chan error)

	go func() {
		err := e.performUpdate()
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-timeout:
		log.Println("Update failed: timeout exceeded.")
		return fmt.Errorf("update failed: timeout exceeded")
	}
}

func (e *PotokenExtractor) performUpdate() error {
	e.ongoingUpdate.Lock()
	defer e.ongoingUpdate.Unlock()

	log.Println("Update started")

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
	)

	// Create a new browser context with the options above
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// Create a new browser instance
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	done := make(chan bool)

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Add handler for network requests
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			//println(ev.Request.URL)
			if ev.Request.Method == "POST" && strings.Contains(ev.Request.URL, "/youtubei/v1/player") {
				tokenInfo := e.extractToken(ev)
				if tokenInfo != nil {
					e.tokenInfo = tokenInfo
				}
				close(done)
			}
		}
	})

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.youtube.com/embed/jNQXAC9IVRw"),
		chromedp.WaitVisible(`#movie_player`),
	)
	if err != nil {
		return fmt.Errorf("Failed trying to load page: %w", err)
	}
	time.Sleep(2 * time.Second)
	err = chromedp.Run(ctx,
		chromedp.Click("#movie_player", chromedp.NodeVisible),
	)
	if err != nil {
		return fmt.Errorf("Failed while trying to simulate click: %w", err)
	}

	<-done

	return nil
}

func (e *PotokenExtractor) extractToken(request *network.EventRequestWillBeSent) *TokenInfo {
	// Get the POST data entries (this is a slice of strings)
	postDataEntries := request.Request.PostDataEntries
	if len(postDataEntries) == 0 {
		log.Println("No PostDataEntries found in the request")
		return nil
	}

	// The postDataEntries contain the body of the POST request
	postData := postDataEntries[0] // Assume the first entry contains the body

	dec, err := base64.StdEncoding.DecodeString(postData.Bytes)
	if err != nil {
		return nil
	}

	var postDataJSON map[string]interface{}
	if err := json.Unmarshal(dec, &postDataJSON); err != nil {
		log.Printf("Failed to extract token from request: %s\n", err)
		return nil
	}

	// Extract the visitorData and potoken from the post data
	visitorData, ok := postDataJSON["context"].(map[string]interface{})["client"].(map[string]interface{})["visitorData"].(string)
	if !ok {
		log.Println("Failed to extract visitor data")
		return nil
	}

	// First, get serviceIntegrityDimensions
	serviceDims, ok := postDataJSON["serviceIntegrityDimensions"].(map[string]interface{})
	if !ok {
		log.Println("Failed to extract serviceIntegrityDimensions")
		return nil
	}

	// Then get poToken
	potoken, ok := serviceDims["poToken"].(string)
	if !ok {
		log.Println("Failed to extract poToken from serviceIntegrityDimensions")
		return nil
	}

	if len(potoken) < 100 {
		fmt.Println("Warnning: likely not a valid po token")
		return nil
	}

	potoken, err = url.QueryUnescape(potoken)
	if err != nil {
		fmt.Println("Error decoding potoken:", err)
		return nil
	}

	tokenInfo := &TokenInfo{
		Updated:     int(time.Now().Unix()),
		Potoken:     potoken,
		VisitorData: visitorData,
	}
	return tokenInfo
}
