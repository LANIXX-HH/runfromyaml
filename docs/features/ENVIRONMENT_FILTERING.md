# Environment Variable Filtering

The environment variable filtering system in runfromyaml intelligently separates system-specific variables from application-relevant ones, ensuring that generated YAML configurations are clean, portable, and focused on the essential environment setup.

## Overview

When using Interactive Shell Mode or processing environment variables, runfromyaml automatically filters out system-specific variables that are not relevant for documentation or reproduction purposes, while preserving variables that are important for application configuration and deployment.

## Filtering Strategy

### System Variables (Excluded)

These variables are automatically filtered out because they are system or session-specific:

#### System Paths and Directories
```bash
HOME                # User home directory
TMPDIR, TMP, TEMP   # Temporary directories
PATH                # System PATH (too system-specific)
LD_LIBRARY_PATH     # Library paths
DYLD_LIBRARY_PATH   # macOS dynamic library paths
PWD, OLDPWD         # Current and previous working directories
```

#### User and Session Information
```bash
USER, USERNAME      # Current user
LOGNAME            # Login name
SHELL              # Current shell
SHLVL              # Shell level
TTY                # Terminal device
SSH_AUTH_SOCK      # SSH authentication socket
SSH_SOCKET_DIR     # SSH socket directory
```

#### Terminal and Display
```bash
TERM               # Terminal type
TERM_PROGRAM       # Terminal program name
TERM_PROGRAM_VERSION # Terminal version
COLORTERM          # Color terminal support
DISPLAY            # X11 display
WARP_*             # Warp terminal specific
```

#### System Internals
```bash
XPC_SERVICE_NAME   # macOS XPC service
XPC_FLAGS          # XPC flags
COMMAND_MODE       # Command mode
LC_CTYPE, LC_ALL   # Locale settings
LANG               # Language settings
__CF_*             # Core Foundation variables
__CFBundleIdentifier # Bundle identifier
_                  # Last command
```

#### Package Managers
```bash
HOMEBREW_*         # Homebrew variables
CONDA_*            # Conda environment variables
INFOPATH           # Info documentation path
```

#### Shell-Specific Variables
```bash
BASH_*             # Bash-specific variables
ZSH_*              # Zsh-specific variables
FISH_*             # Fish shell variables
PS1, PS2, PS3, PS4 # Shell prompts
PROMPT_*           # Prompt configuration
HISTFILE           # History file location
HISTSIZE           # History size
HISTCONTROL        # History control
```

#### Desktop Environment
```bash
XDG_*              # XDG Base Directory Specification
DBUS_*             # D-Bus variables
GNOME_*            # GNOME desktop variables
KDE_*              # KDE desktop variables
QT_*               # Qt framework variables
GTK_*              # GTK variables
```

#### Tool-Specific
```bash
Q_SET_PARENT_CHECK # runfromyaml internal variable
LESS*              # Less pager variables
PAGER              # Default pager
EDITOR, VISUAL     # Default editors
MANPATH            # Manual page path
```

### Relevant Variables (Included)

These variables are preserved because they are important for application configuration:

#### Cloud and Infrastructure
```bash
AWS_*              # Amazon Web Services
  AWS_REGION
  AWS_ACCESS_KEY_ID
  AWS_SECRET_ACCESS_KEY
  AWS_PROFILE
  AWS_DEFAULT_REGION

DOCKER_*           # Docker configuration
  DOCKER_HOST
  DOCKER_TLS_VERIFY
  DOCKER_CERT_PATH

KUBE*, K8S_*       # Kubernetes
  KUBECONFIG
  KUBERNETES_SERVICE_HOST
  K8S_NAMESPACE

HELM_*             # Helm package manager
  HELM_HOME
  HELM_REPOSITORY_CONFIG

TERRAFORM_*        # Terraform
  TERRAFORM_LOG
  TERRAFORM_LOG_PATH
```

#### CI/CD and Version Control
```bash
CI_*               # Continuous Integration
  CI_COMMIT_SHA
  CI_BRANCH
  CI_BUILD_NUMBER

BUILD_*            # Build systems
  BUILD_NUMBER
  BUILD_ID
  BUILD_URL

DEPLOY_*           # Deployment
  DEPLOY_ENV
  DEPLOY_STAGE

JENKINS_*          # Jenkins CI
  JENKINS_URL
  JENKINS_HOME

GIT_*              # Git configuration
  GIT_BRANCH
  GIT_COMMIT
  GIT_AUTHOR_NAME

GITHUB_*           # GitHub
  GITHUB_TOKEN
  GITHUB_REPOSITORY

GITLAB_*           # GitLab
  GITLAB_TOKEN
  GITLAB_CI
```

#### Application and Services
```bash
API_*              # API configuration
  API_KEY
  API_URL
  API_VERSION

DB_*, DATABASE_*   # Database configuration
  DATABASE_URL
  DB_HOST
  DB_PORT
  DB_NAME

REDIS_*            # Redis configuration
  REDIS_URL
  REDIS_HOST

MONGO_*            # MongoDB
  MONGODB_URI
  MONGO_URL

MYSQL_*            # MySQL
  MYSQL_HOST
  MYSQL_PORT
  MYSQL_DATABASE

POSTGRES_*         # PostgreSQL
  POSTGRES_HOST
  POSTGRES_DB
  POSTGRES_USER
```

#### Programming Languages
```bash
NODE_*             # Node.js
  NODE_ENV
  NODE_PATH
  NODE_OPTIONS

PYTHON_*           # Python
  PYTHON_PATH
  PYTHONPATH

JAVA_*             # Java
  JAVA_HOME
  JAVA_OPTS

GO_*               # Go
  GOPATH
  GOPROXY
  GO111MODULE

RUST_*             # Rust
  RUST_LOG
  RUSTUP_HOME

PHP_*              # PHP
  PHP_INI_SCAN_DIR

RUBY_*             # Ruby
  RUBY_VERSION
  BUNDLE_PATH
```

#### Common Application Variables
```bash
PORT               # Application port
HOST               # Application host
DEBUG              # Debug mode
ENVIRONMENT        # Environment name (dev/staging/prod)
STAGE              # Deployment stage
VERSION            # Application version
REGION             # Geographic region
ZONE               # Availability zone
NAMESPACE          # Kubernetes namespace
SERVICE            # Service name
CONFIG             # Configuration file/path
SECRET             # Secret identifier
TOKEN              # Authentication token
KEY                # Generic key
URL                # Service URL
ENDPOINT           # API endpoint
```

## Implementation Details

### Filter Function

The filtering is implemented in two main functions:

```go
func filterRelevantEnvVars(envs map[string]string) map[string]string
func isRelevantEnvVar(key string) bool
```

### Filtering Logic

1. **Exact Match Exclusion**: Variables in the system variables map are excluded
2. **Prefix Exclusion**: Variables starting with excluded prefixes are filtered out
3. **Relevance Check**: Remaining variables are checked against relevant patterns
4. **Prefix Inclusion**: Variables with relevant prefixes are included
5. **Specific Variable Inclusion**: Important individual variables are included

### Example Usage

```go
envs := map[string]string{
    "HOME":       "/home/user",
    "AWS_REGION": "eu-central-1",
    "PATH":       "/usr/bin:/bin",
    "API_KEY":    "secret123",
}

filtered := filterRelevantEnvVars(envs)
// Result: {"AWS_REGION": "eu-central-1", "API_KEY": "secret123"}
```

## Configuration

Currently, the filtering rules are built into the code. Future versions may support:

- Custom filtering rules via configuration files
- User-defined include/exclude patterns
- Environment-specific filtering profiles

## Best Practices

### For Users

1. **Set Relevant Variables Before Recording**: Export application-specific variables before using Interactive Shell Mode
2. **Use Standard Naming Conventions**: Follow common patterns (AWS_*, API_*, DB_*) for automatic inclusion
3. **Review Generated YAML**: Always check that important variables are included

### For Developers

1. **Follow Naming Conventions**: Use recognized prefixes for automatic filtering
2. **Document Custom Variables**: If using non-standard variable names, document their purpose
3. **Test Filtering**: Verify that your variables are correctly classified

## Examples

### Development Environment
```bash
export NODE_ENV=development
export API_KEY=dev-key-123
export DATABASE_URL=postgres://localhost:5432/myapp
export DEBUG=true

runfromyaml --shell
# These variables will be included in the generated YAML
```

### Production Deployment
```bash
export AWS_REGION=eu-central-1
export DOCKER_HOST=tcp://prod-docker:2376
export KUBERNETES_SERVICE_HOST=10.0.0.1
export ENVIRONMENT=production

runfromyaml --shell
# Production-relevant variables are preserved
```

### CI/CD Pipeline
```bash
export CI_COMMIT_SHA=abc123def456
export BUILD_NUMBER=42
export DEPLOY_STAGE=staging
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

runfromyaml --shell
# CI/CD variables are included for reproducibility
```

## Troubleshooting

### Variable Not Included
If an important variable is being filtered out:

1. Check if it matches any exclusion patterns
2. Verify it follows a recognized naming convention
3. Consider renaming to use a standard prefix (API_*, DB_*, etc.)

### System Variable Included
If a system variable is being included:

1. Check if it matches a relevant pattern by mistake
2. Report as a bug if it should be filtered
3. The filtering rules may need adjustment

### Custom Variables
For custom application variables that don't follow standard patterns:

1. Use recognized prefixes when possible (APP_*, ENV_*, etc.)
2. Document the variables and their purpose
3. Consider contributing to the filtering rules if they're commonly used

## Future Enhancements

Planned improvements to the filtering system:

1. **Configurable Rules**: Allow users to define custom filtering rules
2. **Context-Aware Filtering**: Different rules for different contexts (dev/prod)
3. **Interactive Selection**: Allow users to manually select variables during recording
4. **Pattern Learning**: Learn from user selections to improve automatic filtering
5. **Export Profiles**: Predefined filtering profiles for different use cases
