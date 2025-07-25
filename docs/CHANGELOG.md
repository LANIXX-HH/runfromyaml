# Changelog

## Version 0.0.1+ (Current Development)

### Bug Fixes

#### Docker-Compose Environment Variable Expansion Fix
- **Fixed**: Environment variable expansion now works correctly in docker-compose commands
- **Problem**: `expandenv: true` was not being applied to `dcoptions`, `cmdoptions`, and other fields
- **Symptoms**: Errors like `unknown flag: --project-directory $HOME/.tmp/tooling`
- **Root Cause**: Missing environment expansion and incorrect argument splitting
- **Solution**: Added proper `os.ExpandEnv()` calls and `strings.Fields()` for argument splitting
- **Impact**: Docker-compose commands with environment variables now work as expected
- **Documentation**: Added comprehensive fix documentation in `DOCKER_COMPOSE_ENVIRONMENT_EXPANSION_FIX.md`

#### Docker-Compose Empty Values Fix
- **Fixed**: Docker-compose commands with empty values now execute properly
- **Before**: Commands with empty `values` blocks were completely skipped
- **After**: Base docker-compose commands execute even with empty values
- **Impact**: Enables proper use of standard docker-compose operations (up, down, build, pull, etc.)
- **Documentation**: Added comprehensive fix documentation in `DOCKER_COMPOSE_EMPTY_VALUES_FIX.md`
- **Examples**: Added `examples/docker-compose-empty-values.yaml` demonstrating the fix

## Version 0.0.1 (Current)

### New Features

#### AI Integration
- Added OpenAI API integration for command generation and assistance
- New command-line flags:
  - `--ai`: Enable AI mode
  - `--ai-key`: Set OpenAI API key
  - `--ai-in`: Input prompt for AI
  - `--ai-model`: Specify OpenAI model (default: text-davinci-003)
  - `--ai-cmdtype`: Specify command type for AI examples (default: shell)

#### Interactive Shell Mode
- New `--shell` flag to enable interactive shell mode
- `--shell-type` flag to specify shell type (default: bash)
- Automatically records commands and generates YAML structure
- Real-time command capture and documentation

#### Enhanced Configuration
- YAML-based configuration support with options block
- Options can be defined directly in YAML files
- Supports all command-line options in YAML format

#### Improved REST API
- Better authentication handling
- Enhanced output management with `--restout` flag
- Improved error handling and response formatting

### Technical Improvements

#### Code Structure
- Modular package structure:
  - `pkg/config`: Configuration management
  - `pkg/openai`: AI integration
  - `pkg/cli`: Command-line interface
  - `pkg/restapi`: REST API functionality
  - `pkg/functions`: Utility functions
  - `pkg/docker`: Docker integration

#### Dependencies
- Updated to Go 1.19
- Added new dependencies:
  - OpenAI API client libraries
  - Enhanced HTTP authentication
  - Improved logging with logrus

### Bug Fixes
- Improved error handling across all modules
- Better validation for YAML configuration
- Enhanced environment variable expansion

### Breaking Changes
- None - backward compatibility maintained

### Migration Notes
- Existing YAML files continue to work without changes
- New options block is optional
- All existing command-line flags remain functional

## Planned for Next Release

### TODO Items
- [ ] Write comprehensive tests
- [ ] Implement connection between blocks (artifacts)
- [ ] Implement dependency management between blocks
- [ ] Update AI model defaults to newer OpenAI models (GPT-4, etc.)
- [ ] Add support for other AI providers (Claude, Gemini, etc.)
- [ ] Improve error handling and validation
- [ ] Add configuration file validation
- [ ] Implement block execution ordering
- [ ] Add support for conditional execution
- [ ] Enhance logging and debugging capabilities
