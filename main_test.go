package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "gotask_test")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("gotask_test") // Clean up

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "help command",
			args:           []string{"help"},
			expectedOutput: "gotask - A simple task manager",
			expectError:    false,
		},
		{
			name:           "no arguments",
			args:           []string{},
			expectedOutput: "gotask - A simple task manager",
			expectError:    false,
		},
		{
			name:           "unknown command",
			args:           []string{"unknown"},
			expectedOutput: "Unknown command: unknown",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./gotask_test", tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !strings.Contains(string(output), tt.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", tt.expectedOutput, string(output))
			}
		})
	}
}
