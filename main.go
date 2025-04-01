package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var apiKey string  // Store the API key globally
var baseURL string // Store the base API URL globally

func main() {
	var rootCmd = &cobra.Command{
		Use:   "nomi-cli",
		Short: "A CLI client for the Nomi.ai API",
		Long:  `nomi-cli is a command-line client to interact with the Nomi.ai API`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load the API key from the environment variable if not provided as a flag
			if apiKey == "" {
				apiKey = os.Getenv("bb58e912-527d-40af-b585-00bfd684c603")
			}

			// Ensure an API key is available
			if apiKey == "" {
				return fmt.Errorf("API key not found. Please set the NOMI_API_KEY environment variable or use the -k flag")
			}
			// Load the base API URL from the environment variable
			baseURL = os.Getenv("NOMI_API_URL")
			if baseURL == "" {
				baseURL = "https://api.nomi.ai/v1" // Default value if environment variable is not set
			}
			return nil
		},
	}

	// Allow overriding the API key via a flag
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "API key for Nomi.ai (overrides NOMI_API_KEY)")

	// Add commands
	rootCmd.AddCommand(listNomisCmd)
	rootCmd.AddCommand(getNomiCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(listRoomsCmd)
	rootCmd.AddCommand(versionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
