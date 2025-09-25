#!/bin/bash

# Script to verify that all feature files have proper issue references
# This is called from CI but can also be run locally
# Updated: 2025-09-24

echo "Checking feature files for issue references..."
REPO="${GITHUB_REPOSITORY:-f00b455/ts-pure-template}"
ERR=0

# Find all feature files in all packages
while IFS= read -r -d '' f; do
  echo "Checking $f"

  # Debug: Show first few lines
  if [ "${DEBUG:-}" = "1" ]; then
    echo "First 5 lines of $f:"
    head -n 5 "$f"
    echo "---"
  fi

  # Check for Issue comment - be more flexible with regex
  if ! head -n 5 "$f" | grep -qE "^# *Issue: *#[0-9]+"; then
    echo "Error: Missing 'Issue: #<num>' comment"
    if [ "${GITHUB_ACTIONS:-}" = "true" ]; then
      echo "::error file=$f::Missing 'Issue: #<num>' comment"
    fi
    ERR=1
  fi

  # Check for URL comment - be more flexible
  if ! head -n 5 "$f" | grep -qE "^# *URL: *https://github\.com/${REPO}/issues/[0-9]+"; then
    echo "Error: Missing canonical URL comment"
    if [ "${GITHUB_ACTIONS:-}" = "true" ]; then
      echo "::error file=$f::Missing canonical URL comment"
    fi
    ERR=1
  fi

  # Check for issue tag - be more flexible
  if ! head -n 10 "$f" | grep -qE "@.*issue-[0-9]+"; then
    echo "Error: Missing '@issue-<num>' tag"
    if [ "${GITHUB_ACTIONS:-}" = "true" ]; then
      echo "::error file=$f::Missing '@issue-<num>' tag"
    fi
    ERR=1
  fi
done < <(find . -name "*.feature" -path "*/features/*" -print0)

if [ $ERR -eq 0 ]; then
  echo "✅ All feature files have valid issue references"
else
  echo "❌ Some feature files are missing issue references"
  exit $ERR
fi