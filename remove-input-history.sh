#!/bin/bash

# Remove input.txt files from git history while keeping them in working directory
# This is to comply with the AoC FAQ: https://adventofcode.com/2025/about

# This script uses git filter-branch to rewrite history

set -e

echo "Removing day*/input.txt from git history..."
echo "This will rewrite your repository history."
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# Remove the files from history
git filter-branch --tree-filter 'rm -f day*/input.txt' --prune-empty -f -- --all

echo "Removing references to old commits..."
# Clean up refs
rm -rf .git/refs/original/

# Run garbage collection to reclaim space
echo "Running garbage collection..."
git gc --aggressive --prune=now

echo "Done! The working directory files are preserved."
echo "You can now push with: git push --force-with-lease"
