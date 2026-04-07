#!/bin/bash

TYPE=$1
if [[ "$2" == "--file" ]]; then
  DRAFT=$(cat "$3")
else
  DRAFT=$2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

SOURCE="${TYPE%%:*}"
LANGUAGE="${TYPE##*:}"

CONTEXT_FILE="$SCRIPT_DIR/../context/polish/${SOURCE}_${LANGUAGE}_examples.txt"

CONTEXT_SECTION=""
if [ -f "$CONTEXT_FILE" ]; then
  CONTEXT_SECTION="CRITICAL INSTRUCTION: You MUST mirror the tone, structure, and style of my past messages provided below. Do not sound like generic AI.

Examples:
$(cat "$CONTEXT_FILE")

"
fi

cat <<EOF | claude
You are my personal AI editor. Correct the grammar and spelling of my draft and adapt it into professional business language.

${CONTEXT_SECTION}My Rough Draft:
$DRAFT
EOF
