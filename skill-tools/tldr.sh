#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -t 0 ] && [ -z "$1" ]; then
  echo "Usage: ai-assistant tldr <thread_file>  OR  cat thread.txt | ai-assistant tldr"
  exit 1
fi

if [ -n "$1" ]; then
  THREAD=$(cat "$1")
else
  THREAD=$(cat)
fi

CONTEXT_FILE="$(realpath "$SCRIPT_DIR/../context/tldr/default_examples.txt")"

CONTEXT_SECTION=""
if [ -f "$CONTEXT_FILE" ]; then
  CONTEXT_SECTION="Mirror the format and tone of these past summaries of mine:
$(cat "$CONTEXT_FILE")

"
fi

cat <<EOF | claude
Read this long communication thread. Give me a 3-sentence summary of the context, tell me what the core disagreement or hold-up is, and clearly state what action (if any) is required from me.

${CONTEXT_SECTION}Thread:
$THREAD
EOF
