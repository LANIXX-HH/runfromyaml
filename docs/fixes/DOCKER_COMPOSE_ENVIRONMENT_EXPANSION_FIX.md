# Docker Compose Environment Variable Expansion Fix

## Problem Description

The `docker-compose` command type in runfromyaml was not properly expanding environment variables in the `dcoptions`, `cmdoptions`, and other configuration fields, even when `expandenv: true` was explicitly set. This resulted in errors like:

```
unknown flag: --project-directory $HOME/.tmp/tooling
```

Additionally, the command arguments were not being properly split, causing multi-part options to be treated as single arguments.

## Root Cause Analysis

The issue was identified in the `buildDockerComposeArgs` function in `pkg/cli/cli.go`. Two main problems were discovered:

### 1. Missing Environment Variable Expansion

The function was not checking for the `expandenv` option and applying `os.ExpandEnv()` to the configuration strings. This meant that variables like `$HOME` remained unexpanded.

### 2. Incorrect Argument Handling

Options like `-f /path/to/file` were being treated as single arguments instead of being split into separate arguments (`-f` and `/path/to/file`), which caused `docker compose` to misinterpret the command structure.

## Technical Details

### Before the Fix

```go
// Handle dcoptions - BROKEN
if opts, exists := cmd.Options["dcoptions"]; exists {
    if optsSlice, ok := opts.([]interface{}); ok {
        for _, opt := range optsSlice {
            if strOpt, ok := opt.(string); ok {
                args = append(args, strOpt) // No expansion, no splitting
            }
        }
    }
}
```

This resulted in arguments like:
- `args[2] = '-f $HOME/.tmp/tooling/docker-compose.yaml'` (unexpanded)
- `args[3] = '--project-directory $HOME/.tmp/tooling'` (unexpanded)

### After the Fix

```go
// Handle dcoptions - FIXED
if opts, exists := cmd.Options["dcoptions"]; exists {
    if optsSlice, ok := opts.([]interface{}); ok {
        for _, opt := range optsSlice {
            if strOpt, ok := opt.(string); ok {
                if expandenv {
                    strOpt = os.ExpandEnv(strOpt) // Expand environment variables
                }
                // Split the option into separate arguments
                optArgs := strings.Fields(strOpt)
                args = append(args, optArgs...)
            }
        }
    }
}
```

This results in properly formatted arguments:
- `args[2] = '-f'`
- `args[3] = '/Users/anatoli.lichii/.tmp/tooling/docker-compose.yaml'`
- `args[4] = '--project-directory'`
- `args[5] = '/Users/anatoli.lichii/.tmp/tooling'`

## Implementation

### Files Modified

- `pkg/cli/cli.go`: Updated `buildDockerComposeArgs` function

### Changes Made

1. **Added Environment Variable Expansion Check**:
   ```go
   // Check if environment expansion is enabled
   expandenv := false
   if expandenvOpt, exists := cmd.Options["expandenv"]; exists {
       expandenv = expandenvOpt.(bool)
   }
   ```

2. **Applied Expansion to All Option Types**:
   - `dcoptions`: Docker Compose global options
   - `cmdoptions`: Command-specific options
   - `command`: The docker-compose command itself
   - `service`: Service name

3. **Added Argument Splitting**:
   ```go
   // Split the option into separate arguments
   optArgs := strings.Fields(strOpt)
   args = append(args, optArgs...)
   ```

## Testing

### Test Case

YAML configuration:
```yaml
- type: "docker-compose"
  expandenv: true
  name: "build"
  desc: "build tooling container"
  dcoptions:
    - -f $HOME/.tmp/tooling/docker-compose.yaml
    - --project-directory $HOME/.tmp/tooling
  cmdoptions: []
  command: build
  service: ""
  values: []
```

### Before Fix
```bash
$ ./runfromyaml --file commands.yaml --debug
unknown flag: --project-directory $HOME/.tmp/tooling
Error: exit status 1
```

### After Fix
```bash
$ ./runfromyaml --file commands.yaml --debug
docker compose -f /Users/anatoli.lichii/.tmp/tooling/docker-compose.yaml --project-directory /Users/anatoli.lichii/.tmp/tooling build
Compose can now delegate builds to bake for better performance.
...
tooling  Built
```

## Impact

This fix resolves the environment variable expansion issue for all docker-compose command types and ensures proper argument handling. Users can now:

1. Use environment variables in `dcoptions`, `cmdoptions`, and other fields
2. Rely on `expandenv: true` to work correctly for docker-compose commands
3. Use complex command-line options without worrying about argument parsing issues

## Backward Compatibility

This fix is fully backward compatible. Existing YAML configurations will continue to work as before, but now with proper environment variable expansion when `expandenv: true` is set.

## Related Issues

- Environment variables not expanding in docker-compose configurations
- "unknown flag" errors when using paths with environment variables
- Argument parsing issues with multi-part command-line options

## Version

This fix was implemented in version 0.0.1+ and is available in the current build.
