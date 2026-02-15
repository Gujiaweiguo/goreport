#!/bin/bash
# Usage: ./scripts/ci/check-go-coverage.sh [threshold] [coverage_file]
# Default threshold: 50%, Default coverage file: backend/coverage.out

set -e

THRESHOLD=${1:-50}
COVERAGE_FILE=${2:-"backend/coverage.out"}

echo "=== Go Coverage Check ==="
echo "Threshold: ${THRESHOLD}%"
echo "Coverage file: ${COVERAGE_FILE}"

cd backend

if [ ! -f "coverage.out" ]; then
    echo "WARNING: No coverage file found, running tests..."
    go test -coverprofile=coverage.out ./... 2>&1 | grep -v "^?" || true
fi

if [ -f "${COVERAGE_FILE#backend/}" ]; then
    ACTUAL_FILE="${COVERAGE_FILE#backend/}"
else
    ACTUAL_FILE="coverage.out"
fi

if [ ! -f "$ACTUAL_FILE" ]; then
    echo "ERROR: Coverage file not found: $ACTUAL_FILE"
    exit 1
fi

echo ""
echo "=== Coverage by Package ==="
TOTAL_COVERAGE=$(go test ./... -cover 2>&1 | grep "coverage:" | sed 's/.*coverage: \([0-9.]*\)%.*/\1/' | awk '{sum+=$1; count++} END {printf "%.1f", sum/count}')

if [ -z "$TOTAL_COVERAGE" ] || [ "$TOTAL_COVERAGE" = "0.0" ]; then
    echo "ERROR: Could not parse coverage percentage"
    exit 1
fi

echo ""
echo "=== Total Average Coverage: ${TOTAL_COVERAGE}% ==="

if (( $(echo "$TOTAL_COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "FAIL: Coverage ${TOTAL_COVERAGE}% is below threshold ${THRESHOLD}%"
    echo ""
    echo "Low-coverage packages:"
    go test ./... -cover 2>&1 | grep "coverage:" | sort -t: -k2 -n | head -5
    exit 1
else
    echo "PASS: Coverage ${TOTAL_COVERAGE}% meets threshold ${THRESHOLD}%"
fi

HTML_FILE="${ACTUAL_FILE%.out}.html"
go tool cover -html=$ACTUAL_FILE -o $HTML_FILE 2>/dev/null || true
echo "HTML report generated: backend/${HTML_FILE}"
