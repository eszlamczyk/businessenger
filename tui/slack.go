package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

// ── Slack config types ─────────────────────────────────────────────────────────

type SlackWorkspace struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type SlackConfig struct {
	Workspaces []SlackWorkspace `json:"workspaces"`
}

// ── Async message ──────────────────────────────────────────────────────────────

type slackChannelsMsg struct {
	channels []string
	err      error
}

// ── Date options ───────────────────────────────────────────────────────────────

var slackDateOptions = []string{
	"today",
	"yesterday",
	"1 week ago",
	"2 weeks ago",
	"1 month ago",
	"custom...",
}

// ── API fetch ──────────────────────────────────────────────────────────────────

func fetchSlackChannels(token string) tea.Cmd {
	return func() tea.Msg {
		url := "https://slack.com/api/conversations.list?types=public_channel,private_channel&limit=200&exclude_archived=true"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return slackChannelsMsg{err: err}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return slackChannelsMsg{err: err}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return slackChannelsMsg{err: err}
		}

		var result struct {
			OK       bool   `json:"ok"`
			Error    string `json:"error"`
			Channels []struct {
				Name string `json:"name"`
			} `json:"channels"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return slackChannelsMsg{err: err}
		}
		if !result.OK {
			return slackChannelsMsg{err: fmt.Errorf("slack API error: %s", result.Error)}
		}

		channels := make([]string, 0, len(result.Channels))
		for _, ch := range result.Channels {
			channels = append(channels, ch.Name)
		}
		return slackChannelsMsg{channels: channels}
	}
}
