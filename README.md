# businessenger

A small CLI toolkit that pipes developer context into Claude for daily writing tasks.

## Usage

```bash
./ai-assistant <command> [args]
```

## Commands

### `standup`

Generates a standup update from today's git commits.

```bash
./ai-assistant standup
```

Runs `git log` with diffs since midnight (excluding lock files), then passes the output to Claude to produce a concise, professional bulleted summary.

### `polish <type:language> <draft>`

Polishes a rough draft to match your personal writing style.

```bash
./ai-assistant polish slack:english "quick update - shipped the auth fix, waiting on review"
```

Loads style examples from `context/polish/<type>_<language>_examples.txt` and asks Claude to correct grammar/spelling while mirroring your tone.

### `wtf [file]`

Explains a cryptic error or stack trace in plain English with the top two likely causes and fixes.

```bash
./ai-assistant wtf error.log
cat error.log | ./ai-assistant wtf
```

### `tasks <notes_file>`

Extracts actionable tasks from messy meeting notes.

```bash
./ai-assistant tasks meeting_notes.txt
```

Outputs a structured task list. Flags blockers and assigns ownership where names appear in the notes.

### `tldr [file]`

Summarises a long Slack thread or email chain in three sentences.

```bash
./ai-assistant tldr thread.txt
cat thread.txt | ./ai-assistant tldr
```

Tells you the context, the core disagreement or hold-up, and what (if anything) you need to do.

### `diplomat <type:language> <draft>`

Rewrites an angry or blunt draft as a professional, diplomatic message.

```bash
./ai-assistant diplomat email:english "my frustrated draft"
cat draft.txt | ./ai-assistant diplomat slack:english
```

Loads style examples from `context/diplomat/<type>_<language>_examples.txt` (optional) and asks Claude to de-escalate while still defending your position.

### `docgen [file]`

Generates a Markdown README for a source file.

```bash
./ai-assistant docgen src/script.py
cat src/script.py | ./ai-assistant docgen > README.md
```

## Setup

1. Clone the repo and make the entry point executable:
   ```bash
   chmod +x ai-assistant
   ```
2. Ensure `claude` CLI is available in your `PATH`.
3. Optionally add personal style examples under `context/` (gitignored):
   ```
   context/
     polish/<type>_<language>_examples.txt    # required for polish
     diplomat/<type>_<language>_examples.txt  # optional
     standup/default_examples.txt             # optional
     tasks/default_examples.txt               # optional
     tldr/default_examples.txt                # optional
     docgen/default_examples.txt              # optional
   ```

## Project Structure

```
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
