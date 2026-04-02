#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

LOG_SINCE="midnight"

# -p shows the diff.
# -U1 shows only 1 line of context around changes to keep it short.
# The ":(exclude)" syntax prevents huge auto-generated files from flooding the prompt.
GIT_INFO=$(git log \
  --all \
  --no-merges \
  --author="$(git config user.name)" \
  --since="$LOG_SINCE" \
  --stat \
  -p \
  -U1 \
  -- ":(exclude)*lock.json" ":(exclude)*.lock" ":(exclude)*.min.*")


if [ -z "$GIT_INFO" ]; then
	echo "No commits found for today"
	exit 0
fi

CONTEXT_FILE="$(realpath "$SCRIPT_DIR/../context/standup/default_examples.txt")"

CONTEXT_SECTION=""
if [ -f "$CONTEXT_FILE" ]; then
  CONTEXT_SECTION="Mirror the format and tone of these past standup updates of mine:
$(cat "$CONTEXT_FILE")

"
fi

cat <<EOF | claude
You are a helpful assistant writing my standup update. Convert the following git commit messages into
a concise, professional bulleted list of accomplishments for today. Do not use overly technical jargon.

${CONTEXT_SECTION}Git commits (with diffs):
$GIT_INFO
EOF
