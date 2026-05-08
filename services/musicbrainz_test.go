package services

import (
	"context"
	"musiclibmngr/internal/importer"
	"net/http"
	"testing"
	"time"
)

func TestReleaseQuery(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	t.Run("Nevermind", func(t *testing.T) {
		mf := importer.ReleaseInfo{
			Artist:  "nirvana",
			Title:   "nevermind",
			TrackNb: 12,
		}
		mbzs := NewMusicBrainzServirce(client, 5)
		ctx := context.Background()
		defer ctx.Done()
		body, err := mbzs.GetReleaseQuery(ctx, mf)
		if err != nil {
			t.Errorf("Query returned error: %v", err)
		}
		if len(body.Releases) != 4 {
			t.Errorf("Query data mismatch.\nExpected: 4 release\nActual:   %#v", &body)
		}
	})
}

func TestReleaseFetch(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	t.Run("Nevermind", func(t *testing.T) {
		mbzs := NewMusicBrainzServirce(client, 5)
		ctx := context.Background()
		defer ctx.Done()
		body, err := mbzs.GetRelease(ctx, "f922ec87-4758-421d-a839-3193455345ff")
		if err != nil {
			t.Errorf("Query returned error: %v", err)
		}
		if body.Title != "Nevermind" {
			t.Errorf("Expected Nervmind, got %s", body.Title)
		}
		if len(body.Media) != 1 {
			t.Errorf("Expected 1 media, got %v", body.Media)
		}
		if len(body.Media[0].Tracks) != 12 {
			t.Errorf("Expected 12 tracks, got %v", body.Media[0])
		}
		if body.Media[0].Tracks[2].Title != "Come as You Are" {
			t.Errorf("Expected 12 tracks, got %v", body.Media[0].Tracks[2])
		}
	})
}
