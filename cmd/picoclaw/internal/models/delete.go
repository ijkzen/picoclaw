package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/pkg/config"
)

type deleteModelTUI struct {
	cfg      *config.Config
	items    []string
	cursor   int
	selected map[int]bool
	err      error
	state    string // "select" or "confirm"
}

func (m deleteModelTUI) Init() tea.Cmd { return nil }

func (m deleteModelTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ":
			if m.state == "select" && len(m.items) > 0 {
				if m.selected == nil {
					m.selected = make(map[int]bool)
				}
				m.selected[m.cursor] = !m.selected[m.cursor]
			}
		case "enter":
			if m.state == "select" {
				// if none selected, do nothing
				any := false
				for _, v := range m.selected {
					if v {
						any = true
						break
					}
				}
				if !any {
					return m, nil
				}
				m.state = "confirm"
				return m, nil
			} else if m.state == "confirm" {
				// perform deletion
				newList := make([]config.ModelConfig, 0, len(m.cfg.ModelList))
				for i, mm := range m.cfg.ModelList {
					if m.selected[i] {
						continue
					}
					newList = append(newList, mm)
				}
				m.cfg.ModelList = newList
				if err := config.SaveConfig(internal.GetConfigPath(), m.cfg); err != nil {
					m.err = err
				}
				return m, tea.Quit
			}
		case "y", "Y":
			if m.state == "confirm" {
				newList := make([]config.ModelConfig, 0, len(m.cfg.ModelList))
				for i, mm := range m.cfg.ModelList {
					if m.selected[i] {
						continue
					}
					newList = append(newList, mm)
				}
				m.cfg.ModelList = newList
				if err := config.SaveConfig(internal.GetConfigPath(), m.cfg); err != nil {
					m.err = err
				}
				return m, tea.Quit
			}
		case "n", "N":
			if m.state == "confirm" {
				m.state = "select"
			}
		}
	}
	return m, nil
}

func (m deleteModelTUI) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if len(m.items) == 0 {
		return "no models configured\n"
	}

	s := "Select models to delete (Space to toggle, Enter to confirm, q to cancel)\n\n"
	for i, it := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		sel := "[ ]"
		if m.selected[i] {
			sel = "[x]"
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, sel, it)
	}

	if m.state == "confirm" {
		s += "\nConfirm deletion? (y/N)"
	}

	return s
}

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete one or more models from model_list (interactive TUI)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := internal.LoadConfig()
			if err != nil {
				fmt.Println("failed to load config:", err)
				return
			}
			if len(cfg.ModelList) == 0 {
				fmt.Println("no models configured")
				return
			}

			items := make([]string, len(cfg.ModelList))
			for i, m := range cfg.ModelList {
				items[i] = fmt.Sprintf("%s (%s)", m.ModelName, m.Model)
			}

			m := deleteModelTUI{
				cfg:      cfg,
				items:    items,
				cursor:   0,
				selected: make(map[int]bool),
				state:    "select",
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
			fmt.Println("operation completed")
		},
	}
	return cmd
}
