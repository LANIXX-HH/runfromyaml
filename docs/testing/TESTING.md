# Testing Guide for runfromyaml

This document describes the testing setup and how to run tests for the runfromyaml project.

## Test Structure

The project now includes comprehensive tests for all major components:

### Test Files

- **`main_test.go`** - Tests for main application logic
- **`pkg/config/config_test.go`** - Tests for configuration management
- **`pkg/config/yaml_test.go`** - Tests for YAML configuration loading
- **`pkg/cli/cli_test.go`** - Tests for CLI functionality
- **`pkg/functions/functions_test.go`** - Tests for utility functions
- **`pkg/errors/errors_test.go`** - Tests for error handling (already existed)

### Test Configuration Files

- **`test_commands.yaml`** - Sample YAML configuration for integration testing
- **`run_tests.sh`** - Comprehensive test runner script

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./pkg/config/
go test -v ./pkg/cli/
go test -v ./pkg/functions/
go test -v ./pkg/errors/

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make benchmark
```

### Advanced Test Commands

```bash
# Run full test suite with coverage and race detection
make test-full

# Run tests for specific packages
make test-config
make test-cli
make test-functions
make test-errors

# Run comprehensive test script
./run_tests.sh
```

## Test Coverage

Current test coverage by package:

- **config**: ~91% - Excellent coverage of configuration management
- **errors**: ~39% - Good coverage of error handling
- **main**: ~11% - Basic coverage of main application logic
- **functions**: ~8% - Basic coverage of utility functions
- **cli**: ~3% - Basic coverage of CLI functionality

### Improving Coverage

To improve test coverage, consider adding tests for:

1. **CLI Package**: Command execution, YAML parsing, environment handling
2. **Functions Package**: File operations, logging, template processing
3. **Main Package**: Application workflow, error handling, integration scenarios
4. **Missing Packages**: Docker, OpenAI, REST API functionality

## Test Categories

### Unit Tests

- **Configuration Tests**: Validate config parsing, flag handling, YAML loading
- **Function Tests**: Test utility functions like file operations, string manipulation
- **Error Tests**: Validate error creation, formatting, and handling
- **Environment Tests**: Test environment variable management

### Integration Tests

- **Main Workflow Test**: Tests the complete application workflow
- **YAML Processing Test**: Tests loading and processing of YAML configuration files
- **File Operations Test**: Tests file creation, reading, and permission handling

### Benchmark Tests

- **Performance Tests**: Measure performance of critical operations
- **Memory Tests**: Check memory usage patterns
- **Concurrency Tests**: Validate thread safety with race detection

## Test Data

### Test YAML Configuration

The `test_commands.yaml` file contains various command types for testing:

- **exec commands**: Simple command execution
- **shell commands**: Shell script execution
- **config commands**: Configuration file creation
- **docker commands**: Docker container operations
- **ssh commands**: Remote command execution

### Environment Variables

Tests use controlled environment variables to ensure consistent behavior:

- `TEST_VAR`: General test variable
- `HOME`: Test home directory
- `USER`: Test user name

## Continuous Integration

### GitHub Actions

The project includes a comprehensive CI workflow (`.github/workflows/test.yml`) that:

- Tests on multiple Go versions (1.19, 1.20, 1.21)
- Tests on multiple operating systems (Ubuntu, Windows, macOS)
- Runs formatting, vetting, and linting checks
- Generates coverage reports
- Uploads coverage to Codecov

### Local CI Simulation

```bash
# Run the same checks as CI
make quality

# Individual quality checks
make fmt    # Format code
make vet    # Run go vet
make lint   # Run golangci-lint (requires golangci-lint installation)
```

## Writing New Tests

### Test Naming Convention

- Test functions: `TestFunctionName`
- Benchmark functions: `BenchmarkFunctionName`
- Example functions: `ExampleFunctionName`

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "descriptive test case name",
            input:    InputType{...},
            expected: OutputType{...},
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionToTest(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionToTest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if result != tt.expected {
                t.Errorf("FunctionToTest() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Best Practices

1. **Use table-driven tests** for multiple test cases
2. **Test both success and failure scenarios**
3. **Use descriptive test names** that explain what is being tested
4. **Clean up resources** (files, environment variables) after tests
5. **Use temporary directories** for file operations
6. **Mock external dependencies** when possible
7. **Test edge cases** and boundary conditions

## Troubleshooting

### Common Issues

1. **Import Errors**: Ensure all required packages are imported
2. **File Permission Issues**: Use appropriate permissions for test files
3. **Environment Variable Conflicts**: Clean up environment variables after tests
4. **Race Conditions**: Use proper synchronization in concurrent tests

### Debug Tests

```bash
# Run tests with debug output
go test -v -debug ./...

# Run specific test with debug
go test -v -run TestSpecificFunction ./pkg/config/

# Run tests with race detection
go test -race ./...
```

## Future Improvements

### Planned Test Enhancements

1. **Add tests for missing packages**: docker, openai, restapi
2. **Increase coverage** for existing packages
3. **Add integration tests** for complete workflows
4. **Add performance benchmarks** for critical paths
5. **Add property-based tests** for complex logic
6. **Add mutation testing** to validate test quality

### Test Infrastructure

1. **Test containers** for Docker functionality testing
2. **Mock servers** for API testing
3. **Test fixtures** for complex YAML configurations
4. **Parallel test execution** for faster feedback

## Contributing

When contributing new code:

1. **Write tests** for new functionality
2. **Update existing tests** when modifying behavior
3. **Ensure tests pass** before submitting PRs
4. **Maintain or improve** test coverage
5. **Follow testing conventions** established in the project

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TestComments)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify Framework](https://github.com/stretchr/testify) (if we decide to use it)
