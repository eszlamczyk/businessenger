#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "$1" == "--file" ]]; then
  SOURCE=$(cat "$2")
elif [ -n "$1" ]; then
  SOURCE="$1"
elif [ ! -t 0 ]; then
  SOURCE=$(cat)
else
  echo "Usage: ai-assistant docgen \"<code>\"  OR  ai-assistant docgen --file <source_file>  OR  cat script.py | ai-assistant docgen"
  exit 1
fi

CONTEXT_FILE="$(realpath "$SCRIPT_DIR/../context/docgen/default_examples.txt")"

CONTEXT_SECTION=""
if [ -f "$CONTEXT_FILE" ]; then
  CONTEXT_SECTION="Mirror the format and style of these past documentation examples of mine:
$(cat "$CONTEXT_FILE")

"
fi

cat <<EOF | claude
Read the following source code. Generate a clear, developer-friendly Markdown README explaining what this code does, what its inputs and outputs are, and providing one basic usage example.

${CONTEXT_SECTION}Source code:
$SOURCE
EOF
