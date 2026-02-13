#!/bin/bash
# Check Go test coverage threshold
# Usage: ./scripts/ci/check-go-coverage.sh [threshold]
# Default threshold: 45%

set -e

THRESHOLD=${1:-45}
COVERAGE_FILE="backend/coverage.out"

echo "=== Go Coverage Check ==="
echo "Threshold: ${THRESHOLD}%"

cd backend

# Generate coverage report
echo "Generating coverage report..."
go test -coverprofile=coverage.out ./... 2>&1 | grep -v "^?" || true

# Check if coverage file exists
if [ ! -f "coverage.out" ]; then
    echo "ERROR: coverage.out not generated"
    exit 1
fi

# Parse total coverage
TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')

if [ -z "$TOTAL_COVERAGE" ]; then
    echo "ERROR: Could not parse coverage percentage"
    exit 1
fi

echo "Total Coverage: ${TOTAL_COVERAGE}%"

# Compare with threshold
if (( $(echo "$TOTAL_COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "FAIL: Coverage ${TOTAL_COVERAGE}% is below threshold ${THRESHOLD}%"
    exit 1
else
    echo "PASS: Coverage ${TOTAL_COVERAGE}% meets threshold ${THRESHOLD}%"
fi

# Generate HTML report for artifacts
go tool cover -html=coverage.out -o coverage.html
echo "HTML report generated: backend/coverage.html"
