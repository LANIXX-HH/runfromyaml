package config

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		expected Config
		wantErr  bool
	}{
		{
			name: "valid YAML with boolean options",
			yamlData: `
options:
  - key: "debug"
    value: true
  - key: "rest"
    value: false
  - key: "no-auth"
    value: true
`,
			expected: Config{
				Debug:  true,
				Rest:   false,
				NoAuth: true,
			},
			wantErr: false,
		},
		{
			name: "valid YAML with string options",
			yamlData: `
options:
  - key: "file"
    value: "test.yaml"
  - key: "host"
    value: "0.0.0.0"
  - key: "user"
    value: "admin"
  - key: "ai-key"
    value: "sk-test123"
  - key: "ai-model"
    value: "gpt-4"
`,
			expected: Config{
				File:    "test.yaml",
				Host:    "0.0.0.0",
				User:    "admin",
				AIKey:   "sk-test123",
				AIModel: "gpt-4",
			},
			wantErr: false,
		},
		{
			name: "valid YAML with integer options",
			yamlData: `
options:
  - key: "port"
    value: 9090
`,
			expected: Config{
				Port: 9090,
			},
			wantErr: false,
		},
		{
			name: "mixed options",
			yamlData: `
options:
  - key: "debug"
    value: true
  - key: "file"
    value: "custom.yaml"
  - key: "port"
    value: 8888
  - key: "ai"
    value: true
  - key: "ai-model"
    value: "gpt-3.5-turbo"
`,
			expected: Config{
				Debug:   true,
				File:    "custom.yaml",
				Port:    8888,
				AI:      true,
				AIModel: "gpt-3.5-turbo",
			},
			wantErr: false,
		},
		{
			name: "empty options",
			yamlData: `
options: []
`,
			expected: Config{},
			wantErr:  false,
		},
		{
			name: "invalid YAML",
			yamlData: `
options:
  - key: "debug"
    value: true
  - invalid_structure
`,
			expected: Config{},
			wantErr:  true,
		},
		{
			name: "unknown option key (should be ignored)",
			yamlData: `
options:
  - key: "debug"
    value: true
  - key: "unknown-option"
    value: "should-be-ignored"
`,
			expected: Config{
				Debug: true,
			},
			wantErr: false,
		},
		{
			name: "wrong type for boolean option (should be ignored)",
			yamlData: `
options:
  - key: "debug"
    value: "not-a-boolean"
  - key: "rest"
    value: true
`,
			expected: Config{
				Rest: true,
			},
			wantErr: false,
		},
		{
			name: "wrong type for string option (should be ignored)",
			yamlData: `
options:
  - key: "file"
    value: 123
  - key: "host"
    value: "localhost"
`,
			expected: Config{
				Host: "localhost",
			},
			wantErr: false,
		},
		{
			name: "wrong type for integer option (should be ignored)",
			yamlData: `
options:
  - key: "port"
    value: "not-a-number"
  - key: "debug"
    value: true
`,
			expected: Config{
				Debug: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := cfg.LoadFromYAML([]byte(tt.yamlData))

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Compare all relevant fields
				if cfg.Debug != tt.expected.Debug {
					t.Errorf("Debug = %v, want %v", cfg.Debug, tt.expected.Debug)
				}
				if cfg.Rest != tt.expected.Rest {
					t.Errorf("Rest = %v, want %v", cfg.Rest, tt.expected.Rest)
				}
				if cfg.NoAuth != tt.expected.NoAuth {
					t.Errorf("NoAuth = %v, want %v", cfg.NoAuth, tt.expected.NoAuth)
				}
				if cfg.RestOut != tt.expected.RestOut {
					t.Errorf("RestOut = %v, want %v", cfg.RestOut, tt.expected.RestOut)
				}
				if cfg.NoFile != tt.expected.NoFile {
					t.Errorf("NoFile = %v, want %v", cfg.NoFile, tt.expected.NoFile)
				}
				if cfg.AI != tt.expected.AI {
					t.Errorf("AI = %v, want %v", cfg.AI, tt.expected.AI)
				}
				if cfg.File != tt.expected.File {
					t.Errorf("File = %v, want %v", cfg.File, tt.expected.File)
				}
				if cfg.Host != tt.expected.Host {
					t.Errorf("Host = %v, want %v", cfg.Host, tt.expected.Host)
				}
				if cfg.User != tt.expected.User {
					t.Errorf("User = %v, want %v", cfg.User, tt.expected.User)
				}
				if cfg.AIKey != tt.expected.AIKey {
					t.Errorf("AIKey = %v, want %v", cfg.AIKey, tt.expected.AIKey)
				}
				if cfg.AIModel != tt.expected.AIModel {
					t.Errorf("AIModel = %v, want %v", cfg.AIModel, tt.expected.AIModel)
				}
				if cfg.Port != tt.expected.Port {
					t.Errorf("Port = %v, want %v", cfg.Port, tt.expected.Port)
				}
			}
		})
	}
}

func TestYAMLOptionsUnmarshal(t *testing.T) {
	yamlData := `
options:
  - key: "debug"
    value: true
  - key: "file"
    value: "test.yaml"
  - key: "port"
    value: 8080
`

	var yamlOpts YAMLOptions
	err := yaml.Unmarshal([]byte(yamlData), &yamlOpts)

	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(yamlOpts.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(yamlOpts.Options))
	}

	// Check first option
	if yamlOpts.Options[0].Key != "debug" {
		t.Errorf("Expected first option key to be 'debug', got %s", yamlOpts.Options[0].Key)
	}

	if yamlOpts.Options[0].Value != true {
		t.Errorf("Expected first option value to be true, got %v", yamlOpts.Options[0].Value)
	}

	// Check second option
	if yamlOpts.Options[1].Key != "file" {
		t.Errorf("Expected second option key to be 'file', got %s", yamlOpts.Options[1].Key)
	}

	if yamlOpts.Options[1].Value != "test.yaml" {
		t.Errorf("Expected second option value to be 'test.yaml', got %v", yamlOpts.Options[1].Value)
	}

	// Check third option
	if yamlOpts.Options[2].Key != "port" {
		t.Errorf("Expected third option key to be 'port', got %s", yamlOpts.Options[2].Key)
	}

	if yamlOpts.Options[2].Value != 8080 {
		t.Errorf("Expected third option value to be 8080, got %v", yamlOpts.Options[2].Value)
	}
}
