package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"musiclibmngr/internal/importer"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// curl -G "https://musicbrainz.org/ws/2/release" \
// 	--data-urlencode "fmt=json" \
// 	--data-urlencode "query=artist:nirvana AND release-group:nevermindd AND tracks:12" \
// 	--data-urlencode "limit=4"

type MusicBrainzService struct {
	client     *http.Client
	baseURL    string
	maxRetries int
}

func NewMusicBrainzServirce(client *http.Client, maxRetries int) *MusicBrainzService {
	return &MusicBrainzService{
		client:     client,
		maxRetries: maxRetries,
		baseURL:    "https://musicbrainz.org/ws/2/",
	}
}

func (mbs *MusicBrainzService) ReleaseQuery(music importer.ReleaseInfo) string {
	query_slice := []string{}
	if len(music.Artist) > 0 {
		query_slice = append(query_slice, "artist:"+music.Artist)
	}
	if len(music.Title) > 0 {
		query_slice = append(query_slice, "release-group:"+music.Title)
	}
	if music.TrackNb > 0 {
		query_slice = append(query_slice, "tracks:"+strconv.Itoa(music.TrackNb))
	}
	qry := strings.Join(query_slice, " AND ")

	params := url.Values{}
	params.Set("query", qry)
	params.Set("fmt", "json")
	params.Set("limit", "4")
	return mbs.baseURL + "release" + "?" + params.Encode()
}

func (mbs *MusicBrainzService) RunQuery(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := mbs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var body []byte
		body, err = io.ReadAll(resp.Body)
		return body, err
	default:
		return nil, fmt.Errorf("Error: status code %d - %s\n", resp.StatusCode, resp.Status)
	}
}

type ReleaseResult struct {
	Id             string `json:"id"`
	Score          int    `json:"score"`
	ArtistCreditId string `json:"artist-credit-id"`
	Title          string `json:"title"`
}

type QueryResult struct {
	Count    int             `json:"count"`
	Offset   int             `json:"offset"`
	Releases []ReleaseResult `json:"releases"`
}

func (mbs *MusicBrainzService) GetReleaseQuery(ctx context.Context, music importer.ReleaseInfo) (*QueryResult, error) {
	resp, err := mbs.RunQuery(ctx, mbs.ReleaseQuery(music))
	if err != nil {
		return nil, err
	}
	var data QueryResult
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (mbs *MusicBrainzService) ReleaseFetch(id string) string {
	params := url.Values{}
	params.Set("fmt", "json")
	params.Set("inc", "artist-credits+labels+discids+recordings")
	return mbs.baseURL + "release/" + id + "?" + params.Encode()
}

type ReleaseTrack struct {
	Track  int    `json:"position"`
	Title  string `json:"title"`
	Length int    `json:"length"`
}

type ReleaseMedia struct {
	Position int            `json:"position"`
	Tracks   []ReleaseTrack `json:"tracks"`
}

type Release struct {
	Title string         `json:"title"`
	Media []ReleaseMedia `json:"media"`
}

func (mbs *MusicBrainzService) GetRelease(ctx context.Context, id string) (*Release, error) {
	resp, err := mbs.RunQuery(ctx, mbs.ReleaseFetch(id))
	if err != nil {
		return nil, err
	}
	var data Release
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// type healthState struct {
// 	mu             sync.Mutex
// 	rateLimitUntil time.Time
// }

// var health healthState

// func setRateLimited(untill time.Time) {
// 	health.mu.Lock()
// 	health.rateLimitUntil = untill
// 	health.mu.Unlock()
// }

// func canCall() bool {
// 	health.mu.Lock()
// 	defer health.mu.Unlock()
// 	return time.Now().After(health.rateLimitUntil)
// }

// func waitUntil() time.Time {
// 	health.mu.Lock()
// 	defer health.mu.Unlock()
// 	return health.rateLimitUntil
// }

// func (s *MusicBrainzService) Query(ctx context.Context, music importer.ReleaseInfo) ([]byte, error) {
// 	fullURL := buildQuery(s.baseURL, toQuery(music))
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for attempt := 0; attempt <= s.maxRetries; attempt++ {
// 		if !canCall() {
// 			fmt.Printf("Rate limited until %v, skipping\n", waitUntil())
// 			select {
// 			case <-ctx.Done():
// 				return nil, ctx.Err()
// 			case <-time.After(time.Until(waitUntil())):
// 			}
// 			continue
// 		}

// 		if attempt > 0 {
// 			fmt.Println("Waiting before retry...")
// 			select {
// 			case <-ctx.Done():
// 				return nil, ctx.Err()
// 			case <-time.After(6 * time.Second):
// 			}
// 		}

// 		resp, err := s.client.Do(req)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer resp.Body.Close()

// 		switch resp.StatusCode {
// 		case http.StatusOK:
// 			var body []byte
// 			body, err = io.ReadAll(resp.Body)
// 			return body, err
// 		case http.StatusServiceUnavailable:
// 			// Parse Retry-After (seconds or HTTP date)
// 			waitStr := resp.Header.Get("Retry-After")
// 			var until time.Time
// 			if secs, err := strconv.ParseInt(waitStr, 10, 64); err == nil && secs > 0 {
// 				until = time.Now().Add(time.Duration(secs) * time.Second)
// 			} else {
// 				until = time.Now().Add(2 * time.Second) // fallback: MusicBrainz ~1 req/sec
// 			}
// 			setRateLimited(until)
// 			fmt.Printf("Retry-After : %v\n", resp.Header.Get("Retry-After"))
// 			continue
// 		default:
// 			return nil, fmt.Errorf("Error: status code %d - %s\n", resp.StatusCode, resp.Status)
// 		}
// 	}
// 	return nil, fmt.Errorf("max retries exceeded")
// }

// func buildQuery(baseURL string, qry string) string {
// 	params := url.Values{}
// 	params.Set("query", qry)
// 	params.Set("fmt", "json")
// 	params.Set("limit", "4")
// 	return baseURL + "?" + params.Encode()
// }

// func toQuery(music importer.ReleaseInfo) string {
// 	query_slice := []string{}
// 	if len(music.Artist) > 0 {
// 		query_slice = append(query_slice, "artist:"+music.Artist)
// 	}
// 	if len(music.Title) > 0 {
// 		query_slice = append(query_slice, "release-group:"+music.Title)
// 	}
// 	query_slice = append(query_slice, "tracks:"+strconv.Itoa(music.TrackNb))
// 	return strings.Join(query_slice, " AND ")
// }
