package cli

import (
	"os"
	"testing"
)

func TestSSHExpandenvOptions(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_SSH_KEY", "/path/to/test/key")
	defer os.Unsetenv("TEST_SSH_KEY")

	executor := &CommandExecutor{}

	// Test with expandenv enabled
	cmdWithExpandenv := &Command{
		Options: map[string]interface{}{
			"user":      "testuser",
			"host":      "testhost",
			"port":      22,
			"expandenv": true,
			"options": []interface{}{
				"-i $TEST_SSH_KEY",
				"-o ConnectTimeout=5",
				"-o StrictHostKeyChecking=no",
			},
		},
	}

	args := executor.buildSSHArgs(cmdWithExpandenv)

	// Check if the environment variable was expanded in options
	found := false
	for _, arg := range args {
		if arg == "-i /path/to/test/key" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected SSH option with expanded environment variable '-i /path/to/test/key', but got: %v", args)
	}

	// Test with expandenv disabled
	cmdWithoutExpandenv := &Command{
		Options: map[string]interface{}{
			"user":      "testuser",
			"host":      "testhost",
			"port":      22,
			"expandenv": false,
			"options": []interface{}{
				"-i $TEST_SSH_KEY",
				"-o ConnectTimeout=5",
			},
		},
	}

	argsNoExpand := executor.buildSSHArgs(cmdWithoutExpandenv)

	// Check if the environment variable was NOT expanded
	foundLiteral := false
	for _, arg := range argsNoExpand {
		if arg == "-i $TEST_SSH_KEY" {
			foundLiteral = true
			break
		}
	}

	if !foundLiteral {
		t.Errorf("Expected SSH option with literal environment variable '-i $TEST_SSH_KEY', but got: %v", argsNoExpand)
	}
}

func TestSSHExpandenvUserHost(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_USER", "expandeduser")
	os.Setenv("TEST_HOST", "expandedhost")
	defer func() {
		os.Unsetenv("TEST_USER")
		os.Unsetenv("TEST_HOST")
	}()

	executor := &CommandExecutor{}

	// Test with expandenv enabled
	cmd := &Command{
		Options: map[string]interface{}{
			"user":      "$TEST_USER",
			"host":      "$TEST_HOST",
			"port":      22,
			"expandenv": true,
		},
	}

	args := executor.buildSSHArgs(cmd)

	// Check if user and host were expanded
	expectedArgs := []string{"ssh", "-p", "22", "-l", "expandeduser", "expandedhost"}

	if len(args) < len(expectedArgs) {
		t.Errorf("Expected at least %d args, got %d: %v", len(expectedArgs), len(args), args)
		return
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected arg[%d] to be '%s', got '%s'", i, expected, args[i])
		}
	}
}
