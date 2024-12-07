package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListNomis(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Error("Expected Authorization header with API key")
		}

		// Mock response
		response := NomiResponse{
			Nomis: []Nomi{
				{
					UUID:             "test-uuid-1",
					Name:             "John",
					Gender:           "male",
					Created:          "2024-01-01T12:00:00Z",
					RelationshipType: "Friend",
				},
				{
					UUID:             "test-uuid-2",
					Name:             "Alice",
					Gender:           "female",
					Created:          "2024-01-01T12:00:00Z",
					RelationshipType: "Mentor",
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Set test environment
	baseURL = server.URL
	apiKey = "test-api-key"

	// Test list nomis request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL+"/nomis", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result NomiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if len(result.Nomis) != 2 {
		t.Errorf("Expected 2 Nomis, got %d", len(result.Nomis))
	}
}
