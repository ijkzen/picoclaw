package channel

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all message updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case testCompleteMsg:
		m.Testing = false
		m.TestResult = TestResult(msg)
		m.CurrentView = ViewTestResult
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	// Handle form input updates
	if m.CurrentView == ViewForm {
		for i := range m.FormInputs {
			if i < len(m.FormInputs) {
				var cmd tea.Cmd
				m.FormInputs[i], cmd = m.FormInputs[i].Update(msg)
				if cmd != nil {
					return m, cmd
				}
			}
		}
	}

	return m, nil
}

// handleKeyMsg handles keyboard input based on current view
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.CurrentView {
	case ViewList:
		return m.handleListKey(msg)
	case ViewForm:
		return m.handleFormKey(msg)
	case ViewTesting:
		// During testing, only allow quit
		if msg.String() == "q" || msg.String() == "esc" {
			m.CurrentView = ViewList
			return m, nil
		}
		return m, nil
	case ViewTestResult:
		return m.handleTestResultKey(msg)
	}
	return m, nil
}

// handleListKey handles keys in list view
func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Channels)-1 {
			m.Cursor++
		}

	case "enter":
		channel := m.GetCurrentChannel()
		if channel != nil {
			m.SelectedChannel = channel.Name
			m.InitFormForChannel(channel.Name)
			m.CurrentView = ViewForm
		}

	case " ":
		// Toggle enabled status
		channel := m.GetCurrentChannel()
		if channel != nil {
			m.ToggleChannelEnabled(channel.Name)
		}
	}

	return m, nil
}

// handleFormKey handles keys in form view
func (m Model) handleFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentView = ViewList
		m.FormInputs = nil
		return m, nil

	case "tab":
		m.FormCursor++
		if m.FormCursor >= len(m.FormFields) {
			m.FormCursor = 0
		}
		m.UpdateFormFocus()

	case "shift+tab":
		m.FormCursor--
		if m.FormCursor < 0 {
			m.FormCursor = len(m.FormFields) - 1
		}
		m.UpdateFormFocus()

	case "up":
		if m.FormCursor > 0 {
			m.FormCursor--
			m.UpdateFormFocus()
		}

	case "down":
		if m.FormCursor < len(m.FormFields)-1 {
			m.FormCursor++
			m.UpdateFormFocus()
		}

	case " ":
		// Toggle bool fields
		if m.FormCursor < len(m.FormFields) {
			field := &m.FormFields[m.FormCursor]
			if field.Type == FieldTypeBool {
				field.BoolValue = !field.BoolValue
			}
		}

	case "ctrl+t":
		// Test configuration
		m.SaveFormToConfig()
		m.CurrentView = ViewTesting
		m.Testing = true
		return m, m.TestChannel()

	case "ctrl+s":
		// Save and go back
		m.SaveFormToConfig()
		if err := m.SaveConfig(); err != nil {
			m.SetError(err)
		} else {
			m.RefreshChannelStatus()
			m.CurrentView = ViewList
			m.FormInputs = nil
		}
		return m, nil
	}

	// Update form input values
	if m.FormCursor < len(m.FormInputs) && m.FormCursor < len(m.FormFields) {
		if m.FormFields[m.FormCursor].Type != FieldTypeBool {
			var cmd tea.Cmd
			m.FormInputs[m.FormCursor], cmd = m.FormInputs[m.FormCursor].Update(msg)
			m.FormFields[m.FormCursor].Value = m.FormInputs[m.FormCursor].Value()
			return m, cmd
		}
	}

	return m, nil
}

// handleTestResultKey handles keys in test result view
func (m Model) handleTestResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.CurrentView = ViewForm

	case "enter", "y":
		if m.TestResult.Success {
			// Save config and go back to list
			if err := m.SaveConfig(); err != nil {
				m.SetError(err)
			} else {
				m.RefreshChannelStatus()
				m.CurrentView = ViewList
				m.FormInputs = nil
			}
		} else {
			// Retry test
			m.CurrentView = ViewTesting
			m.Testing = true
			return m, m.TestChannel()
		}

	case "n":
		if m.TestResult.Success {
			// Don't save, go back to form
			m.CurrentView = ViewForm
		}

	case "e":
		if !m.TestResult.Success {
			// Edit configuration
			m.CurrentView = ViewForm
		}
	}

	return m, nil
}

// UpdateFormFocus updates the focus state of form inputs
func (m *Model) UpdateFormFocus() {
	for i := range m.FormInputs {
		if i == m.FormCursor {
			m.FormInputs[i].Focus()
		} else {
			m.FormInputs[i].Blur()
		}
	}
}

// ToggleChannelEnabled toggles the enabled state of a channel
func (m *Model) ToggleChannelEnabled(channelName string) {
	switch channelName {
	case "telegram":
		m.Config.Channels.Telegram.Enabled = !m.Config.Channels.Telegram.Enabled
	case "discord":
		m.Config.Channels.Discord.Enabled = !m.Config.Channels.Discord.Enabled
	case "slack":
		m.Config.Channels.Slack.Enabled = !m.Config.Channels.Slack.Enabled
	case "qq":
		m.Config.Channels.QQ.Enabled = !m.Config.Channels.QQ.Enabled
	case "dingtalk":
		m.Config.Channels.DingTalk.Enabled = !m.Config.Channels.DingTalk.Enabled
	case "wecom":
		m.Config.Channels.WeCom.Enabled = !m.Config.Channels.WeCom.Enabled
	case "wecom_app":
		m.Config.Channels.WeComApp.Enabled = !m.Config.Channels.WeComApp.Enabled
	case "feishu":
		m.Config.Channels.Feishu.Enabled = !m.Config.Channels.Feishu.Enabled
	case "line":
		m.Config.Channels.LINE.Enabled = !m.Config.Channels.LINE.Enabled
	case "onebot":
		m.Config.Channels.OneBot.Enabled = !m.Config.Channels.OneBot.Enabled
	case "maixcam":
		m.Config.Channels.MaixCam.Enabled = !m.Config.Channels.MaixCam.Enabled
	case "whatsapp":
		m.Config.Channels.WhatsApp.Enabled = !m.Config.Channels.WhatsApp.Enabled
	}

	// Save config
	_ = m.SaveConfig()
	m.RefreshChannelStatus()
}

// testCompleteMsg is sent when testing is complete
type testCompleteMsg TestResult

// TestChannel tests the current channel configuration
func (m Model) TestChannel() tea.Cmd {
	return func() tea.Msg {
		result := m.PerformTest()
		return testCompleteMsg(result)
	}
}
