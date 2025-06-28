package functions

import (
	"testing"
)

func TestFilterRelevantEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "Filter system variables",
			input: map[string]string{
				"HOME":     "/home/user",
				"PATH":     "/usr/bin:/bin",
				"USER":     "testuser",
				"SHELL":    "/bin/bash",
				"PWD":      "/current/dir",
				"TMPDIR":   "/tmp",
				"AWS_REGION": "eu-central-1",
				"DOCKER_HOST": "tcp://localhost:2376",
			},
			expected: map[string]string{
				"AWS_REGION":  "eu-central-1",
				"DOCKER_HOST": "tcp://localhost:2376",
			},
		},
		{
			name: "Keep relevant development variables",
			input: map[string]string{
				"NODE_ENV":     "production",
				"API_KEY":      "secret123",
				"DATABASE_URL": "postgres://localhost:5432/db",
				"DEBUG":        "true",
				"PORT":         "3000",
				"TERM":         "xterm-256color",
				"EDITOR":       "vim",
			},
			expected: map[string]string{
				"NODE_ENV":     "production",
				"API_KEY":      "secret123",
				"DATABASE_URL": "postgres://localhost:5432/db",
				"DEBUG":        "true",
				"PORT":         "3000",
			},
		},
		{
			name: "Filter shell and terminal variables",
			input: map[string]string{
				"BASH_VERSION":    "5.1.0",
				"ZSH_VERSION":     "5.8",
				"PS1":             "$ ",
				"HISTFILE":        "/home/user/.bash_history",
				"COLORTERM":       "truecolor",
				"KUBERNETES_SERVICE_HOST": "10.0.0.1",
				"CI_COMMIT_SHA":   "abc123",
			},
			expected: map[string]string{
				"KUBERNETES_SERVICE_HOST": "10.0.0.1",
				"CI_COMMIT_SHA":           "abc123",
			},
		},
		{
			name: "Empty input",
			input: map[string]string{},
			expected: map[string]string{},
		},
		{
			name: "Only system variables",
			input: map[string]string{
				"HOME":   "/home/user",
				"PATH":   "/usr/bin:/bin",
				"SHELL":  "/bin/bash",
				"TMPDIR": "/tmp",
			},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterRelevantEnvVars(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d variables, got %d", len(tt.expected), len(result))
			}
			
			for key, expectedValue := range tt.expected {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("Expected key %s not found in result", key)
				} else if actualValue != expectedValue {
					t.Errorf("Expected %s=%s, got %s=%s", key, expectedValue, key, actualValue)
				}
			}
			
			// Check that no unexpected keys are present
			for key := range result {
				if _, expected := tt.expected[key]; !expected {
					t.Errorf("Unexpected key %s found in result", key)
				}
			}
		})
	}
}

func TestIsRelevantEnvVar(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		expected bool
	}{
		// AWS variables
		{"AWS_REGION", "AWS_REGION", true},
		{"AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID", true},
		{"AWS_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", true},
		
		// Docker variables
		{"DOCKER_HOST", "DOCKER_HOST", true},
		{"DOCKER_TLS_VERIFY", "DOCKER_TLS_VERIFY", true},
		
		// Kubernetes variables
		{"KUBECONFIG", "KUBECONFIG", true},
		{"K8S_NAMESPACE", "K8S_NAMESPACE", true},
		
		// CI/CD variables
		{"CI_COMMIT_SHA", "CI_COMMIT_SHA", true},
		{"BUILD_NUMBER", "BUILD_NUMBER", true},
		{"DEPLOY_ENV", "DEPLOY_ENV", true},
		
		// Database variables
		{"DATABASE_URL", "DATABASE_URL", true},
		{"REDIS_URL", "REDIS_URL", true},
		{"MYSQL_HOST", "MYSQL_HOST", true},
		
		// Application variables
		{"API_KEY", "API_KEY", true},
		{"PORT", "PORT", true},
		{"DEBUG", "DEBUG", true},
		{"ENVIRONMENT", "ENVIRONMENT", true},
		
		// Programming language variables
		{"NODE_ENV", "NODE_ENV", true},
		{"PYTHON_PATH", "PYTHON_PATH", true},
		{"JAVA_HOME", "JAVA_HOME", true},
		{"GO_PATH", "GO_PATH", true},
		
		// Git variables
		{"GIT_BRANCH", "GIT_BRANCH", true},
		{"GITHUB_TOKEN", "GITHUB_TOKEN", true},
		
		// Non-relevant variables should return false
		{"HOME", "HOME", false},
		{"PATH", "PATH", false},
		{"USER", "USER", false},
		{"SHELL", "SHELL", false},
		{"RANDOM_VAR", "RANDOM_VAR", false},
		{"CUSTOM_TOOL", "CUSTOM_TOOL", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRelevantEnvVar(tt.envVar)
			if result != tt.expected {
				t.Errorf("isRelevantEnvVar(%s) = %v, expected %v", tt.envVar, result, tt.expected)
			}
		})
	}
}

func TestPrintShellCommandsAsYamlWithFiltering(t *testing.T) {
	commands := []string{"ls -la", "pwd", "echo $AWS_REGION"}
	envs := map[string]string{
		"HOME":       "/home/user",
		"PATH":       "/usr/bin:/bin",
		"AWS_REGION": "eu-central-1",
		"API_KEY":    "secret123",
		"SHELL":      "/bin/bash",
		"DEBUG":      "true",
	}

	result := PrintShellCommandsAsYaml(commands, envs)

	// Check that commands are present
	cmdSection, exists := result["cmd"]
	if !exists {
		t.Error("Expected 'cmd' section in result")
	}

	cmdList, ok := cmdSection.([]map[string]interface{})
	if !ok {
		t.Error("Expected 'cmd' to be a slice of maps")
	}

	if len(cmdList) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(cmdList))
	}

	// Check that only relevant env vars are present
	envSection, exists := result["env"]
	if !exists {
		t.Error("Expected 'env' section in result")
	}

	envList, ok := envSection.([]map[string]interface{})
	if !ok {
		t.Error("Expected 'env' to be a slice of maps")
	}

	// Should only have AWS_REGION, API_KEY, and DEBUG (filtered)
	expectedEnvCount := 3
	if len(envList) != expectedEnvCount {
		t.Errorf("Expected %d environment variables, got %d", expectedEnvCount, len(envList))
	}

	// Check that system variables are filtered out
	foundSystemVar := false
	for _, env := range envList {
		key, _ := env["key"].(string)
		if key == "HOME" || key == "PATH" || key == "SHELL" {
			foundSystemVar = true
			break
		}
	}
	if foundSystemVar {
		t.Error("System variables should be filtered out")
	}

	// Check that relevant variables are present
	relevantVars := map[string]bool{
		"AWS_REGION": false,
		"API_KEY":    false,
		"DEBUG":      false,
	}

	for _, env := range envList {
		key, _ := env["key"].(string)
		if _, exists := relevantVars[key]; exists {
			relevantVars[key] = true
		}
	}

	for key, found := range relevantVars {
		if !found {
			t.Errorf("Expected relevant variable %s not found", key)
		}
	}
}
