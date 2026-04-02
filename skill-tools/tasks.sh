#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -z "$1" ]; then
  echo "Usage: ai-assistant tasks <notes_file>"
  exit 1
fi

NOTES=$(cat "$1")

CONTEXT_FILE="$(realpath "$SCRIPT_DIR/../context/tasks/default_examples.txt")"

CONTEXT_SECTION=""
if [ -f "$CONTEXT_FILE" ]; then
  CONTEXT_SECTION="Mirror the format and style of these past task lists of mine:
$(cat "$CONTEXT_FILE")

"
fi

cat <<EOF | claude
Extract all actionable items from the following meeting notes. Format them as a list of clearly defined tasks. If a task belongs to someone else, note their name. If a task requires a decision first, flag it as a blocker.

${CONTEXT_SECTION}Meeting notes:
$NOTES
EOF
