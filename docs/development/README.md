# Development Documentation

This directory contains technical documentation for developers working on runfromyaml.

## 📋 Contents

### [ARCHITECTURE.md](ARCHITECTURE.md)
**System Architecture and Design**
- Overall system design
- Component relationships
- Code organization
- Design patterns used
- Module structure

### [ERROR_HANDLING.md](ERROR_HANDLING.md)
**Error Handling System**
- Error handling strategies
- Error types and categories
- Error formatting and reporting
- Best practices for error handling
- Validation and recovery patterns

### [ERROR_HANDLING_IMPROVEMENTS.md](ERROR_HANDLING_IMPROVEMENTS.md)
**Error Handling Enhancements**
- Recent improvements to error handling
- New error types and features
- Enhanced error reporting
- Validation improvements
- Implementation details

## 🏗️ Architecture Overview

The runfromyaml project follows a modular architecture with clear separation of concerns:

- **Configuration Management** (`pkg/config/`)
- **Command Line Interface** (`pkg/cli/`)
- **Error Handling** (`pkg/errors/`)
- **Utility Functions** (`pkg/functions/`)
- **Docker Integration** (`pkg/docker/`)
- **REST API** (`pkg/restapi/`)
- **AI Integration** (`pkg/openai/`)

## 🔧 Development Guidelines

### Code Organization
- Follow Go package conventions
- Maintain clear module boundaries
- Use dependency injection where appropriate
- Implement proper error handling

### Error Handling
- Use structured error types
- Provide context and suggestions
- Implement proper error chaining
- Follow validation patterns

### Testing
- Write comprehensive tests for new features
- Maintain good test coverage
- Use table-driven tests
- Include integration tests

## 🎯 Key Design Principles

1. **Modularity**: Clear separation of concerns
2. **Testability**: Easy to test components
3. **Extensibility**: Easy to add new features
4. **Maintainability**: Clean, readable code
5. **Error Resilience**: Robust error handling

For detailed technical information, see the individual documentation files in this directory.
