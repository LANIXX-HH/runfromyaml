package config

import (
	"flag"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := New()

	// Test default values
	if cfg.Debug != false {
		t.Errorf("Expected Debug to be false, got %v", cfg.Debug)
	}

	if cfg.Rest != false {
		t.Errorf("Expected Rest to be false, got %v", cfg.Rest)
	}

	if cfg.File != "commands.yaml" {
		t.Errorf("Expected File to be 'commands.yaml', got %s", cfg.File)
	}

	if cfg.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got %s", cfg.Host)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}

	if cfg.User != "rest" {
		t.Errorf("Expected User to be 'rest', got %s", cfg.User)
	}

	if cfg.AIModel != "text-davinci-003" {
		t.Errorf("Expected AIModel to be 'text-davinci-003', got %s", cfg.AIModel)
	}

	if cfg.AICmdType != "shell" {
		t.Errorf("Expected AICmdType to be 'shell', got %s", cfg.AICmdType)
	}

	if cfg.ShellType != "bash" {
		t.Errorf("Expected ShellType to be 'bash', got %s", cfg.ShellType)
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
		wantErr  bool
	}{
		{
			name: "default values",
			args: []string{},
			expected: Config{
				Debug:     false,
				Rest:      false,
				NoAuth:    false,
				RestOut:   false,
				NoFile:    false,
				AI:        false,
				Shell:     false,
				File:      "commands.yaml",
				Host:      "localhost",
				User:      "rest",
				AIModel:   "text-davinci-003",
				AICmdType: "shell",
				ShellType: "bash",
				Port:      8080,
			},
		},
		{
			name: "debug enabled",
			args: []string{"-debug"},
			expected: Config{
				Debug:     true,
				Rest:      false,
				File:      "commands.yaml",
				Host:      "localhost",
				User:      "rest",
				AIModel:   "text-davinci-003",
				AICmdType: "shell",
				ShellType: "bash",
				Port:      8080,
			},
		},
		{
			name: "custom file",
			args: []string{"-file", "test.yaml"},
			expected: Config{
				Debug:     false,
				Rest:      false,
				File:      "test.yaml",
				Host:      "localhost",
				User:      "rest",
				AIModel:   "text-davinci-003",
				AICmdType: "shell",
				ShellType: "bash",
				Port:      8080,
			},
		},
		{
			name: "rest mode with custom port",
			args: []string{"-rest", "-port", "9090"},
			expected: Config{
				Debug:     false,
				Rest:      true,
				File:      "commands.yaml",
				Host:      "localhost",
				User:      "rest",
				AIModel:   "text-davinci-003",
				AICmdType: "shell",
				ShellType: "bash",
				Port:      9090,
			},
		},
		{
			name: "ai mode with custom model",
			args: []string{"-ai", "-ai-model", "gpt-4", "-ai-key", "test-key"},
			expected: Config{
				Debug:     false,
				Rest:      false,
				AI:        true,
				File:      "commands.yaml",
				Host:      "localhost",
				User:      "rest",
				AIKey:     "test-key",
				AIModel:   "gpt-4",
				AICmdType: "shell",
				ShellType: "bash",
				Port:      8080,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag package for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Save original args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Set test args
			os.Args = append([]string{"test"}, tt.args...)

			cfg := New()
			err := cfg.ParseFlags()

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Compare relevant fields
				if cfg.Debug != tt.expected.Debug {
					t.Errorf("Debug = %v, want %v", cfg.Debug, tt.expected.Debug)
				}
				if cfg.Rest != tt.expected.Rest {
					t.Errorf("Rest = %v, want %v", cfg.Rest, tt.expected.Rest)
				}
				if cfg.AI != tt.expected.AI {
					t.Errorf("AI = %v, want %v", cfg.AI, tt.expected.AI)
				}
				if cfg.File != tt.expected.File {
					t.Errorf("File = %v, want %v", cfg.File, tt.expected.File)
				}
				if cfg.Port != tt.expected.Port {
					t.Errorf("Port = %v, want %v", cfg.Port, tt.expected.Port)
				}
				if cfg.AIModel != tt.expected.AIModel {
					t.Errorf("AIModel = %v, want %v", cfg.AIModel, tt.expected.AIModel)
				}
				if cfg.AIKey != tt.expected.AIKey {
					t.Errorf("AIKey = %v, want %v", cfg.AIKey, tt.expected.AIKey)
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid default config",
			config: Config{
				File: "commands.yaml",
				Host: "localhost",
				Port: 8080,
				User: "rest",
			},
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			config: Config{
				File: "commands.yaml",
				Host: "localhost",
				Port: 0,
				User: "rest",
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: Config{
				File: "commands.yaml",
				Host: "localhost",
				Port: 70000,
				User: "rest",
			},
			wantErr: true,
		},
		{
			name: "empty file when not disabled",
			config: Config{
				File:   "",
				Host:   "localhost",
				Port:   8080,
				User:   "rest",
				NoFile: false,
			},
			wantErr: true,
		},
		{
			name: "empty file when disabled",
			config: Config{
				File:   "",
				Host:   "localhost",
				Port:   8080,
				User:   "rest",
				NoFile: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to validate config (you might need to implement this in your config package)
func validateConfig(cfg *Config) error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return &ConfigError{Message: "Invalid port number"}
	}

	if !cfg.NoFile && cfg.File == "" {
		return &ConfigError{Message: "File cannot be empty when not disabled"}
	}

	return nil
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
