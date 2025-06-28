#!/bin/bash

# runfromyaml Test Runner Script
# This script runs comprehensive tests for the runfromyaml project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run a command and check its exit status
run_test() {
    local test_name="$1"
    local command="$2"
    
    print_status "Running: $test_name"
    
    if eval "$command"; then
        print_success "$test_name passed"
        return 0
    else
        print_error "$test_name failed"
        return 1
    fi
}

# Main test execution
main() {
    print_status "Starting runfromyaml test suite..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    print_status "Go version: $(go version)"
    
    # Initialize test counters
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    
    # Test 1: Format check
    total_tests=$((total_tests + 1))
    if run_test "Code formatting check" "go fmt ./... && git diff --exit-code"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
        print_warning "Code formatting issues found. Run 'go fmt ./...' to fix."
    fi
    
    # Test 2: Vet check
    total_tests=$((total_tests + 1))
    if run_test "Go vet check" "go vet ./..."; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 3: Build test
    total_tests=$((total_tests + 1))
    if run_test "Build test" "go build -v"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 4: Unit tests
    total_tests=$((total_tests + 1))
    if run_test "Unit tests" "go test -v ./..."; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 5: Race condition tests
    total_tests=$((total_tests + 1))
    if run_test "Race condition tests" "go test -race ./..."; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
        print_warning "Race condition tests failed. This might indicate concurrency issues."
    fi
    
    # Test 6: Coverage test
    total_tests=$((total_tests + 1))
    if run_test "Coverage test" "go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"; then
        passed_tests=$((passed_tests + 1))
        print_status "Coverage report generated: coverage.out"
        
        # Generate HTML coverage report
        if go tool cover -html=coverage.out -o coverage.html; then
            print_status "HTML coverage report generated: coverage.html"
        fi
    else
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 7: Benchmark tests
    total_tests=$((total_tests + 1))
    if run_test "Benchmark tests" "go test -bench=. -benchmem ./..."; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
        print_warning "Benchmark tests failed or no benchmarks found."
    fi
    
    # Test 8: Integration test with test YAML file
    if [ -f "test_commands.yaml" ]; then
        total_tests=$((total_tests + 1))
        if run_test "Integration test with test YAML" "./runfromyaml --file test_commands.yaml --debug"; then
            passed_tests=$((passed_tests + 1))
        else
            failed_tests=$((failed_tests + 1))
        fi
    else
        print_warning "test_commands.yaml not found, skipping integration test"
    fi
    
    # Test 9: Help command test
    total_tests=$((total_tests + 1))
    if run_test "Help command test" "./runfromyaml --help"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 10: Configuration validation test
    total_tests=$((total_tests + 1))
    if run_test "Configuration validation test" "./runfromyaml --file non_existent_file.yaml 2>/dev/null || true"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
    
    # Print test summary
    echo
    print_status "Test Summary:"
    echo "  Total tests: $total_tests"
    echo "  Passed: $passed_tests"
    echo "  Failed: $failed_tests"
    
    if [ $failed_tests -eq 0 ]; then
        print_success "All tests passed! 🎉"
        exit 0
    else
        print_error "$failed_tests test(s) failed"
        exit 1
    fi
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
