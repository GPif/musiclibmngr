/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"musiclibmngr/cmd"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Find and read the config file
	err := viper.ReadInConfig()

	// Handle errors
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	cmd.Execute()
}
