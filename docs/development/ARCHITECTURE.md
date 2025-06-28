# Architecture Documentation

## Project Structure

```
runfromyaml/
├── main.go                 # Main entry point
├── pkg/                    # Package modules
│   ├── config/            # Configuration management
│   │   ├── config.go      # Core configuration struct and parsing
│   │   └── yaml.go        # YAML configuration loading
│   ├── cli/               # Command-line interface
│   ├── openai/            # AI integration
│   │   └── openai.go      # OpenAI API client
│   ├── restapi/           # REST API functionality
│   ├── functions/         # Utility functions
│   └── docker/            # Docker integration
├── examples/              # Example YAML files
├── version/               # Version management
└── README.md              # Main documentation
```

## Core Components

### 1. Configuration System (`pkg/config`)

The configuration system handles both command-line flags and YAML-based configuration:

- **config.go**: Defines the main `Config` struct with all application settings
- **yaml.go**: Handles loading configuration from YAML files
- Supports environment variable expansion
- Merges command-line flags with YAML configuration

### 2. Main Application Flow (`main.go`)

The main application follows a modular approach:

1. **Initialization**: Load configuration from flags and YAML
2. **Mode Detection**: Determine which mode to run (AI, Shell, REST, File)
3. **Mode Execution**: Execute the appropriate handler function

#### Execution Modes

- **AI Mode**: `handleAIMode()` - Interacts with OpenAI API
- **File Mode**: `handleFileExecution()` - Processes YAML command files
- **REST Mode**: `handleRestMode()` - Starts REST API server
- **Shell Mode**: `handleShellMode()` - Interactive command recording

### 3. AI Integration (`pkg/openai`)

- OpenAI API client implementation
- Supports multiple models (text-davinci-003, GPT-4, etc.)
- Command generation and assistance
- Configurable shell types for command examples

### 4. Command Processing (`pkg/cli`)

Handles the core YAML command processing:

- YAML parsing and validation
- Command execution (exec, shell, docker, etc.)
- Environment variable expansion
- Interactive shell recording

### 5. REST API (`pkg/restapi`)

- HTTP server for remote command execution
- Authentication support (optional)
- YAML payload processing
- Response formatting

## Configuration Architecture

### Command-Line Configuration

All configuration options are available as command-line flags:

```go
type Config struct {
    Debug     bool
    Rest      bool
    NoAuth    bool
    RestOut   bool
    NoFile    bool
    AI        bool
    Shell     bool
    File      string
    Host      string
    User      string
    AIInput   string
    AIKey     string
    AIModel   string
    AICmdType string
    ShellType string
    Port      int
}
```

### YAML Configuration

Configuration can also be embedded in YAML files:

```yaml
options:
  - key: "rest"
    value: false
  - key: "ai"
    value: true
  - key: "ai-model"
    value: "gpt-4"
```

## YAML Command Structure

### Supported Command Types

1. **exec**: Direct OS command execution
2. **shell**: Shell-wrapped command execution
3. **conf**: Configuration file creation
4. **docker**: Docker container commands
5. **docker-compose**: Docker Compose operations
6. **ssh**: Remote SSH command execution

### Block Structure

Each command block follows this structure:

```yaml
cmd:
  - type: "exec|shell|conf|docker|docker-compose|ssh"
    name: "block_name"
    desc: "description"
    expandenv: true|false
    values:
      - "command1"
      - "command2"
    # Type-specific options...
```

## Extension Points

### Adding New Command Types

1. Define new type in command processing logic
2. Add type-specific configuration options
3. Implement execution handler
4. Update documentation

### Adding New AI Providers

1. Create new package under `pkg/`
2. Implement common AI interface
3. Add configuration options
4. Update command-line flags

### Adding New Output Formats

1. Extend logging configuration
2. Add new output handlers
3. Update REST API response formatting

## Security Considerations

- REST API authentication (optional but recommended)
- Environment variable expansion controls
- File permission handling for configuration files
- SSH key management for remote execution

## Performance Considerations

- Concurrent command execution (planned)
- Memory usage for large YAML files
- Network timeouts for AI and REST operations
- Docker container lifecycle management

## Testing Strategy

Currently missing comprehensive tests. Planned test structure:

- Unit tests for each package
- Integration tests for command execution
- End-to-end tests for complete workflows
- Performance tests for large YAML files
