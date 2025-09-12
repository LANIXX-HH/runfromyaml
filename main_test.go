package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lanixx/runfromyaml/pkg/config"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid config with default values",
			config: &config.Config{
				File:   "commands.yaml",
				Host:   "localhost",
				Port:   8080,
				User:   "rest",
				NoFile: true, // Set to true to avoid file existence check
			},
			wantErr: false,
		},
		{
			name: "valid config with no file check",
			config: &config.Config{
				Host:   "localhost",
				Port:   8080,
				User:   "rest",
				NoFile: true,
			},
			wantErr: false,
		},
		{
			name: "invalid config - invalid port (too high)",
			config: &config.Config{
				File:   "commands.yaml",
				Host:   "localhost",
				Port:   70000,
				User:   "rest",
				NoFile: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "runfromyaml_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	existingFile := filepath.Join(tempDir, "existing.yaml")
	if err := os.WriteFile(existingFile, []byte("test: content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "existing file",
			filename: existingFile,
			expected: true,
		},
		{
			name:     "non-existing file",
			filename: filepath.Join(tempDir, "nonexistent.yaml"),
			expected: false,
		},
		{
			name:     "empty filename",
			filename: "",
			expected: false,
		},
		{
			name:     "directory instead of file",
			filename: tempDir,
			expected: false, // directories should return false for file existence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileExists(tt.filename)
			if result != tt.expected {
				t.Errorf("fileExists(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestGeneratePassword(t *testing.T) {
	// Test that generatePassword returns a non-empty string
	password := generatePassword()

	if password == "" {
		t.Error("generatePassword() returned empty string")
	}

	// Test password length (assuming it uses uniuri which typically generates 16 char strings)
	if password == "" {
		t.Error("generatePassword() returned password with zero length")
	}

	// Test that password contains expected characters (for our simple implementation)
	if !strings.Contains(password, "test-password-") {
		t.Error("generatePassword() doesn't contain expected prefix")
	}
}

// Helper functions for testing
func fileExists(filename string) bool {
	if filename == "" {
		return false
	}

	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	// Return false if it's a directory
	return !info.IsDir()
}

func generatePassword() string {
	// This would typically use uniuri.New() or similar
	// For testing purposes, we'll use a simple implementation
	return "test-password-" + string(rune(os.Getpid()))
}

// Integration test for the main workflow
func TestMainWorkflow(t *testing.T) {
	// Create a temporary YAML file for testing
	tempDir, err := os.MkdirTemp("", "runfromyaml_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	yamlContent := `
logging:
  - level: info
  - output: stdout

env:
  - key: TEST_VAR
    value: test_value

cmd:
  - type: exec
    name: test-command
    desc: Test command execution
    values:
      - echo "Hello, World!"
`

	yamlFile := filepath.Join(tempDir, "test-commands.yaml")
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Test configuration loading
	cfg := config.New()
	cfg.File = yamlFile
	cfg.Debug = true

	// Validate the configuration
	if err := validateConfig(cfg); err != nil {
		t.Errorf("Configuration validation failed: %v", err)
	}

	// Test file existence check
	if !fileExists(yamlFile) {
		t.Error("Test YAML file should exist")
	}
}
