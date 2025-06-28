# Docker-Compose Empty Values Fix

## Issue Summary

Previously, runfromyaml would skip docker-compose commands entirely when the `values` block was empty. This was incorrect behavior because many docker-compose operations (like `up`, `down`, `build`, `pull`, etc.) are standalone commands that don't require additional commands to be executed inside containers.

## Problem Description

### Before the Fix
```yaml
cmd:
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
    values: []  # This would cause the entire command to be skipped
```

**Result**: The command was completely skipped with message:
```
# docker-compose command with empty values - skipping execution
```

### After the Fix
```yaml
cmd:
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
    values: []  # Now executes the base docker-compose command
```

**Result**: The base docker-compose command is executed:
```
# docker-compose command with empty values - executing base command only
docker compose -f docker-compose.yml up -d
```

## Technical Changes

### Code Location
File: `pkg/cli/cli.go`
Function: `executeDockerComposeCommand`

### Before (Problematic Code)
```go
func (e *CommandExecutor) executeDockerComposeCommand(cmd *Command) error {
    // Handle empty values gracefully
    if len(cmd.Values) == 0 {
        functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# docker-compose command with empty values - skipping execution")
        return nil  // ❌ This skips the entire command
    }
    
    args := e.buildDockerComposeArgs(cmd)
    // ... rest of the function
}
```

### After (Fixed Code)
```go
func (e *CommandExecutor) executeDockerComposeCommand(cmd *Command) error {
    args := e.buildDockerComposeArgs(cmd)
    
    // If values are empty, execute the docker-compose command without additional commands
    if len(cmd.Values) == 0 {
        functions.PrintSwitch(color.FgYellow, string(e.config.Level), string(e.config.Output), "# docker-compose command with empty values - executing base command only")
        return e.runCommand(args)  // ✅ This executes the base command
    }
    
    // If values are provided, execute additional commands inside containers
    cmds := splitCommands(cmd.Values)
    // ... rest of the function for handling non-empty values
}
```

## Use Cases Enabled

### 1. Service Management
```yaml
cmd:
  # Build all services
  - type: docker-compose
    name: "build-services"
    desc: "Build all services"
    dcoptions: ["-f", "docker-compose.yml"]
    command: build
    values: []

  # Start services in background
  - type: docker-compose
    name: "start-services"
    desc: "Start services in detached mode"
    dcoptions: ["-f", "docker-compose.yml"]
    command: up
    cmdoptions: ["-d"]
    values: []

  # Stop and remove services
  - type: docker-compose
    name: "cleanup"
    desc: "Stop and remove all services"
    dcoptions: ["-f", "docker-compose.yml"]
    command: down
    cmdoptions: ["--volumes", "--remove-orphans"]
    values: []
```

### 2. Information Commands
```yaml
cmd:
  # Check version
  - type: docker-compose
    name: "version-check"
    desc: "Check docker-compose version"
    command: version
    values: []

  # List running projects
  - type: docker-compose
    name: "list-projects"
    desc: "List all compose projects"
    command: ls
    values: []

  # Show service status
  - type: docker-compose
    name: "service-status"
    desc: "Show status of services"
    dcoptions: ["-f", "docker-compose.yml"]
    command: ps
    values: []
```

### 3. Mixed Workflows
```yaml
cmd:
  # Standard docker-compose command (empty values)
  - type: docker-compose
    name: "start-db"
    desc: "Start database service"
    dcoptions: ["-f", "docker-compose.yml"]
    command: up
    cmdoptions: ["-d"]
    service: db
    values: []

  # Docker-compose with container commands (non-empty values)
  - type: docker-compose
    name: "run-migrations"
    desc: "Run database migrations"
    dcoptions: ["-f", "docker-compose.yml"]
    command: run
    cmdoptions: ["--rm"]
    service: web
    values:
      - python manage.py migrate
      - python manage.py collectstatic --noinput
```

## Testing Results

### Test Configuration
```yaml
logging:
  - level: info
  - output: stdout

cmd:
  - type: docker-compose
    name: "check-version"
    desc: "Check docker-compose version"
    command: version
    values: []

  - type: docker-compose
    name: "show-help"
    desc: "Show docker-compose help"
    command: --help
    values: []

  - type: docker-compose
    name: "list-projects"
    desc: "List compose projects"
    command: ls
    values:  # Completely empty (not even [])
```

### Test Output
```
# docker-compose command with empty values - executing base command only
docker compose version 

Docker Compose version v2.35.1-desktop.1

# docker-compose command with empty values - executing base command only
docker compose --help 

Usage:  docker compose [OPTIONS] COMMAND
[... help output ...]

# docker-compose command with empty values - executing base command only
docker compose ls 

NAME                STATUS              CONFIG FILES
tooling             running(1)          /Users/anatoli.lichii/.tmp/tooling/docker-compose.yaml
```

## Behavior Comparison

| Scenario | Before Fix | After Fix |
|----------|------------|-----------|
| `values: []` | ❌ Skipped entirely | ✅ Executes base command |
| `values:` (empty) | ❌ Skipped entirely | ✅ Executes base command |
| `values: ["cmd1", "cmd2"]` | ✅ Executed with commands | ✅ Executed with commands (unchanged) |

## Impact on Other Command Types

This fix is specific to docker-compose commands. Other command types maintain their existing behavior:

- **exec**: Still skips when values are empty (requires commands to execute)
- **shell**: Still skips when values are empty (requires commands to execute)
- **docker**: Still skips when values are empty (requires commands to execute)
- **ssh**: Still skips when values are empty (requires commands to execute)
- **conf**: Still skips when both data and destination are empty

## Backward Compatibility

This change is fully backward compatible:
- Existing YAML files with non-empty values continue to work unchanged
- No breaking changes to command syntax or structure
- All existing functionality is preserved
- Only affects the previously broken case of empty values

## Documentation Updates

Updated the following documentation files:
- `docs/EMPTY_VALUES_SUPPORT.md` - Added docker-compose specific behavior
- `examples/docker-compose-empty-values.yaml` - Added working examples

## Benefits

1. **Proper Docker-Compose Support**: Can now use standard docker-compose commands without workarounds
2. **Simplified Workflows**: No need to add dummy commands just to make docker-compose blocks work
3. **Better Documentation**: YAML files can serve as proper documentation for container workflows
4. **Consistency**: Behavior now matches user expectations for docker-compose operations
5. **Enhanced Automation**: Enables proper CI/CD workflows with docker-compose

## Example Real-World Usage

```yaml
# Complete application deployment workflow
cmd:
  # Pull latest images
  - type: docker-compose
    name: "pull-images"
    desc: "Pull latest service images"
    dcoptions: ["-f", "docker-compose.prod.yml"]
    command: pull
    values: []

  # Build custom services
  - type: docker-compose
    name: "build-services"
    desc: "Build custom application services"
    dcoptions: ["-f", "docker-compose.prod.yml"]
    command: build
    cmdoptions: ["--no-cache"]
    values: []

  # Start infrastructure services
  - type: docker-compose
    name: "start-infrastructure"
    desc: "Start database and cache services"
    dcoptions: ["-f", "docker-compose.prod.yml"]
    command: up
    cmdoptions: ["-d"]
    service: db redis
    values: []

  # Run database migrations
  - type: docker-compose
    name: "migrate-database"
    desc: "Run database migrations"
    dcoptions: ["-f", "docker-compose.prod.yml"]
    command: run
    cmdoptions: ["--rm"]
    service: web
    values:
      - python manage.py migrate
      - python manage.py check --deploy

  # Start application services
  - type: docker-compose
    name: "start-application"
    desc: "Start all application services"
    dcoptions: ["-f", "docker-compose.prod.yml"]
    command: up
    cmdoptions: ["-d"]
    values: []
```

This fix makes runfromyaml much more practical for real-world container orchestration scenarios.
