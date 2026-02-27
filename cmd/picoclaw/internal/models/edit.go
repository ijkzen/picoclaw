package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/spf13/cobra"
)

type editModelTUI struct {
	cfg        *config.Config
	cursor     int
	index      int
	stage      int // 0: select, 1: edit, 2: confirm
	models     []config.ModelConfig
	inputs     map[string]*textinput.Model
	supported  bool
	orig       config.ModelConfig
	mc         config.ModelConfig
	focusIndex int
	errMsg     string
}

func (m editModelTUI) Init() tea.Cmd { return nil }

func (m editModelTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.stage {
		case 0: // selection
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.models)-1 {
					m.cursor++
				}
			case "enter":
				m.index = m.cursor
				m.orig = m.models[m.index]
				m.mc = m.orig
				// detect vendor
				vendor := ""
				if strings.Contains(m.mc.Model, "/") {
					parts := strings.SplitN(m.mc.Model, "/", 2)
					vendor = parts[0]
				}
				m.supported = false
				for _, v := range supportedVendors {
					if v == vendor {
						m.supported = true
						break
					}
				}
				// prepare inputs
				m.inputs = make(map[string]*textinput.Model)
				newTI := func(placeholder, value string, limit int) *textinput.Model {
					ti := textinput.New()
					ti.Placeholder = placeholder
					ti.SetValue(value)
					ti.CharLimit = limit
					return &ti
				}

				// openai and unknown vendors -> allow editing all fields; supported others -> alias+api_key
				if vendor == "openai" || !m.supported {
					m.inputs["model"] = newTI("model (vendor/model)", m.mc.Model, 512)
					m.inputs["alias"] = newTI("alias", m.mc.ModelName, 256)
					m.inputs["api_base"] = newTI("api_base", m.mc.APIBase, 1024)
					m.inputs["api_key"] = newTI("api_key", m.mc.APIKey, 1024)
				} else {
					m.inputs["alias"] = newTI("alias", m.mc.ModelName, 256)
					m.inputs["api_key"] = newTI("api_key", m.mc.APIKey, 1024)
				}

				// focus first input deterministically
				order := []string{"model", "alias", "api_base", "api_key"}
				for i, k := range order {
					if ti, ok := m.inputs[k]; ok {
						ti.Focus()
						m.focusIndex = i
						break
					}
				}
				m.stage = 1
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		case 1: // edit
			// handle navigation between inputs and forwarding to focused input
			keys := msg.String()
			switch keys {
			case "tab":
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
				i := 0
				for k := range m.inputs {
					if i == m.focusIndex {
						m.inputs[k].Focus()
					} else {
						m.inputs[k].Blur()
					}
					i++
				}
			case "shift+tab":
				m.focusIndex = (m.focusIndex - 1 + len(m.inputs)) % len(m.inputs)
				i := 0
				for k := range m.inputs {
					if i == m.focusIndex {
						m.inputs[k].Focus()
					} else {
						m.inputs[k].Blur()
					}
					i++
				}
			case "enter":
				// if focus is last field, move to confirm
				if m.focusIndex == len(m.inputs)-1 {
					// collect values
					if _, ok := m.inputs["alias"]; ok {
						m.mc.ModelName = m.inputs["alias"].Value()

						if strings.Contains(m.orig.Model, "/") {
							parts := strings.SplitN(m.orig.Model, "/", 2)
							vendor := parts[0]
							// preserve vendor prefix and use new ModelName as suffix
							m.mc.Model = vendor + "/" + m.mc.ModelName
						}
					}

					if _, ok := m.inputs["model"]; ok {
						m.mc.Model = m.inputs["model"].Value()
					}
					if _, ok := m.inputs["api_base"]; ok {
						m.mc.APIBase = m.inputs["api_base"].Value()
					}
					if _, ok := m.inputs["api_key"]; ok {
						m.mc.APIKey = m.inputs["api_key"].Value()
					}
					m.stage = 2
				} else {
					// move focus to next
					m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
					i := 0
					for k := range m.inputs {
						if i == m.focusIndex {
							m.inputs[k].Focus()
						} else {
							m.inputs[k].Blur()
						}
						i++
					}
				}
			case "esc":
				return m, tea.Quit
			}
			// forward message to inputs
			for k := range m.inputs {
				tiPtr := m.inputs[k]
				newTi, cmd := tiPtr.Update(msg)
				*m.inputs[k] = newTi
				if cmd != nil {
					return m, cmd
				}
			}
		case 2: // confirm
			switch msg.String() {
			case "y", "Y", "enter":
				// save: validate and test
				// if the user changed the friendly name, sync it to the model field
				if err := m.mc.Validate(); err != nil {
					m.errMsg = fmt.Sprintf("invalid config: %v", err)
					return m, nil
				}

				m.cfg.ModelList[m.index] = m.mc
				if err := internal.SaveConfigAndRestart(m.cfg); err != nil {
					m.errMsg = fmt.Sprintf("failed to save: %v", err)
					return m, nil
				}
				m.errMsg = "model updated"
				return m, tea.Quit
			case "n", "N":
				m.errMsg = "aborted"
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m editModelTUI) View() string {
	var b strings.Builder
	switch m.stage {
	case 0:
		b.WriteString("Select model to edit:\n\n")
		for i, md := range m.models {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString(fmt.Sprintf(" %s %d) %s (%s)\n", cursor, i+1, md.ModelName, md.Model))
		}
		b.WriteString("\nUse up/down to navigate, enter to select, q to quit\n")
	case 1:
		b.WriteString("Edit fields (tab to switch, enter to confirm):\n\n")
		// render inputs in deterministic order
		order := []string{}
		if _, ok := m.inputs["model"]; ok {
			order = append(order, "model")
		}
		if _, ok := m.inputs["alias"]; ok {
			order = append(order, "alias")
		}
		if _, ok := m.inputs["api_base"]; ok {
			order = append(order, "api_base")
		}
		if _, ok := m.inputs["api_key"]; ok {
			order = append(order, "api_key")
		}
		// render inputs with labels
		labelMap := map[string]string{
			"model":    "Model",
			"alias":    "Alias",
			"api_base": "API Base",
			"api_key":  "API Key",
		}
		for _, k := range order {
			label := labelMap[k]
			b.WriteString(fmt.Sprintf("%s: %s\n", label, m.inputs[k].View()))
		}
		b.WriteString("\nPress Enter on last field to continue to confirmation. Esc to cancel.\n")
	case 2:
		b.WriteString("Confirm changes:\n\n")
		b.WriteString(fmt.Sprintf("Model: %s\n", m.mc.Model))
		b.WriteString(fmt.Sprintf("Alias: %s\n", m.mc.ModelName))
		b.WriteString(fmt.Sprintf("API Base: %s\n", m.mc.APIBase))
		b.WriteString(fmt.Sprintf("API Key: %s\n", m.mc.APIKey))
		b.WriteString("\nSave changes? (y/n)\n")
	}
	if m.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(m.errMsg)
		b.WriteString("\n")
	}
	return b.String()
}

func NewEditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit an existing model in model_list (TUI)",
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
			m := editModelTUI{
				cfg:    cfg,
				models: cfg.ModelList,
				stage:  0,
			}
			p := tea.NewProgram(m, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Println("tui error:", err)
			}
		},
	}
	return cmd
}
