package main

import (
	"context"
	"fmt"
	"io/fs"
	"musiclibmngr/internal/db"
	"musiclibmngr/internal/file"
	"path/filepath"
	"regexp"

	"go.senan.xyz/taglib"
	"gorm.io/gorm"
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

func parseMusicPath(path string) (artist, album, trackNum, title string) {
	re := regexp.MustCompile(`testdata/(.+)/(.+)/(\d+)\s+(.+)\.\w+$`)
	matches := re.FindStringSubmatch(path)
	if len(matches) == 5 {
		return matches[1], matches[2], matches[3], matches[4]
	}
	return "", "", "", ""
}

func fileScan(dbConn *db.DB) error {
	return filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		desc := file.NewDescriptor(path)
		isAudio, err := desc.IsAudio()
		if err != nil || !isAudio {
			return nil
		}

		fmt.Println(path)
		artistName, album, track, title := parseMusicPath(path)
		ctx := context.Background()
		artist, err := gorm.G[db.Artist](dbConn.DB).Where("name = ?", artistName).First(ctx)
		if err == gorm.ErrRecordNotFound {
			artist = db.Artist{Name: artistName}
			if err := gorm.G[db.Artist](dbConn.DB).Create(ctx, &artist); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		fmt.Printf("Artist: %s\nAlbum: %s\nTrack: %s\nTitle: %s\n",
			artistName, album, track, title)

		return nil
	})
}

func main() {

	dbConn, err := db.New("test.db")

	if err != nil {
		panic(err)
	}

	fileScan(dbConn)
	// client := &http.Client{Timeout: 30 * time.Second}
	// err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !d.IsDir() && isAudio(path) { // Only files (-type f)
	// 		fmt.Println("---------------------------")
	// 		fmt.Println(path)
	// 		readFileMeta(path)
	// 		album := filepath.Base(path)
	// 		ctx := context.Background()
	// 		recordService := services.NewMusicBrainzRecordServirce(client, 5)
	// 		buff, err := recordService.Query(ctx, album)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		fmt.Printf(string(buff))
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }
}
