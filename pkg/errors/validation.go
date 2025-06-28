package errors

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Validator provides validation functions
type Validator struct {
	errors []error
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		errors: make([]error, 0),
	}
}

// AddError adds an error to the validator
func (v *Validator) AddError(err error) {
	v.errors = append(v.errors, err)
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []error {
	return v.errors
}

// GetCombinedError returns a single error combining all validation errors
func (v *Validator) GetCombinedError() error {
	if !v.HasErrors() {
		return nil
	}
	
	var messages []string
	for _, err := range v.errors {
		messages = append(messages, err.Error())
	}
	
	return New(ErrorTypeValidation, fmt.Sprintf("Validation failed: %s", strings.Join(messages, "; ")))
}

// ValidateRequired checks if a required field is present
func (v *Validator) ValidateRequired(fieldName string, value interface{}) {
	if value == nil {
		v.AddError(NewValidationError(fmt.Sprintf("Field '%s' is required", fieldName), fieldName, value))
		return
	}
	
	switch val := value.(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			v.AddError(NewValidationError(fmt.Sprintf("Field '%s' cannot be empty", fieldName), fieldName, value))
		}
	case []string:
		if len(val) == 0 {
			v.AddError(NewValidationError(fmt.Sprintf("Field '%s' cannot be empty", fieldName), fieldName, value))
		}
	}
}

// ValidateFileExists checks if a file exists
func (v *Validator) ValidateFileExists(fieldName, filename string) {
	if filename == "" {
		v.AddError(NewValidationError(fmt.Sprintf("Filename for '%s' cannot be empty", fieldName), fieldName, filename))
		return
	}
	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		v.AddError(NewFileError(fmt.Sprintf("File '%s' does not exist", filename), err, filename))
	}
}

// ValidateFilePermissions checks if file permissions are valid
func (v *Validator) ValidateFilePermissions(fieldName string, perm os.FileMode) {
	// Check if permissions are within valid range (0000-0777)
	if perm > 0777 {
		v.AddError(NewValidationError(
			fmt.Sprintf("Invalid file permissions for '%s': %o (must be between 0000-0777)", fieldName, perm),
			fieldName,
			perm,
		))
	}
}

// ValidateCommandType checks if command type is valid
func (v *Validator) ValidateCommandType(cmdType string) {
	validTypes := []string{"exec", "shell", "conf", "docker", "docker-compose", "ssh"}
	
	for _, validType := range validTypes {
		if cmdType == validType {
			return
		}
	}
	
	v.AddError(NewValidationError(
		fmt.Sprintf("Invalid command type '%s'. Valid types: %s", cmdType, strings.Join(validTypes, ", ")),
		"type",
		cmdType,
	))
}

// ValidateDockerCommand checks if docker command is valid
func (v *Validator) ValidateDockerCommand(command string) {
	validCommands := []string{"run", "exec"}
	
	for _, validCmd := range validCommands {
		if command == validCmd {
			return
		}
	}
	
	v.AddError(NewValidationError(
		fmt.Sprintf("Invalid docker command '%s'. Valid commands: %s", command, strings.Join(validCommands, ", ")),
		"command",
		command,
	))
}

// ValidatePort checks if port number is valid
func (v *Validator) ValidatePort(fieldName string, port int) {
	if port < 1 || port > 65535 {
		v.AddError(NewValidationError(
			fmt.Sprintf("Invalid port number for '%s': %d (must be between 1-65535)", fieldName, port),
			fieldName,
			port,
		))
	}
}

// ValidateHostname checks if hostname is valid
func (v *Validator) ValidateHostname(fieldName, hostname string) {
	if hostname == "" {
		v.AddError(NewValidationError(fmt.Sprintf("Hostname for '%s' cannot be empty", fieldName), fieldName, hostname))
		return
	}
	
	// Basic hostname validation
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if hostname != "localhost" && !hostnameRegex.MatchString(hostname) {
		v.AddError(NewValidationError(
			fmt.Sprintf("Invalid hostname format for '%s': %s", fieldName, hostname),
			fieldName,
			hostname,
		))
	}
}

// ValidateLogLevel checks if log level is valid
func (v *Validator) ValidateLogLevel(level string) {
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	
	for _, validLevel := range validLevels {
		if strings.ToLower(level) == validLevel {
			return
		}
	}
	
	v.AddError(NewValidationError(
		fmt.Sprintf("Invalid log level '%s'. Valid levels: %s", level, strings.Join(validLevels, ", ")),
		"level",
		level,
	))
}

// ValidateOutputType checks if output type is valid
func (v *Validator) ValidateOutputType(output string) {
	validOutputs := []string{"stdout", "file", "rest"}
	
	for _, validOutput := range validOutputs {
		if strings.ToLower(output) == validOutput {
			return
		}
	}
	
	v.AddError(NewValidationError(
		fmt.Sprintf("Invalid output type '%s'. Valid types: %s", output, strings.Join(validOutputs, ", ")),
		"output",
		output,
	))
}

// ValidateAIModel checks if AI model name is reasonable
func (v *Validator) ValidateAIModel(model string) {
	if model == "" {
		v.AddError(NewValidationError("AI model cannot be empty", "ai-model", model))
		return
	}
	
	// Basic validation - could be extended with specific model names
	if len(model) < 3 {
		v.AddError(NewValidationError(
			fmt.Sprintf("AI model name '%s' seems too short", model),
			"ai-model",
			model,
		))
	}
}

// ValidateShellType checks if shell type is valid
func (v *Validator) ValidateShellType(shellType string) {
	validShells := []string{"bash", "sh", "zsh", "fish", "powershell", "cmd"}
	
	for _, validShell := range validShells {
		if strings.ToLower(shellType) == validShell {
			return
		}
	}
	
	v.AddError(NewValidationError(
		fmt.Sprintf("Invalid shell type '%s'. Valid types: %s", shellType, strings.Join(validShells, ", ")),
		"shell-type",
		shellType,
	))
}

// ValidateEnvironmentVariable checks if environment variable name is valid
func (v *Validator) ValidateEnvironmentVariable(name, value string) {
	if name == "" {
		v.AddError(NewValidationError("Environment variable name cannot be empty", "env.key", name))
		return
	}
	
	// Environment variable names should follow certain conventions
	envVarRegex := regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
	if !envVarRegex.MatchString(name) {
		v.AddError(NewValidationError(
			fmt.Sprintf("Invalid environment variable name '%s'. Should contain only uppercase letters, numbers, and underscores", name),
			"env.key",
			name,
		))
	}
}

// ValidateDockerContainer checks if container name is valid
func (v *Validator) ValidateDockerContainer(container string) {
	if container == "" {
		v.AddError(NewValidationError("Docker container name cannot be empty", "container", container))
		return
	}
	
	// Docker container names have specific rules
	containerRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
	if !containerRegex.MatchString(container) {
		v.AddError(NewValidationError(
			fmt.Sprintf("Invalid docker container name '%s'. Must start with alphanumeric character and contain only letters, numbers, underscores, periods, and hyphens", container),
			"container",
			container,
		))
	}
}
