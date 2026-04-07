#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET="$HOME/.local/bin/ai-assistant"

mkdir -p "$HOME/.local/bin"

ln -sf "$SCRIPT_DIR/ai-assistant" "$TARGET"

# Create context directories and placeholder files
for dir in standup tasks tldr docgen polish diplomat; do
  mkdir -p "$SCRIPT_DIR/context/$dir"
done

for file in \
  "$SCRIPT_DIR/context/standup/default_examples.txt" \
  "$SCRIPT_DIR/context/tasks/default_examples.txt" \
  "$SCRIPT_DIR/context/tldr/default_examples.txt" \
  "$SCRIPT_DIR/context/docgen/default_examples.txt"
do
  [ -f "$file" ] || touch "$file"
done

if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then
  echo "Warning: $HOME/.local/bin is not in your PATH."
  echo "Add the following to your shell config (~/.bashrc, ~/.zshrc, etc.):"
  echo ""
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  echo ""
fi

echo "Installed: ai-assistant -> $TARGET"
