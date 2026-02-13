#!/bin/bash
# Check Go test coverage threshold
# Usage: ./scripts/ci/check-go-coverage.sh [threshold] [coverage_file]
# Default threshold: 30%
# Default coverage file: backend/coverage.out

set -e

THRESHOLD=${1:-30}
COVERAGE_FILE=${2:-"backend/coverage.out"}

echo "=== Go Coverage Check ==="
echo "Threshold: ${THRESHOLD}%"
echo "Coverage file: ${COVERAGE_FILE}"

cd backend

# Check if coverage file exists
if [ ! -f "coverage.out" ] && [ ! -f "coverage-db.out" ]; then
    echo "WARNING: No coverage file found, running tests..."
    go test -coverprofile=coverage.out ./... 2>&1 | grep -v "^?" || true
fi

# Use specified coverage file or default
if [ -f "${COVERAGE_FILE#backend/}" ]; then
    ACTUAL_FILE="${COVERAGE_FILE#backend/}"
else
    ACTUAL_FILE="coverage.out"
fi

if [ ! -f "$ACTUAL_FILE" ]; then
    echo "ERROR: Coverage file not found: $ACTUAL_FILE"
    exit 1
fi

# Parse total coverage
TOTAL_COVERAGE=$(go tool cover -func=$ACTUAL_FILE | grep "total:" | awk '{print $3}' | sed 's/%//')

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
HTML_FILE="${ACTUAL_FILE%.out}.html"
go tool cover -html=$ACTUAL_FILE -o $HTML_FILE
echo "HTML report generated: backend/${HTML_FILE}"
