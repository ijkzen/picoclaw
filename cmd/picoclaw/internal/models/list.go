package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/spf13/cobra"
)

type listModelTUI struct {
	cfg    *config.Config
	cursor int
	err    error
}

func (m listModelTUI) Init() tea.Cmd { return nil }

func (m listModelTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.cfg.ModelList)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.cfg.ModelList) == 0 {
				return m, tea.Quit
			}
			sel := m.cfg.ModelList[m.cursor]
			name := sel.ModelName
			if strings.TrimSpace(name) == "" {
				name = sel.Model
			}
			if m.cfg.Agents.Defaults.GetModelName() == name {
				return m, tea.Quit
			}
			m.cfg.Agents.Defaults.ModelName = name
			if err := config.SaveConfig(internal.GetConfigPath(), m.cfg); err != nil {
				m.err = err
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m listModelTUI) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error saving config: %v\n", m.err)
	}
	if len(m.cfg.ModelList) == 0 {
		return "no models configured\n"
	}

	var b strings.Builder
	b.WriteString("Select model to set as default (Enter to select, q to cancel)\n\n")
	for i, md := range m.cfg.ModelList {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		status := "unconfigured"
		lowerModel := strings.ToLower(md.Model)
		if strings.HasPrefix(lowerModel, "ollama/") {
			if md.APIKey != "" || md.APIBase != "" || md.Workspace != "" {
				status = "configured"
			}
		} else {
			if md.APIKey != "" {
				status = "configured"
			}
		}
		name := md.ModelName
		if strings.TrimSpace(name) == "" {
			name = md.Model
		}
		b.WriteString(fmt.Sprintf("%s %d) %s (%s) - %s\n", cursor, i+1, name, md.Model, status))
	}
	return b.String()
}

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List models in model_list (shows configured state)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := internal.LoadConfig()
			if err != nil {
				fmt.Println("failed to load config:", err)
				return
			}

			m := listModelTUI{
				cfg:    cfg,
				cursor: 0,
			}

			p := tea.NewProgram(m, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Println("tui error:", err)
				return
			}
			if m.err != nil {
				fmt.Println("failed to save config:", m.err)
				return
			}
		},
	}
	return cmd
}
