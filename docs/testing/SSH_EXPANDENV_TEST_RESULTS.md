# SSH expandenv Test Results

## Test Summary

The `expandenv` functionality for SSH blocks in runfromyaml works as follows:

### When `expandenv: true`
- Environment variables in SSH configuration fields (`user`, `host`, `port`) are expanded
- Environment variables in the `values` commands are expanded
- Both custom environment variables (defined in `env` block) and system variables (like `$USER`) are expanded

### When `expandenv: false`
- Environment variables are NOT expanded and remain as literal strings
- This applies to both SSH configuration fields and command values

### Test Results from ssh_only_test.yaml

#### SSH Block with `expandenv: true`
```yaml
- type: "ssh"
  expandenv: true
  user: $SSH_USER        # Expanded to: testuser
  host: $SSH_HOST        # Expanded to: localhost
  values:
    - echo "USER: $USER"      # $USER expanded to: anatoli.lichii
    - echo "SSH_USER: $SSH_USER"  # $SSH_USER expanded to: testuser
```

**Generated SSH Command:**
```bash
ssh -p 22 -l testuser localhost -o ConnectTimeout=2 -o StrictHostKeyChecking=no echo "SSH expandenv=true test" map[echo "USER:anatoli.lichii"] map[echo "SSH_USER:testuser"]
```

#### SSH Block with `expandenv: false`
```yaml
- type: "ssh"
  expandenv: false
  user: $SSH_USER        # Remains literal: $SSH_USER
  host: $SSH_HOST        # Remains literal: $SSH_HOST
  values:
    - echo "USER: $USER"      # Remains literal: $USER
    - echo "SSH_USER: $SSH_USER"  # Remains literal: $SSH_USER
```

## Environment Variables Used in Test

```yaml
env:
  - key: "SSH_USER"
    value: "testuser"
  - key: "SSH_HOST"
    value: "localhost"
```

## Key Observations

1. **Field Expansion**: The `user` and `host` fields are properly expanded when `expandenv: true`
2. **Command Expansion**: Variables in the `values` commands are expanded
3. **System Variables**: Built-in variables like `$USER` are also expanded
4. **Custom Variables**: Variables defined in the `env` block are expanded
5. **Literal Preservation**: When `expandenv: false`, all variables remain as literal strings

## Practical Use Cases

### Use `expandenv: true` when:
- You want to use environment variables for SSH connection parameters
- You need dynamic SSH commands based on environment
- You want to reuse configuration across different environments

### Use `expandenv: false` when:
- You want to pass literal dollar signs to the remote system
- The remote system should handle the variable expansion
- You need to preserve shell variables for remote execution

## Example Working Configuration

```yaml
env:
  - key: "REMOTE_USER"
    value: "deploy"
  - key: "REMOTE_HOST"
    value: "production.example.com"

cmd:
  - type: "ssh"
    expandenv: true
    name: "deploy-with-expandenv"
    desc: "Deploy using environment variables"
    user: $REMOTE_USER
    host: $REMOTE_HOST
    port: 22
    options:
      - -i $HOME/.ssh/deploy_key
    values:
      - echo "Deploying as user: $REMOTE_USER"
      - cd /var/www/app
      - git pull origin main
```

This would generate:
```bash
ssh -p 22 -l deploy production.example.com -i /Users/anatoli.lichii/.ssh/deploy_key [commands...]
```
