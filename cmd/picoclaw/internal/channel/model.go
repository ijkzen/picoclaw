package channel

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/pkg/config"
)

// ViewState represents the current view in the TUI
type ViewState int

const (
	ViewList ViewState = iota
	ViewForm
	ViewTesting
	ViewTestResult
)

// FieldType represents the type of form field
type FieldType int

const (
	FieldTypeText FieldType = iota
	FieldTypeNumber
	FieldTypeBool
	FieldTypeArray
	FieldTypePassword
)

// FormField represents a single form field
type FormField struct {
	Name        string
	Label       string
	Type        FieldType
	Value       string
	Placeholder string
	Required    bool
	BoolValue   bool
}

// ChannelInfo holds information about a channel
type ChannelInfo struct {
	Name        string
	DisplayName string
	Configured  bool
	Enabled     bool
	Description string
}

// TestResult represents the result of a configuration test
type TestResult struct {
	Success bool
	Message string
	Error   error
}

// Model represents the application state
type Model struct {
	// Current view state
	// Current view state
	CurrentView ViewState

	// Config
	Config     *config.Config
	ConfigPath string

	// Channel list
	Channels []ChannelInfo
	Cursor   int

	// Form state
	FormFields      []FormField
	FormCursor      int
	FormInputs      []textinput.Model
	SelectedChannel string

	// Testing state
	TestResult TestResult
	Testing    bool
	Spinner    spinner.Model

	// Window dimensions
	Width  int
	Height int

	// Error message
	ErrorMsg string
}

// NewModel creates a new model with loaded config
func NewModel() (*Model, error) {
	// Load config
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".picoclaw", "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// Use default config if load fails
		cfg = config.DefaultConfig()
	}

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Create channel list
	channels := []ChannelInfo{
		{
			Name:        "telegram",
			DisplayName: "Telegram",
			Description: "Telegram Bot API",
			Configured:  cfg.Channels.Telegram.Token != "",
			Enabled:     cfg.Channels.Telegram.Enabled,
		},
		{
			Name:        "discord",
			DisplayName: "Discord",
			Description: "Discord Bot",
			Configured:  cfg.Channels.Discord.Token != "",
			Enabled:     cfg.Channels.Discord.Enabled,
		},
		{
			Name:        "slack",
			DisplayName: "Slack",
			Description: "Slack App/Bot",
			Configured:  cfg.Channels.Slack.BotToken != "",
			Enabled:     cfg.Channels.Slack.Enabled,
		},
		{
			Name:        "qq",
			DisplayName: "QQ",
			Description: "QQ Bot",
			Configured:  cfg.Channels.QQ.AppID != "",
			Enabled:     cfg.Channels.QQ.Enabled,
		},
		{
			Name:        "dingtalk",
			DisplayName: "DingTalk",
			Description: "DingTalk Bot",
			Configured:  cfg.Channels.DingTalk.ClientID != "",
			Enabled:     cfg.Channels.DingTalk.Enabled,
		},
		{
			Name:        "wecom",
			DisplayName: "WeCom (Bot)",
			Description: "WeCom Group Bot",
			Configured:  cfg.Channels.WeCom.Token != "",
			Enabled:     cfg.Channels.WeCom.Enabled,
		},
		{
			Name:        "wecom_app",
			DisplayName: "WeCom (App)",
			Description: "WeCom Self-built App",
			Configured:  cfg.Channels.WeComApp.CorpID != "",
			Enabled:     cfg.Channels.WeComApp.Enabled,
		},
		{
			Name:        "feishu",
			DisplayName: "Feishu",
			Description: "Feishu/Lark App",
			Configured:  cfg.Channels.Feishu.AppID != "",
			Enabled:     cfg.Channels.Feishu.Enabled,
		},
		{
			Name:        "line",
			DisplayName: "LINE",
			Description: "LINE Messaging API",
			Configured:  cfg.Channels.LINE.ChannelAccessToken != "",
			Enabled:     cfg.Channels.LINE.Enabled,
		},
		{
			Name:        "onebot",
			DisplayName: "OneBot",
			Description: "OneBot Protocol",
			Configured:  cfg.Channels.OneBot.WSUrl != "",
			Enabled:     cfg.Channels.OneBot.Enabled,
		},
		{
			Name:        "maixcam",
			DisplayName: "MaixCam",
			Description: "MaixCAM Hardware",
			Configured:  cfg.Channels.MaixCam.Host != "",
			Enabled:     cfg.Channels.MaixCam.Enabled,
		},
		{
			Name:        "whatsapp",
			DisplayName: "WhatsApp",
			Description: "WhatsApp Bridge",
			Configured:  cfg.Channels.WhatsApp.BridgeURL != "",
			Enabled:     cfg.Channels.WhatsApp.Enabled,
		},
	}

	return &Model{
		CurrentView: ViewList,
		Config:      cfg,
		ConfigPath:  configPath,
		Channels:    channels,
		Cursor:      0,
		FormCursor:  0,
		Spinner:     s,
		Testing:     false,
		Width:       80,
		Height:      24,
	}, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

// GetCurrentChannel returns the currently selected channel
func (m *Model) GetCurrentChannel() *ChannelInfo {
	if m.Cursor < 0 || m.Cursor >= len(m.Channels) {
		return nil
	}
	return &m.Channels[m.Cursor]
}

// RefreshChannelStatus updates the status of all channels from config
func (m *Model) RefreshChannelStatus() {
	for i := range m.Channels {
		switch m.Channels[i].Name {
		case "telegram":
			m.Channels[i].Configured = m.Config.Channels.Telegram.Token != ""
			m.Channels[i].Enabled = m.Config.Channels.Telegram.Enabled
		case "discord":
			m.Channels[i].Configured = m.Config.Channels.Discord.Token != ""
			m.Channels[i].Enabled = m.Config.Channels.Discord.Enabled
		case "slack":
			m.Channels[i].Configured = m.Config.Channels.Slack.BotToken != ""
			m.Channels[i].Enabled = m.Config.Channels.Slack.Enabled
		case "qq":
			m.Channels[i].Configured = m.Config.Channels.QQ.AppID != ""
			m.Channels[i].Enabled = m.Config.Channels.QQ.Enabled
		case "dingtalk":
			m.Channels[i].Configured = m.Config.Channels.DingTalk.ClientID != ""
			m.Channels[i].Enabled = m.Config.Channels.DingTalk.Enabled
		case "wecom":
			m.Channels[i].Configured = m.Config.Channels.WeCom.Token != ""
			m.Channels[i].Enabled = m.Config.Channels.WeCom.Enabled
		case "wecom_app":
			m.Channels[i].Configured = m.Config.Channels.WeComApp.CorpID != ""
			m.Channels[i].Enabled = m.Config.Channels.WeComApp.Enabled
		case "feishu":
			m.Channels[i].Configured = m.Config.Channels.Feishu.AppID != ""
			m.Channels[i].Enabled = m.Config.Channels.Feishu.Enabled
		case "line":
			m.Channels[i].Configured = m.Config.Channels.LINE.ChannelAccessToken != ""
			m.Channels[i].Enabled = m.Config.Channels.LINE.Enabled
		case "onebot":
			m.Channels[i].Configured = m.Config.Channels.OneBot.WSUrl != ""
			m.Channels[i].Enabled = m.Config.Channels.OneBot.Enabled
		case "maixcam":
			m.Channels[i].Configured = m.Config.Channels.MaixCam.Host != ""
			m.Channels[i].Enabled = m.Config.Channels.MaixCam.Enabled
		case "whatsapp":
			m.Channels[i].Configured = m.Config.Channels.WhatsApp.BridgeURL != ""
			m.Channels[i].Enabled = m.Config.Channels.WhatsApp.Enabled
		}
	}
}

// SaveConfig saves the current configuration to file
func (m *Model) SaveConfig() error {
	return internal.SaveConfigPathAndRestart(m.ConfigPath, m.Config)
}

// GetFieldValue retrieves a field value by name
func (m *Model) GetFieldValue(name string) string {
	for i := range m.FormFields {
		if m.FormFields[i].Name == name {
			if m.FormFields[i].Type == FieldTypeBool {
				return strconv.FormatBool(m.FormFields[i].BoolValue)
			}
			return m.FormFields[i].Value
		}
	}
	return ""
}

// GetFieldValueBool retrieves a bool field value by name
func (m *Model) GetFieldValueBool(name string) bool {
	for i := range m.FormFields {
		if m.FormFields[i].Name == name && m.FormFields[i].Type == FieldTypeBool {
			return m.FormFields[i].BoolValue
		}
	}
	return false
}

// GetFieldValueArray retrieves an array field value by name
func (m *Model) GetFieldValueArray(name string) []string {
	for i := range m.FormFields {
		if m.FormFields[i].Name == name && m.FormFields[i].Type == FieldTypeArray {
			if m.FormFields[i].Value == "" {
				return nil
			}
			parts := strings.Split(m.FormFields[i].Value, ",")
			for j := range parts {
				parts[j] = strings.TrimSpace(parts[j])
			}
			return parts
		}
	}
	return nil
}

// GetFieldValueInt retrieves an int field value by name
func (m *Model) GetFieldValueInt(name string) int {
	for i := range m.FormFields {
		if m.FormFields[i].Name == name && m.FormFields[i].Type == FieldTypeNumber {
			val, _ := strconv.Atoi(m.FormFields[i].Value)
			return val
		}
	}
	return 0
}

// GetFieldValueInt64 retrieves an int64 field value by name
func (m *Model) GetFieldValueInt64(name string) int64 {
	for i := range m.FormFields {
		if m.FormFields[i].Name == name && m.FormFields[i].Type == FieldTypeNumber {
			val, _ := strconv.ParseInt(m.FormFields[i].Value, 10, 64)
			return val
		}
	}
	return 0
}

// SetError sets an error message
func (m *Model) SetError(err error) {
	if err != nil {
		m.ErrorMsg = err.Error()
	} else {
		m.ErrorMsg = ""
	}
}

// StringArrayToString converts a string array to comma-separated string
func StringArrayToString(arr []string) string {
	return strings.Join(arr, ", ")
}
