# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [5.0.0] - 2025-06-28

### 🎉 Major Release - Comprehensive Modernization and Testing Infrastructure

This major release represents a significant modernization of the runfromyaml project with comprehensive testing infrastructure, improved documentation, enhanced error handling, and better cross-platform compatibility.

### ✨ Added

#### **New Features**
- **AI Integration Modernization**: Updated OpenAI integration from deprecated text-davinci-003 to modern gpt-3.5-turbo model
- **Interactive Shell Mode**: Added shell-to-YAML reverse parsing functionality for interactive command recording
- **Enhanced Empty Values Support**: Allow empty values blocks and command blocks for documentation and placeholders
- **Comprehensive Error Handling System**: New custom error types with validation and recovery mechanisms
- **Cross-Platform Compatibility**: Improved Windows, macOS, and Linux support with platform-specific handling

#### **Documentation & Testing**
- **Structured Documentation**: Complete reorganization into docs/ directory with categories:
  - `docs/development/` - Architecture and error handling guides
  - `docs/features/` - Feature documentation and usage examples
  - `docs/fixes/` - Bug fixes and improvement documentation
  - `docs/summaries/` - Project summaries and implementation overviews
  - `docs/testing/` - Testing guides, setup, and results
- **Comprehensive Testing Infrastructure**:
  - Unit tests for CLI, config, functions, and YAML handling
  - Integration tests for SSH expandenv functionality
  - Benchmark tests for performance monitoring
  - GitHub Actions workflow for automated CI/CD
  - golangci-lint configuration for code quality enforcement
- **Enhanced Examples**: Extensive example configurations for various use cases including Docker Compose and SSH scenarios

### 🔧 Improved

#### **Core Functionality**
- **Enhanced SSH Support**: Improved expandenv functionality with proper environment variable expansion
- **Better Docker Compose Handling**: Enhanced environment expansion and empty values support
- **Improved CLI**: Better debugging, logging, and error reporting throughout the application
- **Modern API Usage**: Updated deprecated Docker client API calls and other external dependencies

#### **Code Quality**
- **Error Handling**: Comprehensive error handling patterns throughout codebase
- **Code Structure**: Better separation of concerns with modular architecture
- **Cross-Platform Support**: Platform-specific handling for file permissions and temporary directories
- **Performance**: Optimized file operations and memory usage

### 🐛 Fixed

#### **Critical Fixes**
- **GitHub Actions CI/CD**: Resolved all linting issues and test failures across platforms
- **Windows Compatibility**: Fixed file permission handling and temporary directory cleanup on Windows
- **Benchmark Tests**: Fixed invalid filename generation in performance tests
- **Memory Leaks**: Proper cleanup of resources and file handles

#### **Code Quality Fixes**
- **Linting Issues**: Resolved all golangci-lint warnings and errors
- **Error Handling**: Fixed unchecked error returns throughout codebase
- **API Deprecations**: Updated deprecated Docker client and OpenAI API calls
- **Test Reliability**: Improved test stability across different platforms and environments

### 🔄 Changed

#### **Breaking Changes**
- **OpenAI API**: Migrated from Engines API to Chat Completions API (maintains backward compatibility)
- **Error Handling**: New error types may affect custom error handling implementations

#### **Configuration Changes**
- **Default AI Model**: Changed from `text-davinci-003` to `gpt-3.5-turbo`
- **Linting Configuration**: Updated golangci-lint settings for better code quality
- **CI/CD Pipeline**: Simplified GitHub Actions workflow (removed Go version matrix testing)

### 🏗️ Infrastructure

#### **Development Experience**
- **Simplified CI/CD**: Removed matrix testing, now uses Go 1.21 consistently across all platforms
- **Better Tooling**: Enhanced development tools and scripts for testing and building
- **Documentation**: Comprehensive guides for development, testing, and contribution

#### **Platform Support**
- **Windows**: Full compatibility with Windows-specific file system behavior
- **macOS**: Native support for macOS development and deployment
- **Linux**: Enhanced Linux support with proper permission handling

### 📊 Statistics
- **Commits**: 8 major commits with comprehensive improvements
- **Files Changed**: 50+ files updated across the entire codebase
- **Test Coverage**: Significantly improved with comprehensive test suite
- **Documentation**: 10+ new documentation files added
- **Code Quality**: All linting issues resolved, modern Go practices implemented

### 🚀 Migration Guide

#### **For Users**
- No breaking changes for end users
- New AI features available with updated models
- Enhanced error messages provide better debugging information

#### **For Developers**
- Update any custom error handling to use new error types
- Review new documentation structure for development guidelines
- Use new testing infrastructure for contributions

### 🙏 Acknowledgments
This release represents months of work to modernize and improve the runfromyaml project. Special thanks to the Go community for excellent tooling and best practices that made this comprehensive update possible.

---

## [4.5.2] - 2023-02-26

### Previous Release
For changes prior to v5.0.0, please refer to the git history or previous release notes.

---

## Links
- [GitHub Repository](https://github.com/lanixx-hh/runfromyaml)
- [Release Downloads](https://github.com/lanixx-hh/runfromyaml/releases)
- [Documentation](docs/README.md)
- [Contributing Guidelines](docs/development/ARCHITECTURE.md)
