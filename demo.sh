#!/bin/bash

echo "🚀 LogAid Real-Time Demo"
echo "========================"
echo ""

echo "Testing various command errors that LogAid would catch:"
echo ""

echo "1. Testing 'ls ltr' (missing dash):"
ls ltr 2>&1 | head -1

echo ""
echo "2. Testing 'npm d' (wrong command):"
npm d 2>&1 | head -1

echo ""
echo "3. Testing 'git stauts' (typo):"
git stauts 2>&1 | head -1

echo ""
echo "4. Testing 'docker ps -al' (wrong flag):"
docker ps -al 2>&1 | head -1

echo ""
echo "5. Testing 'kubectl get pod' (missing s):"
kubectl get pod 2>&1 | head -1

echo ""
echo "🔍 LogAid would analyze each of these errors in real-time"
echo "   and provide intelligent suggestions to fix them!"
echo ""
echo "💡 Key Benefits:"
echo "   - Catches typos immediately"
echo "   - Suggests correct commands"
echo "   - Maintains your flow state"
echo "   - Works with any CLI tool"
