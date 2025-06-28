# SSH expandenv Options Fix

**Date:** 2025-06-28  
**Issue:** SSH options array not respecting `expandenv` setting  
**Status:** ✅ Fixed  
**Files Changed:** `pkg/cli/cli.go`, `pkg/cli/ssh_expandenv_test.go`

## Problem Description

The `expandenv` functionality for SSH commands was not working correctly for the `options` array. Environment variables in SSH options (like `-i $HOME/.ssh/key`) were not being expanded even when `expandenv: true` was set.

### Symptoms

```yaml
cmd:
  - type: "ssh"
    expandenv: true
    user: $USER
    host: localhost
    options:
      - -i $HOME/.ssh/id_rsa-localhost  # This was NOT being expanded
    values:
      - echo "test"
```

**Generated SSH command (before fix):**
```bash
ssh -p 22 -l anatoli.lichii localhost -i $HOME/.ssh/id_rsa-localhost echo "test"
Warning: Identity file  $HOME/.ssh/id_rsa-localhost not accessible: No such file or directory.
```

The `$HOME` variable remained literal instead of being expanded to `/Users/anatoli.lichii`.

## Root Cause Analysis

In `pkg/cli/cli.go`, the `buildSSHArgs` function was applying `os.ExpandEnv()` to the `user` and `host` fields when `expandenv: true`, but it was **not** applying the same expansion to the elements in the `options` array.

### Code Analysis

**Before Fix (Lines 353-362):**
```go
// Handle SSH options
if opts, exists := cmd.Options["options"]; exists {
    if optsSlice, ok := opts.([]interface{}); ok {
        for _, opt := range optsSlice {
            if strOpt, ok := opt.(string); ok {
                args = append(args, strOpt)  // ❌ No expansion applied
            }
        }
    }
}
```

The SSH options were being added to the command arguments without checking the `expandenv` setting.

## Solution Implementation

### Code Changes

**After Fix:**
```go
// Handle SSH options
if opts, exists := cmd.Options["options"]; exists {
    if optsSlice, ok := opts.([]interface{}); ok {
        for _, opt := range optsSlice {
            if strOpt, ok := opt.(string); ok {
                // ✅ Apply expandenv to SSH options if enabled
                if expandenv, exists := cmd.Options["expandenv"]; exists && expandenv.(bool) {
                    strOpt = os.ExpandEnv(strOpt)
                }
                args = append(args, strOpt)
            }
        }
    }
}
```

### Logic Flow

1. **Check if options exist**: Verify that the SSH command has an `options` array
2. **Iterate through options**: Loop through each option string
3. **Check expandenv setting**: If `expandenv: true` is set for the command
4. **Apply expansion**: Use `os.ExpandEnv()` to expand environment variables
5. **Add to arguments**: Append the (potentially expanded) option to SSH arguments

## Test Results

### Before Fix
```bash
ssh -p 22 -l anatoli.lichii localhost -i $HOME/.ssh/id_rsa-localhost uname -a
Warning: Identity file  $HOME/.ssh/id_rsa-localhost not accessible: No such file or directory.
```

### After Fix
```bash
ssh -p 22 -l anatoli.lichii localhost -i /Users/anatoli.lichii/.ssh/id_rsa-localhost uname -a
```

✅ The `$HOME` variable is now correctly expanded to `/Users/anatoli.lichii`.

## Unit Tests Added

### Test File: `pkg/cli/ssh_expandenv_test.go`

**Test 1: `TestSSHExpandenvOptions`**
- Tests that SSH options are expanded when `expandenv: true`
- Tests that SSH options remain literal when `expandenv: false`
- Uses environment variable `$TEST_SSH_KEY` to verify expansion

**Test 2: `TestSSHExpandenvUserHost`**
- Tests that user/host expansion still works correctly
- Ensures backward compatibility with existing functionality

### Test Execution
```bash
$ go test ./pkg/cli -v -run TestSSHExpandenv
=== RUN   TestSSHExpandenvOptions
--- PASS: TestSSHExpandenvOptions (0.00s)
=== RUN   TestSSHExpandenvUserHost
--- PASS: TestSSHExpandenvUserHost (0.00s)
PASS
```

## Example Usage

### Working Configuration

```yaml
env:
  - key: "SSH_KEY_PATH"
    value: "$HOME/.ssh/deploy_key"
  - key: "SSH_USER"
    value: "deploy"
  - key: "REMOTE_HOST"
    value: "production.example.com"

cmd:
  - type: "ssh"
    expandenv: true
    name: "deploy-with-env-vars"
    desc: "Deploy using environment variables in SSH options"
    user: $SSH_USER
    host: $REMOTE_HOST
    port: 22
    options:
      - -i $SSH_KEY_PATH                    # ✅ Now properly expanded!
      - -o ConnectTimeout=30
      - -o StrictHostKeyChecking=no
      - -o UserKnownHostsFile=/dev/null
    values:
      - cd /var/www/app
      - git pull origin main
      - sudo systemctl restart nginx
```

### Generated Command
```bash
ssh -p 22 -l deploy production.example.com \
    -i /Users/username/.ssh/deploy_key \
    -o ConnectTimeout=30 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    cd /var/www/app; git pull origin main; sudo systemctl restart nginx
```

## Impact Assessment

### ✅ Benefits
- **Environment variables in SSH options work correctly**: `$HOME`, `$USER`, custom variables
- **Backward compatibility maintained**: `expandenv: false` still works as expected
- **Consistent behavior**: SSH options now behave the same as user/host/values fields
- **Enhanced flexibility**: Enables dynamic SSH configurations across environments

### ✅ No Breaking Changes
- All existing functionality for `user`, `host`, and `values` expansion remains unchanged
- Commands with `expandenv: false` continue to work exactly as before
- No changes to YAML syntax or command structure required

### ✅ Quality Assurance
- Unit tests added to prevent regression
- All existing tests continue to pass
- Code follows existing patterns and conventions

## Related Issues

This fix resolves the issue where the last SSH block in `commands.yaml` was failing due to:
1. ❌ `$HOME` not being expanded in `-i $HOME/.ssh/id_rsa-localhost`
2. ⚠️ SSH daemon not running on localhost (separate configuration issue)

After the fix:
1. ✅ `$HOME` is correctly expanded to `/Users/anatoli.lichii`
2. ⚠️ SSH connection still fails due to daemon not running (expected behavior)

## Files Modified

### `pkg/cli/cli.go`
- **Function:** `buildSSHArgs`
- **Change:** Added environment variable expansion for SSH options array
- **Lines:** ~356-365

### `pkg/cli/ssh_expandenv_test.go` (New File)
- **Tests:** `TestSSHExpandenvOptions`, `TestSSHExpandenvUserHost`
- **Purpose:** Prevent regression and verify fix functionality

## Verification Steps

1. **Build the project:**
   ```bash
   make clean && make
   ```

2. **Test with environment variables:**
   ```yaml
   cmd:
     - type: "ssh"
       expandenv: true
       user: $USER
       host: localhost
       options:
         - -i $HOME/.ssh/test_key
   ```

3. **Verify expansion in debug output:**
   ```bash
   ./runfromyaml --file test.yaml --debug
   ```

4. **Run unit tests:**
   ```bash
   go test ./pkg/cli -v -run TestSSHExpandenv
   ```

## Future Considerations

- Consider adding similar expansion support for other command types if needed
- Monitor for any edge cases with complex environment variable scenarios
- Document best practices for SSH key management in runfromyaml configurations
