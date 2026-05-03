package services

import (
	"bytes"
	"context"
	"encoding/json"
	"musiclibmngr/internal/repo"
	"net/http"
	"testing"
	"time"
)

func TestBuildReleaseUrl(t *testing.T) {
	mf := repo.MusicFile{
		Artist:  "nirvana",
		Release: "nevermind",
	}

	client := &http.Client{
    	Timeout: 10 * time.Second,
	}
	mbzs := NewMusicBrainzRecordServirce(client, 5)
	ctx := context.Background()
	defer ctx.Done()
	body, err := mbzs.Query(ctx, mf)
	if err != nil {
		t.Errorf("Query returned error: %v", err)
	}
	var prettyBody bytes.Buffer
	if err := json.Indent(&prettyBody, body, "", "  "); err != nil {
		t.Errorf("Failed to pretty print JSON: %v", err)
	}
	t.Log(prettyBody.String())

}
