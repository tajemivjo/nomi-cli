package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
	baseURL = server.URL
	apiKey = "test-api-key"

	// Test cases
	tests := []struct {
		name          string
		nomiID        string
		expectedError bool
	}{
		{
			name:          "Valid Nomi ID",
			nomiID:        "test-uuid",
			expectedError: false,
		},
		{
			name:          "Invalid Nomi ID",
			nomiID:        "invalid-uuid",
			expectedError: true,
		},
		{
			name:          "Empty Nomi ID",
			nomiID:        "",
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.nomiID != "" {
				client := &http.Client{}
				url := baseURL + "/nomis/" + tc.nomiID

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				req.Header.Set("Authorization", "Bearer "+apiKey)
				resp, err := client.Do(req)

				if tc.expectedError {
					if err == nil && resp.StatusCode == http.StatusOK {
						t.Error("Expected error but got success")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if resp.StatusCode != http.StatusOK {
						t.Errorf("Expected status OK, got %v", resp.Status)
					}

					var nomi Nomi
					if err := json.NewDecoder(resp.Body).Decode(&nomi); err != nil {
						t.Errorf("Failed to decode response: %v", err)
					}

					// Verify the response data
					if nomi.UUID != tc.nomiID {
						t.Errorf("Expected Nomi UUID %s, got %s", tc.nomiID, nomi.UUID)
					}
				}
			}
		})
	}
}

func TestGetNomiValidation(t *testing.T) {
	// Test argument validation
	cmd := getNomiCmd

	// Test with no arguments
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error for no arguments, got none")
	}

	// Test with too many arguments
	err = cmd.Args(cmd, []string{"id1", "id2"})
	if err == nil {
		t.Error("Expected error for too many arguments, got none")
	}

	// Test with correct number of arguments
	err = cmd.Args(cmd, []string{"id1"})
	if err != nil {
		t.Errorf("Expected no error for single argument, got %v", err)
	}
}
