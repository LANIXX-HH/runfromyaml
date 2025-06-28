package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/fatih/color"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"github.com/lanixx/runfromyaml/pkg/config"
	"github.com/lanixx/runfromyaml/pkg/errors"
	"github.com/lanixx/runfromyaml/pkg/functions"
	"github.com/lanixx/runfromyaml/pkg/openai"
	"github.com/lanixx/runfromyaml/pkg/restapi"
	"gopkg.in/yaml.v2"
)

func init() {
	functions.Config()
}

func main() {
	// Set up error handler and panic recovery
	var errorHandler *errors.ErrorHandler

	defer func() {
		if errorHandler != nil {
			errorHandler.Recovery()
		}
	}()

	cfg := config.New()
	if err := cfg.ParseFlags(); err != nil {
		errorHandler = errors.NewErrorHandler(false) // Default to non-debug for early errors
		errorHandler.Handle(errors.NewConfigError("Failed to parse command line flags", err))
		os.Exit(1)
	}

	// Initialize error handler with debug setting
	errorHandler = errors.NewErrorHandler(cfg.Debug)

	if cfg.Debug {
		functions.PrintColor(color.FgRed, "debug", "\n", os.Args)
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		errorHandler.Handle(err)
		os.Exit(1)
	}

	// Load YAML configuration if file exists
	if !cfg.NoFile {
		if err := loadYAMLConfig(cfg); err != nil {
			errorHandler.Handle(err)
			// Don't exit here, just warn as this is optional
		}
	}

	// Handle different modes
	var err error
	switch {
	case cfg.AI:
		err = handleAIMode(cfg)
	case cfg.Shell:
		err = handleShellMode(cfg)
	case cfg.Rest:
		err = handleRestMode(cfg)
	case !cfg.NoFile:
		err = handleFileExecution(cfg)
	default:
		err = errors.New(errors.ErrorTypeConfig, "No execution mode specified. Use --help for usage information")
	}

	if err != nil {
		errorHandler.Handle(err)
		os.Exit(1)
	}
}

// validateConfig validates the configuration
func validateConfig(cfg *config.Config) error {
	validator := errors.NewValidator()

	// Validate port
	if cfg.Port != 0 {
		validator.ValidatePort("port", cfg.Port)
	}

	// Validate host
	if cfg.Host != "" {
		validator.ValidateHostname("host", cfg.Host)
	}

	// Validate AI model if AI is enabled
	if cfg.AI && cfg.AIModel != "" {
		validator.ValidateAIModel(cfg.AIModel)
	}

	// Validate shell type
	if cfg.ShellType != "" {
		validator.ValidateShellType(cfg.ShellType)
	}

	// Validate file exists if not disabled
	if !cfg.NoFile && cfg.File != "" {
		validator.ValidateFileExists("file", cfg.File)
	}

	return validator.GetCombinedError()
}

// loadYAMLConfig loads configuration from YAML file
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

// handleAIMode handles AI interaction mode
func handleAIMode(cfg *config.Config) error {
	if len(cfg.AIKey) > 0 {
		openai.Key = cfg.AIKey
		openai.IsAiEnabled = true
	} else {
		openai.IsAiEnabled = false
		return errors.NewAIError("OpenAI API key is required for AI mode", nil).
			WithSuggestion("Set the API key using --ai-key flag or in configuration file")
	}

	openai.Model = cfg.AIModel
	openai.ShellType = cfg.AICmdType

	if len(cfg.AIInput) == 0 {
		return errors.NewValidationError("AI input is required", "ai-in", cfg.AIInput).
			WithSuggestion("Provide input using --ai-in flag")
	}

	// Retry logic with better error handling
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		response, err := openai.OpenAI(openai.Key, openai.Model, cfg.AIInput, openai.ShellType)
		if err == nil {
			fmt.Println(openai.PrintAiResponse(response))
			return nil
		}

		if attempt == maxRetries {
			return errors.NewAIError("Failed to get AI response after multiple attempts", err).
				WithContext("attempts", maxRetries).
				WithSuggestion("Check your API key and network connection")
		}

		// Log retry attempt if debug is enabled
		if cfg.Debug {
			fmt.Printf("AI request attempt %d failed: %v\n", attempt, err)
		}
	}

	return nil
}

// handleFileExecution handles YAML file execution
func handleFileExecution(cfg *config.Config) error {
	ydata, err := os.ReadFile(cfg.File)
	if err != nil {
		return errors.NewFileError("Failed to read command file", err, cfg.File).
			WithSuggestion("Ensure the file exists and you have read permissions")
	}

	// Validate YAML structure
	var ydoc map[interface{}][]interface{}
	if err := yaml.Unmarshal(ydata, &ydoc); err != nil {
		return errors.NewYAMLError("Failed to parse YAML structure", err, cfg.File).
			WithSuggestion("Validate your YAML syntax using a YAML validator")
	}

	// Execute commands with error handling
	if err := cli.Runfromyaml(ydata, cfg.Debug); err != nil {
		return errors.NewExecutionError("Failed to execute commands from YAML file", err, cfg.File)
	}

	return nil
}

// handleRestMode handles REST API mode
func handleRestMode(cfg *config.Config) error {
	fmt.Printf("Starting REST API server on %s:%d\n", cfg.Host, cfg.Port)

	if cfg.RestOut {
		restapi.RestOut = cfg.RestOut
		fmt.Println("Output will be redirected to HTTP response")
	}

	if cfg.NoAuth {
		restapi.RestAuth = false
		fmt.Println("⚠️  WARNING: Authentication is disabled - do not use in public networks!")
	} else {
		restapi.RestAuth = true
		restapi.TempPass = uniuri.New()
		restapi.TempUser = cfg.User
		fmt.Printf("Temporary credentials - User: %s, Password: %s\n", restapi.TempUser, restapi.TempPass)
	}

	// Start REST API server (this is blocking)
	restapi.RestApi(cfg.Port, cfg.Host)
	return nil
}

// handleShellMode handles interactive shell mode
func handleShellMode(cfg *config.Config) error {
	fmt.Println("🐚 Interactive Shell Mode")
	fmt.Println("Your input commands will be recorded to generate YAML structure")
	fmt.Println("Enter 'exit' to stop recording and generate YAML")
	fmt.Println()

	// Create a new environment instance
	env := cli.NewEnvironment()
	if env == nil {
		return errors.New(errors.ErrorTypeInternal, "Failed to create environment instance")
	}

	// Add current environment variables
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			env.Set(parts[0], parts[1])
		}
	}

	// Start interactive shell
	commands, err := cli.InteractiveShell(cfg.ShellType)
	if err != nil {
		return errors.NewExecutionError("Failed to run interactive shell", err, cfg.ShellType).
			WithSuggestion("Ensure the specified shell type is available on your system")
	}

	// Generate YAML from recorded commands
	tempmap := functions.PrintShellCommandsAsYaml(commands, env.GetVariables())
	tempyaml, err := yaml.Marshal(tempmap)
	if err != nil {
		return errors.NewYAMLError("Failed to generate YAML from recorded commands", err, "").
			WithSuggestion("Check if the recorded commands contain valid data")
	}

	fmt.Println("\n📄 Generated YAML:")
	fmt.Println("---")
	fmt.Println(string(tempyaml))

	return nil
}
