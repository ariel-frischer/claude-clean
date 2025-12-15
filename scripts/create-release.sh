#!/bin/bash
set -e

# Create a GitHub release with built binaries
# Usage: ./scripts/create-release.sh <version> [release notes]

if [ -z "$1" ]; then
    echo "Error: Version required"
    echo "Usage: ./scripts/create-release.sh <version> [release notes]"
    echo "Example: ./scripts/create-release.sh v0.1.0 'Initial release'"
    exit 1
fi

VERSION=$1
NOTES=${2:-"Release ${VERSION}"}
OUTPUT_DIR="release"

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

# Check if binaries exist
if [ ! -d "${OUTPUT_DIR}" ]; then
    echo "Error: ${OUTPUT_DIR}/ directory not found"
    echo "Run 'make build-release' first to build binaries"
    exit 1
fi

echo "Creating GitHub release ${VERSION}..."
echo "Release notes: ${NOTES}"
echo ""

# Check if tag exists locally
if git rev-parse "${VERSION}" >/dev/null 2>&1; then
    echo "Tag ${VERSION} already exists locally"
else
    echo "Creating git tag ${VERSION}..."
    git tag -a "${VERSION}" -m "Release ${VERSION}"
fi

# Detect GitHub remote
GITHUB_REMOTE="origin"
if git remote -v | grep -q "github.com.*github"; then
    GITHUB_REMOTE="github"
fi

# Check if we need to push the tag to GitHub
if git ls-remote --tags "${GITHUB_REMOTE}" | grep -q "${VERSION}"; then
    echo "Tag ${VERSION} already exists on GitHub remote (${GITHUB_REMOTE})"
else
    echo "Pushing tag ${VERSION} to GitHub remote (${GITHUB_REMOTE})..."
    git push "${GITHUB_REMOTE}" "${VERSION}"
fi

# Create the release
echo "Creating GitHub release..."
gh release create "${VERSION}" \
    ${OUTPUT_DIR}/cclean-* \
    ${OUTPUT_DIR}/SHA256SUMS \
    --title "${VERSION}" \
    --notes "${NOTES}"

echo ""
echo "✓ Release ${VERSION} created successfully!"
echo "View it at: $(gh repo view --json url -q .url)/releases/tag/${VERSION}"
