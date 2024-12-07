package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test API key from environment variable
	os.Setenv("NOMI_API_KEY", "test-api-key")
	os.Setenv("NOMI_API_URL", "https://test.api.nomi.ai/v1")

	// Create a test root command with the same configuration as in main.go
	rootCmd := &cobra.Command{
		Use:   "nomi-cli",
		Short: "A CLI client for the Nomi.ai API",
		Long:  `nomi-cli is a command-line client to interact with the Nomi.ai API`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load the API key from the environment variable if not provided as a flag
			if apiKey == "" {
				apiKey = os.Getenv("NOMI_API_KEY")
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

	// Test with valid API key
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err != nil {
		t.Errorf("Expected no error with valid API key, got %v", err)
	}

	// Test without API key
	os.Unsetenv("NOMI_API_KEY")
	apiKey = "" // Reset the global variable
	err = rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err == nil {
		t.Error("Expected error for missing API key, got none")
	}
}
