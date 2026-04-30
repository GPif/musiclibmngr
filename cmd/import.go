/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"musiclibmngr/internal/db"
	"musiclibmngr/internal/file"
	"musiclibmngr/internal/pathmatcher"
	"os"
	"path/filepath"

	"github.com/jo-hoe/chromaprint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		relativePath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		standardTrackFormat := viper.GetString("standard_track_format")
		multiDiscTrackFormat := viper.GetString("multi_disc_track_format")
		matcher := pathmatcher.NewMatcher()

		fmt.Println("# Matcher")

		var m map[string]string
		for _, format := range []string{multiDiscTrackFormat, standardTrackFormat} {
			m, err = matcher.ExtractData(format, relativePath)
			if !errors.Is(err, pathmatcher.ErrInvalidTemplate) {
				break
			}
		}
		fmt.Printf("%v\n", m)

		fmt.Println("# Extract tag")

		tags, err := file.ExtractTag(path)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", tags)

		fmt.Println("# Fingerprint")

		cp, err := chromaprint.NewBuilder().Build()
		if err != nil {
			return err
		}

		fps, err := cp.CreateFingerprints(path)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v\n", fps)
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

		err = fileScan(filePath, dbConn)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
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
