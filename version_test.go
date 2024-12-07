package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Since Version is a global variable, you can set it here for testing:
func init() {
	Version = "1.2.3"
}

func TestVersionCmd(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set up a root command and add versionCmd
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(versionCmd)

	// Execute the "version" command
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute version command: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)

	// Check if the output contains the expected version string
	expected := "nomi-cli version 1.2.3"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected %q in output, got %q", expected, output)
	}
}
