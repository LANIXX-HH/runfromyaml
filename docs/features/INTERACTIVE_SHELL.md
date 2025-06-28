# Interactive Shell Mode

The Interactive Shell Mode allows you to execute commands in real-time while automatically recording them for YAML generation. This feature bridges the gap between interactive command execution and documentation generation.

## Overview

Interactive Shell Mode provides:
- **Real-time command execution** - See output immediately like a normal shell
- **Automatic command recording** - Commands are captured for YAML generation
- **Smart environment filtering** - Only relevant environment variables are included
- **Multiple shell support** - Works with bash, zsh, sh, fish, PowerShell, and cmd

## Usage

### Basic Usage

```bash
runfromyaml --shell
```

### Specify Shell Type

```bash
runfromyaml --shell --shell-type zsh
```

### Supported Shell Types

- `bash` (default on Unix-like systems)
- `zsh`
- `sh`
- `fish`
- `powershell` (Windows)
- `cmd` (Windows)

## How It Works

1. **Start Interactive Session**: The tool starts an interactive shell session
2. **Execute Commands**: Enter commands as you would in a normal shell
3. **See Real Output**: Command output is displayed immediately
4. **Record for YAML**: Commands (not output) are recorded for documentation
5. **Generate YAML**: Type `exit` to finish and generate the YAML structure

## Example Session

```bash
$ runfromyaml --shell
🐚 Interactive Shell Mode
Your input commands will be recorded to generate YAML structure
Enter 'exit' to stop recording and generate YAML

🐚 Interactive Shell Mode (bash)
Commands will be executed AND recorded for YAML generation
Enter commands (type 'exit' to finish):

> ls -la
total 1024
drwxr-xr-x  10 user  staff   320 Jun 28 22:00 .
drwxr-xr-x  20 user  staff   640 Jun 28 21:00 ..
-rw-r--r--   1 user  staff  1234 Jun 28 22:00 README.md
...

> pwd
/Users/user/project

> echo "Hello World"
Hello World

> exit

📄 Generated YAML:
---
cmd:
- desc: 'Interactive command: ls -la'
  expandenv: true
  name: command-1
  type: shell
  values:
  - ls -la
- desc: 'Interactive command: pwd'
  expandenv: true
  name: command-2
  type: shell
  values:
  - pwd
- desc: 'Interactive command: echo "Hello World"'
  expandenv: true
  name: command-3
  type: shell
  values:
  - echo "Hello World"
logging:
- level: info
- output: stdout
```

## Environment Variable Filtering

The Interactive Shell Mode includes intelligent environment variable filtering to keep the generated YAML clean and relevant.

### System Variables (Filtered Out)

These variables are automatically excluded from the generated YAML:

**System Paths & Directories:**
- `HOME`, `TMPDIR`, `TMP`, `TEMP`
- `PATH`, `LD_LIBRARY_PATH`, `DYLD_LIBRARY_PATH`
- `PWD`, `OLDPWD`

**User & Session Info:**
- `USER`, `USERNAME`, `LOGNAME`
- `SHELL`, `SHLVL`, `TTY`
- `SSH_AUTH_SOCK`, `SSH_SOCKET_DIR`

**Terminal & Display:**
- `TERM`, `TERM_PROGRAM`, `TERM_PROGRAM_VERSION`
- `COLORTERM`, `DISPLAY`

**Package Managers:**
- `HOMEBREW_*`, `CONDA_*`, `INFOPATH`

**Shell-Specific:**
- `BASH_*`, `ZSH_*`, `FISH_*`
- `PS1`, `PS2`, `PS3`, `PS4`
- `HISTFILE`, `HISTSIZE`, `HISTCONTROL`

### Relevant Variables (Included)

These variables are preserved in the generated YAML:

**Cloud & Infrastructure:**
- `AWS_*` - AWS configuration
- `DOCKER_*` - Docker settings
- `KUBE*`, `K8S_*` - Kubernetes configuration
- `HELM_*` - Helm settings
- `TERRAFORM_*` - Terraform variables

**CI/CD & Development:**
- `CI_*`, `BUILD_*`, `DEPLOY_*`
- `JENKINS_*`, `GITHUB_*`, `GITLAB_*`
- `GIT_*` - Git configuration

**Application & Services:**
- `API_*`, `DB_*`, `DATABASE_*`
- `REDIS_*`, `MONGO_*`, `MYSQL_*`, `POSTGRES_*`
- `NODE_*`, `PYTHON_*`, `JAVA_*`, `GO_*`

**Common Application Variables:**
- `PORT`, `HOST`, `DEBUG`, `ENVIRONMENT`
- `VERSION`, `REGION`, `ZONE`, `NAMESPACE`
- `TOKEN`, `KEY`, `URL`, `ENDPOINT`

## Generated YAML Structure

The Interactive Shell Mode generates a complete YAML structure with:

### Command Blocks
```yaml
cmd:
- desc: 'Interactive command: <command>'
  expandenv: true
  name: command-<number>
  type: shell
  values:
  - <command>
```

### Environment Variables (Filtered)
```yaml
env:
- key: AWS_REGION
  value: eu-central-1
- key: API_KEY
  value: secret123
```

### Logging Configuration
```yaml
logging:
- level: info
- output: stdout
```

## Use Cases

### 1. Documentation Generation
Record a series of setup commands and automatically generate documentation:

```bash
runfromyaml --shell
> git clone https://github.com/user/repo.git
> cd repo
> npm install
> npm run build
> exit
```

### 2. Environment Setup Recording
Capture environment configuration steps:

```bash
export AWS_REGION=eu-central-1
runfromyaml --shell
> aws configure list
> kubectl get nodes
> docker ps
> exit
```

### 3. Troubleshooting Documentation
Record troubleshooting steps for future reference:

```bash
runfromyaml --shell
> systemctl status nginx
> tail -f /var/log/nginx/error.log
> nginx -t
> exit
```

## Best Practices

1. **Set Environment Variables First**: Export relevant environment variables before starting the interactive session
2. **Use Descriptive Commands**: Commands will be used as documentation, so make them clear
3. **Group Related Commands**: Start a new session for different logical groups of commands
4. **Review Generated YAML**: Always review the generated YAML before using it in production

## Limitations

1. **No Command Editing**: Once a command is executed, it cannot be edited in the current session
2. **No Command History**: Previous commands from shell history are not accessible
3. **Single Session**: Each interactive session is independent
4. **Output Not Captured**: Only commands are recorded, not their output

## Integration with Other Features

The Interactive Shell Mode works seamlessly with other runfromyaml features:

- **File Execution**: Generated YAML can be saved and executed with `runfromyaml --file`
- **REST API**: Generated YAML can be sent to REST API endpoints
- **AI Integration**: Combine with AI features for command suggestions

## Troubleshooting

### Shell Not Found
If you get a "shell not found" error:
```bash
# Check available shells
cat /etc/shells

# Use a different shell
runfromyaml --shell --shell-type sh
```

### Commands Not Executing
If commands don't execute properly:
1. Check shell permissions
2. Verify the shell is properly installed
3. Try with a different shell type

### Environment Variables Missing
If expected environment variables are missing from the YAML:
1. Check if they match the filtering criteria
2. Verify they are set before starting the interactive session
3. Review the filtering rules in the documentation
