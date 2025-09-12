package cli

import (
	"os"
	"testing"
)

func TestCommandType_String(t *testing.T) {
	tests := []struct {
		name string
		ct   CommandType
		want string
	}{
		{"exec", CommandTypeExec, "exec"},
		{"shell", CommandTypeShell, "shell"},
		{"docker", CommandTypeDocker, "docker"},
		{"docker-compose", CommandTypeDockerCompose, "docker-compose"},
		{"ssh", CommandTypeSSH, "ssh"},
		{"config", CommandTypeConfig, "conf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.ct); got != tt.want {
				t.Errorf("CommandType = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputType_String(t *testing.T) {
	tests := []struct {
		name string
		ot   OutputType
		want string
	}{
		{"rest", OutputTypeRest, "rest"},
		{"file", OutputTypeFile, "file"},
		{"stdout", OutputTypeStdout, "stdout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.ot); got != tt.want {
				t.Errorf("OutputType = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name string
		ll   LogLevel
		want string
	}{
		{"info", LogLevelInfo, "info"},
		{"error", LogLevelError, "error"},
		{"debug", LogLevelDebug, "debug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.ll); got != tt.want {
				t.Errorf("LogLevel = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()

	if env == nil {
		t.Error("NewEnvironment() returned nil")
		return // Early return to avoid nil pointer dereference
	}

	if env.variables == nil {
		t.Error("Environment variables map is nil")
	}

	if env.shell == nil {
		t.Error("Environment shell slice is nil")
	}
}

func TestEnvironment_Set(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"simple variable", "TEST_VAR", "test_value"},
		{"empty value", "EMPTY_VAR", ""},
		{"special characters", "SPECIAL_VAR", "value with spaces & symbols!"},
		{"numeric value", "NUMERIC_VAR", "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env var if it exists
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue != "" {
					_ = os.Setenv(tt.key, originalValue)
				} else {
					_ = os.Unsetenv(tt.key)
				}
			}()

			env.Set(tt.key, tt.value)

			if got := env.variables[tt.key]; got != tt.value {
				t.Errorf("Environment.Set() internal map = %v, want %v", got, tt.value)
			}

			if got := os.Getenv(tt.key); got != tt.value {
				t.Errorf("Environment.Set() OS env = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestEnvironment_Get(t *testing.T) {
	env := NewEnvironment()
	env.Set("TEST_VAR", "test_value")
	env.Set("EMPTY_VAR", "")

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"existing variable", "TEST_VAR", "test_value"},
		{"empty variable", "EMPTY_VAR", ""},
		{"non-existing variable", "NON_EXISTENT", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := env.Get(tt.key)

			if got != tt.expected {
				t.Errorf("Environment.Get() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnvironment_GetVariables(t *testing.T) {
	env := NewEnvironment()
	env.Set("VAR1", "value1")
	env.Set("VAR2", "value2")

	vars := env.GetVariables()

	if len(vars) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(vars))
	}

	if vars["VAR1"] != "value1" {
		t.Errorf("Expected VAR1=value1, got %s", vars["VAR1"])
	}

	if vars["VAR2"] != "value2" {
		t.Errorf("Expected VAR2=value2, got %s", vars["VAR2"])
	}
}

func TestEnvironment_Shell(t *testing.T) {
	env := NewEnvironment()
	env.Set("VAR1", "value1")
	env.Set("VAR2", "value2")

	shell := env.Shell()

	if len(shell) != 2 {
		t.Errorf("Expected 2 shell variables, got %d", len(shell))
	}

	// Check that shell format is correct
	found1, found2 := false, false
	for _, s := range shell {
		if s == "VAR1=value1" {
			found1 = true
		}
		if s == "VAR2=value2" {
			found2 = true
		}
	}

	if !found1 {
		t.Error("VAR1=value1 not found in shell format")
	}

	if !found2 {
		t.Error("VAR2=value2 not found in shell format")
	}
}

func TestCommand_Creation(t *testing.T) {
	env := NewEnvironment()

	cmd := Command{
		Type:        CommandTypeExec,
		Description: "Test command",
		Values:      []string{"echo", "hello"},
		Options:     make(map[string]interface{}),
		Env:         env,
	}

	if cmd.Type != CommandTypeExec {
		t.Errorf("Expected command type exec, got %s", cmd.Type)
	}

	if cmd.Description != "Test command" {
		t.Errorf("Expected description 'Test command', got %s", cmd.Description)
	}

	if len(cmd.Values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(cmd.Values))
	}

	if cmd.Env == nil {
		t.Error("Expected environment to be set")
	}
}
