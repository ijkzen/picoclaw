package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/pkg/config"
)

var supportedVendors = []string{
	"openai",
	"anthropic",
	"openrouter",
	"groq",
	"zhipu",
	"gemini",
	"vllm",
	"nvidia",
	"deepseek",
	"mistral",
	"ollama",
	"moonshot",
	"github_copilot",
	"qwen",
	"cerebras",
	"volcengine",
	"shengsuanyun",
	"antigravity",
}

type addModelTUI struct {
	cfg     *config.Config
	vendors []string
	cursor  int
	step    string // vendor, model, alias, apiKey, apiBase, confirm
	ti      textinput.Model
	inputs  map[string]string
	err     error
}

func newTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 60
	return ti
}

func (m addModelTUI) Init() tea.Cmd {
	if m.step == "vendor" {
		return nil
	}
	return nil
}

func (m addModelTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		// allow escape to quit at any step
		if k == "esc" {
			return m, tea.Quit
		}

		switch m.step {
		case "vendor":
			switch k {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.vendors)-1 {
					m.cursor++
				}
			case "enter":
				sel := m.vendors[m.cursor]
				m.inputs["vendor"] = sel
				m.step = "model"
				m.ti = newTextInput("model (e.g. openai/gpt-5.2)")
				m.ti.Focus()
			case "q", "esc":
				return m, tea.Quit
			}
		case "model", "alias", "apiKey", "apiBase":
			// let textinput handle keys
			var cmd tea.Cmd
			m.ti, cmd = m.ti.Update(msg)
			if msg.String() == "enter" {
				val := strings.TrimSpace(m.ti.Value())
				switch m.step {
				case "model":
					if val == "" {
						m.err = fmt.Errorf("model required")
						return m, tea.Quit
					}
					m.inputs["model"] = val
					m.step = "alias"
					m.ti = newTextInput("model alias (optional)")
					m.ti.Focus()
				case "alias":
					if val != "" {
						m.inputs["alias"] = val
					}
					m.step = "apiKey"
					m.ti = newTextInput("API Key (leave empty for local providers)")
					m.ti.Focus()
				case "apiKey":
					m.inputs["api_key"] = val
					// Show API Base prompt for OpenAI (optional), otherwise confirm
					vendor := m.inputs["vendor"]
					if vendor == "openai" {
						m.step = "apiBase"
						m.ti = newTextInput("API Base (optional, leave empty for default)")
						m.ti.Focus()
					} else {
						m.step = "confirm"
					}
				case "apiBase":
					m.inputs["api_base"] = val
					m.step = "confirm"
				}
			}
			return m, cmd
		case "confirm":
			switch k {
			case "enter", "y", "Y":
				// build ModelConfig and test
				// Ensure model includes vendor prefix when appropriate
				modelVal := m.inputs["model"]
				vendor := m.inputs["vendor"]
				if !strings.Contains(modelVal, "/") && vendor != "" {
					modelVal = vendor + "/" + modelVal
				}

				// Determine model_name: prefer user alias; otherwise use the
				// base model identifier (without vendor prefix).
				origModel := m.inputs["model"]
				modelBase := origModel
				if strings.Contains(origModel, "/") {
					parts := strings.SplitN(origModel, "/", 2)
					modelBase = parts[1]
				}

				modelName := m.inputs["alias"]
				if modelName == "" {
					modelName = modelBase
				}

				mc := config.ModelConfig{
					ModelName: modelName,
					Model:     modelVal,
					APIKey:    m.inputs["api_key"],
					APIBase:   m.inputs["api_base"],
				}
				if err := mc.Validate(); err != nil {
					m.err = err
					return m, tea.Quit
				}

				// save
				cfg, err := internal.LoadConfig()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				cfg.ModelList = append(cfg.ModelList, mc)
				if err := config.SaveConfig(internal.GetConfigPath(), cfg); err != nil {
					m.err = err
					return m, tea.Quit
				}
				return m, tea.Quit
			case "n", "N", "esc", "q":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m addModelTUI) View() string {
	s := ""
	switch m.step {
	case "vendor":
		s += "Select vendor (↑/↓ to navigate, Enter to select, Esc/q to cancel)\n\n"
		for i, v := range m.vendors {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, v)
		}
	case "model", "alias", "apiKey", "apiBase":
		s += fmt.Sprintf("%s\n\n%s", m.ti.Placeholder, m.ti.View())
		s += "\n\nPress Esc to cancel"
	case "confirm":
		s += "Confirm add model:\n\n"
		vendor := m.inputs["vendor"]
		alias := m.inputs["alias"]
		if alias == "" {
			alias = m.inputs["model"]
		}
		s += fmt.Sprintf("Vendor: %s\nModel: %s\nAlias: %s\nAPI Key: %s\nAPI Base: %s\n\n", vendor, m.inputs["model"], alias, m.inputs["api_key"], m.inputs["api_base"])
		s += "Press Enter/Y to save, N or Esc to cancel"
	}
	if m.err != nil {
		s += "\n\nError: " + m.err.Error()
	}
	return s
}

func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a model to model_list (interactive TUI)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := internal.LoadConfig()
			if err != nil {
				fmt.Println("failed to load config:", err)
				return
			}

			m := addModelTUI{
				cfg:     cfg,
				vendors: supportedVendors,
				cursor:  0,
				step:    "vendor",
				ti:      newTextInput(""),
				inputs:  make(map[string]string),
			}
			p := tea.NewProgram(m, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Println("tui error:", err)
				return
			}
		},
	}
	return cmd
}

