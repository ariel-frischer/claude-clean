#!/bin/bash
# Install git hooks for this repository

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
HOOKS_DIR="$SCRIPT_DIR/hooks"
GIT_HOOKS_DIR="$(git rev-parse --git-dir)/hooks"

if [ ! -d "$HOOKS_DIR" ]; then
    echo "Error: hooks directory not found"
    exit 1
fi

for hook in "$HOOKS_DIR"/*; do
    if [ -f "$hook" ]; then
        hookname=$(basename "$hook")
        cp "$hook" "$GIT_HOOKS_DIR/$hookname"
        chmod +x "$GIT_HOOKS_DIR/$hookname"
        echo "✓ Installed $hookname"
    fi
done

echo "Done! Git hooks installed."
