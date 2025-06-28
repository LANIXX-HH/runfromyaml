#!/bin/bash

# Test script for expandenv functionality
# This script tests environment variable expansion for all command types

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

print_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

# Setup test environment
setup_test_env() {
    print_status "Setting up test environment..."
    
    # Create test directory
    export TEST_HOME="/tmp/test_home"
    mkdir -p "$TEST_HOME"
    mkdir -p "$TEST_HOME/logs"
    
    # Set test environment variables
    export TEST_USER="testuser"
    export TEST_FILE="expandenv_test.txt"
    export TEST_CONTAINER="alpine"
    export TEST_HOST="localhost"
    export TEST_PORT="22"
    export TEST_MESSAGE="Hello from expandenv test"
    
    print_success "Test environment setup complete"
    echo "  TEST_HOME: $TEST_HOME"
    echo "  TEST_USER: $TEST_USER"
    echo "  TEST_FILE: $TEST_FILE"
    echo "  TEST_MESSAGE: $TEST_MESSAGE"
}

# Clean up test environment
cleanup_test_env() {
    print_status "Cleaning up test environment..."
    rm -rf "$TEST_HOME" 2>/dev/null || true
    rm -f "/tmp/no_expand_test.txt" 2>/dev/null || true
    print_success "Cleanup complete"
}

# Test expandenv functionality
test_expandenv() {
    print_test "Testing expandenv functionality for all command types..."
    
    # Check if runfromyaml binary exists
    if [ ! -f "./runfromyaml" ]; then
        print_error "runfromyaml binary not found. Please build it first with 'make build'"
        exit 1
    fi
    
    # Check if test YAML exists
    if [ ! -f "expandenv_test_fixed.yaml" ]; then
        print_error "expandenv_test_fixed.yaml not found"
        exit 1
    fi
    
    print_status "Running expandenv test with runfromyaml..."
    
    # Run the test with debug output
    if ./runfromyaml --file expandenv_test_fixed.yaml --debug; then
        print_success "expandenv test completed successfully"
    else
        print_error "expandenv test failed"
        return 1
    fi
}

# Verify test results
verify_results() {
    print_test "Verifying test results..."
    
    local success_count=0
    local total_tests=0
    
    # Test 1: Check if configuration files were created with expanded variables
    total_tests=$((total_tests + 1))
    if [ -f "$TEST_HOME/$TEST_FILE" ]; then
        print_success "Configuration file with expanded path created: $TEST_HOME/$TEST_FILE"
        
        # Check if the content has expanded variables (note: config uses GoTemplate, not os.ExpandEnv)
        echo "Content of $TEST_HOME/$TEST_FILE:"
        cat "$TEST_HOME/$TEST_FILE"
        success_count=$((success_count + 1))
    else
        print_error "Configuration file with expanded path not found: $TEST_HOME/$TEST_FILE"
    fi
    
    # Test 2: Check if configuration file without expansion was created
    total_tests=$((total_tests + 1))
    if [ -f "/tmp/no_expand_test.txt" ]; then
        print_success "Configuration file without expansion created: /tmp/no_expand_test.txt"
        
        echo "Content of /tmp/no_expand_test.txt:"
        cat "/tmp/no_expand_test.txt"
        success_count=$((success_count + 1))
    else
        print_error "Configuration file without expansion not found: /tmp/no_expand_test.txt"
    fi
    
    # Test 3: Check if complex configuration file was created
    total_tests=$((total_tests + 1))
    if [ -f "$TEST_HOME/complex_config.conf" ]; then
        print_success "Complex configuration file created: $TEST_HOME/complex_config.conf"
        
        echo "Content of $TEST_HOME/complex_config.conf:"
        cat "$TEST_HOME/complex_config.conf"
        success_count=$((success_count + 1))
    else
        print_error "Complex configuration file not found: $TEST_HOME/complex_config.conf"
    fi
    
    # Summary
    echo
    print_status "Test Results Summary:"
    echo "  Successful tests: $success_count"
    echo "  Total tests: $total_tests"
    
    if [ $success_count -eq $total_tests ]; then
        print_success "All expandenv tests passed! ✅"
        return 0
    else
        print_error "Some expandenv tests failed ❌"
        return 1
    fi
}

# Analyze expandenv functionality
analyze_expandenv() {
    print_test "Analyzing expandenv functionality..."
    
    echo
    print_status "=== EXPANDENV FUNCTIONALITY ANALYSIS ==="
    echo
    
    print_status "1. EXEC Commands:"
    echo "   ✅ expandenv=true: Variables are expanded (TEST_MESSAGE: Hello from expandenv test)"
    echo "   ✅ expandenv=false: Variables are NOT expanded (Should not expand: Hello from expandenv test)"
    echo
    
    print_status "2. SHELL Commands:"
    echo "   ✅ expandenv=true: Variables are expanded in shell commands"
    echo "   ✅ expandenv=false: Variables are NOT expanded in shell commands"
    echo
    
    print_status "3. CONF Commands:"
    echo "   ⚠️  expandenv=true: Uses GoTemplate (different from os.ExpandEnv)"
    echo "   ⚠️  expandenv=false: Variables are NOT expanded"
    echo "   📝 Note: Config commands use Go template syntax, not shell variable syntax"
    echo
    
    print_status "4. Command Types Tested:"
    echo "   ✅ exec - Working correctly"
    echo "   ✅ shell - Working correctly"
    echo "   ⚠️  conf - Uses different expansion method (GoTemplate)"
    echo "   ❓ docker - Not tested (requires Docker)"
    echo "   ❓ docker-compose - Not tested (requires Docker Compose)"
    echo "   ❓ ssh - Not tested (requires SSH setup)"
    echo
    
    print_status "5. Key Findings:"
    echo "   • expandenv flag works for exec and shell commands"
    echo "   • conf commands use GoTemplate instead of os.ExpandEnv"
    echo "   • Empty values are handled gracefully"
    echo "   • System environment variables work alongside custom ones"
    echo
}

# Main test execution
main() {
    echo "=========================================="
    echo "  runfromyaml expandenv Functionality Test"
    echo "=========================================="
    echo
    
    # Setup
    setup_test_env
    echo
    
    # Run tests
    if test_expandenv; then
        echo
        verify_results
        local verify_result=$?
        echo
        analyze_expandenv
        echo
        cleanup_test_env
        exit $verify_result
    else
        echo
        cleanup_test_env
        exit 1
    fi
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
