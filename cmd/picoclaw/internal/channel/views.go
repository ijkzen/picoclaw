package channel

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View renders the current view
func (m Model) View() string {
	switch m.CurrentView {
	case ViewList:
		return m.ViewList()
	case ViewForm:
		return m.ViewForm()
	case ViewTesting:
		return m.ViewTesting()
	case ViewTestResult:
		return m.ViewTestResult()
	}
	return ""
}

// ViewList renders the channel list view
func (m Model) ViewList() string {
	var b strings.Builder

	// Header
	b.WriteString(Styles.Title.Render("📡 Channel Configuration"))
	b.WriteString("\n")
	b.WriteString(Styles.Subtitle.Render("Configure chat platform integrations"))
	b.WriteString("\n\n")

	// Channel list
	for i, channel := range m.Channels {
		isSelected := i == m.Cursor

		// Status icon
		statusIcon := GetStatusIcon(channel.Configured, channel.Enabled)
		statusText := GetStatusText(channel.Configured, channel.Enabled)

		// Channel line
		var line string
		channelName := Styles.ChannelName.Render(channel.DisplayName)
		desc := Styles.Description.Render(channel.Description)

		if isSelected {
			line = fmt.Sprintf("%s %s %s - %s (%s)",
				Styles.Selected.Render("▶"),
				statusIcon,
				channelName,
				desc,
				statusText,
			)
		} else {
			line = fmt.Sprintf("  %s %s - %s (%s)",
				statusIcon,
				channelName,
				desc,
				statusText,
			)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	// Help footer
	b.WriteString("\n")
	b.WriteString(m.ViewListHelp())

	return b.String()
}

// ViewListHelp renders the help text for list view
func (m Model) ViewListHelp() string {
	help := []string{
		"↑/↓ or j/k: Navigate",
		"Enter: Configure",
		"Space: Toggle enabled",
		"q/Esc: Quit",
	}
	return Styles.Help.Render(strings.Join(help, " • "))
}

// ViewForm renders the form view
func (m Model) ViewForm() string {
	var b strings.Builder

	// Header
	channel := m.GetCurrentChannel()
	if channel != nil {
		b.WriteString(Styles.Title.Render(fmt.Sprintf("⚙️  Configure %s", channel.DisplayName)))
		b.WriteString("\n")
		b.WriteString(Styles.Subtitle.Render(channel.Description))
		b.WriteString("\n\n")
	}

	// Form fields
	for i, field := range m.FormFields {
		isFocused := i == m.FormCursor

		// Field label
		label := Styles.FormLabel.Render(field.Label + ":")

		// Field value
		var value string
		if field.Type == FieldTypeBool {
			if field.BoolValue {
				value = Styles.Enabled.Render("✓ Enabled")
			} else {
				value = Styles.Disabled.Render("✗ Disabled")
			}
		} else if i < len(m.FormInputs) {
			if isFocused {
				value = Styles.FormFocused.Render(m.FormInputs[i].View())
			} else {
				value = Styles.FormInput.Render(m.FormInputs[i].View())
			}
		}

		if isFocused {
			b.WriteString(fmt.Sprintf("▶ %s %s", label, value))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", label, value))
		}
		b.WriteString("\n")
	}

	// Help footer
	b.WriteString("\n")
	b.WriteString(m.ViewFormHelp())

	return b.String()
}

// ViewFormHelp renders the help text for form view
func (m Model) ViewFormHelp() string {
	help := []string{
		"↑/↓: Navigate fields",
		"Tab/Shift+Tab: Next/Prev",
		"Space: Toggle boolean",
		"Ctrl+T: Test",
		"Ctrl+S: Save",
		"Esc: Back",
	}
	return Styles.Help.Render(strings.Join(help, " • "))
}

// ViewTesting renders the testing view
func (m Model) ViewTesting() string {
	var b strings.Builder

	b.WriteString(Styles.Title.Render("🧪 Testing Configuration"))
	b.WriteString("\n\n")

	channel := m.GetCurrentChannel()
	if channel != nil {
		b.WriteString(fmt.Sprintf("Testing %s configuration...", channel.DisplayName))
		b.WriteString("\n\n")
	}

	b.WriteString(m.Spinner.View())
	b.WriteString(" Running tests...")
	b.WriteString("\n\n")

	b.WriteString(Styles.Help.Render("Press Esc to cancel"))

	return b.String()
}

// ViewTestResult renders the test result view
func (m Model) ViewTestResult() string {
	var b strings.Builder

	if m.TestResult.Success {
		b.WriteString(Styles.Success.Render("✓ Configuration Test Successful!"))
	} else {
		b.WriteString(Styles.Error.Render("✗ Configuration Test Failed"))
	}
	b.WriteString("\n\n")

	if m.TestResult.Message != "" {
		b.WriteString(m.TestResult.Message)
		b.WriteString("\n\n")
	}

	if m.TestResult.Error != nil {
		errorBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1).
			Render(m.TestResult.Error.Error())
		b.WriteString(errorBox)
		b.WriteString("\n\n")
	}

	if m.TestResult.Success {
		b.WriteString(Styles.Help.Render("Enter/Y: Save config • N: Don't save • Esc: Back"))
	} else {
		b.WriteString(Styles.Help.Render("Enter: Retry • E: Edit • Esc: Back"))
	}

	return b.String()
}

// runChannelTUI starts the TUI application
func runChannelTUI() error {
	model, err := NewModel()
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
