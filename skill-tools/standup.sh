#!/bin/bash

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

cat <<EOF | claude
You are a helpful assistant writing my standup update. Convert the following git commit messages into
a concise, professional bulleted list of accomplishments for today. Do not use overly technical jargon.
 
Git commits (with diffs):
$GIT_INFO
EOF

