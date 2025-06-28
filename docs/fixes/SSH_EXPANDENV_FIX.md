# SSH expandenv Fix Documentation

## Problem

The `expandenv` functionality for SSH commands was not working correctly for the `options` array. Environment variables in SSH options (like `-i $HOME/.ssh/key`) were not being expanded even when `expandenv: true` was set.

## Root Cause

In `pkg/cli/cli.go`, the `buildSSHArgs` function was applying `os.ExpandEnv()` to the `user` and `host` fields when `expandenv: true`, but it was not applying the same expansion to the elements in the `options` array.

### Before Fix (Lines 353-362)
```go
// Handle SSH options
if opts, exists := cmd.Options["options"]; exists {
    if optsSlice, ok := opts.([]interface{}); ok {
        for _, opt := range optsSlice {
            if strOpt, ok := opt.(string); ok {
                args = append(args, strOpt)  // No expansion applied
            }
        }
    }
}
```

### After Fix
```go
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
```

## Test Results

### Before Fix
```bash
ssh -p 22 -l anatoli.lichii localhost -i $HOME/.ssh/id_rsa-localhost uname -a
Warning: Identity file  $HOME/.ssh/id_rsa-localhost not accessible: No such file or directory.
```

### After Fix
```bash
ssh -p 22 -l anatoli.lichii localhost -i /Users/anatoli.lichii/.ssh/id_rsa-localhost uname -a
Warning: Identity file  /Users/anatoli.lichii/.ssh/id_rsa-localhost not accessible: No such file or directory.
```

## Impact

- ✅ Environment variables in SSH options are now properly expanded when `expandenv: true`
- ✅ Backward compatibility maintained - `expandenv: false` still works as expected
- ✅ All existing functionality for `user`, `host`, and `values` expansion remains unchanged
- ✅ Unit tests added to prevent regression

## Example Usage

```yaml
env:
  - key: "SSH_KEY_PATH"
    value: "$HOME/.ssh/deploy_key"
  - key: "SSH_USER"
    value: "deploy"

cmd:
  - type: "ssh"
    expandenv: true
    name: "deploy-with-env-vars"
    desc: "Deploy using environment variables in SSH options"
    user: $SSH_USER
    host: production.example.com
    port: 22
    options:
      - -i $SSH_KEY_PATH                    # Now properly expanded!
      - -o ConnectTimeout=30
      - -o StrictHostKeyChecking=no
    values:
      - cd /var/www/app
      - git pull origin main
```

This will now correctly generate:
```bash
ssh -p 22 -l deploy production.example.com -i /Users/username/.ssh/deploy_key -o ConnectTimeout=30 -o StrictHostKeyChecking=no [commands...]
```

## Files Changed

- `pkg/cli/cli.go` - Added environment variable expansion for SSH options
- `pkg/cli/ssh_expandenv_test.go` - Added unit tests for the fix

## Tests Added

- `TestSSHExpandenvOptions` - Tests that SSH options are expanded when `expandenv: true`
- `TestSSHExpandenvUserHost` - Tests that user/host expansion still works correctly
