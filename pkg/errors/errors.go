package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeConfig     ErrorType = "CONFIG"
	ErrorTypeFile       ErrorType = "FILE"
	ErrorTypeYAML       ErrorType = "YAML"
	ErrorTypeExecution  ErrorType = "EXECUTION"
	ErrorTypeNetwork    ErrorType = "NETWORK"
	ErrorTypeValidation ErrorType = "VALIDATION"
	ErrorTypeAI         ErrorType = "AI"
	ErrorTypeDocker     ErrorType = "DOCKER"
	ErrorTypeSSH        ErrorType = "SSH"
	ErrorTypeInternal   ErrorType = "INTERNAL"
)

// RunFromYAMLError represents a structured error with context
type RunFromYAMLError struct {
	Type        ErrorType
	Message     string
	Cause       error
	Context     map[string]interface{}
	StackTrace  string
	Suggestions []string
}

// Error implements the error interface
func (e *RunFromYAMLError) Error() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("[%s] %s", e.Type, e.Message))

	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("Caused by: %v", e.Cause))
	}

	if len(e.Context) > 0 {
		var contextParts []string
		for k, v := range e.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("Context: %s", strings.Join(contextParts, ", ")))
	}

	return strings.Join(parts, " | ")
}

// Unwrap returns the underlying error
func (e *RunFromYAMLError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a specific type
func (e *RunFromYAMLError) Is(target error) bool {
	if targetErr, ok := target.(*RunFromYAMLError); ok {
		return e.Type == targetErr.Type
	}
	return false
}

// WithContext adds context information to the error
func (e *RunFromYAMLError) WithContext(key string, value interface{}) *RunFromYAMLError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithSuggestion adds a suggestion for fixing the error
func (e *RunFromYAMLError) WithSuggestion(suggestion string) *RunFromYAMLError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// New creates a new RunFromYAMLError
func New(errorType ErrorType, message string) *RunFromYAMLError {
	return &RunFromYAMLError{
		Type:       errorType,
		Message:    message,
		StackTrace: getStackTrace(),
		Context:    make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, message string) *RunFromYAMLError {
	return &RunFromYAMLError{
		Type:       errorType,
		Message:    message,
		Cause:      err,
		StackTrace: getStackTrace(),
		Context:    make(map[string]interface{}),
	}
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	var stackTrace strings.Builder

	// Skip the first 3 frames (getStackTrace, New/Wrap, caller)
	for i := 3; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		// Only include our package functions
		if strings.Contains(file, "runfromyaml") {
			stackTrace.WriteString(fmt.Sprintf("  %s:%d %s\n", file, line, fn.Name()))
		}
	}

	return stackTrace.String()
}

// Predefined error constructors for common scenarios

// NewConfigError creates a configuration-related error
func NewConfigError(message string, cause error) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeConfig, message)
	return err.WithSuggestion("Check your configuration file syntax and values")
}

// NewFileError creates a file-related error
func NewFileError(message string, cause error, filename string) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeFile, message)
	_ = err.WithContext("filename", filename)
	return err.WithSuggestion("Verify the file exists and you have proper permissions")
}

// NewYAMLError creates a YAML parsing error
func NewYAMLError(message string, cause error, filename string) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeYAML, message)
	_ = err.WithContext("filename", filename)
	return err.WithSuggestion("Check YAML syntax using a YAML validator")
}

// NewExecutionError creates a command execution error
func NewExecutionError(message string, cause error, command string) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeExecution, message)
	_ = err.WithContext("command", command)
	return err.WithSuggestion("Check if the command exists and you have proper permissions")
}

// NewValidationError creates a validation error
func NewValidationError(message string, field string, value interface{}) *RunFromYAMLError {
	err := New(ErrorTypeValidation, message)
	_ = err.WithContext("field", field)
	_ = err.WithContext("value", value)
	return err.WithSuggestion("Review the documentation for valid values")
}

// NewAIError creates an AI-related error
func NewAIError(message string, cause error) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeAI, message)
	return err.WithSuggestion("Check your API key and network connection")
}

// NewDockerError creates a Docker-related error
func NewDockerError(message string, cause error, container string) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeDocker, message)
	_ = err.WithContext("container", container)
	return err.WithSuggestion("Ensure Docker is running and the container exists")
}

// NewSSHError creates an SSH-related error
func NewSSHError(message string, cause error, host string) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeSSH, message)
	_ = err.WithContext("host", host)
	return err.WithSuggestion("Check SSH connectivity and authentication")
}

// NewNetworkError creates a network-related error
func NewNetworkError(message string, cause error) *RunFromYAMLError {
	err := Wrap(cause, ErrorTypeNetwork, message)
	return err.WithSuggestion("Check your network connection and firewall settings")
}

// ErrorHandler provides centralized error handling
type ErrorHandler struct {
	Debug bool
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(debug bool) *ErrorHandler {
	return &ErrorHandler{Debug: debug}
}

// Handle processes and formats errors for output
func (h *ErrorHandler) Handle(err error) {
	if err == nil {
		return
	}

	if rfyErr, ok := err.(*RunFromYAMLError); ok {
		h.handleStructuredError(rfyErr)
	} else {
		h.handleGenericError(err)
	}
}

// handleStructuredError handles RunFromYAMLError instances
func (h *ErrorHandler) handleStructuredError(err *RunFromYAMLError) {
	fmt.Printf("❌ Error: %s\n", err.Message)

	if err.Cause != nil {
		fmt.Printf("   Cause: %v\n", err.Cause)
	}

	if len(err.Context) > 0 {
		fmt.Printf("   Context:\n")
		for k, v := range err.Context {
			fmt.Printf("     %s: %v\n", k, v)
		}
	}

	if len(err.Suggestions) > 0 {
		fmt.Printf("   💡 Suggestions:\n")
		for _, suggestion := range err.Suggestions {
			fmt.Printf("     • %s\n", suggestion)
		}
	}

	if h.Debug && err.StackTrace != "" {
		fmt.Printf("   Stack Trace:\n%s", err.StackTrace)
	}
}

// handleGenericError handles standard Go errors
func (h *ErrorHandler) handleGenericError(err error) {
	fmt.Printf("❌ Error: %v\n", err)

	if h.Debug {
		// Try to get stack trace for debugging
		fmt.Printf("   Stack Trace:\n%s", getStackTrace())
	}
}

// Recovery handles panics and converts them to errors
func (h *ErrorHandler) Recovery() {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case string:
			err = New(ErrorTypeInternal, x)
		case error:
			err = Wrap(x, ErrorTypeInternal, "Panic recovered")
		default:
			err = New(ErrorTypeInternal, fmt.Sprintf("Unknown panic: %v", r))
		}

		h.Handle(err)
	}
}
