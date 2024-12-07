package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version will be set during build time
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of nomi-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nomi-cli version %s\n", Version)
	},
}
