#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
COVERAGE_FILE=${1:-"coverage.out"}
MIN_COVERAGE=${2:-80}

if [[ "$COVERAGE_FILE" != /* ]]; then
    COVERAGE_FILE="$PROJECT_ROOT/$COVERAGE_FILE"
fi

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "Error: Coverage file '$COVERAGE_FILE' not found"
    echo "Usage: $0 [coverage_file] [min_coverage_percent]"
    exit 1
fi

echo "Checking test coverage..."
echo "Coverage file: $COVERAGE_FILE"
echo "Minimum required: ${MIN_COVERAGE}%"

TOTAL_COVERAGE=$(cd "$PROJECT_ROOT/backend" && go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "total:" | awk '{print $3}' | sed 's/%//')

if [ -z "$TOTAL_COVERAGE" ]; then
    echo "Error: Could not parse coverage from file"
    exit 1
fi

COVERAGE_INT=${TOTAL_COVERAGE%.*}
if [ -z "$COVERAGE_INT" ]; then
    COVERAGE_INT=0
fi

echo "Current coverage: ${TOTAL_COVERAGE}%"

if [ "$COVERAGE_INT" -lt "$MIN_COVERAGE" ]; then
    echo ""
    echo "❌ Coverage check FAILED"
    echo "   Required: ${MIN_COVERAGE}%"
    echo "   Actual: ${TOTAL_COVERAGE}%"
    echo ""
    echo "To see detailed coverage report:"
    echo "  go tool cover -html=$COVERAGE_FILE -o coverage.html"
    exit 1
else
    echo ""
    echo "✅ Coverage check PASSED"
    echo "   Required: ${MIN_COVERAGE}%"
    echo "   Actual: ${TOTAL_COVERAGE}%"
    exit 0
fi
