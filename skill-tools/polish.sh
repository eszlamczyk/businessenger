#!/bin/bash

TYPE=$1
DRAFT=$2

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

SOURCE="${TYPE%%:*}"
LANGUAGE="${TYPE##*:}"

CONTEXT_FILE="$(realpath "$SCRIPT_DIR/../context/${SOURCE}_${LANGUAGE}_examples.txt")"

CONTEXT=$(cat "$CONTEXT_FILE")

cat <<EOF | claude
You are my personal AI editor. Correct the grammar and spelling of my draft and adapt it into professional business language.

CRITICAL INSTRUCTION: You MUST mirror the tone, structure, and style of my past messages provided below. Do not sound like genereic AI.

Examples:
$CONTEXT

My Rough Draft:
$DRAFT
EOF
