package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type MusicBrainzService struct {
	client     *http.Client
	baseURL    string
	maxRetries int
}

type healthState struct {
	mu             sync.Mutex
	rateLimitUntil time.Time
}

var health healthState

func setRateLimited(untill time.Time) {
	health.mu.Lock()
	health.rateLimitUntil = untill
	health.mu.Unlock()
}

func canCall() bool {
	health.mu.Lock()
	defer health.mu.Unlock()
	return time.Now().After(health.rateLimitUntil)
}

func waitUntil() time.Time {
	health.mu.Lock()
	defer health.mu.Unlock()
	return health.rateLimitUntil
}

func NewMusicBrainzRecordServirce(client *http.Client, maxRetries int) *MusicBrainzService {
	return &MusicBrainzService{
		client:     client,
		maxRetries: maxRetries,
		baseURL:    "https://musicbrainz.org/ws/2/release",
	}
}

func (s *MusicBrainzService) Query(ctx context.Context, album string) ([]byte, error) {
	fullURL := buildReleaseUrl(s.baseURL, album)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if !canCall() {
			fmt.Printf("%v : Rate limited until %v, skipping\n", time.Now(), waitUntil())
			// Wait exactly the health duration (not s.waitingTime)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Until(waitUntil())):
			}
			continue
		}

		if attempt > 0 {
			fmt.Println("Waiting before retry...")
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(6 * time.Second):
			}
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			var body []byte
			body, err = io.ReadAll(resp.Body)
			return body, err
		case http.StatusServiceUnavailable:
			// Parse Retry-After (seconds or HTTP date)
			waitStr := resp.Header.Get("Retry-After")
			var until time.Time
			if secs, err := strconv.ParseInt(waitStr, 10, 64); err == nil && secs > 0 {
				until = time.Now().Add(time.Duration(secs) * time.Second)
			} else {
				until = time.Now().Add(2 * time.Second) // fallback: MusicBrainz ~1 req/sec
			}
			setRateLimited(until)
			fmt.Printf("Retry-After : %v\n", resp.Header.Get("Retry-After"))
			continue
		default:
			return nil, fmt.Errorf("Error: status code %d - %s\n", resp.StatusCode, resp.Status)
		}
	}
	return nil, fmt.Errorf("max retries exceeded")
}

func buildReleaseUrl(baseURL string, album string) string {
	params := url.Values{}
	params.Set("query", album)
	params.Set("fmt", "json")
	return baseURL + "?" + params.Encode()
}
