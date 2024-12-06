package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var getNomiCmd = &cobra.Command{
	Use:   "get-nomi [id]",
	Short: "Get details of a specific Nomi",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is passed (the Nomi ID)
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		client := &http.Client{}
		url := fmt.Sprintf("%s/nomis/%s", baseURL, id) // Use dynamic baseURL

		// Create the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the Authorization header
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// Perform the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: %s\n", resp.Status)
			return
		}

		// Parse the response body
		var nomi Nomi
		if err := json.NewDecoder(resp.Body).Decode(&nomi); err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		// Print the Nomi details
		fmt.Println("Nomi Details:")
		fmt.Printf("- ID: %s\n- Name: %s\n- Gender: %s\n- Created: %s\n- Relationship Type: %s\n",
			nomi.UUID, nomi.Name, nomi.Gender, nomi.Created, nomi.RelationshipType)
	},
}
