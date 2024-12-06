package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type ChatRequest struct {
	MessageText string `json:"messageText"`
}

type Message struct {
	UUID string `json:"uuid"`
	Text string `json:"text"`
	Sent string `json:"sent"`
}

type ChatResponse struct {
	SentMessage  Message `json:"sentMessage"`
	ReplyMessage Message `json:"replyMessage"`
}

// clearScreen clears the terminal screen based on the OS.
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// findNomiByName retrieves the UUID of a Nomi by its name.
func findNomiByName(name string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/nomis", baseURL) // Use dynamic baseURL

	// Fetch the list of Nomis
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error fetching Nomis: %s", resp.Status)
	}

	var result struct {
		Nomis []struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"nomis"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Search for the Nomi by name
	for _, nomi := range result.Nomis {
		if strings.EqualFold(nomi.Name, name) { // Case-insensitive comparison
			return nomi.UUID, nil
		}
	}

	return "", fmt.Errorf("no Nomi found with the name: %s", name)
}

// spinner displays a spinning wheel animation while waiting for a response.
func spinner(stopChan chan bool) {
	chars := []rune{'|', '/', '-', '\\'}
	for {
		select {
		case <-stopChan:
			return
		default:
			for _, char := range chars {
				select {
				case <-stopChan:
					return
				default:
					fmt.Printf("\r%c", char)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}
}

var chatCmd = &cobra.Command{
	Use:   "chat [id]",
	Short: "Start a live chat session with a specific Nomi",
	Args:  cobra.ExactArgs(1), // Requires exactly one argument: the Nomi Name
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure the screen is cleared when the program exits
		defer clearScreen()
		name := args[0]

		// Find the UUID for the given name
		nomiID, err := findNomiByName(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		client := &http.Client{}
		url := fmt.Sprintf("%s/nomis/%s/chat", baseURL, nomiID) // Use dynamic baseURL

		// Clear the terminal at the start of the chat
		clearScreen()

		fmt.Println("Chat session started with", name+". Type your message and press Enter to send.")
		fmt.Println("Type 'exit' to end the session.")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("You: ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if strings.ToLower(strings.TrimSpace(input)) == "exit" {
				fmt.Println("Chat session ended.")
				break
			}

			// Prepare the request payload
			chatRequest := ChatRequest{MessageText: input}
			requestBody, err := json.Marshal(chatRequest)
			if err != nil {
				fmt.Println("Error encoding request body:", err)
				continue
			}

			// Create the HTTP request
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Set("Authorization", "Bearer "+apiKey)
			req.Header.Set("Content-Type", "application/json")

			// Start the spinner
			stopChan := make(chan bool)
			go spinner(stopChan)

			// Send the request
			resp, err := client.Do(req)

			// Stop the spinner
			close(stopChan)
			fmt.Print("\r") // Clear the spinner line

			if err != nil {
				fmt.Println("Error sending message:", err)
				continue
			}
			defer resp.Body.Close()

			// Check for successful response
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error: %s\n", resp.Status)
				continue
			}

			// Decode the response
			var chatResponse ChatResponse
			if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
				fmt.Println("Error decoding response:", err)
				continue
			}

			// Display the reply
			fmt.Printf("Nomi: %s\n", chatResponse.ReplyMessage.Text)
		}
	},
}
