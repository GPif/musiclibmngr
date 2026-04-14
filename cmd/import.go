/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"musiclibmngr/internal/db"
	"musiclibmngr/internal/file"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func validArgs(cmd *cobra.Command, args []string) error {
	// Check that there is one args and that it is a valid path
	if len(args) != 1 {
		cmd.Help()
		return fmt.Errorf("Error: Must provide exactly one argument (the path to the music file).\n")
	}
	filePath := args[0]
	// A simple check to see if the path exists and is a file.
	// In a real application, you might want more robust path validation (e.g., checking file extensions).
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("Error: The provided path does not exist: %s\n", filePath)
	}
	return nil
}

func parseMusicPath(path string) (artist, album, trackNum, title string) {
	re := regexp.MustCompile(`testdata/(.+)/(.+)/(\d+)\s+(.+)\.\w+$`)
	matches := re.FindStringSubmatch(path)
	if len(matches) == 5 {
		return matches[1], matches[2], matches[3], matches[4]
	}
	return "", "", "", ""
}

func fileScan(baseDir string, dbConn *db.DB) error {
	return filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Scan and import music file into db",
	Long:  `Scan and import music file into db`,
	Run: func(cmd *cobra.Command, args []string) {
		err := validArgs(cmd, args)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		filePath := args[0]
		dbPath := viper.GetString("db")
		dbConn, err := db.New(string(dbPath))
		if err != nil {
			panic(err)
		}

		fileScan(filePath, dbConn)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
