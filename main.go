package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"musiclibmngr/services"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.senan.xyz/taglib"
)

func readFileMeta(path string) {
	tags, err := taglib.ReadTags(path)
	if err != nil {
		fmt.Printf("error reading tags: %v\n", err)
		return
	}

	for k, v := range tags {
		fmt.Printf("%s: %q\n", k, v)
	}
}

func detectType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := io.ReadFull(f, buf)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}

func isAudio(path string) bool {
	t, err := detectType(path)
	if err != nil {
		fmt.Printf("error reading type: %v\n", err)
		return false
	}
	if t == "audio/mpeg" {
		return true
	}
	if strings.HasSuffix(path, ".flac") && t == "application/octet-stream" {
		return true
	}
	return false
}

func main() {
	client := &http.Client{Timeout: 30 * time.Second}
	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isAudio(path) { // Only files (-type f)
			fmt.Println("---------------------------")
			fmt.Println(path)
			readFileMeta(path)
			album := filepath.Base(path)
			ctx := context.Background()
			recordService := services.NewMusicBrainzRecordServirce(client, 5)
			buff, err := recordService.Query(ctx, album)
			if err != nil {
				return err
			}
			fmt.Printf(string(buff))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
