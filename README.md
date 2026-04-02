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

Loads style examples from `context/<source>_<language>_examples.txt` and asks Claude to correct grammar/spelling while mirroring your tone.

## Setup

1. Clone the repo and make the entry point executable:
   ```bash
   chmod +x ai-assistant
   ```
2. Add style examples for `polish` under `context/` (gitignored — create your own):
   ```
   context/<source>_<language>_examples.txt
   ```
3. Ensure `claude` CLI is available in your `PATH`.

## Project Structure

```
ai-assistant          # Entry point
skill-tools/
  standup.sh          # Standup generation skill
  polish.sh           # Message polishing skill
context/              # Your personal style examples (gitignored)
```

## License

MIT
