package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func displayRoom(room Room) {
	name := room.Name
	if name == "" {
		name = "<empty>"
	}

	fmt.Printf("Room: %s\n", name)
	fmt.Printf("- UUID: %s\n", room.UUID)
	fmt.Printf("- Created: %s\n", room.Created)
	fmt.Printf("- Updated: %s\n", room.Updated)
	fmt.Printf("- Status: %s\n", room.Status)
	fmt.Printf("- Backchanneling: %v\n", room.BackchannelingEnabled)

	if room.Note != "" {
		fmt.Printf("- Note: %s\n", room.Note)
	}

	if len(room.Nomis) > 0 {
		fmt.Println("- Nomis:")
		for _, nomi := range room.Nomis {
			fmt.Printf("  • %s (%s, %s)\n",
				nomi.Name,
				nomi.Gender,
				nomi.RelationshipType)
		}
	}
}

var listRoomsCmd = &cobra.Command{
	Use:   "list-rooms",
	Short: "List all rooms",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{}
		url := baseURL + "/rooms"

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
		var response RoomResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		// Print the Rooms
		fmt.Printf("Total Rooms: %d\n\n", len(response.Rooms))
		for _, room := range response.Rooms {
			displayRoom(room)
			fmt.Println()
		}
	},
}
