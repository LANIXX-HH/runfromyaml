package errors

import (
	"errors"
	"testing"
)

func TestRunFromYAMLError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *RunFromYAMLError
		expected string
	}{
		{
			name: "simple error",
			err: &RunFromYAMLError{
				Type:    ErrorTypeConfig,
				Message: "configuration error",
			},
			expected: "[CONFIG] configuration error",
		},
		{
			name: "error with cause",
			err: &RunFromYAMLError{
				Type:    ErrorTypeFile,
				Message: "file not found",
				Cause:   errors.New("no such file"),
			},
			expected: "[FILE] file not found | Caused by: no such file",
		},
		{
			name: "error with context",
			err: &RunFromYAMLError{
				Type:    ErrorTypeValidation,
				Message: "invalid value",
				Context: map[string]interface{}{
					"field": "port",
					"value": 99999,
				},
			},
			expected: "[VALIDATION] invalid value | Context: field=port, value=99999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRunFromYAMLError_WithContext(t *testing.T) {
	err := New(ErrorTypeConfig, "test error")
	err.WithContext("key1", "value1")
	err.WithContext("key2", 42)

	if len(err.Context) != 2 {
		t.Errorf("Expected 2 context entries, got %d", len(err.Context))
	}

	if err.Context["key1"] != "value1" {
		t.Errorf("Expected key1=value1, got %v", err.Context["key1"])
	}

	if err.Context["key2"] != 42 {
		t.Errorf("Expected key2=42, got %v", err.Context["key2"])
	}
}

func TestRunFromYAMLError_WithSuggestion(t *testing.T) {
	err := New(ErrorTypeConfig, "test error")
	err.WithSuggestion("suggestion 1")
	err.WithSuggestion("suggestion 2")

	if len(err.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(err.Suggestions))
	}

	if err.Suggestions[0] != "suggestion 1" {
		t.Errorf("Expected first suggestion to be 'suggestion 1', got %v", err.Suggestions[0])
	}
}

func TestRunFromYAMLError_Is(t *testing.T) {
	err1 := New(ErrorTypeConfig, "config error")
	err2 := New(ErrorTypeConfig, "another config error")
	err3 := New(ErrorTypeFile, "file error")

	if !err1.Is(err2) {
		t.Error("Expected err1.Is(err2) to be true")
	}

	if err1.Is(err3) {
		t.Error("Expected err1.Is(err3) to be false")
	}
}

func TestRunFromYAMLError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	err := Wrap(cause, ErrorTypeConfig, "wrapped error")

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}

func TestNewConfigError(t *testing.T) {
	cause := errors.New("invalid syntax")
	err := NewConfigError("configuration parsing failed", cause)

	if err.Type != ErrorTypeConfig {
		t.Errorf("Expected error type CONFIG, got %v", err.Type)
	}

	if err.Cause != cause {
		t.Errorf("Expected cause to be %v, got %v", cause, err.Cause)
	}

	if len(err.Suggestions) == 0 {
		t.Error("Expected at least one suggestion")
	}
}

func TestNewFileError(t *testing.T) {
	cause := errors.New("permission denied")
	filename := "/etc/passwd"
	err := NewFileError("cannot read file", cause, filename)

	if err.Type != ErrorTypeFile {
		t.Errorf("Expected error type FILE, got %v", err.Type)
	}

	if err.Context["filename"] != filename {
		t.Errorf("Expected filename context to be %v, got %v", filename, err.Context["filename"])
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid port number", "port", 99999)

	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected error type VALIDATION, got %v", err.Type)
	}

	if err.Context["field"] != "port" {
		t.Errorf("Expected field context to be 'port', got %v", err.Context["field"])
	}

	if err.Context["value"] != 99999 {
		t.Errorf("Expected value context to be 99999, got %v", err.Context["value"])
	}
}

func TestValidator(t *testing.T) {
	validator := NewValidator()

	// Test with no errors
	if validator.HasErrors() {
		t.Error("Expected no errors initially")
	}

	if validator.GetCombinedError() != nil {
		t.Error("Expected no combined error initially")
	}

	// Add some errors
	validator.AddError(New(ErrorTypeConfig, "error 1"))
	validator.AddError(New(ErrorTypeFile, "error 2"))

	if !validator.HasErrors() {
		t.Error("Expected to have errors after adding them")
	}

	errors := validator.GetErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	combinedErr := validator.GetCombinedError()
	if combinedErr == nil {
		t.Error("Expected combined error to not be nil")
	}
}

func TestValidator_ValidateRequired(t *testing.T) {
	validator := NewValidator()

	// Test with nil value
	validator.ValidateRequired("test_field", nil)
	if !validator.HasErrors() {
		t.Error("Expected validation error for nil value")
	}

	// Reset validator
	validator = NewValidator()

	// Test with empty string
	validator.ValidateRequired("test_field", "")
	if !validator.HasErrors() {
		t.Error("Expected validation error for empty string")
	}

	// Reset validator
	validator = NewValidator()

	// Test with valid string
	validator.ValidateRequired("test_field", "valid_value")
	if validator.HasErrors() {
		t.Error("Expected no validation error for valid string")
	}
}

func TestValidator_ValidatePort(t *testing.T) {
	validator := NewValidator()

	// Test invalid ports
	testCases := []int{0, -1, 65536, 100000}
	for _, port := range testCases {
		validator = NewValidator()
		validator.ValidatePort("test_port", port)
		if !validator.HasErrors() {
			t.Errorf("Expected validation error for port %d", port)
		}
	}

	// Test valid ports
	validPorts := []int{1, 80, 443, 8080, 65535}
	for _, port := range validPorts {
		validator = NewValidator()
		validator.ValidatePort("test_port", port)
		if validator.HasErrors() {
			t.Errorf("Expected no validation error for port %d", port)
		}
	}
}

func TestValidator_ValidateHostname(t *testing.T) {
	validator := NewValidator()

	// Test invalid hostnames
	invalidHostnames := []string{"", "invalid..hostname", "hostname-", "-hostname"}
	for _, hostname := range invalidHostnames {
		validator = NewValidator()
		validator.ValidateHostname("test_host", hostname)
		if !validator.HasErrors() {
			t.Errorf("Expected validation error for hostname '%s'", hostname)
		}
	}

	// Test valid hostnames
	validHostnames := []string{"localhost", "example.com", "sub.example.com", "192.168.1.1"}
	for _, hostname := range validHostnames {
		validator = NewValidator()
		validator.ValidateHostname("test_host", hostname)
		if validator.HasErrors() {
			t.Errorf("Expected no validation error for hostname '%s'", hostname)
		}
	}
}
