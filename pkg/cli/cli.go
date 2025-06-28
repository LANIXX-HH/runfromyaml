package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"

	functions "github.com/lanixx/runfromyaml/pkg/functions"
)

// CommandType represents the type of command to execute
type CommandType string

const (
	CommandTypeExec          CommandType = "exec"
	CommandTypeShell         CommandType = "shell"
	CommandTypeDocker        CommandType = "docker"
	CommandTypeDockerCompose CommandType = "docker-compose"
	CommandTypeSSH           CommandType = "ssh"
	CommandTypeConfig        CommandType = "conf"
)

// OutputType represents where command output should be directed
type OutputType string

const (
	OutputTypeRest   OutputType = "rest"
	OutputTypeFile   OutputType = "file"
	OutputTypeStdout OutputType = "stdout"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
	LogLevelDebug LogLevel = "debug"
)

// Environment manages environment variables
type Environment struct {
	variables map[string]string
	shell     []string
}

// NewEnvironment creates a new environment manager
func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]string),
		shell:     make([]string, 0),
	}
}

// Set sets an environment variable
func (e *Environment) Set(key, value string) {
	e.variables[key] = value
	e.shell = append(e.shell, key+"="+value)
	_ = os.Setenv(key, value)
}

// Get retrieves an environment variable
func (e *Environment) Get(key string) string {
	return e.variables[key]
}

// GetVariables returns the entire variables map
func (e *Environment) GetVariables() map[string]string {
	return e.variables
}

// Shell returns the environment variables in shell format
func (e *Environment) Shell() []string {
	return e.shell
}

// Command represents a command to be executed
type Command struct {
	Type        CommandType
	Description string
	Values      []string
	Options     map[string]interface{}
	Env         *Environment
}

// CommandConfig holds common configuration for command execution
type CommandConfig struct {
	Env         *Environment
	Level       LogLevel
	Output      OutputType
	Description string
	WaitGroup   *sync.WaitGroup
}

// CommandExecutor handles command execution
type CommandExecutor struct {
	config CommandConfig
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(config CommandConfig) *CommandExecutor {
	return &CommandExecutor{config: config}
}

// Execute runs the command based on its type
func (e *CommandExecutor) Execute(cmd *Command) error {
	switch cmd.Type {
	case CommandTypeExec:
		return e.executeExecCommand(cmd)
	case CommandTypeShell:
		return e.executeShellCommand(cmd)
	case CommandTypeDocker:
		return e.executeDockerCommand(cmd)
	case CommandTypeDockerCompose:
		return e.executeDockerComposeCommand(cmd)
	case CommandTypeSSH:
		return e.executeSSHCommand(cmd)
	case CommandTypeConfig:
		return e.handleConfigCommand(cmd)
	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}

func (e *CommandExecutor) executeExecCommand(cmd *Command) error {
	// Handle empty values gracefully
	if len(cmd.Values) == 0 {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# exec command with empty values - skipping execution")
		return nil
	}

	cmds := splitCommands(cmd.Values)
	for _, cmdStr := range cmds {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue // Skip empty commands
		}
		cmdStr = os.ExpandEnv(cmdStr)
		cmdArgs := strings.Fields(cmdStr)
		if len(cmdArgs) == 0 {
			continue // Skip if no arguments after expansion
		}
		if err := e.runCommand(cmdArgs); err != nil {
			return err
		}
	}
	return nil
}

func (e *CommandExecutor) executeShellCommand(cmd *Command) error {
	// Handle empty values gracefully
	if len(cmd.Values) == 0 {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# shell command with empty values - skipping execution")
		return nil
	}

	// Filter out empty values
	nonEmptyValues := make([]string, 0, len(cmd.Values))
	for _, val := range cmd.Values {
		if strings.TrimSpace(val) != "" {
			nonEmptyValues = append(nonEmptyValues, val)
		}
	}

	if len(nonEmptyValues) == 0 {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# shell command with only empty values - skipping execution")
		return nil
	}

	// Join shell commands - semicolons are already present in the values
	args := append([]string{"bash", "-c"}, strings.Join(nonEmptyValues, " "))
	return e.runCommand(args)
}

func (e *CommandExecutor) executeDockerCommand(cmd *Command) error {
	args := e.buildDockerArgs(cmd)

	// If values are empty, we can't execute docker commands as they require commands to run
	if len(cmd.Values) == 0 {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# docker command with empty values - skipping execution (docker commands require commands to execute)")
		return nil
	}

	cmds := splitCommands(cmd.Values)
	for _, cmdStr := range cmds {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue // Skip empty commands
		}
		cmdStr = os.ExpandEnv(cmdStr)
		cmdArgs := strings.Fields(cmdStr)
		if len(cmdArgs) == 0 {
			continue // Skip if no arguments after expansion
		}
		fullArgs := append(args, cmdArgs...)
		if err := e.runCommand(fullArgs); err != nil {
			return err
		}
	}
	return nil
}

func (e *CommandExecutor) executeDockerComposeCommand(cmd *Command) error {
	args := e.buildDockerComposeArgs(cmd)

	// If values are empty, execute the docker-compose command without additional commands
	if len(cmd.Values) == 0 {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# docker-compose command with empty values - executing base command only")
		return e.runCommand(args)
	}

	// If values are provided, execute additional commands inside containers
	cmds := splitCommands(cmd.Values)
	for _, cmdStr := range cmds {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue // Skip empty commands
		}
		cmdStr = os.ExpandEnv(cmdStr)
		cmdArgs := strings.Fields(cmdStr)
		if len(cmdArgs) == 0 {
			continue // Skip if no arguments after expansion
		}
		fullArgs := append(args, cmdArgs...)
		if err := e.runCommand(fullArgs); err != nil {
			return err
		}
	}
	return nil
}

func (e *CommandExecutor) executeSSHCommand(cmd *Command) error {
	// Handle empty values gracefully
	if len(cmd.Values) == 0 {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# ssh command with empty values - skipping execution")
		return nil
	}

	args := e.buildSSHArgs(cmd)
	cmds := splitCommands(cmd.Values)
	for _, cmdStr := range cmds {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue // Skip empty commands
		}
		cmdStr = os.ExpandEnv(cmdStr)
		cmdArgs := strings.Fields(cmdStr)
		if len(cmdArgs) == 0 {
			continue // Skip if no arguments after expansion
		}
		fullArgs := append(args, cmdArgs...)
		if err := e.runCommand(fullArgs); err != nil {
			return err
		}
	}
	return nil
}

func (e *CommandExecutor) buildDockerArgs(cmd *Command) []string {
	command := cmd.Options["command"].(string)
	container := cmd.Options["container"].(string)
	if command == "run" {
		return []string{"docker", command, "-it", "--rm", container, "sh", "-c"}
	}
	return []string{"docker", command, container, "sh", "-c"}
}

func (e *CommandExecutor) buildDockerComposeArgs(cmd *Command) []string {
	args := []string{"docker", "compose"}

	// Check if environment expansion is enabled
	expandenv := false
	if expandenvOpt, exists := cmd.Options["expandenv"]; exists {
		expandenv = expandenvOpt.(bool)
	}

	// Handle dcoptions
	if opts, exists := cmd.Options["dcoptions"]; exists {
		if optsSlice, ok := opts.([]interface{}); ok {
			for _, opt := range optsSlice {
				if strOpt, ok := opt.(string); ok {
					if expandenv {
						strOpt = os.ExpandEnv(strOpt)
					}
					// Split the option into separate arguments
					optArgs := strings.Fields(strOpt)
					args = append(args, optArgs...)
				}
			}
		}
	}

	// Add command
	if cmd, exists := cmd.Options["command"]; exists {
		if cmdStr, ok := cmd.(string); ok {
			if expandenv {
				cmdStr = os.ExpandEnv(cmdStr)
			}
			args = append(args, cmdStr)
		}
	}

	// Handle cmdoptions
	if opts, exists := cmd.Options["cmdoptions"]; exists {
		if optsSlice, ok := opts.([]interface{}); ok {
			for _, opt := range optsSlice {
				if strOpt, ok := opt.(string); ok {
					if expandenv {
						strOpt = os.ExpandEnv(strOpt)
					}
					// Split the option into separate arguments
					optArgs := strings.Fields(strOpt)
					args = append(args, optArgs...)
				}
			}
		}
	}

	// Add service if it exists
	if service, exists := cmd.Options["service"]; exists {
		if serviceStr, ok := service.(string); ok && serviceStr != "" {
			if expandenv {
				serviceStr = os.ExpandEnv(serviceStr)
			}
			args = append(args, serviceStr)
		}
	}

	return args
}

func (e *CommandExecutor) buildSSHArgs(cmd *Command) []string {
	user := cmd.Options["user"].(string)
	host := cmd.Options["host"].(string)
	port := cmd.Options["port"].(int)

	if expandenv, exists := cmd.Options["expandenv"]; exists {
		if expandenv.(bool) {
			user = os.ExpandEnv(user)
			host = os.ExpandEnv(host)
		}
	}

	args := []string{"ssh", "-p", strconv.Itoa(port), "-l", user, host}

	// Handle SSH options
	if opts, exists := cmd.Options["options"]; exists {
		if optsSlice, ok := opts.([]interface{}); ok {
			for _, opt := range optsSlice {
				if strOpt, ok := opt.(string); ok {
					// Apply expandenv to SSH options if enabled
					if expandenv, exists := cmd.Options["expandenv"]; exists && expandenv.(bool) {
						strOpt = os.ExpandEnv(strOpt)
					}
					args = append(args, strOpt)
				}
			}
		}
	}

	return args
}

func (e *CommandExecutor) handleConfigCommand(cmd *Command) error {
	var confdata, confdest string
	var confperm os.FileMode
	expandenv := false

	if expandenvOpt, exists := cmd.Options["expandenv"]; exists {
		expandenv = expandenvOpt.(bool)
	}

	if data := cmd.Options["confdata"]; data != nil {
		confdata = data.(string)
		if expandenv {
			confdata = functions.GoTemplate(e.config.Env.variables, confdata)
		}
	}

	// Only add description if confdata is not empty
	if confdata != "" {
		confdata = cmd.Description + confdata
	}

	if dest := cmd.Options["confdest"]; dest != nil {
		confdest = dest.(string)
		if expandenv {
			confdest = os.ExpandEnv(confdest)
		}
	}
	if perm := cmd.Options["confperm"]; perm != nil {
		confperm = os.FileMode(int(perm.(int)))
	}

	// Handle empty config gracefully
	if confdata == "" && confdest == "" {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# config command with empty data and destination - skipping")
		return nil
	}

	if confdata != "" && confdest != "" && string(rune(confperm)) != "" {
		functions.WriteFile(confdata, confdest, confperm)
		functions.PrintSwitch(color.FgGreen, string(e.config.Level), string(e.config.Output), "# create ", confdest)
	} else if confdest != "" {
		functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# config command missing data or permissions for ", confdest)
	}

	return nil
}

func (e *CommandExecutor) runCommand(cmd []string) error {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = append(os.Environ(), e.config.Env.shell...)
	functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), strings.Trim(fmt.Sprint(cmd), "[]"), "\n")

	switch e.config.Output {
	case OutputTypeRest:
		out, err := command.CombinedOutput()
		if err != nil {
			functions.PrintRest(color.FgRed, "error", "Error: ", err, string(out))
			return err
		}
		functions.PrintRest(color.FgHiWhite, string(e.config.Level), string(out))
	case OutputTypeFile:
		out, err := command.CombinedOutput()
		if err != nil {
			functions.PrintFile("error", "Error: ", err, string(out))
			return err
		}
		functions.PrintFile(string(e.config.Level), string(out))
	case OutputTypeStdout:
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			functions.PrintColor(color.FgRed, "error", "Error: ", err)
			return err
		}
	}
	return nil
}

func splitCommands(cmd []string) []string {
	return strings.Split(strings.Join(cmd, " "), ";")
}

// Runfromyaml executes commands from YAML file
// Runfromyaml processes and executes commands from YAML data
func Runfromyaml(yamlFile []byte, debug bool) error {
	var yamlDocument map[interface{}]interface{}
	if err := yaml.Unmarshal(yamlFile, &yamlDocument); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	env := NewEnvironment()
	if env == nil {
		return fmt.Errorf("failed to create environment instance")
	}

	parseEnvironmentVariables(yamlDocument, env)
	outputType, outputLevel := parseLoggingSettings(yamlDocument)

	executor := NewCommandExecutor(CommandConfig{
		Env:       env,
		Level:     LogLevel(outputLevel),
		Output:    OutputType(outputType),
		WaitGroup: &sync.WaitGroup{},
	})

	// Process commands
	if cmdBlocks, ok := yamlDocument["cmd"].([]interface{}); ok {
		for i, cmdBlock := range cmdBlocks {
			if cmdMap, ok := cmdBlock.(map[interface{}]interface{}); ok {
				// Validate required fields
				cmdType, ok := cmdMap["type"].(string)
				if !ok {
					return fmt.Errorf("command block %d: missing or invalid 'type' field", i+1)
				}

				cmd := &Command{
					Type:        CommandType(cmdType),
					Description: functions.EvaluateDescription(cmdMap),
					Values:      functions.ExtractAndExpand(cmdMap, "values"),
					Options:     make(map[string]interface{}),
					Env:         env,
				}

				// Copy all options from the YAML block
				for k, v := range cmdMap {
					if k != "type" && k != "values" {
						cmd.Options[k.(string)] = v
					}
				}

				// Validate command before execution
				if err := validateCommand(cmd); err != nil {
					return fmt.Errorf("command block %d validation failed: %w", i+1, err)
				}

				if err := executor.Execute(cmd); err != nil {
					return fmt.Errorf("failed to execute command block %d (%s): %w", i+1, cmdType, err)
				}
			} else {
				return fmt.Errorf("command block %d: invalid format", i+1)
			}
		}
	}

	return nil
}

// validateCommand validates a command before execution
func validateCommand(cmd *Command) error {
	// Validate command type
	validTypes := []CommandType{
		CommandTypeExec,
		CommandTypeShell,
		CommandTypeDocker,
		CommandTypeDockerCompose,
		CommandTypeSSH,
		CommandTypeConfig,
	}

	isValidType := false
	for _, validType := range validTypes {
		if cmd.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid command type: %s", cmd.Type)
	}

	// Validate type-specific requirements only if values are provided
	// This allows empty command blocks for documentation or placeholder purposes
	switch cmd.Type {
	case CommandTypeDocker:
		// Only validate required fields if values are provided
		if len(cmd.Values) > 0 {
			if container, ok := cmd.Options["container"].(string); !ok || container == "" {
				return fmt.Errorf("docker command with values requires 'container' field")
			}
			if command, ok := cmd.Options["command"].(string); !ok || command == "" {
				return fmt.Errorf("docker command with values requires 'command' field")
			}
		}

	case CommandTypeSSH:
		// Only validate required fields if values are provided
		if len(cmd.Values) > 0 {
			if user, ok := cmd.Options["user"].(string); !ok || user == "" {
				return fmt.Errorf("ssh command with values requires 'user' field")
			}
			if host, ok := cmd.Options["host"].(string); !ok || host == "" {
				return fmt.Errorf("ssh command with values requires 'host' field")
			}
		}

	case CommandTypeConfig:
		// Config commands still require destination and data if they exist
		if dest, ok := cmd.Options["confdest"].(string); ok && dest != "" {
			if data, ok := cmd.Options["confdata"].(string); !ok || data == "" {
				return fmt.Errorf("config command with 'confdest' requires 'confdata' field")
			}
		}
		if data, ok := cmd.Options["confdata"].(string); ok && data != "" {
			if dest, ok := cmd.Options["confdest"].(string); !ok || dest == "" {
				return fmt.Errorf("config command with 'confdata' requires 'confdest' field")
			}
		}
	}

	// Allow empty values blocks - useful for documentation, placeholders, or conditional execution
	// No validation required for empty values

	return nil
}

// InteractiveShell provides an interactive shell for command input
func InteractiveShell(shell string) ([]string, error) {
	var commands []string
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter commands (type 'exit' to finish):")
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		if input != "" {
			commands = append(commands, input)
		}
	}

	return commands, nil
}

func parseEnvironmentVariables(yamlDocument map[interface{}]interface{}, env *Environment) {
	// Parse OS environment variables
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		env.Set(parts[0], parts[1])
	}

	// Parse YAML environment variables
	if envVars, ok := yamlDocument["env"].([]interface{}); ok {
		for _, envVar := range envVars {
			if envMap, ok := envVar.(map[interface{}]interface{}); ok {
				if key, ok := envMap["key"].(string); ok {
					if value, ok := envMap["value"].(string); ok {
						env.Set(key, value)
					}
				}
			}
		}
	}
}

func parseLoggingSettings(yamlDocument map[interface{}]interface{}) (string, string) {
	var outputType, outputLevel string

	if logging, ok := yamlDocument["logging"].([]interface{}); ok {
		for _, logEntry := range logging {
			if entry, ok := logEntry.(map[interface{}]interface{}); ok {
				if output, exists := entry["output"]; exists {
					outputType = output.(string)
				}
				if level, exists := entry["level"]; exists {
					outputLevel = level.(string)
				}
			}
		}
	}

	return outputType, outputLevel
}
