#!/bin/bash

echo "=== Testing SSH expandenv functionality ==="
echo "Current USER: $USER"
echo "Current HOME: $HOME"
echo ""

echo "=== Running SSH expandenv test ==="
./runfromyaml --file ssh_expandenv_test.yaml --debug

echo ""
echo "=== Test completed ==="
