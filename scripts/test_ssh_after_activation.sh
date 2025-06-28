#!/bin/bash

echo "=== SSH Activation Test ==="
echo ""

echo "1. Checking Remote Login status:"
sudo systemsetup -getremotelogin

echo ""
echo "2. Checking if SSH daemon is running:"
ps aux | grep sshd | grep -v grep

echo ""
echo "3. Checking if port 22 is open:"
nc -z localhost 22 && echo "✅ Port 22 is open" || echo "❌ Port 22 is closed"

echo ""
echo "4. Testing SSH connection:"
ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa-localhost $USER@localhost echo "SSH works!" 2>&1

echo ""
echo "5. Running runfromyaml SSH test:"
cd /Users/anatoli.lichii/Projects/LANIXX/runfromyaml
./runfromyaml --file commands.yaml --debug 2>&1 | tail -5
