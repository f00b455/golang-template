#!/bin/bash

# Coverage threshold check script for Go projects
# This script calculates test coverage and validates it against a minimum threshold
# Usage: ./scripts/check-coverage.sh [threshold] [coverage-profile]
# Example: ./scripts/check-coverage.sh 95.0 coverage.out

set -e

# Default values
THRESHOLD=${1:-95.0}
COVERAGE_PROFILE=${2:-coverage.out}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    color=$1
    message=$2
    echo -e "${color}${message}${NC}"
}

# Check if coverage profile exists
if [ ! -f "$COVERAGE_PROFILE" ]; then
    print_color "$RED" "❌ Coverage profile not found: $COVERAGE_PROFILE"
    print_color "$YELLOW" "Please run 'go test -coverprofile=$COVERAGE_PROFILE ./...' first"
    exit 1
fi

# Calculate total coverage excluding test files and other non-production code
# The coverage profile already excludes test files, so we just need to filter specific paths
print_color "$YELLOW" "📊 Calculating test coverage..."

# Extract coverage percentage, excluding certain paths
# Paths to exclude from coverage calculation:
# - cmd/: Entry point files with minimal logic
# - docs/: Documentation and generated files
# - scripts/: Build and deployment scripts
# - mocks/: Mock implementations for testing
# - test/: Test utilities and fixtures
# - vendor/: Third-party dependencies
# - .pb.go: Protocol buffer generated files
# - .gen.go: Other generated files

# Get the total coverage from the profile
# Since we're now only testing pkg/ packages, we can use the total directly
COVERAGE_OUTPUT=$(go tool cover -func="$COVERAGE_PROFILE" | tail -n 1)

# Extract the percentage from the output
# The output format is: total: (statements) XX.X%
COVERAGE_PERCENT=$(echo "$COVERAGE_OUTPUT" | awk '{print $NF}' | sed 's/%//')

# Check if we got a valid coverage percentage
if [[ ! "$COVERAGE_PERCENT" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
    print_color "$RED" "❌ Failed to extract coverage percentage"
    echo "Coverage output: $COVERAGE_OUTPUT"
    exit 1
fi

# Print coverage information
echo ""
print_color "$GREEN" "================== COVERAGE REPORT =================="
echo "📈 Total Coverage: ${COVERAGE_PERCENT}%"
echo "🎯 Required Threshold: ${THRESHOLD}%"

# Generate detailed coverage by package
echo ""
echo "📦 Coverage by package:"
echo "---------------------------------------------------"
go tool cover -func="$COVERAGE_PROFILE" | \
    grep -v '^total:' | \
    awk '{
        if (NF >= 3) {
            # Extract package and function name
            split($1, parts, ":")
            pkg = parts[1]
            # Group by package
            if (!(pkg in packages)) {
                packages[pkg] = 0
                count[pkg] = 0
            }
            # Parse percentage (last column)
            pct = $NF
            gsub(/%/, "", pct)
            if (pct != "-" && pct != "") {
                packages[pkg] += pct
                count[pkg]++
            }
        }
    }
    END {
        for (pkg in packages) {
            if (count[pkg] > 0) {
                avg = packages[pkg] / count[pkg]
                printf "  %-50s %6.1f%%\n", pkg, avg
            }
        }
    }' | sort

echo "---------------------------------------------------"

# Compare with threshold using bc for floating point comparison
MEETS_THRESHOLD=$(echo "$COVERAGE_PERCENT >= $THRESHOLD" | bc -l)

if [ "$MEETS_THRESHOLD" -eq 1 ]; then
    echo ""
    print_color "$GREEN" "✅ Coverage check passed! ($COVERAGE_PERCENT% >= $THRESHOLD%)"
    print_color "$GREEN" "================== COVERAGE REPORT =================="
    exit 0
else
    echo ""
    print_color "$RED" "❌ Coverage check failed! ($COVERAGE_PERCENT% < $THRESHOLD%)"

    # Calculate how much coverage is needed
    COVERAGE_GAP=$(echo "$THRESHOLD - $COVERAGE_PERCENT" | bc -l)
    print_color "$YELLOW" "📉 Coverage gap: $COVERAGE_GAP%"

    # Show uncovered lines for improvement hints
    echo ""
    print_color "$YELLOW" "🔍 Files with lowest coverage (focus areas for improvement):"
    echo "---------------------------------------------------"
    go tool cover -func="$COVERAGE_PROFILE" | \
        grep -v '^total:' | \
        grep -E '[0-9]+\.[0-9]+%$' | \
        awk '{
            # Extract percentage
            pct = $NF
            gsub(/%/, "", pct)
            if (pct < 80) {  # Show files with less than 80% coverage
                print $0
            }
        }' | \
        sort -t$'\t' -k3 -n | \
        head -10

    echo "---------------------------------------------------"
    print_color "$RED" "================== COVERAGE REPORT =================="

    exit 1
fi