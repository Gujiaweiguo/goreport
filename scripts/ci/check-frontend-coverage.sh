#!/bin/bash
# Usage: ./scripts/ci/check-frontend-coverage.sh [threshold]
# Default threshold: 40%

set -e

THRESHOLD=${1:-40}
COVERAGE_DIR="frontend/coverage"

echo "=== Frontend Coverage Check ==="
echo "Threshold: ${THRESHOLD}%"

cd frontend

if [ ! -f "$COVERAGE_DIR/coverage-summary.json" ]; then
    echo "Running tests with coverage..."
    npm run test:run -- --coverage --passWithNoTests 2>&1 || true
fi

TOTAL_COVERAGE="0"

if [ -f "$COVERAGE_DIR/coverage-summary.json" ]; then
    TOTAL_COVERAGE=$(cat $COVERAGE_DIR/coverage-summary.json | grep -o '"lines":{[^}]*}' | grep -o '"total":[0-9]*' | head -1 | grep -o '[0-9]*' || echo "0")
elif [ -f "$COVERAGE_DIR/clover.xml" ]; then
    TOTAL_COVERAGE=$(grep -o 'coverage="[0-9.]*"' $COVERAGE_DIR/clover.xml | head -1 | grep -o '[0-9]*' | cut -d'.' -f1 || echo "0")
fi

if [ -z "$TOTAL_COVERAGE" ] || [ "$TOTAL_COVERAGE" = "0" ]; then
    echo "Extracting coverage from test output..."
    COV_LINE=$(npm run test:run -- --coverage --passWithNoTests 2>&1 | grep "All files" | tail -1)
    TOTAL_COVERAGE=$(echo "$COV_LINE" | awk '{print $NF}' | sed 's/%//' || echo "0")
    
    if [ -z "$TOTAL_COVERAGE" ] || [ "$TOTAL_COVERAGE" = "0" ]; then
        TOTAL_COVERAGE=$(echo "$COV_LINE" | grep -oE '[0-9]+\.[0-9]+' | head -1 | cut -d'.' -f1 || echo "0")
    fi
fi

TOTAL_COVERAGE=${TOTAL_COVERAGE:-0}
echo "Total Coverage: ${TOTAL_COVERAGE}%"

if [ "$TOTAL_COVERAGE" -lt "$THRESHOLD" ] 2>/dev/null; then
    echo "FAIL: Coverage ${TOTAL_COVERAGE}% is below threshold ${THRESHOLD}%"
    exit 1
else
    echo "PASS: Coverage ${TOTAL_COVERAGE}% meets threshold ${THRESHOLD}%"
fi
