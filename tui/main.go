package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	highlight = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	dim       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	dimDesc   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	title     = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	label     = lipgloss.NewStyle().Bold(true)
)

// ── Config ────────────────────────────────────────────────────────────────────

type Config struct {
	Channels  []string    `json:"channels"`
	Languages []string    `json:"languages"`
	Slack     SlackConfig `json:"slack"`
}

func loadConfig(dir string) Config {
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return Config{Channels: []string{"slack", "email"}, Languages: []string{"english"}}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg
}

// ── Tools ─────────────────────────────────────────────────────────────────────

type tool struct {
	name          string
	description   string
	needsType     bool
	inputLabel    string
	supportsSlack bool
}

var tools = []tool{
	{"standup", "Generate standup from git commits", false, "", false},
	{"polish", "Polish prose to sound professional", true, "Your draft", false},
	{"wtf", "Explain an error message", false, "Paste your error", false},
	{"tasks", "Extract tasks from meeting notes", false, "Paste your notes", false},
	{"tldr", "Summarize a communication thread", false, "Paste the thread", true},
	{"diplomat", "Rewrite a message diplomatically", true, "Your draft", false},
	{"docgen", "Generate docs for source code", false, "Paste your code", false},
}

// ── Screens ───────────────────────────────────────────────────────────────────

type screen int

const (
	screenMenu screen = iota
	screenChannel
	screenLanguage
	screenInputMode
	screenSlackWorkspace
	screenSlackLoading
	screenSlackChannel
	screenSlackDate
	screenSlackDateCustom
	screenFileBrowser
	screenInput
	screenDone
)

// ── File browser ──────────────────────────────────────────────────────────────

type fileEntry struct {
	name  string
	isDir bool
	isUp  bool
}

func loadEntries(dir string) []fileEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}
		return entries[i].Name() < entries[j].Name()
	})

	var result []fileEntry

	if filepath.Dir(dir) != dir {
		result = append(result, fileEntry{name: "..", isUp: true})
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		result = append(result, fileEntry{name: e.Name(), isDir: e.IsDir()})
	}

	return result
}

// ── Model ─────────────────────────────────────────────────────────────────────

type model struct {
	screen         screen
	cursor         int
	subCursor      int
	channel        string
	language       string
	browserDir     string
	browserEntries []fileEntry
	useFile        bool
	input          string
	config         Config
	useSlack       bool
	slackWorkspace SlackWorkspace
	slackChannels  []string
	slackChannel   string
	slackDateLabel string
	slackDateInput string
	loadErr        string
}

func initialModel() model {
	dir := scriptDir()
	cwd, _ := os.Getwd()
	return model{
		config:         loadConfig(dir),
		browserDir:     cwd,
		browserEntries: loadEntries(cwd),
	}
}

func (m model) Init() tea.Cmd { return tea.EnableBracketedPaste }

// contextHeading builds the title line showing accumulated selections.
func contextHeading(m model) string {
	t := tools[m.cursor]
	if m.channel != "" && m.language != "" {
		return fmt.Sprintf("%s (%s:%s)", t.name, m.channel, m.language)
	}
	return t.name
}

// inputModes returns the available input modes for the current tool.
func inputModes(m model) []struct{ name, desc string } {
	modes := []struct{ name, desc string }{
		{"write", "type or paste content"},
		{"file", "browse and select a file"},
	}
	if tools[m.cursor].supportsSlack && len(m.config.Slack.Workspaces) > 0 {
		modes = append(modes, struct{ name, desc string }{"slack", "fetch from Slack via MCP"})
	}
	return modes
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case slackChannelsMsg:
		if msg.err != nil {
			m.loadErr = msg.err.Error()
			m.screen = screenSlackWorkspace
		} else {
			m.slackChannels = msg.channels
			m.subCursor = 0
			m.loadErr = ""
			m.screen = screenSlackChannel
		}
		return m, nil

	case tea.KeyMsg:
		switch m.screen {

		case screenMenu:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(tools)-1 {
					m.cursor++
				}
			case "enter":
				t := tools[m.cursor]
				if t.needsType {
					m.subCursor = 0
					m.screen = screenChannel
				} else if t.inputLabel != "" || t.supportsSlack {
					m.subCursor = 0
					m.screen = screenInputMode
				} else {
					m.screen = screenDone
					return m, runTool(m)
				}
			}

		case screenChannel:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenMenu
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(m.config.Channels)-1 {
					m.subCursor++
				}
			case "enter":
				m.channel = m.config.Channels[m.subCursor]
				m.subCursor = 0
				m.screen = screenLanguage
			}

		case screenLanguage:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenChannel
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(m.config.Languages)-1 {
					m.subCursor++
				}
			case "enter":
				m.language = m.config.Languages[m.subCursor]
				m.subCursor = 0
				m.screen = screenInputMode
			}

		case screenInputMode:
			modes := inputModes(m)
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				t := tools[m.cursor]
				if t.needsType {
					m.screen = screenLanguage
				} else {
					m.screen = screenMenu
				}
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(modes)-1 {
					m.subCursor++
				}
			case "enter":
				selected := modes[m.subCursor].name
				switch selected {
				case "file":
					m.useFile = true
					m.useSlack = false
					m.subCursor = 0
					m.screen = screenFileBrowser
				case "slack":
					m.useFile = false
					m.useSlack = true
					m.subCursor = 0
					m.screen = screenSlackWorkspace
				default:
					m.useFile = false
					m.useSlack = false
					m.input = ""
					m.screen = screenInput
				}
			}

		case screenSlackWorkspace:
			workspaces := m.config.Slack.Workspaces
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.loadErr = ""
				m.screen = screenInputMode
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(workspaces)-1 {
					m.subCursor++
				}
			case "enter":
				m.slackWorkspace = workspaces[m.subCursor]
				m.loadErr = ""
				m.screen = screenSlackLoading
				return m, fetchSlackChannels(m.slackWorkspace.Token)
			}

		case screenSlackLoading:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenSlackWorkspace
			}

		case screenSlackChannel:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenSlackWorkspace
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(m.slackChannels)-1 {
					m.subCursor++
				}
			case "enter":
				m.slackChannel = m.slackChannels[m.subCursor]
				m.subCursor = 0
				m.screen = screenSlackDate
			}

		case screenSlackDate:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenSlackChannel
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(slackDateOptions)-1 {
					m.subCursor++
				}
			case "enter":
				chosen := slackDateOptions[m.subCursor]
				if chosen == "custom..." {
					m.slackDateInput = ""
					m.screen = screenSlackDateCustom
				} else {
					m.slackDateLabel = chosen
					m.screen = screenDone
					return m, runTool(m)
				}
			}

		case screenSlackDateCustom:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenSlackDate
			case "enter":
				if m.slackDateInput != "" {
					m.slackDateLabel = m.slackDateInput
					m.screen = screenDone
					return m, runTool(m)
				}
			case "backspace":
				if len(m.slackDateInput) > 0 {
					m.slackDateInput = m.slackDateInput[:len(m.slackDateInput)-1]
				}
			default:
				if msg.Type == tea.KeyRunes {
					m.slackDateInput += string(msg.Runes)
				}
			}

		case screenFileBrowser:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+b":
				m.subCursor = 0
				m.screen = screenInputMode
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(m.browserEntries)-1 {
					m.subCursor++
				}
			case "enter":
				entry := m.browserEntries[m.subCursor]
				if entry.isUp {
					newDir := filepath.Dir(m.browserDir)
					m.browserDir = newDir
					m.browserEntries = loadEntries(newDir)
					m.subCursor = 0
				} else if entry.isDir {
					newDir := filepath.Join(m.browserDir, entry.name)
					m.browserDir = newDir
					m.browserEntries = loadEntries(newDir)
					m.subCursor = 0
				} else {
					m.screen = screenDone
					return m, runTool(m)
				}
			}

		case screenInput:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl+b":
				m.input = ""
				m.subCursor = 0
				m.screen = screenInputMode
			case "ctrl+d":
				m.screen = screenDone
				return m, runTool(m)
			case "enter":
				m.input += "\n"
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				if msg.Type == tea.KeyRunes {
					m.input += string(msg.Runes)
				}
			}
		}
	}

	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func selectionList(heading, subheading string, items []string, cursor int) string {
	var b strings.Builder
	b.WriteString(title.Render(heading) + "\n")
	b.WriteString(dimDesc.Render(subheading) + "\n\n")
	for i, item := range items {
		if i == cursor {
			b.WriteString(fmt.Sprintf("%s%s\n", highlight.Render("> "), highlight.Render(item)))
		} else {
			b.WriteString(fmt.Sprintf("  %s\n", dim.Render(item)))
		}
	}
	b.WriteString(dim.Render("\n↑/↓ navigate • enter select • ctrl+b back • q quit"))
	return b.String()
}

func (m model) View() string {
	switch m.screen {

	case screenMenu:
		var b strings.Builder
		b.WriteString(title.Render("businessenger") + "\n")
		for i, t := range tools {
			cursor := "  "
			padded := fmt.Sprintf("%-10s", t.name)
			name := dim.Render(padded)
			desc := dimDesc.Render(t.description)
			if i == m.cursor {
				cursor = highlight.Render("> ")
				name = highlight.Render(padded)
				desc = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(t.description)
			}
			b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, name, desc))
		}
		b.WriteString(dim.Render("\n↑/↓ navigate • enter select • q quit"))
		return b.String()

	case screenChannel:
		t := tools[m.cursor]
		return selectionList(t.name, "select channel", m.config.Channels, m.subCursor)

	case screenLanguage:
		t := tools[m.cursor]
		return selectionList(
			fmt.Sprintf("%s (%s)", t.name, m.channel),
			"select language",
			m.config.Languages,
			m.subCursor,
		)

	case screenInputMode:
		modes := inputModes(m)
		var b strings.Builder
		b.WriteString(title.Render(contextHeading(m)) + "\n")
		b.WriteString(dimDesc.Render("select input mode") + "\n\n")
		for i, mode := range modes {
			padded := fmt.Sprintf("%-6s", mode.name)
			if i == m.subCursor {
				b.WriteString(fmt.Sprintf("%s%s %s\n",
					highlight.Render("> "),
					highlight.Render(padded),
					lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(mode.desc),
				))
			} else {
				b.WriteString(fmt.Sprintf("  %s %s\n", dim.Render(padded), dimDesc.Render(mode.desc)))
			}
		}
		b.WriteString(dim.Render("\n↑/↓ navigate • enter select • ctrl+b back • q quit"))
		return b.String()

	case screenSlackWorkspace:
		workspaces := m.config.Slack.Workspaces
		names := make([]string, len(workspaces))
		for i, ws := range workspaces {
			names[i] = ws.Name
		}
		var b strings.Builder
		b.WriteString(selectionList("tldr — slack", "select workspace", names, m.subCursor))
		if m.loadErr != "" {
			b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("error: "+m.loadErr))
		}
		return b.String()

	case screenSlackLoading:
		var b strings.Builder
		b.WriteString(title.Render("tldr — slack") + "\n")
		b.WriteString(dimDesc.Render(fmt.Sprintf("fetching channels from %s…", m.slackWorkspace.Name)) + "\n\n")
		b.WriteString(dim.Render("ctrl+b cancel • q quit"))
		return b.String()

	case screenSlackChannel:
		return selectionList(
			fmt.Sprintf("tldr — %s", m.slackWorkspace.Name),
			"select channel",
			m.slackChannels,
			m.subCursor,
		)

	case screenSlackDate:
		return selectionList(
			fmt.Sprintf("tldr — #%s", m.slackChannel),
			"messages since",
			slackDateOptions,
			m.subCursor,
		)

	case screenSlackDateCustom:
		var b strings.Builder
		b.WriteString(title.Render(fmt.Sprintf("tldr — #%s", m.slackChannel)) + "\n")
		b.WriteString(dimDesc.Render("enter a date or time reference") + "\n\n")
		b.WriteString(label.Render("since: ") + m.slackDateInput + "\n\n")
		b.WriteString(dim.Render("enter confirm • ctrl+b back • q quit"))
		return b.String()

	case screenFileBrowser:
		var b strings.Builder
		b.WriteString(title.Render(contextHeading(m)) + "\n")
		b.WriteString(dimDesc.Render(m.browserDir) + "\n\n")
		for i, entry := range m.browserEntries {
			var display string
			if entry.isUp {
				display = ".."
			} else if entry.isDir {
				display = entry.name + "/"
			} else {
				display = entry.name
			}
			if i == m.subCursor {
				b.WriteString(fmt.Sprintf("%s%s\n", highlight.Render("> "), highlight.Render(display)))
			} else {
				if entry.isDir || entry.isUp {
					b.WriteString(fmt.Sprintf("  %s\n", dim.Render(display)))
				} else {
					b.WriteString(fmt.Sprintf("  %s\n", dimDesc.Render(display)))
				}
			}
		}
		b.WriteString(dim.Render("\n↑/↓ navigate • enter open/select • ctrl+b back • q quit"))
		return b.String()

	case screenInput:
		t := tools[m.cursor]
		return fmt.Sprintf(
			"%s\n\n%s\n%s\n\n%s",
			title.Render(contextHeading(m)),
			label.Render(t.inputLabel+":"),
			m.input,
			dim.Render("ctrl+d to confirm • ctrl+b back"),
		)

	case screenDone:
		return ""
	}
	return ""
}

// ── Exec ──────────────────────────────────────────────────────────────────────

func runTool(m model) tea.Cmd {
	return tea.ExecProcess(buildCmd(m), func(err error) tea.Msg {
		return tea.Quit()
	})
}

func buildCmd(m model) *exec.Cmd {
	t := tools[m.cursor]
	scriptsDir := scriptDir()
	script := filepath.Join(scriptsDir, "skill-tools", t.name+".sh")

	var args []string
	if m.useSlack {
		args = append(args, "--slack", m.slackWorkspace.Name, "#"+m.slackChannel, m.slackDateLabel)
	} else if t.needsType {
		if m.useFile {
			filePath := filepath.Join(m.browserDir, m.browserEntries[m.subCursor].name)
			args = append(args, m.channel+":"+m.language, "--file", filePath)
		} else {
			args = append(args, m.channel+":"+m.language, m.input)
		}
	} else {
		if m.useFile {
			filePath := filepath.Join(m.browserDir, m.browserEntries[m.subCursor].name)
			args = append(args, "--file", filePath)
		} else if m.input != "" {
			args = append(args, m.input)
		}
	}

	cmd := exec.Command("bash", append([]string{script}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func scriptDir() string {
	exe, _ := os.Executable()
	exe, _ = filepath.EvalSymlinks(exe)
	return filepath.Dir(exe)
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
