#!/bin/bash

set -e

echo "🔧 Setting up pre-commit for runfromyaml..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Installing pre-commit..."

    # Try different installation methods based on OS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install pre-commit
        else
            pip3 install pre-commit
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y pre-commit
        elif command -v yum &> /dev/null; then
            sudo yum install -y pre-commit
        else
            pip3 install pre-commit
        fi
    else
        # Fallback to pip
        pip3 install pre-commit
    fi
fi

# Install golangci-lint if not present
if ! command -v golangci-lint &> /dev/null; then
    echo "📦 Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
fi

# Install the git hook scripts
echo "🔗 Installing pre-commit hooks..."
pre-commit install

# Run against all files to test
echo "🧪 Testing pre-commit setup..."
pre-commit run --all-files

echo "✅ Pre-commit setup complete!"
echo "💡 Hooks will now run automatically on git commit"
