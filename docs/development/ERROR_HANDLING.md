# Error Handling Documentation

## Overview

The runfromyaml project now includes a comprehensive error handling system that provides:

- **Structured Error Types**: Categorized errors for different scenarios
- **Contextual Information**: Rich error context with suggestions
- **Validation Framework**: Input validation with helpful feedback
- **Centralized Handling**: Consistent error processing and formatting
- **Debug Support**: Stack traces and detailed debugging information

## Error Types

### ErrorType Categories

```go
const (
    ErrorTypeConfig     ErrorType = "CONFIG"     // Configuration errors
    ErrorTypeFile       ErrorType = "FILE"       // File system errors
    ErrorTypeYAML       ErrorType = "YAML"       // YAML parsing errors
    ErrorTypeExecution  ErrorType = "EXECUTION"  // Command execution errors
    ErrorTypeNetwork    ErrorType = "NETWORK"    // Network-related errors
    ErrorTypeValidation ErrorType = "VALIDATION" // Input validation errors
    ErrorTypeAI         ErrorType = "AI"         // AI/OpenAI API errors
    ErrorTypeDocker     ErrorType = "DOCKER"     // Docker-related errors
    ErrorTypeSSH        ErrorType = "SSH"        // SSH connection errors
    ErrorTypeInternal   ErrorType = "INTERNAL"   // Internal/panic errors
)
```

## Structured Error Format

### RunFromYAMLError Structure

```go
type RunFromYAMLError struct {
    Type        ErrorType                  // Error category
    Message     string                     // Human-readable message
    Cause       error                      // Underlying error (if any)
    Context     map[string]interface{}     // Additional context
    StackTrace  string                     // Stack trace (debug mode)
    Suggestions []string                   // Helpful suggestions
}
```

### Error Output Format

```
❌ Error: Failed to read configuration file
   Cause: open config.yaml: no such file or directory
   Context:
     filename: config.yaml
   💡 Suggestions:
     • Verify the file exists and you have proper permissions
```

## Error Constructors

### Predefined Error Constructors

```go
// Configuration errors
NewConfigError(message string, cause error) *RunFromYAMLError

// File system errors
NewFileError(message string, cause error, filename string) *RunFromYAMLError

// YAML parsing errors
NewYAMLError(message string, cause error, filename string) *RunFromYAMLError

// Command execution errors
NewExecutionError(message string, cause error, command string) *RunFromYAMLError

// Validation errors
NewValidationError(message string, field string, value interface{}) *RunFromYAMLError

// AI/OpenAI errors
NewAIError(message string, cause error) *RunFromYAMLError

// Docker errors
NewDockerError(message string, cause error, container string) *RunFromYAMLError

// SSH errors
NewSSHError(message string, cause error, host string) *RunFromYAMLError

// Network errors
NewNetworkError(message string, cause error) *RunFromYAMLError
```

### Usage Examples

```go
// File error with context
err := NewFileError("Cannot read configuration", originalErr, "config.yaml")

// Validation error with field context
err := NewValidationError("Invalid port number", "port", 99999)

// AI error with suggestion
err := NewAIError("OpenAI API request failed", apiErr).
    WithSuggestion("Check your API key and network connection")
```

## Validation Framework

### Validator Usage

```go
validator := errors.NewValidator()

// Validate required fields
validator.ValidateRequired("api-key", cfg.AIKey)

// Validate port numbers
validator.ValidatePort("port", cfg.Port)

// Validate hostnames
validator.ValidateHostname("host", cfg.Host)

// Check for errors
if validator.HasErrors() {
    return validator.GetCombinedError()
}
```

### Available Validation Functions

- `ValidateRequired(fieldName, value)` - Check required fields
- `ValidateFileExists(fieldName, filename)` - Verify file existence
- `ValidateFilePermissions(fieldName, perm)` - Check file permissions
- `ValidateCommandType(cmdType)` - Validate command types
- `ValidateDockerCommand(command)` - Validate Docker commands
- `ValidatePort(fieldName, port)` - Check port ranges
- `ValidateHostname(fieldName, hostname)` - Validate hostnames
- `ValidateLogLevel(level)` - Check log levels
- `ValidateOutputType(output)` - Validate output types
- `ValidateAIModel(model)` - Check AI model names
- `ValidateShellType(shellType)` - Validate shell types
- `ValidateEnvironmentVariable(name, value)` - Check env vars
- `ValidateDockerContainer(container)` - Validate container names

## Error Handler

### ErrorHandler Usage

```go
// Create error handler
errorHandler := errors.NewErrorHandler(debug)

// Handle errors
errorHandler.Handle(err)

// Panic recovery
defer errorHandler.Recovery()
```

### Error Handler Features

- **Structured Output**: Formatted error messages with context
- **Debug Mode**: Stack traces and detailed information
- **Panic Recovery**: Converts panics to structured errors
- **Suggestion Display**: Shows helpful suggestions to users

## Integration Examples

### Main Application Integration

```go
func main() {
    var errorHandler *errors.ErrorHandler
    
    defer func() {
        if errorHandler != nil {
            errorHandler.Recovery()
        }
    }()

    cfg := config.New()
    if err := cfg.ParseFlags(); err != nil {
        errorHandler = errors.NewErrorHandler(false)
        errorHandler.Handle(errors.NewConfigError("Failed to parse flags", err))
        os.Exit(1)
    }

    errorHandler = errors.NewErrorHandler(cfg.Debug)
    
    // Validate configuration
    if err := validateConfig(cfg); err != nil {
        errorHandler.Handle(err)
        os.Exit(1)
    }
}
```

### Command Validation

```go
func validateCommand(cmd *Command) error {
    validator := errors.NewValidator()
    
    // Validate command type
    validator.ValidateCommandType(string(cmd.Type))
    
    // Type-specific validation
    switch cmd.Type {
    case CommandTypeDocker:
        if container, ok := cmd.Options["container"].(string); !ok || container == "" {
            validator.AddError(errors.NewValidationError(
                "Docker command requires 'container' field", 
                "container", 
                container,
            ))
        }
    }
    
    return validator.GetCombinedError()
}
```

### File Operations

```go
func loadYAMLConfig(cfg *config.Config) error {
    ydata, err := os.ReadFile(cfg.File)
    if err != nil {
        return errors.NewFileError("Failed to read configuration file", err, cfg.File)
    }
    
    if err := cfg.LoadFromYAML(ydata); err != nil {
        return errors.NewYAMLError("Failed to parse YAML configuration", err, cfg.File)
    }
    
    return nil
}
```

## Best Practices

### 1. Use Appropriate Error Types

Choose the most specific error type for each scenario:

```go
// Good
return errors.NewFileError("Cannot read file", err, filename)

// Less specific
return errors.New(errors.ErrorTypeGeneric, "File error")
```

### 2. Provide Context

Add relevant context information:

```go
err := errors.NewExecutionError("Command failed", originalErr, command).
    WithContext("exit_code", exitCode).
    WithContext("working_dir", workingDir)
```

### 3. Include Helpful Suggestions

Guide users toward solutions:

```go
err := errors.NewValidationError("Invalid port", "port", port).
    WithSuggestion("Port must be between 1 and 65535").
    WithSuggestion("Check if the port is already in use")
```

### 4. Use Validation Early

Validate inputs before processing:

```go
func processConfig(cfg *Config) error {
    validator := errors.NewValidator()
    validator.ValidateRequired("api-key", cfg.APIKey)
    validator.ValidatePort("port", cfg.Port)
    
    if validator.HasErrors() {
        return validator.GetCombinedError()
    }
    
    // Process configuration...
}
```

### 5. Handle Panics Gracefully

Use panic recovery in main functions:

```go
func main() {
    errorHandler := errors.NewErrorHandler(debug)
    defer errorHandler.Recovery()
    
    // Application logic...
}
```

## Testing

### Error Testing Examples

```go
func TestErrorHandling(t *testing.T) {
    // Test error creation
    err := errors.NewFileError("test error", nil, "test.txt")
    assert.Equal(t, errors.ErrorTypeFile, err.Type)
    assert.Equal(t, "test.txt", err.Context["filename"])
    
    // Test validation
    validator := errors.NewValidator()
    validator.ValidatePort("port", 99999)
    assert.True(t, validator.HasErrors())
}
```

## Migration Guide

### From Old Error Handling

**Before:**
```go
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

**After:**
```go
if err != nil {
    structuredErr := errors.NewFileError("Operation failed", err, filename)
    errorHandler.Handle(structuredErr)
    return structuredErr
}
```

### Benefits of New System

1. **Consistent Formatting**: All errors follow the same format
2. **Rich Context**: Errors include relevant context information
3. **User Guidance**: Suggestions help users resolve issues
4. **Better Debugging**: Stack traces and detailed information in debug mode
5. **Type Safety**: Categorized errors for better handling
6. **Validation**: Early input validation prevents runtime errors

## Configuration

### Debug Mode

Enable debug mode for detailed error information:

```bash
runfromyaml --debug --file config.yaml
```

### Error Output Examples

**Production Mode:**
```
❌ Error: Failed to connect to Docker daemon
   💡 Suggestions:
     • Ensure Docker is running and the container exists
```

**Debug Mode:**
```
❌ Error: Failed to connect to Docker daemon
   Cause: dial unix /var/run/docker.sock: connect: no such file or directory
   Context:
     container: alpine:latest
     command: run
   💡 Suggestions:
     • Ensure Docker is running and the container exists
   Stack Trace:
     /path/to/main.go:123 main.handleDockerCommand
     /path/to/cli.go:456 cli.Execute
```
