# businessenger

A small CLI toolkit that pipes developer context into Claude for daily writing tasks.

## Installation

```bash
bash install.sh
```

Symlinks `ai-assistant` into `~/.local/bin`. If that directory isn't in your `PATH`, the script will tell you what to add to your shell config. Also creates all required `context/` directories and empty placeholder files so commands work out of the box.

## Usage

```bash
ai-assistant <command> [args]
ai-assistant --help
```

## Commands

### `standup`

Generates a standup update from today's git commits.

```bash
ai-assistant standup
```

Runs `git log` with diffs since midnight (excluding lock files), then passes the output to Claude to produce a concise, professional bulleted summary.

### `polish <type:language> "<draft>"`

Polishes a rough draft to match your personal writing style.

```bash
ai-assistant polish slack:english "quick update - shipped the auth fix, waiting on review"
ai-assistant polish slack:english --file draft.txt
```

Loads style examples from `context/polish/<type>_<language>_examples.txt` (optional) and asks Claude to correct grammar/spelling while mirroring your tone.

### `wtf "<error>"`

Explains a cryptic error or stack trace in plain English with the top two likely causes and fixes.

```bash
ai-assistant wtf "TypeError: cannot read properties of undefined"
ai-assistant wtf --file error.log
cat error.log | ai-assistant wtf
```

### `tasks "<notes>"`

Extracts actionable tasks from messy meeting notes.

```bash
ai-assistant tasks "alice to fix login bug, bob needs to review PR, blocked on design approval"
ai-assistant tasks --file meeting_notes.txt
```

Outputs a structured task list. Flags blockers and assigns ownership where names appear in the notes.

### `tldr "<thread>"`

Summarises a long Slack thread or email chain in three sentences.

```bash
ai-assistant tldr "long thread text..."
ai-assistant tldr --file thread.txt
cat thread.txt | ai-assistant tldr
```

Tells you the context, the core disagreement or hold-up, and what (if anything) you need to do.

### `diplomat <type:language> "<draft>"`

Rewrites an angry or blunt draft as a professional, diplomatic message.

```bash
ai-assistant diplomat email:english "this is completely broken and i'm frustrated"
ai-assistant diplomat slack:english --file draft.txt
```

Loads style examples from `context/diplomat/<type>_<language>_examples.txt` (optional) and asks Claude to de-escalate while still defending your position.

### `docgen "<code>"`

Generates a Markdown README for source code.

```bash
ai-assistant docgen "$(cat src/script.py)"
ai-assistant docgen --file src/script.py
cat src/script.py | ai-assistant docgen > README.md
```

## Context files

Personal style examples live under `context/` (gitignored — create your own). All commands work without context files; adding examples improves output quality by matching your tone.

```
context/
  polish/<type>_<language>_examples.txt    # optional, improves style matching
  diplomat/<type>_<language>_examples.txt  # optional
  standup/default_examples.txt             # optional
  tasks/default_examples.txt               # optional
  tldr/default_examples.txt                # optional
  docgen/default_examples.txt              # optional
```

## Prerequisites

- [`claude`](https://claude.ai/code) CLI available in your `PATH`

## Project Structure

```
install.sh            # Installs the CLI to ~/.local/bin
ai-assistant          # Entry point
skill-tools/
  standup.sh          # Standup generation
  polish.sh           # Message polishing
  wtf.sh              # Error explanation
  tasks.sh            # Task extraction from meeting notes
  tldr.sh             # Thread summarisation
  diplomat.sh         # Diplomatic rewrite
  docgen.sh           # README generation
context/              # Your personal style examples (gitignored)
```

## License

MIT
