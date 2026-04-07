#!/bin/bash

TYPE=$1
if [[ "$2" == "--file" ]]; then
  DRAFT=$(cat "$3")
else
  DRAFT=$2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -z "$DRAFT" ]; then
  if [ -t 0 ]; then
    echo "Usage: ai-assistant diplomat <type:language> \"<angry draft>\"  OR  cat draft.txt | ai-assistant diplomat <type:language>"
    exit 1
  fi
  DRAFT=$(cat)
fi

CHANNEL="${TYPE%%:*}"
LANGUAGE="${TYPE##*:}"

CONTEXT_FILE="$(realpath "$SCRIPT_DIR/../context/diplomat/${CHANNEL}_${LANGUAGE}_examples.txt")"

CONTEXT_SECTION=""
if [ -f "$CONTEXT_FILE" ]; then
  CONTEXT_SECTION="Mirror the tone and style of these past messages of mine:
$(cat "$CONTEXT_FILE")

"
fi

cat <<EOF | claude
I am providing a very raw, frustrated draft of a message. Rewrite this to be highly professional, diplomatic, and constructive. De-escalate the situation while still firmly protecting my boundaries and addressing the core technical/business issue.

${CONTEXT_SECTION}My draft:
$DRAFT
EOF
