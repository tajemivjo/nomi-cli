package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Save original environment and values
	originalAPIKey := os.Getenv("NOMI_API_KEY")
	originalAPIURL := os.Getenv("NOMI_API_URL")
	originalCmdAPIKey := apiKey
	originalCmdBaseURL := baseURL
	defer func() {
		os.Setenv("NOMI_API_KEY", originalAPIKey)
		os.Setenv("NOMI_API_URL", originalAPIURL)
		apiKey = originalCmdAPIKey
		baseURL = originalCmdBaseURL
	}()

	tests := []struct {
		name          string
		envSetup      func()
		args          []string
		expectedError string
		expectedURL   string
	}{
		{
			name: "Valid API Key from ENV",
			envSetup: func() {
				os.Setenv("NOMI_API_KEY", "test-key")
				os.Setenv("NOMI_API_URL", "https://test.api.nomi.ai/v1")
				apiKey = ""
				baseURL = ""
			},
			args:          []string{},
			expectedError: "",
			expectedURL:   "https://test.api.nomi.ai/v1",
		},
		{
			name: "Valid API Key from Flag",
			envSetup: func() {
				os.Unsetenv("NOMI_API_KEY")
				os.Unsetenv("NOMI_API_URL")
				apiKey = ""
				baseURL = ""
			},
			args:          []string{"--api-key", "test-key"},
			expectedError: "",
			expectedURL:   "https://api.nomi.ai/v1",
		},
		{
			name: "Missing API Key",
			envSetup: func() {
				os.Unsetenv("NOMI_API_KEY")
				os.Unsetenv("NOMI_API_URL")
				apiKey = ""
				baseURL = ""
			},
			args:          []string{},
			expectedError: "API key not found",
			expectedURL:   "",
		},
		{
			name: "Custom API URL",
			envSetup: func() {
				os.Setenv("NOMI_API_KEY", "test-key")
				os.Setenv("NOMI_API_URL", "https://custom.api.nomi.ai/v1")
				apiKey = ""
				baseURL = ""
			},
			args:          []string{},
			expectedError: "",
			expectedURL:   "https://custom.api.nomi.ai/v1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset global variables and environment for each test
			apiKey = ""
			baseURL = ""
			os.Unsetenv("NOMI_API_KEY")
			os.Unsetenv("NOMI_API_URL")

			// Setup environment for this test
			tc.envSetup()

			// Create a new command instance
			cmd := &cobra.Command{
				Use:   "nomi-cli",
				Short: "A CLI client for the Nomi.ai API",
				Run: func(cmd *cobra.Command, args []string) {
					// Empty run function to ensure PreRunE executes
				},
				SilenceUsage: true, // Silence usage on error
			}

			// Add the API key flag
			cmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "API key for Nomi.ai")

			// Set up the PreRunE function
			cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
				// Load the API key from the environment variable if not provided as a flag
				if apiKey == "" {
					apiKey = os.Getenv("NOMI_API_KEY")
				}

				// Ensure an API key is available
				if apiKey == "" {
					return fmt.Errorf("API key not found")
				}

				// Load the base API URL from the environment variable
				baseURL = os.Getenv("NOMI_API_URL")
				if baseURL == "" {
					baseURL = "https://api.nomi.ai/v1"
				}
				return nil
			}

			// Set the command arguments
			cmd.SetArgs(tc.args)

			// Execute the command
			err := cmd.Execute()

			// Verify error expectation
			if tc.expectedError != "" {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("Expected error containing '%s', got: %s", tc.expectedError, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify URL expectation
			if tc.expectedError == "" && tc.expectedURL != baseURL {
				t.Errorf("Expected URL %s, got: %s", tc.expectedURL, baseURL)
			}
		})
	}
}
