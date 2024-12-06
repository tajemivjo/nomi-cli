package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var fullOutput bool // Flag to control output verbosity

var listNomisCmd = &cobra.Command{
	Use:   "list-nomis",
	Short: "List all Nomis",
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/nomis", baseURL), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: %s\n", resp.Status)
			return
		}

		var result NomiResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		// Display the Nomis
		for _, nomi := range result.Nomis {
			if fullOutput {
				// Full output
				fmt.Printf("- ID: %s\n  Name: %s\n  Gender: %s\n  Created: %s\n  Relationship: %s\n\n",
					nomi.UUID, nomi.Name, nomi.Gender, nomi.Created, nomi.RelationshipType)
			} else {
				// Default output (Name and Relationship only)
				fmt.Printf("%s (%s)\n", nomi.Name, nomi.RelationshipType)
			}
		}
	},
}

func init() {
	// Add the --full flag to the list-nomis command
	listNomisCmd.Flags().BoolVarP(&fullOutput, "full", "f", false, "Display full details of each Nomi")
}
