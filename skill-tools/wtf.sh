#!/bin/bash

if [ -t 0 ] && [ -z "$1" ]; then
  echo "Usage: ai-assistant wtf <error_file>  OR  cat error.log | ai-assistant wtf"
  exit 1
fi

if [ -n "$1" ]; then
  ERROR_INPUT=$(cat "$1")
else
  ERROR_INPUT=$(cat)
fi

cat <<EOF | claude
I am getting this error/stack trace. Explain exactly what is breaking in plain English, and give me the top two most likely causes and how to fix them.

Error:
$ERROR_INPUT
EOF
