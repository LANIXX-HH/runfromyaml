# Empty Values and Empty Command Blocks Support

## Overview

This enhancement adds support for empty `values` blocks and empty command blocks in runfromyaml. This functionality is useful for:

- **Documentation**: Creating placeholder commands that document future implementations
- **Conditional Execution**: Commands that may be conditionally populated
- **Template Creation**: YAML templates with placeholder blocks
- **Development Workflow**: Incremental development where commands are added over time

## Changes Made

### 1. Validation Changes

**Before**: Commands were required to have at least one value (except config type)
```go
if cmd.Type != CommandTypeConfig && len(cmd.Values) == 0 {
    return fmt.Errorf("command must have at least one value")
}
```

**After**: Empty values are allowed for all command types
- Validation only checks required fields when values are provided
- Empty command blocks are treated as valid placeholders

### 2. Execution Changes

All command execution methods now handle empty values gracefully:

- **exec**: Skips execution with informative message
- **shell**: Skips execution with informative message  
- **docker**: Skips execution with informative message (docker commands require commands to execute)
- **docker-compose**: **FIXED** - Now executes the base docker-compose command even with empty values
- **ssh**: Skips execution with informative message
- **conf**: Skips creation when both data and destination are empty

#### Docker-Compose Behavior Change

**Before**: Docker-compose commands with empty values were completely skipped
```yaml
# This would be skipped entirely
- type: docker-compose
  command: up
  values: []
```

**After**: Docker-compose commands with empty values execute the base command
```yaml
# This now executes: docker compose up -d
- type: docker-compose
  command: up
  cmdoptions: ["-d"]
  values: []
```

This is important because many docker-compose operations (like `up`, `down`, `build`, `pull`) don't require additional commands to be executed inside containers.

### 3. Value Extraction Improvements

The `ExtractAndExpand` function was enhanced to:
- Return empty slice instead of nil for consistency
- Handle various YAML value types properly
- Support empty arrays and nil values

## Usage Examples

### Empty Values Block

```yaml
cmd:
  # Placeholder for future setup commands
  - type: exec
    name: "future-setup"
    desc: "Setup commands will be added here"
    values: []
```

### Completely Empty Command Block

```yaml
cmd:
  # Documentation placeholder
  - type: shell
    name: "deployment-placeholder"
    desc: "Deployment commands to be implemented"
    values:
```

### Mixed Empty and Working Commands

```yaml
cmd:
  # Empty placeholder
  - type: exec
    name: "placeholder"
    desc: "Future implementation"
    values: []
    
  # Working command
  - type: exec
    name: "working"
    desc: "Current implementation"
    values:
      - echo
      - "This works!"
```

### Docker-Compose with Empty Values

```yaml
cmd:
  # Build services without additional commands
  - type: docker-compose
    name: "build-all"
    desc: "Build all services"
    dcoptions:
      - -f
      - docker-compose.yml
    command: build
    cmdoptions: []
    service: ""
    values: []
    
  # Start services in detached mode
  - type: docker-compose
    name: "start-services"
    desc: "Start all services"
    dcoptions:
      - -f
      - docker-compose.yml
    command: up
    cmdoptions:
      - -d
    service: ""
    values: []
    
  # Run service with commands (non-empty values)
  - type: docker-compose
    name: "run-with-commands"
    desc: "Run service with additional commands"
    dcoptions:
      - -f
      - docker-compose.yml
    command: run
    cmdoptions:
      - --rm
    service: web
    values:
      - echo "Hello from container"
      - pwd
```

### Conditional Config Blocks

```yaml
cmd:
  # Empty config - useful for conditional creation
  - type: conf
    name: "conditional-config"
    desc: "Config created only when needed"
    confdest: ""
    confdata: ""
    confperm: 0644
```

## Output Behavior

When empty command blocks are encountered, runfromyaml will:

1. **Log the behavior**: Display an informative message about the action taken
2. **Continue processing**: Move to the next command block without error
3. **Maintain flow**: Preserve the overall execution flow

### Example Output for Different Command Types

**Most command types (skip execution):**
```
# exec command with empty values - skipping execution
# shell command with empty values - skipping execution
# docker command with empty values - skipping execution (docker commands require commands to execute)
# ssh command with empty values - skipping execution
# config command with empty data and destination - skipping
```

**Docker-compose (executes base command):**
```
# docker-compose command with empty values - executing base command only
docker compose -f docker-compose.yml up -d
```

## Benefits

1. **Better Documentation**: YAML files can serve as living documentation
2. **Incremental Development**: Add commands gradually without breaking existing workflows
3. **Template Support**: Create reusable YAML templates with placeholders
4. **Conditional Logic**: Support for commands that may or may not execute based on conditions
5. **Error Reduction**: No more validation errors for incomplete command blocks

## Backward Compatibility

This change is fully backward compatible:
- Existing YAML files continue to work unchanged
- No breaking changes to command syntax
- All existing functionality preserved

## Testing

The enhancement has been tested with:
- Empty `values: []` arrays
- Missing `values:` keys
- Mixed empty and working commands
- All command types (exec, shell, docker, docker-compose, ssh, conf)
- Environment variable expansion
- Debug mode output

## Example Demo

See `examples/empty-values-demo.yaml` for a comprehensive demonstration of the new functionality.
