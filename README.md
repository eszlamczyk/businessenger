# businessenger

A small CLI toolkit that pipes developer context into Claude for daily writing tasks.

## Installation

```bash
bash install.sh
```

Symlinks `ai-assistant` into `~/.local/bin`. If that directory isn't in your `PATH`, the script will tell you what to add to your shell config. Also creates all required `context/` directories and empty placeholder files so commands work out of the box.

## Configuration

Create `config.json`

```bash
cp config.json.example config.json
```

Then edit `config.json`:

```json
{
  "channels": ["slack", "email"],
  "languages": ["english", "polish"],
  "slack": {
    "workspaces": [
      { "name": "My Team", "token": "xoxp-..." }
    ]
  }
}
```

| Field | Description |
|---|---|
| `channels` | Options shown when selecting a channel (for `polish`, `diplomat`) |
| `languages` | Options shown when selecting a language |
| `slack.workspaces` | Workspaces available in `tldr` Slack mode — add one entry per workspace |

The `slack.workspaces` section is only required if you want the `tldr` Slack mode. See [Slack mode for `tldr`](#slack-mode-for-tldr) for how to get a token.

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

**Slack mode** — fetch and summarise directly from a Slack channel via the TUI (see below). No copy-paste required; the Slack MCP retrieves the messages for you.

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

## TUI

```bash
ai-assistant          # launches the interactive TUI
```

A terminal UI for all commands. Navigate with `↑/↓`, confirm with `enter`, go back with `ctrl+b`.

Most tools offer two input modes — **write** (type/paste) and **file** (browse the filesystem). `tldr` also offers a **slack** mode when Slack workspaces are configured (see below).

### Slack mode for `tldr`

#### 1. Create a Slack app and get a token

1. Go to [api.slack.com/apps](https://api.slack.com/apps) → **Create New App** → **From scratch**.
2. Name it (e.g. "businessenger"), select your workspace, click **Create App**.
3. In the left sidebar go to **OAuth & Permissions** → scroll to **Scopes** → **User Token Scopes** and add:
   - `channels:read` — list public channels
   - `groups:read` — list private channels you're a member of
4. Scroll back to the top of **OAuth & Permissions** and click **Install to Workspace** → **Allow**.
5. Copy the **User OAuth Token** (`xoxp-...`).

#### 2. Add the token to config.json

Add the token to the `slack.workspaces` array in your `config.json` (see [Configuration](#configuration)). Add one entry per workspace — the token is only used locally to fetch channel names.

#### 3. Use it

1. Launch the TUI (`ai-assistant`), select `tldr`, choose **slack** as the input mode.
2. Pick a workspace → channels are fetched live from the Slack API.
3. Pick a channel and a "since" date (relative options or custom text).
4. The Slack MCP retrieves the messages and Claude produces the summary.

> If your `claude` install doesn't auto-discover the Slack MCP, add `--mcp-config <path-to-slack-mcp.json>` to the `claude` call in `skill-tools/tldr.sh`.

## Prerequisites

- [`claude`](https://claude.ai/code) CLI available in your `PATH`
- [Go](https://go.dev/dl/) — required at install time to compile the TUI; not needed afterwards
- Slack MCP configured (only required for `tldr` Slack mode)

## Project Structure

```
install.sh              # Installs the CLI to ~/.local/bin (compiles TUI binary)
ai-assistant            # Entry point (runs TUI when called with no arguments)
ai-assistant-tui        # Compiled TUI binary — built by install.sh (gitignored)
config.json.example     # Config template — copy to config.json and fill in credentials
config.json             # Your config (gitignored — never committed)
tui/
  main.go               # Interactive terminal UI (Go)
  slack.go              # Slack workspace/channel picker and API client
skill-tools/
  standup.sh            # Standup generation
  polish.sh             # Message polishing
  wtf.sh                # Error explanation
  tasks.sh              # Task extraction from meeting notes
  tldr.sh               # Thread summarisation (+ Slack MCP mode)
  diplomat.sh           # Diplomatic rewrite
  docgen.sh             # README generation
context/                # Your personal style examples (gitignored)
```

## License

MIT
