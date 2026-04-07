#!/bin/bash

if [[ "$1" == "--file" ]]; then
  ERROR_INPUT=$(cat "$2")
elif [ -n "$1" ]; then
  ERROR_INPUT="$1"
elif [ ! -t 0 ]; then
  ERROR_INPUT=$(cat)
else
  echo "Usage: ai-assistant wtf \"<error>\"  OR  ai-assistant wtf --file <error_file>  OR  cat error.log | ai-assistant wtf"
  exit 1
fi

cat <<EOF | claude
I am getting this error/stack trace. Explain exactly what is breaking in plain English, and give me the top two most likely causes and how to fix them.

Error:
$ERROR_INPUT
EOF
