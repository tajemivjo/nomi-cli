package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetNomiCommand(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Error("Expected Authorization header with API key")
		}

		// Test different scenarios based on the Nomi ID
		switch r.URL.Path {
		case "/nomis/test-uuid":
			// Return a successful response
			nomi := Nomi{
				UUID:             "test-uuid",
				Name:             "Test Nomi",
				Gender:           "female",
				Created:          "2024-01-01T12:00:00Z",
				RelationshipType: "Friend",
			}
			json.NewEncoder(w).Encode(nomi)

		case "/nomis/invalid-uuid":
			// Return a 404 error
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Nomi not found",
			})

		default:
			// Return a 500 error for unknown paths
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// Set up test environment
	originalBaseURL := baseURL
	originalAPIKey := apiKey
	defer func() {
		baseURL = originalBaseURL
		apiKey = originalAPIKey
	}()

	baseURL = server.URL
	apiKey = "test-api-key"

	// Test cases
	tests := []struct {
		name           string
		nomiID         string
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Valid Nomi ID",
			nomiID:         "test-uuid",
			expectedError:  false,
			expectedOutput: "Nomi Details:\n- ID: test-uuid\n- Name: Test Nomi\n- Gender: female\n- Created: 2024-01-01T12:00:00Z\n- Relationship Type: Friend\n",
		},
		{
			name:           "Invalid Nomi ID",
			nomiID:         "invalid-uuid",
			expectedError:  true,
			expectedOutput: "Error: 404 Not Found\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute command
			getNomiCmd.Run(getNomiCmd, []string{tc.nomiID})

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if output != tc.expectedOutput {
				t.Errorf("Expected output:\n%s\nGot:\n%s", tc.expectedOutput, output)
			}
		})
	}
}

func TestGetNomiValidation(t *testing.T) {
	// Test argument validation
	cmd := getNomiCmd

	tests := []struct {
		name          string
		args          []string
		expectedError bool
	}{
		{
			name:          "No arguments",
			args:          []string{},
			expectedError: true,
		},
		{
			name:          "Too many arguments",
			args:          []string{"id1", "id2"},
			expectedError: true,
		},
		{
			name:          "Correct number of arguments",
			args:          []string{"id1"},
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.Args(cmd, tc.args)
			if tc.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGetNomiRequestErrors(t *testing.T) {
	// Save original values
	originalAPIKey := apiKey
	originalBaseURL := baseURL
	defer func() {
		apiKey = originalAPIKey
		baseURL = originalBaseURL
	}()

	tests := []struct {
		name          string
		setupFunc     func()
		expectedError string
	}{
		{
			name: "Missing API Key",
			setupFunc: func() {
				apiKey = ""
				baseURL = "https://api.example.com"
			},
			expectedError: "Error making request",
		},
		{
			name: "Invalid Base URL",
			setupFunc: func() {
				apiKey = "test-key"
				baseURL = "://invalid-url"
			},
			expectedError: "Error creating request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			tc.setupFunc()

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() {
				w.Close()
				os.Stdout = oldStdout
			}()

			// Create a new command instance for this test
			testCmd := &cobra.Command{
				Use:   "get-nomi [id]",
				Short: "Get details of a specific Nomi",
				Args:  cobra.ExactArgs(1),
				Run:   getNomiCmd.Run,
			}

			// Execute command
			testCmd.Run(testCmd, []string{"test-id"})

			// Read captured output
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if !strings.Contains(output, tc.expectedError) {
				t.Errorf("Expected error containing '%s', got: %s", tc.expectedError, output)
			}
		})
	}
}
